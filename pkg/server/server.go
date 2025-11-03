// Package server implements the replication server.
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
	"github.com/xmtp/xmtpd/pkg/fees"
	"github.com/xmtp/xmtpd/pkg/migrator"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/payerreport/workers"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
	"github.com/xmtp/xmtpd/pkg/utils"

	"github.com/Masterminds/semver/v3"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
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
	Logger        *zap.Logger
	NodeRegistry  registry.NodeRegistry
	Options       *config.ServerOptions
	ServerVersion *semver.Version
	GRPCListener  net.Listener
	FeeCalculator fees.IFeeCalculator
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

func WithLogger(logger *zap.Logger) ReplicationServerOption {
	return func(cfg *ReplicationServerConfig) {
		cfg.Logger = logger
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

func WithFeeCalculator(feeCalculator fees.IFeeCalculator) ReplicationServerOption {
	return func(cfg *ReplicationServerConfig) {
		cfg.FeeCalculator = feeCalculator
	}
}

type ReplicationServer struct {
	ctx    context.Context
	cancel context.CancelFunc

	logger              *zap.Logger
	options             *config.ServerOptions
	metrics             *metrics.Server
	nodeRegistry        registry.NodeRegistry
	registrant          *registrant.Registrant
	validationService   mlsvalidate.MLSValidationService
	indx                *indexer.Indexer
	apiServer           *api.APIServer
	syncServer          *sync.SyncServer
	cursorUpdater       metadata.CursorUpdater
	blockchainPublisher *blockchain.BlockchainPublisher
	migratorServer      *migrator.Migrator
	reportWorkers       *workers.WorkerWrapper
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

	if cfg.Logger == nil {
		return nil, errors.New("logger not provided")
	}

	if cfg.NodeRegistry == nil {
		return nil, errors.New("node registry not provided")
	}

	if cfg.DB == nil {
		return nil, errors.New("database not provided")
	}

	if cfg.FeeCalculator == nil {
		return nil, errors.New("no fee calculator found")
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
			cfg.Logger,
			promReg,
		)
		if err != nil {
			cfg.Logger.Error("initializing metrics server", zap.Error(err))
			return nil, err
		}
	}

	s := &ReplicationServer{
		options:      cfg.Options,
		logger:       cfg.Logger,
		nodeRegistry: cfg.NodeRegistry,
		metrics:      mtcs,
	}
	s.ctx, s.cancel = context.WithCancel(cfg.Ctx)

	if cfg.Options.API.Enable || cfg.Options.Sync.Enable || cfg.Options.PayerReport.Enable {
		s.registrant, err = registrant.NewRegistrant(
			s.ctx,
			cfg.Logger,
			queries.New(cfg.DB),
			cfg.NodeRegistry,
			cfg.Options.Signer.PrivateKey,
			cfg.ServerVersion,
		)
		if err != nil {
			cfg.Logger.Error("failed to initialize registrant", zap.Error(err))
			return nil, err
		}
	}

	if cfg.Options.Indexer.Enable || cfg.Options.API.Enable {
		s.validationService, err = mlsvalidate.NewMlsValidationService(
			cfg.Ctx,
			cfg.Logger,
			cfg.Options.MlsValidation,
			clientMetrics,
		)
		if err != nil {
			cfg.Logger.Error("failed to initialize mls validation service", zap.Error(err))
			return nil, err
		}
	}

	if cfg.Options.Indexer.Enable {
		s.indx, err = indexer.NewIndexer(
			indexer.WithDB(cfg.DB),
			indexer.WithLogger(cfg.Logger),
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

		cfg.Logger.Info("indexer service started")
	}

	if cfg.Options.MigrationServer.Enable {
		s.migratorServer, err = migrator.NewMigrationService(
			migrator.WithContext(cfg.Ctx),
			migrator.WithLogger(cfg.Logger),
			migrator.WithDestinationDB(cfg.DB),
			migrator.WithMigratorConfig(&cfg.Options.MigrationServer),
			migrator.WithContractsOptions(&cfg.Options.Contracts),
		)
		if err != nil {
			cfg.Logger.Error("failed to initialize migrator", zap.Error(err))
			return nil, err
		}

		err = s.migratorServer.Start()
		if err != nil {
			cfg.Logger.Error("failed to start migrator", zap.Error(err))
			return nil, err
		}

		cfg.Logger.Info("migrator service started")
	}

	if cfg.Options.API.Enable {
		if cfg.GRPCListener == nil {
			return nil, errors.New("grpc listener not provided")
		}
		err = startAPIServer(
			s,
			cfg,
			clientMetrics,
			promReg,
		)
		if err != nil {
			cfg.Logger.Error("failed to start api server", zap.Error(err))
			return nil, err
		}

		cfg.Logger.Info("api service started", zap.Int("port", cfg.Options.API.Port))
	}

	if cfg.Options.Sync.Enable {
		domainSeparator, err := getDomainSeparator(cfg.Ctx, cfg.Logger, *cfg.Options)
		if err != nil {
			log.Error("failed to get domain separator", zap.Error(err))
			return nil, err
		}
		s.syncServer, err = sync.NewSyncServer(
			sync.WithContext(s.ctx),
			sync.WithLogger(cfg.Logger),
			sync.WithNodeRegistry(s.nodeRegistry),
			sync.WithRegistrant(s.registrant),
			sync.WithDB(cfg.DB),
			sync.WithFeeCalculator(cfg.FeeCalculator),
			sync.WithPayerReportStore(payerreport.NewStore(cfg.DB, cfg.Logger)),
			sync.WithPayerReportDomainSeparator(domainSeparator),
		)
		if err != nil {
			cfg.Logger.Error("failed to initialize sync server", zap.Error(err))
			return nil, err
		}

		cfg.Logger.Info("sync service started")
	}

	if cfg.Options.PayerReport.Enable {
		domainSeparator, err := getDomainSeparator(cfg.Ctx, cfg.Logger, *cfg.Options)
		if err != nil {
			cfg.Logger.Error(
				"failed to get domain separator for payer report workers",
				zap.Error(err),
			)
			return nil, err
		}

		signer, err := blockchain.NewPrivateKeySigner(
			cfg.Options.Signer.PrivateKey,
			cfg.Options.Contracts.SettlementChain.ChainID,
		)
		if err != nil {
			cfg.Logger.Error("failed to initialize signer for payer report workers", zap.Error(err))
			return nil, err
		}

		settlementChainClient, err := blockchain.NewRPCClient(
			cfg.Ctx,
			cfg.Options.Contracts.SettlementChain.RPCURL,
		)
		if err != nil {
			cfg.Logger.Error(
				"failed to initialize settlement chain client for payer report workers",
				zap.Error(err),
			)
			return nil, err
		}

		reportsManager, err := blockchain.NewReportsManager(
			cfg.Logger,
			settlementChainClient,
			signer,
			cfg.Options.Contracts.SettlementChain,
		)
		if err != nil {
			cfg.Logger.Error(
				"failed to initialize reports manager for payer report workers",
				zap.Error(err),
			)
			return nil, err
		}

		payerReportBaseLogger := cfg.Logger.Named(utils.PayerReportMainLoggerName)

		workerConfig, err := workers.NewWorkerConfigBuilder().
			WithContext(s.ctx).
			WithLogger(payerReportBaseLogger).
			WithRegistrant(s.registrant).
			WithRegistry(s.nodeRegistry).
			WithReportsManager(reportsManager).
			WithStore(payerreport.NewStore(cfg.DB, payerReportBaseLogger)).
			WithDomainSeparator(domainSeparator).
			WithAttestationPollInterval(cfg.Options.PayerReport.AttestationWorkerPollInterval).
			WithGenerationSelfPeriod(cfg.Options.PayerReport.GenerateReportSelfPeriod).
			WithGenerationOthersPeriod(cfg.Options.PayerReport.GenerateReportOthersPeriod).
			Build()
		if err != nil {
			cfg.Logger.Error("failed to build worker config", zap.Error(err))
			return nil, err
		}

		s.reportWorkers = workers.RunWorkers(*workerConfig)
	}

	return s, nil
}

func startAPIServer(
	s *ReplicationServer,
	cfg *ReplicationServerConfig,
	_ *grpcprom.ClientMetrics,
	promReg *prometheus.Registry,
) error {
	var err error

	isMigrationEnabled := cfg.Options.MigrationServer.Enable || cfg.Options.MigrationClient.Enable

	serviceRegistrationFunc := func(grpcServer *grpc.Server) error {
		s.cursorUpdater = metadata.NewCursorUpdater(s.ctx, cfg.Logger, cfg.DB)

		replicationService, err := message.NewReplicationAPIService(
			s.ctx,
			cfg.Logger,
			s.registrant,
			cfg.DB,
			s.validationService,
			s.cursorUpdater,
			cfg.FeeCalculator,
			cfg.Options.API,
			isMigrationEnabled,
			time.Second,
		)
		if err != nil {
			cfg.Logger.Error("failed to initialize replication api service", zap.Error(err))
			return err
		}
		message_api.RegisterReplicationApiServer(grpcServer, replicationService)

		cfg.Logger.Info("replication api registered")

		metadataService, err := metadata.NewMetadataAPIService(
			s.ctx,
			cfg.Logger,
			s.cursorUpdater,
			cfg.ServerVersion,
			metadata.NewPayerInfoFetcher(cfg.DB),
		)
		if err != nil {
			cfg.Logger.Error("failed to initialize metadata api service", zap.Error(err))
			return err
		}
		metadata_api.RegisterMetadataApiServer(grpcServer, metadataService)

		cfg.Logger.Info("metadata api registered")

		return nil
	}

	var jwtVerifier authn.JWTVerifier

	if s.nodeRegistry != nil && s.registrant != nil {
		jwtVerifier, err = authn.NewRegistryVerifier(
			cfg.Logger,
			s.nodeRegistry,
			s.registrant.NodeID(),
			cfg.ServerVersion,
		)
		if err != nil {
			cfg.Logger.Error("failed to initialize jwt verifier", zap.Error(err))
			return err
		}
	}

	apiOpts := []api.APIServerOption{
		api.WithContext(s.ctx),
		api.WithLogger(cfg.Logger),
		api.WithGRPCListener(cfg.GRPCListener),
		api.WithRegistrationFunc(serviceRegistrationFunc),
		api.WithReflection(cfg.Options.Reflection.Enable),
		api.WithPrometheusRegistry(promReg),
	}

	// Add auth interceptors if JWT verifier is available
	if jwtVerifier != nil {
		authInterceptor := server.NewAuthInterceptor(jwtVerifier, cfg.Logger)
		apiOpts = append(apiOpts,
			api.WithUnaryInterceptors(authInterceptor.Unary()),
			api.WithStreamInterceptors(authInterceptor.Stream()),
		)
	}

	s.apiServer, err = api.NewAPIServer(apiOpts...)
	if err != nil {
		cfg.Logger.Error("failed to initialize api server", zap.Error(err))
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
	s.logger.Info("received OS signal, shutting down", zap.String("signal", sig.String()))
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

	if s.reportWorkers != nil {
		s.reportWorkers.Stop()
	}

	if s.migratorServer != nil {
		if err := s.migratorServer.Stop(); err != nil {
			s.logger.Error("failed to stop migrator", zap.Error(err))
		}
	}

	s.cancel()
}

func getDomainSeparator(
	ctx context.Context,
	logger *zap.Logger,
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
		logger,
		settlementChainClient,
		signer,
		options.Contracts.SettlementChain,
	)
	if err != nil {
		return common.Hash{}, err
	}

	return reportsManager.GetDomainSeparator(ctx)
}
