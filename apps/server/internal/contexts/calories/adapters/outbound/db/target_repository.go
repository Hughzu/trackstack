package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/calories/domain"
)

func (r *Repository) FindTarget(ctx context.Context, userID string) (domain.Target, bool, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, target_kcal, target_protein_g, target_carbs_g, target_fat_g, created_at, updated_at
		FROM calorie_targets
		WHERE user_id = ?`,
		userID,
	)

	var target domain.Target
	var carbGrams sql.NullInt64
	var fatGrams sql.NullInt64

	err := row.Scan(
		&target.ID,
		&target.UserID,
		&target.TargetCalories,
		&target.TargetProteinGrams,
		&carbGrams,
		&fatGrams,
		&target.CreatedAt,
		&target.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Target{}, false, nil
		}
		return domain.Target{}, false, fmt.Errorf("find target: %w", err)
	}

	if carbGrams.Valid {
		value := int(carbGrams.Int64)
		target.TargetCarbGrams = &value
	}
	if fatGrams.Valid {
		value := int(fatGrams.Int64)
		target.TargetFatGrams = &value
	}

	return target, true, nil
}

func (r *Repository) CreateTarget(ctx context.Context, target domain.Target) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO calorie_targets (id, user_id, target_kcal, target_protein_g, target_carbs_g, target_fat_g, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		target.ID,
		target.UserID,
		target.TargetCalories,
		target.TargetProteinGrams,
		nullableInt(target.TargetCarbGrams),
		nullableInt(target.TargetFatGrams),
		target.CreatedAt,
		target.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create target: %w", err)
	}

	return nil
}

func (r *Repository) UpdateTarget(ctx context.Context, target domain.Target) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE calorie_targets
		SET target_kcal = ?, target_protein_g = ?, target_carbs_g = ?, target_fat_g = ?, updated_at = ?
		WHERE id = ?`,
		target.TargetCalories,
		target.TargetProteinGrams,
		nullableInt(target.TargetCarbGrams),
		nullableInt(target.TargetFatGrams),
		target.UpdatedAt,
		target.ID,
	)
	if err != nil {
		return fmt.Errorf("update target: %w", err)
	}

	return nil
}
