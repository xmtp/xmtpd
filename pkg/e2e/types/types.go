// Package types defines shared types for the E2E test framework.
package types

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/xmtp/xmtpd/pkg/e2e/chain"
	"github.com/xmtp/xmtpd/pkg/e2e/chaos"
	"github.com/xmtp/xmtpd/pkg/e2e/client"
	"github.com/xmtp/xmtpd/pkg/e2e/gateway"
	"github.com/xmtp/xmtpd/pkg/e2e/keys"
	"github.com/xmtp/xmtpd/pkg/e2e/node"
	"github.com/xmtp/xmtpd/pkg/e2e/observe"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

const (
	// cliImage is used to register nodes on-chain before they start.
	cliImage = "ghcr.io/xmtp/xmtpd-cli:main"
)

// Config holds the configuration for an E2E test run.
type Config struct {
	ChainImage   string
	XmtpdImage   string
	GatewayImage string
	TestFilter   []string
	OutputFormat string
}

// Test is the interface that all E2E tests must implement.
// The runner calls Run with a fully initialized environment.
// Tests use env.T() with testify's require/assert for assertions.
type Test interface {
	// Name returns a short, unique identifier for the test (e.g. "smoke", "chaos-latency").
	Name() string
	// Description returns a human-readable summary of what the test verifies.
	Description() string
	// Run executes the test. Return nil on success or an error on failure.
	// Use require/assert with env.T() for assertions — FailNow panics are
	// caught by the runner and reported as test failures.
	Run(ctx context.Context, env *Environment) error
}

// Environment is the central context for an E2E test. It provides access to
// all infrastructure (nodes, gateways, chain, chaos) and handles for interacting
// with them. The runner creates a fresh environment for each test.
type Environment struct {
	// Logger is the structured logger for this test run.
	Logger *zap.Logger
	// Config holds the test run configuration.
	Config Config
	// Chain provides access to the Anvil blockchain used for on-chain operations.
	Chain *chain.Chain
	// Chaos is the toxiproxy controller for injecting network faults.
	Chaos *chaos.Controller
	// Keys manages private key allocation across roles (admin, client, node, gateway).
	Keys *keys.Manager
	// Network is the Docker network name shared by all containers.
	Network string

	// observer is the database observer for querying node databases.
	// Accessed through NodeHandle methods rather than directly.
	observer *observe.Observer

	Redis       testcontainers.Container
	cleanupFunc func(ctx context.Context) error

	// contracts holds lazily initialized blockchain clients for direct contract calls.
	contracts        *chainClients
	contractsOnce    sync.Once
	contractsInitErr error

	// t is the TestingT adapter set by the runner before each test.
	t *TestingT

	// Node tracking
	nodes      []*NodeHandle
	nodesByID  map[uint32]*NodeHandle
	nextNodeID uint32
	nodesMu    sync.Mutex

	// Gateway tracking
	gateways    []*GatewayHandle
	nextGWIndex int
	gatewaysMu  sync.Mutex

	// Client tracking
	clients   map[uint32]client.Client
	clientsMu sync.Mutex
}

// nodeConfig holds options for AddNode, configured via NodeOption functions.
type nodeConfig struct {
	alias   string
	image   string
	envVars map[string]string
}

// NodeOption configures optional parameters for AddNode.
type NodeOption func(*nodeConfig)

// WithAlias overrides the auto-generated container alias for a node.
// By default, aliases are derived from the nodeID: "node-100", "node-200", etc.
func WithAlias(alias string) NodeOption {
	return func(c *nodeConfig) {
		c.alias = alias
	}
}

// WithNodeImage overrides the default xmtpd Docker image for this node.
// Use this to test a specific version or a locally built image.
func WithNodeImage(image string) NodeOption {
	return func(c *nodeConfig) {
		c.image = image
	}
}

// WithNodeEnvVars sets additional environment variables on the node container.
// These are merged with (and override) the default environment variables.
func WithNodeEnvVars(vars map[string]string) NodeOption {
	return func(c *nodeConfig) {
		c.envVars = vars
	}
}

// gatewayConfig holds options for AddGateway, configured via GatewayOption functions.
type gatewayConfig struct {
	alias   string
	image   string
	envVars map[string]string
}

