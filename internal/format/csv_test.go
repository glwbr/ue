package format

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"uber-extractor/internal/trips"
)

func TestCSVFormatter(t *testing.T) {
	formatter := &CSVFormatter{}

	now := time.Now()

	tripList := []trips.Trip{
		{
			UUID:           "trip-001",
			BeginTime:      now,
			EndTime:        now.Add(30 * time.Minute),
			Status:         trips.StatusCompleted,
			Fare:           25.50,
			Driver:         "John Doe",
			VehicleType:    "UberX",
			Distance:       8.63,
			Duration:       21,
			PickupAddress:  "123 Main St",
			DropoffAddress: "456 Oak Ave",
			PickupLat:      40.7128,
			PickupLon:      -74.0060,
			DropoffLat:     40.7200,
			DropoffLon:     -74.0100,
			Rating:         5,
		},
		{
			UUID:           "trip-002",
			BeginTime:      time.Time{},
			EndTime:        time.Time{},
			Status:         trips.StatusCanceled,
			Fare:           0.00,
			Driver:         "",
			VehicleType:    "",
			Distance:       0,
			Duration:       0,
			PickupAddress:  "",
			DropoffAddress: "",
			PickupLat:      0,
			PickupLon:      0,
			DropoffLat:     0,
			DropoffLon:     0,
			Rating:         0,
		},
	}

	var buf bytes.Buffer
	err := formatter.Format(&buf, tripList)
	if err != nil {
		t.Fatalf("Format() failed: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 3 {
		t.Errorf("expected 3 lines (header + 2 trips), got %d", len(lines))
	}

	expectedHeader := "UUID,BeginTime,EndTime,Status,Fare,Driver,VehicleType,Distance,Duration,PickupAddress,DropoffAddress,PickupLat,PickupLon,DropoffLat,DropoffLon,Rating"
	if lines[0] != expectedHeader {
		t.Errorf("header mismatch\nexpected: %s\ngot: %s", expectedHeader, lines[0])
	}

	if !strings.Contains(lines[1], "trip-001") {
		t.Errorf("expected first trip to contain UUID trip-001")
	}

	if !strings.Contains(lines[1], "25.50") {
		t.Errorf("expected first trip to contain fare 25.50")
	}

	if !strings.Contains(lines[1], "John Doe") {
		t.Errorf("expected first trip to contain driver name")
	}

	if !strings.Contains(lines[2], "trip-002") {
		t.Errorf("expected second trip to contain UUID trip-002")
	}

	if !strings.Contains(lines[2], "CANCELED") {
		t.Errorf("expected second trip to contain CANCELED status")
	}
}

func TestCSVFormatterEmpty(t *testing.T) {
	formatter := &CSVFormatter{}

	var buf bytes.Buffer
	err := formatter.Format(&buf, []trips.Trip{})
	if err != nil {
		t.Fatalf("Format() failed: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 1 {
		t.Errorf("expected 1 line (header only), got %d", len(lines))
	}

	expectedHeader := "UUID,BeginTime,EndTime,Status,Fare,Driver,VehicleType,Distance,Duration,PickupAddress,DropoffAddress,PickupLat,PickupLon,DropoffLat,DropoffLon,Rating"
	if lines[0] != expectedHeader {
		t.Errorf("header mismatch\nexpected: %s\ngot: %s", expectedHeader, lines[0])
	}
}

func TestFormatTime(t *testing.T) {
	t.Run("zero time", func(t *testing.T) {
		result := FormatTime(time.Time{})
		if result != "" {
			t.Errorf("expected empty string for zero time, got %s", result)
		}
	})

	t.Run("valid time", func(t *testing.T) {
		testTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
		result := FormatTime(testTime)
		if result != "2024-01-15T10:30:00Z" {
			t.Errorf("expected RFC3339 format, got %s", result)
		}
	})
}

func TestFormatDuration(t *testing.T) {
	t.Run("zero duration", func(t *testing.T) {
		result := FormatDuration(0)
		if result != "" {
			t.Errorf("expected empty string for zero duration, got %s", result)
		}
	})

	t.Run("integer duration", func(t *testing.T) {
		result := FormatDuration(21.0)
		if result != "21 minutes" {
			t.Errorf("expected '21 minutes', got %s", result)
		}
	})

	t.Run("fractional duration", func(t *testing.T) {
		result := FormatDuration(21.7)
		if result != "22 minutes" {
			t.Errorf("expected '22 minutes' (rounded), got %s", result)
		}
	})
}
