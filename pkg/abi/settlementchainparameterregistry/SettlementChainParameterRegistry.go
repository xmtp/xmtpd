// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package settlementchainparameterregistry

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

// SettlementChainParameterRegistryMetaData contains all meta data concerning the SettlementChainParameterRegistry contract.
var SettlementChainParameterRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"adminParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"get\",\"inputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"value_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"get\",\"inputs\":[{\"name\":\"keys_\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"outputs\":[{\"name\":\"values_\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"implementation_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"admins_\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isAdmin\",\"inputs\":[{\"name\":\"account_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"isAdmin_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"migrate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migratorParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"set\",\"inputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"set\",\"inputs\":[{\"name\":\"keys_\",\"type\":\"string[]\",\"internalType\":\"string[]\"},{\"name\":\"values_\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Migrated\",\"inputs\":[{\"name\":\"migrator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ParameterSet\",\"inputs\":[{\"name\":\"keyHash\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"key\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ArrayLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyAdmins\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyCode\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationFailed\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertData_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"NoKeyComponents\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoKeys\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterOutOfTypeBounds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StringsInsufficientHexLength\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ZeroAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroMigrator\",\"inputs\":[]}]",
}

// SettlementChainParameterRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use SettlementChainParameterRegistryMetaData.ABI instead.
var SettlementChainParameterRegistryABI = SettlementChainParameterRegistryMetaData.ABI

// SettlementChainParameterRegistry is an auto generated Go binding around an Ethereum contract.
type SettlementChainParameterRegistry struct {
	SettlementChainParameterRegistryCaller     // Read-only binding to the contract
	SettlementChainParameterRegistryTransactor // Write-only binding to the contract
	SettlementChainParameterRegistryFilterer   // Log filterer for contract events
}

// SettlementChainParameterRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type SettlementChainParameterRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementChainParameterRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SettlementChainParameterRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementChainParameterRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SettlementChainParameterRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementChainParameterRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SettlementChainParameterRegistrySession struct {
	Contract     *SettlementChainParameterRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                     // Call options to use throughout this session
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// SettlementChainParameterRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SettlementChainParameterRegistryCallerSession struct {
	Contract *SettlementChainParameterRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                           // Call options to use throughout this session
}

// SettlementChainParameterRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SettlementChainParameterRegistryTransactorSession struct {
	Contract     *SettlementChainParameterRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                           // Transaction auth options to use throughout this session
}

// SettlementChainParameterRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type SettlementChainParameterRegistryRaw struct {
	Contract *SettlementChainParameterRegistry // Generic contract binding to access the raw methods on
}

// SettlementChainParameterRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SettlementChainParameterRegistryCallerRaw struct {
	Contract *SettlementChainParameterRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// SettlementChainParameterRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SettlementChainParameterRegistryTransactorRaw struct {
	Contract *SettlementChainParameterRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSettlementChainParameterRegistry creates a new instance of SettlementChainParameterRegistry, bound to a specific deployed contract.
func NewSettlementChainParameterRegistry(address common.Address, backend bind.ContractBackend) (*SettlementChainParameterRegistry, error) {
	contract, err := bindSettlementChainParameterRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SettlementChainParameterRegistry{SettlementChainParameterRegistryCaller: SettlementChainParameterRegistryCaller{contract: contract}, SettlementChainParameterRegistryTransactor: SettlementChainParameterRegistryTransactor{contract: contract}, SettlementChainParameterRegistryFilterer: SettlementChainParameterRegistryFilterer{contract: contract}}, nil
}

// NewSettlementChainParameterRegistryCaller creates a new read-only instance of SettlementChainParameterRegistry, bound to a specific deployed contract.
func NewSettlementChainParameterRegistryCaller(address common.Address, caller bind.ContractCaller) (*SettlementChainParameterRegistryCaller, error) {
	contract, err := bindSettlementChainParameterRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SettlementChainParameterRegistryCaller{contract: contract}, nil
}

