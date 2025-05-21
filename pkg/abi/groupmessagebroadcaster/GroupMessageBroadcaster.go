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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"parameterRegistry_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addMessage\",\"inputs\":[{\"name\":\"groupId_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"message_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"implementation_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"maxPayloadSize\",\"inputs\":[],\"outputs\":[{\"name\":\"size_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxPayloadSizeParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"migrate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migratorParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"minPayloadSize\",\"inputs\":[],\"outputs\":[{\"name\":\"size_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minPayloadSizeParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"parameterRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"paused_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pausedParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"updateMaxPayloadSize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMinPayloadSize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updatePauseStatus\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MaxPayloadSizeUpdated\",\"inputs\":[{\"name\":\"size\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MessageSent\",\"inputs\":[{\"name\":\"groupId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"message\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"sequenceId\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Migrated\",\"inputs\":[{\"name\":\"migrator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinPayloadSizeUpdated\",\"inputs\":[{\"name\":\"size\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PauseStatusUpdated\",\"inputs\":[{\"name\":\"paused\",\"type\":\"bool\",\"indexed\":true,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EmptyCode\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidMaxPayloadSize\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidMinPayloadSize\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPayloadSize\",\"inputs\":[{\"name\":\"actualSize_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minSize_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxSize_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MigrationFailed\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertData_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"NoChange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterOutOfTypeBounds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Paused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroMigrator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroParameterRegistry\",\"inputs\":[]}]",
	Bin: "0x60a060405234801561000f575f5ffd5b5060405161119b38038061119b83398101604081905261002e9161011d565b6001600160a01b0381166080819052819061005c5760405163d973fd8d60e01b815260040160405180910390fd5b61006461006b565b505061014a565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff16156100bb5760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b039081161461011a5780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b5f6020828403121561012d575f5ffd5b81516001600160a01b0381168114610143575f5ffd5b9392505050565b60805161101e61017d5f395f818160f9015281816103750152818161048f0152818161077301526107ec015261101e5ff3fe608060405234801561000f575f5ffd5b50600436106100f0575f3560e01c80635f643f93116100935780639218415d116100635780639218415d1461024d578063cc5999af14610255578063d46153ef1461025d578063f96927ac14610265575f5ffd5b80635f643f931461022d5780638129fc1c146102355780638aab82ba1461023d5780638fd3ab8014610245575f5ffd5b806358e3e94c116100ce57806358e3e94c1461016f57806359d4df41146101b55780635c60da1b146101bd5780635c975abb146101e4575f5ffd5b80630723499e146100f45780630cb858ea146101455780634dff26b51461015a575b5f5ffd5b61011b7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b61014d610292565b60405161013c9190610e28565b61016d610168366004610e3a565b6102b2565b005b7feda186f2b85b2c197e0a3ff15dc0c5c16c74d00b5c7f432acaa215db84203b0054640100000000900463ffffffff165b60405163ffffffff909116815260200161013c565b61016d61036f565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5461011b565b7feda186f2b85b2c197e0a3ff15dc0c5c16c74d00b5c7f432acaa215db84203b0054700100000000000000000000000000000000900460ff16604051901515815260200161013c565b61016d610489565b61016d6105d7565b61014d61074b565b61016d61076b565b61014d6107a6565b61014d6107c6565b61016d6107e6565b7feda186f2b85b2c197e0a3ff15dc0c5c16c74d00b5c7f432acaa215db84203b005463ffffffff166101a0565b60606040518060600160405280602b8152602001610f70602b9139905090565b6102ba610925565b6102c381610995565b7feda186f2b85b2c197e0a3ff15dc0c5c16c74d00b5c7f432acaa215db84203b0080547fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff8116680100000000000000009182900467ffffffffffffffff908116600101169182021790915560405184907f91f47151424884a46811ed593aa8a02ee5012e9332a4dcf1e9236a8ed4443c3e906103629086908690610eb1565b60405180910390a3505050565b5f6103a17f000000000000000000000000000000000000000000000000000000000000000061039c6107c6565b610a3c565b7feda186f2b85b2c197e0a3ff15dc0c5c16c74d00b5c7f432acaa215db84203b00805491925090700100000000000000000000000000000000900460ff1615158215150361041b576040517fa88ee57700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80547fffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffff167001000000000000000000000000000000008315159081029190911782556040517f7c4d1fe30fdbfda9e9c4c43e759ef32e4db5128d4cb58ff3ae9583b89b6242a5905f90a25050565b5f6104bb7f00000000000000000000000000000000000000000000000000000000000000006104b6610292565b610a8e565b7feda186f2b85b2c197e0a3ff15dc0c5c16c74d00b5c7f432acaa215db84203b0080549192509063ffffffff9081169083161015610525576040517f1d8e7a4a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805463ffffffff640100000000909104811690831603610571576040517fa88ee57700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80547fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff1664010000000063ffffffff84169081029190911782556040517f62422e33fcfc9d38acda2bbddab282a9cc6df7e75f88269fd725bef5457b3045905f90a25050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156106215750825b90505f8267ffffffffffffffff16600114801561063d5750303b155b90508115801561064b575080155b15610682576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b84547fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000016600117855583156106e35784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff16680100000000000000001785555b83156107445784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b5050505050565b6060604051806060016040528060258152602001610f4b60259139905090565b6107a461079f7f000000000000000000000000000000000000000000000000000000000000000061079a61074b565b610ae1565b610af4565b565b60606040518060600160405280602b8152602001610f9b602b9139905090565b6060604051806060016040528060238152602001610fc660239139905090565b5f6108137f00000000000000000000000000000000000000000000000000000000000000006104b66107a6565b7feda186f2b85b2c197e0a3ff15dc0c5c16c74d00b5c7f432acaa215db84203b00805491925090640100000000900463ffffffff9081169083161115610885576040517fe219e4f000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805463ffffffff908116908316036108c9576040517fa88ee57700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff831690811782556040517f2caf5b55114860c563b52eba8026a6a0183d9eb1715cbf1c3f8b689f14b5121c905f90a25050565b7feda186f2b85b2c197e0a3ff15dc0c5c16c74d00b5c7f432acaa215db84203b0054700100000000000000000000000000000000900460ff16156107a4576040517f9e87fac800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7feda186f2b85b2c197e0a3ff15dc0c5c16c74d00b5c7f432acaa215db84203b00805463ffffffff168210806109d957508054640100000000900463ffffffff1682115b15610a385780546040517f93b7abe60000000000000000000000000000000000000000000000000000000081526004810184905263ffffffff808316602483015264010000000090920490911660448201526064015b60405180910390fd5b5050565b5f5f610a488484610cf6565b90506001811115610a85576040517f37f4f14800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b15159392505050565b5f5f610a9a8484610cf6565b905063ffffffff811115610ada576040517f37f4f14800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b9392505050565b5f610ada610aef8484610cf6565b610d89565b73ffffffffffffffffffffffffffffffffffffffff8116610b41576040517f0d626a3200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405173ffffffffffffffffffffffffffffffffffffffff8216907fa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098905f90a25f5f8273ffffffffffffffffffffffffffffffffffffffff166040515f60405180830381855af49150503d805f8114610bd5576040519150601f19603f3d011682016040523d82523d5f602084013e610bda565b606091505b509150915081610c1a5782816040517f68b0b16b000000000000000000000000000000000000000000000000000000008152600401610a2f929190610efd565b8051158015610c3e575073ffffffffffffffffffffffffffffffffffffffff83163b155b15610c8d576040517f626c416100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff84166004820152602401610a2f565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5473ffffffffffffffffffffffffffffffffffffffff167fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b60405160405180910390a2505050565b6040517fd6d7d5250000000000000000000000000000000000000000000000000000000081525f9073ffffffffffffffffffffffffffffffffffffffff84169063d6d7d52590610d4a908590600401610e28565b602060405180830381865afa158015610d65573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610ada9190610f33565b5f73ffffffffffffffffffffffffffffffffffffffff821115610dd8576040517f37f4f14800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5090565b5f81518084528060208401602086015e5f6020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081525f610ada6020830184610ddc565b5f5f5f60408486031215610e4c575f5ffd5b83359250602084013567ffffffffffffffff811115610e69575f5ffd5b8401601f81018613610e79575f5ffd5b803567ffffffffffffffff811115610e8f575f5ffd5b866020828401011115610ea0575f5ffd5b939660209190910195509293505050565b60208152816020820152818360408301375f818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b73ffffffffffffffffffffffffffffffffffffffff83168152604060208201525f610f2b6040830184610ddc565b949350505050565b5f60208284031215610f43575f5ffd5b505191905056fe786d74702e67726f75704d65737361676542726f61646361737465722e6d69677261746f72786d74702e67726f75704d65737361676542726f61646361737465722e6d61785061796c6f616453697a65786d74702e67726f75704d65737361676542726f61646361737465722e6d696e5061796c6f616453697a65786d74702e67726f75704d65737361676542726f61646361737465722e706175736564a264697066735822122030d6e7c4151db3a26ccad04a3435ccc5d2493d5036f73bba1c46b03927c82d0e64736f6c634300081c0033",
}

