package builddatetime

import (
	"errors"
	"time"
)

var ErrInvalidDateTime = errors.New("invalid date or time")

func BuildRFC3339DateTime(dateValue string, timeValue string, now time.Time) (string, error) {
	panic("TODO")
}
