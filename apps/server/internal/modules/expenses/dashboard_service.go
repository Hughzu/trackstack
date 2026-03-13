package expenses

import (
	"context"
	"log/slog"
	"strings"
	"sync"
	"time"
)

func safePercent(spent, income float64) int {
	if income == 0 {
		return 0
	}
	return int((spent / income) * 100)
}

func (s *Service) GetDashboard(ctx context.Context, req GetCurrentSheetRequest) (ViewDashboard, error) {
	start := time.Now()
	defer logExpensesTiming(ctx, "dashboard.total", start, nil)

	if strings.TrimSpace(req.UserID) == "" {
		return ViewDashboard{}, ErrInvalidInput
	}
	if req.HistoryLimit <= 0 {
		req.HistoryLimit = 50
	}

	settingsStart := time.Now()
	settings, err := s.getOrCreateSettings(ctx, req.UserID)
	logExpensesTiming(ctx, "dashboard.settings", settingsStart, err)
	if err != nil {
		return ViewDashboard{}, err
	}

	sheetStart := time.Now()
	sheet, err := s.getOrCreateOpenSheet(ctx, req.UserID)
	logExpensesTiming(ctx, "dashboard.sheet", sheetStart, err)
	if err != nil {
		return ViewDashboard{}, err
	}

	totalSpent := 0.0
	spentByCategory := map[Category]float64{}
	pendingChecklist := []ChecklistItem{}
	history := []Entry{}

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		stepStart := time.Now()
		result, queryErr := s.store.GetTotalSpentBySheet(ctx, sheet.ID)
		logExpensesTiming(ctx, "dashboard.total_spent", stepStart, queryErr)
		if queryErr == nil {
			totalSpent = result
		}
	}()

	go func() {
		defer wg.Done()
		stepStart := time.Now()
		result, queryErr := s.store.GetSpentByCategory(ctx, sheet.ID)
		logExpensesTiming(ctx, "dashboard.spent_by_category", stepStart, queryErr)
		if queryErr == nil && result != nil {
			spentByCategory = result
		}
	}()

	go func() {
		defer wg.Done()
		stepStart := time.Now()
		result, queryErr := s.store.GetPendingChecklistItems(ctx, sheet.ID)
		logExpensesTiming(ctx, "dashboard.pending_checklist", stepStart, queryErr)
		if queryErr == nil {
			pendingChecklist = result
		}
	}()

	go func() {
		defer wg.Done()
		stepStart := time.Now()
		result, queryErr := s.store.GetRecentSheetHistory(ctx, sheet.ID, req.HistoryLimit, 0)
		logExpensesTiming(ctx, "dashboard.history", stepStart, queryErr)
		if queryErr == nil {
			history = result
		}
	}()

	wg.Wait()

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

func logExpensesTiming(ctx context.Context, step string, start time.Time, err error) {
	attrs := []any{"step", step, "duration", time.Since(start)}
	if err != nil {
		attrs = append(attrs, "error", err)
	}
	slog.DebugContext(ctx, "expenses timing", attrs...)
}
