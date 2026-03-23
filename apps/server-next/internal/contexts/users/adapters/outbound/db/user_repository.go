package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/users/application/ports"
	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/users/domain"
)

type UserRepository struct {
	db *sql.DB
}

var _ ports.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	row := r.db.QueryRowContext(
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

	var user domain.User
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
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("scan user: %w", err)
	}

	if lastLogin.Valid {
		user.LastLoginAt = &lastLogin.String
	}

	return user, nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID string, lastLoginAt string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET last_login_at = ? WHERE id = ?", lastLoginAt, userID)
	if err != nil {
		return fmt.Errorf("update last login: %w", err)
	}

	return nil
}
