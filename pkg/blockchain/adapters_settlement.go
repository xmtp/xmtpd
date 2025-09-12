package blockchain

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	paramReg "github.com/xmtp/xmtpd/pkg/abi/settlementchainparameterregistry"
)

// Settlement chain adapter for ParameterRegistry.

type settlementRegistryAdapter struct {
	inner *paramReg.SettlementChainParameterRegistry
}

func newSettlementRegistryAdapter(
	addr common.Address,
	backend bind.ContractBackend,
) (*settlementRegistryAdapter, error) {
	inner, err := paramReg.NewSettlementChainParameterRegistry(addr, backend)
	if err != nil {
		return nil, err
	}
	return &settlementRegistryAdapter{inner: inner}, nil
}

func (a *settlementRegistryAdapter) Get(opts *bind.CallOpts, key string) ([32]byte, error) {
	return a.inner.Get(opts, key)
}

func (a *settlementRegistryAdapter) Set(
	opts *bind.TransactOpts,
	key string,
	value [32]byte,
) (*types.Transaction, error) {
	return a.inner.Set(opts, key, value)
}

func (a *settlementRegistryAdapter) SetMany(
	opts *bind.TransactOpts,
	keys []string,
	values [][32]byte,
) (*types.Transaction, error) {
	return a.inner.Set0(opts, keys, values)
}

func (a *settlementRegistryAdapter) ParseParameterSet(l types.Log) (*ParameterSetEvent, error) {
	ev, err := a.inner.ParseParameterSet(l)
	if err != nil {
		return nil, err
	}
	return &ParameterSetEvent{
		Key:   ev.Key.String(),
		Value: ev.Value,
	}, nil
}
