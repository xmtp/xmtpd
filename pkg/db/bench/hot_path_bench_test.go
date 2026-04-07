//go:build bench

package bench

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/db/types"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/utils"
)

const (
	hotPathOriginatorID = int32(1337)
	hotPathPayerCount   = 5
	hotPathBlobSize     = 500

	// Number of staged rows pre-seeded for the select benchmark.
	hotPathStagedSeedRows = 10_000
)

// seedHotPath pre-creates payers and gateway partitions, and seeds staged rows
// for SELECT benchmarks.
func seedHotPath(ctx context.Context) {
	hotPathPayerIDs = make([]int32, hotPathPayerCount)

	for i := range hotPathPayerCount {
		addr := utils.HexEncode(testutils.RandomBytes(20))

		id, err := hotPathQueries.FindOrCreatePayer(ctx, addr)
		if err != nil {
			log.Fatalf("seed hot path payer: %v", err)
		}

		hotPathPayerIDs[i] = id
	}

	// Pre-create gateway partitions so write benchmarks never hit partition-creation overhead.
	for seqID := int64(0); seqID < 50*db.GatewayEnvelopeBandWidth; seqID += db.GatewayEnvelopeBandWidth {
		_ = hotPathQueries.EnsureGatewayPartsV3(ctx, queries.EnsureGatewayPartsV3Params{
			OriginatorNodeID:     hotPathOriginatorID,
			OriginatorSequenceID: seqID,
			BandWidth:            db.GatewayEnvelopeBandWidth,
		})
	}

	// Seed staged rows for the SELECT benchmark.
	var (
		topic = testutils.RandomBytes(32)
		blob  = testutils.RandomBytes(hotPathBlobSize)
	)

	for i := range hotPathStagedSeedRows {
		_, err := hotPathQueries.InsertStagedOriginatorEnvelope(
			ctx,
			queries.InsertStagedOriginatorEnvelopeParams{
				Topic:         topic,
				PayerEnvelope: blob,
			},
		)
		if err != nil {
			log.Fatalf("seed staged envelope %d: %v", i, err)
		}
	}

	log.Printf(
		"seeded hot path: %d payers, %d gateway partitions, %d staged rows",
		hotPathPayerCount, 50, hotPathStagedSeedRows,
	)
}

// BenchmarkHotPathInsertStaged measures INSERT into staged_originator_envelopes.
// The underlying SQL function uses pg_advisory_xact_lock to serialize sequence
// ID assignment, so this reflects the true serialised ingest cost.
func BenchmarkHotPathInsertStaged(b *testing.B) {
	var (
		topic = testutils.RandomBytes(32)
		blob  = testutils.RandomBytes(hotPathBlobSize)
	)

	for b.Loop() {
		_, err := hotPathQueries.InsertStagedOriginatorEnvelope(
			benchCtx,
			queries.InsertStagedOriginatorEnvelopeParams{
				Topic:         topic,
				PayerEnvelope: blob,
			},
		)
		require.NoError(b, err)
	}
}

// BenchmarkHotPathSelectStaged measures SELECT from staged_originator_envelopes
// as the publish worker does — fetching the next batch from a cursor position.
func BenchmarkHotPathSelectStaged(b *testing.B) {
	for b.Loop() {
		_, err := hotPathQueries.SelectStagedOriginatorEnvelopes(
			benchCtx,
			queries.SelectStagedOriginatorEnvelopesParams{
				LastSeenID: 0,
				NumRows:    100,
			},
		)
		require.NoError(b, err)
	}
}

