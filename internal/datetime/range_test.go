package datetime

import (
	"testing"
	"time"
)

func TestParseDateRange(t *testing.T) {
	t.Run("last period format", func(t *testing.T) {
		start, end, err := ParseDateRange("", "", "7d")
		if err != nil {
			t.Fatalf("ParseDateRange() failed: %v", err)
		}

		if start.IsZero() {
			t.Error("expected non-zero start time")
		}

		if end.IsZero() {
			t.Error("expected non-zero end time")
		}

		duration := end.Sub(start)
		expectedDuration := 7 * 24 * time.Hour

		if duration < expectedDuration-24*time.Hour || duration > expectedDuration+24*time.Hour {
			t.Errorf("expected duration ~7 days (within 1 day), got %v", duration)
		}
	})

	t.Run("last 30 days", func(t *testing.T) {
		start, end, err := ParseDateRange("", "", "30d")
		if err != nil {
			t.Fatalf("ParseDateRange() failed: %v", err)
		}

		duration := end.Sub(start)
		expectedDuration := 30 * 24 * time.Hour

		if duration < expectedDuration-24*time.Hour || duration > expectedDuration+24*time.Hour {
			t.Errorf("expected duration ~30 days (within 1 day), got %v", duration)
		}
	})

	t.Run("from and to dates", func(t *testing.T) {
		start, end, err := ParseDateRange("2024-01-01", "2024-01-15", "")
		if err != nil {
			t.Fatalf("ParseDateRange() failed: %v", err)
		}

		expectedStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		expectedEnd := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

		if !start.Equal(expectedStart) {
			t.Errorf("expected start %v, got %v", expectedStart, start)
		}

		if !end.Equal(expectedEnd) {
			t.Errorf("expected end %v, got %v", expectedEnd, end)
		}
	})

	t.Run("from and to dates across months", func(t *testing.T) {
		start, end, err := ParseDateRange("2024-01-25", "2024-02-05", "")
		if err != nil {
			t.Fatalf("ParseDateRange() failed: %v", err)
		}

		expectedStart := time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC)
		expectedEnd := time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC)

		if !start.Equal(expectedStart) {
			t.Errorf("expected start %v, got %v", expectedStart, start)
		}

		if !end.Equal(expectedEnd) {
			t.Errorf("expected end %v, got %v", expectedEnd, end)
		}
	})

	t.Run("last with from and to fails", func(t *testing.T) {
		_, _, err := ParseDateRange("2024-01-01", "2024-01-15", "7d")
		if err == nil {
			t.Error("expected error when using --last with --from/--to")
		}
	})

	t.Run("from without to fails", func(t *testing.T) {
		_, _, err := ParseDateRange("2024-01-01", "", "")
		if err == nil {
			t.Error("expected error when using only --from")
		}
	})

	t.Run("to without from fails", func(t *testing.T) {
		_, _, err := ParseDateRange("", "2024-01-15", "")
		if err == nil {
			t.Error("expected error when using only --to")
		}
	})

	t.Run("from after to fails", func(t *testing.T) {
		_, _, err := ParseDateRange("2024-01-15", "2024-01-01", "")
		if err == nil {
			t.Error("expected error when --from is after --to")
		}
	})

	t.Run("invalid from date format", func(t *testing.T) {
		_, _, err := ParseDateRange("2024/01/01", "2024-01-15", "")
		if err == nil {
			t.Error("expected error for invalid date format")
		}
	})

	t.Run("invalid to date format", func(t *testing.T) {
		_, _, err := ParseDateRange("2024-01-01", "01-15-2024", "")
		if err == nil {
			t.Error("expected error for invalid date format")
		}
	})

	t.Run("no parameters fails", func(t *testing.T) {
		_, _, err := ParseDateRange("", "", "")
		if err == nil {
			t.Error("expected error when no parameters provided")
		}
	})
}

func TestParseLastPeriod(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantDays int
		wantErr  bool
	}{
		{
			name:     "7 days short",
			input:    "7d",
			wantDays: 7,
			wantErr:  false,
		},
		{
			name:     "30 days short",
			input:    "30d",
			wantDays: 30,
			wantErr:  false,
		},
		{
			name:     "365 days short",
			input:    "365d",
			wantDays: 365,
			wantErr:  false,
		},
		{
			name:     "7 days long",
			input:    "7days",
			wantDays: 7,
			wantErr:  false,
		},
		{
			name:     "1 day long",
			input:    "1day",
			wantDays: 1,
			wantErr:  false,
		},
		{
			name:     "with spaces",
			input:    "7 days",
			wantDays: 7,
			wantErr:  false,
		},
		{
			name:     "uppercase",
			input:    "7D",
			wantDays: 7,
			wantErr:  false,
		},
		{
			name:     "missing number",
			input:    "d",
			wantDays: 0,
			wantErr:  true,
		},
		{
			name:     "missing d",
			input:    "7",
			wantDays: 0,
			wantErr:  true,
		},
		{
			name:     "invalid format",
			input:    "7w",
			wantDays: 0,
			wantErr:  true,
		},
		{
			name:     "empty string",
			input:    "",
			wantDays: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			days, err := ParseLastPeriod(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLastPeriod(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if days != tt.wantDays {
				t.Errorf("ParseLastPeriod(%q) = %d, want %d", tt.input, days, tt.wantDays)
			}
		})
	}
}
