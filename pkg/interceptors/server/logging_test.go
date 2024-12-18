package server

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func createTestLogger() (*zap.Logger, *observer.ObservedLogs) {
	core, observedLogs := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)
	return logger, observedLogs
}

type mockServerStream struct {
	grpc.ServerStream
}

func (m *mockServerStream) Context() context.Context {
	return context.Background()
}

func TestUnaryLoggingInterceptor(t *testing.T) {
	logger, logs := createTestLogger()

	interceptor, err := NewLoggingInterceptor(logger)
	require.NoError(t, err)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, status.Errorf(codes.Internal, "mock internal error")
	}

	ctx := context.Background()
	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.TestService/TestMethod",
	}
	req := struct{}{}

	interceptorUnary := interceptor.Unary()
	_, err = interceptorUnary(ctx, req, info, handler)

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
	require.Equal(
		t,
		"mock internal error",
		logEntry.ContextMap()["message"],
		"expected log to contain correct error message",
	)
}
func TestStreamLoggingInterceptor(t *testing.T) {
	logger, logs := createTestLogger()
	interceptor, err := NewLoggingInterceptor(logger)
	require.NoError(t, err)

	handler := func(srv interface{}, ss grpc.ServerStream) error {
		return status.Errorf(codes.NotFound, "mock stream error")
	}

	info := &grpc.StreamServerInfo{
		FullMethod: "/test.TestService/TestStream",
	}

	stream := &mockServerStream{}

	incerceptorStream := interceptor.Stream()
	err = incerceptorStream(nil, stream, info, handler)

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
	require.Equal(
		t,
		"mock stream error",
		logEntry.ContextMap()["message"],
		"expected log to contain correct error message",
	)
}
