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
		minMonthlyFeeMicroDollars int64,
	) error
	AddToNetwork(ctx context.Context, nodeId int64) error
	RemoveFromNetwork(ctx context.Context, nodeId int64) error
	SetHttpAddress(ctx context.Context, nodeId int64, httpAddress string) error
	SetMinMonthlyFee(ctx context.Context, nodeId int64, minMonthlyFeeMicroDollars int64) error
	SetMaxActiveNodes(ctx context.Context, maxActiveNodes uint8) error
	SetNodeOperatorCommissionPercent(ctx context.Context, commissionPercent int64) error
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
	minMonthlyFeeMicroDollars int64,
) error {
	if !common.IsHexAddress(owner) {
		return fmt.Errorf("invalid owner address provided %s", owner)
	}

	if minMonthlyFeeMicroDollars < 0 {
		return fmt.Errorf("invalid min monthly fee provided %d", minMonthlyFeeMicroDollars)
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
				big.NewInt(minMonthlyFeeMicroDollars),
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
				zap.Uint64("node_id", nodeAdded.NodeId.Uint64()),
				zap.String("owner", nodeAdded.Owner.Hex()),
				zap.String("http_address", nodeAdded.HttpAddress),
				zap.String("signing_key_pub", hex.EncodeToString(nodeAdded.SigningKeyPub)),
				zap.String("min_monthly_fee", nodeAdded.MinMonthlyFeeMicroDollars.String()),
			)
		},
	)
}

func (n *nodeRegistryAdmin) AddToNetwork(
	ctx context.Context,
	nodeId int64,
) error {
	return ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.AddToNetwork(
				opts,
				big.NewInt(nodeId),
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
				zap.Uint64("node_id", nodeAdded.NodeId.Uint64()),
			)
		},
	)
}

func (n *nodeRegistryAdmin) RemoveFromNetwork(
	ctx context.Context,
	nodeId int64,
) error {
	return ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.RemoveFromNetwork(
				opts,
				big.NewInt(nodeId),
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
				zap.Uint64("node_id", nodeRemoved.NodeId.Uint64()),
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
			httpAddressUpdated, ok := event.(*noderegistry.NodeRegistryHttpAddressUpdated)
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
	minMonthlyFeeMicroDollars int64,
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
				big.NewInt(minMonthlyFeeMicroDollars),
			)
		},
		func(log *types.Log) (interface{}, error) {
			return n.contract.ParseMinMonthlyFeeUpdated(*log)
		},
		func(event interface{}) {
			minMonthlyFeeUpdated, ok := event.(*noderegistry.NodeRegistryMinMonthlyFeeUpdated)
			if !ok {
				n.logger.Error(
					"min monthly fee updated event is not of type NodesMinMonthlyFeeUpdated",
				)
				return
			}
			n.logger.Info(
				"min monthly fee updated",
				zap.Uint64("node_id", minMonthlyFeeUpdated.NodeId.Uint64()),
				zap.String(
					"min_monthly_fee",
					minMonthlyFeeUpdated.MinMonthlyFeeMicroDollars.String(),
				),
			)
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
			maxActiveNodesUpdated, ok := event.(*noderegistry.NodeRegistryMaxActiveNodesUpdated)
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
			nodeOperatorCommissionPercentUpdated, ok := event.(*noderegistry.NodeRegistryNodeOperatorCommissionPercentUpdated)
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
