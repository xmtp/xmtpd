package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func TestNamespacedDB(t *testing.T) {
	startingDsn := testutils.LocalTestDBDSNPrefix + "/foo" + testutils.LocalTestDBDSNSuffix
	newDBName := "xmtp_" + testutils.RandomString(24)
	// Create namespaced DB
	namespacedDB, err := NewNamespacedDB(
		context.Background(),
		startingDsn,
		newDBName,
		time.Second,
		time.Second,
	)
	t.Cleanup(func() { namespacedDB.Close() })
	require.NoError(t, err)

	result, err := namespacedDB.Query("SELECT current_database();")
	require.NoError(t, err)
	defer result.Close()

	require.True(t, result.Next())
	var dbName string
	err = result.Scan(&dbName)
	require.NoError(t, err)
	require.Equal(t, newDBName, dbName)
}

func TestNamespaceRepeat(t *testing.T) {
	startingDsn := testutils.LocalTestDBDSNPrefix + "/foo" + testutils.LocalTestDBDSNSuffix
	newDBName := "xmtp_" + testutils.RandomString(24)
	// Create namespaced DB
	db1, err := NewNamespacedDB(
		context.Background(),
		startingDsn,
		newDBName,
		time.Second,
		time.Second,
	)
	require.NoError(t, err)
	require.NotNil(t, db1)
	t.Cleanup(func() { db1.Close() })

	// Create again with the same name
	db2, err := NewNamespacedDB(
		context.Background(),
		startingDsn,
		newDBName,
		time.Second,
		time.Second,
	)
	require.NoError(t, err)
	require.NotNil(t, db2)
	t.Cleanup(func() { db2.Close() })
}

func TestNamespacedDBInvalidName(t *testing.T) {
	_, err := NewNamespacedDB(
		context.Background(),
		testutils.LocalTestDBDSNPrefix+"/foo"+testutils.LocalTestDBDSNSuffix,
		"invalid/name",
		time.Second,
		time.Second,
	)
	require.Error(t, err)
}

func TestNamespacedDBInvalidDSN(t *testing.T) {
	_, err := NewNamespacedDB(
		context.Background(),
		"invalid-dsn",
		"dbname",
		time.Second,
		time.Second,
	)
	require.Error(t, err)
}
