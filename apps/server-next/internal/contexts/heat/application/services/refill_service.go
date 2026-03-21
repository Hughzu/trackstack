package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/heat/application/ports"
	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/heat/domain"
	"github.com/google/uuid"
)

type refillService struct {
	repo ports.RefillRepository
}

func NewRefillService(repo ports.RefillRepository) ports.RefillUseCase {
	return &refillService{repo: repo}
}

func (s *refillService) GetRefills(ctx context.Context, query ports.GetRefillsQuery) ([]domain.Refill, error) {
	if strings.TrimSpace(query.UserID) == "" {
		return nil, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}

	return s.repo.GetRefills(ctx, query.UserID, query.From, query.To)
}

func (s *refillService) CreateRefill(ctx context.Context, command ports.CreateRefillCommand) (domain.Refill, error) {
	if strings.TrimSpace(command.UserID) == "" {
		return domain.Refill{}, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}
	if command.WeightKg <= 0 {
		return domain.Refill{}, fmt.Errorf("%w: weight must be greater than zero", domain.ErrInvalidInput)
	}
	if command.Bags <= 0 {
		return domain.Refill{}, fmt.Errorf("%w: bags must be greater than zero", domain.ErrInvalidInput)
	}
	if command.Date.IsZero() {
		return domain.Refill{}, fmt.Errorf("%w: date is required", domain.ErrInvalidInput)
	}

	refillDate := command.Date.UTC()
	seasonLabel := domain.SeasonLabelFor(refillDate)
	refill := domain.Refill{
		ID:          uuid.NewString(),
		UserID:      command.UserID,
		Date:        refillDate.UTC().Format(time.RFC3339),
		WeightKg:    command.WeightKg,
		Bags:        command.Bags,
		Temperature: command.Temperature,
		Season:      &seasonLabel,
	}

	if err := s.repo.CreateRefill(ctx, refill); err != nil {
		return domain.Refill{}, err
	}

	return refill, nil
}

func (s *refillService) DeleteRefill(ctx context.Context, command ports.DeleteRefillCommand) (bool, error) {
	if strings.TrimSpace(command.UserID) == "" {
		return false, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}
	if strings.TrimSpace(command.ID) == "" {
		return false, fmt.Errorf("%w: refill id is required", domain.ErrInvalidInput)
	}

	return s.repo.DeleteRefill(ctx, command.UserID, command.ID)
}

func (s *refillService) GetDashboard(ctx context.Context, query ports.GetDashboardQuery) (ports.DashboardViewModel, error) {
	if strings.TrimSpace(query.UserID) == "" {
		return ports.DashboardViewModel{}, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 20
	}
	offset := (query.Page - 1) * query.Limit

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

	latestDate, seasonToDate, lastSeasonToDate, err := s.repo.GetDashboardStats(ctx, query.UserID, currentSeasonStart, currentEnd, lastSeasonStart, lastEnd)
	if err != nil {
		return ports.DashboardViewModel{}, err
	}

	daysSinceRefill := 0
	if latestDate != nil {
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		lastDate := time.Date(latestDate.Year(), latestDate.Month(), latestDate.Day(), 0, 0, 0, 0, time.UTC)
		daysSinceRefill = int(today.Sub(lastDate).Hours() / 24)
	}

	delta := seasonToDate - lastSeasonToDate
	var deltaPct *int
	if lastSeasonToDate > 0 {
		pct := int((float64(delta) / float64(lastSeasonToDate)) * 100)
		deltaPct = &pct
	}

	snapshot := ports.SeasonSnapshot{
		SeasonLabel:      currentSeasonLabel,
		SeasonToDate:     seasonToDate,
		LastSeasonToDate: lastSeasonToDate,
		Delta:            delta,
		DeltaPct:         deltaPct,
	}

	history, err := s.repo.ListRecentRefills(ctx, query.UserID, query.Limit, offset)
	if err != nil {
		history = []domain.Refill{}
	}

	return ports.DashboardViewModel{
		DaysSinceRefill: daysSinceRefill,
		SeasonSnapshot:  snapshot,
		History:         history,
	}, nil
}

