package contracts

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"github.com/ethereum/go-ethereum/core/types"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	c "github.com/xmtp/xmtpd/pkg/indexer/common"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	re "github.com/xmtp/xmtpd/pkg/utils/retryerrors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
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

var _ c.ILogStorer = &GroupMessageStorer{}

func NewGroupMessageStorer(
	queries *queries.Queries,
	logger *zap.Logger,
	contract *gm.GroupMessageBroadcaster,
) *GroupMessageStorer {
	return &GroupMessageStorer{
		queries:  queries,
		logger:   logger.Named(utils.StorerLoggerName),
		contract: contract,
	}
}

// StoreLog validates and stores a group message log event
func (s *GroupMessageStorer) StoreLog(
	ctx context.Context,
	event types.Log,
) re.RetryableError {
	msgSent, err := s.contract.ParseMessageSent(event)
	if err != nil {
		return re.NewNonRecoverableError(ErrParseGroupMessage, err)
	}

	topicStruct := topic.NewTopic(topic.TopicKindGroupMessagesV1, msgSent.GroupId[:])

	clientEnvelope, err := envelopes.NewClientEnvelopeFromBytes(msgSent.Message)
	if err != nil {
		s.logger.Error(ErrParseClientEnvelope, zap.Error(err))
		return re.NewNonRecoverableError(ErrParseClientEnvelope, err)
	}

	if !clientEnvelope.TopicMatchesPayload() {
		return re.NewNonRecoverableError(ErrTopicDoesNotMatch, errors.New(ErrTopicDoesNotMatch))
	}

	originatorEnvelope, err := buildOriginatorEnvelope(
		constants.GroupMessageOriginatorID,
		msgSent.SequenceId,
		msgSent.Message,
	)
	if err != nil {
		s.logger.Error(ErrBuildOriginatorEnvelope, zap.Error(err))
		return re.NewNonRecoverableError(ErrBuildOriginatorEnvelope, err)
	}

	signedOriginatorEnvelope, err := buildSignedOriginatorEnvelope(
		originatorEnvelope,
		event.TxHash,
	)
	if err != nil {
		s.logger.Error(ErrBuildSignedOriginatorEnvelope, zap.Error(err))
		return re.NewNonRecoverableError(ErrBuildSignedOriginatorEnvelope, err)
	}

	originatorEnvelopeBytes, err := proto.Marshal(signedOriginatorEnvelope)
	if err != nil {
		s.logger.Error(ErrMarshallOriginatorEnvelope, zap.Error(err))
		return re.NewNonRecoverableError(ErrMarshallOriginatorEnvelope, err)
	}

	s.logger.Debug(
		"inserting message from contract",
		utils.TopicField(topicStruct.String()),
	)

	if _, err = s.queries.InsertGatewayEnvelope(ctx, queries.InsertGatewayEnvelopeParams{
		OriginatorNodeID:     constants.GroupMessageOriginatorID,
		OriginatorSequenceID: int64(msgSent.SequenceId),
		Topic:                topicStruct.Bytes(),
		OriginatorEnvelope:   originatorEnvelopeBytes,
		Expiry:               sql.NullInt64{Int64: math.MaxInt64, Valid: true},
	}); err != nil {
		s.logger.Error(ErrInsertEnvelopeFromSmartContract, zap.Error(err))
		return re.NewRecoverableError(ErrInsertEnvelopeFromSmartContract, err)
	}

	return nil
}
