package main

import (
	"strings"
	"testing"
)

func TestCheckSourceAllowsPlainAdapterFile(t *testing.T) {
	src := `package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/expenses"
)

type ExpensesStore struct { db *sql.DB }

func (s *ExpensesStore) CreateExpense(ctx context.Context, entry expenses.Entry) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO expense_entries (id) VALUES (?)", entry.ID)
	if err != nil {
		return fmt.Errorf("insert expense entry: %w", err)
	}
	return nil
}
`

	violations, err := checkSource("internal/modules/expenses/adapters/db/entry_store.go", src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestCheckSourceRejectsThirdPartyImport(t *testing.T) {
	src := `package db

import "github.com/google/uuid"
`

	violations, err := checkSource("internal/modules/expenses/adapters/db/entry_store.go", src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Rule, "third-party") {
		t.Fatalf("unexpected rule: %s", violations[0].Rule)
	}
}

func TestCheckSourceRejectsOtherApplicationImports(t *testing.T) {
	src := `package db

import "github.com/Hughzu/trackstack/apps/server/internal/transport/http"
`

	violations, err := checkSource("internal/modules/expenses/adapters/db/entry_store.go", src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Rule, "own module package") {
		t.Fatalf("unexpected rule: %s", violations[0].Rule)
	}
}

func TestCheckSourceRejectsBranchHeavyAdapterMethod(t *testing.T) {
	src := `package db

func (s *ExpensesStore) Busy() {
	if true {}
	if true {}
	if true {}
	if true {}
	if true {}
	if true {}
	if true {}
	if true {}
	if true {}
}
`

	violations, err := checkSource("internal/modules/expenses/adapters/db/entry_store.go", src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Rule, "branch count") {
		t.Fatalf("unexpected rule: %s", violations[0].Rule)
	}
}

func TestCheckSourceRejectsTimestampAndIDGeneration(t *testing.T) {
	src := `package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (s *ExpensesStore) Create(ctx context.Context) {
	_ = ctx
	_ = time.Now()
	_ = uuid.NewString()
}
`

	violations, err := checkSource("internal/modules/expenses/adapters/db/entry_store.go", src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 3 {
		t.Fatalf("expected 3 violations, got %d", len(violations))
	}
}

func TestCheckSourceRejectsInputNormalization(t *testing.T) {
	src := `package db

import "strings"

func (s *ExpensesStore) Normalize() string {
	return strings.TrimSpace(" demo ")
}
`

	violations, err := checkSource("internal/modules/expenses/adapters/db/entry_store.go", src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Rule, "normalize") {
		t.Fatalf("unexpected rule: %s", violations[0].Rule)
	}
}
