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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/contracts/pkg/nodesv2"
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
	SetMaxActiveNodes(ctx context.Context, maxActiveNodes uint8) error
	SetNodeOperatorCommissionPercent(ctx context.Context, commissionPercent int64) error
	SetBaseURI(ctx context.Context, baseURI string) error
}

type nodeRegistryAdmin struct {
	client   *ethclient.Client
	signer   TransactionSigner
	logger   *zap.Logger
	contract *nodesv2.NodesV2
}

var _ INodeRegistryAdmin = &nodeRegistryAdmin{}

func NewNodeRegistryAdmin(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractsOptions config.ContractsOptions,
) (*nodeRegistryAdmin, error) {
	contract, err := nodesv2.NewNodesV2(
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

	return n.executeTransaction(
		ctx,
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
			nodeAdded := event.(*nodesv2.NodesV2NodeAdded)
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
	return n.executeTransaction(
		ctx,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.DisableNode(opts, big.NewInt(nodeId))
		},
		func(log *types.Log) (interface{}, error) {
			return n.contract.ParseNodeDisabled(*log)
		},
		func(event interface{}) {
			nodeDisabled := event.(*nodesv2.NodesV2NodeDisabled)
			n.logger.Info("node disabled",
				zap.Uint64("node_id", nodeDisabled.NodeId.Uint64()),
			)
		},
	)
}

func (n *nodeRegistryAdmin) EnableNode(ctx context.Context, nodeId int64) error {
	return n.executeTransaction(
		ctx,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.EnableNode(opts, big.NewInt(nodeId))
		},
		func(log *types.Log) (interface{}, error) {
			return n.contract.ParseNodeEnabled(*log)
		},
		func(event interface{}) {
			nodeEnabled := event.(*nodesv2.NodesV2NodeEnabled)
			n.logger.Info("node enabled",
				zap.Uint64("node_id", nodeEnabled.NodeId.Uint64()),
			)
		},
	)
}

func (n *nodeRegistryAdmin) RemoveFromApiNodes(ctx context.Context, nodeId int64) error {
	return n.executeTransaction(
		ctx,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.RemoveFromApiNodes(opts, big.NewInt(nodeId))
		},
		func(log *types.Log) (interface{}, error) {
			return n.contract.ParseApiDisabled(*log)
		},
		func(event interface{}) {
			nodeRemovedFromApiNodes := event.(*nodesv2.NodesV2ApiDisabled)
			n.logger.Info("node removed from active api nodes",
				zap.Uint64("node_id", nodeRemovedFromApiNodes.NodeId.Uint64()),
			)
		},
	)
}

func (n *nodeRegistryAdmin) RemoveFromReplicationNodes(ctx context.Context, nodeId int64) error {
	return n.executeTransaction(
		ctx,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.RemoveFromReplicationNodes(opts, big.NewInt(nodeId))
		},
		func(log *types.Log) (interface{}, error) {
			return n.contract.ParseReplicationDisabled(*log)
		},
		func(event interface{}) {
			nodeRemovedFromReplicationNodes := event.(*nodesv2.NodesV2ReplicationDisabled)
			n.logger.Info("node removed from active replication nodes",
				zap.Uint64("node_id", nodeRemovedFromReplicationNodes.NodeId.Uint64()),
			)
		},
	)
}

func (n *nodeRegistryAdmin) SetMaxActiveNodes(ctx context.Context, maxActiveNodes uint8) error {
	return n.executeTransaction(
		ctx,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.SetMaxActiveNodes(opts, maxActiveNodes)
		},
		func(log *types.Log) (interface{}, error) {
			return n.contract.ParseMaxActiveNodesUpdated(*log)
		},
		func(event interface{}) {
			maxActiveNodesUpdated := event.(*nodesv2.NodesV2MaxActiveNodesUpdated)
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

	return n.executeTransaction(
		ctx,
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
			nodeOperatorCommissionPercentUpdated := event.(*nodesv2.NodesV2NodeOperatorCommissionPercentUpdated)
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

func (n *nodeRegistryAdmin) SetBaseURI(ctx context.Context, baseURI string) error {
	return n.executeTransaction(
		ctx,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.contract.SetBaseURI(opts, baseURI)
		},
		func(log *types.Log) (interface{}, error) {
			return n.contract.ParseBaseURIUpdated(*log)
		},
		func(event interface{}) {
			baseURIUpdated := event.(*nodesv2.NodesV2BaseURIUpdated)
			n.logger.Info("base uri updated",
				zap.String("base_uri", baseURIUpdated.NewBaseURI),
			)
		},
	)
}

// executeTransaction is a helper function that:
// - executes a transaction
// - waits for it to be mined
// - processes the event logs
func (n *nodeRegistryAdmin) executeTransaction(
	ctx context.Context,
	txFunc func(*bind.TransactOpts) (*types.Transaction, error),
	eventParser func(*types.Log) (interface{}, error),
	logHandler func(interface{}),
) error {
	if n.signer == nil {
		return fmt.Errorf("no signer provided")
	}

	tx, err := txFunc(&bind.TransactOpts{
		Context: ctx,
		From:    n.signer.FromAddress(),
		Signer:  n.signer.SignerFunc(),
	})
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
		event, err := eventParser(log)
		if err != nil {
			continue
		}
		logHandler(event)
	}

	return nil
}
