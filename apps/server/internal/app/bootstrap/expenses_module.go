package bootstrap

import (
	"database/sql"

	expenseshttp "github.com/Hughzu/trackstack/apps/server/internal/contexts/expenses/adapters/inbound/http"
	expensesdb "github.com/Hughzu/trackstack/apps/server/internal/contexts/expenses/adapters/outbound/db"
	expensesservice "github.com/Hughzu/trackstack/apps/server/internal/contexts/expenses/application/services"
)

func buildExpensesModule(db *sql.DB) *expenseshttp.ExpensesHandler {
	repository := expensesdb.NewRepository(db)
	sheetManager := expensesservice.NewSheetManager(expensesservice.SheetManagerDeps{
		SheetRepo:             repository,
		ChecklistTemplateRepo: repository,
		RecurringTemplateRepo: repository,
		ChecklistItemRepo:     repository,
		EntryRepo:             repository,
	})

	settingsUseCase := expensesservice.NewSettingsService(expensesservice.SettingsServiceDeps{
		SettingsRepo:          repository,
		ChecklistTemplateRepo: repository,
		RecurringTemplateRepo: repository,
	})
	entryUseCase := expensesservice.NewEntryService(expensesservice.EntryServiceDeps{
		EntryRepo:    repository,
		SheetManager: sheetManager,
	})
	templateUseCase := expensesservice.NewTemplateService(expensesservice.TemplateServiceDeps{
		ChecklistTemplateRepo: repository,
		RecurringTemplateRepo: repository,
		ChecklistItemRepo:     repository,
		EntryRepo:             repository,
		SheetManager:          sheetManager,
	})
	sheetUseCase := expensesservice.NewSheetService(sheetManager)
	dashboardUseCase := expensesservice.NewDashboardService(expensesservice.DashboardServiceDeps{
		SettingsRepo:  repository,
		DashboardRepo: repository,
		SheetManager:  sheetManager,
	})

	handler := expenseshttp.NewExpensesHandler(
		settingsUseCase,
		entryUseCase,
		templateUseCase,
		sheetUseCase,
		dashboardUseCase,
	)

	return handler
}
