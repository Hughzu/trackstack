package monthsnapshot

import (
	"testing"
	"time"
)

func TestBuildMonthSnapshot(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.May, 20, 12, 0, 0, 0, time.UTC)
	entries := []Entry{
		{OccurredAt: time.Date(2026, time.May, 1, 8, 0, 0, 0, time.UTC), Amount: 10},
		{OccurredAt: time.Date(2026, time.May, 18, 8, 0, 0, 0, time.UTC), Amount: 20},
		{OccurredAt: time.Date(2026, time.April, 2, 8, 0, 0, 0, time.UTC), Amount: 15},
		{OccurredAt: time.Date(2026, time.April, 21, 8, 0, 0, 0, time.UTC), Amount: 5},
		{OccurredAt: time.Date(2026, time.March, 9, 8, 0, 0, 0, time.UTC), Amount: 99},
	}

	got := BuildMonthSnapshot(now, entries)
	wantPct := 50
	want := Snapshot{
		CurrentTotal:  30,
		PreviousTotal: 20,
		Delta:         10,
		DeltaPct:      &wantPct,
	}

	assertSnapshot(t, got, want)

	gotNoPrevious := BuildMonthSnapshot(now, []Entry{{OccurredAt: now, Amount: 30}})
	if gotNoPrevious.PreviousTotal != 0 {
		t.Fatalf("PreviousTotal = %d, want 0", gotNoPrevious.PreviousTotal)
	}
	if gotNoPrevious.DeltaPct != nil {
		t.Fatalf("DeltaPct = %v, want nil", *gotNoPrevious.DeltaPct)
	}
}

func assertSnapshot(t *testing.T, got Snapshot, want Snapshot) {
	t.Helper()

	if got.CurrentTotal != want.CurrentTotal {
		t.Fatalf("CurrentTotal = %d, want %d", got.CurrentTotal, want.CurrentTotal)
	}
	if got.PreviousTotal != want.PreviousTotal {
		t.Fatalf("PreviousTotal = %d, want %d", got.PreviousTotal, want.PreviousTotal)
	}
	if got.Delta != want.Delta {
		t.Fatalf("Delta = %d, want %d", got.Delta, want.Delta)
	}
	if !sameIntPtr(got.DeltaPct, want.DeltaPct) {
		t.Fatalf("DeltaPct = %v, want %v", ptrValue(got.DeltaPct), ptrValue(want.DeltaPct))
	}
}

func sameIntPtr(left *int, right *int) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func ptrValue(value *int) any {
	if value == nil {
		return nil
	}
	return *value
}
