package storer

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/xmtp/xmtpd/pkg/abis"
	"github.com/xmtp/xmtpd/pkg/db/queries"
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
	return &GroupMessageStorer{queries: queries, logger: logger, contract: contract}
}

// Validate and store a group message log event
func (s *GroupMessageStorer) StoreLog(ctx context.Context, event types.Log) LogStorageError {
	msgSent, err := s.contract.ParseMessageSent(event)
	if err != nil {
		return NewLogStorageError(err, false)
	}

	// TODO:nm figure out topic structure
	topic := BuildGroupMessageTopic(msgSent.GroupId)

	s.logger.Debug("Inserting message from contract", zap.String("topic", topic))

	if _, err = s.queries.InsertGatewayEnvelope(ctx, queries.InsertGatewayEnvelopeParams{
		// We may not want to hardcode this to 0 and have an originator ID for each smart contract?
		OriginatorID:         0,
		OriginatorSequenceID: int64(msgSent.SequenceId),
		Topic:                []byte(topic),
		OriginatorEnvelope:   msgSent.Message, // TODO:nm parse originator envelope and do some validation
	}); err != nil {
		s.logger.Error("Error inserting envelope from smart contract", zap.Error(err))
		return NewLogStorageError(err, true)
	}

	return nil
}

func BuildGroupMessageTopic(groupId [32]byte) string {
	// We should think about simplifying the topics, since backwards compatibility shouldn't really matter here
	return fmt.Sprintf("/xmtp/1/g-%x/proto", groupId)
}
