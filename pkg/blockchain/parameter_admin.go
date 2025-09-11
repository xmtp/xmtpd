package blockchain

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/xmtp/xmtpd/pkg/utils"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)


type ParameterSetEvent struct {
	Key   string
	Value [32]byte
}

// Registry interface both adapters satisfy.

type ParameterRegistry interface {
	Get(opts *bind.CallOpts, key string) ([32]byte, error)
	Set(opts *bind.TransactOpts, key string, value [32]byte) (*types.Transaction, error)
	SetMany(opts *bind.TransactOpts, keys []string, values [][32]byte) (*types.Transaction, error)
	ParseParameterSet(log types.Log) (*ParameterSetEvent, error)
}

// DTOs / small types

type Uint64Param struct {
	Name  string
	Value uint64
}

// ParameterAdmin is chain-agnostic and works with any ParameterRegistry adapter.
type ParameterAdmin struct {
	client   *ethclient.Client
	signer   TransactionSigner
	logger   *zap.Logger
	registry ParameterRegistry
}

func NewSettlementParameterAdmin(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractsOptions config.ContractsOptions,
) (*ParameterAdmin, error) {
	addr := common.HexToAddress(contractsOptions.SettlementChain.ParameterRegistryAddress)
	reg, err := newSettlementRegistryAdapter(addr, client)
	if err != nil {
		return nil, err
	}
	return &ParameterAdmin{
		client:   client,
		signer:   signer,
		logger:   logger.Named("ParameterAdmin[settlement]"),
		registry: reg,
	}, nil
}

func NewAppchainParameterAdmin(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractsOptions config.ContractsOptions,
) (*ParameterAdmin, error) {
	addr := common.HexToAddress(contractsOptions.AppChain.ParameterRegistryAddress)
	reg, err := newAppchainRegistryAdapter(addr, client)
	if err != nil {
		return nil, err
	}
	return &ParameterAdmin{
		client:   client,
		signer:   signer,
		logger:   logger.Named("ParameterAdmin[appchain]"),
		registry: reg,
	}, nil
}

// -------------------------- Reads --------------------------

func (n *ParameterAdmin) GetParameterAddress(ctx context.Context, paramName string) (common.Address, ProtocolError) {
	payload, err := n.registry.Get(&bind.CallOpts{Context: ctx}, paramName)
	if err != nil {
		return common.Address{}, NewBlockchainError(err)
	}
	return common.BytesToAddress(payload[:]), nil
}

func (n *ParameterAdmin) GetParameterUint8(ctx context.Context, paramName string) (uint8, ProtocolError) {
	payload, err := n.registry.Get(&bind.CallOpts{Context: ctx}, paramName)
	if err != nil {
		return 0, NewBlockchainError(err)
	}
	v, derr := decodeUint8(payload)
	if derr != nil {
		return 0, NewBlockchainError(derr)
	}
	return v, nil
}

func (n *ParameterAdmin) GetParameterUint16(ctx context.Context, paramName string) (uint16, ProtocolError) {
	payload, err := n.registry.Get(&bind.CallOpts{Context: ctx}, paramName)
	if err != nil {
		return 0, NewBlockchainError(err)
	}
	v, derr := decodeUint16(payload)
	if derr != nil {
		return 0, NewBlockchainError(derr)
	}
	return v, nil
}

func (n *ParameterAdmin) GetParameterUint32(ctx context.Context, paramName string) (uint32, ProtocolError) {
	payload, err := n.registry.Get(&bind.CallOpts{Context: ctx}, paramName)
	if err != nil {
		return 0, NewBlockchainError(err)
	}
	v, derr := decodeUint32(payload)
	if derr != nil {
		return 0, NewBlockchainError(derr)
	}
	return v, nil
}

func (n *ParameterAdmin) GetParameterUint64(ctx context.Context, paramName string) (uint64, ProtocolError) {
	payload, err := n.registry.Get(&bind.CallOpts{Context: ctx}, paramName)
	if err != nil {
		return 0, NewBlockchainError(err)
	}
	v, derr := decodeUint64(payload)
	if derr != nil {
		return 0, NewBlockchainError(derr)
	}
	return v, nil
}

func (n *ParameterAdmin) GetParameterUint96(ctx context.Context, paramName string) (*big.Int, ProtocolError) {
	payload, err := n.registry.Get(&bind.CallOpts{Context: ctx}, paramName)
	if err != nil {
		return nil, NewBlockchainError(err)
	}
	u, derr := decodeUint96Big(payload)
	if derr != nil {
		return nil, NewBlockchainError(derr)
	}
	return u, nil
}

func (n *ParameterAdmin) GetParameterBool(ctx context.Context, paramName string) (bool, ProtocolError) {
	payload, err := n.registry.Get(&bind.CallOpts{Context: ctx}, paramName)
	if err != nil {
		return false, NewBlockchainError(err)
	}
	b, derr := decodeBool(payload)
	if derr != nil {
		return false, NewBlockchainError(derr)
	}
	return b, nil
}

// -------------------------- Writes --------------------------

func (n *ParameterAdmin) SetUint8Parameter(ctx context.Context, name string, v uint8) ProtocolError {
	return n.setParameterBytes32(ctx, name, packUint8(v), n.logUint8(name))
}

func (n *ParameterAdmin) SetUint16Parameter(ctx context.Context, name string, v uint16) ProtocolError {
	return n.setParameterBytes32(ctx, name, packUint16(v), n.logUint16(name))
}

