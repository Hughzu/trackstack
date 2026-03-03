package calories

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	defaultTargetKcal     = 2300
	defaultTargetProteinG = 120
)

func (s *Service) GetTarget(ctx context.Context, req GetTargetRequest) (Target, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return Target{}, ErrInvalidInput
	}

	target, err := s.store.GetTarget(ctx, req.UserID)
	if err == nil {
		return target, nil
	}
	if !isNoRows(err) {
		return Target{}, err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	created := Target{
		ID:             uuid.NewString(),
		UserID:         req.UserID,
		TargetKcal:     defaultTargetKcal,
		TargetProteinG: defaultTargetProteinG,
		TargetCarbsG:   nil,
		TargetFatG:     nil,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.store.CreateTarget(ctx, created); err != nil {
		return Target{}, err
	}

	return created, nil
}

func (s *Service) UpdateTarget(ctx context.Context, req UpdateTargetRequest) (Target, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return Target{}, ErrInvalidInput
	}
	if req.TargetKcal == nil || req.TargetProteinG == nil {
		return Target{}, ErrInvalidInput
	}

	existing, err := s.store.GetTarget(ctx, req.UserID)
	if err != nil && !isNoRows(err) {
		return Target{}, err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	if err == nil {
		existing.TargetKcal = *req.TargetKcal
		existing.TargetProteinG = *req.TargetProteinG
		existing.TargetCarbsG = req.TargetCarbsG
		existing.TargetFatG = req.TargetFatG
		existing.UpdatedAt = now
		if err := s.store.UpdateTarget(ctx, existing); err != nil {
			return Target{}, err
		}
		return existing, nil
	}

	created := Target{
		ID:             uuid.NewString(),
		UserID:         req.UserID,
		TargetKcal:     *req.TargetKcal,
		TargetProteinG: *req.TargetProteinG,
		TargetCarbsG:   req.TargetCarbsG,
		TargetFatG:     req.TargetFatG,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := s.store.CreateTarget(ctx, created); err != nil {
		return Target{}, err
	}

	return created, nil
}

func isNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "no rows")
}
