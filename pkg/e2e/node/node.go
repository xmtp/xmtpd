// Package node provides a wrapper around the Node service used for E2E tests.
package node

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
	"go.uber.org/zap"
)

type Options struct {
	Image     string
	Network   string
	Alias     string
	WsURL     string
	RPCURL    string
	SignerKey string
	EnvVars   map[string]string
}

type Node struct {
	logger       *zap.Logger
	container    testcontainers.Container
	opts         Options
	dbConnStr    string
	externalAddr string
	nodeID       uint32
}

func New(ctx context.Context, logger *zap.Logger, opts Options) (*Node, error) {
	n := &Node{
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
	n.container, err = testcontainers.GenericContainer(
		createCtx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Logger:           log.New(io.Discard, "", 0),
		},
	)
	if err != nil {
		if n.container != nil {
			if logs, logErr := n.container.Logs(createCtx); logErr == nil {
				logBytes, _ := io.ReadAll(logs)
				_ = logs.Close()
				logger.Error("container logs on failure", zap.String("logs", string(logBytes)))
			}
		}
		return nil, fmt.Errorf("failed to start xmtpd container: %w", err)
	}

	mappedPort, err := n.container.MappedPort(createCtx, "5050/tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}
	n.externalAddr = "http://localhost:" + mappedPort.Port()

	dbName := "e2e_" + strings.ReplaceAll(opts.Alias, "-", "_")
	n.dbConnStr = fmt.Sprintf(
		"postgres://postgres:xmtp@localhost:8765/%s?sslmode=disable",
		dbName,
	)

	logger.Info("xmtpd node started", zap.String("alias", opts.Alias))

	return n, nil
}

func (n *Node) InternalAddr() string {
	return fmt.Sprintf("http://%s:5050", n.opts.Alias)
}

func (n *Node) ExternalAddr() string {
	return n.externalAddr
}

func (n *Node) Alias() string {
	return n.opts.Alias
}

func (n *Node) DBConnectionString() string {
	return n.dbConnStr
}

func (n *Node) NodeID() uint32 {
	return n.nodeID
}

func (n *Node) SignerKey() string {
	return n.opts.SignerKey
}

func (n *Node) SetNodeID(id uint32) {
	n.nodeID = id
}

func (n *Node) Stop(ctx context.Context) error {
	if n.container == nil {
		return nil
	}
	return n.container.Terminate(ctx)
}

// Restart creates a new container with the same configuration as the original.
// Returns a new Node instance; the caller should replace its reference.
// This is used to bring a stopped node back online while preserving its identity.
func (n *Node) Restart(ctx context.Context) (*Node, error) {
	n.logger.Info("restarting node", zap.String("alias", n.opts.Alias))
	restarted, err := New(ctx, n.logger, n.opts)
	if err != nil {
		return nil, fmt.Errorf("failed to restart node %s: %w", n.opts.Alias, err)
	}
	restarted.nodeID = n.nodeID
	return restarted, nil
}

// Container returns the underlying testcontainer for advanced use cases.
func (n *Node) Container() testcontainers.Container {
	return n.container
}

func buildEnvVars(opts Options) map[string]string {
	env := map[string]string{
		// Core services
		"XMTPD_API_ENABLE":     "true",
		"XMTPD_INDEXER_ENABLE": "true",
		"XMTPD_SYNC_ENABLE":    "true",

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

		// Unique signer per node (allocated from the node key pool)
		"XMTPD_SIGNER_PRIVATE_KEY": opts.SignerKey,

		// Fast registry refresh so nodes discover each other quickly in e2e
		"XMTPD_SETTLEMENT_CHAIN_NODE_REGISTRY_REFRESH_INTERVAL": "2s",

		// Payer report workers with short intervals for e2e testing
		"XMTPD_PAYER_REPORT_RUN_WORKERS":                      "true",
		"XMTPD_PAYER_REPORT_GENERATE_REPORT_SELF_PERIOD":      "2m",
		"XMTPD_PAYER_REPORT_GENERATE_REPORT_OTHERS_PERIOD":    "2m",
		"XMTPD_PAYER_REPORT_ATTESTATION_WORKER_POLL_INTERVAL": "10s",
		"XMTPD_PAYER_REPORT_EXPIRY_SELF_PERIOD":               "5m",
		"XMTPD_PAYER_REPORT_EXPIRY_OTHERS_PERIOD":             "5m",

		// Fast scheduling for e2e: 2-minute spread + 2-minute repeat (instead of 60/60)
		"XMTPD_PAYER_REPORT_WORKER_SPREAD_MINUTES":          "2",
		"XMTPD_PAYER_REPORT_WORKER_REPEAT_INTERVAL_MINUTES": "2",
	}

	// Allow caller overrides
	maps.Copy(env, opts.EnvVars)

	// Use unique DB name per node to isolate state.
	// Replace hyphens with underscores for valid Postgres identifiers.
	env["XMTPD_DB_NAME_OVERRIDE"] = "e2e_" + strings.ReplaceAll(opts.Alias, "-", "_")

	return env
}
