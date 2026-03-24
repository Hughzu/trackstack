package domain

type Template struct {
	ID        string   `json:"id"`
	UserID    string   `json:"userId"`
	Title     string   `json:"title"`
	Amount    float64  `json:"amount"`
	Category  Category `json:"category"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
}
