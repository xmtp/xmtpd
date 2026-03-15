package tests

import (
	"context"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/e2e/types"
	"go.uber.org/zap"
)

// ChaosNetworkPartitionTest verifies that the system correctly handles a
// network partition where a node becomes unreachable by other nodes.
//
// Architecture context: Toxiproxy sits between nodes in the Docker network.
// Each node's on-chain HTTP address points to its toxiproxy proxy. When we
// add a timeout=0 toxic (black hole) on node-300's proxy:
//   - Other nodes CANNOT sync FROM node-300 (they connect via the proxy)
//   - Node-300 CAN still sync FROM others (it connects to THEIR proxies)
//   - Client publishes bypass toxiproxy (direct port mapping to host)
//
// The test:
//  1. Partitions node-300 by black-holing its proxy
//  2. Publishes envelopes to node-300 (direct, bypasses proxy)
//  3. Verifies other nodes do NOT receive node-300's envelopes (they can't sync)
//  4. Removes the partition
//  5. Verifies other nodes catch up by syncing node-300's envelopes
type ChaosNetworkPartitionTest struct{}

func NewChaosNetworkPartitionTest() *ChaosNetworkPartitionTest {
	return &ChaosNetworkPartitionTest{}
}

func (t *ChaosNetworkPartitionTest) Name() string {
	return "chaos-network-partition"
}

func (t *ChaosNetworkPartitionTest) Description() string {
	return "Black-hole a node via toxiproxy, verify isolation and recovery"
}

func (t *ChaosNetworkPartitionTest) Run(
	ctx context.Context,
	env *types.Environment,
) error {
	require := require.New(env.T())

	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddGateway(ctx))

	// Publish baseline envelopes to node-100 and let them replicate to all.
	require.NoError(env.NewClient(100))
	require.NoError(env.Client(100).PublishEnvelopes(ctx, 10))

	baselineCtx, baselineCancel := context.WithTimeout(ctx, 60*time.Second)
	defer baselineCancel()
	for _, n := range env.Nodes() {
		require.NoError(n.WaitForEnvelopes(baselineCtx, 1))
	}

	baselineCount, err := env.Node(100).GetEnvelopeCount(ctx)
	require.NoError(err)
	env.Logger.Info("baseline established",
		zap.Int64("envelope_count", baselineCount),
	)

	// --- INJECT FAULT: Disable node-300's proxy ---
	// Disabling the proxy makes toxiproxy refuse all connections to node-300.
	// This prevents other nodes from syncing FROM node-300.
	env.Logger.Info("injecting network partition on node-300")
	require.NoError(env.Node(300).DisableProxy(ctx))

	// Publish envelopes directly to node-300 (client bypasses toxiproxy).
	// Node-300 stores them locally, but other nodes can't pull them.
	require.NoError(env.NewClient(300))
	const partitionEnvelopes uint = 20
	require.NoError(env.Client(300).PublishEnvelopes(ctx, partitionEnvelopes))

	// Give sync time to attempt (and fail).
	time.Sleep(10 * time.Second)

	// --- VERIFY FAULT EFFECT ---
	// Node-300 should have the envelopes (published directly).
	node300Count, err := env.Node(300).GetEnvelopeCount(ctx)
	require.NoError(err)

	// Node-100 should NOT have node-300's new envelopes (can't sync through proxy).
	node100Count, err := env.Node(100).GetEnvelopeCount(ctx)
	require.NoError(err)

	env.Logger.Info("partition verification",
		zap.Int64("node_300_count", node300Count),
		zap.Int64("node_100_count", node100Count),
	)

	// Node-300 should have more envelopes than node-100 because it has its own
	// published envelopes that couldn't be synced out.
	require.Greater(node300Count, node100Count,
		"partitioned node should have envelopes that others can't sync "+
			"(node_300=%d, node_100=%d)",
		node300Count, node100Count,
	)

	// --- HEAL PARTITION ---
	env.Logger.Info("removing network partition from node-300")
	require.NoError(env.Node(300).EnableProxy(ctx))

	// --- VERIFY RECOVERY ---
	// After the partition is healed, other nodes should sync node-300's envelopes.
	recoveryCtx, recoveryCancel := context.WithTimeout(ctx, 120*time.Second)
	defer recoveryCancel()

	require.NoError(env.Node(100).WaitForEnvelopes(recoveryCtx, node300Count))
	require.NoError(env.Node(200).WaitForEnvelopes(recoveryCtx, node300Count))

	// All nodes should converge.
	for _, n := range env.Nodes() {
		count, err := n.GetEnvelopeCount(ctx)
		require.NoError(err)
		env.Logger.Info("final envelope count",
			zap.Uint32("node_id", n.ID()),
			zap.Int64("count", count),
		)
	}

	env.Logger.Info("chaos network partition test completed successfully")
	return nil
}

var _ types.Test = (*ChaosNetworkPartitionTest)(nil)
