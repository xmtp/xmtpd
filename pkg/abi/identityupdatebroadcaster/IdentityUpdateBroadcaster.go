// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package identityupdatebroadcaster

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

// IdentityUpdateBroadcasterMetaData contains all meta data concerning the IdentityUpdateBroadcaster contract.
var IdentityUpdateBroadcasterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"parameterRegistry_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addIdentityUpdate\",\"inputs\":[{\"name\":\"inboxId_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"identityUpdate_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"implementation_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"maxPayloadSize\",\"inputs\":[],\"outputs\":[{\"name\":\"size_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxPayloadSizeParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"migrate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migratorParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"minPayloadSize\",\"inputs\":[],\"outputs\":[{\"name\":\"size_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minPayloadSizeParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"parameterRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"paused_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pausedParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"updateMaxPayloadSize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMinPayloadSize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updatePauseStatus\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"IdentityUpdateCreated\",\"inputs\":[{\"name\":\"inboxId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"update\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"sequenceId\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MaxPayloadSizeUpdated\",\"inputs\":[{\"name\":\"size\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Migrated\",\"inputs\":[{\"name\":\"migrator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinPayloadSizeUpdated\",\"inputs\":[{\"name\":\"size\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PauseStatusUpdated\",\"inputs\":[{\"name\":\"paused\",\"type\":\"bool\",\"indexed\":true,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EmptyCode\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidMaxPayloadSize\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidMinPayloadSize\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPayloadSize\",\"inputs\":[{\"name\":\"actualSize_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minSize_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxSize_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MigrationFailed\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertData_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"NoChange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterOutOfTypeBounds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Paused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroMigrator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroParameterRegistry\",\"inputs\":[]}]",
}

// IdentityUpdateBroadcasterABI is the input ABI used to generate the binding from.
// Deprecated: Use IdentityUpdateBroadcasterMetaData.ABI instead.
var IdentityUpdateBroadcasterABI = IdentityUpdateBroadcasterMetaData.ABI

// IdentityUpdateBroadcaster is an auto generated Go binding around an Ethereum contract.
type IdentityUpdateBroadcaster struct {
	IdentityUpdateBroadcasterCaller     // Read-only binding to the contract
	IdentityUpdateBroadcasterTransactor // Write-only binding to the contract
	IdentityUpdateBroadcasterFilterer   // Log filterer for contract events
}

// IdentityUpdateBroadcasterCaller is an auto generated read-only Go binding around an Ethereum contract.
type IdentityUpdateBroadcasterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityUpdateBroadcasterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IdentityUpdateBroadcasterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityUpdateBroadcasterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IdentityUpdateBroadcasterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityUpdateBroadcasterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IdentityUpdateBroadcasterSession struct {
	Contract     *IdentityUpdateBroadcaster // Generic contract binding to set the session for
	CallOpts     bind.CallOpts              // Call options to use throughout this session
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// IdentityUpdateBroadcasterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IdentityUpdateBroadcasterCallerSession struct {
	Contract *IdentityUpdateBroadcasterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                    // Call options to use throughout this session
}

// IdentityUpdateBroadcasterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IdentityUpdateBroadcasterTransactorSession struct {
	Contract     *IdentityUpdateBroadcasterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// IdentityUpdateBroadcasterRaw is an auto generated low-level Go binding around an Ethereum contract.
type IdentityUpdateBroadcasterRaw struct {
	Contract *IdentityUpdateBroadcaster // Generic contract binding to access the raw methods on
}

// IdentityUpdateBroadcasterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IdentityUpdateBroadcasterCallerRaw struct {
	Contract *IdentityUpdateBroadcasterCaller // Generic read-only contract binding to access the raw methods on
}

// IdentityUpdateBroadcasterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IdentityUpdateBroadcasterTransactorRaw struct {
	Contract *IdentityUpdateBroadcasterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIdentityUpdateBroadcaster creates a new instance of IdentityUpdateBroadcaster, bound to a specific deployed contract.
