package indexer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/db/queries"
)

/*
*
BlockTracker keeps a database record of the latest block that has been indexed for a contract address
and allows the user to increase the value.
*
*/
type BlockTracker struct {
	latestBlock     *Block
	contractAddress string
	queries         *queries.Queries
	mu              sync.Mutex
}

type Block struct {
	number atomic.Uint64
	hash   common.Hash
}

var (
	ErrEmptyBlockHash = errors.New("block hash is empty")
)

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

	latestBlock, err := loadLatestBlock(ctx, contractAddress, queries)
	if err != nil {
		return nil, err
	}
	bt.latestBlock = latestBlock

	return bt, nil
}

func (bt *BlockTracker) GetLatestBlock() (uint64, []byte) {
	bt.mu.Lock()
	defer bt.mu.Unlock()
	return bt.latestBlock.number.Load(), bt.latestBlock.hash.Bytes()
}

func (bt *BlockTracker) UpdateLatestBlock(
	ctx context.Context,
	block uint64,
	hashBytes []byte,
) error {
	// Quick check without lock
	if block <= bt.latestBlock.number.Load() {
		return nil
	}

	bt.mu.Lock()
	defer bt.mu.Unlock()

	// Re-check after acquiring lock
	if block <= bt.latestBlock.number.Load() {
		return nil
	}

	newHash := common.Hash(hashBytes)

	if newHash == (common.Hash{}) {
		return ErrEmptyBlockHash
	}

	if newHash == bt.latestBlock.hash {
		return nil
	}

	if err := bt.updateDB(ctx, block, newHash.Bytes()); err != nil {
		return err
	}

	bt.latestBlock.number.Store(block)
	bt.latestBlock.hash = newHash

	return nil
}

func (bt *BlockTracker) updateDB(ctx context.Context, block uint64, hash []byte) error {
	return bt.queries.SetLatestBlock(ctx, queries.SetLatestBlockParams{
		ContractAddress: bt.contractAddress,
		BlockNumber:     int64(block),
		BlockHash:       hash,
	})
}

func loadLatestBlock(
	ctx context.Context,
	contractAddress string,
	querier *queries.Queries,
) (*Block, error) {
	block := &Block{
		number: atomic.Uint64{},
		hash:   common.Hash{},
	}

	latestBlock, err := querier.GetLatestBlock(ctx, contractAddress)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return block, nil
		}
		return block, err
	}

	if latestBlock.BlockNumber < 0 {
		return block, fmt.Errorf(
			"invalid block number %d for contract %s",
			latestBlock.BlockNumber,
			contractAddress,
		)
	}

	block.number.Store(uint64(latestBlock.BlockNumber))
	block.hash = common.BytesToHash(latestBlock.BlockHash)

	return block, nil
}
