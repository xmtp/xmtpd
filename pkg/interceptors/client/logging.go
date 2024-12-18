package client

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"

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

		return resp, err
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

		return err
	}
}
