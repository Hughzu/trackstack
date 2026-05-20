package userlookup

import (
	"context"
	"errors"
)

var (
	ErrInvalidEmail = errors.New("email is required")
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID    string
	Email string
}

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (User, bool, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	panic("TODO")
}

func (s *UserService) LookupByEmail(ctx context.Context, email string) (User, error) {
	panic("TODO")
}
