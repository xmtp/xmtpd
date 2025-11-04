package main

import (
	"context"
	"log"
	"time"

	"github.com/xmtp/xmtpd/pkg/gateway"
	"github.com/xmtp/xmtpd/pkg/gateway/authorizers"
)

func main() {
	cfg := gateway.MustLoadConfig()
	redis := gateway.MustSetupRedisClient(context.Background(), cfg.Redis)

	authorizer := authorizers.NewRateLimitBuilder().
		WithLogger(gateway.MustCreateLogger(cfg)).
		WithRedis(redis).
		// Set rate limits to 50 requests/minute and 250 requests/hour
		WithLimits(authorizers.RateLimit{
			Capacity:    50,
			RefillEvery: time.Minute,
		}, authorizers.RateLimit{
			Capacity:    250,
			RefillEvery: time.Hour,
		}).
		MustBuild()

	gatewayService, err := gateway.NewGatewayServiceBuilder(cfg).
		WithRedisClient(redis).
		WithAuthorizers(authorizer).
		Build()
	if err != nil {
		log.Fatalf("Failed to build gateway service: %v", err)
	}

	gatewayService.WaitForShutdown(30 * time.Second)
}
