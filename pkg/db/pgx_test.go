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

func BlackHoleServer(ctx context.Context, port string) error {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("error starting blackhole server: %w", err)
	}
	defer ln.Close()

	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	for ctx.Err() == nil {
		conn, err := ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		// Simulate "black hole" by keeping the connection open without any response.
		go func(c net.Conn) {
			defer c.Close()
			<-ctx.Done()
		}(conn)
	}

	return nil
}

func TestBlackholeDNS(t *testing.T) {
	// Find available port
	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	dsn := fmt.Sprintf("postgres://user:password@localhost:%d/dbname?sslmode=disable", port)
	const testTimeout = 5 * time.Second
	const dbTimeout = 200 * time.Millisecond

	testCtx, cancelTest := context.WithTimeout(context.Background(), testTimeout)
	defer cancelTest()

	// Start server with context
	serverCtx, cancelServer := context.WithCancel(testCtx)
	defer cancelServer()

	serverErrCh := make(chan error, 1)
	go func() {
		serverErrCh <- BlackHoleServer(serverCtx, fmt.Sprintf("%d", port))
	}()
	// Wait for server to start
	time.Sleep(50 * time.Millisecond)

	_, err = NewNamespacedDB(
		testCtx,
		dsn,
		"dbname",
		dbTimeout,
		5*time.Second,
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "database is not ready")

	require.NoError(t, testCtx.Err(), "Test timed out")

	// Cleanup server
	cancelServer()
	require.NoError(t, <-serverErrCh)

}
