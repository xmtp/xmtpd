package blockchain

import (
	"context"
	"fmt"

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
	GROUP_MESSAGE_PAYLOAD_BOOTSTRAPPER_KEY   = "xmtp.groupMessageBroadcaster.payloadBootstrapper"
	IDENTITY_UPDATE_PAUSED_KEY               = "xmtp.identityUpdateBroadcaster.paused"
	GROUP_MESSAGE_PAUSED_KEY                 = "xmtp.groupMessageBroadcaster.paused"
	APP_CHAIN_GATEWAY_PAUSED_KEY             = "xmtp.appChainGateway.paused"
	SETTLEMENT_CHAIN_GATEWAY_PAUSED_KEY      = "xmtp.settlementChainGateway.paused"
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
) (common.Address, ProtocolError) {
	payload, err := n.parameterContract.Get(&bind.CallOpts{
		Context: ctx,
	}, paramName)
	if err != nil {
		return common.Address{}, NewBlockchainError(err)
	}

	return common.BytesToAddress(payload[:]), nil
}

func (n *ParameterAdmin) GetParameterUint8(
	ctx context.Context,
	paramName string,
) (uint8, ProtocolError) {
	payload, err := n.parameterContract.Get(&bind.CallOpts{
		Context: ctx,
	}, paramName)
	if err != nil {
		return 0, NewBlockchainError(err)
	}

	return payload[31], nil
}

func (n *ParameterAdmin) GetParameterBool(
	ctx context.Context,
	paramName string,
) (bool, ProtocolError) {
	payload, err := n.parameterContract.Get(&bind.CallOpts{
		Context: ctx,
	}, paramName)
	if err != nil {
		return false, NewBlockchainError(err)
	}

	b, err := decodeBool(payload)
	if err != nil {
		return false, NewBlockchainError(err)
	}

	return b, nil
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

func packBool(b bool) [32]byte {
	var out [32]byte
	if b {
		out[31] = 1
	}
	return out
}

// decodeBool expects the canonical encoding produced by packBool.
// It returns (bool, nil) for 0x00/0x01 in the last byte and errors otherwise.
func decodeBool(val [32]byte) (bool, error) {
	v := val[31]
	// Ensure normalization: all other bytes should be zero.
	for i := 0; i < 31; i++ {
		if val[i] != 0 {
			return false, fmt.Errorf("non-canonical bool encoding in bytes32 (non-zero prefix)")
		}
	}
	switch v {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, fmt.Errorf("invalid bool encoding: last byte = %d (want 0 or 1)", v)
	}
}

// shared executor -------------------------------------------------------------

func (n *ParameterAdmin) setParameterBytes32(
	ctx context.Context,
	paramName string,
	value [32]byte,
	onEvent func(val [32]byte),
) ProtocolError {
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
) ProtocolError {
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
) ProtocolError {
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

func (n *ParameterAdmin) SetBoolParameter(
	ctx context.Context,
	paramName string,
	paramValue bool,
) ProtocolError {
	return n.setParameterBytes32(ctx, paramName, packBool(paramValue),
		func(val [32]byte) {
			b, err := decodeBool(val)
			if err != nil {
				n.logger.Warn("set bool parameter (non-canonical value observed in event)",
					zap.String("key", paramName),
					zap.Error(err),
				)
				return
			}
			n.logger.Info("set bool parameter",
				zap.String("key", paramName),
				zap.Bool("value", b),
			)
		},
	)
}