func NewIdentityUpdateBroadcaster(address common.Address, backend bind.ContractBackend) (*IdentityUpdateBroadcaster, error) {
	contract, err := bindIdentityUpdateBroadcaster(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdateBroadcaster{IdentityUpdateBroadcasterCaller: IdentityUpdateBroadcasterCaller{contract: contract}, IdentityUpdateBroadcasterTransactor: IdentityUpdateBroadcasterTransactor{contract: contract}, IdentityUpdateBroadcasterFilterer: IdentityUpdateBroadcasterFilterer{contract: contract}}, nil
}

// NewIdentityUpdateBroadcasterCaller creates a new read-only instance of IdentityUpdateBroadcaster, bound to a specific deployed contract.
func NewIdentityUpdateBroadcasterCaller(address common.Address, caller bind.ContractCaller) (*IdentityUpdateBroadcasterCaller, error) {
	contract, err := bindIdentityUpdateBroadcaster(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdateBroadcasterCaller{contract: contract}, nil
}

// NewIdentityUpdateBroadcasterTransactor creates a new write-only instance of IdentityUpdateBroadcaster, bound to a specific deployed contract.
func NewIdentityUpdateBroadcasterTransactor(address common.Address, transactor bind.ContractTransactor) (*IdentityUpdateBroadcasterTransactor, error) {
	contract, err := bindIdentityUpdateBroadcaster(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdateBroadcasterTransactor{contract: contract}, nil
}

// NewIdentityUpdateBroadcasterFilterer creates a new log filterer instance of IdentityUpdateBroadcaster, bound to a specific deployed contract.
func NewIdentityUpdateBroadcasterFilterer(address common.Address, filterer bind.ContractFilterer) (*IdentityUpdateBroadcasterFilterer, error) {
	contract, err := bindIdentityUpdateBroadcaster(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdateBroadcasterFilterer{contract: contract}, nil
}

// bindIdentityUpdateBroadcaster binds a generic wrapper to an already deployed contract.
func bindIdentityUpdateBroadcaster(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IdentityUpdateBroadcasterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IdentityUpdateBroadcaster.Contract.IdentityUpdateBroadcasterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.Contract.IdentityUpdateBroadcasterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.Contract.IdentityUpdateBroadcasterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IdentityUpdateBroadcaster.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.Contract.contract.Transact(opts, method, params...)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IdentityUpdateBroadcaster.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterSession) Implementation() (common.Address, error) {
	return _IdentityUpdateBroadcaster.Contract.Implementation(&_IdentityUpdateBroadcaster.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterCallerSession) Implementation() (common.Address, error) {
	return _IdentityUpdateBroadcaster.Contract.Implementation(&_IdentityUpdateBroadcaster.CallOpts)
}

// MaxPayloadSize is a free data retrieval call binding the contract method 0x58e3e94c.
//
// Solidity: function maxPayloadSize() view returns(uint32 size_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterCaller) MaxPayloadSize(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _IdentityUpdateBroadcaster.contract.Call(opts, &out, "maxPayloadSize")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// MaxPayloadSize is a free data retrieval call binding the contract method 0x58e3e94c.
//
// Solidity: function maxPayloadSize() view returns(uint32 size_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterSession) MaxPayloadSize() (uint32, error) {
	return _IdentityUpdateBroadcaster.Contract.MaxPayloadSize(&_IdentityUpdateBroadcaster.CallOpts)
}

// MaxPayloadSize is a free data retrieval call binding the contract method 0x58e3e94c.
//
// Solidity: function maxPayloadSize() view returns(uint32 size_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterCallerSession) MaxPayloadSize() (uint32, error) {
	return _IdentityUpdateBroadcaster.Contract.MaxPayloadSize(&_IdentityUpdateBroadcaster.CallOpts)
}

// MaxPayloadSizeParameterKey is a free data retrieval call binding the contract method 0x0cb858ea.
//
// Solidity: function maxPayloadSizeParameterKey() pure returns(bytes key_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterCaller) MaxPayloadSizeParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _IdentityUpdateBroadcaster.contract.Call(opts, &out, "maxPayloadSizeParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// MaxPayloadSizeParameterKey is a free data retrieval call binding the contract method 0x0cb858ea.
//
// Solidity: function maxPayloadSizeParameterKey() pure returns(bytes key_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterSession) MaxPayloadSizeParameterKey() ([]byte, error) {
	return _IdentityUpdateBroadcaster.Contract.MaxPayloadSizeParameterKey(&_IdentityUpdateBroadcaster.CallOpts)
}

// MaxPayloadSizeParameterKey is a free data retrieval call binding the contract method 0x0cb858ea.
//
// Solidity: function maxPayloadSizeParameterKey() pure returns(bytes key_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterCallerSession) MaxPayloadSizeParameterKey() ([]byte, error) {
	return _IdentityUpdateBroadcaster.Contract.MaxPayloadSizeParameterKey(&_IdentityUpdateBroadcaster.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(bytes key_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterCaller) MigratorParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _IdentityUpdateBroadcaster.contract.Call(opts, &out, "migratorParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(bytes key_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterSession) MigratorParameterKey() ([]byte, error) {
	return _IdentityUpdateBroadcaster.Contract.MigratorParameterKey(&_IdentityUpdateBroadcaster.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(bytes key_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterCallerSession) MigratorParameterKey() ([]byte, error) {
	return _IdentityUpdateBroadcaster.Contract.MigratorParameterKey(&_IdentityUpdateBroadcaster.CallOpts)
}

// MinPayloadSize is a free data retrieval call binding the contract method 0xf96927ac.
//
// Solidity: function minPayloadSize() view returns(uint32 size_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterCaller) MinPayloadSize(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _IdentityUpdateBroadcaster.contract.Call(opts, &out, "minPayloadSize")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// MinPayloadSize is a free data retrieval call binding the contract method 0xf96927ac.
//
// Solidity: function minPayloadSize() view returns(uint32 size_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterSession) MinPayloadSize() (uint32, error) {
	return _IdentityUpdateBroadcaster.Contract.MinPayloadSize(&_IdentityUpdateBroadcaster.CallOpts)
}

// MinPayloadSize is a free data retrieval call binding the contract method 0xf96927ac.
//
// Solidity: function minPayloadSize() view returns(uint32 size_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterCallerSession) MinPayloadSize() (uint32, error) {
	return _IdentityUpdateBroadcaster.Contract.MinPayloadSize(&_IdentityUpdateBroadcaster.CallOpts)
}

// MinPayloadSizeParameterKey is a free data retrieval call binding the contract method 0x9218415d.
//
// Solidity: function minPayloadSizeParameterKey() pure returns(bytes key_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterCaller) MinPayloadSizeParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _IdentityUpdateBroadcaster.contract.Call(opts, &out, "minPayloadSizeParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// MinPayloadSizeParameterKey is a free data retrieval call binding the contract method 0x9218415d.
//
// Solidity: function minPayloadSizeParameterKey() pure returns(bytes key_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterSession) MinPayloadSizeParameterKey() ([]byte, error) {
	return _IdentityUpdateBroadcaster.Contract.MinPayloadSizeParameterKey(&_IdentityUpdateBroadcaster.CallOpts)
}

// MinPayloadSizeParameterKey is a free data retrieval call binding the contract method 0x9218415d.
//
// Solidity: function minPayloadSizeParameterKey() pure returns(bytes key_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterCallerSession) MinPayloadSizeParameterKey() ([]byte, error) {
	return _IdentityUpdateBroadcaster.Contract.MinPayloadSizeParameterKey(&_IdentityUpdateBroadcaster.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterCaller) ParameterRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IdentityUpdateBroadcaster.contract.Call(opts, &out, "parameterRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterSession) ParameterRegistry() (common.Address, error) {
	return _IdentityUpdateBroadcaster.Contract.ParameterRegistry(&_IdentityUpdateBroadcaster.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterCallerSession) ParameterRegistry() (common.Address, error) {
	return _IdentityUpdateBroadcaster.Contract.ParameterRegistry(&_IdentityUpdateBroadcaster.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool paused_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _IdentityUpdateBroadcaster.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool paused_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterSession) Paused() (bool, error) {
	return _IdentityUpdateBroadcaster.Contract.Paused(&_IdentityUpdateBroadcaster.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool paused_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterCallerSession) Paused() (bool, error) {
	return _IdentityUpdateBroadcaster.Contract.Paused(&_IdentityUpdateBroadcaster.CallOpts)
}

// PausedParameterKey is a free data retrieval call binding the contract method 0xcc5999af.
//
// Solidity: function pausedParameterKey() pure returns(bytes key_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterCaller) PausedParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _IdentityUpdateBroadcaster.contract.Call(opts, &out, "pausedParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// PausedParameterKey is a free data retrieval call binding the contract method 0xcc5999af.
//
// Solidity: function pausedParameterKey() pure returns(bytes key_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterSession) PausedParameterKey() ([]byte, error) {
	return _IdentityUpdateBroadcaster.Contract.PausedParameterKey(&_IdentityUpdateBroadcaster.CallOpts)
}

// PausedParameterKey is a free data retrieval call binding the contract method 0xcc5999af.
//
// Solidity: function pausedParameterKey() pure returns(bytes key_)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterCallerSession) PausedParameterKey() ([]byte, error) {
	return _IdentityUpdateBroadcaster.Contract.PausedParameterKey(&_IdentityUpdateBroadcaster.CallOpts)
}

// AddIdentityUpdate is a paid mutator transaction binding the contract method 0xba74fc7c.
//
// Solidity: function addIdentityUpdate(bytes32 inboxId_, bytes identityUpdate_) returns()
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterTransactor) AddIdentityUpdate(opts *bind.TransactOpts, inboxId_ [32]byte, identityUpdate_ []byte) (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.contract.Transact(opts, "addIdentityUpdate", inboxId_, identityUpdate_)
}

// AddIdentityUpdate is a paid mutator transaction binding the contract method 0xba74fc7c.
//
// Solidity: function addIdentityUpdate(bytes32 inboxId_, bytes identityUpdate_) returns()
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterSession) AddIdentityUpdate(inboxId_ [32]byte, identityUpdate_ []byte) (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.Contract.AddIdentityUpdate(&_IdentityUpdateBroadcaster.TransactOpts, inboxId_, identityUpdate_)
}

// AddIdentityUpdate is a paid mutator transaction binding the contract method 0xba74fc7c.
//
// Solidity: function addIdentityUpdate(bytes32 inboxId_, bytes identityUpdate_) returns()
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterTransactorSession) AddIdentityUpdate(inboxId_ [32]byte, identityUpdate_ []byte) (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.Contract.AddIdentityUpdate(&_IdentityUpdateBroadcaster.TransactOpts, inboxId_, identityUpdate_)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterSession) Initialize() (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.Contract.Initialize(&_IdentityUpdateBroadcaster.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterTransactorSession) Initialize() (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.Contract.Initialize(&_IdentityUpdateBroadcaster.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterTransactor) Migrate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.contract.Transact(opts, "migrate")
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterSession) Migrate() (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.Contract.Migrate(&_IdentityUpdateBroadcaster.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterTransactorSession) Migrate() (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.Contract.Migrate(&_IdentityUpdateBroadcaster.TransactOpts)
}

// UpdateMaxPayloadSize is a paid mutator transaction binding the contract method 0x5f643f93.
//
// Solidity: function updateMaxPayloadSize() returns()
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterTransactor) UpdateMaxPayloadSize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.contract.Transact(opts, "updateMaxPayloadSize")
}

