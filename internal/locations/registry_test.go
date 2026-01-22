package locations

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadRegistry(t *testing.T) {
	t.Run("load existing file", func(t *testing.T) {
		registry, err := Load("testdata/locations_clustering.json")
		if err != nil {
			t.Fatalf("Load() failed: %v", err)
		}

		if len(registry.Locations) != 9 {
			t.Errorf("expected 9 locations, got %d", len(registry.Locations))
		}

		if registry.NextID != 10 {
			t.Errorf("expected NextID 10, got %d", registry.NextID)
		}

		loc := registry.Locations[0]
		if loc.ID != "loc-1" {
			t.Errorf("expected loc-1, got %s", loc.ID)
		}
	})

	t.Run("load non-existent file returns empty registry", func(t *testing.T) {
		registry, err := Load("/tmp/non-existent-locations-xyz.json")
		if err != nil {
			t.Fatalf("Load() failed: %v", err)
		}

		if len(registry.Locations) != 0 {
			t.Errorf("expected 0 locations, got %d", len(registry.Locations))
		}

		if registry.NextID != 1 {
			t.Errorf("expected NextID 1, got %d", registry.NextID)
		}
	})

	t.Run("load invalid JSON fails", func(t *testing.T) {
		tmpDir := t.TempDir()
		invalidFile := filepath.Join(tmpDir, "invalid.json")

		if err := os.WriteFile(invalidFile, []byte("invalid json"), 0644); err != nil {
			t.Fatalf("failed to write invalid file: %v", err)
		}

		_, err := Load(invalidFile)
		if err == nil {
			t.Error("expected error for invalid JSON, got nil")
		}
	})
}

func TestSaveRegistry(t *testing.T) {
	t.Run("save new registry", func(t *testing.T) {
		tmpDir := t.TempDir()
		testPath := filepath.Join(tmpDir, "test-locations.json")

		registry := &Registry{
			Locations: []Location{
				{
					ID:               "loc-1",
					CanonicalAddress: "123 Test Street",
					AddressVariants:  []string{},
					AvgLat:           41.4089,
					AvgLon:           -75.6624,
					VisitCount:       5,
					FirstSeen:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
					LastSeen:         time.Date(2024, 1, 20, 15, 30, 0, 0, time.UTC),
				},
			},
			NextID: 2,
		}

		if err := Save(registry, testPath); err != nil {
			t.Fatalf("Save() failed: %v", err)
		}

		data, err := os.ReadFile(testPath)
		if err != nil {
			t.Fatalf("failed to read saved file: %v", err)
		}

		var loaded Registry
		if err := json.Unmarshal(data, &loaded); err != nil {
			t.Fatalf("failed to unmarshal saved data: %v", err)
		}

		if len(loaded.Locations) != 1 {
			t.Errorf("expected 1 location, got %d", len(loaded.Locations))
		}

		if loaded.Locations[0].ID != "loc-1" {
			t.Errorf("expected loc-1, got %s", loaded.Locations[0].ID)
		}

		if loaded.NextID != 2 {
			t.Errorf("expected NextID 2, got %d", loaded.NextID)
		}
	})

	t.Run("save creates directory if needed", func(t *testing.T) {
		tmpDir := t.TempDir()
		testPath := filepath.Join(tmpDir, "subdir", "nested", "locations.json")

		registry := &Registry{
			Locations: []Location{},
			NextID:    1,
		}

		if err := Save(registry, testPath); err != nil {
			t.Fatalf("Save() failed: %v", err)
		}

		if _, err := os.Stat(testPath); err != nil {
			t.Errorf("expected file to exist at %s, got error: %v", testPath, err)
		}
	})
}

func TestRegistryJSONRoundTrip(t *testing.T) {
	original := &Registry{
		Locations: []Location{
			{
				ID:               "loc-1",
				CanonicalAddress: "1725 Slough Avenue",
				AddressVariants:  []string{"Dunder Mifflin Office"},
				AvgLat:           41.4089,
				AvgLon:           -75.6624,
				VisitCount:       55,
				FirstSeen:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				LastSeen:         time.Date(2024, 1, 21, 23, 59, 59, 0, time.UTC),
			},
			{
				ID:               "loc-2",
				CanonicalAddress: "123 Kellum Court",
				AddressVariants:  []string{},
				AvgLat:           41.4120,
				AvgLon:           -75.6580,
				VisitCount:       68,
				FirstSeen:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				LastSeen:         time.Date(2024, 1, 21, 23, 59, 59, 0, time.UTC),
			},
		},
		NextID: 3,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal() failed: %v", err)
	}

	var loaded Registry
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Unmarshal() failed: %v", err)
	}

	if len(loaded.Locations) != len(original.Locations) {
		t.Errorf("expected %d locations, got %d", len(original.Locations), len(loaded.Locations))
	}

	if loaded.NextID != original.NextID {
		t.Errorf("expected NextID %d, got %d", original.NextID, loaded.NextID)
	}

	for i, loc := range loaded.Locations {
		if loc.ID != original.Locations[i].ID {
			t.Errorf("location %d: expected ID %s, got %s", i, original.Locations[i].ID, loc.ID)
		}
		if loc.CanonicalAddress != original.Locations[i].CanonicalAddress {
			t.Errorf("location %d: expected address %s, got %s", i, original.Locations[i].CanonicalAddress, loc.CanonicalAddress)
		}
		if loc.VisitCount != original.Locations[i].VisitCount {
			t.Errorf("location %d: expected visit count %d, got %d", i, original.Locations[i].VisitCount, loc.VisitCount)
		}
	}
}
