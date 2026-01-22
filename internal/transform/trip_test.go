package transform

import (
	"encoding/json"
	"os"
	"testing"

	"uber-extractor/internal/locations"
	"uber-extractor/internal/trips"
	"uber-extractor/internal/uberapi"
)

func TestProcessTrip(t *testing.T) {
	t.Run("process trip without location processor", func(t *testing.T) {
		data, err := os.ReadFile("testdata/trip_details_obfuscated.json")
		if err != nil {
			t.Fatalf("failed to read test data: %v", err)
		}

		var response uberapi.GetTripResponse
		if err := json.Unmarshal(data, &response); err != nil {
			t.Fatalf("failed to unmarshal trip details: %v", err)
		}

		trip, err := ProcessTrip(&response, nil)
		if err != nil {
			t.Fatalf("ProcessTrip() failed: %v", err)
		}

		if trip.UUID != "6b8dc458-d2ea-42a1-97e9-7db671798503" {
			t.Errorf("expected UUID %s, got %s", "6b8dc458-d2ea-42a1-97e9-7db671798503", trip.UUID)
		}

		if trip.Status != trips.StatusCompleted {
			t.Errorf("expected status %v, got %v", trips.StatusCompleted, trip.Status)
		}

		if trip.Driver != "Michael Scott" {
			t.Errorf("expected driver %s, got %s", "Michael Scott", trip.Driver)
		}

		if trip.VehicleType != "UberX" {
			t.Errorf("expected vehicle type %s, got %s", "UberX", trip.VehicleType)
		}

		if trip.Fare != 17.65 {
			t.Errorf("expected fare %v, got %v", 17.65, trip.Fare)
		}

		if trip.Distance != 8.63 {
			t.Errorf("expected distance %v, got %v", 8.63, trip.Distance)
		}

		if trip.Rating != 0 {
			t.Errorf("expected rating 0, got %d", trip.Rating)
		}

		if trip.Duration != 21 {
			t.Errorf("expected duration 21, got %v", trip.Duration)
		}

		if trip.PickupAddress != "1725 Slough Avenue, Dunder Mifflin, Scranton - PA, 18505" {
			t.Errorf("unexpected pickup address: %s", trip.PickupAddress)
		}

		if trip.DropoffAddress != "123 Kellum Court, Scranton - PA, 18508" {
			t.Errorf("unexpected dropoff address: %s", trip.DropoffAddress)
		}

		if trip.PickupLat != 41.4089 {
			t.Errorf("expected pickup lat %v, got %v", 41.4089, trip.PickupLat)
		}

		if trip.PickupLon != -75.6624 {
			t.Errorf("expected pickup lon %v, got %v", -75.6624, trip.PickupLon)
		}

		if trip.DropoffLat != 41.4120 {
			t.Errorf("expected dropoff lat %v, got %v", 41.4120, trip.DropoffLat)
		}

		if trip.DropoffLon != -75.6580 {
			t.Errorf("expected dropoff lon %v, got %v", -75.6580, trip.DropoffLon)
		}

		if trip.BeginTime.IsZero() {
			t.Error("expected non-zero begin time")
		}

		if trip.EndTime.IsZero() {
			t.Error("expected non-zero end time")
		}

		if trip.MapURL == "" {
			t.Error("expected non-empty map URL")
		}
	})

	t.Run("process trip with location processor", func(t *testing.T) {
		data, err := os.ReadFile("testdata/trip_details_obfuscated.json")
		if err != nil {
			t.Fatalf("failed to read test data: %v", err)
		}

		var response uberapi.GetTripResponse
		if err := json.Unmarshal(data, &response); err != nil {
			t.Fatalf("failed to unmarshal trip details: %v", err)
		}

		registry := &locations.Registry{
			Locations: []locations.Location{},
			NextID:    1,
		}

		lp := locations.NewProcessor(registry)

		trip, err := ProcessTrip(&response, lp)
		if err != nil {
			t.Fatalf("ProcessTrip() failed: %v", err)
		}

		if trip.PickupLocationID == "" {
			t.Error("expected pickup location ID to be set")
		}

		if trip.DropoffLocationID == "" {
			t.Error("expected dropoff location ID to be set")
		}

		if len(registry.Locations) != 2 {
			t.Errorf("expected 2 locations to be created, got %d", len(registry.Locations))
		}
	})
}

func TestProcessCanceledTrip(t *testing.T) {
	canceledTripData := `{
		"data": {
			"getTrip": {
				"trip": {
					"uuid": "canceled-trip-123",
					"status": "CANCELED",
					"driver": "",
					"waypoints": ["Test Location"],
					"beginTripTime": "Tue Jan 20 2026 19:50:26 GMT-0500 (Eastern Standard Time)",
					"dropoffTime": "",
					"fare": "R$0.00",
					"__typename": "Trip"
				},
				"mapURL": "",
				"receipt": {
					"distance": "0",
					"duration": "0 minutes",
					"vehicleType": "",
					"__typename": "Receipt"
				},
				"rating": "0",
				"reviewer": "",
				"__typename": "GetTripResult"
			}
		}
	}`

	var response uberapi.GetTripResponse
	if err := json.Unmarshal([]byte(canceledTripData), &response); err != nil {
		t.Fatalf("failed to unmarshal canceled trip: %v", err)
	}

	registry := &locations.Registry{
		Locations: []locations.Location{},
		NextID:    1,
	}

	lp := locations.NewProcessor(registry)

	trip, err := ProcessTrip(&response, lp)
	if err != nil {
		t.Fatalf("ProcessTrip() failed: %v", err)
	}

	if trip.Status != trips.StatusCanceled {
		t.Errorf("expected status %v, got %v", trips.StatusCanceled, trip.Status)
	}

	if trip.Fare != 0 {
		t.Errorf("expected fare 0 for canceled trip, got %v", trip.Fare)
	}

	if trip.Driver != "" {
		t.Errorf("expected empty driver for canceled trip, got %s", trip.Driver)
	}

	if trip.PickupLocationID != "" {
		t.Error("expected no location ID for canceled trip")
	}

	if trip.DropoffLocationID != "" {
		t.Error("expected no location ID for canceled trip")
	}
}
