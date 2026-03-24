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

const (
	defaultIncome      = 2215
	defaultRatioFund   = 60
	defaultRatioFun    = 20
	defaultRatioFuture = 20
)

type settingsService struct {
	settingsRepo          ports.SettingsRepository
	checklistTemplateRepo ports.ChecklistTemplateRepository
	recurringTemplateRepo ports.RecurringTemplateRepository
}

type SettingsServiceDeps struct {
	SettingsRepo          ports.SettingsRepository
	ChecklistTemplateRepo ports.ChecklistTemplateRepository
	RecurringTemplateRepo ports.RecurringTemplateRepository
}

func NewSettingsService(deps SettingsServiceDeps) ports.SettingsUseCase {
	return &settingsService{
		settingsRepo:          deps.SettingsRepo,
		checklistTemplateRepo: deps.ChecklistTemplateRepo,
		recurringTemplateRepo: deps.RecurringTemplateRepo,
	}
}

func (s *settingsService) GetSettings(ctx context.Context, query ports.GetSettingsQuery) (domain.SettingsView, error) {
	if strings.TrimSpace(query.UserID) == "" {
		return domain.SettingsView{}, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}

	settings, err := getOrCreateSettings(ctx, s.settingsRepo, query.UserID)
	if err != nil {
		return domain.SettingsView{}, err
	}

	checklist, err := s.checklistTemplateRepo.ListChecklistTemplates(ctx, query.UserID)
	if err != nil {
		return domain.SettingsView{}, err
	}
	if checklist == nil {
		checklist = []domain.Template{}
	}

	recurring, err := s.recurringTemplateRepo.ListRecurringTemplates(ctx, query.UserID)
	if err != nil {
		return domain.SettingsView{}, err
	}
	if recurring == nil {
		recurring = []domain.Template{}
	}

	return domain.SettingsView{
		Settings:  settings,
		Checklist: checklist,
		Recurring: recurring,
	}, nil
}

func (s *settingsService) UpdateSettings(ctx context.Context, command ports.UpdateSettingsCommand) (domain.Settings, error) {
	if strings.TrimSpace(command.UserID) == "" {
		return domain.Settings{}, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}

	settings, found, err := s.settingsRepo.FindSettings(ctx, command.UserID)
	if err != nil {
		return domain.Settings{}, err
	}
	if !found {
		return domain.Settings{}, fmt.Errorf("%w: settings not found", domain.ErrNotFound)
	}

	settings.Income = command.Income
	settings.RatioFund = command.RatioFund
	settings.RatioFun = command.RatioFun
	settings.RatioFuture = command.RatioFuture
	settings.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := s.settingsRepo.UpdateSettings(ctx, settings); err != nil {
		return domain.Settings{}, err
	}

	return settings, nil
}

func getOrCreateSettings(ctx context.Context, repo ports.SettingsRepository, userID string) (domain.Settings, error) {
	settings, found, err := repo.FindSettings(ctx, userID)
	if err != nil {
		return domain.Settings{}, err
	}
	if found {
		return settings, nil
	}

	now := time.Now().UTC().Format(time.RFC3339)
	settings = domain.Settings{
		ID:          uuid.NewString(),
		UserID:      userID,
		Income:      defaultIncome,
		RatioFund:   defaultRatioFund,
		RatioFun:    defaultRatioFun,
		RatioFuture: defaultRatioFuture,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := repo.CreateSettings(ctx, settings); err != nil {
		return domain.Settings{}, err
	}

	return settings, nil
}
