//go:build bench

package bench

import (
	"context"
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
	numOriginators    = 3
	numTopics         = 100
	blobSize          = 500
	writeOriginatorID = int32(999) // dedicated originator for write benchmarks
	numBenchPayers    = 5
)

var envelopeOriginators = []int32{100, 200, 300}

// seedEnvelopes populates gateway_envelopes_meta and gateway_envelope_blobs.
func seedEnvelopes(ctx context.Context, tier *envelopeTier) {
	q := queries.New(tier.db)
	tier.originators = envelopeOriginators

	// Generate topics
	tier.topics = make([][]byte, numTopics)
	for i := range numTopics {
		tier.topics[i] = testutils.RandomBytes(32)
	}

	// Create payers for batch insert benchmarks
	tier.payerIDs = make([]int32, numBenchPayers)
	for i := range numBenchPayers {
		addr := utils.HexEncode(testutils.RandomBytes(20))
		id, err := q.FindOrCreatePayer(ctx, addr)
		if err != nil {
			log.Fatalf("seed envelope payer: %v", err)
		}
		tier.payerIDs[i] = id
	}

	// Pre-create partitions for seeded originators
	perOriginator := tier.count / numOriginators
	for seqID := int64(0); seqID < int64(perOriginator)+db.GatewayEnvelopeBandWidth; seqID += db.GatewayEnvelopeBandWidth {
		for _, origID := range tier.originators {
			_ = q.EnsureGatewayParts(ctx, queries.EnsureGatewayPartsParams{
				OriginatorNodeID:     origID,
				OriginatorSequenceID: seqID,
				BandWidth:            db.GatewayEnvelopeBandWidth,
			})
		}
	}
	// Partitions for write benchmark originator
	for seqID := int64(0); seqID < 10*db.GatewayEnvelopeBandWidth; seqID += db.GatewayEnvelopeBandWidth {
		_ = q.EnsureGatewayParts(ctx, queries.EnsureGatewayPartsParams{
			OriginatorNodeID:     writeOriginatorID,
			OriginatorSequenceID: seqID,
			BandWidth:            db.GatewayEnvelopeBandWidth,
		})
	}

	// Seed envelopes distributed across originators and topics
	batchSize := 10_000
	blob := testutils.RandomBytes(blobSize) // reuse same blob for speed
	seqIDs := make([]int64, numOriginators)

	for i := range tier.count {
		origIdx := i % numOriginators
		origID := tier.originators[origIdx]
		seqIDs[origIdx]++
		topicIdx := i % numTopics

		_, err := db.InsertGatewayEnvelopeWithChecksStandalone(
			ctx,
			q,
			queries.InsertGatewayEnvelopeParams{
				OriginatorNodeID:     origID,
				OriginatorSequenceID: seqIDs[origIdx],
				Topic:                tier.topics[topicIdx],
				Expiry:               time.Now().Add(24 * time.Hour).Unix(),
				OriginatorEnvelope:   blob,
			},
		)
		if err != nil {
			log.Fatalf("seed envelope %d: %v", i, err)
		}

		if (i+1)%batchSize == 0 {
			log.Printf(
				"seeded %d/%d envelopes for tier %s",
				i+1, tier.count, tier.name,
			)
		}
	}
	log.Printf("seeded envelopes: %d rows for tier %s", tier.count, tier.name)
}

// --- Read benchmarks ---

func BenchmarkSelectGatewayEnvelopesByTopics(b *testing.B) {
	for _, tier := range envelopeTiers {
		b.Run(tier.name, func(b *testing.B) {
			q := queries.New(tier.db)
			midSeq := int64(tier.count / numOriginators / 2)
			params := queries.SelectGatewayEnvelopesByTopicsParams{
				Topics:            tier.topics[:10],
				RowLimit:          100,
				CursorNodeIds:     tier.originators,
				CursorSequenceIds: []int64{midSeq, midSeq, midSeq},
			}
			for b.Loop() {
				_, err := q.SelectGatewayEnvelopesByTopics(benchCtx, params)
				require.NoError(b, err)
			}
		})
	}
}

func BenchmarkSelectGatewayEnvelopesByOriginators(b *testing.B) {
	for _, tier := range envelopeTiers {
		b.Run(tier.name, func(b *testing.B) {
			q := queries.New(tier.db)
			midSeq := int64(tier.count / numOriginators / 2)
			params := queries.SelectGatewayEnvelopesByOriginatorsParams{
				OriginatorNodeIds: tier.originators,
				RowsPerOriginator: 50,
				RowLimit:          100,
				CursorNodeIds:     tier.originators,
				CursorSequenceIds: []int64{midSeq, midSeq, midSeq},
			}
			for b.Loop() {
				_, err := q.SelectGatewayEnvelopesByOriginators(
					benchCtx,
					params,
				)
				require.NoError(b, err)
			}
		})
	}
}

