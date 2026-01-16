package migrator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/db/types"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/utils"
	re "github.com/xmtp/xmtpd/pkg/utils/retryerrors"
	"go.uber.org/zap"
)

const (
	maxChainMessageSize  = 200 * 1024 // 200KB
	maxChainBatchSize    = 100 * 1024 // 100KB - batch threshold for identity updates
	maxDatabaseBatchSize = 1000       // 1000 messages - batch threshold for database inserts
)

// insertOriginatorEnvelopeDatabaseBatch inserts a batch of originator envelopes into the database.
func (w *Worker) insertOriginatorEnvelopeDatabaseBatch(
	ctx context.Context,
	logger *zap.Logger,
	batch *types.GatewayEnvelopeBatch,
) re.RetryableError {
	if batch == nil {
		return re.NewNonRecoverableError("", errors.New("batch is nil"))
	}

	err := db.RunInTx(
		ctx,
		w.writer.Write(),
		nil,
		func(ctx context.Context, querier *queries.Queries) error {
			_, err := db.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
				ctx,
				w.writer.Write(),
				batch,
			)
			if err != nil {
				logger.Error("insert originator envelope batch failed", zap.Error(err))
				return re.NewRecoverableError("insert originator envelope batch failed", err)
			}

			err = querier.UpdateMigrationProgress(ctx, queries.UpdateMigrationProgressParams{
				LastMigratedID: batch.LastSequenceID(),
				SourceTable:    w.tableName,
			})
			if err != nil {
				logger.Error("update migration progress failed", zap.Error(err))
				return re.NewRecoverableError("update migration progress failed", err)
			}

			return nil
		})
	if err != nil {
		var retryableError re.RetryableError
		if errors.As(err, &retryableError) {
			return retryableError
		}

		return re.NewRecoverableError("database error", err)
	}

	metrics.EmitMigratorTargetLastSequenceID(w.tableName, batch.LastSequenceID())

	return nil
}

/* Commit messages bootstrap functions. */

// flushBroadcasterBatch handles the retry logic for publishing broadcaster batches.
//
//	It tries to bootstrap the batch in a single transaction.
//	If batch fails, it retries the individual messages.
//	On individual retry failure, it inserts a record into the migration dead letter box.
func (w *Worker) flushBroadcasterBatch(
	ctx context.Context,
	logger *zap.Logger,
	batch *BroadcasterBatch,
) error {
	if batch.Len() == 0 {
		return nil
	}

	err := retry(
		ctx,
		50*time.Millisecond,
		w.tableName,
		destinationBlockchain,
		func() re.RetryableError {
			return w.bootstrapBatch(ctx, batch)
		},
	)

	if err == nil || errors.Is(err, ErrMigrationProgressUpdateFailed) {
		return nil
	}

	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	logger.Warn(
		"batch failed, retrying messages individually",
		zap.Error(err),
		zap.Int("count", batch.Len()),
	)

	for message := range batch.All() {
		err := retry(
			ctx,
			50*time.Millisecond,
			w.tableName,
			destinationBlockchain,
			func() re.RetryableError {
				return w.bootstrapBatch(ctx, &BroadcasterBatch{
					identifiers: [][]byte{message.Identifier},
					payloads:    [][]byte{message.Payload},
					sequenceIDs: []uint64{message.SequenceID},
				})
			},
		)
		if err == nil {
			continue
		}

		logger.Warn(
			"individual retry failed, inserting into dead letter box",
			utils.IdentifierField(utils.HexEncode(message.Identifier[:])),
			utils.SequenceIDField(int64(message.SequenceID)),
			zap.Error(err),
		)

		err = insertMigrationDeadLetterBox(
			ctx,
			w.writer.Write(),
			w.tableName,
			int64(message.SequenceID),
			message.Payload[:],
			FailureBlockchainUndetermined,
		)
		if err != nil {
			logger.Error(
				"failed to insert migration dead letter box",
				utils.IdentifierField(utils.HexEncode(message.Identifier[:])),
				utils.SequenceIDField(int64(message.SequenceID)),
				zap.Error(err),
			)

			return err
		}

	}

	return nil
}

