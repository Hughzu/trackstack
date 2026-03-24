package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/calories/application/ports"
	"github.com/Hughzu/trackstack/apps/server/internal/contexts/calories/domain"
	"github.com/google/uuid"
)

const (
	defaultTargetCalories     = 2300
	defaultTargetProteinGrams = 120
)

type targetService struct {
	repo ports.TargetRepository
}

func NewTargetService(repo ports.TargetRepository) ports.TargetUseCase {
	return &targetService{repo: repo}
}

func (s *targetService) GetTarget(ctx context.Context, query ports.GetTargetQuery) (domain.Target, error) {
	if strings.TrimSpace(query.UserID) == "" {
		return domain.Target{}, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}

	target, found, err := s.repo.FindTarget(ctx, query.UserID)
	if err != nil {
		return domain.Target{}, err
	}
	if found {
		return target, nil
	}

	return s.createDefaultTarget(ctx, query.UserID)
}

func (s *targetService) UpdateTarget(ctx context.Context, command ports.UpdateTargetCommand) (domain.Target, error) {
	if strings.TrimSpace(command.UserID) == "" {
		return domain.Target{}, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}
	if command.TargetCalories <= 0 {
		return domain.Target{}, fmt.Errorf("%w: target calories must be greater than zero", domain.ErrInvalidInput)
	}
	if command.TargetProteinGrams <= 0 {
		return domain.Target{}, fmt.Errorf("%w: target protein grams must be greater than zero", domain.ErrInvalidInput)
	}

	existing, found, err := s.repo.FindTarget(ctx, command.UserID)
	if err != nil {
		return domain.Target{}, err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	if found {
		existing.TargetCalories = command.TargetCalories
		existing.TargetProteinGrams = command.TargetProteinGrams
		existing.TargetCarbGrams = command.TargetCarbGrams
		existing.TargetFatGrams = command.TargetFatGrams
		existing.UpdatedAt = now

		if err := s.repo.UpdateTarget(ctx, existing); err != nil {
			return domain.Target{}, err
		}

		return existing, nil
	}

	created := domain.Target{
		ID:                 uuid.NewString(),
		UserID:             command.UserID,
		TargetCalories:     command.TargetCalories,
		TargetProteinGrams: command.TargetProteinGrams,
		TargetCarbGrams:    command.TargetCarbGrams,
		TargetFatGrams:     command.TargetFatGrams,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.repo.CreateTarget(ctx, created); err != nil {
		return domain.Target{}, err
	}

	return created, nil
}

func (s *targetService) createDefaultTarget(ctx context.Context, userID string) (domain.Target, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	target := domain.Target{
		ID:                 uuid.NewString(),
		UserID:             userID,
		TargetCalories:     defaultTargetCalories,
		TargetProteinGrams: defaultTargetProteinGrams,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.repo.CreateTarget(ctx, target); err != nil {
		return domain.Target{}, err
	}

	return target, nil
}
