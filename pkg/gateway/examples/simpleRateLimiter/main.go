package main

import (
	"log"
	"time"

	"github.com/xmtp/xmtpd/pkg/gateway"
)

func main() {
	cfg := gateway.MustLoadConfig()

	gatewayService, err := gateway.NewGatewayServiceBuilder(cfg).
		// Rate limit to allow 50 requests per minute and 10,000 requests per day
		WithAuthorizers(gateway.NewRateLimitAuthorizer(cfg, gateway.RateLimit{
			MaxRequests: 50,
			Window:      time.Minute,
		}), gateway.NewRateLimitAuthorizer(cfg, gateway.RateLimit{
			MaxRequests: 10000,
			Window:      time.Hour * 24,
		})).
		Build()
	if err != nil {
		log.Fatalf("failed to build gateway service: %v", err)
	}

	gatewayService.WaitForShutdown()
}
