package domain

type ChecklistItem struct {
	ID          string   `json:"id"`
	SheetID     string   `json:"sheetId"`
	TemplateID  *string  `json:"templateId,omitempty"`
	Title       string   `json:"title"`
	Amount      float64  `json:"amount"`
	Category    Category `json:"category"`
	CreatedAt   string   `json:"createdAt"`
	CompletedAt *string  `json:"completedAt,omitempty"`
	ExpenseID   *string  `json:"expenseId,omitempty"`
}
