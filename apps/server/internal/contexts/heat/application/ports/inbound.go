package ports

import (
	"context"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/domain"
)

type RefillUseCase interface {
	GetRefills(ctx context.Context, query GetRefillsQuery) ([]domain.Refill, error)
	CreateRefill(ctx context.Context, command CreateRefillCommand) (domain.Refill, error)
	DeleteRefill(ctx context.Context, command DeleteRefillCommand) (bool, error)
	GetDashboard(ctx context.Context, query GetDashboardQuery) (DashboardViewModel, error)
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

type SeasonSnapshot struct {
	SeasonLabel      string `json:"seasonLabel"`
	SeasonToDate     int    `json:"seasonToDate"`
	LastSeasonToDate int    `json:"lastSeasonToDate"`
	Delta            int    `json:"delta"`
	DeltaPct         *int   `json:"deltaPct"`
}

type DashboardViewModel struct {
	DaysSinceRefill int             `json:"daysSinceRefill"`
	SeasonSnapshot  SeasonSnapshot  `json:"seasonSnapshot"`
	History         []domain.Refill `json:"history"`
}

type GetDashboardQuery struct {
	UserID string
	Page   int
	Limit  int
}
