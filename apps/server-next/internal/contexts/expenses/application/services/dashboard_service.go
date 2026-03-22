package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/expenses/application/ports"
	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/expenses/domain"
)

type dashboardService struct {
	settingsRepo  ports.SettingsRepository
	dashboardRepo ports.DashboardRepository
	sheetManager  SheetManager
}

type DashboardServiceDeps struct {
	SettingsRepo  ports.SettingsRepository
	DashboardRepo ports.DashboardRepository
	SheetManager  SheetManager
}

func NewDashboardService(deps DashboardServiceDeps) ports.DashboardUseCase {
	return &dashboardService{
		settingsRepo:  deps.SettingsRepo,
		dashboardRepo: deps.DashboardRepo,
		sheetManager:  deps.SheetManager,
	}
}

func (s *dashboardService) GetDashboard(ctx context.Context, query ports.GetDashboardQuery) (domain.Dashboard, error) {
	if strings.TrimSpace(query.UserID) == "" {
		return domain.Dashboard{}, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}
	if query.HistoryLimit <= 0 {
		query.HistoryLimit = 50
	}

	settings, err := getOrCreateSettings(ctx, s.settingsRepo, query.UserID)
	if err != nil {
		return domain.Dashboard{}, err
	}

	sheet, err := s.sheetManager.GetOrCreateOpenSheet(ctx, query.UserID)
	if err != nil {
		return domain.Dashboard{}, err
	}

	snapshot, err := s.dashboardRepo.GetDashboardSnapshot(ctx, sheet.ID, query.HistoryLimit)
	if err != nil {
		return domain.Dashboard{}, err
	}
	if snapshot.SpentByCategory == nil {
		snapshot.SpentByCategory = map[domain.Category]float64{}
	}
	if snapshot.PendingChecklistItems == nil {
		snapshot.PendingChecklistItems = []domain.ChecklistItem{}
	}
	if snapshot.History == nil {
		snapshot.History = []domain.Entry{}
	}

	balance := domain.DashboardBalance{
		Remaining: settings.Income - snapshot.TotalSpent,
		Income:    settings.Income,
	}

	spent := domain.DashboardSpent{
		Fund:   snapshot.SpentByCategory[domain.CategoryFund],
		Fun:    snapshot.SpentByCategory[domain.CategoryFun],
		Future: snapshot.SpentByCategory[domain.CategoryFuture],
	}

	budget := domain.DashboardBudget{
		Fund:   int((settings.Income * float64(settings.RatioFund)) / 100),
		Fun:    int((settings.Income * float64(settings.RatioFun)) / 100),
		Future: int((settings.Income * float64(settings.RatioFuture)) / 100),
	}

	ratios := []domain.DashboardRatio{
		{
			Percent:    safePercent(spent.Fund, settings.Income),
			CategoryID: string(domain.CategoryFund),
			Label:      "Fund.",
			Value:      spent.Fund,
			Budget:     budget.Fund,
			Target:     settings.RatioFund,
			Over:       spent.Fund > float64(budget.Fund),
		},
		{
			Percent:    safePercent(spent.Fun, settings.Income),
			CategoryID: string(domain.CategoryFun),
			Label:      "Fun",
			Value:      spent.Fun,
			Budget:     budget.Fun,
			Target:     settings.RatioFun,
			Over:       spent.Fun > float64(budget.Fun),
		},
		{
			Percent:    safePercent(spent.Future, settings.Income),
			CategoryID: string(domain.CategoryFuture),
			Label:      "Future",
			Value:      spent.Future,
			Budget:     budget.Future,
			Target:     settings.RatioFuture,
			Over:       spent.Future > float64(budget.Future),
		},
	}

	return domain.Dashboard{
		PeriodKey:          sheet.PeriodKey,
		Balance:            balance,
		Spent:              spent,
		Budget:             budget,
		Ratios:             ratios,
		PendingObligations: snapshot.PendingChecklistItems,
		History:            snapshot.History,
	}, nil
}

func safePercent(spent float64, income float64) int {
	if income == 0 {
		return 0
	}

	return int((spent / income) * 100)
}
