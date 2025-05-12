package indexerpoc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	defaultWaitTime = 1 * time.Second
)

// task represents an indexing task for a specific contract.
type task struct {
	ctx         context.Context
	state       *taskState
	contract    Contract
	filter      *Filter
	src         Source  // Source represents the blockchain data provider.
	db          Storage // Storage represents the storage for the task: postgres, etc.
	batchSize   uint64
	concurrency int
}

// taskState tracks the current state of indexing for a contract.
// TODO: This task doesn't need that many fields.
type taskState struct {
	ContractName string
	NetworkName  string
	ChainID      int64
	BlockNumber  uint64
	BlockHash    common.Hash
}

// Task Options: optional.
type TaskOption func(*task)

func WithConcurrency(n int) TaskOption {
	return func(t *task) {
		if n > 0 {
			t.concurrency = n
		}
	}
}

// getOrCreateTask creates a new indexing task for a contract.
func getOrCreateTask(
	ctx context.Context,
	src Source,
	contract Contract,
	db Storage,
	batchSize uint64,
	opts ...TaskOption,
) (*task, error) {
	if contract.GetChainID() == 0 {
		return nil, fmt.Errorf("contract %s must have a valid chainID", contract.GetName())
	}

	addresses := []string{contract.GetAddress()}
	topics := [][]string{}
	if len(contract.GetTopics()) > 0 {
		topics = append(topics, contract.GetTopics())
	}

	filter := NewFilter(addresses, topics)

	task := &task{
		ctx:         ctx,
		src:         src,
		contract:    contract,
		filter:      filter,
		db:          db,
		batchSize:   batchSize,
		concurrency: 1,
	}

	// Apply options.
	for _, opt := range opts {
		opt(task)
	}

	// Load existing state or create new state.
	// TODO: This would be replaced with a block/task tracker.
	state, err := db.GetState(ctx, contract.GetName(), fmt.Sprintf("%d", contract.GetChainID()))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If no state exists, start from contract-defined block.
			// If state exist we must force the starting block to prevent missing any data.
			initialBlockNum := contract.GetStartBlock()

			if initialBlockNum > 0 {
				initialBlockNum--
			}

			// Get hash for initial block if not starting from genesis.
			var initialHash common.Hash
			if initialBlockNum > 0 {
				initialHash, err = src.GetBlockHash(ctx, initialBlockNum)
				if err != nil {
					return nil, fmt.Errorf("getting hash for initial block on chain %d: %w",
						src.GetChainID(), err)
				}
			}

			task.state = &taskState{
				ContractName: contract.GetName(),
				NetworkName:  src.GetNetworkName(),
				ChainID:      src.GetChainID(),
				BlockNumber:  initialBlockNum,
				BlockHash:    initialHash,
			}
		} else {
			return nil, fmt.Errorf("getting state for %s on chain %d: %w",
				contract.GetName(), src.GetChainID(), err)
		}
	} else {
		task.state = state
	}

	if state.ChainID != src.GetChainID() {
		return nil, fmt.Errorf("chain ID mismatch for %s: stored %d, current %d",
			contract.GetName(), state.ChainID, src.GetChainID())
	}

	return task, nil
}

