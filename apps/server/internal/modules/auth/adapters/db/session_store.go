package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/auth"
)

type SessionStore struct {
	db *sql.DB
}

func NewSessionStore(db *sql.DB) *SessionStore {
	return &SessionStore{db: db}
}

func (s *SessionStore) InsertSession(ctx context.Context, session auth.Session) error {
	_, err := s.db.ExecContext(
		ctx,
		`INSERT INTO sessions (
			id,
			user_id,
			created_at,
			expires_at,
			rotated_at,
			last_seen_at,
			absolute_expires_at,
			parent_id,
			revoked_at,
			user_agent_hash,
			ip_prefix
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		session.ID,
		session.UserID,
		session.CreatedAt,
		session.ExpiresAt,
		session.RotatedAt,
		session.LastSeenAt,
		session.AbsoluteExpiresAt,
		session.ParentID,
		session.RevokedAt,
		session.UserAgentHash,
		session.IPPrefix,
	)
	if err != nil {
		return fmt.Errorf("insert session: %w", err)
	}

	return nil
}

func (s *SessionStore) RevokeSession(ctx context.Context, id string, revokedAt string) error {
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE sessions SET revoked_at = ?, expires_at = ? WHERE id = ?",
		revokedAt,
		revokedAt,
		id,
	)
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}

	return nil
}
