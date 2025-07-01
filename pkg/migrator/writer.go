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
	re "github.com/xmtp/xmtpd/pkg/errors"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (m *Migrator) insertOriginatorEnvelope(
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

	originatorEnvelopeBytes, err := proto.Marshal(env.Proto())
	if err != nil {
		m.log.Error("marshall originator envelope failed", zap.Error(err))
		return re.NewNonRecoverableError("marshall originator envelope failed", err)
	}

	err = db.RunInTx(
		m.ctx,
		m.writer,
		&sql.TxOptions{Isolation: sql.LevelRepeatableRead},
		func(ctx context.Context, querier *queries.Queries) error {
			// When handling identity updates, we need to derive the address log updates from the association state.
			if env.OriginatorNodeID() == InboxLogOriginatorID {
				inboxIDBytes, inboxIDStr, err := getInboxID(
					&env.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope,
				)
				if err != nil {
					return re.NewNonRecoverableError("get inbox ID failed", err)
				}

				associationState, err := m.validateIdentityUpdate(
					ctx,
					querier,
					inboxIDBytes,
					&env.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope,
				)
				if err != nil {
					return re.NewNonRecoverableError("validate identity update failed", err)
				}

				for _, newMember := range associationState.StateDiff.NewMembers {
					m.log.Info("New member", zap.Any("member", newMember))

					if address, ok := newMember.Kind.(*associations.MemberIdentifier_EthereumAddress); ok {
						numRows, err := querier.InsertAddressLog(
							ctx,
							queries.InsertAddressLogParams{
								Address: address.EthereumAddress,
								InboxID: inboxIDStr,
								AssociationSequenceID: sql.NullInt64{
									Valid: true,
									Int64: int64(env.OriginatorSequenceID()),
								},
							},
						)
						if err != nil {
							return re.NewRecoverableError("insert address log failed", err)
						}

						if numRows == 0 {
							m.log.Warn(
								"Could not insert address log",
								zap.String("address", address.EthereumAddress),
								zap.String("inboxID", inboxIDStr),
								zap.Int("sequenceID", int(env.OriginatorSequenceID())),
							)
						}
					}
				}

				for _, removedMember := range associationState.StateDiff.RemovedMembers {
					m.log.Info("Removed member", zap.Any("member", removedMember))

					if address, ok := removedMember.Kind.(*associations.MemberIdentifier_EthereumAddress); ok {
						rows, err := querier.RevokeAddressFromLog(
							ctx,
							queries.RevokeAddressFromLogParams{
								Address: address.EthereumAddress,
								InboxID: inboxIDStr,
								RevocationSequenceID: sql.NullInt64{
									Valid: true,
									Int64: int64(env.OriginatorSequenceID()),
								},
							},
						)
						if err != nil {
							return re.NewRecoverableError("revoke address from log failed", err)
						}

						if rows == 0 {
							m.log.Warn(
								"Could not find address log entry to revoke",
								zap.String("address", address.EthereumAddress),
								zap.String("inboxID", inboxIDStr),
							)
						}
					}
				}
			}

			payerID, err := querier.FindOrCreatePayer(m.ctx, payerAddress.Hex())
			if err != nil {
				m.log.Error("find or create payer failed", zap.Error(err))
				return re.NewNonRecoverableError("find or create payer failed", err)
			}

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

func getInboxID(clientEnvelope *envelopes.ClientEnvelope) ([32]byte, string, error) {
	identityUpdatePayload, ok := clientEnvelope.Payload().(*envelopesProto.ClientEnvelope_IdentityUpdate)
	if !ok {
		return [32]byte{}, "", fmt.Errorf("client envelope payload is not an identity update")
	}

	inboxIDStr := identityUpdatePayload.IdentityUpdate.GetInboxId()

	// TODO: Is this the correct way of getting the inboxID?
	// Do we need bit padding or anything else?
	inboxIDBytes, err := utils.HexDecode(inboxIDStr)
	if err != nil {
		return [32]byte{}, "", fmt.Errorf("invalid inbox ID format: %w", err)
	}

	if len(inboxIDBytes) != 32 {
		return [32]byte{}, "", fmt.Errorf("inbox ID must be 32 bytes")
	}

	var inboxID [32]byte
	copy(inboxID[:], inboxIDBytes)

	return inboxID, inboxIDStr, nil
}

func (m *Migrator) validateIdentityUpdate(
	ctx context.Context,
	querier *queries.Queries,
	inboxID [32]byte,
	clientEnvelope *envelopes.ClientEnvelope,
) (*mlsvalidate.AssociationStateResult, error) {
	// Select gateway envelopes from the originator ID 12.
	// Do not mix with blockchain identity updates.
	gatewayEnvelopes, err := querier.SelectGatewayEnvelopes(
		ctx,
		queries.SelectGatewayEnvelopesParams{
			Topics: []db.Topic{
				db.Topic(topic.NewTopic(topic.TOPIC_KIND_IDENTITY_UPDATES_V1, inboxID[:]).Bytes()),
			},
			OriginatorNodeIds: []int32{int32(InboxLogOriginatorID)},
			RowLimit:          256,
		},
	)
	if err != nil {
		return nil, err
	}

	identityUpdate, ok := clientEnvelope.Payload().(*envelopesProto.ClientEnvelope_IdentityUpdate)
	if !ok {
		return nil, fmt.Errorf("client envelope payload is not an identity update")
	}

	return m.validationService.GetAssociationStateFromEnvelopes(
		ctx,
		gatewayEnvelopes,
		identityUpdate.IdentityUpdate,
	)
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
