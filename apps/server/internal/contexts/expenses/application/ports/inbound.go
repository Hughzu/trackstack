package ports

import (
	"context"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/expenses/domain"
)

var ErrInvalidInput = domain.ErrInvalidInput

type SettingsUseCase interface {
	GetSettings(ctx context.Context, query GetSettingsQuery) (domain.SettingsView, error)
	UpdateSettings(ctx context.Context, command UpdateSettingsCommand) (domain.Settings, error)
}

type EntryUseCase interface {
	AddEntry(ctx context.Context, command AddEntryCommand) (domain.Entry, error)
	DeleteEntry(ctx context.Context, command DeleteEntryCommand) (bool, error)
}

type TemplateUseCase interface {
	UpsertChecklist(ctx context.Context, command UpsertTemplateCommand) (domain.Template, error)
	DeleteChecklist(ctx context.Context, command DeleteTemplateCommand) (bool, error)
	CompleteChecklistItem(ctx context.Context, command CompleteChecklistItemCommand) (domain.Entry, error)
	UpsertRecurring(ctx context.Context, command UpsertTemplateCommand) (domain.Template, error)
	DeleteRecurring(ctx context.Context, command DeleteTemplateCommand) (bool, error)
}

type SheetUseCase interface {
	CloseSheet(ctx context.Context, command CloseSheetCommand) (domain.Sheet, error)
}

type DashboardUseCase interface {
	GetDashboard(ctx context.Context, query GetDashboardQuery) (domain.Dashboard, error)
}

type GetSettingsQuery struct {
	UserID string
}

type UpdateSettingsCommand struct {
	UserID      string
	Income      float64
	RatioFund   int
	RatioFun    int
	RatioFuture int
}

type AddEntryCommand struct {
	UserID   string
	Title    string
	Amount   float64
	Category *string
	Date     *string
}

type DeleteEntryCommand struct {
	UserID string
	ID     string
}

type UpsertTemplateCommand struct {
	ID       *string
	UserID   string
	Title    string
	Amount   float64
	Category *string
}

type DeleteTemplateCommand struct {
	UserID string
	ID     string
}

type CompleteChecklistItemCommand struct {
	UserID string
	ID     string
	Date   *string
}

type CloseSheetCommand struct {
	UserID string
}

type GetDashboardQuery struct {
	UserID       string
	HistoryLimit int
}
