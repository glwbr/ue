package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/spf13/cobra"

	"uber-extractor/internal/auth"
	"uber-extractor/internal/datetime"
	"uber-extractor/internal/format"
	"uber-extractor/internal/locations"
	"uber-extractor/internal/parser"
	"uber-extractor/internal/transform"
	"uber-extractor/internal/trips"
	"uber-extractor/internal/uberapi"
)

var (
	fromDate   string
	toDate     string
	lastPeriod string
	output     string
	summary    bool

	subtitleRegex = regexp.MustCompile(`([A-Za-z]+ \d+) • (\d+:\d+ [AP]M)`)
)

var TripsCmd = &cobra.Command{
	Use:   "trips",
	Short: "Fetch and display trip data from Uber",
	Long:  `Fetch trip data from Uber GraphQL API for a specified date range and display it.`,
	RunE:  runTrips,
	Example: `  # Fetch last 7 days of trips in JSON format
  ue trips --last 7d

  # Fetch trips for a date range in CSV format
  ue trips --from 2024-01-01 --to 2024-01-31 --output csv

  # Show summary without fetching full details
  ue trips --last 30d --summary`,
}

func init() {
	TripsCmd.Flags().StringVar(&fromDate, "from", "", "Start date in YYYY-MM-DD format")
	TripsCmd.Flags().StringVar(&toDate, "to", "", "End date in YYYY-MM-DD format")
	TripsCmd.Flags().StringVar(&lastPeriod, "last", "", "Period in days (e.g., 7d, 3d, 30d)")
	TripsCmd.Flags().StringVarP(&output, "output", "o", "json", "Output format: json, csv (default: json)")
	TripsCmd.Flags().BoolVar(&summary, "summary", false, "Show summary without fetching details")
}

func runTrips(cmd *cobra.Command, args []string) error {
	creds, err := auth.Load()
	if err != nil {
		return err
	}

	startTime, endTime, err := datetime.ParseDateRange(fromDate, toDate, lastPeriod)
	if err != nil {
		return err
	}

	client := uberapi.NewClient(creds.Cookie)
	ctx := context.Background()

	if summary {
		return runSummary(ctx, client, startTime, endTime)
	}

	return runFetch(ctx, client, startTime, endTime)
}

func runSummary(ctx context.Context, client *uberapi.Client, start, end time.Time) error {
	slog.Info("Fetching trip summary", "date_range", fmt.Sprintf("%s to %s", start.Format("2006-01-02"), end.Format("2006-01-02")))

	activities, _, err := client.GetActivities(ctx, start.Unix()*1000, end.Unix()*1000, "")
	if err != nil {
		return fmt.Errorf("failed to fetch activities: %w", err)
	}

	tripCount := len(activities.Data.Activities.Past.Activities)
	slog.Info("Parsing trip summary", "count", tripCount)

	tripSummary := trips.TripSummary{
		Count:      tripCount,
		Activities: activities.Data.Activities.Past.Activities,
	}

	for _, a := range tripSummary.Activities {
		tripSummary.TotalFare += parser.Fare(a.Description)
	}

	slog.Info("Summary calculated", "total_trips", tripSummary.Count, "total_fare", fmt.Sprintf("%.2f", tripSummary.TotalFare))

	fmt.Printf("Found %d trips between %s and %s\n", tripSummary.Count, start.Format("2006-01-02"), end.Format("2006-01-02"))
	fmt.Printf("Total fare: %.2f\n", tripSummary.TotalFare)
	fmt.Println("\nRecent trips:")
	fmt.Println("DATE\tTIME\tFARE\tDESTINATION")

	for _, activity := range tripSummary.Activities {
		timeStr := ""
		if parts := parseSubtitle(activity.Subtitle); len(parts) == 2 {
			timeStr = parts[1]
		}
		fmt.Printf("%s\t%s\t%s\t%s\n", activity.Title, timeStr, activity.Description, truncate(activity.Title, 30))
	}

	return nil
}

func runFetch(ctx context.Context, client *uberapi.Client, start, end time.Time) error {
	slog.Info("Starting trip fetch", "date_range", fmt.Sprintf("%s to %s", start.Format("2006-01-02"), end.Format("2006-01-02")))

	registry, err := locations.Load()
	if err != nil {
		return fmt.Errorf("failed to load locations: %w", err)
	}

	slog.Info("Locations loaded", "count", len(registry.Locations))

	lp := locations.NewProcessor(registry)

	var allTrips []trips.Trip
	pageToken := ""
	pageCount := 0

	for {
		slog.Info("Fetching activities", "page", pageCount+1)

		activities, nextPageToken, err := client.GetActivities(ctx, start.Unix()*1000, end.Unix()*1000, pageToken)
		if err != nil {
			return fmt.Errorf("failed to fetch activities: %w", err)
		}

		activitiesList := activities.Data.Activities.Past.Activities
		slog.Info("Parsing activities", "count", len(activitiesList))

		for i, activity := range activitiesList {
			processedInPage := i + 1
			currentTotal := len(allTrips) + processedInPage
			estimatedTotal := len(allTrips) + len(activitiesList)
			percentage := float64(currentTotal*100) / float64(estimatedTotal)

			slog.Info("Processing trip", "current", currentTotal, "total_estimate", estimatedTotal, "progress", fmt.Sprintf("%.0f%%", percentage))

			slog.Debug("Fetching trip details", "uuid", activity.UUID)

			tripResponse, err := client.GetTrip(ctx, activity.UUID)
			if err != nil {
				slog.Warn("Failed to fetch trip details", "uuid", activity.UUID, "error", err)
				continue
			}

			slog.Debug("Transforming trip", "uuid", activity.UUID)

			trip, err := transform.ProcessTrip(tripResponse, lp)
			if err != nil {
				slog.Warn("Failed to process trip", "uuid", activity.UUID, "error", err)
				continue
			}

			allTrips = append(allTrips, trip)
		}

		pageCount++
		pageToken = nextPageToken

		if pageToken == "" {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	slog.Info("Trips processed successfully", "total", len(allTrips), "pages", pageCount)

	if err := locations.Save(lp.Registry()); err != nil {
		slog.Warn("Failed to save locations", "error", err)
	} else {
		configDir, err := auth.GetConfigDir()
		if err != nil {
			slog.Info("Locations saved", "count", len(lp.Registry().Locations))
		} else {
			path := filepath.Join(configDir, "locations.json")
			slog.Info("Locations saved", "count", len(lp.Registry().Locations), "path", path)
		}
	}

	slog.Info("Formatting output", "format", output, "destination", "stdout")

	f, err := format.GetFormatter(output)
	if err != nil {
		return err
	}

	return f.Format(os.Stdout, allTrips)
}

func parseSubtitle(subtitle string) []string {
	matches := subtitleRegex.FindStringSubmatch(subtitle)
	if len(matches) < 3 {
		return []string{subtitle, ""}
	}
	return []string{matches[1], matches[2]}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
