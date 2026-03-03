package calories

import "context"

type CaloriesStore interface {
	GetTarget(ctx context.Context, userID string) (Target, error)
	CreateTarget(ctx context.Context, target Target) error
	UpdateTarget(ctx context.Context, target Target) error

	CreateLog(ctx context.Context, log Log) error
	DeleteLog(ctx context.Context, userID string, id string) (bool, error)
}
