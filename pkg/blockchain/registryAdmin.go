package blockchain

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/contracts/pkg/nodes"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

type INodeRegistryAdmin interface {
	AddNode(
		ctx context.Context,
		owner string,
		signingKeyPub *ecdsa.PublicKey,
		httpAddress string,
		minMonthlyFee int64,
	) error
	DisableNode(ctx context.Context, nodeId int64) error
	EnableNode(ctx context.Context, nodeId int64) error
	RemoveFromApiNodes(ctx context.Context, nodeId int64) error
	RemoveFromReplicationNodes(ctx context.Context, nodeId int64) error
	SetHttpAddress(ctx context.Context, nodeId int64, httpAddress string) error
	SetMinMonthlyFee(ctx context.Context, nodeId int64, minMonthlyFee int64) error
	SetIsApiEnabled(ctx context.Context, nodeId int64, isApiEnabled bool) error
	SetIsReplicationEnabled(ctx context.Context, nodeId int64, isReplicationEnabled bool) error
	SetMaxActiveNodes(ctx context.Context, maxActiveNodes uint8) error
	SetNodeOperatorCommissionPercent(ctx context.Context, commissionPercent int64) error
}

type nodeRegistryAdmin struct {
	client   *ethclient.Client
	signer   TransactionSigner
	logger   *zap.Logger
	contract *nodes.Nodes
}

var _ INodeRegistryAdmin = &nodeRegistryAdmin{}

func NewNodeRegistryAdmin(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractsOptions config.ContractsOptions,
) (*nodeRegistryAdmin, error) {
	contract, err := nodes.NewNodes(
		common.HexToAddress(contractsOptions.NodesContractAddress),
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
	minMonthlyFee int64,
) error {
	if !common.IsHexAddress(owner) {
		return fmt.Errorf("invalid owner address provided %s", owner)
	}

	if minMonthlyFee < 0 {
		return fmt.Errorf("invalid min monthly fee provided %d", minMonthlyFee)
	}

	ownerAddress := common.HexToAddress(owner)
	signingKey := crypto.FromECDSAPub(signingKeyPub)

	if n.signer == nil {
		return fmt.Errorf("no signer provided")
	}

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
				big.NewInt(minMonthlyFee),
			)
		},
		func(log *types.Log) (interface{}, error) {
			return n.contract.ParseNodeAdded(*log)
		},
		func(event interface{}) {
			nodeAdded, ok := event.(*nodes.NodesNodeAdded)
			if !ok {
				n.logger.Error("node added event is not of type NodesNodeAdded")
				return
			}
			n.logger.Info("node added to registry",
				zap.Uint64("node_id", nodeAdded.NodeId.Uint64()),
				zap.String("owner", nodeAdded.Owner.Hex()),
				zap.String("http_address", nodeAdded.HttpAddress),
				zap.String("signing_key_pub", hex.EncodeToString(nodeAdded.SigningKeyPub)),
				zap.String("min_monthly_fee", nodeAdded.MinMonthlyFee.String()),
			)
		},
	)
}

func (n *nodeRegistryAdmin) DisableNode(ctx context.Context, nodeId int64) error {
	return ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.DisableNode(opts, big.NewInt(nodeId))
		},
		func(log *types.Log) (interface{}, error) {
			return n.contract.ParseNodeDisabled(*log)
		},
		func(event interface{}) {
			nodeDisabled, ok := event.(*nodes.NodesNodeDisabled)
			if !ok {
				n.logger.Error("node disabled event is not of type NodesNodeDisabled")
				return
			}
			n.logger.Info("node disabled",
				zap.Uint64("node_id", nodeDisabled.NodeId.Uint64()),
			)
		},
	)
}

