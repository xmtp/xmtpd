package gateway

import (
	"context"
	"net/http"
	"time"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/blockchain/noncemanager"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"github.com/xmtp/xmtpd/pkg/registry"
	"go.uber.org/zap"
)

type PublishRequest *payer_api.PublishClientEnvelopesRequest

type IdentityKind string

const (
	identityKindIP          IdentityKind = "ip"
	identityKindUserDefined IdentityKind = "user"
)

type Identity struct {
	Kind     IdentityKind
	Identity string
}

type PublishRequestSummary struct {
	TotalEnvelopes       int
	OffchainCostEstimate currency.PicoDollar
	OnchainCostEstimate  currency.PicoDollar
	TotalCostEstimate    currency.PicoDollar
}

type IdentityFn func(headers http.Header, peer string) (Identity, error)

type AuthorizePublishFn func(ctx context.Context, identity Identity, req PublishRequestSummary) (bool, error)

type IGatewayServiceBuilder interface {
	WithIdentityFn(identityFn IdentityFn) IGatewayServiceBuilder
	WithAuthorizers(authorizers ...AuthorizePublishFn) IGatewayServiceBuilder
	WithBlockchainPublisher(
		blockchainPublisher blockchain.IBlockchainPublisher,
	) IGatewayServiceBuilder
	WithNodeRegistry(nodeRegistry registry.NodeRegistry) IGatewayServiceBuilder
	WithLogger(logger *zap.Logger) IGatewayServiceBuilder
	WithMetricsServer(metricsServer *metrics.Server) IGatewayServiceBuilder
	WithContext(ctx context.Context) IGatewayServiceBuilder
	WithPromRegistry(promRegistry *prometheus.Registry) IGatewayServiceBuilder
	WithClientMetrics(clientMetrics *grpcprom.ClientMetrics) IGatewayServiceBuilder
	WithNonceManager(nonceManager noncemanager.NonceManager) IGatewayServiceBuilder
	WithRedisClient(redisClient redis.UniversalClient) IGatewayServiceBuilder
	Build() (GatewayService, error)
}

type GatewayService interface {
	WaitForShutdown(timeout time.Duration)
}
