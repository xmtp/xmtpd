package blockchain

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	paramReg "github.com/xmtp/xmtpd/pkg/abi/settlementchainparameterregistry"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

const (
	NODE_REGISTRY_MAX_CANONICAL_NODES_KEY = "xmtp.nodeRegistry.maxCanonicalNodes"
)

type ParameterAdmin struct {
	client            *ethclient.Client
	signer            TransactionSigner
	logger            *zap.Logger
	parameterContract *paramReg.SettlementChainParameterRegistry
}

func NewParameterAdmin(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractsOptions config.ContractsOptions,
) (*ParameterAdmin, error) {
	paramContract, err := paramReg.NewSettlementChainParameterRegistry(
		common.HexToAddress(contractsOptions.SettlementChain.ParameterRegistryAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	return &ParameterAdmin{
		client:            client,
		signer:            signer,
		logger:            logger.Named("ParameterAdmin"),
		parameterContract: paramContract,
	}, nil
}

func (n *ParameterAdmin) SetParameter(
	ctx context.Context,
	paramName string,
	paramValue uint8,
) error {
	err := ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			var value [32]byte
			// store uint8 in the last byte for big-endian compatibility
			value[31] = paramValue
			return n.parameterContract.Set(opts, paramName, value)
		},
		func(log *types.Log) (interface{}, error) {
			return n.parameterContract.ParseParameterSet(*log)
		},
		func(event interface{}) {
			parameterSet, ok := event.(*paramReg.SettlementChainParameterRegistryParameterSet)
			if !ok {
				n.logger.Error(
					"unexpected event type, not of type SettlementChainParameterRegistryParameterSet",
				)
				return
			}
			n.logger.Info("set parameter",
				zap.String("key", NODE_REGISTRY_MAX_CANONICAL_NODES_KEY),
				zap.Uint64("parameter", utils.DecodeBytes32ToUint64(parameterSet.Value)),
			)
		},
	)
	if err != nil {
		return err
	}

	return nil
}
