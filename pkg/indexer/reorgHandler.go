package indexer

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"

	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/db/queries"
)

type ChainReorgHandler interface {
	FindReorgPoint(detectedAt uint64) (uint64, []byte, error)
}

type ReorgHandler struct {
	ctx     context.Context
	client  blockchain.ChainClient
	queries *queries.Queries
}

var (
	ErrNoBlocksFound = errors.New("no blocks found")
	ErrGetBlock      = errors.New("failed to get block")
)

// TODO: Make this configurable?
const BLOCK_RANGE_SIZE uint64 = 1000

func NewChainReorgHandler(
	ctx context.Context,
	client blockchain.ChainClient,
	queries *queries.Queries,
) *ReorgHandler {
	return &ReorgHandler{
		ctx:     ctx,
		client:  client,
		queries: queries,
	}
}

// TODO: When reorg range has been calculated, alert clients (TBD)
func (r *ReorgHandler) FindReorgPoint(detectedAt uint64) (uint64, []byte, error) {
	startBlock, endBlock := blockRange(detectedAt)

	for {
		storedBlocks, err := r.queries.GetBlocksInRange(
			r.ctx,
			queries.GetBlocksInRangeParams{
				StartBlock: startBlock,
				EndBlock:   endBlock,
			},
		)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return 0, nil, fmt.Errorf("failed to get stored blocks: %w", err)
		}

		if len(storedBlocks) == 0 || errors.Is(err, sql.ErrNoRows) {
			if startBlock == 0 {
				return 0, nil, ErrNoBlocksFound
			}

			startBlock, endBlock = blockRange(startBlock)
			continue
		}

		oldestBlock := storedBlocks[0]
		chainBlock, err := r.client.BlockByNumber(r.ctx, big.NewInt(int64(oldestBlock.BlockNumber)))
		if err != nil {
			return 0, nil, fmt.Errorf("%w %d: %v", ErrGetBlock, oldestBlock.BlockNumber, err)
		}

		// Oldest block doesn't match, reorg happened earlier in the chain
		if !bytes.Equal(oldestBlock.BlockHash, chainBlock.Hash().Bytes()) {
			if startBlock == 0 {
				return 0, nil, ErrNoBlocksFound
			}

			startBlock, endBlock = blockRange(startBlock)
			continue
		}

		// Oldest block matches, reorg happened in this range
		blockNumber, blockHash, err := r.searchInRange(storedBlocks)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to search reorg block in range: %w", err)
		}

		if err := r.queries.UpdateBlocksCanonicalityInRange(r.ctx, blockNumber); err != nil {
			return 0, nil, fmt.Errorf("failed to update block range canonicality: %w", err)
		}

		return blockNumber, blockHash, nil
	}
}

func (r *ReorgHandler) searchInRange(blocks []queries.GetBlocksInRangeRow) (uint64, []byte, error) {
	left, right := 0, len(blocks)-1
	for left <= right {
		mid := (left + right) / 2
		storedBlock := blocks[mid]

		chainBlock, err := r.client.BlockByNumber(
			r.ctx,
			big.NewInt(int64(storedBlock.BlockNumber)),
		)
		if err != nil {
			return 0, nil, fmt.Errorf("%w %d: %v", ErrGetBlock, storedBlock.BlockNumber, err)
		}

		if bytes.Equal(storedBlock.BlockHash, chainBlock.Hash().Bytes()) {
			// Found a matching block, check if next block differs to confirm reorg point
			if mid < len(blocks)-1 {
				nextBlock := blocks[mid+1]
				nextChainBlock, err := r.client.BlockByNumber(
					r.ctx,
					big.NewInt(int64(nextBlock.BlockNumber)),
				)
				if err != nil {
					return 0, nil, fmt.Errorf("%w %d: %v", ErrGetBlock, nextBlock.BlockNumber, err)
				}

				if !bytes.Equal(nextBlock.BlockHash, nextChainBlock.Hash().Bytes()) {
					return storedBlock.BlockNumber, chainBlock.Hash().Bytes(), nil
				}
			}

			// If next block doesn't differ, search upper half
			left = mid + 1
		} else {
			// If chainBlock differs, search lower half
			right = mid - 1
		}
	}

	// This should never happen. If it happens, return the first block in the range.
	block := blocks[0]
	return block.BlockNumber, block.BlockHash, nil
}

func blockRange(from uint64) (startBlock uint64, endBlock uint64) {
	endBlock = from

	if endBlock >= BLOCK_RANGE_SIZE {
		startBlock = endBlock - BLOCK_RANGE_SIZE
	}

	return startBlock, endBlock
}
