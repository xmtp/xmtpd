package migrations

import (
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//go:embed *.sql
var fs embed.FS

func Run(dsn string) error {
	m, err := new(dsn)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func new(dsn string) (*migrate.Migrate, error) {
	var err error
	driver, err := iofs.New(fs, ".")
	if err != nil {
		return nil, err
	}
	return migrate.NewWithSourceInstance("iofs", driver, dsn)
}
