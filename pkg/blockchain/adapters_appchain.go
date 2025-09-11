package blockchain

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	appParamReg "github.com/xmtp/xmtpd/pkg/abi/appchainparameterregistry"
)

type appchainRegistryAdapter struct {
	inner *appParamReg.AppChainParameterRegistry
}

func newAppchainRegistryAdapter(addr common.Address, backend bind.ContractBackend) (*appchainRegistryAdapter, error) {
	inner, err := appParamReg.NewAppChainParameterRegistry(addr, backend)
	if err != nil {
		return nil, err
	}
	return &appchainRegistryAdapter{inner: inner}, nil
}

func (a *appchainRegistryAdapter) Get(opts *bind.CallOpts, key string) ([32]byte, error) {
	return a.inner.Get(opts, key)
}

func (a *appchainRegistryAdapter) Set(opts *bind.TransactOpts, key string, value [32]byte) (*types.Transaction, error) {
	return a.inner.Set(opts, key, value)
}

func (a *appchainRegistryAdapter) SetMany(opts *bind.TransactOpts, keys []string, values [][32]byte) (*types.Transaction, error) {
	return a.inner.Set0(opts, keys, values)
}

func (a *appchainRegistryAdapter) ParseParameterSet(l types.Log) (*ParameterSetEvent, error) {
	ev, err := a.inner.ParseParameterSet(l)
	if err != nil {
		return nil, err
	}
	return &ParameterSetEvent{
		Key:   ev.Key.String(),
		Value: ev.Value,
	}, nil
}
