// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package distributionmanager

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

// DistributionManagerMetaData contains all meta data concerning the DistributionManager contract.
var DistributionManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"parameterRegistry_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nodeRegistry_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"payerReportManager_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"payerRegistry_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeToken_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"areFeesClaimed\",\"inputs\":[{\"name\":\"nodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"originatorNodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"payerReportIndex_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"hasClaimed_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"areProtocolFeesClaimed\",\"inputs\":[{\"name\":\"originatorNodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"payerReportIndex_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"areClaimed_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claim\",\"inputs\":[{\"name\":\"nodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"originatorNodeIds_\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"payerReportIndices_\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"claimed_\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimProtocolFees\",\"inputs\":[{\"name\":\"originatorNodeIds_\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"payerReportIndices_\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"claimed_\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"feeToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOwedFees\",\"inputs\":[{\"name\":\"nodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"owedFees_\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"implementation_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migrate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migratorParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"nodeRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owedProtocolFees\",\"inputs\":[],\"outputs\":[{\"name\":\"owedProtocolFees_\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"parameterRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"paused_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pausedParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"payerRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"payerReportManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"protocolFeesRecipient\",\"inputs\":[],\"outputs\":[{\"name\":\"protocolFeesRecipient_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"protocolFeesRecipientParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"totalOwedFees\",\"inputs\":[],\"outputs\":[{\"name\":\"totalOwedFees_\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updatePauseStatus\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateProtocolFeesRecipient\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"nodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"recipient_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"withdrawn_\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawIntoUnderlying\",\"inputs\":[{\"name\":\"nodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"recipient_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"withdrawn_\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawProtocolFees\",\"inputs\":[],\"outputs\":[{\"name\":\"withdrawn_\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawProtocolFeesIntoUnderlying\",\"inputs\":[],\"outputs\":[{\"name\":\"withdrawn_\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Claim\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"originatorNodeId\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"payerReportIndex\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Migrated\",\"inputs\":[{\"name\":\"migrator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PauseStatusUpdated\",\"inputs\":[{\"name\":\"paused\",\"type\":\"bool\",\"indexed\":true,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProtocolFeesClaim\",\"inputs\":[{\"name\":\"originatorNodeId\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"payerReportIndex\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProtocolFeesRecipientUpdated\",\"inputs\":[{\"name\":\"protocolFeesRecipient\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProtocolFeesWithdrawal\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawal\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyClaimed\",\"inputs\":[{\"name\":\"originatorNodeId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"payerReportIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ArrayLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyCode\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationFailed\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertData_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"NoChange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoFeesOwed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInPayerReport\",\"inputs\":[{\"name\":\"originatorNodeId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"payerReportIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotNodeOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterOutOfTypeBounds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Paused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PayerReportNotSettled\",\"inputs\":[{\"name\":\"originatorNodeId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"payerReportIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ZeroAvailableBalance\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroFeeToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroMigrator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroNodeRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroParameterRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroPayerRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroPayerReportManager\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroProtocolFeeRecipient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroRecipient\",\"inputs\":[]}]",
}

// DistributionManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use DistributionManagerMetaData.ABI instead.
var DistributionManagerABI = DistributionManagerMetaData.ABI

// DistributionManager is an auto generated Go binding around an Ethereum contract.
type DistributionManager struct {
	DistributionManagerCaller     // Read-only binding to the contract
	DistributionManagerTransactor // Write-only binding to the contract
	DistributionManagerFilterer   // Log filterer for contract events
}

// DistributionManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type DistributionManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DistributionManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DistributionManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DistributionManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DistributionManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DistributionManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DistributionManagerSession struct {
	Contract     *DistributionManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// DistributionManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DistributionManagerCallerSession struct {
	Contract *DistributionManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// DistributionManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DistributionManagerTransactorSession struct {
	Contract     *DistributionManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// DistributionManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type DistributionManagerRaw struct {
	Contract *DistributionManager // Generic contract binding to access the raw methods on
}

// DistributionManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DistributionManagerCallerRaw struct {
	Contract *DistributionManagerCaller // Generic read-only contract binding to access the raw methods on
}

// DistributionManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DistributionManagerTransactorRaw struct {
	Contract *DistributionManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDistributionManager creates a new instance of DistributionManager, bound to a specific deployed contract.
func NewDistributionManager(address common.Address, backend bind.ContractBackend) (*DistributionManager, error) {
	contract, err := bindDistributionManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DistributionManager{DistributionManagerCaller: DistributionManagerCaller{contract: contract}, DistributionManagerTransactor: DistributionManagerTransactor{contract: contract}, DistributionManagerFilterer: DistributionManagerFilterer{contract: contract}}, nil
}

// NewDistributionManagerCaller creates a new read-only instance of DistributionManager, bound to a specific deployed contract.
func NewDistributionManagerCaller(address common.Address, caller bind.ContractCaller) (*DistributionManagerCaller, error) {
	contract, err := bindDistributionManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DistributionManagerCaller{contract: contract}, nil
}

// NewDistributionManagerTransactor creates a new write-only instance of DistributionManager, bound to a specific deployed contract.
func NewDistributionManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*DistributionManagerTransactor, error) {
	contract, err := bindDistributionManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DistributionManagerTransactor{contract: contract}, nil
}

// NewDistributionManagerFilterer creates a new log filterer instance of DistributionManager, bound to a specific deployed contract.
func NewDistributionManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*DistributionManagerFilterer, error) {
	contract, err := bindDistributionManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DistributionManagerFilterer{contract: contract}, nil
}

// bindDistributionManager binds a generic wrapper to an already deployed contract.
func bindDistributionManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DistributionManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DistributionManager *DistributionManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DistributionManager.Contract.DistributionManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DistributionManager *DistributionManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DistributionManager.Contract.DistributionManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DistributionManager *DistributionManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DistributionManager.Contract.DistributionManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DistributionManager *DistributionManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DistributionManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DistributionManager *DistributionManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DistributionManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DistributionManager *DistributionManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DistributionManager.Contract.contract.Transact(opts, method, params...)
}

