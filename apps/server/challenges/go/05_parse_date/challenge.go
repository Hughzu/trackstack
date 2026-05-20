package parsedate

import (
	"errors"
	"time"
)

var (
	ErrDateRequired = errors.New("date is required")
	ErrInvalidDate  = errors.New("invalid date")
)

func ParseDate(value string) (time.Time, error) {
	panic("TODO")
}
