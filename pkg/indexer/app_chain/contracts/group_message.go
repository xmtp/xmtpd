// Package contracts implements the GroupMessageBroadcaster and IdentityUpdateBroadcaster contracts.
//
// The GroupMessageBroadcaster contract is responsible for broadcasting group messages to the network.
// The IdentityUpdateBroadcaster contract is responsible for broadcasting identity updates to the network.
// The Solidity implementations are in https://github.com/xmtp/smart-contracts/tree/main/src/app-chain.
package contracts

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	"github.com/xmtp/xmtpd/pkg/db"
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
	db *db.Handler,
	logger *zap.Logger,
	address common.Address,
	chainID int64,
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
		db,
		startBlock,
	)
	if err != nil {
		return nil, err
	}

	topics, err := groupMessageBroadcasterTopic()
	if err != nil {
		return nil, err
	}

	logger = logger.Named(utils.GroupMessageBroadcasterLoggerName).
		With(utils.ContractAddressField(address.Hex()))

	groupMessageStorer := NewGroupMessageStorer(db.Query(), logger, contract)

	reorgHandler := NewGroupMessageReorgHandler(logger)

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

func GroupMessageBroadcasterName(chainID int64) string {
	return fmt.Sprintf("%s-%v", groupMessageName, chainID)
}

func groupMessageBroadcasterTopic() (common.Hash, error) {
	abi, err := gm.GroupMessageBroadcasterMetaData.GetAbi()
	if err != nil {
		return common.Hash{}, err
	}
	return utils.GetEventTopic(abi, groupMessageTopic)
}
