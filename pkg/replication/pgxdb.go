package replication

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func NewPGXDB(dsn string, waitForDB, statementTimeout time.Duration) (*sql.DB, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	config.ConnConfig.RuntimeParams["statement_timeout"] = fmt.Sprint(statementTimeout.Milliseconds())

	dbpool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	db := stdlib.OpenDBFromPool(dbpool)

	waitUntil := time.Now().Add(waitForDB)

	err = db.Ping()
	for err != nil && time.Now().Before(waitUntil) {
		time.Sleep(3 * time.Second)
		err = db.Ping()
	}

	return db, nil
}

func NewBunPGXDb(dsn string, waitForDB, statementTimeout time.Duration) (*bun.DB, error) {
	pgxDb, err := NewPGXDB(dsn, waitForDB, statementTimeout)
	if err != nil {
		return nil, err
	}

	return bun.NewDB(pgxDb, pgdialect.New()), nil
}
