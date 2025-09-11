package gateway

import (
	"context"
	"fmt"
	"net"
	"time"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/redis/go-redis/v9"
	"github.com/xmtp/xmtpd/pkg/api"
	"github.com/xmtp/xmtpd/pkg/api/payer"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/blockchain/noncemanager"
	redisnoncemanager "github.com/xmtp/xmtpd/pkg/blockchain/noncemanager/redis"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	ErrMissingConfig = errors.New("config must be provided and not nil")
	ErrUnauthorized  = errors.New("unauthorized")
)

type GatewayServiceBuilder struct {
	identityFn          IdentityFn
	authorizers         []AuthorizePublishFn
	config              *config.GatewayConfig
	blockchainPublisher blockchain.IBlockchainPublisher
	nodeRegistry        registry.NodeRegistry
	logger              *zap.Logger
	metricsServer       *metrics.Server
	ctx                 context.Context
	promRegistry        *prometheus.Registry
	clientMetrics       *grpcprom.ClientMetrics
	nonceManager        noncemanager.NonceManager
}

func NewGatewayServiceBuilder(config *config.GatewayConfig) IGatewayServiceBuilder {
	return &GatewayServiceBuilder{
		config: config,
	}
}

func (b *GatewayServiceBuilder) WithIdentityFn(identityFn IdentityFn) IGatewayServiceBuilder {
	b.identityFn = identityFn
	return b
}

func (b *GatewayServiceBuilder) WithAuthorizers(
	authorizers ...AuthorizePublishFn,
) IGatewayServiceBuilder {
	b.authorizers = authorizers
	return b
}

func (b *GatewayServiceBuilder) WithNonceManager(
	nonceManager noncemanager.NonceManager,
) IGatewayServiceBuilder {
	b.nonceManager = nonceManager
	return b
}

func (b *GatewayServiceBuilder) WithBlockchainPublisher(
	blockchainPublisher blockchain.IBlockchainPublisher,
) IGatewayServiceBuilder {
	b.blockchainPublisher = blockchainPublisher
	return b
}

func (b *GatewayServiceBuilder) WithNodeRegistry(
	nodeRegistry registry.NodeRegistry,
) IGatewayServiceBuilder {
	b.nodeRegistry = nodeRegistry
	return b
}

func (b *GatewayServiceBuilder) WithLogger(
	logger *zap.Logger,
) IGatewayServiceBuilder {
	b.logger = logger
	return b
}

func (b *GatewayServiceBuilder) WithMetricsServer(
	metricsServer *metrics.Server,
) IGatewayServiceBuilder {
	b.metricsServer = metricsServer
	return b
}

func (b *GatewayServiceBuilder) WithContext(
	ctx context.Context,
) IGatewayServiceBuilder {
	b.ctx = ctx
	return b
}

func (b *GatewayServiceBuilder) WithPromRegistry(
	promRegistry *prometheus.Registry,
) IGatewayServiceBuilder {
	b.promRegistry = promRegistry
	return b
}

func (b *GatewayServiceBuilder) WithClientMetrics(
	clientMetrics *grpcprom.ClientMetrics,
) IGatewayServiceBuilder {
	b.clientMetrics = clientMetrics
	return b
}

