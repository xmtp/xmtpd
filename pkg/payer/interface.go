package payer

import (
	"context"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
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
	IdentityKindIP          IdentityKind = "ip"
	IdentityKindUserDefined IdentityKind = "user"
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

type IdentityFn func(ctx context.Context) (Identity, error)

type AuthorizePublishFn func(ctx context.Context, identity Identity, req PublishRequest) (bool, error)

type IPayerServiceBuilder interface {
	WithIdentityFn(identityFn IdentityFn) IPayerServiceBuilder
	WithAuthorizers(authorizers ...AuthorizePublishFn) IPayerServiceBuilder
	WithBlockchainPublisher(
		blockchainPublisher blockchain.IBlockchainPublisher,
	) IPayerServiceBuilder
	WithNodeRegistry(nodeRegistry registry.NodeRegistry) IPayerServiceBuilder
	WithLogger(logger *zap.Logger) IPayerServiceBuilder
	WithMetricsServer(metricsServer *metrics.Server) IPayerServiceBuilder
	WithContext(ctx context.Context) IPayerServiceBuilder
	WithPromRegistry(promRegistry *prometheus.Registry) IPayerServiceBuilder
	WithClientMetrics(clientMetrics *grpcprom.ClientMetrics) IPayerServiceBuilder
	WithNonceManager(nonceManager noncemanager.NonceManager) IPayerServiceBuilder
	Build() (PayerService, error)
}

type PayerService interface {
	WaitForShutdown()
}
