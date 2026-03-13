package expenses

import (
	"context"
	"strings"
)

func safePercent(spent, income float64) int {
	if income == 0 {
		return 0
	}
	return int((spent / income) * 100)
}

func (s *Service) GetDashboard(ctx context.Context, req GetCurrentSheetRequest) (ViewDashboard, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return ViewDashboard{}, ErrInvalidInput
	}
	if req.HistoryLimit <= 0 {
		req.HistoryLimit = 50
	}

	settingsRes, err := s.GetSettings(ctx, GetSettingsRequest{UserID: req.UserID})
	if err != nil {
		return ViewDashboard{}, err
	}
	settings := settingsRes.Settings

	sheet, err := s.getOrCreateOpenSheet(ctx, req.UserID)
	if err != nil {
		return ViewDashboard{}, err
	}

	totalSpent, _ := s.store.GetTotalSpentBySheet(ctx, sheet.ID)
	spentByCategory, _ := s.store.GetSpentByCategory(ctx, sheet.ID)
	pendingChecklist, err := s.store.GetPendingChecklistItems(ctx, sheet.ID)
	if err != nil {
		pendingChecklist = []ChecklistItem{}
	}
	history, err := s.store.GetRecentSheetHistory(ctx, sheet.ID, req.HistoryLimit, 0)
	if err != nil {
		history = []Entry{}
	}

	balance := DashboardBalance{
		Remaining: settings.Income - totalSpent,
		Income:    settings.Income,
	}
	spent := DashboardSpent{
		Fund:   spentByCategory[CategoryFund],
		Fun:    spentByCategory[CategoryFun],
		Future: spentByCategory[CategoryFuture],
	}
	budget := DashboardBudget{
		Fund:   int((balance.Income * float64(settings.RatioFund)) / 100),
		Fun:    int((balance.Income * float64(settings.RatioFun)) / 100),
		Future: int((balance.Income * float64(settings.RatioFuture)) / 100),
	}

	ratios := []DashboardRatio{
		{
			Percent:    safePercent(spent.Fund, balance.Income),
			CategoryId: "fund",
			Label:      "Fund.",
			Value:      spent.Fund,
			Budget:     budget.Fund,
			Target:     settings.RatioFund,
			Over:       spent.Fund > float64(budget.Fund),
		},
		{
			Percent:    safePercent(spent.Fun, balance.Income),
			CategoryId: "fun",
			Label:      "Fun",
			Value:      spent.Fun,
			Budget:     budget.Fun,
			Target:     settings.RatioFun,
			Over:       spent.Fun > float64(budget.Fun),
		},
		{
			Percent:    safePercent(spent.Future, balance.Income),
			CategoryId: "future",
			Label:      "Future",
			Value:      spent.Future,
			Budget:     budget.Future,
			Target:     settings.RatioFuture,
			Over:       spent.Future > float64(budget.Future),
		},
	}

	return ViewDashboard{
		PeriodKey:          sheet.PeriodKey,
		Balance:            balance,
		Spent:              spent,
		Budget:             budget,
		Ratios:             ratios,
		PendingObligations: pendingChecklist,
		History:            history,
	}, nil
}
