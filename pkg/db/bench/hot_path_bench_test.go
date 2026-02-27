//go:build bench

package bench

import (
	"context"
	"database/sql"
	"log"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
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
func seedHotPath(ctx context.Context, benchDB *sql.DB) {
	q := queries.New(benchDB)
	hotPathPayerIDs = make([]int32, hotPathPayerCount)

	for i := range hotPathPayerCount {
		addr := utils.HexEncode(testutils.RandomBytes(20))

		id, err := q.FindOrCreatePayer(ctx, addr)
		if err != nil {
			log.Fatalf("seed hot path payer: %v", err)
		}

		hotPathPayerIDs[i] = id
	}

	// Pre-create gateway partitions so write benchmarks never hit partition-creation overhead.
	for seqID := int64(0); seqID < 10*db.GatewayEnvelopeBandWidth; seqID += db.GatewayEnvelopeBandWidth {
		_ = q.EnsureGatewayParts(ctx, queries.EnsureGatewayPartsParams{
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
		_, err := q.InsertStagedOriginatorEnvelope(
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
		hotPathPayerCount, 10, hotPathStagedSeedRows,
	)
}

// BenchmarkHotPathInsertStaged measures INSERT into staged_originator_envelopes.
// The underlying SQL function uses pg_advisory_xact_lock to serialize sequence
// ID assignment, so this reflects the true serialised ingest cost.
func BenchmarkHotPathInsertStaged(b *testing.B) {
	var (
		q     = queries.New(hotPathDB)
		topic = testutils.RandomBytes(32)
		blob  = testutils.RandomBytes(hotPathBlobSize)
	)

	for b.Loop() {
		_, err := q.InsertStagedOriginatorEnvelope(
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
	q := queries.New(hotPathDB)
	for b.Loop() {
		_, err := q.SelectStagedOriginatorEnvelopes(
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
//   - INSERT into gateway_envelopes_meta + gateway_envelope_blobs
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
			queries.InsertGatewayEnvelopeParams{
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
		q     = queries.New(hotPathDB)
		topic = testutils.RandomBytes(32)
		blob  = testutils.RandomBytes(hotPathBlobSize)
	)

	for b.Loop() {
		b.StopTimer()
		row, err := q.InsertStagedOriginatorEnvelope(
			benchCtx,
			queries.InsertStagedOriginatorEnvelopeParams{
				Topic:         topic,
				PayerEnvelope: blob,
			},
		)
		require.NoError(b, err)
		b.StartTimer()

		_, err = q.DeleteStagedOriginatorEnvelope(benchCtx, row.ID)
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
		q          = queries.New(hotPathDB)
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
		staged, err := q.InsertStagedOriginatorEnvelope(
			benchCtx,
			queries.InsertStagedOriginatorEnvelopeParams{
				Topic:         topic,
				PayerEnvelope: blob,
			},
		)
		require.NoError(b, err)

		// 2. Publish worker polls and fetches the staged envelope.
		_, err = q.SelectStagedOriginatorEnvelopes(
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
			queries.InsertGatewayEnvelopeParams{
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
		_, err = q.DeleteStagedOriginatorEnvelope(benchCtx, staged.ID)
		require.NoError(b, err)
	}
}
