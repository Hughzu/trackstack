package db

import (
	"context"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/expenses"
)

func (s *ExpensesStore) GetSettings(ctx context.Context, userID string) (expenses.Settings, error) {
	row := s.db.QueryRowContext(
		ctx,
		"SELECT id, user_id, income, ratio_fund, ratio_fun, ratio_future, created_at, updated_at FROM expense_settings WHERE user_id = ?",
		userID,
	)

	var settings expenses.Settings
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
		return expenses.Settings{}, err
	}

	return settings, nil
}

func (s *ExpensesStore) CreateSettings(ctx context.Context, settings expenses.Settings) error {
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO expense_settings (id, user_id, income, ratio_fund, ratio_fun, ratio_future, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		settings.ID, settings.UserID, settings.Income, settings.RatioFund, settings.RatioFun, settings.RatioFuture, settings.CreatedAt, settings.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create settings: %w", err)
	}
	return nil
}

func (s *ExpensesStore) UpdateSettings(ctx context.Context, settings expenses.Settings) error {
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE expense_settings SET income = ?, ratio_fund = ?, ratio_fun = ?, ratio_future = ?, updated_at = ? WHERE id = ?",
		settings.Income, settings.RatioFund, settings.RatioFun, settings.RatioFuture, settings.UpdatedAt, settings.ID,
	)
	if err != nil {
		return fmt.Errorf("update settings: %w", err)
	}

	return nil
}