func (n *nodeRegistryAdmin) EnableNode(ctx context.Context, nodeId int64) error {
	return ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.EnableNode(opts, big.NewInt(nodeId))
		},
		func(log *types.Log) (interface{}, error) {
			return n.contract.ParseNodeEnabled(*log)
		},
		func(event interface{}) {
			nodeEnabled, ok := event.(*nodes.NodesNodeEnabled)
			if !ok {
				n.logger.Error("node enabled event is not of type NodesNodeEnabled")
				return
			}
			n.logger.Info("node enabled",
				zap.Uint64("node_id", nodeEnabled.NodeId.Uint64()),
			)
		},
	)
}

func (n *nodeRegistryAdmin) RemoveFromApiNodes(ctx context.Context, nodeId int64) error {
	return ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.RemoveFromApiNodes(opts, big.NewInt(nodeId))
		},
		func(log *types.Log) (interface{}, error) {
			return n.contract.ParseApiDisabled(*log)
		},
		func(event interface{}) {
			nodeRemovedFromApiNodes, ok := event.(*nodes.NodesApiDisabled)
			if !ok {
				n.logger.Error(
					"node removed from active api nodes event is not of type NodesApiDisabled",
				)
				return
			}
			n.logger.Info("node removed from active api nodes",
				zap.Uint64("node_id", nodeRemovedFromApiNodes.NodeId.Uint64()),
			)
		},
	)
}

func (n *nodeRegistryAdmin) RemoveFromReplicationNodes(ctx context.Context, nodeId int64) error {
	return ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.RemoveFromReplicationNodes(opts, big.NewInt(nodeId))
		},
		func(log *types.Log) (interface{}, error) {
			return n.contract.ParseReplicationDisabled(*log)
		},
		func(event interface{}) {
			nodeRemovedFromReplicationNodes, ok := event.(*nodes.NodesReplicationDisabled)
			if !ok {
				n.logger.Error(
					"node removed from active replication nodes event is not of type NodesReplicationDisabled",
				)
				return
			}
			n.logger.Info("node removed from active replication nodes",
				zap.Uint64("node_id", nodeRemovedFromReplicationNodes.NodeId.Uint64()),
			)
		},
	)
}

func (n *nodeRegistryAdmin) SetHttpAddress(
	ctx context.Context,
	nodeId int64,
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
				big.NewInt(nodeId),
				httpAddress,
			)
		},
		func(log *types.Log) (interface{}, error) {
			return n.contract.ParseHttpAddressUpdated(*log)
		},
		func(event interface{}) {
			httpAddressUpdated, ok := event.(*nodes.NodesHttpAddressUpdated)
			if !ok {
				n.logger.Error(
					"http address updated event is not of type NodesHttpAddressUpdated",
				)
				return
			}
			n.logger.Info("http address updated",
				zap.Uint64("node_id", httpAddressUpdated.NodeId.Uint64()),
				zap.String("http_address", httpAddressUpdated.NewHttpAddress),
			)
		},
	)
}

func (n *nodeRegistryAdmin) SetMinMonthlyFee(
	ctx context.Context,
	nodeId int64,
	minMonthlyFee int64,
) error {
	return ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.SetMinMonthlyFee(
				opts,
				big.NewInt(nodeId),
				big.NewInt(minMonthlyFee),
			)
		},
		func(log *types.Log) (interface{}, error) {
			return n.contract.ParseMinMonthlyFeeUpdated(*log)
		},
		func(event interface{}) {
			minMonthlyFeeUpdated, ok := event.(*nodes.NodesMinMonthlyFeeUpdated)
			if !ok {
				n.logger.Error(
					"min monthly fee updated event is not of type NodesMinMonthlyFeeUpdated",
				)
				return
			}
			n.logger.Info("min monthly fee updated",
				zap.Uint64("node_id", minMonthlyFeeUpdated.NodeId.Uint64()),
				zap.String("min_monthly_fee", minMonthlyFeeUpdated.MinMonthlyFee.String()),
			)
		},
	)
}

