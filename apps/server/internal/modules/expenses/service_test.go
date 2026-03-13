package expenses_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/expenses"
)

type mockStore struct {
	settings       *expenses.Settings
	historyLimit   int
	historyOffset  int
	checklistReads int
	recurringReads int
}

// Settings
func (m *mockStore) GetSettings(ctx context.Context, userID string) (expenses.Settings, error) {
	if m.settings != nil {
		return *m.settings, nil
	}
	return expenses.Settings{}, errors.New("no rows")
}
func (m *mockStore) CreateSettings(ctx context.Context, settings expenses.Settings) error { return nil }
func (m *mockStore) UpdateSettings(ctx context.Context, settings expenses.Settings) error { return nil }

// Checklist Templates
func (m *mockStore) GetChecklistTemplates(ctx context.Context, userID string) ([]expenses.Template, error) {
	m.checklistReads++
	return []expenses.Template{}, nil
}
func (m *mockStore) GetChecklistTemplate(ctx context.Context, id string, userID string) (expenses.Template, error) {
	return expenses.Template{}, errors.New("not found")
}
func (m *mockStore) CreateChecklistTemplate(ctx context.Context, template expenses.Template) error {
	return nil
}
func (m *mockStore) UpdateChecklistTemplate(ctx context.Context, template expenses.Template) error {
	return nil
}
func (m *mockStore) DeleteChecklistTemplate(ctx context.Context, id string, userID string) (bool, error) {
	return true, nil
}

// Checklist Items
func (m *mockStore) GetChecklistItem(ctx context.Context, id string, userID string) (expenses.ChecklistItem, error) {
	return expenses.ChecklistItem{}, errors.New("not found")
}
func (m *mockStore) CreateChecklistItem(ctx context.Context, item expenses.ChecklistItem) error {
	return nil
}
func (m *mockStore) UpdateChecklistItem(ctx context.Context, item expenses.ChecklistItem) error {
	return nil
}
func (m *mockStore) UpdateChecklistItemsByTemplate(ctx context.Context, templateID string, title string, amount float64, category expenses.Category) error {
	return nil
}
func (m *mockStore) DeletePendingChecklistItemsByTemplate(ctx context.Context, templateID string, userID string) error {
	return nil
}
func (m *mockStore) GetPendingChecklistItems(ctx context.Context, sheetID string) ([]expenses.ChecklistItem, error) {
	return nil, nil
}

// Recurring Templates
func (m *mockStore) GetRecurringTemplates(ctx context.Context, userID string) ([]expenses.Template, error) {
	m.recurringReads++
	return []expenses.Template{}, nil
}
func (m *mockStore) GetRecurringTemplate(ctx context.Context, id string, userID string) (expenses.Template, error) {
	return expenses.Template{}, errors.New("not found")
}
func (m *mockStore) CreateRecurringTemplate(ctx context.Context, template expenses.Template) error {
	return nil
}
func (m *mockStore) UpdateRecurringTemplate(ctx context.Context, template expenses.Template) error {
	return nil
}
func (m *mockStore) DeleteRecurringTemplate(ctx context.Context, id string, userID string) (bool, error) {
	return true, nil
}

// Sheets
func (m *mockStore) GetLatestSheet(ctx context.Context, userID string) (*expenses.Sheet, error) {
	return &expenses.Sheet{ID: "sheet-1", UserID: userID, PeriodKey: "2024-01"}, nil
}
func (m *mockStore) GetOpenSheet(ctx context.Context, userID string) (*expenses.Sheet, error) {
	return &expenses.Sheet{ID: "sheet-1", UserID: userID, PeriodKey: "2024-01"}, nil
}
func (m *mockStore) CreateSheet(ctx context.Context, sheet expenses.Sheet) error { return nil }
func (m *mockStore) UpdateSheet(ctx context.Context, sheet expenses.Sheet) error { return nil }

// Dashboard stats
func (m *mockStore) GetTotalSpentBySheet(ctx context.Context, sheetID string) (float64, error) {
	return 0, nil
}
func (m *mockStore) GetSpentByCategory(ctx context.Context, sheetID string) (map[expenses.Category]float64, error) {
	return nil, nil
}
func (m *mockStore) GetSheetHistory(ctx context.Context, sheetID string) ([]expenses.Entry, error) {
	return nil, nil
}
func (m *mockStore) GetRecentSheetHistory(ctx context.Context, sheetID string, limit int, offset int) ([]expenses.Entry, error) {
	m.historyLimit = limit
	m.historyOffset = offset
	return []expenses.Entry{}, nil
}

// Expenses
func (m *mockStore) CreateExpense(ctx context.Context, entry expenses.Entry) error {
	return nil
}
func (m *mockStore) DeleteExpense(ctx context.Context, id string, userID string) (bool, error) {
	return true, nil
}

func TestAddExpense(t *testing.T) {
	store := &mockStore{}
	svc := expenses.NewService(store)

	ctx := context.Background()
	req := expenses.AddExpenseRequest{
		UserID: "user-1",
		Title:  "Groceries",
		Amount: 50.5,
	}

	entry, err := svc.AddExpense(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if entry.Title != "Groceries" {
		t.Errorf("expected Title to be Groceries, got %s", entry.Title)
	}
	if entry.Amount != 50.5 {
		t.Errorf("expected Amount to be 50.5, got %v", entry.Amount)
	}
}

func TestAddExpenseValidation(t *testing.T) {
	store := &mockStore{}
	svc := expenses.NewService(store)

	ctx := context.Background()

	// Missing userID
	_, err := svc.AddExpense(ctx, expenses.AddExpenseRequest{
		UserID: "",
		Title:  "Groceries",
	})
	if !errors.Is(err, expenses.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput for missing UserID, got %v", err)
	}

	// Missing title becomes "Untitled" in service
	entry, err := svc.AddExpense(ctx, expenses.AddExpenseRequest{
		UserID: "user-1",
		Title:  "  ",
		Amount: 10,
	})
	if err != nil {
		t.Fatalf("expected no error for empty title as it sets Untitled, got %v", err)
	}
	if entry.Title != "Untitled" {
		t.Errorf("expected title to fallback to Untitled, got %s", entry.Title)
	}
}

func TestGetDashboardUsesBoundedHistory(t *testing.T) {
	store := &mockStore{settings: &expenses.Settings{Income: 2000, RatioFund: 60, RatioFun: 20, RatioFuture: 20}}
	svc := expenses.NewService(store)

	_, err := svc.GetDashboard(context.Background(), expenses.GetCurrentSheetRequest{
		UserID:       "user-1",
		HistoryLimit: 25,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if store.historyLimit != 25 {
		t.Fatalf("expected bounded history limit 25, got %d", store.historyLimit)
	}
	if store.historyOffset != 0 {
		t.Fatalf("expected bounded history offset 0, got %d", store.historyOffset)
	}
	if store.checklistReads != 0 {
		t.Fatalf("expected dashboard to avoid checklist template reads, got %d", store.checklistReads)
	}
	if store.recurringReads != 0 {
		t.Fatalf("expected dashboard to avoid recurring template reads, got %d", store.recurringReads)
	}
}
