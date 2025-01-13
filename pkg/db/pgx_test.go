package db

import (
	"context"
	"fmt"
	"log"
	"net"
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

func BlackHoleServer(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error starting mock server: %v", err)
	}
	defer ln.Close()

	fmt.Println("Mock server running on port", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		// Simulate "black hole" by keeping the connection open without any response.
		go func(c net.Conn) {
			defer c.Close()
			select {}
		}(conn)
	}
}

func TestBlackholeDNS(t *testing.T) {
	port := "5433"
	dsn := fmt.Sprintf("postgres://user:password@localhost:%s/dbname?sslmode=disable", port)
	// Start the mock server in a goroutine
	go BlackHoleServer(port)

	// Ensure the test doesn't run indefinitely
	testCtx, cancelTest := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelTest()

	_, err := NewNamespacedDB(
		testCtx,
		dsn,
		"dbname",
		200*time.Millisecond,
		200*time.Millisecond,
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "database is not ready")
}
