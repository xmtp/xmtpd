package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/contracts/pkg/nodes"
	"github.com/xmtp/xmtpd/contracts/pkg/nodesv2"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

type RegistryCallerVersion int

const (
	RegistryCallerV1 RegistryCallerVersion = iota
	RegistryCallerV2
)

type INodeRegistryCaller interface {
	GetAllNodesV1(ctx context.Context) ([]nodes.NodesNodeWithId, error)
	GetAllNodesV2(ctx context.Context) ([]nodesv2.INodesNodeWithId, error)
	OwnerOf(ctx context.Context, nodeId int64) (common.Address, error)
}

type baseNodeRegistryCaller struct {
	client *ethclient.Client
	logger *zap.Logger
}

func NewNodeRegistryCaller(
	logger *zap.Logger,
	client *ethclient.Client,
	contractsOptions config.ContractsOptions,
	version RegistryCallerVersion,
) (INodeRegistryCaller, error) {
	switch version {
	case RegistryCallerV1:
		return newNodeRegistryCallerV1(logger, client, contractsOptions)
	case RegistryCallerV2:
		return newNodeRegistryCallerV2(logger, client, contractsOptions)
	default:
		return nil, fmt.Errorf("unsupported registry version: %v", version)
	}
}

/*
*
XMTP Node Registry Caller V1 - Deprecated
*
*/
type nodeRegistryCallerV1 struct {
	baseNodeRegistryCaller
	contract *nodes.NodesCaller
}

func newNodeRegistryCallerV1(
	logger *zap.Logger,
	client *ethclient.Client,
	contractsOptions config.ContractsOptions,
) (*nodeRegistryCallerV1, error) {
	contract, err := nodes.NewNodesCaller(
		common.HexToAddress(contractsOptions.NodesContractAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	return &nodeRegistryCallerV1{
		baseNodeRegistryCaller: baseNodeRegistryCaller{
			client: client,
			logger: logger.Named("NodeRegistryCallerV1"),
		},
		contract: contract,
	}, nil
}

func (n *nodeRegistryCallerV1) GetAllNodesV1(
	ctx context.Context,
) ([]nodes.NodesNodeWithId, error) {

	return n.contract.AllNodes(&bind.CallOpts{
		Context: ctx,
	})
}

func (n *nodeRegistryCallerV1) GetAllNodesV2(
	ctx context.Context,
) ([]nodesv2.INodesNodeWithId, error) {
	return nil, fmt.Errorf("not implemented")
}

func (n *nodeRegistryCallerV1) OwnerOf(
	ctx context.Context,
	nodeId int64,
) (common.Address, error) {
	return n.contract.OwnerOf(&bind.CallOpts{
		Context: ctx,
	}, big.NewInt(nodeId))
}

/*
*
XMTP Node Registry Caller V2
*
*/
type nodeRegistryCallerV2 struct {
	baseNodeRegistryCaller
	contract *nodesv2.NodesV2Caller
}

func newNodeRegistryCallerV2(
	logger *zap.Logger,
	client *ethclient.Client,
	contractsOptions config.ContractsOptions,
) (*nodeRegistryCallerV2, error) {
	contract, err := nodesv2.NewNodesV2Caller(
		common.HexToAddress(contractsOptions.NodesContractAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	return &nodeRegistryCallerV2{
		baseNodeRegistryCaller: baseNodeRegistryCaller{
			client: client,
			logger: logger.Named("NodeRegistryCallerV2"),
		},
		contract: contract,
	}, nil
}

func (n *nodeRegistryCallerV2) GetAllNodesV1(
	ctx context.Context,
) ([]nodes.NodesNodeWithId, error) {
	return nil, fmt.Errorf("not implemented")
}

func (n *nodeRegistryCallerV2) GetAllNodesV2(
	ctx context.Context,
) ([]nodesv2.INodesNodeWithId, error) {

	return n.contract.GetAllNodes(&bind.CallOpts{
		Context: ctx,
	})
}

func (n *nodeRegistryCallerV2) OwnerOf(
	ctx context.Context,
	nodeId int64,
) (common.Address, error) {
	return n.contract.OwnerOf(&bind.CallOpts{
		Context: ctx,
	}, big.NewInt(nodeId))
}