// nextState performs one indexing run, handling blockchain data retrieval and processing.
// It returns an error if there is nothing new to index, or if there is a reorg.
func (t *task) nextState() error {
	ctx := t.ctx

	// Get current latest block from blockchain.
	latestBlockNum, err := t.src.GetLatestBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("getting latest block on %s: %w", t.state.NetworkName, err)
	}

	// Check if we already have the latest block.
	if t.state.BlockNumber == latestBlockNum {
		return ErrNothingNew
	}

	// Calculate how many blocks to retrieve.
	delta := uint64(1)
	if t.batchSize > 1 && latestBlockNum > t.state.BlockNumber {
		delta = min(latestBlockNum-t.state.BlockNumber, t.batchSize)
	}

	// Determine range of blocks to process
	startBlock := t.state.BlockNumber + 1
	endBlock := startBlock + delta - 1

	// TODO: This tx would update the block tracker to the previous block.
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			slog.ErrorContext(ctx, "error rolling back transaction",
				"contract", t.contract.GetName(),
				"network", t.state.NetworkName,
				"error", err,
			)
		}
	}()

	// Retrieve logs.
	logs, isReorg, err := t.nextBatch(ctx, startBlock, endBlock)
	if err != nil {
		return fmt.Errorf("retrieving logs and checking reorgs on %s: %w", t.state.NetworkName, err)
	}

	if isReorg {
		slog.WarnContext(ctx, "detected chain reorganization",
			"contract", t.contract.GetName(),
			"network", t.state.NetworkName,
			"block", t.state.BlockNumber,
		)

		// Reorg handling approach:
		// 1. Detect reorg at block X (current block)
		// 2. Call the contract-specific reorg handler for block X
		// 3. Delete data from block X in block tracker
		// 4. Move back to block X-1 and continue from there
		// 5. If X-1 was also reorged, it will be detected in the next cycle
		//
		// This iterative approach handles reorgs of any depth by backing up
		// one block at a time until we reach a stable block.

		// Handle the reorg with contract's processor.
		// Let the contract's reorg processor handle block X
		reorgErr := t.contract.HandleReorg(ctx, t.state.BlockNumber)
		if reorgErr != nil {
			return fmt.Errorf("handling reorg: %w", reorgErr)
		}

		// Delete stored data from the reorg point
		reorgErr = t.db.DeleteFromBlock(
			ctx,
			t.contract.GetName(),
			fmt.Sprintf("%d", t.contract.GetChainID()),
			t.state.BlockNumber,
		)
		if reorgErr != nil {
			return fmt.Errorf("deleting logs from block: %w", reorgErr)
		}

		// Move back one block
		reorgBlock := t.state.BlockNumber - 1

		// Get hash for the previous block
		reorgHash, err := t.src.GetBlockHash(ctx, reorgBlock)
		if err != nil {
			return fmt.Errorf(
				"getting hash for block before reorg on %s: %w",
				t.state.NetworkName,
				err,
			)
		}

		// Update our state to block X-1
		t.state.BlockNumber = reorgBlock
		t.state.BlockHash = reorgHash

		// Save the updated state
		if err := t.db.SaveState(ctx, t.state); err != nil {
			return fmt.Errorf("saving state after reorg: %w", err)
		}

		// Commit the transaction with our reorg handling
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("committing transaction after reorg: %w", err)
		}

		// Return ErrReorg to signal we need to process again
		// The next cycle will start from block X-1 and will detect if it was also reorged
		return ErrReorg
	}

	// Process logs
	if len(logs) > 0 {
		if err := t.contract.ProcessLogs(ctx, logs); err != nil {
			return fmt.Errorf("processing logs on %s: %w", t.state.NetworkName, err)
		}
	}

	// Get the hash of the last block we processed
	lastBlockHash, err := t.src.GetBlockHash(ctx, endBlock)
	if err != nil {
		return fmt.Errorf("getting hash for last block on %s: %w", t.state.NetworkName, err)
	}

	// Update state with the latest processed block
	t.state.BlockNumber = endBlock
	t.state.BlockHash = lastBlockHash

	// Save updated state
	if err := t.db.SaveState(ctx, t.state); err != nil {
		return fmt.Errorf("saving state: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	slog.InfoContext(ctx, "indexed blocks",
		"contract", t.contract.GetName(),
		"network", t.state.NetworkName,
		"start", startBlock,
		"end", endBlock,
		"logs", len(logs),
	)

	return nil
}

