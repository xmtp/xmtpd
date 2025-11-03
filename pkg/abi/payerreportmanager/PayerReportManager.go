// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package payerreportmanager

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

// IPayerReportManagerPayerReport is an auto generated low-level Go binding around an user-defined struct.
type IPayerReportManagerPayerReport struct {
	StartSequenceId     uint64
	EndSequenceId       uint64
	EndMinuteSinceEpoch uint32
	FeesSettled         *big.Int
	Offset              uint32
	IsSettled           bool
	ProtocolFeeRate     uint16
	PayersMerkleRoot    [32]byte
	NodeIds             []uint32
}

// IPayerReportManagerPayerReportSignature is an auto generated low-level Go binding around an user-defined struct.
type IPayerReportManagerPayerReportSignature struct {
	NodeId    uint32
	Signature []byte
}

// PayerReportManagerMetaData contains all meta data concerning the PayerReportManager contract.
var PayerReportManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"parameterRegistry_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nodeRegistry_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"payerRegistry_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DOMAIN_SEPARATOR\",\"inputs\":[],\"outputs\":[{\"name\":\"domainSeparator_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ONE_HUNDRED_PERCENT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PAYER_REPORT_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields_\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions_\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPayerReport\",\"inputs\":[{\"name\":\"originatorNodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"payerReportIndex_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"payerReport_\",\"type\":\"tuple\",\"internalType\":\"structIPayerReportManager.PayerReport\",\"components\":[{\"name\":\"startSequenceId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endSequenceId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endMinuteSinceEpoch\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"feesSettled\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"offset\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isSettled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"protocolFeeRate\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"payersMerkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nodeIds\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPayerReportDigest\",\"inputs\":[{\"name\":\"originatorNodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"startSequenceId_\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endSequenceId_\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endMinuteSinceEpoch_\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"payersMerkleRoot_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nodeIds_\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}],\"outputs\":[{\"name\":\"digest_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPayerReports\",\"inputs\":[{\"name\":\"originatorNodeIds_\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"payerReportIndices_\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"payerReports_\",\"type\":\"tuple[]\",\"internalType\":\"structIPayerReportManager.PayerReport[]\",\"components\":[{\"name\":\"startSequenceId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endSequenceId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endMinuteSinceEpoch\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"feesSettled\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"offset\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isSettled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"protocolFeeRate\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"payersMerkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nodeIds\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"implementation_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migrate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migratorParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"nodeRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"parameterRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"payerRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"protocolFeeRate\",\"inputs\":[],\"outputs\":[{\"name\":\"protocolFeeRate_\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"protocolFeeRateParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"settle\",\"inputs\":[{\"name\":\"originatorNodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"payerReportIndex_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"payerFees_\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"proofElements_\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submit\",\"inputs\":[{\"name\":\"originatorNodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"startSequenceId_\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endSequenceId_\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endMinuteSinceEpoch_\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"payersMerkleRoot_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nodeIds_\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"signatures_\",\"type\":\"tuple[]\",\"internalType\":\"structIPayerReportManager.PayerReportSignature[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"payerReportIndex_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateProtocolFeeRate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Migrated\",\"inputs\":[{\"name\":\"migrator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayerReportSubmitted\",\"inputs\":[{\"name\":\"originatorNodeId\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"payerReportIndex\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startSequenceId\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"endSequenceId\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"endMinuteSinceEpoch\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"payersMerkleRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"nodeIds\",\"type\":\"uint32[]\",\"indexed\":false,\"internalType\":\"uint32[]\"},{\"name\":\"signingNodeIds\",\"type\":\"uint32[]\",\"indexed\":false,\"internalType\":\"uint32[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayerReportSubsetSettled\",\"inputs\":[{\"name\":\"originatorNodeId\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"payerReportIndex\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"count\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"remaining\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"feesSettled\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProtocolFeeRateUpdated\",\"inputs\":[{\"name\":\"protocolFeeRate\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ArrayLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyCode\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InsufficientSignatures\",\"inputs\":[{\"name\":\"validSignatureCount\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"requiredSignatureCount\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"InvalidBitCount32Input\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidLeafCount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidProtocolFeeRate\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSequenceIds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidStartSequenceId\",\"inputs\":[{\"name\":\"startSequenceId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"lastSequenceId\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MigrationFailed\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertData_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"NoChange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoLeaves\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoProofElements\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoReportsForOriginator\",\"inputs\":[{\"name\":\"originatorNodeId\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"NodeIdAtIndexMismatch\",\"inputs\":[{\"name\":\"expectedId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"actualId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"NodeIdsLengthMismatch\",\"inputs\":[{\"name\":\"expectedCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"providedCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterOutOfTypeBounds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PayerFeesLengthTooLong\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PayerReportAlreadySubmitted\",\"inputs\":[{\"name\":\"originatorNodeId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"startSequenceId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endSequenceId\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"PayerReportEntirelySettled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PayerReportIndexOutOfBounds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SettleUsageFailed\",\"inputs\":[{\"name\":\"returnData_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"UnorderedNodeIds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroMigrator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroNodeRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroParameterRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroPayerRegistry\",\"inputs\":[]}]",
}

// PayerReportManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use PayerReportManagerMetaData.ABI instead.
var PayerReportManagerABI = PayerReportManagerMetaData.ABI

// PayerReportManager is an auto generated Go binding around an Ethereum contract.
type PayerReportManager struct {
	PayerReportManagerCaller     // Read-only binding to the contract
	PayerReportManagerTransactor // Write-only binding to the contract
	PayerReportManagerFilterer   // Log filterer for contract events
}

// PayerReportManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type PayerReportManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PayerReportManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PayerReportManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PayerReportManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PayerReportManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PayerReportManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PayerReportManagerSession struct {
	Contract     *PayerReportManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// PayerReportManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PayerReportManagerCallerSession struct {
	Contract *PayerReportManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// PayerReportManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PayerReportManagerTransactorSession struct {
	Contract     *PayerReportManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// PayerReportManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type PayerReportManagerRaw struct {
	Contract *PayerReportManager // Generic contract binding to access the raw methods on
}

// PayerReportManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PayerReportManagerCallerRaw struct {
	Contract *PayerReportManagerCaller // Generic read-only contract binding to access the raw methods on
}

// PayerReportManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PayerReportManagerTransactorRaw struct {
	Contract *PayerReportManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPayerReportManager creates a new instance of PayerReportManager, bound to a specific deployed contract.
