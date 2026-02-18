package bench

import (
	"context"
	"database/sql"
	"encoding/hex"
	"log"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
)

const (
	numUsagePayers      = 50
	numUsageOriginators = 5
	numUsageMinutes     = 40
)

// seedUsage creates payers and populates unsettled_usage.
func seedUsage(ctx context.Context, db *sql.DB) {
	q := queries.New(db)
	usagePayerIDs = make([]int32, numUsagePayers)
	usageOriginators = make([]int32, numUsageOriginators)

	for i := range numUsageOriginators {
		usageOriginators[i] = int32(600 + i)
	}

	for i := range numUsagePayers {
		addr := hex.EncodeToString(randomBytes(20))
		id, err := q.FindOrCreatePayer(ctx, addr)
		if err != nil {
			log.Fatalf("seed usage payer: %v", err)
		}
		usagePayerIDs[i] = id

		for _, origID := range usageOriginators {
			for minute := range int32(numUsageMinutes) {
				err := q.IncrementUnsettledUsage(
					ctx,
					queries.IncrementUnsettledUsageParams{
						PayerID:           id,
						OriginatorID:      origID,
						MinutesSinceEpoch: minute,
						SpendPicodollars:  1_000_000,
						SequenceID:        int64(minute),
						MessageCount:      1,
					},
				)
				if err != nil {
					log.Fatalf("seed usage: %v", err)
				}
			}
		}
	}
	usageMaxMinute = numUsageMinutes - 1
	log.Printf(
		"seeded usage: %d rows",
		numUsagePayers*numUsageOriginators*numUsageMinutes,
	)
}

func BenchmarkIncrementUnsettledUsage(b *testing.B) {
	q := queries.New(usageDB)
	payerID := usagePayerIDs[0]
	origID := usageOriginators[0]
	var counter atomic.Int32
	counter.Store(100_000) // beyond seeded range
	b.ResetTimer()
	for b.Loop() {
		minute := counter.Add(1)
		err := q.IncrementUnsettledUsage(
			benchCtx,
			queries.IncrementUnsettledUsageParams{
				PayerID:           payerID,
				OriginatorID:      origID,
				MinutesSinceEpoch: minute,
				SpendPicodollars:  1_000_000,
				SequenceID:        int64(minute),
				MessageCount:      1,
			},
		)
		require.NoError(b, err)
	}
}
