package server

import (
	"context"
	"database/sql"
	"errors"
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
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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

type ReplicationServerConfig struct {
	ctx               context.Context
	db                *sql.DB
	log               *zap.Logger
	nodeRegistry      registry.NodeRegistry
	options           *config.ServerOptions
	serverVersion     *semver.Version
	listenAddress     string
	httpListenAddress string
}

func WithContext(ctx context.Context) ReplicationServerOption {
	return func(cfg *ReplicationServerConfig) {
		cfg.ctx = ctx
	}
}

func WithDB(db *sql.DB) ReplicationServerOption {
	return func(cfg *ReplicationServerConfig) {
		cfg.db = db
	}
}

func WithLogger(log *zap.Logger) ReplicationServerOption {
	return func(cfg *ReplicationServerConfig) {
		cfg.log = log
	}
}

func WithNodeRegistry(reg registry.NodeRegistry) ReplicationServerOption {
	return func(cfg *ReplicationServerConfig) {
		cfg.nodeRegistry = reg
	}
}

func WithServerOptions(opts *config.ServerOptions) ReplicationServerOption {
	return func(cfg *ReplicationServerConfig) {
		cfg.options = opts
	}
}

func WithServerVersion(version *semver.Version) ReplicationServerOption {
	return func(cfg *ReplicationServerConfig) {
		cfg.serverVersion = version
	}
}

func WithListenAddress(addr string) ReplicationServerOption {
	return func(cfg *ReplicationServerConfig) {
		cfg.listenAddress = addr
	}
}

func WithHTTPListenAddress(addr string) ReplicationServerOption {
	return func(cfg *ReplicationServerConfig) {
		cfg.httpListenAddress = addr
	}
}

type ReplicationServer struct {
	ctx    context.Context
	cancel context.CancelFunc

	log                 *zap.Logger
	options             *config.ServerOptions
	metrics             *metrics.Server
	nodeRegistry        registry.NodeRegistry
	registrant          *registrant.Registrant
	validationService   mlsvalidate.MLSValidationService
	indx                *indexer.Indexer
	apiServer           *api.ApiServer
	syncServer          *sync.SyncServer
	cursorUpdater       metadata.CursorUpdater
	blockchainPublisher *blockchain.BlockchainPublisher
}

type ReplicationServerOption func(*ReplicationServerConfig)

func NewReplicationServer(
	opts ...ReplicationServerOption,
) (*ReplicationServer, error) {
	var err error

	cfg := &ReplicationServerConfig{}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.options == nil {
		return nil, errors.New("server options not provided")
	}

	if cfg.ctx == nil {
		return nil, errors.New("context not provided")
	}

	if cfg.log == nil {
		return nil, errors.New("logger not provided")
	}

	if cfg.nodeRegistry == nil {
		return nil, errors.New("node registry not provided")
	}

	if cfg.db == nil {
		return nil, errors.New("database not provided")
	}

	if cfg.listenAddress == "" {
		return nil, errors.New("listen address not provided")
	}

	if cfg.httpListenAddress == "" {
		return nil, errors.New("http listen address not provided")
	}

	promReg := prometheus.NewRegistry()

	clientMetrics := grpcprom.NewClientMetrics(
		grpcprom.WithClientHandlingTimeHistogram(),
	)

	var mtcs *metrics.Server
	if cfg.options.Metrics.Enable {
		promReg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		promReg.MustRegister(collectors.NewGoCollector())
		promReg.MustRegister(clientMetrics)

		mtcs, err = metrics.NewMetricsServer(cfg.ctx,
			cfg.options.Metrics.Address,
			cfg.options.Metrics.Port,
			cfg.log,
			promReg,
		)
		if err != nil {
			cfg.log.Error("initializing metrics server", zap.Error(err))
			return nil, err
		}
	}

	s := &ReplicationServer{
		options:      cfg.options,
		log:          cfg.log,
		nodeRegistry: cfg.nodeRegistry,
		metrics:      mtcs,
	}
	s.ctx, s.cancel = context.WithCancel(cfg.ctx)

	if cfg.options.Replication.Enable || cfg.options.Sync.Enable {
		s.registrant, err = registrant.NewRegistrant(
			s.ctx,
			cfg.log,
			queries.New(cfg.db),
			cfg.nodeRegistry,
			cfg.options.Signer.PrivateKey,
			cfg.serverVersion,
		)
		if err != nil {
			return nil, err
		}
	}

	if cfg.options.Indexer.Enable || cfg.options.Replication.Enable {
		s.validationService, err = mlsvalidate.NewMlsValidationService(
			cfg.ctx,
			cfg.log,
			cfg.options.MlsValidation,
			clientMetrics,
		)
		if err != nil {
			return nil, err
		}
	}

	if cfg.options.Indexer.Enable {
		s.indx, err = indexer.NewIndexer(
			indexer.WithDB(cfg.db),
			indexer.WithLogger(cfg.log),
			indexer.WithContext(cfg.ctx),
			indexer.WithValidationService(s.validationService),
			indexer.WithContractsOptions(&cfg.options.Contracts),
		)
		if err != nil {
			return nil, err
		}

		s.indx.StartIndexer()

		cfg.log.Info("Indexer service enabled")
	}

	if cfg.options.Payer.Enable || cfg.options.Replication.Enable {
		err = startAPIServer(
			s,
			cfg,
			clientMetrics,
			promReg,
		)
		if err != nil {
			return nil, err
		}

		cfg.log.Info("API server started", zap.Int("port", cfg.options.API.Port))
	}

	if cfg.options.Sync.Enable {
		s.syncServer, err = sync.NewSyncServer(
			sync.WithContext(s.ctx),
			sync.WithLogger(cfg.log),
			sync.WithNodeRegistry(s.nodeRegistry),
			sync.WithRegistrant(s.registrant),
			sync.WithDB(cfg.db),
			sync.WithFeeCalculator(fees.NewFeeCalculator(getRatesFetcher())),
		)
		if err != nil {
			return nil, err
		}

		cfg.log.Info("Sync service enabled")
	}

	return s, nil
}

func startAPIServer(
	s *ReplicationServer,
	cfg *ReplicationServerConfig,
	clientMetrics *grpcprom.ClientMetrics,
	promReg *prometheus.Registry,
) error {
	var err error

	serviceRegistrationFunc := func(grpcServer *grpc.Server) error {
		if cfg.options.Replication.Enable {

			s.cursorUpdater = metadata.NewCursorUpdater(s.ctx, cfg.log, cfg.db)

			replicationService, err := message.NewReplicationApiService(
				s.ctx,
				cfg.log,
				s.registrant,
				cfg.db,
				s.validationService,
				s.cursorUpdater,
				getRatesFetcher(),
				cfg.options.Replication,
			)
			if err != nil {
				return err
			}
			message_api.RegisterReplicationApiServer(grpcServer, replicationService)

			cfg.log.Info("Replication service enabled")

			metadataService, err := metadata.NewMetadataApiService(
				s.ctx,
				cfg.log,
				s.cursorUpdater,
			)
			if err != nil {
				return err
			}
			metadata_api.RegisterMetadataApiServer(grpcServer, metadataService)

			cfg.log.Info("Metadata service enabled")
		}

		if cfg.options.Payer.Enable {
			payerPrivateKey, err := utils.ParseEcdsaPrivateKey(cfg.options.Payer.PrivateKey)
			if err != nil {
				return err
			}

			signer, err := blockchain.NewPrivateKeySigner(
				cfg.options.Payer.PrivateKey,
				cfg.options.Contracts.AppChain.ChainID,
			)
			if err != nil {
				cfg.log.Fatal("initializing signer", zap.Error(err))
			}

			appChainClient, err := blockchain.NewClient(
				s.ctx,
				cfg.options.Contracts.AppChain.RpcURL,
			)
			if err != nil {
				cfg.log.Fatal("initializing blockchain client", zap.Error(err))
			}

			nonceManager := blockchain.NewSQLBackedNonceManager(cfg.db, cfg.log)

			blockchainPublisher, err := blockchain.NewBlockchainPublisher(
				s.ctx,
				cfg.log,
				appChainClient,
				signer,
				cfg.options.Contracts,
				nonceManager,
			)
			if err != nil {
				cfg.log.Fatal("initializing message publisher", zap.Error(err))
			}

			payerService, err := payer.NewPayerApiService(
				s.ctx,
				cfg.log,
				s.nodeRegistry,
				payerPrivateKey,
				blockchainPublisher,
				nil,
				clientMetrics,
			)
			if err != nil {
				return err
			}
			payer_api.RegisterPayerApiServer(grpcServer, payerService)

			cfg.log.Info("Payer service enabled")
		}

		return nil
	}

	httpRegistrationFunc := func(gwmux *runtime.ServeMux, conn *grpc.ClientConn) error {
		if cfg.options.Replication.Enable {
			err = metadata_api.RegisterMetadataApiHandler(s.ctx, gwmux, conn)
			if err != nil {
				return err
			}

			err = message_api.RegisterReplicationApiHandler(s.ctx, gwmux, conn)
			if err != nil {
				return err
			}
		}

		if cfg.options.Payer.Enable {
			err = payer_api.RegisterPayerApiHandler(s.ctx, gwmux, conn)
			if err != nil {
				return err
			}
		}

		return nil
	}
	var jwtVerifier authn.JWTVerifier

	if s.nodeRegistry != nil && s.registrant != nil {
		jwtVerifier, err = authn.NewRegistryVerifier(
			cfg.log,
			s.nodeRegistry,
			s.registrant.NodeID(),
			cfg.serverVersion,
		)
		if err != nil {
			return err
		}
	}

	s.apiServer, err = api.NewAPIServer(
		api.WithContext(s.ctx),
		api.WithLogger(cfg.log),
		api.WithHTTPListenAddress(cfg.httpListenAddress),
		api.WithListenAddress(cfg.listenAddress),
		api.WithJWTVerifier(jwtVerifier),
		api.WithRegistrationFunc(serviceRegistrationFunc),
		api.WithHTTPRegistrationFunc(httpRegistrationFunc),
		api.WithReflection(cfg.options.Reflection.Enable),
		api.WithPrometheusRegistry(promReg),
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

	if s.blockchainPublisher != nil {
		s.blockchainPublisher.Close()
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
