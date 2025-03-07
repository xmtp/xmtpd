package blockchain

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/contracts/pkg/nodes"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

type INodeRegistryCaller interface {
	GetActiveApiNodes(ctx context.Context) ([]nodes.INodesNodeWithId, error)
	GetActiveReplicationNodes(ctx context.Context) ([]nodes.INodesNodeWithId, error)
	GetAllNodes(ctx context.Context) ([]nodes.INodesNodeWithId, error)
	GetNode(ctx context.Context, nodeId int64) (nodes.INodesNode, error)
	OwnerOf(ctx context.Context, nodeId int64) (common.Address, error)
}

type nodeRegistryCaller struct {
	client   *ethclient.Client
	logger   *zap.Logger
	contract *nodes.NodesCaller
}

func NewNodeRegistryCaller(
	logger *zap.Logger,
	client *ethclient.Client,
	contractsOptions config.ContractsOptions,
) (INodeRegistryCaller, error) {
	contract, err := nodes.NewNodesCaller(
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
) ([]nodes.INodesNodeWithId, error) {
	return n.contract.GetActiveApiNodes(&bind.CallOpts{
		Context: ctx,
	})
}

func (n *nodeRegistryCaller) GetActiveReplicationNodes(
	ctx context.Context,
) ([]nodes.INodesNodeWithId, error) {
	return n.contract.GetActiveReplicationNodes(&bind.CallOpts{
		Context: ctx,
	})
}

func (n *nodeRegistryCaller) GetAllNodes(
	ctx context.Context,
) ([]nodes.INodesNodeWithId, error) {
	return n.contract.GetAllNodes(&bind.CallOpts{
		Context: ctx,
	})
}

func (n *nodeRegistryCaller) GetNode(
	ctx context.Context,
	nodeId int64,
) (nodes.INodesNode, error) {
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
