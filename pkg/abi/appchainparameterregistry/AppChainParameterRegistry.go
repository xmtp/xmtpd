// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package appchainparameterregistry

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

// AppChainParameterRegistryMetaData contains all meta data concerning the AppChainParameterRegistry contract.
var AppChainParameterRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"adminParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"get\",\"inputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"value_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"get\",\"inputs\":[{\"name\":\"keys_\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"outputs\":[{\"name\":\"values_\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"implementation_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"admins_\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isAdmin\",\"inputs\":[{\"name\":\"account_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"isAdmin_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"migrate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migratorParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"set\",\"inputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"set\",\"inputs\":[{\"name\":\"keys_\",\"type\":\"string[]\",\"internalType\":\"string[]\"},{\"name\":\"values_\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Migrated\",\"inputs\":[{\"name\":\"migrator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ParameterSet\",\"inputs\":[{\"name\":\"keyHash\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"key\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ArrayLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyAdmins\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyCode\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationFailed\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertData_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"NoKeyComponents\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoKeys\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterOutOfTypeBounds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StringsInsufficientHexLength\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ZeroAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroMigrator\",\"inputs\":[]}]",
}

// AppChainParameterRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use AppChainParameterRegistryMetaData.ABI instead.
var AppChainParameterRegistryABI = AppChainParameterRegistryMetaData.ABI

// AppChainParameterRegistry is an auto generated Go binding around an Ethereum contract.
type AppChainParameterRegistry struct {
	AppChainParameterRegistryCaller     // Read-only binding to the contract
	AppChainParameterRegistryTransactor // Write-only binding to the contract
	AppChainParameterRegistryFilterer   // Log filterer for contract events
}

// AppChainParameterRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type AppChainParameterRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AppChainParameterRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AppChainParameterRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AppChainParameterRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AppChainParameterRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AppChainParameterRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AppChainParameterRegistrySession struct {
	Contract     *AppChainParameterRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts              // Call options to use throughout this session
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// AppChainParameterRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AppChainParameterRegistryCallerSession struct {
	Contract *AppChainParameterRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                    // Call options to use throughout this session
}

// AppChainParameterRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AppChainParameterRegistryTransactorSession struct {
	Contract     *AppChainParameterRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// AppChainParameterRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type AppChainParameterRegistryRaw struct {
	Contract *AppChainParameterRegistry // Generic contract binding to access the raw methods on
}

// AppChainParameterRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AppChainParameterRegistryCallerRaw struct {
	Contract *AppChainParameterRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// AppChainParameterRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AppChainParameterRegistryTransactorRaw struct {
	Contract *AppChainParameterRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAppChainParameterRegistry creates a new instance of AppChainParameterRegistry, bound to a specific deployed contract.
func NewAppChainParameterRegistry(address common.Address, backend bind.ContractBackend) (*AppChainParameterRegistry, error) {
	contract, err := bindAppChainParameterRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AppChainParameterRegistry{AppChainParameterRegistryCaller: AppChainParameterRegistryCaller{contract: contract}, AppChainParameterRegistryTransactor: AppChainParameterRegistryTransactor{contract: contract}, AppChainParameterRegistryFilterer: AppChainParameterRegistryFilterer{contract: contract}}, nil
}

// NewAppChainParameterRegistryCaller creates a new read-only instance of AppChainParameterRegistry, bound to a specific deployed contract.
func NewAppChainParameterRegistryCaller(address common.Address, caller bind.ContractCaller) (*AppChainParameterRegistryCaller, error) {
	contract, err := bindAppChainParameterRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AppChainParameterRegistryCaller{contract: contract}, nil
}

// NewAppChainParameterRegistryTransactor creates a new write-only instance of AppChainParameterRegistry, bound to a specific deployed contract.
func NewAppChainParameterRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*AppChainParameterRegistryTransactor, error) {
	contract, err := bindAppChainParameterRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AppChainParameterRegistryTransactor{contract: contract}, nil
}

// NewAppChainParameterRegistryFilterer creates a new log filterer instance of AppChainParameterRegistry, bound to a specific deployed contract.
func NewAppChainParameterRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*AppChainParameterRegistryFilterer, error) {
	contract, err := bindAppChainParameterRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AppChainParameterRegistryFilterer{contract: contract}, nil
}

