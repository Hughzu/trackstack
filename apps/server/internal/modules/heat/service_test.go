package heat_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/heat"
)

type mockRefillStore struct {
}

func (m *mockRefillStore) ListByRange(ctx context.Context, userID string, from string, to string) ([]heat.Refill, error) {
	return []heat.Refill{}, nil
}

func (m *mockRefillStore) Create(ctx context.Context, refill heat.Refill) error {
	return nil
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
