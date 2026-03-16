//go:build bench

package bench

import (
	"context"
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
func seedIndexer(ctx context.Context) {
	indexerContracts = make([]string, numContracts)
	for i := range numContracts {
		addr := fmt.Sprintf("0x%040x", i)
		indexerContracts[i] = addr
		err := indexerQueries.SetLatestBlock(ctx, queries.SetLatestBlockParams{
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
	contract := indexerContracts[0]
	for b.Loop() {
		_, err := indexerQueries.GetLatestBlock(benchCtx, contract)
		require.NoError(b, err)
	}
}

func BenchmarkSetLatestBlock(b *testing.B) {
	blockHash := testutils.RandomBytes(32) // pre-generate to avoid crypto/rand in hot path
	var counter atomic.Int64
	counter.Store(100_000)
	for b.Loop() {
		blockNum := counter.Add(1)
		err := indexerQueries.SetLatestBlock(benchCtx, queries.SetLatestBlockParams{
			ContractAddress: indexerContracts[0],
			BlockNumber:     blockNum,
			BlockHash:       blockHash,
		})
		require.NoError(b, err)
	}
}
