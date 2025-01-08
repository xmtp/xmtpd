package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/contracts/pkg/nodes"
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
		logger:   logger.Named("NodeRegistryAdmin"),
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
