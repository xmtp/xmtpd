package server

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func createTestLogger() (*zap.Logger, *observer.ObservedLogs) {
	core, observedLogs := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)
	return logger, observedLogs
}

type mockConnectRequestLoggingInterceptor struct {
	connect.AnyRequest
	procedure string
}

func (m *mockConnectRequestLoggingInterceptor) Spec() connect.Spec {
	return connect.Spec{Procedure: m.procedure}
}

type mockStreamingConnLoggingInterceptor struct {
	connect.StreamingHandlerConn
	procedure string
}

func (m *mockStreamingConnLoggingInterceptor) Spec() connect.Spec {
	return connect.Spec{Procedure: m.procedure}
}

func TestUnaryLoggingInterceptor(t *testing.T) {
	logger, logs := createTestLogger()

	interceptor, err := NewLoggingInterceptor(logger)
	require.NoError(t, err)

	next := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		return nil, status.Errorf(codes.Internal, "mock internal error")
	}

	ctx := context.Background()
	req := &mockConnectRequestLoggingInterceptor{procedure: "/test.TestService/TestMethod"}

	wrappedUnary := interceptor.WrapUnary(next)
	_, err = wrappedUnary(ctx, req)

	require.Error(t, err)
	require.Equal(t, 1, logs.Len(), "expected one log entry but got none")

	logEntry := logs.All()[0]

	require.Equal(t, zapcore.ErrorLevel, logEntry.Level, "expected log level 'Error'")
	require.Contains(t, logEntry.ContextMap(), "method")
	require.Equal(
		t,
		"/test.TestService/TestMethod",
		logEntry.ContextMap()["method"],
		"expected log to contain correct method",
	)
	require.Contains(t, logEntry.ContextMap(), "message")
	require.Contains(
		t,
		logEntry.ContextMap()["message"],
		"mock internal error",
		"expected log to contain correct error message",
	)
}

func TestStreamLoggingInterceptor(t *testing.T) {
	logger, logs := createTestLogger()
	interceptor, err := NewLoggingInterceptor(logger)
	require.NoError(t, err)

	next := func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		return status.Errorf(codes.NotFound, "mock stream error")
	}

	ctx := context.Background()
	conn := &mockStreamingConnLoggingInterceptor{procedure: "/test.TestService/TestStream"}

	wrappedStream := interceptor.WrapStreamingHandler(next)
	err = wrappedStream(ctx, conn)

	require.Error(t, err)
	require.Equal(t, 1, logs.Len(), "expected one log entry but got none")

	logEntry := logs.All()[0]

	require.Equal(t, zapcore.ErrorLevel, logEntry.Level, "expected log level 'Error'")
	require.Contains(t, logEntry.ContextMap(), "method")
	require.Equal(
		t,
		"/test.TestService/TestStream",
		logEntry.ContextMap()["method"],
		"expected log to contain correct method",
	)
	require.Contains(t, logEntry.ContextMap(), "message")
	require.Contains(
		t,
		logEntry.ContextMap()["message"],
		"mock stream error",
		"expected log to contain correct error message",
	)
}
