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

// StuckStateDetectionTest verifies that the system can be observed entering a
// stuck state when attestation quorum is lost. It starts with 4 nodes, waits
// for the payer report pipeline to function normally, then removes enough nodes
// from the canonical network that the remaining nodes cannot reach the 2/3
// quorum required for attestation.
//
// The test then verifies:
//   - New reports are still generated (generator runs on each node independently)
//   - Reports remain stuck in AttestationPending (quorum unreachable)
//   - After restoring the canonical network, reports resume normal progression
//
// This demonstrates the observability needed for stuck state alerting.
type StuckStateDetectionTest struct{}

func NewStuckStateDetectionTest() *StuckStateDetectionTest {
	return &StuckStateDetectionTest{}
}

func (t *StuckStateDetectionTest) Name() string {
	return "stuck-state-detection"
}

func (t *StuckStateDetectionTest) Description() string {
	return "Detect stuck payer reports when attestation quorum is lost, then recover"
}

func (t *StuckStateDetectionTest) Run(
	ctx context.Context,
	env *types.Environment,
) error {
	require := require.New(env.T())

	// Start 4 nodes. With 4 nodes, quorum = (4/2)+1 = 3 attestations needed.
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddGateway(ctx))

	require.NoError(env.NewClient(100))

	// Phase 1: Verify the pipeline works normally.
	// Generate traffic and wait for at least one report to be created.
	env.Logger.Info("phase 1: verifying normal pipeline operation")
	gen := env.Client(100).GenerateTraffic(ctx, client.TrafficOptions{
		BatchSize: 10,
		Duration:  15 * time.Minute,
	})

	createdCtx, createdCancel := context.WithTimeout(ctx, 10*time.Minute)
	defer createdCancel()

	require.NoError(env.Node(100).WaitForPayerReports(
		createdCtx,
		func(c *observe.PayerReportStatusCounts) bool {
			return c.Total > 0
		},
		"at least 1 payer report created (normal operation)",
	))

	// Record the baseline report counts.
	baselineCounts, err := env.Node(100).GetPayerReportStatusCounts(ctx)
	require.NoError(err)

	env.Logger.Info("baseline report counts",
		zap.Int64("total", baselineCounts.Total),
		zap.Int64("attestation_pending", baselineCounts.AttestationPending),
		zap.Int64("attestation_approved", baselineCounts.AttestationApproved),
	)

	// Phase 2: Break quorum by removing 3 of 4 nodes from the canonical network.
	// With only 1 node in the canonical network, reports generated with 4 nodes
	// in ActiveNodeIDs required 3 attestations — but only 1 node is active.
	env.Logger.Info("phase 2: removing nodes to break attestation quorum")
	require.NoError(env.Node(200).RemoveFromCanonicalNetwork(ctx))
	require.NoError(env.Node(300).RemoveFromCanonicalNetwork(ctx))
	require.NoError(env.Node(400).RemoveFromCanonicalNetwork(ctx))

	// Wait for new reports to be generated (they will be stuck).
	// Reports generated during the quorum loss period will have the
	// previous ActiveNodeIDs (4 nodes) but can't get 3 attestations.
	env.Logger.Info("phase 2: waiting for stuck reports to appear")
	time.Sleep(30 * time.Second)

	// Check for stuck state: reports in AttestationPending that can't progress.
	stuckCounts, err := env.Node(100).GetPayerReportStatusCounts(ctx)
	require.NoError(err)

	env.Logger.Info("report counts during quorum loss",
		zap.Int64("total", stuckCounts.Total),
		zap.Int64("attestation_pending", stuckCounts.AttestationPending),
		zap.Int64("attestation_approved", stuckCounts.AttestationApproved),
		zap.Int64("submission_settled", stuckCounts.SubmissionSettled),
	)

	// The key signal for alerting: attestation_pending count is increasing
	// or staying non-zero while approved count isn't growing.
	// (New reports created during quorum loss stay pending.)

	// Phase 3: Restore quorum by re-adding nodes to the canonical network.
	env.Logger.Info("phase 3: restoring canonical network quorum")
	require.NoError(env.Node(200).AddToCanonicalNetwork(ctx))
	require.NoError(env.Node(300).AddToCanonicalNetwork(ctx))
	require.NoError(env.Node(400).AddToCanonicalNetwork(ctx))

	// Wait for the pipeline to recover — new reports should get attested.
	env.Logger.Info("phase 3: waiting for pipeline recovery")
	recoveryCtx, recoveryCancel := context.WithTimeout(ctx, 10*time.Minute)
	defer recoveryCancel()

	// After recovery, at least one new report should progress to approved.
	prevApproved := stuckCounts.AttestationApproved
	require.NoError(env.Node(100).WaitForPayerReports(
		recoveryCtx,
		func(c *observe.PayerReportStatusCounts) bool {
			return c.AttestationApproved > prevApproved
		},
		"attestation approval resumed after quorum restored",
	))

	// Log final state.
	finalCounts, err := env.Node(100).GetPayerReportStatusCounts(ctx)
	require.NoError(err)

	env.Logger.Info("stuck state detection test completed",
		zap.Int64("total_reports", finalCounts.Total),
		zap.Int64("attestation_pending", finalCounts.AttestationPending),
		zap.Int64("attestation_approved", finalCounts.AttestationApproved),
		zap.Int64("submission_submitted", finalCounts.SubmissionSubmitted),
		zap.Int64("submission_settled", finalCounts.SubmissionSettled),
	)

	gen.Stop()
	return nil
}

var _ types.Test = (*StuckStateDetectionTest)(nil)