// NewSettlementChainParameterRegistryTransactor creates a new write-only instance of SettlementChainParameterRegistry, bound to a specific deployed contract.
func NewSettlementChainParameterRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*SettlementChainParameterRegistryTransactor, error) {
	contract, err := bindSettlementChainParameterRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SettlementChainParameterRegistryTransactor{contract: contract}, nil
}

// NewSettlementChainParameterRegistryFilterer creates a new log filterer instance of SettlementChainParameterRegistry, bound to a specific deployed contract.
func NewSettlementChainParameterRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*SettlementChainParameterRegistryFilterer, error) {
	contract, err := bindSettlementChainParameterRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SettlementChainParameterRegistryFilterer{contract: contract}, nil
}

// bindSettlementChainParameterRegistry binds a generic wrapper to an already deployed contract.
func bindSettlementChainParameterRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SettlementChainParameterRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SettlementChainParameterRegistry.Contract.SettlementChainParameterRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.SettlementChainParameterRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.SettlementChainParameterRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SettlementChainParameterRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.contract.Transact(opts, method, params...)
}

// AdminParameterKey is a free data retrieval call binding the contract method 0x9f40b625.
//
// Solidity: function adminParameterKey() pure returns(string key_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCaller) AdminParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SettlementChainParameterRegistry.contract.Call(opts, &out, "adminParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// AdminParameterKey is a free data retrieval call binding the contract method 0x9f40b625.
//
// Solidity: function adminParameterKey() pure returns(string key_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistrySession) AdminParameterKey() (string, error) {
	return _SettlementChainParameterRegistry.Contract.AdminParameterKey(&_SettlementChainParameterRegistry.CallOpts)
}

// AdminParameterKey is a free data retrieval call binding the contract method 0x9f40b625.
//
// Solidity: function adminParameterKey() pure returns(string key_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCallerSession) AdminParameterKey() (string, error) {
	return _SettlementChainParameterRegistry.Contract.AdminParameterKey(&_SettlementChainParameterRegistry.CallOpts)
}

// Get is a free data retrieval call binding the contract method 0x693ec85e.
//
// Solidity: function get(string key_) view returns(bytes32 value_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCaller) Get(opts *bind.CallOpts, key_ string) ([32]byte, error) {
	var out []interface{}
	err := _SettlementChainParameterRegistry.contract.Call(opts, &out, "get", key_)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Get is a free data retrieval call binding the contract method 0x693ec85e.
//
// Solidity: function get(string key_) view returns(bytes32 value_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistrySession) Get(key_ string) ([32]byte, error) {
	return _SettlementChainParameterRegistry.Contract.Get(&_SettlementChainParameterRegistry.CallOpts, key_)
}

// Get is a free data retrieval call binding the contract method 0x693ec85e.
//
// Solidity: function get(string key_) view returns(bytes32 value_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCallerSession) Get(key_ string) ([32]byte, error) {
	return _SettlementChainParameterRegistry.Contract.Get(&_SettlementChainParameterRegistry.CallOpts, key_)
}

// Get0 is a free data retrieval call binding the contract method 0xb5cbae61.
//
// Solidity: function get(string[] keys_) view returns(bytes32[] values_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCaller) Get0(opts *bind.CallOpts, keys_ []string) ([][32]byte, error) {
	var out []interface{}
	err := _SettlementChainParameterRegistry.contract.Call(opts, &out, "get0", keys_)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// Get0 is a free data retrieval call binding the contract method 0xb5cbae61.
//
// Solidity: function get(string[] keys_) view returns(bytes32[] values_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistrySession) Get0(keys_ []string) ([][32]byte, error) {
	return _SettlementChainParameterRegistry.Contract.Get0(&_SettlementChainParameterRegistry.CallOpts, keys_)
}

