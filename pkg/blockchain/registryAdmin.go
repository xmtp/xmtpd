package blockchain

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/abi/noderegistry"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

type INodeRegistryAdmin interface {
	AddNode(
		ctx context.Context,
		owner string,
		signingKeyPub *ecdsa.PublicKey,
		httpAddress string,
	) error
	AddToNetwork(ctx context.Context, nodeId uint32) error
	RemoveFromNetwork(ctx context.Context, nodeId uint32) error
	SetHttpAddress(ctx context.Context, nodeId uint32, httpAddress string) error
	SetMaxCanonical(ctx context.Context, limit uint8) error
}

type nodeRegistryAdmin struct {
	client   *ethclient.Client
	signer   TransactionSigner
	logger   *zap.Logger
	contract *noderegistry.NodeRegistry
}

var _ INodeRegistryAdmin = &nodeRegistryAdmin{}

func NewNodeRegistryAdmin(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractsOptions config.ContractsOptions,
) (*nodeRegistryAdmin, error) {
	contract, err := noderegistry.NewNodeRegistry(
		common.HexToAddress(contractsOptions.SettlementChain.NodeRegistryAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	return &nodeRegistryAdmin{
		client:   client,
		signer:   signer,
		logger:   logger.Named("NodeRegistryAdmin"),
		contract: contract,
	}, nil
}

func (n *nodeRegistryAdmin) AddNode(
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

	return ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.AddNode(
				opts,
				ownerAddress,
				signingKey,
				httpAddress,
			)
		},
		func(log *types.Log) (interface{}, error) {
			return n.contract.ParseNodeAdded(*log)
		},
		func(event interface{}) {
			nodeAdded, ok := event.(*noderegistry.NodeRegistryNodeAdded)
			if !ok {
				n.logger.Error("node added event is not of type NodesNodeAdded")
				return
			}
			n.logger.Info("node added to registry",
				zap.Uint32("node_id", nodeAdded.NodeId),
				zap.String("owner", nodeAdded.Owner.Hex()),
				zap.String("http_address", nodeAdded.HttpAddress),
				zap.String("signing_key_pub", hex.EncodeToString(nodeAdded.SigningPublicKey)),
			)
		},
	)
}

func (n *nodeRegistryAdmin) AddToNetwork(
	ctx context.Context,
	nodeId uint32,
) error {
	return ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.AddToNetwork(
				opts,
				nodeId,
			)
		},
		func(log *types.Log) (interface{}, error) {
			return n.contract.ParseNodeAddedToCanonicalNetwork(*log)
		},
		func(event interface{}) {
			nodeAdded, ok := event.(*noderegistry.NodeRegistryNodeAddedToCanonicalNetwork)
			if !ok {
				n.logger.Error(
					"unexpected event type, not of type NodesNodeAddedToCanonicalNetwork",
				)
				return
			}
			n.logger.Info("node added to canonical network",
				zap.Uint32("node_id", nodeAdded.NodeId),
			)
		},
	)
}

func (n *nodeRegistryAdmin) RemoveFromNetwork(
	ctx context.Context,
	nodeId uint32,
) error {
	return ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.RemoveFromNetwork(
				opts,
				nodeId,
			)
		},
		func(log *types.Log) (interface{}, error) {
			return n.contract.ParseNodeRemovedFromCanonicalNetwork(*log)
		},
		func(event interface{}) {
			nodeRemoved, ok := event.(*noderegistry.NodeRegistryNodeRemovedFromCanonicalNetwork)
			if !ok {
				n.logger.Error(
					"unexpected event type, not of type NodesNodeRemovedFromCanonicalNetwork",
				)
				return
			}
			n.logger.Info("node removed from canonical network",
				zap.Uint32("node_id", nodeRemoved.NodeId),
			)
		},
	)
}

func (n *nodeRegistryAdmin) SetHttpAddress(
	ctx context.Context,
	nodeId uint32,
	httpAddress string,
) error {
	return ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.SetHttpAddress(
				opts,
				nodeId,
				httpAddress,
			)
		},
		func(log *types.Log) (interface{}, error) {
			return n.contract.ParseHttpAddressUpdated(*log)
		},
		func(event interface{}) {
			httpAddressUpdated, ok := event.(*noderegistry.NodeRegistryHttpAddressUpdated)
			if !ok {
				n.logger.Error(
					"http address updated event is not of type NodesHttpAddressUpdated",
				)
				return
			}
			n.logger.Info("http address updated",
				zap.Uint32("node_id", httpAddressUpdated.NodeId),
				zap.String("http_address", httpAddressUpdated.HttpAddress),
			)
		},
	)
}

func (n *nodeRegistryAdmin) SetMaxCanonical(
	ctx context.Context,
	limit uint8,
) error {
	// TODO set the actual limit

	return ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.UpdateMaxCanonicalNodes(
				opts,
			)
		},
		func(log *types.Log) (interface{}, error) {
			return n.contract.ParseMaxCanonicalNodesUpdated(*log)
		},
		func(event interface{}) {
			maxCanonicalUpdated, ok := event.(*noderegistry.NodeRegistryMaxCanonicalNodesUpdated)
			if !ok {
				n.logger.Error(
					"unexpected event type, not of type NodeRegistryMaxCanonicalNodesUpdated",
				)
				return
			}
			n.logger.Info("updated max canonical nodes",
				zap.Uint8("limit", maxCanonicalUpdated.MaxCanonicalNodes),
			)
		},
	)
}
