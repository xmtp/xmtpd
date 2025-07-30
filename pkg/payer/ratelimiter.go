package payer

import (
	"context"
	"time"

	"github.com/xmtp/xmtpd/pkg/config"
)

type RateLimit struct {
	MaxRequests int
	Window      time.Duration
}

func NewRateLimitAuthorizer(config *config.PayerConfig, rateLimit RateLimit) AuthorizePublishFn {
	return func(ctx context.Context, identity Identity, req PublishRequest) (bool, error) {
		// TODO:(nm) Actual implementation
		return true, nil
	}
}