// BenchmarkHotPathInsertGatewayWithUsage measures the atomic transaction that
// the publish worker executes after signing a staged envelope:
//   - INSERT into gateway_envelopes_meta + gateway_envelopes_blob
//   - UPSERT unsettled_usage
//   - UPSERT originator_congestion
func BenchmarkHotPathInsertGatewayWithUsage(b *testing.B) {
	var (
		topic      = testutils.RandomBytes(32)
		blob       = testutils.RandomBytes(hotPathBlobSize)
		payerID    = hotPathPayerIDs[0]
		now        = time.Now()
		expiry     = now.Add(24 * time.Hour).Unix()
		minute     = utils.MinutesSinceEpoch(now)
		sequenceID atomic.Int64
	)

	// Start beyond the seeded range to avoid collisions with other benchmarks.
	sequenceID.Store(10_000_000)

	for b.Loop() {
		seqID := sequenceID.Add(1)
		_, err := db.InsertGatewayEnvelopeAndIncrementUnsettledUsage(
			benchCtx,
			hotPathDB,
			queries.InsertGatewayEnvelopeV3Params{
				OriginatorNodeID:     hotPathOriginatorID,
				OriginatorSequenceID: seqID,
				Topic:                topic,
				OriginatorEnvelope:   blob,
				PayerID:              db.NullInt32(payerID),
				GatewayTime:          now,
				Expiry:               expiry,
			},
			queries.IncrementUnsettledUsageParams{
				PayerID:           payerID,
				OriginatorID:      hotPathOriginatorID,
				MinutesSinceEpoch: minute,
				SpendPicodollars:  1_000_000,
			},
			true,
		)
		require.NoError(b, err)
	}
}

// BenchmarkHotPathDeleteStaged measures DELETE from staged_originator_envelopes.
// A fresh row is inserted (outside timer) before each deletion so that every
// measured delete always hits a real row.
func BenchmarkHotPathDeleteStaged(b *testing.B) {
	var (
		topic = testutils.RandomBytes(32)
		blob  = testutils.RandomBytes(hotPathBlobSize)
	)

	for b.Loop() {
		b.StopTimer()
		row, err := hotPathQueries.InsertStagedOriginatorEnvelope(
			benchCtx,
			queries.InsertStagedOriginatorEnvelopeParams{
				Topic:         topic,
				PayerEnvelope: blob,
			},
		)
		require.NoError(b, err)
		b.StartTimer()

		_, err = hotPathQueries.BulkDeleteStagedOriginatorEnvelopes(benchCtx, []int64{row.ID})
		require.NoError(b, err)
	}
}

// BenchmarkHotPathFullCycle measures the complete publish-worker cycle for a
// single payer envelope, excluding the cryptographic signing step:
//  1. INSERT into staged_originator_envelopes
//  2. SELECT the staged row back (simulates worker polling)
//  3. INSERT gateway envelope + increment unsettled_usage + increment congestion (atomic tx)
//  4. DELETE the staged row
func BenchmarkHotPathFullCycle(b *testing.B) {
	var (
		topic      = testutils.RandomBytes(32)
		blob       = testutils.RandomBytes(hotPathBlobSize)
		payerID    = hotPathPayerIDs[0]
		now        = time.Now()
		expiry     = now.Add(24 * time.Hour).Unix()
		minute     = utils.MinutesSinceEpoch(now)
		gatewaySeq atomic.Int64
	)

	// Start beyond the seeded range to avoid collisions with other benchmarks.
	gatewaySeq.Store(20_000_000)

	for b.Loop() {
		// 1. Payer client stages the envelope.
		staged, err := hotPathQueries.InsertStagedOriginatorEnvelope(
			benchCtx,
			queries.InsertStagedOriginatorEnvelopeParams{
				Topic:         topic,
				PayerEnvelope: blob,
			},
		)
		require.NoError(b, err)

		// 2. Publish worker polls and fetches the staged envelope.
		_, err = hotPathQueries.SelectStagedOriginatorEnvelopes(
			benchCtx,
			queries.SelectStagedOriginatorEnvelopesParams{
				LastSeenID: staged.ID - 1,
				NumRows:    1,
			},
		)
		require.NoError(b, err)

		// 3. Worker inserts the originator envelope and tracks usage/congestion.
		seqID := gatewaySeq.Add(1)
		_, err = db.InsertGatewayEnvelopeAndIncrementUnsettledUsage(
			benchCtx,
			hotPathDB,
			queries.InsertGatewayEnvelopeV3Params{
				OriginatorNodeID:     hotPathOriginatorID,
				OriginatorSequenceID: seqID,
				Topic:                topic,
				OriginatorEnvelope:   blob,
				PayerID:              db.NullInt32(payerID),
				GatewayTime:          now,
				Expiry:               expiry,
			},
			queries.IncrementUnsettledUsageParams{
				PayerID:           payerID,
				OriginatorID:      hotPathOriginatorID,
				MinutesSinceEpoch: minute,
				SpendPicodollars:  1_000_000,
			},
			true,
		)
		require.NoError(b, err)

		// 4. Worker removes the processed staged envelope.
		_, err = hotPathQueries.BulkDeleteStagedOriginatorEnvelopes(benchCtx, []int64{staged.ID})
		require.NoError(b, err)
	}
}

