package auth

import "testing"

func TestUserFullName(t *testing.T) {
	t.Run("full name with first and last", func(t *testing.T) {
		user := &User{
			FirstName: "Michael",
			LastName:  "Scott",
			Email:     "michael.scott@dundermifflin.com",
		}

		fullName := user.FullName()
		if fullName != "Michael Scott" {
			t.Errorf("expected 'Michael Scott', got %s", fullName)
		}
	})

	t.Run("only first name", func(t *testing.T) {
		user := &User{
			FirstName: "Dwight",
			LastName:  "",
			Email:     "dwight@beetfarm.com",
		}

		fullName := user.FullName()
		if fullName != "Dwight " {
			t.Errorf("expected 'Dwight ', got %s", fullName)
		}
	})

	t.Run("only last name", func(t *testing.T) {
		user := &User{
			FirstName: "",
			LastName:  "Schrute",
			Email:     "dwight@beetfarm.com",
		}

		fullName := user.FullName()
		if fullName != " Schrute" {
			t.Errorf("expected ' Schrute', got %s", fullName)
		}
	})

	t.Run("no names returns email", func(t *testing.T) {
		user := &User{
			FirstName: "",
			LastName:  "",
			Email:     "creed@thoughts.com",
		}

		fullName := user.FullName()
		if fullName != "creed@thoughts.com" {
			t.Errorf("expected email 'creed@thoughts.com', got %s", fullName)
		}
	})

	t.Run("empty user returns empty string", func(t *testing.T) {
		user := &User{
			FirstName: "",
			LastName:  "",
			Email:     "",
		}

		fullName := user.FullName()
		if fullName != "" {
			t.Errorf("expected empty string, got %s", fullName)
		}
	})
}
