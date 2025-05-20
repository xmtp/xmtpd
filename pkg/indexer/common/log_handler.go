package common

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	re "github.com/xmtp/xmtpd/pkg/errors"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"go.uber.org/zap"
)

/*
IndexLogs will run until the eventChannel is closed, passing each event to the logStorer.

If an event fails to be stored, and the error is retryable, it will sleep for 100ms and try again.

The only non-retriable errors should be things like malformed events or failed validations.
*/
func IndexLogs(
	ctx context.Context,
	client blockchain.ChainClient,
	eventChannel <-chan types.Log,
	reorgChannel chan<- uint64,
	contract IContract,
) {
	// L3 Orbit works with Arbitrum Elastic Block Time, which under maximum load produces a block every 0.25s.
	// With a maximum throughput of 7M gas per second and a median transaction size of roughly 200k gas,
	// checking for a reorg every 60 blocks (15 seconds) means that, theoretically, a maximum of 495 messages could be affected.
	const reorgCheckInterval = 60

	var (
		storedBlockNumber uint64
		storedBlockHash   []byte
		lastBlockSeen     uint64
		reorgCheckAt      uint64
		reorgDetectedAt   uint64
		reorgBeginsAt     uint64
		reorgFinishesAt   uint64
		reorgInProgress   bool
	)

	// We don't need to listen for the ctx.Done() here, since the eventChannel will be closed when the parent context is canceled
	for event := range eventChannel {
		now := time.Now()
		// 1.1 Handle active reorg state first
		if reorgDetectedAt > 0 {
			// Under a reorg, future events are no-op
			if event.BlockNumber >= reorgDetectedAt {
				contract.Logger().Debug("discarding future event due to reorg",
					zap.Uint64("eventBlockNumber", event.BlockNumber),
					zap.Uint64("reorgBlockNumber", reorgBeginsAt))
				continue
			}
			contract.Logger().Info("starting processing reorg",
				zap.Uint64("eventBlockNumber", event.BlockNumber),
				zap.Uint64("reorgBlockNumber", reorgBeginsAt))

			// When all future events have been discarded, it means we've reached the reorg point
			storedBlockNumber, storedBlockHash = contract.GetLatestBlock()
			lastBlockSeen = event.BlockNumber
			reorgDetectedAt = 0
			reorgInProgress = true
		}

		// 1.2 Handle deactivation of reorg state
		if reorgInProgress && event.BlockNumber > reorgFinishesAt {
			contract.Logger().Info("finished processing reorg",
				zap.Uint64("eventBlockNumber", event.BlockNumber),
				zap.Uint64("reorgFinishesAt", reorgFinishesAt))
			reorgInProgress = false
		}

		// 2. Get the latest block from tracker once per block
		if lastBlockSeen > 0 && lastBlockSeen != event.BlockNumber {
			storedBlockNumber, storedBlockHash = contract.GetLatestBlock()
		}
		lastBlockSeen = event.BlockNumber

		// 3. Check for reorgs, when:
		// - There are no reorgs in progress
		// - There's a stored block
		// - The event block number is greater than the stored block number
		// - The check interval has passed
		skipReorgHandling := false
		if !reorgInProgress &&
			storedBlockNumber > 0 &&
			event.BlockNumber > storedBlockNumber &&
			event.BlockNumber >= reorgCheckAt+reorgCheckInterval {
			onchainBlock, err := client.BlockByNumber(ctx, big.NewInt(int64(storedBlockNumber)))
			if err != nil {
				contract.Logger().
					Warn("error querying block from the blockchain, proceeding with event processing",
						zap.Uint64("blockNumber", storedBlockNumber),
						zap.Error(err),
					)
				skipReorgHandling = true
			}

			if !skipReorgHandling {
				reorgCheckAt = event.BlockNumber
				contract.Logger().Debug("blockchain reorg periodic check",
					zap.Uint64("blockNumber", reorgCheckAt),
				)

				if storedBlockHash != nil &&
					!bytes.Equal(storedBlockHash, onchainBlock.Hash().Bytes()) {
					contract.Logger().Warn("blockchain reorg detected",
						zap.Uint64("storedBlockNumber", storedBlockNumber),
						zap.String("storedBlockHash", hex.EncodeToString(storedBlockHash)),
						zap.String("onchainBlockHash", onchainBlock.Hash().String()),
					)

					reorgBlockNumber, reorgBlockHash, err := contract.FindReorgPoint(
						storedBlockNumber,
					)
					if err != nil && !errors.Is(err, ErrNoBlocksFound) {
						contract.Logger().Error("reorg point not found", zap.Error(err))
						continue
					}

					reorgDetectedAt = storedBlockNumber
					reorgBeginsAt = reorgBlockNumber
					reorgFinishesAt = storedBlockNumber

					if trackerErr := contract.UpdateLatestBlock(ctx, reorgBlockNumber, reorgBlockHash); trackerErr != nil {
						contract.Logger().
							Error("error updating block tracker", zap.Error(trackerErr))
					}

					select {
					case reorgChannel <- reorgBlockNumber:
					default:
						contract.Logger().Warn("reorg signal dropped, channel not ready")
					}

					continue
				}
			}
		}

		err := retry(
			ctx,
			contract.Logger(),
			100*time.Millisecond,
			contract.Address().Hex(),
			func() re.RetryableError {
				return contract.StoreLog(ctx, event)
			},
		)
		if err != nil {
			continue
		}

		contract.Logger().Info("Stored log", zap.Uint64("blockNumber", event.BlockNumber))
		if trackerErr := contract.UpdateLatestBlock(ctx, event.BlockNumber, event.BlockHash.Bytes()); trackerErr != nil {
			contract.Logger().Error("error updating block tracker", zap.Error(trackerErr))
		}
		metrics.EmitIndexerLogProcessingTime(time.Since(now))
	}

	contract.Logger().Debug("exit log handler")
}

func retry(
	ctx context.Context,
	logger *zap.Logger,
	sleep time.Duration,
	address string,
	fn func() re.RetryableError,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := fn(); err != nil {
				logger.Error("error storing log", zap.Error(err))
				if err.ShouldRetry() {
					metrics.EmitIndexerRetryableStorageError(address)

					select {
					case <-ctx.Done():
						return ctx.Err()
					case <-time.After(sleep):
						continue
					}
				}
				return err
			}
			return nil
		}
	}
}
