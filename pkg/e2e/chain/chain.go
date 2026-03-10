// Package chain provides a wrapper around the Anvil blockchain used for on-chain operations.
package chain

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

type ChainOptions struct {
	Image string
}

type Chain struct {
	logger    *zap.Logger
	container testcontainers.Container
	wsURL     string
	rpcURL    string
	network   string
	alias     string
	opts      ChainOptions
}

func New(
	ctx context.Context,
	logger *zap.Logger,
	networkName string,
	opts ChainOptions,
) (*Chain, error) {
	c := &Chain{
		logger:  logger,
		network: networkName,
		alias:   "anvil",
		opts:    opts,
	}

	createCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:        opts.Image,
		ExposedPorts: []string{"8545/tcp"},
		Networks:     []string{networkName},
		NetworkAliases: map[string][]string{
			networkName: {c.alias},
		},
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.AutoRemove = true
		},
		WaitingFor: wait.ForLog("Listening on"),
	}

	var err error
	c.container, err = testcontainers.GenericContainer(
		createCtx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Logger:           log.Default(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start anvil container: %w", err)
	}

	mappedPort, err := c.container.MappedPort(createCtx, "8545/tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	c.wsURL = "ws://localhost:" + mappedPort.Port()
	c.rpcURL = "http://localhost:" + mappedPort.Port()

	logger.Info("anvil chain started",
		zap.String("ws_url", c.wsURL),
		zap.String("rpc_url", c.rpcURL),
	)

	return c, nil
}

func (c *Chain) WsURL() string {
	return c.wsURL
}

func (c *Chain) RPCURL() string {
	return c.rpcURL
}

func (c *Chain) InternalWsURL() string {
	return fmt.Sprintf("ws://%s:8545", c.alias)
}

func (c *Chain) InternalRPCURL() string {
	return fmt.Sprintf("http://%s:8545", c.alias)
}

func (c *Chain) Stop(ctx context.Context) error {
	if c.container == nil {
		return nil
	}
	return c.container.Terminate(ctx)
}
