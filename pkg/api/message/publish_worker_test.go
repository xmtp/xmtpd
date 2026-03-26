package message

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	feeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/fees"
	registryTestUtils "github.com/xmtp/xmtpd/pkg/testutils/registry"
	"github.com/xmtp/xmtpd/pkg/utils"
)

const testNodeID = uint32(100)

func createTestRegistrant(
	t *testing.T,
	ctx context.Context,
	store *db.Handler,
) *registrant.Registrant {
	t.Helper()
	privKey := testutils.RandomPrivateKey(t)
	node := registry.Node{
		NodeID:      testNodeID,
		SigningKey:  &privKey.PublicKey,
		IsCanonical: true,
	}
	mockRegistry := registryTestUtils.CreateMockRegistry(t, []registry.Node{node})

	reg, err := registrant.NewRegistrant(
		ctx,
		testutils.NewLog(t),
		store.WriteQuery(),
		mockRegistry,
		utils.EcdsaPrivateKeyToString(privKey),
		nil,
	)
	require.NoError(t, err)
	return reg
}

func setupPublishWorkerTest(
	t *testing.T,
) (context.Context, *db.Handler, *publishWorker) {
	t.Helper()
	ctx := t.Context()

	store, _ := testutils.NewDB(t, ctx)
	reg := createTestRegistrant(t, ctx, store)

	worker, err := startPublishWorker(
		ctx,
		testutils.NewLog(t),
		reg,
		store,
		feeTestUtils.NewTestFeeCalculator(),
		10*time.Millisecond,
	)
	require.NoError(t, err)

	return ctx, store, worker
}

func stageSingleEnvelope(
	t *testing.T,
	ctx context.Context,
	store *db.Handler,
) queries.StagedOriginatorEnvelope {
	t.Helper()
	clientEnv := envelopeTestUtils.CreateClientEnvelope()
	payerEnv := envelopeTestUtils.CreatePayerEnvelope(t, testNodeID, clientEnv)

	staged, err := store.WriteQuery().InsertStagedOriginatorEnvelope(
		ctx,
		queries.InsertStagedOriginatorEnvelopeParams{
			Topic:         clientEnv.GetAad().GetTargetTopic(),
			PayerEnvelope: testutils.Marshal(t, payerEnv),
		},
	)
	require.NoError(t, err)
	return staged
}

func stageEnvelopes(
	t *testing.T,
	ctx context.Context,
	store *db.Handler,
	count int,
) []queries.StagedOriginatorEnvelope {
	t.Helper()
	staged := make([]queries.StagedOriginatorEnvelope, 0, count)
	for range count {
		staged = append(staged, stageSingleEnvelope(t, ctx, store))
	}
	return staged
}

func waitForProcessed(
	t *testing.T,
	worker *publishWorker,
	minID int64,
) {
	t.Helper()
	require.Eventually(t, func() bool {
		return worker.lastProcessed.Load() >= minID
	}, 5*time.Second, 10*time.Millisecond)
}

func fetchGatewayEnvelopes(
	t *testing.T,
	ctx context.Context,
	store *db.Handler,
	limit int32,
) []queries.GatewayEnvelopesView {
	t.Helper()
	envs, err := store.ReadQuery().SelectGatewayEnvelopesUnfiltered(ctx,
		queries.SelectGatewayEnvelopesUnfilteredParams{
			RowLimit:          limit,
			CursorNodeIds:     []int32{int32(testNodeID)},
			CursorSequenceIds: []int64{0},
		},
	)
	require.NoError(t, err)
	return envs
}

func parseOriginatorEnvelope(
	t *testing.T,
	raw []byte,
) *envelopes.OriginatorEnvelope {
	t.Helper()
	env, err := envelopes.NewOriginatorEnvelopeFromBytes(raw)
	require.NoError(t, err)
	return env
}

func requireNoStagedEnvelopes(
	t *testing.T,
	ctx context.Context,
	store *db.Handler,
) {
	t.Helper()
	remaining, err := store.ReadQuery().SelectStagedOriginatorEnvelopes(ctx,
		queries.SelectStagedOriginatorEnvelopesParams{
			LastSeenID: 0,
			NumRows:    100,
		},
	)
	require.NoError(t, err)
	require.Empty(t, remaining)
}

