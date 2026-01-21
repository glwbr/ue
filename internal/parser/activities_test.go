package parser

import (
	"encoding/json"
	"os"
	"testing"
	"uber-extractor/internal/uberapi"
)

func TestParseActivitiesResponse(t *testing.T) {
	data, err := os.ReadFile("../../testdata/activities.response.json")
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

func TestParseActivitiesFullResponse(t *testing.T) {
	data, err := os.ReadFile("testdata/activities_full.json")
	if err != nil {
		t.Fatalf("failed to read test data: %v", err)
	}

	var response uberapi.ActivitiesResponse

	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatalf("failed to unmarshal activities: %v", err)
	}

	if len(response.Data.Activities.Past.Activities) != 5 {
		t.Errorf("expected 5 activities, got %d", len(response.Data.Activities.Past.Activities))
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

	if completedTrips != 2 {
		t.Errorf("expected 2 completed trips, got %d", completedTrips)
	}

	if canceledTrips != 3 {
		t.Errorf("expected 3 canceled trips, got %d", canceledTrips)
	}

	expectedTotalFare := 22.67
	if totalFare != expectedTotalFare {
		t.Errorf("total fare = %v, want %v", totalFare, expectedTotalFare)
	}
}

func TestParseTimeFromRealData(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "subtitle time format",
			input: "Jan 19 • 4:33 PM",
			valid: false,
		},
		{
			name:  "JavaScript date string",
			input: "Tue Jan 06 2026 21:12:28 GMT+0000 (Coordinated Universal Time)",
			valid: true,
		},
		{
			name:  "RFC3339",
			input: "2024-01-15T10:30:00Z",
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Time(tt.input)
			if (err == nil) != tt.valid {
				t.Errorf("Time(%q) valid = %v, want %v, err = %v", tt.input, err == nil, tt.valid, err)
			}
		})
	}
}

func TestExtractCoordinatesFromRealMapURL(t *testing.T) {
	data, err := os.ReadFile("testdata/activities_full.json")
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

	lat, lon := ExtractCoordinates(mapURL, 0)

	if lat == 0 || lon == 0 {
		t.Errorf("ExtractCoordinates() returned zero values from real map URL")
	}

	expectedLat := -12.26071
	expectedLon := -38.94518

	if lat != expectedLat {
		t.Errorf("latitude = %v, want %v", lat, expectedLat)
	}

	if lon != expectedLon {
		t.Errorf("longitude = %v, want %v", lon, expectedLon)
	}
}

func TestFareFromRealData(t *testing.T) {
	data, err := os.ReadFile("testdata/activities_full.json")
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
			uuid:            "3bd1ebe2-cc40-4fb6-97d5-205e4fb71418",
			expectedFare:    10.13,
			shouldBeNonZero: true,
		},
		{
			uuid:            "a5cc29c7-2c9f-4373-9182-d507780751b9",
			expectedFare:    12.54,
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
		{
			uuid:            "75a7a2c7-c283-4e3e-8bf8-5b6d5137a3cb",
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
