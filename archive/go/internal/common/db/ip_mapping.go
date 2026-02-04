package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// IPMapping represents an IP address to user mapping
type IPMapping struct {
	IPAddress   string `json:"ip_address"`
	UserID      string `json:"user_id"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
}

// GetUserByIP retrieves the user associated with an IP address
func (db *DB) GetUserByIP(ctx context.Context, ipAddress string) (*User, error) {
	query := `
		SELECT u.id, u.created_at, u.last_seen_at
		FROM users u
		JOIN ip_mappings ip ON u.id = ip.user_id
		WHERE ip.ip_address = ?
	`

	var user User
	err := db.QueryRowContext(ctx, query, ipAddress).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.LastSeenAt,
	)

	if err != nil {
		return nil, fmt.Errorf("no user found for IP %s: %w", ipAddress, err)
	}

	return &user, nil
}

// CreateIPMapping creates a new IP to user mapping
func (db *DB) CreateIPMapping(ctx context.Context, mapping *IPMapping) error {
	query := `
		INSERT INTO ip_mappings (ip_address, user_id, description, created_at)
		VALUES (?, ?, ?, ?)
	`

	_, err := db.ExecContext(ctx, query,
		mapping.IPAddress,
		mapping.UserID,
		mapping.Description,
		mapping.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create IP mapping: %w", err)
	}

	return nil
}

// ListIPMappings returns all IP mappings for a user
func (db *DB) ListIPMappings(ctx context.Context, userID string) ([]*IPMapping, error) {
	query := `
		SELECT ip_address, user_id, description, created_at
		FROM ip_mappings
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list IP mappings: %w", err)
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var mappings []*IPMapping
	for rows.Next() {
		var m IPMapping
		if err := rows.Scan(&m.IPAddress, &m.UserID, &m.Description, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan IP mapping: %w", err)
		}
		mappings = append(mappings, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating IP mappings: %w", err)
	}

	return mappings, nil
}

// DeleteIPMapping removes an IP mapping
func (db *DB) DeleteIPMapping(ctx context.Context, ipAddress string) error {
	query := `DELETE FROM ip_mappings WHERE ip_address = ?`

	_, err := db.ExecContext(ctx, query, ipAddress)
	if err != nil {
		return fmt.Errorf("failed to delete IP mapping: %w", err)
	}

	return nil
}

// IPMappingExists checks if an IP mapping exists
func (db *DB) IPMappingExists(ctx context.Context, ipAddress string) (bool, error) {
	query := `SELECT 1 FROM ip_mappings WHERE ip_address = ?`

	var exists int
	err := db.QueryRowContext(ctx, query, ipAddress).Scan(&exists)
	if err != nil {
		return false, nil // Not found is not an error
	}

	return true, nil
}

// NewIPMapping creates a new IP mapping with timestamps
func NewIPMapping(ipAddress, userID, description string) *IPMapping {
	return &IPMapping{
		IPAddress:   ipAddress,
		UserID:      userID,
		Description: description,
		CreatedAt:   time.Now().Unix(),
	}
}
