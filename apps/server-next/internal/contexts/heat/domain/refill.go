package domain

import "time"

type Refill struct {
	ID        string
	Amount    int
	CreatedAt time.Time
}
