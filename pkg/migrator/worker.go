package migrator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/types"
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
	writer              *db.Handler
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
	writer *db.Handler,
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

func (w *Worker) StartReader(ctx context.Context, reader ISourceReader) error {
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
					logger.Debug(contextCancelledMessage)
					return

				case <-ticker.C:
					startE2ELatency := time.Now()

					lastMigratedID, err := metrics.MeasureReaderLatency(
						"migration_tracker",
						func() (int64, error) {
							return w.writer.ReadQuery().GetMigrationProgress(ctx, w.tableName)
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

func (w *Worker) StartTransformer(ctx context.Context, transformer IDataTransformer) error {
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
					logger.Debug(contextCancelledMessage)

					return

				case record, open := <-w.recvChan:
					if !open {
						logger.Debug(channelClosedMessage)

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
							w.writer.Write(),
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

func (w *Worker) StartDatabaseWriter(ctx context.Context) error {
	logger := w.logger.Named(utils.MigratorWriterDatabaseLoggerName).
		With(zap.String(tableField, w.tableName))

	logger.Info("started")

	tracing.GoPanicWrap(
		ctx,
		&w.wg,
		fmt.Sprintf("writer-database-%s", w.tableName),
		func(ctx context.Context) {
			ticker := time.NewTicker(250 * time.Millisecond)
			defer ticker.Stop()

			batch := types.NewGatewayEnvelopeBatch()

			triggerBatchFlush := func() {
				var (
					batchLen            = batch.Len()
					batchLastSequenceID = batch.LastSequenceID()
				)

				logger.Info(
					"publishing batch",
					utils.LengthField(batchLen),
					utils.SequenceIDField(batchLastSequenceID),
				)

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
								return w.insertOriginatorEnvelopeDatabaseBatch(
									ctx,
									logger,
									batch,
								)
							},
						)
					},
				)
				if err != nil {
					if errors.Is(err, context.Canceled) ||
						errors.Is(err, context.DeadlineExceeded) {
						logger.Info(contextCancelledMessage)
						return
					}

					logger.Error(
						"failed to insert batch",
						zap.Error(err),
						utils.LengthField(batchLen),
						utils.SequenceIDField(batchLastSequenceID),
					)

					metrics.EmitMigratorWriterError(
						w.tableName,
						destinationDatabase,
						err.Error(),
					)

					for _, envelope := range batch.All() {
						w.cleanupInflight(ctx, envelope.OriginatorSequenceID)
					}

					batch.Reset()

					return
				}

				logger.Info(
					"batch published successfully",
					utils.LengthField(batchLen),
					utils.SequenceIDField(batchLastSequenceID),
				)

				metrics.EmitMigratorWriterRowsMigrated(w.tableName, int64(batch.Len()))

				for _, envelope := range batch.All() {
					startTime, exists := w.isInflight(envelope.OriginatorSequenceID)

					if exists {
						metrics.EmitMigratorE2ELatency(
							w.tableName,
							destinationDatabase,
							time.Since(startTime).Seconds(),
						)
					}

					metrics.EmitMigratorWriterBytesMigrated(
						w.tableName,
						destinationDatabase,
						len(envelope.OriginatorEnvelope),
					)

					metrics.EmitMigratorDestLastSequenceID(
						w.tableName,
						envelope.OriginatorSequenceID,
					)

					w.cleanupInflight(ctx, envelope.OriginatorSequenceID)
				}

				batch.Reset()
			}

			for {
				select {
				case <-ctx.Done():
					logger.Info(contextCancelledMessage)
					return

				case <-ticker.C:
					if batch.Len() <= maxDatabaseBatchSize/2 {
						continue
					}

					triggerBatchFlush()

				case envelope, open := <-w.wrtrChan:
					if !open {
						logger.Info(channelClosedMessage)

						if batch.Len() <= 0 {
							return
						}

						// Flush remaining messages before exiting.
						triggerBatchFlush()

						return
					}

					if envelope == nil {
						continue
					}

					// Batches should only contain envelopes from the same originator ID.
					if originatorIDToTableName(envelope.OriginatorNodeID()) != w.tableName {
						continue
					}

					if batch.Len() >= maxDatabaseBatchSize {
						triggerBatchFlush()
					}

					payerID, err := w.payerIDFromEnvelope(ctx, envelope)
					if err != nil {
						logger.Error("failed to get payer ID", zap.Error(err))
						continue
					}

					envelopeBytes, err := envelope.Bytes()
					if err != nil {
						logger.Error("failed to get envelope bytes", zap.Error(err))
						continue
					}

					batch.Add(types.GatewayEnvelopeRow{
						OriginatorNodeID:     int32(envelope.OriginatorNodeID()),
						OriginatorSequenceID: int64(envelope.OriginatorSequenceID()),
						Topic:                envelope.TargetTopic().Bytes(),
						PayerID:              payerID,
						GatewayTime:          envelope.OriginatorTime(),
						Expiry: int64(
							envelope.UnsignedOriginatorEnvelope.Proto().GetExpiryUnixtime(),
						),
						OriginatorEnvelope: envelopeBytes,
						SpendPicodollars: int64(
							envelope.UnsignedOriginatorEnvelope.Proto().BaseFeePicodollars,
						),
					})
				}
			}
		},
	)

	return nil
}

func (w *Worker) payerIDFromEnvelope(
	ctx context.Context,
	envelope *envelopes.OriginatorEnvelope,
) (int32, error) {
	payerAddress, err := envelope.UnsignedOriginatorEnvelope.PayerEnvelope.RecoverSigner()
	if err != nil {
		return 0, err
	}

	payerID, err := w.writer.WriteQuery().FindOrCreatePayer(ctx, payerAddress.Hex())
	if err != nil {
		return 0, err
	}

	return payerID, nil
}

func (w *Worker) StartBlockchainWriterBatch(ctx context.Context) error {
	logger := w.logger.Named(utils.MigratorWriterChainLoggerName).
		With(zap.String(tableField, w.tableName))

	if w.tableName != commitMessagesTableName && w.tableName != inboxLogTableName {
		return errors.New(
			"broadcaster batches are only supported for commit messages and inbox log tables",
		)
	}

	if w.blockchainPublisher == nil {
		return errors.New("blockchain publisher is nil")
	}

	logger.Info("started")

	tracing.GoPanicWrap(
		ctx,
		&w.wg,
		fmt.Sprintf("writer-chain-%s", w.tableName),
		func(ctx context.Context) {
			// Flush the batch every 250 milliseconds. Arbitrum Orbit L3 min block time.
			ticker := time.NewTicker(250 * time.Millisecond)
			defer ticker.Stop()

			var batch *BroadcasterBatch

			switch w.tableName {
			case commitMessagesTableName:
				batch = &BroadcasterBatch{
					identifierLength: 16,
				}

			case inboxLogTableName:
				batch = &BroadcasterBatch{
					identifierLength: 32,
				}
			}

			triggerBatchFlush := func() {
				var (
					batchLen            = batch.Len()
					batchLastSequenceID = batch.LastSequenceID()
				)

				logger.Info(
					"publishing batch",
					utils.LengthField(batchLen),
					utils.SequenceIDField(int64(batchLastSequenceID)),
				)

				// flushBroadcasterBatch handles:
				// 1. Batch insert attempt.
				// 2. On batch failure: individual retries.
				// 3. On individual failure: dead letter box insertion.
				// It only returns an error on context cancellation.
				err := metrics.MeasureWriterLatency(
					w.tableName,
					destinationBlockchain,
					func() error {
						return w.flushBroadcasterBatch(
							ctx,
							logger,
							batch,
						)
					},
				)
				if err != nil {
					if errors.Is(err, context.Canceled) ||
						errors.Is(err, context.DeadlineExceeded) {
						logger.Info(contextCancelledMessage)
						return
					}

					logger.Error(
						"failed to publish batch",
						utils.LengthField(batch.Len()),
						utils.SequenceIDField(int64(batchLastSequenceID)),
						zap.Error(err),
					)

					metrics.EmitMigratorWriterError(
						w.tableName,
						destinationBlockchain,
						err.Error(),
					)

					for _, seqID := range batch.sequenceIDs {
						w.cleanupInflight(ctx, int64(seqID))
					}

					batch.Reset()

					return
				}

				logger.Info(
					"batch published successfully",
					utils.LengthField(batch.Len()),
					utils.SequenceIDField(int64(batchLastSequenceID)),
				)

				for item := range batch.All() {
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
						len(item.Payload),
					)

					metrics.EmitMigratorWriterRowsMigrated(w.tableName, 1)

					metrics.EmitMigratorDestLastSequenceID(
						w.tableName,
						int64(item.SequenceID),
					)

					w.cleanupInflight(ctx, int64(item.SequenceID))
				}

				batch.Reset()
			}

			for {
				select {
				case <-ctx.Done():
					logger.Debug(contextCancelledMessage)
					return

				case <-ticker.C:
					if batch.Len() <= 0 {
						continue
					}

					triggerBatchFlush()

				case envelope, open := <-w.wrtrChan:
					if !open {
						logger.Debug(channelClosedMessage)

						if batch.Len() <= 0 {
							return
						}

						// Flush remaining messages before exiting.
						triggerBatchFlush()

						return
					}

					if envelope == nil {
						continue
					}

					// Prepare client envelope. On failure or oversized, insert into dead letter box.
					clientEnvelope, identifier, sequenceID, err := w.prepareClientEnvelope(
						ctx,
						logger,
						envelope,
						w.tableName,
					)
					if err != nil {
						logger.Warn(
							"envelope preparation failed, added to dead letter box",
							zap.Error(err),
						)

						w.cleanupInflight(ctx, int64(envelope.OriginatorSequenceID()))

						continue
					}

					// messages at this point can be up to 200KB.
					messageSize := int64(len(clientEnvelope) + len(identifier) + 8)

					// Only add an element if the resulting batch size is less than the threshold.
					// Otherwise, flush the current batch and start a new one.
					if batch.Size() > 0 &&
						batch.Size()+messageSize >= maxChainBatchSize {
						triggerBatchFlush()
					}

					// Add to batch.
					batch.Add(
						identifier,
						clientEnvelope,
						sequenceID,
					)

					// This path triggers flushing only when the batch size exceeds the threshold.
					if batch.Size() < maxChainBatchSize {
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
