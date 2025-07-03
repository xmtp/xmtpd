package payer

import (
	"context"
	"log"

	"github.com/jessevdk/go-flags"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/grpc/metadata"
)

func IPIdentityFn(ctx context.Context) (Identity, error) {
	return Identity{
		Kind:     IdentityKindIP,
		Identity: ClientIPFromContext(ctx),
	}, nil
}

func NewUserIdentity(userID string) Identity {
	return Identity{
		Kind:     IdentityKindUserDefined,
		Identity: userID,
	}
}

func LoadConfigFromEnv() (*config.PayerConfig, error) {
	var cfg config.PayerConfig
	_, err := flags.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func MustLoadConfig() *config.PayerConfig {
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		log.Fatalf("Failed to load config from environment: %v", err)
	}
	return cfg
}

func ClientIPFromContext(ctx context.Context) string {
	return utils.ClientIPFromContext(ctx)
}

func AuthorizationHeaderFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	auth := md.Get("authorization")
	if len(auth) == 0 {
		return ""
	}
	return auth[0]
}
