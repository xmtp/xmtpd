package postgresstore

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.nhat.io/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type DB struct {
	*sql.DB
	DSN string
}

func NewDB(opts *Options) (*DB, error) {
	driverName, err := otelsql.Register("pgx",
		otelsql.AllowRoot(),
		otelsql.TraceQueryWithoutArgs(),
		otelsql.TraceRowsClose(),
		otelsql.TraceRowsAffected(),
		otelsql.WithSystem(semconv.DBSystemPostgreSQL),
		otelsql.WithDatabaseName("xmtpd"),
	)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open(driverName, opts.DSN)
	if err != nil {
		return nil, err
	}
	return &DB{
		DB:  db,
		DSN: opts.DSN,
	}, nil
}
