package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/application/ports"
	"github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/domain"
)

type SessionRepository struct {
	db *sql.DB
}

var _ ports.SessionRepository = (*SessionRepository)(nil)

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(ctx context.Context, input ports.CreateSessionInput) (domain.RefreshSession, error) {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO sessions (
			id,
			user_id,
			token_hash,
			created_at,
			expires_at,
			rotated_at,
			last_seen_at,
			absolute_expires_at,
			parent_id,
			revoked_at,
			user_agent_hash,
			ip_prefix
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NULL, NULL, ?, ?)`,
		input.ID,
		input.UserID,
		input.TokenHash,
		formatTime(input.CreatedAt),
		formatTime(input.ExpiresAt),
		formatTime(input.CreatedAt),
		formatTime(input.CreatedAt),
		formatTime(input.AbsoluteExpiresAt),
		nullableString(input.UserAgentHash),
		nullableString(input.IPPrefix),
	)
	if err != nil {
		return domain.RefreshSession{}, fmt.Errorf("create session: %w", err)
	}

	return domain.RefreshSession{
		ID:                input.ID,
		UserID:            input.UserID,
		TokenHash:         input.TokenHash,
		CreatedAt:         input.CreatedAt,
		ExpiresAt:         input.ExpiresAt,
		RotatedAt:         input.CreatedAt,
		LastSeenAt:        input.CreatedAt,
		AbsoluteExpiresAt: input.AbsoluteExpiresAt,
		UserAgentHash:     input.UserAgentHash,
		IPPrefix:          input.IPPrefix,
	}, nil
}

func (r *SessionRepository) FindByTokenHash(ctx context.Context, tokenHash string) (domain.RefreshSession, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT
			id,
			user_id,
			token_hash,
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
		WHERE token_hash = ?
		LIMIT 1`,
		tokenHash,
	)

	return scanSession(row)
}

func (r *SessionRepository) RotateSession(ctx context.Context, input ports.RotateSessionInput) (domain.RefreshSession, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.RefreshSession{}, fmt.Errorf("begin rotate session transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	result, err := tx.ExecContext(
		ctx,
		`UPDATE sessions
		SET rotated_at = ?, last_seen_at = ?
		WHERE id = ? AND token_hash = ?`,
		formatTime(input.RotatedAt),
		formatTime(input.NewLastSeenAt),
		input.SessionID,
		input.TokenHash,
	)
	if err != nil {
		return domain.RefreshSession{}, fmt.Errorf("mark session rotated: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.RefreshSession{}, fmt.Errorf("rotate session rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return domain.RefreshSession{}, domain.ErrSessionNotFound
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO sessions (
			id,
			user_id,
			token_hash,
			created_at,
			expires_at,
			rotated_at,
			last_seen_at,
			absolute_expires_at,
			parent_id,
			revoked_at,
			user_agent_hash,
			ip_prefix
		)
		SELECT ?, user_id, ?, ?, ?, ?, ?, ?, id, NULL, ?, ?
		FROM sessions
		WHERE id = ?`,
		input.NewSessionID,
		input.NewTokenHash,
		formatTime(input.RotatedAt),
		formatTime(input.NewExpiresAt),
		formatTime(input.RotatedAt),
		formatTime(input.NewLastSeenAt),
		formatTime(input.AbsoluteExpiresAt),
		nullableString(input.UserAgentHash),
		nullableString(input.IPPrefix),
		input.SessionID,
	)
	if err != nil {
		return domain.RefreshSession{}, fmt.Errorf("insert rotated session: %w", err)
	}

	row := tx.QueryRowContext(
		ctx,
		`SELECT
			id,
			user_id,
			token_hash,
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
		WHERE id = ?`,
		input.NewSessionID,
	)

	session, err := scanSession(row)
	if err != nil {
		return domain.RefreshSession{}, err
	}

	if err := tx.Commit(); err != nil {
		return domain.RefreshSession{}, fmt.Errorf("commit rotate session: %w", err)
	}

	return session, nil
}

func (r *SessionRepository) RevokeSession(ctx context.Context, sessionID string, revokedAt time.Time) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE sessions
		SET revoked_at = COALESCE(revoked_at, ?)
		WHERE id = ?`,
		formatTime(revokedAt),
		sessionID,
	)
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}

	return nil
}

