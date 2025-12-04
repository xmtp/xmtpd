package server

import (
	"context"
	"errors"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
)

type mockConnectRequestGRPCMetrics struct {
	connect.AnyRequest
	procedure string
}

func (m *mockConnectRequestGRPCMetrics) Spec() connect.Spec {
	return connect.Spec{
		Procedure:  m.procedure,
		StreamType: connect.StreamTypeUnary,
	}
}

type mockStreamingConnGRPCMetrics struct {
	connect.StreamingHandlerConn
	procedure  string
	streamType connect.StreamType
	received   int
	sent       int
}

func (m *mockStreamingConnGRPCMetrics) Spec() connect.Spec {
	return connect.Spec{
		Procedure:  m.procedure,
		StreamType: m.streamType,
	}
}

func (m *mockStreamingConnGRPCMetrics) Receive(_ any) error {
	m.received++
	return nil
}

func (m *mockStreamingConnGRPCMetrics) Send(_ any) error {
	m.sent++
	return nil
}

func TestParseProcedure(t *testing.T) {
	tests := []struct {
		name            string
		procedure       string
		expectedService string
		expectedMethod  string
	}{
		{
			name:            "standard procedure path",
			procedure:       "/xmtp.xmtpv4.message_api.ReplicationApi/QueryEnvelopes",
			expectedService: "xmtp.xmtpv4.message_api.ReplicationApi",
			expectedMethod:  "QueryEnvelopes",
		},
		{
			name:            "simple service",
			procedure:       "/test.TestService/TestMethod",
			expectedService: "test.TestService",
			expectedMethod:  "TestMethod",
		},
		{
			name:            "no leading slash",
			procedure:       "test.TestService/TestMethod",
			expectedService: "test.TestService",
			expectedMethod:  "TestMethod",
		},
		{
			name:            "no slash in procedure",
			procedure:       "JustAMethod",
			expectedService: "unknown",
			expectedMethod:  "JustAMethod",
		},
		{
			name:            "empty string",
			procedure:       "",
			expectedService: "unknown",
			expectedMethod:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, method := parseProcedure(tt.procedure)
			require.Equal(t, tt.expectedService, service)
			require.Equal(t, tt.expectedMethod, method)
		})
	}
}

func TestGRPCMetricsInterceptorUnarySuccess(t *testing.T) {
	interceptor := NewGRPCMetricsInterceptor()

	next := func(_ context.Context, _ connect.AnyRequest) (connect.AnyResponse, error) {
		return nil, nil
	}

	ctx := context.Background()
	req := &mockConnectRequestGRPCMetrics{procedure: "/test.TestService/TestMethod"}

	wrappedUnary := interceptor.WrapUnary(next)
	resp, err := wrappedUnary(ctx, req)

	require.NoError(t, err)
	require.Nil(t, resp)
}

func TestGRPCMetricsInterceptorUnaryError(t *testing.T) {
	interceptor := NewGRPCMetricsInterceptor()

	expectedErr := connect.NewError(connect.CodeInvalidArgument, errors.New("invalid argument"))
	next := func(_ context.Context, _ connect.AnyRequest) (connect.AnyResponse, error) {
		return nil, expectedErr
	}

	ctx := context.Background()
	req := &mockConnectRequestGRPCMetrics{procedure: "/test.TestService/TestMethod"}

	wrappedUnary := interceptor.WrapUnary(next)
	resp, err := wrappedUnary(ctx, req)

	require.Error(t, err)
	require.Nil(t, resp)
	require.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
}

func TestGRPCMetricsInterceptorStreamSuccess(t *testing.T) {
	interceptor := NewGRPCMetricsInterceptor()

	next := func(_ context.Context, conn connect.StreamingHandlerConn) error {
		// Simulate receiving and sending messages.
		_ = conn.Receive(nil)
		_ = conn.Receive(nil)
		_ = conn.Send(nil)
		return nil
	}

	ctx := context.Background()
	conn := &mockStreamingConnGRPCMetrics{
		procedure:  "/test.TestService/TestStream",
		streamType: connect.StreamTypeServer,
	}

	wrappedStream := interceptor.WrapStreamingHandler(next)
	err := wrappedStream(ctx, conn)

	require.NoError(t, err)
	require.Equal(t, 2, conn.received, "expected 2 receives")
	require.Equal(t, 1, conn.sent, "expected 1 send")
}

func TestGRPCMetricsInterceptorStreamError(t *testing.T) {
	interceptor := NewGRPCMetricsInterceptor()

	expectedErr := connect.NewError(connect.CodeInternal, errors.New("internal error"))
	next := func(_ context.Context, _ connect.StreamingHandlerConn) error {
		return expectedErr
	}

	ctx := context.Background()
	conn := &mockStreamingConnGRPCMetrics{
		procedure:  "/test.TestService/TestStream",
		streamType: connect.StreamTypeBidi,
	}

	wrappedStream := interceptor.WrapStreamingHandler(next)
	err := wrappedStream(ctx, conn)

	require.Error(t, err)
	require.Equal(t, connect.CodeInternal, connect.CodeOf(err))
}

func TestGRPCMetricsInterceptorWrapStreamingClientNoop(t *testing.T) {
	interceptor := NewGRPCMetricsInterceptor()

	next := func(_ context.Context, _ connect.Spec) connect.StreamingClientConn {
		return nil
	}

	// WrapStreamingClient should return the same function (no-op).
	wrapped := interceptor.WrapStreamingClient(next)
	require.NotNil(t, wrapped)
}
