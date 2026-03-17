package services

import (
	"context"
	"strings"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/application/ports"
	"github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/domain"
)

type DeleteRefillService struct {
	refills ports.RefillDeleter
}

func NewDeleteRefillService(refills ports.RefillDeleter) DeleteRefillService {
	return DeleteRefillService{refills: refills}
}

func (s DeleteRefillService) Execute(ctx context.Context, req DeleteRefillRequest) (bool, error) {
	if strings.TrimSpace(req.UserID) == "" || strings.TrimSpace(req.ID) == "" {
		return false, domain.ErrInvalidInput
	}

	return s.refills.Delete(ctx, req.UserID, req.ID)
}
