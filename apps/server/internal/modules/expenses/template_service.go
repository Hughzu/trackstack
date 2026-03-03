package expenses

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (s *Service) UpsertChecklistTemplate(ctx context.Context, req UpsertTemplateRequest) (Template, error) {
	if strings.TrimSpace(req.UserID) == "" || strings.TrimSpace(req.Title) == "" {
		return Template{}, ErrInvalidInput
	}

	category := CategoryFund
	if req.Category != nil {
		category = Category(*req.Category)
	}

	now := time.Now().UTC().Format(time.RFC3339)

	if req.ID != nil {
		tmpl, err := s.store.GetChecklistTemplate(ctx, *req.ID, req.UserID)
		if err == nil {
			tmpl.Title = req.Title
			tmpl.Amount = req.Amount
			tmpl.Category = category
			tmpl.UpdatedAt = now

			if err := s.store.UpdateChecklistTemplate(ctx, tmpl); err != nil {
				return Template{}, err
			}

			// Sync uncompleted items
			_ = s.store.UpdateChecklistItemsByTemplate(ctx, tmpl.ID, tmpl.Title, tmpl.Amount, tmpl.Category)
			return tmpl, nil
		}
	}

	tmpl := Template{
		ID:        uuid.NewString(),
		UserID:    req.UserID,
		Title:     req.Title,
		Amount:    req.Amount,
		Category:  category,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.store.CreateChecklistTemplate(ctx, tmpl); err != nil {
		return Template{}, err
	}

	// Add to open sheet if it exists
	if sheet, err := s.getOrCreateOpenSheet(ctx, req.UserID); err == nil {
		item := ChecklistItem{
			ID:         uuid.NewString(),
			SheetID:    sheet.ID,
			TemplateID: &tmpl.ID,
			Title:      tmpl.Title,
			Amount:     tmpl.Amount,
			Category:   tmpl.Category,
			CreatedAt:  now,
		}
		_ = s.store.CreateChecklistItem(ctx, item)
	}

	return tmpl, nil
}

func (s *Service) DeleteChecklistTemplate(ctx context.Context, req DeleteTemplateRequest) (bool, error) {
	if strings.TrimSpace(req.UserID) == "" || strings.TrimSpace(req.ID) == "" {
		return false, ErrInvalidInput
	}

	// Clean up pending items first
	_ = s.store.DeletePendingChecklistItemsByTemplate(ctx, req.ID, req.UserID)
	return s.store.DeleteChecklistTemplate(ctx, req.ID, req.UserID)
}

func (s *Service) CompleteChecklistItem(ctx context.Context, req CompleteChecklistItemRequest) (Entry, error) {
	if strings.TrimSpace(req.UserID) == "" || strings.TrimSpace(req.ID) == "" {
		return Entry{}, ErrInvalidInput
	}

	item, err := s.store.GetChecklistItem(ctx, req.ID, req.UserID)
	if err != nil {
		return Entry{}, err
	}

	if item.CompletedAt != nil {
		return Entry{}, errors.New("already completed")
	}

	now := time.Now().UTC().Format(time.RFC3339)
	date := time.Now().UTC().Format("2006-01-02")
	if req.Date != nil && *req.Date != "" {
		date = *req.Date
	}

	entry := Entry{
		ID:        uuid.NewString(),
		SheetID:   item.SheetID,
		UserID:    req.UserID,
		Title:     item.Title,
		Amount:    item.Amount,
		Category:  item.Category,
		Date:      date,
		Type:      EntryTypeChecklist,
		CreatedAt: now,
	}

	if err := s.store.CreateExpense(ctx, entry); err != nil {
		return Entry{}, err
	}

	item.CompletedAt = &now
	item.ExpenseID = &entry.ID
	if err := s.store.UpdateChecklistItem(ctx, item); err != nil {
		return Entry{}, err
	}

	return entry, nil
}

func (s *Service) UpsertRecurringTemplate(ctx context.Context, req UpsertTemplateRequest) (Template, error) {
	if strings.TrimSpace(req.UserID) == "" || strings.TrimSpace(req.Title) == "" {
		return Template{}, ErrInvalidInput
	}

	category := CategoryFund
	if req.Category != nil {
		category = Category(*req.Category)
	}
	now := time.Now().UTC().Format(time.RFC3339)

	if req.ID != nil {
		if tmpl, err := s.store.GetRecurringTemplate(ctx, *req.ID, req.UserID); err == nil {
			tmpl.Title = req.Title
			tmpl.Amount = req.Amount
			tmpl.Category = category
			tmpl.UpdatedAt = now
			if err := s.store.UpdateRecurringTemplate(ctx, tmpl); err != nil {
				return Template{}, err
			}
			return tmpl, nil
		}
	}

	tmpl := Template{
		ID:        uuid.NewString(),
		UserID:    req.UserID,
		Title:     req.Title,
		Amount:    req.Amount,
		Category:  category,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.store.CreateRecurringTemplate(ctx, tmpl); err != nil {
		return Template{}, err
	}

	if sheet, err := s.getOrCreateOpenSheet(ctx, req.UserID); err == nil {
		entry := Entry{
			ID:        uuid.NewString(),
			SheetID:   sheet.ID,
			UserID:    req.UserID,
			Title:     tmpl.Title,
			Amount:    tmpl.Amount,
			Category:  tmpl.Category,
			Date:      time.Now().UTC().Format("2006-01-02"),
			Type:      EntryTypeRecurring,
			CreatedAt: now,
		}
		_ = s.store.CreateExpense(ctx, entry)
	}

	return tmpl, nil
}

func (s *Service) DeleteRecurringTemplate(ctx context.Context, req DeleteTemplateRequest) (bool, error) {
	if strings.TrimSpace(req.UserID) == "" || strings.TrimSpace(req.ID) == "" {
		return false, ErrInvalidInput
	}
	return s.store.DeleteRecurringTemplate(ctx, req.ID, req.UserID)
}