// UpdateMaxPayloadSize is a paid mutator transaction binding the contract method 0x5f643f93.
//
// Solidity: function updateMaxPayloadSize() returns()
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterSession) UpdateMaxPayloadSize() (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.Contract.UpdateMaxPayloadSize(&_IdentityUpdateBroadcaster.TransactOpts)
}

// UpdateMaxPayloadSize is a paid mutator transaction binding the contract method 0x5f643f93.
//
// Solidity: function updateMaxPayloadSize() returns()
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterTransactorSession) UpdateMaxPayloadSize() (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.Contract.UpdateMaxPayloadSize(&_IdentityUpdateBroadcaster.TransactOpts)
}

// UpdateMinPayloadSize is a paid mutator transaction binding the contract method 0xd46153ef.
//
// Solidity: function updateMinPayloadSize() returns()
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterTransactor) UpdateMinPayloadSize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.contract.Transact(opts, "updateMinPayloadSize")
}

// UpdateMinPayloadSize is a paid mutator transaction binding the contract method 0xd46153ef.
//
// Solidity: function updateMinPayloadSize() returns()
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterSession) UpdateMinPayloadSize() (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.Contract.UpdateMinPayloadSize(&_IdentityUpdateBroadcaster.TransactOpts)
}

// UpdateMinPayloadSize is a paid mutator transaction binding the contract method 0xd46153ef.
//
// Solidity: function updateMinPayloadSize() returns()
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterTransactorSession) UpdateMinPayloadSize() (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.Contract.UpdateMinPayloadSize(&_IdentityUpdateBroadcaster.TransactOpts)
}

