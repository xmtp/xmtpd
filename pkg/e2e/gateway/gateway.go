package gateway

import (
	"context"
	"fmt"
	"io"
	"log"
	"maps"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/xmtp/xmtpd/pkg/e2e/chaos"
	"go.uber.org/zap"
)

type Options struct {
	Image        string
	Network      string
	Alias        string
	WsURL        string
	RPCURL       string
	SignerKey    string
	EnvVars      map[string]string
	ChaosControl *chaos.Controller
}

type Gateway struct {
	logger    *zap.Logger
	container testcontainers.Container
	opts      Options
}

func New(ctx context.Context, logger *zap.Logger, opts Options) (*Gateway, error) {
	gw := &Gateway{
		logger: logger,
		opts:   opts,
	}

	createCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	envVars := buildEnvVars(opts)

	req := testcontainers.ContainerRequest{
		Image:        opts.Image,
		ExposedPorts: []string{"5050/tcp"},
		Networks:     []string{opts.Network},
		NetworkAliases: map[string][]string{
			opts.Network: {opts.Alias},
		},
		Env: envVars,
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.ExtraHosts = append(hc.ExtraHosts, "host.docker.internal:host-gateway")
		},
		WaitingFor: wait.ForLog("started api server").WithStartupTimeout(120 * time.Second),
	}

	var err error
	gw.container, err = testcontainers.GenericContainer(
		createCtx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Logger:           log.Default(),
		},
	)
	if err != nil {
		if gw.container != nil {
			if logs, logErr := gw.container.Logs(createCtx); logErr == nil {
				logBytes, _ := io.ReadAll(logs)
				_ = logs.Close()
				logger.Error("container logs on failure", zap.String("logs", string(logBytes)))
			}
		}
		return nil, fmt.Errorf("failed to start gateway container: %w", err)
	}

	if opts.ChaosControl != nil {
		if proxyErr := opts.ChaosControl.RegisterTarget(
			ctx,
			opts.Alias,
			opts.Alias,
			5050,
		); proxyErr != nil {
			logger.Warn("failed to register gateway with chaos controller", zap.Error(proxyErr))
		}
	}

	logger.Info("gateway started", zap.String("alias", opts.Alias))

	return gw, nil
}

func (g *Gateway) InternalAddr() string {
	return fmt.Sprintf("http://%s:5050", g.opts.Alias)
}

func (g *Gateway) Alias() string {
	return g.opts.Alias
}

func (g *Gateway) Stop(ctx context.Context) error {
	if g.container == nil {
		return nil
	}
	return g.container.Terminate(ctx)
}

func (g *Gateway) Container() testcontainers.Container {
	return g.container
}

func buildEnvVars(opts Options) map[string]string {
	env := map[string]string{
		// Core services
		"XMTPD_API_ENABLE":     "true",
		"XMTPD_INDEXER_ENABLE": "true",

		// Contracts config (embedded anvil.json with pre-deployed addresses)
		"XMTPD_CONTRACTS_ENVIRONMENT": "anvil",

		// Chain URLs (pointing to anvil container inside docker network)
		"XMTPD_SETTLEMENT_CHAIN_WSS_URL": opts.WsURL,
		"XMTPD_APP_CHAIN_WSS_URL":        opts.WsURL,
		"XMTPD_SETTLEMENT_CHAIN_RPC_URL": opts.RPCURL,
		"XMTPD_APP_CHAIN_RPC_URL":        opts.RPCURL,

		// Database (host postgres started by dev/up)
		"XMTPD_DB_WRITER_CONNECTION_STRING": "postgres://postgres:xmtp@host.docker.internal:8765/postgres?sslmode=disable",

		// MLS validation service (host service started by dev/up)
		"XMTPD_MLS_VALIDATION_GRPC_ADDRESS": "http://host.docker.internal:60051",

		// Signing keys (allocated from the gateway key pool)
		"XMTPD_SIGNER_PRIVATE_KEY": opts.SignerKey,
		"XMTPD_PAYER_PRIVATE_KEY":  opts.SignerKey,

		// Redis (pointing to redis container in docker network)
		"XMTPD_REDIS_URL": "redis://redis:6379/0",
	}

	maps.Copy(env, opts.EnvVars)

	// Replace hyphens with underscores for valid Postgres identifiers.
	env["XMTPD_DB_NAME_OVERRIDE"] = "e2e_" + strings.ReplaceAll(opts.Alias, "-", "_")

	return env
}
