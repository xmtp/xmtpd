// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package depositsplitter

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

// DepositSplitterMetaData contains all meta data concerning the DepositSplitter contract.
var DepositSplitterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"feeToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"payerRegistry_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"settlementChainGateway_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"appChainId_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"appChainId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"payer_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"payerRegistryAmount_\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"appChainRecipient_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"appChainAmount_\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"appChainGasLimit_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"appChainMaxFeePerGas_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositFromUnderlying\",\"inputs\":[{\"name\":\"payer_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"payerRegistryAmount_\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"appChainRecipient_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"appChainAmount_\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"appChainGasLimit_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"appChainMaxFeePerGas_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositFromUnderlyingWithPermit\",\"inputs\":[{\"name\":\"payer_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"payerRegistryAmount_\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"appChainRecipient_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"appChainAmount_\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"appChainGasLimit_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"appChainMaxFeePerGas_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v_\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositWithPermit\",\"inputs\":[{\"name\":\"payer_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"payerRegistryAmount_\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"appChainRecipient_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"appChainAmount_\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"appChainGasLimit_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"appChainMaxFeePerGas_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v_\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"feeToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"payerRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"settlementChainGateway\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"TransferFromFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAppChainId\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroFeeToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroPayerRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroSettlementChainGateway\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroTotalAmount\",\"inputs\":[]}]",
}

// DepositSplitterABI is the input ABI used to generate the binding from.
// Deprecated: Use DepositSplitterMetaData.ABI instead.
var DepositSplitterABI = DepositSplitterMetaData.ABI

// DepositSplitter is an auto generated Go binding around an Ethereum contract.
type DepositSplitter struct {
	DepositSplitterCaller     // Read-only binding to the contract
	DepositSplitterTransactor // Write-only binding to the contract
	DepositSplitterFilterer   // Log filterer for contract events
}

// DepositSplitterCaller is an auto generated read-only Go binding around an Ethereum contract.
type DepositSplitterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositSplitterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DepositSplitterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositSplitterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DepositSplitterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositSplitterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DepositSplitterSession struct {
	Contract     *DepositSplitter  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DepositSplitterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DepositSplitterCallerSession struct {
	Contract *DepositSplitterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// DepositSplitterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DepositSplitterTransactorSession struct {
	Contract     *DepositSplitterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// DepositSplitterRaw is an auto generated low-level Go binding around an Ethereum contract.
type DepositSplitterRaw struct {
	Contract *DepositSplitter // Generic contract binding to access the raw methods on
}

// DepositSplitterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DepositSplitterCallerRaw struct {
	Contract *DepositSplitterCaller // Generic read-only contract binding to access the raw methods on
}

// DepositSplitterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DepositSplitterTransactorRaw struct {
	Contract *DepositSplitterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDepositSplitter creates a new instance of DepositSplitter, bound to a specific deployed contract.
func NewDepositSplitter(address common.Address, backend bind.ContractBackend) (*DepositSplitter, error) {
	contract, err := bindDepositSplitter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DepositSplitter{DepositSplitterCaller: DepositSplitterCaller{contract: contract}, DepositSplitterTransactor: DepositSplitterTransactor{contract: contract}, DepositSplitterFilterer: DepositSplitterFilterer{contract: contract}}, nil
}

// NewDepositSplitterCaller creates a new read-only instance of DepositSplitter, bound to a specific deployed contract.
func NewDepositSplitterCaller(address common.Address, caller bind.ContractCaller) (*DepositSplitterCaller, error) {
	contract, err := bindDepositSplitter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DepositSplitterCaller{contract: contract}, nil
}

// NewDepositSplitterTransactor creates a new write-only instance of DepositSplitter, bound to a specific deployed contract.
func NewDepositSplitterTransactor(address common.Address, transactor bind.ContractTransactor) (*DepositSplitterTransactor, error) {
	contract, err := bindDepositSplitter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DepositSplitterTransactor{contract: contract}, nil
}

// NewDepositSplitterFilterer creates a new log filterer instance of DepositSplitter, bound to a specific deployed contract.
func NewDepositSplitterFilterer(address common.Address, filterer bind.ContractFilterer) (*DepositSplitterFilterer, error) {
	contract, err := bindDepositSplitter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DepositSplitterFilterer{contract: contract}, nil
}

// bindDepositSplitter binds a generic wrapper to an already deployed contract.
func bindDepositSplitter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DepositSplitterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DepositSplitter *DepositSplitterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DepositSplitter.Contract.DepositSplitterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DepositSplitter *DepositSplitterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositSplitter.Contract.DepositSplitterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DepositSplitter *DepositSplitterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DepositSplitter.Contract.DepositSplitterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DepositSplitter *DepositSplitterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DepositSplitter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DepositSplitter *DepositSplitterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositSplitter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DepositSplitter *DepositSplitterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DepositSplitter.Contract.contract.Transact(opts, method, params...)
}

