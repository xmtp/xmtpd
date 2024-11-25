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
	"github.com/xmtp/xmtpd/pkg/api/message"
	"github.com/xmtp/xmtpd/pkg/api/payer"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/indexer"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"github.com/xmtp/xmtpd/pkg/sync"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/grpc"

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
	indx         *indexer.Indexer
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
	listenAddress string,
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

	if options.Indexer.Enable {
		s.indx = indexer.NewIndexer(ctx, log)
		err = s.indx.StartIndexer(
			s.writerDB,
			options.Contracts,
			validationService,
		)

		if err != nil {
			return nil, err
		}
	}

	serviceRegistrationFunc := func(grpcServer *grpc.Server) error {

		if options.Replication.Enable {
			replicationService, err := message.NewReplicationApiService(
				ctx,
				log,
				s.registrant,
				writerDB,
			)
			if err != nil {
				return err
			}
			message_api.RegisterReplicationApiServer(grpcServer, replicationService)

			log.Info("Replication service enabled")
		}

		if options.Payer.Enable {
			payerPrivateKey, err := utils.ParseEcdsaPrivateKey(options.Payer.PrivateKey)
			if err != nil {
				return err
			}
			payerService, err := payer.NewPayerApiService(
				ctx,
				log,
				s.nodeRegistry,
				payerPrivateKey,
				blockchainPublisher,
			)
			if err != nil {
				return err
			}
			payer_api.RegisterPayerApiServer(grpcServer, payerService)

			log.Info("Payer service enabled")
		}

		return nil
	}

	if options.Payer.Enable || options.Replication.Enable {
		s.apiServer, err = api.NewAPIServer(
			s.ctx,
			log,
			listenAddress,
			options.Reflection.Enable,
			serviceRegistrationFunc,
		)
		if err != nil {
			return nil, err
		}

		log.Info("API server started", zap.Int("port", options.API.Port))
	}

	if options.Sync.Enable {
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

		log.Info("Sync server started")
	}

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
	if s.metrics != nil {
		s.metrics.Close()
	}

	if s.syncServer != nil {
		s.syncServer.Close()
	}

	if s.nodeRegistry != nil {
		s.nodeRegistry.Stop()
	}

	if s.indx != nil {
		s.indx.Close()
	}

	if s.apiServer != nil {
		s.apiServer.Close()
	}

	s.cancel()
}
