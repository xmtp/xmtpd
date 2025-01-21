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
		return r.searchInRange(storedBlocks)
	}
}

func (r *ReorgHandler) searchInRange(blocks []queries.GetBlocksInRangeRow) (uint64, []byte, error) {
	left, right := 0, len(blocks)-1
	for left <= right {
		mid := (left + right) / 2
		block := blocks[mid]

		chainBlock, err := r.client.BlockByNumber(
			r.ctx,
			big.NewInt(int64(block.BlockNumber)),
		)
		if err != nil {
			return 0, nil, fmt.Errorf("%w %d: %v", ErrGetBlock, block.BlockNumber, err)
		}

		if bytes.Equal(block.BlockHash, chainBlock.Hash().Bytes()) {
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
					return block.BlockNumber, chainBlock.Hash().Bytes(), nil
				}
			} else if mid == len(blocks)-1 {
				return block.BlockNumber, chainBlock.Hash().Bytes(), nil
			}

			// If next block doesn't differ, search upper half
			left = mid + 1
		} else {
			// If chainBlock differs, search lower half
			right = mid - 1
		}
	}

	// TODO: This should never happen, start from 0?
	return 0, nil, fmt.Errorf("reorg point not found")
}

func blockRange(from uint64) (startBlock uint64, endBlock uint64) {
	endBlock = from

	if endBlock > BLOCK_RANGE_SIZE {
		startBlock = endBlock - BLOCK_RANGE_SIZE
	}

	return startBlock, endBlock
}
