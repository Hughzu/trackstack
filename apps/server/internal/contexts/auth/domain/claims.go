package domain

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidInput       = errors.New("invalid input")
	ErrSessionNotFound    = errors.New("session not found")
	ErrRefreshTokenReused = errors.New("refresh token reused")
)

const TokenUseAccess = "access"

type SessionClaims struct {
	UserID    string `json:"userId"`
	SessionID string `json:"sessionId,omitempty"`
	TokenUse  string `json:"token_use"`
	jwt.RegisteredClaims
}

type RefreshSession struct {
	ID                string
	UserID            string
	TokenHash         string
	CreatedAt         time.Time
	ExpiresAt         time.Time
	RotatedAt         time.Time
	LastSeenAt        time.Time
	AbsoluteExpiresAt time.Time
	ParentID          *string
	RevokedAt         *time.Time
	UserAgentHash     *string
	IPPrefix          *string
}

func (s RefreshSession) IsRotated() bool {
	return s.RotatedAt.After(s.CreatedAt)
}

func (s RefreshSession) IsRevoked() bool {
	return s.RevokedAt != nil
}

func (s RefreshSession) IsExpired(now time.Time) bool {
	return now.After(s.ExpiresAt) || now.After(s.AbsoluteExpiresAt)
}