// GroupMessageBroadcasterABI is the input ABI used to generate the binding from.
// Deprecated: Use GroupMessageBroadcasterMetaData.ABI instead.
var GroupMessageBroadcasterABI = GroupMessageBroadcasterMetaData.ABI

// GroupMessageBroadcasterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GroupMessageBroadcasterMetaData.Bin instead.
var GroupMessageBroadcasterBin = GroupMessageBroadcasterMetaData.Bin

// DeployGroupMessageBroadcaster deploys a new Ethereum contract, binding an instance of GroupMessageBroadcaster to it.
func DeployGroupMessageBroadcaster(auth *bind.TransactOpts, backend bind.ContractBackend, parameterRegistry_ common.Address) (common.Address, *types.Transaction, *GroupMessageBroadcaster, error) {
	parsed, err := GroupMessageBroadcasterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GroupMessageBroadcasterBin), backend, parameterRegistry_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &GroupMessageBroadcaster{GroupMessageBroadcasterCaller: GroupMessageBroadcasterCaller{contract: contract}, GroupMessageBroadcasterTransactor: GroupMessageBroadcasterTransactor{contract: contract}, GroupMessageBroadcasterFilterer: GroupMessageBroadcasterFilterer{contract: contract}}, nil
}

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
// Solidity: function maxPayloadSizeParameterKey() pure returns(bytes key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) MaxPayloadSizeParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "maxPayloadSizeParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// MaxPayloadSizeParameterKey is a free data retrieval call binding the contract method 0x0cb858ea.
//
// Solidity: function maxPayloadSizeParameterKey() pure returns(bytes key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) MaxPayloadSizeParameterKey() ([]byte, error) {
	return _GroupMessageBroadcaster.Contract.MaxPayloadSizeParameterKey(&_GroupMessageBroadcaster.CallOpts)
}

