package postgresstore

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	*sql.DB
	DSN string
}

func NewDB(opts *Options) (*DB, error) {
	db, err := sql.Open("pgx", opts.DSN)
	if err != nil {
		return nil, err
	}
	return &DB{
		DB:  db,
		DSN: opts.DSN,
	}, nil
}
