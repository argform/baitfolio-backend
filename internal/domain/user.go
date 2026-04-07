package domain

import (
	"time"
)

type User struct {
	UserID uint64
	Username string
	Email string
	PasswordHash string
	FirstName *string
	LastName *string
	About *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u User) Sanitized() User {
	u.PasswordHash = ""
	return u
}
