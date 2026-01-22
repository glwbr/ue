package locations

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestHaversineDistance(t *testing.T) {
	tests := []struct {
		name       string
		lat1       float64
		lon1       float64
		lat2       float64
		lon2       float64
		wantMeters float64
		tolerance  float64
	}{
		{
			name:       "same point",
			lat1:       41.4089,
			lon1:       -75.6624,
			lat2:       41.4089,
			lon2:       -75.6624,
			wantMeters: 0,
			tolerance:  1,
		},
		{
			name:       "Dunder Mifflin to Michael's Condo",
			lat1:       41.4089,
			lon1:       -75.6624,
			lat2:       41.4120,
			lon2:       -75.6580,
			wantMeters: 540,
			tolerance:  100,
		},
		{
			name:       "long distance",
			lat1:       41.4089,
			lon1:       -75.6624,
			lat2:       35.2000,
			lon2:       -80.8500,
			wantMeters: 825000,
			tolerance:  50000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HaversineDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			diff := got - tt.wantMeters
			if diff < 0 {
				diff = -diff
			}
			if diff > tt.tolerance {
				t.Errorf("HaversineDistance() = %v meters, want %v±%v meters", got, tt.wantMeters, tt.tolerance)
			}
		})
	}
}

func TestFindOrCreateLocation(t *testing.T) {
	registry := &Registry{
		Locations: []Location{},
		NextID:    1,
	}

	processor := NewProcessor(registry)

	t.Run("create new location", func(t *testing.T) {
		locID := processor.FindOrCreateLocation("123 Test Street", 41.4089, -75.6624)
		if locID != "loc-1" {
			t.Errorf("expected loc-1, got %s", locID)
		}

		if len(registry.Locations) != 1 {
			t.Errorf("expected 1 location, got %d", len(registry.Locations))
		}

		loc := registry.Locations[0]
		if loc.VisitCount != 1 {
			t.Errorf("expected visit count 1, got %d", loc.VisitCount)
		}
	})

	t.Run("find existing location by address", func(t *testing.T) {
		locID := processor.FindOrCreateLocation("123 test street", 41.4089, -75.6624)
		if locID != "loc-1" {
			t.Errorf("expected loc-1, got %s", locID)
		}

		loc := registry.Locations[0]
		if loc.VisitCount != 2 {
			t.Errorf("expected visit count 2, got %d", loc.VisitCount)
		}
	})

	t.Run("cluster nearby locations within threshold", func(t *testing.T) {
		registry.Locations[0].VisitCount = 1
		registry.Locations[0].AvgLat = 41.4089
		registry.Locations[0].AvgLon = -75.6624

		locID := processor.FindOrCreateLocation("124 Test Street", 41.408905, -75.662405)
		if locID != "loc-1" {
			t.Errorf("expected loc-1 (clustered), got %s", locID)
		}

		loc := registry.Locations[0]
		if loc.VisitCount != 2 {
			t.Errorf("expected visit count 2 after clustering, got %d", loc.VisitCount)
		}

		if len(loc.AddressVariants) == 0 {
			t.Error("expected address variant to be added")
		}
	})

	t.Run("create separate location outside threshold", func(t *testing.T) {
		registry.NextID = 2
		locID := processor.FindOrCreateLocation("456 Far Away Street", 41.4500, -75.6200)
		if locID != "loc-2" {
			t.Errorf("expected loc-2, got %s", locID)
		}

		if len(registry.Locations) != 2 {
			t.Errorf("expected 2 locations, got %d", len(registry.Locations))
		}
	})

	t.Run("ignore zero coordinates", func(t *testing.T) {
		locID := processor.FindOrCreateLocation("Invalid Location", 0, 0)
		if locID != "" {
			t.Errorf("expected empty string for zero coordinates, got %s", locID)
		}
	})
}

func TestFindOrCreateLocationWithRegistry(t *testing.T) {
	data, err := os.ReadFile("testdata/locations_clustering.json")
	if err != nil {
		t.Fatalf("failed to read test data: %v", err)
	}

	var registry Registry
	if err := json.Unmarshal(data, &registry); err != nil {
		t.Fatalf("failed to unmarshal registry: %v", err)
	}

	processor := NewProcessor(&registry)

	t.Run("find exact address match", func(t *testing.T) {
		locID := processor.FindOrCreateLocation("1725 slough avenue - scranton business park - scranton - pa 18505", 41.4089, -75.6624)
		if locID != "loc-1" {
			t.Errorf("expected loc-1, got %s", locID)
		}

		loc := registry.Locations[0]
		if loc.VisitCount != 56 {
			t.Errorf("expected visit count 56, got %d", loc.VisitCount)
		}
	})

	t.Run("find address variant match", func(t *testing.T) {
		locID := processor.FindOrCreateLocation("michael scott residence - 123 kellum court - scranton - pa 18508", 41.4120, -75.6580)
		if locID != "loc-2" {
			t.Errorf("expected loc-2, got %s", locID)
		}

		loc := registry.Locations[1]
		if loc.VisitCount != 69 {
			t.Errorf("expected visit count 69, got %d", loc.VisitCount)
		}
	})

	t.Run("cluster new location near existing", func(t *testing.T) {
		initialCount := registry.Locations[0].VisitCount
		locID := processor.FindOrCreateLocation("1725 Slough Avenue", 41.408901, -75.662401)
		if locID != "loc-1" {
			t.Errorf("expected loc-1 (clustered), got %s", locID)
		}

		loc := registry.Locations[0]
		if loc.VisitCount != initialCount+1 {
			t.Errorf("expected visit count %d, got %d", initialCount+1, loc.VisitCount)
		}
	})

	t.Run("create new location far from existing", func(t *testing.T) {
		registry.NextID = 10
		locID := processor.FindOrCreateLocation("New York Times Square", 40.7580, -73.9855)
		if locID != "loc-10" {
			t.Errorf("expected loc-10, got %s", locID)
		}

		if len(registry.Locations) != 10 {
			t.Errorf("expected 10 locations, got %d", len(registry.Locations))
		}

		newLoc := registry.Locations[9]
		if newLoc.VisitCount != 1 {
			t.Errorf("expected visit count 1, got %d", newLoc.VisitCount)
		}

		if time.Since(newLoc.FirstSeen) > time.Second {
			t.Error("expected FirstSeen to be recent")
		}

		if time.Since(newLoc.LastSeen) > time.Second {
			t.Error("expected LastSeen to be recent")
		}
	})
}

func TestUpdateLocation(t *testing.T) {
	registry := &Registry{
		Locations: []Location{
			{
				ID:               "loc-1",
				CanonicalAddress: "123 Test Street",
				AddressVariants:  []string{},
				AvgLat:           41.4089,
				AvgLon:           -75.6624,
				VisitCount:       1,
				FirstSeen:        time.Now(),
				LastSeen:         time.Now(),
			},
		},
		NextID: 2,
	}

	processor := NewProcessor(registry)

	initialLat := registry.Locations[0].AvgLat
	initialLon := registry.Locations[0].AvgLon

	processor.FindOrCreateLocation("123 TEST STREET", 41.4090, -75.6625)

	loc := registry.Locations[0]

	if loc.VisitCount != 2 {
		t.Errorf("expected visit count 2, got %d", loc.VisitCount)
	}

	if loc.AvgLat == initialLat {
		t.Error("expected average latitude to change")
	}

	if loc.AvgLon == initialLon {
		t.Error("expected average longitude to change")
	}
}
