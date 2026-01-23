package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "ue",
	Short: "CLI tool for extracting and analyzing Uber trip data",
	Long:  `A CLI tool for extracting, analyzing, and exporting Uber trip data from Uber's GraphQL API.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		slog.SetDefault(slog.New(handler))
		return nil
	},
}

func init() {
	cobra.EnableCommandSorting = false

	RootCmd.AddCommand(LoginCmd)
	RootCmd.AddCommand(LogoutCmd)
	RootCmd.AddCommand(StatusCmd)
	RootCmd.AddCommand(TripsCmd)
	RootCmd.AddCommand(LocationsCmd)
}

func Execute() error {
	return RootCmd.Execute()
}
