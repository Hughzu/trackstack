package heat

type Refill struct {
	ID          string   `json:"id"`
	UserID      string   `json:"userId"`
	Date        string   `json:"date"`
	WeightKg    float64  `json:"weightKg"`
	Bags        int      `json:"bags"`
	Temperature *float64 `json:"temperature,omitempty"`
	Season      *string  `json:"season,omitempty"`
}

type ListRefillsRequest struct {
	UserID string
	From   string
	To     string
}

type CreateRefillRequest struct {
	UserID      string
	Date        string
	WeightKg    float64
	Bags        int
	Temperature *float64
}

type CreateRefillInput struct {
	Date        string
	WeightKg    float64
	Bags        int
	Temperature *float64
	Season      *string
}

type DeleteRefillRequest struct {
	UserID string
	ID     string
}
