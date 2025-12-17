package migrator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"
	re "github.com/xmtp/xmtpd/pkg/utils/retryerrors"
	"go.uber.org/zap"
)

type Worker struct {
	// Internals.
	logger *zap.Logger

	// Data management.
	writer              *sql.DB
	wrtrQueries         *queries.Queries
	blockchainPublisher blockchain.IBlockchainPublisher

	// IPC.
	wg         sync.WaitGroup
	recvChan   chan ISourceRecord
	wrtrChan   chan *envelopes.OriginatorEnvelope
	sem        chan struct{}
	inflightMu sync.Mutex
	inflight   map[int64]time.Time

	// Configuration.
	tableName    string
	pollInterval time.Duration
	batchSize    int32
}

func NewWorker(
	tableName string,
	batchSize int32,
	writer *sql.DB,
	blockchainPublisher blockchain.IBlockchainPublisher,
	logger *zap.Logger,
	pollInterval time.Duration,
) *Worker {
	maxInflight := int(batchSize) * 4

	// Pre-fill semaphore with tokens.
	// - Acquire a slot by receiving a token.
	// - Release a slot by sending a token.
	sem := make(chan struct{}, maxInflight)
	for range maxInflight {
		sem <- struct{}{}
	}

	return &Worker{
		logger:              logger,
		writer:              writer,
		wrtrQueries:         queries.New(writer),
		blockchainPublisher: blockchainPublisher,
		recvChan:            make(chan ISourceRecord, batchSize*2),
		wrtrChan:            make(chan *envelopes.OriginatorEnvelope, batchSize*2),
		sem:                 sem,
		inflight:            make(map[int64]time.Time),
		tableName:           tableName,
		pollInterval:        pollInterval,
		batchSize:           batchSize,
	}
}

func (w *Worker) startReader(ctx context.Context, reader ISourceReader) error {
	logger := w.logger.Named(utils.MigratorReaderLoggerName).
		With(zap.String(tableField, w.tableName))

	if reader == nil {
		logger.Error("reader is nil", zap.String(tableField, w.tableName))
		return errors.New("reader is nil")
	}

	logger.Info("started")

	tracing.GoPanicWrap(
		ctx,
		&w.wg,
		fmt.Sprintf("reader-%s", w.tableName),
		func(ctx context.Context) {
			defer close(w.recvChan)

			ticker := time.NewTicker(w.pollInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					logger.Info("context cancelled, stopping")
					return

				case <-ticker.C:
					startE2ELatency := time.Now()

					lastMigratedID, err := metrics.MeasureReaderLatency(
						"migration_tracker",
						func() (int64, error) {
							return w.wrtrQueries.GetMigrationProgress(ctx, w.tableName)
						},
					)
					if err != nil {
						metrics.EmitMigratorReaderError(
							"migration_tracker",
							err.Error(),
						)

						logger.Fatal("failed to get migration progress", zap.Error(err))
					}

					logger.Debug(
						"getting next batch of records",
						zap.Int64(lastMigratedIDField, lastMigratedID),
					)

					records, err := metrics.MeasureReaderLatency(
						w.tableName,
						func() ([]ISourceRecord, error) {
							return reader.Fetch(ctx, lastMigratedID, w.batchSize)
						},
					)
					if err != nil {
						switch err {
						case sql.ErrNoRows:
							logger.Info(noMoreRecordsToMigrateMessage)

							metrics.EmitMigratorReaderNumRowsFound(w.tableName, 0)

							select {
							case <-ctx.Done():
								return
							case <-time.After(sleepTimeOnNoRows):
							}

						default:
							metrics.EmitMigratorReaderError(w.tableName, err.Error())

							logger.Error(
								"getting next batch of records failed, retrying",
								zap.Error(err),
							)

							select {
							case <-ctx.Done():
								return
							case <-time.After(sleepTimeOnError):
							}
						}

						continue
					}

					metrics.EmitMigratorReaderNumRowsFound(w.tableName, int64(len(records)))

					if len(records) == 0 {
						logger.Info(noMoreRecordsToMigrateMessage)

						select {
						case <-ctx.Done():
							return
						case <-time.After(sleepTimeOnNoRows):
						}

						continue
					}

					for _, record := range records {
						id := record.GetID()

						_, seen := w.isInflight(id)

						if seen {
							continue
						}

						select {
						case <-ctx.Done():
							logger.Info(contextCancelledMessage)
							return
						case <-w.sem:
						}

						select {
						case <-ctx.Done():
							logger.Info(contextCancelledMessage)
							return

						case w.recvChan <- record:
							w.addInflightStartTime(id, startE2ELatency)

							logger.Debug(
								"sent record to transformer",
								zap.Int64(idField, id),
							)
						}
					}
				}
			}
		})

	return nil
}

