package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
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
		return fmt.Errorf("not logged in. Please run 'ue login': %w", err)
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
	activities, _, err := client.GetActivities(ctx, start.Unix()*1000, end.Unix()*1000, "")
	if err != nil {
		return fmt.Errorf("failed to fetch activities: %w", err)
	}

	tripSummary := trips.TripSummary{
		Count:      len(activities.Data.Activities.Past.Activities),
		Activities: activities.Data.Activities.Past.Activities,
	}

	for _, a := range tripSummary.Activities {
		tripSummary.TotalFare += parser.Fare(a.Description)
	}

	fmt.Printf("Found %d trips between %s and %s\n", tripSummary.Count, start.Format("2006-01-02"), end.Format("2006-01-02"))
	fmt.Printf("Total fare: %.2f\n", tripSummary.TotalFare)
	fmt.Println("\nRecent trips:")
	fmt.Println("DATE\tTIME\tFARE\tDESTINATION")

	for _, activity := range tripSummary.Activities {
		parts := parseSubtitle(activity.Subtitle)
		timeStr := ""
		if len(parts) == 2 {
			timeStr = parts[1]
		}
		fmt.Printf("%s\t%s\t%s\t%s\n", activity.Title, timeStr, activity.Description, truncate(activity.Title, 30))
	}

	return nil
}

func runFetch(ctx context.Context, client *uberapi.Client, start, end time.Time) error {
	registry, err := locations.Load()
	if err != nil {
		return fmt.Errorf("failed to load locations: %w", err)
	}

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

		slog.Info("Found trips", "count", len(activities.Data.Activities.Past.Activities))

		for _, activity := range activities.Data.Activities.Past.Activities {
			slog.Debug("Fetching trip details", "uuid", activity.UUID)

			tripResponse, err := client.GetTrip(ctx, activity.UUID)
			if err != nil {
				slog.Warn("Failed to fetch trip details", "uuid", activity.UUID, "error", err)
				continue
			}

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

	slog.Info("Trips fetched successfully", "total", len(allTrips), "pages", pageCount)

	if err := locations.Save(lp.Registry()); err != nil {
		slog.Warn("Failed to save locations", "error", err)
	} else {
		slog.Info("Locations saved", "count", len(lp.Registry().Locations))
	}

	f, err := format.GetFormatter(output)
	if err != nil {
		return err
	}

	return f.Format(os.Stdout, allTrips)
}

func parseSubtitle(subtitle string) []string {
	re := regexp.MustCompile(`([A-Za-z]+ \d+) • (\d+:\d+ [AP]M)`)
	matches := re.FindStringSubmatch(subtitle)
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
