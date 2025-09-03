package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pingcap/log"
	"github.com/xmtp/xmtpd/pkg/api/metadata"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/fees"
	"github.com/xmtp/xmtpd/pkg/migrator"
	"github.com/xmtp/xmtpd/pkg/payerreport"
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
	"github.com/xmtp/xmtpd/pkg/authn"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/indexer"
	"github.com/xmtp/xmtpd/pkg/interceptors/server"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/sync"
)

type ReplicationServerConfig struct {
	Ctx           context.Context
	DB            *sql.DB
	Log           *zap.Logger
	NodeRegistry  registry.NodeRegistry
	Options       *config.ServerOptions
	ServerVersion *semver.Version
	GRPCListener  net.Listener
	HTTPListener  net.Listener
}

func WithContext(ctx context.Context) ReplicationServerOption {
	return func(cfg *ReplicationServerConfig) {
		cfg.Ctx = ctx
	}
}

func WithDB(db *sql.DB) ReplicationServerOption {
	return func(cfg *ReplicationServerConfig) {
		cfg.DB = db
	}
}

func WithLogger(log *zap.Logger) ReplicationServerOption {
	return func(cfg *ReplicationServerConfig) {
		cfg.Log = log
	}
}

func WithNodeRegistry(reg registry.NodeRegistry) ReplicationServerOption {
	return func(cfg *ReplicationServerConfig) {
		cfg.NodeRegistry = reg
	}
}

func WithServerOptions(opts *config.ServerOptions) ReplicationServerOption {
	return func(cfg *ReplicationServerConfig) {
		cfg.Options = opts
	}
}

func WithServerVersion(version *semver.Version) ReplicationServerOption {
	return func(cfg *ReplicationServerConfig) {
		cfg.ServerVersion = version
	}
}

func WithGRPCListener(listener net.Listener) ReplicationServerOption {
	return func(cfg *ReplicationServerConfig) {
		cfg.GRPCListener = listener
	}
}

