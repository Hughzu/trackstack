package domain

type Entry struct {
	ID        string    `json:"id"`
	SheetID   string    `json:"sheetId"`
	UserID    string    `json:"userId"`
	Title     string    `json:"title"`
	Amount    float64   `json:"amount"`
	Category  Category  `json:"category"`
	Date      string    `json:"date"`
	Type      EntryType `json:"type"`
	CreatedAt string    `json:"createdAt"`
}
