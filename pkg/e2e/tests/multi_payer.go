package tests

import (
	"context"
	"strings"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/e2e/types"
	"go.uber.org/zap"
)

// MultiPayerTest verifies that traffic from distinct payer addresses is correctly
// attributed in the database. It creates multiple clients with different payer keys
// targeting the same node, generates traffic from each, and asserts that
// per-payer usage records are created with the correct addresses.
type MultiPayerTest struct{}

func NewMultiPayerTest() *MultiPayerTest {
	return &MultiPayerTest{}
}

func (t *MultiPayerTest) Name() string {
	return "multi-payer"
}

func (t *MultiPayerTest) Description() string {
	return "Verify per-payer attribution when multiple payer addresses send traffic"
}

func (t *MultiPayerTest) Run(ctx context.Context, env *types.Environment) error {
	require := require.New(env.T())

	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddGateway(ctx))

	// Create payer A with an explicitly allocated key.
	payerKeyA, err := env.Keys.NextClientKey(ctx)
	require.NoError(err)
	require.NoError(env.NewClient(100,
		types.WithClientName("payer-a"),
		types.WithPayerKey(payerKeyA),
	))
	payerAddrA := env.ClientByName("payer-a").PayerAddress()

	// Create payer B with the next key (guaranteed to be different).
	payerKeyB, err := env.Keys.NextClientKey(ctx)
	require.NoError(err)
	require.NoError(env.NewClient(100,
		types.WithClientName("payer-b"),
		types.WithPayerKey(payerKeyB),
	))
	payerAddrB := env.ClientByName("payer-b").PayerAddress()

	// Sanity: addresses must differ.
	require.NotEqual(payerAddrA, payerAddrB, "payer addresses must be distinct")

	env.Logger.Info("publishing from two distinct payers",
		zap.String("payer_a", payerAddrA),
		zap.String("payer_b", payerAddrB),
	)

	// Publish envelopes from both payers.
	const envelopesPerPayer = 10
	require.NoError(env.ClientByName("payer-a").PublishEnvelopes(ctx, envelopesPerPayer))
	require.NoError(env.ClientByName("payer-b").PublishEnvelopes(ctx, envelopesPerPayer))

	// Wait for envelopes to replicate and usage to be recorded.
	checkCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	for _, n := range env.Nodes() {
		require.NoError(n.WaitForEnvelopes(checkCtx, envelopesPerPayer*2))
	}

	// Verify per-payer attribution on node 100.
	usage, err := env.Node(100).GetUnsettledUsage(ctx)
	require.NoError(err)

	usageByAddr := make(map[string]int64)
	for _, u := range usage {
		usageByAddr[strings.ToLower(u.PayerAddress)] = u.MessageCount
	}

	env.Logger.Info("per-payer usage",
		zap.Any("usage_by_address", usageByAddr),
	)

	require.Contains(usageByAddr, strings.ToLower(payerAddrA), "payer A should have usage records")
	require.Contains(usageByAddr, strings.ToLower(payerAddrB), "payer B should have usage records")
	require.GreaterOrEqual(
		usageByAddr[strings.ToLower(payerAddrA)], int64(1),
		"payer A should have non-zero message count",
	)
	require.GreaterOrEqual(
		usageByAddr[strings.ToLower(payerAddrB)], int64(1),
		"payer B should have non-zero message count",
	)

	env.Logger.Info("multi-payer test completed successfully")
	return nil
}

var _ types.Test = (*MultiPayerTest)(nil)