func (w *Worker) startTransformer(ctx context.Context, transformer IDataTransformer) error {
	logger := w.logger.Named(utils.MigratorTransformerLoggerName).
		With(zap.String(tableField, w.tableName))

	logger.Info("started")

	tracing.GoPanicWrap(
		ctx,
		&w.wg,
		fmt.Sprintf("transformer-%s", w.tableName),
		func(ctx context.Context) {
			defer close(w.wrtrChan)

			for {
				select {
				case <-ctx.Done():
					logger.Info(contextCancelledMessage)
					return

				case record, open := <-w.recvChan:
					if !open {
						logger.Info(channelClosedMessage)
						return
					}

					envelope, err := transformer.Transform(record)
					if err != nil {
						logger.Error(
							"failed to transform",
							zap.Error(err),
							zap.Int64(idField, record.GetID()),
						)

						err := insertMigrationDeadLetterBox(
							ctx,
							w.writer,
							w.tableName,
							record.GetID(),
							nil,
							FailureTransformerError,
						)
						if err != nil {
							logger.Error("failed to insert dead letter box", zap.Error(err))
						}

						metrics.EmitMigratorTransformerError(w.tableName)

						w.cleanupInflight(ctx, record.GetID())
						continue
					}

					select {
					case <-ctx.Done():
						logger.Info(contextCancelledMessage)
						return

					case w.wrtrChan <- envelope:
						logger.Debug(
							"envelope sent to writer",
							utils.OriginatorIDField(envelope.OriginatorNodeID()),
							utils.SequenceIDField(int64(envelope.OriginatorSequenceID())),
						)
					}
				}
			}
		})

	return nil
}

func (w *Worker) startDatabaseWriter(ctx context.Context) error {
	logger := w.logger.Named(utils.MigratorWriterLoggerName).
		With(zap.String(tableField, w.tableName))

	logger.Info("started")

	tracing.GoPanicWrap(
		ctx,
		&w.wg,
		fmt.Sprintf("writer-%s", w.tableName),
		func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					logger.Info(contextCancelledMessage)
					return

				case envelope, open := <-w.wrtrChan:
					if !open {
						logger.Info(channelClosedMessage)
						return
					}

					sequenceID := int64(envelope.OriginatorSequenceID())

					err := metrics.MeasureWriterLatency(
						w.tableName,
						destinationDatabase,
						func() error {
							return retry(
								ctx,
								50*time.Millisecond,
								w.tableName,
								destinationDatabase,
								func() re.RetryableError {
									return w.insertOriginatorEnvelopeDatabase(
										ctx,
										envelope,
									)
								},
							)
						},
					)
					if err != nil {
						metrics.EmitMigratorWriterError(
							w.tableName,
							destinationDatabase,
							err.Error(),
						)

						logger.Error("failed to insert envelope", zap.Error(err))

						w.cleanupInflight(ctx, sequenceID)

						continue
					}

					startTime, exists := w.isInflight(sequenceID)

					if exists {
						metrics.EmitMigratorE2ELatency(
							w.tableName,
							destinationDatabase,
							time.Since(startTime).Seconds(),
						)
					}

					b, err := envelope.Bytes()
					if err != nil {
						logger.Warn("failed to marshal envelope", zap.Error(err))
					} else {
						metrics.EmitMigratorWriterBytesMigrated(w.tableName, destinationDatabase, len(b))
					}

					metrics.EmitMigratorWriterRowsMigrated(w.tableName, 1)

					metrics.EmitMigratorDestLastSequenceID(
						w.tableName,
						sequenceID,
					)

					w.cleanupInflight(ctx, sequenceID)
				}
			}
		},
	)

	return nil
}

