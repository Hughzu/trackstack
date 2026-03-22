package timeutil

import (
	"fmt"
	"strings"
	"time"
)

func ParseDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, fmt.Errorf("date is required")
	}

	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return parsed.UTC(), nil
	}

	if parsed, err := time.Parse("2006-01-02", value); err == nil {
		return time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC), nil
	}

	return time.Time{}, fmt.Errorf("invalid date")
}

func NormalizeDateString(value string) string {
	parsed, err := ParseDate(value)
	if err != nil {
		return value
	}

	return parsed.UTC().Format(time.RFC3339)
}

func BuildRFC3339DateTime(dateValue *string, timeValue *string) (string, error) {
	date := ""
	if dateValue != nil {
		date = strings.TrimSpace(*dateValue)
	}
	if date == "" {
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
		return "", fmt.Errorf("invalid date or time")
	}

	return parsed.UTC().Format(time.RFC3339), nil
}
