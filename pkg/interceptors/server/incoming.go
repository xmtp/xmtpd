package server

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"net"
)

type IncomingInterceptor struct {
	logger *zap.Logger
}

func NewIncomingInterceptor(logger *zap.Logger) (*IncomingInterceptor, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	return &IncomingInterceptor{
		logger: logger,
	}, nil
}

func (i *IncomingInterceptor) logIncomingAddressIfAvailable(ctx context.Context) {
	if i.logger.Core().Enabled(zap.DebugLevel) {
		if p, ok := peer.FromContext(ctx); ok {
			clientAddr := p.Addr.String()
			var dnsName []string
			// Attempt to resolve the DNS name
			host, _, err := net.SplitHostPort(clientAddr)
			if err == nil {
				dnsName, err = net.LookupAddr(host)
				if err != nil || len(dnsName) == 0 {
					dnsName = []string{"Unknown"}
				}
			} else {
				dnsName = []string{"Unknown"}
			}
			i.logger.Debug(
				fmt.Sprintf("Incoming request from %s (DNS: %s)", clientAddr, dnsName[0]),
			)
		}
	}
}

// Unary intercepts unary RPC calls to log errors.
func (i *IncomingInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		i.logIncomingAddressIfAvailable(ctx)

		// Call the handler to complete the RPC
		return handler(ctx, req)
	}
}

// Stream intercepts stream RPC calls to log errors.
func (i *IncomingInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		i.logIncomingAddressIfAvailable(ss.Context())
		// Call the handler to complete the RPC
		return handler(srv, ss)
	}
}