// AppChainId is a free data retrieval call binding the contract method 0x83470923.
//
// Solidity: function appChainId() view returns(uint256)
func (_DepositSplitter *DepositSplitterCaller) AppChainId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DepositSplitter.contract.Call(opts, &out, "appChainId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AppChainId is a free data retrieval call binding the contract method 0x83470923.
//
// Solidity: function appChainId() view returns(uint256)
func (_DepositSplitter *DepositSplitterSession) AppChainId() (*big.Int, error) {
	return _DepositSplitter.Contract.AppChainId(&_DepositSplitter.CallOpts)
}

// AppChainId is a free data retrieval call binding the contract method 0x83470923.
//
// Solidity: function appChainId() view returns(uint256)
func (_DepositSplitter *DepositSplitterCallerSession) AppChainId() (*big.Int, error) {
	return _DepositSplitter.Contract.AppChainId(&_DepositSplitter.CallOpts)
}

// FeeToken is a free data retrieval call binding the contract method 0x647846a5.
//
// Solidity: function feeToken() view returns(address)
func (_DepositSplitter *DepositSplitterCaller) FeeToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DepositSplitter.contract.Call(opts, &out, "feeToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeToken is a free data retrieval call binding the contract method 0x647846a5.
//
// Solidity: function feeToken() view returns(address)
func (_DepositSplitter *DepositSplitterSession) FeeToken() (common.Address, error) {
	return _DepositSplitter.Contract.FeeToken(&_DepositSplitter.CallOpts)
}

// FeeToken is a free data retrieval call binding the contract method 0x647846a5.
//
// Solidity: function feeToken() view returns(address)
func (_DepositSplitter *DepositSplitterCallerSession) FeeToken() (common.Address, error) {
	return _DepositSplitter.Contract.FeeToken(&_DepositSplitter.CallOpts)
}

// PayerRegistry is a free data retrieval call binding the contract method 0x1dc5f4b8.
//
// Solidity: function payerRegistry() view returns(address)
func (_DepositSplitter *DepositSplitterCaller) PayerRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DepositSplitter.contract.Call(opts, &out, "payerRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PayerRegistry is a free data retrieval call binding the contract method 0x1dc5f4b8.
//
// Solidity: function payerRegistry() view returns(address)
func (_DepositSplitter *DepositSplitterSession) PayerRegistry() (common.Address, error) {
	return _DepositSplitter.Contract.PayerRegistry(&_DepositSplitter.CallOpts)
}

// PayerRegistry is a free data retrieval call binding the contract method 0x1dc5f4b8.
//
// Solidity: function payerRegistry() view returns(address)
func (_DepositSplitter *DepositSplitterCallerSession) PayerRegistry() (common.Address, error) {
	return _DepositSplitter.Contract.PayerRegistry(&_DepositSplitter.CallOpts)
}

// SettlementChainGateway is a free data retrieval call binding the contract method 0x801fd7f3.
//
// Solidity: function settlementChainGateway() view returns(address)
func (_DepositSplitter *DepositSplitterCaller) SettlementChainGateway(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DepositSplitter.contract.Call(opts, &out, "settlementChainGateway")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SettlementChainGateway is a free data retrieval call binding the contract method 0x801fd7f3.
//
// Solidity: function settlementChainGateway() view returns(address)
func (_DepositSplitter *DepositSplitterSession) SettlementChainGateway() (common.Address, error) {
	return _DepositSplitter.Contract.SettlementChainGateway(&_DepositSplitter.CallOpts)
}

// SettlementChainGateway is a free data retrieval call binding the contract method 0x801fd7f3.
//
// Solidity: function settlementChainGateway() view returns(address)
func (_DepositSplitter *DepositSplitterCallerSession) SettlementChainGateway() (common.Address, error) {
	return _DepositSplitter.Contract.SettlementChainGateway(&_DepositSplitter.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xc9679bd0.
//
// Solidity: function deposit(address payer_, uint96 payerRegistryAmount_, address appChainRecipient_, uint96 appChainAmount_, uint256 appChainGasLimit_, uint256 appChainMaxFeePerGas_) returns()
func (_DepositSplitter *DepositSplitterTransactor) Deposit(opts *bind.TransactOpts, payer_ common.Address, payerRegistryAmount_ *big.Int, appChainRecipient_ common.Address, appChainAmount_ *big.Int, appChainGasLimit_ *big.Int, appChainMaxFeePerGas_ *big.Int) (*types.Transaction, error) {
	return _DepositSplitter.contract.Transact(opts, "deposit", payer_, payerRegistryAmount_, appChainRecipient_, appChainAmount_, appChainGasLimit_, appChainMaxFeePerGas_)
}

// Deposit is a paid mutator transaction binding the contract method 0xc9679bd0.
//
// Solidity: function deposit(address payer_, uint96 payerRegistryAmount_, address appChainRecipient_, uint96 appChainAmount_, uint256 appChainGasLimit_, uint256 appChainMaxFeePerGas_) returns()
func (_DepositSplitter *DepositSplitterSession) Deposit(payer_ common.Address, payerRegistryAmount_ *big.Int, appChainRecipient_ common.Address, appChainAmount_ *big.Int, appChainGasLimit_ *big.Int, appChainMaxFeePerGas_ *big.Int) (*types.Transaction, error) {
	return _DepositSplitter.Contract.Deposit(&_DepositSplitter.TransactOpts, payer_, payerRegistryAmount_, appChainRecipient_, appChainAmount_, appChainGasLimit_, appChainMaxFeePerGas_)
}

// Deposit is a paid mutator transaction binding the contract method 0xc9679bd0.
//
// Solidity: function deposit(address payer_, uint96 payerRegistryAmount_, address appChainRecipient_, uint96 appChainAmount_, uint256 appChainGasLimit_, uint256 appChainMaxFeePerGas_) returns()
func (_DepositSplitter *DepositSplitterTransactorSession) Deposit(payer_ common.Address, payerRegistryAmount_ *big.Int, appChainRecipient_ common.Address, appChainAmount_ *big.Int, appChainGasLimit_ *big.Int, appChainMaxFeePerGas_ *big.Int) (*types.Transaction, error) {
	return _DepositSplitter.Contract.Deposit(&_DepositSplitter.TransactOpts, payer_, payerRegistryAmount_, appChainRecipient_, appChainAmount_, appChainGasLimit_, appChainMaxFeePerGas_)
}

// DepositFromUnderlying is a paid mutator transaction binding the contract method 0x889e0645.
//
// Solidity: function depositFromUnderlying(address payer_, uint96 payerRegistryAmount_, address appChainRecipient_, uint96 appChainAmount_, uint256 appChainGasLimit_, uint256 appChainMaxFeePerGas_) returns()
func (_DepositSplitter *DepositSplitterTransactor) DepositFromUnderlying(opts *bind.TransactOpts, payer_ common.Address, payerRegistryAmount_ *big.Int, appChainRecipient_ common.Address, appChainAmount_ *big.Int, appChainGasLimit_ *big.Int, appChainMaxFeePerGas_ *big.Int) (*types.Transaction, error) {
	return _DepositSplitter.contract.Transact(opts, "depositFromUnderlying", payer_, payerRegistryAmount_, appChainRecipient_, appChainAmount_, appChainGasLimit_, appChainMaxFeePerGas_)
}

// DepositFromUnderlying is a paid mutator transaction binding the contract method 0x889e0645.
//
// Solidity: function depositFromUnderlying(address payer_, uint96 payerRegistryAmount_, address appChainRecipient_, uint96 appChainAmount_, uint256 appChainGasLimit_, uint256 appChainMaxFeePerGas_) returns()
func (_DepositSplitter *DepositSplitterSession) DepositFromUnderlying(payer_ common.Address, payerRegistryAmount_ *big.Int, appChainRecipient_ common.Address, appChainAmount_ *big.Int, appChainGasLimit_ *big.Int, appChainMaxFeePerGas_ *big.Int) (*types.Transaction, error) {
	return _DepositSplitter.Contract.DepositFromUnderlying(&_DepositSplitter.TransactOpts, payer_, payerRegistryAmount_, appChainRecipient_, appChainAmount_, appChainGasLimit_, appChainMaxFeePerGas_)
}

// DepositFromUnderlying is a paid mutator transaction binding the contract method 0x889e0645.
//
// Solidity: function depositFromUnderlying(address payer_, uint96 payerRegistryAmount_, address appChainRecipient_, uint96 appChainAmount_, uint256 appChainGasLimit_, uint256 appChainMaxFeePerGas_) returns()
func (_DepositSplitter *DepositSplitterTransactorSession) DepositFromUnderlying(payer_ common.Address, payerRegistryAmount_ *big.Int, appChainRecipient_ common.Address, appChainAmount_ *big.Int, appChainGasLimit_ *big.Int, appChainMaxFeePerGas_ *big.Int) (*types.Transaction, error) {
	return _DepositSplitter.Contract.DepositFromUnderlying(&_DepositSplitter.TransactOpts, payer_, payerRegistryAmount_, appChainRecipient_, appChainAmount_, appChainGasLimit_, appChainMaxFeePerGas_)
}

// DepositFromUnderlyingWithPermit is a paid mutator transaction binding the contract method 0x7af40478.
//
// Solidity: function depositFromUnderlyingWithPermit(address payer_, uint96 payerRegistryAmount_, address appChainRecipient_, uint96 appChainAmount_, uint256 appChainGasLimit_, uint256 appChainMaxFeePerGas_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns()
func (_DepositSplitter *DepositSplitterTransactor) DepositFromUnderlyingWithPermit(opts *bind.TransactOpts, payer_ common.Address, payerRegistryAmount_ *big.Int, appChainRecipient_ common.Address, appChainAmount_ *big.Int, appChainGasLimit_ *big.Int, appChainMaxFeePerGas_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _DepositSplitter.contract.Transact(opts, "depositFromUnderlyingWithPermit", payer_, payerRegistryAmount_, appChainRecipient_, appChainAmount_, appChainGasLimit_, appChainMaxFeePerGas_, deadline_, v_, r_, s_)
}

// DepositFromUnderlyingWithPermit is a paid mutator transaction binding the contract method 0x7af40478.
//
// Solidity: function depositFromUnderlyingWithPermit(address payer_, uint96 payerRegistryAmount_, address appChainRecipient_, uint96 appChainAmount_, uint256 appChainGasLimit_, uint256 appChainMaxFeePerGas_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns()
func (_DepositSplitter *DepositSplitterSession) DepositFromUnderlyingWithPermit(payer_ common.Address, payerRegistryAmount_ *big.Int, appChainRecipient_ common.Address, appChainAmount_ *big.Int, appChainGasLimit_ *big.Int, appChainMaxFeePerGas_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _DepositSplitter.Contract.DepositFromUnderlyingWithPermit(&_DepositSplitter.TransactOpts, payer_, payerRegistryAmount_, appChainRecipient_, appChainAmount_, appChainGasLimit_, appChainMaxFeePerGas_, deadline_, v_, r_, s_)
}

// DepositFromUnderlyingWithPermit is a paid mutator transaction binding the contract method 0x7af40478.
//
// Solidity: function depositFromUnderlyingWithPermit(address payer_, uint96 payerRegistryAmount_, address appChainRecipient_, uint96 appChainAmount_, uint256 appChainGasLimit_, uint256 appChainMaxFeePerGas_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns()
func (_DepositSplitter *DepositSplitterTransactorSession) DepositFromUnderlyingWithPermit(payer_ common.Address, payerRegistryAmount_ *big.Int, appChainRecipient_ common.Address, appChainAmount_ *big.Int, appChainGasLimit_ *big.Int, appChainMaxFeePerGas_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _DepositSplitter.Contract.DepositFromUnderlyingWithPermit(&_DepositSplitter.TransactOpts, payer_, payerRegistryAmount_, appChainRecipient_, appChainAmount_, appChainGasLimit_, appChainMaxFeePerGas_, deadline_, v_, r_, s_)
}

// DepositWithPermit is a paid mutator transaction binding the contract method 0x038842b9.
//
// Solidity: function depositWithPermit(address payer_, uint96 payerRegistryAmount_, address appChainRecipient_, uint96 appChainAmount_, uint256 appChainGasLimit_, uint256 appChainMaxFeePerGas_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns()
func (_DepositSplitter *DepositSplitterTransactor) DepositWithPermit(opts *bind.TransactOpts, payer_ common.Address, payerRegistryAmount_ *big.Int, appChainRecipient_ common.Address, appChainAmount_ *big.Int, appChainGasLimit_ *big.Int, appChainMaxFeePerGas_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _DepositSplitter.contract.Transact(opts, "depositWithPermit", payer_, payerRegistryAmount_, appChainRecipient_, appChainAmount_, appChainGasLimit_, appChainMaxFeePerGas_, deadline_, v_, r_, s_)
}

// DepositWithPermit is a paid mutator transaction binding the contract method 0x038842b9.
//
// Solidity: function depositWithPermit(address payer_, uint96 payerRegistryAmount_, address appChainRecipient_, uint96 appChainAmount_, uint256 appChainGasLimit_, uint256 appChainMaxFeePerGas_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns()
func (_DepositSplitter *DepositSplitterSession) DepositWithPermit(payer_ common.Address, payerRegistryAmount_ *big.Int, appChainRecipient_ common.Address, appChainAmount_ *big.Int, appChainGasLimit_ *big.Int, appChainMaxFeePerGas_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _DepositSplitter.Contract.DepositWithPermit(&_DepositSplitter.TransactOpts, payer_, payerRegistryAmount_, appChainRecipient_, appChainAmount_, appChainGasLimit_, appChainMaxFeePerGas_, deadline_, v_, r_, s_)
}

// DepositWithPermit is a paid mutator transaction binding the contract method 0x038842b9.
//
// Solidity: function depositWithPermit(address payer_, uint96 payerRegistryAmount_, address appChainRecipient_, uint96 appChainAmount_, uint256 appChainGasLimit_, uint256 appChainMaxFeePerGas_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns()
func (_DepositSplitter *DepositSplitterTransactorSession) DepositWithPermit(payer_ common.Address, payerRegistryAmount_ *big.Int, appChainRecipient_ common.Address, appChainAmount_ *big.Int, appChainGasLimit_ *big.Int, appChainMaxFeePerGas_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _DepositSplitter.Contract.DepositWithPermit(&_DepositSplitter.TransactOpts, payer_, payerRegistryAmount_, appChainRecipient_, appChainAmount_, appChainGasLimit_, appChainMaxFeePerGas_, deadline_, v_, r_, s_)
}
