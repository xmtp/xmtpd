package tests

import (
	"context"
	"math/big"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/e2e/client"
	"github.com/xmtp/xmtpd/pkg/e2e/observe"
	"github.com/xmtp/xmtpd/pkg/e2e/types"
	"go.uber.org/zap"
)

// PayerLifecycleTest generates traffic and verifies the full payer report lifecycle:
// creation -> attestation -> submission -> settlement -> excess transfer -> claim -> withdraw -> payer withdrawal.
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

	const (
		trafficDuration      = 75 * time.Minute
		generatorTimeout     = 65 * time.Minute
		postGeneratorTimeout = 15 * time.Minute
	)

	// Create nodes 100, 200, 300.
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))

	// Create a gateway.
	require.NoError(env.AddGateway(ctx))

	// Create a client for node 100.
	require.NoError(env.NewClient(100))

	// Record initial fee token balances for all nodes before any traffic.
	allNodes := env.Nodes()
	initialBalances := make(map[uint32]*big.Int, len(allNodes))
	for _, n := range allNodes {
		bal, balErr := n.GetFeeTokenBalance(ctx)
		require.NoError(balErr, "failed to get initial fee token balance for node %d", n.ID())
		initialBalances[n.ID()] = bal
		env.Logger.Info("initial fee token balance",
			zap.Uint32("node_id", n.ID()),
			zap.String("balance", bal.String()))
	}

	// Fund the payer before generating traffic so that settlement can
	// distribute actual tokens to node operators.
	var (
		payer         = env.Client(100)
		node100       = env.Node(100)
		depositAmount = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	)

	env.Logger.Info("funding payer",
		zap.String("payer", payer.Address().Hex()),
		zap.String("amount", depositAmount.String()))

	require.NoError(payer.Deposit(ctx, depositAmount))

	// Verify the deposit landed.
	payerInitialBalance, err := payer.GetPayerBalance(ctx)
	require.NoError(err, "failed to get payer balance")

	env.Logger.Info("payer balance after deposit",
		zap.String("balance", payerInitialBalance.String()))

	require.Positive(payerInitialBalance.Sign(), "payer should have a positive balance")

	// Start background traffic for the duration of the test.
	// Capture the TrafficGenerator to check for errors after stopping.
	traffic := payer.GenerateTraffic(ctx, client.TrafficOptions{
		BatchSize: 10,
		Duration:  trafficDuration,
	})
	defer payer.Stop()

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

	// Stop traffic before verification and check for generation errors.
	payer.Stop()
	require.NoError(traffic.Err(), "background traffic generation failed")

	// Phase 5: Verify payer report consistency across all nodes.
	// Every node should have the same settled reports.
	env.Logger.Info("phase 5: verifying payer reports across all nodes")

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

	// Phase 5b: Verify payer balance decreased after settlement.
	// Settlement deducts from the payer's PayerRegistry balance, so it must be
	// strictly less than the original deposit.
	payerBalanceAfter, err := payer.GetPayerBalance(ctx)
	require.NoError(err, "failed to get payer balance after settlement")

	env.Logger.Info("payer balance after settlement",
		zap.String("before", depositAmount.String()),
		zap.String("after", payerBalanceAfter.String()))

	require.Negative(payerBalanceAfter.Cmp(depositAmount),
		"payer balance should have decreased after settlement (before=%s, after=%s)",
		depositAmount.String(), payerBalanceAfter.String())

	// Phase 6: Transfer excess funds from PayerRegistry to DistributionManager.
	// After settlement, the PayerRegistry holds tokens that are no longer owed to
	// payers (their balances were deducted). This creates "excess" that must be
	// transferred to the DistributionManager before nodes can withdraw.
	env.Logger.Info("phase 6: transferring excess to fee distributor")

	excess, err := env.GetPayerRegistryExcess(ctx)
	require.NoError(err, "failed to get payer registry excess")
	require.Positive(excess.Sign(),
		"payer registry should have excess after settlement (did the payer have tokens deposited?)")

	env.Logger.Info(
		"payer registry excess to be distributed",
		zap.String("excess", excess.String()),
	)

	err = env.SendExcessToFeeDistributor(ctx)
	require.NoError(err, "failed to send excess to fee distributor")

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

	var (
		originatorNodeIDs  = make([]uint32, 0, len(uniqueReports))
		payerReportIndices = make([]*big.Int, 0, len(uniqueReports))
	)

	for k := range uniqueReports {
		originatorNodeIDs = append(originatorNodeIDs, k.originatorNodeID)
		payerReportIndices = append(payerReportIndices, big.NewInt(k.reportIndex))
	}

	env.Logger.Info("unique settled reports for claim",
		zap.Int("count", len(originatorNodeIDs)),
		zap.Any("originator_node_ids", originatorNodeIDs))

	anyWithdrawn := false

	for _, n := range allNodes {
		var (
			nodeID   = n.ID()
			ownerKey = n.SignerKey()
		)

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
	// The sum of all earned fees across nodes must equal the excess that was
	// transferred from the PayerRegistry to the DistributionManager.
	env.Logger.Info("phase 8: verifying fee token balances in node operator wallets")

	totalEarned := new(big.Int)
	for _, n := range allNodes {
		finalBalance, balErr := n.GetFeeTokenBalance(ctx)
		require.NoError(balErr, "failed to get fee token balance for node %d", n.ID())

		initial := initialBalances[n.ID()]
		earned := new(big.Int).Sub(finalBalance, initial)
		totalEarned.Add(totalEarned, earned)

		env.Logger.Info("node operator fee token balance",
			zap.Uint32("node_id", n.ID()),
			zap.String("address", n.Address().Hex()),
			zap.String("initial", initial.String()),
			zap.String("final", finalBalance.String()),
			zap.String("earned", earned.String()))

		require.Positive(earned.Sign(),
			"node %d operator should have earned fee tokens (initial=%s, final=%s)",
			n.ID(), initial.String(), finalBalance.String())
	}

	env.Logger.Info("total fees distributed",
		zap.String("total_earned", totalEarned.String()),
		zap.String("excess", excess.String()))

	// Phase 9: Payer requests withdrawal of remaining balance.
	// After settlement, the payer should still have leftover funds in the
	// PayerRegistry. Verify the withdrawal request flow works.
	payerFinalBalance, err := payer.GetPayerBalance(ctx)
	require.NoError(err, "failed to get remaining payer balance")
	env.Logger.Info("payer remaining balance before withdrawal",
		zap.String("balance", payerFinalBalance.String()))

	require.Equal(
		payerFinalBalance.Add(payerFinalBalance, excess),
		payerInitialBalance,
		"remaining payer balance should equal the sum of the total earned fees and the excess transferred",
	)

	if payerFinalBalance.Sign() > 0 {
		require.NoError(payer.RequestWithdrawal(ctx, payerFinalBalance),
			"payer withdrawal request should succeed")
		env.Logger.Info("payer withdrawal requested",
			zap.String("amount", payerFinalBalance.String()))
	}

	// Log final status.
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
