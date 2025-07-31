// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package groupmessagebroadcaster

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

// GroupMessageBroadcasterMetaData contains all meta data concerning the GroupMessageBroadcaster contract.
var GroupMessageBroadcasterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"parameterRegistry_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addMessage\",\"inputs\":[{\"name\":\"groupId_\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"},{\"name\":\"message_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"bootstrapMessages\",\"inputs\":[{\"name\":\"groupIds_\",\"type\":\"bytes16[]\",\"internalType\":\"bytes16[]\"},{\"name\":\"messages_\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"sequenceIds_\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"implementation_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"maxPayloadSize\",\"inputs\":[],\"outputs\":[{\"name\":\"size_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxPayloadSizeParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"migrate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migratorParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"minPayloadSize\",\"inputs\":[],\"outputs\":[{\"name\":\"size_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minPayloadSizeParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"parameterRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"paused_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pausedParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"payloadBootstrapper\",\"inputs\":[],\"outputs\":[{\"name\":\"payloadBootstrapper_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"payloadBootstrapperParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"updateMaxPayloadSize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMinPayloadSize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updatePauseStatus\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updatePayloadBootstrapper\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MaxPayloadSizeUpdated\",\"inputs\":[{\"name\":\"size\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MessageSent\",\"inputs\":[{\"name\":\"groupId\",\"type\":\"bytes16\",\"indexed\":true,\"internalType\":\"bytes16\"},{\"name\":\"message\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"sequenceId\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Migrated\",\"inputs\":[{\"name\":\"migrator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinPayloadSizeUpdated\",\"inputs\":[{\"name\":\"size\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PauseStatusUpdated\",\"inputs\":[{\"name\":\"paused\",\"type\":\"bool\",\"indexed\":true,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayloadBootstrapperUpdated\",\"inputs\":[{\"name\":\"payloadBootstrapper\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ArrayLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyArray\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyCode\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidMaxPayloadSize\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidMinPayloadSize\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPayloadSize\",\"inputs\":[{\"name\":\"actualSize_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minSize_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxSize_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MigrationFailed\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertData_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"NoChange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotPayloadBootstrapper\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterOutOfTypeBounds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Paused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroMigrator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroParameterRegistry\",\"inputs\":[]}]",
}

// GroupMessageBroadcasterABI is the input ABI used to generate the binding from.
// Deprecated: Use GroupMessageBroadcasterMetaData.ABI instead.
var GroupMessageBroadcasterABI = GroupMessageBroadcasterMetaData.ABI

// GroupMessageBroadcaster is an auto generated Go binding around an Ethereum contract.
type GroupMessageBroadcaster struct {
	GroupMessageBroadcasterCaller     // Read-only binding to the contract
	GroupMessageBroadcasterTransactor // Write-only binding to the contract
	GroupMessageBroadcasterFilterer   // Log filterer for contract events
}

// GroupMessageBroadcasterCaller is an auto generated read-only Go binding around an Ethereum contract.
type GroupMessageBroadcasterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GroupMessageBroadcasterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GroupMessageBroadcasterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GroupMessageBroadcasterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GroupMessageBroadcasterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GroupMessageBroadcasterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GroupMessageBroadcasterSession struct {
	Contract     *GroupMessageBroadcaster // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// GroupMessageBroadcasterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GroupMessageBroadcasterCallerSession struct {
	Contract *GroupMessageBroadcasterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// GroupMessageBroadcasterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GroupMessageBroadcasterTransactorSession struct {
	Contract     *GroupMessageBroadcasterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// GroupMessageBroadcasterRaw is an auto generated low-level Go binding around an Ethereum contract.
type GroupMessageBroadcasterRaw struct {
	Contract *GroupMessageBroadcaster // Generic contract binding to access the raw methods on
}

// GroupMessageBroadcasterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GroupMessageBroadcasterCallerRaw struct {
	Contract *GroupMessageBroadcasterCaller // Generic read-only contract binding to access the raw methods on
}

// GroupMessageBroadcasterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GroupMessageBroadcasterTransactorRaw struct {
	Contract *GroupMessageBroadcasterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGroupMessageBroadcaster creates a new instance of GroupMessageBroadcaster, bound to a specific deployed contract.
