// Package postgres interacts with a postgres database
package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andrew-hayworth22/critiquefi-service/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type DBConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// NewDB creates a new postgres database connection.
func NewDB(ctx context.Context, cfg DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "pgx", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error pinging postgres: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return db, nil
}

// MapError maps postgres errors to our storage layer errors.
// Falls back to the original error if not mapped.
func MapError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return store.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return store.ErrDuplicate
		case "23503":
			return store.ErrForeignKeyViolation
		case "23502":
			return store.ErrNotNullViolation
		}
	}

	return err
}

// CloseRows closes a row object and handles any errors.
// Use this instead of deferring rows.Close() directly.
func CloseRows(rows *sqlx.Rows) {
	if err := rows.Close(); err != nil {
		return
	}
}
