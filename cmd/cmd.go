package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "ue",
	Short: "CLI tool for extracting and analyzing Uber trip data",
	Long:  `A CLI tool for extracting, analyzing, and exporting Uber trip data from Uber's GraphQL API.`,
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
