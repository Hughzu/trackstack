package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/calories"
)

func (s *CaloriesStore) GetTarget(ctx context.Context, userID string) (calories.Target, error) {
	row := s.db.QueryRowContext(
		ctx,
		"SELECT id, user_id, target_kcal, target_protein_g, target_carbs_g, target_fat_g, created_at, updated_at FROM calorie_targets WHERE user_id = ?",
		userID,
	)

	var target calories.Target
	var carbs sql.NullInt64
	var fat sql.NullInt64

	err := row.Scan(
		&target.ID,
		&target.UserID,
		&target.TargetKcal,
		&target.TargetProteinG,
		&carbs,
		&fat,
		&target.CreatedAt,
		&target.UpdatedAt,
	)
	if err != nil {
		return calories.Target{}, err
	}

	if carbs.Valid {
		value := int(carbs.Int64)
		target.TargetCarbsG = &value
	}
	if fat.Valid {
		value := int(fat.Int64)
		target.TargetFatG = &value
	}

	return target, nil
}

func (s *CaloriesStore) CreateTarget(ctx context.Context, target calories.Target) error {
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO calorie_targets (id, user_id, target_kcal, target_protein_g, target_carbs_g, target_fat_g, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		target.ID,
		target.UserID,
		target.TargetKcal,
		target.TargetProteinG,
		nullInt(target.TargetCarbsG),
		nullInt(target.TargetFatG),
		target.CreatedAt,
		target.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create target: %w", err)
	}

	return nil
}

func (s *CaloriesStore) UpdateTarget(ctx context.Context, target calories.Target) error {
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE calorie_targets SET target_kcal = ?, target_protein_g = ?, target_carbs_g = ?, target_fat_g = ?, updated_at = ? WHERE id = ?",
		target.TargetKcal,
		target.TargetProteinG,
		nullInt(target.TargetCarbsG),
		nullInt(target.TargetFatG),
		target.UpdatedAt,
		target.ID,
	)
	if err != nil {
		return fmt.Errorf("update target: %w", err)
	}

	return nil
}

func nullInt(value *int) any {
	if value == nil {
		return nil
	}
	return *value
}
