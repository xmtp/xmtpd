package message

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/tracing"

	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/db/types"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/fees"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	"github.com/xmtp/xmtpd/pkg/utils/retryerrors"
	"go.uber.org/zap"
)

const (
	numRowsPerBatch    = int32(100)
	maxDeadlockRetries = 3
	tickerInterval     = time.Second
)

type publishWorker struct {
	ctx                context.Context
	logger             *zap.Logger
	notifier           chan bool
	registrant         *registrant.Registrant
	store              *db.Handler
	lastProcessed      atomic.Int64
	feeCalculator      fees.IFeeCalculator
	sleepOnFailureTime time.Duration
	traceContextStore  *tracing.TraceContextStore
}

func startPublishWorker(
	ctx context.Context,
	logger *zap.Logger,
	reg *registrant.Registrant,
	store *db.Handler,
	feeCalculator fees.IFeeCalculator,
	sleepOnFailureTime time.Duration,
) (*publishWorker, error) {
	logger = logger.Named(utils.PublishWorkerName)
	logger.Info("starting")

	notifier := make(chan bool, 1)
	worker := &publishWorker{
		ctx:                ctx,
		logger:             logger,
		notifier:           notifier,
		registrant:         reg,
		store:              store,
		feeCalculator:      feeCalculator,
		sleepOnFailureTime: sleepOnFailureTime,
		traceContextStore:  tracing.NewTraceContextStore(),
	}
	go worker.start()

	return worker, nil
}

func (p *publishWorker) notifyStagedPublish() {
	select {
	case p.notifier <- true:
	default:
	}
}

// storeTraceContext saves the span context for a staged envelope ID.
// This enables async trace propagation from the staging request to
// worker processing, creating a complete distributed trace.
func (p *publishWorker) storeTraceContext(stagedID int64, span tracing.Span) {
	p.traceContextStore.Store(stagedID, span)
}

// startEnvelopeSpans creates per-envelope trace spans linked to their API request contexts.
// Retrieve deletes each entry from the store, preventing memory leaks.
func (p *publishWorker) startEnvelopeSpans(
	staged []queries.StagedOriginatorEnvelope,
	originatorID int32,
) []tracing.Span {
	spans := make([]tracing.Span, len(staged))
	for i, stagedEnv := range staged {
		parentCtx := p.traceContextStore.Retrieve(stagedEnv.ID)
		if parentCtx != nil {
			spans[i] = tracing.StartSpanWithParent(
				tracing.SpanPublishWorkerProcess,
				parentCtx,
			)
			tracing.SpanTag(spans[i], tracing.TagTraceLinked, true)
		} else {
			spans[i], _ = tracing.StartSpanFromContext(
				p.ctx,
				tracing.SpanPublishWorkerProcess,
			)
			tracing.SpanTag(spans[i], tracing.TagTraceLinked, false)
		}
		tracing.SpanTag(spans[i], tracing.TagStagedID, stagedEnv.ID)
		tracing.SpanTag(spans[i], tracing.TagOriginatorNode, originatorID)
		tracing.SpanTag(spans[i], tracing.TagTopic, hex.EncodeToString(stagedEnv.Topic))
	}
	return spans
}

func (p *publishWorker) start() {
	p.logger.Info("started")
	defer p.logger.Info("stopped")

	timer := time.NewTimer(tickerInterval)
	defer timer.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-p.notifier:
		case <-timer.C:
		}
		p.pollAndPublish()
		timer.Reset(tickerInterval)
	}
}

// pollAndPublish processes batches in a loop until no more rows are available.
func (p *publishWorker) pollAndPublish() {
	for {
		count, err := p.processBatchWithRetry()
		if err != nil {
			p.logger.Error("failed to process batch", zap.Error(err))
			time.Sleep(p.sleepOnFailureTime)
			return
		}
		if count < numRowsPerBatch {
			return
		}
	}
}

// processBatchWithRetry retries the batch on transient database errors such as deadlocks.
func (p *publishWorker) processBatchWithRetry() (int32, error) {
	for attempt := range maxDeadlockRetries {
		count, err := p.processBatch()
		if err != nil && retryerrors.IsRetryableSQLError(err) && attempt < maxDeadlockRetries-1 {
			p.logger.Warn("retrying batch after transient error",
				zap.Int("attempt", attempt+1),
				zap.Error(err),
			)
			continue
		}
		return count, err
	}
	// unreachable
	return 0, nil
}

type batchResult struct {
	count           int32
	originatorTimes []time.Time
}

