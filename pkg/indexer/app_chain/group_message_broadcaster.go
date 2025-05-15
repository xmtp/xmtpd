package app_chain

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	bt "github.com/xmtp/xmtpd/pkg/indexer/block_tracker"
	"github.com/xmtp/xmtpd/pkg/utils"
)

const (
	groupMessageName  = "groupMessageBroadcaster"
	groupMessageTopic = "MessageSent"
)

type GroupMessageBroadcaster struct {
	address      string
	contract     *gm.GroupMessageBroadcaster
	topics       []common.Hash
	blockTracker *bt.BlockTracker
}

// TODO: Abstract to interface.
// TODO: Include LogStorer in the GroupMessageBroadcaster struct.
func NewGroupMessageBroadcaster(
	ctx context.Context,
	client *ethclient.Client,
	querier *queries.Queries,
	address string,
) (*GroupMessageBroadcaster, error) {
	contract, err := groupMessageBroadcasterContract(address, client)
	if err != nil {
		return nil, err
	}

	groupMessagesTracker, err := bt.NewBlockTracker(
		ctx,
		address,
		querier,
	)
	if err != nil {
		return nil, err
	}

	topics, err := groupMessageBroadcasterTopic()
	if err != nil {
		return nil, err
	}

	return &GroupMessageBroadcaster{
		address:      address,
		contract:     contract,
		blockTracker: groupMessagesTracker,
		topics:       []common.Hash{topics},
	}, nil
}

func (gm *GroupMessageBroadcaster) Address() common.Address {
	return common.HexToAddress(gm.address)
}

func (gm *GroupMessageBroadcaster) Topics() []common.Hash {
	return gm.topics
}

func groupMessageBroadcasterContract(
	address string,
	client *ethclient.Client,
) (*gm.GroupMessageBroadcaster, error) {
	return gm.NewGroupMessageBroadcaster(
		common.HexToAddress(address),
		client,
	)
}

func groupMessageBroadcasterName(chainID int) string {
	return fmt.Sprintf("%s-%v", groupMessageName, chainID)
}

func groupMessageBroadcasterTopic() (common.Hash, error) {
	abi, err := gm.GroupMessageBroadcasterMetaData.GetAbi()
	if err != nil {
		return common.Hash{}, err
	}
	return utils.GetEventTopic(abi, groupMessageTopic)
}
