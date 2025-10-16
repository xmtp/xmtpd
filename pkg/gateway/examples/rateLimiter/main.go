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

	authorizer := authorizers.NewRateLimitBuilder().
		WithLogger(gateway.MustCreateLogger(cfg)).
		WithRedis(gateway.MustSetupRedisClient(context.Background(), cfg.Redis)).
		WithLimits(authorizers.RateLimit{
			Capacity:    50,
			RefillEvery: time.Minute,
		}, authorizers.RateLimit{
			Capacity:    250,
			RefillEvery: time.Hour,
		}).
		MustBuild()

	gatewayService, err := gateway.NewGatewayServiceBuilder(cfg).
		WithAuthorizers(authorizer).
		Build()
	if err != nil {
		log.Fatalf("Failed to build gateway service: %v", err)
	}

	gatewayService.WaitForShutdown()
}
