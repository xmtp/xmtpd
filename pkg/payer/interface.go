package payer

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
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
	Build() (PayerService, error)
}

type PayerService interface {
	Serve(ctx context.Context) error
	ServeUntilShutdown() error
}
