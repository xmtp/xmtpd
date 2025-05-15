package app_chain

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	bt "github.com/xmtp/xmtpd/pkg/indexer/block_tracker"
	"github.com/xmtp/xmtpd/pkg/utils"
)

const (
	identityUpdateName  = "identityUpdateBroadcaster"
	identityUpdateTopic = "IdentityUpdateCreated"
)

// TODO: Abstract to interface.
// TODO: Include LogStorer in the IdentityUpdateBroadcaster struct.
type IdentityUpdateBroadcaster struct {
	address      string
	contract     *iu.IdentityUpdateBroadcaster
	topics       []common.Hash
	blockTracker *bt.BlockTracker
}

func NewIdentityUpdateBroadcaster(
	ctx context.Context,
	client *ethclient.Client,
	querier *queries.Queries,
	address string,
) (*IdentityUpdateBroadcaster, error) {
	contract, err := identityUpdateBroadcasterContract(address, client)
	if err != nil {
		return nil, err
	}

	identityUpdatesTracker, err := bt.NewBlockTracker(
		ctx,
		address,
		querier,
	)
	if err != nil {
		return nil, err
	}

	topics, err := identityUpdateBroadcasterTopic()
	if err != nil {
		return nil, err
	}

	return &IdentityUpdateBroadcaster{
		address:      address,
		contract:     contract,
		blockTracker: identityUpdatesTracker,
		topics:       []common.Hash{topics},
	}, nil
}

func (iu *IdentityUpdateBroadcaster) Address() common.Address {
	return common.HexToAddress(iu.address)
}

func (iu *IdentityUpdateBroadcaster) Topics() []common.Hash {
	return iu.topics
}

func identityUpdateBroadcasterContract(
	address string,
	client *ethclient.Client,
) (*iu.IdentityUpdateBroadcaster, error) {
	return iu.NewIdentityUpdateBroadcaster(
		common.HexToAddress(address),
		client,
	)
}

func identityUpdateBroadcasterName(chainID int) string {
	return fmt.Sprintf("%s-%v", identityUpdateName, chainID)
}

func identityUpdateBroadcasterTopic() (common.Hash, error) {
	abi, err := iu.IdentityUpdateBroadcasterMetaData.GetAbi()
	if err != nil {
		return common.Hash{}, err
	}
	return utils.GetEventTopic(abi, identityUpdateTopic)
}
