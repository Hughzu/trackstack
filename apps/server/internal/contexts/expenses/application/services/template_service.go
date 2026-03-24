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

type templateService struct {
	checklistTemplateRepo ports.ChecklistTemplateRepository
	recurringTemplateRepo ports.RecurringTemplateRepository
	checklistItemRepo     ports.ChecklistItemRepository
	entryRepo             ports.EntryRepository
	sheetManager          SheetManager
}

type TemplateServiceDeps struct {
	ChecklistTemplateRepo ports.ChecklistTemplateRepository
	RecurringTemplateRepo ports.RecurringTemplateRepository
	ChecklistItemRepo     ports.ChecklistItemRepository
	EntryRepo             ports.EntryRepository
	SheetManager          SheetManager
}

func NewTemplateService(deps TemplateServiceDeps) ports.TemplateUseCase {
	return &templateService{
		checklistTemplateRepo: deps.ChecklistTemplateRepo,
		recurringTemplateRepo: deps.RecurringTemplateRepo,
		checklistItemRepo:     deps.ChecklistItemRepo,
		entryRepo:             deps.EntryRepo,
		sheetManager:          deps.SheetManager,
	}
}

func (s *templateService) UpsertChecklist(ctx context.Context, command ports.UpsertTemplateCommand) (domain.Template, error) {
	if strings.TrimSpace(command.UserID) == "" {
		return domain.Template{}, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}
	if strings.TrimSpace(command.Title) == "" {
		return domain.Template{}, fmt.Errorf("%w: title is required", domain.ErrInvalidInput)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	category := resolveCategory(command.Category)

	if command.ID != nil {
		template, found, err := s.checklistTemplateRepo.FindChecklistTemplate(ctx, command.UserID, *command.ID)
		if err != nil {
			return domain.Template{}, err
		}
		if found {
			template.Title = command.Title
			template.Amount = command.Amount
			template.Category = category
			template.UpdatedAt = now

			if err := s.checklistTemplateRepo.UpdateChecklistTemplate(ctx, template); err != nil {
				return domain.Template{}, err
			}
			if err := s.checklistItemRepo.UpdatePendingChecklistItemsByTemplate(ctx, template.ID, template.Title, template.Amount, template.Category); err != nil {
				return domain.Template{}, err
			}

			return template, nil
		}
	}

	sheet, err := s.sheetManager.GetOrCreateOpenSheet(ctx, command.UserID)
	if err != nil {
		return domain.Template{}, err
	}

	template := domain.Template{
		ID:        uuid.NewString(),
		UserID:    command.UserID,
		Title:     command.Title,
		Amount:    command.Amount,
		Category:  category,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.checklistTemplateRepo.CreateChecklistTemplate(ctx, template); err != nil {
		return domain.Template{}, err
	}

	item := domain.ChecklistItem{
		ID:         uuid.NewString(),
		SheetID:    sheet.ID,
		TemplateID: &template.ID,
		Title:      template.Title,
		Amount:     template.Amount,
		Category:   template.Category,
		CreatedAt:  now,
	}
	if err := s.checklistItemRepo.CreateChecklistItem(ctx, item); err != nil {
		return domain.Template{}, err
	}

	return template, nil
}

func (s *templateService) DeleteChecklist(ctx context.Context, command ports.DeleteTemplateCommand) (bool, error) {
	if strings.TrimSpace(command.UserID) == "" {
		return false, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}
	if strings.TrimSpace(command.ID) == "" {
		return false, fmt.Errorf("%w: template id is required", domain.ErrInvalidInput)
	}

	if err := s.checklistItemRepo.DeletePendingChecklistItemsByTemplate(ctx, command.UserID, command.ID); err != nil {
		return false, err
	}

	return s.checklistTemplateRepo.DeleteChecklistTemplate(ctx, command.UserID, command.ID)
}

func (s *templateService) CompleteChecklistItem(ctx context.Context, command ports.CompleteChecklistItemCommand) (domain.Entry, error) {
	if strings.TrimSpace(command.UserID) == "" {
		return domain.Entry{}, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}
	if strings.TrimSpace(command.ID) == "" {
		return domain.Entry{}, fmt.Errorf("%w: checklist item id is required", domain.ErrInvalidInput)
	}

	item, found, err := s.checklistItemRepo.FindChecklistItem(ctx, command.UserID, command.ID)
	if err != nil {
		return domain.Entry{}, err
	}
	if !found {
		return domain.Entry{}, fmt.Errorf("%w: checklist item not found", domain.ErrNotFound)
	}
	if item.CompletedAt != nil {
		return domain.Entry{}, fmt.Errorf("%w: already completed", domain.ErrInvalidInput)
	}

	now := time.Now().UTC()
	date := now.Format("2006-01-02")
	if command.Date != nil && *command.Date != "" {
		date = *command.Date
	}

	entry := domain.Entry{
		ID:        uuid.NewString(),
		SheetID:   item.SheetID,
		UserID:    command.UserID,
		Title:     item.Title,
		Amount:    item.Amount,
		Category:  item.Category,
		Date:      date,
		Type:      domain.EntryTypeChecklist,
		CreatedAt: now.Format(time.RFC3339),
	}
	if err := s.entryRepo.CreateEntry(ctx, entry); err != nil {
		return domain.Entry{}, err
	}

	completedAt := now.Format(time.RFC3339)
	item.CompletedAt = &completedAt
	item.ExpenseID = &entry.ID
	if err := s.checklistItemRepo.UpdateChecklistItem(ctx, item); err != nil {
		return domain.Entry{}, err
	}

	return entry, nil
}

func (s *templateService) UpsertRecurring(ctx context.Context, command ports.UpsertTemplateCommand) (domain.Template, error) {
	if strings.TrimSpace(command.UserID) == "" {
		return domain.Template{}, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}
	if strings.TrimSpace(command.Title) == "" {
		return domain.Template{}, fmt.Errorf("%w: title is required", domain.ErrInvalidInput)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	category := resolveCategory(command.Category)

	if command.ID != nil {
		template, found, err := s.recurringTemplateRepo.FindRecurringTemplate(ctx, command.UserID, *command.ID)
		if err != nil {
			return domain.Template{}, err
		}
		if found {
			template.Title = command.Title
			template.Amount = command.Amount
			template.Category = category
			template.UpdatedAt = now

			if err := s.recurringTemplateRepo.UpdateRecurringTemplate(ctx, template); err != nil {
				return domain.Template{}, err
			}

			return template, nil
		}
	}

	sheet, err := s.sheetManager.GetOrCreateOpenSheet(ctx, command.UserID)
	if err != nil {
		return domain.Template{}, err
	}

	template := domain.Template{
		ID:        uuid.NewString(),
		UserID:    command.UserID,
		Title:     command.Title,
		Amount:    command.Amount,
		Category:  category,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.recurringTemplateRepo.CreateRecurringTemplate(ctx, template); err != nil {
		return domain.Template{}, err
	}

	entry := domain.Entry{
		ID:        uuid.NewString(),
		SheetID:   sheet.ID,
		UserID:    command.UserID,
		Title:     template.Title,
		Amount:    template.Amount,
		Category:  template.Category,
		Date:      time.Now().UTC().Format("2006-01-02"),
		Type:      domain.EntryTypeRecurring,
		CreatedAt: now,
	}
	if err := s.entryRepo.CreateEntry(ctx, entry); err != nil {
		return domain.Template{}, err
	}

	return template, nil
}

func (s *templateService) DeleteRecurring(ctx context.Context, command ports.DeleteTemplateCommand) (bool, error) {
	if strings.TrimSpace(command.UserID) == "" {
		return false, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}
	if strings.TrimSpace(command.ID) == "" {
		return false, fmt.Errorf("%w: template id is required", domain.ErrInvalidInput)
	}

	return s.recurringTemplateRepo.DeleteRecurringTemplate(ctx, command.UserID, command.ID)
}
