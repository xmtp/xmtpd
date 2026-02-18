package server

import (
	"context"
	"errors"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockConnectRequestProtocolValidationInterceptor struct {
	connect.AnyRequest
}

func (m *mockConnectRequestProtocolValidationInterceptor) Peer() connect.Peer {
	return connect.Peer{Protocol: connect.ProtocolConnect}
}

type mockStreamingConnProtocolValidationInterceptor struct {
	connect.StreamingHandlerConn
}

func (m *mockStreamingConnProtocolValidationInterceptor) Peer() connect.Peer {
	return connect.Peer{Protocol: connect.ProtocolConnect}
}

func TestUnaryProtocolValidationInterceptor(t *testing.T) {
	var (
		interceptor = NewProtocolValidationInterceptor()
		req         = &mockConnectRequestProtocolValidationInterceptor{}
		next        = func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			return nil, status.Errorf(codes.Internal, "mock internal error")
		}
	)

	conn, err := interceptor.WrapUnary(next)(t.Context(), req)
	require.Error(t, err)
	require.Nil(t, conn)
	require.Equal(t, connect.CodeFailedPrecondition, func() *connect.Error {
		target := &connect.Error{}
		_ = errors.As(err, &target)
		return target
	}().Code())
	require.Contains(t, err.Error(), errUnsupportedProtocol.Error())
}

func TestStreamProtocolValidationInterceptor(t *testing.T) {
	var (
		interceptor = NewProtocolValidationInterceptor()
		conn        = &mockStreamingConnProtocolValidationInterceptor{}
		next        = func(ctx context.Context, conn connect.StreamingHandlerConn) error {
			return status.Errorf(codes.NotFound, "mock stream error")
		}
	)

	err := interceptor.WrapStreamingHandler(next)(t.Context(), conn)
	require.Error(t, err)
	require.Equal(t, connect.CodeFailedPrecondition, func() *connect.Error {
		target := &connect.Error{}
		_ = errors.As(err, &target)
		return target
	}().Code())
	require.Contains(t, err.Error(), errUnsupportedProtocol.Error())
}
