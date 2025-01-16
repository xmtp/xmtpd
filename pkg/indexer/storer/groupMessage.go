package storer

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/xmtp/xmtpd/contracts/pkg/groupmessages"
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

type GroupMessageStorer struct {
	contract *groupmessages.GroupMessages
	queries  *queries.Queries
	logger   *zap.Logger
}

func NewGroupMessageStorer(
	queries *queries.Queries,
	logger *zap.Logger,
	contract *groupmessages.GroupMessages,
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
	appendLog bool,
) LogStorageError {
	msgSent, err := s.contract.ParseMessageSent(event)
	if err != nil {
		return NewLogStorageError(err, false)
	}

	topicStruct := topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, msgSent.GroupId[:])

	clientEnvelope, err := envelopes.NewClientEnvelopeFromBytes(msgSent.Message)
	if err != nil {
		s.logger.Error("Error parsing client envelope", zap.Error(err))
		return NewLogStorageError(err, false)
	}

	targetTopic := clientEnvelope.TargetTopic()

	if !clientEnvelope.TopicMatchesPayload() {
		s.logger.Error(
			"Client envelope topic does not match payload type",
			zap.Any("targetTopic", targetTopic.String()),
			zap.Any("contractTopic", topicStruct.String()),
		)
		return NewLogStorageError(
			errors.New("client envelope topic does not match payload topic"),
			false,
		)
	}

	signedOriginatorEnvelope, err := buildSignedOriginatorEnvelope(
		buildOriginatorEnvelope(msgSent.SequenceId, msgSent.Message),
		event.TxHash,
	)
	if err != nil {
		s.logger.Error("Error building signed originator envelope", zap.Error(err))
		return NewLogStorageError(err, false)
	}

	originatorEnvelopeBytes, err := proto.Marshal(signedOriginatorEnvelope)
	if err != nil {
		s.logger.Error("Error marshalling originator envelope", zap.Error(err))
		return NewLogStorageError(err, false)
	}

	version := sql.NullInt32{Int32: 1, Valid: true}

	if appendLog {
		version, err = GetVersionForAppend(
			ctx,
			s.queries,
			s.logger,
			GROUP_MESSAGE_ORIGINATOR_ID,
			int64(msgSent.SequenceId),
		)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return NewLogStorageError(err, true)
			}
			if errors.Is(err, sql.ErrNoRows) {
				s.logger.Debug("No rows found for envelope, inserting new",
					zap.Int("originator_node_id", GROUP_MESSAGE_ORIGINATOR_ID),
					zap.Int64("originator_sequence_id", int64(msgSent.SequenceId)),
				)
			}
		}
	}

	s.logger.Debug("Inserting message from contract", zap.String("topic", topicStruct.String()))

	if _, err = s.queries.InsertGatewayEnvelope(ctx, queries.InsertGatewayEnvelopeParams{
		BlockNumber:          sql.NullInt64{Int64: int64(event.BlockNumber), Valid: true},
		BlockHash:            event.BlockHash.Bytes(),
		Version:              version,
		IsCanonical:          sql.NullBool{Bool: true, Valid: true},
		OriginatorNodeID:     GROUP_MESSAGE_ORIGINATOR_ID,
		OriginatorSequenceID: int64(msgSent.SequenceId),
		Topic:                topicStruct.Bytes(),
		OriginatorEnvelope:   originatorEnvelopeBytes,
	}); err != nil {
		s.logger.Error("Error inserting envelope from smart contract", zap.Error(err))
		return NewLogStorageError(err, true)
	}

	return nil
}
