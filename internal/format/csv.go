package format

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"uber-extractor/internal/trips"
)

type CSVFormatter struct{}

func (f *CSVFormatter) Format(w io.Writer, tripList []trips.Trip) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	headers := []string{
		"UUID",
		"BeginTime",
		"EndTime",
		"Status",
		"Fare",
		"Driver",
		"VehicleType",
		"Distance",
		"Duration",
		"PickupAddress",
		"DropoffAddress",
		"PickupLat",
		"PickupLon",
		"DropoffLat",
		"DropoffLon",
		"Rating",
	}

	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, trip := range tripList {
		record := []string{
			trip.UUID,
			FormatTime(trip.BeginTime),
			FormatTime(trip.EndTime),
			trip.Status.String(),
			fmt.Sprintf("%.2f", trip.Fare),
			trip.Driver,
			trip.VehicleType,
			fmt.Sprintf("%.2f", trip.Distance),
			FormatDuration(trip.Duration),
			trip.PickupAddress,
			trip.DropoffAddress,
			fmt.Sprintf("%.6f", trip.PickupLat),
			fmt.Sprintf("%.6f", trip.PickupLon),
			fmt.Sprintf("%.6f", trip.DropoffLat),
			fmt.Sprintf("%.6f", trip.DropoffLon),
			fmt.Sprintf("%d", trip.Rating),
		}

		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func FormatDuration(d float64) string {
	if d == 0 {
		return ""
	}
	return fmt.Sprintf("%.0f minutes", d)
}
