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

// ChaosAttestationFaultTest verifies payer report attestation behavior when
// a node crashes mid-cycle. With 3 nodes, quorum = (3/2)+1 = 2 attestations.
//
// The test runs in three phases:
//
//  1. Start 3 nodes, generate traffic, verify the pipeline works normally
//     (at least one report gets created).
//  2. Stop one node (leaving 2 alive). Since quorum is 2 and 2 nodes remain,
//     attestation should still succeed. Verify that reports progress past
//     AttestationPending even with a dead node.
//  3. Restart the stopped node and verify it rejoins the cluster. The cluster
//     should continue processing reports.
//
// This tests the resilience of the attestation quorum mechanism under
// real node failures (not just network partitions).
type ChaosAttestationFaultTest struct{}

func NewChaosAttestationFaultTest() *ChaosAttestationFaultTest {
	return &ChaosAttestationFaultTest{}
}

func (t *ChaosAttestationFaultTest) Name() string {
	return "chaos-attestation-fault"
}

func (t *ChaosAttestationFaultTest) Description() string {
	return "Kill a node mid-attestation cycle, verify quorum still works with 2 of 3"
}

func (t *ChaosAttestationFaultTest) Run(
	ctx context.Context,
	env *types.Environment,
) error {
	require := require.New(env.T())

	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddGateway(ctx))

	require.NoError(env.NewClient(100))

	// --- PHASE 1: Normal operation, wait for first report ---
	env.Logger.Info("phase 1: verifying normal pipeline operation")

	gen := env.Client(100).GenerateTraffic(ctx, client.TrafficOptions{
		BatchSize: 10,
		Duration:  75 * time.Minute,
	})

	createdCtx, createdCancel := context.WithTimeout(ctx, 75*time.Minute)
	defer createdCancel()

	require.NoError(env.Node(100).WaitForPayerReports(
		createdCtx,
		func(c *observe.PayerReportStatusCounts) bool {
			return c.Total > 0
		},
		"at least 1 payer report created (normal operation)",
	))

	// Record the baseline.
	baselineCounts, err := env.Node(100).GetPayerReportStatusCounts(ctx)
	require.NoError(err)
	env.Logger.Info("phase 1 baseline",
		zap.Int64("total", baselineCounts.Total),
		zap.Int64("attestation_approved", baselineCounts.AttestationApproved),
	)

	// --- PHASE 2: Kill node-300, verify attestation still works ---
	// With 3 nodes, quorum = 2. After killing node-300, 2 nodes remain
	// which is exactly quorum. Attestation should still succeed.
	env.Logger.Info("phase 2: stopping node-300 to test quorum resilience")
	require.NoError(env.Node(300).Stop(ctx))

	// Wait for at least one more report to get attested with only 2 nodes alive.
	// This proves quorum works with the minimum number of nodes.
	attestCtx, attestCancel := context.WithTimeout(ctx, 75*time.Minute)
	defer attestCancel()

	prevApproved := baselineCounts.AttestationApproved
	require.NoError(env.Node(100).WaitForPayerReports(
		attestCtx,
		func(c *observe.PayerReportStatusCounts) bool {
			return c.AttestationApproved > prevApproved
		},
		"attestation succeeded with 2-of-3 quorum (one node down)",
	))

	faultCounts, err := env.Node(100).GetPayerReportStatusCounts(ctx)
	require.NoError(err)
	env.Logger.Info("phase 2 result: attestation succeeded with dead node",
		zap.Int64("total", faultCounts.Total),
		zap.Int64("attestation_approved", faultCounts.AttestationApproved),
		zap.Int64("submission_settled", faultCounts.SubmissionSettled),
	)

	require.Greater(faultCounts.AttestationApproved, prevApproved,
		"should have more approved reports despite dead node",
	)

	// --- PHASE 3: Restart node-300, verify cluster recovers ---
	env.Logger.Info("phase 3: restarting node-300")
	require.NoError(env.Node(300).Start(ctx))

	// Give the restarted node time to rejoin and sync.
	time.Sleep(15 * time.Second)

	// Verify node-300 has synced envelopes from while it was down.
	recoveryCtx, recoveryCancel := context.WithTimeout(ctx, 90*time.Second)
	defer recoveryCancel()
	require.NoError(env.Node(300).WaitForEnvelopes(recoveryCtx, 1))

	gen.Stop()

	// Final state check.
	finalCounts, err := env.Node(100).GetPayerReportStatusCounts(ctx)
	require.NoError(err)
	env.Logger.Info("chaos attestation fault test completed",
		zap.Int64("total_reports", finalCounts.Total),
		zap.Int64("attestation_approved", finalCounts.AttestationApproved),
		zap.Int64("submission_settled", finalCounts.SubmissionSettled),
	)

	return nil
}

var _ types.Test = (*ChaosAttestationFaultTest)(nil)