// GatewayOption configures optional parameters for AddGateway.
type GatewayOption func(*gatewayConfig)

// WithGatewayAlias overrides the auto-generated container alias for a gateway.
// By default, aliases are "gateway-0", "gateway-1", etc.
func WithGatewayAlias(alias string) GatewayOption {
	return func(c *gatewayConfig) {
		c.alias = alias
	}
}

// WithGatewayImage overrides the default gateway Docker image for this gateway.
// Use this to test a specific version or a locally built image.
func WithGatewayImage(image string) GatewayOption {
	return func(c *gatewayConfig) {
		c.image = image
	}
}

// WithGatewayEnvVars sets additional environment variables on the gateway container.
func WithGatewayEnvVars(vars map[string]string) GatewayOption {
	return func(c *gatewayConfig) {
		c.envVars = vars
	}
}

// SetTestingT sets the TestingT adapter for this environment. Called by the runner
// before each test — test authors should not call this directly.
func (e *Environment) SetTestingT(t *TestingT) {
	e.t = t
}

// T returns the TestingT adapter for use with testify's require and assert packages.
// The runner injects this automatically before each test run.
//
// Example:
//
//	require := require.New(env.T())
//	require.NoError(env.AddNode(ctx))
func (e *Environment) T() *TestingT {
	return e.t
}

// SetObserver sets the database observer. Called by the runner during environment setup.
func (e *Environment) SetObserver(obs *observe.Observer) {
	e.observer = obs
}

// Observer returns the database observer for direct queries when needed.
// Prefer using NodeHandle observer methods (e.g. env.Node(100).GetEnvelopeCount)
// which automatically provide the node's connection string.
func (e *Environment) Observer() *observe.Observer {
	return e.observer
}

func (e *Environment) SetCleanupFunc(fn func(ctx context.Context) error) {
	e.cleanupFunc = fn
}

func (e *Environment) Cleanup(ctx context.Context) error {
	if e.cleanupFunc != nil {
		return e.cleanupFunc(ctx)
	}
	return nil
}

// --- Node management ---

// AddNode registers a new node on-chain and starts its container.
// The node gets the next available nodeID (100, 200, 300, ...) from the
// NodeRegistry contract. The alias is auto-generated as "node-{nodeID}"
// unless overridden with WithAlias.
//
// After AddNode returns, the node is accessible via env.Node(nodeID).
func (e *Environment) AddNode(ctx context.Context, opts ...NodeOption) error {
	cfg := &nodeConfig{}
	for _, o := range opts {
		o(cfg)
	}

	signerKey, err := e.Keys.NextNodeKey(ctx)
	if err != nil {
		return fmt.Errorf("failed to allocate node signer key: %w", err)
	}

	// Determine alias — will be finalized after we know the nodeID
	tempAlias := cfg.alias

	// Register the node on-chain before starting
	nodeID, err := e.registerNode(ctx, signerKey)
	if err != nil {
		return fmt.Errorf("failed to register node on-chain: %w", err)
	}

	// Auto-generate alias from nodeID if not overridden
	alias := tempAlias
	if alias == "" {
		alias = fmt.Sprintf("node-%d", nodeID)
	}

	nodeImage := e.Config.XmtpdImage
	if cfg.image != "" {
		nodeImage = cfg.image
	}

	nodeOpts := node.Options{
		Image:     nodeImage,
		Network:   e.Network,
		Alias:     alias,
		WsURL:     e.Chain.InternalWsURL(),
		RPCURL:    e.Chain.InternalRPCURL(),
		SignerKey: signerKey,
		EnvVars:   cfg.envVars,
	}

	// Reset the node's database to ensure clean state (host Postgres persists across runs)
	if err := e.resetNodeDB(alias); err != nil {
		e.Logger.Warn("failed to reset node database", zap.String("alias", alias), zap.Error(err))
	}

	n, err := node.New(ctx, e.Logger.Named(alias), nodeOpts)
	if err != nil {
		return fmt.Errorf("failed to start node %s: %w", alias, err)
	}
	n.SetNodeID(nodeID)

	handle := newNodeHandle(n, e)

	e.nodesMu.Lock()
	e.nodes = append(e.nodes, handle)
	if e.nodesByID == nil {
		e.nodesByID = make(map[uint32]*NodeHandle)
	}
	e.nodesByID[nodeID] = handle
	e.nodesMu.Unlock()

	return nil
}

