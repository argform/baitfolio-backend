package auth

import "testing"

func TestHashPassword(t *testing.T) {
	password := "supersecret123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if hash == "" {
		t.Fatal("expected non-empty hash")
	}

	if hash == password {
		t.Fatal("hash must not equal original password")
	}
}

func TestComparePassword_Success(t *testing.T) {
	password := "supersecret123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	ok := ComparePassword(hash, password)
	if !ok {
		t.Fatal("expected ComparePassword to return true for correct password")
	}
}

func TestComparePassword_Failure(t *testing.T) {
	password := "supersecret123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	ok := ComparePassword(hash, "wrong-password")
	if ok {
		t.Fatal("expected ComparePassword to return false for incorrect password")
	}
}