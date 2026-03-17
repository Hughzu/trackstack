package services

import (
	"context"
	"strings"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/application/ports"
	"github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/domain"
)

type ListRefillsService struct {
	refills ports.RefillRangeLister
}

func NewListRefillsService(refills ports.RefillRangeLister) ListRefillsService {
	return ListRefillsService{refills: refills}
}

func (s ListRefillsService) Execute(ctx context.Context, req ListRefillsRequest) ([]domain.Refill, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return nil, domain.ErrInvalidInput
	}

	from, to, err := domain.NormalizeRange(req.From, req.To)
	if err != nil {
		return nil, err
	}

	refills, err := s.refills.ListByRange(ctx, req.UserID, from, to)
	if err != nil {
		return nil, err
	}

	return refills, nil
}
