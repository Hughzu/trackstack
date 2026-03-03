package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/expenses"
)

func (s *ExpensesStore) GetChecklistTemplates(ctx context.Context, userID string) ([]expenses.Template, error) {
	rows, err := s.db.QueryContext(
		ctx,
		"SELECT id, user_id, title, amount, category, created_at, updated_at FROM expense_checklist_templates WHERE user_id = ? ORDER BY created_at ASC",
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("get checklist templates: %w", err)
	}
	defer rows.Close()

	var templates []expenses.Template
	for rows.Next() {
		var t expenses.Template
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Amount, &t.Category, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan checklist template: %w", err)
		}
		templates = append(templates, t)
	}
	return templates, nil
}

func (s *ExpensesStore) GetChecklistTemplate(ctx context.Context, templateID string, userID string) (expenses.Template, error) {
	row := s.db.QueryRowContext(
		ctx,
		"SELECT id, user_id, title, amount, category, created_at, updated_at FROM expense_checklist_templates WHERE id = ? AND user_id = ?",
		templateID, userID,
	)

	var t expenses.Template
	err := row.Scan(&t.ID, &t.UserID, &t.Title, &t.Amount, &t.Category, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return expenses.Template{}, fmt.Errorf("checklist template not found")
		}
		return expenses.Template{}, fmt.Errorf("get checklist template: %w", err)
	}
	return t, nil
}

func (s *ExpensesStore) CreateChecklistTemplate(ctx context.Context, template expenses.Template) error {
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO expense_checklist_templates (id, user_id, title, amount, category, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		template.ID, template.UserID, template.Title, template.Amount, string(template.Category), template.CreatedAt, template.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert checklist template: %w", err)
	}
	return nil
}

func (s *ExpensesStore) UpdateChecklistTemplate(ctx context.Context, template expenses.Template) error {
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE expense_checklist_templates SET title = ?, amount = ?, category = ?, updated_at = ? WHERE id = ? AND user_id = ?",
		template.Title, template.Amount, string(template.Category), template.UpdatedAt, template.ID, template.UserID,
	)
	if err != nil {
		return fmt.Errorf("update checklist template: %w", err)
	}
	return nil
}

func (s *ExpensesStore) DeleteChecklistTemplate(ctx context.Context, templateID string, userID string) (bool, error) {
	result, err := s.db.ExecContext(ctx, "DELETE FROM expense_checklist_templates WHERE id = ? AND user_id = ?", templateID, userID)
	if err != nil {
		return false, fmt.Errorf("delete checklist template: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0, nil
}

func (s *ExpensesStore) GetRecurringTemplates(ctx context.Context, userID string) ([]expenses.Template, error) {
	rows, err := s.db.QueryContext(
		ctx,
		"SELECT id, user_id, title, amount, category, created_at, updated_at FROM expense_recurring_templates WHERE user_id = ? ORDER BY created_at ASC",
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("get recurring templates: %w", err)
	}
	defer rows.Close()

	var templates []expenses.Template
	for rows.Next() {
		var t expenses.Template
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Amount, &t.Category, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan recurring template: %w", err)
		}
		templates = append(templates, t)
	}
	return templates, nil
}

func (s *ExpensesStore) GetRecurringTemplate(ctx context.Context, templateID string, userID string) (expenses.Template, error) {
	row := s.db.QueryRowContext(
		ctx,
		"SELECT id, user_id, title, amount, category, created_at, updated_at FROM expense_recurring_templates WHERE id = ? AND user_id = ?",
		templateID, userID,
	)

	var t expenses.Template
	err := row.Scan(&t.ID, &t.UserID, &t.Title, &t.Amount, &t.Category, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return expenses.Template{}, fmt.Errorf("recurring template not found")
		}
		return expenses.Template{}, fmt.Errorf("get recurring template: %w", err)
	}
	return t, nil
}