func NewGroupMessageBroadcaster(address common.Address, backend bind.ContractBackend) (*GroupMessageBroadcaster, error) {
	contract, err := bindGroupMessageBroadcaster(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcaster{GroupMessageBroadcasterCaller: GroupMessageBroadcasterCaller{contract: contract}, GroupMessageBroadcasterTransactor: GroupMessageBroadcasterTransactor{contract: contract}, GroupMessageBroadcasterFilterer: GroupMessageBroadcasterFilterer{contract: contract}}, nil
}

// NewGroupMessageBroadcasterCaller creates a new read-only instance of GroupMessageBroadcaster, bound to a specific deployed contract.
func NewGroupMessageBroadcasterCaller(address common.Address, caller bind.ContractCaller) (*GroupMessageBroadcasterCaller, error) {
	contract, err := bindGroupMessageBroadcaster(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterCaller{contract: contract}, nil
}

// NewGroupMessageBroadcasterTransactor creates a new write-only instance of GroupMessageBroadcaster, bound to a specific deployed contract.
func NewGroupMessageBroadcasterTransactor(address common.Address, transactor bind.ContractTransactor) (*GroupMessageBroadcasterTransactor, error) {
	contract, err := bindGroupMessageBroadcaster(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterTransactor{contract: contract}, nil
}

// NewGroupMessageBroadcasterFilterer creates a new log filterer instance of GroupMessageBroadcaster, bound to a specific deployed contract.
func NewGroupMessageBroadcasterFilterer(address common.Address, filterer bind.ContractFilterer) (*GroupMessageBroadcasterFilterer, error) {
	contract, err := bindGroupMessageBroadcaster(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterFilterer{contract: contract}, nil
}

// bindGroupMessageBroadcaster binds a generic wrapper to an already deployed contract.
func bindGroupMessageBroadcaster(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GroupMessageBroadcasterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GroupMessageBroadcaster *GroupMessageBroadcasterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GroupMessageBroadcaster.Contract.GroupMessageBroadcasterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GroupMessageBroadcaster *GroupMessageBroadcasterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.GroupMessageBroadcasterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GroupMessageBroadcaster *GroupMessageBroadcasterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.GroupMessageBroadcasterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GroupMessageBroadcaster.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.contract.Transact(opts, method, params...)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) Implementation() (common.Address, error) {
	return _GroupMessageBroadcaster.Contract.Implementation(&_GroupMessageBroadcaster.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) Implementation() (common.Address, error) {
	return _GroupMessageBroadcaster.Contract.Implementation(&_GroupMessageBroadcaster.CallOpts)
}

// MaxPayloadSize is a free data retrieval call binding the contract method 0x58e3e94c.
//
// Solidity: function maxPayloadSize() view returns(uint32 size_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) MaxPayloadSize(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "maxPayloadSize")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// MaxPayloadSize is a free data retrieval call binding the contract method 0x58e3e94c.
//
// Solidity: function maxPayloadSize() view returns(uint32 size_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) MaxPayloadSize() (uint32, error) {
	return _GroupMessageBroadcaster.Contract.MaxPayloadSize(&_GroupMessageBroadcaster.CallOpts)
}

// MaxPayloadSize is a free data retrieval call binding the contract method 0x58e3e94c.
//
// Solidity: function maxPayloadSize() view returns(uint32 size_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) MaxPayloadSize() (uint32, error) {
	return _GroupMessageBroadcaster.Contract.MaxPayloadSize(&_GroupMessageBroadcaster.CallOpts)
}

// MaxPayloadSizeParameterKey is a free data retrieval call binding the contract method 0x0cb858ea.
//
// Solidity: function maxPayloadSizeParameterKey() pure returns(string key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) MaxPayloadSizeParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "maxPayloadSizeParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// MaxPayloadSizeParameterKey is a free data retrieval call binding the contract method 0x0cb858ea.
//
// Solidity: function maxPayloadSizeParameterKey() pure returns(string key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) MaxPayloadSizeParameterKey() (string, error) {
	return _GroupMessageBroadcaster.Contract.MaxPayloadSizeParameterKey(&_GroupMessageBroadcaster.CallOpts)
}

// MaxPayloadSizeParameterKey is a free data retrieval call binding the contract method 0x0cb858ea.
//
// Solidity: function maxPayloadSizeParameterKey() pure returns(string key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) MaxPayloadSizeParameterKey() (string, error) {
	return _GroupMessageBroadcaster.Contract.MaxPayloadSizeParameterKey(&_GroupMessageBroadcaster.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) MigratorParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "migratorParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) MigratorParameterKey() (string, error) {
	return _GroupMessageBroadcaster.Contract.MigratorParameterKey(&_GroupMessageBroadcaster.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) MigratorParameterKey() (string, error) {
	return _GroupMessageBroadcaster.Contract.MigratorParameterKey(&_GroupMessageBroadcaster.CallOpts)
}

// MinPayloadSize is a free data retrieval call binding the contract method 0xf96927ac.
//
// Solidity: function minPayloadSize() view returns(uint32 size_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) MinPayloadSize(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "minPayloadSize")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// MinPayloadSize is a free data retrieval call binding the contract method 0xf96927ac.
//
// Solidity: function minPayloadSize() view returns(uint32 size_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) MinPayloadSize() (uint32, error) {
	return _GroupMessageBroadcaster.Contract.MinPayloadSize(&_GroupMessageBroadcaster.CallOpts)
}

// MinPayloadSize is a free data retrieval call binding the contract method 0xf96927ac.
//
// Solidity: function minPayloadSize() view returns(uint32 size_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) MinPayloadSize() (uint32, error) {
	return _GroupMessageBroadcaster.Contract.MinPayloadSize(&_GroupMessageBroadcaster.CallOpts)
}