// processBatch runs the entire publish pipeline inside a single transaction:
// lock rows, compute fees, insert gateway envelopes, delete staged.
func (p *publishWorker) processBatch() (int32, error) {
	originatorID := int32(p.registrant.NodeID())

	var spans []tracing.Span

	result, err := db.RunInTxWithResult(
		p.ctx, p.store.DB(), &sql.TxOptions{},
		func(ctx context.Context, txQueries *queries.Queries) (batchResult, error) {
			staged, err := txQueries.SelectAndLockStagedEnvelopes(
				ctx, numRowsPerBatch,
			)
			if err != nil {
				return batchResult{}, fmt.Errorf("select and lock staged envelopes: %w", err)
			}
			if len(staged) == 0 {
				return batchResult{}, nil
			}

			spans = p.startEnvelopeSpans(staged, originatorID)

			p.logger.Debug(
				"processing batch", zap.Int("batch_size", len(staged)),
			)

			prepared, err := p.prepareEnvelopes(staged, txQueries)
			if err != nil {
				return batchResult{}, err
			}

			return p.persistBatch(ctx, txQueries, prepared)
		},
	)

	// Finish per-envelope spans outside the transaction so we capture
	// errors from the tx body, commit failures, and rollbacks.
	for _, span := range spans {
		if err != nil {
			span.Finish(tracing.WithError(err))
		} else {
			span.Finish()
		}
	}

	if err != nil {
		return 0, err
	}

	// Emit metrics after the transaction commits successfully.
	for _, t := range result.originatorTimes {
		metrics.EmitAPIStagedEnvelopeProcessingDelay(time.Since(t))
	}

	// Update lastProcessed from the DB to cover both our inserts and other workers'
	latestSeq, err := p.store.WriteQuery().GetLatestSequenceId(
		p.ctx, originatorID,
	)
	if err == nil && latestSeq > 0 {
		p.lastProcessed.Store(latestSeq)
	}

	if result.count > 0 {
		p.logger.Info("batch published", zap.Int32("batch_size", result.count))
	}

	return result.count, nil
}

type preparedEnvelope struct {
	staged          queries.StagedOriginatorEnvelope
	originatorBytes []byte
	payerAddress    string
	isReserved      bool
	baseFee         currency.PicoDollar
	congestionFee   currency.PicoDollar
	expiry          int64
}

// prepareEnvelopes signs and validates each staged envelope, computing fees along the way.
func (p *publishWorker) prepareEnvelopes(
	batch []queries.StagedOriginatorEnvelope,
	txQueries *queries.Queries,
) ([]preparedEnvelope, error) {
	originatorID := uint32(p.registrant.NodeID())
	prepared := make([]preparedEnvelope, 0, len(batch))
	batchCalc := p.feeCalculator.NewBatchFeeCalculator(
		p.ctx, txQueries, originatorID,
	)

	for _, stagedEnv := range batch {
		prep, err := p.prepareSingleEnvelope(stagedEnv, batchCalc)
		if err != nil {
			return nil, fmt.Errorf(
				"prepare envelope %d: %w", stagedEnv.ID, err,
			)
		}
		prepared = append(prepared, *prep)
	}

	return prepared, nil
}

// prepareSingleEnvelope parses, computes fees, signs, and validates a single staged envelope.
func (p *publishWorker) prepareSingleEnvelope(
	stagedEnv queries.StagedOriginatorEnvelope,
	batchCalc *fees.BatchFeeCalculator,
) (*preparedEnvelope, error) {
	env, err := envelopes.NewPayerEnvelopeFromBytes(stagedEnv.PayerEnvelope)
	if err != nil {
		return nil, err
	}

	parsedTopic, err := topic.ParseTopic(stagedEnv.Topic)
	if err != nil {
		return nil, err
	}

	isReserved := parsedTopic.IsReserved()
	retentionDays := env.RetentionDays()
	var baseFee, congestionFee currency.PicoDollar

	if !isReserved {
		baseFee, err = p.feeCalculator.CalculateBaseFee(
			stagedEnv.OriginatorTime,
			int64(len(stagedEnv.PayerEnvelope)),
			retentionDays,
		)
		if err != nil {
			return nil, err
		}

		congestionFee, err = batchCalc.CalculateCongestionFee(stagedEnv.OriginatorTime)
		if err != nil {
			return nil, err
		}
	}

	originatorEnv, err := p.registrant.SignStagedEnvelope(
		stagedEnv, baseFee, congestionFee, retentionDays,
	)
	if err != nil {
		return nil, err
	}

	validatedEnvelope, err := envelopes.NewOriginatorEnvelope(originatorEnv)
	if err != nil {
		return nil, err
	}

	originatorBytes, err := validatedEnvelope.Bytes()
	if err != nil {
		return nil, err
	}

	payerAddress, err := validatedEnvelope.UnsignedOriginatorEnvelope.PayerEnvelope.RecoverSigner()
	if err != nil {
		return nil, err
	}

	return &preparedEnvelope{
		staged:          stagedEnv,
		originatorBytes: originatorBytes,
		payerAddress:    payerAddress.Hex(),
		isReserved:      isReserved,
		baseFee:         baseFee,
		congestionFee:   congestionFee,
		expiry: int64(
			validatedEnvelope.UnsignedOriginatorEnvelope.Proto().GetExpiryUnixtime(),
		),
	}, nil
}

