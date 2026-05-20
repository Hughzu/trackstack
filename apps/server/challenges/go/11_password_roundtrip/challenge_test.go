package passwordroundtrip

import (
	"errors"
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	t.Parallel()

	_, err := HashPassword("")
	if !errors.Is(err, ErrPasswordRequired) {
		t.Fatalf("HashPassword() error = %v, want %v", err, ErrPasswordRequired)
	}

	hash, err := HashPassword("super-secret")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword() returned empty hash")
	}
	if hash == "super-secret" {
		t.Fatal("HashPassword() returned raw password")
	}
	if !strings.HasPrefix(hash, HashPrefix) {
		t.Fatalf("HashPassword() = %q, want prefix %q", hash, HashPrefix)
	}

	if !VerifyPassword("super-secret", hash) {
		t.Fatal("VerifyPassword() = false, want true for matching password")
	}
	if VerifyPassword("wrong-password", hash) {
		t.Fatal("VerifyPassword() = true, want false for wrong password")
	}
	if VerifyPassword("super-secret", "not-a-valid-hash") {
		t.Fatal("VerifyPassword() = true, want false for invalid hash")
	}
}
