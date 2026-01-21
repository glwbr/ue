package datetime

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var lastPeriodRegex = regexp.MustCompile(`^(\d+)\s*d(ays?)?$`)

func ParseDateRange(from, to, last string) (time.Time, time.Time, error) {
	now := time.Now()

	if last != "" {
		if from != "" || to != "" {
			return time.Time{}, time.Time{}, fmt.Errorf("cannot use --last with --from or --to")
		}

		days, err := ParseLastPeriod(last)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}

		startTime := now.AddDate(0, 0, -days).Truncate(24 * time.Hour)
		endTime := now
		return startTime, endTime, nil
	}

	if from != "" && to != "" {
		startTime, err := time.Parse("2006-01-02", from)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid from date format: %w", err)
		}

		endTime, err := time.Parse("2006-01-02", to)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid to date format: %w", err)
		}

		if startTime.After(endTime) {
			return time.Time{}, time.Time{}, fmt.Errorf("from date cannot be after to date")
		}

		return startTime, endTime, nil
	}

	return time.Time{}, time.Time{}, fmt.Errorf("must specify both --from and --to, or use --last")
}

func ParseLastPeriod(s string) (int, error) {
	s = strings.TrimSpace(strings.ToLower(s))

	matches := lastPeriodRegex.FindStringSubmatch(s)
	if len(matches) < 2 {
		return 0, fmt.Errorf("invalid last period format. Expected format: 7d, 30d, 7days")
	}

	days, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("invalid number of days: %w", err)
	}

	return days, nil
}
