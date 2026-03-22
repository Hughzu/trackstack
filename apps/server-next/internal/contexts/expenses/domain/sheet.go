package domain

type Sheet struct {
	ID        string  `json:"id"`
	UserID    string  `json:"userId"`
	PeriodKey string  `json:"periodKey"`
	CreatedAt string  `json:"createdAt"`
	ClosedAt  *string `json:"closedAt,omitempty"`
}
