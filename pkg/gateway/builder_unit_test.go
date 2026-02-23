package gateway

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/config"
	mockblockchain "github.com/xmtp/xmtpd/pkg/mocks/blockchain"
	mocknoncemanager "github.com/xmtp/xmtpd/pkg/mocks/noncemanager"
	mockregistry "github.com/xmtp/xmtpd/pkg/mocks/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
	"go.uber.org/zap"
)

// Using mockery-generated mocks from pkg/mocks

func createMinimalTestConfig(t *testing.T) *config.GatewayConfig {
	wsURL, rpcURL := anvil.StartAnvil(t, false)
	return &config.GatewayConfig{
		Payer: testutils.GetPayerOptions(t),
		API: config.APIOptions{
			Port: 0,
		},
		Contracts: *testutils.NewContractsOptions(t, wsURL, rpcURL),
		Log: config.LogOptions{
			LogEncoding: "console",
			LogLevel:    "info",
		},
		Redis: config.RedisOptions{
			RedisURL:  "redis://localhost:6379",
			KeyPrefix: fmt.Sprintf("xmtpd:test:%s:", t.Name()),
		},
		Metrics: config.MetricsOptions{
			Enable:  false,
			Address: "localhost",
			Port:    0,
		},
	}
}

func TestBuilderBasicFunctionality(t *testing.T) {
	cfg := createMinimalTestConfig(t)
	builder := NewGatewayServiceBuilder(cfg)

	assert.NotNil(t, builder)
	assert.Equal(t, cfg, builder.(*GatewayServiceBuilder).config)
}

func TestBuilderMethodChaining(t *testing.T) {
	cfg := createMinimalTestConfig(t)
	builder := NewGatewayServiceBuilder(cfg)

	// Test WithIdentityFn
	identityFn := func(headers http.Header, peer string) (Identity, error) {
		return Identity{Kind: identityKindIP, Identity: "test"}, nil
	}
	result := builder.WithIdentityFn(identityFn)
	assert.Equal(t, builder, result) // Should return self for chaining

	// Test WithAuthorizers
	authFn := func(ctx context.Context, identity Identity, req PublishRequestSummary) (bool, error) {
		return true, nil
	}
	result = builder.WithAuthorizers(authFn)
	assert.Equal(t, builder, result)

	// Test WithNonceManager
	mockNonceManager := mocknoncemanager.NewMockNonceManager(t)
	result = builder.WithNonceManager(mockNonceManager)
	assert.Equal(t, builder, result)

	// Test WithContext
	ctx := context.Background()
	result = builder.WithContext(ctx)
	assert.Equal(t, builder, result)

	// Test WithLogger
	logger := zap.NewNop()
	result = builder.WithLogger(logger)
	assert.Equal(t, builder, result)
}

func TestBuilderDependencyStorage(t *testing.T) {
	cfg := createMinimalTestConfig(t)

	// Create mock dependencies
	mockBlockchainPublisher := mockblockchain.NewMockIBlockchainPublisher(t)
	mockNodeRegistry := mockregistry.NewMockNodeRegistry(t)
	mockNonceManager := mocknoncemanager.NewMockNonceManager(t)
	mockLogger := zap.NewNop()
	customPromRegistry := prometheus.NewRegistry()
	customClientMetrics := grpcprom.NewClientMetrics()

	builder := NewGatewayServiceBuilder(cfg).
		WithBlockchainPublisher(mockBlockchainPublisher).
		WithNodeRegistry(mockNodeRegistry).
		WithNonceManager(mockNonceManager).
		WithLogger(mockLogger).
		WithPromRegistry(customPromRegistry).
		WithClientMetrics(customClientMetrics)

	// Verify dependencies are set
	builderImpl := builder.(*GatewayServiceBuilder)
	assert.Equal(t, mockBlockchainPublisher, builderImpl.blockchainPublisher)
	assert.Equal(t, mockNodeRegistry, builderImpl.nodeRegistry)
	assert.Equal(t, mockNonceManager, builderImpl.nonceManager)
	assert.Equal(t, mockLogger, builderImpl.logger)
	assert.Equal(t, customPromRegistry, builderImpl.promRegistry)
	assert.Equal(t, customClientMetrics, builderImpl.clientMetrics)
}

