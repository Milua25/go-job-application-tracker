package main

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var fs embed.FS

func applyDBMigrations(dsn string) error {
	slog.Info("applying database migrations")

	// Dedicated connection for migrations — m.Close() will close this, not the app's pool.
	migrationDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open migration database connection: %w", err)
	}

	driver, err := migratepg.WithInstance(migrationDB, &migratepg.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	d, err := iofs.New(fs, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create iofs driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("failed to get migration version: %w", err)
	}
	if dirty {
		return fmt.Errorf("database is in a dirty state at version %d", version)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration up failed: %w", err)
	}

	newVersion, _, _ := m.Version()
	slog.Info("migrations applied", "from_version", version, "to_version", newVersion)

	return nil
}
