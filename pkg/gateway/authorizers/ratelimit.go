package authorizers

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/xmtp/xmtpd/pkg/gateway"
	"github.com/xmtp/xmtpd/pkg/ratelimiter"
	"go.uber.org/zap"
)

type RateLimit = ratelimiter.Limit

func NewRateLimitAuthorizer(
	logger *zap.Logger,
	redisClient redis.UniversalClient,
	limits []RateLimit,
	keyPrefix string,
) (gateway.AuthorizePublishFn, error) {
	limiter, err := ratelimiter.NewRedisLimiter(redisClient, keyPrefix, limits)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context, identity gateway.Identity, req gateway.PublishRequestSummary) (bool, error) {
		result, err := limiter.Allow(ctx, identity.Identity, uint64(req.TotalEnvelopes))
		if err != nil {
			return false, err
		}

		if !result.Allowed {
			if result.RetryAfter == nil {
				return false, errors.New("expected retry after to be set")
			}

			errorMessage := fmt.Errorf(
				"rate limit exceeded for %s. Limit hit: %d requests per %.0f minutes",
				identity.Identity,
				result.FailedLimit.Capacity,
				result.FailedLimit.RefillEvery.Minutes(),
			)

			return false, gateway.NewRateLimitExceededError(errorMessage, *result.RetryAfter)
		}

		return true, nil
	}, nil
}
