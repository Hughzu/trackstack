package ports

import (
	"context"
	"time"
)

type UserProvider interface {
	VerifyCredentials(ctx context.Context, email string, password string) (string, error)
	UpdateLastLogin(ctx context.Context, userID string, timestamp string) error
}

type IssuedToken struct {
	Value     string
	ExpiresAt time.Time
}

type TokenIssuer interface {
	IssueToken(userID string) (IssuedToken, error)
}
