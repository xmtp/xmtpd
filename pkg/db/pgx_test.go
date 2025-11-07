package db_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/xmtp/xmtpd/pkg/db"
	"go.uber.org/zap"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func TestNamespacedDB(t *testing.T) {
	startingDsn := testutils.LocalTestDBDSNPrefix + "/foo" + testutils.LocalTestDBDSNSuffix
	newDBName := "xmtp_" + testutils.RandomString(24)
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	// Create namespaced DB
	namespacedDB, err := db.NewNamespacedDB(
		context.Background(),
		logger,
		startingDsn,
		newDBName,
		time.Second,
		time.Second,
		nil,
	)
	t.Cleanup(func() { _ = namespacedDB.Close() })
	require.NoError(t, err)

	result, err := namespacedDB.Query("SELECT current_database();")
	require.NoError(t, err)
	defer func() {
		_ = result.Close()
	}()

	require.True(t, result.Next())
	var dbName string
	err = result.Scan(&dbName)
	require.NoError(t, err)
	require.Equal(t, newDBName, dbName)
}

func TestNamespaceRepeat(t *testing.T) {
	startingDsn := testutils.LocalTestDBDSNPrefix + "/foo" + testutils.LocalTestDBDSNSuffix
	newDBName := "xmtp_" + testutils.RandomString(24)
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	// Create namespaced DB
	db1, err := db.NewNamespacedDB(
		context.Background(),
		logger,
		startingDsn,
		newDBName,
		time.Second,
		time.Second,
		nil,
	)
	require.NoError(t, err)
	require.NotNil(t, db1)
	t.Cleanup(func() { _ = db1.Close() })

	// Create again with the same name
	db2, err := db.NewNamespacedDB(
		context.Background(),
		logger,
		startingDsn,
		newDBName,
		time.Second,
		time.Second,
		nil,
	)
	require.NoError(t, err)
	require.NotNil(t, db2)
	t.Cleanup(func() { _ = db2.Close() })
}

func TestNamespacedDBInvalidName(t *testing.T) {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	_, err = db.NewNamespacedDB(
		context.Background(),
		logger,
		testutils.LocalTestDBDSNPrefix+"/foo"+testutils.LocalTestDBDSNSuffix,
		"invalid/name",
		time.Second,
		time.Second,
		nil,
	)
	require.Error(t, err)
}

func TestNamespacedDBInvalidDSN(t *testing.T) {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	_, err = db.NewNamespacedDB(
		context.Background(),
		logger,
		"invalid-dsn",
		"dbname",
		time.Second,
		time.Second,
		nil,
	)
	require.Error(t, err)
}

func BlackHoleServer(ctx context.Context, port string) error {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("error starting blackhole server: %w", err)
	}
	defer func() {
		_ = ln.Close()
	}()

	go func() {
		<-ctx.Done()
		_ = ln.Close()
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
			defer func() {
				_ = c.Close()
			}()
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
	_ = listener.Close()

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

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

	_, err = db.NewNamespacedDB(
		testCtx,
		logger,
		dsn,
		"dbname",
		dbTimeout,
		5*time.Second,
		nil,
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "database is not ready")

	require.NoError(t, testCtx.Err(), "Test timed out")

	// Cleanup server
	cancelServer()
	require.NoError(t, <-serverErrCh)
}
