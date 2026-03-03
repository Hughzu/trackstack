package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/expenses"
)

func (s *ExpensesStore) GetTotalSpentBySheet(ctx context.Context, sheetID string) (float64, error) {
	var total sql.NullFloat64
	err := s.db.QueryRowContext(ctx, "SELECT SUM(amount) FROM expense_entries WHERE sheet_id = ?", sheetID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("get total spent: %w", err)
	}
	if total.Valid {
		return total.Float64, nil
	}
	return 0, nil
}

func (s *ExpensesStore) GetSpentByCategory(ctx context.Context, sheetID string) (map[expenses.Category]float64, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT category, SUM(amount) FROM expense_entries WHERE sheet_id = ? GROUP BY category", sheetID)
	if err != nil {
		return nil, fmt.Errorf("get spent by category: %w", err)
	}
	defer rows.Close()

	res := map[expenses.Category]float64{
		expenses.CategoryFund:   0,
		expenses.CategoryFun:    0,
		expenses.CategoryFuture: 0,
	}

	for rows.Next() {
		var cat string
		var amount float64
		if err := rows.Scan(&cat, &amount); err != nil {
			return nil, err
		}
		res[expenses.Category(cat)] = amount
	}
	return res, nil
}

func (s *ExpensesStore) GetPendingChecklistItems(ctx context.Context, sheetID string) ([]expenses.ChecklistItem, error) {
	rows, err := s.db.QueryContext(
		ctx,
		"SELECT id, sheet_id, template_id, title, amount, category, created_at, completed_at, expense_id FROM expense_checklist_items WHERE sheet_id = ? AND completed_at IS NULL ORDER BY created_at ASC",
		sheetID,
	)
	if err != nil {
		return nil, fmt.Errorf("get pending checklist: %w", err)
	}
	defer rows.Close()

	var items []expenses.ChecklistItem
	for rows.Next() {
		var item expenses.ChecklistItem
		var tId sql.NullString
		var completed sql.NullString
		var eId sql.NullString
		if err := rows.Scan(&item.ID, &item.SheetID, &tId, &item.Title, &item.Amount, &item.Category, &item.CreatedAt, &completed, &eId); err != nil {
			return nil, err
		}
		if tId.Valid {
			item.TemplateID = &tId.String
		}
		if completed.Valid {
			item.CompletedAt = &completed.String
		}
		if eId.Valid {
			item.ExpenseID = &eId.String
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *ExpensesStore) GetSheetHistory(ctx context.Context, sheetID string) ([]expenses.Entry, error) {
	rows, err := s.db.QueryContext(
		ctx,
		"SELECT id, sheet_id, user_id, title, amount, category, date, type, created_at FROM expense_entries WHERE sheet_id = ? ORDER BY date DESC, created_at DESC",
		sheetID,
	)
	if err != nil {
		return nil, fmt.Errorf("get sheet history: %w", err)
	}
	defer rows.Close()

	var history []expenses.Entry
	for rows.Next() {
		var entry expenses.Entry
		if err := rows.Scan(&entry.ID, &entry.SheetID, &entry.UserID, &entry.Title, &entry.Amount, &entry.Category, &entry.Date, &entry.Type, &entry.CreatedAt); err != nil {
			return nil, err
		}
		history = append(history, entry)
	}
	if history == nil {
		history = []expenses.Entry{}
	}
	return history, nil
}