// Get0 is a free data retrieval call binding the contract method 0xb5cbae61.
//
// Solidity: function get(string[] keys_) view returns(bytes32[] values_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCallerSession) Get0(keys_ []string) ([][32]byte, error) {
	return _SettlementChainParameterRegistry.Contract.Get0(&_SettlementChainParameterRegistry.CallOpts, keys_)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SettlementChainParameterRegistry.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistrySession) Implementation() (common.Address, error) {
	return _SettlementChainParameterRegistry.Contract.Implementation(&_SettlementChainParameterRegistry.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCallerSession) Implementation() (common.Address, error) {
	return _SettlementChainParameterRegistry.Contract.Implementation(&_SettlementChainParameterRegistry.CallOpts)
}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address account_) view returns(bool isAdmin_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCaller) IsAdmin(opts *bind.CallOpts, account_ common.Address) (bool, error) {
	var out []interface{}
	err := _SettlementChainParameterRegistry.contract.Call(opts, &out, "isAdmin", account_)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address account_) view returns(bool isAdmin_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistrySession) IsAdmin(account_ common.Address) (bool, error) {
	return _SettlementChainParameterRegistry.Contract.IsAdmin(&_SettlementChainParameterRegistry.CallOpts, account_)
}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address account_) view returns(bool isAdmin_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCallerSession) IsAdmin(account_ common.Address) (bool, error) {
	return _SettlementChainParameterRegistry.Contract.IsAdmin(&_SettlementChainParameterRegistry.CallOpts, account_)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCaller) MigratorParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SettlementChainParameterRegistry.contract.Call(opts, &out, "migratorParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistrySession) MigratorParameterKey() (string, error) {
	return _SettlementChainParameterRegistry.Contract.MigratorParameterKey(&_SettlementChainParameterRegistry.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCallerSession) MigratorParameterKey() (string, error) {
	return _SettlementChainParameterRegistry.Contract.MigratorParameterKey(&_SettlementChainParameterRegistry.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xa224cee7.
//
// Solidity: function initialize(address[] admins_) returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryTransactor) Initialize(opts *bind.TransactOpts, admins_ []common.Address) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.contract.Transact(opts, "initialize", admins_)
}

// Initialize is a paid mutator transaction binding the contract method 0xa224cee7.
//
// Solidity: function initialize(address[] admins_) returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistrySession) Initialize(admins_ []common.Address) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.Initialize(&_SettlementChainParameterRegistry.TransactOpts, admins_)
}

// Initialize is a paid mutator transaction binding the contract method 0xa224cee7.
//
// Solidity: function initialize(address[] admins_) returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryTransactorSession) Initialize(admins_ []common.Address) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.Initialize(&_SettlementChainParameterRegistry.TransactOpts, admins_)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryTransactor) Migrate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.contract.Transact(opts, "migrate")
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistrySession) Migrate() (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.Migrate(&_SettlementChainParameterRegistry.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryTransactorSession) Migrate() (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.Migrate(&_SettlementChainParameterRegistry.TransactOpts)
}

// Set is a paid mutator transaction binding the contract method 0x2e3196a5.
//
// Solidity: function set(string key_, bytes32 value_) returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryTransactor) Set(opts *bind.TransactOpts, key_ string, value_ [32]byte) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.contract.Transact(opts, "set", key_, value_)
}

// Set is a paid mutator transaction binding the contract method 0x2e3196a5.
//
// Solidity: function set(string key_, bytes32 value_) returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistrySession) Set(key_ string, value_ [32]byte) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.Set(&_SettlementChainParameterRegistry.TransactOpts, key_, value_)
}

// Set is a paid mutator transaction binding the contract method 0x2e3196a5.
//
// Solidity: function set(string key_, bytes32 value_) returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryTransactorSession) Set(key_ string, value_ [32]byte) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.Set(&_SettlementChainParameterRegistry.TransactOpts, key_, value_)
}

// Set0 is a paid mutator transaction binding the contract method 0xfa482768.
//
// Solidity: function set(string[] keys_, bytes32[] values_) returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryTransactor) Set0(opts *bind.TransactOpts, keys_ []string, values_ [][32]byte) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.contract.Transact(opts, "set0", keys_, values_)
}