// Node returns the NodeHandle for the node with the given on-chain nodeID.
// Panics if no node with that ID exists — use this after AddNode has succeeded.
//
// Example:
//
//	env.Node(100).AddLatency(ctx, 500)
//	env.Node(200).Stop(ctx)
//	env.Node(100).WaitForEnvelopes(ctx, 10)
func (e *Environment) Node(nodeID uint32) *NodeHandle {
	e.nodesMu.Lock()
	defer e.nodesMu.Unlock()
	h, ok := e.nodesByID[nodeID]
	if !ok {
		panic(fmt.Sprintf("no node with ID %d", nodeID))
	}
	return h
}

// Nodes returns all registered node handles in creation order.
func (e *Environment) Nodes() []*NodeHandle {
	e.nodesMu.Lock()
	defer e.nodesMu.Unlock()
	result := make([]*NodeHandle, len(e.nodes))
	copy(result, e.nodes)
	return result
}

// --- Gateway management ---

// AddGateway starts a new gateway container. Gateways are indexed 0, 1, 2, ...
// in creation order. The alias is auto-generated as "gateway-{index}" unless
// overridden with WithGatewayAlias.
//
// After AddGateway returns, the gateway is accessible via env.Gateway(index).
func (e *Environment) AddGateway(ctx context.Context, opts ...GatewayOption) error {
	cfg := &gatewayConfig{}
	for _, o := range opts {
		o(cfg)
	}

	e.gatewaysMu.Lock()
	idx := e.nextGWIndex
	e.nextGWIndex++
	e.gatewaysMu.Unlock()

	alias := cfg.alias
	if alias == "" {
		alias = fmt.Sprintf("gateway-%d", idx)
	}

	gwImage := e.Config.GatewayImage
	if cfg.image != "" {
		gwImage = cfg.image
	}

	gwKey, err := e.Keys.NextGatewayKey(ctx)
	if err != nil {
		return fmt.Errorf("failed to allocate gateway key: %w", err)
	}

	// Reset the gateway's database to ensure clean state
	if err := e.resetNodeDB(alias); err != nil {
		e.Logger.Warn(
			"failed to reset gateway database",
			zap.String("alias", alias),
			zap.Error(err),
		)
	}

	gw, err := gateway.New(ctx, e.Logger.Named(alias), gateway.Options{
		Image:        gwImage,
		Network:      e.Network,
		Alias:        alias,
		WsURL:        e.Chain.InternalWsURL(),
		RPCURL:       e.Chain.InternalRPCURL(),
		SignerKey:    gwKey,
		ChaosControl: e.Chaos,
		EnvVars:      cfg.envVars,
	})
	if err != nil {
		return fmt.Errorf("failed to start gateway %s: %w", alias, err)
	}

	handle := newGatewayHandle(gw, idx, e.Chaos)

	e.gatewaysMu.Lock()
	e.gateways = append(e.gateways, handle)
	e.gatewaysMu.Unlock()

	return nil
}

// Gateway returns the GatewayHandle at the given creation index (0, 1, 2, ...).
// Panics if no gateway with that index exists.
//
// Example:
//
//	env.Gateway(0).AddLatency(ctx, 500)
//	env.Gateway(1).Stop(ctx)
func (e *Environment) Gateway(index int) *GatewayHandle {
	e.gatewaysMu.Lock()
	defer e.gatewaysMu.Unlock()
	if index < 0 || index >= len(e.gateways) {
		panic(fmt.Sprintf("no gateway with index %d", index))
	}
	return e.gateways[index]
}

// Gateways returns all registered gateway handles in creation order.
func (e *Environment) Gateways() []*GatewayHandle {
	e.gatewaysMu.Lock()
	defer e.gatewaysMu.Unlock()
	result := make([]*GatewayHandle, len(e.gateways))
	copy(result, e.gateways)
	return result
}

// --- Client management ---

