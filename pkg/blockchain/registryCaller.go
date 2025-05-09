package blockchain

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/abi/noderegistry"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

type INodeRegistryCaller interface {
	GetAllNodes(ctx context.Context) ([]noderegistry.INodeRegistryNodeWithId, error)
	GetNode(ctx context.Context, nodeId int64) (noderegistry.INodeRegistryNode, error)
	OwnerOf(ctx context.Context, nodeId int64) (common.Address, error)
}

type nodeRegistryCaller struct {
	client   *ethclient.Client
	logger   *zap.Logger
	contract *noderegistry.NodeRegistryCaller
}

func NewNodeRegistryCaller(
	logger *zap.Logger,
	client *ethclient.Client,
	contractsOptions config.ContractsOptions,
) (INodeRegistryCaller, error) {
	contract, err := noderegistry.NewNodeRegistryCaller(
		common.HexToAddress(contractsOptions.SettlementChain.NodeRegistryAddress),
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

func (n *nodeRegistryCaller) GetAllNodes(
	ctx context.Context,
) ([]noderegistry.INodeRegistryNodeWithId, error) {
	return n.contract.GetAllNodes(&bind.CallOpts{
		Context: ctx,
	})
}

func (n *nodeRegistryCaller) GetNode(
	ctx context.Context,
	nodeId int64,
) (noderegistry.INodeRegistryNode, error) {
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
