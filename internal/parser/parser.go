package parser

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidTime     = errors.New("invalid time format")
	ErrInvalidDuration = errors.New("invalid duration format")
)

var (
	CurrencyRegex = regexp.MustCompile(`^[A-Z]{3}\$([\d.]+)$`)
	durationRegex = regexp.MustCompile(`^(\d+)\s+minutes?$`)
)

func Time(s string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05.000Z07:00",
		time.RFC1123,
		"Mon Jan 02 2006 15:04:05 GMT+0000 (Coordinated Universal Time)",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, ErrInvalidTime
}

func Duration(s string) (time.Duration, error) {
	if s == "" {
		return 0, ErrInvalidDuration
	}

	matches := durationRegex.FindStringSubmatch(s)
	if len(matches) < 2 {
		return 0, ErrInvalidDuration
	}

	minutes, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}

	return time.Duration(minutes) * time.Minute, nil
}

func Distance(s string) float64 {
	if s == "" {
		return 0
	}

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return val
}

func Fare(s string) float64 {
	if s == "" {
		return 0
	}

	matches := CurrencyRegex.FindStringSubmatch(s)
	if len(matches) < 2 {
		val, err := strconv.ParseFloat(strings.TrimLeft(s, "R$"), 64)
		if err != nil {
			return 0
		}
		return val
	}

	val, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0
	}
	return val
}

func Rating(s string) int {
	if s == "" {
		return 0
	}

	rating, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return rating
}

func ExtractCoordinates(mapURL string, markerIndex int) (lat, lon float64) {
	if mapURL == "" {
		return 0, 0
	}

	parsedURL, err := url.Parse(mapURL)
	if err != nil {
		return 0, 0
	}

	query := parsedURL.Query()
	markers := query["marker"]
	if len(markers) == 0 {
		return 0, 0
	}

	if markerIndex >= len(markers) {
		return 0, 0
	}

	marker := markers[markerIndex]

	marker = strings.ReplaceAll(marker, "%24", "$")
	marker = strings.ReplaceAll(marker, "%3A", ":")

	latRegex := regexp.MustCompile(`lat:([\d.-]+)`)
	lonRegex := regexp.MustCompile(`lng:([\d.-]+)`)

	latMatches := latRegex.FindStringSubmatch(marker)
	lonMatches := lonRegex.FindStringSubmatch(marker)

	if len(latMatches) < 2 || len(lonMatches) < 2 {
		return 0, 0
	}

	lat, _ = strconv.ParseFloat(latMatches[1], 64)
	lon, _ = strconv.ParseFloat(lonMatches[1], 64)

	return lat, lon
}
