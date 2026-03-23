package ports

import "context"

type UserProvider interface {
	VerifyCredentials(ctx context.Context, email string, password string) (string, error)
	UpdateLastLogin(ctx context.Context, userID string, timestamp string) error
}

type TokenIssuer interface {
	IssueToken(userID string) (string, error)
}
