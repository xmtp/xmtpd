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
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
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
	proxies   sync.Map
	nextPort  int
	portsMu   sync.Mutex
}

// NewController starts a toxiproxy container on the given Docker network
// and returns a controller for managing proxies and toxics.
func NewController(
	ctx context.Context,
	logger *zap.Logger,
	id string,
) (*Controller, error) {
	c := &Controller{
		logger:   logger,
		network:  id,
		proxies:  sync.Map{},
		nextPort: baseProxyPort,
	}

	createCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:        toxiproxyImage,
		ExposedPorts: []string{fmt.Sprintf("%d/tcp", toxiproxyAPI)},
		Networks:     []string{id},
		NetworkAliases: map[string][]string{
			id: {"toxiproxy"},
		},
		Labels: map[string]string{
			"com.docker.compose.project": id,
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
			Logger:           log.New(io.Discard, "", 0),
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
	c.portsMu.Lock()
	defer c.portsMu.Unlock()
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

	c.proxies.Store(name, &ProxyTarget{
		Name:         name,
		Upstream:     upstream,
		UpstreamPort: port,
		ListenPort:   listenPort,
	})

	// Create the proxy in toxiproxy: listen on 0.0.0.0:<listenPort>,
	// forward to <upstream>:<port>
	err := c.execToxiproxyCmd(ctx, "create",
		"-l", fmt.Sprintf("0.0.0.0:%d", listenPort),
		"-u", fmt.Sprintf("%s:%d", upstream, port),
		name,
	)
	if err != nil {
		c.proxies.Delete(name)
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
func (c *Controller) ProxyAddress(name string) (string, error) {
	proxy, ok := c.proxies.Load(name)
	if !ok {
		return "", fmt.Errorf("no proxy registered for %s", name)
	}
	return fmt.Sprintf("http://toxiproxy:%d", proxy.(*ProxyTarget).ListenPort), nil
}

// addToxic is the common implementation for all toxic injection methods.
func (c *Controller) addToxic(
	ctx context.Context,
	targetName string,
	toxicSuffix string,
	toxicType string,
	attribute string,
) error {
	proxy, ok := c.proxies.Load(targetName)
	if !ok {
		return fmt.Errorf("unknown proxy target: %s", targetName)
	}

	c.logger.Info("adding toxic",
		zap.String("target", proxy.(*ProxyTarget).Name),
		zap.String("type", toxicType),
		zap.String("attribute", attribute),
	)

	return c.execToxiproxyCmd(ctx, "toxic", "add",
		"-n", targetName+"_"+toxicSuffix,
		"-t", toxicType,
		"-a", attribute,
		proxy.(*ProxyTarget).Name,
	)
}

// AddLatency injects a network latency toxic on the named proxy.
func (c *Controller) AddLatency(ctx context.Context, targetName string, latencyMs int) error {
	return c.addToxic(ctx, targetName, "latency", "latency", fmt.Sprintf("latency=%d", latencyMs))
}

// AddBandwidthLimit restricts the proxy's throughput to the specified rate in KB/s.
func (c *Controller) AddBandwidthLimit(ctx context.Context, targetName string, rateKB int) error {
	return c.addToxic(ctx, targetName, "bandwidth", "bandwidth", fmt.Sprintf("rate=%d", rateKB))
}

// AddConnectionReset simulates TCP connection resets (RST) on the proxy.
// Connections are reset after the specified timeout in milliseconds.
func (c *Controller) AddConnectionReset(
	ctx context.Context,
	targetName string,
	timeoutMs int,
) error {
	return c.addToxic(ctx, targetName, "reset", "reset_peer", fmt.Sprintf("timeout=%d", timeoutMs))
}

// AddTimeout stops all data from getting through and closes the connection after timeout.
// If timeoutMs is 0, data is dropped indefinitely (black hole / network partition).
func (c *Controller) AddTimeout(ctx context.Context, targetName string, timeoutMs int) error {
	return c.addToxic(ctx, targetName, "timeout", "timeout", fmt.Sprintf("timeout=%d", timeoutMs))
}

// DisableProxy completely disables the named proxy, refusing all connections.
// This is a stronger isolation than toxics — no data flows at all.
//
// Note: toxiproxy-cli only supports a "toggle" command (no separate enable/disable),
// so this call is non-idempotent. Callers must always pair DisableProxy with a
// corresponding EnableProxy call to restore proxy state.
func (c *Controller) DisableProxy(ctx context.Context, targetName string) error {
	proxy, ok := c.proxies.Load(targetName)
	if !ok {
		return fmt.Errorf("unknown proxy target: %s", targetName)
	}

	c.logger.Info("disabling proxy", zap.String("target", proxy.(*ProxyTarget).Name))

	return c.execToxiproxyCmd(ctx, "toggle", proxy.(*ProxyTarget).Name)
}

// EnableProxy re-enables a previously disabled proxy, restoring connectivity.
func (c *Controller) EnableProxy(ctx context.Context, targetName string) error {
	proxy, ok := c.proxies.Load(targetName)
	if !ok {
		return fmt.Errorf("unknown proxy target: %s", targetName)
	}

	c.logger.Info("enabling proxy", zap.String("target", proxy.(*ProxyTarget).Name))

	return c.execToxiproxyCmd(ctx, "toggle", proxy.(*ProxyTarget).Name)
}

// RemoveAllToxics removes all active toxics from the named proxy,
// restoring normal network conditions for that service.
// Errors from removing non-existent toxics are ignored; other errors are returned.
func (c *Controller) RemoveAllToxics(ctx context.Context, targetName string) error {
	c.logger.Info("removing all toxics", zap.String("target", targetName))
	// toxiproxy-cli doesn't have a bulk remove, so we remove known toxic names
	var errs []error
	for _, suffix := range []string{"latency", "bandwidth", "reset", "timeout"} {
		toxicName := fmt.Sprintf("%s_%s", targetName, suffix)
		err := c.execToxiproxyCmd(ctx, "toxic", "remove",
			"-n", toxicName,
			targetName,
		)
		if err != nil && !strings.Contains(err.Error(), "not found") {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
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
