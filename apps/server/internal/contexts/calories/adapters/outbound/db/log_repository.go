package db

import (
	"context"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/calories/domain"
)

func (r *Repository) CreateLog(ctx context.Context, log domain.Log) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO calorie_logs (id, user_id, date_time, calories, protein_g, carbs_g, fat_g, title)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		log.ID,
		log.UserID,
		log.DateTime,
		log.Calories,
		log.ProteinGrams,
		nullableInt(log.CarbGrams),
		nullableInt(log.FatGrams),
		log.Title,
	)
	if err != nil {
		return fmt.Errorf("create log: %w", err)
	}

	return nil
}

func (r *Repository) DeleteLog(ctx context.Context, userID string, id string) (bool, error) {
	result, err := r.db.ExecContext(
		ctx,
		"DELETE FROM calorie_logs WHERE id = ? AND user_id = ?",
		id,
		userID,
	)
	if err != nil {
		return false, fmt.Errorf("delete log: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("delete log rows: %w", err)
	}

	return rowsAffected > 0, nil
}
