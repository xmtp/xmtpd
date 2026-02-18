package bench

import (
	"context"
	"database/sql"
	"log"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
)

const (
	numCongestionOriginators = 5
	numCongestionMinutes     = 2000 // per originator
)

// seedCongestion populates originator_congestion with time-bucketed message counts.
func seedCongestion(ctx context.Context, db *sql.DB) {
	q := queries.New(db)
	congestionOriginators = make([]int32, numCongestionOriginators)
	for i := range numCongestionOriginators {
		origID := int32(500 + i)
		congestionOriginators[i] = origID
		for minute := range int32(numCongestionMinutes) {
			err := q.IncrementOriginatorCongestion(ctx, queries.IncrementOriginatorCongestionParams{
				OriginatorID:      origID,
				MinutesSinceEpoch: minute,
			})
			if err != nil {
				log.Fatalf("seed congestion: %v", err)
			}
		}
	}
	congestionMaxMinute = numCongestionMinutes - 1
	log.Printf("seeded congestion: %d rows", numCongestionOriginators*numCongestionMinutes)
}

func BenchmarkIncrementOriginatorCongestion(b *testing.B) {
	q := queries.New(congestionDB)
	origID := congestionOriginators[0]
	var counter atomic.Int32
	counter.Store(100_000) // start beyond seeded range
	b.ResetTimer()
	for b.Loop() {
		minute := counter.Add(1)
		err := q.IncrementOriginatorCongestion(
			benchCtx,
			queries.IncrementOriginatorCongestionParams{
				OriginatorID:      origID,
				MinutesSinceEpoch: minute,
			},
		)
		require.NoError(b, err)
	}
}

func BenchmarkGetRecentOriginatorCongestion(b *testing.B) {
	q := queries.New(congestionDB)
	params := queries.GetRecentOriginatorCongestionParams{
		OriginatorID: congestionOriginators[0],
		EndMinute:    congestionMaxMinute,
		NumMinutes:   60, // last hour
	}
	b.ResetTimer()
	for b.Loop() {
		_, err := q.GetRecentOriginatorCongestion(benchCtx, params)
		require.NoError(b, err)
	}
}

func BenchmarkSumOriginatorCongestion(b *testing.B) {
	q := queries.New(congestionDB)
	params := queries.SumOriginatorCongestionParams{
		OriginatorID:        congestionOriginators[0],
		MinutesSinceEpochGt: 0,
		MinutesSinceEpochLt: int64(congestionMaxMinute),
	}
	b.ResetTimer()
	for b.Loop() {
		_, err := q.SumOriginatorCongestion(benchCtx, params)
		require.NoError(b, err)
	}
}
