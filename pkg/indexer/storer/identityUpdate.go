package storer

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pingcap/log"
	"github.com/xmtp/xmtpd/pkg/abis"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	IDENTITY_UPDATE_ORIGINATOR_ID = 1
)

type IdentityUpdateStorer struct {
	contract          *abis.IdentityUpdates
	queries           *queries.Queries
	logger            *zap.Logger
	validationService mlsvalidate.MLSValidationService
}

func NewIdentityUpdateStorer(
	queries *queries.Queries,
	logger *zap.Logger,
	contract *abis.IdentityUpdates,
	validationService mlsvalidate.MLSValidationService,
) *IdentityUpdateStorer {
	return &IdentityUpdateStorer{
		queries:           queries,
		logger:            logger.Named("IdentityUpdateStorer"),
		contract:          contract,
		validationService: validationService,
	}
}

// Validate and store an identity update log event
func (s *IdentityUpdateStorer) StoreLog(ctx context.Context, event types.Log) LogStorageError {
	msgSent, err := s.contract.ParseIdentityUpdateCreated(event)
	if err != nil {
		return NewLogStorageError(err, false)
	}

	// TODO:nm figure out topic structure
	topic := BuildInboxTopic(msgSent.InboxId)

	s.logger.Info("Inserting identity update from contract", zap.String("topic", topic))

	associationState, err := s.validateIdentityUpdate(ctx, msgSent.InboxId, msgSent.Update)
	if err != nil {
		log.Error("Error validating identity update", zap.Error(err), zap.Any("msgSent", msgSent))
		return NewLogStorageError(err, true)
	}

	inboxId := utils.HexEncode(msgSent.InboxId[:])

	for _, new_member := range associationState.StateDiff.NewMembers {
		s.logger.Info("New member", zap.Any("member", new_member))
		if address, ok := new_member.Kind.(*associations.MemberIdentifier_Address); ok {
			_, err = s.queries.InsertAddressLog(ctx, queries.InsertAddressLogParams{
				Address:               address.Address,
				InboxID:               inboxId,
				AssociationSequenceID: sql.NullInt64{Valid: true, Int64: int64(msgSent.SequenceId)},
				RevocationSequenceID:  sql.NullInt64{Valid: false},
			})
			if err != nil {
				return NewLogStorageError(err, true)
			}
		}
	}

	for _, removed_member := range associationState.StateDiff.RemovedMembers {
		log.Info("Removed member", zap.Any("member", removed_member))
		if address, ok := removed_member.Kind.(*associations.MemberIdentifier_Address); ok {
			err = s.queries.RevokeAddressFromLog(ctx, queries.RevokeAddressFromLogParams{
				Address:              address.Address,
				InboxID:              inboxId,
				RevocationSequenceID: sql.NullInt64{Valid: true, Int64: int64(msgSent.SequenceId)},
			})
			if err != nil {
				return NewLogStorageError(err, true)
			}
		}
	}

	originatorEnvelope, err := buildOriginatorEnvelope(topic, msgSent.SequenceId, msgSent.Update)
	if err != nil {
		s.logger.Error("Error building originator envelope", zap.Error(err))
		return NewLogStorageError(err, true)
	}

	originatorEnvelopeBytes, err := proto.Marshal(originatorEnvelope)
	if err != nil {
		s.logger.Error("Error marshalling originator envelope", zap.Error(err))
		return NewLogStorageError(err, true)
	}

	if _, err = s.queries.InsertGatewayEnvelope(ctx, queries.InsertGatewayEnvelopeParams{
		// We may not want to hardcode this to 1 and have an originator ID for each smart contract?
		OriginatorNodeID:     IDENTITY_UPDATE_ORIGINATOR_ID,
		OriginatorSequenceID: int64(msgSent.SequenceId),
		Topic:                []byte(topic),
		OriginatorEnvelope:   originatorEnvelopeBytes,
	}); err != nil {
		s.logger.Error("Error inserting envelope from smart contract", zap.Error(err))
		return NewLogStorageError(err, true)
	}

	return nil
}

func (s *IdentityUpdateStorer) validateIdentityUpdate(
	ctx context.Context,
	inboxId [32]byte,
	update []byte,
) (*mlsvalidate.AssociationStateResult, error) {
	gatewayEnvelopes, err := s.queries.SelectGatewayEnvelopes(
		ctx,
		queries.SelectGatewayEnvelopesParams{
			Topic:            []byte(BuildInboxTopic(inboxId)),
			OriginatorNodeID: sql.NullInt32{Int32: IDENTITY_UPDATE_ORIGINATOR_ID, Valid: true},
			RowLimit:         sql.NullInt32{Int32: 256, Valid: true},
		},
	)
	if err != nil {
		return nil, err
	}

	return s.validationService.GetAssociationStateFromEnvelopes(ctx, gatewayEnvelopes, update)
}

func BuildInboxTopic(inboxId [32]byte) string {
	return fmt.Sprintf("i/%x", inboxId)
}

func buildOriginatorEnvelope(
	topic string,
	sequenceId uint64,
	update []byte,
) (*message_api.UnsignedOriginatorEnvelope, error) {
	clientEnvelopeBytes, err := proto.Marshal(buildClientEnvelope(topic, update))
	if err != nil {
		return nil, err
	}
	return &message_api.UnsignedOriginatorEnvelope{
		OriginatorNodeId:     IDENTITY_UPDATE_ORIGINATOR_ID,
		OriginatorSequenceId: sequenceId,
		OriginatorNs:         time.Now().UnixNano(),
		PayerEnvelope: &message_api.PayerEnvelope{
			UnsignedClientEnvelope: clientEnvelopeBytes,
		},
	}, nil
}

func buildClientEnvelope(topic string, update []byte) *message_api.ClientEnvelope {
	var identityUpdate *associations.IdentityUpdate
	proto.Unmarshal(update, identityUpdate)
	return &message_api.ClientEnvelope{
		Aad: nil,
		Payload: &message_api.ClientEnvelope_IdentityUpdate{
			IdentityUpdate: identityUpdate,
		},
	}
}