// Set0 is a paid mutator transaction binding the contract method 0xfa482768.
//
// Solidity: function set(string[] keys_, bytes32[] values_) returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistrySession) Set0(keys_ []string, values_ [][32]byte) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.Set0(&_SettlementChainParameterRegistry.TransactOpts, keys_, values_)
}

// Set0 is a paid mutator transaction binding the contract method 0xfa482768.
//
// Solidity: function set(string[] keys_, bytes32[] values_) returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryTransactorSession) Set0(keys_ []string, values_ [][32]byte) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.Set0(&_SettlementChainParameterRegistry.TransactOpts, keys_, values_)
}

// SettlementChainParameterRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the SettlementChainParameterRegistry contract.
type SettlementChainParameterRegistryInitializedIterator struct {
	Event *SettlementChainParameterRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *SettlementChainParameterRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementChainParameterRegistryInitialized)
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
		it.Event = new(SettlementChainParameterRegistryInitialized)
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
func (it *SettlementChainParameterRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementChainParameterRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementChainParameterRegistryInitialized represents a Initialized event raised by the SettlementChainParameterRegistry contract.
type SettlementChainParameterRegistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*SettlementChainParameterRegistryInitializedIterator, error) {

	logs, sub, err := _SettlementChainParameterRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &SettlementChainParameterRegistryInitializedIterator{contract: _SettlementChainParameterRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *SettlementChainParameterRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _SettlementChainParameterRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementChainParameterRegistryInitialized)
				if err := _SettlementChainParameterRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) ParseInitialized(log types.Log) (*SettlementChainParameterRegistryInitialized, error) {
	event := new(SettlementChainParameterRegistryInitialized)
	if err := _SettlementChainParameterRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementChainParameterRegistryMigratedIterator is returned from FilterMigrated and is used to iterate over the raw logs and unpacked data for Migrated events raised by the SettlementChainParameterRegistry contract.
type SettlementChainParameterRegistryMigratedIterator struct {
	Event *SettlementChainParameterRegistryMigrated // Event containing the contract specifics and raw log

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
func (it *SettlementChainParameterRegistryMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementChainParameterRegistryMigrated)
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
		it.Event = new(SettlementChainParameterRegistryMigrated)
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
func (it *SettlementChainParameterRegistryMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementChainParameterRegistryMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementChainParameterRegistryMigrated represents a Migrated event raised by the SettlementChainParameterRegistry contract.
type SettlementChainParameterRegistryMigrated struct {
	Migrator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMigrated is a free log retrieval operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) FilterMigrated(opts *bind.FilterOpts, migrator []common.Address) (*SettlementChainParameterRegistryMigratedIterator, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _SettlementChainParameterRegistry.contract.FilterLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return &SettlementChainParameterRegistryMigratedIterator{contract: _SettlementChainParameterRegistry.contract, event: "Migrated", logs: logs, sub: sub}, nil
}

// WatchMigrated is a free log subscription operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) WatchMigrated(opts *bind.WatchOpts, sink chan<- *SettlementChainParameterRegistryMigrated, migrator []common.Address) (event.Subscription, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _SettlementChainParameterRegistry.contract.WatchLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementChainParameterRegistryMigrated)
				if err := _SettlementChainParameterRegistry.contract.UnpackLog(event, "Migrated", log); err != nil {
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
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) ParseMigrated(log types.Log) (*SettlementChainParameterRegistryMigrated, error) {
	event := new(SettlementChainParameterRegistryMigrated)
	if err := _SettlementChainParameterRegistry.contract.UnpackLog(event, "Migrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementChainParameterRegistryParameterSetIterator is returned from FilterParameterSet and is used to iterate over the raw logs and unpacked data for ParameterSet events raised by the SettlementChainParameterRegistry contract.
type SettlementChainParameterRegistryParameterSetIterator struct {
	Event *SettlementChainParameterRegistryParameterSet // Event containing the contract specifics and raw log

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
func (it *SettlementChainParameterRegistryParameterSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementChainParameterRegistryParameterSet)
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
		it.Event = new(SettlementChainParameterRegistryParameterSet)
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
func (it *SettlementChainParameterRegistryParameterSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementChainParameterRegistryParameterSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementChainParameterRegistryParameterSet represents a ParameterSet event raised by the SettlementChainParameterRegistry contract.
type SettlementChainParameterRegistryParameterSet struct {
	KeyHash common.Hash
	Key     string
	Value   [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterParameterSet is a free log retrieval operation binding the contract event 0xe5fb4bb3c30225c1b5468d38191df0224e5e7d3d7e1486e66aed36e313b5069b.
//
// Solidity: event ParameterSet(string indexed keyHash, string key, bytes32 value)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) FilterParameterSet(opts *bind.FilterOpts, keyHash []string) (*SettlementChainParameterRegistryParameterSetIterator, error) {

	var keyHashRule []interface{}
	for _, keyHashItem := range keyHash {
		keyHashRule = append(keyHashRule, keyHashItem)
	}

	logs, sub, err := _SettlementChainParameterRegistry.contract.FilterLogs(opts, "ParameterSet", keyHashRule)
	if err != nil {
		return nil, err
	}
	return &SettlementChainParameterRegistryParameterSetIterator{contract: _SettlementChainParameterRegistry.contract, event: "ParameterSet", logs: logs, sub: sub}, nil
}

// WatchParameterSet is a free log subscription operation binding the contract event 0xe5fb4bb3c30225c1b5468d38191df0224e5e7d3d7e1486e66aed36e313b5069b.
//
// Solidity: event ParameterSet(string indexed keyHash, string key, bytes32 value)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) WatchParameterSet(opts *bind.WatchOpts, sink chan<- *SettlementChainParameterRegistryParameterSet, keyHash []string) (event.Subscription, error) {

	var keyHashRule []interface{}
	for _, keyHashItem := range keyHash {
		keyHashRule = append(keyHashRule, keyHashItem)
	}

	logs, sub, err := _SettlementChainParameterRegistry.contract.WatchLogs(opts, "ParameterSet", keyHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementChainParameterRegistryParameterSet)
				if err := _SettlementChainParameterRegistry.contract.UnpackLog(event, "ParameterSet", log); err != nil {
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
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) ParseParameterSet(log types.Log) (*SettlementChainParameterRegistryParameterSet, error) {
	event := new(SettlementChainParameterRegistryParameterSet)
	if err := _SettlementChainParameterRegistry.contract.UnpackLog(event, "ParameterSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementChainParameterRegistryUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the SettlementChainParameterRegistry contract.
type SettlementChainParameterRegistryUpgradedIterator struct {
	Event *SettlementChainParameterRegistryUpgraded // Event containing the contract specifics and raw log

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
func (it *SettlementChainParameterRegistryUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementChainParameterRegistryUpgraded)
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
		it.Event = new(SettlementChainParameterRegistryUpgraded)
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
func (it *SettlementChainParameterRegistryUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementChainParameterRegistryUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementChainParameterRegistryUpgraded represents a Upgraded event raised by the SettlementChainParameterRegistry contract.
type SettlementChainParameterRegistryUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*SettlementChainParameterRegistryUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _SettlementChainParameterRegistry.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &SettlementChainParameterRegistryUpgradedIterator{contract: _SettlementChainParameterRegistry.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *SettlementChainParameterRegistryUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _SettlementChainParameterRegistry.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementChainParameterRegistryUpgraded)
				if err := _SettlementChainParameterRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) ParseUpgraded(log types.Log) (*SettlementChainParameterRegistryUpgraded, error) {
	event := new(SettlementChainParameterRegistryUpgraded)
	if err := _SettlementChainParameterRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
