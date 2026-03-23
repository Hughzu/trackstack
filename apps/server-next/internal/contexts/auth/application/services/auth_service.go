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

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	email = strings.TrimSpace(email)
	if email == "" || password == "" {
		return "", domain.ErrInvalidInput
	}

	userID, err := s.userProvider.VerifyCredentials(ctx, email, password)
	if err != nil {
		return "", domain.ErrUnauthorized
	}

	token, err := s.tokenIssuer.IssueToken(userID)
	if err != nil {
		return "", err
	}

	asyncCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
	go func() {
		defer cancel()
		_ = s.userProvider.UpdateLastLogin(asyncCtx, userID, time.Now().UTC().Format(time.RFC3339))
	}()

	return token, nil
}
