package heat

import "context"

type RefillStore interface {
	ListByRange(ctx context.Context, userID string, from string, to string) ([]Refill, error)
	Create(ctx context.Context, refill Refill) error
	Delete(ctx context.Context, userID string, id string) (bool, error)
}