// MinPayloadSizeParameterKey is a free data retrieval call binding the contract method 0x9218415d.
//
// Solidity: function minPayloadSizeParameterKey() pure returns(string key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) MinPayloadSizeParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "minPayloadSizeParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// MinPayloadSizeParameterKey is a free data retrieval call binding the contract method 0x9218415d.
//
// Solidity: function minPayloadSizeParameterKey() pure returns(string key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) MinPayloadSizeParameterKey() (string, error) {
	return _GroupMessageBroadcaster.Contract.MinPayloadSizeParameterKey(&_GroupMessageBroadcaster.CallOpts)
}

// MinPayloadSizeParameterKey is a free data retrieval call binding the contract method 0x9218415d.
//
// Solidity: function minPayloadSizeParameterKey() pure returns(string key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) MinPayloadSizeParameterKey() (string, error) {
	return _GroupMessageBroadcaster.Contract.MinPayloadSizeParameterKey(&_GroupMessageBroadcaster.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) ParameterRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "parameterRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) ParameterRegistry() (common.Address, error) {
	return _GroupMessageBroadcaster.Contract.ParameterRegistry(&_GroupMessageBroadcaster.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) ParameterRegistry() (common.Address, error) {
	return _GroupMessageBroadcaster.Contract.ParameterRegistry(&_GroupMessageBroadcaster.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool paused_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool paused_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) Paused() (bool, error) {
	return _GroupMessageBroadcaster.Contract.Paused(&_GroupMessageBroadcaster.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool paused_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) Paused() (bool, error) {
	return _GroupMessageBroadcaster.Contract.Paused(&_GroupMessageBroadcaster.CallOpts)
}

// PausedParameterKey is a free data retrieval call binding the contract method 0xcc5999af.
//
// Solidity: function pausedParameterKey() pure returns(string key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) PausedParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "pausedParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// PausedParameterKey is a free data retrieval call binding the contract method 0xcc5999af.
//
// Solidity: function pausedParameterKey() pure returns(string key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) PausedParameterKey() (string, error) {
	return _GroupMessageBroadcaster.Contract.PausedParameterKey(&_GroupMessageBroadcaster.CallOpts)
}

// PausedParameterKey is a free data retrieval call binding the contract method 0xcc5999af.
//
// Solidity: function pausedParameterKey() pure returns(string key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) PausedParameterKey() (string, error) {
	return _GroupMessageBroadcaster.Contract.PausedParameterKey(&_GroupMessageBroadcaster.CallOpts)
}

// PayloadBootstrapper is a free data retrieval call binding the contract method 0x405a11fc.
//
// Solidity: function payloadBootstrapper() view returns(address payloadBootstrapper_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) PayloadBootstrapper(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "payloadBootstrapper")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PayloadBootstrapper is a free data retrieval call binding the contract method 0x405a11fc.
//
// Solidity: function payloadBootstrapper() view returns(address payloadBootstrapper_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) PayloadBootstrapper() (common.Address, error) {
	return _GroupMessageBroadcaster.Contract.PayloadBootstrapper(&_GroupMessageBroadcaster.CallOpts)
}

// PayloadBootstrapper is a free data retrieval call binding the contract method 0x405a11fc.
//
// Solidity: function payloadBootstrapper() view returns(address payloadBootstrapper_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) PayloadBootstrapper() (common.Address, error) {
	return _GroupMessageBroadcaster.Contract.PayloadBootstrapper(&_GroupMessageBroadcaster.CallOpts)
}

// PayloadBootstrapperParameterKey is a free data retrieval call binding the contract method 0x4600f300.
//
// Solidity: function payloadBootstrapperParameterKey() pure returns(string key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) PayloadBootstrapperParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "payloadBootstrapperParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// PayloadBootstrapperParameterKey is a free data retrieval call binding the contract method 0x4600f300.
//
// Solidity: function payloadBootstrapperParameterKey() pure returns(string key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) PayloadBootstrapperParameterKey() (string, error) {
	return _GroupMessageBroadcaster.Contract.PayloadBootstrapperParameterKey(&_GroupMessageBroadcaster.CallOpts)
}

// PayloadBootstrapperParameterKey is a free data retrieval call binding the contract method 0x4600f300.
//
// Solidity: function payloadBootstrapperParameterKey() pure returns(string key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) PayloadBootstrapperParameterKey() (string, error) {
	return _GroupMessageBroadcaster.Contract.PayloadBootstrapperParameterKey(&_GroupMessageBroadcaster.CallOpts)
}

// AddMessage is a paid mutator transaction binding the contract method 0x7e4af76c.
//
// Solidity: function addMessage(bytes16 groupId_, bytes message_) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactor) AddMessage(opts *bind.TransactOpts, groupId_ [16]byte, message_ []byte) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.contract.Transact(opts, "addMessage", groupId_, message_)
}

