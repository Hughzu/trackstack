package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/users"
)

type UsersStore struct {
	db *sql.DB
}

func NewUsersStore(db *sql.DB) *UsersStore {
	return &UsersStore{db: db}
}

func (s *UsersStore) FindByEmail(ctx context.Context, email string) (users.User, error) {
	row := s.db.QueryRowContext(
		ctx,
		`SELECT
			id,
			email,
			password_hash,
			session_version,
			created_at,
			last_login_at
		FROM users
		WHERE email = ?
		LIMIT 1`,
		email,
	)

	var user users.User
	var lastLogin sql.NullString
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.SessionVersion,
		&user.CreatedAt,
		&lastLogin,
	)
	if err != nil {
		return users.User{}, err
	}

	if lastLogin.Valid {
		user.LastLoginAt = &lastLogin.String
	}

	return user, nil
}

func (s *UsersStore) UpdateLastLogin(ctx context.Context, userID string, lastLoginAt string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE users SET last_login_at = ? WHERE id = ?", lastLoginAt, userID)
	if err != nil {
		return fmt.Errorf("update last login: %w", err)
	}

	return nil
}
