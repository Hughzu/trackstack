package services

import (
	"context"
	"time"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/heat/application/ports"
	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/heat/domain"
)

type refillService struct {
	repo ports.RefillRepository
}

func NewRefillService(repo ports.RefillRepository) ports.RefillUseCase {
	return &refillService{repo: repo}
}

func (s *refillService) GetRefills(ctx context.Context, userID string, from, to time.Time) ([]domain.Refill, error) {
	return s.repo.GetRefills(ctx, userID, from, to)
}
