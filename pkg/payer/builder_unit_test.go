package payer

import (
	"context"
	"testing"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/xmtp/xmtpd/pkg/config"
	mockblockchain "github.com/xmtp/xmtpd/pkg/mocks/blockchain"
	mockregistry "github.com/xmtp/xmtpd/pkg/mocks/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
	"go.uber.org/zap"
)

// Using mockery-generated mocks from pkg/mocks

func createMinimalTestConfig(t *testing.T) *config.PayerConfig {
	wsUrl := anvil.StartAnvil(t, false)
	return &config.PayerConfig{
		Payer: testutils.GetPayerOptions(t),
		API: config.ApiOptions{
			Port:     0,
			HTTPPort: 0,
		},
		Contracts: testutils.NewContractsOptions(t, wsUrl),
		Log: config.LogOptions{
			LogEncoding: "console",
			LogLevel:    "info",
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
	builder := NewPayerServiceBuilder(cfg)

	assert.NotNil(t, builder)
	assert.Equal(t, cfg, builder.(*PayerServiceBuilder).config)
}

func TestBuilderMethodChaining(t *testing.T) {
	cfg := createMinimalTestConfig(t)
	builder := NewPayerServiceBuilder(cfg)

	// Test WithIdentityFn
	identityFn := func(ctx context.Context) (Identity, error) {
		return Identity{Kind: IdentityKindIP, Identity: "test"}, nil
	}
	result := builder.WithIdentityFn(identityFn)
	assert.Equal(t, builder, result) // Should return self for chaining

	// Test WithAuthorizers
	authFn := func(ctx context.Context, identity Identity, req PublishRequest) (bool, error) {
		return true, nil
	}
	result = builder.WithAuthorizers(authFn)
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
	mockLogger := zap.NewNop()
	customPromRegistry := prometheus.NewRegistry()
	customClientMetrics := grpcprom.NewClientMetrics()

	builder := NewPayerServiceBuilder(cfg).
		WithBlockchainPublisher(mockBlockchainPublisher).
		WithNodeRegistry(mockNodeRegistry).
		WithLogger(mockLogger).
		WithPromRegistry(customPromRegistry).
		WithClientMetrics(customClientMetrics)

	// Verify dependencies are set
	builderImpl := builder.(*PayerServiceBuilder)
	assert.Equal(t, mockBlockchainPublisher, builderImpl.blockchainPublisher)
	assert.Equal(t, mockNodeRegistry, builderImpl.nodeRegistry)
	assert.Equal(t, mockLogger, builderImpl.logger)
	assert.Equal(t, customPromRegistry, builderImpl.promRegistry)
	assert.Equal(t, customClientMetrics, builderImpl.clientMetrics)
}

func TestBuilderConfigValidation(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		builder := NewPayerServiceBuilder(nil)
		service, err := builder.Build()

		assert.Error(t, err)
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

		service, err := NewPayerServiceBuilder(cfg).
			WithBlockchainPublisher(mockBlockchainPublisher).
			WithNodeRegistry(mockNodeRegistry).
			WithLogger(mockLogger).
			Build()

		assert.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "failed to parse payer private key")
	})
}

func TestBuilderDefaultIdentityFunction(t *testing.T) {
	cfg := createMinimalTestConfig(t)
	builder := NewPayerServiceBuilder(cfg)

	// Don't set identity function - should default to IPIdentityFn
	builderImpl := builder.(*PayerServiceBuilder)
	assert.Nil(t, builderImpl.identityFn) // Should be nil before build

	// After setting, verify it's stored
	builder.WithIdentityFn(IPIdentityFn)
	assert.NotNil(t, builderImpl.identityFn)
}

func TestBuilderContextHandling(t *testing.T) {
	cfg := createMinimalTestConfig(t)

	// Test custom context
	customCtx := context.Background()

	builder := NewPayerServiceBuilder(cfg).WithContext(customCtx)
	builderImpl := builder.(*PayerServiceBuilder)

	assert.Equal(t, customCtx, builderImpl.ctx)
}

func TestBuilderPrometheusMetricsHandling(t *testing.T) {
	cfg := createMinimalTestConfig(t)

	// Create custom prometheus registry and client metrics
	customPromRegistry := prometheus.NewRegistry()
	customClientMetrics := grpcprom.NewClientMetrics()

	builder := NewPayerServiceBuilder(cfg).
		WithPromRegistry(customPromRegistry).
		WithClientMetrics(customClientMetrics)

	builderImpl := builder.(*PayerServiceBuilder)
	assert.Equal(t, customPromRegistry, builderImpl.promRegistry)
	assert.Equal(t, customClientMetrics, builderImpl.clientMetrics)
}

func TestBuilderAllMethodsReturnBuilder(t *testing.T) {
	cfg := createMinimalTestConfig(t)

	// Test method chaining returns proper interface
	builder := NewPayerServiceBuilder(cfg).
		WithIdentityFn(IPIdentityFn).
		WithAuthorizers().
		WithContext(context.Background()).
		WithLogger(zap.NewNop()).
		WithPromRegistry(prometheus.NewRegistry()).
		WithClientMetrics(grpcprom.NewClientMetrics())

	assert.NotNil(t, builder)

	// Each method should return the builder interface
	assert.Implements(t, (*IPayerServiceBuilder)(nil), builder)
}
