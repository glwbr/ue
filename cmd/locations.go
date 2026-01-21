package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"uber-extractor/internal/locations"
)

var LocationsCmd = &cobra.Command{
	Use:   "locations",
	Short: "List all saved locations",
	Long:  `List all locations that have been clustered and saved to the locations registry.`,
	Example: `  # List all saved locations
  ue locations

  # Show details about location clustering
  ue locations`,
	RunE: runLocations,
}

func runLocations(cmd *cobra.Command, args []string) error {
	registry, err := locations.Load()
	if err != nil {
		return fmt.Errorf("failed to load locations: %w", err)
	}

	if len(registry.Locations) == 0 {
		fmt.Println("No locations saved yet.")
		fmt.Println("\nRun 'ue trips' to fetch and cluster locations from your trip data.")
		return nil
	}

	tw := tabwriter.NewWriter(os.Stdout, 10, 8, 3, ' ', 0)
	defer tw.Flush()

	fmt.Fprintln(tw, "ID\t\tVISITS\tCOORDINATES\t\tADDRESS")
	for _, loc := range registry.Locations {
		coords := fmt.Sprintf("%.6f, %.6f", loc.AvgLat, loc.AvgLon)
		fmt.Fprintf(tw, "%s\t\t%d\t%s\t\t%s\n",
			loc.ID,
			loc.VisitCount,
			coords,
			truncate(loc.CanonicalAddress, 40),
		)
	}

	fmt.Printf("\nTotal locations: %d\n", len(registry.Locations))

	return nil
}
