package tests

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/e2e/client"
	"github.com/xmtp/xmtpd/pkg/e2e/observe"
	"github.com/xmtp/xmtpd/pkg/e2e/types"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

// PayerLifecycleTest generates traffic and verifies the full payer report lifecycle:
// creation -> attestation -> submission -> settlement -> excess transfer -> claim -> withdraw.
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
	return "Generate traffic and verify the full payer report lifecycle including fee distribution"
}

func (t *PayerLifecycleTest) Run(ctx context.Context, env *types.Environment) error {
	require := require.New(env.T())

	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddGateway(ctx))

	require.NoError(env.NewClient(100))

	// Fund the payer before generating traffic so that settlement can
	// distribute actual tokens to node operators.
	payerKey := env.Client(100).PayerKey()
	payerPrivKey, err := utils.ParseEcdsaPrivateKey(payerKey)
	require.NoError(err, "failed to parse payer key")
	payerAddr := crypto.PubkeyToAddress(payerPrivKey.PublicKey)

	// Deposit a generous amount — 10^18 picodollars (~$1M).
	depositAmount := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	env.Logger.Info("funding payer",
		zap.String("payer", payerAddr.Hex()),
		zap.String("amount", depositAmount.String()))
	require.NoError(env.FundPayer(ctx, payerAddr, depositAmount))

	// Verify the deposit landed
	payerBalance, err := env.GetPayerBalance(ctx, payerAddr)
	require.NoError(err, "failed to get payer balance")
	env.Logger.Info("payer balance after deposit",
		zap.String("balance", payerBalance.String()))
	require.Positive(payerBalance.Sign(), "payer should have a positive balance")

	const trafficDuration = 75 * time.Minute
	const generatorTimeout = 65 * time.Minute
	const postGeneratorTimeout = 15 * time.Minute

	// Start background traffic for the duration of the test
	env.Client(100).GenerateTraffic(ctx, client.TrafficOptions{
		BatchSize: 10,
		Duration:  trafficDuration,
	})

	node100 := env.Node(100)

	// Phase 1: Wait for payer reports to be created
	// This is the long wait — up to 60 min for the generator's scheduled minute.
	env.Logger.Info("phase 1: waiting for payer reports to be created")
	createdCtx, createdCancel := context.WithTimeout(ctx, generatorTimeout)
	defer createdCancel()

	require.NoError(node100.WaitForPayerReports(
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

	require.NoError(node100.WaitForPayerReports(
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

	require.NoError(node100.WaitForPayerReports(
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

	require.NoError(node100.WaitForPayerReports(
		settledCtx,
		func(c *observe.PayerReportStatusCounts) bool {
			return c.SubmissionSettled > 0
		},
		"at least 1 payer report settled",
	))

	// Stop traffic before verification
	env.Client(100).Stop()

	// Phase 5: Verify payer report consistency across all nodes.
	// Every node should have the same settled reports.
	env.Logger.Info("phase 5: verifying payer reports across all nodes")

	allNodes := env.Nodes()
	for _, n := range allNodes {
		counts, err := n.GetPayerReportStatusCounts(ctx)
		require.NoError(err, "failed to get payer report counts from node %d", n.ID())

		env.Logger.Info("node payer report status",
			zap.Uint32("node_id", n.ID()),
			zap.Int64("total", counts.Total),
			zap.Int64("attestation_approved", counts.AttestationApproved),
			zap.Int64("submission_submitted", counts.SubmissionSubmitted),
			zap.Int64("submission_settled", counts.SubmissionSettled),
			zap.Int64("submission_rejected", counts.SubmissionRejected),
		)

		require.Positive(counts.Total,
			"node %d should have payer reports", n.ID())
		require.Positive(counts.SubmissionSettled,
			"node %d should have settled payer reports", n.ID())
		require.Equal(int64(0), counts.SubmissionRejected,
			"node %d should have no rejected payer reports", n.ID())
	}

	// Phase 6: Transfer excess funds from PayerRegistry to DistributionManager.
	// After settlement, the PayerRegistry holds tokens that are no longer owed to
	// payers (their balances were deducted). This creates "excess" that must be
	// transferred to the DistributionManager before nodes can withdraw.
	env.Logger.Info("phase 6: transferring excess to fee distributor")
	excess, err := env.GetPayerRegistryExcess(ctx)
	require.NoError(err, "failed to get payer registry excess")
	env.Logger.Info("payer registry excess", zap.String("excess", excess.String()))

	require.Positive(excess.Sign(),
		"payer registry should have excess after settlement (did the payer have tokens deposited?)")
	require.NoError(env.SendExcessToFeeDistributor(ctx))

	// Phase 7: Each node claims and withdraws their owed fees.
	env.Logger.Info("phase 7: claiming and withdrawing owed fees for each node")

	// Get settled reports from one node to build claim parameters.
	// All nodes should have the same settled reports.
	settledReports, err := node100.GetSettledPayerReports(ctx)
	require.NoError(err, "failed to get settled payer reports")
	require.NotEmpty(settledReports, "should have settled payer reports")

	env.Logger.Info("settled payer reports found",
		zap.Int("count", len(settledReports)))

	// Group settled reports by originator node ID.
	// Not all originators may have reports on-chain (e.g. a node didn't generate
	// reports during the test window), so we need to deduplicate the originator list.
	type reportKey struct {
		originatorNodeID uint32
		reportIndex      int64
	}
	uniqueReports := make(map[reportKey]struct{})
	for _, r := range settledReports {
		uniqueReports[reportKey{
			originatorNodeID: uint32(r.OriginatorNodeID),
			reportIndex:      int64(r.SubmittedReportIndex),
		}] = struct{}{}
	}

	originatorNodeIDs := make([]uint32, 0, len(uniqueReports))
	payerReportIndices := make([]*big.Int, 0, len(uniqueReports))
	for k := range uniqueReports {
		originatorNodeIDs = append(originatorNodeIDs, k.originatorNodeID)
		payerReportIndices = append(payerReportIndices, big.NewInt(k.reportIndex))
	}

	env.Logger.Info("unique settled reports for claim",
		zap.Int("count", len(originatorNodeIDs)),
		zap.Any("originator_node_ids", originatorNodeIDs))

	anyWithdrawn := false

	for _, n := range allNodes {
		nodeID := n.ID()
		ownerKey := n.SignerKey()

		// Attempt to claim fees for this node.
		// The claim may fail with NoReportsForOriginator if the on-chain contract
		// doesn't have reports for some originators. This is expected when not all
		// nodes generated reports during the test window.
		err = env.ClaimFromDistributionManager(
			ctx, ownerKey, nodeID, originatorNodeIDs, payerReportIndices,
		)
		if err != nil {
			env.Logger.Warn("claim failed for node (may be expected if no reports for originator)",
				zap.Uint32("node_id", nodeID),
				zap.Error(err))
			continue
		}

		// Check owed fees after claim to verify they were credited
		owedAfter, err := env.GetDistributionManagerOwedFees(ctx, nodeID)
		require.NoError(err, "failed to get owed fees after claim for node %d", nodeID)
		env.Logger.Info("owed fees after claim",
			zap.Uint32("node_id", nodeID),
			zap.String("owed", owedAfter.String()))

		// Withdraw owed fees (must be signed by node owner).
		if owedAfter.Sign() > 0 {
			err = env.WithdrawFromDistributionManager(ctx, ownerKey, nodeID)
			require.NoError(err, "withdraw failed for node %d", nodeID)

			env.Logger.Info("node claimed and withdrew fees",
				zap.Uint32("node_id", nodeID),
				zap.String("amount", owedAfter.String()))
			anyWithdrawn = true
		} else {
			env.Logger.Info("node has no owed fees to withdraw",
				zap.Uint32("node_id", nodeID))
		}
	}

	require.True(anyWithdrawn, "at least one node should have withdrawn fees")

	// Phase 8: Verify that fee tokens arrived in node operator wallets.
	env.Logger.Info("phase 8: verifying fee token balances in node operator wallets")

	for _, n := range allNodes {
		ownerKey := n.SignerKey()
		ownerPrivKey, keyErr := utils.ParseEcdsaPrivateKey(ownerKey)
		require.NoError(keyErr, "failed to parse owner key for node %d", n.ID())
		ownerAddr := crypto.PubkeyToAddress(ownerPrivKey.PublicKey)

		balance, balErr := env.GetFeeTokenBalance(ctx, ownerAddr)
		require.NoError(balErr, "failed to get fee token balance for node %d", n.ID())
		env.Logger.Info("node operator fee token balance",
			zap.Uint32("node_id", n.ID()),
			zap.String("address", ownerAddr.Hex()),
			zap.String("balance", balance.String()))

		require.Positive(balance.Sign(),
			"node %d operator should have received fee tokens", n.ID())
	}

	// Log final status
	finalCounts, err := node100.GetPayerReportStatusCounts(ctx)
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
