package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/grpc/codes"

	"go.uber.org/zap"

	"google.golang.org/grpc"
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

// NewLoggingInterceptor creates a new instance of LoggingInterceptor.
func NewLoggingInterceptor(logger *zap.Logger) (*LoggingInterceptor, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	return &LoggingInterceptor{
		logger: logger,
	}, nil
}

// Unary intercepts unary RPC calls to log errors.
func (i *LoggingInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req) // Call the actual RPC handler
		duration := time.Since(start)

		if err != nil {
			st, _ := status.FromError(err)
			i.logger.Error(
				"client unary RPC error",
				utils.MethodField(info.FullMethod),
				utils.DurationMsField(duration),
				zap.String(codeField, st.Code().String()),
				zap.String(messageField, st.Message()),
			)
		}

		return resp, sanitizeError(err)
	}
}

// Stream intercepts stream RPC calls to log errors.
func (i *LoggingInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()
		err := handler(srv, ss) // Call the actual stream handler
		duration := time.Since(start)

		if err != nil {
			st, _ := status.FromError(err)
			i.logger.Error(
				"stream client RPC error",
				utils.MethodField(info.FullMethod),
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
		finalCode codes.Code
		finalMsg  string
	)
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		finalCode = codes.DeadlineExceeded
		finalMsg = "request timed out"
	case errors.Is(err, context.Canceled):
		finalCode = codes.Canceled
		finalMsg = "request was canceled"
	default:
		if st, ok := status.FromError(err); ok {
			finalCode = st.Code()
			switch finalCode {
			case codes.InvalidArgument, codes.Unimplemented, codes.NotFound:
				finalMsg = st.Message()
			case codes.Internal:
				finalMsg = "internal server error"
			default:
				finalMsg = "request has failed"
			}
		} else {
			finalCode = codes.Internal
			finalMsg = "internal server error"
		}
	}
	// Emit metric for every non-nil error path
	metrics.EmitNewFailedGRPCRequest(finalCode)
	return status.Error(finalCode, finalMsg)
}
