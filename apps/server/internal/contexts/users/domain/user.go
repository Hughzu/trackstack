package domain

import "errors"

var (
	ErrNotFound     = errors.New("user not found")
	ErrInvalidInput = errors.New("invalid input")
)

type User struct {
	ID             string
	Email          string
	PasswordHash   string
	SessionVersion int
	CreatedAt      string
	LastLoginAt    *string
}
