package ports

import (
	"context"
	"time"
)

type LoginResult struct {
	AccessToken string
	TokenType   string
	ExpiresAt   time.Time
	UserID      string
}

type AuthUseCase interface {
	Login(ctx context.Context, email string, password string) (LoginResult, error)
}