// NewClient creates a traffic generation client bound to the node with the given
// nodeID. The client is registered in the environment and accessible via Client(nodeID).
// Only one client per nodeID is allowed; creating a second client for the same
// nodeID replaces the previous one (stopping its traffic first).
//
// Example:
//
//	env.NewClient(100)
//	env.Client(100).PublishEnvelopes(ctx, 10)
func (e *Environment) NewClient(nodeID uint32) error {
	n := e.Node(nodeID) // panics if node doesn't exist

	c := client.New(e.Logger.Named(fmt.Sprintf("client-%d", nodeID)), client.Options{
		NodeAddr:     n.Address(),
		PayerKey:     keys.ClientKey(),
		OriginatorID: nodeID,
	})

	e.clientsMu.Lock()
	defer e.clientsMu.Unlock()
	if e.clients == nil {
		e.clients = make(map[uint32]client.Client)
	}
	// Stop existing client if any
	if existing, ok := e.clients[nodeID]; ok {
		existing.Stop()
	}
	e.clients[nodeID] = c

	return nil
}

// Client returns the traffic client bound to the node with the given nodeID.
// Panics if no client for that nodeID has been created via NewClient.
//
// Example:
//
//	env.Client(100).PublishEnvelopes(ctx, 10)
//	env.Client(100).GenerateTraffic(ctx, client.TrafficOptions{BatchSize: 5, Duration: 60*time.Second})
//	env.Client(100).Stop()
func (e *Environment) Client(nodeID uint32) client.Client {
	e.clientsMu.Lock()
	defer e.clientsMu.Unlock()
	c, ok := e.clients[nodeID]
	if !ok {
		panic(fmt.Sprintf("no client for node %d — call env.NewClient(%d) first", nodeID, nodeID))
	}
	return c
}

// Clients returns all registered clients.
func (e *Environment) Clients() []client.Client {
	e.clientsMu.Lock()
	defer e.clientsMu.Unlock()
	result := make([]client.Client, 0, len(e.clients))
	for _, c := range e.clients {
		result = append(result, c)
	}
	return result
}

// --- On-chain operations ---

func (e *Environment) allocateNodeID() uint32 {
	e.nodesMu.Lock()
	defer e.nodesMu.Unlock()
	if e.nextNodeID == 0 {
		e.nextNodeID = 100
	}
	id := e.nextNodeID
	e.nextNodeID += 100
	return id
}

// registerNode registers and enables a node on-chain using the CLI container.
// Returns the allocated node ID.
func (e *Environment) registerNode(ctx context.Context, signerKey string) (uint32, error) {
	privateKey, err := utils.ParseEcdsaPrivateKey(signerKey)
	if err != nil {
		return 0, fmt.Errorf("failed to parse signer key: %w", err)
	}

	nodeID := e.allocateNodeID()
	pubKeyHex := "0x" + hex.EncodeToString(crypto.CompressPubkey(&privateKey.PublicKey))
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	alias := fmt.Sprintf("node-%d", nodeID)

	// If chaos is enabled, register a toxiproxy proxy for this node so that
	// inter-node traffic flows through toxiproxy. The on-chain HTTP address
	// points to the proxy instead of the real container.
	httpAddress := fmt.Sprintf("http://%s:5050", alias)
	if e.Chaos != nil {
		if proxyErr := e.Chaos.RegisterTarget(ctx, alias, alias, 5050); proxyErr != nil {
			return 0, fmt.Errorf("failed to register chaos proxy for %s: %w", alias, proxyErr)
		}
		httpAddress = e.Chaos.ProxyAddress(alias)
	}
	rpcURL := e.Chain.InternalRPCURL()

	e.Logger.Info("registering node on-chain",
		zap.Uint32("node_id", nodeID),
		zap.String("address", address),
		zap.String("http_address", httpAddress),
	)

	// Register the node
	registerCmd := []string{
		"--environment=anvil",
		"--private-key=" + keys.AdminKey(),
		"--settlement-rpc-url=" + rpcURL,
		"nodes", "register",
		"--owner-address=" + address,
		"--signing-key-pub=" + pubKeyHex,
		"--http-address=" + httpAddress,
	}

	if err := e.runCLI(ctx, registerCmd); err != nil {
		return 0, fmt.Errorf("register node failed: %w", err)
	}

	// Enable the node in the canonical network
	if err := e.AddNodeToCanonicalNetwork(ctx, nodeID); err != nil {
		return 0, err
	}

	e.Logger.Info("node registered and enabled",
		zap.Uint32("node_id", nodeID),
	)

	return nodeID, nil
}

