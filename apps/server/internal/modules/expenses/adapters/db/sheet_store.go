package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/expenses"
)

func (s *ExpensesStore) GetLatestSheet(ctx context.Context, userID string) (*expenses.Sheet, error) {
	row := s.db.QueryRowContext(ctx, "SELECT id, user_id, period_key, created_at, closed_at FROM expense_sheets WHERE user_id = ? ORDER BY created_at DESC LIMIT 1", userID)
	var sheet expenses.Sheet
	var closedAt sql.NullString
	err := row.Scan(&sheet.ID, &sheet.UserID, &sheet.PeriodKey, &sheet.CreatedAt, &closedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get latest sheet: %w", err)
	}
	if closedAt.Valid {
		sheet.ClosedAt = &closedAt.String
	}
	return &sheet, nil
}

func (s *ExpensesStore) GetOpenSheet(ctx context.Context, userID string) (*expenses.Sheet, error) {
	row := s.db.QueryRowContext(ctx, "SELECT id, user_id, period_key, created_at, closed_at FROM expense_sheets WHERE user_id = ? AND closed_at IS NULL ORDER BY created_at DESC LIMIT 1", userID)
	var sheet expenses.Sheet
	var closedAt sql.NullString
	err := row.Scan(&sheet.ID, &sheet.UserID, &sheet.PeriodKey, &sheet.CreatedAt, &closedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get open sheet: %w", err)
	}
	if closedAt.Valid {
		sheet.ClosedAt = &closedAt.String
	}
	return &sheet, nil
}

func (s *ExpensesStore) CreateSheet(ctx context.Context, sheet expenses.Sheet) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO expense_sheets (id, user_id, period_key, created_at, closed_at) VALUES (?, ?, ?, ?, ?)", sheet.ID, sheet.UserID, sheet.PeriodKey, sheet.CreatedAt, sheet.ClosedAt)
	if err != nil {
		return fmt.Errorf("insert new sheet: %w", err)
	}
	return nil
}

func (s *ExpensesStore) UpdateSheet(ctx context.Context, sheet expenses.Sheet) error {
	_, err := s.db.ExecContext(ctx, "UPDATE expense_sheets SET closed_at = ? WHERE id = ?", sheet.ClosedAt, sheet.ID)
	if err != nil {
		return fmt.Errorf("update sheet: %w", err)
	}
	return nil
}
