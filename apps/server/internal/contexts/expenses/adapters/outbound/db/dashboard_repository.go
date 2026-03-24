package db

import (
	"context"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/expenses/application/ports"
	"github.com/Hughzu/trackstack/apps/server/internal/contexts/expenses/domain"
)

func (r *Repository) GetDashboardSnapshot(ctx context.Context, sheetID string, historyLimit int) (ports.DashboardSnapshot, error) {
	spentByCategory, err := r.getSpentByCategory(ctx, sheetID)
	if err != nil {
		return ports.DashboardSnapshot{}, err
	}

	pendingChecklistItems, err := r.ListPendingChecklistItemsBySheet(ctx, sheetID)
	if err != nil {
		return ports.DashboardSnapshot{}, err
	}

	history, err := r.listRecentEntriesBySheet(ctx, sheetID, historyLimit, 0)
	if err != nil {
		return ports.DashboardSnapshot{}, err
	}

	totalSpent := 0.0
	for _, amount := range spentByCategory {
		totalSpent += amount
	}

	return ports.DashboardSnapshot{
		TotalSpent:            totalSpent,
		SpentByCategory:       spentByCategory,
		PendingChecklistItems: pendingChecklistItems,
		History:               history,
	}, nil
}

func (r *Repository) getSpentByCategory(ctx context.Context, sheetID string) (map[domain.Category]float64, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT category, SUM(amount)
		FROM expense_entries
		WHERE sheet_id = ?
		GROUP BY category`,
		sheetID,
	)
	if err != nil {
		return nil, fmt.Errorf("get spent by category: %w", err)
	}
	defer rows.Close()

	spentByCategory := map[domain.Category]float64{
		domain.CategoryFund:   0,
		domain.CategoryFun:    0,
		domain.CategoryFuture: 0,
	}

	for rows.Next() {
		var category string
		var amount float64
		if err := rows.Scan(&category, &amount); err != nil {
			return nil, fmt.Errorf("scan spent by category: %w", err)
		}
		spentByCategory[domain.Category(category)] = amount
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("spent by category rows: %w", err)
	}

	return spentByCategory, nil
}

func (r *Repository) listRecentEntriesBySheet(ctx context.Context, sheetID string, limit int, offset int) ([]domain.Entry, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, sheet_id, user_id, title, amount, category, date, type, created_at
		FROM expense_entries
		WHERE sheet_id = ?
		ORDER BY date DESC, created_at DESC
		LIMIT ? OFFSET ?`,
		sheetID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list recent entries by sheet: %w", err)
	}
	defer rows.Close()

	var entries []domain.Entry
	for rows.Next() {
		entry, scanErr := scanEntry(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list recent entries rows: %w", err)
	}

	if entries == nil {
		entries = []domain.Entry{}
	}

	return entries, nil
}
