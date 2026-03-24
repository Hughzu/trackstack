package services

import (
	"context"
	"log/slog"
	"strings"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/users/application/ports"
	"github.com/Hughzu/trackstack/apps/server/internal/contexts/users/domain"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) VerifyCredentials(ctx context.Context, email string, password string) (domain.User, error) {
	normalized := strings.ToLower(strings.TrimSpace(email))
	if normalized == "" || password == "" {
		return domain.User{}, domain.ErrInvalidInput
	}

	user, err := s.repo.FindByEmail(ctx, normalized)
	if err != nil {
		return domain.User{}, err // Passes ErrNotFound up
	}

	if !VerifyPassword(password, user.PasswordHash) {
		slog.Warn("credentials rejected", "email", normalized)
		return domain.User{}, domain.ErrNotFound
	}

	return user, nil
}

func (s *UserService) UpdateLastLogin(ctx context.Context, userID string, lastLoginAt string) error {
	if strings.TrimSpace(userID) == "" || strings.TrimSpace(lastLoginAt) == "" {
		return domain.ErrInvalidInput
	}

	return s.repo.UpdateLastLogin(ctx, userID, lastLoginAt)
}
