package storer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pingcap/log"
	"github.com/xmtp/xmtpd/contracts/pkg/identityupdates"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	// We may not want to hardcode this to 1 and have an originator ID for each smart contract?
	IDENTITY_UPDATE_ORIGINATOR_ID = 1
)

type IdentityUpdateStorer struct {
	contract          *identityupdates.IdentityUpdates
	db                *sql.DB
	logger            *zap.Logger
	validationService mlsvalidate.MLSValidationService
}

func NewIdentityUpdateStorer(
	db *sql.DB,
	logger *zap.Logger,
	contract *identityupdates.IdentityUpdates,
	validationService mlsvalidate.MLSValidationService,
) *IdentityUpdateStorer {
	return &IdentityUpdateStorer{
		db:                db,
		logger:            logger.Named("IdentityUpdateStorer"),
		contract:          contract,
		validationService: validationService,
	}
}

// Validate and store an identity update log event
func (s *IdentityUpdateStorer) StoreLog(
	ctx context.Context,
	event types.Log,
	appendLog bool,
) LogStorageError {
	msgSent, err := s.contract.ParseIdentityUpdateCreated(event)
	if err != nil {
		return NewLogStorageError(err, false)
	}

	err = db.RunInTx(
		ctx,
		s.db,
		&sql.TxOptions{Isolation: sql.LevelRepeatableRead},
		func(ctx context.Context, querier *queries.Queries) error {
			latestSequenceId, err := querier.GetLatestSequenceId(ctx, IDENTITY_UPDATE_ORIGINATOR_ID)
			if err != nil {
				return NewLogStorageError(err, true)
			}

			if uint64(latestSequenceId) >= msgSent.SequenceId {
				s.logger.Warn(
					"Duplicate identity update",
					zap.Uint64("latest_sequence_id", uint64(latestSequenceId)),
					zap.Uint64("msg_sequence_id", msgSent.SequenceId),
				)
				return nil
			}

			messageTopic := topic.NewTopic(topic.TOPIC_KIND_IDENTITY_UPDATES_V1, msgSent.InboxId[:])

			s.logger.Info(
				"Inserting identity update from contract",
				zap.String("topic", messageTopic.String()),
			)

			clientEnvelope, err := envelopes.NewClientEnvelopeFromBytes(msgSent.Update)
			if err != nil {
				s.logger.Error("Error parsing client envelope", zap.Error(err))
				return NewLogStorageError(err, false)
			}

			associationState, err := s.validateIdentityUpdate(
				ctx,
				querier,
				msgSent.InboxId,
				clientEnvelope,
			)
			if err != nil {
				log.Error("Error validating identity update", zap.Error(err))
				return NewLogStorageError(err, true)
			}

			inboxId := utils.HexEncode(msgSent.InboxId[:])

			for _, new_member := range associationState.StateDiff.NewMembers {
				s.logger.Info("New member", zap.Any("member", new_member))
				if address, ok := new_member.Kind.(*associations.MemberIdentifier_Address); ok {
					numRows, err := querier.InsertAddressLog(ctx, queries.InsertAddressLogParams{
						Address: address.Address,
						InboxID: inboxId,
						AssociationSequenceID: sql.NullInt64{
							Valid: true,
							Int64: int64(msgSent.SequenceId),
						},
					})
					if err != nil {
						return NewLogStorageError(err, true)
					}
					if numRows == 0 {
						s.logger.Warn(
							"Could not insert address log",
							zap.String("address", address.Address),
							zap.String("inbox_id", inboxId),
							zap.Int("sequence_id", int(msgSent.SequenceId)),
						)
					}
				}
			}

			for _, removed_member := range associationState.StateDiff.RemovedMembers {
				log.Info("Removed member", zap.Any("member", removed_member))
				if address, ok := removed_member.Kind.(*associations.MemberIdentifier_Address); ok {
					rows, err := querier.RevokeAddressFromLog(
						ctx,
						queries.RevokeAddressFromLogParams{
							Address: address.Address,
							InboxID: inboxId,
							RevocationSequenceID: sql.NullInt64{
								Valid: true,
								Int64: int64(msgSent.SequenceId),
							},
						},
					)
					if err != nil {
						return NewLogStorageError(err, true)
					}
					if rows == 0 {
						s.logger.Warn(
							"Could not find address log entry to revoke",
							zap.String("address", address.Address),
							zap.String("inbox_id", inboxId),
						)
					}
				}
			}

			signedOriginatorEnvelope, err := buildSignedOriginatorEnvelope(
				buildOriginatorEnvelope(msgSent.SequenceId, msgSent.Update),
				event.TxHash,
			)
			if err != nil {
				s.logger.Error("Error building signed originator envelope", zap.Error(err))
				return NewLogStorageError(err, true)
			}

			originatorEnvelopeBytes, err := proto.Marshal(signedOriginatorEnvelope)
			if err != nil {
				s.logger.Error("Error marshalling originator envelope", zap.Error(err))
				return NewLogStorageError(err, true)
			}

			version := sql.NullInt32{Int32: 1, Valid: true}

			if appendLog {
				version, err = GetVersionForAppend(
					ctx,
					querier,
					s.logger,
					IDENTITY_UPDATE_ORIGINATOR_ID,
					int64(msgSent.SequenceId),
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return NewLogStorageError(err, true)
					}
					if errors.Is(err, sql.ErrNoRows) {
						s.logger.Debug("No rows found for envelope, inserting new",
							zap.Int("originator_node_id", IDENTITY_UPDATE_ORIGINATOR_ID),
							zap.Int64("originator_sequence_id", int64(msgSent.SequenceId)),
						)
					}
				}
			}

			if _, err = querier.InsertGatewayEnvelope(ctx, queries.InsertGatewayEnvelopeParams{
				BlockNumber:          sql.NullInt64{Int64: int64(event.BlockNumber), Valid: true},
				BlockHash:            event.BlockHash.Bytes(),
				Version:              version,
				IsCanonical:          sql.NullBool{Bool: true, Valid: true},
				OriginatorNodeID:     IDENTITY_UPDATE_ORIGINATOR_ID,
				OriginatorSequenceID: int64(msgSent.SequenceId),
				Topic:                messageTopic.Bytes(),
				OriginatorEnvelope:   originatorEnvelopeBytes,
			}); err != nil {
				s.logger.Error("Error inserting envelope from smart contract", zap.Error(err))
				return NewLogStorageError(err, true)
			}

			return nil
		},
	)

	if err != nil {
		if logStorageErr, ok := err.(LogStorageError); ok {
			return logStorageErr
		}
		// If the error was not a LogStorageError we can assume it's a DB error and it should be retried
		return NewLogStorageError(err, true)
	}

	return nil
}

