package sync

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	envUtils "github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/fees"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

type originatorStream struct {
	ctx                        context.Context
	db                         *sql.DB
	log                        *zap.Logger
	node                       *registry.Node
	queries                    *queries.Queries
	lastEnvelope               *envUtils.OriginatorEnvelope
	stream                     message_api.ReplicationApi_SubscribeEnvelopesClient
	feeCalculator              fees.IFeeCalculator
	payerReportStore           payerreport.IPayerReportStore
	payerReportDomainSeparator common.Hash
}

func newOriginatorStream(
	ctx context.Context,
	db *sql.DB,
	log *zap.Logger,
	node *registry.Node,
	lastEnvelope *envUtils.OriginatorEnvelope,
	stream message_api.ReplicationApi_SubscribeEnvelopesClient,
	feeCalculator fees.IFeeCalculator,
	payerReportStore payerreport.IPayerReportStore,
	payerReportDomainSeparator common.Hash,
) *originatorStream {
	return &originatorStream{
		ctx: ctx,
		db:  db,
		log: log.With(
			zap.Uint32("originator_id", node.NodeID),
			zap.String("http_address", node.HttpAddress),
		),
		node:                       node,
		queries:                    queries.New(db),
		lastEnvelope:               lastEnvelope,
		stream:                     stream,
		feeCalculator:              feeCalculator,
		payerReportStore:           payerReportStore,
		payerReportDomainSeparator: payerReportDomainSeparator,
	}
}

func (s *originatorStream) listen() error {
	recvChan := make(chan *message_api.SubscribeEnvelopesResponse)
	errChan := make(chan error)
	writeQueue := make(chan *envUtils.OriginatorEnvelope, 10)

	go func() {
		for {
			envs, err := s.stream.Recv()
			if err != nil {
				errChan <- err
				return
			}
			recvChan <- envs
		}
	}()

	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			case env, ok := <-writeQueue:
				if !ok {
					s.log.Error("writeQueue is closed")
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
	}()

	defer close(writeQueue)

	for {
		select {
		case <-s.ctx.Done():
			s.log.Info("Context canceled, stopping stream listener")
			return backoff.Permanent(s.ctx.Err())

		case envs, ok := <-recvChan:
			if !ok {
				s.log.Error("recvChan is closed")
				return backoff.Permanent(errors.New("recvChan is closed"))
			}

			if envs == nil || len(envs.Envelopes) == 0 {
				continue
			}
			s.log.Debug(
				"Received envelopes",
				zap.Any("numEnvelopes", len(envs.Envelopes)),
			)

			for _, env := range envs.Envelopes {
				// Any message that fails validation here will be dropped permanently
				parsedEnv, err := s.validateEnvelope(env)
				if err != nil {
					s.log.Error("discarding envelope after validation failed", zap.Error(err))
					continue
				}
				writeQueue <- parsedEnv
			}

		case err, ok := <-errChan:
			if !ok {
				s.log.Error("errChan is closed")
				return backoff.Permanent(errors.New("errChan is closed"))
			}

			if err == io.EOF {
				s.log.Info("Stream closed with EOF")
				// reset backoff to 1 second
				return backoff.RetryAfter(1)
			}
			s.log.Error(
				"Stream closed with error",
				zap.Error(err),
			)

			if strings.Contains(err.Error(), "is not compatible") {
				// the node won't accept our version
				// try again in an hour in case their config has changed
				return backoff.RetryAfter(3600)
			}

			// keep existing backoff
			return err
		}
	}
}

// validateEnvelope performs all static validation on an envelope
// if an error is encountered, the envelope will be dropped and the stream will continue
func (s *originatorStream) validateEnvelope(
	envProto *envelopes.OriginatorEnvelope,
) (*envUtils.OriginatorEnvelope, error) {
	var err error
	defer func() {
		if err != nil {
			metrics.EmitSyncOriginatorErrorMessages(s.node.NodeID, 1)
		}
	}()

	var env *envUtils.OriginatorEnvelope
	env, err = envUtils.NewOriginatorEnvelope(envProto)
	if err != nil {
		s.log.Error("Failed to unmarshal originator envelope", zap.Error(err))
		return nil, err
	}

	// TODO:(nm) Handle fetching envelopes from other nodes
	if env.OriginatorNodeID() != s.node.NodeID {
		s.log.Error("Received envelope from wrong node",
			zap.Any("nodeID", env.OriginatorNodeID()),
			zap.Any("expectedNodeId", s.node.NodeID),
		)
		err = errors.New("originator ID does not match envelope")
		return nil, err
	}

	metrics.EmitSyncLastSeenOriginatorSequenceId(env.OriginatorNodeID(), env.OriginatorSequenceID())
	metrics.EmitSyncOriginatorReceivedMessagesCount(env.OriginatorNodeID(), 1)

	var lastSequenceID uint64 = 0
	var lastNs int64 = 0
	if s.lastEnvelope != nil {
		lastSequenceID = s.lastEnvelope.OriginatorSequenceID()
		lastNs = s.lastEnvelope.OriginatorNs()
	}

	if env.OriginatorSequenceID() != lastSequenceID+1 || env.OriginatorNs() < lastNs {
		// TODO(rich) Submit misbehavior report and continue
		s.log.Error(
			"Received out-of-order envelope",
			zap.Uint64("expectedSequenceID", lastSequenceID+1),
			zap.Uint64("actualSequenceID", env.OriginatorSequenceID()),
			zap.Int64("lastTimestampNs", lastNs),
			zap.Int64("actualTimestampNs", env.OriginatorNs()),
			zap.Uint32("originatorId", env.OriginatorNodeID()),
		)
	}

	if env.OriginatorSequenceID() > lastSequenceID {
		s.lastEnvelope = env
	}

	// Validate that there is a valid payer signature
	_, err = env.UnsignedOriginatorEnvelope.PayerEnvelope.RecoverSigner()
	if err != nil {
		s.log.Error("Failed to recover payer address", zap.Error(err))
		return nil, err
	}

	return env, nil
}

func (s *originatorStream) storeEnvelope(env *envUtils.OriginatorEnvelope) error {
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

func (s *originatorStream) storeReservedEnvelope(env *envUtils.OriginatorEnvelope) error {
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

func (s *originatorStream) calculateFees(
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

func (s *originatorStream) getPayerID(env *envUtils.OriginatorEnvelope) (int32, error) {
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
