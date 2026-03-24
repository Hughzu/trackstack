package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func Open(url, token string, pool PoolConfig) (*sql.DB, error) {
	dsn := url
	if token != "" && token != "local" {
		dsn = fmt.Sprintf("%s?authToken=%s", url, token)
	}

	db, err := sql.Open("libsql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if pool.MaxOpenConns > 0 {
		db.SetMaxOpenConns(pool.MaxOpenConns)
	}
	if pool.MaxIdleConns > 0 {
		db.SetMaxIdleConns(pool.MaxIdleConns)
	}
	if pool.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(pool.ConnMaxLifetime)
	}
	if pool.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(pool.ConnMaxIdleTime)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return db, nil
}
