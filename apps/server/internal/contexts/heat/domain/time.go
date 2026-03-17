package domain

import (
	"fmt"
	"strings"
	"time"
)

func NormalizeRange(from string, to string) (string, string, error) {
	from = strings.TrimSpace(from)
	to = strings.TrimSpace(to)

	var fromValue string
	var toValue string
	if from != "" {
		parsed, err := ParseDate(from)
		if err != nil {
			return "", "", err
		}
		fromValue = parsed.UTC().Format(time.RFC3339)
	}
	if to != "" {
		parsed, err := ParseDate(to)
		if err != nil {
			return "", "", err
		}
		if IsDateOnly(to) {
			parsed = parsed.Add(24 * time.Hour)
		}
		toValue = parsed.UTC().Format(time.RFC3339)
	}

	return fromValue, toValue, nil
}

func ParseDate(value string) (time.Time, error) {
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

func IsDateOnly(value string) bool {
	return !strings.Contains(value, "T")
}

func SeasonLabelFor(date time.Time) string {
	startYear := date.Year()
	if date.Month() < time.September {
		startYear--
	}
	return fmt.Sprintf("%d-%d", startYear, startYear+1)
}
