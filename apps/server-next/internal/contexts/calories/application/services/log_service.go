package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/calories/application/ports"
	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/calories/domain"
	"github.com/Hughzu/trackstack/apps/server-next/internal/platform/timeutil"
	"github.com/google/uuid"
)

type logService struct {
	repo ports.LogRepository
}

func NewLogService(repo ports.LogRepository) ports.LogUseCase {
	return &logService{repo: repo}
}

func (s *logService) AddLog(ctx context.Context, command ports.AddLogCommand) (domain.Log, error) {
	if strings.TrimSpace(command.UserID) == "" {
		return domain.Log{}, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}
	if command.Calories <= 0 {
		return domain.Log{}, fmt.Errorf("%w: calories must be greater than zero", domain.ErrInvalidInput)
	}
	if command.ProteinGrams <= 0 {
		return domain.Log{}, fmt.Errorf("%w: protein grams must be greater than zero", domain.ErrInvalidInput)
	}

	dateTime, err := timeutil.BuildRFC3339DateTime(command.Date, command.Time)
	if err != nil {
		return domain.Log{}, fmt.Errorf("%w: %v", domain.ErrInvalidInput, err)
	}

	log := domain.Log{
		ID:           uuid.NewString(),
		UserID:       command.UserID,
		DateTime:     dateTime,
		Calories:     command.Calories,
		ProteinGrams: command.ProteinGrams,
		CarbGrams:    command.CarbGrams,
		FatGrams:     command.FatGrams,
		Title:        normalizeOptionalString(command.Title),
	}

	if err := s.repo.CreateLog(ctx, log); err != nil {
		return domain.Log{}, err
	}

	return log, nil
}

func (s *logService) DeleteLog(ctx context.Context, command ports.DeleteLogCommand) (bool, error) {
	if strings.TrimSpace(command.UserID) == "" {
		return false, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}
	if strings.TrimSpace(command.ID) == "" {
		return false, fmt.Errorf("%w: log id is required", domain.ErrInvalidInput)
	}

	return s.repo.DeleteLog(ctx, command.UserID, command.ID)
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}