func WithHTTPListener(listener net.Listener) ReplicationServerOption {
	return func(cfg *ReplicationServerConfig) {
		cfg.HTTPListener = listener
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
	migratorServer      *migrator.Migrator
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

	if cfg.Options == nil {
		return nil, errors.New("server Options not provided")
	}

	if cfg.Ctx == nil {
		return nil, errors.New("context not provided")
	}

	if cfg.Log == nil {
		return nil, errors.New("logger not provided")
	}

	if cfg.NodeRegistry == nil {
		return nil, errors.New("node registry not provided")
	}

	if cfg.DB == nil {
		return nil, errors.New("database not provided")
	}

	if cfg.GRPCListener == nil {
		return nil, errors.New("GRPC listener not provided")
	}

	if cfg.HTTPListener == nil {
		return nil, errors.New("http listener not provided")
	}

	promReg := prometheus.NewRegistry()

	clientMetrics := grpcprom.NewClientMetrics(
		grpcprom.WithClientHandlingTimeHistogram(),
	)

	var mtcs *metrics.Server
	if cfg.Options.Metrics.Enable {
		promReg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		promReg.MustRegister(collectors.NewGoCollector())
		promReg.MustRegister(clientMetrics)

		mtcs, err = metrics.NewMetricsServer(cfg.Ctx,
			cfg.Options.Metrics.Address,
			cfg.Options.Metrics.Port,
			cfg.Log,
			promReg,
		)
		if err != nil {
			cfg.Log.Error("initializing metrics server", zap.Error(err))
			return nil, err
		}
	}

	s := &ReplicationServer{
		options:      cfg.Options,
		log:          cfg.Log,
		nodeRegistry: cfg.NodeRegistry,
		metrics:      mtcs,
	}
	s.ctx, s.cancel = context.WithCancel(cfg.Ctx)

	if cfg.Options.Replication.Enable || cfg.Options.Sync.Enable {
		s.registrant, err = registrant.NewRegistrant(
			s.ctx,
			cfg.Log,
			queries.New(cfg.DB),
			cfg.NodeRegistry,
			cfg.Options.Signer.PrivateKey,
			cfg.ServerVersion,
		)
		if err != nil {
			return nil, err
		}
	}

	if cfg.Options.Indexer.Enable || cfg.Options.Replication.Enable {
		s.validationService, err = mlsvalidate.NewMlsValidationService(
			cfg.Ctx,
			cfg.Log,
			cfg.Options.MlsValidation,
			clientMetrics,
		)
		if err != nil {
			return nil, err
		}
	}

	if cfg.Options.Indexer.Enable {
		s.indx, err = indexer.NewIndexer(
			indexer.WithDB(cfg.DB),
			indexer.WithLogger(cfg.Log),
			indexer.WithContext(cfg.Ctx),
			indexer.WithValidationService(s.validationService),
			indexer.WithContractsOptions(&cfg.Options.Contracts),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize indexer: %w", err)
		}

		err = s.indx.StartIndexer()
		if err != nil {
			return nil, fmt.Errorf("failed to start indexer: %w", err)
		}

		cfg.Log.Info("Indexer service enabled")
	}

	if cfg.Options.MigrationServer.Enable {
		s.migratorServer, err = migrator.NewMigrationService(
			migrator.WithContext(cfg.Ctx),
			migrator.WithLogger(cfg.Log),
			migrator.WithDestinationDB(cfg.DB),
			migrator.WithMigratorConfig(&cfg.Options.MigrationServer),
			migrator.WithContractsOptions(&cfg.Options.Contracts),
		)
		if err != nil {
			return nil, err
		}

		err = s.migratorServer.Start()
		if err != nil {
			return nil, err
		}

		cfg.Log.Info("Migration service enabled")
	}

	err = startAPIServer(
		s,
		cfg,
		clientMetrics,
		promReg,
	)
	if err != nil {
		return nil, err
	}

	cfg.Log.Info("API server started", zap.Int("port", cfg.Options.API.Port))

	if cfg.Options.Sync.Enable {
		domainSeparator, err := getDomainSeparator(cfg.Ctx, cfg.Log, *cfg.Options)
		if err != nil {
			log.Error("failed to get domain separator", zap.Error(err))
			return nil, err
		}
		s.syncServer, err = sync.NewSyncServer(
			sync.WithContext(s.ctx),
			sync.WithLogger(cfg.Log),
			sync.WithNodeRegistry(s.nodeRegistry),
			sync.WithRegistrant(s.registrant),
			sync.WithDB(cfg.DB),
			sync.WithFeeCalculator(fees.NewFeeCalculator(getRatesFetcher())),
			sync.WithPayerReportStore(payerreport.NewStore(cfg.DB, cfg.Log)),
			sync.WithPayerReportDomainSeparator(domainSeparator),
		)
		if err != nil {
			return nil, err
		}

		cfg.Log.Info("Sync service enabled")
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

	isMigrationEnabled := cfg.Options.MigrationServer.Enable || cfg.Options.MigrationClient.Enable

	serviceRegistrationFunc := func(grpcServer *grpc.Server) error {
		if cfg.Options.Replication.Enable {

			s.cursorUpdater = metadata.NewCursorUpdater(s.ctx, cfg.Log, cfg.DB)

			replicationService, err := message.NewReplicationApiService(
				s.ctx,
				cfg.Log,
				s.registrant,
				cfg.DB,
				s.validationService,
				s.cursorUpdater,
				getRatesFetcher(),
				cfg.Options.Replication,
				isMigrationEnabled,
			)
			if err != nil {
				return err
			}
			message_api.RegisterReplicationApiServer(grpcServer, replicationService)

			cfg.Log.Info("Replication service enabled")

			metadataService, err := metadata.NewMetadataApiService(
				s.ctx,
				cfg.Log,
				s.cursorUpdater,
				cfg.ServerVersion,
				metadata.NewPayerInfoFetcher(cfg.DB),
			)
			if err != nil {
				return err
			}
			metadata_api.RegisterMetadataApiServer(grpcServer, metadataService)

			cfg.Log.Info("Metadata service enabled")
		}

		return nil
	}

	httpRegistrationFunc := func(gwmux *runtime.ServeMux, conn *grpc.ClientConn) error {
		if cfg.Options.Replication.Enable {
			err = metadata_api.RegisterMetadataApiHandler(s.ctx, gwmux, conn)
			if err != nil {
				return err
			}

			err = message_api.RegisterReplicationApiHandler(s.ctx, gwmux, conn)
			if err != nil {
				return err
			}
		}

		return nil
	}
	var jwtVerifier authn.JWTVerifier

	if s.nodeRegistry != nil && s.registrant != nil {
		jwtVerifier, err = authn.NewRegistryVerifier(
			cfg.Log,
			s.nodeRegistry,
			s.registrant.NodeID(),
			cfg.ServerVersion,
		)
		if err != nil {
			return err
		}
	}

	apiOpts := []api.ApiServerOption{
		api.WithContext(s.ctx),
		api.WithLogger(cfg.Log),
		api.WithGRPCListener(cfg.GRPCListener),
		api.WithHTTPListener(cfg.HTTPListener),
		api.WithRegistrationFunc(serviceRegistrationFunc),
		api.WithHTTPRegistrationFunc(httpRegistrationFunc),
		api.WithReflection(cfg.Options.Reflection.Enable),
		api.WithPrometheusRegistry(promReg),
	}

	// Add auth interceptors if JWT verifier is available
	if jwtVerifier != nil {
		authInterceptor := server.NewAuthInterceptor(jwtVerifier, cfg.Log)
		apiOpts = append(apiOpts,
			api.WithUnaryInterceptors(authInterceptor.Unary()),
			api.WithStreamInterceptors(authInterceptor.Stream()),
		)
	}

	s.apiServer, err = api.NewAPIServer(apiOpts...)
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
	sig := <-termChannel
	s.log.Info("Received OS signal, shutting down", zap.String("signal", sig.String()))
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

	if s.migratorServer != nil {
		if err := s.migratorServer.Stop(); err != nil {
			s.log.Error("failed to stop migration service", zap.Error(err))
		}
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

func getDomainSeparator(
	ctx context.Context,
	log *zap.Logger,
	options config.ServerOptions,
) (common.Hash, error) {
	signer, err := blockchain.NewPrivateKeySigner(
		options.Signer.PrivateKey,
		options.Contracts.SettlementChain.ChainID,
	)
	if err != nil {
		return common.Hash{}, err
	}

	settlementChainClient, err := blockchain.NewRPCClient(
		ctx,
		options.Contracts.SettlementChain.RPCURL,
	)
	if err != nil {
		return common.Hash{}, err
	}

	reportsManager, err := blockchain.NewReportsManager(
		log,
		settlementChainClient,
		signer,
		options.Contracts.SettlementChain,
	)
	if err != nil {
		return common.Hash{}, err
	}

	return reportsManager.GetDomainSeparator(ctx)
}
