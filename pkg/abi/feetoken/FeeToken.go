// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package feetoken

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

// FeeTokenMetaData contains all meta data concerning the FeeToken contract.
var FeeTokenMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"parameterRegistry_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"underlying_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DOMAIN_SEPARATOR\",\"inputs\":[],\"outputs\":[{\"name\":\"domainSeparator_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PERMIT_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"contractName\",\"inputs\":[],\"outputs\":[{\"name\":\"contractName_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"decimals_\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"amount_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositFor\",\"inputs\":[{\"name\":\"recipient_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"success_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositForWithPermit\",\"inputs\":[{\"name\":\"recipient_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v_\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"success_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositWithPermit\",\"inputs\":[{\"name\":\"amount_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v_\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPermitDigest\",\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"nonce_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"digest_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"implementation_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migrate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migratorParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nonces\",\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"nonce_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"parameterRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"permit\",\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v_\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"underlying\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"version_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"amount_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawTo\",\"inputs\":[{\"name\":\"recipient_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"success_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Migrated\",\"inputs\":[{\"name\":\"migrator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC2612ExpiredSignature\",\"inputs\":[{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC2612InvalidSigner\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"EmptyCode\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidAccountNonce\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"currentNonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationFailed\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertData_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterOutOfTypeBounds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransferFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransferFromFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroMigrator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroParameterRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroRecipient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroUnderlying\",\"inputs\":[]}]",
}

// FeeTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use FeeTokenMetaData.ABI instead.
var FeeTokenABI = FeeTokenMetaData.ABI

// FeeToken is an auto generated Go binding around an Ethereum contract.
type FeeToken struct {
	FeeTokenCaller     // Read-only binding to the contract
	FeeTokenTransactor // Write-only binding to the contract
	FeeTokenFilterer   // Log filterer for contract events
}

// FeeTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type FeeTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FeeTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FeeTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FeeTokenSession struct {
	Contract     *FeeToken         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FeeTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FeeTokenCallerSession struct {
	Contract *FeeTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// FeeTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FeeTokenTransactorSession struct {
	Contract     *FeeTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// FeeTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type FeeTokenRaw struct {
	Contract *FeeToken // Generic contract binding to access the raw methods on
}

// FeeTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FeeTokenCallerRaw struct {
	Contract *FeeTokenCaller // Generic read-only contract binding to access the raw methods on
}

// FeeTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FeeTokenTransactorRaw struct {
	Contract *FeeTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFeeToken creates a new instance of FeeToken, bound to a specific deployed contract.
func NewFeeToken(address common.Address, backend bind.ContractBackend) (*FeeToken, error) {
	contract, err := bindFeeToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FeeToken{FeeTokenCaller: FeeTokenCaller{contract: contract}, FeeTokenTransactor: FeeTokenTransactor{contract: contract}, FeeTokenFilterer: FeeTokenFilterer{contract: contract}}, nil
}

// NewFeeTokenCaller creates a new read-only instance of FeeToken, bound to a specific deployed contract.
func NewFeeTokenCaller(address common.Address, caller bind.ContractCaller) (*FeeTokenCaller, error) {
	contract, err := bindFeeToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeeTokenCaller{contract: contract}, nil
}

// NewFeeTokenTransactor creates a new write-only instance of FeeToken, bound to a specific deployed contract.
func NewFeeTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*FeeTokenTransactor, error) {
	contract, err := bindFeeToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeeTokenTransactor{contract: contract}, nil
}

// NewFeeTokenFilterer creates a new log filterer instance of FeeToken, bound to a specific deployed contract.
func NewFeeTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*FeeTokenFilterer, error) {
	contract, err := bindFeeToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeeTokenFilterer{contract: contract}, nil
}

