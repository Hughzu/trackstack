package calories

import (
	"context"
	"strings"
	"sync"
	"time"
)

type DashboardSummary struct {
	Consumed int    `json:"consumed"`
	Protein  int    `json:"protein"`
	Carbs    int    `json:"carbs"`
	Fat      int    `json:"fat"`
	Target   Target `json:"target"`
}

type DashboardViewModel struct {
	Summary     DashboardSummary `json:"summary"`
	Logs        []Log            `json:"logs"`
	RecentMeals []Log            `json:"recentMeals"`
}

type GetDashboardRequest struct {
	UserID      string
	RecentLimit int
	LogsLimit   int
}

func (s *Service) GetDashboard(ctx context.Context, req GetDashboardRequest) (DashboardViewModel, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return DashboardViewModel{}, ErrInvalidInput
	}

	if req.RecentLimit <= 0 {
		req.RecentLimit = 8
	}
	if req.LogsLimit <= 0 {
		req.LogsLimit = 50
	}

	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	from := startOfDay.Format(time.RFC3339)
	to := endOfDay.Format(time.RFC3339)

	var target Target
	var targetErr error
	var summaryData Summary
	logs := []Log{}
	recentMeals := []Log{}

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		target, targetErr = s.GetTarget(ctx, GetTargetRequest{UserID: req.UserID})
	}()

	go func() {
		defer wg.Done()
		summary, err := s.store.GetSummaryByRange(ctx, req.UserID, from, to)
		if err == nil {
			summaryData = summary
		}
	}()

	go func() {
		defer wg.Done()
		result, err := s.store.GetLogsByRangeLimited(ctx, req.UserID, from, to, req.LogsLimit)
		if err == nil {
			logs = result
		}
	}()

	go func() {
		defer wg.Done()
		result, err := s.store.GetRecentLogs(ctx, req.UserID, req.RecentLimit)
		if err == nil {
			recentMeals = result
		}
	}()

	wg.Wait()
	if targetErr != nil {
		return DashboardViewModel{}, targetErr
	}

	summary := DashboardSummary{
		Consumed: summaryData.Calories,
		Protein:  summaryData.ProteinG,
		Carbs:    summaryData.CarbsG,
		Fat:      summaryData.FatG,
		Target:   target,
	}

	return DashboardViewModel{
		Summary:     summary,
		Logs:        logs,
		RecentMeals: recentMeals,
	}, nil
}
