package sync

import (
	"context"
	"database/sql"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	envUtils "github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/fees"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

type EnvelopeSink struct {
	ctx                        context.Context
	db                         *sql.DB
	logger                     *zap.Logger
	queries                    *queries.Queries
	feeCalculator              fees.IFeeCalculator
	payerReportStore           payerreport.IPayerReportStore
	payerReportDomainSeparator common.Hash
	writeQueue                 chan *envUtils.OriginatorEnvelope
	errorRetrySleepTime        time.Duration
}

func newEnvelopeSink(
	ctx context.Context,
	db *sql.DB,
	logger *zap.Logger,
	feeCalculator fees.IFeeCalculator,
	payerReportStore payerreport.IPayerReportStore,
	payerReportDomainSeparator common.Hash,
	writeQueue chan *envUtils.OriginatorEnvelope,
	errorRetrySleepTime time.Duration,
) *EnvelopeSink {
	return &EnvelopeSink{
		ctx:                        ctx,
		db:                         db,
		logger:                     logger,
		queries:                    queries.New(db),
		feeCalculator:              feeCalculator,
		payerReportStore:           payerReportStore,
		payerReportDomainSeparator: payerReportDomainSeparator,
		writeQueue:                 writeQueue,
		errorRetrySleepTime:        errorRetrySleepTime,
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

		storeLoop:
			for {
				select {
				case <-s.ctx.Done():
					return
				default:
					err := s.storeEnvelope(env)
					if err != nil {
						s.logger.Error("error storing envelope", zap.Error(err))
						time.Sleep(s.errorRetrySleepTime)
						continue
					}
					break storeLoop
				}
			}
		}
	}
}

func (s *EnvelopeSink) storeEnvelope(env *envUtils.OriginatorEnvelope) error {
	if env.TargetTopic().IsReserved() {
		s.logger.Info(
			"found envelope with reserved topic",
			utils.TopicField(env.TargetTopic().String()),
		)
		return s.storeReservedEnvelope(env)
	}

	// Calculate the fees independently to verify the originator's calculation
	ourFeeCalculation, err := s.calculateFees(env)
	if err != nil {
		s.logger.Error("failed to calculate fees", zap.Error(err))
		return err
	}
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

	inserted, err := db.InsertGatewayEnvelopeAndIncrementUnsettledUsage(
		s.ctx,
		s.db,
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
	)

	if err != nil {
		s.logger.Error("failed to insert gateway envelope", zap.Error(err))
		return err
	} else if inserted == 0 {
		// Envelope was already inserted by another worker
		s.logger.Debug("envelope already inserted",
			utils.OriginatorIDField(env.OriginatorNodeID()),
			utils.SequenceIDField(int64(env.OriginatorSequenceID())),
		)

		return nil
	}

	return nil
}

func (s *EnvelopeSink) storeReservedEnvelope(env *envUtils.OriginatorEnvelope) error {
	payerID, err := s.getPayerID(env)
	if err != nil {
		s.logger.Error("failed to get payer ID", zap.Error(err))
		return err
	}

	switch env.TargetTopic().Kind() {
	case topic.TopicKindPayerReportsV1:
		err := s.payerReportStore.StoreSyncedReport(
			s.ctx,
			env,
			payerID,
			s.payerReportDomainSeparator,
		)
		if err != nil {
			s.logger.Error("failed to store synced report", zap.Error(err))
			// Return nil here to avoid infinite retries
		}
		return nil
	case topic.TopicKindPayerReportAttestationsV1:
		err := s.payerReportStore.StoreSyncedAttestation(
			s.ctx,
			env,
			payerID,
		)
		if err != nil {
			s.logger.Error("failed to store synced attestation", zap.Error(err))
			// Return nil here to avoid infinite retries
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
	payerEnvelopeLength := len(env.UnsignedOriginatorEnvelope.PayerEnvelopeBytes())
	messageTime := utils.NsToDate(env.OriginatorNs())

	baseFee, err := s.feeCalculator.CalculateBaseFee(
		messageTime,
		int64(payerEnvelopeLength),
		env.UnsignedOriginatorEnvelope.PayerEnvelope.RetentionDays(),
	)
	if err != nil {
		return 0, err
	}

	congestionFee, err := s.feeCalculator.CalculateCongestionFee(
		s.ctx,
		s.queries,
		messageTime,
		env.OriginatorNodeID(),
	)
	if err != nil {
		return 0, err
	}

	return baseFee + congestionFee, nil
}

func (s *EnvelopeSink) getPayerID(env *envUtils.OriginatorEnvelope) (int32, error) {
	payerAddress, err := env.UnsignedOriginatorEnvelope.PayerEnvelope.RecoverSigner()
	if err != nil {
		return 0, err
	}

	payerID, err := s.queries.FindOrCreatePayer(s.ctx, payerAddress.Hex())
	if err != nil {
		return 0, err
	}

	return payerID, nil
}
