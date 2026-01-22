package format

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"uber-extractor/internal/trips"
)

func TestJSONFormatter(t *testing.T) {
	formatter := &JSONFormatter{}

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
	}

	var buf bytes.Buffer
	err := formatter.Format(&buf, tripList)
	if err != nil {
		t.Fatalf("Format() failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "trip-001") {
		t.Errorf("expected output to contain trip UUID")
	}

	if !strings.Contains(output, "25.5") {
		t.Errorf("expected output to contain fare")
	}

	var result []trips.Trip
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("expected 1 trip, got %d", len(result))
	}

	if result[0].UUID != "trip-001" {
		t.Errorf("expected UUID trip-001, got %s", result[0].UUID)
	}

	if result[0].Fare != 25.50 {
		t.Errorf("expected fare 25.50, got %v", result[0].Fare)
	}
}

func TestJSONFormatterEmpty(t *testing.T) {
	formatter := &JSONFormatter{}

	var buf bytes.Buffer
	err := formatter.Format(&buf, []trips.Trip{})
	if err != nil {
		t.Fatalf("Format() failed: %v", err)
	}

	output := buf.String()

	if !strings.HasPrefix(output, "[") || !strings.HasSuffix(output, "]\n") {
		t.Errorf("expected JSON array format, got: %s", output)
	}

	var result []trips.Trip
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected 0 trips, got %d", len(result))
	}
}

func TestJSONFormatterMultipleTrips(t *testing.T) {
	formatter := &JSONFormatter{}

	tripList := []trips.Trip{
		{
			UUID:   "trip-001",
			Status: trips.StatusCompleted,
			Fare:   25.50,
			Driver: "John Doe",
		},
		{
			UUID:   "trip-002",
			Status: trips.StatusCompleted,
			Fare:   15.75,
			Driver: "Jane Smith",
		},
		{
			UUID:   "trip-003",
			Status: trips.StatusCanceled,
			Fare:   0.00,
			Driver: "",
		},
	}

	var buf bytes.Buffer
	err := formatter.Format(&buf, tripList)
	if err != nil {
		t.Fatalf("Format() failed: %v", err)
	}

	var result []trips.Trip
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("expected 3 trips, got %d", len(result))
	}

	if result[0].UUID != "trip-001" {
		t.Errorf("expected first UUID trip-001, got %s", result[0].UUID)
	}

	if result[1].UUID != "trip-002" {
		t.Errorf("expected second UUID trip-002, got %s", result[1].UUID)
	}

	if result[2].UUID != "trip-003" {
		t.Errorf("expected third UUID trip-003, got %s", result[2].UUID)
	}
}

func TestJSONFormatterPrettyPrint(t *testing.T) {
	formatter := &JSONFormatter{}

	tripList := []trips.Trip{
		{
			UUID:   "trip-001",
			Status: trips.StatusCompleted,
			Fare:   25.50,
		},
	}

	var buf bytes.Buffer
	err := formatter.Format(&buf, tripList)
	if err != nil {
		t.Fatalf("Format() failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "\n") {
		t.Error("expected pretty-printed JSON with newlines")
	}

	if !strings.Contains(output, "  ") {
		t.Error("expected pretty-printed JSON with indentation")
	}
}
