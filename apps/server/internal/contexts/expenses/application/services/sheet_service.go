package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/expenses/application/ports"
	"github.com/Hughzu/trackstack/apps/server/internal/contexts/expenses/domain"
	"github.com/google/uuid"
)

type sheetService struct {
	sheetManager SheetManager
}

type SheetManager interface {
	GetOrCreateOpenSheet(ctx context.Context, userID string) (domain.Sheet, error)
	CloseAndCreateNextSheet(ctx context.Context, userID string) (domain.Sheet, error)
}

type SheetManagerDeps struct {
	SheetRepo             ports.SheetRepository
	ChecklistTemplateRepo ports.ChecklistTemplateRepository
	RecurringTemplateRepo ports.RecurringTemplateRepository
	ChecklistItemRepo     ports.ChecklistItemRepository
	EntryRepo             ports.EntryRepository
}

type sheetManager struct {
	sheetRepo             ports.SheetRepository
	checklistTemplateRepo ports.ChecklistTemplateRepository
	recurringTemplateRepo ports.RecurringTemplateRepository
	checklistItemRepo     ports.ChecklistItemRepository
	entryRepo             ports.EntryRepository
}

func NewSheetManager(deps SheetManagerDeps) SheetManager {
	return &sheetManager{
		sheetRepo:             deps.SheetRepo,
		checklistTemplateRepo: deps.ChecklistTemplateRepo,
		recurringTemplateRepo: deps.RecurringTemplateRepo,
		checklistItemRepo:     deps.ChecklistItemRepo,
		entryRepo:             deps.EntryRepo,
	}
}

func NewSheetService(sheetManager SheetManager) ports.SheetUseCase {
	return &sheetService{
		sheetManager: sheetManager,
	}
}

func (s *sheetService) CloseSheet(ctx context.Context, command ports.CloseSheetCommand) (domain.Sheet, error) {
	if strings.TrimSpace(command.UserID) == "" {
		return domain.Sheet{}, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}

	return s.sheetManager.CloseAndCreateNextSheet(ctx, command.UserID)
}

func (m *sheetManager) CloseAndCreateNextSheet(ctx context.Context, userID string) (domain.Sheet, error) {
	openSheet, err := m.sheetRepo.FindOpenSheet(ctx, userID)
	if err != nil {
		return domain.Sheet{}, err
	}

	basePeriodKey := formatPeriodKey(time.Now().UTC())
	if openSheet != nil {
		now := time.Now().UTC().Format(time.RFC3339)
		openSheet.ClosedAt = &now
		if err := m.sheetRepo.UpdateSheet(ctx, *openSheet); err != nil {
			return domain.Sheet{}, err
		}
		basePeriodKey = openSheet.PeriodKey
	} else {
		latestSheet, err := m.sheetRepo.FindLatestSheet(ctx, userID)
		if err != nil {
			return domain.Sheet{}, err
		}
		if latestSheet != nil {
			basePeriodKey = latestSheet.PeriodKey
		}
	}

	return m.createSheetWithBootstrap(ctx, userID, addPeriodKeyMonth(basePeriodKey))
}

func (m *sheetManager) GetOrCreateOpenSheet(ctx context.Context, userID string) (domain.Sheet, error) {
	openSheet, err := m.sheetRepo.FindOpenSheet(ctx, userID)
	if err != nil {
		return domain.Sheet{}, err
	}
	if openSheet != nil {
		return *openSheet, nil
	}

	latestSheet, err := m.sheetRepo.FindLatestSheet(ctx, userID)
	if err != nil {
		return domain.Sheet{}, err
	}

	periodKey := formatPeriodKey(time.Now().UTC())
	if latestSheet != nil {
		periodKey = addPeriodKeyMonth(latestSheet.PeriodKey)
	}

	return m.createSheetWithBootstrap(ctx, userID, periodKey)
}

func (m *sheetManager) createSheetWithBootstrap(ctx context.Context, userID string, periodKey string) (domain.Sheet, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	sheet := domain.Sheet{
		ID:        uuid.NewString(),
		UserID:    userID,
		PeriodKey: periodKey,
		CreatedAt: now,
	}

	if err := m.sheetRepo.CreateSheet(ctx, sheet); err != nil {
		return domain.Sheet{}, err
	}

	checklistTemplates, err := m.checklistTemplateRepo.ListChecklistTemplates(ctx, userID)
	if err != nil {
		return domain.Sheet{}, err
	}
	for _, template := range checklistTemplates {
		item := domain.ChecklistItem{
			ID:         uuid.NewString(),
			SheetID:    sheet.ID,
			TemplateID: &template.ID,
			Title:      template.Title,
			Amount:     template.Amount,
			Category:   template.Category,
			CreatedAt:  now,
		}
		if err := m.checklistItemRepo.CreateChecklistItem(ctx, item); err != nil {
			return domain.Sheet{}, err
		}
	}

	recurringTemplates, err := m.recurringTemplateRepo.ListRecurringTemplates(ctx, userID)
	if err != nil {
		return domain.Sheet{}, err
	}
	recurringDate := fmt.Sprintf("%s-01", periodKey)
	for _, template := range recurringTemplates {
		entry := domain.Entry{
			ID:        uuid.NewString(),
			SheetID:   sheet.ID,
			UserID:    userID,
			Title:     template.Title,
			Amount:    template.Amount,
			Category:  template.Category,
			Date:      recurringDate,
			Type:      domain.EntryTypeRecurring,
			CreatedAt: now,
		}
		if err := m.entryRepo.CreateEntry(ctx, entry); err != nil {
			return domain.Sheet{}, err
		}
	}

	return sheet, nil
}

func formatPeriodKey(date time.Time) string {
	return fmt.Sprintf("%04d-%02d", date.Year(), date.Month())
}

func addPeriodKeyMonth(periodKey string) string {
	var year int
	var month int
	_, _ = fmt.Sscanf(periodKey, "%d-%d", &year, &month)
	base := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	return formatPeriodKey(base.AddDate(0, 1, 0))
}
