package users

import "context"

type UsersStore interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	UpdateLastLogin(ctx context.Context, userID string, lastLoginAt string) error
}
