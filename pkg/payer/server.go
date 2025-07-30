package payer

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xmtp/xmtpd/pkg/api"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/registry"
	"go.uber.org/zap"
)

type payerServiceImpl struct {
	log                 *zap.Logger
	identityFn          IdentityFn
	authorizers         []AuthorizePublishFn
	config              *config.PayerConfig
	metrics             *metrics.Server
	ctx                 context.Context
	cancel              context.CancelFunc
	apiServer           *api.ApiServer
	nodeRegistry        registry.NodeRegistry
	blockchainPublisher blockchain.IBlockchainPublisher
}

// Shutdown gracefully stops the API server and cleans up resources.
func (s *payerServiceImpl) Shutdown(timeout time.Duration) error {
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

	if s.apiServer != nil {
		s.apiServer.Close(timeout)
	}

	return nil
}

func (s *payerServiceImpl) WaitForShutdown() {
	timeout := 5 * time.Second
	termChannel := make(chan os.Signal, 1)
	signal.Notify(termChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-termChannel
	_ = s.Shutdown(timeout)
}