// AddMessage is a paid mutator transaction binding the contract method 0x7e4af76c.
//
// Solidity: function addMessage(bytes16 groupId_, bytes message_) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) AddMessage(groupId_ [16]byte, message_ []byte) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.AddMessage(&_GroupMessageBroadcaster.TransactOpts, groupId_, message_)
}

// AddMessage is a paid mutator transaction binding the contract method 0x7e4af76c.
//
// Solidity: function addMessage(bytes16 groupId_, bytes message_) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorSession) AddMessage(groupId_ [16]byte, message_ []byte) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.AddMessage(&_GroupMessageBroadcaster.TransactOpts, groupId_, message_)
}

// BootstrapMessages is a paid mutator transaction binding the contract method 0xcbc3a5ea.
//
// Solidity: function bootstrapMessages(bytes16[] groupIds_, bytes[] messages_, uint64[] sequenceIds_) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactor) BootstrapMessages(opts *bind.TransactOpts, groupIds_ [][16]byte, messages_ [][]byte, sequenceIds_ []uint64) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.contract.Transact(opts, "bootstrapMessages", groupIds_, messages_, sequenceIds_)
}

// BootstrapMessages is a paid mutator transaction binding the contract method 0xcbc3a5ea.
//
// Solidity: function bootstrapMessages(bytes16[] groupIds_, bytes[] messages_, uint64[] sequenceIds_) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) BootstrapMessages(groupIds_ [][16]byte, messages_ [][]byte, sequenceIds_ []uint64) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.BootstrapMessages(&_GroupMessageBroadcaster.TransactOpts, groupIds_, messages_, sequenceIds_)
}

// BootstrapMessages is a paid mutator transaction binding the contract method 0xcbc3a5ea.
//
// Solidity: function bootstrapMessages(bytes16[] groupIds_, bytes[] messages_, uint64[] sequenceIds_) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorSession) BootstrapMessages(groupIds_ [][16]byte, messages_ [][]byte, sequenceIds_ []uint64) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.BootstrapMessages(&_GroupMessageBroadcaster.TransactOpts, groupIds_, messages_, sequenceIds_)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) Initialize() (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.Initialize(&_GroupMessageBroadcaster.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorSession) Initialize() (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.Initialize(&_GroupMessageBroadcaster.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactor) Migrate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.contract.Transact(opts, "migrate")
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) Migrate() (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.Migrate(&_GroupMessageBroadcaster.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorSession) Migrate() (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.Migrate(&_GroupMessageBroadcaster.TransactOpts)
}

// UpdateMaxPayloadSize is a paid mutator transaction binding the contract method 0x5f643f93.
//
// Solidity: function updateMaxPayloadSize() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactor) UpdateMaxPayloadSize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.contract.Transact(opts, "updateMaxPayloadSize")
}

// UpdateMaxPayloadSize is a paid mutator transaction binding the contract method 0x5f643f93.
//
// Solidity: function updateMaxPayloadSize() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) UpdateMaxPayloadSize() (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.UpdateMaxPayloadSize(&_GroupMessageBroadcaster.TransactOpts)
}

// UpdateMaxPayloadSize is a paid mutator transaction binding the contract method 0x5f643f93.
//
// Solidity: function updateMaxPayloadSize() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorSession) UpdateMaxPayloadSize() (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.UpdateMaxPayloadSize(&_GroupMessageBroadcaster.TransactOpts)
}

// UpdateMinPayloadSize is a paid mutator transaction binding the contract method 0xd46153ef.
//
// Solidity: function updateMinPayloadSize() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactor) UpdateMinPayloadSize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.contract.Transact(opts, "updateMinPayloadSize")
}

// UpdateMinPayloadSize is a paid mutator transaction binding the contract method 0xd46153ef.
//
// Solidity: function updateMinPayloadSize() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) UpdateMinPayloadSize() (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.UpdateMinPayloadSize(&_GroupMessageBroadcaster.TransactOpts)
}

// UpdateMinPayloadSize is a paid mutator transaction binding the contract method 0xd46153ef.
//
// Solidity: function updateMinPayloadSize() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorSession) UpdateMinPayloadSize() (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.UpdateMinPayloadSize(&_GroupMessageBroadcaster.TransactOpts)
}

// UpdatePauseStatus is a paid mutator transaction binding the contract method 0x59d4df41.
//
// Solidity: function updatePauseStatus() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactor) UpdatePauseStatus(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.contract.Transact(opts, "updatePauseStatus")
}

// UpdatePauseStatus is a paid mutator transaction binding the contract method 0x59d4df41.
//
// Solidity: function updatePauseStatus() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) UpdatePauseStatus() (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.UpdatePauseStatus(&_GroupMessageBroadcaster.TransactOpts)
}

// UpdatePauseStatus is a paid mutator transaction binding the contract method 0x59d4df41.
//
// Solidity: function updatePauseStatus() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorSession) UpdatePauseStatus() (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.UpdatePauseStatus(&_GroupMessageBroadcaster.TransactOpts)
}

