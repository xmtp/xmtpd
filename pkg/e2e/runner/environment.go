package runner

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types/network"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/xmtp/xmtpd/pkg/e2e/chain"
	"github.com/xmtp/xmtpd/pkg/e2e/chaos"
	"github.com/xmtp/xmtpd/pkg/e2e/keys"
	"github.com/xmtp/xmtpd/pkg/e2e/observe"
	"github.com/xmtp/xmtpd/pkg/e2e/types"
	"go.uber.org/zap"
)

type Environment = types.Environment

// hostDBConnStr is the connection string for the host Postgres instance
// used by all e2e node containers. Nodes create per-node namespaces (e2e_node_100, etc.)
// that must be cleaned up between test runs to avoid stale state.
const hostDBConnStr = "postgres://postgres:xmtp@localhost:8765/postgres?sslmode=disable"

func NewEnvironment(
	ctx context.Context,
	logger *zap.Logger,
	cfg Config,
	test string,
) (*types.Environment, error) {
	// Clean up stale resources from previous runs.
	cleanupStaleNetworks(ctx, logger)
	if err := dropE2EDatabases(ctx, logger); err != nil {
		logger.Warn("failed to clean up e2e databases (non-fatal)", zap.Error(err))
	}

	id := fmt.Sprintf("xmtpd-e2e-%s-%d", strings.ToLower(test), time.Now().Unix())

	env := &types.Environment{
		ID:      id,
		Logger:  logger,
		Config:  cfg,
		Network: id,
	}

	env.SetCleanupFunc(func(cleanCtx context.Context) error {
		return cleanupEnvironment(cleanCtx, env)
	})

	var err error

	err = createDockerNetwork(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker network: %w", err)
	}
	env.Network = id

	env.Chaos, err = chaos.NewController(ctx, logger.Named("chaos"), id)
	if err != nil {
		_ = env.Cleanup(ctx)
		return nil, fmt.Errorf("failed to start chaos controller: %w", err)
	}

	env.Chain, err = chain.New(ctx, logger.Named("chain"), id, chain.ChainOptions{
		Image: cfg.ChainImage,
	})
	if err != nil {
		_ = env.Cleanup(ctx)
		return nil, fmt.Errorf("failed to start chain: %w", err)
	}

	env.Keys = keys.NewManager(logger.Named("keys"), env.Chain.RPCURL())

	env.Contracts, err = chain.NewContracts(ctx, logger.Named("contracts"), env.Chain.RPCURL())
	if err != nil {
		_ = env.Cleanup(ctx)
		return nil, fmt.Errorf("failed to initialize contracts reader: %w", err)
	}

	env.SetObserver(observe.New(logger.Named("observer")))

	env.Redis, err = startRedis(ctx, id)
	if err != nil {
		_ = env.Cleanup(ctx)
		return nil, fmt.Errorf("failed to start redis: %w", err)
	}

	env.SetTestingT(types.NewTestingT(logger))

	return env, nil
}

func startRedis(ctx context.Context, id string) (testcontainers.Container, error) {
	createCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		Labels: map[string]string{
			"com.docker.compose.project": id,
		},
		Networks: []string{id},
		NetworkAliases: map[string][]string{
			id: {"redis"},
		},
		WaitingFor: wait.ForLog("Ready to accept connections").WithStartupTimeout(30 * time.Second),
	}

	return testcontainers.GenericContainer(
		createCtx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Logger:           log.New(io.Discard, "", 0),
		},
	)
}

func cleanupEnvironment(ctx context.Context, e *types.Environment) error {
	var firstErr error

	capture := func(err error) {
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}

	for _, gw := range e.Gateways() {
		capture(gw.Stop(ctx))
	}

	for _, n := range e.Nodes() {
		capture(n.Stop(ctx))
	}

	if e.Redis != nil {
		capture(e.Redis.Terminate(ctx))
	}

	if obs := e.Observer(); obs != nil {
		obs.Close()
	}

	if e.Contracts != nil {
		e.Contracts.Close()
	}

	if e.Chain != nil {
		capture(e.Chain.Stop(ctx))
	}

	if e.Chaos != nil {
		capture(e.Chaos.Stop(ctx))
	}

	if e.Network != "" {
		capture(removeDockerNetwork(ctx, e.Network))
	}

	return firstErr
}

func createDockerNetwork(ctx context.Context, id string) error {
	cli, err := dockerClient()
	if err != nil {
		return err
	}
	defer func() {
		_ = cli.Close()
	}()

	_, err = cli.NetworkCreate(ctx, id, dockerNetworkCreateOptions())
	if err != nil {
		return err
	}

	return nil
}

func removeDockerNetwork(ctx context.Context, name string) error {
	cli, err := dockerClient()
	if err != nil {
		return err
	}

	defer func() {
		_ = cli.Close()
	}()

	return cli.NetworkRemove(ctx, name)
}

// cleanupStaleNetworks removes Docker networks from previous e2e runs.
// This prevents Docker from exhausting its bridge subnet address space,
// which causes container-to-host networking (host.docker.internal) to fail.
func cleanupStaleNetworks(ctx context.Context, logger *zap.Logger) {
	cli, err := dockerClient()
	if err != nil {
		logger.Warn("failed to create docker client for network cleanup", zap.Error(err))
		return
	}
	defer func() { _ = cli.Close() }()

	networks, err := cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		logger.Warn("failed to list docker networks", zap.Error(err))
		return
	}

	for _, n := range networks {
		if strings.HasPrefix(n.Name, "xmtpd-e2e-") {
			if err := cli.NetworkRemove(ctx, n.Name); err != nil {
				logger.Warn("failed to remove stale network",
					zap.String("network", n.Name), zap.Error(err))
			} else {
				logger.Info("removed stale e2e network", zap.String("network", n.Name))
			}
		}
	}
}

// dropE2EDatabases drops idle databases matching the e2e_* pattern from the host
// Postgres. This prevents stale state from previous test runs from interfering
// with new runs (e.g., settled payer reports on a chain that no longer exists).
// Databases with active connections are skipped to avoid interfering with
// concurrently running test processes.
func dropE2EDatabases(ctx context.Context, logger *zap.Logger) error {
	db, err := sql.Open("postgres", hostDBConnStr)
	if err != nil {
		return fmt.Errorf("failed to connect to host postgres: %w", err)
	}
	defer func() { _ = db.Close() }()

	// Only select databases with NO active connections (excluding our own session).
	rows, err := db.QueryContext(ctx, `
		SELECT d.datname
		FROM pg_database d
		WHERE d.datname LIKE 'e2e_%'
		  AND NOT EXISTS (
			SELECT 1 FROM pg_stat_activity a
			WHERE a.datname = d.datname
			  AND a.pid != pg_backend_pid()
		  )
	`)
	if err != nil {
		return fmt.Errorf("failed to list e2e databases: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var dbNames []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("failed to scan database name: %w", err)
		}
		dbNames = append(dbNames, name)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for _, name := range dbNames {
		_, err := db.ExecContext(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %q", name))
		if err != nil {
			logger.Warn("failed to drop e2e database",
				zap.String("database", name), zap.Error(err),
			)
		} else {
			logger.Info("dropped stale e2e database", zap.String("database", name))
		}
	}

	return nil
}