func (n *ParameterAdmin) SetUint32Parameter(ctx context.Context, name string, v uint32) ProtocolError {
	return n.setParameterBytes32(ctx, name, packUint32(v), n.logUint32(name))
}

func (n *ParameterAdmin) SetUint64Parameter(ctx context.Context, name string, v uint64) ProtocolError {
	return n.setParameterBytes32(ctx, name, packUint64(v), n.logUint64(name))
}

func (n *ParameterAdmin) SetUint96Parameter(ctx context.Context, name string, v *big.Int) ProtocolError {
	enc, err := packUint96Big(v)
	if err != nil {
		return NewBlockchainError(err)
	}
	return n.setParameterBytes32(ctx, name, enc, n.logUint96(name))
}

func (n *ParameterAdmin) SetAddressParameter(ctx context.Context, name string, addr common.Address) ProtocolError {
	return n.setParameterBytes32(ctx, name, packAddress(addr), n.logAddress(name))
}

func (n *ParameterAdmin) SetBoolParameter(ctx context.Context, name string, v bool) ProtocolError {
	return n.setParameterBytes32(ctx, name, packBool(v), n.logBool(name))
}

func (n *ParameterAdmin) SetManyUint64Parameters(ctx context.Context, items []Uint64Param) ProtocolError {
	keys := make([]string, len(items))
	vals := make([][32]byte, len(items))
	for i, it := range items {
		keys[i] = it.Name
		vals[i] = packUint64(it.Value)
	}
	return n.setParametersBytes32Many(ctx, keys, vals)
}

func (n *ParameterAdmin) logUint8(paramName string) func([32]byte) {
	return func(val [32]byte) {
		u8, err := decodeUint8(val)
		if err != nil {
			n.logger.Warn("set uint8 parameter (non-canonical value observed in event)",
				zap.String("key", paramName),
				zap.Error(err),
			)
			return
		}
		n.logger.Info("set uint8 parameter",
			zap.String("key", paramName),
			zap.Uint8("value", u8),
		)
	}
}

func (n *ParameterAdmin) logUint16(paramName string) func([32]byte) {
	return func(val [32]byte) {
		u16, err := decodeUint16(val)
		if err != nil {
			n.logger.Warn("set uint16 parameter (non-canonical value observed in event)",
				zap.String("key", paramName),
				zap.Error(err),
			)
			return
		}
		n.logger.Info("set uint16 parameter",
			zap.String("key", paramName),
			zap.Uint16("value", u16),
		)
	}
}

func (n *ParameterAdmin) logUint32(paramName string) func([32]byte) {
	return func(val [32]byte) {
		u32, err := decodeUint32(val)
		if err != nil {
			n.logger.Warn("set uint32 parameter (non-canonical value observed in event)",
				zap.String("key", paramName),
				zap.Error(err),
			)
			return
		}
		n.logger.Info("set uint32 parameter",
			zap.String("key", paramName),
			zap.Uint32("value", u32),
		)
	}
}

func (n *ParameterAdmin) logUint64(paramName string) func([32]byte) {
	return func(val [32]byte) {
		u64, err := decodeUint64(val)
		if err != nil {
			n.logger.Warn("set uint64 parameter (non-canonical value observed in event)",
				zap.String("key", paramName),
				zap.Error(err),
			)
			return
		}
		n.logger.Info("set uint64 parameter",
			zap.String("key", paramName),
			zap.Uint64("value", u64),
		)
	}
}

func (n *ParameterAdmin) logUint96(paramName string) func([32]byte) {
	return func(val [32]byte) {
		u, err := decodeUint96Big(val)
		if err != nil {
			n.logger.Warn("set uint96 parameter (non-canonical value observed in event)",
				zap.String("key", paramName),
				zap.Error(err),
			)
			return
		}
		n.logger.Info("set uint96 parameter",
			zap.String("key", paramName),
			zap.String("value", u.String()),
		)
	}
}

func (n *ParameterAdmin) logAddress(paramName string) func([32]byte) {
	return func(val [32]byte) {
		addr := common.BytesToAddress(val[12:])
		n.logger.Info("set address parameter",
			zap.String("key", paramName),
			zap.String("address", addr.Hex()),
		)
	}
}

func (n *ParameterAdmin) logBool(paramName string) func([32]byte) {
	return func(val [32]byte) {
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
	}
}


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
			return n.registry.Set(opts, paramName, value)
		},
		func(log *types.Log) (interface{}, error) {
			return n.registry.ParseParameterSet(*log)
		},
		func(event interface{}) {
			ev, ok := event.(*ParameterSetEvent)
			if !ok {
				n.logger.Error("unexpected event type, want ParameterSetEvent")
				return
			}
			if onEvent != nil {
				onEvent(ev.Value)
			}
		},
	)
}

func (n *ParameterAdmin) setParametersBytes32Many(
	ctx context.Context,
	keys []string,
	values [][32]byte,
) ProtocolError {
	return ExecuteTransaction(
		ctx,
		n.signer,
		n.logger,
		n.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return n.registry.SetMany(opts, keys, values)
		},
		func(log *types.Log) (interface{}, error) {
			return n.registry.ParseParameterSet(*log)
		},
		func(event interface{}) {
			ev, ok := event.(*ParameterSetEvent)
			if !ok {
				n.logger.Error("unexpected event type, want ParameterSetEvent")
				return
			}
			n.logger.Info("set parameter (batch)",
				zap.String("key", ev.Key),
				zap.Uint64("value", utils.DecodeBytes32ToUint64(ev.Value)),
			)
		},
	)
}