// TestPublishWorkerProcessesSingleEnvelope verifies that the publish worker
// picks up a staged envelope, inserts it into gateway_envelopes with valid
// fees, and deletes it from the staging table.
func TestPublishWorkerProcessesSingleEnvelope(t *testing.T) {
	ctx, store, worker := setupPublishWorkerTest(t)

	staged := stageSingleEnvelope(t, ctx, store)
	worker.notifyStagedPublish()
	waitForProcessed(t, worker, staged.ID)

	// Verify gateway envelope was inserted with correct sequence ID
	latestSeq, err := store.ReadQuery().GetLatestSequenceId(ctx, int32(testNodeID))
	require.NoError(t, err)
	require.Equal(t, staged.ID, latestSeq)

	// Verify the originator envelope is valid and has non-zero base fee
	gatewayEnvs := fetchGatewayEnvelopes(t, ctx, store, 10)
	require.Len(t, gatewayEnvs, 1)
	origEnv := parseOriginatorEnvelope(t, gatewayEnvs[0].OriginatorEnvelope)
	require.Greater(t, origEnv.UnsignedOriginatorEnvelope.BaseFee(), currency.PicoDollar(0),
		"non-reserved topic should have non-zero base fee")

	requireNoStagedEnvelopes(t, ctx, store)
}

// TestPublishWorkerLastProcessedAccuracy verifies that lastProcessed is set
// correctly after processing, reflecting the latest gateway sequence ID.
func TestPublishWorkerLastProcessedAccuracy(t *testing.T) {
	ctx, store, worker := setupPublishWorkerTest(t)

	require.Equal(t, int64(0), worker.lastProcessed.Load())

	staged1 := stageSingleEnvelope(t, ctx, store)
	worker.notifyStagedPublish()
	waitForProcessed(t, worker, staged1.ID)
	require.Equal(t, staged1.ID, worker.lastProcessed.Load())

	staged2 := stageSingleEnvelope(t, ctx, store)
	worker.notifyStagedPublish()
	waitForProcessed(t, worker, staged2.ID)
	require.Equal(t, staged2.ID, worker.lastProcessed.Load())
}

// TestPublishWorkerHighThroughput verifies that when more rows than
// numRowsPerBatch are staged, the worker processes all of them in multiple
// batches without waiting for additional triggers.
func TestPublishWorkerHighThroughput(t *testing.T) {
	ctx, store, worker := setupPublishWorkerTest(t)

	totalEnvelopes := int(numRowsPerBatch) + 50
	staged := stageEnvelopes(t, ctx, store, totalEnvelopes)
	lastID := staged[len(staged)-1].ID

	worker.notifyStagedPublish()

	require.Eventually(t, func() bool {
		return worker.lastProcessed.Load() >= lastID
	}, 10*time.Second, 10*time.Millisecond,
		"worker should process all %d envelopes across multiple batches", totalEnvelopes)

	latestSeq, err := store.ReadQuery().GetLatestSequenceId(ctx, int32(testNodeID))
	require.NoError(t, err)
	require.Equal(t, lastID, latestSeq)

	requireNoStagedEnvelopes(t, ctx, store)
}

// TestPublishWorkerOrderingGuarantee verifies that envelopes are inserted
// in gateway_envelopes in the same order they were staged (by ID).
func TestPublishWorkerOrderingGuarantee(t *testing.T) {
	ctx, store, worker := setupPublishWorkerTest(t)

	count := 10
	staged := stageEnvelopes(t, ctx, store, count)
	lastID := staged[len(staged)-1].ID

	worker.notifyStagedPublish()
	waitForProcessed(t, worker, lastID)

	envelopes, err := store.ReadQuery().SelectGatewayEnvelopesUnfiltered(ctx,
		queries.SelectGatewayEnvelopesUnfilteredParams{
			RowLimit:          int32(count + 1),
			CursorNodeIds:     []int32{int32(testNodeID)},
			CursorSequenceIds: []int64{0},
		},
	)
	require.NoError(t, err)
	require.Len(t, envelopes, count)

	for i, env := range envelopes {
		require.Equal(t, staged[i].ID, env.OriginatorSequenceID,
			"envelope %d should have sequence_id matching staged ID", i)
	}
}