// bindAppChainParameterRegistry binds a generic wrapper to an already deployed contract.
func bindAppChainParameterRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AppChainParameterRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AppChainParameterRegistry *AppChainParameterRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AppChainParameterRegistry.Contract.AppChainParameterRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AppChainParameterRegistry *AppChainParameterRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AppChainParameterRegistry.Contract.AppChainParameterRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AppChainParameterRegistry *AppChainParameterRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AppChainParameterRegistry.Contract.AppChainParameterRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AppChainParameterRegistry *AppChainParameterRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AppChainParameterRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AppChainParameterRegistry *AppChainParameterRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AppChainParameterRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AppChainParameterRegistry *AppChainParameterRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AppChainParameterRegistry.Contract.contract.Transact(opts, method, params...)
}

// AdminParameterKey is a free data retrieval call binding the contract method 0x9f40b625.
//
// Solidity: function adminParameterKey() pure returns(string key_)
func (_AppChainParameterRegistry *AppChainParameterRegistryCaller) AdminParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AppChainParameterRegistry.contract.Call(opts, &out, "adminParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// AdminParameterKey is a free data retrieval call binding the contract method 0x9f40b625.
//
// Solidity: function adminParameterKey() pure returns(string key_)
func (_AppChainParameterRegistry *AppChainParameterRegistrySession) AdminParameterKey() (string, error) {
	return _AppChainParameterRegistry.Contract.AdminParameterKey(&_AppChainParameterRegistry.CallOpts)
}

// AdminParameterKey is a free data retrieval call binding the contract method 0x9f40b625.
//
// Solidity: function adminParameterKey() pure returns(string key_)
func (_AppChainParameterRegistry *AppChainParameterRegistryCallerSession) AdminParameterKey() (string, error) {
	return _AppChainParameterRegistry.Contract.AdminParameterKey(&_AppChainParameterRegistry.CallOpts)
}

// Get is a free data retrieval call binding the contract method 0x693ec85e.
//
// Solidity: function get(string key_) view returns(bytes32 value_)
func (_AppChainParameterRegistry *AppChainParameterRegistryCaller) Get(opts *bind.CallOpts, key_ string) ([32]byte, error) {
	var out []interface{}
	err := _AppChainParameterRegistry.contract.Call(opts, &out, "get", key_)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Get is a free data retrieval call binding the contract method 0x693ec85e.
//
// Solidity: function get(string key_) view returns(bytes32 value_)
func (_AppChainParameterRegistry *AppChainParameterRegistrySession) Get(key_ string) ([32]byte, error) {
	return _AppChainParameterRegistry.Contract.Get(&_AppChainParameterRegistry.CallOpts, key_)
}

// Get is a free data retrieval call binding the contract method 0x693ec85e.
//
// Solidity: function get(string key_) view returns(bytes32 value_)
func (_AppChainParameterRegistry *AppChainParameterRegistryCallerSession) Get(key_ string) ([32]byte, error) {
	return _AppChainParameterRegistry.Contract.Get(&_AppChainParameterRegistry.CallOpts, key_)
}

// Get0 is a free data retrieval call binding the contract method 0xb5cbae61.
//
// Solidity: function get(string[] keys_) view returns(bytes32[] values_)
func (_AppChainParameterRegistry *AppChainParameterRegistryCaller) Get0(opts *bind.CallOpts, keys_ []string) ([][32]byte, error) {
	var out []interface{}
	err := _AppChainParameterRegistry.contract.Call(opts, &out, "get0", keys_)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// Get0 is a free data retrieval call binding the contract method 0xb5cbae61.
//
// Solidity: function get(string[] keys_) view returns(bytes32[] values_)
func (_AppChainParameterRegistry *AppChainParameterRegistrySession) Get0(keys_ []string) ([][32]byte, error) {
	return _AppChainParameterRegistry.Contract.Get0(&_AppChainParameterRegistry.CallOpts, keys_)
}

// Get0 is a free data retrieval call binding the contract method 0xb5cbae61.
//
// Solidity: function get(string[] keys_) view returns(bytes32[] values_)
func (_AppChainParameterRegistry *AppChainParameterRegistryCallerSession) Get0(keys_ []string) ([][32]byte, error) {
	return _AppChainParameterRegistry.Contract.Get0(&_AppChainParameterRegistry.CallOpts, keys_)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_AppChainParameterRegistry *AppChainParameterRegistryCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AppChainParameterRegistry.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_AppChainParameterRegistry *AppChainParameterRegistrySession) Implementation() (common.Address, error) {
	return _AppChainParameterRegistry.Contract.Implementation(&_AppChainParameterRegistry.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_AppChainParameterRegistry *AppChainParameterRegistryCallerSession) Implementation() (common.Address, error) {
	return _AppChainParameterRegistry.Contract.Implementation(&_AppChainParameterRegistry.CallOpts)
}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address account_) view returns(bool isAdmin_)
func (_AppChainParameterRegistry *AppChainParameterRegistryCaller) IsAdmin(opts *bind.CallOpts, account_ common.Address) (bool, error) {
	var out []interface{}
	err := _AppChainParameterRegistry.contract.Call(opts, &out, "isAdmin", account_)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address account_) view returns(bool isAdmin_)
func (_AppChainParameterRegistry *AppChainParameterRegistrySession) IsAdmin(account_ common.Address) (bool, error) {
	return _AppChainParameterRegistry.Contract.IsAdmin(&_AppChainParameterRegistry.CallOpts, account_)
}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address account_) view returns(bool isAdmin_)
func (_AppChainParameterRegistry *AppChainParameterRegistryCallerSession) IsAdmin(account_ common.Address) (bool, error) {
	return _AppChainParameterRegistry.Contract.IsAdmin(&_AppChainParameterRegistry.CallOpts, account_)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_AppChainParameterRegistry *AppChainParameterRegistryCaller) MigratorParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AppChainParameterRegistry.contract.Call(opts, &out, "migratorParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_AppChainParameterRegistry *AppChainParameterRegistrySession) MigratorParameterKey() (string, error) {
	return _AppChainParameterRegistry.Contract.MigratorParameterKey(&_AppChainParameterRegistry.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_AppChainParameterRegistry *AppChainParameterRegistryCallerSession) MigratorParameterKey() (string, error) {
	return _AppChainParameterRegistry.Contract.MigratorParameterKey(&_AppChainParameterRegistry.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xa224cee7.
//
// Solidity: function initialize(address[] admins_) returns()
func (_AppChainParameterRegistry *AppChainParameterRegistryTransactor) Initialize(opts *bind.TransactOpts, admins_ []common.Address) (*types.Transaction, error) {
	return _AppChainParameterRegistry.contract.Transact(opts, "initialize", admins_)
}

// Initialize is a paid mutator transaction binding the contract method 0xa224cee7.
//
// Solidity: function initialize(address[] admins_) returns()
func (_AppChainParameterRegistry *AppChainParameterRegistrySession) Initialize(admins_ []common.Address) (*types.Transaction, error) {
	return _AppChainParameterRegistry.Contract.Initialize(&_AppChainParameterRegistry.TransactOpts, admins_)
}

// Initialize is a paid mutator transaction binding the contract method 0xa224cee7.
//
// Solidity: function initialize(address[] admins_) returns()
func (_AppChainParameterRegistry *AppChainParameterRegistryTransactorSession) Initialize(admins_ []common.Address) (*types.Transaction, error) {
	return _AppChainParameterRegistry.Contract.Initialize(&_AppChainParameterRegistry.TransactOpts, admins_)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_AppChainParameterRegistry *AppChainParameterRegistryTransactor) Migrate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AppChainParameterRegistry.contract.Transact(opts, "migrate")
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_AppChainParameterRegistry *AppChainParameterRegistrySession) Migrate() (*types.Transaction, error) {
	return _AppChainParameterRegistry.Contract.Migrate(&_AppChainParameterRegistry.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_AppChainParameterRegistry *AppChainParameterRegistryTransactorSession) Migrate() (*types.Transaction, error) {
	return _AppChainParameterRegistry.Contract.Migrate(&_AppChainParameterRegistry.TransactOpts)
}

// Set is a paid mutator transaction binding the contract method 0x2e3196a5.
//
// Solidity: function set(string key_, bytes32 value_) returns()
func (_AppChainParameterRegistry *AppChainParameterRegistryTransactor) Set(opts *bind.TransactOpts, key_ string, value_ [32]byte) (*types.Transaction, error) {
	return _AppChainParameterRegistry.contract.Transact(opts, "set", key_, value_)
}

// Set is a paid mutator transaction binding the contract method 0x2e3196a5.
//
// Solidity: function set(string key_, bytes32 value_) returns()
func (_AppChainParameterRegistry *AppChainParameterRegistrySession) Set(key_ string, value_ [32]byte) (*types.Transaction, error) {
	return _AppChainParameterRegistry.Contract.Set(&_AppChainParameterRegistry.TransactOpts, key_, value_)
}

// Set is a paid mutator transaction binding the contract method 0x2e3196a5.
//
// Solidity: function set(string key_, bytes32 value_) returns()
func (_AppChainParameterRegistry *AppChainParameterRegistryTransactorSession) Set(key_ string, value_ [32]byte) (*types.Transaction, error) {
	return _AppChainParameterRegistry.Contract.Set(&_AppChainParameterRegistry.TransactOpts, key_, value_)
}

// Set0 is a paid mutator transaction binding the contract method 0xfa482768.
//
// Solidity: function set(string[] keys_, bytes32[] values_) returns()
func (_AppChainParameterRegistry *AppChainParameterRegistryTransactor) Set0(opts *bind.TransactOpts, keys_ []string, values_ [][32]byte) (*types.Transaction, error) {
	return _AppChainParameterRegistry.contract.Transact(opts, "set0", keys_, values_)
}

// Set0 is a paid mutator transaction binding the contract method 0xfa482768.
//
// Solidity: function set(string[] keys_, bytes32[] values_) returns()
func (_AppChainParameterRegistry *AppChainParameterRegistrySession) Set0(keys_ []string, values_ [][32]byte) (*types.Transaction, error) {
	return _AppChainParameterRegistry.Contract.Set0(&_AppChainParameterRegistry.TransactOpts, keys_, values_)
}

// Set0 is a paid mutator transaction binding the contract method 0xfa482768.
//
// Solidity: function set(string[] keys_, bytes32[] values_) returns()
func (_AppChainParameterRegistry *AppChainParameterRegistryTransactorSession) Set0(keys_ []string, values_ [][32]byte) (*types.Transaction, error) {
	return _AppChainParameterRegistry.Contract.Set0(&_AppChainParameterRegistry.TransactOpts, keys_, values_)
}

// AppChainParameterRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the AppChainParameterRegistry contract.
type AppChainParameterRegistryInitializedIterator struct {
	Event *AppChainParameterRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *AppChainParameterRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AppChainParameterRegistryInitialized)
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
		it.Event = new(AppChainParameterRegistryInitialized)
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
func (it *AppChainParameterRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AppChainParameterRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AppChainParameterRegistryInitialized represents a Initialized event raised by the AppChainParameterRegistry contract.
type AppChainParameterRegistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AppChainParameterRegistry *AppChainParameterRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*AppChainParameterRegistryInitializedIterator, error) {

	logs, sub, err := _AppChainParameterRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &AppChainParameterRegistryInitializedIterator{contract: _AppChainParameterRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AppChainParameterRegistry *AppChainParameterRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *AppChainParameterRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _AppChainParameterRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AppChainParameterRegistryInitialized)
				if err := _AppChainParameterRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_AppChainParameterRegistry *AppChainParameterRegistryFilterer) ParseInitialized(log types.Log) (*AppChainParameterRegistryInitialized, error) {
	event := new(AppChainParameterRegistryInitialized)
	if err := _AppChainParameterRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AppChainParameterRegistryMigratedIterator is returned from FilterMigrated and is used to iterate over the raw logs and unpacked data for Migrated events raised by the AppChainParameterRegistry contract.
type AppChainParameterRegistryMigratedIterator struct {
	Event *AppChainParameterRegistryMigrated // Event containing the contract specifics and raw log

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
func (it *AppChainParameterRegistryMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AppChainParameterRegistryMigrated)
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
		it.Event = new(AppChainParameterRegistryMigrated)
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
func (it *AppChainParameterRegistryMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AppChainParameterRegistryMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AppChainParameterRegistryMigrated represents a Migrated event raised by the AppChainParameterRegistry contract.
type AppChainParameterRegistryMigrated struct {
	Migrator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMigrated is a free log retrieval operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_AppChainParameterRegistry *AppChainParameterRegistryFilterer) FilterMigrated(opts *bind.FilterOpts, migrator []common.Address) (*AppChainParameterRegistryMigratedIterator, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _AppChainParameterRegistry.contract.FilterLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return &AppChainParameterRegistryMigratedIterator{contract: _AppChainParameterRegistry.contract, event: "Migrated", logs: logs, sub: sub}, nil
}

// WatchMigrated is a free log subscription operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_AppChainParameterRegistry *AppChainParameterRegistryFilterer) WatchMigrated(opts *bind.WatchOpts, sink chan<- *AppChainParameterRegistryMigrated, migrator []common.Address) (event.Subscription, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _AppChainParameterRegistry.contract.WatchLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AppChainParameterRegistryMigrated)
				if err := _AppChainParameterRegistry.contract.UnpackLog(event, "Migrated", log); err != nil {
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
func (_AppChainParameterRegistry *AppChainParameterRegistryFilterer) ParseMigrated(log types.Log) (*AppChainParameterRegistryMigrated, error) {
	event := new(AppChainParameterRegistryMigrated)
	if err := _AppChainParameterRegistry.contract.UnpackLog(event, "Migrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AppChainParameterRegistryParameterSetIterator is returned from FilterParameterSet and is used to iterate over the raw logs and unpacked data for ParameterSet events raised by the AppChainParameterRegistry contract.
type AppChainParameterRegistryParameterSetIterator struct {
	Event *AppChainParameterRegistryParameterSet // Event containing the contract specifics and raw log

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
func (it *AppChainParameterRegistryParameterSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AppChainParameterRegistryParameterSet)
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
		it.Event = new(AppChainParameterRegistryParameterSet)
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
func (it *AppChainParameterRegistryParameterSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AppChainParameterRegistryParameterSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AppChainParameterRegistryParameterSet represents a ParameterSet event raised by the AppChainParameterRegistry contract.
type AppChainParameterRegistryParameterSet struct {
	KeyHash common.Hash
	Key     string
	Value   [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterParameterSet is a free log retrieval operation binding the contract event 0xe5fb4bb3c30225c1b5468d38191df0224e5e7d3d7e1486e66aed36e313b5069b.
//
// Solidity: event ParameterSet(string indexed keyHash, string key, bytes32 value)
func (_AppChainParameterRegistry *AppChainParameterRegistryFilterer) FilterParameterSet(opts *bind.FilterOpts, keyHash []string) (*AppChainParameterRegistryParameterSetIterator, error) {

	var keyHashRule []interface{}
	for _, keyHashItem := range keyHash {
		keyHashRule = append(keyHashRule, keyHashItem)
	}

	logs, sub, err := _AppChainParameterRegistry.contract.FilterLogs(opts, "ParameterSet", keyHashRule)
	if err != nil {
		return nil, err
	}
	return &AppChainParameterRegistryParameterSetIterator{contract: _AppChainParameterRegistry.contract, event: "ParameterSet", logs: logs, sub: sub}, nil
}

// WatchParameterSet is a free log subscription operation binding the contract event 0xe5fb4bb3c30225c1b5468d38191df0224e5e7d3d7e1486e66aed36e313b5069b.
//
// Solidity: event ParameterSet(string indexed keyHash, string key, bytes32 value)
func (_AppChainParameterRegistry *AppChainParameterRegistryFilterer) WatchParameterSet(opts *bind.WatchOpts, sink chan<- *AppChainParameterRegistryParameterSet, keyHash []string) (event.Subscription, error) {

	var keyHashRule []interface{}
	for _, keyHashItem := range keyHash {
		keyHashRule = append(keyHashRule, keyHashItem)
	}

	logs, sub, err := _AppChainParameterRegistry.contract.WatchLogs(opts, "ParameterSet", keyHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AppChainParameterRegistryParameterSet)
				if err := _AppChainParameterRegistry.contract.UnpackLog(event, "ParameterSet", log); err != nil {
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

// ParseParameterSet is a log parse operation binding the contract event 0xe5fb4bb3c30225c1b5468d38191df0224e5e7d3d7e1486e66aed36e313b5069b.
//
// Solidity: event ParameterSet(string indexed keyHash, string key, bytes32 value)
func (_AppChainParameterRegistry *AppChainParameterRegistryFilterer) ParseParameterSet(log types.Log) (*AppChainParameterRegistryParameterSet, error) {
	event := new(AppChainParameterRegistryParameterSet)
	if err := _AppChainParameterRegistry.contract.UnpackLog(event, "ParameterSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AppChainParameterRegistryUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the AppChainParameterRegistry contract.
type AppChainParameterRegistryUpgradedIterator struct {
	Event *AppChainParameterRegistryUpgraded // Event containing the contract specifics and raw log

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
func (it *AppChainParameterRegistryUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AppChainParameterRegistryUpgraded)
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
		it.Event = new(AppChainParameterRegistryUpgraded)
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
func (it *AppChainParameterRegistryUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AppChainParameterRegistryUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AppChainParameterRegistryUpgraded represents a Upgraded event raised by the AppChainParameterRegistry contract.
type AppChainParameterRegistryUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_AppChainParameterRegistry *AppChainParameterRegistryFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*AppChainParameterRegistryUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _AppChainParameterRegistry.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &AppChainParameterRegistryUpgradedIterator{contract: _AppChainParameterRegistry.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_AppChainParameterRegistry *AppChainParameterRegistryFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *AppChainParameterRegistryUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _AppChainParameterRegistry.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AppChainParameterRegistryUpgraded)
				if err := _AppChainParameterRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_AppChainParameterRegistry *AppChainParameterRegistryFilterer) ParseUpgraded(log types.Log) (*AppChainParameterRegistryUpgraded, error) {
	event := new(AppChainParameterRegistryUpgraded)
	if err := _AppChainParameterRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
