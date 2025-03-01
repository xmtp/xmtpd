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

type RegistryAdminVersion int

const (
	RegistryAdminV1 RegistryAdminVersion = iota
	RegistryAdminV2
)

type INodeRegistryAdmin interface {
	AddNode(
		ctx context.Context,
		owner string,
		signingKeyPub *ecdsa.PublicKey,
		httpAddress string,
	) error
	UpdateHealth(ctx context.Context, nodeId int64, health bool) error
	UpdateHttpAddress(ctx context.Context, nodeId int64, address string) error
}

type baseNodeRegistryAdmin struct {
	client *ethclient.Client
	signer TransactionSigner
	logger *zap.Logger
}

func NewNodeRegistryAdmin(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractsOptions config.ContractsOptions,
	version RegistryAdminVersion,
) (INodeRegistryAdmin, error) {
	switch version {
	case RegistryAdminV1:
		return newNodeRegistryAdminV1(logger, client, signer, contractsOptions)
	case RegistryAdminV2:
		return newNodeRegistryAdminV2(logger, client, signer, contractsOptions)
	default:
		return nil, fmt.Errorf("unsupported registry version: %v", version)
	}
}

/*
*
XMTP Node Registry Admin V1 - Deprecated
*
*/
type nodeRegistryAdminV1 struct {
	baseNodeRegistryAdmin
	contract *nodes.Nodes
}

var _ INodeRegistryAdmin = &nodeRegistryAdminV1{}

func newNodeRegistryAdminV1(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractsOptions config.ContractsOptions,
) (*nodeRegistryAdminV1, error) {
	contract, err := nodes.NewNodes(
		common.HexToAddress(contractsOptions.NodesContractAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	return &nodeRegistryAdminV1{
		baseNodeRegistryAdmin: baseNodeRegistryAdmin{
			signer: signer,
			client: client,
			logger: logger.Named("NodeRegistryAdminV1"),
		},
		contract: contract,
	}, nil
}

func (n *nodeRegistryAdminV1) AddNode(
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

func (n *nodeRegistryAdminV1) UpdateHealth(
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

func (n *nodeRegistryAdminV1) UpdateHttpAddress(
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
XMTP Node Registry Admin V2
*
*/
type nodeRegistryAdminV2 struct {
	baseNodeRegistryAdmin
	contract *nodesv2.NodesV2
}

var _ INodeRegistryAdmin = &nodeRegistryAdminV2{}

func newNodeRegistryAdminV2(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractsOptions config.ContractsOptions,
) (*nodeRegistryAdminV2, error) {
	contract, err := nodesv2.NewNodesV2(
		common.HexToAddress(contractsOptions.NodesContractAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	return &nodeRegistryAdminV2{
		baseNodeRegistryAdmin: baseNodeRegistryAdmin{
			signer: signer,
			client: client,
			logger: logger.Named("NodeRegistryAdminV2"),
		},
		contract: contract,
	}, nil
}

func (n *nodeRegistryAdminV2) AddNode(
	ctx context.Context,
	owner string,
	signingKeyPub *ecdsa.PublicKey,
	httpAddress string,
) error {
	if !common.IsHexAddress(owner) {
		return fmt.Errorf("invalid owner address provided %s", owner)
	}

	minMonthlyFee := big.NewInt(0)
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

func (n *nodeRegistryAdminV2) UpdateHealth(
	ctx context.Context, nodeId int64, health bool,
) error {
	return fmt.Errorf("not implemented")
}

func (n *nodeRegistryAdminV2) UpdateHttpAddress(
	ctx context.Context, nodeId int64, address string,
) error {
	return nil
}
