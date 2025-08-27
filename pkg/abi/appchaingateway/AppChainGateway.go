// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package appchaingateway

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

// AppChainGatewayMetaData contains all meta data concerning the AppChainGateway contract.
var AppChainGatewayMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"parameterRegistry_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"settlementChainGateway_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"implementation_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migrate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migratorParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"parameterRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"paused_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pausedParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"receiveDeposit\",\"inputs\":[{\"name\":\"recipient_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"receiveParameters\",\"inputs\":[{\"name\":\"nonce_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"keys_\",\"type\":\"string[]\",\"internalType\":\"string[]\"},{\"name\":\"values_\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"settlementChainGateway\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"settlementChainGatewayAlias\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updatePauseStatus\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"recipient_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"withdrawIntoUnderlying\",\"inputs\":[{\"name\":\"recipient_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"DepositReceived\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Migrated\",\"inputs\":[{\"name\":\"migrator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ParametersReceived\",\"inputs\":[{\"name\":\"nonce\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"keys\",\"type\":\"string[]\",\"indexed\":false,\"internalType\":\"string[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PauseStatusUpdated\",\"inputs\":[{\"name\":\"paused\",\"type\":\"bool\",\"indexed\":true,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawal\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"messageId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EmptyCode\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationFailed\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertData_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"NoChange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotSettlementChainGateway\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterOutOfTypeBounds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Paused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransferFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroMigrator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroParameterRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroRecipient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroSettlementChainGateway\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroWithdrawalAmount\",\"inputs\":[]}]",
}

// AppChainGatewayABI is the input ABI used to generate the binding from.
// Deprecated: Use AppChainGatewayMetaData.ABI instead.
var AppChainGatewayABI = AppChainGatewayMetaData.ABI

// AppChainGateway is an auto generated Go binding around an Ethereum contract.
type AppChainGateway struct {
	AppChainGatewayCaller     // Read-only binding to the contract
	AppChainGatewayTransactor // Write-only binding to the contract
	AppChainGatewayFilterer   // Log filterer for contract events
}

