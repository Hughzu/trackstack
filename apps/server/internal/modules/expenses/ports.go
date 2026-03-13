package expenses

import "context"

type ExpensesStore interface {
	// Settings
	GetSettings(ctx context.Context, userID string) (Settings, error)
	CreateSettings(ctx context.Context, settings Settings) error
	UpdateSettings(ctx context.Context, settings Settings) error

	// Checklist Templates
	GetChecklistTemplates(ctx context.Context, userID string) ([]Template, error)
	GetChecklistTemplate(ctx context.Context, templateID string, userID string) (Template, error)
	CreateChecklistTemplate(ctx context.Context, template Template) error
	UpdateChecklistTemplate(ctx context.Context, template Template) error
	DeleteChecklistTemplate(ctx context.Context, templateID string, userID string) (bool, error)

	// Recurring Templates
	GetRecurringTemplates(ctx context.Context, userID string) ([]Template, error)
	GetRecurringTemplate(ctx context.Context, templateID string, userID string) (Template, error)
	CreateRecurringTemplate(ctx context.Context, template Template) error
	UpdateRecurringTemplate(ctx context.Context, template Template) error
	DeleteRecurringTemplate(ctx context.Context, templateID string, userID string) (bool, error)

	// Sheets
	GetOpenSheet(ctx context.Context, userID string) (*Sheet, error)
	GetLatestSheet(ctx context.Context, userID string) (*Sheet, error)
	CreateSheet(ctx context.Context, sheet Sheet) error
	UpdateSheet(ctx context.Context, sheet Sheet) error

	// Checklist Items
	GetChecklistItem(ctx context.Context, itemID string, userID string) (ChecklistItem, error)
	GetPendingChecklistItems(ctx context.Context, sheetID string) ([]ChecklistItem, error)
	CreateChecklistItem(ctx context.Context, item ChecklistItem) error
	UpdateChecklistItem(ctx context.Context, item ChecklistItem) error
	UpdateChecklistItemsByTemplate(ctx context.Context, templateID string, title string, amount float64, category Category) error
	DeletePendingChecklistItemsByTemplate(ctx context.Context, templateID string, userID string) error

	// Entries
	CreateExpense(ctx context.Context, entry Entry) error
	DeleteExpense(ctx context.Context, entryID string, userID string) (bool, error)
	GetSheetHistory(ctx context.Context, sheetID string) ([]Entry, error)
	GetRecentSheetHistory(ctx context.Context, sheetID string, limit int, offset int) ([]Entry, error)

	// Aggregation
	GetTotalSpentBySheet(ctx context.Context, sheetID string) (float64, error)
	GetSpentByCategory(ctx context.Context, sheetID string) (map[Category]float64, error)
}
