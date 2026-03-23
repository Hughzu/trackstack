package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/users/application/services"
	platformdb "github.com/Hughzu/trackstack/apps/server-next/internal/platform/db"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type commandConfig struct {
	TursoUsersURLHTTP string
	TursoUsersToken   string

	DBMaxOpenConns           int
	DBMaxIdleConns           int
	DBConnMaxLifetimeSeconds int
	DBConnMaxIdleTimeSeconds int
}

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
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	db, err := platformdb.Open(cfg.TursoUsersURLHTTP, cfg.TursoUsersToken, platformdb.PoolConfig{
		MaxOpenConns:    cfg.DBMaxOpenConns,
		MaxIdleConns:    cfg.DBMaxIdleConns,
		ConnMaxLifetime: time.Duration(cfg.DBConnMaxLifetimeSeconds) * time.Second,
		ConnMaxIdleTime: time.Duration(cfg.DBConnMaxIdleTimeSeconds) * time.Second,
	})
	if err != nil {
		return fmt.Errorf("open users db: %w", err)
	}
	defer func() {
		_ = db.Close()
	}()

	passwordHash, err := services.HashPassword(password)
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

func loadConfig() (commandConfig, error) {
	maxOpenConns, err := getEnvInt("DB_MAX_OPEN_CONNS", 10)
	if err != nil {
		return commandConfig{}, err
	}

	maxIdleConns, err := getEnvInt("DB_MAX_IDLE_CONNS", 5)
	if err != nil {
		return commandConfig{}, err
	}

	connMaxLifetimeSeconds, err := getEnvInt("DB_CONN_MAX_LIFETIME_SECONDS", 300)
	if err != nil {
		return commandConfig{}, err
	}

	connMaxIdleTimeSeconds, err := getEnvInt("DB_CONN_MAX_IDLE_TIME_SECONDS", 60)
	if err != nil {
		return commandConfig{}, err
	}

	cfg := commandConfig{
		TursoUsersURLHTTP:        strings.TrimSpace(os.Getenv("TURSO_USERS_URL_HTTP")),
		TursoUsersToken:          strings.TrimSpace(os.Getenv("TURSO_USERS_TOKEN")),
		DBMaxOpenConns:           maxOpenConns,
		DBMaxIdleConns:           maxIdleConns,
		DBConnMaxLifetimeSeconds: connMaxLifetimeSeconds,
		DBConnMaxIdleTimeSeconds: connMaxIdleTimeSeconds,
	}

	if cfg.TursoUsersURLHTTP == "" {
		return commandConfig{}, fmt.Errorf("TURSO_USERS_URL_HTTP must not be empty")
	}

	return cfg, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func getEnvInt(key string, fallback int) (int, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}

	return parsed, nil
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
