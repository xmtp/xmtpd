package message

import (
	"context"
	"errors"
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
	"go.uber.org/zap"
)

var errPublishFailed = errors.New("publish batch failed")

type publishWorker struct {
	ctx                context.Context
	logger             *zap.Logger
	listener           <-chan []queries.StagedOriginatorEnvelope
	notifier           chan<- bool
	registrant         *registrant.Registrant
	store              *db.Handler
	subscription       db.DBSubscription[queries.StagedOriginatorEnvelope, int64]
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

	query := func(ctx context.Context, lastSeenID int64, numRows int32) ([]queries.StagedOriginatorEnvelope, int64, error) {
		results, err := store.WriteQuery().SelectStagedOriginatorEnvelopes(
			ctx,
			queries.SelectStagedOriginatorEnvelopesParams{
				LastSeenID: lastSeenID,
				NumRows:    numRows,
			},
		)
		if err != nil {
			return nil, 0, err
		}
		if len(results) > 0 {
			lastSeenID = results[len(results)-1].ID
		}
		return results, lastSeenID, nil
	}

	notifier := make(chan bool, 1)
	subscription := db.NewDBSubscription(
		ctx,
		logger,
		query,
		0, // lastSeenID
		db.PollingOptions{Interval: time.Second, Notifier: notifier, NumRows: 100},
	)

	listener, err := subscription.Start()
	if err != nil {
		return nil, err
	}

	worker := &publishWorker{
		ctx:                ctx,
		logger:             logger,
		notifier:           notifier,
		subscription:       *subscription,
		listener:           listener,
		registrant:         reg,
		store:              store,
		feeCalculator:      feeCalculator,
		sleepOnFailureTime: sleepOnFailureTime,
		traceContextStore:  tracing.NewTraceContextStore(),
	}
	go worker.start()

	logger.Debug("started")

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

func (p *publishWorker) start() {
	for {
		select {
		case <-p.ctx.Done():
			return
		case batch, ok := <-p.listener:
			if !ok {
				p.logger.Error("listener is closed")
				return
			}

			p.logger.Info("processing batch", zap.Int("batch_size", len(batch)))
			for !p.publishBatch(batch) {
				time.Sleep(p.sleepOnFailureTime)
			}

			if len(batch) > 0 {
				p.lastProcessed.Store(batch[len(batch)-1].ID)
			}
		}
	}
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

func (p *publishWorker) publishBatch(batch []queries.StagedOriginatorEnvelope) (success bool) {
	originatorID := int32(p.registrant.NodeID())

	// Drain all stored trace contexts for this batch.
	// Link to the first available parent context for distributed tracing.
	var span tracing.Span
	var ctx context.Context
	traceLinked := false
	for _, stagedEnv := range batch {
		parentCtx := p.traceContextStore.Retrieve(stagedEnv.ID)
		if parentCtx != nil && !traceLinked {
			span = tracing.StartSpanWithParent(tracing.SpanPublishWorkerProcess, parentCtx)
			ctx = tracing.ContextWithSpan(p.ctx, span)
			traceLinked = true
		}
	}
	if !traceLinked {
		span, ctx = tracing.StartSpanFromContext(p.ctx, tracing.SpanPublishWorkerProcess)
	}
	tracing.SpanTag(span, tracing.TagTraceLinked, traceLinked)
	tracing.SpanTag(span, tracing.TagBatchSize, len(batch))
	tracing.SpanTag(span, tracing.TagOriginatorNode, originatorID)
	defer func() {
		if !success {
			span.Finish(tracing.WithError(errPublishFailed))
		} else {
			span.Finish()
		}
	}()

	prepared := make([]preparedEnvelope, 0, len(batch))
	var additionalMessages int32

	// Phase 1: CPU prep (per-envelope)
	for _, stagedEnv := range batch {
		logger := p.logger.With(utils.SequenceIDField(stagedEnv.ID))

		env, err := envelopes.NewPayerEnvelopeFromBytes(stagedEnv.PayerEnvelope)
		if err != nil {
			logger.Warn("failed to unmarshall originator envelope", zap.Error(err))
			return false
		}

		parsedTopic, err := topic.ParseTopic(stagedEnv.Topic)
		if err != nil {
			logger.Error("failed to parse topic", zap.Error(err))
			return false
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
				logger.Error("failed to calculate base fee", zap.Error(err))
				return false
			}

			congestionFee, err = p.feeCalculator.CalculateCongestionFee(
				p.ctx,
				p.store.Query(),
				stagedEnv.OriginatorTime,
				uint32(originatorID),
				additionalMessages,
			)
			if err != nil {
				logger.Error("failed to calculate congestion fee", zap.Error(err))
				return false
			}

			additionalMessages++
		}

		originatorEnv, err := p.registrant.SignStagedEnvelope(
			stagedEnv, baseFee, congestionFee, retentionDays,
		)
		if err != nil {
			logger.Error("failed to sign staged envelope", zap.Error(err))
			return false
		}

		validatedEnvelope, err := envelopes.NewOriginatorEnvelope(originatorEnv)
		if err != nil {
			logger.Error("failed to validate originator envelope", zap.Error(err))
			return false
		}

		originatorBytes, err := validatedEnvelope.Bytes()
		if err != nil {
			logger.Error("failed to marshal originator envelope", zap.Error(err))
			return false
		}

		payerAddress, err := validatedEnvelope.UnsignedOriginatorEnvelope.PayerEnvelope.RecoverSigner()
		if err != nil {
			logger.Error("failed to recover payer address", zap.Error(err))
			return false
		}

		prepared = append(prepared, preparedEnvelope{
			staged:          stagedEnv,
			originatorBytes: originatorBytes,
			payerAddress:    payerAddress.Hex(),
			isReserved:      isReserved,
			baseFee:         baseFee,
			congestionFee:   congestionFee,
			expiry: int64(
				validatedEnvelope.UnsignedOriginatorEnvelope.Proto().GetExpiryUnixtime(),
			),
		})
	}

	// Phase 2: Batch DB ops

	// 2a. Bulk find/create payers (deduplicated)
	uniqueAddresses := deduplicatePayerAddresses(prepared)
	payerRows, err := p.store.WriteQuery().BulkFindOrCreatePayers(p.ctx, uniqueAddresses)
	if err != nil {
		p.logger.Error("failed to bulk find/create payers", zap.Error(err))
		return false
	}
	payerMap := make(map[string]int32, len(payerRows))
	for _, row := range payerRows {
		payerMap[row.Address] = row.ID
	}

	// 2b. Build batch and insert
	insertSpan, _ := tracing.StartSpanFromContext(ctx, tracing.SpanPublishWorkerInsertGateway)
	batchInput := types.NewGatewayEnvelopeBatch()
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
			IsReserved:           prep.isReserved,
		})
	}

	inserted, err := db.InsertGatewayEnvelopeBatchV2AndIncrementUnsettledUsage(
		p.ctx, p.store.DB(), batchInput,
	)
	if p.ctx.Err() != nil {
		insertSpan.Finish()
		return false
	}
	if err != nil {
		insertSpan.Finish(tracing.WithError(err))
		p.logger.Error("failed to batch insert gateway envelopes", zap.Error(err))
		return false
	}
	tracing.SpanTag(insertSpan, "inserted_rows", inserted)
	insertSpan.Finish()

	// 2c. Bulk delete staged envelopes
	deleteSpan, _ := tracing.StartSpanFromContext(ctx, tracing.SpanPublishWorkerDeleteStaged)
	stagedIDs := make([]int64, len(prepared))
	for i, prep := range prepared {
		stagedIDs[i] = prep.staged.ID
	}
	deletedCount, err := p.store.WriteQuery().BulkDeleteStagedOriginatorEnvelopes(p.ctx, stagedIDs)
	if p.ctx.Err() != nil {
		deleteSpan.Finish()
		return true
	}
	if err != nil {
		deleteSpan.Finish(tracing.WithError(err))
		p.logger.Error("failed to bulk delete staged envelopes", zap.Error(err))
		return true
	}
	tracing.SpanTag(deleteSpan, "deleted_rows", deletedCount)
	deleteSpan.Finish()

	// Emit metrics
	for _, prep := range prepared {
		metrics.EmitAPIStagedEnvelopeProcessingDelay(time.Since(prep.staged.OriginatorTime))
	}

	if inserted > 0 {
		p.logger.Info("batch published",
			zap.Int("batch_size", len(prepared)),
			zap.Int64("inserted", inserted),
		)
	}

	return true
}

// calculateFees computes the base and congestion fees for a staged envelope.
// Used by the API handler to estimate fees before the publish worker processes the batch.
func (p *publishWorker) calculateFees(
	stagedEnv *queries.StagedOriginatorEnvelope,
	retentionDays uint32,
) (currency.PicoDollar, currency.PicoDollar, error) {
	baseFee, err := p.feeCalculator.CalculateBaseFee(
		stagedEnv.OriginatorTime,
		int64(len(stagedEnv.PayerEnvelope)),
		retentionDays,
	)
	if err != nil {
		return 0, 0, err
	}

	congestionFee, err := p.feeCalculator.CalculateCongestionFee(
		p.ctx,
		p.store.Query(),
		stagedEnv.OriginatorTime,
		p.registrant.NodeID(),
		0,
	)
	if err != nil {
		return 0, 0, err
	}

	return baseFee, congestionFee, nil
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
