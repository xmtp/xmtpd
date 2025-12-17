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
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/utils"
	re "github.com/xmtp/xmtpd/pkg/utils/retryerrors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	maxChainMessageSize = 200 * 1024 // 200KB
	maxBatchSize        = 100 * 1024 // 100KB - batch threshold for identity updates
)

// insertOriginatorEnvelopeDatabase inserts an originator envelope into the database.
func (w *Worker) insertOriginatorEnvelopeDatabase(
	ctx context.Context,
	env *envelopes.OriginatorEnvelope,
) re.RetryableError {
	if env == nil {
		return re.NewNonRecoverableError("", errors.New("envelope is nil"))
	}

	payerAddress, err := env.UnsignedOriginatorEnvelope.PayerEnvelope.RecoverSigner()
	if err != nil {
		w.logger.Error("recover payer address failed", zap.Error(err))
		return re.NewNonRecoverableError("recover payer address failed", err)
	}

	querier := queries.New(w.writer)
	payerID, err := querier.FindOrCreatePayer(ctx, payerAddress.Hex())
	if err != nil {
		w.logger.Error("find or create payer failed", zap.Error(err))
		return re.NewRecoverableError("find or create payer failed", err)
	}

	originatorEnvelopeBytes, err := proto.Marshal(env.Proto())
	if err != nil {
		w.logger.Error("marshall originator envelope failed", zap.Error(err))
		return re.NewNonRecoverableError("marshall originator envelope failed", err)
	}

	err = db.RunInTx(
		ctx,
		w.writer,
		nil,
		func(ctx context.Context, querier *queries.Queries) error {
			_, err := db.InsertGatewayEnvelopeWithChecksTransactional(
				ctx,
				querier,
				queries.InsertGatewayEnvelopeParams{
					OriginatorNodeID:     int32(env.OriginatorNodeID()),
					OriginatorSequenceID: int64(env.OriginatorSequenceID()),
					Topic:                env.TargetTopic().Bytes(),
					OriginatorEnvelope:   originatorEnvelopeBytes,
					PayerID:              db.NullInt32(payerID),
					GatewayTime:          env.OriginatorTime(),
					Expiry: int64(
						env.UnsignedOriginatorEnvelope.Proto().GetExpiryUnixtime(),
					),
				},
			)
			if err != nil {
				w.logger.Error("insert originator envelope failed", zap.Error(err))
				return re.NewRecoverableError("insert originator envelope failed", err)
			}

			err = querier.IncrementUnsettledUsage(ctx, queries.IncrementUnsettledUsageParams{
				PayerID:           payerID,
				OriginatorID:      int32(env.OriginatorNodeID()),
				MinutesSinceEpoch: utils.MinutesSinceEpoch(env.OriginatorTime()),
				SpendPicodollars: int64(env.UnsignedOriginatorEnvelope.BaseFee()) +
					int64(env.UnsignedOriginatorEnvelope.CongestionFee()),
				SequenceID: int64(env.OriginatorSequenceID()),
			})
			if err != nil {
				w.logger.Error("increment unsettled usage failed", zap.Error(err))
				return re.NewRecoverableError("increment unsettled usage failed", err)
			}

			err = querier.UpdateMigrationProgress(ctx, queries.UpdateMigrationProgressParams{
				LastMigratedID: int64(env.OriginatorSequenceID()),
				SourceTable:    w.tableName,
			})
			if err != nil {
				w.logger.Error("update migration progress failed", zap.Error(err))
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

	return nil
}

// insertOriginatorEnvelopeBlockchainUnary is a generic function that inserts an originator envelope into the blockchain.
// On failure, it inserts a record into the migration dead letter box.
func (w *Worker) insertOriginatorEnvelopeBlockchainUnary(
	ctx context.Context,
	env *envelopes.OriginatorEnvelope,
) re.RetryableError {
	var (
		identifier   = env.TargetTopic().Identifier()
		sequenceID   = env.OriginatorSequenceID()
		originatorID = env.OriginatorNodeID()
	)

	tableName, ok := originatorIDToTableName[originatorID]
	if !ok {
		return re.NewNonRecoverableError("", errors.New("invalid originator id"))
	}

	clientEnvelopeBytes, err := env.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope.Bytes()
	if err != nil {
		w.logger.Error("failed to get payer envelope bytes", zap.Error(err))
		return re.NewNonRecoverableError("failed to get payer envelope bytes", err)
	}

	totalSize := len(clientEnvelopeBytes) + len(identifier) + 8

	querier := queries.New(w.writer)

	switch originatorID {
	case CommitMessageOriginatorID:
		groupID, err := utils.ParseGroupID(identifier)
		if err != nil {
			return re.NewNonRecoverableError("error converting identifier to group ID", err)
		}

		if totalSize > maxChainMessageSize {
			err := insertMigrationDeadLetterBox(
				ctx,
				w.writer,
				tableName,
				int64(sequenceID),
				groupID[:],
				FailureOversizedChainMessage,
			)
			if err != nil {
				// Ensure dead letter box is inserted before returning an error.
				return re.NewRecoverableError("insert migration dead letter box failed", err)
			}

			w.logger.Warn(
				"oversized commit message, skipped and added to dead letter box",
				utils.GroupIDField(utils.HexEncode(groupID[:])),
				utils.SequenceIDField(int64(sequenceID)),
				zap.Int("size", totalSize),
			)

			// Return a non-recoverable error to ensure the message is not retried.
			return re.NewNonRecoverableError(
				"oversized commit message, skipped and added to dead letter box",
				ErrDeadLetterBox,
			)
		}

		_, err = w.blockchainPublisher.BootstrapGroupMessages(
			ctx,
			[][16]byte{groupID},
			[][]byte{clientEnvelopeBytes},
			[]uint64{sequenceID},
		)
		if err != nil {
			// If the broadcaster reverts with NotPaused() or NotPayloadBootstrapper(), wait and try again until resolved.
			if strings.Contains(err.Error(), "NotPaused()") ||
				strings.Contains(err.Error(), "NotPayloadBootstrapper()") {
				return re.NewRecoverableError(
					err.Error(),
					err,
				)
			}

			if strings.Contains(err.Error(), "InvalidPayloadSize()") {
				errInsert := insertMigrationDeadLetterBox(
					ctx,
					w.writer,
					tableName,
					int64(sequenceID),
					groupID[:],
					FailureOversizedChainMessage,
				)
				if errInsert != nil {
					// Ensure dead letter box is inserted before returning an error.
					return re.NewRecoverableError(
						"invalid payload size, inserted migration dead letter box failed",
						errInsert,
					)
				}

				// Return a non-recoverable error to ensure the message is not retried.
				return re.NewNonRecoverableError("invalid payload size", err)
			}

			// Transient errors - recoverable (timeout, network issues).
			if strings.Contains(err.Error(), "timed out") ||
				errors.Is(err, context.DeadlineExceeded) ||
				errors.Is(err, context.Canceled) {
				return re.NewRecoverableError("transient blockchain error", err)
			}

			// Unknown errors - default to recoverable.
			return re.NewRecoverableError(
				fmt.Sprintf("error publishing group message %d", sequenceID),
				err,
			)
		}

		err = querier.UpdateMigrationProgress(ctx, queries.UpdateMigrationProgressParams{
			LastMigratedID: int64(sequenceID),
			SourceTable:    tableName,
		})
		if err != nil {
			w.logger.Error("update migration progress failed", zap.Error(err))

			// If we reached this point, the message has been published and the log emitted.
			// Therefore, we can return a non-recoverable error to ensure the message is not retried.
			return re.NewNonRecoverableError("update migration progress failed", err)
		}

		w.logger.Debug(
			"published group message",
			utils.GroupIDField(utils.HexEncode(groupID[:])),
			utils.SequenceIDField(int64(sequenceID)),
		)

	case InboxLogOriginatorID:
		inboxID, err := utils.ParseInboxID(identifier)
		if err != nil {
			return re.NewNonRecoverableError("error converting identifier to inbox ID", err)
		}

		w.logger.Debug(
			"publishing identity update",
			utils.InboxIDField(utils.HexEncode(inboxID[:])),
			utils.SequenceIDField(int64(sequenceID)),
		)

		_, err = w.blockchainPublisher.BootstrapIdentityUpdates(
			ctx,
			[][32]byte{inboxID},
			[][]byte{clientEnvelopeBytes},
			[]uint64{sequenceID},
		)
		if err != nil {
			// If the broadcaster reverts with NotPaused() or NotPayloadBootstrapper(), wait and try again until resolved.
			if strings.Contains(err.Error(), "NotPaused()") ||
				strings.Contains(err.Error(), "NotPayloadBootstrapper()") {
				return re.NewRecoverableError(err.Error(), err)
			}

			if strings.Contains(err.Error(), "InvalidPayloadSize()") {
				errInsert := insertMigrationDeadLetterBox(
					ctx,
					w.writer,
					tableName,
					int64(sequenceID),
					inboxID[:],
					FailureOversizedChainMessage,
				)
				if errInsert != nil {
					// Ensure dead letter box is inserted before returning an error.
					return re.NewRecoverableError(
						"invalid payload size, inserted migration dead letter box failed",
						errInsert,
					)
				}

				// Return a non-recoverable error to ensure the message is not retried.
				return re.NewNonRecoverableError("invalid payload size", err)
			}

			// Transient errors - recoverable (timeout, network issues).
			if strings.Contains(err.Error(), "timed out") ||
				errors.Is(err, context.DeadlineExceeded) ||
				errors.Is(err, context.Canceled) {
				return re.NewRecoverableError("transient blockchain error", err)
			}

			// Unknown errors - default to recoverable.
			return re.NewRecoverableError(
				fmt.Sprintf("error publishing identity update %d", sequenceID),
				err,
			)
		}

		err = querier.UpdateMigrationProgress(ctx, queries.UpdateMigrationProgressParams{
			LastMigratedID: int64(sequenceID),
			SourceTable:    tableName,
		})
		if err != nil {
			w.logger.Error("update migration progress failed", zap.Error(err))

			// If we reached this point, the message has been published and the log emitted.
			// Therefore, we can return a non-recoverable error to ensure the message is not retried.
			return re.NewNonRecoverableError("update migration progress failed", err)
		}

		w.logger.Debug(
			"published identity update",
			utils.InboxIDField(utils.HexEncode(inboxID[:])),
			utils.SequenceIDField(int64(sequenceID)),
		)
	}

	return nil
}

/* Identity updates bootstrap functions. */

// flushIdentityUpdatesBatch handles the retry logic for the identity updates batch.
// - It tries to bootstrap the batch in a single transaction.
// - If batch fails, it retries the individual identity updates.
// - On individual retry failure, it inserts a record into the migration dead letter box.
func (w *Worker) flushIdentityUpdatesBatch(
	ctx context.Context,
	logger *zap.Logger,
	batch *IdentityUpdateBatch,
) error {
	if batch.Len() == 0 {
		return nil
	}

	err := retry(
		ctx,
		50*time.Millisecond,
		inboxLogTableName,
		destinationBlockchain,
		func() re.RetryableError {
			return w.bootstrapIdentityUpdates(ctx, batch)
		},
	)

	if err == nil {
		return nil
	}

	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	logger.Warn(
		"batch failed, retrying identity updates individually",
		zap.Error(err),
		zap.Int("count", batch.Len()),
	)

	for identityUpdate := range batch.All() {
		err := retry(
			ctx,
			50*time.Millisecond,
			inboxLogTableName,
			destinationBlockchain,
			func() re.RetryableError {
				return w.bootstrapIdentityUpdates(ctx, &IdentityUpdateBatch{
					inboxIDs:        [][32]byte{identityUpdate.InboxID},
					identityUpdates: [][]byte{identityUpdate.IdentityUpdate},
					sequenceIDs:     []uint64{identityUpdate.SequenceID},
				})
			},
		)
		if err == nil {
			logger.Debug(
				"individual retry succeeded",
				utils.InboxIDField(utils.HexEncode(identityUpdate.InboxID[:])),
				utils.SequenceIDField(int64(identityUpdate.SequenceID)),
			)

			continue
		}

		logger.Error(
			"individual retry failed, inserting into dead letter box",
			utils.InboxIDField(utils.HexEncode(identityUpdate.InboxID[:])),
			utils.SequenceIDField(int64(identityUpdate.SequenceID)),
			zap.Error(err),
		)

		err = insertMigrationDeadLetterBox(
			ctx,
			w.writer,
			inboxLogTableName,
			int64(identityUpdate.SequenceID),
			identityUpdate.InboxID[:],
			FailureBlockchainUndetermined,
		)
		if err != nil {
			logger.Error(
				"failed to insert migration dead letter box",
				utils.InboxIDField(utils.HexEncode(identityUpdate.InboxID[:])),
				utils.SequenceIDField(int64(identityUpdate.SequenceID)),
				zap.Error(err),
			)
		}

	}

	return nil
}

// bootstrapIdentityUpdates bootstraps a batch of identity updates to the blockchain.
// On failure, it inserts a record into the migration dead letter box.
func (w *Worker) bootstrapIdentityUpdates(
	ctx context.Context,
	batch *IdentityUpdateBatch,
) re.RetryableError {
	// Should never happen.
	if len(batch.inboxIDs) != len(batch.identityUpdates) ||
		len(batch.inboxIDs) != len(batch.sequenceIDs) {
		return re.NewNonRecoverableError("array mismatch", errors.New("array mismatch"))
	}

	querier := queries.New(w.writer)

	_, err := w.blockchainPublisher.BootstrapIdentityUpdates(
		ctx,
		batch.inboxIDs,
		batch.identityUpdates,
		batch.sequenceIDs,
	)
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
			"error publishing identity update batch",
			err,
		)
	}

	err = querier.UpdateMigrationProgress(ctx, queries.UpdateMigrationProgressParams{
		LastMigratedID: int64(batch.LastSequenceID()),
		SourceTable:    inboxLogTableName,
	})
	if err != nil {
		w.logger.Error("update migration progress failed", zap.Error(err))

		// If we reached this point, the message has been published and the log emitted.
		// Therefore, we can return a non-recoverable error to ensure the message is not retried.
		return re.NewNonRecoverableError("update migration progress failed", err)
	}

	w.logger.Debug(
		"published identity update batch",
		zap.Int("length", batch.Len()),
		utils.SequenceIDField(int64(batch.LastSequenceID())),
	)

	return nil
}

// prepareClientEnvelope prepares the client envelope for the identity update.
// On failure, it inserts a record into the migration dead letter box.
func (w *Worker) prepareClientEnvelope(
	ctx context.Context,
	env *envelopes.OriginatorEnvelope,
	tableName string,
) (clientEnvelopeBytes []byte, identifier []byte, sequenceID uint64, err error) {
	sequenceID = env.OriginatorSequenceID()
	identifier = env.TargetTopic().Identifier()

	clientEnvelopeBytes, err = env.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope.Bytes()
	if err != nil {
		w.logger.Error("failed to get payer envelope bytes", zap.Error(err))
		return nil, nil, 0, fmt.Errorf("failed to get payer envelope bytes: %w", err)
	}

	// 8 bytes for the sequence ID.
	totalSize := len(clientEnvelopeBytes) + len(identifier) + 8

	if totalSize > maxChainMessageSize {
		err := insertMigrationDeadLetterBox(
			ctx,
			w.writer,
			tableName,
			int64(sequenceID),
			identifier[:],
			FailureOversizedChainMessage,
		)
		if err != nil {
			// Ensure dead letter box is inserted before returning an error.
			return nil, nil, 0, fmt.Errorf("insert migration dead letter box failed: %w", err)
		}

		w.logger.Warn(
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
