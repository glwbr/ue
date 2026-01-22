package locations

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"uber-extractor/internal/auth"
)

type Registry struct {
	Locations []Location `json:"locations"`
	NextID    int        `json:"nextID"`
}

type Location struct {
	ID               string    `json:"id"`
	CanonicalAddress string    `json:"canonicalAddress"`
	AddressVariants  []string  `json:"addressVariants"`
	AvgLat           float64   `json:"avgLat"`
	AvgLon           float64   `json:"avgLon"`
	VisitCount       int       `json:"visitCount"`
	FirstSeen        time.Time `json:"firstSeen"`
	LastSeen         time.Time `json:"lastSeen"`
}

func getDefaultPath() string {
	dir, err := auth.GetConfigDir()
	if err != nil {
		return "locations.json"
	}
	return filepath.Join(dir, "locations.json")
}

func Load(path ...string) (*Registry, error) {
	p := getDefaultPath()
	if len(path) > 0 && path[0] != "" {
		p = path[0]
	}

	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return &Registry{Locations: []Location{}, NextID: 1}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load registry: %w", err)
	}

	var reg Registry
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return &reg, nil
}

func Save(reg *Registry, path ...string) error {
	p := getDefaultPath()
	if len(path) > 0 && path[0] != "" {
		p = path[0]
	}

	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(reg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return os.WriteFile(p, data, 0644)
}
