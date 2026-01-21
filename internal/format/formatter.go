package format

import (
	"fmt"
	"io"

	"uber-extractor/internal/trips"
)

type Formatter interface {
	Format(w io.Writer, tripList []trips.Trip) error
}

func GetFormatter(format string) (Formatter, error) {
	switch format {
	case "json":
		return &JSONFormatter{}, nil
	case "csv":
		return &CSVFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
