package trips

import (
	"encoding/json"
	"testing"
)

func TestParseTripStatus(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  TripStatus
	}{
		{
			name:  "completed uppercase",
			input: "COMPLETED",
			want:  StatusCompleted,
		},
		{
			name:  "completed mixed case",
			input: "Completed",
			want:  StatusCompleted,
		},
		{
			name:  "completed lowercase",
			input: "completed",
			want:  StatusCompleted,
		},
		{
			name:  "canceled uppercase",
			input: "CANCELED",
			want:  StatusCanceled,
		},
		{
			name:  "canceled mixed case",
			input: "Canceled",
			want:  StatusCanceled,
		},
		{
			name:  "unknown status",
			input: "UNKNOWN",
			want:  StatusUnknown,
		},
		{
			name:  "invalid status",
			input: "INVALID",
			want:  StatusUnknown,
		},
		{
			name:  "status with whitespace",
			input: "  COMPLETED  ",
			want:  StatusCompleted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseTripStatus(tt.input)
			if got != tt.want {
				t.Errorf("ParseTripStatus(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestTripStatusString(t *testing.T) {
	tests := []struct {
		name  string
		input TripStatus
		want  string
	}{
		{
			name:  "completed",
			input: StatusCompleted,
			want:  "COMPLETED",
		},
		{
			name:  "canceled",
			input: StatusCanceled,
			want:  "CANCELED",
		},
		{
			name:  "unknown",
			input: StatusUnknown,
			want:  "UNKNOWN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.String()
			if got != tt.want {
				t.Errorf("TripStatus.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTripStatusMarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input TripStatus
		want  string
	}{
		{
			name:  "completed",
			input: StatusCompleted,
			want:  `"COMPLETED"`,
		},
		{
			name:  "canceled",
			input: StatusCanceled,
			want:  `"CANCELED"`,
		},
		{
			name:  "unknown",
			input: StatusUnknown,
			want:  `"UNKNOWN"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.MarshalJSON()
			if err != nil {
				t.Errorf("TripStatus.MarshalJSON() error = %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("TripStatus.MarshalJSON() = %q, want %q", string(got), tt.want)
			}
		})
	}
}

func TestTripStatusUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  TripStatus
	}{
		{
			name:  "completed",
			input: `"COMPLETED"`,
			want:  StatusCompleted,
		},
		{
			name:  "canceled",
			input: `"CANCELED"`,
			want:  StatusCanceled,
		},
		{
			name:  "unknown",
			input: `"UNKNOWN"`,
			want:  StatusUnknown,
		},
		{
			name:  "invalid",
			input: `"INVALID"`,
			want:  StatusUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got TripStatus
			err := got.UnmarshalJSON([]byte(tt.input))
			if err != nil {
				t.Errorf("TripStatus.UnmarshalJSON() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("TripStatus.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTripWithStatusJSON(t *testing.T) {
	trip := Trip{
		UUID:   "test-uuid",
		Status: StatusCompleted,
		Fare:   10.50,
	}

	data, err := json.Marshal(trip)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var unmarshaledTrip Trip
	err = json.Unmarshal(data, &unmarshaledTrip)
	if err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if unmarshaledTrip.Status != StatusCompleted {
		t.Errorf("Status = %v, want %v", unmarshaledTrip.Status, StatusCompleted)
	}

	var result map[string]interface{}
	json.Unmarshal(data, &result)
	if result["status"] != "COMPLETED" {
		t.Errorf("JSON status = %v, want COMPLETED", result["status"])
	}
}
