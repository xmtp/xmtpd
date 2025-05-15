package storer

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"github.com/ethereum/go-ethereum/core/types"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	// We may not want to hardcode this to 0 and have an originator ID for each smart contract?
	GROUP_MESSAGE_ORIGINATOR_ID = 0
)

var (
	ErrParseGroupMessage               = "error parsing group message"
	ErrParseClientEnvelope             = "error parsing client envelope"
	ErrTopicDoesNotMatch               = "client envelope topic does not match payload topic"
	ErrBuildOriginatorEnvelope         = "error building originator envelope"
	ErrBuildSignedOriginatorEnvelope   = "error building signed originator envelope"
	ErrMarshallOriginatorEnvelope      = "error marshalling originator envelope"
	ErrInsertEnvelopeFromSmartContract = "error inserting envelope from smart contract"
	ErrInsertBlockchainMessage         = "error inserting blockchain message"
)

type GroupMessageStorer struct {
	contract *gm.GroupMessageBroadcaster
	queries  *queries.Queries
	logger   *zap.Logger
}

func NewGroupMessageStorer(
	queries *queries.Queries,
	logger *zap.Logger,
	contract *gm.GroupMessageBroadcaster,
) *GroupMessageStorer {
	return &GroupMessageStorer{
		queries:  queries,
		logger:   logger.Named("GroupMessageStorer"),
		contract: contract,
	}
}

// Validate and store a group message log event
func (s *GroupMessageStorer) StoreLog(
	ctx context.Context,
	event types.Log,
) LogStorageError {
	msgSent, err := s.contract.ParseMessageSent(event)
	if err != nil {
		return NewUnrecoverableLogStorageError(ErrParseGroupMessage, err)
	}

	topicStruct := topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, msgSent.GroupId[:])

	clientEnvelope, err := envelopes.NewClientEnvelopeFromBytes(msgSent.Message)
	if err != nil {
		s.logger.Error(ErrParseClientEnvelope, zap.Error(err))
		return NewUnrecoverableLogStorageError(ErrParseClientEnvelope, err)
	}

	targetTopic := clientEnvelope.TargetTopic()

	if !clientEnvelope.TopicMatchesPayload() {
		s.logger.Error(
			ErrTopicDoesNotMatch,
			zap.Any("targetTopic", targetTopic.String()),
			zap.Any("contractTopic", topicStruct.String()),
		)
		return NewUnrecoverableLogStorageError(
			ErrTopicDoesNotMatch,
			errors.New(ErrTopicDoesNotMatch),
		)
	}

	originatorEnvelope, err := buildOriginatorEnvelope(msgSent.SequenceId, msgSent.Message)
	if err != nil {
		s.logger.Error(ErrBuildOriginatorEnvelope, zap.Error(err))
		return NewUnrecoverableLogStorageError(ErrBuildOriginatorEnvelope, err)
	}

	signedOriginatorEnvelope, err := buildSignedOriginatorEnvelope(
		originatorEnvelope,
		event.TxHash,
	)
	if err != nil {
		s.logger.Error(ErrBuildSignedOriginatorEnvelope, zap.Error(err))
		return NewUnrecoverableLogStorageError(ErrBuildSignedOriginatorEnvelope, err)
	}

	originatorEnvelopeBytes, err := proto.Marshal(signedOriginatorEnvelope)
	if err != nil {
		s.logger.Error(ErrMarshallOriginatorEnvelope, zap.Error(err))
		return NewUnrecoverableLogStorageError(ErrMarshallOriginatorEnvelope, err)
	}

	s.logger.Info("Inserting message from contract", zap.String("topic", topicStruct.String()))

	if _, err = s.queries.InsertGatewayEnvelope(ctx, queries.InsertGatewayEnvelopeParams{
		OriginatorNodeID:     GROUP_MESSAGE_ORIGINATOR_ID,
		OriginatorSequenceID: int64(msgSent.SequenceId),
		Topic:                topicStruct.Bytes(),
		OriginatorEnvelope:   originatorEnvelopeBytes,
		Expiry:               sql.NullInt64{Int64: math.MaxInt64, Valid: true},
	}); err != nil {
		s.logger.Error(ErrInsertEnvelopeFromSmartContract, zap.Error(err))
		return NewRetryableLogStorageError(ErrInsertEnvelopeFromSmartContract, err)
	}

	if err = s.queries.InsertBlockchainMessage(ctx, queries.InsertBlockchainMessageParams{
		BlockNumber:          event.BlockNumber,
		BlockHash:            event.BlockHash.Bytes(),
		OriginatorNodeID:     GROUP_MESSAGE_ORIGINATOR_ID,
		OriginatorSequenceID: int64(msgSent.SequenceId),
		IsCanonical:          true, // New messages are always canonical
	}); err != nil {
		s.logger.Error(ErrInsertBlockchainMessage, zap.Error(err))
		return NewRetryableLogStorageError(ErrInsertBlockchainMessage, err)
	}

	return nil
}
