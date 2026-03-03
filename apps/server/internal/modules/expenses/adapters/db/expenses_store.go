package db

import (
	"database/sql"
)

type ExpensesStore struct {
	db *sql.DB
}

func NewExpensesStore(db *sql.DB) *ExpensesStore {
	return &ExpensesStore{db: db}
}
