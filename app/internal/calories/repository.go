package calories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/23St/trackstack/internal/common/db"
)

// Repository handles database operations for the calories domain
type Repository interface {
	GetDailyTotals(ctx context.Context, userID string, start, end time.Time) (calories, protein int, err error)
	GetUserTargets(ctx context.Context, userID string) (*UserTargets, error)
	CreateMeal(ctx context.Context, meal *Meal) error
}

// SQLRepository implements Repository using SQL database
type SQLRepository struct {
	db *db.DB
}

// NewRepository creates a new SQL repository
func NewRepository(database *db.DB) Repository {
	return &SQLRepository{db: database}
}

// GetDailyTotals returns the sum of calories and protein for a given date range
func (r *SQLRepository) GetDailyTotals(ctx context.Context, userID string, start, end time.Time) (calories, protein int, err error) {
	query := `
		SELECT 
			COALESCE(SUM(calories), 0) as total_calories,
			COALESCE(SUM(protein), 0) as total_protein
		FROM meals
		WHERE user_id = ?
		  AND logged_at >= ?
		  AND logged_at < ?
	`

	startUnix := start.Unix()
	endUnix := end.Unix()

	err = r.db.QueryRowContext(ctx, query, userID, startUnix, endUnix).Scan(&calories, &protein)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get daily totals: %w", err)
	}

	return calories, protein, nil
}

// GetUserTargets retrieves user's daily targets, returning defaults if not found
func (r *SQLRepository) GetUserTargets(ctx context.Context, userID string) (*UserTargets, error) {
	query := `
		SELECT user_id, calorie_target, protein_target, updated_at
		FROM user_targets
		WHERE user_id = ?
	`

	var targets UserTargets
	var updatedAtUnix int64

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&targets.UserID,
		&targets.CalorieTarget,
		&targets.ProteinTarget,
		&updatedAtUnix,
	)

	if errors.Is(err, sql.ErrNoRows) {
		// Return defaults if no record exists
		return &UserTargets{
			UserID:        userID,
			CalorieTarget: 2300,
			ProteinTarget: 120,
			UpdatedAt:     time.Now(),
		}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user targets: %w", err)
	}

	targets.UpdatedAt = time.Unix(updatedAtUnix, 0)
	return &targets, nil
}

// CreateMeal inserts a new meal record into the database
func (r *SQLRepository) CreateMeal(ctx context.Context, meal *Meal) error {
	query := `
		INSERT INTO meals (id, user_id, name, calories, protein, carbs, fat, logged_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		meal.ID,
		meal.UserID,
		meal.Name,
		meal.Calories,
		meal.Protein,
		meal.Carbs,
		meal.Fat,
		meal.LoggedAt.Unix(),
	)
	if err != nil {
		return fmt.Errorf("failed to create meal: %w", err)
	}

	return nil
}
