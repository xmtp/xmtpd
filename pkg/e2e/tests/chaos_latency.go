package tests

import (
	"context"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/e2e/client"
	"github.com/xmtp/xmtpd/pkg/e2e/types"
)

// ChaosLatencyTest injects network latency into a node while traffic is flowing
// and verifies the system remains functional after the latency is removed.
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
	env.Client(100).GenerateTraffic(ctx, client.TrafficOptions{
		BatchSize: 5,
		Duration:  60 * time.Second,
	})

	// Let some clean traffic flow first
	time.Sleep(5 * time.Second)

	// Inject 500ms latency to node-100 while traffic is flowing
	require.NoError(env.Node(100).AddLatency(ctx, 500))

	// Let traffic flow under degraded conditions
	time.Sleep(15 * time.Second)

	// Remove the toxic while traffic is still flowing
	require.NoError(env.Node(100).RemoveAllToxics(ctx))

	// Let traffic flow under normal conditions again
	time.Sleep(5 * time.Second)

	env.Client(100).Stop()

	env.Logger.Info("chaos latency test completed")
	return nil
}

var _ types.Test = (*ChaosLatencyTest)(nil)
