package userlookup

import (
	"context"
	"errors"
	"testing"
)

func TestUserServiceLookupByEmail(t *testing.T) {
	t.Parallel()

	t.Run("blank email", func(t *testing.T) {
		t.Parallel()

		svc := NewUserService(fakeUserRepo{})
		_, err := svc.LookupByEmail(context.Background(), "   ")
		if !errors.Is(err, ErrInvalidEmail) {
			t.Fatalf("LookupByEmail() error = %v, want %v", err, ErrInvalidEmail)
		}
	})

	t.Run("repo error", func(t *testing.T) {
		t.Parallel()

		wantErr := errors.New("boom")
		svc := NewUserService(fakeUserRepo{err: wantErr})
		_, err := svc.LookupByEmail(context.Background(), "alice@example.com")
		if !errors.Is(err, wantErr) {
			t.Fatalf("LookupByEmail() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		t.Parallel()

		svc := NewUserService(fakeUserRepo{})
		_, err := svc.LookupByEmail(context.Background(), "alice@example.com")
		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("LookupByEmail() error = %v, want %v", err, ErrUserNotFound)
		}
	})

	t.Run("success with normalized email", func(t *testing.T) {
		t.Parallel()

		repo := fakeUserRepo{user: User{ID: "user-1", Email: "alice@example.com"}, found: true}
		svc := NewUserService(repo)
		user, err := svc.LookupByEmail(context.Background(), "  Alice@Example.com ")
		if err != nil {
			t.Fatalf("LookupByEmail() error = %v", err)
		}
		if user != repo.user {
			t.Fatalf("LookupByEmail() = %#v, want %#v", user, repo.user)
		}
	})
}

type fakeUserRepo struct {
	user  User
	found bool
	err   error
}

func (f fakeUserRepo) FindByEmail(_ context.Context, _ string) (User, bool, error) {
	if f.err != nil {
		return User{}, false, f.err
	}
	return f.user, f.found, nil
}
