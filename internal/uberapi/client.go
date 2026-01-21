package uberapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
	uberEndpoint   = "https://riders.uber.com/graphql"
)

type Client struct {
	httpClient *http.Client
	cookie     string
	endpoint   string
}

func NewClient(cookie string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		cookie:   cookie,
		endpoint: uberEndpoint,
	}
}

func (c *Client) GetActivities(ctx context.Context, startTime, endTime int64, pageToken string) (*ActivitiesResponse, string, error) {
	query := `
	query Activities($startTimeMs: Float, $endTimeMs: Float, $limit: Int = 50, $nextPageToken: String) {
		activities {
			past(
				startTimeMs: $startTimeMs
				endTimeMs: $endTimeMs
				nextPageToken: $nextPageToken
				limit: $limit
				orderTypes: [RIDES]
				profileType: PERSONAL
			) {
				activities {
					uuid
					title
					subtitle
					description
					__typename
				}
				nextPageToken
				__typename
			}
			__typename
		}
	}
	`

	request := map[string]interface{}{
		"query": query,
		"variables": map[string]interface{}{
			"startTimeMs":   float64(startTime),
			"endTimeMs":     float64(endTime),
			"nextPageToken": pageToken,
		},
	}

	body, err := c.makeRequest(ctx, request)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get activities: %w", err)
	}

	var response ActivitiesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, "", fmt.Errorf("failed to unmarshal activities response: %w", err)
	}

	return &response, response.Data.Activities.Past.NextPageToken, nil
}

func (c *Client) GetTrip(ctx context.Context, uuid string) (*GetTripResponse, error) {
	query := `
	query GetTrip($tripUUID: String!) {
		getTrip(tripUUID: $tripUUID) {
			trip {
				beginTripTime
				dropoffTime
				cityID
				countryID
				status
				fare
				driver
				uuid
				vehicleDisplayName
				waypoints
				marketplace
				__typename
			}
			mapURL
			receipt {
				distance
				distanceLabel
				duration
				vehicleType
				__typename
			}
			rating
			__typename
		}
	}
	`

	request := map[string]interface{}{
		"query": query,
		"variables": map[string]interface{}{
			"tripUUID": uuid,
		},
	}

	body, err := c.makeRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}

	var response GetTripResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal trip response: %w", err)
	}

	return &response, nil
}

func (c *Client) GetCurrentUser(ctx context.Context) (*UserResponse, error) {
	query := `
	query CurrentUser {
		currentUser {
			firstName
			lastName
			email
			__typename
		}
	}
	`

	request := map[string]interface{}{
		"query": query,
	}

	body, err := c.makeRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}

	var response UserResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user response: %w", err)
	}

	return &response, nil
}

func (c *Client) makeRequest(ctx context.Context, request map[string]interface{}) ([]byte, error) {
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", c.cookie)
	req.Header.Set("x-csrf-token", "x")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
