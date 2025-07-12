package main

import (
	"context"
	"log"

	"github.com/xmtp/xmtpd/pkg/payer"
)

func main() {
	payerService, err := payer.NewPayerServiceBuilder(payer.MustLoadConfig()).
		Build() // This will gather all the config from environment variables and command line flags
	if err != nil {
		log.Fatalf("Failed to build payer service: %v", err)
	}

	err = payerService.Serve(context.Background())
	if err != nil {
		log.Fatalf("Failed to serve payer service: %v", err)
	}
}
