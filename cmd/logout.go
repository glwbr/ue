package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"uber-extractor/internal/auth"
)

var LogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Uber",
	Long:  `Remove stored credentials from ~/.ue/credentials.`,
	RunE:  runLogout,
}

func runLogout(cmd *cobra.Command, args []string) error {
	if err := auth.Remove(); err != nil {
		return fmt.Errorf("failed to remove credentials: %w", err)
	}

	fmt.Println("Logged out successfully")
	return nil
}
