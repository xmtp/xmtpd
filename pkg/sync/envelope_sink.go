package sync

import (
	"context"
	"encoding/hex"
	"errors"
	"time"

	"github.com/cenkalti/backoff/v4"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/xmtp/xmtpd/pkg/migrator"
	"github.com/xmtp/xmtpd/pkg/utils/retryerrors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	envUtils "github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/fees"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

type EnvelopeSink struct {
	ctx                        context.Context
	db                         *db.Handler
	logger                     *zap.Logger
	feeCalculator              fees.IFeeCalculator
	payerReportStore           payerreport.IPayerReportStore
	payerReportDomainSeparator common.Hash
	writeQueue                 chan *envUtils.OriginatorEnvelope
	errorRetrySleepTime        time.Duration
	payerCache                 *lru.Cache[string, int32]
}

func newEnvelopeSink(
	ctx context.Context,
	db *db.Handler,
	logger *zap.Logger,
	feeCalculator fees.IFeeCalculator,
	payerReportStore payerreport.IPayerReportStore,
	payerReportDomainSeparator common.Hash,
	writeQueue chan *envUtils.OriginatorEnvelope,
	errorRetrySleepTime time.Duration,
) *EnvelopeSink {
	// Per-entry cost in hashicorp/golang-lru/v2:
	// - Doubly-linked list node: 3 pointers × 8 B = 24 B
	// - Map entry: string header (16 B) + hex address (42 B) + int32 (4 B) + map overhead (~20 B) = 82 B
	// - Total per entry: ~106 B
	// At 1000 entries: ~106 KB
	payerCache, _ := lru.New[string, int32](1000)

	return &EnvelopeSink{
		ctx:                        ctx,
		db:                         db,
		logger:                     logger,
		feeCalculator:              feeCalculator,
		payerReportStore:           payerReportStore,
		payerReportDomainSeparator: payerReportDomainSeparator,
		writeQueue:                 writeQueue,
		errorRetrySleepTime:        errorRetrySleepTime,
		payerCache:                 payerCache,
	}
}

func (s *EnvelopeSink) Start() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case env, ok := <-s.writeQueue:
			if !ok {
				s.logger.Debug("writeQueue is closed")
				return
			}

			if env == nil {
				continue
			}

			boCtx := backoff.WithContext(
				utils.NewBackoff(10*time.Millisecond, s.errorRetrySleepTime, 0), s.ctx,
			)

			operation := func() error {
				select {
				case <-s.ctx.Done():
					return backoff.Permanent(errors.New("shutting down"))
				default:
					err := s.storeEnvelope(env)
					if err != nil {
						s.logger.Error("error storing envelope", zap.Error(err))

						if !retryerrors.IsRetryableSQLError(err) {
							s.logger.Error("Unexpected runtime error. Retry might be indefinite.")
						}

						return err
					}

					return nil
				}
			}

			err := backoff.Retry(operation, boCtx)
			if err != nil {
				return
			}
		}
	}
}

func (s *EnvelopeSink) storeEnvelope(env *envUtils.OriginatorEnvelope) error {
	// Create APM span for sync worker storing envelope from another node
	span, ctx := tracing.StartSpanFromContext(s.ctx, tracing.SpanSyncWorkerStoreEnvelope)
	defer span.Finish()

	// Tag with envelope info for debugging
	tracing.SpanTag(span, "source_node", env.OriginatorNodeID())
	tracing.SpanTag(span, "sequence_id", env.OriginatorSequenceID())
	tracing.SpanTag(span, "topic", hex.EncodeToString(env.TargetTopic().Bytes()))
	tracing.SpanTag(span, "is_reserved", env.TargetTopic().IsReserved())

	if env.TargetTopic().IsReserved() {
		s.logger.Info(
			"found envelope with reserved topic",
			utils.TopicField(env.TargetTopic().String()),
		)
		return s.storeReservedEnvelope(env, ctx)
	}

	// Calculate the fees independently to verify the originator's calculation
	feeSpan, _ := tracing.StartSpanFromContext(ctx, tracing.SpanSyncWorkerVerifyFees)
	ourFeeCalculation, err := s.calculateFees(env)
	if err != nil {
		feeSpan.Finish(tracing.WithError(err))
		s.logger.Error("failed to calculate fees", zap.Error(err))
		return err
	}
	feeSpan.Finish()
	originatorsFeeCalculation := env.UnsignedOriginatorEnvelope.BaseFee() +
		env.UnsignedOriginatorEnvelope.CongestionFee()

	if ourFeeCalculation != originatorsFeeCalculation {
		s.logger.Warn(
			"fee calculation mismatch",
			zap.String("our_fee", ourFeeCalculation.String()),
			zap.String("originator_fee", originatorsFeeCalculation.String()),
		)
	}

	// If for some reason the envelope is not able to marshal (but has made it this far)
	// the node will retry indefinitely.
	// I don't think this will ever happen
	originatorBytes, err := env.Bytes()
	if err != nil {
		s.logger.Error("failed to marshal originator envelope", zap.Error(err))
		return err
	}

	// The payer address has already been validated, so any errors here should be transient
	payerID, err := s.getPayerID(env)
	if err != nil {
		s.logger.Error("failed to get payer ID", zap.Error(err))
		return err
	}

	originatorID := int32(env.OriginatorNodeID())
	originatorTime := utils.NsToDate(env.OriginatorNs())
	expiry := env.UnsignedOriginatorEnvelope.Proto().GetExpiryUnixtime()

	insertSpan, _ := tracing.StartSpanFromContext(ctx, tracing.SpanSyncWorkerInsertGateway)
	inserted, err := db.InsertGatewayEnvelopeAndIncrementUnsettledUsage(
		ctx,
		s.db.Write(),
		queries.InsertGatewayEnvelopeParams{
			OriginatorNodeID:     int32(env.OriginatorNodeID()),
			OriginatorSequenceID: int64(env.OriginatorSequenceID()),
			Topic:                env.TargetTopic().Bytes(),
			OriginatorEnvelope:   originatorBytes,
			PayerID:              db.NullInt32(payerID),
			Expiry:               int64(expiry),
		},
		queries.IncrementUnsettledUsageParams{
			PayerID:           payerID,
			OriginatorID:      originatorID,
			MinutesSinceEpoch: utils.MinutesSinceEpoch(originatorTime),
			SpendPicodollars:  int64(ourFeeCalculation),
			MessageCount:      1,
		},
		!migrator.IsMigratorOriginatorID(env.OriginatorNodeID()),
	)

	if err != nil {
		insertSpan.Finish(tracing.WithError(err))
		s.logger.Error("failed to insert gateway envelope", zap.Error(err))
		return err
	} else if inserted == 0 {
		// Envelope was already inserted by another worker
		tracing.SpanTag(insertSpan, "already_inserted", true)
		s.logger.Debug("envelope already inserted",
			utils.OriginatorIDField(env.OriginatorNodeID()),
			utils.SequenceIDField(int64(env.OriginatorSequenceID())),
		)
		insertSpan.Finish()
		return nil
	}
	tracing.SpanTag(insertSpan, "inserted_rows", inserted)
	insertSpan.Finish()

	return nil
}

