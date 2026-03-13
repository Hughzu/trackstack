package expenses

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (s *Service) getOrCreateSettings(ctx context.Context, userID string) (Settings, error) {
	settings, err := s.store.GetSettings(ctx, userID)
	if err == nil {
		return settings, nil
	}
	if !strings.Contains(err.Error(), "no rows") {
		return Settings{}, err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	settings = Settings{
		ID:          uuid.NewString(),
		UserID:      userID,
		Income:      2215,
		RatioFund:   60,
		RatioFun:    20,
		RatioFuture: 20,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.store.CreateSettings(ctx, settings); err != nil {
		return Settings{}, err
	}

	return settings, nil
}

func (s *Service) GetSettings(ctx context.Context, req GetSettingsRequest) (ViewSettings, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return ViewSettings{}, ErrInvalidInput
	}

	settings, err := s.getOrCreateSettings(ctx, req.UserID)
	if err != nil {
		return ViewSettings{}, err
	}

	checklist, err := s.store.GetChecklistTemplates(ctx, req.UserID)
	if err != nil {
		return ViewSettings{}, err
	}

	recurring, err := s.store.GetRecurringTemplates(ctx, req.UserID)
	if err != nil {
		return ViewSettings{}, err
	}

	return ViewSettings{
		Settings:  settings,
		Checklist: checklist,
		Recurring: recurring,
	}, nil
}

func (s *Service) UpdateSettings(ctx context.Context, req UpdateSettingsRequest) (Settings, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return Settings{}, ErrInvalidInput
	}
	if req.Income == nil || req.RatioFund == nil || req.RatioFun == nil || req.RatioFuture == nil {
		return Settings{}, ErrInvalidInput
	}

	settings, err := s.store.GetSettings(ctx, req.UserID)
	if err != nil {
		return Settings{}, err
	}

	settings.Income = *req.Income
	settings.RatioFund = *req.RatioFund
	settings.RatioFun = *req.RatioFun
	settings.RatioFuture = *req.RatioFuture
	settings.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := s.store.UpdateSettings(ctx, settings); err != nil {
		return Settings{}, err
	}
	return settings, nil
}
