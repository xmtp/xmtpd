package gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jessevdk/go-flags"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/utils"

	"go.uber.org/zap"
)

/* Identity functions */

func IPIdentityFn(headers http.Header, peer string) (Identity, error) {
	return Identity{
		Kind:     identityKindIP,
		Identity: ClientIPFromHeaderOrPeer(headers, peer),
	}, nil
}

func NewUserIdentity(userID string) Identity {
	return Identity{
		Kind:     identityKindUserDefined,
		Identity: userID,
	}
}

func AuthorizationTokenFromContext(ctx context.Context) string {
	return utils.AuthorizationHeaderFromContext(ctx)
}

func AuthorizationTokenFromHeader(headers http.Header) string {
	return utils.AuthorizationTokenFromHeader(headers)
}

func ClientIPFromContext(ctx context.Context) string {
	return utils.ClientIPFromContext(ctx)
}

func ClientIPFromHeaderOrPeer(headers http.Header, peer string) string {
	return utils.ClientIPFromHeaderOrPeer(headers, peer)
}

func IdentityFromContext(ctx context.Context) (Identity, bool) {
	identity, ok := ctx.Value(identityCtxKey{}).(Identity)
	return identity, ok
}

/* Config helpers */

func LoadConfigFromEnv() (*config.GatewayConfig, error) {
	var cfg config.GatewayConfig
	_, err := flags.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	logger, _, err := utils.BuildLogger(cfg.Log)
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %v", err)
	}

	validator := config.NewOptionsValidator(logger)
	if err = validator.ParseJSONConfig(&cfg.Contracts); err != nil {
		return nil, fmt.Errorf("could not parse contracts JSON config: %s", err)
	}

	return &cfg, nil
}

func MustLoadConfig() *config.GatewayConfig {
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		log.Fatalf("failed to load config from environment: %v", err)
	}
	return cfg
}

func MustCreateLogger(cfg *config.GatewayConfig) *zap.Logger {
	logger, _, err := utils.BuildLogger(cfg.Log)
	if err != nil {
		log.Fatalf("Failed to setup logger: %v", err)
	}
	return logger
}
