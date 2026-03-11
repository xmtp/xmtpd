package tests

import (
	"context"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/e2e/client"
	"github.com/xmtp/xmtpd/pkg/e2e/types"
	"go.uber.org/zap"
)

// ChaosCompoundFaultTest applies multiple different faults across different
// nodes simultaneously to stress the system's ability to handle compound
// failure scenarios. Specifically:
//
//   - Node 200: high latency (1000ms) — simulates a geographically distant node
//   - Node 300: bandwidth throttle (10 KB/s) — simulates a congested link
//
// While these faults are active, traffic is generated and the test verifies:
//  1. Traffic still flows to the unaffected node (100)
//  2. Both degraded nodes eventually receive envelopes (even if slowly)
//  3. After all faults are removed, the cluster converges
type ChaosCompoundFaultTest struct{}

func NewChaosCompoundFaultTest() *ChaosCompoundFaultTest {
	return &ChaosCompoundFaultTest{}
}

func (t *ChaosCompoundFaultTest) Name() string {
	return "chaos-compound-fault"
}

func (t *ChaosCompoundFaultTest) Description() string {
	return "Apply latency + bandwidth faults across multiple nodes simultaneously"
}

func (t *ChaosCompoundFaultTest) Run(
	ctx context.Context,
	env *types.Environment,
) error {
	require := require.New(env.T())

	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddGateway(ctx))

	require.NoError(env.NewClient(100))

	// Establish baseline — publish some envelopes and let them replicate.
	require.NoError(env.Client(100).PublishEnvelopes(ctx, 10))
	baselineCtx, baselineCancel := context.WithTimeout(ctx, 60*time.Second)
	defer baselineCancel()
	for _, n := range env.Nodes() {
		require.NoError(n.WaitForEnvelopes(baselineCtx, 1))
	}

	env.Logger.Info("baseline established, injecting compound faults")

	// --- INJECT COMPOUND FAULTS ---
	// Node 200: 1000ms latency (high but not a partition)
	require.NoError(env.Node(200).AddLatency(ctx, 1000))
	// Node 300: 10 KB/s bandwidth limit (severely throttled)
	require.NoError(env.Node(300).AddBandwidthLimit(ctx, 10))

	// Start traffic generation while faults are active.
	gen := env.Client(100).GenerateTraffic(ctx, client.TrafficOptions{
		BatchSize: 5,
		Duration:  30 * time.Second,
	})

	// Let traffic flow under compound faults.
	time.Sleep(30 * time.Second)
	gen.Stop()

	// --- VERIFY NODE-100 (HEALTHY) RECEIVED TRAFFIC ---
	healthyCount, err := env.Node(100).GetEnvelopeCount(ctx)
	require.NoError(err)
	env.Logger.Info("healthy node envelope count",
		zap.Int64("node_100_count", healthyCount),
	)
	require.Greater(healthyCount, int64(10),
		"healthy node should have received traffic beyond baseline",
	)

	// --- CHECK DEGRADED NODES ---
	// They should have SOME envelopes (latency and throttling don't prevent
	// delivery, they just slow it down). But they may lag behind node-100.
	latencyNodeCount, err := env.Node(200).GetEnvelopeCount(ctx)
	require.NoError(err)
	throttledNodeCount, err := env.Node(300).GetEnvelopeCount(ctx)
	require.NoError(err)

	env.Logger.Info("degraded node envelope counts during faults",
		zap.Int64("node_200_latency", latencyNodeCount),
		zap.Int64("node_300_throttled", throttledNodeCount),
	)

	// Both degraded nodes should have at least some envelopes.
	require.Positive(latencyNodeCount,
		"latency-affected node should have some envelopes",
	)
	require.Positive(throttledNodeCount,
		"throttled node should have some envelopes",
	)

	// --- REMOVE ALL FAULTS ---
	env.Logger.Info("removing all faults")
	require.NoError(env.Node(200).RemoveAllToxics(ctx))
	require.NoError(env.Node(300).RemoveAllToxics(ctx))

	// --- VERIFY CONVERGENCE ---
	// After faults are removed, all nodes should converge.
	convergenceCtx, convergenceCancel := context.WithTimeout(ctx, 120*time.Second)
	defer convergenceCancel()

	require.NoError(env.Node(200).WaitForEnvelopes(convergenceCtx, healthyCount))
	require.NoError(env.Node(300).WaitForEnvelopes(convergenceCtx, healthyCount))

	// Final verification — all nodes should agree.
	for _, n := range env.Nodes() {
		count, err := n.GetEnvelopeCount(ctx)
		require.NoError(err)
		env.Logger.Info("final envelope count",
			zap.Uint32("node_id", n.ID()),
			zap.Int64("count", count),
		)
		require.GreaterOrEqual(count, healthyCount,
			"node %d should have converged to at least %d envelopes",
			n.ID(), healthyCount,
		)
	}

	env.Logger.Info("chaos compound fault test completed successfully")
	return nil
}

var _ types.Test = (*ChaosCompoundFaultTest)(nil)
