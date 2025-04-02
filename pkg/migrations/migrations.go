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

	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
