package calories

import (
	"context"
	"fmt"
	"time"
)

// Service defines the interface for calorie business logic
type Service interface {
	CalculateDailySummary(ctx context.Context, userID string, date time.Time) (*DailySummary, error)
}

// service implements Service interface
type service struct {
	repo Repository
}

// NewService creates a new Service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CalculateDailySummary computes the dashboard view for a given user and date
func (s *service) CalculateDailySummary(ctx context.Context, userID string, date time.Time) (*DailySummary, error) {
	// Get date range for "today" (midnight to midnight)
	start, end := getDayRange(date)

	// Fetch consumed totals from database
	totalCalories, totalProtein, err := s.repo.GetDailyTotals(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily totals: %w", err)
	}

	// Fetch user targets
	targets, err := s.repo.GetUserTargets(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user targets: %w", err)
	}

	// Calculate metric statuses
	return &DailySummary{
		Calories: CalculateMetricStatus(totalCalories, targets.CalorieTarget),
		Protein:  CalculateMetricStatus(totalProtein, targets.ProteinTarget),
	}, nil
}

// CalculateMetricStatus determines the status for a single metric
func CalculateMetricStatus(consumed, target int) MetricStatus {
	remaining := target - consumed
	percentage := 0
	if target > 0 {
		percentage = (consumed * 100) / target
	}

	// Determine status based on remaining percentage
	remainingPct := 100 - percentage
	var status string

	if remainingPct > 20 {
		status = "safe" // Green: >20% budget remaining
	} else if remainingPct >= 5 {
		status = "warn" // Yellow: 5-20% remaining
	} else {
		status = "danger" // Red: <5% or over budget
	}

	return MetricStatus{
		Consumed:   consumed,
		Target:     target,
		Remaining:  remaining,
		Percentage: percentage,
		Status:     status,
	}
}

// getDayRange returns the start and end time for a given date (midnight to midnight)
func getDayRange(date time.Time) (start, end time.Time) {
	start = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end = start.Add(24 * time.Hour)
	return
}