func NewPayerReportManager(address common.Address, backend bind.ContractBackend) (*PayerReportManager, error) {
	contract, err := bindPayerReportManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PayerReportManager{PayerReportManagerCaller: PayerReportManagerCaller{contract: contract}, PayerReportManagerTransactor: PayerReportManagerTransactor{contract: contract}, PayerReportManagerFilterer: PayerReportManagerFilterer{contract: contract}}, nil
}

// NewPayerReportManagerCaller creates a new read-only instance of PayerReportManager, bound to a specific deployed contract.
func NewPayerReportManagerCaller(address common.Address, caller bind.ContractCaller) (*PayerReportManagerCaller, error) {
	contract, err := bindPayerReportManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PayerReportManagerCaller{contract: contract}, nil
}

// NewPayerReportManagerTransactor creates a new write-only instance of PayerReportManager, bound to a specific deployed contract.
func NewPayerReportManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*PayerReportManagerTransactor, error) {
	contract, err := bindPayerReportManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PayerReportManagerTransactor{contract: contract}, nil
}

// NewPayerReportManagerFilterer creates a new log filterer instance of PayerReportManager, bound to a specific deployed contract.
func NewPayerReportManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*PayerReportManagerFilterer, error) {
	contract, err := bindPayerReportManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PayerReportManagerFilterer{contract: contract}, nil
}

