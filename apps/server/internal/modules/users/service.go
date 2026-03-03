package users

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

var ErrInvalidInput = errors.New("invalid input")
var ErrNotFound = errors.New("not found")

type Service struct {
	store UsersStore
}

func NewService(store UsersStore) *Service {
	return &Service{store: store}
}

func (s *Service) FindByEmail(ctx context.Context, email string) (User, error) {
	normalized := strings.ToLower(strings.TrimSpace(email))
	if normalized == "" {
		return User{}, ErrInvalidInput
	}

	user, err := s.store.FindByEmail(ctx, normalized)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}

	return user, nil
}

func (s *Service) UpdateLastLogin(ctx context.Context, userID string, lastLoginAt string) error {
	if strings.TrimSpace(userID) == "" || strings.TrimSpace(lastLoginAt) == "" {
		return ErrInvalidInput
	}

	return s.store.UpdateLastLogin(ctx, userID, lastLoginAt)
}
