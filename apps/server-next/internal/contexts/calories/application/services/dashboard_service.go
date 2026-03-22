package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/calories/application/ports"
	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/calories/domain"
)

type dashboardService struct {
	dashboardRepo ports.DashboardRepository
	targetRepo    ports.TargetRepository
}

func NewDashboardService(dashboardRepo ports.DashboardRepository, targetRepo ports.TargetRepository) ports.DashboardUseCase {
	return &dashboardService{dashboardRepo: dashboardRepo, targetRepo: targetRepo}
}

func (s *dashboardService) GetDashboard(ctx context.Context, query ports.GetDashboardQuery) (domain.Dashboard, error) {
	if strings.TrimSpace(query.UserID) == "" {
		return domain.Dashboard{}, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}
	if query.RecentLimit <= 0 {
		query.RecentLimit = 8
	}
	if query.LogsLimit <= 0 {
		query.LogsLimit = 50
	}

	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	from := startOfDay.Format(time.RFC3339)
	to := endOfDay.Format(time.RFC3339)

	target, err := NewTargetService(s.targetRepo).GetTarget(ctx, ports.GetTargetQuery{UserID: query.UserID})
	if err != nil {
		return domain.Dashboard{}, err
	}

	summary, err := s.dashboardRepo.GetNutritionSummaryByRange(ctx, query.UserID, from, to)
	if err != nil {
		return domain.Dashboard{}, err
	}

	logs, err := s.dashboardRepo.GetLogsByRange(ctx, query.UserID, from, to, query.LogsLimit)
	if err != nil {
		return domain.Dashboard{}, err
	}

	recentMeals, err := s.dashboardRepo.GetRecentMeals(ctx, query.UserID, query.RecentLimit)
	if err != nil {
		return domain.Dashboard{}, err
	}

	return domain.Dashboard{
		Summary: domain.DashboardNutritionSummary{
			ConsumedCalories: summary.Calories,
			ProteinGrams:     summary.ProteinGrams,
			CarbGrams:        summary.CarbGrams,
			FatGrams:         summary.FatGrams,
			Target:           target,
		},
		Logs:        logs,
		RecentMeals: recentMeals,
	}, nil
}