// persistBatch performs all DB operations within the transaction:
// payer resolution, batch insert, and staged deletion.
func (p *publishWorker) persistBatch(
	ctx context.Context,
	txQueries *queries.Queries,
	prepared []preparedEnvelope,
) (batchResult, error) {
	originatorID := int32(p.registrant.NodeID())

	// Bulk find/create payers (deduplicated)
	uniqueAddresses := deduplicatePayerAddresses(prepared)
	payerRows, err := txQueries.BulkFindOrCreatePayers(ctx, uniqueAddresses)
	if err != nil {
		return batchResult{}, fmt.Errorf("bulk find/create payers: %w", err)
	}
	payerMap := make(map[string]int32, len(payerRows))
	for _, row := range payerRows {
		payerMap[row.Address] = row.ID
	}

	// Build batch and collect staged IDs for deletion
	batchInput := types.NewGatewayEnvelopeBatch()
	stagedIDs := make([]int64, 0, len(prepared))
	originatorTimes := make([]time.Time, 0, len(prepared))
	for _, prep := range prepared {
		batchInput.Add(types.GatewayEnvelopeRow{
			OriginatorNodeID:     originatorID,
			OriginatorSequenceID: prep.staged.ID,
			Topic:                prep.staged.Topic,
			PayerID:              payerMap[prep.payerAddress],
			GatewayTime:          prep.staged.OriginatorTime,
			Expiry:               prep.expiry,
			OriginatorEnvelope:   prep.originatorBytes,
			SpendPicodollars:     int64(prep.baseFee) + int64(prep.congestionFee),
			CountUsage:           !prep.isReserved,
			CountCongestion:      !prep.isReserved,
		})
		stagedIDs = append(stagedIDs, prep.staged.ID)
		originatorTimes = append(originatorTimes, prep.staged.OriginatorTime)
	}

	insertSpan, _ := tracing.StartSpanFromContext(ctx, tracing.SpanPublishWorkerInsertGateway)
	inserted, err := db.InsertGatewayEnvelopeBatchV2Transactional(
		ctx, txQueries, p.logger, batchInput,
	)
	if err != nil {
		insertSpan.Finish(tracing.WithError(err))
		return batchResult{}, fmt.Errorf("batch insert gateway envelopes: %w", err)
	}
	tracing.SpanTag(insertSpan, "inserted_rows", inserted)
	insertSpan.Finish()

	deleteSpan, _ := tracing.StartSpanFromContext(ctx, tracing.SpanPublishWorkerDeleteStaged)
	_, err = txQueries.BulkDeleteStagedOriginatorEnvelopes(ctx, stagedIDs)
	if err != nil {
		deleteSpan.Finish(tracing.WithError(err))
		return batchResult{}, fmt.Errorf("bulk delete staged envelopes: %w", err)
	}
	deleteSpan.Finish()

	if inserted > 0 {
		p.logger.Info("batch inserted",
			zap.Int("batch_size", len(prepared)),
			zap.Int64("inserted", inserted),
		)
	}

	return batchResult{
		count:           int32(len(prepared)),
		originatorTimes: originatorTimes,
	}, nil
}

func deduplicatePayerAddresses(prepared []preparedEnvelope) []string {
	seen := make(map[string]struct{}, len(prepared))
	result := make([]string, 0, len(prepared))
	for _, p := range prepared {
		if _, ok := seen[p.payerAddress]; !ok {
			seen[p.payerAddress] = struct{}{}
			result = append(result, p.payerAddress)
		}
	}
	return result
}
