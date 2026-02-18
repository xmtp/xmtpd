// Package server implements the base server that manages all the other services.
// Conceptually it's the server that represents the entire xmtpd node.
package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"github.com/DataDog/datadog-agent/pkg/trace/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/api/metadata"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/fees"
	"github.com/xmtp/xmtpd/pkg/migrator"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/payerreport/workers"

	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api/metadata_apiconnect"
	"github.com/xmtp/xmtpd/pkg/sync"
	"github.com/xmtp/xmtpd/pkg/utils"

	"github.com/Masterminds/semver/v3"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/api"
	"github.com/xmtp/xmtpd/pkg/api/message"
	"github.com/xmtp/xmtpd/pkg/authn"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/indexer"
	"github.com/xmtp/xmtpd/pkg/interceptors/server"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
)

type BaseServerConfig struct {
	Ctx           context.Context
	DB            *db.Handler
	Logger        *zap.Logger
	NodeRegistry  registry.NodeRegistry
	Options       *config.ServerOptions
	ServerVersion *semver.Version
	FeeCalculator fees.IFeeCalculator
	PromReg       *prometheus.Registry
}

func (cfg BaseServerConfig) Valid() error {
	var errs []error

	if cfg.Options == nil {
		errs = append(errs, errors.New("server options not provided"))
	}

	if cfg.Ctx == nil {
		errs = append(errs, errors.New("context not provided"))
	}

	if cfg.Logger == nil {
		errs = append(errs, errors.New("logger not provided"))
	}

	if cfg.NodeRegistry == nil {
		errs = append(errs, errors.New("node registry not provided"))
	}

	if cfg.DB == nil {
		errs = append(errs, errors.New("database handler not provided"))
	}

	if cfg.FeeCalculator == nil {
		errs = append(errs, errors.New("fee calculator not provided"))
	}

	if cfg.PromReg == nil {
		errs = append(errs, errors.New("prometheus registry not provided"))
	}

	return errors.Join(errs...)
}

func WithContext(ctx context.Context) BaseServerOption {
	return func(cfg *BaseServerConfig) {
		cfg.Ctx = ctx
	}
}

func WithDB(db *db.Handler) BaseServerOption {
	return func(cfg *BaseServerConfig) {
		cfg.DB = db
	}
}

func WithLogger(logger *zap.Logger) BaseServerOption {
	return func(cfg *BaseServerConfig) {
		cfg.Logger = logger
	}
}

func WithNodeRegistry(reg registry.NodeRegistry) BaseServerOption {
	return func(cfg *BaseServerConfig) {
		cfg.NodeRegistry = reg
	}
}

func WithServerOptions(opts *config.ServerOptions) BaseServerOption {
	return func(cfg *BaseServerConfig) {
		cfg.Options = opts
	}
}

func WithServerVersion(version *semver.Version) BaseServerOption {
	return func(cfg *BaseServerConfig) {
		cfg.ServerVersion = version
	}
}

func WithFeeCalculator(feeCalculator fees.IFeeCalculator) BaseServerOption {
	return func(cfg *BaseServerConfig) {
		cfg.FeeCalculator = feeCalculator
	}
}

func WithPromReg(promReg *prometheus.Registry) BaseServerOption {
	return func(cfg *BaseServerConfig) {
		cfg.PromReg = promReg
	}
}

type BaseServer struct {
	// Control mechanisms.
	ctx    context.Context
	cancel context.CancelFunc

	// Configuration.
	logger  *zap.Logger
	options *config.ServerOptions

	// Services
	api           *api.APIServer
	sync          *sync.SyncServer
	indexer       *indexer.Indexer
	mlsValidation mlsvalidate.MLSValidationService
	migrator      *migrator.Migrator
	metrics       *metrics.Server

	// Node identity.
	registrant   *registrant.Registrant
	nodeRegistry registry.NodeRegistry

	// Dependencies.
	db                  *db.Handler
	cursorUpdater       metadata.CursorUpdater
	blockchainPublisher *blockchain.BlockchainPublisher
	reportWorkers       *workers.WorkerWrapper
}

