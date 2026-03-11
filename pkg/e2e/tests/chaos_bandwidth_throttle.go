package tests

import (
	"context"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/e2e/types"
	"go.uber.org/zap"
)

// ChaosBandwidthThrottleTest verifies graceful degradation under severe
// bandwidth constraints. It applies aggressive throttling to ALL nodes
// simultaneously, then verifies:
//
//  1. Envelopes still get through (just slowly)
//  2. No nodes crash or become permanently stuck
//  3. After throttling is removed, the system returns to normal throughput
//
// This simulates a scenario like a cloud provider network incident where
// all inter-node bandwidth is severely constrained.
type ChaosBandwidthThrottleTest struct{}

func NewChaosBandwidthThrottleTest() *ChaosBandwidthThrottleTest {
	return &ChaosBandwidthThrottleTest{}
}

func (t *ChaosBandwidthThrottleTest) Name() string {
	return "chaos-bandwidth-throttle"
}

func (t *ChaosBandwidthThrottleTest) Description() string {
	return "Throttle bandwidth on all nodes, verify graceful degradation and recovery"
}

func (t *ChaosBandwidthThrottleTest) Run(
	ctx context.Context,
	env *types.Environment,
) error {
	require := require.New(env.T())

	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddGateway(ctx))

	require.NoError(env.NewClient(100))

	// Establish baseline — publish and verify normal replication.
	require.NoError(env.Client(100).PublishEnvelopes(ctx, 10))
	baselineCtx, baselineCancel := context.WithTimeout(ctx, 60*time.Second)
	defer baselineCancel()
	for _, n := range env.Nodes() {
		require.NoError(n.WaitForEnvelopes(baselineCtx, 1))
	}

	env.Logger.Info("baseline established, applying bandwidth throttle to all nodes")

	// --- INJECT FAULT: Throttle ALL nodes to 5 KB/s ---
	for _, n := range env.Nodes() {
		require.NoError(n.AddBandwidthLimit(ctx, 5))
	}

	// Publish envelopes under throttled conditions.
	// Use a smaller batch since bandwidth is severely limited.
	require.NoError(env.Client(100).PublishEnvelopes(ctx, 5))

	// Wait for envelopes to reach at least one other node.
	// Use a generous timeout — at 5 KB/s this will be slow.
	throttledCtx, throttledCancel := context.WithTimeout(ctx, 120*time.Second)
	defer throttledCancel()

	// Node-100 should have the envelopes locally (published directly to it).
	node100Count, err := env.Node(100).GetEnvelopeCount(ctx)
	require.NoError(err)
	require.Positive(node100Count,
		"publishing node should have envelopes despite throttle",
	)

	// At least one other node should eventually get the envelopes (slowly).
	require.NoError(env.Node(200).WaitForEnvelopes(throttledCtx, 1))

	env.Logger.Info("envelopes replicated under throttle (slowly)")

	// Record throttled counts.
	throttledCounts := make(map[uint32]int64)
	for _, n := range env.Nodes() {
		count, err := n.GetEnvelopeCount(ctx)
		require.NoError(err)
		throttledCounts[n.ID()] = count
		env.Logger.Info("throttled envelope count",
			zap.Uint32("node_id", n.ID()),
			zap.Int64("count", count),
		)
	}

	// --- REMOVE FAULT ---
	env.Logger.Info("removing bandwidth throttle from all nodes")
	for _, n := range env.Nodes() {
		require.NoError(n.RemoveAllToxics(ctx))
	}

	// --- VERIFY RECOVERY ---
	// Publish more envelopes at full speed.
	require.NoError(env.Client(100).PublishEnvelopes(ctx, 20))

	recoveryCtx, recoveryCancel := context.WithTimeout(ctx, 60*time.Second)
	defer recoveryCancel()

	// All nodes should converge quickly after throttle is removed.
	targetCount := throttledCounts[uint32(100)] + 20
	for _, n := range env.Nodes() {
		require.NoError(n.WaitForEnvelopes(recoveryCtx, targetCount))
	}

	// Final check — all nodes should have the same count.
	var finalCounts []int64
	for _, n := range env.Nodes() {
		count, err := n.GetEnvelopeCount(ctx)
		require.NoError(err)
		finalCounts = append(finalCounts, count)
		env.Logger.Info("post-recovery envelope count",
			zap.Uint32("node_id", n.ID()),
			zap.Int64("count", count),
		)
	}

	// All nodes should agree (or be very close).
	for i := 1; i < len(finalCounts); i++ {
		require.Equal(finalCounts[0], finalCounts[i],
			"all nodes should converge after throttle removed",
		)
	}

	env.Logger.Info("chaos bandwidth throttle test completed successfully")
	return nil
}

var _ types.Test = (*ChaosBandwidthThrottleTest)(nil)