func (b *GatewayServiceBuilder) Build() (GatewayService, error) {
	if b.config == nil {
		return nil, ErrMissingConfig
	}

	if b.identityFn == nil {
		b.identityFn = IPIdentityFn
	}

	// Use injected context or default to background context
	ctx := b.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	// Create logger if not provided
	if b.logger == nil {
		logger, _, err := utils.BuildLogger(b.config.Log)
		if err != nil {
			return nil, errors.Wrap(err, "failed to build logger")
		}
		b.logger = logger
	}

	if b.nonceManager == nil {
		nonceManager, err := setupNonceManager(ctx, b.logger, b.config)
		if err != nil {
			return nil, errors.Wrap(err, "failed to setup nonce manager")
		}
		b.nonceManager = nonceManager
	}

	// Create blockchain publisher if not provided
	if b.blockchainPublisher == nil {
		blockchainPublisher, err := setupBlockchainPublisher(
			ctx,
			b.logger,
			b.config,
			b.nonceManager,
		)
		if err != nil {
			return nil, err
		}
		b.blockchainPublisher = blockchainPublisher
	}

	// Create node registry if not provided
	if b.nodeRegistry == nil {
		nodeRegistry, err := setupNodeRegistry(ctx, b.logger, b.config)
		if err != nil {
			return nil, errors.Wrap(err, "failed to setup node registry")
		}
		b.nodeRegistry = nodeRegistry
	}

	// Create metrics server if not provided and metrics are enabled
	promRegistry := b.promRegistry
	clientMetrics := b.clientMetrics
	if b.config.Metrics.Enable && b.metricsServer == nil {
		metricsServer, promReg, clientMet, err := setupMetrics(
			ctx,
			b.logger,
			&b.config.Metrics,
			b.promRegistry,
			b.clientMetrics,
		)
		if err != nil {
			return nil, err
		}
		b.metricsServer = metricsServer
		promRegistry = promReg
		clientMetrics = clientMet
	}

	return b.buildGatewayService(ctx, promRegistry, clientMetrics)
}

func (b *GatewayServiceBuilder) buildGatewayService(
	ctx context.Context,
	promRegistry *prometheus.Registry,
	clientMetrics *grpcprom.ClientMetrics,
) (GatewayService, error) {
	ctx, cancel := context.WithCancel(ctx)

	gatewayPrivateKey, err := utils.ParseEcdsaPrivateKey(b.config.Payer.PrivateKey)
	if err != nil {
		cancel()
		return nil, errors.Wrap(err, "failed to parse gateway private key")
	}

	serviceRegistrationFunc := func(grpcServer *grpc.Server) error {
		gatewayAPIService, err := payer.NewPayerAPIService(
			ctx,
			b.logger,
			b.nodeRegistry,
			gatewayPrivateKey,
			b.blockchainPublisher,
			nil,
			clientMetrics,
		)
		if err != nil {
			return err
		}
		payer_api.RegisterPayerApiServer(grpcServer, gatewayAPIService)

		return nil
	}

	httpRegistrationFunc := func(gwmux *runtime.ServeMux, conn *grpc.ClientConn) error {
		return payer_api.RegisterPayerApiHandler(ctx, gwmux, conn)
	}

	httpListener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", b.config.API.HTTPPort))
	if err != nil {
		cancel()
		return nil, err
	}

	grpcListener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", b.config.API.Port))
	if err != nil {
		_ = httpListener.Close()
		cancel()
		return nil, err
	}

	// Create gateway interceptor
	gatewayInterceptor := NewGatewayInterceptor(b.logger, b.identityFn, b.authorizers)

	apiServer, err := api.NewAPIServer(
		api.WithContext(ctx),
		api.WithLogger(b.logger),
		api.WithHTTPListener(httpListener),
		api.WithGRPCListener(grpcListener),
		api.WithRegistrationFunc(serviceRegistrationFunc),
		api.WithHTTPRegistrationFunc(httpRegistrationFunc),
		api.WithPrometheusRegistry(promRegistry),
		api.WithUnaryInterceptors(gatewayInterceptor.Unary()),
		api.WithStreamInterceptors(gatewayInterceptor.Stream()),
	)
	if err != nil {
		cancel()
		return nil, errors.Wrap(err, "failed to initialize api server")
	}

	return &gatewayServiceImpl{
		apiServer:           apiServer,
		ctx:                 ctx,
		cancel:              cancel,
		log:                 b.logger,
		identityFn:          b.identityFn,
		authorizers:         b.authorizers,
		metrics:             b.metricsServer,
		config:              b.config,
		blockchainPublisher: b.blockchainPublisher,
		nodeRegistry:        b.nodeRegistry,
	}, nil
}

