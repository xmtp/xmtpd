package tests

import (
	"context"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/e2e/types"
)

// GatewayScaleTest verifies that gateways can be dynamically added and removed
// while traffic is being generated without errors.
type GatewayScaleTest struct{}

func NewGatewayScaleTest() *GatewayScaleTest {
	return &GatewayScaleTest{}
}

func (t *GatewayScaleTest) Name() string {
	return "gateway-scale"
}

func (t *GatewayScaleTest) Description() string {
	return "Add and remove gateways dynamically while generating traffic"
}

func (t *GatewayScaleTest) Run(ctx context.Context, env *types.Environment) error {
	require := require.New(env.T())

	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddGateway(ctx))

	require.NoError(env.NewClient(100))

	// Generate some initial traffic
	require.NoError(env.Client(100).PublishEnvelopes(ctx, 10))

	// Scale up: add 2 extra gateways
	require.NoError(env.AddGateway(ctx))
	require.NoError(env.AddGateway(ctx))

	// Generate traffic with more gateways
	require.NoError(env.Client(100).PublishEnvelopes(ctx, 10))

	// Scale down: stop the extra gateways
	require.NoError(env.Gateway(1).Stop(ctx))
	require.NoError(env.Gateway(2).Stop(ctx))

	// Generate traffic again after scale-down
	require.NoError(env.Client(100).PublishEnvelopes(ctx, 10))

	env.Logger.Info("gateway scale test completed")
	return nil
}

var _ types.Test = (*GatewayScaleTest)(nil)