// UpdatePauseStatus is a paid mutator transaction binding the contract method 0x59d4df41.
//
// Solidity: function updatePauseStatus() returns()
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterTransactor) UpdatePauseStatus(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.contract.Transact(opts, "updatePauseStatus")
}

// UpdatePauseStatus is a paid mutator transaction binding the contract method 0x59d4df41.
//
// Solidity: function updatePauseStatus() returns()
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterSession) UpdatePauseStatus() (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.Contract.UpdatePauseStatus(&_IdentityUpdateBroadcaster.TransactOpts)
}

// UpdatePauseStatus is a paid mutator transaction binding the contract method 0x59d4df41.
//
// Solidity: function updatePauseStatus() returns()
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterTransactorSession) UpdatePauseStatus() (*types.Transaction, error) {
	return _IdentityUpdateBroadcaster.Contract.UpdatePauseStatus(&_IdentityUpdateBroadcaster.TransactOpts)
}

// IdentityUpdateBroadcasterIdentityUpdateCreatedIterator is returned from FilterIdentityUpdateCreated and is used to iterate over the raw logs and unpacked data for IdentityUpdateCreated events raised by the IdentityUpdateBroadcaster contract.
type IdentityUpdateBroadcasterIdentityUpdateCreatedIterator struct {
	Event *IdentityUpdateBroadcasterIdentityUpdateCreated // Event containing the contract specifics and raw log

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
func (it *IdentityUpdateBroadcasterIdentityUpdateCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityUpdateBroadcasterIdentityUpdateCreated)
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
		it.Event = new(IdentityUpdateBroadcasterIdentityUpdateCreated)
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
func (it *IdentityUpdateBroadcasterIdentityUpdateCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityUpdateBroadcasterIdentityUpdateCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityUpdateBroadcasterIdentityUpdateCreated represents a IdentityUpdateCreated event raised by the IdentityUpdateBroadcaster contract.
type IdentityUpdateBroadcasterIdentityUpdateCreated struct {
	InboxId    [32]byte
	Update     []byte
	SequenceId uint64
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterIdentityUpdateCreated is a free log retrieval operation binding the contract event 0xc1a40f292090ec0435e939cdfe248e0322a88566679a90a50c4e9e5ef762dbd5.
//
// Solidity: event IdentityUpdateCreated(bytes32 indexed inboxId, bytes update, uint64 indexed sequenceId)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterFilterer) FilterIdentityUpdateCreated(opts *bind.FilterOpts, inboxId [][32]byte, sequenceId []uint64) (*IdentityUpdateBroadcasterIdentityUpdateCreatedIterator, error) {

	var inboxIdRule []interface{}
	for _, inboxIdItem := range inboxId {
		inboxIdRule = append(inboxIdRule, inboxIdItem)
	}

	var sequenceIdRule []interface{}
	for _, sequenceIdItem := range sequenceId {
		sequenceIdRule = append(sequenceIdRule, sequenceIdItem)
	}

	logs, sub, err := _IdentityUpdateBroadcaster.contract.FilterLogs(opts, "IdentityUpdateCreated", inboxIdRule, sequenceIdRule)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdateBroadcasterIdentityUpdateCreatedIterator{contract: _IdentityUpdateBroadcaster.contract, event: "IdentityUpdateCreated", logs: logs, sub: sub}, nil
}

// WatchIdentityUpdateCreated is a free log subscription operation binding the contract event 0xc1a40f292090ec0435e939cdfe248e0322a88566679a90a50c4e9e5ef762dbd5.
//
// Solidity: event IdentityUpdateCreated(bytes32 indexed inboxId, bytes update, uint64 indexed sequenceId)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterFilterer) WatchIdentityUpdateCreated(opts *bind.WatchOpts, sink chan<- *IdentityUpdateBroadcasterIdentityUpdateCreated, inboxId [][32]byte, sequenceId []uint64) (event.Subscription, error) {

	var inboxIdRule []interface{}
	for _, inboxIdItem := range inboxId {
		inboxIdRule = append(inboxIdRule, inboxIdItem)
	}

	var sequenceIdRule []interface{}
	for _, sequenceIdItem := range sequenceId {
		sequenceIdRule = append(sequenceIdRule, sequenceIdItem)
	}

	logs, sub, err := _IdentityUpdateBroadcaster.contract.WatchLogs(opts, "IdentityUpdateCreated", inboxIdRule, sequenceIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityUpdateBroadcasterIdentityUpdateCreated)
				if err := _IdentityUpdateBroadcaster.contract.UnpackLog(event, "IdentityUpdateCreated", log); err != nil {
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

// ParseIdentityUpdateCreated is a log parse operation binding the contract event 0xc1a40f292090ec0435e939cdfe248e0322a88566679a90a50c4e9e5ef762dbd5.
//
// Solidity: event IdentityUpdateCreated(bytes32 indexed inboxId, bytes update, uint64 indexed sequenceId)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterFilterer) ParseIdentityUpdateCreated(log types.Log) (*IdentityUpdateBroadcasterIdentityUpdateCreated, error) {
	event := new(IdentityUpdateBroadcasterIdentityUpdateCreated)
	if err := _IdentityUpdateBroadcaster.contract.UnpackLog(event, "IdentityUpdateCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityUpdateBroadcasterInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the IdentityUpdateBroadcaster contract.
type IdentityUpdateBroadcasterInitializedIterator struct {
	Event *IdentityUpdateBroadcasterInitialized // Event containing the contract specifics and raw log

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
func (it *IdentityUpdateBroadcasterInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityUpdateBroadcasterInitialized)
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
		it.Event = new(IdentityUpdateBroadcasterInitialized)
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
func (it *IdentityUpdateBroadcasterInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityUpdateBroadcasterInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityUpdateBroadcasterInitialized represents a Initialized event raised by the IdentityUpdateBroadcaster contract.
type IdentityUpdateBroadcasterInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterFilterer) FilterInitialized(opts *bind.FilterOpts) (*IdentityUpdateBroadcasterInitializedIterator, error) {

	logs, sub, err := _IdentityUpdateBroadcaster.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &IdentityUpdateBroadcasterInitializedIterator{contract: _IdentityUpdateBroadcaster.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *IdentityUpdateBroadcasterInitialized) (event.Subscription, error) {

	logs, sub, err := _IdentityUpdateBroadcaster.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityUpdateBroadcasterInitialized)
				if err := _IdentityUpdateBroadcaster.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterFilterer) ParseInitialized(log types.Log) (*IdentityUpdateBroadcasterInitialized, error) {
	event := new(IdentityUpdateBroadcasterInitialized)
	if err := _IdentityUpdateBroadcaster.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityUpdateBroadcasterMaxPayloadSizeUpdatedIterator is returned from FilterMaxPayloadSizeUpdated and is used to iterate over the raw logs and unpacked data for MaxPayloadSizeUpdated events raised by the IdentityUpdateBroadcaster contract.
type IdentityUpdateBroadcasterMaxPayloadSizeUpdatedIterator struct {
	Event *IdentityUpdateBroadcasterMaxPayloadSizeUpdated // Event containing the contract specifics and raw log

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
func (it *IdentityUpdateBroadcasterMaxPayloadSizeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityUpdateBroadcasterMaxPayloadSizeUpdated)
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
		it.Event = new(IdentityUpdateBroadcasterMaxPayloadSizeUpdated)
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
func (it *IdentityUpdateBroadcasterMaxPayloadSizeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityUpdateBroadcasterMaxPayloadSizeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityUpdateBroadcasterMaxPayloadSizeUpdated represents a MaxPayloadSizeUpdated event raised by the IdentityUpdateBroadcaster contract.
type IdentityUpdateBroadcasterMaxPayloadSizeUpdated struct {
	Size *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterMaxPayloadSizeUpdated is a free log retrieval operation binding the contract event 0x62422e33fcfc9d38acda2bbddab282a9cc6df7e75f88269fd725bef5457b3045.
//
// Solidity: event MaxPayloadSizeUpdated(uint256 indexed size)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterFilterer) FilterMaxPayloadSizeUpdated(opts *bind.FilterOpts, size []*big.Int) (*IdentityUpdateBroadcasterMaxPayloadSizeUpdatedIterator, error) {

	var sizeRule []interface{}
	for _, sizeItem := range size {
		sizeRule = append(sizeRule, sizeItem)
	}

	logs, sub, err := _IdentityUpdateBroadcaster.contract.FilterLogs(opts, "MaxPayloadSizeUpdated", sizeRule)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdateBroadcasterMaxPayloadSizeUpdatedIterator{contract: _IdentityUpdateBroadcaster.contract, event: "MaxPayloadSizeUpdated", logs: logs, sub: sub}, nil
}

// WatchMaxPayloadSizeUpdated is a free log subscription operation binding the contract event 0x62422e33fcfc9d38acda2bbddab282a9cc6df7e75f88269fd725bef5457b3045.
//
// Solidity: event MaxPayloadSizeUpdated(uint256 indexed size)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterFilterer) WatchMaxPayloadSizeUpdated(opts *bind.WatchOpts, sink chan<- *IdentityUpdateBroadcasterMaxPayloadSizeUpdated, size []*big.Int) (event.Subscription, error) {

	var sizeRule []interface{}
	for _, sizeItem := range size {
		sizeRule = append(sizeRule, sizeItem)
	}

	logs, sub, err := _IdentityUpdateBroadcaster.contract.WatchLogs(opts, "MaxPayloadSizeUpdated", sizeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityUpdateBroadcasterMaxPayloadSizeUpdated)
				if err := _IdentityUpdateBroadcaster.contract.UnpackLog(event, "MaxPayloadSizeUpdated", log); err != nil {
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

// ParseMaxPayloadSizeUpdated is a log parse operation binding the contract event 0x62422e33fcfc9d38acda2bbddab282a9cc6df7e75f88269fd725bef5457b3045.
//
// Solidity: event MaxPayloadSizeUpdated(uint256 indexed size)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterFilterer) ParseMaxPayloadSizeUpdated(log types.Log) (*IdentityUpdateBroadcasterMaxPayloadSizeUpdated, error) {
	event := new(IdentityUpdateBroadcasterMaxPayloadSizeUpdated)
	if err := _IdentityUpdateBroadcaster.contract.UnpackLog(event, "MaxPayloadSizeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityUpdateBroadcasterMigratedIterator is returned from FilterMigrated and is used to iterate over the raw logs and unpacked data for Migrated events raised by the IdentityUpdateBroadcaster contract.
type IdentityUpdateBroadcasterMigratedIterator struct {
	Event *IdentityUpdateBroadcasterMigrated // Event containing the contract specifics and raw log

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
func (it *IdentityUpdateBroadcasterMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityUpdateBroadcasterMigrated)
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
		it.Event = new(IdentityUpdateBroadcasterMigrated)
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
func (it *IdentityUpdateBroadcasterMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityUpdateBroadcasterMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityUpdateBroadcasterMigrated represents a Migrated event raised by the IdentityUpdateBroadcaster contract.
type IdentityUpdateBroadcasterMigrated struct {
	Migrator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMigrated is a free log retrieval operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterFilterer) FilterMigrated(opts *bind.FilterOpts, migrator []common.Address) (*IdentityUpdateBroadcasterMigratedIterator, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _IdentityUpdateBroadcaster.contract.FilterLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdateBroadcasterMigratedIterator{contract: _IdentityUpdateBroadcaster.contract, event: "Migrated", logs: logs, sub: sub}, nil
}

// WatchMigrated is a free log subscription operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterFilterer) WatchMigrated(opts *bind.WatchOpts, sink chan<- *IdentityUpdateBroadcasterMigrated, migrator []common.Address) (event.Subscription, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _IdentityUpdateBroadcaster.contract.WatchLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityUpdateBroadcasterMigrated)
				if err := _IdentityUpdateBroadcaster.contract.UnpackLog(event, "Migrated", log); err != nil {
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
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterFilterer) ParseMigrated(log types.Log) (*IdentityUpdateBroadcasterMigrated, error) {
	event := new(IdentityUpdateBroadcasterMigrated)
	if err := _IdentityUpdateBroadcaster.contract.UnpackLog(event, "Migrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityUpdateBroadcasterMinPayloadSizeUpdatedIterator is returned from FilterMinPayloadSizeUpdated and is used to iterate over the raw logs and unpacked data for MinPayloadSizeUpdated events raised by the IdentityUpdateBroadcaster contract.
type IdentityUpdateBroadcasterMinPayloadSizeUpdatedIterator struct {
	Event *IdentityUpdateBroadcasterMinPayloadSizeUpdated // Event containing the contract specifics and raw log

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
func (it *IdentityUpdateBroadcasterMinPayloadSizeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityUpdateBroadcasterMinPayloadSizeUpdated)
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
		it.Event = new(IdentityUpdateBroadcasterMinPayloadSizeUpdated)
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
func (it *IdentityUpdateBroadcasterMinPayloadSizeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityUpdateBroadcasterMinPayloadSizeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityUpdateBroadcasterMinPayloadSizeUpdated represents a MinPayloadSizeUpdated event raised by the IdentityUpdateBroadcaster contract.
type IdentityUpdateBroadcasterMinPayloadSizeUpdated struct {
	Size *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterMinPayloadSizeUpdated is a free log retrieval operation binding the contract event 0x2caf5b55114860c563b52eba8026a6a0183d9eb1715cbf1c3f8b689f14b5121c.
//
// Solidity: event MinPayloadSizeUpdated(uint256 indexed size)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterFilterer) FilterMinPayloadSizeUpdated(opts *bind.FilterOpts, size []*big.Int) (*IdentityUpdateBroadcasterMinPayloadSizeUpdatedIterator, error) {

	var sizeRule []interface{}
	for _, sizeItem := range size {
		sizeRule = append(sizeRule, sizeItem)
	}

	logs, sub, err := _IdentityUpdateBroadcaster.contract.FilterLogs(opts, "MinPayloadSizeUpdated", sizeRule)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdateBroadcasterMinPayloadSizeUpdatedIterator{contract: _IdentityUpdateBroadcaster.contract, event: "MinPayloadSizeUpdated", logs: logs, sub: sub}, nil
}

// WatchMinPayloadSizeUpdated is a free log subscription operation binding the contract event 0x2caf5b55114860c563b52eba8026a6a0183d9eb1715cbf1c3f8b689f14b5121c.
//
// Solidity: event MinPayloadSizeUpdated(uint256 indexed size)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterFilterer) WatchMinPayloadSizeUpdated(opts *bind.WatchOpts, sink chan<- *IdentityUpdateBroadcasterMinPayloadSizeUpdated, size []*big.Int) (event.Subscription, error) {

	var sizeRule []interface{}
	for _, sizeItem := range size {
		sizeRule = append(sizeRule, sizeItem)
	}

	logs, sub, err := _IdentityUpdateBroadcaster.contract.WatchLogs(opts, "MinPayloadSizeUpdated", sizeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityUpdateBroadcasterMinPayloadSizeUpdated)
				if err := _IdentityUpdateBroadcaster.contract.UnpackLog(event, "MinPayloadSizeUpdated", log); err != nil {
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

// ParseMinPayloadSizeUpdated is a log parse operation binding the contract event 0x2caf5b55114860c563b52eba8026a6a0183d9eb1715cbf1c3f8b689f14b5121c.
//
// Solidity: event MinPayloadSizeUpdated(uint256 indexed size)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterFilterer) ParseMinPayloadSizeUpdated(log types.Log) (*IdentityUpdateBroadcasterMinPayloadSizeUpdated, error) {
	event := new(IdentityUpdateBroadcasterMinPayloadSizeUpdated)
	if err := _IdentityUpdateBroadcaster.contract.UnpackLog(event, "MinPayloadSizeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityUpdateBroadcasterPauseStatusUpdatedIterator is returned from FilterPauseStatusUpdated and is used to iterate over the raw logs and unpacked data for PauseStatusUpdated events raised by the IdentityUpdateBroadcaster contract.
type IdentityUpdateBroadcasterPauseStatusUpdatedIterator struct {
	Event *IdentityUpdateBroadcasterPauseStatusUpdated // Event containing the contract specifics and raw log

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
func (it *IdentityUpdateBroadcasterPauseStatusUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityUpdateBroadcasterPauseStatusUpdated)
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
		it.Event = new(IdentityUpdateBroadcasterPauseStatusUpdated)
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
func (it *IdentityUpdateBroadcasterPauseStatusUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityUpdateBroadcasterPauseStatusUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityUpdateBroadcasterPauseStatusUpdated represents a PauseStatusUpdated event raised by the IdentityUpdateBroadcaster contract.
type IdentityUpdateBroadcasterPauseStatusUpdated struct {
	Paused bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPauseStatusUpdated is a free log retrieval operation binding the contract event 0x7c4d1fe30fdbfda9e9c4c43e759ef32e4db5128d4cb58ff3ae9583b89b6242a5.
//
// Solidity: event PauseStatusUpdated(bool indexed paused)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterFilterer) FilterPauseStatusUpdated(opts *bind.FilterOpts, paused []bool) (*IdentityUpdateBroadcasterPauseStatusUpdatedIterator, error) {

	var pausedRule []interface{}
	for _, pausedItem := range paused {
		pausedRule = append(pausedRule, pausedItem)
	}

	logs, sub, err := _IdentityUpdateBroadcaster.contract.FilterLogs(opts, "PauseStatusUpdated", pausedRule)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdateBroadcasterPauseStatusUpdatedIterator{contract: _IdentityUpdateBroadcaster.contract, event: "PauseStatusUpdated", logs: logs, sub: sub}, nil
}

// WatchPauseStatusUpdated is a free log subscription operation binding the contract event 0x7c4d1fe30fdbfda9e9c4c43e759ef32e4db5128d4cb58ff3ae9583b89b6242a5.
//
// Solidity: event PauseStatusUpdated(bool indexed paused)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterFilterer) WatchPauseStatusUpdated(opts *bind.WatchOpts, sink chan<- *IdentityUpdateBroadcasterPauseStatusUpdated, paused []bool) (event.Subscription, error) {

	var pausedRule []interface{}
	for _, pausedItem := range paused {
		pausedRule = append(pausedRule, pausedItem)
	}

	logs, sub, err := _IdentityUpdateBroadcaster.contract.WatchLogs(opts, "PauseStatusUpdated", pausedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityUpdateBroadcasterPauseStatusUpdated)
				if err := _IdentityUpdateBroadcaster.contract.UnpackLog(event, "PauseStatusUpdated", log); err != nil {
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
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterFilterer) ParsePauseStatusUpdated(log types.Log) (*IdentityUpdateBroadcasterPauseStatusUpdated, error) {
	event := new(IdentityUpdateBroadcasterPauseStatusUpdated)
	if err := _IdentityUpdateBroadcaster.contract.UnpackLog(event, "PauseStatusUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityUpdateBroadcasterUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the IdentityUpdateBroadcaster contract.
type IdentityUpdateBroadcasterUpgradedIterator struct {
	Event *IdentityUpdateBroadcasterUpgraded // Event containing the contract specifics and raw log

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
func (it *IdentityUpdateBroadcasterUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityUpdateBroadcasterUpgraded)
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
		it.Event = new(IdentityUpdateBroadcasterUpgraded)
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
func (it *IdentityUpdateBroadcasterUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityUpdateBroadcasterUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityUpdateBroadcasterUpgraded represents a Upgraded event raised by the IdentityUpdateBroadcaster contract.
type IdentityUpdateBroadcasterUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*IdentityUpdateBroadcasterUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _IdentityUpdateBroadcaster.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdateBroadcasterUpgradedIterator{contract: _IdentityUpdateBroadcaster.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *IdentityUpdateBroadcasterUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _IdentityUpdateBroadcaster.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityUpdateBroadcasterUpgraded)
				if err := _IdentityUpdateBroadcaster.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_IdentityUpdateBroadcaster *IdentityUpdateBroadcasterFilterer) ParseUpgraded(log types.Log) (*IdentityUpdateBroadcasterUpgraded, error) {
	event := new(IdentityUpdateBroadcasterUpgraded)
	if err := _IdentityUpdateBroadcaster.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
