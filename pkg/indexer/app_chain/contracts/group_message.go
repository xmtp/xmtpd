package contracts

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	c "github.com/xmtp/xmtpd/pkg/indexer/common"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

const (
	groupMessageName  = "group-message-broadcaster"
	groupMessageTopic = "MessageSent"
)

type GroupMessageBroadcaster struct {
	address common.Address
	topics  []common.Hash
	logger  *zap.Logger
	c.IBlockTracker
	c.IReorgHandler
	c.ILogStorer
}

var _ c.IContract = &GroupMessageBroadcaster{}

func NewGroupMessageBroadcaster(
	ctx context.Context,
	client *ethclient.Client,
	querier *queries.Queries,
	logger *zap.Logger,
	address common.Address,
	chainID int,
	startBlock uint64,
) (*GroupMessageBroadcaster, error) {
	contract, err := groupMessageBroadcasterContract(address, client)
	if err != nil {
		return nil, err
	}

	groupMessagesTracker, err := c.NewBlockTracker(
		ctx,
		client,
		address,
		querier,
		startBlock,
	)
	if err != nil {
		return nil, err
	}

	topics, err := groupMessageBroadcasterTopic()
	if err != nil {
		return nil, err
	}

	logger = logger.Named("group-message-broadcaster").
		With(zap.String("contractAddress", address.Hex()))

	groupMessageStorer := NewGroupMessageStorer(querier, logger, contract)

	reorgHandler := c.NewChainReorgHandler(ctx, client, querier)

	return &GroupMessageBroadcaster{
		address:       address,
		topics:        []common.Hash{topics},
		logger:        logger,
		IBlockTracker: groupMessagesTracker,
		IReorgHandler: reorgHandler,
		ILogStorer:    groupMessageStorer,
	}, nil
}

func (gm *GroupMessageBroadcaster) Address() common.Address {
	return gm.address
}

func (gm *GroupMessageBroadcaster) Topics() []common.Hash {
	return gm.topics
}

func (gm *GroupMessageBroadcaster) Logger() *zap.Logger {
	return gm.logger
}

func groupMessageBroadcasterContract(
	address common.Address,
	client *ethclient.Client,
) (*gm.GroupMessageBroadcaster, error) {
	return gm.NewGroupMessageBroadcaster(
		address,
		client,
	)
}

func GroupMessageBroadcasterName(chainID int) string {
	return fmt.Sprintf("%s-%v", groupMessageName, chainID)
}

func groupMessageBroadcasterTopic() (common.Hash, error) {
	abi, err := gm.GroupMessageBroadcasterMetaData.GetAbi()
	if err != nil {
		return common.Hash{}, err
	}
	return utils.GetEventTopic(abi, groupMessageTopic)
}
