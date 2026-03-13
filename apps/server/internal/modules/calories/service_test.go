package calories_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/calories"
)

type mockStore struct {
	target        *calories.Target
	getTargetErr  error
	createdTarget *calories.Target
	logsLimit     int
}

func (m *mockStore) GetTarget(ctx context.Context, userID string) (calories.Target, error) {
	if m.getTargetErr != nil {
		return calories.Target{}, m.getTargetErr
	}
	if m.target != nil {
		return *m.target, nil
	}
	return calories.Target{}, sql.ErrNoRows
}

func (m *mockStore) CreateTarget(ctx context.Context, target calories.Target) error {
	m.createdTarget = &target
	return nil
}

func (m *mockStore) UpdateTarget(ctx context.Context, target calories.Target) error {
	m.target = &target
	return nil
}

func (m *mockStore) CreateLog(ctx context.Context, log calories.Log) error {
	return nil
}

func (m *mockStore) DeleteLog(ctx context.Context, userID string, id string) (bool, error) {
	return true, nil
}

func (m *mockStore) GetSummaryByRange(ctx context.Context, userID string, from string, to string) (calories.Summary, error) {
	return calories.Summary{}, nil
}

func (m *mockStore) GetLogsByRange(ctx context.Context, userID string, from string, to string) ([]calories.Log, error) {
	return nil, nil
}

func (m *mockStore) GetLogsByRangeLimited(ctx context.Context, userID string, from string, to string, limit int) ([]calories.Log, error) {
	m.logsLimit = limit
	return []calories.Log{}, nil
}

func (m *mockStore) GetRecentLogs(ctx context.Context, userID string, limit int) ([]calories.Log, error) {
	return nil, nil
}

func TestAddLogValidation(t *testing.T) {
	svc := calories.NewService(&mockStore{})

	_, err := svc.AddLog(context.Background(), calories.AddLogRequest{UserID: ""})
	if !errors.Is(err, calories.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for missing user id, got %v", err)
	}

	userID := "user-1"
	_, err = svc.AddLog(context.Background(), calories.AddLogRequest{UserID: userID})
	if !errors.Is(err, calories.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for missing required fields, got %v", err)
	}
}

func TestGetTargetCreatesDefaultWhenMissing(t *testing.T) {
	store := &mockStore{getTargetErr: sql.ErrNoRows}
	svc := calories.NewService(store)

	target, err := svc.GetTarget(context.Background(), calories.GetTargetRequest{UserID: "user-1"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if target.TargetKcal != 2300 {
		t.Fatalf("expected default target kcal 2300, got %d", target.TargetKcal)
	}
	if target.TargetProteinG != 120 {
		t.Fatalf("expected default target protein 120, got %d", target.TargetProteinG)
	}
	if store.createdTarget == nil {
		t.Fatalf("expected default target to be persisted")
	}
}

func TestUpdateTargetValidation(t *testing.T) {
	svc := calories.NewService(&mockStore{})

	errCase := func(req calories.UpdateTargetRequest) {
		t.Helper()
		_, err := svc.UpdateTarget(context.Background(), req)
		if !errors.Is(err, calories.ErrInvalidInput) {
			t.Fatalf("expected ErrInvalidInput, got %v", err)
		}
	}

	errCase(calories.UpdateTargetRequest{})
	errCase(calories.UpdateTargetRequest{UserID: "user-1"})
}

func TestGetDashboardUsesBoundedLogs(t *testing.T) {
	store := &mockStore{target: &calories.Target{ID: "target-1", UserID: "user-1", TargetKcal: 2300, TargetProteinG: 120}}
	svc := calories.NewService(store)

	_, err := svc.GetDashboard(context.Background(), calories.GetDashboardRequest{UserID: "user-1", LogsLimit: 12})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if store.logsLimit != 12 {
		t.Fatalf("expected logs limit 12, got %d", store.logsLimit)
	}
}
