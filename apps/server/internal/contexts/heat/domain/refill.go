package domain

type Refill struct {
	ID          string   `json:"id"`
	UserID      string   `json:"userId"`
	Date        string   `json:"date"`
	WeightKg    float64  `json:"weightKg"`
	Bags        int      `json:"bags"`
	Temperature *float64 `json:"temperature,omitempty"`
	Season      *string  `json:"season,omitempty"`
}