// AreFeesClaimed is a free data retrieval call binding the contract method 0x1617d097.
//
// Solidity: function areFeesClaimed(uint32 nodeId_, uint32 originatorNodeId_, uint256 payerReportIndex_) view returns(bool hasClaimed_)
func (_DistributionManager *DistributionManagerCaller) AreFeesClaimed(opts *bind.CallOpts, nodeId_ uint32, originatorNodeId_ uint32, payerReportIndex_ *big.Int) (bool, error) {
	var out []interface{}
	err := _DistributionManager.contract.Call(opts, &out, "areFeesClaimed", nodeId_, originatorNodeId_, payerReportIndex_)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AreFeesClaimed is a free data retrieval call binding the contract method 0x1617d097.
//
// Solidity: function areFeesClaimed(uint32 nodeId_, uint32 originatorNodeId_, uint256 payerReportIndex_) view returns(bool hasClaimed_)
func (_DistributionManager *DistributionManagerSession) AreFeesClaimed(nodeId_ uint32, originatorNodeId_ uint32, payerReportIndex_ *big.Int) (bool, error) {
	return _DistributionManager.Contract.AreFeesClaimed(&_DistributionManager.CallOpts, nodeId_, originatorNodeId_, payerReportIndex_)
}

// AreFeesClaimed is a free data retrieval call binding the contract method 0x1617d097.
//
// Solidity: function areFeesClaimed(uint32 nodeId_, uint32 originatorNodeId_, uint256 payerReportIndex_) view returns(bool hasClaimed_)
func (_DistributionManager *DistributionManagerCallerSession) AreFeesClaimed(nodeId_ uint32, originatorNodeId_ uint32, payerReportIndex_ *big.Int) (bool, error) {
	return _DistributionManager.Contract.AreFeesClaimed(&_DistributionManager.CallOpts, nodeId_, originatorNodeId_, payerReportIndex_)
}

// AreProtocolFeesClaimed is a free data retrieval call binding the contract method 0xdc17fb75.
//
// Solidity: function areProtocolFeesClaimed(uint32 originatorNodeId_, uint256 payerReportIndex_) view returns(bool areClaimed_)
func (_DistributionManager *DistributionManagerCaller) AreProtocolFeesClaimed(opts *bind.CallOpts, originatorNodeId_ uint32, payerReportIndex_ *big.Int) (bool, error) {
	var out []interface{}
	err := _DistributionManager.contract.Call(opts, &out, "areProtocolFeesClaimed", originatorNodeId_, payerReportIndex_)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AreProtocolFeesClaimed is a free data retrieval call binding the contract method 0xdc17fb75.
//
// Solidity: function areProtocolFeesClaimed(uint32 originatorNodeId_, uint256 payerReportIndex_) view returns(bool areClaimed_)
func (_DistributionManager *DistributionManagerSession) AreProtocolFeesClaimed(originatorNodeId_ uint32, payerReportIndex_ *big.Int) (bool, error) {
	return _DistributionManager.Contract.AreProtocolFeesClaimed(&_DistributionManager.CallOpts, originatorNodeId_, payerReportIndex_)
}

// AreProtocolFeesClaimed is a free data retrieval call binding the contract method 0xdc17fb75.
//
// Solidity: function areProtocolFeesClaimed(uint32 originatorNodeId_, uint256 payerReportIndex_) view returns(bool areClaimed_)
func (_DistributionManager *DistributionManagerCallerSession) AreProtocolFeesClaimed(originatorNodeId_ uint32, payerReportIndex_ *big.Int) (bool, error) {
	return _DistributionManager.Contract.AreProtocolFeesClaimed(&_DistributionManager.CallOpts, originatorNodeId_, payerReportIndex_)
}

// FeeToken is a free data retrieval call binding the contract method 0x647846a5.
//
// Solidity: function feeToken() view returns(address)
func (_DistributionManager *DistributionManagerCaller) FeeToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DistributionManager.contract.Call(opts, &out, "feeToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeToken is a free data retrieval call binding the contract method 0x647846a5.
//
// Solidity: function feeToken() view returns(address)
func (_DistributionManager *DistributionManagerSession) FeeToken() (common.Address, error) {
	return _DistributionManager.Contract.FeeToken(&_DistributionManager.CallOpts)
}

// FeeToken is a free data retrieval call binding the contract method 0x647846a5.
//
// Solidity: function feeToken() view returns(address)
func (_DistributionManager *DistributionManagerCallerSession) FeeToken() (common.Address, error) {
	return _DistributionManager.Contract.FeeToken(&_DistributionManager.CallOpts)
}

// GetOwedFees is a free data retrieval call binding the contract method 0x5b8e9219.
//
// Solidity: function getOwedFees(uint32 nodeId_) view returns(uint96 owedFees_)
func (_DistributionManager *DistributionManagerCaller) GetOwedFees(opts *bind.CallOpts, nodeId_ uint32) (*big.Int, error) {
	var out []interface{}
	err := _DistributionManager.contract.Call(opts, &out, "getOwedFees", nodeId_)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetOwedFees is a free data retrieval call binding the contract method 0x5b8e9219.
//
// Solidity: function getOwedFees(uint32 nodeId_) view returns(uint96 owedFees_)
func (_DistributionManager *DistributionManagerSession) GetOwedFees(nodeId_ uint32) (*big.Int, error) {
	return _DistributionManager.Contract.GetOwedFees(&_DistributionManager.CallOpts, nodeId_)
}

// GetOwedFees is a free data retrieval call binding the contract method 0x5b8e9219.
//
// Solidity: function getOwedFees(uint32 nodeId_) view returns(uint96 owedFees_)
func (_DistributionManager *DistributionManagerCallerSession) GetOwedFees(nodeId_ uint32) (*big.Int, error) {
	return _DistributionManager.Contract.GetOwedFees(&_DistributionManager.CallOpts, nodeId_)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_DistributionManager *DistributionManagerCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DistributionManager.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_DistributionManager *DistributionManagerSession) Implementation() (common.Address, error) {
	return _DistributionManager.Contract.Implementation(&_DistributionManager.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_DistributionManager *DistributionManagerCallerSession) Implementation() (common.Address, error) {
	return _DistributionManager.Contract.Implementation(&_DistributionManager.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_DistributionManager *DistributionManagerCaller) MigratorParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _DistributionManager.contract.Call(opts, &out, "migratorParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_DistributionManager *DistributionManagerSession) MigratorParameterKey() (string, error) {
	return _DistributionManager.Contract.MigratorParameterKey(&_DistributionManager.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_DistributionManager *DistributionManagerCallerSession) MigratorParameterKey() (string, error) {
	return _DistributionManager.Contract.MigratorParameterKey(&_DistributionManager.CallOpts)
}

// NodeRegistry is a free data retrieval call binding the contract method 0xd9b5c4a5.
//
// Solidity: function nodeRegistry() view returns(address)
func (_DistributionManager *DistributionManagerCaller) NodeRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DistributionManager.contract.Call(opts, &out, "nodeRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NodeRegistry is a free data retrieval call binding the contract method 0xd9b5c4a5.
//
// Solidity: function nodeRegistry() view returns(address)
func (_DistributionManager *DistributionManagerSession) NodeRegistry() (common.Address, error) {
	return _DistributionManager.Contract.NodeRegistry(&_DistributionManager.CallOpts)
}

// NodeRegistry is a free data retrieval call binding the contract method 0xd9b5c4a5.
//
// Solidity: function nodeRegistry() view returns(address)
func (_DistributionManager *DistributionManagerCallerSession) NodeRegistry() (common.Address, error) {
	return _DistributionManager.Contract.NodeRegistry(&_DistributionManager.CallOpts)
}

// OwedProtocolFees is a free data retrieval call binding the contract method 0xbda4f788.
//
// Solidity: function owedProtocolFees() view returns(uint96 owedProtocolFees_)
func (_DistributionManager *DistributionManagerCaller) OwedProtocolFees(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DistributionManager.contract.Call(opts, &out, "owedProtocolFees")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OwedProtocolFees is a free data retrieval call binding the contract method 0xbda4f788.
//
// Solidity: function owedProtocolFees() view returns(uint96 owedProtocolFees_)
func (_DistributionManager *DistributionManagerSession) OwedProtocolFees() (*big.Int, error) {
	return _DistributionManager.Contract.OwedProtocolFees(&_DistributionManager.CallOpts)
}

// OwedProtocolFees is a free data retrieval call binding the contract method 0xbda4f788.
//
// Solidity: function owedProtocolFees() view returns(uint96 owedProtocolFees_)
func (_DistributionManager *DistributionManagerCallerSession) OwedProtocolFees() (*big.Int, error) {
	return _DistributionManager.Contract.OwedProtocolFees(&_DistributionManager.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_DistributionManager *DistributionManagerCaller) ParameterRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DistributionManager.contract.Call(opts, &out, "parameterRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_DistributionManager *DistributionManagerSession) ParameterRegistry() (common.Address, error) {
	return _DistributionManager.Contract.ParameterRegistry(&_DistributionManager.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_DistributionManager *DistributionManagerCallerSession) ParameterRegistry() (common.Address, error) {
	return _DistributionManager.Contract.ParameterRegistry(&_DistributionManager.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool paused_)
func (_DistributionManager *DistributionManagerCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _DistributionManager.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool paused_)
func (_DistributionManager *DistributionManagerSession) Paused() (bool, error) {
	return _DistributionManager.Contract.Paused(&_DistributionManager.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool paused_)
func (_DistributionManager *DistributionManagerCallerSession) Paused() (bool, error) {
	return _DistributionManager.Contract.Paused(&_DistributionManager.CallOpts)
}

// PausedParameterKey is a free data retrieval call binding the contract method 0xcc5999af.
//
// Solidity: function pausedParameterKey() pure returns(string key_)
func (_DistributionManager *DistributionManagerCaller) PausedParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _DistributionManager.contract.Call(opts, &out, "pausedParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// PausedParameterKey is a free data retrieval call binding the contract method 0xcc5999af.
//
// Solidity: function pausedParameterKey() pure returns(string key_)
func (_DistributionManager *DistributionManagerSession) PausedParameterKey() (string, error) {
	return _DistributionManager.Contract.PausedParameterKey(&_DistributionManager.CallOpts)
}

// PausedParameterKey is a free data retrieval call binding the contract method 0xcc5999af.
//
// Solidity: function pausedParameterKey() pure returns(string key_)
func (_DistributionManager *DistributionManagerCallerSession) PausedParameterKey() (string, error) {
	return _DistributionManager.Contract.PausedParameterKey(&_DistributionManager.CallOpts)
}

// PayerRegistry is a free data retrieval call binding the contract method 0x1dc5f4b8.
//
// Solidity: function payerRegistry() view returns(address)
func (_DistributionManager *DistributionManagerCaller) PayerRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DistributionManager.contract.Call(opts, &out, "payerRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PayerRegistry is a free data retrieval call binding the contract method 0x1dc5f4b8.
//
// Solidity: function payerRegistry() view returns(address)
func (_DistributionManager *DistributionManagerSession) PayerRegistry() (common.Address, error) {
	return _DistributionManager.Contract.PayerRegistry(&_DistributionManager.CallOpts)
}

// PayerRegistry is a free data retrieval call binding the contract method 0x1dc5f4b8.
//
// Solidity: function payerRegistry() view returns(address)
func (_DistributionManager *DistributionManagerCallerSession) PayerRegistry() (common.Address, error) {
	return _DistributionManager.Contract.PayerRegistry(&_DistributionManager.CallOpts)
}

// PayerReportManager is a free data retrieval call binding the contract method 0x5a9b918c.
//
// Solidity: function payerReportManager() view returns(address)
func (_DistributionManager *DistributionManagerCaller) PayerReportManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DistributionManager.contract.Call(opts, &out, "payerReportManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PayerReportManager is a free data retrieval call binding the contract method 0x5a9b918c.
//
// Solidity: function payerReportManager() view returns(address)
func (_DistributionManager *DistributionManagerSession) PayerReportManager() (common.Address, error) {
	return _DistributionManager.Contract.PayerReportManager(&_DistributionManager.CallOpts)
}

// PayerReportManager is a free data retrieval call binding the contract method 0x5a9b918c.
//
// Solidity: function payerReportManager() view returns(address)
func (_DistributionManager *DistributionManagerCallerSession) PayerReportManager() (common.Address, error) {
	return _DistributionManager.Contract.PayerReportManager(&_DistributionManager.CallOpts)
}

// ProtocolFeesRecipient is a free data retrieval call binding the contract method 0x68930637.
//
// Solidity: function protocolFeesRecipient() view returns(address protocolFeesRecipient_)
func (_DistributionManager *DistributionManagerCaller) ProtocolFeesRecipient(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DistributionManager.contract.Call(opts, &out, "protocolFeesRecipient")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProtocolFeesRecipient is a free data retrieval call binding the contract method 0x68930637.
//
// Solidity: function protocolFeesRecipient() view returns(address protocolFeesRecipient_)
func (_DistributionManager *DistributionManagerSession) ProtocolFeesRecipient() (common.Address, error) {
	return _DistributionManager.Contract.ProtocolFeesRecipient(&_DistributionManager.CallOpts)
}

// ProtocolFeesRecipient is a free data retrieval call binding the contract method 0x68930637.
//
// Solidity: function protocolFeesRecipient() view returns(address protocolFeesRecipient_)
func (_DistributionManager *DistributionManagerCallerSession) ProtocolFeesRecipient() (common.Address, error) {
	return _DistributionManager.Contract.ProtocolFeesRecipient(&_DistributionManager.CallOpts)
}

// ProtocolFeesRecipientParameterKey is a free data retrieval call binding the contract method 0xd106d5c1.
//
// Solidity: function protocolFeesRecipientParameterKey() pure returns(string key_)
func (_DistributionManager *DistributionManagerCaller) ProtocolFeesRecipientParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _DistributionManager.contract.Call(opts, &out, "protocolFeesRecipientParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ProtocolFeesRecipientParameterKey is a free data retrieval call binding the contract method 0xd106d5c1.
//
// Solidity: function protocolFeesRecipientParameterKey() pure returns(string key_)
func (_DistributionManager *DistributionManagerSession) ProtocolFeesRecipientParameterKey() (string, error) {
	return _DistributionManager.Contract.ProtocolFeesRecipientParameterKey(&_DistributionManager.CallOpts)
}

// ProtocolFeesRecipientParameterKey is a free data retrieval call binding the contract method 0xd106d5c1.
//
// Solidity: function protocolFeesRecipientParameterKey() pure returns(string key_)
func (_DistributionManager *DistributionManagerCallerSession) ProtocolFeesRecipientParameterKey() (string, error) {
	return _DistributionManager.Contract.ProtocolFeesRecipientParameterKey(&_DistributionManager.CallOpts)
}

// TotalOwedFees is a free data retrieval call binding the contract method 0x80dadbce.
//
// Solidity: function totalOwedFees() view returns(uint96 totalOwedFees_)
func (_DistributionManager *DistributionManagerCaller) TotalOwedFees(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DistributionManager.contract.Call(opts, &out, "totalOwedFees")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalOwedFees is a free data retrieval call binding the contract method 0x80dadbce.
//
// Solidity: function totalOwedFees() view returns(uint96 totalOwedFees_)
func (_DistributionManager *DistributionManagerSession) TotalOwedFees() (*big.Int, error) {
	return _DistributionManager.Contract.TotalOwedFees(&_DistributionManager.CallOpts)
}

// TotalOwedFees is a free data retrieval call binding the contract method 0x80dadbce.
//
// Solidity: function totalOwedFees() view returns(uint96 totalOwedFees_)
func (_DistributionManager *DistributionManagerCallerSession) TotalOwedFees() (*big.Int, error) {
	return _DistributionManager.Contract.TotalOwedFees(&_DistributionManager.CallOpts)
}

// Claim is a paid mutator transaction binding the contract method 0x0da80b8a.
//
// Solidity: function claim(uint32 nodeId_, uint32[] originatorNodeIds_, uint256[] payerReportIndices_) returns(uint96 claimed_)
func (_DistributionManager *DistributionManagerTransactor) Claim(opts *bind.TransactOpts, nodeId_ uint32, originatorNodeIds_ []uint32, payerReportIndices_ []*big.Int) (*types.Transaction, error) {
	return _DistributionManager.contract.Transact(opts, "claim", nodeId_, originatorNodeIds_, payerReportIndices_)
}

// Claim is a paid mutator transaction binding the contract method 0x0da80b8a.
//
// Solidity: function claim(uint32 nodeId_, uint32[] originatorNodeIds_, uint256[] payerReportIndices_) returns(uint96 claimed_)
func (_DistributionManager *DistributionManagerSession) Claim(nodeId_ uint32, originatorNodeIds_ []uint32, payerReportIndices_ []*big.Int) (*types.Transaction, error) {
	return _DistributionManager.Contract.Claim(&_DistributionManager.TransactOpts, nodeId_, originatorNodeIds_, payerReportIndices_)
}

// Claim is a paid mutator transaction binding the contract method 0x0da80b8a.
//
// Solidity: function claim(uint32 nodeId_, uint32[] originatorNodeIds_, uint256[] payerReportIndices_) returns(uint96 claimed_)
func (_DistributionManager *DistributionManagerTransactorSession) Claim(nodeId_ uint32, originatorNodeIds_ []uint32, payerReportIndices_ []*big.Int) (*types.Transaction, error) {
	return _DistributionManager.Contract.Claim(&_DistributionManager.TransactOpts, nodeId_, originatorNodeIds_, payerReportIndices_)
}

// ClaimProtocolFees is a paid mutator transaction binding the contract method 0x6d399ad4.
//
// Solidity: function claimProtocolFees(uint32[] originatorNodeIds_, uint256[] payerReportIndices_) returns(uint96 claimed_)
func (_DistributionManager *DistributionManagerTransactor) ClaimProtocolFees(opts *bind.TransactOpts, originatorNodeIds_ []uint32, payerReportIndices_ []*big.Int) (*types.Transaction, error) {
	return _DistributionManager.contract.Transact(opts, "claimProtocolFees", originatorNodeIds_, payerReportIndices_)
}

// ClaimProtocolFees is a paid mutator transaction binding the contract method 0x6d399ad4.
//
// Solidity: function claimProtocolFees(uint32[] originatorNodeIds_, uint256[] payerReportIndices_) returns(uint96 claimed_)
func (_DistributionManager *DistributionManagerSession) ClaimProtocolFees(originatorNodeIds_ []uint32, payerReportIndices_ []*big.Int) (*types.Transaction, error) {
	return _DistributionManager.Contract.ClaimProtocolFees(&_DistributionManager.TransactOpts, originatorNodeIds_, payerReportIndices_)
}

// ClaimProtocolFees is a paid mutator transaction binding the contract method 0x6d399ad4.
//
// Solidity: function claimProtocolFees(uint32[] originatorNodeIds_, uint256[] payerReportIndices_) returns(uint96 claimed_)
func (_DistributionManager *DistributionManagerTransactorSession) ClaimProtocolFees(originatorNodeIds_ []uint32, payerReportIndices_ []*big.Int) (*types.Transaction, error) {
	return _DistributionManager.Contract.ClaimProtocolFees(&_DistributionManager.TransactOpts, originatorNodeIds_, payerReportIndices_)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_DistributionManager *DistributionManagerTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DistributionManager.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_DistributionManager *DistributionManagerSession) Initialize() (*types.Transaction, error) {
	return _DistributionManager.Contract.Initialize(&_DistributionManager.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_DistributionManager *DistributionManagerTransactorSession) Initialize() (*types.Transaction, error) {
	return _DistributionManager.Contract.Initialize(&_DistributionManager.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_DistributionManager *DistributionManagerTransactor) Migrate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DistributionManager.contract.Transact(opts, "migrate")
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_DistributionManager *DistributionManagerSession) Migrate() (*types.Transaction, error) {
	return _DistributionManager.Contract.Migrate(&_DistributionManager.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_DistributionManager *DistributionManagerTransactorSession) Migrate() (*types.Transaction, error) {
	return _DistributionManager.Contract.Migrate(&_DistributionManager.TransactOpts)
}

// UpdatePauseStatus is a paid mutator transaction binding the contract method 0x59d4df41.
//
// Solidity: function updatePauseStatus() returns()
func (_DistributionManager *DistributionManagerTransactor) UpdatePauseStatus(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DistributionManager.contract.Transact(opts, "updatePauseStatus")
}

// UpdatePauseStatus is a paid mutator transaction binding the contract method 0x59d4df41.
//
// Solidity: function updatePauseStatus() returns()
func (_DistributionManager *DistributionManagerSession) UpdatePauseStatus() (*types.Transaction, error) {
	return _DistributionManager.Contract.UpdatePauseStatus(&_DistributionManager.TransactOpts)
}

// UpdatePauseStatus is a paid mutator transaction binding the contract method 0x59d4df41.
//
// Solidity: function updatePauseStatus() returns()
func (_DistributionManager *DistributionManagerTransactorSession) UpdatePauseStatus() (*types.Transaction, error) {
	return _DistributionManager.Contract.UpdatePauseStatus(&_DistributionManager.TransactOpts)
}

// UpdateProtocolFeesRecipient is a paid mutator transaction binding the contract method 0x73e1f6d4.
//
// Solidity: function updateProtocolFeesRecipient() returns()
func (_DistributionManager *DistributionManagerTransactor) UpdateProtocolFeesRecipient(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DistributionManager.contract.Transact(opts, "updateProtocolFeesRecipient")
}

// UpdateProtocolFeesRecipient is a paid mutator transaction binding the contract method 0x73e1f6d4.
//
// Solidity: function updateProtocolFeesRecipient() returns()
func (_DistributionManager *DistributionManagerSession) UpdateProtocolFeesRecipient() (*types.Transaction, error) {
	return _DistributionManager.Contract.UpdateProtocolFeesRecipient(&_DistributionManager.TransactOpts)
}

// UpdateProtocolFeesRecipient is a paid mutator transaction binding the contract method 0x73e1f6d4.
//
// Solidity: function updateProtocolFeesRecipient() returns()
func (_DistributionManager *DistributionManagerTransactorSession) UpdateProtocolFeesRecipient() (*types.Transaction, error) {
	return _DistributionManager.Contract.UpdateProtocolFeesRecipient(&_DistributionManager.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x8316e5ae.
//
// Solidity: function withdraw(uint32 nodeId_, address recipient_) returns(uint96 withdrawn_)
func (_DistributionManager *DistributionManagerTransactor) Withdraw(opts *bind.TransactOpts, nodeId_ uint32, recipient_ common.Address) (*types.Transaction, error) {
	return _DistributionManager.contract.Transact(opts, "withdraw", nodeId_, recipient_)
}

// Withdraw is a paid mutator transaction binding the contract method 0x8316e5ae.
//
// Solidity: function withdraw(uint32 nodeId_, address recipient_) returns(uint96 withdrawn_)
func (_DistributionManager *DistributionManagerSession) Withdraw(nodeId_ uint32, recipient_ common.Address) (*types.Transaction, error) {
	return _DistributionManager.Contract.Withdraw(&_DistributionManager.TransactOpts, nodeId_, recipient_)
}

// Withdraw is a paid mutator transaction binding the contract method 0x8316e5ae.
//
// Solidity: function withdraw(uint32 nodeId_, address recipient_) returns(uint96 withdrawn_)
func (_DistributionManager *DistributionManagerTransactorSession) Withdraw(nodeId_ uint32, recipient_ common.Address) (*types.Transaction, error) {
	return _DistributionManager.Contract.Withdraw(&_DistributionManager.TransactOpts, nodeId_, recipient_)
}

// WithdrawIntoUnderlying is a paid mutator transaction binding the contract method 0x21c5bf0e.
//
// Solidity: function withdrawIntoUnderlying(uint32 nodeId_, address recipient_) returns(uint96 withdrawn_)
func (_DistributionManager *DistributionManagerTransactor) WithdrawIntoUnderlying(opts *bind.TransactOpts, nodeId_ uint32, recipient_ common.Address) (*types.Transaction, error) {
	return _DistributionManager.contract.Transact(opts, "withdrawIntoUnderlying", nodeId_, recipient_)
}

// WithdrawIntoUnderlying is a paid mutator transaction binding the contract method 0x21c5bf0e.
//
// Solidity: function withdrawIntoUnderlying(uint32 nodeId_, address recipient_) returns(uint96 withdrawn_)
func (_DistributionManager *DistributionManagerSession) WithdrawIntoUnderlying(nodeId_ uint32, recipient_ common.Address) (*types.Transaction, error) {
	return _DistributionManager.Contract.WithdrawIntoUnderlying(&_DistributionManager.TransactOpts, nodeId_, recipient_)
}

// WithdrawIntoUnderlying is a paid mutator transaction binding the contract method 0x21c5bf0e.
//
// Solidity: function withdrawIntoUnderlying(uint32 nodeId_, address recipient_) returns(uint96 withdrawn_)
func (_DistributionManager *DistributionManagerTransactorSession) WithdrawIntoUnderlying(nodeId_ uint32, recipient_ common.Address) (*types.Transaction, error) {
	return _DistributionManager.Contract.WithdrawIntoUnderlying(&_DistributionManager.TransactOpts, nodeId_, recipient_)
}

// WithdrawProtocolFees is a paid mutator transaction binding the contract method 0x8795cccb.
//
// Solidity: function withdrawProtocolFees() returns(uint96 withdrawn_)
func (_DistributionManager *DistributionManagerTransactor) WithdrawProtocolFees(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DistributionManager.contract.Transact(opts, "withdrawProtocolFees")
}

// WithdrawProtocolFees is a paid mutator transaction binding the contract method 0x8795cccb.
//
// Solidity: function withdrawProtocolFees() returns(uint96 withdrawn_)
func (_DistributionManager *DistributionManagerSession) WithdrawProtocolFees() (*types.Transaction, error) {
	return _DistributionManager.Contract.WithdrawProtocolFees(&_DistributionManager.TransactOpts)
}

// WithdrawProtocolFees is a paid mutator transaction binding the contract method 0x8795cccb.
//
// Solidity: function withdrawProtocolFees() returns(uint96 withdrawn_)
func (_DistributionManager *DistributionManagerTransactorSession) WithdrawProtocolFees() (*types.Transaction, error) {
	return _DistributionManager.Contract.WithdrawProtocolFees(&_DistributionManager.TransactOpts)
}

// WithdrawProtocolFeesIntoUnderlying is a paid mutator transaction binding the contract method 0x7d776753.
//
// Solidity: function withdrawProtocolFeesIntoUnderlying() returns(uint96 withdrawn_)
func (_DistributionManager *DistributionManagerTransactor) WithdrawProtocolFeesIntoUnderlying(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DistributionManager.contract.Transact(opts, "withdrawProtocolFeesIntoUnderlying")
}

// WithdrawProtocolFeesIntoUnderlying is a paid mutator transaction binding the contract method 0x7d776753.
//
// Solidity: function withdrawProtocolFeesIntoUnderlying() returns(uint96 withdrawn_)
func (_DistributionManager *DistributionManagerSession) WithdrawProtocolFeesIntoUnderlying() (*types.Transaction, error) {
	return _DistributionManager.Contract.WithdrawProtocolFeesIntoUnderlying(&_DistributionManager.TransactOpts)
}

// WithdrawProtocolFeesIntoUnderlying is a paid mutator transaction binding the contract method 0x7d776753.
//
// Solidity: function withdrawProtocolFeesIntoUnderlying() returns(uint96 withdrawn_)
func (_DistributionManager *DistributionManagerTransactorSession) WithdrawProtocolFeesIntoUnderlying() (*types.Transaction, error) {
	return _DistributionManager.Contract.WithdrawProtocolFeesIntoUnderlying(&_DistributionManager.TransactOpts)
}

// DistributionManagerClaimIterator is returned from FilterClaim and is used to iterate over the raw logs and unpacked data for Claim events raised by the DistributionManager contract.
type DistributionManagerClaimIterator struct {
	Event *DistributionManagerClaim // Event containing the contract specifics and raw log

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
func (it *DistributionManagerClaimIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DistributionManagerClaim)
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
		it.Event = new(DistributionManagerClaim)
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
func (it *DistributionManagerClaimIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DistributionManagerClaimIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DistributionManagerClaim represents a Claim event raised by the DistributionManager contract.
type DistributionManagerClaim struct {
	NodeId           uint32
	OriginatorNodeId uint32
	PayerReportIndex *big.Int
	Amount           *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterClaim is a free log retrieval operation binding the contract event 0x1345eb089f453d59f04c34856473ba8e5890ea1fb18e6b0414af0527ff986a0d.
//
// Solidity: event Claim(uint32 indexed nodeId, uint32 indexed originatorNodeId, uint256 indexed payerReportIndex, uint96 amount)
func (_DistributionManager *DistributionManagerFilterer) FilterClaim(opts *bind.FilterOpts, nodeId []uint32, originatorNodeId []uint32, payerReportIndex []*big.Int) (*DistributionManagerClaimIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var originatorNodeIdRule []interface{}
	for _, originatorNodeIdItem := range originatorNodeId {
		originatorNodeIdRule = append(originatorNodeIdRule, originatorNodeIdItem)
	}
	var payerReportIndexRule []interface{}
	for _, payerReportIndexItem := range payerReportIndex {
		payerReportIndexRule = append(payerReportIndexRule, payerReportIndexItem)
	}

	logs, sub, err := _DistributionManager.contract.FilterLogs(opts, "Claim", nodeIdRule, originatorNodeIdRule, payerReportIndexRule)
	if err != nil {
		return nil, err
	}
	return &DistributionManagerClaimIterator{contract: _DistributionManager.contract, event: "Claim", logs: logs, sub: sub}, nil
}

// WatchClaim is a free log subscription operation binding the contract event 0x1345eb089f453d59f04c34856473ba8e5890ea1fb18e6b0414af0527ff986a0d.
//
// Solidity: event Claim(uint32 indexed nodeId, uint32 indexed originatorNodeId, uint256 indexed payerReportIndex, uint96 amount)
func (_DistributionManager *DistributionManagerFilterer) WatchClaim(opts *bind.WatchOpts, sink chan<- *DistributionManagerClaim, nodeId []uint32, originatorNodeId []uint32, payerReportIndex []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var originatorNodeIdRule []interface{}
	for _, originatorNodeIdItem := range originatorNodeId {
		originatorNodeIdRule = append(originatorNodeIdRule, originatorNodeIdItem)
	}
	var payerReportIndexRule []interface{}
	for _, payerReportIndexItem := range payerReportIndex {
		payerReportIndexRule = append(payerReportIndexRule, payerReportIndexItem)
	}

	logs, sub, err := _DistributionManager.contract.WatchLogs(opts, "Claim", nodeIdRule, originatorNodeIdRule, payerReportIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DistributionManagerClaim)
				if err := _DistributionManager.contract.UnpackLog(event, "Claim", log); err != nil {
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

// ParseClaim is a log parse operation binding the contract event 0x1345eb089f453d59f04c34856473ba8e5890ea1fb18e6b0414af0527ff986a0d.
//
// Solidity: event Claim(uint32 indexed nodeId, uint32 indexed originatorNodeId, uint256 indexed payerReportIndex, uint96 amount)
func (_DistributionManager *DistributionManagerFilterer) ParseClaim(log types.Log) (*DistributionManagerClaim, error) {
	event := new(DistributionManagerClaim)
	if err := _DistributionManager.contract.UnpackLog(event, "Claim", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DistributionManagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the DistributionManager contract.
type DistributionManagerInitializedIterator struct {
	Event *DistributionManagerInitialized // Event containing the contract specifics and raw log

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
func (it *DistributionManagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DistributionManagerInitialized)
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
		it.Event = new(DistributionManagerInitialized)
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
func (it *DistributionManagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DistributionManagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DistributionManagerInitialized represents a Initialized event raised by the DistributionManager contract.
type DistributionManagerInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_DistributionManager *DistributionManagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*DistributionManagerInitializedIterator, error) {

	logs, sub, err := _DistributionManager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &DistributionManagerInitializedIterator{contract: _DistributionManager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_DistributionManager *DistributionManagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *DistributionManagerInitialized) (event.Subscription, error) {

	logs, sub, err := _DistributionManager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DistributionManagerInitialized)
				if err := _DistributionManager.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_DistributionManager *DistributionManagerFilterer) ParseInitialized(log types.Log) (*DistributionManagerInitialized, error) {
	event := new(DistributionManagerInitialized)
	if err := _DistributionManager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DistributionManagerMigratedIterator is returned from FilterMigrated and is used to iterate over the raw logs and unpacked data for Migrated events raised by the DistributionManager contract.
type DistributionManagerMigratedIterator struct {
	Event *DistributionManagerMigrated // Event containing the contract specifics and raw log

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
func (it *DistributionManagerMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DistributionManagerMigrated)
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
		it.Event = new(DistributionManagerMigrated)
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
func (it *DistributionManagerMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DistributionManagerMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DistributionManagerMigrated represents a Migrated event raised by the DistributionManager contract.
type DistributionManagerMigrated struct {
	Migrator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMigrated is a free log retrieval operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_DistributionManager *DistributionManagerFilterer) FilterMigrated(opts *bind.FilterOpts, migrator []common.Address) (*DistributionManagerMigratedIterator, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _DistributionManager.contract.FilterLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return &DistributionManagerMigratedIterator{contract: _DistributionManager.contract, event: "Migrated", logs: logs, sub: sub}, nil
}

// WatchMigrated is a free log subscription operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_DistributionManager *DistributionManagerFilterer) WatchMigrated(opts *bind.WatchOpts, sink chan<- *DistributionManagerMigrated, migrator []common.Address) (event.Subscription, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _DistributionManager.contract.WatchLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DistributionManagerMigrated)
				if err := _DistributionManager.contract.UnpackLog(event, "Migrated", log); err != nil {
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
func (_DistributionManager *DistributionManagerFilterer) ParseMigrated(log types.Log) (*DistributionManagerMigrated, error) {
	event := new(DistributionManagerMigrated)
	if err := _DistributionManager.contract.UnpackLog(event, "Migrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DistributionManagerPauseStatusUpdatedIterator is returned from FilterPauseStatusUpdated and is used to iterate over the raw logs and unpacked data for PauseStatusUpdated events raised by the DistributionManager contract.
type DistributionManagerPauseStatusUpdatedIterator struct {
	Event *DistributionManagerPauseStatusUpdated // Event containing the contract specifics and raw log

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
func (it *DistributionManagerPauseStatusUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DistributionManagerPauseStatusUpdated)
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
		it.Event = new(DistributionManagerPauseStatusUpdated)
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
func (it *DistributionManagerPauseStatusUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DistributionManagerPauseStatusUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DistributionManagerPauseStatusUpdated represents a PauseStatusUpdated event raised by the DistributionManager contract.
type DistributionManagerPauseStatusUpdated struct {
	Paused bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPauseStatusUpdated is a free log retrieval operation binding the contract event 0x7c4d1fe30fdbfda9e9c4c43e759ef32e4db5128d4cb58ff3ae9583b89b6242a5.
//
// Solidity: event PauseStatusUpdated(bool indexed paused)
func (_DistributionManager *DistributionManagerFilterer) FilterPauseStatusUpdated(opts *bind.FilterOpts, paused []bool) (*DistributionManagerPauseStatusUpdatedIterator, error) {

	var pausedRule []interface{}
	for _, pausedItem := range paused {
		pausedRule = append(pausedRule, pausedItem)
	}

	logs, sub, err := _DistributionManager.contract.FilterLogs(opts, "PauseStatusUpdated", pausedRule)
	if err != nil {
		return nil, err
	}
	return &DistributionManagerPauseStatusUpdatedIterator{contract: _DistributionManager.contract, event: "PauseStatusUpdated", logs: logs, sub: sub}, nil
}

// WatchPauseStatusUpdated is a free log subscription operation binding the contract event 0x7c4d1fe30fdbfda9e9c4c43e759ef32e4db5128d4cb58ff3ae9583b89b6242a5.
//
// Solidity: event PauseStatusUpdated(bool indexed paused)
func (_DistributionManager *DistributionManagerFilterer) WatchPauseStatusUpdated(opts *bind.WatchOpts, sink chan<- *DistributionManagerPauseStatusUpdated, paused []bool) (event.Subscription, error) {

	var pausedRule []interface{}
	for _, pausedItem := range paused {
		pausedRule = append(pausedRule, pausedItem)
	}

	logs, sub, err := _DistributionManager.contract.WatchLogs(opts, "PauseStatusUpdated", pausedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DistributionManagerPauseStatusUpdated)
				if err := _DistributionManager.contract.UnpackLog(event, "PauseStatusUpdated", log); err != nil {
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
func (_DistributionManager *DistributionManagerFilterer) ParsePauseStatusUpdated(log types.Log) (*DistributionManagerPauseStatusUpdated, error) {
	event := new(DistributionManagerPauseStatusUpdated)
	if err := _DistributionManager.contract.UnpackLog(event, "PauseStatusUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DistributionManagerProtocolFeesClaimIterator is returned from FilterProtocolFeesClaim and is used to iterate over the raw logs and unpacked data for ProtocolFeesClaim events raised by the DistributionManager contract.
type DistributionManagerProtocolFeesClaimIterator struct {
	Event *DistributionManagerProtocolFeesClaim // Event containing the contract specifics and raw log

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
func (it *DistributionManagerProtocolFeesClaimIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DistributionManagerProtocolFeesClaim)
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
		it.Event = new(DistributionManagerProtocolFeesClaim)
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
func (it *DistributionManagerProtocolFeesClaimIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DistributionManagerProtocolFeesClaimIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DistributionManagerProtocolFeesClaim represents a ProtocolFeesClaim event raised by the DistributionManager contract.
type DistributionManagerProtocolFeesClaim struct {
	OriginatorNodeId uint32
	PayerReportIndex *big.Int
	Amount           *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterProtocolFeesClaim is a free log retrieval operation binding the contract event 0x4b13292cd0c7c3d0915e067520d2d9a9281868666614ff9ce7fe127e8d029a17.
//
// Solidity: event ProtocolFeesClaim(uint32 indexed originatorNodeId, uint256 indexed payerReportIndex, uint96 amount)
func (_DistributionManager *DistributionManagerFilterer) FilterProtocolFeesClaim(opts *bind.FilterOpts, originatorNodeId []uint32, payerReportIndex []*big.Int) (*DistributionManagerProtocolFeesClaimIterator, error) {

	var originatorNodeIdRule []interface{}
	for _, originatorNodeIdItem := range originatorNodeId {
		originatorNodeIdRule = append(originatorNodeIdRule, originatorNodeIdItem)
	}
	var payerReportIndexRule []interface{}
	for _, payerReportIndexItem := range payerReportIndex {
		payerReportIndexRule = append(payerReportIndexRule, payerReportIndexItem)
	}

	logs, sub, err := _DistributionManager.contract.FilterLogs(opts, "ProtocolFeesClaim", originatorNodeIdRule, payerReportIndexRule)
	if err != nil {
		return nil, err
	}
	return &DistributionManagerProtocolFeesClaimIterator{contract: _DistributionManager.contract, event: "ProtocolFeesClaim", logs: logs, sub: sub}, nil
}

// WatchProtocolFeesClaim is a free log subscription operation binding the contract event 0x4b13292cd0c7c3d0915e067520d2d9a9281868666614ff9ce7fe127e8d029a17.
//
// Solidity: event ProtocolFeesClaim(uint32 indexed originatorNodeId, uint256 indexed payerReportIndex, uint96 amount)
func (_DistributionManager *DistributionManagerFilterer) WatchProtocolFeesClaim(opts *bind.WatchOpts, sink chan<- *DistributionManagerProtocolFeesClaim, originatorNodeId []uint32, payerReportIndex []*big.Int) (event.Subscription, error) {

	var originatorNodeIdRule []interface{}
	for _, originatorNodeIdItem := range originatorNodeId {
		originatorNodeIdRule = append(originatorNodeIdRule, originatorNodeIdItem)
	}
	var payerReportIndexRule []interface{}
	for _, payerReportIndexItem := range payerReportIndex {
		payerReportIndexRule = append(payerReportIndexRule, payerReportIndexItem)
	}

	logs, sub, err := _DistributionManager.contract.WatchLogs(opts, "ProtocolFeesClaim", originatorNodeIdRule, payerReportIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DistributionManagerProtocolFeesClaim)
				if err := _DistributionManager.contract.UnpackLog(event, "ProtocolFeesClaim", log); err != nil {
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

// ParseProtocolFeesClaim is a log parse operation binding the contract event 0x4b13292cd0c7c3d0915e067520d2d9a9281868666614ff9ce7fe127e8d029a17.
//
// Solidity: event ProtocolFeesClaim(uint32 indexed originatorNodeId, uint256 indexed payerReportIndex, uint96 amount)
func (_DistributionManager *DistributionManagerFilterer) ParseProtocolFeesClaim(log types.Log) (*DistributionManagerProtocolFeesClaim, error) {
	event := new(DistributionManagerProtocolFeesClaim)
	if err := _DistributionManager.contract.UnpackLog(event, "ProtocolFeesClaim", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DistributionManagerProtocolFeesRecipientUpdatedIterator is returned from FilterProtocolFeesRecipientUpdated and is used to iterate over the raw logs and unpacked data for ProtocolFeesRecipientUpdated events raised by the DistributionManager contract.
type DistributionManagerProtocolFeesRecipientUpdatedIterator struct {
	Event *DistributionManagerProtocolFeesRecipientUpdated // Event containing the contract specifics and raw log

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
func (it *DistributionManagerProtocolFeesRecipientUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DistributionManagerProtocolFeesRecipientUpdated)
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
		it.Event = new(DistributionManagerProtocolFeesRecipientUpdated)
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
func (it *DistributionManagerProtocolFeesRecipientUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DistributionManagerProtocolFeesRecipientUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DistributionManagerProtocolFeesRecipientUpdated represents a ProtocolFeesRecipientUpdated event raised by the DistributionManager contract.
type DistributionManagerProtocolFeesRecipientUpdated struct {
	ProtocolFeesRecipient common.Address
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterProtocolFeesRecipientUpdated is a free log retrieval operation binding the contract event 0x946e24a22a02b6e595bf703f75e3d98a755497d17cdc57ba967888c80106ce1f.
//
// Solidity: event ProtocolFeesRecipientUpdated(address protocolFeesRecipient)
func (_DistributionManager *DistributionManagerFilterer) FilterProtocolFeesRecipientUpdated(opts *bind.FilterOpts) (*DistributionManagerProtocolFeesRecipientUpdatedIterator, error) {

	logs, sub, err := _DistributionManager.contract.FilterLogs(opts, "ProtocolFeesRecipientUpdated")
	if err != nil {
		return nil, err
	}
	return &DistributionManagerProtocolFeesRecipientUpdatedIterator{contract: _DistributionManager.contract, event: "ProtocolFeesRecipientUpdated", logs: logs, sub: sub}, nil
}

// WatchProtocolFeesRecipientUpdated is a free log subscription operation binding the contract event 0x946e24a22a02b6e595bf703f75e3d98a755497d17cdc57ba967888c80106ce1f.
//
// Solidity: event ProtocolFeesRecipientUpdated(address protocolFeesRecipient)
func (_DistributionManager *DistributionManagerFilterer) WatchProtocolFeesRecipientUpdated(opts *bind.WatchOpts, sink chan<- *DistributionManagerProtocolFeesRecipientUpdated) (event.Subscription, error) {

	logs, sub, err := _DistributionManager.contract.WatchLogs(opts, "ProtocolFeesRecipientUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DistributionManagerProtocolFeesRecipientUpdated)
				if err := _DistributionManager.contract.UnpackLog(event, "ProtocolFeesRecipientUpdated", log); err != nil {
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

// ParseProtocolFeesRecipientUpdated is a log parse operation binding the contract event 0x946e24a22a02b6e595bf703f75e3d98a755497d17cdc57ba967888c80106ce1f.
//
// Solidity: event ProtocolFeesRecipientUpdated(address protocolFeesRecipient)
func (_DistributionManager *DistributionManagerFilterer) ParseProtocolFeesRecipientUpdated(log types.Log) (*DistributionManagerProtocolFeesRecipientUpdated, error) {
	event := new(DistributionManagerProtocolFeesRecipientUpdated)
	if err := _DistributionManager.contract.UnpackLog(event, "ProtocolFeesRecipientUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DistributionManagerProtocolFeesWithdrawalIterator is returned from FilterProtocolFeesWithdrawal and is used to iterate over the raw logs and unpacked data for ProtocolFeesWithdrawal events raised by the DistributionManager contract.
type DistributionManagerProtocolFeesWithdrawalIterator struct {
	Event *DistributionManagerProtocolFeesWithdrawal // Event containing the contract specifics and raw log

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
func (it *DistributionManagerProtocolFeesWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DistributionManagerProtocolFeesWithdrawal)
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
		it.Event = new(DistributionManagerProtocolFeesWithdrawal)
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
func (it *DistributionManagerProtocolFeesWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DistributionManagerProtocolFeesWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DistributionManagerProtocolFeesWithdrawal represents a ProtocolFeesWithdrawal event raised by the DistributionManager contract.
type DistributionManagerProtocolFeesWithdrawal struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterProtocolFeesWithdrawal is a free log retrieval operation binding the contract event 0x0fea118a9a3a3c942b3df91cd1b2b781ecebf3583b59e40a7e51b9e7b95abb05.
//
// Solidity: event ProtocolFeesWithdrawal(uint96 amount)
func (_DistributionManager *DistributionManagerFilterer) FilterProtocolFeesWithdrawal(opts *bind.FilterOpts) (*DistributionManagerProtocolFeesWithdrawalIterator, error) {

	logs, sub, err := _DistributionManager.contract.FilterLogs(opts, "ProtocolFeesWithdrawal")
	if err != nil {
		return nil, err
	}
	return &DistributionManagerProtocolFeesWithdrawalIterator{contract: _DistributionManager.contract, event: "ProtocolFeesWithdrawal", logs: logs, sub: sub}, nil
}

// WatchProtocolFeesWithdrawal is a free log subscription operation binding the contract event 0x0fea118a9a3a3c942b3df91cd1b2b781ecebf3583b59e40a7e51b9e7b95abb05.
//
// Solidity: event ProtocolFeesWithdrawal(uint96 amount)
func (_DistributionManager *DistributionManagerFilterer) WatchProtocolFeesWithdrawal(opts *bind.WatchOpts, sink chan<- *DistributionManagerProtocolFeesWithdrawal) (event.Subscription, error) {

	logs, sub, err := _DistributionManager.contract.WatchLogs(opts, "ProtocolFeesWithdrawal")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DistributionManagerProtocolFeesWithdrawal)
				if err := _DistributionManager.contract.UnpackLog(event, "ProtocolFeesWithdrawal", log); err != nil {
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

// ParseProtocolFeesWithdrawal is a log parse operation binding the contract event 0x0fea118a9a3a3c942b3df91cd1b2b781ecebf3583b59e40a7e51b9e7b95abb05.
//
// Solidity: event ProtocolFeesWithdrawal(uint96 amount)
func (_DistributionManager *DistributionManagerFilterer) ParseProtocolFeesWithdrawal(log types.Log) (*DistributionManagerProtocolFeesWithdrawal, error) {
	event := new(DistributionManagerProtocolFeesWithdrawal)
	if err := _DistributionManager.contract.UnpackLog(event, "ProtocolFeesWithdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DistributionManagerUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the DistributionManager contract.
type DistributionManagerUpgradedIterator struct {
	Event *DistributionManagerUpgraded // Event containing the contract specifics and raw log

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
func (it *DistributionManagerUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DistributionManagerUpgraded)
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
		it.Event = new(DistributionManagerUpgraded)
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
func (it *DistributionManagerUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DistributionManagerUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DistributionManagerUpgraded represents a Upgraded event raised by the DistributionManager contract.
type DistributionManagerUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_DistributionManager *DistributionManagerFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*DistributionManagerUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _DistributionManager.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &DistributionManagerUpgradedIterator{contract: _DistributionManager.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_DistributionManager *DistributionManagerFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *DistributionManagerUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _DistributionManager.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DistributionManagerUpgraded)
				if err := _DistributionManager.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_DistributionManager *DistributionManagerFilterer) ParseUpgraded(log types.Log) (*DistributionManagerUpgraded, error) {
	event := new(DistributionManagerUpgraded)
	if err := _DistributionManager.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DistributionManagerWithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the DistributionManager contract.
type DistributionManagerWithdrawalIterator struct {
	Event *DistributionManagerWithdrawal // Event containing the contract specifics and raw log

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
func (it *DistributionManagerWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DistributionManagerWithdrawal)
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
		it.Event = new(DistributionManagerWithdrawal)
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
func (it *DistributionManagerWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DistributionManagerWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DistributionManagerWithdrawal represents a Withdrawal event raised by the DistributionManager contract.
type DistributionManagerWithdrawal struct {
	NodeId uint32
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0x26bb85bbef39e470a6e9ff0c8c5b3e74ae1e80bc66f4be2072177949295e656f.
//
// Solidity: event Withdrawal(uint32 indexed nodeId, uint96 amount)
func (_DistributionManager *DistributionManagerFilterer) FilterWithdrawal(opts *bind.FilterOpts, nodeId []uint32) (*DistributionManagerWithdrawalIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _DistributionManager.contract.FilterLogs(opts, "Withdrawal", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &DistributionManagerWithdrawalIterator{contract: _DistributionManager.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0x26bb85bbef39e470a6e9ff0c8c5b3e74ae1e80bc66f4be2072177949295e656f.
//
// Solidity: event Withdrawal(uint32 indexed nodeId, uint96 amount)
func (_DistributionManager *DistributionManagerFilterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *DistributionManagerWithdrawal, nodeId []uint32) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _DistributionManager.contract.WatchLogs(opts, "Withdrawal", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DistributionManagerWithdrawal)
				if err := _DistributionManager.contract.UnpackLog(event, "Withdrawal", log); err != nil {
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

// ParseWithdrawal is a log parse operation binding the contract event 0x26bb85bbef39e470a6e9ff0c8c5b3e74ae1e80bc66f4be2072177949295e656f.
//
// Solidity: event Withdrawal(uint32 indexed nodeId, uint96 amount)
func (_DistributionManager *DistributionManagerFilterer) ParseWithdrawal(log types.Log) (*DistributionManagerWithdrawal, error) {
	event := new(DistributionManagerWithdrawal)
	if err := _DistributionManager.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
