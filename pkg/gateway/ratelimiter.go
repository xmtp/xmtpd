package gateway

import (
	"context"
	"time"

	"github.com/xmtp/xmtpd/pkg/config"
)

type RateLimit struct {
	Limit  int
	Window time.Duration
}

func NewRateLimitAuthorizer(config *config.GatewayConfig, rateLimit RateLimit) AuthorizePublishFn {
	return func(ctx context.Context, identity Identity, req PublishRequestSummary) (bool, error) {
		// TODO:(nm) Actual implementation
		return true, nil
	}
}
