package parser

import (
	"encoding/json"
	"os"
	"testing"
	"uber-extractor/internal/uberapi"
)

func TestParseActivitiesResponse(t *testing.T) {
	data, err := os.ReadFile("testdata/activities_response_simple.json")
	if err != nil {
		t.Fatalf("failed to read test data: %v", err)
	}

	var response uberapi.ActivitiesResponse

	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatalf("failed to unmarshal activities: %v", err)
	}

	if len(response.Data.Activities.Past.Activities) != 2 {
		t.Errorf("expected 2 activities, got %d", len(response.Data.Activities.Past.Activities))
	}

	if response.Data.Activities.Past.Activities[0].UUID != "00000000-0000-0000-0000-000000000001" {
		t.Errorf("unexpected first UUID: %s", response.Data.Activities.Past.Activities[0].UUID)
	}

	if response.Data.Activities.Past.NextPageToken == "" {
		t.Error("nextPageToken should not be empty")
	}
}

func TestParseObfuscatedActivitiesResponse(t *testing.T) {
	data, err := os.ReadFile("testdata/activities_obfuscated.json")
	if err != nil {
		t.Fatalf("failed to read test data: %v", err)
	}

	var response uberapi.ActivitiesResponse

	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatalf("failed to unmarshal activities: %v", err)
	}

	if len(response.Data.Activities.Past.Activities) != 9 {
		t.Errorf("expected 9 activities, got %d", len(response.Data.Activities.Past.Activities))
	}

	completedTrips := 0
	canceledTrips := 0
	var totalFare float64

	for _, activity := range response.Data.Activities.Past.Activities {
		fare := Fare(activity.Description)
		totalFare += fare

		if fare > 0 {
			completedTrips++
		} else {
			canceledTrips++
		}
	}

	if completedTrips != 6 {
		t.Errorf("expected 6 completed trips, got %d", completedTrips)
	}

	if canceledTrips != 3 {
		t.Errorf("expected 3 canceled trips, got %d", canceledTrips)
	}

	expectedTotalFare := 134.95
	if totalFare != expectedTotalFare {
		t.Errorf("total fare = %v, want %v", totalFare, expectedTotalFare)
	}
}

func TestExtractCoordinatesFromObfuscatedMapURL(t *testing.T) {
	data, err := os.ReadFile("testdata/activities_obfuscated.json")
	if err != nil {
		t.Fatalf("failed to read test data: %v", err)
	}

	var response struct {
		Data struct {
			Activities struct {
				Past struct {
					Activities []struct {
						UUID     string `json:"uuid"`
						Title    string `json:"title"`
						ImageURL struct {
							Light string `json:"light"`
						} `json:"imageURL"`
					} `json:"activities"`
				} `json:"past"`
			} `json:"activities"`
		} `json:"data"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatalf("failed to unmarshal activities: %v", err)
	}

	activityWithMap := response.Data.Activities.Past.Activities[0]
	mapURL := activityWithMap.ImageURL.Light

	lat, lon, err := ExtractCoordinates(mapURL, 0)
	if err != nil {
		t.Fatalf("ExtractCoordinates() failed: %v", err)
	}

	if lat == 0 || lon == 0 {
		t.Errorf("ExtractCoordinates() returned zero values from map URL")
	}

	expectedLat := 41.4089
	expectedLon := -75.6624

	if lat != expectedLat {
		t.Errorf("latitude = %v, want %v", lat, expectedLat)
	}

	if lon != expectedLon {
		t.Errorf("longitude = %v, want %v", lon, expectedLon)
	}
}

func TestFareFromObfuscatedData(t *testing.T) {
	data, err := os.ReadFile("testdata/activities_obfuscated.json")
	if err != nil {
		t.Fatalf("failed to read test data: %v", err)
	}

	var response uberapi.ActivitiesResponse

	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatalf("failed to unmarshal activities: %v", err)
	}

	testCases := []struct {
		uuid            string
		expectedFare    float64
		shouldBeNonZero bool
	}{
		{
			uuid:            "6b8dc458-d2ea-42a1-97e9-7db671798503",
			expectedFare:    15.50,
			shouldBeNonZero: true,
		},
		{
			uuid:            "3bd1ebe2-cc40-4fb6-97d5-205e4fb71418",
			expectedFare:    10.75,
			shouldBeNonZero: true,
		},
		{
			uuid:            "2608cdbb-ee14-4752-9a31-2b7e7cb084c3",
			expectedFare:    0.00,
			shouldBeNonZero: false,
		},
		{
			uuid:            "7fecd5c2-ff22-4123-8666-fe286c38133d",
			expectedFare:    0.00,
			shouldBeNonZero: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.uuid[:8], func(t *testing.T) {
			var description string
			for _, activity := range response.Data.Activities.Past.Activities {
				if activity.UUID == tc.uuid {
					description = activity.Description
					break
				}
			}

			if description == "" {
				t.Fatalf("activity %s not found", tc.uuid)
			}

			fare := Fare(description)

			if fare != tc.expectedFare {
				t.Errorf("Fare() = %v, want %v", fare, tc.expectedFare)
			}

			if tc.shouldBeNonZero && fare == 0 {
				t.Errorf("Fare() returned 0 but expected non-zero")
			}

			if !tc.shouldBeNonZero && fare != 0 {
				t.Errorf("Fare() returned non-zero %v but expected 0 (canceled trip)", fare)
			}
		})
	}
}
