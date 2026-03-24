package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/expenses/domain"
)

func (r *Repository) FindChecklistTemplate(ctx context.Context, userID string, id string) (domain.Template, bool, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, title, amount, category, created_at, updated_at
		FROM expense_checklist_templates
		WHERE id = ? AND user_id = ?`,
		id,
		userID,
	)

	template, found, err := scanTemplate(row)
	if err != nil {
		return domain.Template{}, false, fmt.Errorf("find checklist template: %w", err)
	}

	return template, found, nil
}

func (r *Repository) ListChecklistTemplates(ctx context.Context, userID string) ([]domain.Template, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, title, amount, category, created_at, updated_at
		FROM expense_checklist_templates
		WHERE user_id = ?
		ORDER BY created_at ASC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list checklist templates: %w", err)
	}
	defer rows.Close()

	var templates []domain.Template
	for rows.Next() {
		template, _, scanErr := scanTemplate(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		templates = append(templates, template)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list checklist templates rows: %w", err)
	}

	return templates, nil
}

func (r *Repository) CreateChecklistTemplate(ctx context.Context, template domain.Template) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO expense_checklist_templates (id, user_id, title, amount, category, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		template.ID,
		template.UserID,
		template.Title,
		template.Amount,
		string(template.Category),
		template.CreatedAt,
		template.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create checklist template: %w", err)
	}

	return nil
}

func (r *Repository) UpdateChecklistTemplate(ctx context.Context, template domain.Template) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE expense_checklist_templates
		SET title = ?, amount = ?, category = ?, updated_at = ?
		WHERE id = ? AND user_id = ?`,
		template.Title,
		template.Amount,
		string(template.Category),
		template.UpdatedAt,
		template.ID,
		template.UserID,
	)
	if err != nil {
		return fmt.Errorf("update checklist template: %w", err)
	}

	return nil
}

func (r *Repository) DeleteChecklistTemplate(ctx context.Context, userID string, id string) (bool, error) {
	result, err := r.db.ExecContext(
		ctx,
		"DELETE FROM expense_checklist_templates WHERE id = ? AND user_id = ?",
		id,
		userID,
	)
	if err != nil {
		return false, fmt.Errorf("delete checklist template: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("delete checklist template rows: %w", err)
	}

	return rowsAffected > 0, nil
}

func (r *Repository) FindRecurringTemplate(ctx context.Context, userID string, id string) (domain.Template, bool, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, title, amount, category, created_at, updated_at
		FROM expense_recurring_templates
		WHERE id = ? AND user_id = ?`,
		id,
		userID,
	)

	template, found, err := scanTemplate(row)
	if err != nil {
		return domain.Template{}, false, fmt.Errorf("find recurring template: %w", err)
	}

	return template, found, nil
}

func (r *Repository) ListRecurringTemplates(ctx context.Context, userID string) ([]domain.Template, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, title, amount, category, created_at, updated_at
		FROM expense_recurring_templates
		WHERE user_id = ?
		ORDER BY created_at ASC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list recurring templates: %w", err)
	}
	defer rows.Close()

	var templates []domain.Template
	for rows.Next() {
		template, _, scanErr := scanTemplate(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		templates = append(templates, template)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list recurring templates rows: %w", err)
	}

	return templates, nil
}

func (r *Repository) CreateRecurringTemplate(ctx context.Context, template domain.Template) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO expense_recurring_templates (id, user_id, title, amount, category, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		template.ID,
		template.UserID,
		template.Title,
		template.Amount,
		string(template.Category),
		template.CreatedAt,
		template.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create recurring template: %w", err)
	}

	return nil
}

func (r *Repository) UpdateRecurringTemplate(ctx context.Context, template domain.Template) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE expense_recurring_templates
		SET title = ?, amount = ?, category = ?, updated_at = ?
		WHERE id = ? AND user_id = ?`,
		template.Title,
		template.Amount,
		string(template.Category),
		template.UpdatedAt,
		template.ID,
		template.UserID,
	)
	if err != nil {
		return fmt.Errorf("update recurring template: %w", err)
	}

	return nil
}

func (r *Repository) DeleteRecurringTemplate(ctx context.Context, userID string, id string) (bool, error) {
	result, err := r.db.ExecContext(
		ctx,
		"DELETE FROM expense_recurring_templates WHERE id = ? AND user_id = ?",
		id,
		userID,
	)
	if err != nil {
		return false, fmt.Errorf("delete recurring template: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("delete recurring template rows: %w", err)
	}

	return rowsAffected > 0, nil
}

func scanTemplate(scanner interface{ Scan(dest ...any) error }) (domain.Template, bool, error) {
	var template domain.Template

	err := scanner.Scan(
		&template.ID,
		&template.UserID,
		&template.Title,
		&template.Amount,
		&template.Category,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Template{}, false, nil
		}
		return domain.Template{}, false, fmt.Errorf("scan template: %w", err)
	}

	return template, true, nil
}
