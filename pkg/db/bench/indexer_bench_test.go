package bench

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

const numContracts = 100

// seedIndexer populates the latest_block table with contract addresses.
func seedIndexer(ctx context.Context, db *sql.DB) {
	q := queries.New(db)
	indexerContracts = make([]string, numContracts)
	for i := range numContracts {
		addr := fmt.Sprintf("0x%040x", i)
		indexerContracts[i] = addr
		err := q.SetLatestBlock(ctx, queries.SetLatestBlockParams{
			ContractAddress: addr,
			BlockNumber:     int64(1000 + i),
			BlockHash:       testutils.RandomBytes(32),
		})
		if err != nil {
			log.Fatalf("seed indexer: %v", err)
		}
	}
	log.Printf("seeded indexer: %d contracts", numContracts)
}

func BenchmarkGetLatestBlock(b *testing.B) {
	q := queries.New(indexerDB)
	contract := indexerContracts[0]
	for b.Loop() {
		_, err := q.GetLatestBlock(benchCtx, contract)
		require.NoError(b, err)
	}
}

func BenchmarkSetLatestBlock(b *testing.B) {
	q := queries.New(indexerDB)
	blockHash := testutils.RandomBytes(32) // pre-generate to avoid crypto/rand in hot path
	var counter atomic.Int64
	counter.Store(100_000)
	for b.Loop() {
		blockNum := counter.Add(1)
		err := q.SetLatestBlock(benchCtx, queries.SetLatestBlockParams{
			ContractAddress: indexerContracts[0],
			BlockNumber:     blockNum,
			BlockHash:       blockHash,
		})
		require.NoError(b, err)
	}
}