// AddNodeToCanonicalNetwork adds a node to the canonical network by its node ID.
func (e *Environment) AddNodeToCanonicalNetwork(ctx context.Context, nodeID uint32) error {
	cmd := []string{
		"--environment=anvil",
		"--private-key=" + keys.AdminKey(),
		"--settlement-rpc-url=" + e.Chain.InternalRPCURL(),
		"nodes", "canonical-network",
		"--add",
		fmt.Sprintf("--node-id=%d", nodeID),
	}
	if err := e.runCLI(ctx, cmd); err != nil {
		return fmt.Errorf("add node %d to canonical network failed: %w", nodeID, err)
	}
	e.Logger.Info("node added to canonical network", zap.Uint32("node_id", nodeID))
	return nil
}

// RemoveNodeFromCanonicalNetwork removes a node from the canonical network by its node ID.
func (e *Environment) RemoveNodeFromCanonicalNetwork(ctx context.Context, nodeID uint32) error {
	cmd := []string{
		"--environment=anvil",
		"--private-key=" + keys.AdminKey(),
		"--settlement-rpc-url=" + e.Chain.InternalRPCURL(),
		"nodes", "canonical-network",
		"--remove",
		fmt.Sprintf("--node-id=%d", nodeID),
	}
	if err := e.runCLI(ctx, cmd); err != nil {
		return fmt.Errorf("remove node %d from canonical network failed: %w", nodeID, err)
	}
	e.Logger.Info("node removed from canonical network", zap.Uint32("node_id", nodeID))
	return nil
}

// resetNodeDB drops and recreates the per-node database to ensure a clean state.
// This is necessary because the host Postgres persists across test runs, while the
// anvil chain is fresh each time. Stale DB state (e.g. settled payer reports) would
// cause on-chain queries to fail.
func (e *Environment) resetNodeDB(alias string) error {
	dbName := "e2e_" + strings.ReplaceAll(alias, "-", "_")
	adminConnStr := "postgres://postgres:xmtp@localhost:8765/postgres?sslmode=disable"

	db, err := sql.Open("postgres", adminConnStr)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres for DB reset: %w", err)
	}
	defer func() {
		_ = db.Close()
	}()

	// Terminate active connections and drop the database
	_, _ = db.Exec(fmt.Sprintf(
		"SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '%s'",
		dbName,
	))
	_, _ = db.Exec("DROP DATABASE IF EXISTS " + dbName)

	e.Logger.Info("reset node database", zap.String("db_name", dbName))
	return nil
}

// runCLI runs an xmtpd-cli container with the given command and waits for it to exit.
func (e *Environment) runCLI(ctx context.Context, cmd []string) error {
	createCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:    cliImage,
		Networks: []string{e.Network},
		Cmd:      cmd,
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.ExtraHosts = append(hc.ExtraHosts, "host.docker.internal:host-gateway")
		},
		WaitingFor: wait.ForExit(),
	}

	cliContainer, err := testcontainers.GenericContainer(
		createCtx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Logger:           log.Default(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to start cli container: %w", err)
	}
	defer func() {
		_ = cliContainer.Terminate(ctx)
	}()

	state, err := cliContainer.State(createCtx)
	if err != nil {
		return fmt.Errorf("failed to get cli container state: %w", err)
	}

	if state.ExitCode != 0 {
		if logs, logErr := cliContainer.Logs(createCtx); logErr == nil {
			logBytes, _ := io.ReadAll(logs)
			_ = logs.Close()
			e.Logger.Error("cli container failed",
				zap.Int("exit_code", state.ExitCode),
				zap.String("logs", strings.TrimSpace(string(logBytes))),
			)
		}
		return fmt.Errorf("cli exited with code %d", state.ExitCode)
	}

	return nil
}
