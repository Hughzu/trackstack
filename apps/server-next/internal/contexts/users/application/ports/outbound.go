package ports

import (
	"context"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/users/domain"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	UpdateLastLogin(ctx context.Context, userID string, timestamp string) error
}
