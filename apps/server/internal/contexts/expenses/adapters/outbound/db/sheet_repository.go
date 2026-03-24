package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/expenses/domain"
)

func (r *Repository) FindLatestSheet(ctx context.Context, userID string) (*domain.Sheet, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, period_key, created_at, closed_at
		FROM expense_sheets
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT 1`,
		userID,
	)

	sheet, found, err := scanSheet(row)
	if err != nil {
		return nil, fmt.Errorf("find latest sheet: %w", err)
	}
	if !found {
		return nil, nil
	}

	return &sheet, nil
}

func (r *Repository) FindOpenSheet(ctx context.Context, userID string) (*domain.Sheet, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, period_key, created_at, closed_at
		FROM expense_sheets
		WHERE user_id = ? AND closed_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1`,
		userID,
	)

	sheet, found, err := scanSheet(row)
	if err != nil {
		return nil, fmt.Errorf("find open sheet: %w", err)
	}
	if !found {
		return nil, nil
	}

	return &sheet, nil
}

func (r *Repository) CreateSheet(ctx context.Context, sheet domain.Sheet) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO expense_sheets (id, user_id, period_key, created_at, closed_at)
		VALUES (?, ?, ?, ?, ?)`,
		sheet.ID,
		sheet.UserID,
		sheet.PeriodKey,
		sheet.CreatedAt,
		nullableString(sheet.ClosedAt),
	)
	if err != nil {
		return fmt.Errorf("create sheet: %w", err)
	}

	return nil
}

func (r *Repository) UpdateSheet(ctx context.Context, sheet domain.Sheet) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE expense_sheets
		SET closed_at = ?
		WHERE id = ?`,
		nullableString(sheet.ClosedAt),
		sheet.ID,
	)
	if err != nil {
		return fmt.Errorf("update sheet: %w", err)
	}

	return nil
}

func scanSheet(scanner interface{ Scan(dest ...any) error }) (domain.Sheet, bool, error) {
	var sheet domain.Sheet
	var closedAt sql.NullString

	err := scanner.Scan(&sheet.ID, &sheet.UserID, &sheet.PeriodKey, &sheet.CreatedAt, &closedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Sheet{}, false, nil
		}
		return domain.Sheet{}, false, err
	}

	if closedAt.Valid {
		sheet.ClosedAt = &closedAt.String
	}

	return sheet, true, nil
}
