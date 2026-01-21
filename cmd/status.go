package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"uber-extractor/internal/auth"
	"uber-extractor/internal/uberapi"
)

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check authentication status",
	Long:  `Validate stored credentials and show login status.`,
	RunE:  runStatus,
}

func runStatus(cmd *cobra.Command, args []string) error {
	creds, err := auth.Load()
	if err != nil {
		fmt.Println("Not logged in")
		return nil
	}

	fmt.Println("Validating session...")
	client := uberapi.NewClient(creds.Cookie)
	ctx := context.Background()

	resp, err := client.GetCurrentUser(ctx)
	if err != nil {
		fmt.Printf("Authentication failed: %v\n", err)
		fmt.Println("Please run 'ue login' to re-authenticate")
		return nil
	}

	if resp.Data.CurrentUser == nil {
		fmt.Println("Failed to get user information")
		return nil
	}

	user := &auth.User{
		FirstName: resp.Data.CurrentUser.FirstName,
		LastName:  resp.Data.CurrentUser.LastName,
		Email:     resp.Data.CurrentUser.Email,
	}

	fmt.Printf("Logged in as: %s\n", user.FullName())
	return nil
}
