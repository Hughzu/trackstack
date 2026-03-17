package services

import (
	"context"
	"strings"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/application/ports"
	"github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/domain"
	"github.com/google/uuid"
)

type CreateRefillService struct {
	refills ports.RefillCreator
}

func NewCreateRefillService(refills ports.RefillCreator) CreateRefillService {
	return CreateRefillService{refills: refills}
}

func (s CreateRefillService) Execute(ctx context.Context, req CreateRefillRequest) (domain.Refill, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return domain.Refill{}, domain.ErrInvalidInput
	}
	if strings.TrimSpace(req.Date) == "" || req.WeightKg <= 0 || req.Bags <= 0 {
		return domain.Refill{}, domain.ErrInvalidInput
	}

	refillDate, err := domain.ParseDate(req.Date)
	if err != nil {
		return domain.Refill{}, err
	}

	seasonLabel := domain.SeasonLabelFor(refillDate)
	refill := domain.Refill{
		ID:          uuid.NewString(),
		UserID:      req.UserID,
		Date:        refillDate.UTC().Format(time.RFC3339),
		WeightKg:    req.WeightKg,
		Bags:        req.Bags,
		Temperature: req.Temperature,
		Season:      &seasonLabel,
	}

	if err := s.refills.Create(ctx, refill); err != nil {
		return domain.Refill{}, err
	}

	return refill, nil
}
