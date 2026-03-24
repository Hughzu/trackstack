package ports

import (
	"context"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/domain"
)

type RefillRepository interface {
	GetRefills(ctx context.Context, userID string, from, to time.Time) ([]domain.Refill, error)
	CreateRefill(ctx context.Context, refill domain.Refill) error
	DeleteRefill(ctx context.Context, userID string, id string) (bool, error)
	GetDashboardStats(ctx context.Context, userID string, currentSeasonStart, currentSeasonEnd, lastSeasonStart, lastSeasonEnd time.Time) (*time.Time, int, int, error)
	ListRecentRefills(ctx context.Context, userID string, limit, offset int) ([]domain.Refill, error)
}
