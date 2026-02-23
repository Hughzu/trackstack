package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/heat"
	"github.com/google/uuid"
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

func (s *RefillStore) Create(ctx context.Context, userID string, input heat.CreateRefillInput) (heat.Refill, error) {
	refillID := uuid.NewString()

	_, err := s.db.ExecContext(
		ctx,
		`INSERT INTO refills (id, user_id, date, weight_kg, bags, temperature, season)
VALUES (?, ?, ?, ?, ?, ?, ?)`,
		refillID,
		userID,
		input.Date,
		input.WeightKg,
		input.Bags,
		input.Temperature,
		input.Season,
	)
	if err != nil {
		return heat.Refill{}, fmt.Errorf("create refill: %w", err)
	}

	return heat.Refill{
		ID:          refillID,
		UserID:      userID,
		Date:        input.Date,
		WeightKg:    input.WeightKg,
		Bags:        input.Bags,
		Temperature: input.Temperature,
		Season:      input.Season,
	}, nil
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
