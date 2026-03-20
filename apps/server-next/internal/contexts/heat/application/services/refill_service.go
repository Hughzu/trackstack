package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/heat/application/ports"
	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/heat/domain"
	"github.com/google/uuid"
)

type refillService struct {
	repo ports.RefillRepository
}

func NewRefillService(repo ports.RefillRepository) ports.RefillUseCase {
	return &refillService{repo: repo}
}

func (s *refillService) GetRefills(ctx context.Context, query ports.GetRefillsQuery) ([]domain.Refill, error) {
	if strings.TrimSpace(query.UserID) == "" {
		return nil, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}

	return s.repo.GetRefills(ctx, query.UserID, query.From, query.To)
}

func (s *refillService) CreateRefill(ctx context.Context, command ports.CreateRefillCommand) (domain.Refill, error) {
	if strings.TrimSpace(command.UserID) == "" {
		return domain.Refill{}, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}
	if command.WeightKg <= 0 {
		return domain.Refill{}, fmt.Errorf("%w: weight must be greater than zero", domain.ErrInvalidInput)
	}
	if command.Bags <= 0 {
		return domain.Refill{}, fmt.Errorf("%w: bags must be greater than zero", domain.ErrInvalidInput)
	}
	if command.Date.IsZero() {
		return domain.Refill{}, fmt.Errorf("%w: date is required", domain.ErrInvalidInput)
	}

	refillDate := command.Date.UTC()
	seasonLabel := domain.SeasonLabelFor(refillDate)
	refill := domain.Refill{
		ID:          uuid.NewString(),
		UserID:      command.UserID,
		Date:        refillDate.UTC().Format(time.RFC3339),
		WeightKg:    command.WeightKg,
		Bags:        command.Bags,
		Temperature: command.Temperature,
		Season:      &seasonLabel,
	}

	if err := s.repo.CreateRefill(ctx, refill); err != nil {
		return domain.Refill{}, err
	}

	return refill, nil
}

func (s *refillService) DeleteRefill(ctx context.Context, command ports.DeleteRefillCommand) (bool, error) {
	if strings.TrimSpace(command.UserID) == "" {
		return false, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}
	if strings.TrimSpace(command.ID) == "" {
		return false, fmt.Errorf("%w: refill id is required", domain.ErrInvalidInput)
	}

	return s.repo.DeleteRefill(ctx, command.UserID, command.ID)
}
