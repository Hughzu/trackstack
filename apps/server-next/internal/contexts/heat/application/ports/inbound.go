package ports

import (
	"context"
	"time"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/heat/domain"
)

type RefillUseCase interface {
	GetRefills(ctx context.Context, userID string, from, to time.Time) ([]domain.Refill, error)
}
