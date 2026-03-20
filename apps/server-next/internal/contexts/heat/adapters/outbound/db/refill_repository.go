package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/heat/domain"
	"github.com/Hughzu/trackstack/apps/server-next/internal/platform/timeutil"
)

type RefillRepository struct {
	db *sql.DB
}

func NewRefillRepository(db *sql.DB) *RefillRepository {
	return &RefillRepository{db: db}
}

func (r *RefillRepository) GetRefills(ctx context.Context, userID string, from, to time.Time) ([]domain.Refill, error) {
	query := `
		SELECT id, user_id, weight_kg, bags, temperature, season, date 
		FROM refills 
		WHERE user_id = ? 
		AND date >= ? 
		AND date <= ?
		ORDER BY date DESC
	`

	fromStr := from.UTC().Format(time.RFC3339)
	toStr := to.UTC().Format(time.RFC3339)

	rows, err := r.db.QueryContext(ctx, query, userID, fromStr, toStr)
	if err != nil {
		return nil, fmt.Errorf("failed to query refills, %w", err)
	}

	defer rows.Close()

	var refills []domain.Refill
	for rows.Next() {
		var refill domain.Refill
		var temperature sql.NullFloat64
		var season sql.NullString

		if err := rows.Scan(&refill.ID, &refill.UserID, &refill.WeightKg, &refill.Bags, &temperature, &season, &refill.Date); err != nil {
			return nil, fmt.Errorf("failed to scan refill, %w", err)
		}
		if temperature.Valid {
			refill.Temperature = &temperature.Float64
		}
		if season.Valid {
			refill.Season = &season.String
		}
		refill.Date = timeutil.NormalizeDateString(refill.Date)

		refills = append(refills, refill)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return refills, nil

}

func (r *RefillRepository) CreateRefill(ctx context.Context, refill domain.Refill) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO refills (id, user_id, date, weight_kg, bags, temperature, season)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		refill.ID,
		refill.UserID,
		refill.Date,
		refill.WeightKg,
		refill.Bags,
		refill.Temperature,
		refill.Season,
	)
	if err != nil {
		return fmt.Errorf("failed to create refill: %w", err)
	}

	return nil
}

func (r *RefillRepository) DeleteRefill(ctx context.Context, userID string, id string) (bool, error) {
	result, err := r.db.ExecContext(
		ctx,
		"DELETE FROM refills WHERE id = ? AND user_id = ?",
		id,
		userID,
	)
	if err != nil {
		return false, fmt.Errorf("failed to delete refill: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to read deleted refill rows: %w", err)
	}

	return rows > 0, nil
}
