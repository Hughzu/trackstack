package ports

import (
	"context"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/calories/domain"
)

var ErrInvalidInput = domain.ErrInvalidInput

type TargetUseCase interface {
	GetTarget(ctx context.Context, query GetTargetQuery) (domain.Target, error)
	UpdateTarget(ctx context.Context, command UpdateTargetCommand) (domain.Target, error)
}

type LogUseCase interface {
	AddLog(ctx context.Context, command AddLogCommand) (domain.Log, error)
	DeleteLog(ctx context.Context, command DeleteLogCommand) (bool, error)
}

type DashboardUseCase interface {
	GetDashboard(ctx context.Context, query GetDashboardQuery) (domain.Dashboard, error)
}

type GetTargetQuery struct {
	UserID string
}

type UpdateTargetCommand struct {
	UserID             string
	TargetCalories     int
	TargetProteinGrams int
	TargetCarbGrams    *int
	TargetFatGrams     *int
}

type AddLogCommand struct {
	UserID       string
	Calories     int
	ProteinGrams int
	CarbGrams    *int
	FatGrams     *int
	Title        *string
	Date         *string
	Time         *string
}

type DeleteLogCommand struct {
	UserID string
	ID     string
}

type GetDashboardQuery struct {
	UserID      string
	RecentLimit int
	LogsLimit   int
}
