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

func (s *SessionStore) FindSessionByID(ctx context.Context, id string) (auth.Session, error) {
	row := s.db.QueryRowContext(
		ctx,
		`SELECT
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
		FROM sessions
		WHERE id = ?
		LIMIT 1`,
		id,
	)

	var session auth.Session
	var parentID sql.NullString
	var revokedAt sql.NullString
	var userAgentHash sql.NullString
	var ipPrefix sql.NullString

	err := row.Scan(
		&session.ID,
		&session.UserID,
		&session.CreatedAt,
		&session.ExpiresAt,
		&session.RotatedAt,
		&session.LastSeenAt,
		&session.AbsoluteExpiresAt,
		&parentID,
		&revokedAt,
		&userAgentHash,
		&ipPrefix,
	)
	if err != nil {
		return auth.Session{}, err
	}

	if parentID.Valid {
		session.ParentID = &parentID.String
	}
	if revokedAt.Valid {
		session.RevokedAt = &revokedAt.String
	}
	if userAgentHash.Valid {
		session.UserAgentHash = &userAgentHash.String
	}
	if ipPrefix.Valid {
		session.IPPrefix = &ipPrefix.String
	}

	return session, nil
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

func (s *SessionStore) TouchSession(ctx context.Context, id string, lastSeenAt string, expiresAt string) error {
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE sessions SET last_seen_at = ?, expires_at = ? WHERE id = ?",
		lastSeenAt,
		expiresAt,
		id,
	)
	if err != nil {
		return fmt.Errorf("touch session: %w", err)
	}

	return nil
}

func (s *SessionStore) RotateOutSession(ctx context.Context, id string, revokedAt string, expiresAt string, rotatedAt string) error {
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE sessions SET revoked_at = ?, expires_at = ?, rotated_at = ? WHERE id = ?",
		revokedAt,
		expiresAt,
		rotatedAt,
		id,
	)
	if err != nil {
		return fmt.Errorf("rotate out session: %w", err)
	}

	return nil
}