type BaseServerOption func(*BaseServerConfig)

// NewBaseServer creates a new base server.
// The Base server is the root service that manages the other services:
// - API server: Replication and metadata APIs.
// - Sync service: internal sync and communication between nodes.
// - Indexer service: indexes the blockchain and provides the data to the APIs.
// - Migration service: migrates an old V3 database to the D14N network.
// - Payer report service: generates and sends payer reports to the nodes.
func NewBaseServer(
	opts ...BaseServerOption,
) (*BaseServer, error) {
	var err error

	cfg := &BaseServerConfig{}

	for _, opt := range opts {
		opt(cfg)
	}

	err = cfg.Valid()
	if err != nil {
		return nil, fmt.Errorf("invalid base server configuration: %w", err)
	}

	promReg := cfg.PromReg

	// Client metrics are registered into the metrics server,
	// and also passed to the MLS validation service.
	clientMetrics := grpcprom.NewClientMetrics(
		grpcprom.WithClientHandlingTimeHistogram(),
	)

	// Setup metrics server.
	var metricsServer *metrics.Server
	if cfg.Options.Metrics.Enable {
		promReg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		promReg.MustRegister(collectors.NewGoCollector())
		promReg.MustRegister(clientMetrics)

		metricsServer, err = metrics.NewMetricsServer(cfg.Ctx,
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

	// Initialize services.
	svc := &BaseServer{
		options:      cfg.Options,
		logger:       cfg.Logger,
		nodeRegistry: cfg.NodeRegistry,
		metrics:      metricsServer,
		db:           cfg.DB,
	}

	svc.ctx, svc.cancel = context.WithCancel(cfg.Ctx)

	// Initialize registrant if needed, which is required for the API, sync and payer report services.
	// When the node runs as an indexer, it doesn't require an identity.
	if cfg.Options.API.Enable || cfg.Options.Sync.Enable || cfg.Options.PayerReport.Enable {
		svc.registrant, err = registrant.NewRegistrant(
			svc.ctx,
			cfg.Logger,
			cfg.DB.Query(),
			cfg.NodeRegistry,
			cfg.Options.Signer.PrivateKey,
			cfg.ServerVersion,
		)
		if err != nil {
			cfg.Logger.Error("failed to initialize registrant", zap.Error(err))
			return nil, err
		}
	}

	// Initialize MLS validation service if needed, which is required for the API and indexer services.
	// Sync and payer report services don't perform any MLS validation.
	if cfg.Options.Indexer.Enable || cfg.Options.API.Enable {
		svc.mlsValidation, err = mlsvalidate.NewMLSValidationService(
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

	// Maybe initialize indexer.
	if cfg.Options.Indexer.Enable {
		svc.indexer, err = indexer.NewIndexer(
			indexer.WithDB(cfg.DB),
			indexer.WithLogger(cfg.Logger),
			indexer.WithContext(cfg.Ctx),
			indexer.WithValidationService(svc.mlsValidation),
			indexer.WithContractsOptions(&cfg.Options.Contracts),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize indexer: %w", err)
		}

		err = svc.indexer.Start()
		if err != nil {
			return nil, fmt.Errorf("failed to start indexer: %w", err)
		}

		cfg.Logger.Info("indexer service started")
	}

	// Maybe initialize migration service.
	if cfg.Options.MigrationServer.Enable {
		svc.migrator, err = migrator.NewMigrationService(
			migrator.WithContext(cfg.Ctx),
			migrator.WithLogger(cfg.Logger),
			migrator.WithDestinationDB(cfg.DB),
			migrator.WithMigratorConfig(&cfg.Options.MigrationServer),
			migrator.WithContractsOptions(&cfg.Options.Contracts),
			migrator.WithFeeCalculator(cfg.FeeCalculator),
		)
		if err != nil {
			cfg.Logger.Error("failed to initialize migrator", zap.Error(err))
			return nil, err
		}

		err = svc.migrator.Start()
		if err != nil {
			cfg.Logger.Error("failed to start migrator", zap.Error(err))
			return nil, err
		}

		cfg.Logger.Info("migrator service started")
	}

	// Maybe initialize API server.
	// The API serves the replication and metadata APIs.
	if cfg.Options.API.Enable {
		svc.cursorUpdater = metadata.NewCursorUpdater(svc.ctx, cfg.Logger, cfg.DB)

		err = startAPIServer(
			svc,
			cfg,
			promReg,
		)
		if err != nil {
			cfg.Logger.Error("failed to start api server", zap.Error(err))
			return nil, err
		}
	}

	// Maybe initialize sync service.
	// The sync service is responsible for syncing nodes between each other.
	if cfg.Options.Sync.Enable {
		domainSeparator, err := getDomainSeparator(cfg.Ctx, cfg.Logger, *cfg.Options)
		if err != nil {
			log.Error("failed to get domain separator", zap.Error(err))
			return nil, err
		}

		svc.sync, err = sync.NewSyncServer(
			sync.WithContext(svc.ctx),
			sync.WithLogger(cfg.Logger),
			sync.WithNodeRegistry(svc.nodeRegistry),
			sync.WithRegistrant(svc.registrant),
			sync.WithDB(cfg.DB),
			sync.WithFeeCalculator(cfg.FeeCalculator),
			sync.WithPayerReportStore(
				payerreport.NewStore(
					cfg.Logger.Named(utils.PayerReportMainLoggerName).
						With(utils.WorkerNodeIDField(svc.registrant.NodeID())),
					cfg.DB),
			),
			sync.WithPayerReportDomainSeparator(domainSeparator),
			sync.WithClientMetrics(clientMetrics),
		)
		if err != nil {
			cfg.Logger.Error("failed to initialize sync server", zap.Error(err))
			return nil, err
		}

		cfg.Logger.Info("sync service started")
	}

	// Maybe initialize payer report service.
	// The payer report service is responsible for generating, attesting and submitting payer reports to the settlement chain.
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

		payerReportBaseLogger := cfg.Logger.Named(utils.PayerReportMainLoggerName).
			With(utils.WorkerNodeIDField(svc.registrant.NodeID()))

		workerConfig, err := workers.NewWorkerConfigBuilder().
			WithLogger(payerReportBaseLogger).
			WithContext(svc.ctx).
			WithRegistrant(svc.registrant).
			WithRegistry(svc.nodeRegistry).
			WithReportsManager(reportsManager).
			WithStore(payerreport.NewStore(payerReportBaseLogger, cfg.DB)).
			WithDomainSeparator(domainSeparator).
			WithAttestationPollInterval(cfg.Options.PayerReport.AttestationWorkerPollInterval).
			WithGenerationSelfPeriod(cfg.Options.PayerReport.GenerateReportSelfPeriod).
			WithGenerationOthersPeriod(cfg.Options.PayerReport.GenerateReportOthersPeriod).
			WithExpirySelfPeriod(cfg.Options.PayerReport.ExpirySelfPeriod).
			WithExpiryOthersPeriod(cfg.Options.PayerReport.ExpiryOthersPeriod).
			Build()
		if err != nil {
			cfg.Logger.Error("failed to build worker config", zap.Error(err))
			return nil, err
		}

		svc.reportWorkers = workers.RunWorkers(*workerConfig)
	}

	return svc, nil
}

func startAPIServer(
	svc *BaseServer,
	cfg *BaseServerConfig,
	promReg *prometheus.Registry,
) (err error) {
	isMigrationEnabled := cfg.Options.MigrationServer.Enable || cfg.Options.MigrationClient.Enable

	// Add auth interceptors if JWT verifier is available.
	var (
		jwtVerifier     authn.JWTVerifier
		authInterceptor *server.ServerAuthInterceptor
		apiOpts         = make([]api.APIServerOption, 0)
	)

	if svc.nodeRegistry != nil && svc.registrant != nil {
		jwtVerifier, err = authn.NewRegistryVerifier(
			cfg.Logger,
			svc.nodeRegistry,
			svc.registrant.NodeID(),
			cfg.ServerVersion,
		)
		if err != nil {
			cfg.Logger.Error("failed to initialize jwt verifier", zap.Error(err))
			return err
		}
	}

	if jwtVerifier != nil {
		authInterceptor = server.NewServerAuthInterceptor(
			cfg.Logger,
			jwtVerifier,
			server.RequireToken(cfg.Options.API.RequireJWTToken),
			server.DoDNSLookup(cfg.Options.API.AuthLoggingDNSLookup),
		)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Options.API.Port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", cfg.Options.API.Port, err)
	}

	registrationFunc := func(mux *http.ServeMux, interceptors ...connect.Interceptor) (servicePaths []string, err error) {
		if jwtVerifier != nil && authInterceptor != nil {
			interceptors = append(
				interceptors,
				authInterceptor,
			)
		}

		// Register replication API.
		replicationService, err := message.NewReplicationAPIService(
			svc.ctx,
			cfg.Logger,
			svc.registrant,
			cfg.NodeRegistry,
			cfg.DB,
			svc.mlsValidation,
			svc.cursorUpdater,
			cfg.FeeCalculator,
			cfg.Options.API,
			isMigrationEnabled,
			10*time.Millisecond,
		)
		if err != nil {
			return nil, err
		}

		if replicationService == nil {
			svc.logger.Error("replication service is nil")
			return nil, fmt.Errorf("replication service is nil")
		}

		replicationPath, replicationHandler := message_apiconnect.NewReplicationApiHandler(
			replicationService,
			connect.WithInterceptors(interceptors...),
		)

		mux.Handle(replicationPath, replicationHandler)

		svc.logger.Info("replication api registered")

		// Register metadata API.
		metadataService, err := metadata.NewMetadataAPIService(
			svc.ctx,
			cfg.Logger,
			svc.cursorUpdater,
			cfg.ServerVersion,
			metadata.NewPayerInfoFetcher(cfg.DB),
		)
		if err != nil {
			return nil, err
		}

		if metadataService == nil {
			svc.logger.Error("metadata service is nil")
			return nil, fmt.Errorf("metadata service is nil")
		}

		metadataPath, metadataHandler := metadata_apiconnect.NewMetadataApiHandler(
			metadataService,
			connect.WithInterceptors(interceptors...),
		)

		mux.Handle(metadataPath, metadataHandler)

		svc.logger.Info("metadata api registered")

		return []string{
			metadata_apiconnect.MetadataApiName,
			message_apiconnect.ReplicationApiName,
		}, nil
	}

	apiOpts = append(apiOpts, []api.APIServerOption{
		api.WithContext(svc.ctx),
		api.WithLogger(cfg.Logger),
		api.WithListener(listener),
		api.WithPrometheusRegistry(promReg),
		api.WithReflection(cfg.Options.Reflection.Enable),
		api.WithRegistrationFunc(registrationFunc),
	}...)

	svc.api, err = api.NewAPIServer(apiOpts...)
	if err != nil {
		cfg.Logger.Error("failed to initialize api server", zap.Error(err))
		return err
	}

	svc.api.Start()

	return nil
}

func (s *BaseServer) Addr() string {
	return s.api.Addr()
}

func (s *BaseServer) WaitForShutdown(timeout time.Duration) {
	termChannel := make(chan os.Signal, 1)
	signal.Notify(termChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	sig := <-termChannel
	s.logger.Info("received OS signal, shutting down", zap.String("signal", sig.String()))
	s.Shutdown(timeout)
}

func (s *BaseServer) Shutdown(timeout time.Duration) {
	if s.api != nil {
		s.api.Close(timeout)
	}

	if s.metrics != nil {
		s.metrics.Close()
	}

	if s.sync != nil {
		s.sync.Close()
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

	if s.indexer != nil {
		s.indexer.Close()
	}

	if s.reportWorkers != nil {
		s.reportWorkers.Stop()
	}

	if s.migrator != nil {
		if err := s.migrator.Stop(); err != nil {
			s.logger.Error("failed to stop migator", zap.Error(err))
		}
	}

	if s.db != nil {
		if err := s.db.Close(); err != nil {
			s.logger.Error("failed to close database connections", zap.Error(err))
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
