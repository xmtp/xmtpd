// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rateregistry

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// IRateRegistryRates is an auto generated low-level Go binding around an user-defined struct.
type IRateRegistryRates struct {
	MessageFee          uint64
	StorageFee          uint64
	CongestionFee       uint64
	TargetRatePerMinute uint64
	StartTime           uint64
}

// RateRegistryMetaData contains all meta data concerning the RateRegistry contract.
var RateRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"parameterRegistry_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"congestionFeeParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"contractName\",\"inputs\":[],\"outputs\":[{\"name\":\"contractName_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getRates\",\"inputs\":[{\"name\":\"fromIndex_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"count_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"rates_\",\"type\":\"tuple[]\",\"internalType\":\"structIRateRegistry.Rates[]\",\"components\":[{\"name\":\"messageFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"storageFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"congestionFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"targetRatePerMinute\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRatesCount\",\"inputs\":[],\"outputs\":[{\"name\":\"count_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"implementation_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"messageFeeParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"migrate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migratorParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"parameterRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ratesInEffectAfterParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"storageFeeParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"targetRatePerMinuteParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"updateRates\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"version_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Migrated\",\"inputs\":[{\"name\":\"migrator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RatesUpdated\",\"inputs\":[{\"name\":\"messageFee\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"storageFee\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"congestionFee\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"targetRatePerMinute\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"startTime\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EmptyCode\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"EndIndexOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FromIndexOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidStartTime\",\"inputs\":[{\"name\":\"startTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"lastStartTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MigrationFailed\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertData_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"NoChange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterOutOfTypeBounds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroCount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroMigrator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroParameterRegistry\",\"inputs\":[]}]",
}

// RateRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use RateRegistryMetaData.ABI instead.
var RateRegistryABI = RateRegistryMetaData.ABI

// RateRegistry is an auto generated Go binding around an Ethereum contract.
type RateRegistry struct {
	RateRegistryCaller     // Read-only binding to the contract
	RateRegistryTransactor // Write-only binding to the contract
	RateRegistryFilterer   // Log filterer for contract events
}

// RateRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type RateRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RateRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RateRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RateRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RateRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RateRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RateRegistrySession struct {
	Contract     *RateRegistry     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RateRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RateRegistryCallerSession struct {
	Contract *RateRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// RateRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RateRegistryTransactorSession struct {
	Contract     *RateRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// RateRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type RateRegistryRaw struct {
	Contract *RateRegistry // Generic contract binding to access the raw methods on
}

// RateRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RateRegistryCallerRaw struct {
	Contract *RateRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// RateRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RateRegistryTransactorRaw struct {
	Contract *RateRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRateRegistry creates a new instance of RateRegistry, bound to a specific deployed contract.
func NewRateRegistry(address common.Address, backend bind.ContractBackend) (*RateRegistry, error) {
	contract, err := bindRateRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RateRegistry{RateRegistryCaller: RateRegistryCaller{contract: contract}, RateRegistryTransactor: RateRegistryTransactor{contract: contract}, RateRegistryFilterer: RateRegistryFilterer{contract: contract}}, nil
}

// NewRateRegistryCaller creates a new read-only instance of RateRegistry, bound to a specific deployed contract.
func NewRateRegistryCaller(address common.Address, caller bind.ContractCaller) (*RateRegistryCaller, error) {
	contract, err := bindRateRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RateRegistryCaller{contract: contract}, nil
}

// NewRateRegistryTransactor creates a new write-only instance of RateRegistry, bound to a specific deployed contract.
func NewRateRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*RateRegistryTransactor, error) {
	contract, err := bindRateRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RateRegistryTransactor{contract: contract}, nil
}

// NewRateRegistryFilterer creates a new log filterer instance of RateRegistry, bound to a specific deployed contract.
func NewRateRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*RateRegistryFilterer, error) {
	contract, err := bindRateRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RateRegistryFilterer{contract: contract}, nil
}

// bindRateRegistry binds a generic wrapper to an already deployed contract.
func bindRateRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RateRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RateRegistry *RateRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RateRegistry.Contract.RateRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RateRegistry *RateRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RateRegistry.Contract.RateRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RateRegistry *RateRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RateRegistry.Contract.RateRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RateRegistry *RateRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RateRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RateRegistry *RateRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RateRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RateRegistry *RateRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RateRegistry.Contract.contract.Transact(opts, method, params...)
}