func TestBuilderConfigValidation(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		builder := NewGatewayServiceBuilder(nil)
		service, err := builder.Build()

		require.Error(t, err)
		assert.Equal(t, ErrMissingConfig, err)
		assert.Nil(t, service)
	})

	t.Run("invalid private key", func(t *testing.T) {
		cfg := createMinimalTestConfig(t)
		cfg.Payer.PrivateKey = "invalid_key"
		// Mock dependencies to avoid other errors
		mockBlockchainPublisher := mockblockchain.NewMockIBlockchainPublisher(t)
		mockNodeRegistry := mockregistry.NewMockNodeRegistry(t)
		mockLogger := zap.NewNop()

		service, err := NewGatewayServiceBuilder(cfg).
			WithBlockchainPublisher(mockBlockchainPublisher).
			WithNodeRegistry(mockNodeRegistry).
			WithLogger(mockLogger).
			Build()

		require.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "failed to parse gateway private key")
	})

	t.Run("invalid redis url", func(t *testing.T) {
		cfg := createMinimalTestConfig(t)
		cfg.Redis.RedisURL = "invalid://redis-url"

		// Mock dependencies to avoid other errors
		mockBlockchainPublisher := mockblockchain.NewMockIBlockchainPublisher(t)
		mockNodeRegistry := mockregistry.NewMockNodeRegistry(t)
		mockLogger := zap.NewNop()

		service, err := NewGatewayServiceBuilder(cfg).
			WithBlockchainPublisher(mockBlockchainPublisher).
			WithNodeRegistry(mockNodeRegistry).
			WithLogger(mockLogger).
			Build()

		require.Error(t, err)
		assert.Nil(t, service)
		// The error could be from parsing or connection, check that it's Redis-related
		assert.Contains(t, err.Error(), "redis")
	})

	t.Run("empty redis url", func(t *testing.T) {
		cfg := createMinimalTestConfig(t)
		cfg.Redis.RedisURL = ""

		// Mock dependencies to avoid other errors
		mockBlockchainPublisher := mockblockchain.NewMockIBlockchainPublisher(t)
		mockNodeRegistry := mockregistry.NewMockNodeRegistry(t)
		mockLogger := zap.NewNop()

		service, err := NewGatewayServiceBuilder(cfg).
			WithBlockchainPublisher(mockBlockchainPublisher).
			WithNodeRegistry(mockNodeRegistry).
			WithLogger(mockLogger).
			Build()

		require.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "redis URL is empty")
	})
}

func TestBuilderDefaultIdentityFunction(t *testing.T) {
	cfg := createMinimalTestConfig(t)
	builder := NewGatewayServiceBuilder(cfg)

	// Don't set identity function - should default to IPIdentityFn
	builderImpl := builder.(*GatewayServiceBuilder)
	assert.Nil(t, builderImpl.identityFn) // Should be nil before build

	// After setting, verify it's stored
	builder.WithIdentityFn(IPIdentityFn)
	assert.NotNil(t, builderImpl.identityFn)
}

func TestBuilderContextHandling(t *testing.T) {
	cfg := createMinimalTestConfig(t)

	// Test custom context
	customCtx := context.Background()

	builder := NewGatewayServiceBuilder(cfg).WithContext(customCtx)
	builderImpl := builder.(*GatewayServiceBuilder)

	assert.Equal(t, customCtx, builderImpl.ctx)
}

func TestBuilderPrometheusMetricsHandling(t *testing.T) {
	cfg := createMinimalTestConfig(t)

	// Create custom prometheus registry and client metrics
	customPromRegistry := prometheus.NewRegistry()
	customClientMetrics := grpcprom.NewClientMetrics()

	builder := NewGatewayServiceBuilder(cfg).
		WithPromRegistry(customPromRegistry).
		WithClientMetrics(customClientMetrics)

	builderImpl := builder.(*GatewayServiceBuilder)
	assert.Equal(t, customPromRegistry, builderImpl.promRegistry)
	assert.Equal(t, customClientMetrics, builderImpl.clientMetrics)
}

func TestBuilderAllMethodsReturnBuilder(t *testing.T) {
	cfg := createMinimalTestConfig(t)

	// Create mock dependencies for testing
	mockNonceManager := mocknoncemanager.NewMockNonceManager(t)

	// Test method chaining returns proper interface
	builder := NewGatewayServiceBuilder(cfg).
		WithIdentityFn(IPIdentityFn).
		WithAuthorizers().
		WithNonceManager(mockNonceManager).
		WithContext(context.Background()).
		WithLogger(zap.NewNop()).
		WithPromRegistry(prometheus.NewRegistry()).
		WithClientMetrics(grpcprom.NewClientMetrics())

	assert.NotNil(t, builder)

	// Each method should return the builder interface
	assert.Implements(t, (*IGatewayServiceBuilder)(nil), builder)
}
