package db

import (
	"context"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/calories"
)

func (s *CaloriesStore) CreateLog(ctx context.Context, log calories.Log) error {
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO calorie_logs (id, user_id, date_time, calories, protein_g, carbs_g, fat_g, title) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		log.ID,
		log.UserID,
		log.DateTime,
		log.Calories,
		log.ProteinG,
		nullInt(log.CarbsG),
		nullInt(log.FatG),
		log.Title,
	)
	if err != nil {
		return fmt.Errorf("create log: %w", err)
	}

	return nil
}

func (s *CaloriesStore) DeleteLog(ctx context.Context, userID string, id string) (bool, error) {
	result, err := s.db.ExecContext(
		ctx,
		"DELETE FROM calorie_logs WHERE id = ? AND user_id = ?",
		id,
		userID,
	)
	if err != nil {
		return false, fmt.Errorf("delete log: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("delete log rows: %w", err)
	}

	return rows > 0, nil
}
