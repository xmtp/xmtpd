package main

import (
	"context"
	"log"
	"slices"

	"github.com/xmtp/xmtpd/pkg/gateway"
)

func main() {
	gatewayService, err := gateway.NewGatewayServiceBuilder(gateway.MustLoadConfig()).
		WithAuthorizers(func(ctx context.Context, identity gateway.Identity, req gateway.PublishRequestSummary) (bool, error) {
			// A simple authorization function that allows only the IP 127.0.0.1
			allowedIPs := []string{"127.0.0.1"}
			if !slices.Contains(allowedIPs, identity.Identity) {
				return false, gateway.ErrUnauthorized
			}
			return true, nil
		}).
		Build() // This will gather all the config from environment variables and flags
	if err != nil {
		log.Fatalf("failed to build gateway service: %v", err)
	}

	gatewayService.WaitForShutdown()
}
