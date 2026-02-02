package server

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
)

// TracingInterceptor creates DD APM spans for all Connect RPC calls.
// Provides automatic instrumentation with meaningful span names for flamegraphs.
type TracingInterceptor struct{}

var _ connect.Interceptor = (*TracingInterceptor)(nil)

// NewTracingInterceptor creates a new instance of TracingInterceptor.
func NewTracingInterceptor() *TracingInterceptor {
	return &TracingInterceptor{}
}

func (i *TracingInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if !tracing.IsEnabled() {
			return next(ctx, req)
		}

		procedure := req.Spec().Procedure
		method := extractMethodName(procedure)
		service := extractServiceName(procedure)

		// Clean span name for better readability in flamegraphs
		// e.g., "xmtpd.api.PublishPayerEnvelopes" instead of "grpc.unary /xmtp..."
		operationName := "xmtpd.api." + method

		span, ctx := tracing.StartSpanFromContext(ctx, operationName)
		defer span.Finish()

		// Set standard tags for filtering and grouping
		span.SetTag(ext.SpanType, "web")
		span.SetTag(ext.RPCSystem, "grpc")
		tracing.SpanResource(span, method) // Shows nicely in DD UI
		tracing.SpanTag(span, "rpc.method", method)
		tracing.SpanTag(span, "rpc.service", service)
		tracing.SpanTag(span, "rpc.procedure", procedure)
		tracing.SpanTag(span, "rpc.type", "unary")

		resp, err := next(ctx, req)

		if err != nil {
			span.SetTag(ext.Error, true)
			span.SetTag(ext.ErrorMsg, err.Error())
			tracing.SpanTag(span, "rpc.status", connect.CodeOf(err).String())
		} else {
			tracing.SpanTag(span, "rpc.status", "OK")
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
		if !tracing.IsEnabled() {
			return next(ctx, conn)
		}

		procedure := conn.Spec().Procedure
		method := extractMethodName(procedure)
		service := extractServiceName(procedure)

		// Clean span name for streams
		operationName := "xmtpd.api." + method + ".stream"

		span, ctx := tracing.StartSpanFromContext(ctx, operationName)
		defer span.Finish()

		span.SetTag(ext.SpanType, "web")
		span.SetTag(ext.RPCSystem, "grpc")
		tracing.SpanResource(span, method)
		tracing.SpanTag(span, "rpc.method", method)
		tracing.SpanTag(span, "rpc.service", service)
		tracing.SpanTag(span, "rpc.procedure", procedure)
		tracing.SpanTag(span, "rpc.type", "stream")

		err := next(ctx, conn)

		if err != nil {
			span.SetTag(ext.Error, true)
			span.SetTag(ext.ErrorMsg, err.Error())
			tracing.SpanTag(span, "rpc.status", connect.CodeOf(err).String())
		} else {
			tracing.SpanTag(span, "rpc.status", "OK")
		}

		return err
	}
}

// extractMethodName gets the RPC method from the procedure path.
// "/xmtp.xmtpv4.ReplicationApi/PublishPayerEnvelopes" -> "PublishPayerEnvelopes"
func extractMethodName(procedure string) string {
	parts := strings.Split(procedure, "/")
	if len(parts) >= 3 {
		return parts[2]
	}
	if len(parts) >= 2 {
		return parts[1]
	}
	return procedure
}

// extractServiceName gets the service from the procedure path.
// "/xmtp.xmtpv4.ReplicationApi/PublishPayerEnvelopes" -> "ReplicationApi"
func extractServiceName(procedure string) string {
	parts := strings.Split(procedure, "/")
	if len(parts) >= 2 {
		// parts[1] is like "xmtp.xmtpv4.ReplicationApi"
		serviceParts := strings.Split(parts[1], ".")
		if len(serviceParts) > 0 {
			return serviceParts[len(serviceParts)-1]
		}
		return parts[1]
	}
	return "unknown"
}
