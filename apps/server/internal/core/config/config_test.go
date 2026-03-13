package config

import "testing"

func TestValidateRequiresHTTPURLsWhenHTTPMode(t *testing.T) {
	cfg := Config{DBConnectionMode: "HTTP"}

	err := cfg.Validate()
	if err == nil {
		t.Fatalf("expected validation error for missing HTTP URLs")
	}
}

func TestValidateAcceptsExplicitHTTPURLs(t *testing.T) {
	cfg := Config{
		DBConnectionMode:     "HTTP",
		TursoCaloriesURLHTTP: "https://calories.example.turso.io",
		TursoExpensesURLHTTP: "https://expenses.example.turso.io",
		TursoHeatURLHTTP:     "https://heat.example.turso.io",
		TursoUsersURLHTTP:    "https://users.example.turso.io",
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid config, got %v", err)
	}
}

func TestValidateRequiresWSURLsWhenWSMode(t *testing.T) {
	cfg := Config{
		DBConnectionMode:   "WS",
		TursoCaloriesURLWS: "wss://calories.example.turso.io",
		TursoExpensesURLWS: "wss://expenses.example.turso.io",
		TursoHeatURLWS:     "wss://heat.example.turso.io",
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatalf("expected validation error for missing users WS URL")
	}
}
