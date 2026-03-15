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

// SettlementVerificationTest verifies end-to-end settlement by cross-referencing
// database state against on-chain contract state. After traffic generates payer
// reports that settle, the test reads the on-chain report via the PayerReportManager
// contract and asserts that sequence IDs, settlement status, and node IDs match
// the database records.
type SettlementVerificationTest struct{}

func NewSettlementVerificationTest() *SettlementVerificationTest {
	return &SettlementVerificationTest{}
}

func (t *SettlementVerificationTest) Name() string {
	return "settlement-verification"
}

func (t *SettlementVerificationTest) Description() string {
	return "Verify on-chain settlement matches database state for payer reports"
}

func (t *SettlementVerificationTest) Run(
	ctx context.Context,
	env *types.Environment,
) error {
	require := require.New(env.T())

	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddGateway(ctx))

	require.NoError(env.NewClient(100))

	// Generate traffic long enough for the payer report cycle to complete.
	// The generator worker fires at minute 0 of a 60-min cycle (worst case ~60 min),
	// followed by attestation (~10s), submission (up to 60 min), and settlement (up to 60 min).
	// Submitter and settlement workers use findNextRunTime with a 60-minute repeat interval.
	const trafficDuration = 75 * time.Minute
	const generatorTimeout = 75 * time.Minute
	const postGeneratorTimeout = 65 * time.Minute

	gen := env.Client(100).GenerateTraffic(ctx, client.TrafficOptions{
		BatchSize: 10,
		Duration:  trafficDuration,
	})

	observeNode := env.Node(100)

	// Phase 1: Wait for at least one report to be settled in the DB.
	env.Logger.Info("waiting for payer reports to settle in database")
	settledCtx, settledCancel := context.WithTimeout(ctx, generatorTimeout+postGeneratorTimeout)
	defer settledCancel()

	require.NoError(observeNode.WaitForPayerReports(
		settledCtx,
		func(c *observe.PayerReportStatusCounts) bool {
			return c.SubmissionSettled > 0
		},
		"at least 1 payer report settled",
	))

	// Phase 2: Get settled reports from the DB (includes the on-chain index).
	env.Logger.Info("fetching settled reports from database")
	settledReports, err := observeNode.GetSettledPayerReports(ctx)
	require.NoError(err)
	require.NotEmpty(settledReports, "expected at least one settled report")

	// Phase 3: Cross-reference each settled report against on-chain state.
	for _, dbReport := range settledReports {
		env.Logger.Info("verifying on-chain state for settled report",
			zap.Int32("originator_node_id", dbReport.OriginatorNodeID),
			zap.Int32("report_index", dbReport.SubmittedReportIndex),
		)

		onChain, err := env.Contracts.GetPayerReport(
			ctx,
			uint32(dbReport.OriginatorNodeID),
			uint64(dbReport.SubmittedReportIndex),
		)
		require.NoError(err, "failed to read on-chain report")

		// Verify settlement status.
		require.True(
			onChain.IsSettled,
			"on-chain report should be settled",
		)

		env.Logger.Info("on-chain state verified",
			zap.Uint64("on_chain_start_seq", onChain.StartSequenceID),
			zap.Uint64("on_chain_end_seq", onChain.EndSequenceID),
			zap.Bool("on_chain_settled", onChain.IsSettled),
			zap.Uint32s("on_chain_node_ids", onChain.NodeIDs),
		)
	}

	gen.Stop()
	env.Logger.Info("settlement verification test completed successfully")
	return nil
}

var _ types.Test = (*SettlementVerificationTest)(nil)
