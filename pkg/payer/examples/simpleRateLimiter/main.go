package main

import (
	"context"
	"log"
	"time"

	"github.com/xmtp/xmtpd/pkg/payer"
)

func main() {
	cfg := payer.MustLoadConfig()

	payerService, err := payer.NewPayerServiceBuilder(cfg).
		// Rate limit to allow 50 requests per minute and 10,000 requests per day
		WithAuthorizers(payer.NewRateLimitAuthorizer(cfg, payer.RateLimit{
			MaxRequests: 50,
			Window:      time.Minute,
		}), payer.NewRateLimitAuthorizer(cfg, payer.RateLimit{
			MaxRequests: 10000,
			Window:      time.Hour * 24,
		})).
		Build()
	if err != nil {
		log.Fatalf("Failed to build payer service: %v", err)
	}

	err = payerService.Serve(context.Background())
	if err != nil {
		log.Fatalf("Failed to serve payer service: %v", err)
	}
}
