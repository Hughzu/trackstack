package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/expenses/application/ports"
	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/expenses/domain"
	"github.com/google/uuid"
)

type entryService struct {
	entryRepo    ports.EntryRepository
	sheetManager SheetManager
}

type EntryServiceDeps struct {
	EntryRepo    ports.EntryRepository
	SheetManager SheetManager
}

func NewEntryService(deps EntryServiceDeps) ports.EntryUseCase {
	return &entryService{
		entryRepo:    deps.EntryRepo,
		sheetManager: deps.SheetManager,
	}
}

func (s *entryService) AddEntry(ctx context.Context, command ports.AddEntryCommand) (domain.Entry, error) {
	if strings.TrimSpace(command.UserID) == "" {
		return domain.Entry{}, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}

	title := strings.TrimSpace(command.Title)
	if title == "" {
		title = "Untitled"
	}

	sheet, err := s.sheetManager.GetOrCreateOpenSheet(ctx, command.UserID)
	if err != nil {
		return domain.Entry{}, err
	}

	now := time.Now().UTC()
	date := now.Format("2006-01-02")
	if command.Date != nil && *command.Date != "" {
		date = *command.Date
	}

	entry := domain.Entry{
		ID:        uuid.NewString(),
		SheetID:   sheet.ID,
		UserID:    command.UserID,
		Title:     title,
		Amount:    command.Amount,
		Category:  resolveCategory(command.Category),
		Date:      date,
		Type:      domain.EntryTypeManual,
		CreatedAt: now.Format(time.RFC3339),
	}

	if err := s.entryRepo.CreateEntry(ctx, entry); err != nil {
		return domain.Entry{}, err
	}

	return entry, nil
}

func (s *entryService) DeleteEntry(ctx context.Context, command ports.DeleteEntryCommand) (bool, error) {
	if strings.TrimSpace(command.UserID) == "" {
		return false, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}
	if strings.TrimSpace(command.ID) == "" {
		return false, fmt.Errorf("%w: entry id is required", domain.ErrInvalidInput)
	}

	return s.entryRepo.DeleteEntry(ctx, command.UserID, command.ID)
}

func resolveCategory(category *string) domain.Category {
	if category == nil {
		return domain.CategoryFund
	}

	return domain.Category(*category)
}
