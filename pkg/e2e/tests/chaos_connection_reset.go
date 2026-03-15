package tests

import (
	"context"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/e2e/client"
	"github.com/xmtp/xmtpd/pkg/e2e/types"
	"go.uber.org/zap"
)

// ChaosConnectionResetTest injects TCP connection resets (RST) on a node while
// traffic is flowing and verifies:
//
//  1. Traffic generation continues despite RSTs (the client retries or the
//     remaining nodes still accept traffic)
//  2. After RSTs are removed, the affected node recovers and catches up
//  3. All nodes eventually converge on envelope count
//
// This is a more aggressive fault than latency — it actively kills TCP connections,
// forcing reconnection and retry logic to engage.
type ChaosConnectionResetTest struct{}

func NewChaosConnectionResetTest() *ChaosConnectionResetTest {
	return &ChaosConnectionResetTest{}
}

func (t *ChaosConnectionResetTest) Name() string {
	return "chaos-connection-reset"
}

func (t *ChaosConnectionResetTest) Description() string {
	return "Inject TCP RSTs under load, verify system recovers and nodes converge"
}

func (t *ChaosConnectionResetTest) Run(
	ctx context.Context,
	env *types.Environment,
) error {
	require := require.New(env.T())

	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddGateway(ctx))

	require.NoError(env.NewClient(100))

	// Start background traffic.
	gen := env.Client(100).GenerateTraffic(ctx, client.TrafficOptions{
		BatchSize: 5,
		Duration:  45 * time.Second,
	})

	// Let clean traffic flow first so we have a baseline.
	time.Sleep(5 * time.Second)

	preResetCount, err := env.Node(200).GetEnvelopeCount(ctx)
	require.NoError(err)
	env.Logger.Info("pre-reset baseline",
		zap.Int64("node_200_envelopes", preResetCount),
	)
	require.Positive(preResetCount, "should have envelopes before reset injection")

	// --- INJECT FAULT: TCP RSTs on node-200 ---
	// Connections are reset after 500ms, forcing reconnection.
	env.Logger.Info("injecting connection resets on node-200")
	require.NoError(env.Node(200).AddConnectionReset(ctx, 500))

	// Let traffic flow under RST conditions for 15 seconds.
	// The traffic generator publishes to node-100, so it should continue working.
	// But node-200's sync connections will get RST'd.
	time.Sleep(15 * time.Second)

	// --- VERIFY FAULT EFFECT ---
	// Node-200 may have fallen behind due to RSTs disrupting sync.
	midFaultCount200, err := env.Node(200).GetEnvelopeCount(ctx)
	require.NoError(err)
	midFaultCount100, err := env.Node(100).GetEnvelopeCount(ctx)
	require.NoError(err)

	env.Logger.Info("mid-fault state",
		zap.Int64("node_100_envelopes", midFaultCount100),
		zap.Int64("node_200_envelopes", midFaultCount200),
	)

	// Traffic should still be flowing (node-100 should have more envelopes).
	require.Greater(midFaultCount100, preResetCount,
		"node-100 should still be accepting traffic during RSTs on node-200",
	)

	// --- REMOVE FAULT ---
	env.Logger.Info("removing connection resets from node-200")
	require.NoError(env.Node(200).RemoveAllToxics(ctx))

	// Let remaining traffic flow cleanly.
	time.Sleep(10 * time.Second)

	// Stop traffic generator and check for errors.
	gen.Stop()
	// Traffic generator may have encountered transient errors from RSTs —
	// this is acceptable as long as the system recovered.

	// --- VERIFY RECOVERY ---
	// Wait for node-200 to catch up with node-100 via sync.
	recoveryCtx, recoveryCancel := context.WithTimeout(ctx, 90*time.Second)
	defer recoveryCancel()

	// Get target count from node-100.
	targetCount, err := env.Node(100).GetEnvelopeCount(ctx)
	require.NoError(err)

	require.NoError(env.Node(200).WaitForEnvelopes(recoveryCtx, targetCount))

	// All three nodes should converge.
	for _, n := range env.Nodes() {
		count, err := n.GetEnvelopeCount(ctx)
		require.NoError(err)
		env.Logger.Info("final envelope count",
			zap.Uint32("node_id", n.ID()),
			zap.Int64("count", count),
		)
		require.GreaterOrEqual(count, targetCount,
			"node %d should have at least %d envelopes",
			n.ID(), targetCount,
		)
	}

	env.Logger.Info("chaos connection reset test completed successfully")
	return nil
}

var _ types.Test = (*ChaosConnectionResetTest)(nil)