// BenchmarkHotPathBatchCycle measures the batched publish-worker cycle using
// the V2 batch operations: BulkFindOrCreatePayers, InsertGatewayEnvelopeBatchV2,
// and BulkDeleteStagedOriginatorEnvelopes. This is the new hot path after the
// batch refactor, reducing DB round-trips from ~4*N to 3 per batch.
//
// Each iteration:
//
//	[untimed] Seeds N staged envelopes and selects them as a batch.
//	[timed]   BulkFindOrCreatePayers → InsertGatewayEnvelopeBatchV2 → BulkDeleteStaged.
//
// Sub-benchmarks: batch=1, batch=10, batch=100, batch=500
func BenchmarkHotPathBatchCycle(b *testing.B) {
	batchSizes := []int{1, 10, 100, 500}

	for _, batchSize := range batchSizes {
		b.Run(fmt.Sprintf("batch=%d", batchSize), func(b *testing.B) {
			var (
				topic      = testutils.RandomBytes(32)
				blob       = testutils.RandomBytes(hotPathBlobSize)
				payerAddr  = utils.HexEncode(testutils.RandomBytes(20))
				now        = time.Now()
				expiry     = now.Add(24 * time.Hour).Unix()
				gatewaySeq atomic.Int64
			)

			// Start at 40M to avoid collisions with other hot path benchmarks.
			gatewaySeq.Store(40_000_000)

			for b.Loop() {
				// --- Untimed: seed N staged envelopes ---
				b.StopTimer()
				var lastSeenID int64
				for range batchSize {
					staged, err := hotPathQueries.InsertStagedOriginatorEnvelope(
						benchCtx,
						queries.InsertStagedOriginatorEnvelopeParams{
							Topic:         topic,
							PayerEnvelope: blob,
						},
					)
					require.NoError(b, err)
					if lastSeenID == 0 {
						lastSeenID = staged.ID - 1
					}
				}

				// Fetch the batch (single SELECT, untimed).
				batch, err := hotPathQueries.SelectStagedOriginatorEnvelopes(
					benchCtx,
					queries.SelectStagedOriginatorEnvelopesParams{
						LastSeenID: lastSeenID,
						NumRows:    int32(batchSize),
					},
				)
				require.NoError(b, err)
				require.Len(b, batch, batchSize)
				b.StartTimer()

				// --- Timed: batch DB operations (3 round-trips) ---

				// 1. Bulk find/create payers.
				_, err = hotPathQueries.BulkFindOrCreatePayers(benchCtx, []string{payerAddr})
				require.NoError(b, err)

				// 2. Build batch and insert via V2 function.
				batchInput := types.NewGatewayEnvelopeBatch()
				stagedIDs := make([]int64, len(batch))
				for i, stagedEnv := range batch {
					seqID := gatewaySeq.Add(1)
					batchInput.Add(types.GatewayEnvelopeRow{
						OriginatorNodeID:     hotPathOriginatorID,
						OriginatorSequenceID: seqID,
						Topic:                topic,
						PayerID:              hotPathPayerIDs[0],
						GatewayTime:          now,
						Expiry:               expiry,
						OriginatorEnvelope:   blob,
						SpendPicodollars:     1_000_000,
						CountUsage:           true,
						CountCongestion:      true,
					})
					stagedIDs[i] = stagedEnv.ID
				}

				_, err = db.InsertGatewayEnvelopeBatchV2AndIncrementUnsettledUsage(
					benchCtx, hotPathDB, testutils.NewLog(b), batchInput,
				)
				require.NoError(b, err)

				// 3. Bulk delete staged envelopes.
				_, err = hotPathQueries.BulkDeleteStagedOriginatorEnvelopes(benchCtx, stagedIDs)
				require.NoError(b, err)
			}
		})
	}
}
