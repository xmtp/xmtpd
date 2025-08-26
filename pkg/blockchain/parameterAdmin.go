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
	NODE_REGISTRY_MAX_CANONICAL_NODES_KEY    = "xmtp.nodeRegistry.maxCanonicalNodes"
	IDENTITY_UPDATE_PAYLOAD_BOOTSTRAPPER_KEY = "xmtp.identityUpdateBroadcaster.payloadBootstrapper"
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

func (n *ParameterAdmin) GetParameterAddress(
	ctx context.Context,
	paramName string,
) (common.Address, error) {
	payload, err := n.parameterContract.Get(&bind.CallOpts{
		Context: ctx,
	}, paramName)
	if err != nil {
		return common.Address{}, err
	}

	return common.BytesToAddress(payload[:]), nil
}

func (n *ParameterAdmin) GetParameterUint8(
	ctx context.Context,
	paramName string,
) (uint8, error) {
	payload, err := n.parameterContract.Get(&bind.CallOpts{
		Context: ctx,
	}, paramName)
	if err != nil {
		return 0, err
	}

	return payload[31], nil
}

// Param helpers ---------------------------------------------------------------

func packUint8(v uint8) [32]byte {
	var out [32]byte
	out[31] = v // big-endian placement
	return out
}

func packAddress(a common.Address) [32]byte {
	var out [32]byte
	copy(out[12:], a.Bytes()) // right-align to 32 bytes
	return out
}

// shared executor -------------------------------------------------------------

func (n *ParameterAdmin) setParameterBytes32(
	ctx context.Context,
	paramName string,
	value [32]byte,
	onEvent func(val [32]byte),
) error {
	return ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
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
			if onEvent != nil {
				onEvent(parameterSet.Value)
			}
		},
	)
}

// typed wrappers --------------------------------------------------------------

func (n *ParameterAdmin) SetUint8Parameter(
	ctx context.Context,
	paramName string,
	paramValue uint8,
) error {
	return n.setParameterBytes32(ctx, paramName, packUint8(paramValue),
		func(val [32]byte) {
			n.logger.Info("set uint8 parameter",
				zap.String("key", paramName),
				zap.Uint64("value", utils.DecodeBytes32ToUint64(val)),
			)
		},
	)
}

func (n *ParameterAdmin) SetAddressParameter(
	ctx context.Context,
	paramName string,
	paramValue common.Address,
) error {
	return n.setParameterBytes32(ctx, paramName, packAddress(paramValue),
		func(val [32]byte) {
			addr := common.BytesToAddress(val[12:])
			n.logger.Info("set address parameter",
				zap.String("key", paramName),
				zap.String("address", addr.Hex()),
			)
		},
	)
}
