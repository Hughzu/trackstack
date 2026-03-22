package ports

import (
	"context"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/calories/domain"
)

type TargetRepository interface {
	FindTarget(ctx context.Context, userID string) (domain.Target, bool, error)
	CreateTarget(ctx context.Context, target domain.Target) error
	UpdateTarget(ctx context.Context, target domain.Target) error
}

type LogRepository interface {
	CreateLog(ctx context.Context, log domain.Log) error
	DeleteLog(ctx context.Context, userID string, id string) (bool, error)
}

type DashboardRepository interface {
	GetNutritionSummaryByRange(ctx context.Context, userID string, from string, to string) (domain.NutritionSummary, error)
	GetLogsByRange(ctx context.Context, userID string, from string, to string, limit int) ([]domain.Log, error)
	GetRecentMeals(ctx context.Context, userID string, limit int) ([]domain.Log, error)
}
