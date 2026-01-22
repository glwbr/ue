package auth

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"
)

func TestCredentialsJSON(t *testing.T) {
	t.Run("credentials marshal and unmarshal", func(t *testing.T) {
		testTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

		creds := &Credentials{
			Cookie:    "test-cookie-12345",
			Email:     "test@example.com",
			LoginTime: testTime,
		}

		data, err := json.Marshal(creds)
		if err != nil {
			t.Fatalf("Marshal() failed: %v", err)
		}

		var unmarshaled Credentials
		if err := json.Unmarshal(data, &unmarshaled); err != nil {
			t.Fatalf("Unmarshal() failed: %v", err)
		}

		if unmarshaled.Cookie != creds.Cookie {
			t.Errorf("expected cookie %s, got %s", creds.Cookie, unmarshaled.Cookie)
		}

		if unmarshaled.Email != creds.Email {
			t.Errorf("expected email %s, got %s", creds.Email, unmarshaled.Email)
		}

		if !unmarshaled.LoginTime.Equal(creds.LoginTime) {
			t.Errorf("expected login time %v, got %v", creds.LoginTime, unmarshaled.LoginTime)
		}
	})
}

func TestGetConfigDir(t *testing.T) {
	t.Run("returns user home directory with .ue", func(t *testing.T) {
		dir, err := GetConfigDir()
		if err != nil {
			t.Fatalf("GetConfigDir() failed: %v", err)
		}

		if dir == "" {
			t.Error("expected non-empty config directory")
		}

		if filepath.Base(dir) != ".ue" {
			t.Errorf("expected .ue directory, got %s", filepath.Base(dir))
		}

		parentDir := filepath.Dir(dir)
		if parentDir == "" || parentDir == "." {
			t.Error("expected parent directory to be user home")
		}
	})
}