func (r *SessionRepository) RevokeFamily(ctx context.Context, sessionID string, revokedAt time.Time) error {
	row := r.db.QueryRowContext(
		ctx,
		`WITH RECURSIVE ancestors(id, parent_id) AS (
			SELECT id, parent_id FROM sessions WHERE id = ?
			UNION ALL
			SELECT s.id, s.parent_id
			FROM sessions s
			JOIN ancestors a ON a.parent_id = s.id
		)
		SELECT COALESCE((SELECT id FROM ancestors WHERE parent_id IS NULL LIMIT 1), ?)`,
		sessionID,
		sessionID,
	)

	var rootID string
	if err := row.Scan(&rootID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("find session family root: %w", err)
	}

	_, err := r.db.ExecContext(
		ctx,
		`WITH RECURSIVE family(id) AS (
			SELECT id FROM sessions WHERE id = ?
			UNION ALL
			SELECT s.id
			FROM sessions s
			JOIN family f ON s.parent_id = f.id
		)
		UPDATE sessions
		SET revoked_at = COALESCE(revoked_at, ?)
		WHERE id IN (SELECT id FROM family)`,
		rootID,
		formatTime(revokedAt),
	)
	if err != nil {
		return fmt.Errorf("revoke session family: %w", err)
	}

	return nil
}

func scanSession(scanner interface{ Scan(dest ...any) error }) (domain.RefreshSession, error) {
	var (
		session         domain.RefreshSession
		createdAt       string
		expiresAt       string
		rotatedAt       string
		lastSeenAt      string
		absoluteExpires string
		parentID        sql.NullString
		revokedAt       sql.NullString
		userAgentHash   sql.NullString
		ipPrefix        sql.NullString
	)

	err := scanner.Scan(
		&session.ID,
		&session.UserID,
		&session.TokenHash,
		&createdAt,
		&expiresAt,
		&rotatedAt,
		&lastSeenAt,
		&absoluteExpires,
		&parentID,
		&revokedAt,
		&userAgentHash,
		&ipPrefix,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.RefreshSession{}, domain.ErrSessionNotFound
		}
		return domain.RefreshSession{}, fmt.Errorf("scan session: %w", err)
	}

	if session.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
		return domain.RefreshSession{}, fmt.Errorf("parse created_at: %w", err)
	}
	if session.ExpiresAt, err = time.Parse(time.RFC3339, expiresAt); err != nil {
		return domain.RefreshSession{}, fmt.Errorf("parse expires_at: %w", err)
	}
	if session.RotatedAt, err = time.Parse(time.RFC3339, rotatedAt); err != nil {
		return domain.RefreshSession{}, fmt.Errorf("parse rotated_at: %w", err)
	}
	if session.LastSeenAt, err = time.Parse(time.RFC3339, lastSeenAt); err != nil {
		return domain.RefreshSession{}, fmt.Errorf("parse last_seen_at: %w", err)
	}
	if session.AbsoluteExpiresAt, err = time.Parse(time.RFC3339, absoluteExpires); err != nil {
		return domain.RefreshSession{}, fmt.Errorf("parse absolute_expires_at: %w", err)
	}
	if parentID.Valid {
		session.ParentID = &parentID.String
	}
	if revokedAt.Valid {
		parsed, parseErr := time.Parse(time.RFC3339, revokedAt.String)
		if parseErr != nil {
			return domain.RefreshSession{}, fmt.Errorf("parse revoked_at: %w", parseErr)
		}
		session.RevokedAt = &parsed
	}
	if userAgentHash.Valid {
		session.UserAgentHash = &userAgentHash.String
	}
	if ipPrefix.Valid {
		session.IPPrefix = &ipPrefix.String
	}

	return session, nil
}

func formatTime(value time.Time) string {
	return value.UTC().Format(time.RFC3339)
}

func nullableString(value *string) any {
	if value == nil {
		return nil
	}

	return *value
}
