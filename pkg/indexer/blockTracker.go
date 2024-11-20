package indexer

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/xmtp/xmtpd/pkg/db/queries"
)

/*
*
BlockTracker keeps a database record of the latest block that has been indexed for a contract address
and allows the user to increase the value.
*
*/
type BlockTracker struct {
	latestBlock     atomic.Uint64
	contractAddress string
	queries         *queries.Queries
	mu              sync.Mutex
}

// Return a new BlockTracker initialized to the latest block from the DB
func NewBlockTracker(
	ctx context.Context,
	contractAddress string,
	queries *queries.Queries,
) (*BlockTracker, error) {
	bt := &BlockTracker{
		contractAddress: contractAddress,
		queries:         queries,
	}

	latestBlock, err := getLatestBlock(ctx, contractAddress, queries)
	if err != nil {
		return nil, err
	}
	bt.latestBlock.Store(latestBlock)

	return bt, nil
}

func (bt *BlockTracker) GetLatestBlock() uint64 {
	return bt.latestBlock.Load()
}

func (bt *BlockTracker) UpdateLatestBlock(ctx context.Context, block uint64) error {
	// Quick check without lock
	if block <= bt.latestBlock.Load() {
		return nil
	}

	bt.mu.Lock()
	defer bt.mu.Unlock()

	// Re-check after acquiring lock
	if block <= bt.latestBlock.Load() {
		return nil
	}

	if err := bt.updateDB(ctx, block); err != nil {
		return err
	}

	bt.latestBlock.Store(block)
	return nil
}

func (bt *BlockTracker) updateDB(ctx context.Context, block uint64) error {
	return bt.queries.SetLatestBlock(ctx, queries.SetLatestBlockParams{
		ContractAddress: bt.contractAddress,
		BlockNumber:     int64(block),
	})
}

func getLatestBlock(
	ctx context.Context,
	contractAddress string,
	querier *queries.Queries,
) (uint64, error) {
	latestBlock, err := querier.GetLatestBlock(ctx, contractAddress)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return uint64(latestBlock), nil
}