func SetupRedisClient(
	ctx context.Context,
	redisURL string,
	timeout time.Duration,
) (redis.UniversalClient, error) {
	if redisURL == "" {
		return nil, fmt.Errorf("redis URL is empty")
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	deadline := time.Now().Add(timeout)
	for {
		if ctx.Err() != nil {
			_ = client.Close()
			return nil, fmt.Errorf(
				"context canceled while connecting to Redis at %s: %w",
				redisURL,
				ctx.Err(),
			)
		}
		if _, err := client.Ping(ctx).Result(); err == nil {
			break
		} else if time.Now().After(deadline) {
			_ = client.Close()
			return nil, fmt.Errorf("failed to connect to Redis at %s within %s: %w", redisURL, timeout, err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	return client, nil
}

func setupNonceManager(
	ctx context.Context,
	log *zap.Logger,
	cfg *config.GatewayConfig,
) (noncemanager.NonceManager, error) {
	redisClient, err := SetupRedisClient(ctx, cfg.Redis.RedisURL, 10*time.Second)
	if err != nil {
		return nil, err
	}

	return redisnoncemanager.NewRedisBackedNonceManager(redisClient, log, cfg.Redis.KeyPrefix)
}

func setupNodeRegistry(
	ctx context.Context,
	log *zap.Logger,
	cfg *config.GatewayConfig,
) (registry.NodeRegistry, error) {
	settlementChainClient, err := blockchain.NewRPCClient(
		ctx,
		cfg.Contracts.SettlementChain.RPCURL,
	)
	if err != nil {
		return nil, err
	}

	chainRegistry, err := registry.NewSmartContractRegistry(
		ctx,
		settlementChainClient,
		log,
		cfg.Contracts,
	)
	if err != nil {
		return nil, err
	}
	err = chainRegistry.Start()
	if err != nil {
		return nil, err
	}

	return chainRegistry, nil
}

func setupBlockchainPublisher(
	ctx context.Context,
	log *zap.Logger,
	cfg *config.GatewayConfig,
	nonceManager noncemanager.NonceManager,
) (*blockchain.BlockchainPublisher, error) {
	signer, err := blockchain.NewPrivateKeySigner(
		cfg.Payer.PrivateKey,
		cfg.Contracts.AppChain.ChainID,
	)
	if err != nil {
		return nil, err
	}

	appChainClient, err := blockchain.NewRPCClient(
		ctx,
		cfg.Contracts.AppChain.RPCURL,
	)
	if err != nil {
		return nil, err
	}

	return blockchain.NewBlockchainPublisher(
		ctx,
		log,
		appChainClient,
		signer,
		cfg.Contracts,
		nonceManager,
	)
}

// If metrics are enabled, sets them up
func setupMetrics(
	ctx context.Context,
	log *zap.Logger,
	metricsOptions *config.MetricsOptions,
	promRegistry *prometheus.Registry,
	clientMetrics *grpcprom.ClientMetrics,
) (*metrics.Server, *prometheus.Registry, *grpcprom.ClientMetrics, error) {
	// Use provided registry or create new one
	promReg := promRegistry
	if promReg == nil {
		promReg = prometheus.NewRegistry()
		promReg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		promReg.MustRegister(collectors.NewGoCollector())
	}

	// Use provided client metrics or create new ones
	clientMet := clientMetrics
	if clientMet == nil {
		clientMet = grpcprom.NewClientMetrics(
			grpcprom.WithClientHandlingTimeHistogram(),
		)
	}

	// Register client metrics if we have a registry
	if clientMet != nil {
		promReg.MustRegister(clientMet)
	}

	mtcs, err := metrics.NewMetricsServer(ctx,
		metricsOptions.Address,
		metricsOptions.Port,
		log,
		promReg,
	)
	if err != nil {
		log.Error("initializing metrics server", zap.Error(err))
		return nil, nil, nil, err
	}

	return mtcs, promReg, clientMet, nil
}
