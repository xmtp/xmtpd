package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/utils"

	"go.uber.org/zap"

	"google.golang.org/grpc/status"
)

const (
	unknownDNSName = "unknown"

	codeField    = "code"
	messageField = "message"
)

// LoggingInterceptor logs errors for unary and stream RPCs.
type LoggingInterceptor struct {
	logger *zap.Logger
}

// Connect-go interceptor interface implementation.
var _ connect.Interceptor = (*LoggingInterceptor)(nil)

// NewLoggingInterceptor creates a new instance of LoggingInterceptor.
func NewLoggingInterceptor(logger *zap.Logger) (*LoggingInterceptor, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	return &LoggingInterceptor{
		logger: logger,
	}, nil
}

func (i *LoggingInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		start := time.Now()

		// Call the actual RPC handler.
		resp, err := next(ctx, req)

		duration := time.Since(start)

		if err != nil {
			st, _ := status.FromError(err)
			i.logger.Error(
				"client unary RPC error",
				utils.MethodField(req.Spec().Procedure),
				utils.DurationMsField(duration),
				zap.String(codeField, st.Code().String()),
				zap.String(messageField, st.Message()),
			)
		}

		return resp, sanitizeError(err)
	}
}

// WrapStreamingClient is a no-op. Interface requirement.
func (i *LoggingInterceptor) WrapStreamingClient(
	next connect.StreamingClientFunc,
) connect.StreamingClientFunc {
	return next
}

func (i *LoggingInterceptor) WrapStreamingHandler(
	next connect.StreamingHandlerFunc,
) connect.StreamingHandlerFunc {
	return func(
		ctx context.Context,
		conn connect.StreamingHandlerConn,
	) error {
		start := time.Now()

		// Call the actual stream handler.
		err := next(ctx, conn)

		duration := time.Since(start)

		if err != nil {
			st, _ := status.FromError(err)
			i.logger.Error(
				"stream client RPC error",
				utils.MethodField(conn.Spec().Procedure),
				utils.DurationMsField(duration),
				zap.String(codeField, st.Code().String()),
				zap.String(messageField, st.Message()),
			)
		}

		return sanitizeError(err)
	}
}

// sanitizeError standardizes and sanitizes gRPC errors to prevent internal details
// from being exposed to clients. It maps known internal errors (e.g., context cancellations
// or timeouts) to appropriate gRPC status codes and safely wraps unknown errors.
//
// If the error is a known Go context error (e.g., context.DeadlineExceeded or context.Canceled),
// it returns a corresponding gRPC error with a sanitized message.
//
// If the error is already a gRPC status error, it inspects the code:
//   - For InvalidArgument, NotFound and Unimplemented, it preserves the original code and message.
//   - For Internal errors, it replaces the message with a generic "internal server error".
//   - For all other codes, it replaces the message with a generic "request has failed".
func sanitizeError(err error) error {
	if err == nil {
		return nil
	}

	var (
		finalCode connect.Code
		finalMsg  string
	)

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		finalCode = connect.CodeDeadlineExceeded
		finalMsg = "request timed out"

	case errors.Is(err, context.Canceled):
		finalCode = connect.CodeCanceled
		finalMsg = "request was canceled"

	default:
		if st, ok := status.FromError(err); ok {
			finalCode = connect.Code(st.Code())
			switch finalCode {
			case connect.CodeInvalidArgument, connect.CodeUnimplemented, connect.CodeNotFound:
				finalMsg = st.Message()
			case connect.CodeInternal:
				finalMsg = "internal server error"
			default:
				finalMsg = "request has failed"
			}
		} else {
			finalCode = connect.CodeInternal
			finalMsg = "internal server error"
		}
	}

	// Emit metric for every non-nil error path
	metrics.EmitNewFailedGRPCRequest(connect.Code(finalCode))

	return connect.NewError(finalCode, errors.New(finalMsg))
}