// TestPublishWorkerConcurrentWorkers verifies that two workers processing the
// same staging table produce no duplicates, no gaps, and correct ordering.
func TestPublishWorkerConcurrentWorkers(t *testing.T) {
	ctx := t.Context()

	store, _ := testutils.NewDB(t, ctx)
	log := testutils.NewLog(t)
	feeCalc := feeTestUtils.NewTestFeeCalculator()

	// Both workers share the same private key so they match the node_info row
	privKey := testutils.RandomPrivateKey(t)
	privKeyStr := utils.EcdsaPrivateKeyToString(privKey)

	createWorker := func() *publishWorker {
		node := registry.Node{
			NodeID:      testNodeID,
			SigningKey:  &privKey.PublicKey,
			IsCanonical: true,
		}
		mockReg := registryTestUtils.CreateMockRegistry(t, []registry.Node{node})

		reg, err := registrant.NewRegistrant(
			ctx, log, store.WriteQuery(), mockReg, privKeyStr, nil,
		)
		require.NoError(t, err)

		w, err := startPublishWorker(ctx, log, reg, store, feeCalc, 10*time.Millisecond)
		require.NoError(t, err)
		return w
	}

	worker1 := createWorker()
	worker2 := createWorker()

	count := 50
	staged := stageEnvelopes(t, ctx, store, count)
	lastID := staged[len(staged)-1].ID

	worker1.notifyStagedPublish()
	worker2.notifyStagedPublish()

	require.Eventually(t, func() bool {
		return worker1.lastProcessed.Load() >= lastID ||
			worker2.lastProcessed.Load() >= lastID
	}, 10*time.Second, 10*time.Millisecond,
		"at least one worker should process all envelopes")

	envelopes, err := store.ReadQuery().SelectGatewayEnvelopesUnfiltered(ctx,
		queries.SelectGatewayEnvelopesUnfilteredParams{
			RowLimit:          int32(count + 10),
			CursorNodeIds:     []int32{int32(testNodeID)},
			CursorSequenceIds: []int64{0},
		},
	)
	require.NoError(t, err)
	require.Len(t, envelopes, count, "should have exactly %d gateway envelopes (no dups)", count)

	for i := 1; i < len(envelopes); i++ {
		require.Greater(t, envelopes[i].OriginatorSequenceID,
			envelopes[i-1].OriginatorSequenceID,
			"gateway envelopes should be in ascending sequence_id order")
	}

	requireNoStagedEnvelopes(t, ctx, store)
}

// TestPublishWorkerEmptyTable verifies that the worker handles an empty
// staging table gracefully (no errors, no panics).
func TestPublishWorkerEmptyTable(t *testing.T) {
	_, _, worker := setupPublishWorkerTest(t)

	count, err := worker.processBatch()
	require.NoError(t, err)
	require.Equal(t, int32(0), count)
	require.Equal(t, int64(0), worker.lastProcessed.Load())
}

// TestPublishWorkerTimerPolling verifies that the worker picks up envelopes
// via the ticker even without an explicit notification.
func TestPublishWorkerTimerPolling(t *testing.T) {
	ctx, store, worker := setupPublishWorkerTest(t)

	staged := stageSingleEnvelope(t, ctx, store)
	// Do NOT call notifyStagedPublish — rely on the 1-second ticker

	require.Eventually(t, func() bool {
		return worker.lastProcessed.Load() >= staged.ID
	}, 5*time.Second, 50*time.Millisecond,
		"worker should pick up the envelope via ticker polling")
}

