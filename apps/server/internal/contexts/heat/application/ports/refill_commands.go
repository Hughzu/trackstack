package ports

import (
	"context"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/domain"
)

type RefillCreator interface {
	Create(ctx context.Context, refill domain.Refill) error
}

type RefillDeleter interface {
	Delete(ctx context.Context, userID string, id string) (bool, error)
}