func (n *nodeRegistryAdmin) SetIsApiEnabled(
	ctx context.Context,
	nodeId int64,
	isApiEnabled bool,
) error {
	return ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.SetIsApiEnabled(opts, big.NewInt(nodeId), isApiEnabled)
		},
		func(log *types.Log) (interface{}, error) {
			if isApiEnabled {
				return n.contract.ParseApiEnabled(*log)
			}
			return n.contract.ParseApiDisabled(*log)
		},
		func(event interface{}) {
			if isApiEnabled {
				apiEnabled := event.(*nodes.NodesApiEnabled)
				n.logger.Info("api enabled",
					zap.Uint64("node_id", apiEnabled.NodeId.Uint64()),
				)
			} else {
				apiDisabled := event.(*nodes.NodesApiDisabled)
				n.logger.Info("api disabled",
					zap.Uint64("node_id", apiDisabled.NodeId.Uint64()),
				)
			}
		},
	)
}

func (n *nodeRegistryAdmin) SetIsReplicationEnabled(
	ctx context.Context,
	nodeId int64,
	isReplicationEnabled bool,
) error {
	return ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.SetIsReplicationEnabled(
				opts,
				big.NewInt(nodeId),
				isReplicationEnabled,
			)
		},
		func(log *types.Log) (interface{}, error) {
			if isReplicationEnabled {
				return n.contract.ParseReplicationEnabled(*log)
			}
			return n.contract.ParseReplicationDisabled(*log)
		},
		func(event interface{}) {
			if isReplicationEnabled {
				replicationEnabled := event.(*nodes.NodesReplicationEnabled)
				n.logger.Info("replication enabled",
					zap.Uint64("node_id", replicationEnabled.NodeId.Uint64()),
				)
			} else {
				replicationDisabled := event.(*nodes.NodesReplicationDisabled)
				n.logger.Info("replication disabled",
					zap.Uint64("node_id", replicationDisabled.NodeId.Uint64()),
				)
			}
		},
	)
}

func (n *nodeRegistryAdmin) SetMaxActiveNodes(ctx context.Context, maxActiveNodes uint8) error {
	return ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.SetMaxActiveNodes(opts, maxActiveNodes)
		},
		func(log *types.Log) (interface{}, error) {
			return n.contract.ParseMaxActiveNodesUpdated(*log)
		},
		func(event interface{}) {
			maxActiveNodesUpdated, ok := event.(*nodes.NodesMaxActiveNodesUpdated)
			if !ok {
				n.logger.Error(
					"max active nodes updated event is not of type NodesMaxActiveNodesUpdated",
				)
				return
			}
			n.logger.Info("max active nodes set",
				zap.Uint8("max_active_nodes", maxActiveNodesUpdated.NewMaxActiveNodes),
			)
		},
	)
}

func (n *nodeRegistryAdmin) SetNodeOperatorCommissionPercent(
	ctx context.Context,
	commissionPercent int64,
) error {
	if commissionPercent < 0 || commissionPercent > 10000 {
		return fmt.Errorf("invalid commission percent provided %d", commissionPercent)
	}

	return ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.SetNodeOperatorCommissionPercent(
				opts,
				big.NewInt(commissionPercent),
			)
		},
		func(log *types.Log) (interface{}, error) {
			return n.contract.ParseNodeOperatorCommissionPercentUpdated(*log)
		},
		func(event interface{}) {
			nodeOperatorCommissionPercentUpdated, ok := event.(*nodes.NodesNodeOperatorCommissionPercentUpdated)
			if !ok {
				n.logger.Error(
					"node operator commission percent updated event is not of type NodesNodeOperatorCommissionPercentUpdated",
				)
				return
			}
			n.logger.Info(
				"node operator commission percent updated",
				zap.Uint64(
					"node_operator_commission_percent",
					nodeOperatorCommissionPercentUpdated.NewCommissionPercent.Uint64(),
				),
			)
		},
	)
}
