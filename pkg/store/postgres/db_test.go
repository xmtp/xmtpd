package postgresstore_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
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

func newTestDB(t testing.TB) (*postgresstore.DB, func()) {
	t.Helper()
	opts := &postgresstore.Options{
		DSN: localTestDBDSNPrefix + localTestDBDSNSuffix,
	}
	ctlDB, err := postgresstore.NewDB(opts)
	require.NoError(t, err)

	dbName := "test_" + test.RandomStringLower(13)
	_, err = ctlDB.Exec("CREATE DATABASE " + dbName)
	require.NoError(t, err)

	opts2 := &postgresstore.Options{
		DSN: localTestDBDSNPrefix + "/" + dbName + localTestDBDSNSuffix,
	}
	db, err := postgresstore.NewDB(opts2)
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
