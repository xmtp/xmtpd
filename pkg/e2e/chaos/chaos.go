// Package chaos provides a toxiproxy-based chaos injection controller for E2E tests.
//
// The controller creates toxiproxy proxies that sit between nodes in the Docker
// network. When a node is registered, a proxy is created that forwards traffic
// from toxiproxy:<listen-port> to <upstream>:<port>. Nodes are registered on-chain
// with the proxy address so that all inter-node traffic flows through toxiproxy.
//
// Toxics (latency, bandwidth limits, connection resets, timeouts) can then be
// applied to these proxies to simulate network faults during E2E tests.
package chaos

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

const (
	toxiproxyImage = "ghcr.io/shopify/toxiproxy:2.12.0"
	toxiproxyAPI   = 8474
	// baseProxyPort is the starting port for toxiproxy listen addresses.
	// Each proxy gets the next available port (20000, 20001, 20002, ...).
	// These ports are only used within the Docker network (no host mapping needed).
	baseProxyPort = 20000
)

// ProxyTarget represents a toxiproxy proxy created for a service (node or gateway).
type ProxyTarget struct {
	// Name is the proxy name in toxiproxy (e.g. "node-100").
	Name string
	// Upstream is the container alias of the real service (e.g. "node-100").
	Upstream string
	// UpstreamPort is the port the real service listens on (e.g. 5050).
	UpstreamPort int
	// ListenPort is the port toxiproxy listens on for this proxy.
	ListenPort int
}

// Controller manages a toxiproxy container and its proxies for chaos injection.
type Controller struct {
	logger    *zap.Logger
	container testcontainers.Container
	apiURL    string
	network   string
	proxies   map[string]*ProxyTarget

	nextPort int
	mu       sync.Mutex
}

// NewController starts a toxiproxy container on the given Docker network
// and returns a controller for managing proxies and toxics.
func NewController(
	ctx context.Context,
	logger *zap.Logger,
	networkName string,
) (*Controller, error) {
	c := &Controller{
		logger:   logger,
		network:  networkName,
		proxies:  make(map[string]*ProxyTarget),
		nextPort: baseProxyPort,
	}

	createCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:        toxiproxyImage,
		ExposedPorts: []string{fmt.Sprintf("%d/tcp", toxiproxyAPI)},
		Networks:     []string{networkName},
		NetworkAliases: map[string][]string{
			networkName: {"toxiproxy"},
		},
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.ExtraHosts = append(hc.ExtraHosts, "host.docker.internal:host-gateway")
		},
		WaitingFor: wait.ForHTTP("/version").
			WithPort(nat.Port(fmt.Sprintf("%d/tcp", toxiproxyAPI))),
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
		return nil, fmt.Errorf("failed to start toxiproxy container: %w", err)
	}

	mappedPort, err := c.container.MappedPort(
		createCtx,
		nat.Port(fmt.Sprintf("%d/tcp", toxiproxyAPI)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get toxiproxy api port: %w", err)
	}

	c.apiURL = "http://localhost:" + mappedPort.Port()

	logger.Info("toxiproxy started", zap.String("api_url", c.apiURL))

	return c, nil
}

// allocatePort returns the next available port for a proxy listener.
func (c *Controller) allocatePort() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	port := c.nextPort
	c.nextPort++
	return port
}

// RegisterTarget creates a toxiproxy proxy for the given service.
// The proxy listens on toxiproxy:<allocated-port> and forwards to upstream:<port>.
// Other containers should connect to the proxy address (returned by ProxyAddress)
// instead of directly to the service.
func (c *Controller) RegisterTarget(ctx context.Context, name, upstream string, port int) error {
	listenPort := c.allocatePort()

	c.proxies[name] = &ProxyTarget{
		Name:         name,
		Upstream:     upstream,
		UpstreamPort: port,
		ListenPort:   listenPort,
	}

	// Create the proxy in toxiproxy: listen on 0.0.0.0:<listenPort>,
	// forward to <upstream>:<port>
	err := c.execToxiproxyCmd(ctx, "create",
		"-l", fmt.Sprintf("0.0.0.0:%d", listenPort),
		"-u", fmt.Sprintf("%s:%d", upstream, port),
		name,
	)
	if err != nil {
		delete(c.proxies, name)
		return fmt.Errorf("failed to create toxiproxy proxy %s: %w", name, err)
	}

	c.logger.Info("registered proxy target",
		zap.String("name", name),
		zap.String("upstream", fmt.Sprintf("%s:%d", upstream, port)),
		zap.String("proxy_address", fmt.Sprintf("toxiproxy:%d", listenPort)),
	)
	return nil
}

// ProxyAddress returns the Docker-network-accessible address for the given proxy.
// This is the address that should be used for on-chain registration so that
// inter-node traffic flows through toxiproxy.
//
// Returns "http://toxiproxy:<listen-port>" for the named proxy.
// Panics if no proxy with that name exists.
func (c *Controller) ProxyAddress(name string) string {
	proxy, ok := c.proxies[name]
	if !ok {
		panic("no proxy registered for " + name)
	}
	return fmt.Sprintf("http://toxiproxy:%d", proxy.ListenPort)
}