func (w *Worker) startBlockchainWriterUnary(ctx context.Context) error {
	logger := w.logger.Named(utils.MigratorWriterLoggerName).
		With(zap.String(tableField, w.tableName))

	if w.blockchainPublisher == nil {
		logger.Error("blockchain publisher is nil")
		return errors.New("blockchain publisher is nil")
	}

	logger.Info("started")

	tracing.GoPanicWrap(
		ctx,
		&w.wg,
		fmt.Sprintf("writer-blockchain-%s", w.tableName),
		func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					logger.Info(contextCancelledMessage)
					return

				case envelope, open := <-w.wrtrChan:
					if !open {
						logger.Info(channelClosedMessage)
						return
					}

					sequenceID := int64(envelope.OriginatorSequenceID())

					err := metrics.MeasureWriterLatency(
						w.tableName,
						destinationBlockchain,
						func() error {
							return retry(
								ctx,
								50*time.Millisecond,
								w.tableName,
								destinationBlockchain,
								func() re.RetryableError {
									return w.insertOriginatorEnvelopeBlockchainUnary(
										ctx,
										envelope,
									)
								},
							)
						},
					)
					if err != nil {
						if !errors.Is(err, ErrDeadLetterBox) {
							logger.Error(
								"error publishing blockchain message",
								zap.Error(err),
							)
						}

						metrics.EmitMigratorWriterError(
							w.tableName,
							destinationBlockchain,
							err.Error(),
						)

						w.cleanupInflight(ctx, sequenceID)

						continue
					}

					startTime, exists := w.isInflight(sequenceID)

					if exists {
						metrics.EmitMigratorE2ELatency(
							w.tableName,
							destinationBlockchain,
							time.Since(startTime).Seconds(),
						)
					}

					b, err := envelope.Bytes()
					if err != nil {
						logger.Warn("failed to marshal envelope", zap.Error(err))
					} else {
						metrics.EmitMigratorWriterBytesMigrated(w.tableName, destinationBlockchain, len(b))
					}

					metrics.EmitMigratorWriterRowsMigrated(w.tableName, 1)

					metrics.EmitMigratorDestLastSequenceID(
						w.tableName,
						sequenceID,
					)

					w.cleanupInflight(ctx, sequenceID)
				}
			}
		})

	return nil
}

