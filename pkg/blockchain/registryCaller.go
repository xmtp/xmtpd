package blockchain

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/contracts/pkg/nodesv2"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

type INodeRegistryCaller interface {
	GetAllNodes(ctx context.Context) ([]nodesv2.INodesNodeWithId, error)
	OwnerOf(ctx context.Context, nodeId int64) (common.Address, error)
}

type nodeRegistryCaller struct {
	client   *ethclient.Client
	logger   *zap.Logger
	contract *nodesv2.NodesV2Caller
}

func NewNodeRegistryCaller(
	logger *zap.Logger,
	client *ethclient.Client,
	contractsOptions config.ContractsOptions,
) (INodeRegistryCaller, error) {
	contract, err := nodesv2.NewNodesV2Caller(
		common.HexToAddress(contractsOptions.NodesContractAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	return &nodeRegistryCaller{
		client:   client,
		logger:   logger.Named("NodeRegistryCaller"),
		contract: contract,
	}, nil
}

func (n *nodeRegistryCaller) GetActiveApiNodes(
	ctx context.Context,
) ([]nodesv2.INodesNodeWithId, error) {
	return n.contract.GetActiveApiNodes(&bind.CallOpts{
		Context: ctx,
	})
}

func (n *nodeRegistryCaller) GetActiveReplicationNodes(
	ctx context.Context,
) ([]nodesv2.INodesNodeWithId, error) {
	return n.contract.GetActiveReplicationNodes(&bind.CallOpts{
		Context: ctx,
	})
}

func (n *nodeRegistryCaller) GetAllNodes(
	ctx context.Context,
) ([]nodesv2.INodesNodeWithId, error) {
	return n.contract.GetAllNodes(&bind.CallOpts{
		Context: ctx,
	})
}

func (n *nodeRegistryCaller) GetAllNodesCount(
	ctx context.Context,
) (uint64, error) {
	count, err := n.contract.GetAllNodesCount(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return 0, err
	}
	return count.Uint64(), nil
}

func (n *nodeRegistryCaller) GetNode(
	ctx context.Context,
	nodeId int64,
) (nodesv2.INodesNode, error) {
	return n.contract.GetNode(&bind.CallOpts{
		Context: ctx,
	}, big.NewInt(nodeId))
}

func (n *nodeRegistryCaller) OwnerOf(
	ctx context.Context,
	nodeId int64,
) (common.Address, error) {
	return n.contract.OwnerOf(&bind.CallOpts{
		Context: ctx,
	}, big.NewInt(nodeId))
}
