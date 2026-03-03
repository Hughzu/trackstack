package calories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (s *Service) AddLog(ctx context.Context, req AddLogRequest) (Log, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return Log{}, ErrInvalidInput
	}
	if req.Calories == nil || req.ProteinG == nil {
		return Log{}, ErrInvalidInput
	}

	date := ""
	if req.Date != nil {
		date = strings.TrimSpace(*req.Date)
	}
	dateTime, err := buildDateTimeISO(date, req.Time)
	if err != nil {
		return Log{}, err
	}

	var title *string
	if req.Title != nil {
		trimmed := strings.TrimSpace(*req.Title)
		if trimmed != "" {
			title = &trimmed
		}
	}

	log := Log{
		ID:       uuid.NewString(),
		UserID:   req.UserID,
		DateTime: dateTime,
		Calories: *req.Calories,
		ProteinG: *req.ProteinG,
		CarbsG:   req.CarbsG,
		FatG:     req.FatG,
		Title:    title,
	}

	if err := s.store.CreateLog(ctx, log); err != nil {
		return Log{}, err
	}

	return log, nil
}

func (s *Service) DeleteLog(ctx context.Context, req DeleteLogRequest) (bool, error) {
	if strings.TrimSpace(req.UserID) == "" || strings.TrimSpace(req.ID) == "" {
		return false, ErrInvalidInput
	}
	return s.store.DeleteLog(ctx, req.UserID, req.ID)
}

func buildDateTimeISO(date string, timeValue *string) (string, error) {
	if strings.TrimSpace(date) == "" {
		date = time.Now().UTC().Format("2006-01-02")
	}

	clock := ""
	if timeValue != nil {
		clock = strings.TrimSpace(*timeValue)
	}
	if clock == "" {
		clock = time.Now().Format("15:04")
	}

	parsed, err := time.ParseInLocation("2006-01-02T15:04", fmt.Sprintf("%sT%s", date, clock), time.Local)
	if err != nil {
		return "", ErrInvalidInput
	}

	return parsed.UTC().Format(time.RFC3339), nil
}
