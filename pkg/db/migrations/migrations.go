package migrations

import (
	"context"
	"database/sql"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var migrationFs embed.FS

func Migrate(ctx context.Context, db *sql.DB) error {
	return withMigrator(ctx, db, func(m *migrate.Migrate) error {
		return m.Up()
	})
}

// MigrateTo migrates the database to the specified schema version. It can move
// the schema either up or down depending on the current state. Useful for
// testing migration behavior against pre-existing data.
func MigrateTo(ctx context.Context, db *sql.DB, version uint) error {
	return withMigrator(ctx, db, func(m *migrate.Migrate) error {
		return m.Migrate(version)
	})
}

func withMigrator(
	ctx context.Context,
	db *sql.DB,
	fn func(*migrate.Migrate) error,
) error {
	migrationFs, err := iofs.New(migrationFs, ".")
	if err != nil {
		return err
	}
	defer func() {
		_ = migrationFs.Close()
	}()

	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	driver, err := postgres.WithConnection(ctx, conn, &postgres.Config{})
	if err != nil {
		return err
	}
	defer func() {
		_ = driver.Close()
	}()

	migrator, err := migrate.NewWithInstance("iofs", migrationFs, "postgres", driver)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = migrator.Close()
	}()

	if err := fn(migrator); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
