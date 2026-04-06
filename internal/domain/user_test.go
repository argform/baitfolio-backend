package domain

import "testing"

func TestUserSanitized(t *testing.T) {
	user := User{
		UserID: 1,
		Username: "tester",
		Email: "tester@example.com",
		PasswordHash: "secret-hash",
	}

	sanitized := user.Sanitized()

	if sanitized.PasswordHash != "" {
		t.Fatal("expected sanitized user without password hash")
	}

	if user.PasswordHash != "secret-hash" {
		t.Fatal("expected sanitized copy to not mutate original user")
	}
}
