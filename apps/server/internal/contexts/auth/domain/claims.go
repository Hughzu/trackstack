package domain

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidInput = errors.New("invalid input")
)

type SessionClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}
