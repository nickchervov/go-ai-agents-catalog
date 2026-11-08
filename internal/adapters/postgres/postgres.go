package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var migrationPath = "pkg/migrations"

type PostgresRepository struct {
	db *sqlx.DB
}

func migrateUP(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		db.Close()
		return fmt.Errorf("creating db driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationPath, "pgx", driver)
	if err != nil {
		db.Close()
		return fmt.Errorf("creating migrate: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		db.Close()
		return fmt.Errorf("up migrate: %w", err)
	}
	return nil
}

func NewPostgresRepository(ctx context.Context, dsn string) (*PostgresRepository, error) {
	db, err := sqlx.ConnectContext(ctx, "pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := migrateUP(db); err != nil {
		return nil, fmt.Errorf("up migrate: %w", err)
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) Close() {
	r.db.Close()
}
