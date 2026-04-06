package ports

import (
	"context"
	"time"
)

type LoginCommand struct {
	Email     string
	Password  string
	UserAgent string
	IP        string
}

type RefreshCommand struct {
	RefreshToken string
	UserAgent    string
	IP           string
}

type LogoutCommand struct {
	RefreshToken string
}

type AuthTokens struct {
	AccessToken      string
	TokenType        string
	ExpiresAt        time.Time
	UserID           string
	SessionID        string
	RefreshToken     string
	RefreshExpiresAt time.Time
}

type AuthUseCase interface {
	Login(ctx context.Context, command LoginCommand) (AuthTokens, error)
	Refresh(ctx context.Context, command RefreshCommand) (AuthTokens, error)
	Logout(ctx context.Context, command LogoutCommand) error
}