// bindFeeToken binds a generic wrapper to an already deployed contract.
func bindFeeToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FeeTokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeeToken *FeeTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeToken.Contract.FeeTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeeToken *FeeTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeToken.Contract.FeeTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeeToken *FeeTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeToken.Contract.FeeTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeeToken *FeeTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeeToken *FeeTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeeToken *FeeTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeToken.Contract.contract.Transact(opts, method, params...)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32 domainSeparator_)
func (_FeeToken *FeeTokenCaller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FeeToken.contract.Call(opts, &out, "DOMAIN_SEPARATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32 domainSeparator_)
func (_FeeToken *FeeTokenSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _FeeToken.Contract.DOMAINSEPARATOR(&_FeeToken.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32 domainSeparator_)
func (_FeeToken *FeeTokenCallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _FeeToken.Contract.DOMAINSEPARATOR(&_FeeToken.CallOpts)
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() view returns(bytes32)
func (_FeeToken *FeeTokenCaller) PERMITTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FeeToken.contract.Call(opts, &out, "PERMIT_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() view returns(bytes32)
func (_FeeToken *FeeTokenSession) PERMITTYPEHASH() ([32]byte, error) {
	return _FeeToken.Contract.PERMITTYPEHASH(&_FeeToken.CallOpts)
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() view returns(bytes32)
func (_FeeToken *FeeTokenCallerSession) PERMITTYPEHASH() ([32]byte, error) {
	return _FeeToken.Contract.PERMITTYPEHASH(&_FeeToken.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_FeeToken *FeeTokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FeeToken.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_FeeToken *FeeTokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _FeeToken.Contract.Allowance(&_FeeToken.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_FeeToken *FeeTokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _FeeToken.Contract.Allowance(&_FeeToken.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_FeeToken *FeeTokenCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FeeToken.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_FeeToken *FeeTokenSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _FeeToken.Contract.BalanceOf(&_FeeToken.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_FeeToken *FeeTokenCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _FeeToken.Contract.BalanceOf(&_FeeToken.CallOpts, account)
}

// ContractName is a free data retrieval call binding the contract method 0x75d0c0dc.
//
// Solidity: function contractName() pure returns(string contractName_)
func (_FeeToken *FeeTokenCaller) ContractName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FeeToken.contract.Call(opts, &out, "contractName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ContractName is a free data retrieval call binding the contract method 0x75d0c0dc.
//
// Solidity: function contractName() pure returns(string contractName_)
func (_FeeToken *FeeTokenSession) ContractName() (string, error) {
	return _FeeToken.Contract.ContractName(&_FeeToken.CallOpts)
}

// ContractName is a free data retrieval call binding the contract method 0x75d0c0dc.
//
// Solidity: function contractName() pure returns(string contractName_)
func (_FeeToken *FeeTokenCallerSession) ContractName() (string, error) {
	return _FeeToken.Contract.ContractName(&_FeeToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8 decimals_)
func (_FeeToken *FeeTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _FeeToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8 decimals_)
func (_FeeToken *FeeTokenSession) Decimals() (uint8, error) {
	return _FeeToken.Contract.Decimals(&_FeeToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8 decimals_)
func (_FeeToken *FeeTokenCallerSession) Decimals() (uint8, error) {
	return _FeeToken.Contract.Decimals(&_FeeToken.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_FeeToken *FeeTokenCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _FeeToken.contract.Call(opts, &out, "eip712Domain")

	outstruct := new(struct {
		Fields            [1]byte
		Name              string
		Version           string
		ChainId           *big.Int
		VerifyingContract common.Address
		Salt              [32]byte
		Extensions        []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_FeeToken *FeeTokenSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _FeeToken.Contract.Eip712Domain(&_FeeToken.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_FeeToken *FeeTokenCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _FeeToken.Contract.Eip712Domain(&_FeeToken.CallOpts)
}

// GetPermitDigest is a free data retrieval call binding the contract method 0x2c19e8b5.
//
// Solidity: function getPermitDigest(address owner_, address spender_, uint256 value_, uint256 nonce_, uint256 deadline_) view returns(bytes32 digest_)
func (_FeeToken *FeeTokenCaller) GetPermitDigest(opts *bind.CallOpts, owner_ common.Address, spender_ common.Address, value_ *big.Int, nonce_ *big.Int, deadline_ *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _FeeToken.contract.Call(opts, &out, "getPermitDigest", owner_, spender_, value_, nonce_, deadline_)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetPermitDigest is a free data retrieval call binding the contract method 0x2c19e8b5.
//
// Solidity: function getPermitDigest(address owner_, address spender_, uint256 value_, uint256 nonce_, uint256 deadline_) view returns(bytes32 digest_)
func (_FeeToken *FeeTokenSession) GetPermitDigest(owner_ common.Address, spender_ common.Address, value_ *big.Int, nonce_ *big.Int, deadline_ *big.Int) ([32]byte, error) {
	return _FeeToken.Contract.GetPermitDigest(&_FeeToken.CallOpts, owner_, spender_, value_, nonce_, deadline_)
}

// GetPermitDigest is a free data retrieval call binding the contract method 0x2c19e8b5.
//
// Solidity: function getPermitDigest(address owner_, address spender_, uint256 value_, uint256 nonce_, uint256 deadline_) view returns(bytes32 digest_)
func (_FeeToken *FeeTokenCallerSession) GetPermitDigest(owner_ common.Address, spender_ common.Address, value_ *big.Int, nonce_ *big.Int, deadline_ *big.Int) ([32]byte, error) {
	return _FeeToken.Contract.GetPermitDigest(&_FeeToken.CallOpts, owner_, spender_, value_, nonce_, deadline_)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_FeeToken *FeeTokenCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FeeToken.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_FeeToken *FeeTokenSession) Implementation() (common.Address, error) {
	return _FeeToken.Contract.Implementation(&_FeeToken.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_FeeToken *FeeTokenCallerSession) Implementation() (common.Address, error) {
	return _FeeToken.Contract.Implementation(&_FeeToken.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_FeeToken *FeeTokenCaller) MigratorParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FeeToken.contract.Call(opts, &out, "migratorParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_FeeToken *FeeTokenSession) MigratorParameterKey() (string, error) {
	return _FeeToken.Contract.MigratorParameterKey(&_FeeToken.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_FeeToken *FeeTokenCallerSession) MigratorParameterKey() (string, error) {
	return _FeeToken.Contract.MigratorParameterKey(&_FeeToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_FeeToken *FeeTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FeeToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_FeeToken *FeeTokenSession) Name() (string, error) {
	return _FeeToken.Contract.Name(&_FeeToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_FeeToken *FeeTokenCallerSession) Name() (string, error) {
	return _FeeToken.Contract.Name(&_FeeToken.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner_) view returns(uint256 nonce_)
func (_FeeToken *FeeTokenCaller) Nonces(opts *bind.CallOpts, owner_ common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FeeToken.contract.Call(opts, &out, "nonces", owner_)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner_) view returns(uint256 nonce_)
func (_FeeToken *FeeTokenSession) Nonces(owner_ common.Address) (*big.Int, error) {
	return _FeeToken.Contract.Nonces(&_FeeToken.CallOpts, owner_)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner_) view returns(uint256 nonce_)
func (_FeeToken *FeeTokenCallerSession) Nonces(owner_ common.Address) (*big.Int, error) {
	return _FeeToken.Contract.Nonces(&_FeeToken.CallOpts, owner_)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_FeeToken *FeeTokenCaller) ParameterRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FeeToken.contract.Call(opts, &out, "parameterRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_FeeToken *FeeTokenSession) ParameterRegistry() (common.Address, error) {
	return _FeeToken.Contract.ParameterRegistry(&_FeeToken.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_FeeToken *FeeTokenCallerSession) ParameterRegistry() (common.Address, error) {
	return _FeeToken.Contract.ParameterRegistry(&_FeeToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_FeeToken *FeeTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FeeToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_FeeToken *FeeTokenSession) Symbol() (string, error) {
	return _FeeToken.Contract.Symbol(&_FeeToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_FeeToken *FeeTokenCallerSession) Symbol() (string, error) {
	return _FeeToken.Contract.Symbol(&_FeeToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_FeeToken *FeeTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FeeToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_FeeToken *FeeTokenSession) TotalSupply() (*big.Int, error) {
	return _FeeToken.Contract.TotalSupply(&_FeeToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_FeeToken *FeeTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _FeeToken.Contract.TotalSupply(&_FeeToken.CallOpts)
}

// Underlying is a free data retrieval call binding the contract method 0x6f307dc3.
//
// Solidity: function underlying() view returns(address)
func (_FeeToken *FeeTokenCaller) Underlying(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FeeToken.contract.Call(opts, &out, "underlying")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Underlying is a free data retrieval call binding the contract method 0x6f307dc3.
//
// Solidity: function underlying() view returns(address)
func (_FeeToken *FeeTokenSession) Underlying() (common.Address, error) {
	return _FeeToken.Contract.Underlying(&_FeeToken.CallOpts)
}

// Underlying is a free data retrieval call binding the contract method 0x6f307dc3.
//
// Solidity: function underlying() view returns(address)
func (_FeeToken *FeeTokenCallerSession) Underlying() (common.Address, error) {
	return _FeeToken.Contract.Underlying(&_FeeToken.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(string version_)
func (_FeeToken *FeeTokenCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FeeToken.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(string version_)
func (_FeeToken *FeeTokenSession) Version() (string, error) {
	return _FeeToken.Contract.Version(&_FeeToken.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(string version_)
func (_FeeToken *FeeTokenCallerSession) Version() (string, error) {
	return _FeeToken.Contract.Version(&_FeeToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_FeeToken *FeeTokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _FeeToken.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_FeeToken *FeeTokenSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _FeeToken.Contract.Approve(&_FeeToken.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_FeeToken *FeeTokenTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _FeeToken.Contract.Approve(&_FeeToken.TransactOpts, spender, value)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount_) returns()
func (_FeeToken *FeeTokenTransactor) Deposit(opts *bind.TransactOpts, amount_ *big.Int) (*types.Transaction, error) {
	return _FeeToken.contract.Transact(opts, "deposit", amount_)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount_) returns()
func (_FeeToken *FeeTokenSession) Deposit(amount_ *big.Int) (*types.Transaction, error) {
	return _FeeToken.Contract.Deposit(&_FeeToken.TransactOpts, amount_)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount_) returns()
func (_FeeToken *FeeTokenTransactorSession) Deposit(amount_ *big.Int) (*types.Transaction, error) {
	return _FeeToken.Contract.Deposit(&_FeeToken.TransactOpts, amount_)
}

// DepositFor is a paid mutator transaction binding the contract method 0x2f4f21e2.
//
// Solidity: function depositFor(address recipient_, uint256 amount_) returns(bool success_)
func (_FeeToken *FeeTokenTransactor) DepositFor(opts *bind.TransactOpts, recipient_ common.Address, amount_ *big.Int) (*types.Transaction, error) {
	return _FeeToken.contract.Transact(opts, "depositFor", recipient_, amount_)
}

// DepositFor is a paid mutator transaction binding the contract method 0x2f4f21e2.
//
// Solidity: function depositFor(address recipient_, uint256 amount_) returns(bool success_)
func (_FeeToken *FeeTokenSession) DepositFor(recipient_ common.Address, amount_ *big.Int) (*types.Transaction, error) {
	return _FeeToken.Contract.DepositFor(&_FeeToken.TransactOpts, recipient_, amount_)
}

// DepositFor is a paid mutator transaction binding the contract method 0x2f4f21e2.
//
// Solidity: function depositFor(address recipient_, uint256 amount_) returns(bool success_)
func (_FeeToken *FeeTokenTransactorSession) DepositFor(recipient_ common.Address, amount_ *big.Int) (*types.Transaction, error) {
	return _FeeToken.Contract.DepositFor(&_FeeToken.TransactOpts, recipient_, amount_)
}

// DepositForWithPermit is a paid mutator transaction binding the contract method 0xa58f33d3.
//
// Solidity: function depositForWithPermit(address recipient_, uint256 amount_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns(bool success_)
func (_FeeToken *FeeTokenTransactor) DepositForWithPermit(opts *bind.TransactOpts, recipient_ common.Address, amount_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _FeeToken.contract.Transact(opts, "depositForWithPermit", recipient_, amount_, deadline_, v_, r_, s_)
}

// DepositForWithPermit is a paid mutator transaction binding the contract method 0xa58f33d3.
//
// Solidity: function depositForWithPermit(address recipient_, uint256 amount_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns(bool success_)
func (_FeeToken *FeeTokenSession) DepositForWithPermit(recipient_ common.Address, amount_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _FeeToken.Contract.DepositForWithPermit(&_FeeToken.TransactOpts, recipient_, amount_, deadline_, v_, r_, s_)
}

// DepositForWithPermit is a paid mutator transaction binding the contract method 0xa58f33d3.
//
// Solidity: function depositForWithPermit(address recipient_, uint256 amount_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns(bool success_)
func (_FeeToken *FeeTokenTransactorSession) DepositForWithPermit(recipient_ common.Address, amount_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _FeeToken.Contract.DepositForWithPermit(&_FeeToken.TransactOpts, recipient_, amount_, deadline_, v_, r_, s_)
}

// DepositWithPermit is a paid mutator transaction binding the contract method 0x4a970be7.
//
// Solidity: function depositWithPermit(uint256 amount_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns()
func (_FeeToken *FeeTokenTransactor) DepositWithPermit(opts *bind.TransactOpts, amount_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _FeeToken.contract.Transact(opts, "depositWithPermit", amount_, deadline_, v_, r_, s_)
}

// DepositWithPermit is a paid mutator transaction binding the contract method 0x4a970be7.
//
// Solidity: function depositWithPermit(uint256 amount_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns()
func (_FeeToken *FeeTokenSession) DepositWithPermit(amount_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _FeeToken.Contract.DepositWithPermit(&_FeeToken.TransactOpts, amount_, deadline_, v_, r_, s_)
}

// DepositWithPermit is a paid mutator transaction binding the contract method 0x4a970be7.
//
// Solidity: function depositWithPermit(uint256 amount_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns()
func (_FeeToken *FeeTokenTransactorSession) DepositWithPermit(amount_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _FeeToken.Contract.DepositWithPermit(&_FeeToken.TransactOpts, amount_, deadline_, v_, r_, s_)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_FeeToken *FeeTokenTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeToken.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_FeeToken *FeeTokenSession) Initialize() (*types.Transaction, error) {
	return _FeeToken.Contract.Initialize(&_FeeToken.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_FeeToken *FeeTokenTransactorSession) Initialize() (*types.Transaction, error) {
	return _FeeToken.Contract.Initialize(&_FeeToken.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_FeeToken *FeeTokenTransactor) Migrate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeToken.contract.Transact(opts, "migrate")
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_FeeToken *FeeTokenSession) Migrate() (*types.Transaction, error) {
	return _FeeToken.Contract.Migrate(&_FeeToken.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_FeeToken *FeeTokenTransactorSession) Migrate() (*types.Transaction, error) {
	return _FeeToken.Contract.Migrate(&_FeeToken.TransactOpts)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner_, address spender_, uint256 value_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns()
func (_FeeToken *FeeTokenTransactor) Permit(opts *bind.TransactOpts, owner_ common.Address, spender_ common.Address, value_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _FeeToken.contract.Transact(opts, "permit", owner_, spender_, value_, deadline_, v_, r_, s_)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner_, address spender_, uint256 value_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns()
func (_FeeToken *FeeTokenSession) Permit(owner_ common.Address, spender_ common.Address, value_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _FeeToken.Contract.Permit(&_FeeToken.TransactOpts, owner_, spender_, value_, deadline_, v_, r_, s_)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner_, address spender_, uint256 value_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns()
func (_FeeToken *FeeTokenTransactorSession) Permit(owner_ common.Address, spender_ common.Address, value_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _FeeToken.Contract.Permit(&_FeeToken.TransactOpts, owner_, spender_, value_, deadline_, v_, r_, s_)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_FeeToken *FeeTokenTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _FeeToken.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_FeeToken *FeeTokenSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _FeeToken.Contract.Transfer(&_FeeToken.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_FeeToken *FeeTokenTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _FeeToken.Contract.Transfer(&_FeeToken.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_FeeToken *FeeTokenTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _FeeToken.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_FeeToken *FeeTokenSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _FeeToken.Contract.TransferFrom(&_FeeToken.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_FeeToken *FeeTokenTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _FeeToken.Contract.TransferFrom(&_FeeToken.TransactOpts, from, to, value)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount_) returns()
func (_FeeToken *FeeTokenTransactor) Withdraw(opts *bind.TransactOpts, amount_ *big.Int) (*types.Transaction, error) {
	return _FeeToken.contract.Transact(opts, "withdraw", amount_)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount_) returns()
func (_FeeToken *FeeTokenSession) Withdraw(amount_ *big.Int) (*types.Transaction, error) {
	return _FeeToken.Contract.Withdraw(&_FeeToken.TransactOpts, amount_)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount_) returns()
func (_FeeToken *FeeTokenTransactorSession) Withdraw(amount_ *big.Int) (*types.Transaction, error) {
	return _FeeToken.Contract.Withdraw(&_FeeToken.TransactOpts, amount_)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x205c2878.
//
// Solidity: function withdrawTo(address recipient_, uint256 amount_) returns(bool success_)
func (_FeeToken *FeeTokenTransactor) WithdrawTo(opts *bind.TransactOpts, recipient_ common.Address, amount_ *big.Int) (*types.Transaction, error) {
	return _FeeToken.contract.Transact(opts, "withdrawTo", recipient_, amount_)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x205c2878.
//
// Solidity: function withdrawTo(address recipient_, uint256 amount_) returns(bool success_)
func (_FeeToken *FeeTokenSession) WithdrawTo(recipient_ common.Address, amount_ *big.Int) (*types.Transaction, error) {
	return _FeeToken.Contract.WithdrawTo(&_FeeToken.TransactOpts, recipient_, amount_)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x205c2878.
//
// Solidity: function withdrawTo(address recipient_, uint256 amount_) returns(bool success_)
func (_FeeToken *FeeTokenTransactorSession) WithdrawTo(recipient_ common.Address, amount_ *big.Int) (*types.Transaction, error) {
	return _FeeToken.Contract.WithdrawTo(&_FeeToken.TransactOpts, recipient_, amount_)
}

// FeeTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the FeeToken contract.
type FeeTokenApprovalIterator struct {
	Event *FeeTokenApproval // Event containing the contract specifics and raw log

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
func (it *FeeTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeTokenApproval)
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
		it.Event = new(FeeTokenApproval)
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
func (it *FeeTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeTokenApproval represents a Approval event raised by the FeeToken contract.
type FeeTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_FeeToken *FeeTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*FeeTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _FeeToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &FeeTokenApprovalIterator{contract: _FeeToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_FeeToken *FeeTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *FeeTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _FeeToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeTokenApproval)
				if err := _FeeToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_FeeToken *FeeTokenFilterer) ParseApproval(log types.Log) (*FeeTokenApproval, error) {
	event := new(FeeTokenApproval)
	if err := _FeeToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeeTokenEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the FeeToken contract.
type FeeTokenEIP712DomainChangedIterator struct {
	Event *FeeTokenEIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *FeeTokenEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeTokenEIP712DomainChanged)
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
		it.Event = new(FeeTokenEIP712DomainChanged)
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
func (it *FeeTokenEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeTokenEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeTokenEIP712DomainChanged represents a EIP712DomainChanged event raised by the FeeToken contract.
type FeeTokenEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_FeeToken *FeeTokenFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*FeeTokenEIP712DomainChangedIterator, error) {

	logs, sub, err := _FeeToken.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &FeeTokenEIP712DomainChangedIterator{contract: _FeeToken.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_FeeToken *FeeTokenFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *FeeTokenEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _FeeToken.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeTokenEIP712DomainChanged)
				if err := _FeeToken.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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

// ParseEIP712DomainChanged is a log parse operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_FeeToken *FeeTokenFilterer) ParseEIP712DomainChanged(log types.Log) (*FeeTokenEIP712DomainChanged, error) {
	event := new(FeeTokenEIP712DomainChanged)
	if err := _FeeToken.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeeTokenInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the FeeToken contract.
type FeeTokenInitializedIterator struct {
	Event *FeeTokenInitialized // Event containing the contract specifics and raw log

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
func (it *FeeTokenInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeTokenInitialized)
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
		it.Event = new(FeeTokenInitialized)
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
func (it *FeeTokenInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeTokenInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeTokenInitialized represents a Initialized event raised by the FeeToken contract.
type FeeTokenInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_FeeToken *FeeTokenFilterer) FilterInitialized(opts *bind.FilterOpts) (*FeeTokenInitializedIterator, error) {

	logs, sub, err := _FeeToken.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &FeeTokenInitializedIterator{contract: _FeeToken.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_FeeToken *FeeTokenFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *FeeTokenInitialized) (event.Subscription, error) {

	logs, sub, err := _FeeToken.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeTokenInitialized)
				if err := _FeeToken.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_FeeToken *FeeTokenFilterer) ParseInitialized(log types.Log) (*FeeTokenInitialized, error) {
	event := new(FeeTokenInitialized)
	if err := _FeeToken.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeeTokenMigratedIterator is returned from FilterMigrated and is used to iterate over the raw logs and unpacked data for Migrated events raised by the FeeToken contract.
type FeeTokenMigratedIterator struct {
	Event *FeeTokenMigrated // Event containing the contract specifics and raw log

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
func (it *FeeTokenMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeTokenMigrated)
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
		it.Event = new(FeeTokenMigrated)
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
func (it *FeeTokenMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeTokenMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeTokenMigrated represents a Migrated event raised by the FeeToken contract.
type FeeTokenMigrated struct {
	Migrator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMigrated is a free log retrieval operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_FeeToken *FeeTokenFilterer) FilterMigrated(opts *bind.FilterOpts, migrator []common.Address) (*FeeTokenMigratedIterator, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _FeeToken.contract.FilterLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return &FeeTokenMigratedIterator{contract: _FeeToken.contract, event: "Migrated", logs: logs, sub: sub}, nil
}

// WatchMigrated is a free log subscription operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_FeeToken *FeeTokenFilterer) WatchMigrated(opts *bind.WatchOpts, sink chan<- *FeeTokenMigrated, migrator []common.Address) (event.Subscription, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _FeeToken.contract.WatchLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeTokenMigrated)
				if err := _FeeToken.contract.UnpackLog(event, "Migrated", log); err != nil {
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
func (_FeeToken *FeeTokenFilterer) ParseMigrated(log types.Log) (*FeeTokenMigrated, error) {
	event := new(FeeTokenMigrated)
	if err := _FeeToken.contract.UnpackLog(event, "Migrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeeTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the FeeToken contract.
type FeeTokenTransferIterator struct {
	Event *FeeTokenTransfer // Event containing the contract specifics and raw log

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
func (it *FeeTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeTokenTransfer)
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
		it.Event = new(FeeTokenTransfer)
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
func (it *FeeTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeTokenTransfer represents a Transfer event raised by the FeeToken contract.
type FeeTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_FeeToken *FeeTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FeeTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeeToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FeeTokenTransferIterator{contract: _FeeToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_FeeToken *FeeTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *FeeTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeeToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeTokenTransfer)
				if err := _FeeToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_FeeToken *FeeTokenFilterer) ParseTransfer(log types.Log) (*FeeTokenTransfer, error) {
	event := new(FeeTokenTransfer)
	if err := _FeeToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeeTokenUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the FeeToken contract.
type FeeTokenUpgradedIterator struct {
	Event *FeeTokenUpgraded // Event containing the contract specifics and raw log

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
func (it *FeeTokenUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeTokenUpgraded)
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
		it.Event = new(FeeTokenUpgraded)
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
func (it *FeeTokenUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeTokenUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeTokenUpgraded represents a Upgraded event raised by the FeeToken contract.
type FeeTokenUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_FeeToken *FeeTokenFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*FeeTokenUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _FeeToken.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &FeeTokenUpgradedIterator{contract: _FeeToken.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_FeeToken *FeeTokenFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *FeeTokenUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _FeeToken.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeTokenUpgraded)
				if err := _FeeToken.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_FeeToken *FeeTokenFilterer) ParseUpgraded(log types.Log) (*FeeTokenUpgraded, error) {
	event := new(FeeTokenUpgraded)
	if err := _FeeToken.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
