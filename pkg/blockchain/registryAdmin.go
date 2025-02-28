package blockchain

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/contracts/pkg/nodes"
	"github.com/xmtp/xmtpd/contracts/pkg/nodesv2"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

/*
*
A NodeRegistryAdmin is a struct responsible for calling admin functions on the node registry
*
*/
type NodeRegistryAdmin struct {
	client   *ethclient.Client
	signer   TransactionSigner
	contract *nodes.Nodes
	logger   *zap.Logger
}

func NewNodeRegistryAdmin(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractsOptions config.ContractsOptions,
) (*NodeRegistryAdmin, error) {
	contract, err := nodes.NewNodes(
		common.HexToAddress(contractsOptions.NodesContractAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	return &NodeRegistryAdmin{
		signer:   signer,
		client:   client,
		logger:   logger.Named("NodeRegistryAdmin"),
		contract: contract,
	}, nil
}

// XMTP Node Registry V1 - Deprecated
// TODO(borja): Remove once migration V1 -> V2 is done
func (n *NodeRegistryAdmin) AddNode(
	ctx context.Context,
	owner string,
	signingKeyPub *ecdsa.PublicKey,
	httpAddress string,
) error {
	if !common.IsHexAddress(owner) {
		return fmt.Errorf("invalid owner address provided %s", owner)
	}

	ownerAddress := common.HexToAddress(owner)
	signingKey := crypto.FromECDSAPub(signingKeyPub)

	if n.signer == nil {
		return fmt.Errorf("no signer provided")
	}
	tx, err := n.contract.AddNode(&bind.TransactOpts{
		Context: ctx,
		From:    n.signer.FromAddress(),
		Signer:  n.signer.SignerFunc(),
	}, ownerAddress, signingKey, httpAddress)

	if err != nil {
		return err
	}

	receipt, err := WaitForTransaction(
		ctx,
		n.logger,
		n.client,
		2*time.Second,
		250*time.Millisecond,
		tx.Hash(),
	)

	if err != nil {
		return err
	}

	for _, log := range receipt.Logs {
		nodeUpdated, err := n.contract.ParseNodeUpdated(*log)
		if err != nil {
			continue
		}

		n.logger.Info("node added to registry V1",
			zap.Uint64("node_id", nodeUpdated.NodeId.Uint64()),
			zap.String("http_address", nodeUpdated.Node.HttpAddress),
			zap.String("signing_key_pub", hex.EncodeToString(nodeUpdated.Node.SigningKeyPub)),
		)
	}

	return err
}

/*
*
A NodeRegistryCaller is a struct responsible for calling public functions on the node registry
*
*/
type NodeRegistryCaller struct {
	client   *ethclient.Client
	contract *nodes.NodesCaller
	logger   *zap.Logger
}

func NewNodeRegistryCaller(
	logger *zap.Logger,
	client *ethclient.Client,
	contractsOptions config.ContractsOptions,
) (*NodeRegistryCaller, error) {
	contract, err := nodes.NewNodesCaller(
		common.HexToAddress(contractsOptions.NodesContractAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	return &NodeRegistryCaller{
		client:   client,
		logger:   logger.Named("NodeRegistryCaller"),
		contract: contract,
	}, nil
}

func (n *NodeRegistryCaller) GetAllNodes(
	ctx context.Context,
) ([]nodes.NodesNodeWithId, error) {

	return n.contract.AllNodes(&bind.CallOpts{
		Context: ctx,
	})
}

func (n *NodeRegistryCaller) OwnerOf(
	ctx context.Context,
	nodeId int64,
) (common.Address, error) {
	return n.contract.OwnerOf(&bind.CallOpts{
		Context: ctx,
	}, big.NewInt(nodeId))
}

func (n *NodeRegistryAdmin) UpdateHealth(
	ctx context.Context, nodeId int64, health bool,
) error {
	tx, err := n.contract.UpdateHealth(
		&bind.TransactOpts{
			Context: ctx,
			From:    n.signer.FromAddress(),
			Signer:  n.signer.SignerFunc(),
		},
		big.NewInt(nodeId),
		health,
	)

	if err != nil {
		return err
	}

	_, err = WaitForTransaction(
		ctx,
		n.logger,
		n.client,
		2*time.Second,
		250*time.Millisecond,
		tx.Hash(),
	)

	return err
}

func (n *NodeRegistryAdmin) UpdateHttpAddress(
	ctx context.Context, nodeId int64, address string,
) error {
	tx, err := n.contract.UpdateHttpAddress(
		&bind.TransactOpts{
			Context: ctx,
			From:    n.signer.FromAddress(),
			Signer:  n.signer.SignerFunc(),
		},
		big.NewInt(nodeId),
		address,
	)

	if err != nil {
		return err
	}

	_, err = WaitForTransaction(
		ctx,
		n.logger,
		n.client,
		2*time.Second,
		250*time.Millisecond,
		tx.Hash(),
	)

	return err
}

/*
*
XMTP Node Registry V2
TODO(borja): Remove once migration is complete.
*
*/
type NodeRegistryAdminV2 struct {
	client   *ethclient.Client
	signer   TransactionSigner
	contract *nodesv2.NodesV2
	logger   *zap.Logger
}

func NewNodeRegistryAdminV2(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractsOptions config.ContractsOptions,
) (*NodeRegistryAdminV2, error) {
	contract, err := nodesv2.NewNodesV2(
		common.HexToAddress(contractsOptions.NodesContractAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	return &NodeRegistryAdminV2{
		signer:   signer,
		client:   client,
		logger:   logger.Named("NodeRegistryAdminV2"),
		contract: contract,
	}, nil
}

func (n *NodeRegistryAdminV2) AddNodeV2(
	ctx context.Context,
	owner string,
	signingKeyPub *ecdsa.PublicKey,
	httpAddress string,
	minMonthlyFee *big.Int,
) error {
	if !common.IsHexAddress(owner) {
		return fmt.Errorf("invalid owner address provided %s", owner)
	}

	if minMonthlyFee == nil {
		minMonthlyFee = big.NewInt(0)
	}

	ownerAddress := common.HexToAddress(owner)
	signingKey := crypto.FromECDSAPub(signingKeyPub)

	if n.signer == nil {
		return fmt.Errorf("no signer provided")
	}

	tx, err := n.contract.AddNode(&bind.TransactOpts{
		Context: ctx,
		From:    n.signer.FromAddress(),
		Signer:  n.signer.SignerFunc(),
	}, ownerAddress, signingKey, httpAddress, minMonthlyFee)
	if err != nil {
		fmt.Println("error adding node")
		return err
	}

	receipt, err := WaitForTransaction(
		ctx,
		n.logger,
		n.client,
		2*time.Second,
		250*time.Millisecond,
		tx.Hash(),
	)

	if err != nil {
		return err
	}

	for _, log := range receipt.Logs {
		nodeAdded, err := n.contract.ParseNodeAdded(*log)
		if err != nil {
			continue
		}
		n.logger.Info("node added to registry V2",
			zap.Uint64("node_id", nodeAdded.NodeId.Uint64()),
			zap.String("owner", nodeAdded.Owner.Hex()),
			zap.String("http_address", nodeAdded.HttpAddress),
			zap.String("signing_key_pub", hex.EncodeToString(nodeAdded.SigningKeyPub)),
			zap.String("min_monthly_fee", nodeAdded.MinMonthlyFee.String()),
		)
	}

	return nil
}

/*
*
A NodeRegistryCaller is a struct responsible for calling public functions on the node registry
*
*/
type NodeRegistryCallerV2 struct {
	client   *ethclient.Client
	contract *nodesv2.NodesV2Caller
	logger   *zap.Logger
}

func NewNodeRegistryCallerV2(
	logger *zap.Logger,
	client *ethclient.Client,
	contractsOptions config.ContractsOptions,
) (*NodeRegistryCallerV2, error) {
	contract, err := nodesv2.NewNodesV2Caller(
		common.HexToAddress(contractsOptions.NodesContractAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	return &NodeRegistryCallerV2{
		client:   client,
		logger:   logger.Named("NodeRegistryCallerV2"),
		contract: contract,
	}, nil
}

func (n *NodeRegistryCallerV2) GetAllNodes(
	ctx context.Context,
) ([]nodesv2.INodesNodeWithId, error) {

	return n.contract.GetAllNodes(&bind.CallOpts{
		Context: ctx,
	})
}

func (n *NodeRegistryCallerV2) OwnerOf(
	ctx context.Context,
	nodeId int64,
) (common.Address, error) {
	return n.contract.OwnerOf(&bind.CallOpts{
		Context: ctx,
	}, big.NewInt(nodeId))
}