// MaxPayloadSizeParameterKey is a free data retrieval call binding the contract method 0x0cb858ea.
//
// Solidity: function maxPayloadSizeParameterKey() pure returns(bytes key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) MaxPayloadSizeParameterKey() ([]byte, error) {
	return _GroupMessageBroadcaster.Contract.MaxPayloadSizeParameterKey(&_GroupMessageBroadcaster.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(bytes key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) MigratorParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "migratorParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(bytes key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) MigratorParameterKey() ([]byte, error) {
	return _GroupMessageBroadcaster.Contract.MigratorParameterKey(&_GroupMessageBroadcaster.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(bytes key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) MigratorParameterKey() ([]byte, error) {
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
// Solidity: function minPayloadSizeParameterKey() pure returns(bytes key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) MinPayloadSizeParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "minPayloadSizeParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// MinPayloadSizeParameterKey is a free data retrieval call binding the contract method 0x9218415d.
//
// Solidity: function minPayloadSizeParameterKey() pure returns(bytes key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) MinPayloadSizeParameterKey() ([]byte, error) {
	return _GroupMessageBroadcaster.Contract.MinPayloadSizeParameterKey(&_GroupMessageBroadcaster.CallOpts)
}

// MinPayloadSizeParameterKey is a free data retrieval call binding the contract method 0x9218415d.
//
// Solidity: function minPayloadSizeParameterKey() pure returns(bytes key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) MinPayloadSizeParameterKey() ([]byte, error) {
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
// Solidity: function pausedParameterKey() pure returns(bytes key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) PausedParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "pausedParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// PausedParameterKey is a free data retrieval call binding the contract method 0xcc5999af.
//
// Solidity: function pausedParameterKey() pure returns(bytes key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) PausedParameterKey() ([]byte, error) {
	return _GroupMessageBroadcaster.Contract.PausedParameterKey(&_GroupMessageBroadcaster.CallOpts)
}

// PausedParameterKey is a free data retrieval call binding the contract method 0xcc5999af.
//
// Solidity: function pausedParameterKey() pure returns(bytes key_)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) PausedParameterKey() ([]byte, error) {
	return _GroupMessageBroadcaster.Contract.PausedParameterKey(&_GroupMessageBroadcaster.CallOpts)
}