// CongestionFeeParameterKey is a free data retrieval call binding the contract method 0xed7e6986.
//
// Solidity: function congestionFeeParameterKey() pure returns(string key_)
func (_RateRegistry *RateRegistryCaller) CongestionFeeParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "congestionFeeParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// CongestionFeeParameterKey is a free data retrieval call binding the contract method 0xed7e6986.
//
// Solidity: function congestionFeeParameterKey() pure returns(string key_)
func (_RateRegistry *RateRegistrySession) CongestionFeeParameterKey() (string, error) {
	return _RateRegistry.Contract.CongestionFeeParameterKey(&_RateRegistry.CallOpts)
}

// CongestionFeeParameterKey is a free data retrieval call binding the contract method 0xed7e6986.
//
// Solidity: function congestionFeeParameterKey() pure returns(string key_)
func (_RateRegistry *RateRegistryCallerSession) CongestionFeeParameterKey() (string, error) {
	return _RateRegistry.Contract.CongestionFeeParameterKey(&_RateRegistry.CallOpts)
}

// ContractName is a free data retrieval call binding the contract method 0x75d0c0dc.
//
// Solidity: function contractName() pure returns(string contractName_)
func (_RateRegistry *RateRegistryCaller) ContractName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "contractName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ContractName is a free data retrieval call binding the contract method 0x75d0c0dc.
//
// Solidity: function contractName() pure returns(string contractName_)
func (_RateRegistry *RateRegistrySession) ContractName() (string, error) {
	return _RateRegistry.Contract.ContractName(&_RateRegistry.CallOpts)
}

// ContractName is a free data retrieval call binding the contract method 0x75d0c0dc.
//
// Solidity: function contractName() pure returns(string contractName_)
func (_RateRegistry *RateRegistryCallerSession) ContractName() (string, error) {
	return _RateRegistry.Contract.ContractName(&_RateRegistry.CallOpts)
}

// GetRates is a free data retrieval call binding the contract method 0x45b05a43.
//
// Solidity: function getRates(uint256 fromIndex_, uint256 count_) view returns((uint64,uint64,uint64,uint64,uint64)[] rates_)
func (_RateRegistry *RateRegistryCaller) GetRates(opts *bind.CallOpts, fromIndex_ *big.Int, count_ *big.Int) ([]IRateRegistryRates, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "getRates", fromIndex_, count_)

	if err != nil {
		return *new([]IRateRegistryRates), err
	}

	out0 := *abi.ConvertType(out[0], new([]IRateRegistryRates)).(*[]IRateRegistryRates)

	return out0, err

}

// GetRates is a free data retrieval call binding the contract method 0x45b05a43.
//
// Solidity: function getRates(uint256 fromIndex_, uint256 count_) view returns((uint64,uint64,uint64,uint64,uint64)[] rates_)
func (_RateRegistry *RateRegistrySession) GetRates(fromIndex_ *big.Int, count_ *big.Int) ([]IRateRegistryRates, error) {
	return _RateRegistry.Contract.GetRates(&_RateRegistry.CallOpts, fromIndex_, count_)
}

// GetRates is a free data retrieval call binding the contract method 0x45b05a43.
//
// Solidity: function getRates(uint256 fromIndex_, uint256 count_) view returns((uint64,uint64,uint64,uint64,uint64)[] rates_)
func (_RateRegistry *RateRegistryCallerSession) GetRates(fromIndex_ *big.Int, count_ *big.Int) ([]IRateRegistryRates, error) {
	return _RateRegistry.Contract.GetRates(&_RateRegistry.CallOpts, fromIndex_, count_)
}

// GetRatesCount is a free data retrieval call binding the contract method 0x2da72291.
//
// Solidity: function getRatesCount() view returns(uint256 count_)
func (_RateRegistry *RateRegistryCaller) GetRatesCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "getRatesCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRatesCount is a free data retrieval call binding the contract method 0x2da72291.
//
// Solidity: function getRatesCount() view returns(uint256 count_)
func (_RateRegistry *RateRegistrySession) GetRatesCount() (*big.Int, error) {
	return _RateRegistry.Contract.GetRatesCount(&_RateRegistry.CallOpts)
}

