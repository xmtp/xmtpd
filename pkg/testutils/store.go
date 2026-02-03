package testutils

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/db/vectorclock"
	"github.com/xmtp/xmtpd/pkg/migrations"
)

const (
	LocalTestDBDSNPrefix = "postgres://postgres:xmtp@localhost:8765"
	LocalTestDBDSNSuffix = "?sslmode=disable"

	envRunVectorClockCheck = "XMTP_RUN_VECTOR_CLOCK_CHECK"
)

var runVectorClockCheck = parseBoolConfig(os.Getenv(envRunVectorClockCheck))

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

	return NewDBWithLogger(t, ctx, zap.NewNop())
}

func NewDBWithLogger(t *testing.T, ctx context.Context, log *zap.Logger) (*db.Handler, string) {
	t.Helper()

	dbh, dsn := NewRawDB(t, ctx)

	readFunc := db.GetVectorClockReader(dbh)
	vc := vectorclock.New(log, readFunc)

	if runVectorClockCheck {
		// This is an opt-in vector clock integrity check ran on test shutdown.
		// Downside is that this will need to establish a new db connection and do one additional query per test which is non-negligent.
		t.Cleanup(func() {
			// Since t.Cleanup() is called AFTER t.Context() is cancelled,
			// by now our DB is already closed.
			// This is overhead, but lets create a new connection to the DB while it's alive and
			// run our sanity check.

			newConn, _ := openDB(t, dsn)
			defer func() {
				_ = newConn.Close()
			}()

			dbState, err := db.GetVectorClockReader(newConn)(context.Background())
			require.NoError(t, err)
			require.EqualValuesf(t, dbState, vc.Values(), "vector clock does not match DB state")
		})
	}

	return db.NewDBHandler(dbh, vc), dsn
}

func NewDBs(t *testing.T, ctx context.Context, count int) []*db.Handler {
	out := make([]*db.Handler, count)
	for i := range count {
		db, _ := NewDB(t, ctx)
		out[i] = db
	}

	return out
}

func NewRawDBs(t *testing.T, ctx context.Context, count int) []*sql.DB {
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
	dbh *db.Handler,
	rows []queries.InsertGatewayEnvelopeParams,
	notifyChan ...chan bool,
) {
	ctx := t.Context()
	for _, row := range rows {
		inserted, err := db.InsertGatewayEnvelopeWithChecksStandalone(ctx, dbh, row)
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

func CreatePayer(t *testing.T, db *db.Handler, address ...string) int32 {
	var payerAddress string
	if len(address) > 0 {
		payerAddress = address[0]
	} else {
		payerAddress = RandomString(42)
	}

	id, err := db.Query().FindOrCreatePayer(context.Background(), payerAddress)
	require.NoError(t, err)

	return id
}
