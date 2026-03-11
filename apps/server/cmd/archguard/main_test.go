package main

import "testing"

func TestCheckPackagesAllowsWiringToComposeAdapters(t *testing.T) {
	packages := []listedPackage{
		{
			ImportPath: "github.com/Hughzu/trackstack/apps/server/internal/wiring/expenses",
			Imports: []string{
				"github.com/Hughzu/trackstack/apps/server/internal/modules/expenses",
				"github.com/Hughzu/trackstack/apps/server/internal/modules/expenses/adapters/db",
			},
			Module: &listedModule{Path: "github.com/Hughzu/trackstack/apps/server"},
		},
	}

	violations := checkPackages(packages)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestCheckPackagesRejectsCrossModuleImport(t *testing.T) {
	packages := []listedPackage{
		{
			ImportPath: "github.com/Hughzu/trackstack/apps/server/internal/modules/expenses",
			Imports: []string{
				"github.com/Hughzu/trackstack/apps/server/internal/modules/users",
			},
			Module: &listedModule{Path: "github.com/Hughzu/trackstack/apps/server"},
		},
	}

	violations := checkPackages(packages)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}

	if violations[0].Rule != "modules must not import other modules directly" {
		t.Fatalf("unexpected rule: %s", violations[0].Rule)
	}
}

func TestCheckPackagesRejectsAdapterImportOutsideWiring(t *testing.T) {
	packages := []listedPackage{
		{
			ImportPath: "github.com/Hughzu/trackstack/apps/server/internal/transport/http",
			Imports: []string{
				"github.com/Hughzu/trackstack/apps/server/internal/modules/expenses/adapters/db",
			},
			Module: &listedModule{Path: "github.com/Hughzu/trackstack/apps/server"},
		},
	}

	violations := checkPackages(packages)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}

	if violations[0].Rule != "only wiring packages may import concrete adapters" {
		t.Fatalf("unexpected rule: %s", violations[0].Rule)
	}
}

func TestCheckPackagesRejectsModuleRuntimeImports(t *testing.T) {
	packages := []listedPackage{
		{
			ImportPath: "github.com/Hughzu/trackstack/apps/server/internal/modules/auth",
			Imports: []string{
				"github.com/Hughzu/trackstack/apps/server/internal/core/config",
				"github.com/Hughzu/trackstack/apps/server/internal/transport/http",
				"github.com/Hughzu/trackstack/apps/server/internal/wiring/auth",
				"github.com/Hughzu/trackstack/apps/server/cmd/server",
			},
			Module: &listedModule{Path: "github.com/Hughzu/trackstack/apps/server"},
		},
	}

	violations := checkPackages(packages)
	if len(violations) != 4 {
		t.Fatalf("expected 4 violations, got %d", len(violations))
	}
}
