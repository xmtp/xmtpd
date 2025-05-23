package common

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/db/queries"
)

/*
*
BlockTracker keeps a database record of the latest block that has been indexed for a contract address
and allows the user to increase the value.
*
*/
type BlockTracker struct {
	latest  *Block
	address common.Address
	queries *queries.Queries
	mu      sync.Mutex
}

type Block struct {
	number atomic.Uint64
	hash   common.Hash
}

func (b *Block) save(number uint64, hash []byte) {
	b.number.Store(number)
	b.hash = common.BytesToHash(hash)
}

var ErrEmptyBlockHash = errors.New("block hash is empty")

var _ IBlockTracker = &BlockTracker{}

// Return a new BlockTracker initialized to the latest block from the DB
func NewBlockTracker(
	ctx context.Context,
	client blockchain.ChainClient,
	address common.Address,
	queries *queries.Queries,
	startBlock uint64,
) (*BlockTracker, error) {
	bt := &BlockTracker{
		address: address,
		queries: queries,
	}

	latest, err := loadLatestBlock(ctx, client, address, queries, startBlock)
	if err != nil {
		return nil, err
	}
	bt.latest = latest

	return bt, nil
}

func (bt *BlockTracker) GetLatestBlock() (uint64, []byte) {
	bt.mu.Lock()
	defer bt.mu.Unlock()
	return bt.latest.number.Load(), bt.latest.hash.Bytes()
}

func (bt *BlockTracker) UpdateLatestBlock(
	ctx context.Context,
	block uint64,
	hashBytes []byte,
) error {
	// Quick check without lock
	if block <= bt.latest.number.Load() {
		return nil
	}

	bt.mu.Lock()
	defer bt.mu.Unlock()

	// Re-check after acquiring lock
	if block <= bt.latest.number.Load() {
		return nil
	}

	newHash := common.Hash(hashBytes)

	if newHash == (common.Hash{}) {
		return ErrEmptyBlockHash
	}

	if newHash == bt.latest.hash {
		return nil
	}

	if err := bt.updateDB(ctx, block, newHash.Bytes()); err != nil {
		return err
	}

	bt.latest.number.Store(block)
	bt.latest.hash = newHash

	return nil
}

func (bt *BlockTracker) updateDB(ctx context.Context, block uint64, hash []byte) error {
	return bt.queries.SetLatestBlock(ctx, queries.SetLatestBlockParams{
		ContractAddress: bt.address.Hex(),
		BlockNumber:     int64(block),
		BlockHash:       hash,
	})
}

// loadLatestBlock returns the latest block for an address.
// - returns an error if querying the database fails.
// - returns the stored block number and hash if the db contains a row for the address.
// - returns the start block number and hash if the db does not contain a row for the address and there is a start block number.
func loadLatestBlock(
	ctx context.Context,
	client blockchain.ChainClient,
	address common.Address,
	querier *queries.Queries,
	startBlock uint64,
) (*Block, error) {
	block := &Block{
		number: atomic.Uint64{},
		hash:   common.Hash{},
	}

	storedBlock, err := querier.GetLatestBlock(ctx, address.Hex())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return block, err
	}

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		onchainBlock, err := client.BlockByNumber(ctx, big.NewInt(int64(startBlock)))
		if err != nil {
			return nil, err
		}

		block.save(uint64(onchainBlock.NumberU64()), onchainBlock.Hash().Bytes())

		return block, nil
	}

	if storedBlock.BlockNumber < 0 {
		return block, fmt.Errorf(
			"invalid block number %d for contract %s",
			storedBlock.BlockNumber,
			address,
		)
	}

	block.save(uint64(storedBlock.BlockNumber), storedBlock.BlockHash)

	return block, nil
}
