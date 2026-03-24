package ports

import (
	"context"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/expenses/domain"
)

type SettingsRepository interface {
	FindSettings(ctx context.Context, userID string) (domain.Settings, bool, error)
	CreateSettings(ctx context.Context, settings domain.Settings) error
	UpdateSettings(ctx context.Context, settings domain.Settings) error
}

type SheetRepository interface {
	FindOpenSheet(ctx context.Context, userID string) (*domain.Sheet, error)
	FindLatestSheet(ctx context.Context, userID string) (*domain.Sheet, error)
	CreateSheet(ctx context.Context, sheet domain.Sheet) error
	UpdateSheet(ctx context.Context, sheet domain.Sheet) error
}

type EntryRepository interface {
	CreateEntry(ctx context.Context, entry domain.Entry) error
	DeleteEntry(ctx context.Context, userID string, id string) (bool, error)
}

type ChecklistTemplateRepository interface {
	FindChecklistTemplate(ctx context.Context, userID string, id string) (domain.Template, bool, error)
	ListChecklistTemplates(ctx context.Context, userID string) ([]domain.Template, error)
	CreateChecklistTemplate(ctx context.Context, template domain.Template) error
	UpdateChecklistTemplate(ctx context.Context, template domain.Template) error
	DeleteChecklistTemplate(ctx context.Context, userID string, id string) (bool, error)
}

type RecurringTemplateRepository interface {
	FindRecurringTemplate(ctx context.Context, userID string, id string) (domain.Template, bool, error)
	ListRecurringTemplates(ctx context.Context, userID string) ([]domain.Template, error)
	CreateRecurringTemplate(ctx context.Context, template domain.Template) error
	UpdateRecurringTemplate(ctx context.Context, template domain.Template) error
	DeleteRecurringTemplate(ctx context.Context, userID string, id string) (bool, error)
}

type ChecklistItemRepository interface {
	FindChecklistItem(ctx context.Context, userID string, id string) (domain.ChecklistItem, bool, error)
	ListPendingChecklistItemsBySheet(ctx context.Context, sheetID string) ([]domain.ChecklistItem, error)
	CreateChecklistItem(ctx context.Context, item domain.ChecklistItem) error
	UpdateChecklistItem(ctx context.Context, item domain.ChecklistItem) error
	UpdatePendingChecklistItemsByTemplate(ctx context.Context, templateID string, title string, amount float64, category domain.Category) error
	DeletePendingChecklistItemsByTemplate(ctx context.Context, userID string, templateID string) error
}

type DashboardSnapshot struct {
	TotalSpent            float64
	SpentByCategory       map[domain.Category]float64
	PendingChecklistItems []domain.ChecklistItem
	History               []domain.Entry
}

type DashboardRepository interface {
	GetDashboardSnapshot(ctx context.Context, sheetID string, historyLimit int) (DashboardSnapshot, error)
}
