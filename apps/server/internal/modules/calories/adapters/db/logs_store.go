package db

import (
	"context"
	"database/sql"
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

func (s *CaloriesStore) GetSummaryByRange(ctx context.Context, userID string, from string, to string) (calories.Summary, error) {
	row := s.db.QueryRowContext(
		ctx,
		"SELECT COALESCE(SUM(calories), 0), COALESCE(SUM(protein_g), 0), COALESCE(SUM(carbs_g), 0), COALESCE(SUM(fat_g), 0) FROM calorie_logs WHERE user_id = ? AND date_time >= ? AND date_time < ?",
		userID, from, to,
	)

	var sum calories.Summary
	err := row.Scan(&sum.Calories, &sum.ProteinG, &sum.CarbsG, &sum.FatG)
	if err != nil && err != sql.ErrNoRows {
		return sum, fmt.Errorf("scan summary: %w", err)
	}

	return sum, nil
}

func (s *CaloriesStore) GetLogsByRange(ctx context.Context, userID string, from string, to string) ([]calories.Log, error) {
	rows, err := s.db.QueryContext(
		ctx,
		"SELECT id, user_id, date_time, calories, protein_g, carbs_g, fat_g, title FROM calorie_logs WHERE user_id = ? AND date_time >= ? AND date_time < ? ORDER BY date_time DESC",
		userID, from, to,
	)
	if err != nil {
		return nil, fmt.Errorf("query logs: %w", err)
	}
	defer rows.Close()

	var logs []calories.Log
	for rows.Next() {
		var log calories.Log
		var carbs sql.NullInt64
		var fat sql.NullInt64
		var title sql.NullString
		if err := rows.Scan(&log.ID, &log.UserID, &log.DateTime, &log.Calories, &log.ProteinG, &carbs, &fat, &title); err != nil {
			return nil, fmt.Errorf("scan log: %w", err)
		}
		if carbs.Valid {
			val := int(carbs.Int64)
			log.CarbsG = &val
		}
		if fat.Valid {
			val := int(fat.Int64)
			log.FatG = &val
		}
		if title.Valid {
			val := title.String
			log.Title = &val
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (s *CaloriesStore) GetLogsByRangeLimited(ctx context.Context, userID string, from string, to string, limit int) ([]calories.Log, error) {
	if limit <= 0 {
		limit = 50
	}

	rows, err := s.db.QueryContext(
		ctx,
		"SELECT id, user_id, date_time, calories, protein_g, carbs_g, fat_g, title FROM calorie_logs WHERE user_id = ? AND date_time >= ? AND date_time < ? ORDER BY date_time DESC LIMIT ?",
		userID, from, to, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query limited logs: %w", err)
	}
	defer rows.Close()

	var logs []calories.Log
	for rows.Next() {
		var log calories.Log
		var carbs sql.NullInt64
		var fat sql.NullInt64
		var title sql.NullString
		if err := rows.Scan(&log.ID, &log.UserID, &log.DateTime, &log.Calories, &log.ProteinG, &carbs, &fat, &title); err != nil {
			return nil, fmt.Errorf("scan limited log: %w", err)
		}
		if carbs.Valid {
			val := int(carbs.Int64)
			log.CarbsG = &val
		}
		if fat.Valid {
			val := int(fat.Int64)
			log.FatG = &val
		}
		if title.Valid {
			val := title.String
			log.Title = &val
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (s *CaloriesStore) GetRecentLogs(ctx context.Context, userID string, limit int) ([]calories.Log, error) {
	rows, err := s.db.QueryContext(
		ctx,
		"WITH ranked AS (SELECT id, user_id, date_time, calories, protein_g, carbs_g, fat_g, title, ROW_NUMBER() OVER (PARTITION BY LOWER(TRIM(title)) ORDER BY date_time DESC) AS rn FROM calorie_logs WHERE user_id = ? AND title IS NOT NULL AND TRIM(title) != '') SELECT id, user_id, date_time, calories, protein_g, carbs_g, fat_g, title FROM ranked WHERE rn = 1 ORDER BY date_time DESC LIMIT ?",
		userID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query recent logs: %w", err)
	}
	defer rows.Close()

	var logs []calories.Log
	for rows.Next() {
		var log calories.Log
		var carbs sql.NullInt64
		var fat sql.NullInt64
		var title sql.NullString
		if err := rows.Scan(&log.ID, &log.UserID, &log.DateTime, &log.Calories, &log.ProteinG, &carbs, &fat, &title); err != nil {
			return nil, fmt.Errorf("scan recent log: %w", err)
		}
		if carbs.Valid {
			val := int(carbs.Int64)
			log.CarbsG = &val
		}
		if fat.Valid {
			val := int(fat.Int64)
			log.FatG = &val
		}
		if title.Valid {
			val := title.String
			log.Title = &val
		}
		logs = append(logs, log)
	}

	return logs, nil
}
