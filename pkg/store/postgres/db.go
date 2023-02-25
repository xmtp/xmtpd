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

func NewDB(dsn string) (*DB, error) {
	db, err := otelsql.Open("pgx", dsn,
		otelsql.WithAttributes(semconv.DBSystemSqlite),
		otelsql.WithDBName("mydb"),
	)
	if err != nil {
		return nil, err
	}
	return &DB{
		DB:  db,
		DSN: dsn,
	}, nil
}
