package targetservice

import (
	"context"
	"errors"
	"testing"
)

func TestServiceGetOrCreateTarget(t *testing.T) {
	t.Parallel()

	t.Run("missing user id", func(t *testing.T) {
		t.Parallel()

		svc := NewService(fakeTargetRepo{}, func() string { return "target-1" })
		_, err := svc.GetOrCreateTarget(context.Background(), "  ")
		if !errors.Is(err, ErrUserIDRequired) {
			t.Fatalf("GetOrCreateTarget() error = %v, want %v", err, ErrUserIDRequired)
		}
	})

	t.Run("existing target", func(t *testing.T) {
		t.Parallel()

		existing := Target{ID: "target-1", UserID: "user-1", TargetCalories: 2500, TargetProteinGrams: 140}
		repo := fakeTargetRepo{target: existing, found: true}
		svc := NewService(repo, func() string { return "unused" })

		got, err := svc.GetOrCreateTarget(context.Background(), "user-1")
		if err != nil {
			t.Fatalf("GetOrCreateTarget() error = %v", err)
		}
		if got != existing {
			t.Fatalf("GetOrCreateTarget() = %#v, want %#v", got, existing)
		}
	})

	t.Run("create default target", func(t *testing.T) {
		t.Parallel()

		repo := fakeTargetRepo{}
		svc := NewService(repo, func() string { return "target-42" })

		got, err := svc.GetOrCreateTarget(context.Background(), "user-42")
		if err != nil {
			t.Fatalf("GetOrCreateTarget() error = %v", err)
		}

		want := Target{
			ID:                 "target-42",
			UserID:             "user-42",
			TargetCalories:     DefaultTargetCalories,
			TargetProteinGrams: DefaultTargetProteinGrams,
		}

		if got != want {
			t.Fatalf("GetOrCreateTarget() = %#v, want %#v", got, want)
		}
	})
}

type fakeTargetRepo struct {
	target Target
	found  bool
	err    error
}

func (f fakeTargetRepo) FindByUserID(_ context.Context, _ string) (Target, bool, error) {
	if f.err != nil {
		return Target{}, false, f.err
	}
	return f.target, f.found, nil
}

func (f fakeTargetRepo) Create(_ context.Context, _ Target) error {
	return f.err
}