func BenchmarkSelectGatewayEnvelopesBySingleOriginator(b *testing.B) {
	for _, tier := range envelopeTiers {
		b.Run(tier.name, func(b *testing.B) {
			q := queries.New(tier.db)
			midSeq := int64(tier.count / numOriginators / 2)
			params := queries.SelectGatewayEnvelopesBySingleOriginatorParams{
				OriginatorNodeID: tier.originators[0],
				CursorSequenceID: midSeq,
				RowLimit:         100,
			}
			for b.Loop() {
				_, err := q.SelectGatewayEnvelopesBySingleOriginator(
					benchCtx,
					params,
				)
				require.NoError(b, err)
			}
		})
	}
}

func BenchmarkSelectGatewayEnvelopesUnfiltered(b *testing.B) {
	for _, tier := range envelopeTiers {
		b.Run(tier.name, func(b *testing.B) {
			q := queries.New(tier.db)
			midSeq := int64(tier.count / numOriginators / 2)
			params := queries.SelectGatewayEnvelopesUnfilteredParams{
				RowLimit:          100,
				CursorNodeIds:     tier.originators,
				CursorSequenceIds: []int64{midSeq, midSeq, midSeq},
			}
			for b.Loop() {
				_, err := q.SelectGatewayEnvelopesUnfiltered(
					benchCtx,
					params,
				)
				require.NoError(b, err)
			}
		})
	}
}

func BenchmarkSelectNewestFromTopics(b *testing.B) {
	for _, tier := range envelopeTiers {
		b.Run(tier.name, func(b *testing.B) {
			q := queries.New(tier.db)
			topics := tier.topics[:10]
			for b.Loop() {
				_, err := q.SelectNewestFromTopics(benchCtx, topics)
				require.NoError(b, err)
			}
		})
	}
}

// --- Write benchmarks ---

func BenchmarkInsertGatewayEnvelope(b *testing.B) {
	for _, tier := range envelopeTiers {
		b.Run(tier.name, func(b *testing.B) {
			q := queries.New(tier.db)
			blob := testutils.RandomBytes(blobSize)
			topic := tier.topics[0]
			expiry := time.Now().Add(24 * time.Hour).Unix()
			var counter atomic.Int64
			counter.Store(1_000_000)
			for b.Loop() {
				seqID := counter.Add(1)
				_, err := q.InsertGatewayEnvelope(
					benchCtx,
					queries.InsertGatewayEnvelopeParams{
						OriginatorNodeID:     writeOriginatorID,
						OriginatorSequenceID: seqID,
						Topic:                topic,
						Expiry:               expiry,
						OriginatorEnvelope:   blob,
					},
				)
				require.NoError(b, err)
			}
		})
	}
}

func BenchmarkInsertGatewayEnvelopeBatch(b *testing.B) {
	for _, tier := range envelopeTiers {
		b.Run(tier.name, func(b *testing.B) {
			q := queries.New(tier.db)
			blob := testutils.RandomBytes(blobSize)
			batchLen := 10

			// Pre-allocate slices to avoid allocations in hot path.
			nodeIDs := make([]int32, batchLen)
			seqIDs := make([]int64, batchLen)
			topics := make([][]byte, batchLen)
			payerIDs := make([]int32, batchLen)
			times := make([]time.Time, batchLen)
			expiries := make([]int64, batchLen)
			blobs := make([][]byte, batchLen)
			spends := make([]int64, batchLen)
			for j := range batchLen {
				nodeIDs[j] = writeOriginatorID
				topics[j] = tier.topics[j%numTopics]
				payerIDs[j] = tier.payerIDs[j%numBenchPayers]
				blobs[j] = blob
				spends[j] = 1_000_000
			}

			var counter atomic.Int64
			counter.Store(5_000_000)
			for b.Loop() {
				baseSeq := counter.Add(int64(batchLen))
				now := time.Now()
				exp := now.Add(24 * time.Hour).Unix()
				for j := range batchLen {
					seqIDs[j] = baseSeq + int64(j)
					times[j] = now
					expiries[j] = exp
				}

				_, err := q.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
					benchCtx,
					queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams{
						OriginatorNodeIds:     nodeIDs,
						OriginatorSequenceIds: seqIDs,
						Topics:                topics,
						PayerIds:              payerIDs,
						GatewayTimes:          times,
						Expiries:              expiries,
						OriginatorEnvelopes:   blobs,
						SpendPicodollars:      spends,
					},
				)
				require.NoError(b, err)
			}
		})
	}
}
