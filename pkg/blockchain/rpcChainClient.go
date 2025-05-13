package blockchain

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
)

// RpcChainClient wraps ethclient.Client and implements ChainClient
// It also provides contract log parsing methods
type RpcChainClient struct {
	*ethclient.Client
	groupMessageContract   *gm.GroupMessageBroadcaster
	identityUpdateContract *iu.IdentityUpdateBroadcaster
}

func NewRpcChainClient(
	ctx context.Context,
	rpcUrl string,
	messagesContractAddr, identityUpdatesContractAddr common.Address,
) (*RpcChainClient, error) {
	client, err := ethclient.DialContext(ctx, rpcUrl)
	if err != nil {
		return nil, err
	}
	groupMessageContract, err := gm.NewGroupMessageBroadcaster(messagesContractAddr, client)
	if err != nil {
		return nil, err
	}
	identityUpdateContract, err := iu.NewIdentityUpdateBroadcaster(
		identityUpdatesContractAddr,
		client,
	)
	if err != nil {
		return nil, err
	}
	return &RpcChainClient{
		Client:                 client,
		groupMessageContract:   groupMessageContract,
		identityUpdateContract: identityUpdateContract,
	}, nil
}

func (r *RpcChainClient) ParseMessageSent(
	log types.Log,
) (*gm.GroupMessageBroadcasterMessageSent, error) {
	return r.groupMessageContract.ParseMessageSent(log)
}

func (r *RpcChainClient) ParseIdentityUpdateCreated(
	log types.Log,
) (*iu.IdentityUpdateBroadcasterIdentityUpdateCreated, error) {
	return r.identityUpdateContract.ParseIdentityUpdateCreated(log)
}
