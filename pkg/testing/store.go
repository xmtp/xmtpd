package testing

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

func newPGXDB(t testing.TB) (*sql.DB, string, func()) {
	dsn := localTestDBDSNPrefix + localTestDBDSNSuffix
	config, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	ctlDB := stdlib.OpenDB(*config)
	dbName := "test_" + RandomStringLower(12)
	_, err = ctlDB.Exec("CREATE DATABASE " + dbName)
	require.NoError(t, err)

	dsn = localTestDBDSNPrefix + "/" + dbName + localTestDBDSNSuffix
	config2, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	db := stdlib.OpenDB(*config2)
	return db, dsn, func() {
		err := db.Close()
		require.NoError(t, err)
		_, err = ctlDB.Exec("DROP DATABASE " + dbName)
		require.NoError(t, err)
		ctlDB.Close()
	}
}

func NewDB(t *testing.T, ctx context.Context) (*sql.DB, string, func()) {
	db, dsn, cleanup := newPGXDB(t)
	require.NoError(t, migrations.Migrate(ctx, db))

	return db, dsn, cleanup
}
