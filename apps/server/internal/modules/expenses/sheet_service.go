package expenses

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func formatPeriodKey(date time.Time) string {
	return fmt.Sprintf("%04d-%02d", date.Year(), date.Month())
}

func addPeriodKeyMonth(periodKey string) string {
	var year, month int
	fmt.Sscanf(periodKey, "%d-%d", &year, &month)
	base := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	base = base.AddDate(0, 1, 0)
	return formatPeriodKey(base)
}

func (s *Service) getOrCreateOpenSheet(ctx context.Context, userID string) (Sheet, error) {
	open, err := s.store.GetOpenSheet(ctx, userID)
	if err != nil {
		return Sheet{}, err
	}
	if open != nil {
		return *open, nil
	}

	last, err := s.store.GetLatestSheet(ctx, userID)
	if err != nil {
		return Sheet{}, err
	}

	periodKey := formatPeriodKey(time.Now())
	if last != nil {
		periodKey = addPeriodKeyMonth(last.PeriodKey)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	sheet := Sheet{
		ID:        uuid.NewString(),
		UserID:    userID,
		PeriodKey: periodKey,
		CreatedAt: now,
	}

	if err := s.store.CreateSheet(ctx, sheet); err != nil {
		return Sheet{}, err
	}

	// Bootstrap checklist requirements
	templates, _ := s.store.GetChecklistTemplates(ctx, userID)
	for _, t := range templates {
		item := ChecklistItem{
			ID:         uuid.NewString(),
			SheetID:    sheet.ID,
			TemplateID: &t.ID,
			Title:      t.Title,
			Amount:     t.Amount,
			Category:   t.Category,
			CreatedAt:  now,
		}
		_ = s.store.CreateChecklistItem(ctx, item)
	}

	// Bootstrap recurring expenses
	recurring, _ := s.store.GetRecurringTemplates(ctx, userID)
	recurringDate := fmt.Sprintf("%s-01", periodKey)
	for _, t := range recurring {
		entry := Entry{
			ID:        uuid.NewString(),
			SheetID:   sheet.ID,
			UserID:    userID,
			Title:     t.Title,
			Amount:    t.Amount,
			Category:  t.Category,
			Date:      recurringDate,
			Type:      EntryTypeRecurring,
			CreatedAt: now,
		}
		_ = s.store.CreateExpense(ctx, entry)
	}

	return sheet, nil
}

func (s *Service) CloseSheet(ctx context.Context, req CloseSheetRequest) (Sheet, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return Sheet{}, ErrInvalidInput
	}

	open, err := s.store.GetOpenSheet(ctx, req.UserID)
	if err != nil {
		return Sheet{}, err
	}

	basePeriodKey := formatPeriodKey(time.Now())
	if open != nil {
		now := time.Now().UTC().Format(time.RFC3339)
		open.ClosedAt = &now
		if err := s.store.UpdateSheet(ctx, *open); err != nil {
			return Sheet{}, err
		}
		basePeriodKey = open.PeriodKey
	} else {
		last, err := s.store.GetLatestSheet(ctx, req.UserID)
		if err == nil && last != nil {
			basePeriodKey = last.PeriodKey
		}
	}

	// Create next sheet via bootstrap
	nextPeriodKey := addPeriodKeyMonth(basePeriodKey)
	now := time.Now().UTC().Format(time.RFC3339)
	ns := Sheet{
		ID:        uuid.NewString(),
		UserID:    req.UserID,
		PeriodKey: nextPeriodKey,
		CreatedAt: now,
	}
	if err := s.store.CreateSheet(ctx, ns); err != nil {
		return Sheet{}, err
	}

	return s.getOrCreateOpenSheet(ctx, req.UserID)
}
