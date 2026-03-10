package tests

import (
	"context"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/e2e/client"
	"github.com/xmtp/xmtpd/pkg/e2e/types"
)

// ChaosNodeDownTest stops a node while traffic is flowing, restarts it,
// and verifies it catches up with the rest of the cluster.
type ChaosNodeDownTest struct{}

func NewChaosNodeDownTest() *ChaosNodeDownTest {
	return &ChaosNodeDownTest{}
}

func (t *ChaosNodeDownTest) Name() string {
	return "chaos-node-down"
}

func (t *ChaosNodeDownTest) Description() string {
	return "Stop a node while traffic is flowing, bring it back, verify it catches up"
}

func (t *ChaosNodeDownTest) Run(ctx context.Context, env *types.Environment) error {
	require := require.New(env.T())

	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddGateway(ctx))

	require.NoError(env.NewClient(100))
	env.Client(100).GenerateTraffic(ctx, client.TrafficOptions{
		BatchSize: 5,
		Duration:  60 * time.Second,
	})

	// Let some traffic flow first
	time.Sleep(5 * time.Second)

	// Kill node-300 while traffic flows
	require.NoError(env.Node(300).Stop(ctx))

	// Traffic keeps flowing to the 2 remaining nodes
	time.Sleep(10 * time.Second)

	// Restart node-300
	require.NoError(env.Node(300).Start(ctx))

	env.Client(100).Stop()

	// Verify healthy nodes received envelopes first to establish a baseline
	checkCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()
	require.NoError(env.Node(100).WaitForEnvelopes(checkCtx, 1))

	// Get envelope count from a healthy node as the catch-up target
	expectedCount, err := env.Node(100).GetEnvelopeCount(ctx)
	require.NoError(err, "failed to get envelope count from node 100")
	require.Positive(expectedCount, "node 100 should have envelopes after traffic")

	// Verify node-300 caught up to the same count as healthy nodes
	require.NoError(env.Node(300).WaitForEnvelopes(checkCtx, expectedCount))

	env.Logger.Info("chaos node down test completed")
	return nil
}

var _ types.Test = (*ChaosNodeDownTest)(nil)
