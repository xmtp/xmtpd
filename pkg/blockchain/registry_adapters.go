package blockchain

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	appReg "github.com/xmtp/xmtpd/pkg/abi/appchainparameterregistry"
	settleReg "github.com/xmtp/xmtpd/pkg/abi/settlementchainparameterregistry"
)

// -----------------------------
// SettlementChain adapter
// -----------------------------

type settlementRegistryAdapter struct {
	inner *settleReg.SettlementChainParameterRegistry
}

func NewSettlementRegistryAdapter(
	client *ethclient.Client,
	addressHex string,
) (IParameterRegistry, error) {
	inner, err := settleReg.NewSettlementChainParameterRegistry(
		common.HexToAddress(addressHex),
		client,
	)
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

func (a *settlementRegistryAdapter) ParseParameterSet(log types.Log) (string, [32]byte, error) {
	ev, err := a.inner.ParseParameterSet(log)
	if err != nil {
		return "", [32]byte{}, err
	}
	return ev.Key, ev.Value, nil
}

type appChainRegistryAdapter struct {
	inner *appReg.AppChainParameterRegistry
}

func NewAppChainRegistryAdapter(
	client *ethclient.Client,
	addressHex string,
) (IParameterRegistry, error) {
	inner, err := appReg.NewAppChainParameterRegistry(common.HexToAddress(addressHex), client)
	if err != nil {
		return nil, err
	}
	return &appChainRegistryAdapter{inner: inner}, nil
}

func (a *appChainRegistryAdapter) Get(opts *bind.CallOpts, key string) ([32]byte, error) {
	return a.inner.Get(opts, key)
}

func (a *appChainRegistryAdapter) Set(
	opts *bind.TransactOpts,
	key string,
	value [32]byte,
) (*types.Transaction, error) {
	return a.inner.Set(opts, key, value)
}

func (a *appChainRegistryAdapter) SetMany(
	opts *bind.TransactOpts,
	keys []string,
	values [][32]byte,
) (*types.Transaction, error) {
	return a.inner.Set0(opts, keys, values)
}

func (a *appChainRegistryAdapter) ParseParameterSet(log types.Log) (string, [32]byte, error) {
	ev, err := a.inner.ParseParameterSet(log)
	if err != nil {
		return "", [32]byte{}, err
	}
	return ev.Key, ev.Value, nil
}

var (
	_ IParameterRegistry = (*settlementRegistryAdapter)(nil)
	_ IParameterRegistry = (*appChainRegistryAdapter)(nil)
)
