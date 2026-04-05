package auth

import (
	"testing"
	"time"
)

func TestJWTManager_GenerateAndParseToken(t *testing.T) {
	manager := NewJWTManager("test-secret", time.Hour)

	token, err := manager.GenerateToken(42)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := manager.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken returned error: %v", err)
	}

	if claims.UserID != 42 {
		t.Fatalf("expected userID 42, got %d", claims.UserID)
	}
}

func TestJWTManager_ParseToken_InvalidToken(t *testing.T) {
	manager := NewJWTManager("test-secret", time.Hour)

	_, err := manager.ParseToken("this-is-not-a-valid-token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestJWTManager_ParseToken_ExpiredToken(t *testing.T) {
	manager := NewJWTManager("test-secret", -time.Hour)

	token, err := manager.GenerateToken(42)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	_, err = manager.ParseToken(token)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}