func (s *IdentityUpdateStorer) validateIdentityUpdate(
	ctx context.Context,
	querier *queries.Queries,
	inboxId [32]byte,
	clientEnvelope *envelopes.ClientEnvelope,
) (*mlsvalidate.AssociationStateResult, error) {
	gatewayEnvelopes, err := querier.SelectGatewayEnvelopes(
		ctx,
		queries.SelectGatewayEnvelopesParams{
			Topics:            []db.Topic{db.Topic(BuildInboxTopic(inboxId))},
			OriginatorNodeIds: []int32{IDENTITY_UPDATE_ORIGINATOR_ID},
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

	return s.validationService.GetAssociationStateFromEnvelopes(
		ctx,
		gatewayEnvelopes,
		identityUpdate.IdentityUpdate,
	)
}

func BuildInboxTopic(inboxId [32]byte) string {
	return fmt.Sprintf("i/%x", inboxId)
}

func buildOriginatorEnvelope(
	sequenceId uint64,
	clientEnvelopeBytes []byte,
) *envelopesProto.UnsignedOriginatorEnvelope {
	return &envelopesProto.UnsignedOriginatorEnvelope{
		OriginatorNodeId:     IDENTITY_UPDATE_ORIGINATOR_ID,
		OriginatorSequenceId: sequenceId,
		OriginatorNs:         time.Now().UnixNano(),
		PayerEnvelope: &envelopesProto.PayerEnvelope{
			UnsignedClientEnvelope: clientEnvelopeBytes,
		},
	}
}

func buildSignedOriginatorEnvelope(
	originatorEnvelope *envelopesProto.UnsignedOriginatorEnvelope,
	transactionHash common.Hash,
) (*envelopesProto.OriginatorEnvelope, error) {
	envelopeBytes, err := proto.Marshal(originatorEnvelope)
	if err != nil {
		return nil, err
	}

	return &envelopesProto.OriginatorEnvelope{
		UnsignedOriginatorEnvelope: envelopeBytes,
		Proof: &envelopesProto.OriginatorEnvelope_BlockchainProof{
			BlockchainProof: &envelopesProto.BlockchainProof{
				TransactionHash: transactionHash[:],
			},
		},
	}, nil
}
