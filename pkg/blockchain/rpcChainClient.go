package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
	"github.com/xmtp/xmtpd/pkg/utils"
)

func buildMessagesTopic() (common.Hash, error) {
	abi, err := gm.GroupMessageBroadcasterMetaData.GetAbi()
	if err != nil {
		return common.Hash{}, err
	}
	return utils.GetEventTopic(abi, "MessageSent")
}

func buildIdentityUpdatesTopic() (common.Hash, error) {
	abi, err := iu.IdentityUpdateBroadcasterMetaData.GetAbi()
	if err != nil {
		return common.Hash{}, err
	}
	return utils.GetEventTopic(abi, "IdentityUpdateCreated")
}

type RpcChainClient struct {
	*ethclient.Client
	messagesContractAddress        common.Address
	identityUpdatesContractAddress common.Address
	groupMessageContract           *gm.GroupMessageBroadcaster
	identityUpdateContract         *iu.IdentityUpdateBroadcaster
}

// TODO(rich): rename
func NewRpcChainClient(
	ctx context.Context,
	rpcUrl string,
	messagesContractAddr, identityUpdatesContractAddr common.Address,
) (*RpcChainClient, error) {
	client, err := ethclient.DialContext(ctx, rpcUrl)
	if err != nil {
		return nil, err
	}
	groupMessageContract, err := gm.NewGroupMessageBroadcaster(
		messagesContractAddr,
		client,
	)
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
		Client:                         client,
		messagesContractAddress:        messagesContractAddr,
		identityUpdatesContractAddress: identityUpdatesContractAddr,
		groupMessageContract:           groupMessageContract,
		identityUpdateContract:         identityUpdateContract,
	}, nil
}

func (r *RpcChainClient) FilterLogs(
	ctx context.Context,
	eventType EventType,
	fromBlock, toBlock uint64,
) ([]types.Log, error) {
	query, err := r.buildFilterQuery(eventType, fromBlock, toBlock)
	if err != nil {
		return nil, err
	}
	return r.Client.FilterLogs(ctx, query)
}

func (r *RpcChainClient) buildFilterQuery(
	eventType EventType,
	fromBlock uint64,
	toBlock uint64,
) (ethereum.FilterQuery, error) {
	var contractAddress common.Address
	var topic common.Hash
	var err error

	switch eventType {
	case EventTypeMessageSent:
		contractAddress = r.messagesContractAddress
		topic, err = buildMessagesTopic()
	case EventTypeIdentityUpdateCreated:
		contractAddress = r.identityUpdatesContractAddress
		topic, err = buildIdentityUpdatesTopic()
	default:
		err = fmt.Errorf("unknown event type: %v", eventType)
	}
	if err != nil {
		return ethereum.FilterQuery{}, err
	}

	return ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		ToBlock:   new(big.Int).SetUint64(toBlock),
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{topic}},
	}, nil
}

func (r *RpcChainClient) ContractAddress(eventType EventType) (string, error) {
	switch eventType {
	case EventTypeMessageSent:
		return r.messagesContractAddress.Hex(), nil
	case EventTypeIdentityUpdateCreated:
		return r.identityUpdatesContractAddress.Hex(), nil
	default:
		return "", fmt.Errorf("unknown event type: %v", eventType)
	}
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
