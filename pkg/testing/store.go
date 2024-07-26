package testing

import (
	"database/sql"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun/driver/pgdriver"
)

const (
	localTestDBDSNPrefix = "postgres://postgres:xmtp@localhost:15432"
	localTestDBDSNSuffix = "?sslmode=disable"
)

func NewDB(t testing.TB) (*sql.DB, string, func()) {
	dsn := localTestDBDSNPrefix + localTestDBDSNSuffix
	ctlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	dbName := "test_" + RandomStringLower(12)
	_, err := ctlDB.Exec("CREATE DATABASE " + dbName)
	require.NoError(t, err)

	dsn = localTestDBDSNPrefix + "/" + dbName + localTestDBDSNSuffix
	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	return db, dsn, func() {
		db.Close()
		_, err = ctlDB.Exec("DROP DATABASE " + dbName)
		require.NoError(t, err)
		ctlDB.Close()
	}
}

func NewPGXDB(t testing.TB) (*sql.DB, string, func()) {
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
		db.Close()
		_, err = ctlDB.Exec("DROP DATABASE " + dbName)
		require.NoError(t, err)
		ctlDB.Close()
	}
}
