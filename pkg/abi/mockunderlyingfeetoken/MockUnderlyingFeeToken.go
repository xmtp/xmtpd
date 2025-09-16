// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mockunderlyingfeetoken

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

// MockUnderlyingFeeTokenMetaData contains all meta data concerning the MockUnderlyingFeeToken contract.
var MockUnderlyingFeeTokenMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"parameterRegistry_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DOMAIN_SEPARATOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"decimals_\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"implementation_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migrate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migratorParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"to_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nonces\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"parameterRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"permit\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Migrated\",\"inputs\":[{\"name\":\"migrator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC2612ExpiredSignature\",\"inputs\":[{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC2612InvalidSigner\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"EmptyCode\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidAccountNonce\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"currentNonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationFailed\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertData_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterOutOfTypeBounds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroMigrator\",\"inputs\":[]}]",
}

// MockUnderlyingFeeTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use MockUnderlyingFeeTokenMetaData.ABI instead.
var MockUnderlyingFeeTokenABI = MockUnderlyingFeeTokenMetaData.ABI

// MockUnderlyingFeeToken is an auto generated Go binding around an Ethereum contract.
type MockUnderlyingFeeToken struct {
	MockUnderlyingFeeTokenCaller     // Read-only binding to the contract
	MockUnderlyingFeeTokenTransactor // Write-only binding to the contract
	MockUnderlyingFeeTokenFilterer   // Log filterer for contract events
}

// MockUnderlyingFeeTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type MockUnderlyingFeeTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockUnderlyingFeeTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MockUnderlyingFeeTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockUnderlyingFeeTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MockUnderlyingFeeTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockUnderlyingFeeTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MockUnderlyingFeeTokenSession struct {
	Contract     *MockUnderlyingFeeToken // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// MockUnderlyingFeeTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MockUnderlyingFeeTokenCallerSession struct {
	Contract *MockUnderlyingFeeTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// MockUnderlyingFeeTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MockUnderlyingFeeTokenTransactorSession struct {
	Contract     *MockUnderlyingFeeTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// MockUnderlyingFeeTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type MockUnderlyingFeeTokenRaw struct {
	Contract *MockUnderlyingFeeToken // Generic contract binding to access the raw methods on
}

// MockUnderlyingFeeTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MockUnderlyingFeeTokenCallerRaw struct {
	Contract *MockUnderlyingFeeTokenCaller // Generic read-only contract binding to access the raw methods on
}

// MockUnderlyingFeeTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MockUnderlyingFeeTokenTransactorRaw struct {
	Contract *MockUnderlyingFeeTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMockUnderlyingFeeToken creates a new instance of MockUnderlyingFeeToken, bound to a specific deployed contract.
func NewMockUnderlyingFeeToken(address common.Address, backend bind.ContractBackend) (*MockUnderlyingFeeToken, error) {
	contract, err := bindMockUnderlyingFeeToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockUnderlyingFeeToken{MockUnderlyingFeeTokenCaller: MockUnderlyingFeeTokenCaller{contract: contract}, MockUnderlyingFeeTokenTransactor: MockUnderlyingFeeTokenTransactor{contract: contract}, MockUnderlyingFeeTokenFilterer: MockUnderlyingFeeTokenFilterer{contract: contract}}, nil
}

// NewMockUnderlyingFeeTokenCaller creates a new read-only instance of MockUnderlyingFeeToken, bound to a specific deployed contract.
func NewMockUnderlyingFeeTokenCaller(address common.Address, caller bind.ContractCaller) (*MockUnderlyingFeeTokenCaller, error) {
	contract, err := bindMockUnderlyingFeeToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockUnderlyingFeeTokenCaller{contract: contract}, nil
}

// NewMockUnderlyingFeeTokenTransactor creates a new write-only instance of MockUnderlyingFeeToken, bound to a specific deployed contract.
func NewMockUnderlyingFeeTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*MockUnderlyingFeeTokenTransactor, error) {
	contract, err := bindMockUnderlyingFeeToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockUnderlyingFeeTokenTransactor{contract: contract}, nil
}

