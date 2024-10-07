package server

import (
	"context"
	"database/sql"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/indexer"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/sync"

	"github.com/xmtp/xmtpd/pkg/api"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"go.uber.org/zap"
)

type ReplicationServer struct {
	apiServer  *api.ApiServer
	syncServer *sync.SyncServer

	ctx          context.Context
	cancel       context.CancelFunc
	log          *zap.Logger
	registrant   *registrant.Registrant
	nodeRegistry registry.NodeRegistry
	options      config.ServerOptions
	metrics      *metrics.Server
	writerDB     *sql.DB
	// Can add reader DB later if needed
}

func NewReplicationServer(
	ctx context.Context,
	log *zap.Logger,
	options config.ServerOptions,
	nodeRegistry registry.NodeRegistry,
	writerDB *sql.DB,
	blockchainPublisher blockchain.IBlockchainPublisher,
) (*ReplicationServer, error) {
	var err error

	var mtcs *metrics.Server
	if options.Metrics.Enable {
		promReg := prometheus.NewRegistry()
		promReg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		promReg.MustRegister(collectors.NewGoCollector())

		mtcs, err = metrics.NewMetricsServer(ctx,
			options.Metrics.Address,
			options.Metrics.Port,
			log,
			promReg,
		)
		if err != nil {

			log.Error("initializing metrics server", zap.Error(err))
			return nil, err
		}
	}

	s := &ReplicationServer{
		options:      options,
		log:          log,
		nodeRegistry: nodeRegistry,
		writerDB:     writerDB,
		metrics:      mtcs,
	}
	s.ctx, s.cancel = context.WithCancel(ctx)

	s.registrant, err = registrant.NewRegistrant(
		s.ctx,
		log,
		queries.New(s.writerDB),
		nodeRegistry,
		options.Signer.PrivateKey,
	)
	if err != nil {
		return nil, err
	}

	validationService, err := mlsvalidate.NewMlsValidationService(ctx, options.MlsValidation)
	if err != nil {
		return nil, err
	}
	err = indexer.StartIndexer(
		s.ctx,
		log,
		s.writerDB,
		options.Contracts,
		validationService,
	)
	if err != nil {
		return nil, err
	}

	// TODO(rich): Add configuration to specify whether to run API/sync server
	s.apiServer, err = api.NewAPIServer(
		s.ctx,
		s.writerDB,
		log,
		options.API.Port,
		s.registrant,
		options.Reflection.Enable,
		blockchainPublisher,
	)
	if err != nil {
		return nil, err
	}

	s.syncServer, err = sync.NewSyncServer(
		s.ctx,
		log,
		s.nodeRegistry,
		s.registrant,
		s.writerDB,
	)
	if err != nil {
		return nil, err
	}

	log.Info("Replication server started", zap.Int("port", options.API.Port))
	return s, nil
}

func (s *ReplicationServer) Addr() net.Addr {
	return s.apiServer.Addr()
}

func (s *ReplicationServer) WaitForShutdown() {
	termChannel := make(chan os.Signal, 1)
	signal.Notify(termChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-termChannel
	s.Shutdown()
}

func (s *ReplicationServer) Shutdown() {
	// Close metrics server.
	if s.metrics != nil {
		if err := s.metrics.Close(); err != nil {
			s.log.Error("stopping metrics", zap.Error(err))
		}
	}

	if s.apiServer != nil {
		s.apiServer.Close()
	}
	s.cancel()
}
