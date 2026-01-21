package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"uber-extractor/internal/auth"
	"uber-extractor/internal/uberapi"
)

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Uber",
	Long:  `Login to Uber by providing your authentication cookie from browser dev tools.`,
	RunE:  runLogin,
}

func runLogin(cmd *cobra.Command, args []string) error {
	fmt.Println("Paste your Uber cookie (from browser dev tools):")
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	cookie, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read cookie: %w", err)
	}

	cookie = strings.TrimSpace(cookie)
	if cookie == "" {
		return fmt.Errorf("cookie cannot be empty")
	}

	fmt.Println("Validating session...")
	client := uberapi.NewClient(cookie)
	ctx := context.Background()

	resp, err := client.GetCurrentUser(ctx)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if resp.Data.CurrentUser == nil {
		return fmt.Errorf("failed to get user information")
	}

	user := &auth.User{
		FirstName: resp.Data.CurrentUser.FirstName,
		LastName:  resp.Data.CurrentUser.LastName,
		Email:     resp.Data.CurrentUser.Email,
	}

	if err := auth.Save(cookie, user.Email); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	fmt.Printf("Session valid. Logged in as: %s\n", user.FullName())
	return nil
}
