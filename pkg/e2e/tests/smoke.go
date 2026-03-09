// Package tests contains the E2E test implementations for xmtpd.
package tests

import (
	"context"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/e2e/types"
)

// SmokeTest verifies basic cluster functionality: starts nodes and a gateway,
// publishes envelopes, and checks they replicate to all nodes.
type SmokeTest struct{}

func NewSmokeTest() *SmokeTest {
	return &SmokeTest{}
}

func (t *SmokeTest) Name() string {
	return "smoke"
}

func (t *SmokeTest) Description() string {
	return "Start nodes and gateways, publish envelopes, verify they are replicated"
}

func (t *SmokeTest) Run(ctx context.Context, env *types.Environment) error {
	require := require.New(env.T())

	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddGateway(ctx))

	require.NoError(env.NewClient(100))
	require.NoError(env.Client(100).PublishEnvelopes(ctx, 10))

	checkCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	for _, n := range env.Nodes() {
		require.NoError(n.WaitForEnvelopes(checkCtx, 1))
	}

	env.Logger.Info("smoke test completed successfully")
	return nil
}

var _ types.Test = (*SmokeTest)(nil)
