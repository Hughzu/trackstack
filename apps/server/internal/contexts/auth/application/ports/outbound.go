package ports

import (
	"context"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/domain"
)

type UserProvider interface {
	VerifyCredentials(ctx context.Context, email string, password string) (string, error)
	UpdateLastLogin(ctx context.Context, userID string, timestamp string) error
}

type IssueTokenInput struct {
	UserID    string
	SessionID string
}

type IssuedToken struct {
	Value     string
	ExpiresAt time.Time
}

type TokenIssuer interface {
	IssueToken(input IssueTokenInput) (IssuedToken, error)
}

type CreateSessionInput struct {
	ID                string
	UserID            string
	TokenHash         string
	CreatedAt         time.Time
	ExpiresAt         time.Time
	AbsoluteExpiresAt time.Time
	UserAgentHash     *string
	IPPrefix          *string
}

type RotateSessionInput struct {
	SessionID         string
	TokenHash         string
	RotatedAt         time.Time
	NewSessionID      string
	NewTokenHash      string
	NewExpiresAt      time.Time
	NewLastSeenAt     time.Time
	AbsoluteExpiresAt time.Time
	UserAgentHash     *string
	IPPrefix          *string
}

type SessionRepository interface {
	CreateSession(ctx context.Context, input CreateSessionInput) (domain.RefreshSession, error)
	FindByTokenHash(ctx context.Context, tokenHash string) (domain.RefreshSession, error)
	RotateSession(ctx context.Context, input RotateSessionInput) (domain.RefreshSession, error)
	RevokeSession(ctx context.Context, sessionID string, revokedAt time.Time) error
	RevokeFamily(ctx context.Context, sessionID string, revokedAt time.Time) error
}

type RefreshTokenManager interface {
	GenerateToken() (string, error)
	HashToken(token string) string
}
