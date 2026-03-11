package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/core/config"
	coredb "github.com/Hughzu/trackstack/apps/server/internal/core/db"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/auth"
	"github.com/Hughzu/trackstack/apps/server/internal/wiring/common"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type userRecord struct {
	ID string
}

func init() {
	_ = godotenv.Load("../web/.env", ".env")
}

func main() {
	emailArg := flag.String("email", "", "email for the seeded user")
	passwordArg := flag.String("password", "", "password for the seeded user")
	flag.Parse()

	email := firstNonEmpty(*emailArg, os.Getenv("E2E_TEST_EMAIL"))
	password := firstNonEmpty(*passwordArg, os.Getenv("E2E_TEST_PASSWORD"))

	if strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
		_, _ = os.Stderr.WriteString("Usage: go run ./cmd/seed-user --email you@example.com --password yourpass\n")
		_, _ = os.Stderr.WriteString("Or set E2E_TEST_EMAIL and E2E_TEST_PASSWORD in the environment.\n")
		os.Exit(1)
	}

	if err := run(strings.ToLower(strings.TrimSpace(email)), password); err != nil {
		_, _ = os.Stderr.WriteString("Failed to seed user: " + err.Error() + "\n")
		os.Exit(1)
	}
}

func run(email string, password string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	dsn, err := common.BuildTursoDSN(common.TursoConfig{
		Mode:    cfg.DBConnectionMode,
		URL:     cfg.TursoUsersURL,
		URLHTTP: cfg.TursoUsersURLHTTP,
		URLWS:   cfg.TursoUsersURLWS,
		Token:   cfg.TursoUsersToken,
	})
	if err != nil {
		return fmt.Errorf("build users dsn: %w", err)
	}

	db, err := coredb.OpenLibSQL(dsn)
	if err != nil {
		return fmt.Errorf("open users db: %w", err)
	}
	defer func() {
		_ = db.Close()
	}()

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	ctx := context.Background()
	existingID, err := findUserByEmail(ctx, db, email)
	if err != nil {
		return err
	}

	if existingID != "" {
		if err := updatePassword(ctx, db, existingID, passwordHash); err != nil {
			return err
		}
		_, _ = os.Stdout.WriteString("Updated password for existing user " + email + "\n")
		return nil
	}

	id := uuid.NewString()
	if err := insertUser(ctx, db, id, email, passwordHash, time.Now().UTC().Format(time.RFC3339)); err != nil {
		return err
	}

	_, _ = os.Stdout.WriteString("Created user " + email + " with id " + id + "\n")
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func findUserByEmail(ctx context.Context, db *sql.DB, email string) (string, error) {
	row := db.QueryRowContext(ctx, "SELECT id FROM users WHERE email = ? LIMIT 1", email)
	var record userRecord
	if err := row.Scan(&record.ID); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("find user by email: %w", err)
	}
	return record.ID, nil
}

func updatePassword(ctx context.Context, db *sql.DB, id string, passwordHash string) error {
	if _, err := db.ExecContext(ctx, "UPDATE users SET password_hash = ? WHERE id = ?", passwordHash, id); err != nil {
		return fmt.Errorf("update user password: %w", err)
	}
	return nil
}

func insertUser(ctx context.Context, db *sql.DB, id string, email string, passwordHash string, createdAt string) error {
	if _, err := db.ExecContext(ctx, "INSERT INTO users (id, email, password_hash, created_at) VALUES (?, ?, ?, ?)", id, email, passwordHash, createdAt); err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}
