package migrator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/utils"
	re "github.com/xmtp/xmtpd/pkg/utils/retryerrors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (m *Migrator) insertOriginatorEnvelopeDatabase(
	env *envelopes.OriginatorEnvelope,
) re.RetryableError {
	if env == nil {
		return re.NewNonRecoverableError("", errors.New("envelope is nil"))
	}

	tableName, ok := originatorIDToTableName[env.OriginatorNodeID()]
	if !ok {
		return re.NewNonRecoverableError("", errors.New("invalid originator id"))
	}

	payerAddress, err := env.UnsignedOriginatorEnvelope.PayerEnvelope.RecoverSigner()
	if err != nil {
		m.log.Error("recover payer address failed", zap.Error(err))
		return re.NewNonRecoverableError("recover payer address failed", err)
	}

	querier := queries.New(m.writer)
	payerID, err := querier.FindOrCreatePayer(m.ctx, payerAddress.Hex())
	if err != nil {
		m.log.Error("find or create payer failed", zap.Error(err))
		return re.NewRecoverableError("find or create payer failed", err)
	}

	originatorEnvelopeBytes, err := proto.Marshal(env.Proto())
	if err != nil {
		m.log.Error("marshall originator envelope failed", zap.Error(err))
		return re.NewNonRecoverableError("marshall originator envelope failed", err)
	}

	err = db.RunInTx(
		m.ctx,
		m.writer,
		nil,
		func(ctx context.Context, querier *queries.Queries) error {
			_, err = querier.InsertGatewayEnvelope(ctx, queries.InsertGatewayEnvelopeParams{
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
			})
			if err != nil {
				m.log.Error("insert originator envelope failed", zap.Error(err))
				return re.NewRecoverableError("insert originator envelope failed", err)
			}

			err = querier.IncrementUnsettledUsage(ctx, queries.IncrementUnsettledUsageParams{
				PayerID:           payerID,
				OriginatorID:      int32(env.OriginatorNodeID()),
				MinutesSinceEpoch: utils.MinutesSinceEpoch(env.OriginatorTime()),
				SpendPicodollars: int64(env.UnsignedOriginatorEnvelope.BaseFee()) +
					int64(env.UnsignedOriginatorEnvelope.CongestionFee()),
				SequenceID: int64(env.OriginatorSequenceID()),
			})
			if err != nil {
				m.log.Error("increment unsettled usage failed", zap.Error(err))
				return re.NewRecoverableError("increment unsettled usage failed", err)
			}

			err = querier.UpdateMigrationProgress(ctx, queries.UpdateMigrationProgressParams{
				LastMigratedID: int64(env.OriginatorSequenceID()),
				SourceTable:    tableName,
			})
			if err != nil {
				m.log.Error("update migration progress failed", zap.Error(err))
				return re.NewRecoverableError("update migration progress failed", err)
			}

			return nil
		})
	if err != nil {
		var retryableError re.RetryableError
		if errors.As(err, &retryableError) {
			return retryableError
		}

		return re.NewRecoverableError("db error", err)
	}

	return nil
}

func (m *Migrator) insertOriginatorEnvelopeBlockchain(
	env *envelopes.OriginatorEnvelope,
) error {
	var (
		identifier = env.TargetTopic().Identifier()
		sequenceID = env.OriginatorSequenceID()
	)

	tableName, ok := originatorIDToTableName[env.OriginatorNodeID()]
	if !ok {
		return fmt.Errorf("invalid originator id: %d", env.OriginatorNodeID())
	}

	clientEnvelopeBytes, err := env.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope.Bytes()
	if err != nil {
		m.log.Error("failed to get payer envelope bytes", zap.Error(err))
		return fmt.Errorf("failed to get payer envelope bytes: %w", err)
	}

	querier := queries.New(m.writer)

	switch env.OriginatorNodeID() {
	case GroupMessageOriginatorID:
		groupID, err := utils.ParseGroupID(identifier)
		if err != nil {
			return fmt.Errorf("error converting identifier to group ID: %w", err)
		}

		log, err := m.blockchainPublisher.BootstrapGroupMessages(
			m.ctx,
			[][16]byte{groupID},
			[][]byte{clientEnvelopeBytes},
			[]uint64{sequenceID},
		)
		if err != nil {
			return fmt.Errorf(
				"error publishing group message with sequence ID %d: %w",
				sequenceID,
				err,
			)
		}

		if len(log) == 0 {
			return fmt.Errorf(
				"received nil log publishing group message with sequence ID %d",
				sequenceID,
			)
		}

		err = querier.UpdateMigrationProgress(m.ctx, queries.UpdateMigrationProgressParams{
			LastMigratedID: int64(env.OriginatorSequenceID()),
			SourceTable:    tableName,
		})
		if err != nil {
			m.log.Error("update migration progress failed", zap.Error(err))
			return fmt.Errorf("update migration progress failed: %w", err)
		}

		m.log.Debug(
			"published group message",
			zap.String("group_id", utils.HexEncode(groupID[:])),
			zap.Uint64("sequence_id", sequenceID),
		)

	case InboxLogOriginatorID:
		inboxID, err := utils.ParseInboxID(identifier)
		if err != nil {
			return fmt.Errorf("error converting identifier to inbox ID: %w", err)
		}

		m.log.Debug(
			"publishing identity update",
			zap.String("inbox_id", utils.HexEncode(inboxID[:])),
			zap.Uint64("sequence_id", sequenceID),
		)

		log, err := m.blockchainPublisher.BootstrapIdentityUpdates(
			m.ctx,
			[][32]byte{inboxID},
			[][]byte{clientEnvelopeBytes},
			[]uint64{sequenceID},
		)
		if err != nil {
			return fmt.Errorf(
				"error publishing identity update with sequence ID %d: %w",
				sequenceID,
				err,
			)
		}

		if len(log) == 0 {
			return fmt.Errorf(
				"received nil log publishing identity update with sequence ID %d",
				sequenceID,
			)
		}

		err = querier.UpdateMigrationProgress(m.ctx, queries.UpdateMigrationProgressParams{
			LastMigratedID: int64(env.OriginatorSequenceID()),
			SourceTable:    tableName,
		})
		if err != nil {
			m.log.Error("update migration progress failed", zap.Error(err))
			return fmt.Errorf("update migration progress failed: %w", err)
		}

		m.log.Debug(
			"published identity update",
			zap.String("inbox_id", utils.HexEncode(inboxID[:])),
			zap.Uint64("sequence_id", sequenceID),
		)
	}

	return nil
}

func retry(
	ctx context.Context,
	logger *zap.Logger,
	sleep time.Duration,
	tableName string,
	destination string,
	fn func() re.RetryableError,
) error {
	attempts := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			if err := fn(); err != nil {
				logger.Error("error storing log", zap.Error(err))

				if err.ShouldRetry() {
					attempts++
					metrics.EmitMigratorWriterRetryAttempts(tableName, destination, attempts)
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
