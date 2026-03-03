package db

import (
	"context"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/expenses"
)

func (s *ExpensesStore) CreateExpense(ctx context.Context, entry expenses.Entry) error {
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO expense_entries (id, sheet_id, user_id, title, amount, category, date, type, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		entry.ID, entry.SheetID, entry.UserID, entry.Title, entry.Amount, string(entry.Category), entry.Date, string(entry.Type), entry.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert expense entry: %w", err)
	}
	return nil
}

func (s *ExpensesStore) DeleteExpense(ctx context.Context, entryID string, userID string) (bool, error) {
	result, err := s.db.ExecContext(ctx, "DELETE FROM expense_entries WHERE id = ? AND user_id = ?", entryID, userID)
	if err != nil {
		return false, fmt.Errorf("delete expense: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0, nil
}
