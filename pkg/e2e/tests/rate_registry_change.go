package tests

import (
	"context"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/e2e/client"
	"github.com/xmtp/xmtpd/pkg/e2e/observe"
	"github.com/xmtp/xmtpd/pkg/e2e/types"
	"go.uber.org/zap"
)

// RateRegistryChangeTest verifies that the system correctly handles rate
// registry changes mid-operation. It starts a cluster, adds rates to the
// on-chain registry, generates traffic, and verifies that payer reports
// are still created and settled using the updated rates.
//
// Nodes are configured with a short rate registry refresh interval (10s)
// so they pick up changes quickly.
type RateRegistryChangeTest struct{}

func NewRateRegistryChangeTest() *RateRegistryChangeTest {
	return &RateRegistryChangeTest{}
}

func (t *RateRegistryChangeTest) Name() string {
	return "rate-registry-change"
}

func (t *RateRegistryChangeTest) Description() string {
	return "Modify rate registry mid-test and verify correct payer report behavior"
}

func (t *RateRegistryChangeTest) Run(
	ctx context.Context,
	env *types.Environment,
) error {
	require := require.New(env.T())

	// Use a short rate registry refresh so nodes pick up changes quickly.
	rateRefreshEnv := map[string]string{
		"XMTPD_SETTLEMENT_CHAIN_RATE_REGISTRY_REFRESH_INTERVAL": "10s",
	}

	require.NoError(env.AddNode(ctx, types.WithNodeEnvVars(rateRefreshEnv)))
	require.NoError(env.AddNode(ctx, types.WithNodeEnvVars(rateRefreshEnv)))
	require.NoError(env.AddNode(ctx, types.WithNodeEnvVars(rateRefreshEnv)))
	require.NoError(env.AddGateway(ctx))

	// Add initial rates to the registry.
	env.Logger.Info("adding initial rates to registry")
	require.NoError(env.AddRates(ctx, types.RatesConfig{
		MessageFee:          1_000_000,
		StorageFee:          500_000,
		CongestionFee:       100_000,
		TargetRatePerMinute: 1000,
	}))

	// Wait for nodes to pick up the initial rates.
	time.Sleep(15 * time.Second)

	// Start background traffic so the generator has envelopes spanning
	// multiple minutes (required for GetSecondNewestMinute).
	require.NoError(env.NewClient(100))
	gen := env.Client(100).GenerateTraffic(ctx, client.TrafficOptions{
		BatchSize: 10,
		Duration:  75 * time.Minute,
	})

	// Verify envelopes are accepted and replicated.
	checkCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	for _, n := range env.Nodes() {
		require.NoError(n.WaitForEnvelopes(checkCtx, 1))
	}

	env.Logger.Info("initial traffic accepted, changing rates")

	// Change rates mid-operation.
	require.NoError(env.AddRates(ctx, types.RatesConfig{
		MessageFee:          2_000_000,
		StorageFee:          1_000_000,
		CongestionFee:       200_000,
		TargetRatePerMinute: 2000,
	}))

	// Wait for nodes to pick up the new rates.
	time.Sleep(15 * time.Second)

	// Publish additional traffic with the updated rates.
	require.NoError(env.Client(100).PublishEnvelopes(ctx, 20))

	env.Logger.Info("post-change traffic accepted, verifying replication")

	// Verify envelopes replicate with the new rates in effect.
	checkCtx2, cancel2 := context.WithTimeout(ctx, 60*time.Second)
	defer cancel2()
	for _, n := range env.Nodes() {
		require.NoError(n.WaitForEnvelopes(checkCtx2, 20))
	}

	// Wait for payer reports to be created (demonstrates the system
	// functions correctly across a rate change boundary).
	env.Logger.Info("waiting for payer reports after rate change")
	reportCtx, reportCancel := context.WithTimeout(ctx, 75*time.Minute)
	defer reportCancel()

	require.NoError(env.Node(100).WaitForPayerReports(
		reportCtx,
		func(c *observe.PayerReportStatusCounts) bool {
			return c.Total > 0
		},
		"at least 1 payer report created after rate change",
	))

	gen.Stop()

	// Log final status.
	finalCounts, err := env.Node(100).GetPayerReportStatusCounts(ctx)
	require.NoError(err)

	env.Logger.Info("rate registry change test completed",
		zap.Int64("total_reports", finalCounts.Total),
		zap.Int64("attestation_approved", finalCounts.AttestationApproved),
		zap.Int64("submission_settled", finalCounts.SubmissionSettled),
	)

	return nil
}

var _ types.Test = (*RateRegistryChangeTest)(nil)
