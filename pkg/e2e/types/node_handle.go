package types

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/e2e/node"
	"github.com/xmtp/xmtpd/pkg/e2e/observe"
	"github.com/xmtp/xmtpd/pkg/utils"
)

// NodeHandle provides a fluent API for interacting with an xmtpd node.
// It wraps the underlying node container and composes lifecycle management,
// chaos injection, and database observation into a single addressable handle.
//
// Nodes are addressed by their on-chain nodeID (100, 200, 300, ...):
//
//	env.Node(100).AddLatency(ctx, 500)
//	env.Node(200).Stop(ctx)
//	env.Node(100).WaitForEnvelopes(ctx, 10)
type NodeHandle struct {
	node *node.Node
	env  *Environment
}

// newNodeHandle creates a NodeHandle wrapping the given node with references
// to the environment's chaos controller and observer.
func newNodeHandle(n *node.Node, env *Environment) *NodeHandle {
	return &NodeHandle{
		node: n,
		env:  env,
	}
}

// --- Identity ---

// ID returns the on-chain nodeID assigned by the NodeRegistry contract.
func (h *NodeHandle) ID() uint32 {
	return h.node.NodeID()
}

// Alias returns the container's network alias (e.g. "node-100").
func (h *NodeHandle) Alias() string {
	return h.node.Alias()
}

// Endpoint returns the host-accessible gRPC address (e.g. "http://localhost:XXXXX").
// Use this to create clients that publish to this node.
func (h *NodeHandle) Endpoint() string {
	return h.node.ExternalAddr()
}

// Address returns the Ethereum address derived from this node's signer key.
// This is the node owner's address used for on-chain operations and fee distribution.
func (h *NodeHandle) Address() common.Address {
	privKey, err := utils.ParseEcdsaPrivateKey(h.node.SignerKey())
	if err != nil {
		panic("failed to parse node signer key: " + err.Error())
	}
	return crypto.PubkeyToAddress(privKey.PublicKey)
}

// SignerKey returns the private key hex string used by this node for signing.
// This is also the key associated with the node's owner address on-chain.
func (h *NodeHandle) SignerKey() string {
	return h.node.SignerKey()
}

// DBConnectionString returns the Postgres connection string for this node's database.
func (h *NodeHandle) DBConnectionString() string {
	return h.node.DBConnectionString()
}

// --- Lifecycle ---

// Stop terminates the node's container. The node can be restarted with Start.
func (h *NodeHandle) Stop(ctx context.Context) error {
	return h.node.Stop(ctx)
}

// Start restarts a previously stopped node by creating a new container
// with the same configuration. The node retains its nodeID and alias.
func (h *NodeHandle) Start(ctx context.Context) error {
	n, err := h.node.Restart(ctx)
	if err != nil {
		return err
	}
	h.node = n
	return nil
}

// --- On-chain network management ---

// AddToCanonicalNetwork adds this node to the on-chain canonical network,
// allowing it to participate in replication and consensus.
func (h *NodeHandle) AddToCanonicalNetwork(ctx context.Context) error {
	return h.env.AddNodeToCanonicalNetwork(ctx, h.node.NodeID())
}

// RemoveFromCanonicalNetwork removes this node from the on-chain canonical network.
// The node container keeps running but is excluded from the active set.
func (h *NodeHandle) RemoveFromCanonicalNetwork(ctx context.Context) error {
	return h.env.RemoveNodeFromCanonicalNetwork(ctx, h.node.NodeID())
}

// --- Chaos injection ---
// All chaos methods delegate to the toxiproxy controller using this node's alias.
// Toxics affect all network traffic to/from the node container, including traffic
// from clients publishing to it.

// AddLatency injects a network latency toxic on this node.
// All connections to/from the node will experience the specified delay in milliseconds.
func (h *NodeHandle) AddLatency(ctx context.Context, ms int) error {
	return h.env.Chaos.AddLatency(ctx, h.node.Alias(), ms)
}

// AddBandwidthLimit restricts the node's network throughput to the specified rate in KB/s.
func (h *NodeHandle) AddBandwidthLimit(ctx context.Context, kbps int) error {
	return h.env.Chaos.AddBandwidthLimit(ctx, h.node.Alias(), kbps)
}

// AddConnectionReset simulates TCP connection resets (RST) on the node's connections.
// Connections are reset after the specified timeout in milliseconds.
func (h *NodeHandle) AddConnectionReset(ctx context.Context, timeoutMs int) error {
	return h.env.Chaos.AddConnectionReset(ctx, h.node.Alias(), timeoutMs)
}

// AddTimeout blocks all data and closes connections after the specified timeout.
// If timeoutMs is 0, data is dropped indefinitely without closing the connection,
// effectively simulating a network partition (black hole).
func (h *NodeHandle) AddTimeout(ctx context.Context, timeoutMs int) error {
	return h.env.Chaos.AddTimeout(ctx, h.node.Alias(), timeoutMs)
}

