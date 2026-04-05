package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, error := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), error
}

func ComparePassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidatePassword(password string, lenght int) (bool, string) {
	panic("TODO: check password safety")
}