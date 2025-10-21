package blockchain

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/abi/noderegistry"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

type INodeRegistryCaller interface {
	GetAllNodes(ctx context.Context) ([]noderegistry.INodeRegistryNodeWithId, error)
	GetNode(ctx context.Context, nodeID uint32) (noderegistry.INodeRegistryNode, error)
	OwnerOf(ctx context.Context, nodeID uint32) (common.Address, error)
	GetMaxCanonicalNodes(
		ctx context.Context,
	) (uint8, error)
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

	nodeRegistryCallerLogger := logger.Named(utils.NodeRegistryCallerLoggerName).With(
		utils.SettlementChainChainIDField(contractsOptions.SettlementChain.ChainID),
	)

	return &nodeRegistryCaller{
		client:   client,
		logger:   nodeRegistryCallerLogger,
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
	nodeID uint32,
) (noderegistry.INodeRegistryNode, error) {
	return n.contract.GetNode(&bind.CallOpts{
		Context: ctx,
	}, nodeID)
}

func (n *nodeRegistryCaller) OwnerOf(
	ctx context.Context,
	nodeID uint32,
) (common.Address, error) {
	return n.contract.OwnerOf(&bind.CallOpts{
		Context: ctx,
	}, big.NewInt(int64(nodeID)))
}

func (n *nodeRegistryCaller) GetMaxCanonicalNodes(
	ctx context.Context,
) (uint8, error) {
	return n.contract.MaxCanonicalNodes(&bind.CallOpts{
		Context: ctx,
	})
}
