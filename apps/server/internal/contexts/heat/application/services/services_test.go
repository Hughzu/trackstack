package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/application/services"
	"github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/domain"
)

type mockRangeLister struct {
	lastUserID string
	lastFrom   string
	lastTo     string
	refills    []domain.Refill
}

func (m *mockRangeLister) ListByRange(ctx context.Context, userID string, from string, to string) ([]domain.Refill, error) {
	m.lastUserID = userID
	m.lastFrom = from
	m.lastTo = to
	if m.refills == nil {
		return []domain.Refill{}, nil
	}
	return m.refills, nil
}

type mockCreator struct {
	created domain.Refill
	err     error
}

func (m *mockCreator) Create(ctx context.Context, refill domain.Refill) error {
	m.created = refill
	return m.err
}

type mockDashboardStore struct {
	recent       []domain.Refill
	recentLimit  int
	recentOffset int
	latest       *domain.Refill
	sums         []int
	sumCalls     int
}

func (m *mockDashboardStore) ListRecent(ctx context.Context, userID string, limit int, offset int) ([]domain.Refill, error) {
	m.recentLimit = limit
	m.recentOffset = offset
	if m.recent == nil {
		return []domain.Refill{}, nil
	}
	return m.recent, nil
}

func (m *mockDashboardStore) GetLatest(ctx context.Context, userID string) (*domain.Refill, error) {
	return m.latest, nil
}

func (m *mockDashboardStore) GetSumByRange(ctx context.Context, userID string, from string, to string) (int, error) {
	if m.sumCalls >= len(m.sums) {
		m.sumCalls++
		return 0, nil
	}
	value := m.sums[m.sumCalls]
	m.sumCalls++
	return value, nil
}

func TestCreateRefillService(t *testing.T) {
	creator := &mockCreator{}
	svc := services.NewCreateRefillService(creator)

	_, err := svc.Execute(context.Background(), services.CreateRefillRequest{
		UserID:   "",
		Date:     "2023-10-15",
		WeightKg: 15,
		Bags:     1,
	})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for missing user, got %v", err)
	}

	refill, err := svc.Execute(context.Background(), services.CreateRefillRequest{
		UserID:   "user-1",
		Date:     "2023-10-15",
		WeightKg: 15,
		Bags:     1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if refill.ID == "" {
		t.Fatal("expected generated refill id")
	}
	if refill.Season == nil || *refill.Season != "2023-2024" {
		t.Fatalf("expected season 2023-2024, got %+v", refill.Season)
	}
	if creator.created.UserID != "user-1" {
		t.Fatalf("expected created refill user user-1, got %q", creator.created.UserID)
	}

	refillJan, err := svc.Execute(context.Background(), services.CreateRefillRequest{
		UserID:   "user-1",
		Date:     "2024-01-15",
		WeightKg: 15,
		Bags:     1,
	})
	if err != nil {
		t.Fatalf("expected no error for january refill, got %v", err)
	}
	if refillJan.Season == nil || *refillJan.Season != "2023-2024" {
		t.Fatalf("expected season 2023-2024 for january refill, got %+v", refillJan.Season)
	}
}

func TestListRefillsService(t *testing.T) {
	lister := &mockRangeLister{}
	svc := services.NewListRefillsService(lister)

	_, err := svc.Execute(context.Background(), services.ListRefillsRequest{
		UserID: "user-1",
		From:   "invalid-date",
		To:     "2023-10-15",
	})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for invalid date, got %v", err)
	}

	_, err = svc.Execute(context.Background(), services.ListRefillsRequest{
		UserID: "user-1",
		From:   "2023-10-15",
		To:     "2023-10-16T15:00:00Z",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if lister.lastUserID != "user-1" {
		t.Fatalf("expected user-1, got %q", lister.lastUserID)
	}
	expectedFrom := "2023-10-15T00:00:00Z"
	if lister.lastFrom != expectedFrom {
		t.Fatalf("expected normalized from %q, got %q", expectedFrom, lister.lastFrom)
	}
	if lister.lastTo != "2023-10-16T15:00:00Z" {
		t.Fatalf("expected normalized to remain RFC3339, got %q", lister.lastTo)
	}
}

func TestGetDashboardServiceUsesRecentSlice(t *testing.T) {
	latestDate := time.Now().UTC().Add(-48 * time.Hour).Format(time.RFC3339)
	store := &mockDashboardStore{
		recent: []domain.Refill{{ID: "refill-1"}},
		latest: &domain.Refill{ID: "latest-1", Date: latestDate},
		sums:   []int{12, 10},
	}
	svc := services.NewGetDashboardService(store, store, store)

	viewModel, err := svc.Execute(context.Background(), services.GetDashboardRequest{
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
		t.Fatalf("expected history from recent query, got %+v", viewModel.History)
	}
	if viewModel.SeasonSnapshot.SeasonToDate != 12 {
		t.Fatalf("expected season-to-date 12, got %d", viewModel.SeasonSnapshot.SeasonToDate)
	}
	if viewModel.SeasonSnapshot.LastSeasonToDate != 10 {
		t.Fatalf("expected last-season-to-date 10, got %d", viewModel.SeasonSnapshot.LastSeasonToDate)
	}
	if viewModel.SeasonSnapshot.Delta != 2 {
		t.Fatalf("expected delta 2, got %d", viewModel.SeasonSnapshot.Delta)
	}
	if viewModel.SeasonSnapshot.DeltaPct == nil || *viewModel.SeasonSnapshot.DeltaPct != 20 {
		t.Fatalf("expected delta pct 20, got %+v", viewModel.SeasonSnapshot.DeltaPct)
	}
	if viewModel.DaysSinceRefill < 1 {
		t.Fatalf("expected days since refill to be positive, got %d", viewModel.DaysSinceRefill)
	}
}
