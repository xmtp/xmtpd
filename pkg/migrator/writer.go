package migrator

import (
	"context"
	"database/sql"
	"time"

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	re "github.com/xmtp/xmtpd/pkg/errors"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var (
	ErrInsertOriginatorEnvelope   = "insert originator envelope failed"
	ErrMarshallOriginatorEnvelope = "marshall originator envelope failed"
)

func (s *dbMigrator) insertOriginatorEnvelope(
	env *envelopes.OriginatorEnvelope,
	querier *queries.Queries,
) re.RetryableError {
	originatorEnvelopeBytes, err := proto.Marshal(env.Proto())
	if err != nil {
		s.log.Error(ErrMarshallOriginatorEnvelope, zap.Error(err))
		return re.NewNonRecoverableError(ErrMarshallOriginatorEnvelope, err)
	}

	// TODO: Derive AddressLog from IdentityUpdate.

	// TODO: Get AssociationStates from originator envelope.

	payerAddress, err := env.UnsignedOriginatorEnvelope.PayerEnvelope.RecoverSigner()
	if err != nil {
		s.log.Error("failed to recover payer address", zap.Error(err))
		return re.NewNonRecoverableError("failed to recover payer address", err)
	}

	payerID, err := querier.FindOrCreatePayer(s.ctx, payerAddress.Hex())
	if err != nil {
		s.log.Error("failed to find or create payer", zap.Error(err))
		return re.NewNonRecoverableError("failed to find or create payer", err)
	}

	_, err = db.InsertGatewayEnvelopeAndIncrementUnsettledUsage(
		s.ctx,
		s.writer,
		queries.InsertGatewayEnvelopeParams{
			OriginatorNodeID:     int32(env.OriginatorNodeID()),
			OriginatorSequenceID: int64(env.OriginatorSequenceID()),
			Topic:                env.TargetTopic().Bytes(),
			OriginatorEnvelope:   originatorEnvelopeBytes,
			PayerID:              db.NullInt32(payerID),
			GatewayTime:          env.OriginatorTime(),
			Expiry: sql.NullInt64{
				Int64: int64(env.UnsignedOriginatorEnvelope.Proto().GetExpiryUnixtime()),
				Valid: true,
			},
		},
		queries.IncrementUnsettledUsageParams{
			PayerID:           payerID,
			OriginatorID:      int32(env.OriginatorNodeID()),
			MinutesSinceEpoch: utils.MinutesSinceEpoch(env.OriginatorTime()),
			SpendPicodollars: int64(env.UnsignedOriginatorEnvelope.BaseFee()) +
				int64(env.UnsignedOriginatorEnvelope.CongestionFee()),
		},
	)
	if err != nil {
		s.log.Error(ErrInsertOriginatorEnvelope, zap.Error(err))
		return re.NewRecoverableError(ErrInsertOriginatorEnvelope, err)
	}

	return nil
}

func retry(
	ctx context.Context,
	logger *zap.Logger,
	sleep time.Duration,
	fn func() re.RetryableError,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			if err := fn(); err != nil {
				logger.Error("error storing log", zap.Error(err))

				if err.ShouldRetry() {
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
