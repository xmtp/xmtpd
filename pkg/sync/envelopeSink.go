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
	log                        *zap.Logger
	queries                    *queries.Queries
	feeCalculator              fees.IFeeCalculator
	payerReportStore           payerreport.IPayerReportStore
	payerReportDomainSeparator common.Hash
	writeQueue                 chan *envUtils.OriginatorEnvelope
}

func newEnvelopeSink(
	ctx context.Context,
	db *sql.DB,
	log *zap.Logger,
	feeCalculator fees.IFeeCalculator,
	payerReportStore payerreport.IPayerReportStore,
	payerReportDomainSeparator common.Hash,
	writeQueue chan *envUtils.OriginatorEnvelope,
) *EnvelopeSink {
	return &EnvelopeSink{
		ctx:                        ctx,
		db:                         db,
		log:                        log,
		queries:                    queries.New(db),
		feeCalculator:              feeCalculator,
		payerReportStore:           payerReportStore,
		payerReportDomainSeparator: payerReportDomainSeparator,
		writeQueue:                 writeQueue,
	}
}

func (s *EnvelopeSink) Start() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case env, ok := <-s.writeQueue:
			if !ok {
				s.log.Debug("writeQueue is closed")
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
						s.log.Error("error storing envelope", zap.Error(err))
						time.Sleep(1 * time.Second)
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
		s.log.Info(
			"Found envelope with reserved topic",
			zap.String("topic", env.TargetTopic().String()),
		)
		return s.storeReservedEnvelope(env)
	}

	// Calculate the fees independently to verify the originator's calculation
	ourFeeCalculation, err := s.calculateFees(env)
	if err != nil {
		s.log.Error("Failed to calculate fees", zap.Error(err))
		return err
	}
	originatorsFeeCalculation := env.UnsignedOriginatorEnvelope.BaseFee() +
		env.UnsignedOriginatorEnvelope.CongestionFee()

	if ourFeeCalculation != originatorsFeeCalculation {
		s.log.Warn(
			"Fee calculation mismatch",
			zap.Any("ourFee", ourFeeCalculation),
			zap.Any("originatorsFee", originatorsFeeCalculation),
		)
	}

	// If for some reason the envelope is not able to marshal (but has made it this far)
	// the node will retry indefinitely.
	// I don't think this will ever happen
	originatorBytes, err := env.Bytes()
	if err != nil {
		s.log.Error("Failed to marshal originator envelope", zap.Error(err))
		return err
	}

	// The payer address has already been validated, so any errors here should be transient
	payerId, err := s.getPayerID(env)
	if err != nil {
		s.log.Error("Failed to get payer ID", zap.Error(err))
		return err
	}

	originatorID := int32(env.OriginatorNodeID())
	originatorTime := utils.NsToDate(env.OriginatorNs())
	expiry := env.UnsignedOriginatorEnvelope.Proto().GetExpiryUnixtime()
	var expiryToSave sql.NullInt64
	if expiry > 0 {
		expiryToSave = db.NullInt64(int64(expiry))
	}

	inserted, err := db.InsertGatewayEnvelopeAndIncrementUnsettledUsage(
		s.ctx,
		s.db,
		queries.InsertGatewayEnvelopeParams{
			OriginatorNodeID:     int32(env.OriginatorNodeID()),
			OriginatorSequenceID: int64(env.OriginatorSequenceID()),
			Topic:                env.TargetTopic().Bytes(),
			OriginatorEnvelope:   originatorBytes,
			PayerID:              db.NullInt32(payerId),
			Expiry:               expiryToSave,
		},
		queries.IncrementUnsettledUsageParams{
			PayerID:           payerId,
			OriginatorID:      originatorID,
			MinutesSinceEpoch: utils.MinutesSinceEpoch(originatorTime),
			SpendPicodollars:  int64(ourFeeCalculation),
			MessageCount:      1,
		},
	)

	if err != nil {
		s.log.Error("Failed to insert gateway envelope", zap.Error(err))
		return err
	} else if inserted == 0 {
		// Envelope was already inserted by another worker
		s.log.Debug("Envelope already inserted", zap.Uint32("originatorID", env.OriginatorNodeID()), zap.Uint64("sequenceID", env.OriginatorSequenceID()))
		return nil
	}

	return nil
}

func (s *EnvelopeSink) storeReservedEnvelope(env *envUtils.OriginatorEnvelope) error {
	payerID, err := s.getPayerID(env)
	if err != nil {
		s.log.Error("Failed to get payer ID", zap.Error(err))
		return err
	}

	switch env.TargetTopic().Kind() {
	case topic.TOPIC_KIND_PAYER_REPORTS_V1:
		err := s.payerReportStore.StoreSyncedReport(
			s.ctx,
			env,
			payerID,
			s.payerReportDomainSeparator,
		)
		if err != nil {
			s.log.Error("Failed to store synced report", zap.Error(err))
			// Return nil here to avoid infinite retries
		}
		return nil
	case topic.TOPIC_KIND_PAYER_REPORT_ATTESTATIONS_V1:
		err := s.payerReportStore.StoreSyncedAttestation(
			s.ctx,
			env,
			payerID,
		)
		if err != nil {
			s.log.Error("Failed to store synced attestation", zap.Error(err))
			// Return nil here to avoid infinite retries
		}
		return nil
	default:
		s.log.Info(
			"Received unknown reserved topic",
			zap.String("topic", env.TargetTopic().String()),
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

	payerId, err := s.queries.FindOrCreatePayer(s.ctx, payerAddress.Hex())
	if err != nil {
		return 0, err
	}

	return payerId, nil
}
