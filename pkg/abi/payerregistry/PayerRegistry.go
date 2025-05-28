// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package payerregistry

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

// IPayerRegistryPayerFee is an auto generated low-level Go binding around an user-defined struct.
type IPayerRegistryPayerFee struct {
	Payer common.Address
	Fee   *big.Int
}

// PayerRegistryMetaData contains all meta data concerning the PayerRegistry contract.
var PayerRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"parameterRegistry_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"token_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelWithdrawal\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"payer_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount_\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositWithPermit\",\"inputs\":[{\"name\":\"payer_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount_\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"deadline_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v_\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"excess\",\"inputs\":[],\"outputs\":[{\"name\":\"excess_\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"feeDistributor\",\"inputs\":[],\"outputs\":[{\"name\":\"feeDistributor_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"feeDistributorParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"finalizeWithdrawal\",\"inputs\":[{\"name\":\"recipient_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getBalance\",\"inputs\":[{\"name\":\"payer_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"balance_\",\"type\":\"int104\",\"internalType\":\"int104\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBalances\",\"inputs\":[{\"name\":\"payers_\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"balances_\",\"type\":\"int104[]\",\"internalType\":\"int104[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPendingWithdrawal\",\"inputs\":[{\"name\":\"payer_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"pendingWithdrawal_\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"withdrawableTimestamp_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"implementation_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migrate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migratorParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"minimumDeposit\",\"inputs\":[],\"outputs\":[{\"name\":\"minimumDeposit_\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minimumDepositParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"parameterRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"paused_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pausedParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"requestWithdrawal\",\"inputs\":[{\"name\":\"amount_\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendExcessToFeeDistributor\",\"inputs\":[],\"outputs\":[{\"name\":\"excess_\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"settleUsage\",\"inputs\":[{\"name\":\"payerFees_\",\"type\":\"tuple[]\",\"internalType\":\"structIPayerRegistry.PayerFee[]\",\"components\":[{\"name\":\"payer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fee\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"outputs\":[{\"name\":\"feesSettled_\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"settler\",\"inputs\":[],\"outputs\":[{\"name\":\"settler_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"settlerParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"token\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalDebt\",\"inputs\":[],\"outputs\":[{\"name\":\"totalDebt_\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalDeposits\",\"inputs\":[],\"outputs\":[{\"name\":\"totalDeposits_\",\"type\":\"int104\",\"internalType\":\"int104\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalWithdrawable\",\"inputs\":[],\"outputs\":[{\"name\":\"totalWithdrawable_\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateFeeDistributor\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMinimumDeposit\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updatePauseStatus\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateSettler\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateWithdrawLockPeriod\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawLockPeriod\",\"inputs\":[],\"outputs\":[{\"name\":\"withdrawLockPeriod_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawLockPeriodParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"payer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ExcessTransferred\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeDistributorUpdated\",\"inputs\":[{\"name\":\"feeDistributor\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Migrated\",\"inputs\":[{\"name\":\"migrator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinimumDepositUpdated\",\"inputs\":[{\"name\":\"minimumDeposit\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PauseStatusUpdated\",\"inputs\":[{\"name\":\"paused\",\"type\":\"bool\",\"indexed\":true,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SettlerUpdated\",\"inputs\":[{\"name\":\"settler\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UsageSettled\",\"inputs\":[{\"name\":\"payer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawLockPeriodUpdated\",\"inputs\":[{\"name\":\"withdrawLockPeriod\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawalCancelled\",\"inputs\":[{\"name\":\"payer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawalFinalized\",\"inputs\":[{\"name\":\"payer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawalRequested\",\"inputs\":[{\"name\":\"payer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"},{\"name\":\"withdrawableTimestamp\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EmptyCode\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InsufficientBalance\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientDeposit\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"minimumDeposit\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationFailed\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertData_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"NoChange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoExcess\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoPendingWithdrawal\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotSettler\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterOutOfTypeBounds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Paused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PayerInDebt\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PendingWithdrawalExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransferFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransferFromFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WithdrawalNotReady\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"withdrawableTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"ZeroFeeDistributor\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroMigrator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroMinimumDeposit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroParameterRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroSettler\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroWithdrawalAmount\",\"inputs\":[]}]",
}

// PayerRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use PayerRegistryMetaData.ABI instead.
var PayerRegistryABI = PayerRegistryMetaData.ABI

// PayerRegistry is an auto generated Go binding around an Ethereum contract.
type PayerRegistry struct {
	PayerRegistryCaller     // Read-only binding to the contract
	PayerRegistryTransactor // Write-only binding to the contract
	PayerRegistryFilterer   // Log filterer for contract events
}

// PayerRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type PayerRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PayerRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PayerRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PayerRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PayerRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PayerRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PayerRegistrySession struct {
	Contract     *PayerRegistry    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PayerRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PayerRegistryCallerSession struct {
	Contract *PayerRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// PayerRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PayerRegistryTransactorSession struct {
	Contract     *PayerRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// PayerRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type PayerRegistryRaw struct {
	Contract *PayerRegistry // Generic contract binding to access the raw methods on
}

// PayerRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PayerRegistryCallerRaw struct {
	Contract *PayerRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// PayerRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PayerRegistryTransactorRaw struct {
	Contract *PayerRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPayerRegistry creates a new instance of PayerRegistry, bound to a specific deployed contract.
func NewPayerRegistry(address common.Address, backend bind.ContractBackend) (*PayerRegistry, error) {
	contract, err := bindPayerRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PayerRegistry{PayerRegistryCaller: PayerRegistryCaller{contract: contract}, PayerRegistryTransactor: PayerRegistryTransactor{contract: contract}, PayerRegistryFilterer: PayerRegistryFilterer{contract: contract}}, nil
}

// NewPayerRegistryCaller creates a new read-only instance of PayerRegistry, bound to a specific deployed contract.
func NewPayerRegistryCaller(address common.Address, caller bind.ContractCaller) (*PayerRegistryCaller, error) {
	contract, err := bindPayerRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PayerRegistryCaller{contract: contract}, nil
}

// NewPayerRegistryTransactor creates a new write-only instance of PayerRegistry, bound to a specific deployed contract.
func NewPayerRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*PayerRegistryTransactor, error) {
	contract, err := bindPayerRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PayerRegistryTransactor{contract: contract}, nil
}

// NewPayerRegistryFilterer creates a new log filterer instance of PayerRegistry, bound to a specific deployed contract.
func NewPayerRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*PayerRegistryFilterer, error) {
	contract, err := bindPayerRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PayerRegistryFilterer{contract: contract}, nil
}

// bindPayerRegistry binds a generic wrapper to an already deployed contract.
func bindPayerRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PayerRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PayerRegistry *PayerRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PayerRegistry.Contract.PayerRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PayerRegistry *PayerRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayerRegistry.Contract.PayerRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PayerRegistry *PayerRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PayerRegistry.Contract.PayerRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PayerRegistry *PayerRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PayerRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PayerRegistry *PayerRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayerRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PayerRegistry *PayerRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PayerRegistry.Contract.contract.Transact(opts, method, params...)
}