// NewMockUnderlyingFeeTokenFilterer creates a new log filterer instance of MockUnderlyingFeeToken, bound to a specific deployed contract.
func NewMockUnderlyingFeeTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*MockUnderlyingFeeTokenFilterer, error) {
	contract, err := bindMockUnderlyingFeeToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockUnderlyingFeeTokenFilterer{contract: contract}, nil
}

// bindMockUnderlyingFeeToken binds a generic wrapper to an already deployed contract.
func bindMockUnderlyingFeeToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockUnderlyingFeeTokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockUnderlyingFeeToken.Contract.MockUnderlyingFeeTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.Contract.MockUnderlyingFeeTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.Contract.MockUnderlyingFeeTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockUnderlyingFeeToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.Contract.contract.Transact(opts, method, params...)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCaller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _MockUnderlyingFeeToken.contract.Call(opts, &out, "DOMAIN_SEPARATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _MockUnderlyingFeeToken.Contract.DOMAINSEPARATOR(&_MockUnderlyingFeeToken.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _MockUnderlyingFeeToken.Contract.DOMAINSEPARATOR(&_MockUnderlyingFeeToken.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MockUnderlyingFeeToken.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _MockUnderlyingFeeToken.Contract.Allowance(&_MockUnderlyingFeeToken.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _MockUnderlyingFeeToken.Contract.Allowance(&_MockUnderlyingFeeToken.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MockUnderlyingFeeToken.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _MockUnderlyingFeeToken.Contract.BalanceOf(&_MockUnderlyingFeeToken.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _MockUnderlyingFeeToken.Contract.BalanceOf(&_MockUnderlyingFeeToken.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8 decimals_)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _MockUnderlyingFeeToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8 decimals_)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenSession) Decimals() (uint8, error) {
	return _MockUnderlyingFeeToken.Contract.Decimals(&_MockUnderlyingFeeToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8 decimals_)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCallerSession) Decimals() (uint8, error) {
	return _MockUnderlyingFeeToken.Contract.Decimals(&_MockUnderlyingFeeToken.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _MockUnderlyingFeeToken.contract.Call(opts, &out, "eip712Domain")

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
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _MockUnderlyingFeeToken.Contract.Eip712Domain(&_MockUnderlyingFeeToken.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _MockUnderlyingFeeToken.Contract.Eip712Domain(&_MockUnderlyingFeeToken.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MockUnderlyingFeeToken.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenSession) Implementation() (common.Address, error) {
	return _MockUnderlyingFeeToken.Contract.Implementation(&_MockUnderlyingFeeToken.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCallerSession) Implementation() (common.Address, error) {
	return _MockUnderlyingFeeToken.Contract.Implementation(&_MockUnderlyingFeeToken.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCaller) MigratorParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockUnderlyingFeeToken.contract.Call(opts, &out, "migratorParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenSession) MigratorParameterKey() (string, error) {
	return _MockUnderlyingFeeToken.Contract.MigratorParameterKey(&_MockUnderlyingFeeToken.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCallerSession) MigratorParameterKey() (string, error) {
	return _MockUnderlyingFeeToken.Contract.MigratorParameterKey(&_MockUnderlyingFeeToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockUnderlyingFeeToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenSession) Name() (string, error) {
	return _MockUnderlyingFeeToken.Contract.Name(&_MockUnderlyingFeeToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCallerSession) Name() (string, error) {
	return _MockUnderlyingFeeToken.Contract.Name(&_MockUnderlyingFeeToken.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCaller) Nonces(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MockUnderlyingFeeToken.contract.Call(opts, &out, "nonces", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenSession) Nonces(owner common.Address) (*big.Int, error) {
	return _MockUnderlyingFeeToken.Contract.Nonces(&_MockUnderlyingFeeToken.CallOpts, owner)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCallerSession) Nonces(owner common.Address) (*big.Int, error) {
	return _MockUnderlyingFeeToken.Contract.Nonces(&_MockUnderlyingFeeToken.CallOpts, owner)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCaller) ParameterRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MockUnderlyingFeeToken.contract.Call(opts, &out, "parameterRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenSession) ParameterRegistry() (common.Address, error) {
	return _MockUnderlyingFeeToken.Contract.ParameterRegistry(&_MockUnderlyingFeeToken.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCallerSession) ParameterRegistry() (common.Address, error) {
	return _MockUnderlyingFeeToken.Contract.ParameterRegistry(&_MockUnderlyingFeeToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockUnderlyingFeeToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenSession) Symbol() (string, error) {
	return _MockUnderlyingFeeToken.Contract.Symbol(&_MockUnderlyingFeeToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCallerSession) Symbol() (string, error) {
	return _MockUnderlyingFeeToken.Contract.Symbol(&_MockUnderlyingFeeToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockUnderlyingFeeToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenSession) TotalSupply() (*big.Int, error) {
	return _MockUnderlyingFeeToken.Contract.TotalSupply(&_MockUnderlyingFeeToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _MockUnderlyingFeeToken.Contract.TotalSupply(&_MockUnderlyingFeeToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.Contract.Approve(&_MockUnderlyingFeeToken.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.Contract.Approve(&_MockUnderlyingFeeToken.TransactOpts, spender, value)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenSession) Initialize() (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.Contract.Initialize(&_MockUnderlyingFeeToken.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenTransactorSession) Initialize() (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.Contract.Initialize(&_MockUnderlyingFeeToken.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenTransactor) Migrate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.contract.Transact(opts, "migrate")
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenSession) Migrate() (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.Contract.Migrate(&_MockUnderlyingFeeToken.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenTransactorSession) Migrate() (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.Contract.Migrate(&_MockUnderlyingFeeToken.TransactOpts)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to_, uint256 amount_) returns()
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenTransactor) Mint(opts *bind.TransactOpts, to_ common.Address, amount_ *big.Int) (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.contract.Transact(opts, "mint", to_, amount_)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to_, uint256 amount_) returns()
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenSession) Mint(to_ common.Address, amount_ *big.Int) (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.Contract.Mint(&_MockUnderlyingFeeToken.TransactOpts, to_, amount_)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to_, uint256 amount_) returns()
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenTransactorSession) Mint(to_ common.Address, amount_ *big.Int) (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.Contract.Mint(&_MockUnderlyingFeeToken.TransactOpts, to_, amount_)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenTransactor) Permit(opts *bind.TransactOpts, owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.contract.Transact(opts, "permit", owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.Contract.Permit(&_MockUnderlyingFeeToken.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenTransactorSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.Contract.Permit(&_MockUnderlyingFeeToken.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.Contract.Transfer(&_MockUnderlyingFeeToken.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.Contract.Transfer(&_MockUnderlyingFeeToken.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.Contract.TransferFrom(&_MockUnderlyingFeeToken.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockUnderlyingFeeToken.Contract.TransferFrom(&_MockUnderlyingFeeToken.TransactOpts, from, to, value)
}

// MockUnderlyingFeeTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the MockUnderlyingFeeToken contract.
type MockUnderlyingFeeTokenApprovalIterator struct {
	Event *MockUnderlyingFeeTokenApproval // Event containing the contract specifics and raw log

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
func (it *MockUnderlyingFeeTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockUnderlyingFeeTokenApproval)
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
		it.Event = new(MockUnderlyingFeeTokenApproval)
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
func (it *MockUnderlyingFeeTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockUnderlyingFeeTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockUnderlyingFeeTokenApproval represents a Approval event raised by the MockUnderlyingFeeToken contract.
type MockUnderlyingFeeTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*MockUnderlyingFeeTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _MockUnderlyingFeeToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &MockUnderlyingFeeTokenApprovalIterator{contract: _MockUnderlyingFeeToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *MockUnderlyingFeeTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _MockUnderlyingFeeToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockUnderlyingFeeTokenApproval)
				if err := _MockUnderlyingFeeToken.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenFilterer) ParseApproval(log types.Log) (*MockUnderlyingFeeTokenApproval, error) {
	event := new(MockUnderlyingFeeTokenApproval)
	if err := _MockUnderlyingFeeToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockUnderlyingFeeTokenEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the MockUnderlyingFeeToken contract.
type MockUnderlyingFeeTokenEIP712DomainChangedIterator struct {
	Event *MockUnderlyingFeeTokenEIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *MockUnderlyingFeeTokenEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockUnderlyingFeeTokenEIP712DomainChanged)
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
		it.Event = new(MockUnderlyingFeeTokenEIP712DomainChanged)
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
func (it *MockUnderlyingFeeTokenEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockUnderlyingFeeTokenEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockUnderlyingFeeTokenEIP712DomainChanged represents a EIP712DomainChanged event raised by the MockUnderlyingFeeToken contract.
type MockUnderlyingFeeTokenEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*MockUnderlyingFeeTokenEIP712DomainChangedIterator, error) {

	logs, sub, err := _MockUnderlyingFeeToken.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &MockUnderlyingFeeTokenEIP712DomainChangedIterator{contract: _MockUnderlyingFeeToken.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *MockUnderlyingFeeTokenEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _MockUnderlyingFeeToken.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockUnderlyingFeeTokenEIP712DomainChanged)
				if err := _MockUnderlyingFeeToken.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenFilterer) ParseEIP712DomainChanged(log types.Log) (*MockUnderlyingFeeTokenEIP712DomainChanged, error) {
	event := new(MockUnderlyingFeeTokenEIP712DomainChanged)
	if err := _MockUnderlyingFeeToken.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockUnderlyingFeeTokenInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the MockUnderlyingFeeToken contract.
type MockUnderlyingFeeTokenInitializedIterator struct {
	Event *MockUnderlyingFeeTokenInitialized // Event containing the contract specifics and raw log

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
func (it *MockUnderlyingFeeTokenInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockUnderlyingFeeTokenInitialized)
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
		it.Event = new(MockUnderlyingFeeTokenInitialized)
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
func (it *MockUnderlyingFeeTokenInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockUnderlyingFeeTokenInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockUnderlyingFeeTokenInitialized represents a Initialized event raised by the MockUnderlyingFeeToken contract.
type MockUnderlyingFeeTokenInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenFilterer) FilterInitialized(opts *bind.FilterOpts) (*MockUnderlyingFeeTokenInitializedIterator, error) {

	logs, sub, err := _MockUnderlyingFeeToken.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &MockUnderlyingFeeTokenInitializedIterator{contract: _MockUnderlyingFeeToken.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *MockUnderlyingFeeTokenInitialized) (event.Subscription, error) {

	logs, sub, err := _MockUnderlyingFeeToken.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockUnderlyingFeeTokenInitialized)
				if err := _MockUnderlyingFeeToken.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenFilterer) ParseInitialized(log types.Log) (*MockUnderlyingFeeTokenInitialized, error) {
	event := new(MockUnderlyingFeeTokenInitialized)
	if err := _MockUnderlyingFeeToken.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockUnderlyingFeeTokenMigratedIterator is returned from FilterMigrated and is used to iterate over the raw logs and unpacked data for Migrated events raised by the MockUnderlyingFeeToken contract.
type MockUnderlyingFeeTokenMigratedIterator struct {
	Event *MockUnderlyingFeeTokenMigrated // Event containing the contract specifics and raw log

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
func (it *MockUnderlyingFeeTokenMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockUnderlyingFeeTokenMigrated)
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
		it.Event = new(MockUnderlyingFeeTokenMigrated)
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
func (it *MockUnderlyingFeeTokenMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockUnderlyingFeeTokenMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockUnderlyingFeeTokenMigrated represents a Migrated event raised by the MockUnderlyingFeeToken contract.
type MockUnderlyingFeeTokenMigrated struct {
	Migrator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMigrated is a free log retrieval operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenFilterer) FilterMigrated(opts *bind.FilterOpts, migrator []common.Address) (*MockUnderlyingFeeTokenMigratedIterator, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _MockUnderlyingFeeToken.contract.FilterLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return &MockUnderlyingFeeTokenMigratedIterator{contract: _MockUnderlyingFeeToken.contract, event: "Migrated", logs: logs, sub: sub}, nil
}

// WatchMigrated is a free log subscription operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenFilterer) WatchMigrated(opts *bind.WatchOpts, sink chan<- *MockUnderlyingFeeTokenMigrated, migrator []common.Address) (event.Subscription, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _MockUnderlyingFeeToken.contract.WatchLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockUnderlyingFeeTokenMigrated)
				if err := _MockUnderlyingFeeToken.contract.UnpackLog(event, "Migrated", log); err != nil {
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
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenFilterer) ParseMigrated(log types.Log) (*MockUnderlyingFeeTokenMigrated, error) {
	event := new(MockUnderlyingFeeTokenMigrated)
	if err := _MockUnderlyingFeeToken.contract.UnpackLog(event, "Migrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockUnderlyingFeeTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the MockUnderlyingFeeToken contract.
type MockUnderlyingFeeTokenTransferIterator struct {
	Event *MockUnderlyingFeeTokenTransfer // Event containing the contract specifics and raw log

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
func (it *MockUnderlyingFeeTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockUnderlyingFeeTokenTransfer)
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
		it.Event = new(MockUnderlyingFeeTokenTransfer)
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
func (it *MockUnderlyingFeeTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockUnderlyingFeeTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockUnderlyingFeeTokenTransfer represents a Transfer event raised by the MockUnderlyingFeeToken contract.
type MockUnderlyingFeeTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockUnderlyingFeeTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockUnderlyingFeeToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MockUnderlyingFeeTokenTransferIterator{contract: _MockUnderlyingFeeToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *MockUnderlyingFeeTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockUnderlyingFeeToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockUnderlyingFeeTokenTransfer)
				if err := _MockUnderlyingFeeToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenFilterer) ParseTransfer(log types.Log) (*MockUnderlyingFeeTokenTransfer, error) {
	event := new(MockUnderlyingFeeTokenTransfer)
	if err := _MockUnderlyingFeeToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockUnderlyingFeeTokenUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the MockUnderlyingFeeToken contract.
type MockUnderlyingFeeTokenUpgradedIterator struct {
	Event *MockUnderlyingFeeTokenUpgraded // Event containing the contract specifics and raw log

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
func (it *MockUnderlyingFeeTokenUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockUnderlyingFeeTokenUpgraded)
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
		it.Event = new(MockUnderlyingFeeTokenUpgraded)
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
func (it *MockUnderlyingFeeTokenUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockUnderlyingFeeTokenUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockUnderlyingFeeTokenUpgraded represents a Upgraded event raised by the MockUnderlyingFeeToken contract.
type MockUnderlyingFeeTokenUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*MockUnderlyingFeeTokenUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _MockUnderlyingFeeToken.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &MockUnderlyingFeeTokenUpgradedIterator{contract: _MockUnderlyingFeeToken.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *MockUnderlyingFeeTokenUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _MockUnderlyingFeeToken.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockUnderlyingFeeTokenUpgraded)
				if err := _MockUnderlyingFeeToken.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_MockUnderlyingFeeToken *MockUnderlyingFeeTokenFilterer) ParseUpgraded(log types.Log) (*MockUnderlyingFeeTokenUpgraded, error) {
	event := new(MockUnderlyingFeeTokenUpgraded)
	if err := _MockUnderlyingFeeToken.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
