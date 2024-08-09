package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/xmtp/xmtpd/pkg/migrations"
)

func newPGXDB(
	ctx context.Context,
	dsn string,
	waitForDB, statementTimeout time.Duration,
) (*sql.DB, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	config.ConnConfig.RuntimeParams["statement_timeout"] = fmt.Sprint(
		statementTimeout.Milliseconds(),
	)

	dbPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	db := stdlib.OpenDBFromPool(dbPool)

	waitUntil := time.Now().Add(waitForDB)

	err = db.Ping()
	for err != nil && time.Now().Before(waitUntil) {
		time.Sleep(3 * time.Second)
		err = db.Ping()
	}

	return db, err
}

func NewDB(
	ctx context.Context,
	dsn string,
	waitForDB, statementTimeout time.Duration,
) (*sql.DB, error) {
	db, err := newPGXDB(ctx, dsn, waitForDB, statementTimeout)
	if err != nil {
		return nil, err
	}
	err = migrations.Migrate(ctx, db)
	if err != nil {
		return nil, err
	}
	return db, nil
}