// UpdatePayloadBootstrapper is a paid mutator transaction binding the contract method 0x886bd989.
//
// Solidity: function updatePayloadBootstrapper() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactor) UpdatePayloadBootstrapper(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.contract.Transact(opts, "updatePayloadBootstrapper")
}

// UpdatePayloadBootstrapper is a paid mutator transaction binding the contract method 0x886bd989.
//
// Solidity: function updatePayloadBootstrapper() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) UpdatePayloadBootstrapper() (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.UpdatePayloadBootstrapper(&_GroupMessageBroadcaster.TransactOpts)
}

// UpdatePayloadBootstrapper is a paid mutator transaction binding the contract method 0x886bd989.
//
// Solidity: function updatePayloadBootstrapper() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorSession) UpdatePayloadBootstrapper() (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.UpdatePayloadBootstrapper(&_GroupMessageBroadcaster.TransactOpts)
}

// GroupMessageBroadcasterInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterInitializedIterator struct {
	Event *GroupMessageBroadcasterInitialized // Event containing the contract specifics and raw log

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
func (it *GroupMessageBroadcasterInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessageBroadcasterInitialized)
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
		it.Event = new(GroupMessageBroadcasterInitialized)
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
func (it *GroupMessageBroadcasterInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessageBroadcasterInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessageBroadcasterInitialized represents a Initialized event raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) FilterInitialized(opts *bind.FilterOpts) (*GroupMessageBroadcasterInitializedIterator, error) {

	logs, sub, err := _GroupMessageBroadcaster.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterInitializedIterator{contract: _GroupMessageBroadcaster.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *GroupMessageBroadcasterInitialized) (event.Subscription, error) {

	logs, sub, err := _GroupMessageBroadcaster.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessageBroadcasterInitialized)
				if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) ParseInitialized(log types.Log) (*GroupMessageBroadcasterInitialized, error) {
	event := new(GroupMessageBroadcasterInitialized)
	if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessageBroadcasterMaxPayloadSizeUpdatedIterator is returned from FilterMaxPayloadSizeUpdated and is used to iterate over the raw logs and unpacked data for MaxPayloadSizeUpdated events raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterMaxPayloadSizeUpdatedIterator struct {
	Event *GroupMessageBroadcasterMaxPayloadSizeUpdated // Event containing the contract specifics and raw log

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
func (it *GroupMessageBroadcasterMaxPayloadSizeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessageBroadcasterMaxPayloadSizeUpdated)
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
		it.Event = new(GroupMessageBroadcasterMaxPayloadSizeUpdated)
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
func (it *GroupMessageBroadcasterMaxPayloadSizeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessageBroadcasterMaxPayloadSizeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessageBroadcasterMaxPayloadSizeUpdated represents a MaxPayloadSizeUpdated event raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterMaxPayloadSizeUpdated struct {
	Size *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterMaxPayloadSizeUpdated is a free log retrieval operation binding the contract event 0x62422e33fcfc9d38acda2bbddab282a9cc6df7e75f88269fd725bef5457b3045.
//
// Solidity: event MaxPayloadSizeUpdated(uint256 indexed size)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) FilterMaxPayloadSizeUpdated(opts *bind.FilterOpts, size []*big.Int) (*GroupMessageBroadcasterMaxPayloadSizeUpdatedIterator, error) {

	var sizeRule []interface{}
	for _, sizeItem := range size {
		sizeRule = append(sizeRule, sizeItem)
	}

	logs, sub, err := _GroupMessageBroadcaster.contract.FilterLogs(opts, "MaxPayloadSizeUpdated", sizeRule)
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterMaxPayloadSizeUpdatedIterator{contract: _GroupMessageBroadcaster.contract, event: "MaxPayloadSizeUpdated", logs: logs, sub: sub}, nil
}

// WatchMaxPayloadSizeUpdated is a free log subscription operation binding the contract event 0x62422e33fcfc9d38acda2bbddab282a9cc6df7e75f88269fd725bef5457b3045.
//
// Solidity: event MaxPayloadSizeUpdated(uint256 indexed size)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) WatchMaxPayloadSizeUpdated(opts *bind.WatchOpts, sink chan<- *GroupMessageBroadcasterMaxPayloadSizeUpdated, size []*big.Int) (event.Subscription, error) {

	var sizeRule []interface{}
	for _, sizeItem := range size {
		sizeRule = append(sizeRule, sizeItem)
	}

	logs, sub, err := _GroupMessageBroadcaster.contract.WatchLogs(opts, "MaxPayloadSizeUpdated", sizeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessageBroadcasterMaxPayloadSizeUpdated)
				if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "MaxPayloadSizeUpdated", log); err != nil {
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
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) ParseMaxPayloadSizeUpdated(log types.Log) (*GroupMessageBroadcasterMaxPayloadSizeUpdated, error) {
	event := new(GroupMessageBroadcasterMaxPayloadSizeUpdated)
	if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "MaxPayloadSizeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessageBroadcasterMessageSentIterator is returned from FilterMessageSent and is used to iterate over the raw logs and unpacked data for MessageSent events raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterMessageSentIterator struct {
	Event *GroupMessageBroadcasterMessageSent // Event containing the contract specifics and raw log

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
func (it *GroupMessageBroadcasterMessageSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessageBroadcasterMessageSent)
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
		it.Event = new(GroupMessageBroadcasterMessageSent)
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
func (it *GroupMessageBroadcasterMessageSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessageBroadcasterMessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessageBroadcasterMessageSent represents a MessageSent event raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterMessageSent struct {
	GroupId    [16]byte
	Message    []byte
	SequenceId uint64
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterMessageSent is a free log retrieval operation binding the contract event 0xe69329a8fa6c24860e47ba6211a332cf49c3e692bdbcc4bf5500d724bf9ccd05.
//
// Solidity: event MessageSent(bytes16 indexed groupId, bytes message, uint64 indexed sequenceId)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) FilterMessageSent(opts *bind.FilterOpts, groupId [][16]byte, sequenceId []uint64) (*GroupMessageBroadcasterMessageSentIterator, error) {

	var groupIdRule []interface{}
	for _, groupIdItem := range groupId {
		groupIdRule = append(groupIdRule, groupIdItem)
	}

	var sequenceIdRule []interface{}
	for _, sequenceIdItem := range sequenceId {
		sequenceIdRule = append(sequenceIdRule, sequenceIdItem)
	}

	logs, sub, err := _GroupMessageBroadcaster.contract.FilterLogs(opts, "MessageSent", groupIdRule, sequenceIdRule)
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterMessageSentIterator{contract: _GroupMessageBroadcaster.contract, event: "MessageSent", logs: logs, sub: sub}, nil
}

// WatchMessageSent is a free log subscription operation binding the contract event 0xe69329a8fa6c24860e47ba6211a332cf49c3e692bdbcc4bf5500d724bf9ccd05.
//
// Solidity: event MessageSent(bytes16 indexed groupId, bytes message, uint64 indexed sequenceId)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) WatchMessageSent(opts *bind.WatchOpts, sink chan<- *GroupMessageBroadcasterMessageSent, groupId [][16]byte, sequenceId []uint64) (event.Subscription, error) {

	var groupIdRule []interface{}
	for _, groupIdItem := range groupId {
		groupIdRule = append(groupIdRule, groupIdItem)
	}

	var sequenceIdRule []interface{}
	for _, sequenceIdItem := range sequenceId {
		sequenceIdRule = append(sequenceIdRule, sequenceIdItem)
	}

	logs, sub, err := _GroupMessageBroadcaster.contract.WatchLogs(opts, "MessageSent", groupIdRule, sequenceIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessageBroadcasterMessageSent)
				if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "MessageSent", log); err != nil {
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

// ParseMessageSent is a log parse operation binding the contract event 0xe69329a8fa6c24860e47ba6211a332cf49c3e692bdbcc4bf5500d724bf9ccd05.
//
// Solidity: event MessageSent(bytes16 indexed groupId, bytes message, uint64 indexed sequenceId)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) ParseMessageSent(log types.Log) (*GroupMessageBroadcasterMessageSent, error) {
	event := new(GroupMessageBroadcasterMessageSent)
	if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "MessageSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessageBroadcasterMigratedIterator is returned from FilterMigrated and is used to iterate over the raw logs and unpacked data for Migrated events raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterMigratedIterator struct {
	Event *GroupMessageBroadcasterMigrated // Event containing the contract specifics and raw log

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
func (it *GroupMessageBroadcasterMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessageBroadcasterMigrated)
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
		it.Event = new(GroupMessageBroadcasterMigrated)
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
func (it *GroupMessageBroadcasterMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessageBroadcasterMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessageBroadcasterMigrated represents a Migrated event raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterMigrated struct {
	Migrator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMigrated is a free log retrieval operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) FilterMigrated(opts *bind.FilterOpts, migrator []common.Address) (*GroupMessageBroadcasterMigratedIterator, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _GroupMessageBroadcaster.contract.FilterLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterMigratedIterator{contract: _GroupMessageBroadcaster.contract, event: "Migrated", logs: logs, sub: sub}, nil
}

// WatchMigrated is a free log subscription operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) WatchMigrated(opts *bind.WatchOpts, sink chan<- *GroupMessageBroadcasterMigrated, migrator []common.Address) (event.Subscription, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _GroupMessageBroadcaster.contract.WatchLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessageBroadcasterMigrated)
				if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "Migrated", log); err != nil {
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
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) ParseMigrated(log types.Log) (*GroupMessageBroadcasterMigrated, error) {
	event := new(GroupMessageBroadcasterMigrated)
	if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "Migrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessageBroadcasterMinPayloadSizeUpdatedIterator is returned from FilterMinPayloadSizeUpdated and is used to iterate over the raw logs and unpacked data for MinPayloadSizeUpdated events raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterMinPayloadSizeUpdatedIterator struct {
	Event *GroupMessageBroadcasterMinPayloadSizeUpdated // Event containing the contract specifics and raw log

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
func (it *GroupMessageBroadcasterMinPayloadSizeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessageBroadcasterMinPayloadSizeUpdated)
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
		it.Event = new(GroupMessageBroadcasterMinPayloadSizeUpdated)
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
func (it *GroupMessageBroadcasterMinPayloadSizeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessageBroadcasterMinPayloadSizeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessageBroadcasterMinPayloadSizeUpdated represents a MinPayloadSizeUpdated event raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterMinPayloadSizeUpdated struct {
	Size *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterMinPayloadSizeUpdated is a free log retrieval operation binding the contract event 0x2caf5b55114860c563b52eba8026a6a0183d9eb1715cbf1c3f8b689f14b5121c.
//
// Solidity: event MinPayloadSizeUpdated(uint256 indexed size)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) FilterMinPayloadSizeUpdated(opts *bind.FilterOpts, size []*big.Int) (*GroupMessageBroadcasterMinPayloadSizeUpdatedIterator, error) {

	var sizeRule []interface{}
	for _, sizeItem := range size {
		sizeRule = append(sizeRule, sizeItem)
	}

	logs, sub, err := _GroupMessageBroadcaster.contract.FilterLogs(opts, "MinPayloadSizeUpdated", sizeRule)
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterMinPayloadSizeUpdatedIterator{contract: _GroupMessageBroadcaster.contract, event: "MinPayloadSizeUpdated", logs: logs, sub: sub}, nil
}

// WatchMinPayloadSizeUpdated is a free log subscription operation binding the contract event 0x2caf5b55114860c563b52eba8026a6a0183d9eb1715cbf1c3f8b689f14b5121c.
//
// Solidity: event MinPayloadSizeUpdated(uint256 indexed size)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) WatchMinPayloadSizeUpdated(opts *bind.WatchOpts, sink chan<- *GroupMessageBroadcasterMinPayloadSizeUpdated, size []*big.Int) (event.Subscription, error) {

	var sizeRule []interface{}
	for _, sizeItem := range size {
		sizeRule = append(sizeRule, sizeItem)
	}

	logs, sub, err := _GroupMessageBroadcaster.contract.WatchLogs(opts, "MinPayloadSizeUpdated", sizeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessageBroadcasterMinPayloadSizeUpdated)
				if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "MinPayloadSizeUpdated", log); err != nil {
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
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) ParseMinPayloadSizeUpdated(log types.Log) (*GroupMessageBroadcasterMinPayloadSizeUpdated, error) {
	event := new(GroupMessageBroadcasterMinPayloadSizeUpdated)
	if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "MinPayloadSizeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessageBroadcasterPauseStatusUpdatedIterator is returned from FilterPauseStatusUpdated and is used to iterate over the raw logs and unpacked data for PauseStatusUpdated events raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterPauseStatusUpdatedIterator struct {
	Event *GroupMessageBroadcasterPauseStatusUpdated // Event containing the contract specifics and raw log

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
func (it *GroupMessageBroadcasterPauseStatusUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessageBroadcasterPauseStatusUpdated)
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
		it.Event = new(GroupMessageBroadcasterPauseStatusUpdated)
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
func (it *GroupMessageBroadcasterPauseStatusUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessageBroadcasterPauseStatusUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessageBroadcasterPauseStatusUpdated represents a PauseStatusUpdated event raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterPauseStatusUpdated struct {
	Paused bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPauseStatusUpdated is a free log retrieval operation binding the contract event 0x7c4d1fe30fdbfda9e9c4c43e759ef32e4db5128d4cb58ff3ae9583b89b6242a5.
//
// Solidity: event PauseStatusUpdated(bool indexed paused)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) FilterPauseStatusUpdated(opts *bind.FilterOpts, paused []bool) (*GroupMessageBroadcasterPauseStatusUpdatedIterator, error) {

	var pausedRule []interface{}
	for _, pausedItem := range paused {
		pausedRule = append(pausedRule, pausedItem)
	}

	logs, sub, err := _GroupMessageBroadcaster.contract.FilterLogs(opts, "PauseStatusUpdated", pausedRule)
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterPauseStatusUpdatedIterator{contract: _GroupMessageBroadcaster.contract, event: "PauseStatusUpdated", logs: logs, sub: sub}, nil
}

// WatchPauseStatusUpdated is a free log subscription operation binding the contract event 0x7c4d1fe30fdbfda9e9c4c43e759ef32e4db5128d4cb58ff3ae9583b89b6242a5.
//
// Solidity: event PauseStatusUpdated(bool indexed paused)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) WatchPauseStatusUpdated(opts *bind.WatchOpts, sink chan<- *GroupMessageBroadcasterPauseStatusUpdated, paused []bool) (event.Subscription, error) {

	var pausedRule []interface{}
	for _, pausedItem := range paused {
		pausedRule = append(pausedRule, pausedItem)
	}

	logs, sub, err := _GroupMessageBroadcaster.contract.WatchLogs(opts, "PauseStatusUpdated", pausedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessageBroadcasterPauseStatusUpdated)
				if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "PauseStatusUpdated", log); err != nil {
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
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) ParsePauseStatusUpdated(log types.Log) (*GroupMessageBroadcasterPauseStatusUpdated, error) {
	event := new(GroupMessageBroadcasterPauseStatusUpdated)
	if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "PauseStatusUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessageBroadcasterPayloadBootstrapperUpdatedIterator is returned from FilterPayloadBootstrapperUpdated and is used to iterate over the raw logs and unpacked data for PayloadBootstrapperUpdated events raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterPayloadBootstrapperUpdatedIterator struct {
	Event *GroupMessageBroadcasterPayloadBootstrapperUpdated // Event containing the contract specifics and raw log

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
func (it *GroupMessageBroadcasterPayloadBootstrapperUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessageBroadcasterPayloadBootstrapperUpdated)
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
		it.Event = new(GroupMessageBroadcasterPayloadBootstrapperUpdated)
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
func (it *GroupMessageBroadcasterPayloadBootstrapperUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessageBroadcasterPayloadBootstrapperUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessageBroadcasterPayloadBootstrapperUpdated represents a PayloadBootstrapperUpdated event raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterPayloadBootstrapperUpdated struct {
	PayloadBootstrapper common.Address
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterPayloadBootstrapperUpdated is a free log retrieval operation binding the contract event 0x38ecae7c300c129540d5181b5e16ec68d73e388d6add9ad70e63307f6794e2a9.
//
// Solidity: event PayloadBootstrapperUpdated(address indexed payloadBootstrapper)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) FilterPayloadBootstrapperUpdated(opts *bind.FilterOpts, payloadBootstrapper []common.Address) (*GroupMessageBroadcasterPayloadBootstrapperUpdatedIterator, error) {

	var payloadBootstrapperRule []interface{}
	for _, payloadBootstrapperItem := range payloadBootstrapper {
		payloadBootstrapperRule = append(payloadBootstrapperRule, payloadBootstrapperItem)
	}

	logs, sub, err := _GroupMessageBroadcaster.contract.FilterLogs(opts, "PayloadBootstrapperUpdated", payloadBootstrapperRule)
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterPayloadBootstrapperUpdatedIterator{contract: _GroupMessageBroadcaster.contract, event: "PayloadBootstrapperUpdated", logs: logs, sub: sub}, nil
}

// WatchPayloadBootstrapperUpdated is a free log subscription operation binding the contract event 0x38ecae7c300c129540d5181b5e16ec68d73e388d6add9ad70e63307f6794e2a9.
//
// Solidity: event PayloadBootstrapperUpdated(address indexed payloadBootstrapper)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) WatchPayloadBootstrapperUpdated(opts *bind.WatchOpts, sink chan<- *GroupMessageBroadcasterPayloadBootstrapperUpdated, payloadBootstrapper []common.Address) (event.Subscription, error) {

	var payloadBootstrapperRule []interface{}
	for _, payloadBootstrapperItem := range payloadBootstrapper {
		payloadBootstrapperRule = append(payloadBootstrapperRule, payloadBootstrapperItem)
	}

	logs, sub, err := _GroupMessageBroadcaster.contract.WatchLogs(opts, "PayloadBootstrapperUpdated", payloadBootstrapperRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessageBroadcasterPayloadBootstrapperUpdated)
				if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "PayloadBootstrapperUpdated", log); err != nil {
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

// ParsePayloadBootstrapperUpdated is a log parse operation binding the contract event 0x38ecae7c300c129540d5181b5e16ec68d73e388d6add9ad70e63307f6794e2a9.
//
// Solidity: event PayloadBootstrapperUpdated(address indexed payloadBootstrapper)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) ParsePayloadBootstrapperUpdated(log types.Log) (*GroupMessageBroadcasterPayloadBootstrapperUpdated, error) {
	event := new(GroupMessageBroadcasterPayloadBootstrapperUpdated)
	if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "PayloadBootstrapperUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessageBroadcasterUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterUpgradedIterator struct {
	Event *GroupMessageBroadcasterUpgraded // Event containing the contract specifics and raw log

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
func (it *GroupMessageBroadcasterUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessageBroadcasterUpgraded)
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
		it.Event = new(GroupMessageBroadcasterUpgraded)
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
func (it *GroupMessageBroadcasterUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessageBroadcasterUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessageBroadcasterUpgraded represents a Upgraded event raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*GroupMessageBroadcasterUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _GroupMessageBroadcaster.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterUpgradedIterator{contract: _GroupMessageBroadcaster.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *GroupMessageBroadcasterUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _GroupMessageBroadcaster.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessageBroadcasterUpgraded)
				if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) ParseUpgraded(log types.Log) (*GroupMessageBroadcasterUpgraded, error) {
	event := new(GroupMessageBroadcasterUpgraded)
	if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
