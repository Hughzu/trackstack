package calories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Service defines the interface for calorie business logic
type Service interface {
	CalculateDailySummary(ctx context.Context, userID string, date time.Time) (*DailySummary, error)
	LogMeal(ctx context.Context, userID string, name string, calories, protein, carbs, fat int) (*DailySummary, error)
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

	// Fetch recent meals
	meals, err := s.repo.GetDailyMeals(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily meals: %w", err)
	}

	// Calculate metric statuses
	return &DailySummary{
		Calories:    CalculateMetricStatus(totalCalories, targets.CalorieTarget),
		Protein:     CalculateMetricStatus(totalProtein, targets.ProteinTarget),
		RecentMeals: meals,
	}, nil
}

// LogMeal creates a meal entry and returns the updated daily summary
func (s *service) LogMeal(ctx context.Context, userID string, name string, calories, protein, carbs, fat int) (*DailySummary, error) {
	// Validate inputs
	if calories <= 0 {
		return nil, fmt.Errorf("calories must be greater than 0")
	}
	if protein < 0 {
		return nil, fmt.Errorf("protein cannot be negative")
	}
	if carbs < 0 {
		return nil, fmt.Errorf("carbs cannot be negative")
	}
	if fat < 0 {
		return nil, fmt.Errorf("fat cannot be negative")
	}

	// Create meal
	meal := &Meal{
		ID:       uuid.New().String(),
		UserID:   userID,
		Name:     name,
		Calories: calories,
		Protein:  protein,
		Carbs:    carbs,
		Fat:      fat,
		LoggedAt: time.Now(),
	}

	if err := s.repo.CreateMeal(ctx, meal); err != nil {
		return nil, fmt.Errorf("failed to log meal: %w", err)
	}

	// Return updated summary
	return s.CalculateDailySummary(ctx, userID, time.Now())
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
