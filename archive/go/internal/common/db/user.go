package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID         string `json:"id"`
	CreatedAt  int64  `json:"created_at"`
	LastSeenAt int64  `json:"last_seen_at"`
}

// Session represents a user session
type Session struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
	ExpiresAt int64  `json:"expires_at"`
}

// CreateUser inserts a new user
func (db *DB) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, created_at, last_seen_at)
		VALUES (?, ?, ?)
	`

	_, err := db.ExecContext(ctx, query, user.ID, user.CreatedAt, user.LastSeenAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByID retrieves a user by ID
func (db *DB) GetUserByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT id, created_at, last_seen_at
		FROM users
		WHERE id = ?
	`

	var user User
	err := db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.LastSeenAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// UpdateUserLastSeen updates the last_seen_at timestamp
func (db *DB) UpdateUserLastSeen(ctx context.Context, userID string) error {
	query := `
		UPDATE users
		SET last_seen_at = ?
		WHERE id = ?
	`

	now := time.Now().Unix()
	_, err := db.ExecContext(ctx, query, now, userID)
	if err != nil {
		return fmt.Errorf("failed to update last seen: %w", err)
	}

	return nil
}

// CreateSession creates a new session
func (db *DB) CreateSession(ctx context.Context, session *Session) error {
	query := `
		INSERT INTO sessions (id, user_id, created_at, expires_at)
		VALUES (?, ?, ?, ?)
	`

	_, err := db.ExecContext(ctx, query,
		session.ID,
		session.UserID,
		session.CreatedAt,
		session.ExpiresAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// GetSessionByID retrieves a session and validates expiry
func (db *DB) GetSessionByID(ctx context.Context, sessionID string) (*Session, error) {
	query := `
		SELECT id, user_id, created_at, expires_at
		FROM sessions
		WHERE id = ?
	`

	var session Session
	err := db.QueryRowContext(ctx, query, sessionID).Scan(
		&session.ID,
		&session.UserID,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Check if session is expired
	if time.Now().Unix() > session.ExpiresAt {
		return nil, fmt.Errorf("session expired")
	}

	return &session, nil
}

// DeleteExpiredSessions cleans up old sessions
func (db *DB) DeleteExpiredSessions(ctx context.Context) error {
	query := `
		DELETE FROM sessions
		WHERE expires_at < ?
	`

	now := time.Now().Unix()
	_, err := db.ExecContext(ctx, query, now)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	return nil
}

// NewUser creates a new user with generated ID and timestamps
func NewUser() *User {
	now := time.Now().Unix()
	return &User{
		ID:         uuid.New().String(),
		CreatedAt:  now,
		LastSeenAt: now,
	}
}

// NewSession creates a new session with generated ID and 30-day expiry
func NewSession(userID string) *Session {
	now := time.Now()
	return &Session{
		ID:        uuid.New().String(),
		UserID:    userID,
		CreatedAt: now.Unix(),
		ExpiresAt: now.Add(30 * 24 * time.Hour).Unix(), // 30 days
	}
}
