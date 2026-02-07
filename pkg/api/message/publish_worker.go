package message

import (
	"context"
	"encoding/hex"
	"errors"
	"sync/atomic"
	"time"

	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/tracing"

	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/fees"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

var errPublishFailed = errors.New("publish staged envelope failed")

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
	// traceContextStore enables async trace context propagation from staging
	// requests to worker processing, allowing end-to-end distributed tracing
	traceContextStore *tracing.TraceContextStore
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

			for _, stagedEnv := range batch {
				p.logger.Info("publishing envelope", utils.SequenceIDField(stagedEnv.ID))
				for !p.publishStagedEnvelope(stagedEnv) {
					// Infinite retry on failure to publish; we cannot
					// continue to the next envelope until this one is processed
					time.Sleep(p.sleepOnFailureTime)
				}
				p.lastProcessed.Store(stagedEnv.ID)
				metrics.EmitApiStagedEnvelopeProcessingDelay(time.Since(stagedEnv.OriginatorTime))
			}
		}
	}
}

func (p *publishWorker) publishStagedEnvelope(
	stagedEnv queries.StagedOriginatorEnvelope,
) (success bool) {
	// Retrieve parent span context from async trace propagation
	// This links the worker processing to the original staging request
	parentCtx := p.traceContextStore.Retrieve(stagedEnv.ID)

	// Create APM span, linked to parent if available
	var span tracing.Span
	var ctx context.Context
	if parentCtx != nil {
		// Linked to original staging request - full distributed trace!
		span = tracing.StartSpanWithParent(tracing.SpanPublishWorkerProcess, parentCtx)
		ctx = tracing.ContextWithSpan(p.ctx, span)
		tracing.SpanTag(span, tracing.TagTraceLinked, true)
	} else {
		// No parent context - timer fallback or context expired
		span, ctx = tracing.StartSpanFromContext(p.ctx, tracing.SpanPublishWorkerProcess)
		tracing.SpanTag(span, tracing.TagTraceLinked, false)
	}
	defer func() {
		if !success {
			span.Finish(tracing.WithError(errPublishFailed))
		} else {
			span.Finish()
		}
	}()

	originatorID := int32(p.registrant.NodeID())

	// Add searchable tags for debugging (same info as our PERF_TRACE logs)
	tracing.SpanTag(span, tracing.TagStagedID, stagedEnv.ID)
	tracing.SpanTag(span, tracing.TagOriginatorNode, originatorID)
	tracing.SpanTag(span, tracing.TagTopic, hex.EncodeToString(stagedEnv.Topic))

	logger := p.logger.With(
		utils.SequenceIDField(stagedEnv.ID),
		utils.OriginatorIDField(uint32(originatorID)),
	)
	logger = tracing.Link(span, logger)

	env, err := envelopes.NewPayerEnvelopeFromBytes(stagedEnv.PayerEnvelope)
	if err != nil {
		logger.Warn("failed to unmarshall originator envelope", zap.Error(err))
		return false
	}

	parsedTopic, err := topic.ParseTopic(stagedEnv.Topic)
	if err != nil {
		return false
	}

	var (
		baseFee         currency.PicoDollar
		congestionFee   currency.PicoDollar
		isReservedTopic = parsedTopic.IsReserved()
		retentionDays   = env.RetentionDays()
	)

	if !isReservedTopic {
		feeSpan, _ := tracing.StartSpanFromContext(ctx, tracing.SpanPublishWorkerCalculateFees)
		baseFee, congestionFee, err = p.calculateFees(&stagedEnv, retentionDays)
		if err != nil {
			feeSpan.Finish(tracing.WithError(err))
			logger.Error("failed to calculate fees", zap.Error(err))
			return false
		}
		tracing.SpanTag(feeSpan, "base_fee", int64(baseFee))
		tracing.SpanTag(feeSpan, "congestion_fee", int64(congestionFee))
		feeSpan.Finish()
	}

	signSpan, _ := tracing.StartSpanFromContext(ctx, tracing.SpanPublishWorkerSignEnvelope)
	originatorEnv, err := p.registrant.SignStagedEnvelope(
		stagedEnv,
		baseFee,
		congestionFee,
		retentionDays,
	)
	if err != nil {
		signSpan.Finish(tracing.WithError(err))
		logger.Error(
			"failed to sign staged envelope",
			zap.Error(err),
		)
		return false
	}
	signSpan.Finish()

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

	payerID, err := p.store.WriteQuery().FindOrCreatePayer(ctx, payerAddress.Hex())
	if err != nil {
		logger.Error("failed to find or create payer", zap.Error(err))
		return false
	}

	// On unique constraint conflicts, no error is thrown, but numRows is 0
	var inserted int64

	insertSpan, _ := tracing.StartSpanFromContext(ctx, tracing.SpanPublishWorkerInsertGateway)
	tracing.SpanTag(insertSpan, "is_reserved_topic", isReservedTopic)

	if isReservedTopic {

		// Reserved topics are not charged fees, so we only need to insert the envelope into the database.
		rows, err := db.InsertGatewayEnvelopeWithChecksStandalone(
			ctx,
			p.store.WriteQuery(),
			queries.InsertGatewayEnvelopeParams{
				OriginatorNodeID:     originatorID,
				OriginatorSequenceID: stagedEnv.ID,
				Topic:                stagedEnv.Topic,
				OriginatorEnvelope:   originatorBytes,
				PayerID:              db.NullInt32(payerID),
				Expiry: int64(
					validatedEnvelope.UnsignedOriginatorEnvelope.Proto().GetExpiryUnixtime(),
				),
			})
		if err != nil {
			insertSpan.Finish(tracing.WithError(err))
			logger.Error("failed to insert gateway envelope with reserved topic", zap.Error(err))
			return false
		}

		if rows.InsertedMetaRows > 0 {
			inserted = rows.InsertedMetaRows
		}

	} else {
		inserted, err = db.InsertGatewayEnvelopeAndIncrementUnsettledUsage(
			ctx,
			p.store.DB(),
			queries.InsertGatewayEnvelopeParams{
				OriginatorNodeID:     originatorID,
				OriginatorSequenceID: stagedEnv.ID,
				Topic:                stagedEnv.Topic,
				OriginatorEnvelope:   originatorBytes,
				PayerID:              db.NullInt32(payerID),
				GatewayTime:          stagedEnv.OriginatorTime,
				Expiry: int64(
					validatedEnvelope.UnsignedOriginatorEnvelope.Proto().GetExpiryUnixtime(),
				),
			},
			queries.IncrementUnsettledUsageParams{
				PayerID:           payerID,
				OriginatorID:      originatorID,
				MinutesSinceEpoch: utils.MinutesSinceEpoch(stagedEnv.OriginatorTime),
				SpendPicodollars:  int64(baseFee) + int64(congestionFee),
			},
		)
	}

	if ctx.Err() != nil {
		insertSpan.Finish()
		return false
	} else if err != nil {
		insertSpan.Finish(tracing.WithError(err))
		logger.Error("failed to insert gateway envelope", zap.Error(err))
		return false
	} else if inserted == 0 {
		// Envelope was already inserted by another worker
		tracing.SpanTag(insertSpan, "already_inserted", true)
		logger.Debug("envelope already inserted")
	}
	tracing.SpanTag(insertSpan, "inserted_rows", inserted)
	insertSpan.Finish()

	// Try to delete the row regardless of if the gateway envelope was inserted elsewhere
	deleteSpan, _ := tracing.StartSpanFromContext(ctx, tracing.SpanPublishWorkerDeleteStaged)
	deleted, err := p.store.WriteQuery().DeleteStagedOriginatorEnvelope(ctx, stagedEnv.ID)
	if ctx.Err() != nil {
		deleteSpan.Finish()
		return true
	} else if err != nil {
		deleteSpan.Finish(tracing.WithError(err))
		logger.Error("failed to delete staged envelope", zap.Error(err))
		// Envelope is already inserted, so it is safe to continue
		return true
	} else if deleted == 0 {
		// Envelope was already deleted by another worker
		tracing.SpanTag(deleteSpan, "already_deleted", true)
		logger.Debug("envelope already deleted")
	}
	tracing.SpanTag(deleteSpan, "deleted_rows", deleted)
	deleteSpan.Finish()

	return true
}

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

	// TODO:nm: Set this to the actual congestion fee
	// For now we are setting congestion to 0
	congestionFee, err := p.feeCalculator.CalculateCongestionFee(
		p.ctx,
		p.store.Query(),
		stagedEnv.OriginatorTime,
		p.registrant.NodeID(),
	)
	if err != nil {
		return 0, 0, err
	}

	return baseFee, congestionFee, nil
}
