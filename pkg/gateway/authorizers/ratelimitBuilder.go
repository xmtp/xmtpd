package authorizers

import (
	"log"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/xmtp/xmtpd/pkg/gateway"
	"go.uber.org/zap"
)

type RateLimitBuilder struct {
	logger      *zap.Logger
	redisClient redis.UniversalClient
	limits      []RateLimit
	keyPrefix   string
}

func NewRateLimitBuilder() *RateLimitBuilder {
	return &RateLimitBuilder{
		logger:      zap.NewNop(),
		redisClient: nil,
		limits:      []RateLimit{},
		keyPrefix:   "rl:",
	}
}

func (r *RateLimitBuilder) WithRedis(
	client redis.UniversalClient,
) *RateLimitBuilder {
	r.redisClient = client
	return r
}

func (r *RateLimitBuilder) WithLogger(logger *zap.Logger) *RateLimitBuilder {
	r.logger = logger
	return r
}

func (r *RateLimitBuilder) WithLimits(limits ...RateLimit) *RateLimitBuilder {
	r.limits = limits
	return r
}

func (r *RateLimitBuilder) WithKeyPrefix(prefix string) *RateLimitBuilder {
	r.keyPrefix = prefix
	return r
}

func (r *RateLimitBuilder) Build() (gateway.AuthorizePublishFn, error) {
	if r.redisClient == nil {
		return nil, errors.New("redis client is not set")
	}

	if len(r.limits) == 0 {
		return nil, errors.New("no rate limits configured")
	}

	return NewRateLimitAuthorizer(r.logger, r.redisClient, r.limits, r.keyPrefix)
}

func (r *RateLimitBuilder) MustBuild() gateway.AuthorizePublishFn {
	authorizer, err := r.Build()
	if err != nil {
		log.Fatalf("failed to build rate limit authorizer: %v", err)
	}

	return authorizer
}
