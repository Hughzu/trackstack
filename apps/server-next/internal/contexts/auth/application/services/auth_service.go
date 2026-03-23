package services

import (
	"context"
	"strings"
	"time"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/auth/application/ports"
	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/auth/domain"
)

type AuthService struct {
	userProvider ports.UserProvider
	tokenIssuer  ports.TokenIssuer
}

var _ ports.AuthUseCase = (*AuthService)(nil)

func NewAuthService(userProvider ports.UserProvider, tokenIssuer ports.TokenIssuer) *AuthService {
	return &AuthService{
		userProvider: userProvider,
		tokenIssuer:  tokenIssuer,
	}
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (ports.LoginResult, error) {
	email = strings.TrimSpace(email)
	if email == "" || password == "" {
		return ports.LoginResult{}, domain.ErrInvalidInput
	}

	userID, err := s.userProvider.VerifyCredentials(ctx, email, password)
	if err != nil {
		return ports.LoginResult{}, domain.ErrUnauthorized
	}

	issuedToken, err := s.tokenIssuer.IssueToken(userID)
	if err != nil {
		return ports.LoginResult{}, err
	}

	asyncCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
	go func() {
		defer cancel()
		_ = s.userProvider.UpdateLastLogin(asyncCtx, userID, time.Now().UTC().Format(time.RFC3339))
	}()

	return ports.LoginResult{
		AccessToken: issuedToken.Value,
		TokenType:   "Bearer",
		ExpiresAt:   issuedToken.ExpiresAt,
		UserID:      userID,
	}, nil
}
