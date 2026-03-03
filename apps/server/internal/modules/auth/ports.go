package auth

import "context"

type SessionStore interface {
	FindSessionByID(ctx context.Context, id string) (Session, error)
	InsertSession(ctx context.Context, session Session) error
	TouchSession(ctx context.Context, id string, lastSeenAt string, expiresAt string) error
	RotateOutSession(ctx context.Context, id string, revokedAt string, expiresAt string, rotatedAt string) error
	RevokeSession(ctx context.Context, id string, revokedAt string) error
}
