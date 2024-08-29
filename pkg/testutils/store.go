package testutils

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/migrations"
)

const (
	localTestDBDSNPrefix = "postgres://postgres:xmtp@localhost:8765"
	localTestDBDSNSuffix = "?sslmode=disable"
)

func openDB(t testing.TB, dsn string) (*sql.DB, string, func()) {
	config, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	db := stdlib.OpenDB(*config)
	return db, dsn, func() {
		err := db.Close()
		require.NoError(t, err)
	}
}

func newCtlDB(t testing.TB) (*sql.DB, string, func()) {
	return openDB(t, localTestDBDSNPrefix+localTestDBDSNSuffix)
}

func newInstanceDB(t testing.TB, ctx context.Context, ctlDB *sql.DB) (*sql.DB, string, func()) {
	dbName := "test_" + RandomStringLower(12)
	_, err := ctlDB.Exec("CREATE DATABASE " + dbName)
	require.NoError(t, err)

	db, dsn, cleanup := openDB(t, localTestDBDSNPrefix+"/"+dbName+localTestDBDSNSuffix)
	require.NoError(t, migrations.Migrate(ctx, db))

	return db, dsn, func() {
		cleanup()
		_, err = ctlDB.Exec("DROP DATABASE " + dbName)
		require.NoError(t, err)
	}
}

func NewDB(t *testing.T, ctx context.Context) (*sql.DB, string, func()) {
	ctlDB, _, ctlCleanup := newCtlDB(t)
	db, dsn, cleanup := newInstanceDB(t, ctx, ctlDB)

	return db, dsn, func() {
		cleanup()
		ctlCleanup()
	}
}

func NewDBs(t *testing.T, ctx context.Context, count int) ([]*sql.DB, func()) {
	ctlDB, _, ctlCleanup := newCtlDB(t)
	dbs := []*sql.DB{}
	cleanups := []func(){}

	for i := 0; i < count; i++ {
		db, _, cleanup := newInstanceDB(t, ctx, ctlDB)
		dbs = append(dbs, db)
		cleanups = append(cleanups, cleanup)
	}

	return dbs, func() {
		for i := 0; i < count; i++ {
			cleanups[i]()
		}
		ctlCleanup()
	}
}
