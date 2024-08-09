package migrations

import (
	"context"
	"database/sql"
	"embed"

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
	defer migrationFs.Close()

	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	driver, err := postgres.WithConnection(ctx, conn, &postgres.Config{})
	if err != nil {
		return err
	}
	defer driver.Close()

	migrator, err := migrate.NewWithInstance("iofs", migrationFs, "postgres", driver)
	if err != nil {
		return err
	}
	defer migrator.Close()

	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
