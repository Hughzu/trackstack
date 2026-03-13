package expenses

type Category string

const (
	CategoryFund   Category = "fund"
	CategoryFun    Category = "fun"
	CategoryFuture Category = "future"
)

type EntryType string

const (
	EntryTypeManual    EntryType = "manual"
	EntryTypeRecurring EntryType = "recurring"
	EntryTypeChecklist EntryType = "checklist"
)

type Settings struct {
	ID          string  `json:"id"`
	UserID      string  `json:"userId"`
	Income      float64 `json:"income"`
	RatioFund   int     `json:"ratioFund"`
	RatioFun    int     `json:"ratioFun"`
	RatioFuture int     `json:"ratioFuture"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

type Template struct {
	ID        string   `json:"id"`
	UserID    string   `json:"userId"`
	Title     string   `json:"title"`
	Amount    float64  `json:"amount"`
	Category  Category `json:"category"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
}

type Sheet struct {
	ID        string  `json:"id"`
	UserID    string  `json:"userId"`
	PeriodKey string  `json:"periodKey"`
	CreatedAt string  `json:"createdAt"`
	ClosedAt  *string `json:"closedAt,omitempty"`
}

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

type ViewSettings struct {
	Settings  Settings   `json:"settings"`
	Checklist []Template `json:"checklist"`
	Recurring []Template `json:"recurring"`
}

type DashboardBalance struct {
	Remaining float64 `json:"remaining"`
	Income    float64 `json:"income"`
}

type DashboardSpent struct {
	Fund   float64 `json:"fund"`
	Fun    float64 `json:"fun"`
	Future float64 `json:"future"`
}

type DashboardBudget struct {
	Fund   int `json:"fund"`
	Fun    int `json:"fun"`
	Future int `json:"future"`
}

type DashboardRatio struct {
	Percent    int     `json:"percent"`
	CategoryId string  `json:"categoryId"` // e.g. "fund", "fun", "future"
	Label      string  `json:"label"`
	Value      float64 `json:"value"`
	Budget     int     `json:"budget"`
	Target     int     `json:"target"`
	Over       bool    `json:"over"`
}

type ViewDashboard struct {
	PeriodKey          string           `json:"periodKey"`
	Balance            DashboardBalance `json:"balance"`
	Spent              DashboardSpent   `json:"spent"`
	Budget             DashboardBudget  `json:"budget"`
	Ratios             []DashboardRatio `json:"ratios"`
	PendingObligations []ChecklistItem  `json:"pendingObligations"`
	History            []Entry          `json:"history"`
}

type GetSettingsRequest struct {
	UserID string
}

type UpdateSettingsRequest struct {
	UserID      string
	Income      *float64
	RatioFund   *int
	RatioFun    *int
	RatioFuture *int
}

type UpsertTemplateRequest struct {
	ID       *string
	UserID   string
	Title    string
	Amount   float64
	Category *string
}

type CompleteChecklistItemRequest struct {
	ID     string
	UserID string
	Date   *string
}

type AddExpenseRequest struct {
	UserID   string
	Title    string
	Amount   float64
	Category *string
	Date     *string
}

type DeleteTemplateRequest struct {
	ID     string
	UserID string
}

type DeleteExpenseRequest struct {
	ID     string
	UserID string
}

type CloseSheetRequest struct {
	UserID string
}

type GetCurrentSheetRequest struct {
	UserID       string
	HistoryLimit int
}