// AddLatency injects a network latency toxic on the named proxy.
// All connections through the proxy will experience the specified delay in milliseconds.
func (c *Controller) AddLatency(ctx context.Context, targetName string, latencyMs int) error {
	proxy, ok := c.proxies[targetName]
	if !ok {
		return fmt.Errorf("unknown proxy target: %s", targetName)
	}

	c.logger.Info("adding latency",
		zap.String("target", proxy.Name),
		zap.Int("latency_ms", latencyMs),
	)

	return c.execToxiproxyCmd(ctx, "toxic", "add",
		"-n", targetName+"_latency",
		"-t", "latency",
		"-a", fmt.Sprintf("latency=%d", latencyMs),
		proxy.Name,
	)
}

// AddBandwidthLimit restricts the proxy's throughput to the specified rate in KB/s.
func (c *Controller) AddBandwidthLimit(ctx context.Context, targetName string, rateKB int) error {
	proxy, ok := c.proxies[targetName]
	if !ok {
		return fmt.Errorf("unknown proxy target: %s", targetName)
	}

	c.logger.Info("adding bandwidth limit",
		zap.String("target", proxy.Name),
		zap.Int("rate_kb", rateKB),
	)

	return c.execToxiproxyCmd(ctx, "toxic", "add",
		"-n", targetName+"_bandwidth",
		"-t", "bandwidth",
		"-a", fmt.Sprintf("rate=%d", rateKB),
		proxy.Name,
	)
}

// AddConnectionReset simulates TCP connection resets (RST) on the proxy.
// Connections are reset after the specified timeout in milliseconds.
func (c *Controller) AddConnectionReset(
	ctx context.Context,
	targetName string,
	timeoutMs int,
) error {
	proxy, ok := c.proxies[targetName]
	if !ok {
		return fmt.Errorf("unknown proxy target: %s", targetName)
	}

	c.logger.Info("adding connection reset",
		zap.String("target", proxy.Name),
		zap.Int("timeout_ms", timeoutMs),
	)

	return c.execToxiproxyCmd(ctx, "toxic", "add",
		"-n", targetName+"_reset",
		"-t", "reset_peer",
		"-a", fmt.Sprintf("timeout=%d", timeoutMs),
		proxy.Name,
	)
}

// AddTimeout stops all data from getting through and closes the connection after
// the specified timeout in milliseconds. If timeoutMs is 0, data is dropped
// indefinitely without closing the connection (black hole / network partition).
func (c *Controller) AddTimeout(ctx context.Context, targetName string, timeoutMs int) error {
	proxy, ok := c.proxies[targetName]
	if !ok {
		return fmt.Errorf("unknown proxy target: %s", targetName)
	}

	c.logger.Info("adding timeout",
		zap.String("target", proxy.Name),
		zap.Int("timeout_ms", timeoutMs),
	)

	return c.execToxiproxyCmd(ctx, "toxic", "add",
		"-n", targetName+"_timeout",
		"-t", "timeout",
		"-a", fmt.Sprintf("timeout=%d", timeoutMs),
		proxy.Name,
	)
}

// RemoveAllToxics removes all active toxics from the named proxy,
// restoring normal network conditions for that service.
func (c *Controller) RemoveAllToxics(ctx context.Context, targetName string) error {
	c.logger.Info("removing all toxics", zap.String("target", targetName))
	// toxiproxy-cli doesn't have a bulk remove, so we remove known toxic names
	for _, suffix := range []string{"latency", "bandwidth", "reset", "timeout"} {
		toxicName := fmt.Sprintf("%s_%s", targetName, suffix)
		// Ignore errors for toxics that don't exist
		_ = c.execToxiproxyCmd(ctx, "toxic", "remove",
			"-n", toxicName,
			targetName,
		)
	}
	return nil
}

// Stop terminates the toxiproxy container.
func (c *Controller) Stop(ctx context.Context) error {
	if c.container == nil {
		return nil
	}
	return c.container.Terminate(ctx)
}

// APIURL returns the host-accessible URL for the toxiproxy HTTP API.
func (c *Controller) APIURL() string {
	return c.apiURL
}

func (c *Controller) execToxiproxyCmd(ctx context.Context, args ...string) error {
	cmd := append([]string{"/toxiproxy-cli"}, args...)
	exitCode, output, err := c.container.Exec(ctx, cmd)
	if err != nil {
		return fmt.Errorf("toxiproxy exec failed: %w", err)
	}
	if exitCode != 0 {
		var outStr string
		if output != nil {
			outBytes, _ := io.ReadAll(output)
			outStr = string(outBytes)
		}
		return fmt.Errorf("toxiproxy command %v exited with code %d: %s", args, exitCode, outStr)
	}
	return nil
}
