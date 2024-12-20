package server

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"net"
	"testing"
)

func mockUnaryHandler(ctx context.Context, req interface{}) (interface{}, error) {
	return "response", nil
}

func mockStreamHandler(srv interface{}, ss grpc.ServerStream) error {
	return nil
}

type mockServerStreamWithContext struct {
	grpc.ServerStream
	ctx context.Context
}

func (m *mockServerStreamWithContext) Context() context.Context {
	return m.ctx
}

func TestIncomingInterceptor_Unary(t *testing.T) {
	logger, logs := createTestLogger()

	interceptor, err := NewIncomingInterceptor(logger)
	if err != nil {
		t.Fatalf("failed to create interceptor: %v", err)
	}

	ctx := peer.NewContext(context.Background(), &peer.Peer{
		Addr: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345},
	})

	// Call the unary interceptor
	_, _ = interceptor.Unary()(ctx, nil, nil, mockUnaryHandler)

	require.NoError(t, err)
	require.Equal(t, 1, logs.Len(), "expected one log entry but got none")

	logEntry := logs.All()[0]

	require.Equal(t, zapcore.DebugLevel, logEntry.Level)
	require.Contains(t, logEntry.Message, "Incoming request")

}

func TestIncomingInterceptor_Stream(t *testing.T) {
	logger, logs := createTestLogger()

	interceptor, err := NewIncomingInterceptor(logger)
	if err != nil {
		t.Fatalf("failed to create interceptor: %v", err)
	}

	ctx := peer.NewContext(context.Background(), &peer.Peer{
		Addr: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345},
	})

	// Create a mock server stream
	stream := &mockServerStreamWithContext{ctx: ctx}

	// Call the stream interceptor
	_ = interceptor.Stream()(nil, stream, nil, mockStreamHandler)

	require.NoError(t, err)
	require.Equal(t, 1, logs.Len(), "expected one log entry but got none")

	logEntry := logs.All()[0]

	require.Equal(t, zapcore.DebugLevel, logEntry.Level)
	require.Contains(t, logEntry.Message, "Incoming request")
}
