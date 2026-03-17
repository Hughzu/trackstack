package ports

import (
	"context"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/domain"
)

type RefillRangeLister interface {
	ListByRange(ctx context.Context, userID string, from string, to string) ([]domain.Refill, error)
}

type RecentRefillLister interface {
	ListRecent(ctx context.Context, userID string, limit int, offset int) ([]domain.Refill, error)
}

type LatestRefillGetter interface {
	GetLatest(ctx context.Context, userID string) (*domain.Refill, error)
}

type RefillRangeSummarizer interface {
	GetSumByRange(ctx context.Context, userID string, from string, to string) (int, error)
}
