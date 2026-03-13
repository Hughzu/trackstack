package heat_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/heat"
)

type mockRefillStore struct {
	recent       []heat.Refill
	recentLimit  int
	recentOffset int
}

func (m *mockRefillStore) ListByRange(ctx context.Context, userID string, from string, to string) ([]heat.Refill, error) {
	return []heat.Refill{}, nil
}

func (m *mockRefillStore) ListRecent(ctx context.Context, userID string, limit int, offset int) ([]heat.Refill, error) {
	m.recentLimit = limit
	m.recentOffset = offset
	if m.recent == nil {
		return []heat.Refill{}, nil
	}
	return m.recent, nil
}

func (m *mockRefillStore) Create(ctx context.Context, refill heat.Refill) error {
	return nil
}

func (m *mockRefillStore) GetLatest(ctx context.Context, userID string) (*heat.Refill, error) {
	return nil, nil
}

func (m *mockRefillStore) GetSumByRange(ctx context.Context, userID string, fromDate string, toDate string) (int, error) {
	return 0, nil
}

func (m *mockRefillStore) Delete(ctx context.Context, userID string, id string) (bool, error) {
	return true, nil
}

func TestCreateRefill(t *testing.T) {
	store := &mockRefillStore{}
	svc := heat.NewService(store)

	ctx := context.Background()

	// Should fail on empty user
	_, err := svc.CreateRefill(ctx, heat.CreateRefillRequest{
		UserID:   "",
		Date:     "2023-10-15",
		WeightKg: 15.0,
		Bags:     1,
	})
	if !errors.Is(err, heat.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput for missing UserID, got %v", err)
	}

	// Should succeed on valid input
	refill, err := svc.CreateRefill(ctx, heat.CreateRefillRequest{
		UserID:   "user-1",
		Date:     "2023-10-15",
		WeightKg: 15.0,
		Bags:     1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if refill.ID == "" {
		t.Errorf("expected ID to be generated, got empty")
	}

	// verify season logic based on October 2023 date -> 2023-2024
	if refill.Season == nil || *refill.Season != "2023-2024" {
		season := "nil"
		if refill.Season != nil {
			season = *refill.Season
		}
		t.Errorf("expected season 2023-2024, got %s", season)
	}

	// verify season logic based on Jan 2024 date -> 2023-2024
	refill2, _ := svc.CreateRefill(ctx, heat.CreateRefillRequest{
		UserID:   "user-1",
		Date:     "2024-01-15",
		WeightKg: 15,
		Bags:     1,
	})
	if refill2.Season == nil || *refill2.Season != "2023-2024" {
		season := "nil"
		if refill2.Season != nil {
			season = *refill2.Season
		}
		t.Errorf("expected season 2023-2024 for Jan 2024, got %s", season)
	}
}

func TestListRefills(t *testing.T) {
	store := &mockRefillStore{}
	svc := heat.NewService(store)

	ctx := context.Background()

	// Invalid date format
	_, err := svc.ListRefills(ctx, heat.ListRefillsRequest{
		UserID: "user-1",
		From:   "invalid-date",
		To:     "2023-10-15",
	})
	if !errors.Is(err, heat.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput for invalid format, got %v", err)
	}

	// Valid dates
	_, err = svc.ListRefills(ctx, heat.ListRefillsRequest{
		UserID: "user-1",
		From:   "2023-10-15",
		To:     "2023-10-16T15:00:00Z",
	})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestGetDashboardUsesRecentSlice(t *testing.T) {
	store := &mockRefillStore{
		recent: []heat.Refill{{ID: "refill-1"}},
	}
	svc := heat.NewService(store)

	viewModel, err := svc.GetDashboard(context.Background(), heat.GetDashboardRequest{
		UserID: "user-1",
		Page:   2,
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if store.recentLimit != 10 {
		t.Fatalf("expected recent limit 10, got %d", store.recentLimit)
	}
	if store.recentOffset != 10 {
		t.Fatalf("expected recent offset 10, got %d", store.recentOffset)
	}
	if len(viewModel.History) != 1 || viewModel.History[0].ID != "refill-1" {
		t.Fatalf("expected dashboard history to come from recent query, got %+v", viewModel.History)
	}
}
