package heat

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrInvalidInput = errors.New("invalid input")

type Service struct {
	store RefillStore
}

func NewService(store RefillStore) *Service {
	return &Service{store: store}
}

func (s *Service) ListRefills(ctx context.Context, req ListRefillsRequest) ([]Refill, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return nil, ErrInvalidInput
	}

	from, to, err := normalizeRange(req.From, req.To)
	if err != nil {
		return nil, err
	}

	refills, err := s.store.ListByRange(ctx, req.UserID, from, to)
	if err != nil {
		return nil, err
	}

	return refills, nil
}

func (s *Service) CreateRefill(ctx context.Context, req CreateRefillRequest) (Refill, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return Refill{}, ErrInvalidInput
	}
	if strings.TrimSpace(req.Date) == "" || req.WeightKg <= 0 || req.Bags <= 0 {
		return Refill{}, ErrInvalidInput
	}

	refillDate, err := parseDate(req.Date)
	if err != nil {
		return Refill{}, err
	}

	seasonLabel := seasonLabelFor(refillDate)
	input := CreateRefillInput{
		Date:        refillDate.UTC().Format(time.RFC3339),
		WeightKg:    req.WeightKg,
		Bags:        req.Bags,
		Temperature: req.Temperature,
		Season:      &seasonLabel,
	}

	refill, err := s.store.Create(ctx, req.UserID, input)
	if err != nil {
		return Refill{}, err
	}

	return refill, nil
}

func (s *Service) DeleteRefill(ctx context.Context, req DeleteRefillRequest) (bool, error) {
	if strings.TrimSpace(req.UserID) == "" || strings.TrimSpace(req.ID) == "" {
		return false, ErrInvalidInput
	}

	return s.store.Delete(ctx, req.UserID, req.ID)
}

func normalizeRange(from string, to string) (string, string, error) {
	from = strings.TrimSpace(from)
	to = strings.TrimSpace(to)

	var fromValue string
	var toValue string
	if from != "" {
		parsed, err := parseDate(from)
		if err != nil {
			return "", "", err
		}
		fromValue = parsed.UTC().Format(time.RFC3339)
	}
	if to != "" {
		parsed, err := parseDate(to)
		if err != nil {
			return "", "", err
		}
		if isDateOnly(to) {
			parsed = parsed.Add(24 * time.Hour)
		}
		toValue = parsed.UTC().Format(time.RFC3339)
	}

	return fromValue, toValue, nil
}

func parseDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("%w: date is required", ErrInvalidInput)
	}

	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return parsed, nil
	}

	if parsed, err := time.Parse("2006-01-02", value); err == nil {
		return time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC), nil
	}

	return time.Time{}, fmt.Errorf("%w: invalid date", ErrInvalidInput)
}

func isDateOnly(value string) bool {
	return !strings.Contains(value, "T")
}

func seasonLabelFor(date time.Time) string {
	startYear := date.Year()
	if date.Month() < time.September {
		startYear--
	}
	return fmt.Sprintf("%d-%d", startYear, startYear+1)
}
