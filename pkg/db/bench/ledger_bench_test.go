package bench

import (
	"context"
	"database/sql"
	"log"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/utils"
)

const (
	numLedgerPayers         = 50
	numLedgerEventsPerPayer = 100
)

// seedLedger creates payers and populates payer_ledger_events.
func seedLedger(ctx context.Context, db *sql.DB) {
	q := queries.New(db)
	ledgerPayerIDs = make([]int32, numLedgerPayers)

	for i := range numLedgerPayers {
		addr := utils.HexEncode(testutils.RandomBytes(20))
		id, err := q.FindOrCreatePayer(ctx, addr)
		if err != nil {
			log.Fatalf("seed ledger payer: %v", err)
		}
		ledgerPayerIDs[i] = id

		// Insert events: mix of deposits (1), withdrawals (2), settlements (3)
		for j := range numLedgerEventsPerPayer {
			eventType := int16((j % 3) + 1)
			amount := int64(1_000_000) // 1M picodollars
			if eventType == 2 {
				amount = -amount
			}
			err := q.InsertPayerLedgerEvent(
				ctx,
				queries.InsertPayerLedgerEventParams{
					EventID:           testutils.RandomBytes(32),
					PayerID:           id,
					AmountPicodollars: amount,
					EventType:         eventType,
				},
			)
			if err != nil {
				log.Fatalf("seed ledger event: %v", err)
			}
		}
	}
	log.Printf(
		"seeded ledger: %d payers, %d events",
		numLedgerPayers,
		numLedgerPayers*numLedgerEventsPerPayer,
	)
}

func BenchmarkInsertPayerLedgerEvent(b *testing.B) {
	q := queries.New(ledgerDB)
	payerID := ledgerPayerIDs[0]
	// Pre-generate a pool of unique event IDs to avoid crypto/rand in hot path.
	const poolSize = 10_000
	eventIDs := make([][]byte, poolSize)
	for i := range poolSize {
		eventIDs[i] = testutils.RandomBytes(32)
	}
	var counter atomic.Int64
	for b.Loop() {
		idx := counter.Add(1)
		err := q.InsertPayerLedgerEvent(
			benchCtx,
			queries.InsertPayerLedgerEventParams{
				EventID:           eventIDs[idx%poolSize],
				PayerID:           payerID,
				AmountPicodollars: 1_000_000,
				EventType:         1,
			},
		)
		require.NoError(b, err)
	}
}

func BenchmarkGetPayerBalance(b *testing.B) {
	q := queries.New(ledgerDB)
	payerID := ledgerPayerIDs[0]
	for b.Loop() {
		_, err := q.GetPayerBalance(benchCtx, payerID)
		require.NoError(b, err)
	}
}

func BenchmarkGetLastEvent(b *testing.B) {
	q := queries.New(ledgerDB)
	params := queries.GetLastEventParams{
		PayerID:   ledgerPayerIDs[0],
		EventType: 1, // deposits
	}
	for b.Loop() {
		_, err := q.GetLastEvent(benchCtx, params)
		require.NoError(b, err)
	}
}