// TestPublishWorkerContextCancellation verifies that the worker stops
// cleanly when the context is cancelled.
func TestPublishWorkerContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())

	store, _ := testutils.NewDB(t, ctx)
	privKey := testutils.RandomPrivateKey(t)
	privKeyStr := utils.EcdsaPrivateKeyToString(privKey)

	node := registry.Node{
		NodeID:      testNodeID,
		SigningKey:  &privKey.PublicKey,
		IsCanonical: true,
	}
	mockRegistry := registryTestUtils.CreateMockRegistry(t, []registry.Node{node})

	reg, err := registrant.NewRegistrant(
		ctx, testutils.NewLog(t), store.WriteQuery(), mockRegistry, privKeyStr, nil,
	)
	require.NoError(t, err)

	worker, err := startPublishWorker(
		ctx, testutils.NewLog(t), reg, store,
		feeTestUtils.NewTestFeeCalculator(), 10*time.Millisecond,
	)
	require.NoError(t, err)

	stageSingleEnvelope(t, ctx, store)
	cancel()

	// processBatch should return promptly with a context error, not hang.
	_, err = worker.processBatch()
	require.Error(t, err)
}

// TestPublishWorkerPollAndPublishExitsOnCancel verifies that pollAndPublish
// returns promptly when the context is cancelled, rather than continuing to
// retry in the inner loop.
func TestPublishWorkerPollAndPublishExitsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())

	// Cancel the context before calling pollAndPublish
	cancel()

	worker := &publishWorker{
		ctx:                ctx,
		logger:             testutils.NewLog(t),
		sleepOnFailureTime: 10 * time.Millisecond,
		notifier:           make(chan bool, 1),
	}

	// pollAndPublish should return immediately instead of looping
	done := make(chan struct{})
	go func() {
		worker.pollAndPublish()
		close(done)
	}()

	select {
	case <-done:
		// Success: pollAndPublish exited promptly
	case <-time.After(2 * time.Second):
		t.Fatal("pollAndPublish did not exit after context cancellation")
	}
}

// TestPublishWorkerProcessBatchWithRetryExitsOnCancel verifies that
// processBatchWithRetry returns promptly when the context is cancelled.
func TestPublishWorkerProcessBatchWithRetryExitsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())

	cancel()

	worker := &publishWorker{
		ctx:                ctx,
		logger:             testutils.NewLog(t),
		sleepOnFailureTime: 10 * time.Millisecond,
	}

	count, err := worker.processBatchWithRetry()
	require.Equal(t, int32(0), count)
	require.ErrorIs(t, err, context.Canceled)
}

// TestPublishWorkerReservedTopicNoFees verifies that envelopes on reserved topics
// (e.g., payer reports) are processed with zero base fee and zero congestion fee.
func TestPublishWorkerReservedTopicNoFees(t *testing.T) {
	ctx, store, worker := setupPublishWorkerTest(t)

	clientEnv := envelopeTestUtils.CreatePayerReportClientEnvelope(testNodeID)
	payerEnv := envelopeTestUtils.CreatePayerEnvelope(t, testNodeID, clientEnv)

	staged, err := store.WriteQuery().InsertStagedOriginatorEnvelope(
		ctx,
		queries.InsertStagedOriginatorEnvelopeParams{
			Topic:         clientEnv.GetAad().GetTargetTopic(),
			PayerEnvelope: testutils.Marshal(t, payerEnv),
		},
	)
	require.NoError(t, err)

	worker.notifyStagedPublish()
	waitForProcessed(t, worker, staged.ID)

	// Verify the originator envelope has zero fees
	gatewayEnvs := fetchGatewayEnvelopes(t, ctx, store, 10)
	require.Len(t, gatewayEnvs, 1)
	origEnv := parseOriginatorEnvelope(t, gatewayEnvs[0].OriginatorEnvelope)
	require.Equal(t, currency.PicoDollar(0), origEnv.UnsignedOriginatorEnvelope.BaseFee(),
		"reserved topic should have zero base fee")
	require.Equal(t, currency.PicoDollar(0), origEnv.UnsignedOriginatorEnvelope.CongestionFee(),
		"reserved topic should have zero congestion fee")
}
