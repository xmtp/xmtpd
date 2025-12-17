package testutils

import (
	"context"
	"database/sql"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/xmtp/xmtpd/pkg/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/migrations"
)

const (
	LocalTestDBDSNPrefix = "postgres://postgres:xmtp@localhost:8765"
	LocalTestDBDSNSuffix = "?sslmode=disable"
)

func GetCallerName(depth int) string {
	pc, _, _, ok := runtime.Caller(depth)
	if !ok {
		return "unknown"
	}
	nameFull := runtime.FuncForPC(pc).Name()
	nameEnd := filepath.Ext(nameFull)
	name := strings.TrimPrefix(nameEnd, ".")
	return strings.ToLower(name)
}

func openDB(t testing.TB, dsn string) (*sql.DB, string) {
	config, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	dbInstance := stdlib.OpenDB(*config)
	t.Cleanup(func() {
		require.NoError(t, dbInstance.Close())
	})
	return dbInstance, dsn
}

func newCtlDB(t testing.TB) (*sql.DB, string) {
	return openDB(t, LocalTestDBDSNPrefix+LocalTestDBDSNSuffix)
}

func newInstanceDB(t testing.TB, ctx context.Context, ctlDB *sql.DB) (*sql.DB, string) {
	dbName := "test_" + GetCallerName(3) + "_" + RandomStringLower(12)
	t.Logf("creating database %s ...", dbName)
	_, err := ctlDB.Exec("CREATE DATABASE " + dbName)
	require.NoError(t, err)

	t.Cleanup(func() {
		_, err := ctlDB.Exec("DROP DATABASE " + dbName)
		require.NoError(t, err)
	})

	dbInstance, dsn := openDB(t, LocalTestDBDSNPrefix+"/"+dbName+LocalTestDBDSNSuffix)
	require.NoError(t, migrations.Migrate(ctx, dbInstance))

	return dbInstance, dsn
}

func NewRawDB(t *testing.T, ctx context.Context) (*sql.DB, string) {
	ctlDB, _ := newCtlDB(t)
	dbInstance, dsn := newInstanceDB(t, ctx, ctlDB)

	return dbInstance, dsn
}

func NewDB(t *testing.T, ctx context.Context) (*db.Handler, string) {
	t.Helper()

	dbh, dsn := NewRawDB(t, ctx)
	return db.NewDBHandler(dbh), dsn
}

func NewDBs(t *testing.T, ctx context.Context, count int) []*sql.DB {
	ctlDB, _ := newCtlDB(t)
	dbs := []*sql.DB{}

	for i := 0; i < count; i++ {
		dbInstance, _ := newInstanceDB(t, ctx, ctlDB)
		dbs = append(dbs, dbInstance)
	}

	return dbs
}

func InsertGatewayEnvelopes(
	t *testing.T,
	dbInstance *sql.DB,
	rows []queries.InsertGatewayEnvelopeParams,
	notifyChan ...chan bool,
) {
	ctx := t.Context()
	q := queries.New(dbInstance)
	for _, row := range rows {
		inserted, err := db.InsertGatewayEnvelopeWithChecksStandalone(ctx, q, row)
		require.NoError(t, err)
		require.EqualValues(t, int64(1), inserted.InsertedMetaRows)

		if len(notifyChan) > 0 {
			select {
			case notifyChan[0] <- true:
			default:
			}
		}
	}
}

func CreatePayer(t *testing.T, db *sql.DB, address ...string) int32 {
	q := queries.New(db)
	var payerAddress string
	if len(address) > 0 {
		payerAddress = address[0]
	} else {
		payerAddress = RandomString(42)
	}

	id, err := q.FindOrCreatePayer(context.Background(), payerAddress)
	require.NoError(t, err)

	return id
}