// AppChainGatewayCaller is an auto generated read-only Go binding around an Ethereum contract.
type AppChainGatewayCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AppChainGatewayTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AppChainGatewayTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AppChainGatewayFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AppChainGatewayFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AppChainGatewaySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AppChainGatewaySession struct {
	Contract     *AppChainGateway  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AppChainGatewayCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AppChainGatewayCallerSession struct {
	Contract *AppChainGatewayCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// AppChainGatewayTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AppChainGatewayTransactorSession struct {
	Contract     *AppChainGatewayTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// AppChainGatewayRaw is an auto generated low-level Go binding around an Ethereum contract.
type AppChainGatewayRaw struct {
	Contract *AppChainGateway // Generic contract binding to access the raw methods on
}

// AppChainGatewayCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AppChainGatewayCallerRaw struct {
	Contract *AppChainGatewayCaller // Generic read-only contract binding to access the raw methods on
}

// AppChainGatewayTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AppChainGatewayTransactorRaw struct {
	Contract *AppChainGatewayTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAppChainGateway creates a new instance of AppChainGateway, bound to a specific deployed contract.
func NewAppChainGateway(address common.Address, backend bind.ContractBackend) (*AppChainGateway, error) {
	contract, err := bindAppChainGateway(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AppChainGateway{AppChainGatewayCaller: AppChainGatewayCaller{contract: contract}, AppChainGatewayTransactor: AppChainGatewayTransactor{contract: contract}, AppChainGatewayFilterer: AppChainGatewayFilterer{contract: contract}}, nil
}

// NewAppChainGatewayCaller creates a new read-only instance of AppChainGateway, bound to a specific deployed contract.
func NewAppChainGatewayCaller(address common.Address, caller bind.ContractCaller) (*AppChainGatewayCaller, error) {
	contract, err := bindAppChainGateway(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AppChainGatewayCaller{contract: contract}, nil
}

// NewAppChainGatewayTransactor creates a new write-only instance of AppChainGateway, bound to a specific deployed contract.
func NewAppChainGatewayTransactor(address common.Address, transactor bind.ContractTransactor) (*AppChainGatewayTransactor, error) {
	contract, err := bindAppChainGateway(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AppChainGatewayTransactor{contract: contract}, nil
}

// NewAppChainGatewayFilterer creates a new log filterer instance of AppChainGateway, bound to a specific deployed contract.
func NewAppChainGatewayFilterer(address common.Address, filterer bind.ContractFilterer) (*AppChainGatewayFilterer, error) {
	contract, err := bindAppChainGateway(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AppChainGatewayFilterer{contract: contract}, nil
}

// bindAppChainGateway binds a generic wrapper to an already deployed contract.
func bindAppChainGateway(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AppChainGatewayMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AppChainGateway *AppChainGatewayRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AppChainGateway.Contract.AppChainGatewayCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AppChainGateway *AppChainGatewayRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AppChainGateway.Contract.AppChainGatewayTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AppChainGateway *AppChainGatewayRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AppChainGateway.Contract.AppChainGatewayTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AppChainGateway *AppChainGatewayCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AppChainGateway.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AppChainGateway *AppChainGatewayTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AppChainGateway.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AppChainGateway *AppChainGatewayTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AppChainGateway.Contract.contract.Transact(opts, method, params...)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_AppChainGateway *AppChainGatewayCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AppChainGateway.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_AppChainGateway *AppChainGatewaySession) Implementation() (common.Address, error) {
	return _AppChainGateway.Contract.Implementation(&_AppChainGateway.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_AppChainGateway *AppChainGatewayCallerSession) Implementation() (common.Address, error) {
	return _AppChainGateway.Contract.Implementation(&_AppChainGateway.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_AppChainGateway *AppChainGatewayCaller) MigratorParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AppChainGateway.contract.Call(opts, &out, "migratorParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_AppChainGateway *AppChainGatewaySession) MigratorParameterKey() (string, error) {
	return _AppChainGateway.Contract.MigratorParameterKey(&_AppChainGateway.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_AppChainGateway *AppChainGatewayCallerSession) MigratorParameterKey() (string, error) {
	return _AppChainGateway.Contract.MigratorParameterKey(&_AppChainGateway.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_AppChainGateway *AppChainGatewayCaller) ParameterRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AppChainGateway.contract.Call(opts, &out, "parameterRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_AppChainGateway *AppChainGatewaySession) ParameterRegistry() (common.Address, error) {
	return _AppChainGateway.Contract.ParameterRegistry(&_AppChainGateway.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_AppChainGateway *AppChainGatewayCallerSession) ParameterRegistry() (common.Address, error) {
	return _AppChainGateway.Contract.ParameterRegistry(&_AppChainGateway.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool paused_)
func (_AppChainGateway *AppChainGatewayCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _AppChainGateway.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool paused_)
func (_AppChainGateway *AppChainGatewaySession) Paused() (bool, error) {
	return _AppChainGateway.Contract.Paused(&_AppChainGateway.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool paused_)
func (_AppChainGateway *AppChainGatewayCallerSession) Paused() (bool, error) {
	return _AppChainGateway.Contract.Paused(&_AppChainGateway.CallOpts)
}

// PausedParameterKey is a free data retrieval call binding the contract method 0xcc5999af.
//
// Solidity: function pausedParameterKey() pure returns(string key_)
func (_AppChainGateway *AppChainGatewayCaller) PausedParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AppChainGateway.contract.Call(opts, &out, "pausedParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// PausedParameterKey is a free data retrieval call binding the contract method 0xcc5999af.
//
// Solidity: function pausedParameterKey() pure returns(string key_)
func (_AppChainGateway *AppChainGatewaySession) PausedParameterKey() (string, error) {
	return _AppChainGateway.Contract.PausedParameterKey(&_AppChainGateway.CallOpts)
}

// PausedParameterKey is a free data retrieval call binding the contract method 0xcc5999af.
//
// Solidity: function pausedParameterKey() pure returns(string key_)
func (_AppChainGateway *AppChainGatewayCallerSession) PausedParameterKey() (string, error) {
	return _AppChainGateway.Contract.PausedParameterKey(&_AppChainGateway.CallOpts)
}

// SettlementChainGateway is a free data retrieval call binding the contract method 0x801fd7f3.
//
// Solidity: function settlementChainGateway() view returns(address)
func (_AppChainGateway *AppChainGatewayCaller) SettlementChainGateway(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AppChainGateway.contract.Call(opts, &out, "settlementChainGateway")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SettlementChainGateway is a free data retrieval call binding the contract method 0x801fd7f3.
//
// Solidity: function settlementChainGateway() view returns(address)
func (_AppChainGateway *AppChainGatewaySession) SettlementChainGateway() (common.Address, error) {
	return _AppChainGateway.Contract.SettlementChainGateway(&_AppChainGateway.CallOpts)
}

// SettlementChainGateway is a free data retrieval call binding the contract method 0x801fd7f3.
//
// Solidity: function settlementChainGateway() view returns(address)
func (_AppChainGateway *AppChainGatewayCallerSession) SettlementChainGateway() (common.Address, error) {
	return _AppChainGateway.Contract.SettlementChainGateway(&_AppChainGateway.CallOpts)
}

// SettlementChainGatewayAlias is a free data retrieval call binding the contract method 0x06646743.
//
// Solidity: function settlementChainGatewayAlias() view returns(address)
func (_AppChainGateway *AppChainGatewayCaller) SettlementChainGatewayAlias(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AppChainGateway.contract.Call(opts, &out, "settlementChainGatewayAlias")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SettlementChainGatewayAlias is a free data retrieval call binding the contract method 0x06646743.
//
// Solidity: function settlementChainGatewayAlias() view returns(address)
func (_AppChainGateway *AppChainGatewaySession) SettlementChainGatewayAlias() (common.Address, error) {
	return _AppChainGateway.Contract.SettlementChainGatewayAlias(&_AppChainGateway.CallOpts)
}

// SettlementChainGatewayAlias is a free data retrieval call binding the contract method 0x06646743.
//
// Solidity: function settlementChainGatewayAlias() view returns(address)
func (_AppChainGateway *AppChainGatewayCallerSession) SettlementChainGatewayAlias() (common.Address, error) {
	return _AppChainGateway.Contract.SettlementChainGatewayAlias(&_AppChainGateway.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_AppChainGateway *AppChainGatewayTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AppChainGateway.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_AppChainGateway *AppChainGatewaySession) Initialize() (*types.Transaction, error) {
	return _AppChainGateway.Contract.Initialize(&_AppChainGateway.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_AppChainGateway *AppChainGatewayTransactorSession) Initialize() (*types.Transaction, error) {
	return _AppChainGateway.Contract.Initialize(&_AppChainGateway.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_AppChainGateway *AppChainGatewayTransactor) Migrate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AppChainGateway.contract.Transact(opts, "migrate")
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_AppChainGateway *AppChainGatewaySession) Migrate() (*types.Transaction, error) {
	return _AppChainGateway.Contract.Migrate(&_AppChainGateway.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_AppChainGateway *AppChainGatewayTransactorSession) Migrate() (*types.Transaction, error) {
	return _AppChainGateway.Contract.Migrate(&_AppChainGateway.TransactOpts)
}

// ReceiveDeposit is a paid mutator transaction binding the contract method 0x4578474d.
//
// Solidity: function receiveDeposit(address recipient_) payable returns()
func (_AppChainGateway *AppChainGatewayTransactor) ReceiveDeposit(opts *bind.TransactOpts, recipient_ common.Address) (*types.Transaction, error) {
	return _AppChainGateway.contract.Transact(opts, "receiveDeposit", recipient_)
}

// ReceiveDeposit is a paid mutator transaction binding the contract method 0x4578474d.
//
// Solidity: function receiveDeposit(address recipient_) payable returns()
func (_AppChainGateway *AppChainGatewaySession) ReceiveDeposit(recipient_ common.Address) (*types.Transaction, error) {
	return _AppChainGateway.Contract.ReceiveDeposit(&_AppChainGateway.TransactOpts, recipient_)
}

// ReceiveDeposit is a paid mutator transaction binding the contract method 0x4578474d.
//
// Solidity: function receiveDeposit(address recipient_) payable returns()
func (_AppChainGateway *AppChainGatewayTransactorSession) ReceiveDeposit(recipient_ common.Address) (*types.Transaction, error) {
	return _AppChainGateway.Contract.ReceiveDeposit(&_AppChainGateway.TransactOpts, recipient_)
}

// ReceiveParameters is a paid mutator transaction binding the contract method 0x333b1e2f.
//
// Solidity: function receiveParameters(uint256 nonce_, string[] keys_, bytes32[] values_) returns()
func (_AppChainGateway *AppChainGatewayTransactor) ReceiveParameters(opts *bind.TransactOpts, nonce_ *big.Int, keys_ []string, values_ [][32]byte) (*types.Transaction, error) {
	return _AppChainGateway.contract.Transact(opts, "receiveParameters", nonce_, keys_, values_)
}

// ReceiveParameters is a paid mutator transaction binding the contract method 0x333b1e2f.
//
// Solidity: function receiveParameters(uint256 nonce_, string[] keys_, bytes32[] values_) returns()
func (_AppChainGateway *AppChainGatewaySession) ReceiveParameters(nonce_ *big.Int, keys_ []string, values_ [][32]byte) (*types.Transaction, error) {
	return _AppChainGateway.Contract.ReceiveParameters(&_AppChainGateway.TransactOpts, nonce_, keys_, values_)
}

// ReceiveParameters is a paid mutator transaction binding the contract method 0x333b1e2f.
//
// Solidity: function receiveParameters(uint256 nonce_, string[] keys_, bytes32[] values_) returns()
func (_AppChainGateway *AppChainGatewayTransactorSession) ReceiveParameters(nonce_ *big.Int, keys_ []string, values_ [][32]byte) (*types.Transaction, error) {
	return _AppChainGateway.Contract.ReceiveParameters(&_AppChainGateway.TransactOpts, nonce_, keys_, values_)
}

// UpdatePauseStatus is a paid mutator transaction binding the contract method 0x59d4df41.
//
// Solidity: function updatePauseStatus() returns()
func (_AppChainGateway *AppChainGatewayTransactor) UpdatePauseStatus(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AppChainGateway.contract.Transact(opts, "updatePauseStatus")
}

// UpdatePauseStatus is a paid mutator transaction binding the contract method 0x59d4df41.
//
// Solidity: function updatePauseStatus() returns()
func (_AppChainGateway *AppChainGatewaySession) UpdatePauseStatus() (*types.Transaction, error) {
	return _AppChainGateway.Contract.UpdatePauseStatus(&_AppChainGateway.TransactOpts)
}

// UpdatePauseStatus is a paid mutator transaction binding the contract method 0x59d4df41.
//
// Solidity: function updatePauseStatus() returns()
func (_AppChainGateway *AppChainGatewayTransactorSession) UpdatePauseStatus() (*types.Transaction, error) {
	return _AppChainGateway.Contract.UpdatePauseStatus(&_AppChainGateway.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x51cff8d9.
//
// Solidity: function withdraw(address recipient_) payable returns()
func (_AppChainGateway *AppChainGatewayTransactor) Withdraw(opts *bind.TransactOpts, recipient_ common.Address) (*types.Transaction, error) {
	return _AppChainGateway.contract.Transact(opts, "withdraw", recipient_)
}

// Withdraw is a paid mutator transaction binding the contract method 0x51cff8d9.
//
// Solidity: function withdraw(address recipient_) payable returns()
func (_AppChainGateway *AppChainGatewaySession) Withdraw(recipient_ common.Address) (*types.Transaction, error) {
	return _AppChainGateway.Contract.Withdraw(&_AppChainGateway.TransactOpts, recipient_)
}

// Withdraw is a paid mutator transaction binding the contract method 0x51cff8d9.
//
// Solidity: function withdraw(address recipient_) payable returns()
func (_AppChainGateway *AppChainGatewayTransactorSession) Withdraw(recipient_ common.Address) (*types.Transaction, error) {
	return _AppChainGateway.Contract.Withdraw(&_AppChainGateway.TransactOpts, recipient_)
}

// WithdrawIntoUnderlying is a paid mutator transaction binding the contract method 0xae8c3d65.
//
// Solidity: function withdrawIntoUnderlying(address recipient_) payable returns()
func (_AppChainGateway *AppChainGatewayTransactor) WithdrawIntoUnderlying(opts *bind.TransactOpts, recipient_ common.Address) (*types.Transaction, error) {
	return _AppChainGateway.contract.Transact(opts, "withdrawIntoUnderlying", recipient_)
}

// WithdrawIntoUnderlying is a paid mutator transaction binding the contract method 0xae8c3d65.
//
// Solidity: function withdrawIntoUnderlying(address recipient_) payable returns()
func (_AppChainGateway *AppChainGatewaySession) WithdrawIntoUnderlying(recipient_ common.Address) (*types.Transaction, error) {
	return _AppChainGateway.Contract.WithdrawIntoUnderlying(&_AppChainGateway.TransactOpts, recipient_)
}

// WithdrawIntoUnderlying is a paid mutator transaction binding the contract method 0xae8c3d65.
//
// Solidity: function withdrawIntoUnderlying(address recipient_) payable returns()
func (_AppChainGateway *AppChainGatewayTransactorSession) WithdrawIntoUnderlying(recipient_ common.Address) (*types.Transaction, error) {
	return _AppChainGateway.Contract.WithdrawIntoUnderlying(&_AppChainGateway.TransactOpts, recipient_)
}

// AppChainGatewayDepositReceivedIterator is returned from FilterDepositReceived and is used to iterate over the raw logs and unpacked data for DepositReceived events raised by the AppChainGateway contract.
type AppChainGatewayDepositReceivedIterator struct {
	Event *AppChainGatewayDepositReceived // Event containing the contract specifics and raw log

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
func (it *AppChainGatewayDepositReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AppChainGatewayDepositReceived)
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
		it.Event = new(AppChainGatewayDepositReceived)
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
func (it *AppChainGatewayDepositReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AppChainGatewayDepositReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AppChainGatewayDepositReceived represents a DepositReceived event raised by the AppChainGateway contract.
type AppChainGatewayDepositReceived struct {
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDepositReceived is a free log retrieval operation binding the contract event 0x9936746a4565f9766fa768f88f56a7487c78780ac179562773d1c75c5269537e.
//
// Solidity: event DepositReceived(address indexed recipient, uint256 amount)
func (_AppChainGateway *AppChainGatewayFilterer) FilterDepositReceived(opts *bind.FilterOpts, recipient []common.Address) (*AppChainGatewayDepositReceivedIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AppChainGateway.contract.FilterLogs(opts, "DepositReceived", recipientRule)
	if err != nil {
		return nil, err
	}
	return &AppChainGatewayDepositReceivedIterator{contract: _AppChainGateway.contract, event: "DepositReceived", logs: logs, sub: sub}, nil
}

// WatchDepositReceived is a free log subscription operation binding the contract event 0x9936746a4565f9766fa768f88f56a7487c78780ac179562773d1c75c5269537e.
//
// Solidity: event DepositReceived(address indexed recipient, uint256 amount)
func (_AppChainGateway *AppChainGatewayFilterer) WatchDepositReceived(opts *bind.WatchOpts, sink chan<- *AppChainGatewayDepositReceived, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AppChainGateway.contract.WatchLogs(opts, "DepositReceived", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AppChainGatewayDepositReceived)
				if err := _AppChainGateway.contract.UnpackLog(event, "DepositReceived", log); err != nil {
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

// ParseDepositReceived is a log parse operation binding the contract event 0x9936746a4565f9766fa768f88f56a7487c78780ac179562773d1c75c5269537e.
//
// Solidity: event DepositReceived(address indexed recipient, uint256 amount)
func (_AppChainGateway *AppChainGatewayFilterer) ParseDepositReceived(log types.Log) (*AppChainGatewayDepositReceived, error) {
	event := new(AppChainGatewayDepositReceived)
	if err := _AppChainGateway.contract.UnpackLog(event, "DepositReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AppChainGatewayInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the AppChainGateway contract.
type AppChainGatewayInitializedIterator struct {
	Event *AppChainGatewayInitialized // Event containing the contract specifics and raw log

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
func (it *AppChainGatewayInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AppChainGatewayInitialized)
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
		it.Event = new(AppChainGatewayInitialized)
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
func (it *AppChainGatewayInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AppChainGatewayInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AppChainGatewayInitialized represents a Initialized event raised by the AppChainGateway contract.
type AppChainGatewayInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AppChainGateway *AppChainGatewayFilterer) FilterInitialized(opts *bind.FilterOpts) (*AppChainGatewayInitializedIterator, error) {

	logs, sub, err := _AppChainGateway.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &AppChainGatewayInitializedIterator{contract: _AppChainGateway.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AppChainGateway *AppChainGatewayFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *AppChainGatewayInitialized) (event.Subscription, error) {

	logs, sub, err := _AppChainGateway.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AppChainGatewayInitialized)
				if err := _AppChainGateway.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_AppChainGateway *AppChainGatewayFilterer) ParseInitialized(log types.Log) (*AppChainGatewayInitialized, error) {
	event := new(AppChainGatewayInitialized)
	if err := _AppChainGateway.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AppChainGatewayMigratedIterator is returned from FilterMigrated and is used to iterate over the raw logs and unpacked data for Migrated events raised by the AppChainGateway contract.
type AppChainGatewayMigratedIterator struct {
	Event *AppChainGatewayMigrated // Event containing the contract specifics and raw log

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
func (it *AppChainGatewayMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AppChainGatewayMigrated)
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
		it.Event = new(AppChainGatewayMigrated)
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
func (it *AppChainGatewayMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AppChainGatewayMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AppChainGatewayMigrated represents a Migrated event raised by the AppChainGateway contract.
type AppChainGatewayMigrated struct {
	Migrator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMigrated is a free log retrieval operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_AppChainGateway *AppChainGatewayFilterer) FilterMigrated(opts *bind.FilterOpts, migrator []common.Address) (*AppChainGatewayMigratedIterator, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _AppChainGateway.contract.FilterLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return &AppChainGatewayMigratedIterator{contract: _AppChainGateway.contract, event: "Migrated", logs: logs, sub: sub}, nil
}

// WatchMigrated is a free log subscription operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_AppChainGateway *AppChainGatewayFilterer) WatchMigrated(opts *bind.WatchOpts, sink chan<- *AppChainGatewayMigrated, migrator []common.Address) (event.Subscription, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _AppChainGateway.contract.WatchLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AppChainGatewayMigrated)
				if err := _AppChainGateway.contract.UnpackLog(event, "Migrated", log); err != nil {
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
func (_AppChainGateway *AppChainGatewayFilterer) ParseMigrated(log types.Log) (*AppChainGatewayMigrated, error) {
	event := new(AppChainGatewayMigrated)
	if err := _AppChainGateway.contract.UnpackLog(event, "Migrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AppChainGatewayParametersReceivedIterator is returned from FilterParametersReceived and is used to iterate over the raw logs and unpacked data for ParametersReceived events raised by the AppChainGateway contract.
type AppChainGatewayParametersReceivedIterator struct {
	Event *AppChainGatewayParametersReceived // Event containing the contract specifics and raw log

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
func (it *AppChainGatewayParametersReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AppChainGatewayParametersReceived)
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
		it.Event = new(AppChainGatewayParametersReceived)
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
func (it *AppChainGatewayParametersReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AppChainGatewayParametersReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AppChainGatewayParametersReceived represents a ParametersReceived event raised by the AppChainGateway contract.
type AppChainGatewayParametersReceived struct {
	Nonce *big.Int
	Keys  []string
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterParametersReceived is a free log retrieval operation binding the contract event 0x7fb103ff88be3fb0c788a3ee63031c896dfd0695b79a8061c138e01fd1224b76.
//
// Solidity: event ParametersReceived(uint256 indexed nonce, string[] keys)
func (_AppChainGateway *AppChainGatewayFilterer) FilterParametersReceived(opts *bind.FilterOpts, nonce []*big.Int) (*AppChainGatewayParametersReceivedIterator, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}

	logs, sub, err := _AppChainGateway.contract.FilterLogs(opts, "ParametersReceived", nonceRule)
	if err != nil {
		return nil, err
	}
	return &AppChainGatewayParametersReceivedIterator{contract: _AppChainGateway.contract, event: "ParametersReceived", logs: logs, sub: sub}, nil
}

// WatchParametersReceived is a free log subscription operation binding the contract event 0x7fb103ff88be3fb0c788a3ee63031c896dfd0695b79a8061c138e01fd1224b76.
//
// Solidity: event ParametersReceived(uint256 indexed nonce, string[] keys)
func (_AppChainGateway *AppChainGatewayFilterer) WatchParametersReceived(opts *bind.WatchOpts, sink chan<- *AppChainGatewayParametersReceived, nonce []*big.Int) (event.Subscription, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}

	logs, sub, err := _AppChainGateway.contract.WatchLogs(opts, "ParametersReceived", nonceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AppChainGatewayParametersReceived)
				if err := _AppChainGateway.contract.UnpackLog(event, "ParametersReceived", log); err != nil {
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

// ParseParametersReceived is a log parse operation binding the contract event 0x7fb103ff88be3fb0c788a3ee63031c896dfd0695b79a8061c138e01fd1224b76.
//
// Solidity: event ParametersReceived(uint256 indexed nonce, string[] keys)
func (_AppChainGateway *AppChainGatewayFilterer) ParseParametersReceived(log types.Log) (*AppChainGatewayParametersReceived, error) {
	event := new(AppChainGatewayParametersReceived)
	if err := _AppChainGateway.contract.UnpackLog(event, "ParametersReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AppChainGatewayPauseStatusUpdatedIterator is returned from FilterPauseStatusUpdated and is used to iterate over the raw logs and unpacked data for PauseStatusUpdated events raised by the AppChainGateway contract.
type AppChainGatewayPauseStatusUpdatedIterator struct {
	Event *AppChainGatewayPauseStatusUpdated // Event containing the contract specifics and raw log

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
func (it *AppChainGatewayPauseStatusUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AppChainGatewayPauseStatusUpdated)
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
		it.Event = new(AppChainGatewayPauseStatusUpdated)
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
func (it *AppChainGatewayPauseStatusUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AppChainGatewayPauseStatusUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AppChainGatewayPauseStatusUpdated represents a PauseStatusUpdated event raised by the AppChainGateway contract.
type AppChainGatewayPauseStatusUpdated struct {
	Paused bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPauseStatusUpdated is a free log retrieval operation binding the contract event 0x7c4d1fe30fdbfda9e9c4c43e759ef32e4db5128d4cb58ff3ae9583b89b6242a5.
//
// Solidity: event PauseStatusUpdated(bool indexed paused)
func (_AppChainGateway *AppChainGatewayFilterer) FilterPauseStatusUpdated(opts *bind.FilterOpts, paused []bool) (*AppChainGatewayPauseStatusUpdatedIterator, error) {

	var pausedRule []interface{}
	for _, pausedItem := range paused {
		pausedRule = append(pausedRule, pausedItem)
	}

	logs, sub, err := _AppChainGateway.contract.FilterLogs(opts, "PauseStatusUpdated", pausedRule)
	if err != nil {
		return nil, err
	}
	return &AppChainGatewayPauseStatusUpdatedIterator{contract: _AppChainGateway.contract, event: "PauseStatusUpdated", logs: logs, sub: sub}, nil
}

// WatchPauseStatusUpdated is a free log subscription operation binding the contract event 0x7c4d1fe30fdbfda9e9c4c43e759ef32e4db5128d4cb58ff3ae9583b89b6242a5.
//
// Solidity: event PauseStatusUpdated(bool indexed paused)
func (_AppChainGateway *AppChainGatewayFilterer) WatchPauseStatusUpdated(opts *bind.WatchOpts, sink chan<- *AppChainGatewayPauseStatusUpdated, paused []bool) (event.Subscription, error) {

	var pausedRule []interface{}
	for _, pausedItem := range paused {
		pausedRule = append(pausedRule, pausedItem)
	}

	logs, sub, err := _AppChainGateway.contract.WatchLogs(opts, "PauseStatusUpdated", pausedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AppChainGatewayPauseStatusUpdated)
				if err := _AppChainGateway.contract.UnpackLog(event, "PauseStatusUpdated", log); err != nil {
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

// ParsePauseStatusUpdated is a log parse operation binding the contract event 0x7c4d1fe30fdbfda9e9c4c43e759ef32e4db5128d4cb58ff3ae9583b89b6242a5.
//
// Solidity: event PauseStatusUpdated(bool indexed paused)
func (_AppChainGateway *AppChainGatewayFilterer) ParsePauseStatusUpdated(log types.Log) (*AppChainGatewayPauseStatusUpdated, error) {
	event := new(AppChainGatewayPauseStatusUpdated)
	if err := _AppChainGateway.contract.UnpackLog(event, "PauseStatusUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AppChainGatewayUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the AppChainGateway contract.
type AppChainGatewayUpgradedIterator struct {
	Event *AppChainGatewayUpgraded // Event containing the contract specifics and raw log

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
func (it *AppChainGatewayUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AppChainGatewayUpgraded)
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
		it.Event = new(AppChainGatewayUpgraded)
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
func (it *AppChainGatewayUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AppChainGatewayUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AppChainGatewayUpgraded represents a Upgraded event raised by the AppChainGateway contract.
type AppChainGatewayUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_AppChainGateway *AppChainGatewayFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*AppChainGatewayUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _AppChainGateway.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &AppChainGatewayUpgradedIterator{contract: _AppChainGateway.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_AppChainGateway *AppChainGatewayFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *AppChainGatewayUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _AppChainGateway.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AppChainGatewayUpgraded)
				if err := _AppChainGateway.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_AppChainGateway *AppChainGatewayFilterer) ParseUpgraded(log types.Log) (*AppChainGatewayUpgraded, error) {
	event := new(AppChainGatewayUpgraded)
	if err := _AppChainGateway.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AppChainGatewayWithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the AppChainGateway contract.
type AppChainGatewayWithdrawalIterator struct {
	Event *AppChainGatewayWithdrawal // Event containing the contract specifics and raw log

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
func (it *AppChainGatewayWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AppChainGatewayWithdrawal)
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
		it.Event = new(AppChainGatewayWithdrawal)
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
func (it *AppChainGatewayWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AppChainGatewayWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AppChainGatewayWithdrawal represents a Withdrawal event raised by the AppChainGateway contract.
type AppChainGatewayWithdrawal struct {
	Account   common.Address
	MessageId *big.Int
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0xa708f6433a1b53b1e6af0c278ad548516ef5eab45716a7f85657ee720cd2ece0.
//
// Solidity: event Withdrawal(address indexed account, uint256 indexed messageId, address indexed recipient, uint256 amount)
func (_AppChainGateway *AppChainGatewayFilterer) FilterWithdrawal(opts *bind.FilterOpts, account []common.Address, messageId []*big.Int, recipient []common.Address) (*AppChainGatewayWithdrawalIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AppChainGateway.contract.FilterLogs(opts, "Withdrawal", accountRule, messageIdRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &AppChainGatewayWithdrawalIterator{contract: _AppChainGateway.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0xa708f6433a1b53b1e6af0c278ad548516ef5eab45716a7f85657ee720cd2ece0.
//
// Solidity: event Withdrawal(address indexed account, uint256 indexed messageId, address indexed recipient, uint256 amount)
func (_AppChainGateway *AppChainGatewayFilterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *AppChainGatewayWithdrawal, account []common.Address, messageId []*big.Int, recipient []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AppChainGateway.contract.WatchLogs(opts, "Withdrawal", accountRule, messageIdRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AppChainGatewayWithdrawal)
				if err := _AppChainGateway.contract.UnpackLog(event, "Withdrawal", log); err != nil {
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

// ParseWithdrawal is a log parse operation binding the contract event 0xa708f6433a1b53b1e6af0c278ad548516ef5eab45716a7f85657ee720cd2ece0.
//
// Solidity: event Withdrawal(address indexed account, uint256 indexed messageId, address indexed recipient, uint256 amount)
func (_AppChainGateway *AppChainGatewayFilterer) ParseWithdrawal(log types.Log) (*AppChainGatewayWithdrawal, error) {
	event := new(AppChainGatewayWithdrawal)
	if err := _AppChainGateway.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
