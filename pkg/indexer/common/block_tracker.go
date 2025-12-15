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
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
)

// BlockTracker keeps a database record of the latest block that has been indexed for a contract address
// and allows the user to increase the value.
type BlockTracker struct {
	latest  *Block
	address common.Address
	db      *db.Handler
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

// NewBlockTracker returns a new BlockTracker initialized to the latest block from the DB.
func NewBlockTracker(
	ctx context.Context,
	client blockchain.ChainClient,
	address common.Address,
	db *db.Handler,
	deploymentBlock uint64,
) (*BlockTracker, error) {
	bt := &BlockTracker{
		address: address,
		db:      db,
	}

	latest, err := bt.loadLatestBlock(ctx, client, address, deploymentBlock)
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
	return bt.db.WriteQuery().SetLatestBlock(ctx, queries.SetLatestBlockParams{
		ContractAddress: bt.address.Hex(),
		BlockNumber:     int64(block),
		BlockHash:       hash,
	})
}

// loadLatestBlock returns the latest block for an address.
// - returns an error if querying the database fails.
// - returns the stored block number and hash if the db contains a valid block number for the address.
// - returns the deployment block number and hash if the db does not contain a block number for the address.
func (bt *BlockTracker) loadLatestBlock(
	ctx context.Context,
	client blockchain.ChainClient,
	address common.Address,
	deploymentBlock uint64,
) (*Block, error) {
	block := &Block{
		number: atomic.Uint64{},
		hash:   common.Hash{},
	}

	storedBlock, err := bt.db.ReadQuery().GetLatestBlock(ctx, address.Hex())
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			header, err := client.HeaderByNumber(ctx, big.NewInt(int64(deploymentBlock)))
			if err != nil {
				return nil, err
			}

			block.save(deploymentBlock, header.Hash().Bytes())

			return block, nil

		default:
			return nil, err
		}
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