// AddMessage is a paid mutator transaction binding the contract method 0x4dff26b5.
//
// Solidity: function addMessage(bytes32 groupId_, bytes message_) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactor) AddMessage(opts *bind.TransactOpts, groupId_ [32]byte, message_ []byte) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.contract.Transact(opts, "addMessage", groupId_, message_)
}

// AddMessage is a paid mutator transaction binding the contract method 0x4dff26b5.
//
// Solidity: function addMessage(bytes32 groupId_, bytes message_) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) AddMessage(groupId_ [32]byte, message_ []byte) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.AddMessage(&_GroupMessageBroadcaster.TransactOpts, groupId_, message_)
}

// AddMessage is a paid mutator transaction binding the contract method 0x4dff26b5.
//
// Solidity: function addMessage(bytes32 groupId_, bytes message_) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorSession) AddMessage(groupId_ [32]byte, message_ []byte) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.AddMessage(&_GroupMessageBroadcaster.TransactOpts, groupId_, message_)
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
	GroupId    [32]byte
	Message    []byte
	SequenceId uint64
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterMessageSent is a free log retrieval operation binding the contract event 0x91f47151424884a46811ed593aa8a02ee5012e9332a4dcf1e9236a8ed4443c3e.
//
// Solidity: event MessageSent(bytes32 indexed groupId, bytes message, uint64 indexed sequenceId)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) FilterMessageSent(opts *bind.FilterOpts, groupId [][32]byte, sequenceId []uint64) (*GroupMessageBroadcasterMessageSentIterator, error) {

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

// WatchMessageSent is a free log subscription operation binding the contract event 0x91f47151424884a46811ed593aa8a02ee5012e9332a4dcf1e9236a8ed4443c3e.
//
// Solidity: event MessageSent(bytes32 indexed groupId, bytes message, uint64 indexed sequenceId)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) WatchMessageSent(opts *bind.WatchOpts, sink chan<- *GroupMessageBroadcasterMessageSent, groupId [][32]byte, sequenceId []uint64) (event.Subscription, error) {

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

// ParseMessageSent is a log parse operation binding the contract event 0x91f47151424884a46811ed593aa8a02ee5012e9332a4dcf1e9236a8ed4443c3e.
//
// Solidity: event MessageSent(bytes32 indexed groupId, bytes message, uint64 indexed sequenceId)
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
