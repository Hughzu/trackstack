package heat

import "context"

type RefillStore interface {
	ListByRange(ctx context.Context, userID string, from string, to string) ([]Refill, error)
	Create(ctx context.Context, userID string, input CreateRefillInput) (Refill, error)
}
