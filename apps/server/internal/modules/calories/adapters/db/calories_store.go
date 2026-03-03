package db

import "database/sql"

type CaloriesStore struct {
	db *sql.DB
}

func NewCaloriesStore(db *sql.DB) *CaloriesStore {
	return &CaloriesStore{db: db}
}
