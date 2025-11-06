// Package gateway implements the gateway service.
package gateway

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/xmtp/xmtpd/pkg/api"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/registry"
	"go.uber.org/zap"
)

type gatewayServiceImpl struct {
	logger              *zap.Logger
	identityFn          IdentityFn
	authorizers         []AuthorizePublishFn
	config              *config.GatewayConfig
	metrics             *metrics.Server
	ctx                 context.Context
	cancel              context.CancelFunc
	apiServer           *api.APIServer
	nodeRegistry        registry.NodeRegistry
	blockchainPublisher blockchain.IBlockchainPublisher
	redisClient         redis.UniversalClient
}

// Shutdown gracefully stops the API server and cleans up resources.
func (s *gatewayServiceImpl) Shutdown(timeout time.Duration) error {
	s.logger.Info("shutting down gateway service")

	if s.metrics != nil {
		s.metrics.Close()
	}

	if s.nodeRegistry != nil {
		s.nodeRegistry.Stop()
	}

	if s.blockchainPublisher != nil {
		s.blockchainPublisher.Close()
	}

	if s.cancel != nil {
		s.cancel()
	}

	if s.redisClient != nil {
		_ = s.redisClient.Close()
	}

	if s.apiServer != nil {
		s.apiServer.Close(timeout)
	}

	s.logger.Info("gateway service stopped")

	return nil
}

func (s *gatewayServiceImpl) WaitForShutdown(timeout time.Duration) {
	termChannel := make(chan os.Signal, 1)
	signal.Notify(termChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	sig := <-termChannel
	s.logger.Info("received OS signal, shutting down", zap.String("signal", sig.String()))
	_ = s.Shutdown(timeout)
}
