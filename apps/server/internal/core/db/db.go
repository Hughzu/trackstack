package db

import (
	"database/sql"
	"fmt"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func OpenLibSQL(dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("missing database dsn")
	}

	db, err := sql.Open("libsql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open libsql db: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping libsql db: %w", err)
	}

	return db, nil
}
