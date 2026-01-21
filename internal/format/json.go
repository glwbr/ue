package format

import (
	"encoding/json"
	"io"

	"uber-extractor/internal/trips"
)

type JSONFormatter struct{}

func (f *JSONFormatter) Format(w io.Writer, tripList []trips.Trip) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(tripList)
}
