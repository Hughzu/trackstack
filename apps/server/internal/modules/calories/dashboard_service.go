package calories

import (
	"context"
	"strings"
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
}

func (s *Service) GetDashboard(ctx context.Context, req GetDashboardRequest) (DashboardViewModel, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return DashboardViewModel{}, ErrInvalidInput
	}

	if req.RecentLimit <= 0 {
		req.RecentLimit = 8
	}

	target, err := s.GetTarget(ctx, GetTargetRequest{UserID: req.UserID})
	if err != nil {
		return DashboardViewModel{}, err
	}

	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	from := startOfDay.Format(time.RFC3339)
	to := endOfDay.Format(time.RFC3339)

	summaryData, _ := s.store.GetSummaryByRange(ctx, req.UserID, from, to)
	logs, err := s.store.GetLogsByRange(ctx, req.UserID, from, to)
	if err != nil {
		logs = []Log{}
	}
	recentMeals, err := s.store.GetRecentLogs(ctx, req.UserID, req.RecentLimit)
	if err != nil {
		recentMeals = []Log{}
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
