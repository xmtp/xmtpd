package runner

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

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

// NewEnvironment creates a new environment for the test run.
// Each test run gets a new environment, completely isolated from other test runs.
func NewEnvironment(
	ctx context.Context,
	logger *zap.Logger,
	cfg Config,
	test string,
) (*types.Environment, error) {
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
