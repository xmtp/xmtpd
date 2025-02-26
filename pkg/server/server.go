package server

import (
	"context"
	"database/sql"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xmtp/xmtpd/pkg/api/metadata"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/fees"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"

	"github.com/Masterminds/semver/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/xmtp/xmtpd/pkg/api"
	"github.com/xmtp/xmtpd/pkg/api/message"
	"github.com/xmtp/xmtpd/pkg/api/payer"
	"github.com/xmtp/xmtpd/pkg/authn"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/indexer"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/sync"
	"github.com/xmtp/xmtpd/pkg/utils"
)

type ReplicationServer struct {
	apiServer  *api.ApiServer
	syncServer *sync.SyncServer

	ctx               context.Context
	cancel            context.CancelFunc
	log               *zap.Logger
	registrant        *registrant.Registrant
	nodeRegistry      registry.NodeRegistry
	indx              *indexer.Indexer
	options           config.ServerOptions
	metrics           *metrics.Server
	validationService mlsvalidate.MLSValidationService
	cursorUpdater     metadata.CursorUpdater
}

func NewReplicationServer(
	ctx context.Context,
	log *zap.Logger,
	options config.ServerOptions,
	nodeRegistry registry.NodeRegistry,
	writerDB *sql.DB,
	listenAddress string,
	serverVersion *semver.Version,
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
		metrics:      mtcs,
	}
	s.ctx, s.cancel = context.WithCancel(ctx)

	if options.Replication.Enable || options.Sync.Enable {
		s.registrant, err = registrant.NewRegistrant(
			s.ctx,
			log,
			queries.New(writerDB),
			nodeRegistry,
			options.Signer.PrivateKey,
			serverVersion,
		)
		if err != nil {
			return nil, err
		}
	}

	if options.Indexer.Enable {
		s.validationService, err = mlsvalidate.NewMlsValidationService(
			ctx,
			log,
			options.MlsValidation,
		)
		if err != nil {
			return nil, err
		}

		s.indx = indexer.NewIndexer(ctx, log)
		err = s.indx.StartIndexer(
			writerDB,
			options.Contracts,
			s.validationService,
		)

		if err != nil {
			return nil, err
		}

		log.Info("Indexer service enabled")
	}

	if options.Payer.Enable || options.Replication.Enable {
		err = startAPIServer(
			s.ctx,
			log,
			options,
			s,
			writerDB,
			listenAddress,
			serverVersion,
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
			writerDB,
		)
		if err != nil {
			return nil, err
		}

		log.Info("Sync service enabled")
	}

	return s, nil
}

func startAPIServer(
	ctx context.Context,
	log *zap.Logger,
	options config.ServerOptions,
	s *ReplicationServer,
	writerDB *sql.DB,
	listenAddress string,
	serverVersion *semver.Version,
) error {
	var err error

	serviceRegistrationFunc := func(grpcServer *grpc.Server) error {
		if options.Replication.Enable {
			if s.validationService == nil {
				s.validationService, err = mlsvalidate.NewMlsValidationService(
					ctx,
					log,
					options.MlsValidation,
				)
				if err != nil {
					return err
				}
			}
			s.cursorUpdater = metadata.NewCursorUpdater(ctx, log, writerDB)

			replicationService, err := message.NewReplicationApiService(
				ctx,
				log,
				s.registrant,
				writerDB,
				s.validationService,
				s.cursorUpdater,
				getRatesFetcher(),
			)
			if err != nil {
				return err
			}
			message_api.RegisterReplicationApiServer(grpcServer, replicationService)

			log.Info("Replication service enabled")

			metadataService, err := metadata.NewMetadataApiService(
				ctx,
				log,
				s.cursorUpdater,
			)
			if err != nil {
				return err
			}
			metadata_api.RegisterMetadataApiServer(grpcServer, metadataService)

			log.Info("Metadata service enabled")
		}

		if options.Payer.Enable {
			payerPrivateKey, err := utils.ParseEcdsaPrivateKey(options.Payer.PrivateKey)
			if err != nil {
				return err
			}

			signer, err := blockchain.NewPrivateKeySigner(
				options.Payer.PrivateKey,
				options.Contracts.ChainID,
			)
			if err != nil {
				log.Fatal("initializing signer", zap.Error(err))
			}

			ethclient, err := blockchain.NewClient(ctx, options.Contracts.RpcUrl)
			if err != nil {
				log.Fatal("initializing blockchain client", zap.Error(err))
			}

			blockchainPublisher, err := blockchain.NewBlockchainPublisher(
				ctx,
				log,
				ethclient,
				signer,
				options.Contracts,
			)
			if err != nil {
				log.Fatal("initializing message publisher", zap.Error(err))
			}

			payerService, err := payer.NewPayerApiService(
				ctx,
				log,
				s.nodeRegistry,
				payerPrivateKey,
				blockchainPublisher,
				nil,
			)
			if err != nil {
				return err
			}
			payer_api.RegisterPayerApiServer(grpcServer, payerService)

			log.Info("Payer service enabled")
		}

		return nil
	}

	var jwtVerifier authn.JWTVerifier

	if s.nodeRegistry != nil && s.registrant != nil {
		jwtVerifier, err = authn.NewRegistryVerifier(
			s.nodeRegistry,
			s.registrant.NodeID(),
			serverVersion,
		)
		if err != nil {
			return err
		}
	}

	s.apiServer, err = api.NewAPIServer(
		s.ctx,
		log,
		listenAddress,
		options.Reflection.Enable,
		serviceRegistrationFunc,
		jwtVerifier,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *ReplicationServer) Addr() net.Addr {
	return s.apiServer.Addr()
}

func (s *ReplicationServer) WaitForShutdown(timeout time.Duration) {
	termChannel := make(chan os.Signal, 1)
	signal.Notify(termChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-termChannel
	s.Shutdown(timeout)
}

func (s *ReplicationServer) Shutdown(timeout time.Duration) {
	if s.metrics != nil {
		s.metrics.Close()
	}

	if s.syncServer != nil {
		s.syncServer.Close()
	}

	if s.nodeRegistry != nil {
		s.nodeRegistry.Stop()
	}

	if s.cursorUpdater != nil {
		s.cursorUpdater.Stop()
	}

	if s.indx != nil {
		s.indx.Close()
	}
	if s.apiServer != nil {
		s.apiServer.Close(timeout)
	}

	s.cancel()
}

// TODO:nm Replace this with something that fetches rates from the blockchain
// Will need a rates smart contract first
func getRatesFetcher() fees.IRatesFetcher {
	return fees.NewFixedRatesFetcher(&fees.Rates{
		MessageFee:    currency.PicoDollar(100),
		StorageFee:    currency.PicoDollar(100),
		CongestionFee: currency.PicoDollar(100),
	})
}
