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

// PayerLifecycleTest generates traffic and verifies the full payer report lifecycle:
// creation -> attestation -> submission -> settlement.
//
// Worker scheduling:
//
//	Generator (workerID=1), Submitter (workerID=2), Settlement (workerID=3)
//	fire at minute offsets 0, +5, +10 within a 60-minute cycle based on
//	Knuth hash of the node ID. Attestation polls every 10s (env var).
//
// The worst-case wait is for the generator's first fire (~60 min).
// Total test runtime: ~90 minutes.
type PayerLifecycleTest struct{}

func NewPayerLifecycleTest() *PayerLifecycleTest {
	return &PayerLifecycleTest{}
}

func (t *PayerLifecycleTest) Name() string {
	return "payer-lifecycle"
}

func (t *PayerLifecycleTest) Description() string {
	return "Generate traffic and verify payer reports are created, attested, and submitted"
}

func (t *PayerLifecycleTest) Run(ctx context.Context, env *types.Environment) error {
	require := require.New(env.T())

	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddGateway(ctx))

	require.NoError(env.NewClient(100))

	const trafficDuration = 75 * time.Minute
	const generatorTimeout = 65 * time.Minute
	const postGeneratorTimeout = 15 * time.Minute

	// Start background traffic for the duration of the test
	env.Client(100).GenerateTraffic(ctx, client.TrafficOptions{
		BatchSize: 10,
		Duration:  trafficDuration,
	})

	observeNode := env.Node(100)

	// Phase 1: Wait for payer reports to be created
	// This is the long wait — up to 60 min for the generator's scheduled minute.
	env.Logger.Info("phase 1: waiting for payer reports to be created")
	createdCtx, createdCancel := context.WithTimeout(ctx, generatorTimeout)
	defer createdCancel()

	require.NoError(observeNode.WaitForPayerReports(
		createdCtx,
		func(c *observe.PayerReportStatusCounts) bool {
			return c.Total > 0
		},
		"at least 1 payer report created",
	))

	// Phase 2: Wait for attestation approval
	// Attestation worker polls every 10s, so this should be fast.
	env.Logger.Info("phase 2: waiting for payer reports to be attested")
	attestedCtx, attestedCancel := context.WithTimeout(ctx, postGeneratorTimeout)
	defer attestedCancel()

	require.NoError(observeNode.WaitForPayerReports(
		attestedCtx,
		func(c *observe.PayerReportStatusCounts) bool {
			return c.AttestationApproved > 0
		},
		"at least 1 payer report attested",
	))

	// Phase 3: Wait for submission (or settlement, since settled implies submitted).
	// The submitter fires 5 min after the generator. With anvil's instant mining,
	// reports can transition from submitted -> settled before the observer polls.
	env.Logger.Info("phase 3: waiting for payer reports to be submitted")
	submittedCtx, submittedCancel := context.WithTimeout(ctx, postGeneratorTimeout)
	defer submittedCancel()

	require.NoError(observeNode.WaitForPayerReports(
		submittedCtx,
		func(c *observe.PayerReportStatusCounts) bool {
			return c.SubmissionSubmitted > 0 || c.SubmissionSettled > 0
		},
		"at least 1 payer report submitted or settled",
	))

	// Phase 4: Wait for settlement
	// The node that submitted the report immediately tries to settle it.
	// With anvil's instant mining this should complete quickly.
	env.Logger.Info("phase 4: waiting for payer reports to be settled")
	settledCtx, settledCancel := context.WithTimeout(ctx, postGeneratorTimeout)
	defer settledCancel()

	require.NoError(observeNode.WaitForPayerReports(
		settledCtx,
		func(c *observe.PayerReportStatusCounts) bool {
			return c.SubmissionSettled > 0
		},
		"at least 1 payer report settled",
	))

	env.Client(100).Stop()

	// Log final status
	finalCounts, err := observeNode.GetPayerReportStatusCounts(ctx)
	if err != nil {
		env.Logger.Warn("failed to get final payer report counts", zap.Error(err))
	} else {
		env.Logger.Info("payer lifecycle test completed",
			zap.Int64("total_reports", finalCounts.Total),
			zap.Int64("attestation_approved", finalCounts.AttestationApproved),
			zap.Int64("submission_submitted", finalCounts.SubmissionSubmitted),
			zap.Int64("submission_settled", finalCounts.SubmissionSettled),
		)
	}

	return nil
}

var _ types.Test = (*PayerLifecycleTest)(nil)
