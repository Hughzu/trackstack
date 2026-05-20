package targetservice

import (
	"context"
	"errors"
)

const (
	DefaultTargetCalories     = 2300
	DefaultTargetProteinGrams = 120
)

var ErrUserIDRequired = errors.New("user id is required")

type Target struct {
	ID                 string
	UserID             string
	TargetCalories     int
	TargetProteinGrams int
}

type Repository interface {
	FindByUserID(ctx context.Context, userID string) (Target, bool, error)
	Create(ctx context.Context, target Target) error
}

type Service struct {
	repo  Repository
	newID func() string
}

func NewService(repo Repository, newID func() string) *Service {
	panic("TODO")
}

func (s *Service) GetOrCreateTarget(ctx context.Context, userID string) (Target, error) {
	panic("TODO")
}
