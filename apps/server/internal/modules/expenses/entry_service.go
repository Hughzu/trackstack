package expenses

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (s *Service) AddExpense(ctx context.Context, req AddExpenseRequest) (Entry, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return Entry{}, ErrInvalidInput
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		req.Title = "Untitled"
	}

	sheet, err := s.getOrCreateOpenSheet(ctx, req.UserID)
	if err != nil {
		return Entry{}, err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	date := time.Now().UTC().Format("2006-01-02")
	if req.Date != nil && *req.Date != "" {
		date = *req.Date
	}

	category := CategoryFund
	if req.Category != nil {
		category = Category(*req.Category)
	}

	entry := Entry{
		ID:        uuid.NewString(),
		SheetID:   sheet.ID,
		UserID:    req.UserID,
		Title:     req.Title,
		Amount:    req.Amount,
		Category:  category,
		Date:      date,
		Type:      EntryTypeManual,
		CreatedAt: now,
	}

	if err := s.store.CreateExpense(ctx, entry); err != nil {
		return Entry{}, err
	}

	return entry, nil
}

func (s *Service) DeleteExpense(ctx context.Context, req DeleteExpenseRequest) (bool, error) {
	if strings.TrimSpace(req.UserID) == "" || strings.TrimSpace(req.ID) == "" {
		return false, ErrInvalidInput
	}
	return s.store.DeleteExpense(ctx, req.ID, req.UserID)
}