// RemoveAllToxics removes all active toxics from this node, restoring normal
// network conditions.
func (h *NodeHandle) RemoveAllToxics(ctx context.Context) error {
	return h.env.Chaos.RemoveAllToxics(ctx, h.node.Alias())
}

// --- Balance queries ---

// GetFeeTokenBalance returns the fee token (xUSD) balance for this node's owner address.
func (h *NodeHandle) GetFeeTokenBalance(ctx context.Context) (*big.Int, error) {
	return h.env.GetFeeTokenBalance(ctx, h.Address())
}

// GetGasBalance returns the native ETH balance for this node's owner address.
func (h *NodeHandle) GetGasBalance(ctx context.Context) (*big.Int, error) {
	return h.env.GetGasBalance(ctx, h.Address())
}

// DisableProxy completely disables this node's proxy, refusing all connections.
// This is a stronger isolation than toxics — no data flows at all.
func (h *NodeHandle) DisableProxy(ctx context.Context) error {
	if h.env.Chaos == nil {
		return errors.New("chaos controller not available")
	}
	return h.env.Chaos.DisableProxy(ctx, h.node.Alias())
}

// EnableProxy re-enables this node's proxy after it was disabled.
func (h *NodeHandle) EnableProxy(ctx context.Context) error {
	if h.env.Chaos == nil {
		return errors.New("chaos controller not available")
	}
	return h.env.Chaos.EnableProxy(ctx, h.node.Alias())
}

// --- Observer (database queries) ---
// All observer methods query this node's own Postgres database.
// The node handle provides the connection string automatically.

// GetEnvelopeCount returns the total number of envelopes stored in this node's database.
func (h *NodeHandle) GetEnvelopeCount(ctx context.Context) (int64, error) {
	return h.env.Observer().GetEnvelopeCount(ctx, h.node.DBConnectionString())
}

// GetVectorClock returns the vector clock entries from this node's database,
// showing the latest sequence ID received from each originator node.
func (h *NodeHandle) GetVectorClock(ctx context.Context) ([]observe.VectorClockEntry, error) {
	return h.env.Observer().GetVectorClock(ctx, h.node.DBConnectionString())
}

// GetStagedEnvelopeCount returns the number of staged originator envelopes
// waiting to be processed in this node's database.
func (h *NodeHandle) GetStagedEnvelopeCount(ctx context.Context) (int64, error) {
	return h.env.Observer().GetStagedEnvelopeCount(ctx, h.node.DBConnectionString())
}

// GetPayerReportCount returns the total number of payer reports in this node's database.
func (h *NodeHandle) GetPayerReportCount(ctx context.Context) (int64, error) {
	return h.env.Observer().GetPayerReportCount(ctx, h.node.DBConnectionString())
}

// GetPayerReportStatusCounts returns a breakdown of payer report statuses
// (attestation and submission) from this node's database.
func (h *NodeHandle) GetPayerReportStatusCounts(
	ctx context.Context,
) (*observe.PayerReportStatusCounts, error) {
	return h.env.Observer().GetPayerReportStatusCounts(ctx, h.node.DBConnectionString())
}

// GetUnsettledUsage returns per-payer spending stats for unsettled usage
// from this node's database.
func (h *NodeHandle) GetUnsettledUsage(ctx context.Context) ([]observe.PayerUsageStats, error) {
	return h.env.Observer().GetUnsettledUsage(ctx, h.node.DBConnectionString())
}

// GetSettledPayerReports returns settled payer reports with their originator node ID
// and submitted report index, needed for claiming from the DistributionManager.
func (h *NodeHandle) GetSettledPayerReports(
	ctx context.Context,
) ([]observe.SettledPayerReport, error) {
	return h.env.Observer().GetSettledPayerReports(ctx, h.node.DBConnectionString())
}

// GetNodeInfo returns the node_id stored in this node's database.
func (h *NodeHandle) GetNodeInfo(ctx context.Context) (int32, error) {
	return h.env.Observer().GetNodeInfo(ctx, h.node.DBConnectionString())
}

// WaitForEnvelopes polls this node's database until at least minCount envelopes
// are present, or the context is cancelled. Polls every 2 seconds.
func (h *NodeHandle) WaitForEnvelopes(ctx context.Context, minCount int64) error {
	return h.env.Observer().WaitForEnvelopes(ctx, h.node.DBConnectionString(), minCount)
}

// WaitForPayerReports polls this node's database until the checkFn returns true
// for the current payer report status counts, or the context is cancelled.
// The description is used in timeout error messages for debugging.
func (h *NodeHandle) WaitForPayerReports(
	ctx context.Context,
	checkFn func(*observe.PayerReportStatusCounts) bool,
	description string,
) error {
	return h.env.Observer().
		WaitForPayerReports(ctx, h.node.DBConnectionString(), checkFn, description)
}
