package domain

import (
	"fmt"
	"time"
)

func SeasonLabelFor(date time.Time) string {
	startYear := date.Year()
	if date.Month() < time.September {
		startYear--
	}

	return fmt.Sprintf("%d-%d", startYear, startYear+1)
}
