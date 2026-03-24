package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/expenses/domain"
)

func (r *Repository) FindChecklistItem(ctx context.Context, userID string, id string) (domain.ChecklistItem, bool, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT i.id, i.sheet_id, i.template_id, i.title, i.amount, i.category, i.created_at, i.completed_at, i.expense_id
		FROM expense_checklist_items i
		JOIN expense_sheets s ON i.sheet_id = s.id
		WHERE i.id = ? AND s.user_id = ?`,
		id,
		userID,
	)

	item, found, err := scanChecklistItem(row)
	if err != nil {
		return domain.ChecklistItem{}, false, fmt.Errorf("find checklist item: %w", err)
	}

	return item, found, nil
}

func (r *Repository) ListPendingChecklistItemsBySheet(ctx context.Context, sheetID string) ([]domain.ChecklistItem, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, sheet_id, template_id, title, amount, category, created_at, completed_at, expense_id
		FROM expense_checklist_items
		WHERE sheet_id = ? AND completed_at IS NULL
		ORDER BY created_at ASC`,
		sheetID,
	)
	if err != nil {
		return nil, fmt.Errorf("list pending checklist items: %w", err)
	}
	defer rows.Close()

	var items []domain.ChecklistItem
	for rows.Next() {
		item, _, scanErr := scanChecklistItem(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list pending checklist items rows: %w", err)
	}

	return items, nil
}

func (r *Repository) CreateChecklistItem(ctx context.Context, item domain.ChecklistItem) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO expense_checklist_items (id, sheet_id, template_id, title, amount, category, created_at, completed_at, expense_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		item.ID,
		item.SheetID,
		nullableString(item.TemplateID),
		item.Title,
		item.Amount,
		string(item.Category),
		item.CreatedAt,
		nullableString(item.CompletedAt),
		nullableString(item.ExpenseID),
	)
	if err != nil {
		return fmt.Errorf("create checklist item: %w", err)
	}

	return nil
}

func (r *Repository) UpdateChecklistItem(ctx context.Context, item domain.ChecklistItem) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE expense_checklist_items
		SET title = ?, amount = ?, category = ?, completed_at = ?, expense_id = ?
		WHERE id = ?`,
		item.Title,
		item.Amount,
		string(item.Category),
		nullableString(item.CompletedAt),
		nullableString(item.ExpenseID),
		item.ID,
	)
	if err != nil {
		return fmt.Errorf("update checklist item: %w", err)
	}

	return nil
}

func (r *Repository) UpdatePendingChecklistItemsByTemplate(ctx context.Context, templateID string, title string, amount float64, category domain.Category) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE expense_checklist_items
		SET title = ?, amount = ?, category = ?
		WHERE template_id = ? AND completed_at IS NULL`,
		title,
		amount,
		string(category),
		templateID,
	)
	if err != nil {
		return fmt.Errorf("update pending checklist items by template: %w", err)
	}

	return nil
}

func (r *Repository) DeletePendingChecklistItemsByTemplate(ctx context.Context, userID string, templateID string) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM expense_checklist_items
		WHERE template_id = ?
		AND completed_at IS NULL
		AND sheet_id IN (SELECT id FROM expense_sheets WHERE user_id = ? AND closed_at IS NULL)`,
		templateID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("delete pending checklist items by template: %w", err)
	}

	return nil
}

func scanChecklistItem(scanner interface{ Scan(dest ...any) error }) (domain.ChecklistItem, bool, error) {
	var item domain.ChecklistItem
	var templateID sql.NullString
	var completedAt sql.NullString
	var expenseID sql.NullString

	err := scanner.Scan(
		&item.ID,
		&item.SheetID,
		&templateID,
		&item.Title,
		&item.Amount,
		&item.Category,
		&item.CreatedAt,
		&completedAt,
		&expenseID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ChecklistItem{}, false, nil
		}
		return domain.ChecklistItem{}, false, fmt.Errorf("scan checklist item: %w", err)
	}

	if templateID.Valid {
		item.TemplateID = &templateID.String
	}
	if completedAt.Valid {
		item.CompletedAt = &completedAt.String
	}
	if expenseID.Valid {
		item.ExpenseID = &expenseID.String
	}

	return item, true, nil
}
