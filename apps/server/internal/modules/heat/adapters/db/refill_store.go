package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/heat"
)

type RefillStore struct {
	db *sql.DB
}

func NewRefillStore(db *sql.DB) *RefillStore {
	return &RefillStore{db: db}
}

func (s *RefillStore) ListByRange(ctx context.Context, userID string, from string, to string) ([]heat.Refill, error) {
	query := `
SELECT id, user_id, date, weight_kg, bags, temperature, season
FROM refills
WHERE user_id = ?`

	args := []any{userID}
	if from != "" {
		query += " AND date >= ?"
		args = append(args, from)
	}
	if to != "" {
		query += " AND date < ?"
		args = append(args, to)
	}

	query += " ORDER BY date DESC"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list refills: %w", err)
	}
	defer rows.Close()

	var refills []heat.Refill
	for rows.Next() {
		var refill heat.Refill
		var temperature sql.NullFloat64
		var season sql.NullString
		if err := rows.Scan(
			&refill.ID,
			&refill.UserID,
			&refill.Date,
			&refill.WeightKg,
			&refill.Bags,
			&temperature,
			&season,
		); err != nil {
			return nil, fmt.Errorf("scan refill: %w", err)
		}

		if temperature.Valid {
			value := temperature.Float64
			refill.Temperature = &value
		}
		if season.Valid {
			value := season.String
			refill.Season = &value
		}

		refills = append(refills, refill)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list refills rows: %w", err)
	}

	return refills, nil
}

func (s *RefillStore) Create(ctx context.Context, refill heat.Refill) error {
	_, err := s.db.ExecContext(
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
		return fmt.Errorf("create refill: %w", err)
	}

	return nil
}

func (s *RefillStore) Delete(ctx context.Context, userID string, id string) (bool, error) {
	result, err := s.db.ExecContext(
		ctx,
		"DELETE FROM refills WHERE id = ? AND user_id = ?",
		id,
		userID,
	)
	if err != nil {
		return false, fmt.Errorf("delete refill: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("delete refill rows: %w", err)
	}

	return rows > 0, nil
}