func (w *Worker) startBlockchainWriterIdentityUpdateBatches(ctx context.Context) error {
	logger := w.logger.Named(utils.MigratorWriterBatchLoggerName).
		With(zap.String(tableField, w.tableName))

	if w.tableName != inboxLogTableName {
		logger.Error("identity update batches are only supported for inbox log table")
		return errors.New("identity update batches are only supported for inbox log table")
	}

	if w.blockchainPublisher == nil {
		logger.Error("blockchain publisher is nil")
		return errors.New("blockchain publisher is nil")
	}

	logger.Info("started")

	tracing.GoPanicWrap(
		ctx,
		&w.wg,
		fmt.Sprintf("writer-identity-update-batches-%s", w.tableName),
		func(ctx context.Context) {
			// Flush the batch every 250 milliseconds. Arbitrum Orbit L3 min block time.
			ticker := time.NewTicker(250 * time.Millisecond)
			defer ticker.Stop()

			identityUpdateBatch := &IdentityUpdateBatch{}

			triggerBatchFlush := func() {
				lastSequenceID := identityUpdateBatch.LastSequenceID()

				logger.Info(
					"flushing identity update batch",
					zap.Int("length", identityUpdateBatch.Len()),
					zap.Uint64("last_sequence_id", lastSequenceID),
				)

				err := metrics.MeasureWriterLatency(
					w.tableName,
					destinationBlockchain,
					func() error {
						return w.flushIdentityUpdatesBatch(
							ctx,
							logger,
							identityUpdateBatch,
						)
					},
				)
				if err != nil {
					logger.Error(
						"failed to flush identity update batch",
						zap.Int("length", identityUpdateBatch.Len()),
						zap.Uint64("last_sequence_id", lastSequenceID),
						zap.Error(err),
					)

					metrics.EmitMigratorWriterError(
						w.tableName,
						destinationBlockchain,
						err.Error(),
					)

					for _, seqID := range identityUpdateBatch.sequenceIDs {
						w.cleanupInflight(ctx, int64(seqID))
					}

					identityUpdateBatch.Reset()

					return
				}

				logger.Info(
					"identity update batch flushed successfully",
					zap.Int("length", identityUpdateBatch.Len()),
					zap.Uint64("last_sequence_id", lastSequenceID),
				)

				for item := range identityUpdateBatch.All() {
					startTime, exists := w.isInflight(int64(item.SequenceID))

					if exists {
						metrics.EmitMigratorE2ELatency(
							w.tableName,
							destinationBlockchain,
							time.Since(startTime).Seconds(),
						)
					}

					metrics.EmitMigratorWriterBytesMigrated(
						w.tableName,
						destinationBlockchain,
						len(item.IdentityUpdate),
					)

					metrics.EmitMigratorWriterRowsMigrated(w.tableName, 1)

					metrics.EmitMigratorDestLastSequenceID(
						w.tableName,
						int64(item.SequenceID),
					)

					w.cleanupInflight(ctx, int64(item.SequenceID))
				}

				identityUpdateBatch.Reset()
			}

			for {
				select {
				case <-ctx.Done():
					logger.Info(contextCancelledMessage)

					return

				case <-ticker.C:
					if identityUpdateBatch.Len() <= 0 {
						continue
					}

					triggerBatchFlush()

				case envelope, open := <-w.wrtrChan:
					if !open {
						logger.Info(channelClosedMessage)

						if identityUpdateBatch.Len() <= 0 {
							return
						}

						// Flush remaining identity updates before exiting.
						triggerBatchFlush()

						return
					}

					// Prepare client envelope. On failure or oversized, insert into dead letter box.
					clientEnvelopeBytes, identifier, sequenceID, err := w.prepareClientEnvelope(
						ctx,
						envelope,
						w.tableName,
					)
					if err != nil {
						logger.Error(
							"failed to prepare identity update envelope",
							utils.InboxIDField(utils.HexEncode(identifier[:])),
							utils.SequenceIDField(int64(sequenceID)),
							zap.Error(err),
						)

						w.cleanupInflight(ctx, int64(envelope.OriginatorSequenceID()))

						continue
					}

					inboxID, err := utils.ParseInboxID(identifier)
					if err != nil {
						logger.Error(
							"failed to parse inbox ID",
							utils.InboxIDField(utils.HexEncode(identifier[:])),
							utils.SequenceIDField(int64(sequenceID)),
							zap.Error(err),
						)

						w.cleanupInflight(ctx, int64(envelope.OriginatorSequenceID()))

						continue
					}

					// Add to batch.
					identityUpdateBatch.Add(
						inboxID,
						clientEnvelopeBytes,
						sequenceID,
					)

					// This path triggers flushing only when the batch size exceeds the threshold.
					if identityUpdateBatch.Size() < maxBatchSize {
						continue
					}

					triggerBatchFlush()
				}
			}
		})

	return nil
}

/* Inflight management. */

func (w *Worker) isInflight(id int64) (time.Time, bool) {
	w.inflightMu.Lock()
	startTime, seen := w.inflight[id]
	w.inflightMu.Unlock()

	return startTime, seen
}

func (w *Worker) addInflightStartTime(id int64, startTime time.Time) {
	w.inflightMu.Lock()
	w.inflight[id] = startTime
	w.inflightMu.Unlock()
}

func (w *Worker) cleanupInflight(ctx context.Context, id int64) {
	w.inflightMu.Lock()
	delete(w.inflight, id)
	w.inflightMu.Unlock()

	select {
	case w.sem <- struct{}{}:
	case <-ctx.Done():
	default:
		// semaphore already full => double-release attempt; don't block
	}
}
