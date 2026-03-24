package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/calories/domain"
	"github.com/Hughzu/trackstack/apps/server/internal/platform/timeutil"
)

func (r *Repository) GetNutritionSummaryByRange(ctx context.Context, userID string, from string, to string) (domain.NutritionSummary, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT COALESCE(SUM(calories), 0), COALESCE(SUM(protein_g), 0), COALESCE(SUM(carbs_g), 0), COALESCE(SUM(fat_g), 0)
		FROM calorie_logs
		WHERE user_id = ? AND date_time >= ? AND date_time < ?`,
		userID,
		from,
		to,
	)

	var summary domain.NutritionSummary
	err := row.Scan(&summary.Calories, &summary.ProteinGrams, &summary.CarbGrams, &summary.FatGrams)
	if err != nil && err != sql.ErrNoRows {
		return domain.NutritionSummary{}, fmt.Errorf("scan nutrition summary: %w", err)
	}

	return summary, nil
}

func (r *Repository) GetLogsByRange(ctx context.Context, userID string, from string, to string, limit int) ([]domain.Log, error) {
	if limit <= 0 {
		limit = 50
	}

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, date_time, calories, protein_g, carbs_g, fat_g, title
		FROM calorie_logs
		WHERE user_id = ? AND date_time >= ? AND date_time < ?
		ORDER BY date_time DESC
		LIMIT ?`,
		userID,
		from,
		to,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query logs by range: %w", err)
	}
	defer rows.Close()

	logs, err := scanLogs(rows)
	if err != nil {
		return nil, err
	}

	if logs == nil {
		return []domain.Log{}, nil
	}

	return logs, nil
}

func (r *Repository) GetRecentMeals(ctx context.Context, userID string, limit int) ([]domain.Log, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`WITH ranked AS (
			SELECT id, user_id, date_time, calories, protein_g, carbs_g, fat_g, title,
			ROW_NUMBER() OVER (PARTITION BY LOWER(TRIM(title)) ORDER BY date_time DESC) AS rn
			FROM calorie_logs
			WHERE user_id = ? AND title IS NOT NULL AND TRIM(title) != ''
		)
		SELECT id, user_id, date_time, calories, protein_g, carbs_g, fat_g, title
		FROM ranked
		WHERE rn = 1
		ORDER BY date_time DESC
		LIMIT ?`,
		userID,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query recent meals: %w", err)
	}
	defer rows.Close()

	logs, err := scanLogs(rows)
	if err != nil {
		return nil, err
	}

	if logs == nil {
		return []domain.Log{}, nil
	}

	return logs, nil
}

func scanLogs(rows *sql.Rows) ([]domain.Log, error) {
	var logs []domain.Log
	for rows.Next() {
		log, err := scanLog(rows)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return logs, nil
}

func scanLog(scanner interface{ Scan(dest ...any) error }) (domain.Log, error) {
	var log domain.Log
	var carbGrams sql.NullInt64
	var fatGrams sql.NullInt64
	var title sql.NullString

	err := scanner.Scan(
		&log.ID,
		&log.UserID,
		&log.DateTime,
		&log.Calories,
		&log.ProteinGrams,
		&carbGrams,
		&fatGrams,
		&title,
	)
	if err != nil {
		return domain.Log{}, fmt.Errorf("scan log: %w", err)
	}

	if carbGrams.Valid {
		value := int(carbGrams.Int64)
		log.CarbGrams = &value
	}
	if fatGrams.Valid {
		value := int(fatGrams.Int64)
		log.FatGrams = &value
	}
	if title.Valid {
		value := title.String
		log.Title = &value
	}

	log.DateTime = timeutil.NormalizeDateString(log.DateTime)

	return log, nil
}