// GetRatesCount is a free data retrieval call binding the contract method 0x2da72291.
//
// Solidity: function getRatesCount() view returns(uint256 count_)
func (_RateRegistry *RateRegistryCallerSession) GetRatesCount() (*big.Int, error) {
	return _RateRegistry.Contract.GetRatesCount(&_RateRegistry.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_RateRegistry *RateRegistryCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_RateRegistry *RateRegistrySession) Implementation() (common.Address, error) {
	return _RateRegistry.Contract.Implementation(&_RateRegistry.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_RateRegistry *RateRegistryCallerSession) Implementation() (common.Address, error) {
	return _RateRegistry.Contract.Implementation(&_RateRegistry.CallOpts)
}

// MessageFeeParameterKey is a free data retrieval call binding the contract method 0x63c03291.
//
// Solidity: function messageFeeParameterKey() pure returns(string key_)
func (_RateRegistry *RateRegistryCaller) MessageFeeParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "messageFeeParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// MessageFeeParameterKey is a free data retrieval call binding the contract method 0x63c03291.
//
// Solidity: function messageFeeParameterKey() pure returns(string key_)
func (_RateRegistry *RateRegistrySession) MessageFeeParameterKey() (string, error) {
	return _RateRegistry.Contract.MessageFeeParameterKey(&_RateRegistry.CallOpts)
}

// MessageFeeParameterKey is a free data retrieval call binding the contract method 0x63c03291.
//
// Solidity: function messageFeeParameterKey() pure returns(string key_)
func (_RateRegistry *RateRegistryCallerSession) MessageFeeParameterKey() (string, error) {
	return _RateRegistry.Contract.MessageFeeParameterKey(&_RateRegistry.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_RateRegistry *RateRegistryCaller) MigratorParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "migratorParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_RateRegistry *RateRegistrySession) MigratorParameterKey() (string, error) {
	return _RateRegistry.Contract.MigratorParameterKey(&_RateRegistry.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_RateRegistry *RateRegistryCallerSession) MigratorParameterKey() (string, error) {
	return _RateRegistry.Contract.MigratorParameterKey(&_RateRegistry.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_RateRegistry *RateRegistryCaller) ParameterRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "parameterRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_RateRegistry *RateRegistrySession) ParameterRegistry() (common.Address, error) {
	return _RateRegistry.Contract.ParameterRegistry(&_RateRegistry.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_RateRegistry *RateRegistryCallerSession) ParameterRegistry() (common.Address, error) {
	return _RateRegistry.Contract.ParameterRegistry(&_RateRegistry.CallOpts)
}

// RatesInEffectAfterParameterKey is a free data retrieval call binding the contract method 0xf392c428.
//
// Solidity: function ratesInEffectAfterParameterKey() pure returns(string key_)
func (_RateRegistry *RateRegistryCaller) RatesInEffectAfterParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "ratesInEffectAfterParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// RatesInEffectAfterParameterKey is a free data retrieval call binding the contract method 0xf392c428.
//
// Solidity: function ratesInEffectAfterParameterKey() pure returns(string key_)
func (_RateRegistry *RateRegistrySession) RatesInEffectAfterParameterKey() (string, error) {
	return _RateRegistry.Contract.RatesInEffectAfterParameterKey(&_RateRegistry.CallOpts)
}

// RatesInEffectAfterParameterKey is a free data retrieval call binding the contract method 0xf392c428.
//
// Solidity: function ratesInEffectAfterParameterKey() pure returns(string key_)
func (_RateRegistry *RateRegistryCallerSession) RatesInEffectAfterParameterKey() (string, error) {
	return _RateRegistry.Contract.RatesInEffectAfterParameterKey(&_RateRegistry.CallOpts)
}

// StorageFeeParameterKey is a free data retrieval call binding the contract method 0xba3261d5.
//
// Solidity: function storageFeeParameterKey() pure returns(string key_)
func (_RateRegistry *RateRegistryCaller) StorageFeeParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "storageFeeParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// StorageFeeParameterKey is a free data retrieval call binding the contract method 0xba3261d5.
//
// Solidity: function storageFeeParameterKey() pure returns(string key_)
func (_RateRegistry *RateRegistrySession) StorageFeeParameterKey() (string, error) {
	return _RateRegistry.Contract.StorageFeeParameterKey(&_RateRegistry.CallOpts)
}

// StorageFeeParameterKey is a free data retrieval call binding the contract method 0xba3261d5.
//
// Solidity: function storageFeeParameterKey() pure returns(string key_)
func (_RateRegistry *RateRegistryCallerSession) StorageFeeParameterKey() (string, error) {
	return _RateRegistry.Contract.StorageFeeParameterKey(&_RateRegistry.CallOpts)
}

// TargetRatePerMinuteParameterKey is a free data retrieval call binding the contract method 0x49156e09.
//
// Solidity: function targetRatePerMinuteParameterKey() pure returns(string key_)
func (_RateRegistry *RateRegistryCaller) TargetRatePerMinuteParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "targetRatePerMinuteParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TargetRatePerMinuteParameterKey is a free data retrieval call binding the contract method 0x49156e09.
//
// Solidity: function targetRatePerMinuteParameterKey() pure returns(string key_)
func (_RateRegistry *RateRegistrySession) TargetRatePerMinuteParameterKey() (string, error) {
	return _RateRegistry.Contract.TargetRatePerMinuteParameterKey(&_RateRegistry.CallOpts)
}

// TargetRatePerMinuteParameterKey is a free data retrieval call binding the contract method 0x49156e09.
//
// Solidity: function targetRatePerMinuteParameterKey() pure returns(string key_)
func (_RateRegistry *RateRegistryCallerSession) TargetRatePerMinuteParameterKey() (string, error) {
	return _RateRegistry.Contract.TargetRatePerMinuteParameterKey(&_RateRegistry.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(string version_)
func (_RateRegistry *RateRegistryCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(string version_)
func (_RateRegistry *RateRegistrySession) Version() (string, error) {
	return _RateRegistry.Contract.Version(&_RateRegistry.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(string version_)
func (_RateRegistry *RateRegistryCallerSession) Version() (string, error) {
	return _RateRegistry.Contract.Version(&_RateRegistry.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_RateRegistry *RateRegistryTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RateRegistry.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_RateRegistry *RateRegistrySession) Initialize() (*types.Transaction, error) {
	return _RateRegistry.Contract.Initialize(&_RateRegistry.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_RateRegistry *RateRegistryTransactorSession) Initialize() (*types.Transaction, error) {
	return _RateRegistry.Contract.Initialize(&_RateRegistry.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_RateRegistry *RateRegistryTransactor) Migrate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RateRegistry.contract.Transact(opts, "migrate")
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_RateRegistry *RateRegistrySession) Migrate() (*types.Transaction, error) {
	return _RateRegistry.Contract.Migrate(&_RateRegistry.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_RateRegistry *RateRegistryTransactorSession) Migrate() (*types.Transaction, error) {
	return _RateRegistry.Contract.Migrate(&_RateRegistry.TransactOpts)
}

// UpdateRates is a paid mutator transaction binding the contract method 0x3c3821f4.
//
// Solidity: function updateRates() returns()
func (_RateRegistry *RateRegistryTransactor) UpdateRates(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RateRegistry.contract.Transact(opts, "updateRates")
}

// UpdateRates is a paid mutator transaction binding the contract method 0x3c3821f4.
//
// Solidity: function updateRates() returns()
func (_RateRegistry *RateRegistrySession) UpdateRates() (*types.Transaction, error) {
	return _RateRegistry.Contract.UpdateRates(&_RateRegistry.TransactOpts)
}

// UpdateRates is a paid mutator transaction binding the contract method 0x3c3821f4.
//
// Solidity: function updateRates() returns()
func (_RateRegistry *RateRegistryTransactorSession) UpdateRates() (*types.Transaction, error) {
	return _RateRegistry.Contract.UpdateRates(&_RateRegistry.TransactOpts)
}

// RateRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the RateRegistry contract.
type RateRegistryInitializedIterator struct {
	Event *RateRegistryInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RateRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RateRegistryInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RateRegistryInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RateRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RateRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RateRegistryInitialized represents a Initialized event raised by the RateRegistry contract.
type RateRegistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_RateRegistry *RateRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*RateRegistryInitializedIterator, error) {

	logs, sub, err := _RateRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &RateRegistryInitializedIterator{contract: _RateRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_RateRegistry *RateRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *RateRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _RateRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RateRegistryInitialized)
				if err := _RateRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_RateRegistry *RateRegistryFilterer) ParseInitialized(log types.Log) (*RateRegistryInitialized, error) {
	event := new(RateRegistryInitialized)
	if err := _RateRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RateRegistryMigratedIterator is returned from FilterMigrated and is used to iterate over the raw logs and unpacked data for Migrated events raised by the RateRegistry contract.
type RateRegistryMigratedIterator struct {
	Event *RateRegistryMigrated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RateRegistryMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RateRegistryMigrated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RateRegistryMigrated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RateRegistryMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RateRegistryMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RateRegistryMigrated represents a Migrated event raised by the RateRegistry contract.
type RateRegistryMigrated struct {
	Migrator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMigrated is a free log retrieval operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_RateRegistry *RateRegistryFilterer) FilterMigrated(opts *bind.FilterOpts, migrator []common.Address) (*RateRegistryMigratedIterator, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _RateRegistry.contract.FilterLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return &RateRegistryMigratedIterator{contract: _RateRegistry.contract, event: "Migrated", logs: logs, sub: sub}, nil
}

// WatchMigrated is a free log subscription operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_RateRegistry *RateRegistryFilterer) WatchMigrated(opts *bind.WatchOpts, sink chan<- *RateRegistryMigrated, migrator []common.Address) (event.Subscription, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _RateRegistry.contract.WatchLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RateRegistryMigrated)
				if err := _RateRegistry.contract.UnpackLog(event, "Migrated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMigrated is a log parse operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_RateRegistry *RateRegistryFilterer) ParseMigrated(log types.Log) (*RateRegistryMigrated, error) {
	event := new(RateRegistryMigrated)
	if err := _RateRegistry.contract.UnpackLog(event, "Migrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RateRegistryRatesUpdatedIterator is returned from FilterRatesUpdated and is used to iterate over the raw logs and unpacked data for RatesUpdated events raised by the RateRegistry contract.
type RateRegistryRatesUpdatedIterator struct {
	Event *RateRegistryRatesUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RateRegistryRatesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RateRegistryRatesUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RateRegistryRatesUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RateRegistryRatesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RateRegistryRatesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RateRegistryRatesUpdated represents a RatesUpdated event raised by the RateRegistry contract.
type RateRegistryRatesUpdated struct {
	MessageFee          uint64
	StorageFee          uint64
	CongestionFee       uint64
	TargetRatePerMinute uint64
	StartTime           uint64
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterRatesUpdated is a free log retrieval operation binding the contract event 0x8aa9960c80aa047e81f0b89c422f689e9b6adc187f3f2b2ccb957baf2e6f761b.
//
// Solidity: event RatesUpdated(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 targetRatePerMinute, uint64 startTime)
func (_RateRegistry *RateRegistryFilterer) FilterRatesUpdated(opts *bind.FilterOpts) (*RateRegistryRatesUpdatedIterator, error) {

	logs, sub, err := _RateRegistry.contract.FilterLogs(opts, "RatesUpdated")
	if err != nil {
		return nil, err
	}
	return &RateRegistryRatesUpdatedIterator{contract: _RateRegistry.contract, event: "RatesUpdated", logs: logs, sub: sub}, nil
}

// WatchRatesUpdated is a free log subscription operation binding the contract event 0x8aa9960c80aa047e81f0b89c422f689e9b6adc187f3f2b2ccb957baf2e6f761b.
//
// Solidity: event RatesUpdated(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 targetRatePerMinute, uint64 startTime)
func (_RateRegistry *RateRegistryFilterer) WatchRatesUpdated(opts *bind.WatchOpts, sink chan<- *RateRegistryRatesUpdated) (event.Subscription, error) {

	logs, sub, err := _RateRegistry.contract.WatchLogs(opts, "RatesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RateRegistryRatesUpdated)
				if err := _RateRegistry.contract.UnpackLog(event, "RatesUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRatesUpdated is a log parse operation binding the contract event 0x8aa9960c80aa047e81f0b89c422f689e9b6adc187f3f2b2ccb957baf2e6f761b.
//
// Solidity: event RatesUpdated(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 targetRatePerMinute, uint64 startTime)
func (_RateRegistry *RateRegistryFilterer) ParseRatesUpdated(log types.Log) (*RateRegistryRatesUpdated, error) {
	event := new(RateRegistryRatesUpdated)
	if err := _RateRegistry.contract.UnpackLog(event, "RatesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RateRegistryUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the RateRegistry contract.
type RateRegistryUpgradedIterator struct {
	Event *RateRegistryUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RateRegistryUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RateRegistryUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RateRegistryUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RateRegistryUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RateRegistryUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RateRegistryUpgraded represents a Upgraded event raised by the RateRegistry contract.
type RateRegistryUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_RateRegistry *RateRegistryFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*RateRegistryUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _RateRegistry.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &RateRegistryUpgradedIterator{contract: _RateRegistry.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_RateRegistry *RateRegistryFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *RateRegistryUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _RateRegistry.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RateRegistryUpgraded)
				if err := _RateRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_RateRegistry *RateRegistryFilterer) ParseUpgraded(log types.Log) (*RateRegistryUpgraded, error) {
	event := new(RateRegistryUpgraded)
	if err := _RateRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
