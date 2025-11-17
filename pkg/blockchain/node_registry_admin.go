package blockchain

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/abi/noderegistry"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

type INodeRegistryAdmin interface {
	AddNode(
		ctx context.Context,
		ownerAddress common.Address,
		signingKeyPub *ecdsa.PublicKey,
		httpAddress string,
	) (uint32, error)
	AddToNetwork(ctx context.Context, nodeID uint32) error
	RemoveFromNetwork(ctx context.Context, nodeID uint32) error
	SetHTTPAddress(ctx context.Context, nodeID uint32, httpAddress string) error
	SetMaxCanonical(ctx context.Context, limit uint8) error
}

type nodeRegistryAdmin struct {
	client         *ethclient.Client
	signer         TransactionSigner
	logger         *zap.Logger
	nodeContract   *noderegistry.NodeRegistry
	parameterAdmin IParameterAdmin
}

var _ INodeRegistryAdmin = &nodeRegistryAdmin{}

func NewNodeRegistryAdmin(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractsOptions config.ContractsOptions,
	parameterAdmin IParameterAdmin,
) (*nodeRegistryAdmin, error) {
	nodeContract, err := noderegistry.NewNodeRegistry(
		common.HexToAddress(contractsOptions.SettlementChain.NodeRegistryAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	nodeRegistryAdminLogger := logger.Named(utils.NodeRegistryAdminLoggerName).With(
		utils.SettlementChainChainIDField(contractsOptions.SettlementChain.ChainID),
	)

	return &nodeRegistryAdmin{
		client:         client,
		signer:         signer,
		logger:         nodeRegistryAdminLogger,
		nodeContract:   nodeContract,
		parameterAdmin: parameterAdmin,
	}, nil
}

func (n *nodeRegistryAdmin) AddNode(
	ctx context.Context,
	ownerAddress common.Address,
	signingKeyPub *ecdsa.PublicKey,
	httpAddress string,
) (uint32, error) {
	signingKey := crypto.FromECDSAPub(signingKeyPub)

	nodeID := uint32(0)
	err := ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.nodeContract.AddNode(
				opts,
				ownerAddress,
				signingKey,
				httpAddress,
			)
		},
		func(log *types.Log) (any, error) {
			return n.nodeContract.ParseNodeAdded(*log)
		},
		func(event any) {
			nodeAdded, ok := event.(*noderegistry.NodeRegistryNodeAdded)
			if !ok {
				n.logger.Error("node added event is not of type NodesNodeAdded")
				return
			}
			n.logger.Info("node added to registry",
				utils.OriginatorIDField(nodeAdded.NodeId),
				utils.NodeOwnerField(nodeAdded.Owner.Hex()),
				utils.NodeHTTPAddressField(nodeAdded.HttpAddress),
				utils.NodeSigningPublicKeyField(hex.EncodeToString(nodeAdded.SigningPublicKey)),
			)
			nodeID = nodeAdded.NodeId
		},
	)
	if err != nil {
		return 0, err
	}

	return nodeID, nil
}

func (n *nodeRegistryAdmin) AddToNetwork(
	ctx context.Context,
	nodeID uint32,
) error {
	return ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.nodeContract.AddToNetwork(
				opts,
				nodeID,
			)
		},
		func(log *types.Log) (any, error) {
			return n.nodeContract.ParseNodeAddedToCanonicalNetwork(*log)
		},
		func(event any) {
			nodeAdded, ok := event.(*noderegistry.NodeRegistryNodeAddedToCanonicalNetwork)
			if !ok {
				n.logger.Error(
					"unexpected event type, not of type NodesNodeAddedToCanonicalNetwork",
				)
				return
			}
			n.logger.Info("node added to canonical network",
				utils.OriginatorIDField(nodeAdded.NodeId),
			)
		},
	)
}

func (n *nodeRegistryAdmin) RemoveFromNetwork(
	ctx context.Context,
	nodeID uint32,
) error {
	return ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.nodeContract.RemoveFromNetwork(
				opts,
				nodeID,
			)
		},
		func(log *types.Log) (any, error) {
			return n.nodeContract.ParseNodeRemovedFromCanonicalNetwork(*log)
		},
		func(event any) {
			nodeRemoved, ok := event.(*noderegistry.NodeRegistryNodeRemovedFromCanonicalNetwork)
			if !ok {
				n.logger.Error(
					"unexpected event type, not of type NodesNodeRemovedFromCanonicalNetwork",
				)
				return
			}
			n.logger.Info("node removed from canonical network",
				utils.OriginatorIDField(nodeRemoved.NodeId),
			)
		},
	)
}

func (n *nodeRegistryAdmin) SetHTTPAddress(
	ctx context.Context,
	nodeID uint32,
	httpAddress string,
) error {
	return ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.nodeContract.SetHttpAddress(
				opts,
				nodeID,
				httpAddress,
			)
		},
		func(log *types.Log) (any, error) {
			return n.nodeContract.ParseHttpAddressUpdated(*log)
		},
		func(event any) {
			httpAddressUpdated, ok := event.(*noderegistry.NodeRegistryHttpAddressUpdated)
			if !ok {
				n.logger.Error(
					"http address updated event is not of type NodesHttpAddressUpdated",
				)
				return
			}
			n.logger.Info("http address updated",
				utils.OriginatorIDField(httpAddressUpdated.NodeId),
				utils.NodeHTTPAddressField(httpAddressUpdated.HttpAddress),
			)
		},
	)
}

func (n *nodeRegistryAdmin) SetMaxCanonical(
	ctx context.Context,
	limit uint8,
) error {
	err := n.parameterAdmin.SetUint8Parameter(ctx, NodeRegistryMaxCanonicalNodesKey, limit)
	if err != nil {
		return errors.Wrap(err, "failed to update max canonical nodes parameter")
	}

	err = ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.nodeContract.UpdateMaxCanonicalNodes(
				opts,
			)
		},
		func(log *types.Log) (any, error) {
			return n.nodeContract.ParseMaxCanonicalNodesUpdated(*log)
		},
		func(event any) {
			maxCanonicalUpdated, ok := event.(*noderegistry.NodeRegistryMaxCanonicalNodesUpdated)
			if !ok {
				n.logger.Error(
					"unexpected event type, not of type NodeRegistryMaxCanonicalNodesUpdated",
				)
				return
			}
			n.logger.Info("updated max canonical nodes",
				utils.LimitField(maxCanonicalUpdated.MaxCanonicalNodes),
			)
		},
	)
	if err != nil {
		if err.IsNoChange() {
			n.logger.Info("no update needed",
				utils.LimitField(limit),
			)
			return nil
		}
		return err
	}
	return nil
}