// bootstrapBatch bootstraps a batch of messages to the blockchain.
// On failure, it inserts a record into the migration dead letter box.
func (w *Worker) bootstrapBatch(
	ctx context.Context,
	batch *BroadcasterBatch,
) re.RetryableError {
	// Should never happen.
	if len(batch.identifiers) != len(batch.payloads) ||
		len(batch.identifiers) != len(batch.sequenceIDs) {
		return re.NewNonRecoverableError("array mismatch", errors.New("array mismatch"))
	}

	var publishFn func(context.Context, [][]byte, [][]byte, []uint64) error

	switch w.tableName {
	case commitMessagesTableName:
		publishFn = func(ctx context.Context, ids [][]byte, payloads [][]byte, seqs []uint64) error {
			groupIDs := make([][16]byte, len(ids))
			for i, id := range ids {
				copy(groupIDs[i][:], id)
			}
			_, err := w.blockchainPublisher.BootstrapGroupMessages(ctx, groupIDs, payloads, seqs)
			return err
		}

	case inboxLogTableName:
		publishFn = func(ctx context.Context, ids [][]byte, payloads [][]byte, seqs []uint64) error {
			inboxIDs := make([][32]byte, len(ids))
			for i, id := range ids {
				copy(inboxIDs[i][:], id)
			}
			_, err := w.blockchainPublisher.BootstrapIdentityUpdates(ctx, inboxIDs, payloads, seqs)
			return err
		}

	default:
		return re.NewNonRecoverableError("invalid table name", errors.New("invalid table name"))
	}

	err := publishFn(ctx, batch.identifiers, batch.payloads, batch.sequenceIDs)
	if err != nil {
		// If the broadcaster reverts with NotPaused() or NotPayloadBootstrapper(), wait and try again until resolved.
		if strings.Contains(err.Error(), "NotPaused()") ||
			strings.Contains(err.Error(), "NotPayloadBootstrapper()") {
			return re.NewRecoverableError(err.Error(), err)
		}

		// Transient errors - recoverable (timeout, network issues).
		if strings.Contains(err.Error(), "timed out") ||
			errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) {
			return re.NewRecoverableError("transient blockchain error", err)
		}

		// Unknown errors - default to recoverable.
		return re.NewRecoverableError(
			"error publishing batch",
			err,
		)
	}

	err = w.writer.WriteQuery().UpdateMigrationProgress(ctx, queries.UpdateMigrationProgressParams{
		LastMigratedID: int64(batch.LastSequenceID()),
		SourceTable:    w.tableName,
	})
	if err != nil {
		// If we reached this point, the message has been published and the log emitted.
		// Therefore, we can return a non-recoverable error to ensure the message is not retried.
		return re.NewNonRecoverableError(
			"update migration progress failed",
			ErrMigrationProgressUpdateFailed,
		)
	}

	metrics.EmitMigratorTargetLastSequenceID(w.tableName, int64(batch.LastSequenceID()))

	return nil
}

// prepareClientEnvelope extracts a client envelope from an originator envelope.
// On failure, it inserts a record into the migration dead letter box.
func (w *Worker) prepareClientEnvelope(
	ctx context.Context,
	logger *zap.Logger,
	env *envelopes.OriginatorEnvelope,
	tableName string,
) (clientEnvelopeBytes []byte, identifier []byte, sequenceID uint64, err error) {
	sequenceID = env.OriginatorSequenceID()
	identifier = env.TargetTopic().Identifier()

	clientEnvelopeBytes, err = env.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope.Bytes()
	if err != nil {
		logger.Error("failed to get payer envelope bytes", zap.Error(err))
		return nil, nil, 0, fmt.Errorf("failed to get payer envelope bytes: %w", err)
	}

	// 8 bytes for the sequence ID.
	totalSize := len(clientEnvelopeBytes) + len(identifier) + 8

	if totalSize > maxChainMessageSize {
		err := insertMigrationDeadLetterBox(
			ctx,
			w.writer.Write(),
			tableName,
			int64(sequenceID),
			identifier[:],
			FailureOversizedChainMessage,
		)
		if err != nil {
			// Ensure dead letter box is inserted before returning an error.
			return nil, nil, 0, fmt.Errorf("insert migration dead letter box failed: %w", err)
		}

		logger.Warn(
			"oversized blockchain payload, skipped and added to dead letter box",
			zap.String("identifier", utils.HexEncode(identifier[:])),
			utils.SequenceIDField(int64(sequenceID)),
			zap.Int("size", totalSize),
		)

		// Return a non-recoverable error to ensure the message is not retried.
		return nil, nil, 0, fmt.Errorf(
			"oversized blockchain payload, skipped and added to dead letter box: %w",
			ErrDeadLetterBox,
		)
	}

	return clientEnvelopeBytes, identifier, sequenceID, nil
}

/* Database helper functions. */

// retry implements the retry logic for insert (db or chain) operations.
func retry(
	ctx context.Context,
	sleep time.Duration,
	tableName string,
	destination string,
	fn func() re.RetryableError,
) error {
	attempts := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			if err := fn(); err != nil {
				if err.ShouldRetry() {
					attempts++
					metrics.EmitMigratorWriterRetryAttempts(tableName, destination, attempts)
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

// insertMigrationDeadLetterBox inserts a record into the migration dead letter box.
func insertMigrationDeadLetterBox(
	ctx context.Context,
	database *sql.DB,
	sourceTable string,
	sequenceID int64,
	payload []byte,
	reason FailureReason,
) error {
	return db.RunInTx(
		ctx,
		database,
		nil,
		func(ctx context.Context, querier *queries.Queries) error {
			_, err := querier.InsertMigrationDeadLetterBox(
				ctx,
				queries.InsertMigrationDeadLetterBoxParams{
					SourceTable: sourceTable,
					SequenceID:  sequenceID,
					Payload:     payload,
					Reason:      reason.String(),
					Retryable:   reason.ShouldRetry(),
				},
			)
			if err != nil {
				return fmt.Errorf("insert migration dead letter box failed: %w", err)
			}

			// Skip this record by advancing migration progress.
			err = querier.UpdateMigrationProgress(ctx, queries.UpdateMigrationProgressParams{
				LastMigratedID: sequenceID,
				SourceTable:    sourceTable,
			})
			if err != nil {
				return fmt.Errorf("update migration progress failed: %w", err)
			}

			return nil
		},
	)
}