// bindPayerReportManager binds a generic wrapper to an already deployed contract.
func bindPayerReportManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PayerReportManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PayerReportManager *PayerReportManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PayerReportManager.Contract.PayerReportManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PayerReportManager *PayerReportManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayerReportManager.Contract.PayerReportManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PayerReportManager *PayerReportManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PayerReportManager.Contract.PayerReportManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PayerReportManager *PayerReportManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PayerReportManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PayerReportManager *PayerReportManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayerReportManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PayerReportManager *PayerReportManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PayerReportManager.Contract.contract.Transact(opts, method, params...)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32 domainSeparator_)
func (_PayerReportManager *PayerReportManagerCaller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _PayerReportManager.contract.Call(opts, &out, "DOMAIN_SEPARATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32 domainSeparator_)
func (_PayerReportManager *PayerReportManagerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _PayerReportManager.Contract.DOMAINSEPARATOR(&_PayerReportManager.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32 domainSeparator_)
func (_PayerReportManager *PayerReportManagerCallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _PayerReportManager.Contract.DOMAINSEPARATOR(&_PayerReportManager.CallOpts)
}

// ONEHUNDREDPERCENT is a free data retrieval call binding the contract method 0xdd0081c7.
//
// Solidity: function ONE_HUNDRED_PERCENT() view returns(uint16)
func (_PayerReportManager *PayerReportManagerCaller) ONEHUNDREDPERCENT(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _PayerReportManager.contract.Call(opts, &out, "ONE_HUNDRED_PERCENT")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// ONEHUNDREDPERCENT is a free data retrieval call binding the contract method 0xdd0081c7.
//
// Solidity: function ONE_HUNDRED_PERCENT() view returns(uint16)
func (_PayerReportManager *PayerReportManagerSession) ONEHUNDREDPERCENT() (uint16, error) {
	return _PayerReportManager.Contract.ONEHUNDREDPERCENT(&_PayerReportManager.CallOpts)
}

// ONEHUNDREDPERCENT is a free data retrieval call binding the contract method 0xdd0081c7.
//
// Solidity: function ONE_HUNDRED_PERCENT() view returns(uint16)
func (_PayerReportManager *PayerReportManagerCallerSession) ONEHUNDREDPERCENT() (uint16, error) {
	return _PayerReportManager.Contract.ONEHUNDREDPERCENT(&_PayerReportManager.CallOpts)
}

// PAYERREPORTTYPEHASH is a free data retrieval call binding the contract method 0x3d8fcde2.
//
// Solidity: function PAYER_REPORT_TYPEHASH() view returns(bytes32)
func (_PayerReportManager *PayerReportManagerCaller) PAYERREPORTTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _PayerReportManager.contract.Call(opts, &out, "PAYER_REPORT_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PAYERREPORTTYPEHASH is a free data retrieval call binding the contract method 0x3d8fcde2.
//
// Solidity: function PAYER_REPORT_TYPEHASH() view returns(bytes32)
func (_PayerReportManager *PayerReportManagerSession) PAYERREPORTTYPEHASH() ([32]byte, error) {
	return _PayerReportManager.Contract.PAYERREPORTTYPEHASH(&_PayerReportManager.CallOpts)
}

// PAYERREPORTTYPEHASH is a free data retrieval call binding the contract method 0x3d8fcde2.
//
// Solidity: function PAYER_REPORT_TYPEHASH() view returns(bytes32)
func (_PayerReportManager *PayerReportManagerCallerSession) PAYERREPORTTYPEHASH() ([32]byte, error) {
	return _PayerReportManager.Contract.PAYERREPORTTYPEHASH(&_PayerReportManager.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields_, string name_, string version_, uint256 chainId_, address verifyingContract_, bytes32 salt_, uint256[] extensions_)
func (_PayerReportManager *PayerReportManagerCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _PayerReportManager.contract.Call(opts, &out, "eip712Domain")

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
// Solidity: function eip712Domain() view returns(bytes1 fields_, string name_, string version_, uint256 chainId_, address verifyingContract_, bytes32 salt_, uint256[] extensions_)
func (_PayerReportManager *PayerReportManagerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _PayerReportManager.Contract.Eip712Domain(&_PayerReportManager.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields_, string name_, string version_, uint256 chainId_, address verifyingContract_, bytes32 salt_, uint256[] extensions_)
func (_PayerReportManager *PayerReportManagerCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _PayerReportManager.Contract.Eip712Domain(&_PayerReportManager.CallOpts)
}

// GetPayerReport is a free data retrieval call binding the contract method 0x22ccd722.
//
// Solidity: function getPayerReport(uint32 originatorNodeId_, uint256 payerReportIndex_) view returns((uint64,uint64,uint32,uint96,uint32,bool,uint16,bytes32,uint32[]) payerReport_)
func (_PayerReportManager *PayerReportManagerCaller) GetPayerReport(opts *bind.CallOpts, originatorNodeId_ uint32, payerReportIndex_ *big.Int) (IPayerReportManagerPayerReport, error) {
	var out []interface{}
	err := _PayerReportManager.contract.Call(opts, &out, "getPayerReport", originatorNodeId_, payerReportIndex_)

	if err != nil {
		return *new(IPayerReportManagerPayerReport), err
	}

	out0 := *abi.ConvertType(out[0], new(IPayerReportManagerPayerReport)).(*IPayerReportManagerPayerReport)

	return out0, err

}

// GetPayerReport is a free data retrieval call binding the contract method 0x22ccd722.
//
// Solidity: function getPayerReport(uint32 originatorNodeId_, uint256 payerReportIndex_) view returns((uint64,uint64,uint32,uint96,uint32,bool,uint16,bytes32,uint32[]) payerReport_)
func (_PayerReportManager *PayerReportManagerSession) GetPayerReport(originatorNodeId_ uint32, payerReportIndex_ *big.Int) (IPayerReportManagerPayerReport, error) {
	return _PayerReportManager.Contract.GetPayerReport(&_PayerReportManager.CallOpts, originatorNodeId_, payerReportIndex_)
}

// GetPayerReport is a free data retrieval call binding the contract method 0x22ccd722.
//
// Solidity: function getPayerReport(uint32 originatorNodeId_, uint256 payerReportIndex_) view returns((uint64,uint64,uint32,uint96,uint32,bool,uint16,bytes32,uint32[]) payerReport_)
func (_PayerReportManager *PayerReportManagerCallerSession) GetPayerReport(originatorNodeId_ uint32, payerReportIndex_ *big.Int) (IPayerReportManagerPayerReport, error) {
	return _PayerReportManager.Contract.GetPayerReport(&_PayerReportManager.CallOpts, originatorNodeId_, payerReportIndex_)
}

// GetPayerReportDigest is a free data retrieval call binding the contract method 0x356e1189.
//
// Solidity: function getPayerReportDigest(uint32 originatorNodeId_, uint64 startSequenceId_, uint64 endSequenceId_, uint32 endMinuteSinceEpoch_, bytes32 payersMerkleRoot_, uint32[] nodeIds_) view returns(bytes32 digest_)
func (_PayerReportManager *PayerReportManagerCaller) GetPayerReportDigest(opts *bind.CallOpts, originatorNodeId_ uint32, startSequenceId_ uint64, endSequenceId_ uint64, endMinuteSinceEpoch_ uint32, payersMerkleRoot_ [32]byte, nodeIds_ []uint32) ([32]byte, error) {
	var out []interface{}
	err := _PayerReportManager.contract.Call(opts, &out, "getPayerReportDigest", originatorNodeId_, startSequenceId_, endSequenceId_, endMinuteSinceEpoch_, payersMerkleRoot_, nodeIds_)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetPayerReportDigest is a free data retrieval call binding the contract method 0x356e1189.
//
// Solidity: function getPayerReportDigest(uint32 originatorNodeId_, uint64 startSequenceId_, uint64 endSequenceId_, uint32 endMinuteSinceEpoch_, bytes32 payersMerkleRoot_, uint32[] nodeIds_) view returns(bytes32 digest_)
func (_PayerReportManager *PayerReportManagerSession) GetPayerReportDigest(originatorNodeId_ uint32, startSequenceId_ uint64, endSequenceId_ uint64, endMinuteSinceEpoch_ uint32, payersMerkleRoot_ [32]byte, nodeIds_ []uint32) ([32]byte, error) {
	return _PayerReportManager.Contract.GetPayerReportDigest(&_PayerReportManager.CallOpts, originatorNodeId_, startSequenceId_, endSequenceId_, endMinuteSinceEpoch_, payersMerkleRoot_, nodeIds_)
}

// GetPayerReportDigest is a free data retrieval call binding the contract method 0x356e1189.
//
// Solidity: function getPayerReportDigest(uint32 originatorNodeId_, uint64 startSequenceId_, uint64 endSequenceId_, uint32 endMinuteSinceEpoch_, bytes32 payersMerkleRoot_, uint32[] nodeIds_) view returns(bytes32 digest_)
func (_PayerReportManager *PayerReportManagerCallerSession) GetPayerReportDigest(originatorNodeId_ uint32, startSequenceId_ uint64, endSequenceId_ uint64, endMinuteSinceEpoch_ uint32, payersMerkleRoot_ [32]byte, nodeIds_ []uint32) ([32]byte, error) {
	return _PayerReportManager.Contract.GetPayerReportDigest(&_PayerReportManager.CallOpts, originatorNodeId_, startSequenceId_, endSequenceId_, endMinuteSinceEpoch_, payersMerkleRoot_, nodeIds_)
}

// GetPayerReports is a free data retrieval call binding the contract method 0xb881bca0.
//
// Solidity: function getPayerReports(uint32[] originatorNodeIds_, uint256[] payerReportIndices_) view returns((uint64,uint64,uint32,uint96,uint32,bool,uint16,bytes32,uint32[])[] payerReports_)
func (_PayerReportManager *PayerReportManagerCaller) GetPayerReports(opts *bind.CallOpts, originatorNodeIds_ []uint32, payerReportIndices_ []*big.Int) ([]IPayerReportManagerPayerReport, error) {
	var out []interface{}
	err := _PayerReportManager.contract.Call(opts, &out, "getPayerReports", originatorNodeIds_, payerReportIndices_)

	if err != nil {
		return *new([]IPayerReportManagerPayerReport), err
	}

	out0 := *abi.ConvertType(out[0], new([]IPayerReportManagerPayerReport)).(*[]IPayerReportManagerPayerReport)

	return out0, err

}

// GetPayerReports is a free data retrieval call binding the contract method 0xb881bca0.
//
// Solidity: function getPayerReports(uint32[] originatorNodeIds_, uint256[] payerReportIndices_) view returns((uint64,uint64,uint32,uint96,uint32,bool,uint16,bytes32,uint32[])[] payerReports_)
func (_PayerReportManager *PayerReportManagerSession) GetPayerReports(originatorNodeIds_ []uint32, payerReportIndices_ []*big.Int) ([]IPayerReportManagerPayerReport, error) {
	return _PayerReportManager.Contract.GetPayerReports(&_PayerReportManager.CallOpts, originatorNodeIds_, payerReportIndices_)
}

// GetPayerReports is a free data retrieval call binding the contract method 0xb881bca0.
//
// Solidity: function getPayerReports(uint32[] originatorNodeIds_, uint256[] payerReportIndices_) view returns((uint64,uint64,uint32,uint96,uint32,bool,uint16,bytes32,uint32[])[] payerReports_)
func (_PayerReportManager *PayerReportManagerCallerSession) GetPayerReports(originatorNodeIds_ []uint32, payerReportIndices_ []*big.Int) ([]IPayerReportManagerPayerReport, error) {
	return _PayerReportManager.Contract.GetPayerReports(&_PayerReportManager.CallOpts, originatorNodeIds_, payerReportIndices_)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_PayerReportManager *PayerReportManagerCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PayerReportManager.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_PayerReportManager *PayerReportManagerSession) Implementation() (common.Address, error) {
	return _PayerReportManager.Contract.Implementation(&_PayerReportManager.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_PayerReportManager *PayerReportManagerCallerSession) Implementation() (common.Address, error) {
	return _PayerReportManager.Contract.Implementation(&_PayerReportManager.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_PayerReportManager *PayerReportManagerCaller) MigratorParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _PayerReportManager.contract.Call(opts, &out, "migratorParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_PayerReportManager *PayerReportManagerSession) MigratorParameterKey() (string, error) {
	return _PayerReportManager.Contract.MigratorParameterKey(&_PayerReportManager.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(string key_)
func (_PayerReportManager *PayerReportManagerCallerSession) MigratorParameterKey() (string, error) {
	return _PayerReportManager.Contract.MigratorParameterKey(&_PayerReportManager.CallOpts)
}

// NodeRegistry is a free data retrieval call binding the contract method 0xd9b5c4a5.
//
// Solidity: function nodeRegistry() view returns(address)
func (_PayerReportManager *PayerReportManagerCaller) NodeRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PayerReportManager.contract.Call(opts, &out, "nodeRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NodeRegistry is a free data retrieval call binding the contract method 0xd9b5c4a5.
//
// Solidity: function nodeRegistry() view returns(address)
func (_PayerReportManager *PayerReportManagerSession) NodeRegistry() (common.Address, error) {
	return _PayerReportManager.Contract.NodeRegistry(&_PayerReportManager.CallOpts)
}

// NodeRegistry is a free data retrieval call binding the contract method 0xd9b5c4a5.
//
// Solidity: function nodeRegistry() view returns(address)
func (_PayerReportManager *PayerReportManagerCallerSession) NodeRegistry() (common.Address, error) {
	return _PayerReportManager.Contract.NodeRegistry(&_PayerReportManager.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_PayerReportManager *PayerReportManagerCaller) ParameterRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PayerReportManager.contract.Call(opts, &out, "parameterRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_PayerReportManager *PayerReportManagerSession) ParameterRegistry() (common.Address, error) {
	return _PayerReportManager.Contract.ParameterRegistry(&_PayerReportManager.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_PayerReportManager *PayerReportManagerCallerSession) ParameterRegistry() (common.Address, error) {
	return _PayerReportManager.Contract.ParameterRegistry(&_PayerReportManager.CallOpts)
}

// PayerRegistry is a free data retrieval call binding the contract method 0x1dc5f4b8.
//
// Solidity: function payerRegistry() view returns(address)
func (_PayerReportManager *PayerReportManagerCaller) PayerRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PayerReportManager.contract.Call(opts, &out, "payerRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PayerRegistry is a free data retrieval call binding the contract method 0x1dc5f4b8.
//
// Solidity: function payerRegistry() view returns(address)
func (_PayerReportManager *PayerReportManagerSession) PayerRegistry() (common.Address, error) {
	return _PayerReportManager.Contract.PayerRegistry(&_PayerReportManager.CallOpts)
}

// PayerRegistry is a free data retrieval call binding the contract method 0x1dc5f4b8.
//
// Solidity: function payerRegistry() view returns(address)
func (_PayerReportManager *PayerReportManagerCallerSession) PayerRegistry() (common.Address, error) {
	return _PayerReportManager.Contract.PayerRegistry(&_PayerReportManager.CallOpts)
}

// ProtocolFeeRate is a free data retrieval call binding the contract method 0x58f85880.
//
// Solidity: function protocolFeeRate() view returns(uint16 protocolFeeRate_)
func (_PayerReportManager *PayerReportManagerCaller) ProtocolFeeRate(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _PayerReportManager.contract.Call(opts, &out, "protocolFeeRate")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// ProtocolFeeRate is a free data retrieval call binding the contract method 0x58f85880.
//
// Solidity: function protocolFeeRate() view returns(uint16 protocolFeeRate_)
func (_PayerReportManager *PayerReportManagerSession) ProtocolFeeRate() (uint16, error) {
	return _PayerReportManager.Contract.ProtocolFeeRate(&_PayerReportManager.CallOpts)
}

// ProtocolFeeRate is a free data retrieval call binding the contract method 0x58f85880.
//
// Solidity: function protocolFeeRate() view returns(uint16 protocolFeeRate_)
func (_PayerReportManager *PayerReportManagerCallerSession) ProtocolFeeRate() (uint16, error) {
	return _PayerReportManager.Contract.ProtocolFeeRate(&_PayerReportManager.CallOpts)
}

// ProtocolFeeRateParameterKey is a free data retrieval call binding the contract method 0xacc48f0f.
//
// Solidity: function protocolFeeRateParameterKey() pure returns(string key_)
func (_PayerReportManager *PayerReportManagerCaller) ProtocolFeeRateParameterKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _PayerReportManager.contract.Call(opts, &out, "protocolFeeRateParameterKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ProtocolFeeRateParameterKey is a free data retrieval call binding the contract method 0xacc48f0f.
//
// Solidity: function protocolFeeRateParameterKey() pure returns(string key_)
func (_PayerReportManager *PayerReportManagerSession) ProtocolFeeRateParameterKey() (string, error) {
	return _PayerReportManager.Contract.ProtocolFeeRateParameterKey(&_PayerReportManager.CallOpts)
}

// ProtocolFeeRateParameterKey is a free data retrieval call binding the contract method 0xacc48f0f.
//
// Solidity: function protocolFeeRateParameterKey() pure returns(string key_)
func (_PayerReportManager *PayerReportManagerCallerSession) ProtocolFeeRateParameterKey() (string, error) {
	return _PayerReportManager.Contract.ProtocolFeeRateParameterKey(&_PayerReportManager.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_PayerReportManager *PayerReportManagerTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayerReportManager.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_PayerReportManager *PayerReportManagerSession) Initialize() (*types.Transaction, error) {
	return _PayerReportManager.Contract.Initialize(&_PayerReportManager.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_PayerReportManager *PayerReportManagerTransactorSession) Initialize() (*types.Transaction, error) {
	return _PayerReportManager.Contract.Initialize(&_PayerReportManager.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_PayerReportManager *PayerReportManagerTransactor) Migrate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayerReportManager.contract.Transact(opts, "migrate")
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_PayerReportManager *PayerReportManagerSession) Migrate() (*types.Transaction, error) {
	return _PayerReportManager.Contract.Migrate(&_PayerReportManager.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_PayerReportManager *PayerReportManagerTransactorSession) Migrate() (*types.Transaction, error) {
	return _PayerReportManager.Contract.Migrate(&_PayerReportManager.TransactOpts)
}

// Settle is a paid mutator transaction binding the contract method 0x6576143c.
//
// Solidity: function settle(uint32 originatorNodeId_, uint256 payerReportIndex_, bytes[] payerFees_, bytes32[] proofElements_) returns()
func (_PayerReportManager *PayerReportManagerTransactor) Settle(opts *bind.TransactOpts, originatorNodeId_ uint32, payerReportIndex_ *big.Int, payerFees_ [][]byte, proofElements_ [][32]byte) (*types.Transaction, error) {
	return _PayerReportManager.contract.Transact(opts, "settle", originatorNodeId_, payerReportIndex_, payerFees_, proofElements_)
}

// Settle is a paid mutator transaction binding the contract method 0x6576143c.
//
// Solidity: function settle(uint32 originatorNodeId_, uint256 payerReportIndex_, bytes[] payerFees_, bytes32[] proofElements_) returns()
func (_PayerReportManager *PayerReportManagerSession) Settle(originatorNodeId_ uint32, payerReportIndex_ *big.Int, payerFees_ [][]byte, proofElements_ [][32]byte) (*types.Transaction, error) {
	return _PayerReportManager.Contract.Settle(&_PayerReportManager.TransactOpts, originatorNodeId_, payerReportIndex_, payerFees_, proofElements_)
}

// Settle is a paid mutator transaction binding the contract method 0x6576143c.
//
// Solidity: function settle(uint32 originatorNodeId_, uint256 payerReportIndex_, bytes[] payerFees_, bytes32[] proofElements_) returns()
func (_PayerReportManager *PayerReportManagerTransactorSession) Settle(originatorNodeId_ uint32, payerReportIndex_ *big.Int, payerFees_ [][]byte, proofElements_ [][32]byte) (*types.Transaction, error) {
	return _PayerReportManager.Contract.Settle(&_PayerReportManager.TransactOpts, originatorNodeId_, payerReportIndex_, payerFees_, proofElements_)
}

// Submit is a paid mutator transaction binding the contract method 0x844446cd.
//
// Solidity: function submit(uint32 originatorNodeId_, uint64 startSequenceId_, uint64 endSequenceId_, uint32 endMinuteSinceEpoch_, bytes32 payersMerkleRoot_, uint32[] nodeIds_, (uint32,bytes)[] signatures_) returns(uint256 payerReportIndex_)
func (_PayerReportManager *PayerReportManagerTransactor) Submit(opts *bind.TransactOpts, originatorNodeId_ uint32, startSequenceId_ uint64, endSequenceId_ uint64, endMinuteSinceEpoch_ uint32, payersMerkleRoot_ [32]byte, nodeIds_ []uint32, signatures_ []IPayerReportManagerPayerReportSignature) (*types.Transaction, error) {
	return _PayerReportManager.contract.Transact(opts, "submit", originatorNodeId_, startSequenceId_, endSequenceId_, endMinuteSinceEpoch_, payersMerkleRoot_, nodeIds_, signatures_)
}

// Submit is a paid mutator transaction binding the contract method 0x844446cd.
//
// Solidity: function submit(uint32 originatorNodeId_, uint64 startSequenceId_, uint64 endSequenceId_, uint32 endMinuteSinceEpoch_, bytes32 payersMerkleRoot_, uint32[] nodeIds_, (uint32,bytes)[] signatures_) returns(uint256 payerReportIndex_)
func (_PayerReportManager *PayerReportManagerSession) Submit(originatorNodeId_ uint32, startSequenceId_ uint64, endSequenceId_ uint64, endMinuteSinceEpoch_ uint32, payersMerkleRoot_ [32]byte, nodeIds_ []uint32, signatures_ []IPayerReportManagerPayerReportSignature) (*types.Transaction, error) {
	return _PayerReportManager.Contract.Submit(&_PayerReportManager.TransactOpts, originatorNodeId_, startSequenceId_, endSequenceId_, endMinuteSinceEpoch_, payersMerkleRoot_, nodeIds_, signatures_)
}

// Submit is a paid mutator transaction binding the contract method 0x844446cd.
//
// Solidity: function submit(uint32 originatorNodeId_, uint64 startSequenceId_, uint64 endSequenceId_, uint32 endMinuteSinceEpoch_, bytes32 payersMerkleRoot_, uint32[] nodeIds_, (uint32,bytes)[] signatures_) returns(uint256 payerReportIndex_)
func (_PayerReportManager *PayerReportManagerTransactorSession) Submit(originatorNodeId_ uint32, startSequenceId_ uint64, endSequenceId_ uint64, endMinuteSinceEpoch_ uint32, payersMerkleRoot_ [32]byte, nodeIds_ []uint32, signatures_ []IPayerReportManagerPayerReportSignature) (*types.Transaction, error) {
	return _PayerReportManager.Contract.Submit(&_PayerReportManager.TransactOpts, originatorNodeId_, startSequenceId_, endSequenceId_, endMinuteSinceEpoch_, payersMerkleRoot_, nodeIds_, signatures_)
}

// UpdateProtocolFeeRate is a paid mutator transaction binding the contract method 0xd3444eee.
//
// Solidity: function updateProtocolFeeRate() returns()
func (_PayerReportManager *PayerReportManagerTransactor) UpdateProtocolFeeRate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayerReportManager.contract.Transact(opts, "updateProtocolFeeRate")
}

// UpdateProtocolFeeRate is a paid mutator transaction binding the contract method 0xd3444eee.
//
// Solidity: function updateProtocolFeeRate() returns()
func (_PayerReportManager *PayerReportManagerSession) UpdateProtocolFeeRate() (*types.Transaction, error) {
	return _PayerReportManager.Contract.UpdateProtocolFeeRate(&_PayerReportManager.TransactOpts)
}

// UpdateProtocolFeeRate is a paid mutator transaction binding the contract method 0xd3444eee.
//
// Solidity: function updateProtocolFeeRate() returns()
func (_PayerReportManager *PayerReportManagerTransactorSession) UpdateProtocolFeeRate() (*types.Transaction, error) {
	return _PayerReportManager.Contract.UpdateProtocolFeeRate(&_PayerReportManager.TransactOpts)
}

// PayerReportManagerEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the PayerReportManager contract.
type PayerReportManagerEIP712DomainChangedIterator struct {
	Event *PayerReportManagerEIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *PayerReportManagerEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayerReportManagerEIP712DomainChanged)
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
		it.Event = new(PayerReportManagerEIP712DomainChanged)
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
func (it *PayerReportManagerEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayerReportManagerEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayerReportManagerEIP712DomainChanged represents a EIP712DomainChanged event raised by the PayerReportManager contract.
type PayerReportManagerEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_PayerReportManager *PayerReportManagerFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*PayerReportManagerEIP712DomainChangedIterator, error) {

	logs, sub, err := _PayerReportManager.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &PayerReportManagerEIP712DomainChangedIterator{contract: _PayerReportManager.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_PayerReportManager *PayerReportManagerFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *PayerReportManagerEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _PayerReportManager.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayerReportManagerEIP712DomainChanged)
				if err := _PayerReportManager.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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
func (_PayerReportManager *PayerReportManagerFilterer) ParseEIP712DomainChanged(log types.Log) (*PayerReportManagerEIP712DomainChanged, error) {
	event := new(PayerReportManagerEIP712DomainChanged)
	if err := _PayerReportManager.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayerReportManagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the PayerReportManager contract.
type PayerReportManagerInitializedIterator struct {
	Event *PayerReportManagerInitialized // Event containing the contract specifics and raw log

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
func (it *PayerReportManagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayerReportManagerInitialized)
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
		it.Event = new(PayerReportManagerInitialized)
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
func (it *PayerReportManagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayerReportManagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayerReportManagerInitialized represents a Initialized event raised by the PayerReportManager contract.
type PayerReportManagerInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_PayerReportManager *PayerReportManagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*PayerReportManagerInitializedIterator, error) {

	logs, sub, err := _PayerReportManager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &PayerReportManagerInitializedIterator{contract: _PayerReportManager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_PayerReportManager *PayerReportManagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *PayerReportManagerInitialized) (event.Subscription, error) {

	logs, sub, err := _PayerReportManager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayerReportManagerInitialized)
				if err := _PayerReportManager.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_PayerReportManager *PayerReportManagerFilterer) ParseInitialized(log types.Log) (*PayerReportManagerInitialized, error) {
	event := new(PayerReportManagerInitialized)
	if err := _PayerReportManager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayerReportManagerMigratedIterator is returned from FilterMigrated and is used to iterate over the raw logs and unpacked data for Migrated events raised by the PayerReportManager contract.
type PayerReportManagerMigratedIterator struct {
	Event *PayerReportManagerMigrated // Event containing the contract specifics and raw log

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
func (it *PayerReportManagerMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayerReportManagerMigrated)
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
		it.Event = new(PayerReportManagerMigrated)
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
func (it *PayerReportManagerMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayerReportManagerMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayerReportManagerMigrated represents a Migrated event raised by the PayerReportManager contract.
type PayerReportManagerMigrated struct {
	Migrator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMigrated is a free log retrieval operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_PayerReportManager *PayerReportManagerFilterer) FilterMigrated(opts *bind.FilterOpts, migrator []common.Address) (*PayerReportManagerMigratedIterator, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _PayerReportManager.contract.FilterLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return &PayerReportManagerMigratedIterator{contract: _PayerReportManager.contract, event: "Migrated", logs: logs, sub: sub}, nil
}

// WatchMigrated is a free log subscription operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_PayerReportManager *PayerReportManagerFilterer) WatchMigrated(opts *bind.WatchOpts, sink chan<- *PayerReportManagerMigrated, migrator []common.Address) (event.Subscription, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _PayerReportManager.contract.WatchLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayerReportManagerMigrated)
				if err := _PayerReportManager.contract.UnpackLog(event, "Migrated", log); err != nil {
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
func (_PayerReportManager *PayerReportManagerFilterer) ParseMigrated(log types.Log) (*PayerReportManagerMigrated, error) {
	event := new(PayerReportManagerMigrated)
	if err := _PayerReportManager.contract.UnpackLog(event, "Migrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayerReportManagerPayerReportSubmittedIterator is returned from FilterPayerReportSubmitted and is used to iterate over the raw logs and unpacked data for PayerReportSubmitted events raised by the PayerReportManager contract.
type PayerReportManagerPayerReportSubmittedIterator struct {
	Event *PayerReportManagerPayerReportSubmitted // Event containing the contract specifics and raw log

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
func (it *PayerReportManagerPayerReportSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayerReportManagerPayerReportSubmitted)
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
		it.Event = new(PayerReportManagerPayerReportSubmitted)
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
func (it *PayerReportManagerPayerReportSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayerReportManagerPayerReportSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayerReportManagerPayerReportSubmitted represents a PayerReportSubmitted event raised by the PayerReportManager contract.
type PayerReportManagerPayerReportSubmitted struct {
	OriginatorNodeId    uint32
	PayerReportIndex    *big.Int
	StartSequenceId     uint64
	EndSequenceId       uint64
	EndMinuteSinceEpoch uint32
	PayersMerkleRoot    [32]byte
	NodeIds             []uint32
	SigningNodeIds      []uint32
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterPayerReportSubmitted is a free log retrieval operation binding the contract event 0xc51249b9e6e78c1a6b57bbca4dd3ff04bff48ba6f7246ab4d109a89e8a824dda.
//
// Solidity: event PayerReportSubmitted(uint32 indexed originatorNodeId, uint256 indexed payerReportIndex, uint64 startSequenceId, uint64 indexed endSequenceId, uint32 endMinuteSinceEpoch, bytes32 payersMerkleRoot, uint32[] nodeIds, uint32[] signingNodeIds)
func (_PayerReportManager *PayerReportManagerFilterer) FilterPayerReportSubmitted(opts *bind.FilterOpts, originatorNodeId []uint32, payerReportIndex []*big.Int, endSequenceId []uint64) (*PayerReportManagerPayerReportSubmittedIterator, error) {

	var originatorNodeIdRule []interface{}
	for _, originatorNodeIdItem := range originatorNodeId {
		originatorNodeIdRule = append(originatorNodeIdRule, originatorNodeIdItem)
	}
	var payerReportIndexRule []interface{}
	for _, payerReportIndexItem := range payerReportIndex {
		payerReportIndexRule = append(payerReportIndexRule, payerReportIndexItem)
	}

	var endSequenceIdRule []interface{}
	for _, endSequenceIdItem := range endSequenceId {
		endSequenceIdRule = append(endSequenceIdRule, endSequenceIdItem)
	}

	logs, sub, err := _PayerReportManager.contract.FilterLogs(opts, "PayerReportSubmitted", originatorNodeIdRule, payerReportIndexRule, endSequenceIdRule)
	if err != nil {
		return nil, err
	}
	return &PayerReportManagerPayerReportSubmittedIterator{contract: _PayerReportManager.contract, event: "PayerReportSubmitted", logs: logs, sub: sub}, nil
}

// WatchPayerReportSubmitted is a free log subscription operation binding the contract event 0xc51249b9e6e78c1a6b57bbca4dd3ff04bff48ba6f7246ab4d109a89e8a824dda.
//
// Solidity: event PayerReportSubmitted(uint32 indexed originatorNodeId, uint256 indexed payerReportIndex, uint64 startSequenceId, uint64 indexed endSequenceId, uint32 endMinuteSinceEpoch, bytes32 payersMerkleRoot, uint32[] nodeIds, uint32[] signingNodeIds)
func (_PayerReportManager *PayerReportManagerFilterer) WatchPayerReportSubmitted(opts *bind.WatchOpts, sink chan<- *PayerReportManagerPayerReportSubmitted, originatorNodeId []uint32, payerReportIndex []*big.Int, endSequenceId []uint64) (event.Subscription, error) {

	var originatorNodeIdRule []interface{}
	for _, originatorNodeIdItem := range originatorNodeId {
		originatorNodeIdRule = append(originatorNodeIdRule, originatorNodeIdItem)
	}
	var payerReportIndexRule []interface{}
	for _, payerReportIndexItem := range payerReportIndex {
		payerReportIndexRule = append(payerReportIndexRule, payerReportIndexItem)
	}

	var endSequenceIdRule []interface{}
	for _, endSequenceIdItem := range endSequenceId {
		endSequenceIdRule = append(endSequenceIdRule, endSequenceIdItem)
	}

	logs, sub, err := _PayerReportManager.contract.WatchLogs(opts, "PayerReportSubmitted", originatorNodeIdRule, payerReportIndexRule, endSequenceIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayerReportManagerPayerReportSubmitted)
				if err := _PayerReportManager.contract.UnpackLog(event, "PayerReportSubmitted", log); err != nil {
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

// ParsePayerReportSubmitted is a log parse operation binding the contract event 0xc51249b9e6e78c1a6b57bbca4dd3ff04bff48ba6f7246ab4d109a89e8a824dda.
//
// Solidity: event PayerReportSubmitted(uint32 indexed originatorNodeId, uint256 indexed payerReportIndex, uint64 startSequenceId, uint64 indexed endSequenceId, uint32 endMinuteSinceEpoch, bytes32 payersMerkleRoot, uint32[] nodeIds, uint32[] signingNodeIds)
func (_PayerReportManager *PayerReportManagerFilterer) ParsePayerReportSubmitted(log types.Log) (*PayerReportManagerPayerReportSubmitted, error) {
	event := new(PayerReportManagerPayerReportSubmitted)
	if err := _PayerReportManager.contract.UnpackLog(event, "PayerReportSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayerReportManagerPayerReportSubsetSettledIterator is returned from FilterPayerReportSubsetSettled and is used to iterate over the raw logs and unpacked data for PayerReportSubsetSettled events raised by the PayerReportManager contract.
type PayerReportManagerPayerReportSubsetSettledIterator struct {
	Event *PayerReportManagerPayerReportSubsetSettled // Event containing the contract specifics and raw log

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
func (it *PayerReportManagerPayerReportSubsetSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayerReportManagerPayerReportSubsetSettled)
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
		it.Event = new(PayerReportManagerPayerReportSubsetSettled)
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
func (it *PayerReportManagerPayerReportSubsetSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayerReportManagerPayerReportSubsetSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayerReportManagerPayerReportSubsetSettled represents a PayerReportSubsetSettled event raised by the PayerReportManager contract.
type PayerReportManagerPayerReportSubsetSettled struct {
	OriginatorNodeId uint32
	PayerReportIndex *big.Int
	Count            uint32
	Remaining        uint32
	FeesSettled      *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterPayerReportSubsetSettled is a free log retrieval operation binding the contract event 0x3e2ef0c87bba9a992cfdc5189b278bd9852605eaca1ee05c64f285cff2c07691.
//
// Solidity: event PayerReportSubsetSettled(uint32 indexed originatorNodeId, uint256 indexed payerReportIndex, uint32 count, uint32 remaining, uint96 feesSettled)
func (_PayerReportManager *PayerReportManagerFilterer) FilterPayerReportSubsetSettled(opts *bind.FilterOpts, originatorNodeId []uint32, payerReportIndex []*big.Int) (*PayerReportManagerPayerReportSubsetSettledIterator, error) {

	var originatorNodeIdRule []interface{}
	for _, originatorNodeIdItem := range originatorNodeId {
		originatorNodeIdRule = append(originatorNodeIdRule, originatorNodeIdItem)
	}
	var payerReportIndexRule []interface{}
	for _, payerReportIndexItem := range payerReportIndex {
		payerReportIndexRule = append(payerReportIndexRule, payerReportIndexItem)
	}

	logs, sub, err := _PayerReportManager.contract.FilterLogs(opts, "PayerReportSubsetSettled", originatorNodeIdRule, payerReportIndexRule)
	if err != nil {
		return nil, err
	}
	return &PayerReportManagerPayerReportSubsetSettledIterator{contract: _PayerReportManager.contract, event: "PayerReportSubsetSettled", logs: logs, sub: sub}, nil
}

// WatchPayerReportSubsetSettled is a free log subscription operation binding the contract event 0x3e2ef0c87bba9a992cfdc5189b278bd9852605eaca1ee05c64f285cff2c07691.
//
// Solidity: event PayerReportSubsetSettled(uint32 indexed originatorNodeId, uint256 indexed payerReportIndex, uint32 count, uint32 remaining, uint96 feesSettled)
func (_PayerReportManager *PayerReportManagerFilterer) WatchPayerReportSubsetSettled(opts *bind.WatchOpts, sink chan<- *PayerReportManagerPayerReportSubsetSettled, originatorNodeId []uint32, payerReportIndex []*big.Int) (event.Subscription, error) {

	var originatorNodeIdRule []interface{}
	for _, originatorNodeIdItem := range originatorNodeId {
		originatorNodeIdRule = append(originatorNodeIdRule, originatorNodeIdItem)
	}
	var payerReportIndexRule []interface{}
	for _, payerReportIndexItem := range payerReportIndex {
		payerReportIndexRule = append(payerReportIndexRule, payerReportIndexItem)
	}

	logs, sub, err := _PayerReportManager.contract.WatchLogs(opts, "PayerReportSubsetSettled", originatorNodeIdRule, payerReportIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayerReportManagerPayerReportSubsetSettled)
				if err := _PayerReportManager.contract.UnpackLog(event, "PayerReportSubsetSettled", log); err != nil {
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

// ParsePayerReportSubsetSettled is a log parse operation binding the contract event 0x3e2ef0c87bba9a992cfdc5189b278bd9852605eaca1ee05c64f285cff2c07691.
//
// Solidity: event PayerReportSubsetSettled(uint32 indexed originatorNodeId, uint256 indexed payerReportIndex, uint32 count, uint32 remaining, uint96 feesSettled)
func (_PayerReportManager *PayerReportManagerFilterer) ParsePayerReportSubsetSettled(log types.Log) (*PayerReportManagerPayerReportSubsetSettled, error) {
	event := new(PayerReportManagerPayerReportSubsetSettled)
	if err := _PayerReportManager.contract.UnpackLog(event, "PayerReportSubsetSettled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayerReportManagerProtocolFeeRateUpdatedIterator is returned from FilterProtocolFeeRateUpdated and is used to iterate over the raw logs and unpacked data for ProtocolFeeRateUpdated events raised by the PayerReportManager contract.
type PayerReportManagerProtocolFeeRateUpdatedIterator struct {
	Event *PayerReportManagerProtocolFeeRateUpdated // Event containing the contract specifics and raw log

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
func (it *PayerReportManagerProtocolFeeRateUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayerReportManagerProtocolFeeRateUpdated)
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
		it.Event = new(PayerReportManagerProtocolFeeRateUpdated)
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
func (it *PayerReportManagerProtocolFeeRateUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayerReportManagerProtocolFeeRateUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayerReportManagerProtocolFeeRateUpdated represents a ProtocolFeeRateUpdated event raised by the PayerReportManager contract.
type PayerReportManagerProtocolFeeRateUpdated struct {
	ProtocolFeeRate uint16
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterProtocolFeeRateUpdated is a free log retrieval operation binding the contract event 0x1b89f70e2fb3771e6a9f5cd2d2ad3439b823b67e0e69bdeb89b7d68a35fdd5c9.
//
// Solidity: event ProtocolFeeRateUpdated(uint16 protocolFeeRate)
func (_PayerReportManager *PayerReportManagerFilterer) FilterProtocolFeeRateUpdated(opts *bind.FilterOpts) (*PayerReportManagerProtocolFeeRateUpdatedIterator, error) {

	logs, sub, err := _PayerReportManager.contract.FilterLogs(opts, "ProtocolFeeRateUpdated")
	if err != nil {
		return nil, err
	}
	return &PayerReportManagerProtocolFeeRateUpdatedIterator{contract: _PayerReportManager.contract, event: "ProtocolFeeRateUpdated", logs: logs, sub: sub}, nil
}

// WatchProtocolFeeRateUpdated is a free log subscription operation binding the contract event 0x1b89f70e2fb3771e6a9f5cd2d2ad3439b823b67e0e69bdeb89b7d68a35fdd5c9.
//
// Solidity: event ProtocolFeeRateUpdated(uint16 protocolFeeRate)
func (_PayerReportManager *PayerReportManagerFilterer) WatchProtocolFeeRateUpdated(opts *bind.WatchOpts, sink chan<- *PayerReportManagerProtocolFeeRateUpdated) (event.Subscription, error) {

	logs, sub, err := _PayerReportManager.contract.WatchLogs(opts, "ProtocolFeeRateUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayerReportManagerProtocolFeeRateUpdated)
				if err := _PayerReportManager.contract.UnpackLog(event, "ProtocolFeeRateUpdated", log); err != nil {
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

// ParseProtocolFeeRateUpdated is a log parse operation binding the contract event 0x1b89f70e2fb3771e6a9f5cd2d2ad3439b823b67e0e69bdeb89b7d68a35fdd5c9.
//
// Solidity: event ProtocolFeeRateUpdated(uint16 protocolFeeRate)
func (_PayerReportManager *PayerReportManagerFilterer) ParseProtocolFeeRateUpdated(log types.Log) (*PayerReportManagerProtocolFeeRateUpdated, error) {
	event := new(PayerReportManagerProtocolFeeRateUpdated)
	if err := _PayerReportManager.contract.UnpackLog(event, "ProtocolFeeRateUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayerReportManagerUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the PayerReportManager contract.
type PayerReportManagerUpgradedIterator struct {
	Event *PayerReportManagerUpgraded // Event containing the contract specifics and raw log

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
func (it *PayerReportManagerUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayerReportManagerUpgraded)
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
		it.Event = new(PayerReportManagerUpgraded)
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
func (it *PayerReportManagerUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayerReportManagerUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayerReportManagerUpgraded represents a Upgraded event raised by the PayerReportManager contract.
type PayerReportManagerUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_PayerReportManager *PayerReportManagerFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*PayerReportManagerUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _PayerReportManager.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &PayerReportManagerUpgradedIterator{contract: _PayerReportManager.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_PayerReportManager *PayerReportManagerFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *PayerReportManagerUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _PayerReportManager.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayerReportManagerUpgraded)
				if err := _PayerReportManager.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_PayerReportManager *PayerReportManagerFilterer) ParseUpgraded(log types.Log) (*PayerReportManagerUpgraded, error) {
	event := new(PayerReportManagerUpgraded)
	if err := _PayerReportManager.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