func (s *EnvelopeSink) storeReservedEnvelope(
	env *envUtils.OriginatorEnvelope,
	ctx context.Context,
) error {
	// Create APM span for reserved envelope processing
	span, ctx := tracing.StartSpanFromContext(ctx, tracing.SpanSyncWorkerStoreReservedEnvelope)
	defer span.Finish()

	tracing.SpanTag(span, "topic_kind", env.TargetTopic().Kind().String())

	payerID, err := s.getPayerID(env)
	if err != nil {
		s.logger.Error("failed to get payer ID", zap.Error(err))
		return err
	}

	switch env.TargetTopic().Kind() {
	case topic.TopicKindPayerReportsV1:
		reportSpan, reportCtx := tracing.StartSpanFromContext(
			ctx,
			tracing.SpanSyncWorkerStorePayerReport,
		)
		err := s.payerReportStore.StoreSyncedReport(
			reportCtx,
			env,
			payerID,
			s.payerReportDomainSeparator,
		)
		if err != nil {
			reportSpan.Finish(tracing.WithError(err))
			s.logger.Error("failed to store synced report", zap.Error(err))
			// Return nil here to avoid infinite retries
		} else {
			reportSpan.Finish()
		}
		return nil
	case topic.TopicKindPayerReportAttestationsV1:
		attestSpan, attestCtx := tracing.StartSpanFromContext(
			ctx,
			tracing.SpanSyncWorkerStoreAttestation,
		)
		err := s.payerReportStore.StoreSyncedAttestation(
			attestCtx,
			env,
			payerID,
		)
		if err != nil {
			attestSpan.Finish(tracing.WithError(err))
			s.logger.Error("failed to store synced attestation", zap.Error(err))
			// Return nil here to avoid infinite retries
		} else {
			attestSpan.Finish()
		}
		return nil
	default:
		s.logger.Info(
			"received unknown reserved topic",
			utils.TopicField(env.TargetTopic().String()),
		)
		return nil
	}
}

func (s *EnvelopeSink) calculateFees(
	env *envUtils.OriginatorEnvelope,
) (currency.PicoDollar, error) {
	var (
		payerEnvelopeLength = len(env.UnsignedOriginatorEnvelope.PayerEnvelopeBytes())
		messageTime         = utils.NsToDate(env.OriginatorNs())
		congestionFee       currency.PicoDollar
	)

	baseFee, err := s.feeCalculator.CalculateBaseFee(
		messageTime,
		int64(payerEnvelopeLength),
		env.UnsignedOriginatorEnvelope.PayerEnvelope.RetentionDays(),
	)
	if err != nil {
		return 0, err
	}

	if migrator.IsMigratorOriginatorID(env.OriginatorNodeID()) {
		return baseFee + 0, nil
	}

	// NOTE: This is code smell IMO. We have a function that is (by name) a reader function,
	// but it feels wrong to IMPOSE read limitation on it this way. However, if the goal is to
	// have read queries work on a db read replica, then this should operate on the read db.
	congestionFee, err = s.feeCalculator.CalculateCongestionFee(
		s.ctx,
		s.db.ReadQuery(),
		messageTime,
		env.OriginatorNodeID(),
	)
	if err != nil {
		return 0, err
	}

	return baseFee + congestionFee, nil
}

// getPayerID resolves a payer address to its database ID.
// Don't ask the DB if you don't have to: the LRU cache short-circuits the
// FindOrCreatePayer write round-trip (~1 ms) for already-seen payers, which
// shaves roughly 1 s off a full batch backfill.
func (s *EnvelopeSink) getPayerID(env *envUtils.OriginatorEnvelope) (int32, error) {
	payerAddress, err := env.UnsignedOriginatorEnvelope.PayerEnvelope.RecoverSigner()
	if err != nil {
		return 0, err
	}

	hex := payerAddress.Hex()
	if id, ok := s.payerCache.Get(hex); ok {
		return id, nil
	}

	id, err := s.db.WriteQuery().FindOrCreatePayer(s.ctx, hex)
	if err != nil {
		return 0, err
	}

	s.payerCache.Add(hex, id)
	return id, nil
}
