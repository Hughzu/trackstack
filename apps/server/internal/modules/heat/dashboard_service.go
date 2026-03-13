package heat

import (
	"context"
	"strings"
	"time"
)

type SeasonSnapshot struct {
	SeasonLabel      string `json:"seasonLabel"`
	SeasonToDate     int    `json:"seasonToDate"`
	LastSeasonToDate int    `json:"lastSeasonToDate"`
	Delta            int    `json:"delta"`
	DeltaPct         *int   `json:"deltaPct"`
}

type DashboardViewModel struct {
	DaysSinceRefill int            `json:"daysSinceRefill"`
	SeasonSnapshot  SeasonSnapshot `json:"seasonSnapshot"`
	History         []Refill       `json:"history"`
}

type GetDashboardRequest struct {
	UserID string
	Page   int
	Limit  int
}

func (s *Service) GetDashboard(ctx context.Context, req GetDashboardRequest) (DashboardViewModel, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return DashboardViewModel{}, ErrInvalidInput
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	offset := (req.Page - 1) * req.Limit

	daysSinceRefill := 0
	latest, err := s.store.GetLatest(ctx, req.UserID)
	if err != nil {
		return DashboardViewModel{}, err
	}
	if latest != nil {
		latestDate, _ := time.Parse(time.RFC3339, latest.Date)
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		lastDate := time.Date(latestDate.Year(), latestDate.Month(), latestDate.Day(), 0, 0, 0, 0, time.UTC)
		daysSinceRefill = int(today.Sub(lastDate).Hours() / 24)
	}

	// Calculate season snapshot
	now := time.Now().UTC()
	todayUtcStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	todayUtcEndExclusive := todayUtcStart.Add(24 * time.Hour)

	seasonStartYear := todayUtcStart.Year()
	if todayUtcStart.Month() < time.September {
		seasonStartYear--
	}

	currentSeasonStart := time.Date(seasonStartYear, time.September, 1, 0, 0, 0, 0, time.UTC)
	currentSeasonEnd := time.Date(seasonStartYear+1, time.September, 1, 0, 0, 0, 0, time.UTC)
	currentSeasonLabel := seasonLabelFor(currentSeasonStart)

	currentEnd := currentSeasonEnd
	if todayUtcEndExclusive.Before(currentSeasonEnd) {
		currentEnd = todayUtcEndExclusive
	}

	lastSeasonStart := time.Date(seasonStartYear-1, time.September, 1, 0, 0, 0, 0, time.UTC)
	lastSeasonEnd := time.Date(seasonStartYear, time.September, 1, 0, 0, 0, 0, time.UTC)

	offsetDuration := currentEnd.Sub(currentSeasonStart)
	lastSeasonEndSamePeriod := lastSeasonStart.Add(offsetDuration)
	lastEnd := lastSeasonEnd
	if lastSeasonEndSamePeriod.Before(lastSeasonEnd) {
		lastEnd = lastSeasonEndSamePeriod
	}

	seasonToDate, _ := s.store.GetSumByRange(ctx, req.UserID, currentSeasonStart.Format(time.RFC3339), currentEnd.Format(time.RFC3339))
	lastSeasonToDate, _ := s.store.GetSumByRange(ctx, req.UserID, lastSeasonStart.Format(time.RFC3339), lastEnd.Format(time.RFC3339))

	delta := seasonToDate - lastSeasonToDate
	var deltaPct *int
	if lastSeasonToDate > 0 {
		pct := int((float64(delta) / float64(lastSeasonToDate)) * 100)
		deltaPct = &pct
	}

	snapshot := SeasonSnapshot{
		SeasonLabel:      currentSeasonLabel,
		SeasonToDate:     seasonToDate,
		LastSeasonToDate: lastSeasonToDate,
		Delta:            delta,
		DeltaPct:         deltaPct,
	}

	history, err := s.store.ListRecent(ctx, req.UserID, req.Limit, offset)
	if err != nil {
		history = []Refill{}
	}

	return DashboardViewModel{
		DaysSinceRefill: daysSinceRefill,
		SeasonSnapshot:  snapshot,
		History:         history,
	}, nil
}
