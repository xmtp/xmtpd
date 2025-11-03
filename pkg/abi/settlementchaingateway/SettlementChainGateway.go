// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package settlementchaingateway

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

// SettlementChainGatewayMetaData contains all meta data concerning the SettlementChainGateway contract.
var SettlementChainGatewayMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"parameterRegistry_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"appChainGateway_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeToken_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"appChainAlias\",\"inputs\":[],\"outputs\":[{\"name\":\"alias_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"appChainGateway\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"calculateMaxDepositFee\",\"inputs\":[{\"name\":\"chainId_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasLimit_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxFeePerGas_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxBaseFee_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"fees_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"chainId_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"recipient_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasLimit_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxFeePerGas_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositFromUnderlying\",\"inputs\":[{\"name\":\"chainId_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"recipient_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasLimit_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxFeePerGas_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositFromUnderlyingWithPermit\",\"inputs\":[{\"name\":\"chainId_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"recipient_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasLimit_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxFeePerGas_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v_\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositWithPermit\",\"inputs\":[{\"name\":\"chainId_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"recipient_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasLimit_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxFeePerGas_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v_\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"feeToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getInbox\",\"inputs\":[{\"name\":\"chainId_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"inbox_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"implementation_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"inboxParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migrate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migratorParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"parameterRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"paused_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pausedParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"receiveWithdrawal\",\"inputs\":[{\"name\":\"recipient_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"amount_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"receiveWithdrawalIntoUnderlying\",\"inputs\":[{\"name\":\"recipient_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"amount_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendParameters\",\"inputs\":[{\"name\":\"chainIds_\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"keys_\",\"type\":\"string[]\",\"internalType\":\"string[]\"},{\"name\":\"gasLimit_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxFeePerGas_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountToSend_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"totalSent_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendParametersFromUnderlying\",\"inputs\":[{\"name\":\"chainIds_\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"keys_\",\"type\":\"string[]\",\"internalType\":\"string[]\"},{\"name\":\"gasLimit_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxFeePerGas_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountToSend_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"totalSent_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendParametersFromUnderlyingWithPermit\",\"inputs\":[{\"name\":\"chainIds_\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"keys_\",\"type\":\"string[]\",\"internalType\":\"string[]\"},{\"name\":\"gasLimit_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxFeePerGas_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountToSend_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v_\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"totalSent_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendParametersWithPermit\",\"inputs\":[{\"name\":\"chainIds_\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"keys_\",\"type\":\"string[]\",\"internalType\":\"string[]\"},{\"name\":\"gasLimit_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxFeePerGas_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountToSend_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v_\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"totalSent_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateInbox\",\"inputs\":[{\"name\":\"chainId_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updatePauseStatus\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"messageNumber\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"maxFees\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InboxUpdated\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"inbox\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Migrated\",\"inputs\":[{\"name\":\"migrator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ParametersSent\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"messageNumber\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"keys\",\"type\":\"string[]\",\"indexed\":false,\"internalType\":\"string[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PauseStatusUpdated\",\"inputs\":[{\"name\":\"paused\",\"type\":\"bool\",\"indexed\":true,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawalReceived\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EmptyCode\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InsufficientAmount\",\"inputs\":[{\"name\":\"appChainAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxTotalCosts\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidFeeTokenDecimals\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationFailed\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertData_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"NoChainIds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoChange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoKeys\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterOutOfTypeBounds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Paused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransferFromFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnsupportedChainId\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAppChainGateway\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroBalance\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroFeeToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroMigrator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroParameterRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroRecipient\",\"inputs\":[]}]",
}

// SettlementChainGatewayABI is the input ABI used to generate the binding from.
// Deprecated: Use SettlementChainGatewayMetaData.ABI instead.
var SettlementChainGatewayABI = SettlementChainGatewayMetaData.ABI

// SettlementChainGateway is an auto generated Go binding around an Ethereum contract.
type SettlementChainGateway struct {
	SettlementChainGatewayCaller     // Read-only binding to the contract
	SettlementChainGatewayTransactor // Write-only binding to the contract
	SettlementChainGatewayFilterer   // Log filterer for contract events
}

// SettlementChainGatewayCaller is an auto generated read-only Go binding around an Ethereum contract.
type SettlementChainGatewayCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementChainGatewayTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SettlementChainGatewayTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementChainGatewayFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SettlementChainGatewayFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementChainGatewaySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SettlementChainGatewaySession struct {
	Contract     *SettlementChainGateway // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// SettlementChainGatewayCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SettlementChainGatewayCallerSession struct {
	Contract *SettlementChainGatewayCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// SettlementChainGatewayTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SettlementChainGatewayTransactorSession struct {
	Contract     *SettlementChainGatewayTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// SettlementChainGatewayRaw is an auto generated low-level Go binding around an Ethereum contract.
type SettlementChainGatewayRaw struct {
	Contract *SettlementChainGateway // Generic contract binding to access the raw methods on
}

// SettlementChainGatewayCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SettlementChainGatewayCallerRaw struct {
	Contract *SettlementChainGatewayCaller // Generic read-only contract binding to access the raw methods on
}

// SettlementChainGatewayTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SettlementChainGatewayTransactorRaw struct {
	Contract *SettlementChainGatewayTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSettlementChainGateway creates a new instance of SettlementChainGateway, bound to a specific deployed contract.
