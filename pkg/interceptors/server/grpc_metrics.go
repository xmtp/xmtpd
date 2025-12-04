package server

import (
	"context"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/xmtp/xmtpd/pkg/metrics"
)

// codeOK represents a successful RPC, matching gRPC's OK status (code 0).
// As extracted from https://github.com/connectrpc/connect-go/blob/main/code.go#L35:
// ConnectRPC intentionally doesn't export this constant because success is
// represented by a nil error in Go, not by an error with an OK code.
// We define it here for Prometheus metric label compatibility with gRPC dashboards.
const codeOK = connect.Code(0)

// GRPCMetricsInterceptor provides Prometheus metrics for ConnectRPC server calls.
// It emits metrics compatible with the standard grpc-ecosystem prometheus middleware.
type GRPCMetricsInterceptor struct{}

var _ connect.Interceptor = (*GRPCMetricsInterceptor)(nil)

func NewGRPCMetricsInterceptor() *GRPCMetricsInterceptor {
	return &GRPCMetricsInterceptor{}
}

func (i *GRPCMetricsInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		var (
			service, method = parseProcedure(req.Spec().Procedure)
			grpcType        = req.Spec().StreamType.String()
			start           = time.Now()
		)

		// Record that the RPC has started.
		metrics.EmitGRPCServerStarted(grpcType, service, method)

		// For unary calls, we count 1 message received (the request).
		metrics.EmitGRPCServerMsgReceived(grpcType, service, method)

		// Call the next handler.
		resp, err := next(ctx, req)

		// Determine the response code. See codeOK for more details.
		code := codeOK
		if err != nil {
			code = connect.CodeOf(err)
		} else {
			// For unary calls that succeed, count 1 message sent (the response).
			metrics.EmitGRPCServerMsgSent(grpcType, service, method)
		}

		// Record completion metrics.
		metrics.EmitGRPCServerHandled(grpcType, service, method, code)

		duration := time.Since(start)
		metrics.EmitGRPCServerHandlingTime(grpcType, service, method, duration)

		return resp, err
	}
}

// WrapStreamingClient is a no-op for server-side interceptors.
func (i *GRPCMetricsInterceptor) WrapStreamingClient(
	next connect.StreamingClientFunc,
) connect.StreamingClientFunc {
	return next
}

func (i *GRPCMetricsInterceptor) WrapStreamingHandler(
	next connect.StreamingHandlerFunc,
) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		var (
			service, method = parseProcedure(conn.Spec().Procedure)
			grpcType        = conn.Spec().StreamType.String()
			start           = time.Now()
		)

		// Record that the RPC has started.
		metrics.EmitGRPCServerStarted(grpcType, service, method)

		// Wrap the connection to count messages.
		wrappedConn := &metricsStreamingHandlerConn{
			StreamingHandlerConn: conn,
			grpcType:             grpcType,
			service:              service,
			method:               method,
		}

		// Call the next handler.
		err := next(ctx, wrappedConn)

		// Determine the response code. See codeOK for more details.
		code := codeOK
		if err != nil {
			code = connect.CodeOf(err)
		}

		// Record completion metrics.
		metrics.EmitGRPCServerHandled(grpcType, service, method, code)

		duration := time.Since(start)
		metrics.EmitGRPCServerHandlingTime(grpcType, service, method, duration)

		return err
	}
}

// metricsStreamingHandlerConn wraps a StreamingHandlerConn to count sent and received messages.
type metricsStreamingHandlerConn struct {
	connect.StreamingHandlerConn
	grpcType string
	service  string
	method   string
}

func (c *metricsStreamingHandlerConn) Receive(msg any) error {
	err := c.StreamingHandlerConn.Receive(msg)
	if err == nil {
		metrics.EmitGRPCServerMsgReceived(c.grpcType, c.service, c.method)
	}
	return err
}

func (c *metricsStreamingHandlerConn) Send(msg any) error {
	err := c.StreamingHandlerConn.Send(msg)
	if err == nil {
		metrics.EmitGRPCServerMsgSent(c.grpcType, c.service, c.method)
	}
	return err
}

// parseProcedure extracts service and method from a ConnectRPC procedure path.
// The procedure path format is: /package.service/method
// Example: /xmtp.xmtpv4.message_api.ReplicationApi/QueryEnvelopes
// Returns service="xmtp.xmtpv4.message_api.ReplicationApi", method="QueryEnvelopes"
func parseProcedure(procedure string) (string, string) {
	// Trim leading slash without allocating.
	if len(procedure) > 0 && procedure[0] == '/' {
		procedure = procedure[1:]
	}

	// Find last slash and return the service and method.
	if i := strings.LastIndexByte(procedure, '/'); i != -1 {
		return procedure[:i], procedure[i+1:]
	}

	return "unknown", procedure
}
