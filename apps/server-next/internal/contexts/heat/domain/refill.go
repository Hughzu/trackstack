package domain

import "time"

type Refill struct {
	ID          string
	WeightKg    float64
	Bags        int
	Temperature *float64
	Season      *string
	Date        time.Time
}
