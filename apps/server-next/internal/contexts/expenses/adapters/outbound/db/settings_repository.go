package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/expenses/domain"
)

func (r *Repository) FindSettings(ctx context.Context, userID string) (domain.Settings, bool, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, income, ratio_fund, ratio_fun, ratio_future, created_at, updated_at
		FROM expense_settings
		WHERE user_id = ?`,
		userID,
	)

	var settings domain.Settings
	err := row.Scan(
		&settings.ID,
		&settings.UserID,
		&settings.Income,
		&settings.RatioFund,
		&settings.RatioFun,
		&settings.RatioFuture,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Settings{}, false, nil
		}
		return domain.Settings{}, false, fmt.Errorf("find settings: %w", err)
	}

	return settings, true, nil
}

func (r *Repository) CreateSettings(ctx context.Context, settings domain.Settings) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO expense_settings (id, user_id, income, ratio_fund, ratio_fun, ratio_future, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		settings.ID,
		settings.UserID,
		settings.Income,
		settings.RatioFund,
		settings.RatioFun,
		settings.RatioFuture,
		settings.CreatedAt,
		settings.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create settings: %w", err)
	}

	return nil
}

func (r *Repository) UpdateSettings(ctx context.Context, settings domain.Settings) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE expense_settings
		SET income = ?, ratio_fund = ?, ratio_fun = ?, ratio_future = ?, updated_at = ?
		WHERE id = ?`,
		settings.Income,
		settings.RatioFund,
		settings.RatioFun,
		settings.RatioFuture,
		settings.UpdatedAt,
		settings.ID,
	)
	if err != nil {
		return fmt.Errorf("update settings: %w", err)
	}

	return nil
}
