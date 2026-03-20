package ports

import (
	"context"
	"time"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/heat/domain"
)

type RefillRepository interface {
	GetRefills(ctx context.Context, userID string, from, to time.Time) ([]domain.Refill, error)
	CreateRefill(ctx context.Context, refill domain.Refill) error
}
