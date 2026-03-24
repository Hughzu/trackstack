package db

import (
	"context"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/expenses/domain"
)

func (r *Repository) CreateEntry(ctx context.Context, entry domain.Entry) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO expense_entries (id, sheet_id, user_id, title, amount, category, date, type, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		entry.ID,
		entry.SheetID,
		entry.UserID,
		entry.Title,
		entry.Amount,
		string(entry.Category),
		entry.Date,
		string(entry.Type),
		entry.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create entry: %w", err)
	}

	return nil
}

func (r *Repository) DeleteEntry(ctx context.Context, userID string, id string) (bool, error) {
	result, err := r.db.ExecContext(
		ctx,
		"DELETE FROM expense_entries WHERE id = ? AND user_id = ?",
		id,
		userID,
	)
	if err != nil {
		return false, fmt.Errorf("delete entry: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("delete entry rows: %w", err)
	}

	return rowsAffected > 0, nil
}

func scanEntry(scanner interface{ Scan(dest ...any) error }) (domain.Entry, error) {
	var entry domain.Entry

	err := scanner.Scan(
		&entry.ID,
		&entry.SheetID,
		&entry.UserID,
		&entry.Title,
		&entry.Amount,
		&entry.Category,
		&entry.Date,
		&entry.Type,
		&entry.CreatedAt,
	)
	if err != nil {
		return domain.Entry{}, fmt.Errorf("scan entry: %w", err)
	}

	return entry, nil
}
