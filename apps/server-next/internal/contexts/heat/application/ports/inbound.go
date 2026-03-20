package ports

import (
	"context"
	"time"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/heat/domain"
)

type RefillUseCase interface {
	GetRefills(ctx context.Context, query GetRefillsQuery) ([]domain.Refill, error)
	CreateRefill(ctx context.Context, command CreateRefillCommand) (domain.Refill, error)
	DeleteRefill(ctx context.Context, command DeleteRefillCommand) (bool, error)
}

var ErrInvalidInput = domain.ErrInvalidInput

type GetRefillsQuery struct {
	UserID string
	From   time.Time
	To     time.Time
}

type CreateRefillCommand struct {
	UserID      string
	Date        time.Time
	WeightKg    float64
	Bags        int
	Temperature *float64
}

type DeleteRefillCommand struct {
	UserID string
	ID     string
}
