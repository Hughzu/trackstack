package heat

import "context"

type RefillStore interface {
	ListByRange(ctx context.Context, userID string, from string, to string) ([]Refill, error)
	ListRecent(ctx context.Context, userID string, limit int, offset int) ([]Refill, error)
	GetLatest(ctx context.Context, userID string) (*Refill, error)
	GetSumByRange(ctx context.Context, userID string, from string, to string) (int, error)
	Create(ctx context.Context, refill Refill) error
	Delete(ctx context.Context, userID string, id string) (bool, error)
}