// Excess is a free data retrieval call binding the contract method 0x1ae2379c.
//
// Solidity: function excess() view returns(uint96 excess_)
func (_PayerRegistry *PayerRegistryCaller) Excess(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PayerRegistry.contract.Call(opts, &out, "excess")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Excess is a free data retrieval call binding the contract method 0x1ae2379c.
//
// Solidity: function excess() view returns(uint96 excess_)
func (_PayerRegistry *PayerRegistrySession) Excess() (*big.Int, error) {
	return _PayerRegistry.Contract.Excess(&_PayerRegistry.CallOpts)
}

// Excess is a free data retrieval call binding the contract method 0x1ae2379c.
//
// Solidity: function excess() view returns(uint96 excess_)
func (_PayerRegistry *PayerRegistryCallerSession) Excess() (*big.Int, error) {
	return _PayerRegistry.Contract.Excess(&_PayerRegistry.CallOpts)
}

// FeeDistributor is a free data retrieval call binding the contract method 0x0d43e8ad.
//
// Solidity: function feeDistributor() view returns(address feeDistributor_)
func (_PayerRegistry *PayerRegistryCaller) FeeDistributor(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PayerRegistry.contract.Call(opts, &out, "feeDistributor")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeDistributor is a free data retrieval call binding the contract method 0x0d43e8ad.
//
// Solidity: function feeDistributor() view returns(address feeDistributor_)
func (_PayerRegistry *PayerRegistrySession) FeeDistributor() (common.Address, error) {
	return _PayerRegistry.Contract.FeeDistributor(&_PayerRegistry.CallOpts)
}

// FeeDistributor is a free data retrieval call binding the contract method 0x0d43e8ad.
//
// Solidity: function feeDistributor() view returns(address feeDistributor_)
func (_PayerRegistry *PayerRegistryCallerSession) FeeDistributor() (common.Address, error) {
	return _PayerRegistry.Contract.FeeDistributor(&_PayerRegistry.CallOpts)
}

// FeeDistributorParameterKey is a free data retrieval call binding the contract method 0x18daabc5.
//
// Solidity: function feeDistributorParameterKey() pure returns(bytes key_)
func (_PayerRegistry *PayerRegistryCaller) FeeDistributorParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _PayerRegistry.contract.Call(opts, &out, "feeDistributorParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// FeeDistributorParameterKey is a free data retrieval call binding the contract method 0x18daabc5.
//
// Solidity: function feeDistributorParameterKey() pure returns(bytes key_)
func (_PayerRegistry *PayerRegistrySession) FeeDistributorParameterKey() ([]byte, error) {
	return _PayerRegistry.Contract.FeeDistributorParameterKey(&_PayerRegistry.CallOpts)
}

// FeeDistributorParameterKey is a free data retrieval call binding the contract method 0x18daabc5.
//
// Solidity: function feeDistributorParameterKey() pure returns(bytes key_)
func (_PayerRegistry *PayerRegistryCallerSession) FeeDistributorParameterKey() ([]byte, error) {
	return _PayerRegistry.Contract.FeeDistributorParameterKey(&_PayerRegistry.CallOpts)
}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(address payer_) view returns(int104 balance_)
func (_PayerRegistry *PayerRegistryCaller) GetBalance(opts *bind.CallOpts, payer_ common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PayerRegistry.contract.Call(opts, &out, "getBalance", payer_)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(address payer_) view returns(int104 balance_)
func (_PayerRegistry *PayerRegistrySession) GetBalance(payer_ common.Address) (*big.Int, error) {
	return _PayerRegistry.Contract.GetBalance(&_PayerRegistry.CallOpts, payer_)
}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(address payer_) view returns(int104 balance_)
func (_PayerRegistry *PayerRegistryCallerSession) GetBalance(payer_ common.Address) (*big.Int, error) {
	return _PayerRegistry.Contract.GetBalance(&_PayerRegistry.CallOpts, payer_)
}

// GetBalances is a free data retrieval call binding the contract method 0x2d2ae1c1.
//
// Solidity: function getBalances(address[] payers_) view returns(int104[] balances_)
func (_PayerRegistry *PayerRegistryCaller) GetBalances(opts *bind.CallOpts, payers_ []common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _PayerRegistry.contract.Call(opts, &out, "getBalances", payers_)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetBalances is a free data retrieval call binding the contract method 0x2d2ae1c1.
//
// Solidity: function getBalances(address[] payers_) view returns(int104[] balances_)
func (_PayerRegistry *PayerRegistrySession) GetBalances(payers_ []common.Address) ([]*big.Int, error) {
	return _PayerRegistry.Contract.GetBalances(&_PayerRegistry.CallOpts, payers_)
}

// GetBalances is a free data retrieval call binding the contract method 0x2d2ae1c1.
//
// Solidity: function getBalances(address[] payers_) view returns(int104[] balances_)
func (_PayerRegistry *PayerRegistryCallerSession) GetBalances(payers_ []common.Address) ([]*big.Int, error) {
	return _PayerRegistry.Contract.GetBalances(&_PayerRegistry.CallOpts, payers_)
}

// GetPendingWithdrawal is a free data retrieval call binding the contract method 0x7ee8b2f8.
//
// Solidity: function getPendingWithdrawal(address payer_) view returns(uint96 pendingWithdrawal_, uint32 withdrawableTimestamp_)
func (_PayerRegistry *PayerRegistryCaller) GetPendingWithdrawal(opts *bind.CallOpts, payer_ common.Address) (struct {
	PendingWithdrawal     *big.Int
	WithdrawableTimestamp uint32
}, error) {
	var out []interface{}
	err := _PayerRegistry.contract.Call(opts, &out, "getPendingWithdrawal", payer_)

	outstruct := new(struct {
		PendingWithdrawal     *big.Int
		WithdrawableTimestamp uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.PendingWithdrawal = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.WithdrawableTimestamp = *abi.ConvertType(out[1], new(uint32)).(*uint32)

	return *outstruct, err

}

// GetPendingWithdrawal is a free data retrieval call binding the contract method 0x7ee8b2f8.
//
// Solidity: function getPendingWithdrawal(address payer_) view returns(uint96 pendingWithdrawal_, uint32 withdrawableTimestamp_)
func (_PayerRegistry *PayerRegistrySession) GetPendingWithdrawal(payer_ common.Address) (struct {
	PendingWithdrawal     *big.Int
	WithdrawableTimestamp uint32
}, error) {
	return _PayerRegistry.Contract.GetPendingWithdrawal(&_PayerRegistry.CallOpts, payer_)
}

// GetPendingWithdrawal is a free data retrieval call binding the contract method 0x7ee8b2f8.
//
// Solidity: function getPendingWithdrawal(address payer_) view returns(uint96 pendingWithdrawal_, uint32 withdrawableTimestamp_)
func (_PayerRegistry *PayerRegistryCallerSession) GetPendingWithdrawal(payer_ common.Address) (struct {
	PendingWithdrawal     *big.Int
	WithdrawableTimestamp uint32
}, error) {
	return _PayerRegistry.Contract.GetPendingWithdrawal(&_PayerRegistry.CallOpts, payer_)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_PayerRegistry *PayerRegistryCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PayerRegistry.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_PayerRegistry *PayerRegistrySession) Implementation() (common.Address, error) {
	return _PayerRegistry.Contract.Implementation(&_PayerRegistry.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_PayerRegistry *PayerRegistryCallerSession) Implementation() (common.Address, error) {
	return _PayerRegistry.Contract.Implementation(&_PayerRegistry.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(bytes key_)
func (_PayerRegistry *PayerRegistryCaller) MigratorParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _PayerRegistry.contract.Call(opts, &out, "migratorParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(bytes key_)
func (_PayerRegistry *PayerRegistrySession) MigratorParameterKey() ([]byte, error) {
	return _PayerRegistry.Contract.MigratorParameterKey(&_PayerRegistry.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(bytes key_)
func (_PayerRegistry *PayerRegistryCallerSession) MigratorParameterKey() ([]byte, error) {
	return _PayerRegistry.Contract.MigratorParameterKey(&_PayerRegistry.CallOpts)
}

// MinimumDeposit is a free data retrieval call binding the contract method 0x636bfbab.
//
// Solidity: function minimumDeposit() view returns(uint96 minimumDeposit_)
func (_PayerRegistry *PayerRegistryCaller) MinimumDeposit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PayerRegistry.contract.Call(opts, &out, "minimumDeposit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinimumDeposit is a free data retrieval call binding the contract method 0x636bfbab.
//
// Solidity: function minimumDeposit() view returns(uint96 minimumDeposit_)
func (_PayerRegistry *PayerRegistrySession) MinimumDeposit() (*big.Int, error) {
	return _PayerRegistry.Contract.MinimumDeposit(&_PayerRegistry.CallOpts)
}

// MinimumDeposit is a free data retrieval call binding the contract method 0x636bfbab.
//
// Solidity: function minimumDeposit() view returns(uint96 minimumDeposit_)
func (_PayerRegistry *PayerRegistryCallerSession) MinimumDeposit() (*big.Int, error) {
	return _PayerRegistry.Contract.MinimumDeposit(&_PayerRegistry.CallOpts)
}

// MinimumDepositParameterKey is a free data retrieval call binding the contract method 0x7a303d33.
//
// Solidity: function minimumDepositParameterKey() pure returns(bytes key_)
func (_PayerRegistry *PayerRegistryCaller) MinimumDepositParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _PayerRegistry.contract.Call(opts, &out, "minimumDepositParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// MinimumDepositParameterKey is a free data retrieval call binding the contract method 0x7a303d33.
//
// Solidity: function minimumDepositParameterKey() pure returns(bytes key_)
func (_PayerRegistry *PayerRegistrySession) MinimumDepositParameterKey() ([]byte, error) {
	return _PayerRegistry.Contract.MinimumDepositParameterKey(&_PayerRegistry.CallOpts)
}

// MinimumDepositParameterKey is a free data retrieval call binding the contract method 0x7a303d33.
//
// Solidity: function minimumDepositParameterKey() pure returns(bytes key_)
func (_PayerRegistry *PayerRegistryCallerSession) MinimumDepositParameterKey() ([]byte, error) {
	return _PayerRegistry.Contract.MinimumDepositParameterKey(&_PayerRegistry.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_PayerRegistry *PayerRegistryCaller) ParameterRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PayerRegistry.contract.Call(opts, &out, "parameterRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_PayerRegistry *PayerRegistrySession) ParameterRegistry() (common.Address, error) {
	return _PayerRegistry.Contract.ParameterRegistry(&_PayerRegistry.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_PayerRegistry *PayerRegistryCallerSession) ParameterRegistry() (common.Address, error) {
	return _PayerRegistry.Contract.ParameterRegistry(&_PayerRegistry.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool paused_)
func (_PayerRegistry *PayerRegistryCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _PayerRegistry.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool paused_)
func (_PayerRegistry *PayerRegistrySession) Paused() (bool, error) {
	return _PayerRegistry.Contract.Paused(&_PayerRegistry.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool paused_)
func (_PayerRegistry *PayerRegistryCallerSession) Paused() (bool, error) {
	return _PayerRegistry.Contract.Paused(&_PayerRegistry.CallOpts)
}

// PausedParameterKey is a free data retrieval call binding the contract method 0xcc5999af.
//
// Solidity: function pausedParameterKey() pure returns(bytes key_)
func (_PayerRegistry *PayerRegistryCaller) PausedParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _PayerRegistry.contract.Call(opts, &out, "pausedParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// PausedParameterKey is a free data retrieval call binding the contract method 0xcc5999af.
//
// Solidity: function pausedParameterKey() pure returns(bytes key_)
func (_PayerRegistry *PayerRegistrySession) PausedParameterKey() ([]byte, error) {
	return _PayerRegistry.Contract.PausedParameterKey(&_PayerRegistry.CallOpts)
}

// PausedParameterKey is a free data retrieval call binding the contract method 0xcc5999af.
//
// Solidity: function pausedParameterKey() pure returns(bytes key_)
func (_PayerRegistry *PayerRegistryCallerSession) PausedParameterKey() ([]byte, error) {
	return _PayerRegistry.Contract.PausedParameterKey(&_PayerRegistry.CallOpts)
}

// Settler is a free data retrieval call binding the contract method 0xab221a76.
//
// Solidity: function settler() view returns(address settler_)
func (_PayerRegistry *PayerRegistryCaller) Settler(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PayerRegistry.contract.Call(opts, &out, "settler")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Settler is a free data retrieval call binding the contract method 0xab221a76.
//
// Solidity: function settler() view returns(address settler_)
func (_PayerRegistry *PayerRegistrySession) Settler() (common.Address, error) {
	return _PayerRegistry.Contract.Settler(&_PayerRegistry.CallOpts)
}

// Settler is a free data retrieval call binding the contract method 0xab221a76.
//
// Solidity: function settler() view returns(address settler_)
func (_PayerRegistry *PayerRegistryCallerSession) Settler() (common.Address, error) {
	return _PayerRegistry.Contract.Settler(&_PayerRegistry.CallOpts)
}

// SettlerParameterKey is a free data retrieval call binding the contract method 0x028d1638.
//
// Solidity: function settlerParameterKey() pure returns(bytes key_)
func (_PayerRegistry *PayerRegistryCaller) SettlerParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _PayerRegistry.contract.Call(opts, &out, "settlerParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// SettlerParameterKey is a free data retrieval call binding the contract method 0x028d1638.
//
// Solidity: function settlerParameterKey() pure returns(bytes key_)
func (_PayerRegistry *PayerRegistrySession) SettlerParameterKey() ([]byte, error) {
	return _PayerRegistry.Contract.SettlerParameterKey(&_PayerRegistry.CallOpts)
}

// SettlerParameterKey is a free data retrieval call binding the contract method 0x028d1638.
//
// Solidity: function settlerParameterKey() pure returns(bytes key_)
func (_PayerRegistry *PayerRegistryCallerSession) SettlerParameterKey() ([]byte, error) {
	return _PayerRegistry.Contract.SettlerParameterKey(&_PayerRegistry.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_PayerRegistry *PayerRegistryCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PayerRegistry.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_PayerRegistry *PayerRegistrySession) Token() (common.Address, error) {
	return _PayerRegistry.Contract.Token(&_PayerRegistry.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_PayerRegistry *PayerRegistryCallerSession) Token() (common.Address, error) {
	return _PayerRegistry.Contract.Token(&_PayerRegistry.CallOpts)
}

// TotalDebt is a free data retrieval call binding the contract method 0xfc7b9c18.
//
// Solidity: function totalDebt() view returns(uint96 totalDebt_)
func (_PayerRegistry *PayerRegistryCaller) TotalDebt(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PayerRegistry.contract.Call(opts, &out, "totalDebt")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalDebt is a free data retrieval call binding the contract method 0xfc7b9c18.
//
// Solidity: function totalDebt() view returns(uint96 totalDebt_)
func (_PayerRegistry *PayerRegistrySession) TotalDebt() (*big.Int, error) {
	return _PayerRegistry.Contract.TotalDebt(&_PayerRegistry.CallOpts)
}

// TotalDebt is a free data retrieval call binding the contract method 0xfc7b9c18.
//
// Solidity: function totalDebt() view returns(uint96 totalDebt_)
func (_PayerRegistry *PayerRegistryCallerSession) TotalDebt() (*big.Int, error) {
	return _PayerRegistry.Contract.TotalDebt(&_PayerRegistry.CallOpts)
}

// TotalDeposits is a free data retrieval call binding the contract method 0x7d882097.
//
// Solidity: function totalDeposits() view returns(int104 totalDeposits_)
func (_PayerRegistry *PayerRegistryCaller) TotalDeposits(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PayerRegistry.contract.Call(opts, &out, "totalDeposits")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalDeposits is a free data retrieval call binding the contract method 0x7d882097.
//
// Solidity: function totalDeposits() view returns(int104 totalDeposits_)
func (_PayerRegistry *PayerRegistrySession) TotalDeposits() (*big.Int, error) {
	return _PayerRegistry.Contract.TotalDeposits(&_PayerRegistry.CallOpts)
}

// TotalDeposits is a free data retrieval call binding the contract method 0x7d882097.
//
// Solidity: function totalDeposits() view returns(int104 totalDeposits_)
func (_PayerRegistry *PayerRegistryCallerSession) TotalDeposits() (*big.Int, error) {
	return _PayerRegistry.Contract.TotalDeposits(&_PayerRegistry.CallOpts)
}

// TotalWithdrawable is a free data retrieval call binding the contract method 0x0600a865.
//
// Solidity: function totalWithdrawable() view returns(uint96 totalWithdrawable_)
func (_PayerRegistry *PayerRegistryCaller) TotalWithdrawable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PayerRegistry.contract.Call(opts, &out, "totalWithdrawable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalWithdrawable is a free data retrieval call binding the contract method 0x0600a865.
//
// Solidity: function totalWithdrawable() view returns(uint96 totalWithdrawable_)
func (_PayerRegistry *PayerRegistrySession) TotalWithdrawable() (*big.Int, error) {
	return _PayerRegistry.Contract.TotalWithdrawable(&_PayerRegistry.CallOpts)
}

// TotalWithdrawable is a free data retrieval call binding the contract method 0x0600a865.
//
// Solidity: function totalWithdrawable() view returns(uint96 totalWithdrawable_)
func (_PayerRegistry *PayerRegistryCallerSession) TotalWithdrawable() (*big.Int, error) {
	return _PayerRegistry.Contract.TotalWithdrawable(&_PayerRegistry.CallOpts)
}

// WithdrawLockPeriod is a free data retrieval call binding the contract method 0x2628490f.
//
// Solidity: function withdrawLockPeriod() view returns(uint32 withdrawLockPeriod_)
func (_PayerRegistry *PayerRegistryCaller) WithdrawLockPeriod(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _PayerRegistry.contract.Call(opts, &out, "withdrawLockPeriod")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// WithdrawLockPeriod is a free data retrieval call binding the contract method 0x2628490f.
//
// Solidity: function withdrawLockPeriod() view returns(uint32 withdrawLockPeriod_)
func (_PayerRegistry *PayerRegistrySession) WithdrawLockPeriod() (uint32, error) {
	return _PayerRegistry.Contract.WithdrawLockPeriod(&_PayerRegistry.CallOpts)
}

// WithdrawLockPeriod is a free data retrieval call binding the contract method 0x2628490f.
//
// Solidity: function withdrawLockPeriod() view returns(uint32 withdrawLockPeriod_)
func (_PayerRegistry *PayerRegistryCallerSession) WithdrawLockPeriod() (uint32, error) {
	return _PayerRegistry.Contract.WithdrawLockPeriod(&_PayerRegistry.CallOpts)
}

// WithdrawLockPeriodParameterKey is a free data retrieval call binding the contract method 0xaf1a77c5.
//
// Solidity: function withdrawLockPeriodParameterKey() pure returns(bytes key_)
func (_PayerRegistry *PayerRegistryCaller) WithdrawLockPeriodParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _PayerRegistry.contract.Call(opts, &out, "withdrawLockPeriodParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// WithdrawLockPeriodParameterKey is a free data retrieval call binding the contract method 0xaf1a77c5.
//
// Solidity: function withdrawLockPeriodParameterKey() pure returns(bytes key_)
func (_PayerRegistry *PayerRegistrySession) WithdrawLockPeriodParameterKey() ([]byte, error) {
	return _PayerRegistry.Contract.WithdrawLockPeriodParameterKey(&_PayerRegistry.CallOpts)
}

// WithdrawLockPeriodParameterKey is a free data retrieval call binding the contract method 0xaf1a77c5.
//
// Solidity: function withdrawLockPeriodParameterKey() pure returns(bytes key_)
func (_PayerRegistry *PayerRegistryCallerSession) WithdrawLockPeriodParameterKey() ([]byte, error) {
	return _PayerRegistry.Contract.WithdrawLockPeriodParameterKey(&_PayerRegistry.CallOpts)
}

// CancelWithdrawal is a paid mutator transaction binding the contract method 0x22611280.
//
// Solidity: function cancelWithdrawal() returns()
func (_PayerRegistry *PayerRegistryTransactor) CancelWithdrawal(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayerRegistry.contract.Transact(opts, "cancelWithdrawal")
}

// CancelWithdrawal is a paid mutator transaction binding the contract method 0x22611280.
//
// Solidity: function cancelWithdrawal() returns()
func (_PayerRegistry *PayerRegistrySession) CancelWithdrawal() (*types.Transaction, error) {
	return _PayerRegistry.Contract.CancelWithdrawal(&_PayerRegistry.TransactOpts)
}

// CancelWithdrawal is a paid mutator transaction binding the contract method 0x22611280.
//
// Solidity: function cancelWithdrawal() returns()
func (_PayerRegistry *PayerRegistryTransactorSession) CancelWithdrawal() (*types.Transaction, error) {
	return _PayerRegistry.Contract.CancelWithdrawal(&_PayerRegistry.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0x3ae50b73.
//
// Solidity: function deposit(address payer_, uint96 amount_) returns()
func (_PayerRegistry *PayerRegistryTransactor) Deposit(opts *bind.TransactOpts, payer_ common.Address, amount_ *big.Int) (*types.Transaction, error) {
	return _PayerRegistry.contract.Transact(opts, "deposit", payer_, amount_)
}

// Deposit is a paid mutator transaction binding the contract method 0x3ae50b73.
//
// Solidity: function deposit(address payer_, uint96 amount_) returns()
func (_PayerRegistry *PayerRegistrySession) Deposit(payer_ common.Address, amount_ *big.Int) (*types.Transaction, error) {
	return _PayerRegistry.Contract.Deposit(&_PayerRegistry.TransactOpts, payer_, amount_)
}

// Deposit is a paid mutator transaction binding the contract method 0x3ae50b73.
//
// Solidity: function deposit(address payer_, uint96 amount_) returns()
func (_PayerRegistry *PayerRegistryTransactorSession) Deposit(payer_ common.Address, amount_ *big.Int) (*types.Transaction, error) {
	return _PayerRegistry.Contract.Deposit(&_PayerRegistry.TransactOpts, payer_, amount_)
}

// DepositWithPermit is a paid mutator transaction binding the contract method 0x7e9b9b18.
//
// Solidity: function depositWithPermit(address payer_, uint96 amount_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns()
func (_PayerRegistry *PayerRegistryTransactor) DepositWithPermit(opts *bind.TransactOpts, payer_ common.Address, amount_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _PayerRegistry.contract.Transact(opts, "depositWithPermit", payer_, amount_, deadline_, v_, r_, s_)
}

// DepositWithPermit is a paid mutator transaction binding the contract method 0x7e9b9b18.
//
// Solidity: function depositWithPermit(address payer_, uint96 amount_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns()
func (_PayerRegistry *PayerRegistrySession) DepositWithPermit(payer_ common.Address, amount_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _PayerRegistry.Contract.DepositWithPermit(&_PayerRegistry.TransactOpts, payer_, amount_, deadline_, v_, r_, s_)
}

// DepositWithPermit is a paid mutator transaction binding the contract method 0x7e9b9b18.
//
// Solidity: function depositWithPermit(address payer_, uint96 amount_, uint256 deadline_, uint8 v_, bytes32 r_, bytes32 s_) returns()
func (_PayerRegistry *PayerRegistryTransactorSession) DepositWithPermit(payer_ common.Address, amount_ *big.Int, deadline_ *big.Int, v_ uint8, r_ [32]byte, s_ [32]byte) (*types.Transaction, error) {
	return _PayerRegistry.Contract.DepositWithPermit(&_PayerRegistry.TransactOpts, payer_, amount_, deadline_, v_, r_, s_)
}

// FinalizeWithdrawal is a paid mutator transaction binding the contract method 0x4abf24cb.
//
// Solidity: function finalizeWithdrawal(address recipient_) returns()
func (_PayerRegistry *PayerRegistryTransactor) FinalizeWithdrawal(opts *bind.TransactOpts, recipient_ common.Address) (*types.Transaction, error) {
	return _PayerRegistry.contract.Transact(opts, "finalizeWithdrawal", recipient_)
}

// FinalizeWithdrawal is a paid mutator transaction binding the contract method 0x4abf24cb.
//
// Solidity: function finalizeWithdrawal(address recipient_) returns()
func (_PayerRegistry *PayerRegistrySession) FinalizeWithdrawal(recipient_ common.Address) (*types.Transaction, error) {
	return _PayerRegistry.Contract.FinalizeWithdrawal(&_PayerRegistry.TransactOpts, recipient_)
}

// FinalizeWithdrawal is a paid mutator transaction binding the contract method 0x4abf24cb.
//
// Solidity: function finalizeWithdrawal(address recipient_) returns()
func (_PayerRegistry *PayerRegistryTransactorSession) FinalizeWithdrawal(recipient_ common.Address) (*types.Transaction, error) {
	return _PayerRegistry.Contract.FinalizeWithdrawal(&_PayerRegistry.TransactOpts, recipient_)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_PayerRegistry *PayerRegistryTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayerRegistry.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_PayerRegistry *PayerRegistrySession) Initialize() (*types.Transaction, error) {
	return _PayerRegistry.Contract.Initialize(&_PayerRegistry.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_PayerRegistry *PayerRegistryTransactorSession) Initialize() (*types.Transaction, error) {
	return _PayerRegistry.Contract.Initialize(&_PayerRegistry.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_PayerRegistry *PayerRegistryTransactor) Migrate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayerRegistry.contract.Transact(opts, "migrate")
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_PayerRegistry *PayerRegistrySession) Migrate() (*types.Transaction, error) {
	return _PayerRegistry.Contract.Migrate(&_PayerRegistry.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_PayerRegistry *PayerRegistryTransactorSession) Migrate() (*types.Transaction, error) {
	return _PayerRegistry.Contract.Migrate(&_PayerRegistry.TransactOpts)
}

// RequestWithdrawal is a paid mutator transaction binding the contract method 0xb15780a0.
//
// Solidity: function requestWithdrawal(uint96 amount_) returns()
func (_PayerRegistry *PayerRegistryTransactor) RequestWithdrawal(opts *bind.TransactOpts, amount_ *big.Int) (*types.Transaction, error) {
	return _PayerRegistry.contract.Transact(opts, "requestWithdrawal", amount_)
}

// RequestWithdrawal is a paid mutator transaction binding the contract method 0xb15780a0.
//
// Solidity: function requestWithdrawal(uint96 amount_) returns()
func (_PayerRegistry *PayerRegistrySession) RequestWithdrawal(amount_ *big.Int) (*types.Transaction, error) {
	return _PayerRegistry.Contract.RequestWithdrawal(&_PayerRegistry.TransactOpts, amount_)
}

// RequestWithdrawal is a paid mutator transaction binding the contract method 0xb15780a0.
//
// Solidity: function requestWithdrawal(uint96 amount_) returns()
func (_PayerRegistry *PayerRegistryTransactorSession) RequestWithdrawal(amount_ *big.Int) (*types.Transaction, error) {
	return _PayerRegistry.Contract.RequestWithdrawal(&_PayerRegistry.TransactOpts, amount_)
}

// SendExcessToFeeDistributor is a paid mutator transaction binding the contract method 0xe2813ab4.
//
// Solidity: function sendExcessToFeeDistributor() returns(uint96 excess_)
func (_PayerRegistry *PayerRegistryTransactor) SendExcessToFeeDistributor(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayerRegistry.contract.Transact(opts, "sendExcessToFeeDistributor")
}

// SendExcessToFeeDistributor is a paid mutator transaction binding the contract method 0xe2813ab4.
//
// Solidity: function sendExcessToFeeDistributor() returns(uint96 excess_)
func (_PayerRegistry *PayerRegistrySession) SendExcessToFeeDistributor() (*types.Transaction, error) {
	return _PayerRegistry.Contract.SendExcessToFeeDistributor(&_PayerRegistry.TransactOpts)
}

// SendExcessToFeeDistributor is a paid mutator transaction binding the contract method 0xe2813ab4.
//
// Solidity: function sendExcessToFeeDistributor() returns(uint96 excess_)
func (_PayerRegistry *PayerRegistryTransactorSession) SendExcessToFeeDistributor() (*types.Transaction, error) {
	return _PayerRegistry.Contract.SendExcessToFeeDistributor(&_PayerRegistry.TransactOpts)
}

// SettleUsage is a paid mutator transaction binding the contract method 0xd219be1b.
//
// Solidity: function settleUsage((address,uint96)[] payerFees_) returns(uint96 feesSettled_)
func (_PayerRegistry *PayerRegistryTransactor) SettleUsage(opts *bind.TransactOpts, payerFees_ []IPayerRegistryPayerFee) (*types.Transaction, error) {
	return _PayerRegistry.contract.Transact(opts, "settleUsage", payerFees_)
}

// SettleUsage is a paid mutator transaction binding the contract method 0xd219be1b.
//
// Solidity: function settleUsage((address,uint96)[] payerFees_) returns(uint96 feesSettled_)
func (_PayerRegistry *PayerRegistrySession) SettleUsage(payerFees_ []IPayerRegistryPayerFee) (*types.Transaction, error) {
	return _PayerRegistry.Contract.SettleUsage(&_PayerRegistry.TransactOpts, payerFees_)
}

// SettleUsage is a paid mutator transaction binding the contract method 0xd219be1b.
//
// Solidity: function settleUsage((address,uint96)[] payerFees_) returns(uint96 feesSettled_)
func (_PayerRegistry *PayerRegistryTransactorSession) SettleUsage(payerFees_ []IPayerRegistryPayerFee) (*types.Transaction, error) {
	return _PayerRegistry.Contract.SettleUsage(&_PayerRegistry.TransactOpts, payerFees_)
}

// UpdateFeeDistributor is a paid mutator transaction binding the contract method 0xdd215de7.
//
// Solidity: function updateFeeDistributor() returns()
func (_PayerRegistry *PayerRegistryTransactor) UpdateFeeDistributor(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayerRegistry.contract.Transact(opts, "updateFeeDistributor")
}

// UpdateFeeDistributor is a paid mutator transaction binding the contract method 0xdd215de7.
//
// Solidity: function updateFeeDistributor() returns()
func (_PayerRegistry *PayerRegistrySession) UpdateFeeDistributor() (*types.Transaction, error) {
	return _PayerRegistry.Contract.UpdateFeeDistributor(&_PayerRegistry.TransactOpts)
}

// UpdateFeeDistributor is a paid mutator transaction binding the contract method 0xdd215de7.
//
// Solidity: function updateFeeDistributor() returns()
func (_PayerRegistry *PayerRegistryTransactorSession) UpdateFeeDistributor() (*types.Transaction, error) {
	return _PayerRegistry.Contract.UpdateFeeDistributor(&_PayerRegistry.TransactOpts)
}

// UpdateMinimumDeposit is a paid mutator transaction binding the contract method 0xd490fae9.
//
// Solidity: function updateMinimumDeposit() returns()
func (_PayerRegistry *PayerRegistryTransactor) UpdateMinimumDeposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayerRegistry.contract.Transact(opts, "updateMinimumDeposit")
}

// UpdateMinimumDeposit is a paid mutator transaction binding the contract method 0xd490fae9.
//
// Solidity: function updateMinimumDeposit() returns()
func (_PayerRegistry *PayerRegistrySession) UpdateMinimumDeposit() (*types.Transaction, error) {
	return _PayerRegistry.Contract.UpdateMinimumDeposit(&_PayerRegistry.TransactOpts)
}

// UpdateMinimumDeposit is a paid mutator transaction binding the contract method 0xd490fae9.
//
// Solidity: function updateMinimumDeposit() returns()
func (_PayerRegistry *PayerRegistryTransactorSession) UpdateMinimumDeposit() (*types.Transaction, error) {
	return _PayerRegistry.Contract.UpdateMinimumDeposit(&_PayerRegistry.TransactOpts)
}

// UpdatePauseStatus is a paid mutator transaction binding the contract method 0x59d4df41.
//
// Solidity: function updatePauseStatus() returns()
func (_PayerRegistry *PayerRegistryTransactor) UpdatePauseStatus(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayerRegistry.contract.Transact(opts, "updatePauseStatus")
}

// UpdatePauseStatus is a paid mutator transaction binding the contract method 0x59d4df41.
//
// Solidity: function updatePauseStatus() returns()
func (_PayerRegistry *PayerRegistrySession) UpdatePauseStatus() (*types.Transaction, error) {
	return _PayerRegistry.Contract.UpdatePauseStatus(&_PayerRegistry.TransactOpts)
}

// UpdatePauseStatus is a paid mutator transaction binding the contract method 0x59d4df41.
//
// Solidity: function updatePauseStatus() returns()
func (_PayerRegistry *PayerRegistryTransactorSession) UpdatePauseStatus() (*types.Transaction, error) {
	return _PayerRegistry.Contract.UpdatePauseStatus(&_PayerRegistry.TransactOpts)
}

// UpdateSettler is a paid mutator transaction binding the contract method 0x9d5619da.
//
// Solidity: function updateSettler() returns()
func (_PayerRegistry *PayerRegistryTransactor) UpdateSettler(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayerRegistry.contract.Transact(opts, "updateSettler")
}

// UpdateSettler is a paid mutator transaction binding the contract method 0x9d5619da.
//
// Solidity: function updateSettler() returns()
func (_PayerRegistry *PayerRegistrySession) UpdateSettler() (*types.Transaction, error) {
	return _PayerRegistry.Contract.UpdateSettler(&_PayerRegistry.TransactOpts)
}

// UpdateSettler is a paid mutator transaction binding the contract method 0x9d5619da.
//
// Solidity: function updateSettler() returns()
func (_PayerRegistry *PayerRegistryTransactorSession) UpdateSettler() (*types.Transaction, error) {
	return _PayerRegistry.Contract.UpdateSettler(&_PayerRegistry.TransactOpts)
}

// UpdateWithdrawLockPeriod is a paid mutator transaction binding the contract method 0xd384ee35.
//
// Solidity: function updateWithdrawLockPeriod() returns()
func (_PayerRegistry *PayerRegistryTransactor) UpdateWithdrawLockPeriod(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayerRegistry.contract.Transact(opts, "updateWithdrawLockPeriod")
}

// UpdateWithdrawLockPeriod is a paid mutator transaction binding the contract method 0xd384ee35.
//
// Solidity: function updateWithdrawLockPeriod() returns()
func (_PayerRegistry *PayerRegistrySession) UpdateWithdrawLockPeriod() (*types.Transaction, error) {
	return _PayerRegistry.Contract.UpdateWithdrawLockPeriod(&_PayerRegistry.TransactOpts)
}

// UpdateWithdrawLockPeriod is a paid mutator transaction binding the contract method 0xd384ee35.
//
// Solidity: function updateWithdrawLockPeriod() returns()
func (_PayerRegistry *PayerRegistryTransactorSession) UpdateWithdrawLockPeriod() (*types.Transaction, error) {
	return _PayerRegistry.Contract.UpdateWithdrawLockPeriod(&_PayerRegistry.TransactOpts)
}

// PayerRegistryDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the PayerRegistry contract.
type PayerRegistryDepositIterator struct {
	Event *PayerRegistryDeposit // Event containing the contract specifics and raw log

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
func (it *PayerRegistryDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayerRegistryDeposit)
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
		it.Event = new(PayerRegistryDeposit)
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
func (it *PayerRegistryDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayerRegistryDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayerRegistryDeposit represents a Deposit event raised by the PayerRegistry contract.
type PayerRegistryDeposit struct {
	Payer  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x5047c753a53960392b00d7af1a52e5e9ddfba5fd85f8c61391736813f9ec7e29.
//
// Solidity: event Deposit(address indexed payer, uint96 amount)
func (_PayerRegistry *PayerRegistryFilterer) FilterDeposit(opts *bind.FilterOpts, payer []common.Address) (*PayerRegistryDepositIterator, error) {

	var payerRule []interface{}
	for _, payerItem := range payer {
		payerRule = append(payerRule, payerItem)
	}

	logs, sub, err := _PayerRegistry.contract.FilterLogs(opts, "Deposit", payerRule)
	if err != nil {
		return nil, err
	}
	return &PayerRegistryDepositIterator{contract: _PayerRegistry.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x5047c753a53960392b00d7af1a52e5e9ddfba5fd85f8c61391736813f9ec7e29.
//
// Solidity: event Deposit(address indexed payer, uint96 amount)
func (_PayerRegistry *PayerRegistryFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *PayerRegistryDeposit, payer []common.Address) (event.Subscription, error) {

	var payerRule []interface{}
	for _, payerItem := range payer {
		payerRule = append(payerRule, payerItem)
	}

	logs, sub, err := _PayerRegistry.contract.WatchLogs(opts, "Deposit", payerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayerRegistryDeposit)
				if err := _PayerRegistry.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0x5047c753a53960392b00d7af1a52e5e9ddfba5fd85f8c61391736813f9ec7e29.
//
// Solidity: event Deposit(address indexed payer, uint96 amount)
func (_PayerRegistry *PayerRegistryFilterer) ParseDeposit(log types.Log) (*PayerRegistryDeposit, error) {
	event := new(PayerRegistryDeposit)
	if err := _PayerRegistry.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayerRegistryExcessTransferredIterator is returned from FilterExcessTransferred and is used to iterate over the raw logs and unpacked data for ExcessTransferred events raised by the PayerRegistry contract.
type PayerRegistryExcessTransferredIterator struct {
	Event *PayerRegistryExcessTransferred // Event containing the contract specifics and raw log

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
func (it *PayerRegistryExcessTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayerRegistryExcessTransferred)
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
		it.Event = new(PayerRegistryExcessTransferred)
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
func (it *PayerRegistryExcessTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayerRegistryExcessTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayerRegistryExcessTransferred represents a ExcessTransferred event raised by the PayerRegistry contract.
type PayerRegistryExcessTransferred struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterExcessTransferred is a free log retrieval operation binding the contract event 0x5fbb96aee4bddf19c82360ce325a9d8f3fbb84374f5466eb466d4fe93a768d02.
//
// Solidity: event ExcessTransferred(uint96 amount)
func (_PayerRegistry *PayerRegistryFilterer) FilterExcessTransferred(opts *bind.FilterOpts) (*PayerRegistryExcessTransferredIterator, error) {

	logs, sub, err := _PayerRegistry.contract.FilterLogs(opts, "ExcessTransferred")
	if err != nil {
		return nil, err
	}
	return &PayerRegistryExcessTransferredIterator{contract: _PayerRegistry.contract, event: "ExcessTransferred", logs: logs, sub: sub}, nil
}

// WatchExcessTransferred is a free log subscription operation binding the contract event 0x5fbb96aee4bddf19c82360ce325a9d8f3fbb84374f5466eb466d4fe93a768d02.
//
// Solidity: event ExcessTransferred(uint96 amount)
func (_PayerRegistry *PayerRegistryFilterer) WatchExcessTransferred(opts *bind.WatchOpts, sink chan<- *PayerRegistryExcessTransferred) (event.Subscription, error) {

	logs, sub, err := _PayerRegistry.contract.WatchLogs(opts, "ExcessTransferred")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayerRegistryExcessTransferred)
				if err := _PayerRegistry.contract.UnpackLog(event, "ExcessTransferred", log); err != nil {
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

// ParseExcessTransferred is a log parse operation binding the contract event 0x5fbb96aee4bddf19c82360ce325a9d8f3fbb84374f5466eb466d4fe93a768d02.
//
// Solidity: event ExcessTransferred(uint96 amount)
func (_PayerRegistry *PayerRegistryFilterer) ParseExcessTransferred(log types.Log) (*PayerRegistryExcessTransferred, error) {
	event := new(PayerRegistryExcessTransferred)
	if err := _PayerRegistry.contract.UnpackLog(event, "ExcessTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayerRegistryFeeDistributorUpdatedIterator is returned from FilterFeeDistributorUpdated and is used to iterate over the raw logs and unpacked data for FeeDistributorUpdated events raised by the PayerRegistry contract.
type PayerRegistryFeeDistributorUpdatedIterator struct {
	Event *PayerRegistryFeeDistributorUpdated // Event containing the contract specifics and raw log

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
func (it *PayerRegistryFeeDistributorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayerRegistryFeeDistributorUpdated)
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
		it.Event = new(PayerRegistryFeeDistributorUpdated)
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
func (it *PayerRegistryFeeDistributorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayerRegistryFeeDistributorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayerRegistryFeeDistributorUpdated represents a FeeDistributorUpdated event raised by the PayerRegistry contract.
type PayerRegistryFeeDistributorUpdated struct {
	FeeDistributor common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterFeeDistributorUpdated is a free log retrieval operation binding the contract event 0x98b244c47458fedc88c0cf9958073e565e971d8ae9d8d1c50dfc6920fe939cbb.
//
// Solidity: event FeeDistributorUpdated(address indexed feeDistributor)
func (_PayerRegistry *PayerRegistryFilterer) FilterFeeDistributorUpdated(opts *bind.FilterOpts, feeDistributor []common.Address) (*PayerRegistryFeeDistributorUpdatedIterator, error) {

	var feeDistributorRule []interface{}
	for _, feeDistributorItem := range feeDistributor {
		feeDistributorRule = append(feeDistributorRule, feeDistributorItem)
	}

	logs, sub, err := _PayerRegistry.contract.FilterLogs(opts, "FeeDistributorUpdated", feeDistributorRule)
	if err != nil {
		return nil, err
	}
	return &PayerRegistryFeeDistributorUpdatedIterator{contract: _PayerRegistry.contract, event: "FeeDistributorUpdated", logs: logs, sub: sub}, nil
}

// WatchFeeDistributorUpdated is a free log subscription operation binding the contract event 0x98b244c47458fedc88c0cf9958073e565e971d8ae9d8d1c50dfc6920fe939cbb.
//
// Solidity: event FeeDistributorUpdated(address indexed feeDistributor)
func (_PayerRegistry *PayerRegistryFilterer) WatchFeeDistributorUpdated(opts *bind.WatchOpts, sink chan<- *PayerRegistryFeeDistributorUpdated, feeDistributor []common.Address) (event.Subscription, error) {

	var feeDistributorRule []interface{}
	for _, feeDistributorItem := range feeDistributor {
		feeDistributorRule = append(feeDistributorRule, feeDistributorItem)
	}

	logs, sub, err := _PayerRegistry.contract.WatchLogs(opts, "FeeDistributorUpdated", feeDistributorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayerRegistryFeeDistributorUpdated)
				if err := _PayerRegistry.contract.UnpackLog(event, "FeeDistributorUpdated", log); err != nil {
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

// ParseFeeDistributorUpdated is a log parse operation binding the contract event 0x98b244c47458fedc88c0cf9958073e565e971d8ae9d8d1c50dfc6920fe939cbb.
//
// Solidity: event FeeDistributorUpdated(address indexed feeDistributor)
func (_PayerRegistry *PayerRegistryFilterer) ParseFeeDistributorUpdated(log types.Log) (*PayerRegistryFeeDistributorUpdated, error) {
	event := new(PayerRegistryFeeDistributorUpdated)
	if err := _PayerRegistry.contract.UnpackLog(event, "FeeDistributorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayerRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the PayerRegistry contract.
type PayerRegistryInitializedIterator struct {
	Event *PayerRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *PayerRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayerRegistryInitialized)
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
		it.Event = new(PayerRegistryInitialized)
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
func (it *PayerRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayerRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayerRegistryInitialized represents a Initialized event raised by the PayerRegistry contract.
type PayerRegistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_PayerRegistry *PayerRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*PayerRegistryInitializedIterator, error) {

	logs, sub, err := _PayerRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &PayerRegistryInitializedIterator{contract: _PayerRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_PayerRegistry *PayerRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *PayerRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _PayerRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayerRegistryInitialized)
				if err := _PayerRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_PayerRegistry *PayerRegistryFilterer) ParseInitialized(log types.Log) (*PayerRegistryInitialized, error) {
	event := new(PayerRegistryInitialized)
	if err := _PayerRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayerRegistryMigratedIterator is returned from FilterMigrated and is used to iterate over the raw logs and unpacked data for Migrated events raised by the PayerRegistry contract.
type PayerRegistryMigratedIterator struct {
	Event *PayerRegistryMigrated // Event containing the contract specifics and raw log

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
func (it *PayerRegistryMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayerRegistryMigrated)
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
		it.Event = new(PayerRegistryMigrated)
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
func (it *PayerRegistryMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayerRegistryMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayerRegistryMigrated represents a Migrated event raised by the PayerRegistry contract.
type PayerRegistryMigrated struct {
	Migrator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMigrated is a free log retrieval operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_PayerRegistry *PayerRegistryFilterer) FilterMigrated(opts *bind.FilterOpts, migrator []common.Address) (*PayerRegistryMigratedIterator, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _PayerRegistry.contract.FilterLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return &PayerRegistryMigratedIterator{contract: _PayerRegistry.contract, event: "Migrated", logs: logs, sub: sub}, nil
}

// WatchMigrated is a free log subscription operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_PayerRegistry *PayerRegistryFilterer) WatchMigrated(opts *bind.WatchOpts, sink chan<- *PayerRegistryMigrated, migrator []common.Address) (event.Subscription, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _PayerRegistry.contract.WatchLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayerRegistryMigrated)
				if err := _PayerRegistry.contract.UnpackLog(event, "Migrated", log); err != nil {
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
func (_PayerRegistry *PayerRegistryFilterer) ParseMigrated(log types.Log) (*PayerRegistryMigrated, error) {
	event := new(PayerRegistryMigrated)
	if err := _PayerRegistry.contract.UnpackLog(event, "Migrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayerRegistryMinimumDepositUpdatedIterator is returned from FilterMinimumDepositUpdated and is used to iterate over the raw logs and unpacked data for MinimumDepositUpdated events raised by the PayerRegistry contract.
type PayerRegistryMinimumDepositUpdatedIterator struct {
	Event *PayerRegistryMinimumDepositUpdated // Event containing the contract specifics and raw log

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
func (it *PayerRegistryMinimumDepositUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayerRegistryMinimumDepositUpdated)
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
		it.Event = new(PayerRegistryMinimumDepositUpdated)
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
func (it *PayerRegistryMinimumDepositUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayerRegistryMinimumDepositUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayerRegistryMinimumDepositUpdated represents a MinimumDepositUpdated event raised by the PayerRegistry contract.
type PayerRegistryMinimumDepositUpdated struct {
	MinimumDeposit *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterMinimumDepositUpdated is a free log retrieval operation binding the contract event 0x79e23d5ce9133842ec1b1a05e78704d86f2a9f499ddc0eab991061cee393c527.
//
// Solidity: event MinimumDepositUpdated(uint96 minimumDeposit)
func (_PayerRegistry *PayerRegistryFilterer) FilterMinimumDepositUpdated(opts *bind.FilterOpts) (*PayerRegistryMinimumDepositUpdatedIterator, error) {

	logs, sub, err := _PayerRegistry.contract.FilterLogs(opts, "MinimumDepositUpdated")
	if err != nil {
		return nil, err
	}
	return &PayerRegistryMinimumDepositUpdatedIterator{contract: _PayerRegistry.contract, event: "MinimumDepositUpdated", logs: logs, sub: sub}, nil
}

// WatchMinimumDepositUpdated is a free log subscription operation binding the contract event 0x79e23d5ce9133842ec1b1a05e78704d86f2a9f499ddc0eab991061cee393c527.
//
// Solidity: event MinimumDepositUpdated(uint96 minimumDeposit)
func (_PayerRegistry *PayerRegistryFilterer) WatchMinimumDepositUpdated(opts *bind.WatchOpts, sink chan<- *PayerRegistryMinimumDepositUpdated) (event.Subscription, error) {

	logs, sub, err := _PayerRegistry.contract.WatchLogs(opts, "MinimumDepositUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayerRegistryMinimumDepositUpdated)
				if err := _PayerRegistry.contract.UnpackLog(event, "MinimumDepositUpdated", log); err != nil {
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

// ParseMinimumDepositUpdated is a log parse operation binding the contract event 0x79e23d5ce9133842ec1b1a05e78704d86f2a9f499ddc0eab991061cee393c527.
//
// Solidity: event MinimumDepositUpdated(uint96 minimumDeposit)
func (_PayerRegistry *PayerRegistryFilterer) ParseMinimumDepositUpdated(log types.Log) (*PayerRegistryMinimumDepositUpdated, error) {
	event := new(PayerRegistryMinimumDepositUpdated)
	if err := _PayerRegistry.contract.UnpackLog(event, "MinimumDepositUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayerRegistryPauseStatusUpdatedIterator is returned from FilterPauseStatusUpdated and is used to iterate over the raw logs and unpacked data for PauseStatusUpdated events raised by the PayerRegistry contract.
type PayerRegistryPauseStatusUpdatedIterator struct {
	Event *PayerRegistryPauseStatusUpdated // Event containing the contract specifics and raw log

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
func (it *PayerRegistryPauseStatusUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayerRegistryPauseStatusUpdated)
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
		it.Event = new(PayerRegistryPauseStatusUpdated)
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
func (it *PayerRegistryPauseStatusUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayerRegistryPauseStatusUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayerRegistryPauseStatusUpdated represents a PauseStatusUpdated event raised by the PayerRegistry contract.
type PayerRegistryPauseStatusUpdated struct {
	Paused bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPauseStatusUpdated is a free log retrieval operation binding the contract event 0x7c4d1fe30fdbfda9e9c4c43e759ef32e4db5128d4cb58ff3ae9583b89b6242a5.
//
// Solidity: event PauseStatusUpdated(bool indexed paused)
func (_PayerRegistry *PayerRegistryFilterer) FilterPauseStatusUpdated(opts *bind.FilterOpts, paused []bool) (*PayerRegistryPauseStatusUpdatedIterator, error) {

	var pausedRule []interface{}
	for _, pausedItem := range paused {
		pausedRule = append(pausedRule, pausedItem)
	}

	logs, sub, err := _PayerRegistry.contract.FilterLogs(opts, "PauseStatusUpdated", pausedRule)
	if err != nil {
		return nil, err
	}
	return &PayerRegistryPauseStatusUpdatedIterator{contract: _PayerRegistry.contract, event: "PauseStatusUpdated", logs: logs, sub: sub}, nil
}

// WatchPauseStatusUpdated is a free log subscription operation binding the contract event 0x7c4d1fe30fdbfda9e9c4c43e759ef32e4db5128d4cb58ff3ae9583b89b6242a5.
//
// Solidity: event PauseStatusUpdated(bool indexed paused)
func (_PayerRegistry *PayerRegistryFilterer) WatchPauseStatusUpdated(opts *bind.WatchOpts, sink chan<- *PayerRegistryPauseStatusUpdated, paused []bool) (event.Subscription, error) {

	var pausedRule []interface{}
	for _, pausedItem := range paused {
		pausedRule = append(pausedRule, pausedItem)
	}

	logs, sub, err := _PayerRegistry.contract.WatchLogs(opts, "PauseStatusUpdated", pausedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayerRegistryPauseStatusUpdated)
				if err := _PayerRegistry.contract.UnpackLog(event, "PauseStatusUpdated", log); err != nil {
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
func (_PayerRegistry *PayerRegistryFilterer) ParsePauseStatusUpdated(log types.Log) (*PayerRegistryPauseStatusUpdated, error) {
	event := new(PayerRegistryPauseStatusUpdated)
	if err := _PayerRegistry.contract.UnpackLog(event, "PauseStatusUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayerRegistrySettlerUpdatedIterator is returned from FilterSettlerUpdated and is used to iterate over the raw logs and unpacked data for SettlerUpdated events raised by the PayerRegistry contract.
type PayerRegistrySettlerUpdatedIterator struct {
	Event *PayerRegistrySettlerUpdated // Event containing the contract specifics and raw log

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
func (it *PayerRegistrySettlerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayerRegistrySettlerUpdated)
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
		it.Event = new(PayerRegistrySettlerUpdated)
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
func (it *PayerRegistrySettlerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayerRegistrySettlerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayerRegistrySettlerUpdated represents a SettlerUpdated event raised by the PayerRegistry contract.
type PayerRegistrySettlerUpdated struct {
	Settler common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterSettlerUpdated is a free log retrieval operation binding the contract event 0x1f53e003aaf46af23aeb50e85f438d6f0c33618ce44e3545f1ec030c79b17729.
//
// Solidity: event SettlerUpdated(address indexed settler)
func (_PayerRegistry *PayerRegistryFilterer) FilterSettlerUpdated(opts *bind.FilterOpts, settler []common.Address) (*PayerRegistrySettlerUpdatedIterator, error) {

	var settlerRule []interface{}
	for _, settlerItem := range settler {
		settlerRule = append(settlerRule, settlerItem)
	}

	logs, sub, err := _PayerRegistry.contract.FilterLogs(opts, "SettlerUpdated", settlerRule)
	if err != nil {
		return nil, err
	}
	return &PayerRegistrySettlerUpdatedIterator{contract: _PayerRegistry.contract, event: "SettlerUpdated", logs: logs, sub: sub}, nil
}

// WatchSettlerUpdated is a free log subscription operation binding the contract event 0x1f53e003aaf46af23aeb50e85f438d6f0c33618ce44e3545f1ec030c79b17729.
//
// Solidity: event SettlerUpdated(address indexed settler)
func (_PayerRegistry *PayerRegistryFilterer) WatchSettlerUpdated(opts *bind.WatchOpts, sink chan<- *PayerRegistrySettlerUpdated, settler []common.Address) (event.Subscription, error) {

	var settlerRule []interface{}
	for _, settlerItem := range settler {
		settlerRule = append(settlerRule, settlerItem)
	}

	logs, sub, err := _PayerRegistry.contract.WatchLogs(opts, "SettlerUpdated", settlerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayerRegistrySettlerUpdated)
				if err := _PayerRegistry.contract.UnpackLog(event, "SettlerUpdated", log); err != nil {
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

// ParseSettlerUpdated is a log parse operation binding the contract event 0x1f53e003aaf46af23aeb50e85f438d6f0c33618ce44e3545f1ec030c79b17729.
//
// Solidity: event SettlerUpdated(address indexed settler)
func (_PayerRegistry *PayerRegistryFilterer) ParseSettlerUpdated(log types.Log) (*PayerRegistrySettlerUpdated, error) {
	event := new(PayerRegistrySettlerUpdated)
	if err := _PayerRegistry.contract.UnpackLog(event, "SettlerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayerRegistryUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the PayerRegistry contract.
type PayerRegistryUpgradedIterator struct {
	Event *PayerRegistryUpgraded // Event containing the contract specifics and raw log

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
func (it *PayerRegistryUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayerRegistryUpgraded)
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
		it.Event = new(PayerRegistryUpgraded)
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
func (it *PayerRegistryUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayerRegistryUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayerRegistryUpgraded represents a Upgraded event raised by the PayerRegistry contract.
type PayerRegistryUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_PayerRegistry *PayerRegistryFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*PayerRegistryUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _PayerRegistry.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &PayerRegistryUpgradedIterator{contract: _PayerRegistry.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_PayerRegistry *PayerRegistryFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *PayerRegistryUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _PayerRegistry.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayerRegistryUpgraded)
				if err := _PayerRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_PayerRegistry *PayerRegistryFilterer) ParseUpgraded(log types.Log) (*PayerRegistryUpgraded, error) {
	event := new(PayerRegistryUpgraded)
	if err := _PayerRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayerRegistryUsageSettledIterator is returned from FilterUsageSettled and is used to iterate over the raw logs and unpacked data for UsageSettled events raised by the PayerRegistry contract.
type PayerRegistryUsageSettledIterator struct {
	Event *PayerRegistryUsageSettled // Event containing the contract specifics and raw log

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
func (it *PayerRegistryUsageSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayerRegistryUsageSettled)
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
		it.Event = new(PayerRegistryUsageSettled)
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
func (it *PayerRegistryUsageSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayerRegistryUsageSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayerRegistryUsageSettled represents a UsageSettled event raised by the PayerRegistry contract.
type PayerRegistryUsageSettled struct {
	Payer  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterUsageSettled is a free log retrieval operation binding the contract event 0xd642d7d9645fa79e7c4ace29981b4671bad3d6bf2b40433e73f886deda9b7512.
//
// Solidity: event UsageSettled(address indexed payer, uint96 amount)
func (_PayerRegistry *PayerRegistryFilterer) FilterUsageSettled(opts *bind.FilterOpts, payer []common.Address) (*PayerRegistryUsageSettledIterator, error) {

	var payerRule []interface{}
	for _, payerItem := range payer {
		payerRule = append(payerRule, payerItem)
	}

	logs, sub, err := _PayerRegistry.contract.FilterLogs(opts, "UsageSettled", payerRule)
	if err != nil {
		return nil, err
	}
	return &PayerRegistryUsageSettledIterator{contract: _PayerRegistry.contract, event: "UsageSettled", logs: logs, sub: sub}, nil
}

// WatchUsageSettled is a free log subscription operation binding the contract event 0xd642d7d9645fa79e7c4ace29981b4671bad3d6bf2b40433e73f886deda9b7512.
//
// Solidity: event UsageSettled(address indexed payer, uint96 amount)
func (_PayerRegistry *PayerRegistryFilterer) WatchUsageSettled(opts *bind.WatchOpts, sink chan<- *PayerRegistryUsageSettled, payer []common.Address) (event.Subscription, error) {

	var payerRule []interface{}
	for _, payerItem := range payer {
		payerRule = append(payerRule, payerItem)
	}

	logs, sub, err := _PayerRegistry.contract.WatchLogs(opts, "UsageSettled", payerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayerRegistryUsageSettled)
				if err := _PayerRegistry.contract.UnpackLog(event, "UsageSettled", log); err != nil {
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

// ParseUsageSettled is a log parse operation binding the contract event 0xd642d7d9645fa79e7c4ace29981b4671bad3d6bf2b40433e73f886deda9b7512.
//
// Solidity: event UsageSettled(address indexed payer, uint96 amount)
func (_PayerRegistry *PayerRegistryFilterer) ParseUsageSettled(log types.Log) (*PayerRegistryUsageSettled, error) {
	event := new(PayerRegistryUsageSettled)
	if err := _PayerRegistry.contract.UnpackLog(event, "UsageSettled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayerRegistryWithdrawLockPeriodUpdatedIterator is returned from FilterWithdrawLockPeriodUpdated and is used to iterate over the raw logs and unpacked data for WithdrawLockPeriodUpdated events raised by the PayerRegistry contract.
type PayerRegistryWithdrawLockPeriodUpdatedIterator struct {
	Event *PayerRegistryWithdrawLockPeriodUpdated // Event containing the contract specifics and raw log

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
func (it *PayerRegistryWithdrawLockPeriodUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayerRegistryWithdrawLockPeriodUpdated)
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
		it.Event = new(PayerRegistryWithdrawLockPeriodUpdated)
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
func (it *PayerRegistryWithdrawLockPeriodUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayerRegistryWithdrawLockPeriodUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayerRegistryWithdrawLockPeriodUpdated represents a WithdrawLockPeriodUpdated event raised by the PayerRegistry contract.
type PayerRegistryWithdrawLockPeriodUpdated struct {
	WithdrawLockPeriod uint32
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterWithdrawLockPeriodUpdated is a free log retrieval operation binding the contract event 0x775897ecd039cde44004e6e6e5950c773b83326fd56bde68c5985503279cef46.
//
// Solidity: event WithdrawLockPeriodUpdated(uint32 withdrawLockPeriod)
func (_PayerRegistry *PayerRegistryFilterer) FilterWithdrawLockPeriodUpdated(opts *bind.FilterOpts) (*PayerRegistryWithdrawLockPeriodUpdatedIterator, error) {

	logs, sub, err := _PayerRegistry.contract.FilterLogs(opts, "WithdrawLockPeriodUpdated")
	if err != nil {
		return nil, err
	}
	return &PayerRegistryWithdrawLockPeriodUpdatedIterator{contract: _PayerRegistry.contract, event: "WithdrawLockPeriodUpdated", logs: logs, sub: sub}, nil
}

// WatchWithdrawLockPeriodUpdated is a free log subscription operation binding the contract event 0x775897ecd039cde44004e6e6e5950c773b83326fd56bde68c5985503279cef46.
//
// Solidity: event WithdrawLockPeriodUpdated(uint32 withdrawLockPeriod)
func (_PayerRegistry *PayerRegistryFilterer) WatchWithdrawLockPeriodUpdated(opts *bind.WatchOpts, sink chan<- *PayerRegistryWithdrawLockPeriodUpdated) (event.Subscription, error) {

	logs, sub, err := _PayerRegistry.contract.WatchLogs(opts, "WithdrawLockPeriodUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayerRegistryWithdrawLockPeriodUpdated)
				if err := _PayerRegistry.contract.UnpackLog(event, "WithdrawLockPeriodUpdated", log); err != nil {
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

// ParseWithdrawLockPeriodUpdated is a log parse operation binding the contract event 0x775897ecd039cde44004e6e6e5950c773b83326fd56bde68c5985503279cef46.
//
// Solidity: event WithdrawLockPeriodUpdated(uint32 withdrawLockPeriod)
func (_PayerRegistry *PayerRegistryFilterer) ParseWithdrawLockPeriodUpdated(log types.Log) (*PayerRegistryWithdrawLockPeriodUpdated, error) {
	event := new(PayerRegistryWithdrawLockPeriodUpdated)
	if err := _PayerRegistry.contract.UnpackLog(event, "WithdrawLockPeriodUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayerRegistryWithdrawalCancelledIterator is returned from FilterWithdrawalCancelled and is used to iterate over the raw logs and unpacked data for WithdrawalCancelled events raised by the PayerRegistry contract.
type PayerRegistryWithdrawalCancelledIterator struct {
	Event *PayerRegistryWithdrawalCancelled // Event containing the contract specifics and raw log

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
func (it *PayerRegistryWithdrawalCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayerRegistryWithdrawalCancelled)
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
		it.Event = new(PayerRegistryWithdrawalCancelled)
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
func (it *PayerRegistryWithdrawalCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayerRegistryWithdrawalCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayerRegistryWithdrawalCancelled represents a WithdrawalCancelled event raised by the PayerRegistry contract.
type PayerRegistryWithdrawalCancelled struct {
	Payer common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalCancelled is a free log retrieval operation binding the contract event 0xc51fdb96728de385ec7859819e3997bc618362ef0dbca0ad051d856866cda3db.
//
// Solidity: event WithdrawalCancelled(address indexed payer)
func (_PayerRegistry *PayerRegistryFilterer) FilterWithdrawalCancelled(opts *bind.FilterOpts, payer []common.Address) (*PayerRegistryWithdrawalCancelledIterator, error) {

	var payerRule []interface{}
	for _, payerItem := range payer {
		payerRule = append(payerRule, payerItem)
	}

	logs, sub, err := _PayerRegistry.contract.FilterLogs(opts, "WithdrawalCancelled", payerRule)
	if err != nil {
		return nil, err
	}
	return &PayerRegistryWithdrawalCancelledIterator{contract: _PayerRegistry.contract, event: "WithdrawalCancelled", logs: logs, sub: sub}, nil
}

// WatchWithdrawalCancelled is a free log subscription operation binding the contract event 0xc51fdb96728de385ec7859819e3997bc618362ef0dbca0ad051d856866cda3db.
//
// Solidity: event WithdrawalCancelled(address indexed payer)
func (_PayerRegistry *PayerRegistryFilterer) WatchWithdrawalCancelled(opts *bind.WatchOpts, sink chan<- *PayerRegistryWithdrawalCancelled, payer []common.Address) (event.Subscription, error) {

	var payerRule []interface{}
	for _, payerItem := range payer {
		payerRule = append(payerRule, payerItem)
	}

	logs, sub, err := _PayerRegistry.contract.WatchLogs(opts, "WithdrawalCancelled", payerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayerRegistryWithdrawalCancelled)
				if err := _PayerRegistry.contract.UnpackLog(event, "WithdrawalCancelled", log); err != nil {
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

// ParseWithdrawalCancelled is a log parse operation binding the contract event 0xc51fdb96728de385ec7859819e3997bc618362ef0dbca0ad051d856866cda3db.
//
// Solidity: event WithdrawalCancelled(address indexed payer)
func (_PayerRegistry *PayerRegistryFilterer) ParseWithdrawalCancelled(log types.Log) (*PayerRegistryWithdrawalCancelled, error) {
	event := new(PayerRegistryWithdrawalCancelled)
	if err := _PayerRegistry.contract.UnpackLog(event, "WithdrawalCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayerRegistryWithdrawalFinalizedIterator is returned from FilterWithdrawalFinalized and is used to iterate over the raw logs and unpacked data for WithdrawalFinalized events raised by the PayerRegistry contract.
type PayerRegistryWithdrawalFinalizedIterator struct {
	Event *PayerRegistryWithdrawalFinalized // Event containing the contract specifics and raw log

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
func (it *PayerRegistryWithdrawalFinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayerRegistryWithdrawalFinalized)
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
		it.Event = new(PayerRegistryWithdrawalFinalized)
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
func (it *PayerRegistryWithdrawalFinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayerRegistryWithdrawalFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayerRegistryWithdrawalFinalized represents a WithdrawalFinalized event raised by the PayerRegistry contract.
type PayerRegistryWithdrawalFinalized struct {
	Payer common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalFinalized is a free log retrieval operation binding the contract event 0xa637395731e680948e2ce15889fef15308910430bf7a63a535288ed4604f0333.
//
// Solidity: event WithdrawalFinalized(address indexed payer)
func (_PayerRegistry *PayerRegistryFilterer) FilterWithdrawalFinalized(opts *bind.FilterOpts, payer []common.Address) (*PayerRegistryWithdrawalFinalizedIterator, error) {

	var payerRule []interface{}
	for _, payerItem := range payer {
		payerRule = append(payerRule, payerItem)
	}

	logs, sub, err := _PayerRegistry.contract.FilterLogs(opts, "WithdrawalFinalized", payerRule)
	if err != nil {
		return nil, err
	}
	return &PayerRegistryWithdrawalFinalizedIterator{contract: _PayerRegistry.contract, event: "WithdrawalFinalized", logs: logs, sub: sub}, nil
}

// WatchWithdrawalFinalized is a free log subscription operation binding the contract event 0xa637395731e680948e2ce15889fef15308910430bf7a63a535288ed4604f0333.
//
// Solidity: event WithdrawalFinalized(address indexed payer)
func (_PayerRegistry *PayerRegistryFilterer) WatchWithdrawalFinalized(opts *bind.WatchOpts, sink chan<- *PayerRegistryWithdrawalFinalized, payer []common.Address) (event.Subscription, error) {

	var payerRule []interface{}
	for _, payerItem := range payer {
		payerRule = append(payerRule, payerItem)
	}

	logs, sub, err := _PayerRegistry.contract.WatchLogs(opts, "WithdrawalFinalized", payerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayerRegistryWithdrawalFinalized)
				if err := _PayerRegistry.contract.UnpackLog(event, "WithdrawalFinalized", log); err != nil {
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

// ParseWithdrawalFinalized is a log parse operation binding the contract event 0xa637395731e680948e2ce15889fef15308910430bf7a63a535288ed4604f0333.
//
// Solidity: event WithdrawalFinalized(address indexed payer)
func (_PayerRegistry *PayerRegistryFilterer) ParseWithdrawalFinalized(log types.Log) (*PayerRegistryWithdrawalFinalized, error) {
	event := new(PayerRegistryWithdrawalFinalized)
	if err := _PayerRegistry.contract.UnpackLog(event, "WithdrawalFinalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayerRegistryWithdrawalRequestedIterator is returned from FilterWithdrawalRequested and is used to iterate over the raw logs and unpacked data for WithdrawalRequested events raised by the PayerRegistry contract.
type PayerRegistryWithdrawalRequestedIterator struct {
	Event *PayerRegistryWithdrawalRequested // Event containing the contract specifics and raw log

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
func (it *PayerRegistryWithdrawalRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayerRegistryWithdrawalRequested)
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
		it.Event = new(PayerRegistryWithdrawalRequested)
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
func (it *PayerRegistryWithdrawalRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayerRegistryWithdrawalRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayerRegistryWithdrawalRequested represents a WithdrawalRequested event raised by the PayerRegistry contract.
type PayerRegistryWithdrawalRequested struct {
	Payer                 common.Address
	Amount                *big.Int
	WithdrawableTimestamp uint32
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalRequested is a free log retrieval operation binding the contract event 0x44767ee038d5f04d9489720720f411e526a07e4790957f953141be69ce080502.
//
// Solidity: event WithdrawalRequested(address indexed payer, uint96 amount, uint32 withdrawableTimestamp)
func (_PayerRegistry *PayerRegistryFilterer) FilterWithdrawalRequested(opts *bind.FilterOpts, payer []common.Address) (*PayerRegistryWithdrawalRequestedIterator, error) {

	var payerRule []interface{}
	for _, payerItem := range payer {
		payerRule = append(payerRule, payerItem)
	}

	logs, sub, err := _PayerRegistry.contract.FilterLogs(opts, "WithdrawalRequested", payerRule)
	if err != nil {
		return nil, err
	}
	return &PayerRegistryWithdrawalRequestedIterator{contract: _PayerRegistry.contract, event: "WithdrawalRequested", logs: logs, sub: sub}, nil
}

// WatchWithdrawalRequested is a free log subscription operation binding the contract event 0x44767ee038d5f04d9489720720f411e526a07e4790957f953141be69ce080502.
//
// Solidity: event WithdrawalRequested(address indexed payer, uint96 amount, uint32 withdrawableTimestamp)
func (_PayerRegistry *PayerRegistryFilterer) WatchWithdrawalRequested(opts *bind.WatchOpts, sink chan<- *PayerRegistryWithdrawalRequested, payer []common.Address) (event.Subscription, error) {

	var payerRule []interface{}
	for _, payerItem := range payer {
		payerRule = append(payerRule, payerItem)
	}

	logs, sub, err := _PayerRegistry.contract.WatchLogs(opts, "WithdrawalRequested", payerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayerRegistryWithdrawalRequested)
				if err := _PayerRegistry.contract.UnpackLog(event, "WithdrawalRequested", log); err != nil {
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

// ParseWithdrawalRequested is a log parse operation binding the contract event 0x44767ee038d5f04d9489720720f411e526a07e4790957f953141be69ce080502.
//
// Solidity: event WithdrawalRequested(address indexed payer, uint96 amount, uint32 withdrawableTimestamp)
func (_PayerRegistry *PayerRegistryFilterer) ParseWithdrawalRequested(log types.Log) (*PayerRegistryWithdrawalRequested, error) {
	event := new(PayerRegistryWithdrawalRequested)
	if err := _PayerRegistry.contract.UnpackLog(event, "WithdrawalRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
