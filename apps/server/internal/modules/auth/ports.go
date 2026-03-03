package auth

import "context"

type SessionStore interface {
	InsertSession(ctx context.Context, session Session) error
	RevokeSession(ctx context.Context, id string, revokedAt string) error
}
