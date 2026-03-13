package common

import "testing"

func TestResolveTursoSelectionHTTP(t *testing.T) {
	selection, err := ResolveTursoSelection(TursoConfig{
		Mode:    "HTTP",
		URLHTTP: "https://calories.example.turso.io",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if selection.Scheme != "https" {
		t.Fatalf("expected https scheme, got %s", selection.Scheme)
	}
	if selection.Mode != "HTTP" {
		t.Fatalf("expected HTTP mode, got %s", selection.Mode)
	}
}

func TestResolveTursoSelectionWS(t *testing.T) {
	selection, err := ResolveTursoSelection(TursoConfig{
		Mode:  "WS",
		URLWS: "wss://calories.example.turso.io",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if selection.Scheme != "wss" {
		t.Fatalf("expected wss scheme, got %s", selection.Scheme)
	}
}

func TestBuildTursoDSNAppendsAuthToken(t *testing.T) {
	dsn, err := BuildTursoDSN(TursoConfig{
		Mode:    "HTTP",
		URLHTTP: "https://calories.example.turso.io",
		Token:   "secret-token",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if dsn != "https://calories.example.turso.io?authToken=secret-token" {
		t.Fatalf("unexpected dsn %q", dsn)
	}
}
