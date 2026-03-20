package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/heat/domain"
)

type RefillRepository struct {
	db *sql.DB
}

func NewRefillRepository(db *sql.DB) *RefillRepository {
	return &RefillRepository{db: db}
}

func (r *RefillRepository) GetRefills(ctx context.Context, userID string, from, to time.Time) ([]domain.Refill, error) {
	query := `
		SELECT id, weight_kg, bags, temperature, season, date 
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
		var id string
		var weightKg float64
		var bags int
		var temperature *float64
		var season *string
		var dateStr string

		if err := rows.Scan(&id, &weightKg, &bags, &temperature, &season, &dateStr); err != nil {
			return nil, fmt.Errorf("failed to scan refill, %w", err)
		}

		parsedTime, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			parsedTime, _ = time.Parse("2006-01-02", dateStr)
		}

		refills = append(refills, domain.Refill{
			ID:          id,
			WeightKg:    weightKg,
			Bags:        bags,
			Temperature: temperature,
			Season:      season,
			Date:        parsedTime,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return refills, nil

}