// nextBatch retrieves a batch of blocks and logs while checking for reorgs.
func (t *task) nextBatch(
	ctx context.Context,
	startBlock, endBlock uint64,
) ([]types.Log, bool, error) {
	// Special case - if we have no previous state, no reorg check needed
	if t.state.BlockNumber == 0 || t.state.BlockHash == (common.Hash{}) {
		logs, err := t.src.GetLogs(ctx, startBlock, endBlock, t.filter)
		return logs, false, err
	}

	// We need to check for reorgs first by verifying the parent hash of the first block
	firstBlock, err := t.src.GetBlockByNumber(ctx, startBlock)
	if err != nil {
		return nil, false, fmt.Errorf("getting first block %d for reorg check on %s: %w",
			startBlock, t.state.NetworkName, err)
	}

	// Check if the parent hash matches our stored hash.
	if firstBlock.ParentHash() != t.state.BlockHash {
		// Reorg detected!
		return nil, true, nil
	}

	// No reorg, retrieve logs
	// For small ranges or if concurrency is 1, use a simple approach
	if t.concurrency <= 1 || endBlock-startBlock < 10 {
		logs, err := t.src.GetLogs(ctx, startBlock, endBlock, t.filter)
		if err != nil {
			return nil, false, fmt.Errorf("getting logs on %s: %w", t.state.NetworkName, err)
		}
		return logs, false, nil
	}

	// For larger ranges, use concurrent processing
	return t.getLogs(ctx, startBlock, endBlock)
}

// getLogs retrieves logs in parallel for better performance with large batches.
// This is likely going to be used only when backfilling.
func (t *task) getLogs(
	ctx context.Context,
	startBlock, endBlock uint64,
) ([]types.Log, bool, error) {
	var (
		allLogs []types.Log
		mu      sync.Mutex
		wg      sync.WaitGroup
		errCh   = make(chan error, t.concurrency)
	)

	totalBlocks := endBlock - startBlock + 1
	blocksPerWorker := totalBlocks / uint64(t.concurrency)
	if blocksPerWorker == 0 {
		blocksPerWorker = 1
	}

	for i := 0; i < t.concurrency; i++ {
		wg.Add(1)

		workerStart := startBlock + uint64(i)*blocksPerWorker
		workerEnd := workerStart + blocksPerWorker - 1

		// Adjust the last worker to include any remainder
		if i == t.concurrency-1 {
			workerEnd = endBlock
		}

		// Skip if we've already gone past the end
		if workerStart > endBlock {
			wg.Done()
			continue
		}

		go func(start, end uint64) {
			defer wg.Done()

			logs, err := t.src.GetLogs(ctx, start, end, t.filter)
			if err != nil {
				errCh <- fmt.Errorf("worker getting logs %d-%d on %s: %w",
					start, end, t.state.NetworkName, err)
				return
			}

			mu.Lock()
			allLogs = append(allLogs, logs...)
			mu.Unlock()
		}(workerStart, workerEnd)
	}

	// Wait for all workers
	wg.Wait()
	close(errCh)

	// Check for errors
	if len(errCh) > 0 {
		return nil, false, <-errCh // Return the first error
	}

	// Sort logs by block number and index for consistency
	sort.SliceStable(allLogs, func(i, j int) bool {
		if allLogs[i].BlockNumber != allLogs[j].BlockNumber {
			return allLogs[i].BlockNumber < allLogs[j].BlockNumber
		}
		return allLogs[i].Index < allLogs[j].Index
	})

	return allLogs, false, nil
}

// run starts the indexer and continuously processes new blocks.
func (t *task) run() {
	reorgCount := 0
	maxReorgs := 1000 // Add a safety limit similar to shovel

	for {
		select {
		case <-t.ctx.Done():
			return
		default:
			err := t.nextState()
			switch {
			case err == nil:
				reorgCount = 0
				continue
			case err == ErrNothingNew:
				reorgCount = 0
				time.Sleep(defaultWaitTime)
			case err == ErrReorg:
				reorgCount++
				if reorgCount > maxReorgs {
					slog.ErrorContext(
						t.ctx,
						"too many consecutive reorgs detected - possible deep chain reorganization",
						"contract",
						t.contract.GetName(),
						"network",
						t.state.NetworkName,
						"count",
						reorgCount,
					)

					// Back off to avoid hammering the node and allow the network to stabilize
					time.Sleep(defaultWaitTime)
					reorgCount = 0
				} else {
					// Continue immediately to process the next block
					// (which will be the block before the reorg)
					continue
				}
			default:
				// Other error, log and wait before retrying
				slog.ErrorContext(t.ctx, "error processing blocks",
					"contract", t.contract.GetName(),
					"network", t.state.NetworkName,
					"error", err,
				)
				time.Sleep(defaultWaitTime)
			}
		}
	}
}
