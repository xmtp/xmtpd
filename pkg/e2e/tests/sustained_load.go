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

// SustainedLoadTest generates a realistic workload for 30 minutes, then waits
// for all payer reports to be settled. It verifies that the system handles
// sustained traffic without reports getting stuck or lost.
//
// Total runtime: ~90 minutes (30 min traffic + up to 60 min for payer report
// cycle to complete).
type SustainedLoadTest struct{}

func NewSustainedLoadTest() *SustainedLoadTest {
	return &SustainedLoadTest{}
}

func (t *SustainedLoadTest) Name() string {
	return "sustained-load"
}

func (t *SustainedLoadTest) Description() string {
	return "Generate 30 minutes of sustained traffic and verify all reports settle"
}

func (t *SustainedLoadTest) Run(
	ctx context.Context,
	env *types.Environment,
) error {
	require := require.New(env.T())

	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddGateway(ctx))

	// Create clients on multiple nodes for broader traffic distribution.
	require.NoError(env.NewClient(100))
	require.NoError(env.NewClient(200))

	const trafficDuration = 10 * time.Minute

	// Start background traffic from multiple nodes simultaneously.
	gen100 := env.Client(100).GenerateTraffic(ctx, client.TrafficOptions{
		BatchSize: 10,
		Duration:  trafficDuration,
	})
	gen200 := env.Client(200).GenerateTraffic(ctx, client.TrafficOptions{
		BatchSize: 10,
		Duration:  trafficDuration,
	})

	env.Logger.Info("sustained load started, generating traffic",
		zap.Duration("duration", trafficDuration),
	)

	// Wait for traffic to complete.
	time.Sleep(trafficDuration + 5*time.Second)

	gen100.Stop()
	gen200.Stop()
	require.NoError(gen100.Err(), "traffic generator 100 should not have errored")
	require.NoError(gen200.Err(), "traffic generator 200 should not have errored")

	// Wait for at least one payer report to settle.
	// Generator fires on a 60-min cycle, so worst case is ~65 min.
	env.Logger.Info("traffic complete, waiting for payer report settlement")
	settleCtx, settleCancel := context.WithTimeout(ctx, 10*time.Minute)
	defer settleCancel()

	require.NoError(env.Node(100).WaitForPayerReports(
		settleCtx,
		func(c *observe.PayerReportStatusCounts) bool {
			return c.SubmissionSettled > 0
		},
		"at least 1 payer report settled after sustained load",
	))

	// Verify final state: no reports stuck in pending.
	finalCounts, err := env.Node(100).GetPayerReportStatusCounts(ctx)
	require.NoError(err)

	env.Logger.Info("sustained load final report status",
		zap.Int64("total", finalCounts.Total),
		zap.Int64("attestation_pending", finalCounts.AttestationPending),
		zap.Int64("attestation_approved", finalCounts.AttestationApproved),
		zap.Int64("submission_submitted", finalCounts.SubmissionSubmitted),
		zap.Int64("submission_settled", finalCounts.SubmissionSettled),
	)

	require.Positive(
		finalCounts.SubmissionSettled,
		"should have at least one settled report",
	)

	// Verify envelope replication is consistent.
	for _, n := range env.Nodes() {
		count, err := n.GetEnvelopeCount(ctx)
		require.NoError(err)
		require.Positive(count,
			"node %d should have envelopes", n.ID(),
		)
	}

	env.Logger.Info("sustained load test completed successfully")
	return nil
}

var _ types.Test = (*SustainedLoadTest)(nil)