func (s *ExpensesStore) CreateRecurringTemplate(ctx context.Context, template expenses.Template) error {
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO expense_recurring_templates (id, user_id, title, amount, category, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		template.ID, template.UserID, template.Title, template.Amount, string(template.Category), template.CreatedAt, template.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert recurring template: %w", err)
	}
	return nil
}

func (s *ExpensesStore) UpdateRecurringTemplate(ctx context.Context, template expenses.Template) error {
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE expense_recurring_templates SET title = ?, amount = ?, category = ?, updated_at = ? WHERE id = ? AND user_id = ?",
		template.Title, template.Amount, string(template.Category), template.UpdatedAt, template.ID, template.UserID,
	)
	if err != nil {
		return fmt.Errorf("update recurring template: %w", err)
	}
	return nil
}

func (s *ExpensesStore) DeleteRecurringTemplate(ctx context.Context, templateID string, userID string) (bool, error) {
	result, err := s.db.ExecContext(ctx, "DELETE FROM expense_recurring_templates WHERE id = ? AND user_id = ?", templateID, userID)
	if err != nil {
		return false, fmt.Errorf("delete recurring template: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0, nil
}

func (s *ExpensesStore) GetChecklistItem(ctx context.Context, itemID string, userID string) (expenses.ChecklistItem, error) {
	var item expenses.ChecklistItem
	var tId sql.NullString
	var completed sql.NullString
	var eId sql.NullString

	err := s.db.QueryRowContext(
		ctx,
		`SELECT i.id, i.sheet_id, i.template_id, i.title, i.amount, i.category, i.created_at, i.completed_at, i.expense_id 
		 FROM expense_checklist_items i 
		 JOIN expense_sheets s ON i.sheet_id = s.id 
		 WHERE i.id = ? AND s.user_id = ?`,
		itemID, userID).
		Scan(&item.ID, &item.SheetID, &tId, &item.Title, &item.Amount, &item.Category, &item.CreatedAt, &completed, &eId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return expenses.ChecklistItem{}, fmt.Errorf("checklist item not found")
		}
		return expenses.ChecklistItem{}, fmt.Errorf("get checklist item: %w", err)
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
	return item, nil
}

func (s *ExpensesStore) CreateChecklistItem(ctx context.Context, item expenses.ChecklistItem) error {
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO expense_checklist_items (id, sheet_id, template_id, title, amount, category, created_at, completed_at, expense_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		item.ID, item.SheetID, item.TemplateID, item.Title, item.Amount, string(item.Category), item.CreatedAt, item.CompletedAt, item.ExpenseID,
	)
	if err != nil {
		return fmt.Errorf("insert checklist item: %w", err)
	}
	return nil
}

func (s *ExpensesStore) UpdateChecklistItem(ctx context.Context, item expenses.ChecklistItem) error {
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE expense_checklist_items SET title = ?, amount = ?, category = ?, completed_at = ?, expense_id = ? WHERE id = ?",
		item.Title, item.Amount, string(item.Category), item.CompletedAt, item.ExpenseID, item.ID,
	)
	if err != nil {
		return fmt.Errorf("update checklist item: %w", err)
	}
	return nil
}

func (s *ExpensesStore) UpdateChecklistItemsByTemplate(ctx context.Context, templateID string, title string, amount float64, category expenses.Category) error {
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE expense_checklist_items SET title = ?, amount = ?, category = ? WHERE template_id = ? AND completed_at IS NULL",
		title, amount, string(category), templateID,
	)
	if err != nil {
		return fmt.Errorf("update checklist items by template: %w", err)
	}
	return nil
}

func (s *ExpensesStore) DeletePendingChecklistItemsByTemplate(ctx context.Context, templateID string, userID string) error {
	_, err := s.db.ExecContext(
		ctx,
		"DELETE FROM expense_checklist_items WHERE template_id = ? AND completed_at IS NULL AND sheet_id IN (SELECT id FROM expense_sheets WHERE user_id = ? AND closed_at IS NULL)",
		templateID, userID,
	)
	if err != nil {
		return fmt.Errorf("delete pending checklist items: %w", err)
	}
	return nil
}
