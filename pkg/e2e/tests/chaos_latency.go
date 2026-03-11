package tests

import (
	"context"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/e2e/client"
	"github.com/xmtp/xmtpd/pkg/e2e/types"
	"go.uber.org/zap"
)

// ChaosLatencyTest injects network latency into a node while traffic is flowing
// and verifies:
//  1. Traffic still flows to the affected node (latency doesn't prevent delivery)
//  2. Envelopes replicate across all nodes (including the latency-affected one)
//  3. After latency is removed, the system returns to normal
type ChaosLatencyTest struct{}

func NewChaosLatencyTest() *ChaosLatencyTest {
	return &ChaosLatencyTest{}
}

func (t *ChaosLatencyTest) Name() string {
	return "chaos-latency"
}

func (t *ChaosLatencyTest) Description() string {
	return "Inject latency into a node while traffic is flowing, verify system remains functional"
}

func (t *ChaosLatencyTest) Run(ctx context.Context, env *types.Environment) error {
	require := require.New(env.T())

	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddGateway(ctx))

	require.NoError(env.NewClient(100))
	gen := env.Client(100).GenerateTraffic(ctx, client.TrafficOptions{
		BatchSize: 5,
		Duration:  60 * time.Second,
	})

	// Let some clean traffic flow first.
	time.Sleep(5 * time.Second)

	preLatencyCount, err := env.Node(100).GetEnvelopeCount(ctx)
	require.NoError(err)
	require.Positive(preLatencyCount, "should have envelopes before latency injection")

	// Inject 500ms latency to node-100 while traffic is flowing.
	env.Logger.Info("injecting 500ms latency on node-100")
	require.NoError(env.Node(100).AddLatency(ctx, 500))

	// Let traffic flow under degraded conditions.
	time.Sleep(15 * time.Second)

	// Verify traffic is still flowing despite latency.
	midLatencyCount, err := env.Node(100).GetEnvelopeCount(ctx)
	require.NoError(err)
	env.Logger.Info("mid-latency envelope count",
		zap.Int64("node_100_count", midLatencyCount),
		zap.Int64("pre_latency_count", preLatencyCount),
	)
	require.Greater(midLatencyCount, preLatencyCount,
		"node should still be receiving envelopes under latency",
	)

	// Remove the toxic while traffic is still flowing.
	env.Logger.Info("removing latency from node-100")
	require.NoError(env.Node(100).RemoveAllToxics(ctx))

	// Let traffic flow under normal conditions again.
	time.Sleep(5 * time.Second)

	gen.Stop()
	require.NoError(gen.Err(), "traffic generator should not have errored")

	// Verify all nodes have envelopes and have converged.
	checkCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	targetCount, err := env.Node(100).GetEnvelopeCount(ctx)
	require.NoError(err)
	for _, n := range env.Nodes() {
		require.NoError(n.WaitForEnvelopes(checkCtx, targetCount))
		count, err := n.GetEnvelopeCount(ctx)
		require.NoError(err)
		env.Logger.Info("final envelope count",
			zap.Uint32("node_id", n.ID()),
			zap.Int64("count", count),
		)
	}

	env.Logger.Info("chaos latency test completed successfully")
	return nil
}

var _ types.Test = (*ChaosLatencyTest)(nil)
