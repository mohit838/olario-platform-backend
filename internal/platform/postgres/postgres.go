package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mohit838/olario-platform-backend/internal/config"
)

// Open creates a Postgres connection pool and verifies it with Ping.
// The pool is returned as *sql.DB because database/sql is stable and works well
// with golang-migrate and repository adapters.
func Open(ctx context.Context, cfg config.DatabaseConfig) (*sql.DB, error) {
	if cfg.DSN == "" {
		return nil, fmt.Errorf("database dsn is required")
	}
	if strings.Contains(cfg.DSN, "<") || strings.Contains(cfg.DSN, ">") {
		return nil, fmt.Errorf("database dsn still contains placeholder values")
	}

	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.MaxIdleTime)

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return db, nil
}
