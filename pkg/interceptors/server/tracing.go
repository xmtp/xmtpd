package server

import (
	"context"

	"connectrpc.com/connect"
	"github.com/xmtp/xmtpd/pkg/tracing"
	grpcUtils "github.com/xmtp/xmtpd/pkg/utils/grpc"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
)

// TracingInterceptor creates DD APM spans for all gRPC and gRPC-Web calls.
// Provides automatic instrumentation with meaningful span names for flamegraphs.
type TracingInterceptor struct{}

var _ connect.Interceptor = (*TracingInterceptor)(nil)

// NewTracingInterceptor creates a new instance of TracingInterceptor.
func NewTracingInterceptor() *TracingInterceptor {
	return &TracingInterceptor{}
}

func (i *TracingInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		span, ctx := startRPCSpan(ctx, req.Spec().Procedure, "unary")
		defer span.Finish()

		resp, err := next(ctx, req)
		tagRPCResult(span, err)

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
		span, ctx := startRPCSpan(ctx, conn.Spec().Procedure, "stream")
		defer span.Finish()

		err := next(ctx, conn)
		tagRPCResult(span, err)

		return err
	}
}

// startRPCSpan creates a traced span for an RPC call and sets standard tags
// for Datadog APM filtering and grouping. The caller must call span.Finish().
func startRPCSpan(
	ctx context.Context,
	procedure string,
	rpcType string,
) (tracing.Span, context.Context) {
	_, service, method := grpcUtils.ParseProcedure(procedure)

	// Clean span name for better readability in flamegraphs
	// e.g., "xmtpd.api.PublishPayerEnvelopes" instead of "grpc.unary /xmtp..."
	operationName := "xmtpd.api." + method
	if rpcType == "stream" {
		operationName += ".stream"
	}

	span, ctx := tracing.StartSpanFromContext(ctx, operationName)

	// Set standard tags for filtering and grouping
	span.SetTag(ext.SpanType, "web")
	span.SetTag(ext.RPCSystem, "grpc")
	tracing.SpanResource(span, method) // Shows nicely in DD UI
	tracing.SpanTag(span, "rpc.method", method)
	tracing.SpanTag(span, "rpc.service", service)
	tracing.SpanTag(span, "rpc.procedure", procedure)
	tracing.SpanTag(span, "rpc.type", rpcType)

	return span, ctx
}

// tagRPCResult sets error or success status tags on the span.
func tagRPCResult(span tracing.Span, err error) {
	if err != nil {
		span.SetTag(ext.Error, true)
		span.SetTag(ext.ErrorMsg, err.Error())
		tracing.SpanTag(span, "rpc.status", connect.CodeOf(err).String())
	} else {
		tracing.SpanTag(span, "rpc.status", "OK")
	}
}
