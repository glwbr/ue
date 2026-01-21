package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

const (
	credentialsFileName = "credentials"
	configDirName       = ".ue"
)

type Credentials struct {
	Cookie    string    `json:"cookie"`
	Email     string    `json:"email"`
	LoginTime time.Time `json:"loginTime"`
}

func GetConfigDir() (string, error) {
	return getCredentialsDir()
}

func getCredentialsDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(usr.HomeDir, configDirName), nil
}

func getCredentialsPath() (string, error) {
	dir, err := getCredentialsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, credentialsFileName), nil
}

func Save(cookie, email string) error {
	dir, err := getCredentialsDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create credentials directory: %w", err)
	}

	creds := &Credentials{
		Cookie:    cookie,
		Email:     email,
		LoginTime: time.Now(),
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	path, err := getCredentialsPath()
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write credentials file: %w", err)
	}

	return nil
}

func Load() (*Credentials, error) {
	path, err := getCredentialsPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no credentials found. Please run 'ue auth login'")
		}
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal credentials: %w", err)
	}

	return &creds, nil
}

func Remove() error {
	path, err := getCredentialsPath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove credentials file: %w", err)
	}

	return nil
}
