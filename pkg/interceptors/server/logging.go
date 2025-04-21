package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/xmtp/xmtpd/pkg/metrics"
	"google.golang.org/grpc/codes"

	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
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
				"Client Unary RPC Error",
				zap.String("method", info.FullMethod),
				zap.Duration("duration", duration),
				zap.Any("code", st.Code()),
				zap.String("message", st.Message()),
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
				"Stream Client RPC Error",
				zap.String("method", info.FullMethod),
				zap.Duration("duration", duration),
				zap.Any("code", st.Code()),
				zap.String("message", st.Message()),
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

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "request timed out")
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "request was canceled")
	default:
		st, ok := status.FromError(err)
		if !ok {
			return status.Error(codes.Internal, "internal server error")
		}
		msg := "request has failed"
		switch st.Code() {
		case codes.InvalidArgument, codes.Unimplemented, codes.NotFound:
			msg = st.Message()
		case codes.Internal:
			msg = "internal server error"
		}
		metrics.EmitNewFailedGRPCRequest(st.Code().String())

		return status.Error(st.Code(), msg)
	}
}