func NewSettlementChainGateway(address common.Address, backend bind.ContractBackend) (*SettlementChainGateway, error) {
	contract, err := bindSettlementChainGateway(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SettlementChainGateway{SettlementChainGatewayCaller: SettlementChainGatewayCaller{contract: contract}, SettlementChainGatewayTransactor: SettlementChainGatewayTransactor{contract: contract}, SettlementChainGatewayFilterer: SettlementChainGatewayFilterer{contract: contract}}, nil
}

// NewSettlementChainGatewayCaller creates a new read-only instance of SettlementChainGateway, bound to a specific deployed contract.
func NewSettlementChainGatewayCaller(address common.Address, caller bind.ContractCaller) (*SettlementChainGatewayCaller, error) {
	contract, err := bindSettlementChainGateway(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SettlementChainGatewayCaller{contract: contract}, nil
}

// NewSettlementChainGatewayTransactor creates a new write-only instance of SettlementChainGateway, bound to a specific deployed contract.
func NewSettlementChainGatewayTransactor(address common.Address, transactor bind.ContractTransactor) (*SettlementChainGatewayTransactor, error) {
	contract, err := bindSettlementChainGateway(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SettlementChainGatewayTransactor{contract: contract}, nil
}

// NewSettlementChainGatewayFilterer creates a new log filterer instance of SettlementChainGateway, bound to a specific deployed contract.
func NewSettlementChainGatewayFilterer(address common.Address, filterer bind.ContractFilterer) (*SettlementChainGatewayFilterer, error) {
	contract, err := bindSettlementChainGateway(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SettlementChainGatewayFilterer{contract: contract}, nil
}

// bindSettlementChainGateway binds a generic wrapper to an already deployed contract.
func bindSettlementChainGateway(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SettlementChainGatewayMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SettlementChainGateway *SettlementChainGatewayRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SettlementChainGateway.Contract.SettlementChainGatewayCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SettlementChainGateway *SettlementChainGatewayRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.SettlementChainGatewayTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SettlementChainGateway *SettlementChainGatewayRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.SettlementChainGatewayTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SettlementChainGateway *SettlementChainGatewayCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SettlementChainGateway.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SettlementChainGateway *SettlementChainGatewayTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SettlementChainGateway *SettlementChainGatewayTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.contract.Transact(opts, method, params...)
}

// AppChainAlias is a free data retrieval call binding the contract method 0xc5b2989c.
//
// Solidity: function appChainAlias() view returns(address alias_)
func (_SettlementChainGateway *SettlementChainGatewayCaller) AppChainAlias(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SettlementChainGateway.contract.Call(opts, &out, "appChainAlias")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AppChainAlias is a free data retrieval call binding the contract method 0xc5b2989c.
//
// Solidity: function appChainAlias() view returns(address alias_)
func (_SettlementChainGateway *SettlementChainGatewaySession) AppChainAlias() (common.Address, error) {
	return _SettlementChainGateway.Contract.AppChainAlias(&_SettlementChainGateway.CallOpts)
}

// AppChainAlias is a free data retrieval call binding the contract method 0xc5b2989c.
//
// Solidity: function appChainAlias() view returns(address alias_)
func (_SettlementChainGateway *SettlementChainGatewayCallerSession) AppChainAlias() (common.Address, error) {
	return _SettlementChainGateway.Contract.AppChainAlias(&_SettlementChainGateway.CallOpts)
}

// AppChainGateway is a free data retrieval call binding the contract method 0x4623bf41.
//
// Solidity: function appChainGateway() view returns(address)
func (_SettlementChainGateway *SettlementChainGatewayCaller) AppChainGateway(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SettlementChainGateway.contract.Call(opts, &out, "appChainGateway")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AppChainGateway is a free data retrieval call binding the contract method 0x4623bf41.
//
// Solidity: function appChainGateway() view returns(address)
func (_SettlementChainGateway *SettlementChainGatewaySession) AppChainGateway() (common.Address, error) {
	return _SettlementChainGateway.Contract.AppChainGateway(&_SettlementChainGateway.CallOpts)
}

// AppChainGateway is a free data retrieval call binding the contract method 0x4623bf41.
//
// Solidity: function appChainGateway() view returns(address)
func (_SettlementChainGateway *SettlementChainGatewayCallerSession) AppChainGateway() (common.Address, error) {
	return _SettlementChainGateway.Contract.AppChainGateway(&_SettlementChainGateway.CallOpts)
}

// CalculateMaxDepositFee is a free data retrieval call binding the contract method 0xcc7e8ab0.
//
// Solidity: function calculateMaxDepositFee(uint256 chainId_, uint256 gasLimit_, uint256 maxFeePerGas_, uint256 maxBaseFee_) view returns(uint256 fees_)
func (_SettlementChainGateway *SettlementChainGatewayCaller) CalculateMaxDepositFee(opts *bind.CallOpts, chainId_ *big.Int, gasLimit_ *big.Int, maxFeePerGas_ *big.Int, maxBaseFee_ *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _SettlementChainGateway.contract.Call(opts, &out, "calculateMaxDepositFee", chainId_, gasLimit_, maxFeePerGas_, maxBaseFee_)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateMaxDepositFee is a free data retrieval call binding the contract method 0xcc7e8ab0.
//
// Solidity: function calculateMaxDepositFee(uint256 chainId_, uint256 gasLimit_, uint256 maxFeePerGas_, uint256 maxBaseFee_) view returns(uint256 fees_)
func (_SettlementChainGateway *SettlementChainGatewaySession) CalculateMaxDepositFee(chainId_ *big.Int, gasLimit_ *big.Int, maxFeePerGas_ *big.Int, maxBaseFee_ *big.Int) (*big.Int, error) {
	return _SettlementChainGateway.Contract.CalculateMaxDepositFee(&_SettlementChainGateway.CallOpts, chainId_, gasLimit_, maxFeePerGas_, maxBaseFee_)
}

// CalculateMaxDepositFee is a free data retrieval call binding the contract method 0xcc7e8ab0.
//
// Solidity: function calculateMaxDepositFee(uint256 chainId_, uint256 gasLimit_, uint256 maxFeePerGas_, uint256 maxBaseFee_) view returns(uint256 fees_)
func (_SettlementChainGateway *SettlementChainGatewayCallerSession) CalculateMaxDepositFee(chainId_ *big.Int, gasLimit_ *big.Int, maxFeePerGas_ *big.Int, maxBaseFee_ *big.Int) (*big.Int, error) {
	return _SettlementChainGateway.Contract.CalculateMaxDepositFee(&_SettlementChainGateway.CallOpts, chainId_, gasLimit_, maxFeePerGas_, maxBaseFee_)
}

// FeeToken is a free data retrieval call binding the contract method 0x647846a5.
//
// Solidity: function feeToken() view returns(address)
func (_SettlementChainGateway *SettlementChainGatewayCaller) FeeToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SettlementChainGateway.contract.Call(opts, &out, "feeToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeToken is a free data retrieval call binding the contract method 0x647846a5.
//
// Solidity: function feeToken() view returns(address)
func (_SettlementChainGateway *SettlementChainGatewaySession) FeeToken() (common.Address, error) {
	return _SettlementChainGateway.Contract.FeeToken(&_SettlementChainGateway.CallOpts)
}

// FeeToken is a free data retrieval call binding the contract method 0x647846a5.
//
// Solidity: function feeToken() view returns(address)
func (_SettlementChainGateway *SettlementChainGatewayCallerSession) FeeToken() (common.Address, error) {
	return _SettlementChainGateway.Contract.FeeToken(&_SettlementChainGateway.CallOpts)
}

// GetInbox is a free data retrieval call binding the contract method 0xf9a89218.
//
// Solidity: function getInbox(uint256 chainId_) view returns(address inbox_)
func (_SettlementChainGateway *SettlementChainGatewayCaller) GetInbox(opts *bind.CallOpts, chainId_ *big.Int) (common.Address, error) {
	var out []interface{}
	err := _SettlementChainGateway.contract.Call(opts, &out, "getInbox", chainId_)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetInbox is a free data retrieval call binding the contract method 0xf9a89218.
//
// Solidity: function getInbox(uint256 chainId_) view returns(address inbox_)
func (_SettlementChainGateway *SettlementChainGatewaySession) GetInbox(chainId_ *big.Int) (common.Address, error) {
	return _SettlementChainGateway.Contract.GetInbox(&_SettlementChainGateway.CallOpts, chainId_)
}

// GetInbox is a free data retrieval call binding the contract method 0xf9a89218.
//
// Solidity: function getInbox(uint256 chainId_) view returns(address inbox_)
func (_SettlementChainGateway *SettlementChainGatewayCallerSession) GetInbox(chainId_ *big.Int) (common.Address, error) {
	return _SettlementChainGateway.Contract.GetInbox(&_SettlementChainGateway.CallOpts, chainId_)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_SettlementChainGateway *SettlementChainGatewayCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SettlementChainGateway.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_SettlementChainGateway *SettlementChainGatewaySession) Implementation() (common.Address, error) {
	return _SettlementChainGateway.Contract.Implementation(&_SettlementChainGateway.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_SettlementChainGateway *SettlementChainGatewayCallerSession) Implementation() (common.Address, error) {
	return _SettlementChainGateway.Contract.Implementation(&_SettlementChainGateway.CallOpts)
}

// InboxParameterKey is a free data retrieval call binding the contract method 0xe92824d0.
//
// Solidity: function inboxParameterKey() pure returns(string key_)
func (_SettlementChainGateway *SettlementChainGatewayCaller) InboxParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SettlementChainGateway.contract.Call(opts, &out, "inboxParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// InboxParameterKey is a free data retrieval call binding the contract method 0xe92824d0.
//
// Solidity: function inboxParameterKey() pure returns(string key_)
func (_SettlementChainGateway *SettlementChainGatewaySession) InboxParameterKey() (string, error) {
	return _SettlementChainGateway.Contract.InboxParameterKey(&_SettlementChainGateway.CallOpts)
}

// InboxParameterKey is a free data retrieval call binding the contract method 0xe92824d0.
//
// Solidity: function inboxParameterKey() pure returns(string key_)
func (_SettlementChainGateway *SettlementChainGatewayCallerSession) InboxParameterKey() (string, error) {
	return _SettlementChainGateway.Contract.InboxParameterKey(&_SettlementChainGateway.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_SettlementChainGateway *SettlementChainGatewayCaller) MigratorParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SettlementChainGateway.contract.Call(opts, &out, "migratorParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_SettlementChainGateway *SettlementChainGatewaySession) MigratorParameterKey() (string, error) {
	return _SettlementChainGateway.Contract.MigratorParameterKey(&_SettlementChainGateway.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_SettlementChainGateway *SettlementChainGatewayCallerSession) MigratorParameterKey() (string, error) {
	return _SettlementChainGateway.Contract.MigratorParameterKey(&_SettlementChainGateway.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_SettlementChainGateway *SettlementChainGatewayCaller) ParameterRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SettlementChainGateway.contract.Call(opts, &out, "parameterRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_SettlementChainGateway *SettlementChainGatewaySession) ParameterRegistry() (common.Address, error) {
	return _SettlementChainGateway.Contract.ParameterRegistry(&_SettlementChainGateway.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_SettlementChainGateway *SettlementChainGatewayCallerSession) ParameterRegistry() (common.Address, error) {
	return _SettlementChainGateway.Contract.ParameterRegistry(&_SettlementChainGateway.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool paused_)
func (_SettlementChainGateway *SettlementChainGatewayCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SettlementChainGateway.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool paused_)
func (_SettlementChainGateway *SettlementChainGatewaySession) Paused() (bool, error) {
	return _SettlementChainGateway.Contract.Paused(&_SettlementChainGateway.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool paused_)
func (_SettlementChainGateway *SettlementChainGatewayCallerSession) Paused() (bool, error) {
	return _SettlementChainGateway.Contract.Paused(&_SettlementChainGateway.CallOpts)
}

// PausedParameterKey is a free data retrieval call binding the contract method 0xcc5999af.
//
// Solidity: function pausedParameterKey() pure returns(string key_)
func (_SettlementChainGateway *SettlementChainGatewayCaller) PausedParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SettlementChainGateway.contract.Call(opts, &out, "pausedParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// PausedParameterKey is a free data retrieval call binding the contract method 0xcc5999af.
//
// Solidity: function pausedParameterKey() pure returns(string key_)
func (_SettlementChainGateway *SettlementChainGatewaySession) PausedParameterKey() (string, error) {
	return _SettlementChainGateway.Contract.PausedParameterKey(&_SettlementChainGateway.CallOpts)
}

// PausedParameterKey is a free data retrieval call binding the contract method 0xcc5999af.
//
// Solidity: function pausedParameterKey() pure returns(string key_)
func (_SettlementChainGateway *SettlementChainGatewayCallerSession) PausedParameterKey() (string, error) {
	return _SettlementChainGateway.Contract.PausedParameterKey(&_SettlementChainGateway.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0x57223aab.
//
// Solidity: function deposit(uint256 chainId_, address recipient_, uint256 amount_, uint256 gasLimit_, uint256 maxFeePerGas_) returns()
func (_SettlementChainGateway *SettlementChainGatewayTransactor) Deposit(opts *bind.TransactOpts, chainId_ *big.Int, recipient_ common.Address, amount_ *big.Int, gasLimit_ *big.Int, maxFeePerGas_ *big.Int) (*types.Transaction, error) {
	return _SettlementChainGateway.contract.Transact(opts, "deposit", chainId_, recipient_, amount_, gasLimit_, maxFeePerGas_)
}

// Deposit is a paid mutator transaction binding the contract method 0x57223aab.
//
// Solidity: function deposit(uint256 chainId_, address recipient_, uint256 amount_, uint256 gasLimit_, uint256 maxFeePerGas_) returns()
func (_SettlementChainGateway *SettlementChainGatewaySession) Deposit(chainId_ *big.Int, recipient_ common.Address, amount_ *big.Int, gasLimit_ *big.Int, maxFeePerGas_ *big.Int) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.Deposit(&_SettlementChainGateway.TransactOpts, chainId_, recipient_, amount_, gasLimit_, maxFeePerGas_)
}

// Deposit is a paid mutator transaction binding the contract method 0x57223aab.
//
// Solidity: function deposit(uint256 chainId_, address recipient_, uint256 amount_, uint256 gasLimit_, uint256 maxFeePerGas_) returns()
func (_SettlementChainGateway *SettlementChainGatewayTransactorSession) Deposit(chainId_ *big.Int, recipient_ common.Address, amount_ *big.Int, gasLimit_ *big.Int, maxFeePerGas_ *big.Int) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.Deposit(&_SettlementChainGateway.TransactOpts, chainId_, recipient_, amount_, gasLimit_, maxFeePerGas_)
}

// DepositFromUnderlying is a paid mutator transaction binding the contract method 0x62e63622.
//
// Solidity: function depositFromUnderlying(uint256 chainId_, address recipient_, uint256 amount_, uint256 gasLimit_, uint256 maxFeePerGas_) returns()
func (_SettlementChainGateway *SettlementChainGatewayTransactor) DepositFromUnderlying(opts *bind.TransactOpts, chainId_ *big.Int, recipient_ common.Address, amount_ *big.Int, gasLimit_ *big.Int, maxFeePerGas_ *big.Int) (*types.Transaction, error) {
	return _SettlementChainGateway.contract.Transact(opts, "depositFromUnderlying", chainId_, recipient_, amount_, gasLimit_, maxFeePerGas_)
}

// DepositFromUnderlying is a paid mutator transaction binding the contract method 0x62e63622.
//
// Solidity: function depositFromUnderlying(uint256 chainId_, address recipient_, uint256 amount_, uint256 gasLimit_, uint256 maxFeePerGas_) returns()
func (_SettlementChainGateway *SettlementChainGatewaySession) DepositFromUnderlying(chainId_ *big.Int, recipient_ common.Address, amount_ *big.Int, gasLimit_ *big.Int, maxFeePerGas_ *big.Int) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.DepositFromUnderlying(&_SettlementChainGateway.TransactOpts, chainId_, recipient_, amount_, gasLimit_, maxFeePerGas_)
}

// DepositFromUnderlying is a paid mutator transaction binding the contract method 0x62e63622.
//
// Solidity: function depositFromUnderlying(uint256 chainId_, address recipient_, uint256 amount_, uint256 gasLimit_, uint256 maxFeePerGas_) returns()
func (_SettlementChainGateway *SettlementChainGatewayTransactorSession) DepositFromUnderlying(chainId_ *big.Int, recipient_ common.Address, amount_ *big.Int, gasLimit_ *big.Int, maxFeePerGas_ *big.Int) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.DepositFromUnderlying(&_SettlementChainGateway.TransactOpts, chainId_, recipient_, amount_, gasLimit_, maxFeePerGas_)
}

// DepositFromUnderlyingWithPermit is a paid mutator transaction binding the contract method 0x94341bd5.
//
// Solidity: function depositFromUnderlyingWithPermit(uint256 chainId_, address recipient_, uint256 amount_, uint256 gasLimit_, uint256 maxFeePerGas_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns()
func (_SettlementChainGateway *SettlementChainGatewayTransactor) DepositFromUnderlyingWithPermit(opts *bind.TransactOpts, chainId_ *big.Int, recipient_ common.Address, amount_ *big.Int, gasLimit_ *big.Int, maxFeePerGas_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _SettlementChainGateway.contract.Transact(opts, "depositFromUnderlyingWithPermit", chainId_, recipient_, amount_, gasLimit_, maxFeePerGas_, deadline_, v_, r_, s_)
}

// DepositFromUnderlyingWithPermit is a paid mutator transaction binding the contract method 0x94341bd5.
//
// Solidity: function depositFromUnderlyingWithPermit(uint256 chainId_, address recipient_, uint256 amount_, uint256 gasLimit_, uint256 maxFeePerGas_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns()
func (_SettlementChainGateway *SettlementChainGatewaySession) DepositFromUnderlyingWithPermit(chainId_ *big.Int, recipient_ common.Address, amount_ *big.Int, gasLimit_ *big.Int, maxFeePerGas_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.DepositFromUnderlyingWithPermit(&_SettlementChainGateway.TransactOpts, chainId_, recipient_, amount_, gasLimit_, maxFeePerGas_, deadline_, v_, r_, s_)
}

// DepositFromUnderlyingWithPermit is a paid mutator transaction binding the contract method 0x94341bd5.
//
// Solidity: function depositFromUnderlyingWithPermit(uint256 chainId_, address recipient_, uint256 amount_, uint256 gasLimit_, uint256 maxFeePerGas_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns()
func (_SettlementChainGateway *SettlementChainGatewayTransactorSession) DepositFromUnderlyingWithPermit(chainId_ *big.Int, recipient_ common.Address, amount_ *big.Int, gasLimit_ *big.Int, maxFeePerGas_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.DepositFromUnderlyingWithPermit(&_SettlementChainGateway.TransactOpts, chainId_, recipient_, amount_, gasLimit_, maxFeePerGas_, deadline_, v_, r_, s_)
}

// DepositWithPermit is a paid mutator transaction binding the contract method 0x59e79866.
//
// Solidity: function depositWithPermit(uint256 chainId_, address recipient_, uint256 amount_, uint256 gasLimit_, uint256 maxFeePerGas_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns()
func (_SettlementChainGateway *SettlementChainGatewayTransactor) DepositWithPermit(opts *bind.TransactOpts, chainId_ *big.Int, recipient_ common.Address, amount_ *big.Int, gasLimit_ *big.Int, maxFeePerGas_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _SettlementChainGateway.contract.Transact(opts, "depositWithPermit", chainId_, recipient_, amount_, gasLimit_, maxFeePerGas_, deadline_, v_, r_, s_)
}

// DepositWithPermit is a paid mutator transaction binding the contract method 0x59e79866.
//
// Solidity: function depositWithPermit(uint256 chainId_, address recipient_, uint256 amount_, uint256 gasLimit_, uint256 maxFeePerGas_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns()
func (_SettlementChainGateway *SettlementChainGatewaySession) DepositWithPermit(chainId_ *big.Int, recipient_ common.Address, amount_ *big.Int, gasLimit_ *big.Int, maxFeePerGas_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.DepositWithPermit(&_SettlementChainGateway.TransactOpts, chainId_, recipient_, amount_, gasLimit_, maxFeePerGas_, deadline_, v_, r_, s_)
}

// DepositWithPermit is a paid mutator transaction binding the contract method 0x59e79866.
//
// Solidity: function depositWithPermit(uint256 chainId_, address recipient_, uint256 amount_, uint256 gasLimit_, uint256 maxFeePerGas_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns()
func (_SettlementChainGateway *SettlementChainGatewayTransactorSession) DepositWithPermit(chainId_ *big.Int, recipient_ common.Address, amount_ *big.Int, gasLimit_ *big.Int, maxFeePerGas_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.DepositWithPermit(&_SettlementChainGateway.TransactOpts, chainId_, recipient_, amount_, gasLimit_, maxFeePerGas_, deadline_, v_, r_, s_)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_SettlementChainGateway *SettlementChainGatewayTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SettlementChainGateway.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_SettlementChainGateway *SettlementChainGatewaySession) Initialize() (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.Initialize(&_SettlementChainGateway.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_SettlementChainGateway *SettlementChainGatewayTransactorSession) Initialize() (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.Initialize(&_SettlementChainGateway.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_SettlementChainGateway *SettlementChainGatewayTransactor) Migrate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SettlementChainGateway.contract.Transact(opts, "migrate")
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_SettlementChainGateway *SettlementChainGatewaySession) Migrate() (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.Migrate(&_SettlementChainGateway.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_SettlementChainGateway *SettlementChainGatewayTransactorSession) Migrate() (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.Migrate(&_SettlementChainGateway.TransactOpts)
}

// ReceiveWithdrawal is a paid mutator transaction binding the contract method 0x573b8299.
//
// Solidity: function receiveWithdrawal(address recipient_) returns(uint256 amount_)
func (_SettlementChainGateway *SettlementChainGatewayTransactor) ReceiveWithdrawal(opts *bind.TransactOpts, recipient_ common.Address) (*types.Transaction, error) {
	return _SettlementChainGateway.contract.Transact(opts, "receiveWithdrawal", recipient_)
}

// ReceiveWithdrawal is a paid mutator transaction binding the contract method 0x573b8299.
//
// Solidity: function receiveWithdrawal(address recipient_) returns(uint256 amount_)
func (_SettlementChainGateway *SettlementChainGatewaySession) ReceiveWithdrawal(recipient_ common.Address) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.ReceiveWithdrawal(&_SettlementChainGateway.TransactOpts, recipient_)
}

// ReceiveWithdrawal is a paid mutator transaction binding the contract method 0x573b8299.
//
// Solidity: function receiveWithdrawal(address recipient_) returns(uint256 amount_)
func (_SettlementChainGateway *SettlementChainGatewayTransactorSession) ReceiveWithdrawal(recipient_ common.Address) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.ReceiveWithdrawal(&_SettlementChainGateway.TransactOpts, recipient_)
}

// ReceiveWithdrawalIntoUnderlying is a paid mutator transaction binding the contract method 0xa0f94c84.
//
// Solidity: function receiveWithdrawalIntoUnderlying(address recipient_) returns(uint256 amount_)
func (_SettlementChainGateway *SettlementChainGatewayTransactor) ReceiveWithdrawalIntoUnderlying(opts *bind.TransactOpts, recipient_ common.Address) (*types.Transaction, error) {
	return _SettlementChainGateway.contract.Transact(opts, "receiveWithdrawalIntoUnderlying", recipient_)
}

// ReceiveWithdrawalIntoUnderlying is a paid mutator transaction binding the contract method 0xa0f94c84.
//
// Solidity: function receiveWithdrawalIntoUnderlying(address recipient_) returns(uint256 amount_)
func (_SettlementChainGateway *SettlementChainGatewaySession) ReceiveWithdrawalIntoUnderlying(recipient_ common.Address) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.ReceiveWithdrawalIntoUnderlying(&_SettlementChainGateway.TransactOpts, recipient_)
}

// ReceiveWithdrawalIntoUnderlying is a paid mutator transaction binding the contract method 0xa0f94c84.
//
// Solidity: function receiveWithdrawalIntoUnderlying(address recipient_) returns(uint256 amount_)
func (_SettlementChainGateway *SettlementChainGatewayTransactorSession) ReceiveWithdrawalIntoUnderlying(recipient_ common.Address) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.ReceiveWithdrawalIntoUnderlying(&_SettlementChainGateway.TransactOpts, recipient_)
}

// SendParameters is a paid mutator transaction binding the contract method 0xf93e63ac.
//
// Solidity: function sendParameters(uint256[] chainIds_, string[] keys_, uint256 gasLimit_, uint256 maxFeePerGas_, uint256 amountToSend_) returns(uint256 totalSent_)
func (_SettlementChainGateway *SettlementChainGatewayTransactor) SendParameters(opts *bind.TransactOpts, chainIds_ []*big.Int, keys_ []string, gasLimit_ *big.Int, maxFeePerGas_ *big.Int, amountToSend_ *big.Int) (*types.Transaction, error) {
	return _SettlementChainGateway.contract.Transact(opts, "sendParameters", chainIds_, keys_, gasLimit_, maxFeePerGas_, amountToSend_)
}

// SendParameters is a paid mutator transaction binding the contract method 0xf93e63ac.
//
// Solidity: function sendParameters(uint256[] chainIds_, string[] keys_, uint256 gasLimit_, uint256 maxFeePerGas_, uint256 amountToSend_) returns(uint256 totalSent_)
func (_SettlementChainGateway *SettlementChainGatewaySession) SendParameters(chainIds_ []*big.Int, keys_ []string, gasLimit_ *big.Int, maxFeePerGas_ *big.Int, amountToSend_ *big.Int) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.SendParameters(&_SettlementChainGateway.TransactOpts, chainIds_, keys_, gasLimit_, maxFeePerGas_, amountToSend_)
}

// SendParameters is a paid mutator transaction binding the contract method 0xf93e63ac.
//
// Solidity: function sendParameters(uint256[] chainIds_, string[] keys_, uint256 gasLimit_, uint256 maxFeePerGas_, uint256 amountToSend_) returns(uint256 totalSent_)
func (_SettlementChainGateway *SettlementChainGatewayTransactorSession) SendParameters(chainIds_ []*big.Int, keys_ []string, gasLimit_ *big.Int, maxFeePerGas_ *big.Int, amountToSend_ *big.Int) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.SendParameters(&_SettlementChainGateway.TransactOpts, chainIds_, keys_, gasLimit_, maxFeePerGas_, amountToSend_)
}

// SendParametersFromUnderlying is a paid mutator transaction binding the contract method 0x967aa768.
//
// Solidity: function sendParametersFromUnderlying(uint256[] chainIds_, string[] keys_, uint256 gasLimit_, uint256 maxFeePerGas_, uint256 amountToSend_) returns(uint256 totalSent_)
func (_SettlementChainGateway *SettlementChainGatewayTransactor) SendParametersFromUnderlying(opts *bind.TransactOpts, chainIds_ []*big.Int, keys_ []string, gasLimit_ *big.Int, maxFeePerGas_ *big.Int, amountToSend_ *big.Int) (*types.Transaction, error) {
	return _SettlementChainGateway.contract.Transact(opts, "sendParametersFromUnderlying", chainIds_, keys_, gasLimit_, maxFeePerGas_, amountToSend_)
}

// SendParametersFromUnderlying is a paid mutator transaction binding the contract method 0x967aa768.
//
// Solidity: function sendParametersFromUnderlying(uint256[] chainIds_, string[] keys_, uint256 gasLimit_, uint256 maxFeePerGas_, uint256 amountToSend_) returns(uint256 totalSent_)
func (_SettlementChainGateway *SettlementChainGatewaySession) SendParametersFromUnderlying(chainIds_ []*big.Int, keys_ []string, gasLimit_ *big.Int, maxFeePerGas_ *big.Int, amountToSend_ *big.Int) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.SendParametersFromUnderlying(&_SettlementChainGateway.TransactOpts, chainIds_, keys_, gasLimit_, maxFeePerGas_, amountToSend_)
}

// SendParametersFromUnderlying is a paid mutator transaction binding the contract method 0x967aa768.
//
// Solidity: function sendParametersFromUnderlying(uint256[] chainIds_, string[] keys_, uint256 gasLimit_, uint256 maxFeePerGas_, uint256 amountToSend_) returns(uint256 totalSent_)
func (_SettlementChainGateway *SettlementChainGatewayTransactorSession) SendParametersFromUnderlying(chainIds_ []*big.Int, keys_ []string, gasLimit_ *big.Int, maxFeePerGas_ *big.Int, amountToSend_ *big.Int) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.SendParametersFromUnderlying(&_SettlementChainGateway.TransactOpts, chainIds_, keys_, gasLimit_, maxFeePerGas_, amountToSend_)
}

// SendParametersFromUnderlyingWithPermit is a paid mutator transaction binding the contract method 0xb74d7f01.
//
// Solidity: function sendParametersFromUnderlyingWithPermit(uint256[] chainIds_, string[] keys_, uint256 gasLimit_, uint256 maxFeePerGas_, uint256 amountToSend_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns(uint256 totalSent_)
func (_SettlementChainGateway *SettlementChainGatewayTransactor) SendParametersFromUnderlyingWithPermit(opts *bind.TransactOpts, chainIds_ []*big.Int, keys_ []string, gasLimit_ *big.Int, maxFeePerGas_ *big.Int, amountToSend_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _SettlementChainGateway.contract.Transact(opts, "sendParametersFromUnderlyingWithPermit", chainIds_, keys_, gasLimit_, maxFeePerGas_, amountToSend_, deadline_, v_, r_, s_)
}

// SendParametersFromUnderlyingWithPermit is a paid mutator transaction binding the contract method 0xb74d7f01.
//
// Solidity: function sendParametersFromUnderlyingWithPermit(uint256[] chainIds_, string[] keys_, uint256 gasLimit_, uint256 maxFeePerGas_, uint256 amountToSend_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns(uint256 totalSent_)
func (_SettlementChainGateway *SettlementChainGatewaySession) SendParametersFromUnderlyingWithPermit(chainIds_ []*big.Int, keys_ []string, gasLimit_ *big.Int, maxFeePerGas_ *big.Int, amountToSend_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.SendParametersFromUnderlyingWithPermit(&_SettlementChainGateway.TransactOpts, chainIds_, keys_, gasLimit_, maxFeePerGas_, amountToSend_, deadline_, v_, r_, s_)
}

// SendParametersFromUnderlyingWithPermit is a paid mutator transaction binding the contract method 0xb74d7f01.
//
// Solidity: function sendParametersFromUnderlyingWithPermit(uint256[] chainIds_, string[] keys_, uint256 gasLimit_, uint256 maxFeePerGas_, uint256 amountToSend_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns(uint256 totalSent_)
func (_SettlementChainGateway *SettlementChainGatewayTransactorSession) SendParametersFromUnderlyingWithPermit(chainIds_ []*big.Int, keys_ []string, gasLimit_ *big.Int, maxFeePerGas_ *big.Int, amountToSend_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.SendParametersFromUnderlyingWithPermit(&_SettlementChainGateway.TransactOpts, chainIds_, keys_, gasLimit_, maxFeePerGas_, amountToSend_, deadline_, v_, r_, s_)
}

// SendParametersWithPermit is a paid mutator transaction binding the contract method 0x6147efb6.
//
// Solidity: function sendParametersWithPermit(uint256[] chainIds_, string[] keys_, uint256 gasLimit_, uint256 maxFeePerGas_, uint256 amountToSend_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns(uint256 totalSent_)
func (_SettlementChainGateway *SettlementChainGatewayTransactor) SendParametersWithPermit(opts *bind.TransactOpts, chainIds_ []*big.Int, keys_ []string, gasLimit_ *big.Int, maxFeePerGas_ *big.Int, amountToSend_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _SettlementChainGateway.contract.Transact(opts, "sendParametersWithPermit", chainIds_, keys_, gasLimit_, maxFeePerGas_, amountToSend_, deadline_, v_, r_, s_)
}

// SendParametersWithPermit is a paid mutator transaction binding the contract method 0x6147efb6.
//
// Solidity: function sendParametersWithPermit(uint256[] chainIds_, string[] keys_, uint256 gasLimit_, uint256 maxFeePerGas_, uint256 amountToSend_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns(uint256 totalSent_)
func (_SettlementChainGateway *SettlementChainGatewaySession) SendParametersWithPermit(chainIds_ []*big.Int, keys_ []string, gasLimit_ *big.Int, maxFeePerGas_ *big.Int, amountToSend_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.SendParametersWithPermit(&_SettlementChainGateway.TransactOpts, chainIds_, keys_, gasLimit_, maxFeePerGas_, amountToSend_, deadline_, v_, r_, s_)
}

// SendParametersWithPermit is a paid mutator transaction binding the contract method 0x6147efb6.
//
// Solidity: function sendParametersWithPermit(uint256[] chainIds_, string[] keys_, uint256 gasLimit_, uint256 maxFeePerGas_, uint256 amountToSend_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns(uint256 totalSent_)
func (_SettlementChainGateway *SettlementChainGatewayTransactorSession) SendParametersWithPermit(chainIds_ []*big.Int, keys_ []string, gasLimit_ *big.Int, maxFeePerGas_ *big.Int, amountToSend_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.SendParametersWithPermit(&_SettlementChainGateway.TransactOpts, chainIds_, keys_, gasLimit_, maxFeePerGas_, amountToSend_, deadline_, v_, r_, s_)
}

// UpdateInbox is a paid mutator transaction binding the contract method 0xe7a9e711.
//
// Solidity: function updateInbox(uint256 chainId_) returns()
func (_SettlementChainGateway *SettlementChainGatewayTransactor) UpdateInbox(opts *bind.TransactOpts, chainId_ *big.Int) (*types.Transaction, error) {
	return _SettlementChainGateway.contract.Transact(opts, "updateInbox", chainId_)
}

// UpdateInbox is a paid mutator transaction binding the contract method 0xe7a9e711.
//
// Solidity: function updateInbox(uint256 chainId_) returns()
func (_SettlementChainGateway *SettlementChainGatewaySession) UpdateInbox(chainId_ *big.Int) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.UpdateInbox(&_SettlementChainGateway.TransactOpts, chainId_)
}

// UpdateInbox is a paid mutator transaction binding the contract method 0xe7a9e711.
//
// Solidity: function updateInbox(uint256 chainId_) returns()
func (_SettlementChainGateway *SettlementChainGatewayTransactorSession) UpdateInbox(chainId_ *big.Int) (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.UpdateInbox(&_SettlementChainGateway.TransactOpts, chainId_)
}

// UpdatePauseStatus is a paid mutator transaction binding the contract method 0x59d4df41.
//
// Solidity: function updatePauseStatus() returns()
func (_SettlementChainGateway *SettlementChainGatewayTransactor) UpdatePauseStatus(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SettlementChainGateway.contract.Transact(opts, "updatePauseStatus")
}

// UpdatePauseStatus is a paid mutator transaction binding the contract method 0x59d4df41.
//
// Solidity: function updatePauseStatus() returns()
func (_SettlementChainGateway *SettlementChainGatewaySession) UpdatePauseStatus() (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.UpdatePauseStatus(&_SettlementChainGateway.TransactOpts)
}

// UpdatePauseStatus is a paid mutator transaction binding the contract method 0x59d4df41.
//
// Solidity: function updatePauseStatus() returns()
func (_SettlementChainGateway *SettlementChainGatewayTransactorSession) UpdatePauseStatus() (*types.Transaction, error) {
	return _SettlementChainGateway.Contract.UpdatePauseStatus(&_SettlementChainGateway.TransactOpts)
}

// SettlementChainGatewayDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the SettlementChainGateway contract.
type SettlementChainGatewayDepositIterator struct {
	Event *SettlementChainGatewayDeposit // Event containing the contract specifics and raw log

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
func (it *SettlementChainGatewayDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementChainGatewayDeposit)
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
		it.Event = new(SettlementChainGatewayDeposit)
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
func (it *SettlementChainGatewayDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementChainGatewayDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementChainGatewayDeposit represents a Deposit event raised by the SettlementChainGateway contract.
type SettlementChainGatewayDeposit struct {
	ChainId       *big.Int
	MessageNumber *big.Int
	Recipient     common.Address
	Amount        *big.Int
	MaxFees       *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x565bc26c555ccdedc59901ea62fd4da0e6c846441449fd5ce2b87acd8d63bc01.
//
// Solidity: event Deposit(uint256 indexed chainId, uint256 indexed messageNumber, address indexed recipient, uint256 amount, uint256 maxFees)
func (_SettlementChainGateway *SettlementChainGatewayFilterer) FilterDeposit(opts *bind.FilterOpts, chainId []*big.Int, messageNumber []*big.Int, recipient []common.Address) (*SettlementChainGatewayDepositIterator, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}
	var messageNumberRule []interface{}
	for _, messageNumberItem := range messageNumber {
		messageNumberRule = append(messageNumberRule, messageNumberItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _SettlementChainGateway.contract.FilterLogs(opts, "Deposit", chainIdRule, messageNumberRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &SettlementChainGatewayDepositIterator{contract: _SettlementChainGateway.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x565bc26c555ccdedc59901ea62fd4da0e6c846441449fd5ce2b87acd8d63bc01.
//
// Solidity: event Deposit(uint256 indexed chainId, uint256 indexed messageNumber, address indexed recipient, uint256 amount, uint256 maxFees)
func (_SettlementChainGateway *SettlementChainGatewayFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *SettlementChainGatewayDeposit, chainId []*big.Int, messageNumber []*big.Int, recipient []common.Address) (event.Subscription, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}
	var messageNumberRule []interface{}
	for _, messageNumberItem := range messageNumber {
		messageNumberRule = append(messageNumberRule, messageNumberItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _SettlementChainGateway.contract.WatchLogs(opts, "Deposit", chainIdRule, messageNumberRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementChainGatewayDeposit)
				if err := _SettlementChainGateway.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0x565bc26c555ccdedc59901ea62fd4da0e6c846441449fd5ce2b87acd8d63bc01.
//
// Solidity: event Deposit(uint256 indexed chainId, uint256 indexed messageNumber, address indexed recipient, uint256 amount, uint256 maxFees)
func (_SettlementChainGateway *SettlementChainGatewayFilterer) ParseDeposit(log types.Log) (*SettlementChainGatewayDeposit, error) {
	event := new(SettlementChainGatewayDeposit)
	if err := _SettlementChainGateway.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementChainGatewayInboxUpdatedIterator is returned from FilterInboxUpdated and is used to iterate over the raw logs and unpacked data for InboxUpdated events raised by the SettlementChainGateway contract.
type SettlementChainGatewayInboxUpdatedIterator struct {
	Event *SettlementChainGatewayInboxUpdated // Event containing the contract specifics and raw log

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
func (it *SettlementChainGatewayInboxUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementChainGatewayInboxUpdated)
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
		it.Event = new(SettlementChainGatewayInboxUpdated)
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
func (it *SettlementChainGatewayInboxUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementChainGatewayInboxUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementChainGatewayInboxUpdated represents a InboxUpdated event raised by the SettlementChainGateway contract.
type SettlementChainGatewayInboxUpdated struct {
	ChainId *big.Int
	Inbox   common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInboxUpdated is a free log retrieval operation binding the contract event 0x658d5a189c8f6b7d7a3dad10bb354506270575b40200c1871cf44216ebc99685.
//
// Solidity: event InboxUpdated(uint256 indexed chainId, address indexed inbox)
func (_SettlementChainGateway *SettlementChainGatewayFilterer) FilterInboxUpdated(opts *bind.FilterOpts, chainId []*big.Int, inbox []common.Address) (*SettlementChainGatewayInboxUpdatedIterator, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}
	var inboxRule []interface{}
	for _, inboxItem := range inbox {
		inboxRule = append(inboxRule, inboxItem)
	}

	logs, sub, err := _SettlementChainGateway.contract.FilterLogs(opts, "InboxUpdated", chainIdRule, inboxRule)
	if err != nil {
		return nil, err
	}
	return &SettlementChainGatewayInboxUpdatedIterator{contract: _SettlementChainGateway.contract, event: "InboxUpdated", logs: logs, sub: sub}, nil
}

// WatchInboxUpdated is a free log subscription operation binding the contract event 0x658d5a189c8f6b7d7a3dad10bb354506270575b40200c1871cf44216ebc99685.
//
// Solidity: event InboxUpdated(uint256 indexed chainId, address indexed inbox)
func (_SettlementChainGateway *SettlementChainGatewayFilterer) WatchInboxUpdated(opts *bind.WatchOpts, sink chan<- *SettlementChainGatewayInboxUpdated, chainId []*big.Int, inbox []common.Address) (event.Subscription, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}
	var inboxRule []interface{}
	for _, inboxItem := range inbox {
		inboxRule = append(inboxRule, inboxItem)
	}

	logs, sub, err := _SettlementChainGateway.contract.WatchLogs(opts, "InboxUpdated", chainIdRule, inboxRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementChainGatewayInboxUpdated)
				if err := _SettlementChainGateway.contract.UnpackLog(event, "InboxUpdated", log); err != nil {
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

// ParseInboxUpdated is a log parse operation binding the contract event 0x658d5a189c8f6b7d7a3dad10bb354506270575b40200c1871cf44216ebc99685.
//
// Solidity: event InboxUpdated(uint256 indexed chainId, address indexed inbox)
func (_SettlementChainGateway *SettlementChainGatewayFilterer) ParseInboxUpdated(log types.Log) (*SettlementChainGatewayInboxUpdated, error) {
	event := new(SettlementChainGatewayInboxUpdated)
	if err := _SettlementChainGateway.contract.UnpackLog(event, "InboxUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementChainGatewayInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the SettlementChainGateway contract.
type SettlementChainGatewayInitializedIterator struct {
	Event *SettlementChainGatewayInitialized // Event containing the contract specifics and raw log

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
func (it *SettlementChainGatewayInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementChainGatewayInitialized)
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
		it.Event = new(SettlementChainGatewayInitialized)
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
func (it *SettlementChainGatewayInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementChainGatewayInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementChainGatewayInitialized represents a Initialized event raised by the SettlementChainGateway contract.
type SettlementChainGatewayInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_SettlementChainGateway *SettlementChainGatewayFilterer) FilterInitialized(opts *bind.FilterOpts) (*SettlementChainGatewayInitializedIterator, error) {

	logs, sub, err := _SettlementChainGateway.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &SettlementChainGatewayInitializedIterator{contract: _SettlementChainGateway.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_SettlementChainGateway *SettlementChainGatewayFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *SettlementChainGatewayInitialized) (event.Subscription, error) {

	logs, sub, err := _SettlementChainGateway.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementChainGatewayInitialized)
				if err := _SettlementChainGateway.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_SettlementChainGateway *SettlementChainGatewayFilterer) ParseInitialized(log types.Log) (*SettlementChainGatewayInitialized, error) {
	event := new(SettlementChainGatewayInitialized)
	if err := _SettlementChainGateway.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementChainGatewayMigratedIterator is returned from FilterMigrated and is used to iterate over the raw logs and unpacked data for Migrated events raised by the SettlementChainGateway contract.
type SettlementChainGatewayMigratedIterator struct {
	Event *SettlementChainGatewayMigrated // Event containing the contract specifics and raw log

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
func (it *SettlementChainGatewayMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementChainGatewayMigrated)
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
		it.Event = new(SettlementChainGatewayMigrated)
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
func (it *SettlementChainGatewayMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementChainGatewayMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementChainGatewayMigrated represents a Migrated event raised by the SettlementChainGateway contract.
type SettlementChainGatewayMigrated struct {
	Migrator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMigrated is a free log retrieval operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_SettlementChainGateway *SettlementChainGatewayFilterer) FilterMigrated(opts *bind.FilterOpts, migrator []common.Address) (*SettlementChainGatewayMigratedIterator, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _SettlementChainGateway.contract.FilterLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return &SettlementChainGatewayMigratedIterator{contract: _SettlementChainGateway.contract, event: "Migrated", logs: logs, sub: sub}, nil
}

// WatchMigrated is a free log subscription operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_SettlementChainGateway *SettlementChainGatewayFilterer) WatchMigrated(opts *bind.WatchOpts, sink chan<- *SettlementChainGatewayMigrated, migrator []common.Address) (event.Subscription, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _SettlementChainGateway.contract.WatchLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementChainGatewayMigrated)
				if err := _SettlementChainGateway.contract.UnpackLog(event, "Migrated", log); err != nil {
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
func (_SettlementChainGateway *SettlementChainGatewayFilterer) ParseMigrated(log types.Log) (*SettlementChainGatewayMigrated, error) {
	event := new(SettlementChainGatewayMigrated)
	if err := _SettlementChainGateway.contract.UnpackLog(event, "Migrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementChainGatewayParametersSentIterator is returned from FilterParametersSent and is used to iterate over the raw logs and unpacked data for ParametersSent events raised by the SettlementChainGateway contract.
type SettlementChainGatewayParametersSentIterator struct {
	Event *SettlementChainGatewayParametersSent // Event containing the contract specifics and raw log

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
func (it *SettlementChainGatewayParametersSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementChainGatewayParametersSent)
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
		it.Event = new(SettlementChainGatewayParametersSent)
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
func (it *SettlementChainGatewayParametersSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementChainGatewayParametersSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementChainGatewayParametersSent represents a ParametersSent event raised by the SettlementChainGateway contract.
type SettlementChainGatewayParametersSent struct {
	ChainId       *big.Int
	MessageNumber *big.Int
	Nonce         *big.Int
	Keys          []string
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterParametersSent is a free log retrieval operation binding the contract event 0xf384f39d85facefd419c3628bcac831dcc97846c721065ec7b6bffae10f8445b.
//
// Solidity: event ParametersSent(uint256 indexed chainId, uint256 indexed messageNumber, uint256 nonce, string[] keys)
func (_SettlementChainGateway *SettlementChainGatewayFilterer) FilterParametersSent(opts *bind.FilterOpts, chainId []*big.Int, messageNumber []*big.Int) (*SettlementChainGatewayParametersSentIterator, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}
	var messageNumberRule []interface{}
	for _, messageNumberItem := range messageNumber {
		messageNumberRule = append(messageNumberRule, messageNumberItem)
	}

	logs, sub, err := _SettlementChainGateway.contract.FilterLogs(opts, "ParametersSent", chainIdRule, messageNumberRule)
	if err != nil {
		return nil, err
	}
	return &SettlementChainGatewayParametersSentIterator{contract: _SettlementChainGateway.contract, event: "ParametersSent", logs: logs, sub: sub}, nil
}

// WatchParametersSent is a free log subscription operation binding the contract event 0xf384f39d85facefd419c3628bcac831dcc97846c721065ec7b6bffae10f8445b.
//
// Solidity: event ParametersSent(uint256 indexed chainId, uint256 indexed messageNumber, uint256 nonce, string[] keys)
func (_SettlementChainGateway *SettlementChainGatewayFilterer) WatchParametersSent(opts *bind.WatchOpts, sink chan<- *SettlementChainGatewayParametersSent, chainId []*big.Int, messageNumber []*big.Int) (event.Subscription, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}
	var messageNumberRule []interface{}
	for _, messageNumberItem := range messageNumber {
		messageNumberRule = append(messageNumberRule, messageNumberItem)
	}

	logs, sub, err := _SettlementChainGateway.contract.WatchLogs(opts, "ParametersSent", chainIdRule, messageNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementChainGatewayParametersSent)
				if err := _SettlementChainGateway.contract.UnpackLog(event, "ParametersSent", log); err != nil {
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

// ParseParametersSent is a log parse operation binding the contract event 0xf384f39d85facefd419c3628bcac831dcc97846c721065ec7b6bffae10f8445b.
//
// Solidity: event ParametersSent(uint256 indexed chainId, uint256 indexed messageNumber, uint256 nonce, string[] keys)
func (_SettlementChainGateway *SettlementChainGatewayFilterer) ParseParametersSent(log types.Log) (*SettlementChainGatewayParametersSent, error) {
	event := new(SettlementChainGatewayParametersSent)
	if err := _SettlementChainGateway.contract.UnpackLog(event, "ParametersSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementChainGatewayPauseStatusUpdatedIterator is returned from FilterPauseStatusUpdated and is used to iterate over the raw logs and unpacked data for PauseStatusUpdated events raised by the SettlementChainGateway contract.
type SettlementChainGatewayPauseStatusUpdatedIterator struct {
	Event *SettlementChainGatewayPauseStatusUpdated // Event containing the contract specifics and raw log

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
func (it *SettlementChainGatewayPauseStatusUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementChainGatewayPauseStatusUpdated)
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
		it.Event = new(SettlementChainGatewayPauseStatusUpdated)
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
func (it *SettlementChainGatewayPauseStatusUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementChainGatewayPauseStatusUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementChainGatewayPauseStatusUpdated represents a PauseStatusUpdated event raised by the SettlementChainGateway contract.
type SettlementChainGatewayPauseStatusUpdated struct {
	Paused bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPauseStatusUpdated is a free log retrieval operation binding the contract event 0x7c4d1fe30fdbfda9e9c4c43e759ef32e4db5128d4cb58ff3ae9583b89b6242a5.
//
// Solidity: event PauseStatusUpdated(bool indexed paused)
func (_SettlementChainGateway *SettlementChainGatewayFilterer) FilterPauseStatusUpdated(opts *bind.FilterOpts, paused []bool) (*SettlementChainGatewayPauseStatusUpdatedIterator, error) {

	var pausedRule []interface{}
	for _, pausedItem := range paused {
		pausedRule = append(pausedRule, pausedItem)
	}

	logs, sub, err := _SettlementChainGateway.contract.FilterLogs(opts, "PauseStatusUpdated", pausedRule)
	if err != nil {
		return nil, err
	}
	return &SettlementChainGatewayPauseStatusUpdatedIterator{contract: _SettlementChainGateway.contract, event: "PauseStatusUpdated", logs: logs, sub: sub}, nil
}

// WatchPauseStatusUpdated is a free log subscription operation binding the contract event 0x7c4d1fe30fdbfda9e9c4c43e759ef32e4db5128d4cb58ff3ae9583b89b6242a5.
//
// Solidity: event PauseStatusUpdated(bool indexed paused)
func (_SettlementChainGateway *SettlementChainGatewayFilterer) WatchPauseStatusUpdated(opts *bind.WatchOpts, sink chan<- *SettlementChainGatewayPauseStatusUpdated, paused []bool) (event.Subscription, error) {

	var pausedRule []interface{}
	for _, pausedItem := range paused {
		pausedRule = append(pausedRule, pausedItem)
	}

	logs, sub, err := _SettlementChainGateway.contract.WatchLogs(opts, "PauseStatusUpdated", pausedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementChainGatewayPauseStatusUpdated)
				if err := _SettlementChainGateway.contract.UnpackLog(event, "PauseStatusUpdated", log); err != nil {
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
func (_SettlementChainGateway *SettlementChainGatewayFilterer) ParsePauseStatusUpdated(log types.Log) (*SettlementChainGatewayPauseStatusUpdated, error) {
	event := new(SettlementChainGatewayPauseStatusUpdated)
	if err := _SettlementChainGateway.contract.UnpackLog(event, "PauseStatusUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementChainGatewayUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the SettlementChainGateway contract.
type SettlementChainGatewayUpgradedIterator struct {
	Event *SettlementChainGatewayUpgraded // Event containing the contract specifics and raw log

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
func (it *SettlementChainGatewayUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementChainGatewayUpgraded)
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
		it.Event = new(SettlementChainGatewayUpgraded)
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
func (it *SettlementChainGatewayUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementChainGatewayUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementChainGatewayUpgraded represents a Upgraded event raised by the SettlementChainGateway contract.
type SettlementChainGatewayUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_SettlementChainGateway *SettlementChainGatewayFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*SettlementChainGatewayUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _SettlementChainGateway.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &SettlementChainGatewayUpgradedIterator{contract: _SettlementChainGateway.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_SettlementChainGateway *SettlementChainGatewayFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *SettlementChainGatewayUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _SettlementChainGateway.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementChainGatewayUpgraded)
				if err := _SettlementChainGateway.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_SettlementChainGateway *SettlementChainGatewayFilterer) ParseUpgraded(log types.Log) (*SettlementChainGatewayUpgraded, error) {
	event := new(SettlementChainGatewayUpgraded)
	if err := _SettlementChainGateway.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementChainGatewayWithdrawalReceivedIterator is returned from FilterWithdrawalReceived and is used to iterate over the raw logs and unpacked data for WithdrawalReceived events raised by the SettlementChainGateway contract.
type SettlementChainGatewayWithdrawalReceivedIterator struct {
	Event *SettlementChainGatewayWithdrawalReceived // Event containing the contract specifics and raw log

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
func (it *SettlementChainGatewayWithdrawalReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementChainGatewayWithdrawalReceived)
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
		it.Event = new(SettlementChainGatewayWithdrawalReceived)
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
func (it *SettlementChainGatewayWithdrawalReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementChainGatewayWithdrawalReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementChainGatewayWithdrawalReceived represents a WithdrawalReceived event raised by the SettlementChainGateway contract.
type SettlementChainGatewayWithdrawalReceived struct {
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalReceived is a free log retrieval operation binding the contract event 0xc161e691576b6bf95af8b4cc817d4a54a0304daf970d2ad185e36b7bdebd3f1d.
//
// Solidity: event WithdrawalReceived(address indexed recipient, uint256 amount)
func (_SettlementChainGateway *SettlementChainGatewayFilterer) FilterWithdrawalReceived(opts *bind.FilterOpts, recipient []common.Address) (*SettlementChainGatewayWithdrawalReceivedIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _SettlementChainGateway.contract.FilterLogs(opts, "WithdrawalReceived", recipientRule)
	if err != nil {
		return nil, err
	}
	return &SettlementChainGatewayWithdrawalReceivedIterator{contract: _SettlementChainGateway.contract, event: "WithdrawalReceived", logs: logs, sub: sub}, nil
}

// WatchWithdrawalReceived is a free log subscription operation binding the contract event 0xc161e691576b6bf95af8b4cc817d4a54a0304daf970d2ad185e36b7bdebd3f1d.
//
// Solidity: event WithdrawalReceived(address indexed recipient, uint256 amount)
func (_SettlementChainGateway *SettlementChainGatewayFilterer) WatchWithdrawalReceived(opts *bind.WatchOpts, sink chan<- *SettlementChainGatewayWithdrawalReceived, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _SettlementChainGateway.contract.WatchLogs(opts, "WithdrawalReceived", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementChainGatewayWithdrawalReceived)
				if err := _SettlementChainGateway.contract.UnpackLog(event, "WithdrawalReceived", log); err != nil {
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

// ParseWithdrawalReceived is a log parse operation binding the contract event 0xc161e691576b6bf95af8b4cc817d4a54a0304daf970d2ad185e36b7bdebd3f1d.
//
// Solidity: event WithdrawalReceived(address indexed recipient, uint256 amount)
func (_SettlementChainGateway *SettlementChainGatewayFilterer) ParseWithdrawalReceived(log types.Log) (*SettlementChainGatewayWithdrawalReceived, error) {
	event := new(SettlementChainGatewayWithdrawalReceived)
	if err := _SettlementChainGateway.contract.UnpackLog(event, "WithdrawalReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
