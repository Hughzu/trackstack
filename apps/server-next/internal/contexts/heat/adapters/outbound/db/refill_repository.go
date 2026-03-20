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
		SELECT id, amount, created_at 
		FROM refills 
		WHERE user_id = ? 
		AND created_at >= ? 
		AND created_at <= ?
		ORDER BY created_at DESC
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
		var amount int
		var createdAtStr string

		if err := rows.Scan(&id, &amount, &createdAtStr); err != nil {
			return nil, fmt.Errorf("failed to scan refill, %w", err)
		}

		parsedTime, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse refill time, %w", err)
		}

		refills = append(refills, domain.Refill{
			ID:        id,
			Amount:    amount,
			CreatedAt: parsedTime,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return refills, nil

}
