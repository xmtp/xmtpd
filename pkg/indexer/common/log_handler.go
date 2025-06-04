package common

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	re "github.com/xmtp/xmtpd/pkg/errors"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"go.uber.org/zap"
)

/*
- IndexLogs will run until the eventChannel is closed, passing each event to the logStorer.
- If an event fails to be stored, and the error is retryable, it will sleep for 100ms and try again.
- The only non-retriable errors should be things like malformed events or failed validations.
*/
func IndexLogs(
	ctx context.Context,
	eventChannel <-chan types.Log,
	contract IContract,
) {
	for {
		select {
		case <-ctx.Done():
			contract.Logger().Debug("IndexLogs context cancelled, exiting log handler")
			return

		case event, open := <-eventChannel:
			if !open {
				contract.Logger().Debug("IndexLogs event channel closed, exiting log handler")
				return
			}

			// TODO: Handle reorged event in future PR.
			// This should be handled by the IReorgHandler, and have a different implementaton per contract.
			// Backfilled logs always have Removed = false. Only subscription logs can be reorged.
			if event.Removed {
				if err := contract.HandleLog(ctx, event); err != nil {
					contract.Logger().
						Error("reorg handling failed",
							zap.Error(err),
							zap.Uint64("blockNumber", event.BlockNumber),
							zap.String("blockHash", event.BlockHash.Hex()),
						)
				}
			}

			now := time.Now()

			err := retry(
				ctx,
				contract.Logger(),
				100*time.Millisecond,
				contract.Address().Hex(),
				event.BlockNumber,
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
	}
}

func retry(
	ctx context.Context,
	logger *zap.Logger,
	sleep time.Duration,
	address string,
	blockNumber uint64,
	fn func() re.RetryableError,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			if err := fn(); err != nil {
				logger.Error("error storing log",
					zap.Uint64("blockNumber", blockNumber),
					zap.Error(err),
				)

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
