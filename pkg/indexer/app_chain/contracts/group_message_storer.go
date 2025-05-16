package contracts

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"github.com/ethereum/go-ethereum/core/types"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	re "github.com/xmtp/xmtpd/pkg/errors"
	"github.com/xmtp/xmtpd/pkg/topic"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	// We may not want to hardcode this to 0 and have an originator ID for each smart contract?
	GROUP_MESSAGE_ORIGINATOR_ID = 0
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
		logger:   logger.Named("storer"),
		contract: contract,
	}
}

// Validate and store a group message log event
func (s *GroupMessageStorer) StoreLog(
	ctx context.Context,
	event types.Log,
) re.RetryableError {
	msgSent, err := s.contract.ParseMessageSent(event)
	if err != nil {
		return re.NewUnrecoverableLogStorageError(err)
	}

	topicStruct := topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, msgSent.GroupId[:])

	clientEnvelope, err := envelopes.NewClientEnvelopeFromBytes(msgSent.Message)
	if err != nil {
		s.logger.Error("Error parsing client envelope", zap.Error(err))
		return re.NewUnrecoverableLogStorageError(err)
	}

	targetTopic := clientEnvelope.TargetTopic()

	if !clientEnvelope.TopicMatchesPayload() {
		s.logger.Error(
			"Client envelope topic does not match payload type",
			zap.Any("targetTopic", targetTopic.String()),
			zap.Any("contractTopic", topicStruct.String()),
		)
		return re.NewUnrecoverableLogStorageError(
			errors.New("client envelope topic does not match payload topic"),
		)
	}

	originatorEnvelope, err := buildOriginatorEnvelope(msgSent.SequenceId, msgSent.Message)
	if err != nil {
		s.logger.Error("Error building originator envelope", zap.Error(err))
		return re.NewUnrecoverableLogStorageError(err)
	}

	signedOriginatorEnvelope, err := buildSignedOriginatorEnvelope(
		originatorEnvelope,
		event.TxHash,
	)
	if err != nil {
		s.logger.Error("Error building signed originator envelope", zap.Error(err))
		return re.NewUnrecoverableLogStorageError(err)
	}

	originatorEnvelopeBytes, err := proto.Marshal(signedOriginatorEnvelope)
	if err != nil {
		s.logger.Error("Error marshalling originator envelope", zap.Error(err))
		return re.NewUnrecoverableLogStorageError(err)
	}

	s.logger.Info("Inserting message from contract", zap.String("topic", topicStruct.String()))

	if _, err = s.queries.InsertGatewayEnvelope(ctx, queries.InsertGatewayEnvelopeParams{
		OriginatorNodeID:     GROUP_MESSAGE_ORIGINATOR_ID,
		OriginatorSequenceID: int64(msgSent.SequenceId),
		Topic:                topicStruct.Bytes(),
		OriginatorEnvelope:   originatorEnvelopeBytes,
		Expiry:               sql.NullInt64{Int64: math.MaxInt64, Valid: true},
	}); err != nil {
		s.logger.Error("Error inserting envelope from smart contract", zap.Error(err))
		return re.NewRetryableLogStorageError(err)
	}

	if err = s.queries.InsertBlockchainMessage(ctx, queries.InsertBlockchainMessageParams{
		BlockNumber:          event.BlockNumber,
		BlockHash:            event.BlockHash.Bytes(),
		OriginatorNodeID:     GROUP_MESSAGE_ORIGINATOR_ID,
		OriginatorSequenceID: int64(msgSent.SequenceId),
		IsCanonical:          true, // New messages are always canonical
	}); err != nil {
		s.logger.Error("Error inserting blockchain message", zap.Error(err))
		return re.NewRetryableLogStorageError(err)
	}

	return nil
}
