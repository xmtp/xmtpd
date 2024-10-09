package storer

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/xmtp/xmtpd/pkg/abis"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/topic"
	"go.uber.org/zap"
)

type GroupMessageStorer struct {
	contract *abis.GroupMessages
	queries  *queries.Queries
	logger   *zap.Logger
}

func NewGroupMessageStorer(
	queries *queries.Queries,
	logger *zap.Logger,
	contract *abis.GroupMessages,
) *GroupMessageStorer {
	return &GroupMessageStorer{
		queries:  queries,
		logger:   logger.Named("GroupMessageStorer"),
		contract: contract,
	}
}

// Validate and store a group message log event
func (s *GroupMessageStorer) StoreLog(ctx context.Context, event types.Log) LogStorageError {
	msgSent, err := s.contract.ParseMessageSent(event)
	if err != nil {
		return NewLogStorageError(err, false)
	}

	topicStruct := topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, msgSent.GroupId[:])

	s.logger.Debug("Inserting message from contract", zap.String("topic", topicStruct.String()))

	if _, err = s.queries.InsertGatewayEnvelope(ctx, queries.InsertGatewayEnvelopeParams{
		// We may not want to hardcode this to 0 and have an originator ID for each smart contract?
		OriginatorNodeID:     0,
		OriginatorSequenceID: int64(msgSent.SequenceId),
		Topic:                topicStruct.Bytes(),
		OriginatorEnvelope:   msgSent.Message, // TODO:nm parse originator envelope and do some validation
	}); err != nil {
		s.logger.Error("Error inserting envelope from smart contract", zap.Error(err))
		return NewLogStorageError(err, true)
	}

	return nil
}
