package postgresstore

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type DB struct {
	*sql.DB
	DSN string
}

func NewDB(opts *Options) (*DB, error) {
	db, err := otelsql.Open("pgx", opts.DSN,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithDBName("xmtpd"),
	)
	if err != nil {
		return nil, err
	}
	return &DB{
		DB:  db,
		DSN: opts.DSN,
	}, nil
}
