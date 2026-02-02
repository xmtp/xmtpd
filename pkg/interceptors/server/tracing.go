package server

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
)

// TracingInterceptor creates DD APM spans for all Connect RPC calls.
type TracingInterceptor struct{}

// Connect-go interceptor interface implementation.
var _ connect.Interceptor = (*TracingInterceptor)(nil)

// NewTracingInterceptor creates a new instance of TracingInterceptor.
func NewTracingInterceptor() *TracingInterceptor {
	return &TracingInterceptor{}
}

func (i *TracingInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		// Extract operation name from procedure (e.g., "/xmtp.xmtpd.api.v1.ReplicationApi/PublishPayerEnvelopes")
		procedure := req.Spec().Procedure
		operationName := fmt.Sprintf("grpc.unary %s", procedure)

		span, ctx := tracing.StartSpanFromContext(ctx, operationName)
		defer span.Finish()

		// Set standard gRPC tags
		span.SetTag(ext.SpanType, "rpc")
		span.SetTag(ext.RPCSystem, "grpc")
		span.SetTag(ext.RPCService, req.Spec().Procedure)
		span.SetTag("grpc.method_type", "unary")

		// Call the actual RPC handler
		resp, err := next(ctx, req)

		if err != nil {
			span.SetTag(ext.Error, true)
			span.SetTag(ext.ErrorMsg, err.Error())
			span.SetTag("grpc.code", connect.CodeOf(err).String())
		} else {
			span.SetTag("grpc.code", "OK")
		}

		return resp, err
	}
}

// WrapStreamingClient is a no-op for server-side interceptor.
func (i *TracingInterceptor) WrapStreamingClient(
	next connect.StreamingClientFunc,
) connect.StreamingClientFunc {
	return next
}

func (i *TracingInterceptor) WrapStreamingHandler(
	next connect.StreamingHandlerFunc,
) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		procedure := conn.Spec().Procedure
		operationName := fmt.Sprintf("grpc.stream %s", procedure)

		span, ctx := tracing.StartSpanFromContext(ctx, operationName)
		defer span.Finish()

		// Set standard gRPC tags
		span.SetTag(ext.SpanType, "rpc")
		span.SetTag(ext.RPCSystem, "grpc")
		span.SetTag(ext.RPCService, procedure)
		span.SetTag("grpc.method_type", "stream")

		// Call the actual stream handler
		err := next(ctx, conn)

		if err != nil {
			span.SetTag(ext.Error, true)
			span.SetTag(ext.ErrorMsg, err.Error())
			span.SetTag("grpc.code", connect.CodeOf(err).String())
		} else {
			span.SetTag("grpc.code", "OK")
		}

		return err
	}
}
