package db

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func nullableString(value *string) any {
	if value == nil {
		return nil
	}

	return *value
}
