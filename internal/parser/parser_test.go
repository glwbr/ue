package parser

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		validate func(t *testing.T, result time.Time)
	}{
		{
			name:    "RFC3339 format",
			input:   "2024-01-15T10:30:00Z",
			wantErr: false,
			validate: func(t *testing.T, result time.Time) {
				if result.Year() != 2024 || result.Month() != 1 || result.Day() != 15 {
					t.Errorf("unexpected time: %v", result)
				}
			},
		},
		{
			name:    "long format with timezone",
			input:   "Mon Jan 15 2026 21:29:14 GMT+0000 (Coordinated Universal Time)",
			wantErr: false,
			validate: func(t *testing.T, result time.Time) {
				if result.Year() != 2026 || result.Month() != 1 || result.Day() != 15 {
					t.Errorf("unexpected time: %v", result)
				}
			},
		},
		{
			name:    "invalid format",
			input:   "invalid time",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Time(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Time() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestDuration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		{
			name:    "16 minutes",
			input:   "16 minutes",
			want:    16 * time.Minute,
			wantErr: false,
		},
		{
			name:    "1 minute",
			input:   "1 minute",
			want:    1 * time.Minute,
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid format",
			input:   "invalid",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Duration(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Duration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFare(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  float64
	}{
		{
			name:  "Brazilian Real format",
			input: "R$10.84",
			want:  10.84,
		},
		{
			name:  "USD format",
			input: "USD$25.50",
			want:  25.50,
		},
		{
			name:  "simple number",
			input: "15.99",
			want:  15.99,
		},
		{
			name:  "empty string",
			input: "",
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Fare(tt.input); got != tt.want {
				t.Errorf("Fare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRating(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "valid rating",
			input: "5",
			want:  5,
		},
		{
			name:  "zero rating",
			input: "0",
			want:  0,
		},
		{
			name:  "empty string",
			input: "",
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Rating(tt.input); got != tt.want {
				t.Errorf("Rating() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractCoordinates(t *testing.T) {
	tests := []struct {
		name        string
		mapURL      string
		markerIndex int
		wantLat     float64
		wantLon     float64
	}{
		{
			name:        "valid map URL",
			mapURL:      "https://static-maps.uber.com/map?marker=lat%3A-12.26071%24lng%3A-38.9452%24icon%3Apickup.png",
			markerIndex: 0,
			wantLat:     -12.26071,
			wantLon:     -38.9452,
		},
		{
			name:        "empty URL",
			mapURL:      "",
			markerIndex: 0,
			wantLat:     0,
			wantLon:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lat, lon := ExtractCoordinates(tt.mapURL, tt.markerIndex)
			if lat != tt.wantLat || lon != tt.wantLon {
				t.Errorf("ExtractCoordinates() = (%v, %v), want (%v, %v)", lat, lon, tt.wantLat, tt.wantLon)
			}
		})
	}
}
