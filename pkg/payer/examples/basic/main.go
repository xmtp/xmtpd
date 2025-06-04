package main

import (
	"context"
	"log"
	"slices"

	"github.com/xmtp/xmtpd/pkg/payer"
)

func main() {
	payerService, err := payer.NewPayerServiceBuilder(payer.MustLoadConfig()).
		WithAuthorizers(func(ctx context.Context, identity payer.Identity, req payer.PublishRequest) (bool, error) {
			// A simple authorization function that allows only the IP 127.0.0.1
			allowedIPs := []string{"127.0.0.1"}
			if !slices.Contains(allowedIPs, identity.Identity) {
				return false, payer.ErrUnauthorized
			}
			return true, nil
		}).
		Build() // This will gather all the config from environment variables and flags
	if err != nil {
		log.Fatalf("Failed to build payer service: %v", err)
	}

	err = payerService.Serve(context.Background())
	if err != nil {
		log.Fatalf("Failed to serve payer service: %v", err)
	}
}
