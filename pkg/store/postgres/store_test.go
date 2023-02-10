package postgresstore_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	crdttest "github.com/xmtp/xmtpd/pkg/crdt/testing"
	postgresstore "github.com/xmtp/xmtpd/pkg/store/postgres"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

const (
	localTestDBDSNSuffix = "?sslmode=disable"
)

var localTestDBDSNPrefix string

func init() {
	dir, err := findProjectDir()
	if err != nil {
		panic(err)
	}

	err = godotenv.Load(filepath.Join(dir, ".env.local"))
	if err != nil {
		panic(err)
	}

	localTestDBDSNPrefix = fmt.Sprintf("postgres://%s:%s@%s:%s", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"))
}

func TestScopedPostgresStore(t *testing.T) {
	ctx := context.Background()
	log := test.NewLogger(t)
	topic := "topic-" + test.RandomStringLower(13)

	t.Run("event", func(t *testing.T) {
		t.Parallel()
		crdttest.RunStoreEventTests(t, topic, func(t *testing.T) *crdttest.TestStore {
			store := newTestScopedStore(t, topic)
			return crdttest.NewTestStore(ctx, log, store)
		})
	})

	t.Run("query", func(t *testing.T) {
		t.Parallel()
		crdttest.RunStoreQueryTests(t, topic, func(t *testing.T) *crdttest.TestStore {
			store := newTestScopedStore(t, topic)
			return crdttest.NewTestStore(ctx, log, store)
		})
	})
}

type testScopedStore struct {
	*postgresstore.Store
	db        *postgresstore.DB
	dbCleanup func()
}

func newTestScopedStore(t *testing.T, topic string) *testScopedStore {
	t.Helper()
	log := test.NewLogger(t)

	db, cleanup := newTestDB(t)
	store, err := postgresstore.NewNodeStore(log, db)
	require.NoError(t, err)
	scopedStore, err := store.NewTopic(topic)
	require.NoError(t, err)

	return &testScopedStore{
		Store:     scopedStore.(*postgresstore.Store),
		db:        db,
		dbCleanup: cleanup,
	}
}

func (s *testScopedStore) Close() error {
	s.dbCleanup()
	return s.Store.Close()
}

func newTestDB(t *testing.T) (*postgresstore.DB, func()) {
	t.Helper()
	dsn := localTestDBDSNPrefix + localTestDBDSNSuffix
	ctlDB, err := postgresstore.NewDB(dsn)
	require.NoError(t, err)

	dbName := "test_" + test.RandomStringLower(13)
	_, err = ctlDB.Exec("CREATE DATABASE " + dbName)
	require.NoError(t, err)

	dsn = localTestDBDSNPrefix + "/" + dbName + localTestDBDSNSuffix
	db, err := postgresstore.NewDB(dsn)
	require.NoError(t, err)

	return db, func() {
		db.Close()
		_, err = ctlDB.Exec(fmt.Sprintf("REVOKE CONNECT ON DATABASE %s FROM public", dbName))
		require.NoError(t, err)
		_, err = ctlDB.Exec(fmt.Sprintf("SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = '%s'", dbName))
		require.NoError(t, err)
		_, err = ctlDB.Exec("DROP DATABASE " + dbName)
		require.NoError(t, err)
		ctlDB.Close()
	}
}

func findProjectDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for dir != "/" {
		dir = filepath.Dir(dir)
		info, err := os.Stat(filepath.Join(dir, "go.mod"))
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", err
		}
		if info.IsDir() {
			continue
		}
		break
	}
	return dir, nil
}
