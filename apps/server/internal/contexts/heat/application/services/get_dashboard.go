package services

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/application/ports"
	"github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/domain"
)

type GetDashboardService struct {
	recentRefills ports.RecentRefillLister
	latestRefill  ports.LatestRefillGetter
	seasonTotals  ports.RefillRangeSummarizer
}

func NewGetDashboardService(recentRefills ports.RecentRefillLister, latestRefill ports.LatestRefillGetter, seasonTotals ports.RefillRangeSummarizer) GetDashboardService {
	return GetDashboardService{
		recentRefills: recentRefills,
		latestRefill:  latestRefill,
		seasonTotals:  seasonTotals,
	}
}

func (s GetDashboardService) Execute(ctx context.Context, req GetDashboardRequest) (DashboardViewModel, error) {
	start := time.Now()
	defer logHeatTiming(ctx, "dashboard.total", start, nil)

	if strings.TrimSpace(req.UserID) == "" {
		return DashboardViewModel{}, domain.ErrInvalidInput
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	offset := (req.Page - 1) * req.Limit

	daysSinceRefill := 0
	latestStart := time.Now()
	latest, err := s.latestRefill.GetLatest(ctx, req.UserID)
	logHeatTiming(ctx, "dashboard.latest", latestStart, err)
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

	now := time.Now().UTC()
	todayUTCStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	todayUTCEndExclusive := todayUTCStart.Add(24 * time.Hour)

	seasonStartYear := todayUTCStart.Year()
	if todayUTCStart.Month() < time.September {
		seasonStartYear--
	}

	currentSeasonStart := time.Date(seasonStartYear, time.September, 1, 0, 0, 0, 0, time.UTC)
	currentSeasonEnd := time.Date(seasonStartYear+1, time.September, 1, 0, 0, 0, 0, time.UTC)
	currentSeasonLabel := domain.SeasonLabelFor(currentSeasonStart)

	currentEnd := currentSeasonEnd
	if todayUTCEndExclusive.Before(currentSeasonEnd) {
		currentEnd = todayUTCEndExclusive
	}

	lastSeasonStart := time.Date(seasonStartYear-1, time.September, 1, 0, 0, 0, 0, time.UTC)
	lastSeasonEnd := time.Date(seasonStartYear, time.September, 1, 0, 0, 0, 0, time.UTC)

	offsetDuration := currentEnd.Sub(currentSeasonStart)
	lastSeasonEndSamePeriod := lastSeasonStart.Add(offsetDuration)
	lastEnd := lastSeasonEnd
	if lastSeasonEndSamePeriod.Before(lastSeasonEnd) {
		lastEnd = lastSeasonEndSamePeriod
	}

	currentSeasonStartTime := time.Now()
	seasonToDate, err := s.seasonTotals.GetSumByRange(ctx, req.UserID, currentSeasonStart.Format(time.RFC3339), currentEnd.Format(time.RFC3339))
	logHeatTiming(ctx, "dashboard.current_season_sum", currentSeasonStartTime, err)
	if err != nil {
		seasonToDate = 0
	}
	lastSeasonStartTime := time.Now()
	lastSeasonToDate, err := s.seasonTotals.GetSumByRange(ctx, req.UserID, lastSeasonStart.Format(time.RFC3339), lastEnd.Format(time.RFC3339))
	logHeatTiming(ctx, "dashboard.last_season_sum", lastSeasonStartTime, err)
	if err != nil {
		lastSeasonToDate = 0
	}

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

	historyStart := time.Now()
	history, err := s.recentRefills.ListRecent(ctx, req.UserID, req.Limit, offset)
	logHeatTiming(ctx, "dashboard.history", historyStart, err)
	if err != nil {
		history = []domain.Refill{}
	}

	return DashboardViewModel{
		DaysSinceRefill: daysSinceRefill,
		SeasonSnapshot:  snapshot,
		History:         history,
	}, nil
}

func logHeatTiming(ctx context.Context, step string, start time.Time, err error) {
	attrs := []any{"step", step, "duration", time.Since(start)}
	if err != nil {
		attrs = append(attrs, "error", err)
	}
	if logger := slog.Default(); logger != nil {
		logger.DebugContext(ctx, "heat timing", attrs...)
	}
}
