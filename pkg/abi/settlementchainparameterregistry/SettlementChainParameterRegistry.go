// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package settlementchainparameterregistry

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

// SettlementChainParameterRegistryMetaData contains all meta data concerning the SettlementChainParameterRegistry contract.
var SettlementChainParameterRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"adminParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"get\",\"inputs\":[{\"name\":\"keys_\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"values_\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"get\",\"inputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"value_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"implementation_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"admins_\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isAdmin\",\"inputs\":[{\"name\":\"account_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"isAdmin_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"migrate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migratorParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"set\",\"inputs\":[{\"name\":\"keys_\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"values_\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"set\",\"inputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Migrated\",\"inputs\":[{\"name\":\"migrator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ParameterSet\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ArrayLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyCode\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationFailed\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertData_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"NoKeyComponents\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoKeys\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterOutOfTypeBounds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StringsInsufficientHexLength\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ZeroMigrator\",\"inputs\":[]}]",
	Bin: "0x6080604052348015600e575f5ffd5b5060156019565b60c9565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff161560685760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b039081161460c65780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b61123f806100d65f395ff3fe608060405234801561000f575f5ffd5b50600436106100b9575f3560e01c80638aab82ba116100725780639f40b625116100585780639f40b62514610191578063a224cee714610199578063d6d7d525146101ac575f5ffd5b80638aab82ba146101745780638fd3ab8014610189575f5ffd5b806324d7806c116100a257806324d7806c146100e55780635c60da1b1461010d5780637257a3e314610154575f5ffd5b80631df893cc146100bd57806323d56420146100d2575b5f5ffd5b6100d06100cb366004610dac565b6101cd565b005b6100d06100e0366004610e56565b6102ff565b6100f86100f3366004610e9e565b61036d565b60405190151581526020015b60405180910390f35b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610104565b610167610162366004610ed1565b610390565b6040516101049190610f10565b61017c6104ab565b6040516101049190610f9e565b6100d06104cb565b61017c6104e8565b6100d06101a7366004610ed1565b610508565b6101bf6101ba366004610fb0565b6106e5565b604051908152602001610104565b6101d561072f565b5f83900361020f576040517f0190983500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b828114610248576040517fa24a13a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7fefab3f4eb315eafaa267b58974a509c07c739fbfe8e62b4eff49c4ced69850005f5b848110156102f7576102ef8287878481811061028957610289610fe3565b905060200281019061029b9190611010565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284375f920191909152508992508891508690508181106102e3576102e3610fe3565b9050602002013561076e565b60010161026b565b505050505050565b61030761072f565b6103687fefab3f4eb315eafaa267b58974a509c07c739fbfe8e62b4eff49c4ced698500084848080601f0160208091040260200160405190810160405280939291908181526020018383808284375f9201919091525086925061076e915050565b505050565b5f8061038861038361037d6104e8565b856107d1565b6107ec565b141592915050565b60605f8290036103cc576040517f0190983500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8167ffffffffffffffff8111156103e5576103e5611071565b60405190808252806020026020018201604052801561040e578160200160208202803683370190505b5090507fefab3f4eb315eafaa267b58974a509c07c739fbfe8e62b4eff49c4ced69850005f5b838110156104a3578185858381811061044f5761044f610fe3565b90506020028101906104619190611010565b60405161046f92919061109e565b90815260200160405180910390205483828151811061049057610490610fe3565b6020908102919091010152600101610434565b505092915050565b60606040518060600160405280602e81526020016111dc602e9139905090565b6104e66104e16104dc6103836104ab565b610832565b610885565b565b60606040518060600160405280602d81526020016111af602d9139905090565b5f610511610a90565b805490915060ff68010000000000000000820416159067ffffffffffffffff165f8115801561053d5750825b90505f8267ffffffffffffffff1660011480156105595750303b155b905081158015610567575080155b1561059e576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b84547fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000016600117855583156105ff5784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff16680100000000000000001785555b7fefab3f4eb315eafaa267b58974a509c07c739fbfe8e62b4eff49c4ced69850005f6106296104e8565b90505f5b888110156106785761067083610669848d8d8681811061064f5761064f610fe3565b90506020020160208101906106649190610e9e565b6107d1565b600161076e565b60010161062d565b50505083156106dc5784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50505050505050565b5f7fefab3f4eb315eafaa267b58974a509c07c739fbfe8e62b4eff49c4ced6985000604051610717908590859061109e565b90815260200160405180910390205490505b92915050565b6107383361036d565b6104e6576040517f7bfa4b9f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80835f018360405161078091906110c4565b90815260405190819003602001812082905561079d9084906110c4565b604051908190038120907f9577bebd11e9d897c6432d7db290e4a2101d3e13e93e0bb00ca291c37ff6bc54905f90a3505050565b60606107e5836107e084610ab8565b610ac3565b9392505050565b5f7fefab3f4eb315eafaa267b58974a509c07c739fbfe8e62b4eff49c4ced698500060405161081c9084906110c4565b9081526020016040518091039020549050919050565b5f73ffffffffffffffffffffffffffffffffffffffff821115610881576040517f37f4f14800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5090565b73ffffffffffffffffffffffffffffffffffffffff81166108d2576040517f0d626a3200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405173ffffffffffffffffffffffffffffffffffffffff8216907fa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098905f90a25f5f8273ffffffffffffffffffffffffffffffffffffffff166040515f60405180830381855af49150503d805f8114610966576040519150601f19603f3d011682016040523d82523d5f602084013e61096b565b606091505b5091509150816109b45782816040517f68b0b16b0000000000000000000000000000000000000000000000000000000081526004016109ab9291906110cf565b60405180910390fd5b80511580156109d8575073ffffffffffffffffffffffffffffffffffffffff83163b155b15610a27576040517f626c416100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff841660048201526024016109ab565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5473ffffffffffffffffffffffffffffffffffffffff167fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b60405160405180910390a2505050565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00610729565b606061072982610b26565b6060826040518060400160405280600181526020017f2e0000000000000000000000000000000000000000000000000000000000000081525083604051602001610b0f939291906110fd565b604051602081830303815290604052905092915050565b606061072973ffffffffffffffffffffffffffffffffffffffff831660146060825f610b53846002611150565b610b5e906002611167565b67ffffffffffffffff811115610b7657610b76611071565b6040519080825280601f01601f191660200182016040528015610ba0576020820181803683370190505b5090507f3000000000000000000000000000000000000000000000000000000000000000815f81518110610bd657610bd6610fe3565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191690815f1a9053507f780000000000000000000000000000000000000000000000000000000000000081600181518110610c3857610c38610fe3565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191690815f1a9053505f610c72856002611150565b610c7d906001611167565b90505b6001811115610d19577f303132333435363738396162636465660000000000000000000000000000000083600f1660108110610cbe57610cbe610fe3565b1a60f81b828281518110610cd457610cd4610fe3565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191690815f1a90535060049290921c91610d128161117a565b9050610c80565b508115610d5c576040517fe22e27eb00000000000000000000000000000000000000000000000000000000815260048101869052602481018590526044016109ab565b949350505050565b5f5f83601f840112610d74575f5ffd5b50813567ffffffffffffffff811115610d8b575f5ffd5b6020830191508360208260051b8501011115610da5575f5ffd5b9250929050565b5f5f5f5f60408587031215610dbf575f5ffd5b843567ffffffffffffffff811115610dd5575f5ffd5b610de187828801610d64565b909550935050602085013567ffffffffffffffff811115610e00575f5ffd5b610e0c87828801610d64565b95989497509550505050565b5f5f83601f840112610e28575f5ffd5b50813567ffffffffffffffff811115610e3f575f5ffd5b602083019150836020828501011115610da5575f5ffd5b5f5f5f60408486031215610e68575f5ffd5b833567ffffffffffffffff811115610e7e575f5ffd5b610e8a86828701610e18565b909790965060209590950135949350505050565b5f60208284031215610eae575f5ffd5b813573ffffffffffffffffffffffffffffffffffffffff811681146107e5575f5ffd5b5f5f60208385031215610ee2575f5ffd5b823567ffffffffffffffff811115610ef8575f5ffd5b610f0485828601610d64565b90969095509350505050565b602080825282518282018190525f918401906040840190835b81811015610f47578351835260209384019390920191600101610f29565b509095945050505050565b5f81518084528060208401602086015e5f6020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081525f6107e56020830184610f52565b5f5f60208385031215610fc1575f5ffd5b823567ffffffffffffffff811115610fd7575f5ffd5b610f0485828601610e18565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f5f83357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112611043575f5ffd5b83018035915067ffffffffffffffff82111561105d575f5ffd5b602001915036819003821315610da5575f5ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b818382375f9101908152919050565b5f81518060208401855e5f93019283525090919050565b5f6107e582846110ad565b73ffffffffffffffffffffffffffffffffffffffff83168152604060208201525f610d5c6040830184610f52565b5f61111a61111461110e84886110ad565b866110ad565b846110ad565b95945050505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b808202811582820484141761072957610729611123565b8082018082111561072957610729611123565b5f8161118857611188611123565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff019056fe786d74702e736574746c656d656e74436861696e506172616d6574657252656769737472792e697341646d696e786d74702e736574746c656d656e74436861696e506172616d6574657252656769737472792e6d69677261746f72a26469706673582212202f8bde0545d0b7eadc02733bcbd2c9be69fb5c6615f405b50bb69d1dd919d3b764736f6c634300081c0033",
}

// SettlementChainParameterRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use SettlementChainParameterRegistryMetaData.ABI instead.
var SettlementChainParameterRegistryABI = SettlementChainParameterRegistryMetaData.ABI

// SettlementChainParameterRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SettlementChainParameterRegistryMetaData.Bin instead.
var SettlementChainParameterRegistryBin = SettlementChainParameterRegistryMetaData.Bin

// DeploySettlementChainParameterRegistry deploys a new Ethereum contract, binding an instance of SettlementChainParameterRegistry to it.
func DeploySettlementChainParameterRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SettlementChainParameterRegistry, error) {
	parsed, err := SettlementChainParameterRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SettlementChainParameterRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SettlementChainParameterRegistry{SettlementChainParameterRegistryCaller: SettlementChainParameterRegistryCaller{contract: contract}, SettlementChainParameterRegistryTransactor: SettlementChainParameterRegistryTransactor{contract: contract}, SettlementChainParameterRegistryFilterer: SettlementChainParameterRegistryFilterer{contract: contract}}, nil
}

// SettlementChainParameterRegistry is an auto generated Go binding around an Ethereum contract.
type SettlementChainParameterRegistry struct {
	SettlementChainParameterRegistryCaller     // Read-only binding to the contract
	SettlementChainParameterRegistryTransactor // Write-only binding to the contract
	SettlementChainParameterRegistryFilterer   // Log filterer for contract events
}

// SettlementChainParameterRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type SettlementChainParameterRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementChainParameterRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SettlementChainParameterRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementChainParameterRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SettlementChainParameterRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementChainParameterRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SettlementChainParameterRegistrySession struct {
	Contract     *SettlementChainParameterRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                     // Call options to use throughout this session
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// SettlementChainParameterRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SettlementChainParameterRegistryCallerSession struct {
	Contract *SettlementChainParameterRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                           // Call options to use throughout this session
}

// SettlementChainParameterRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SettlementChainParameterRegistryTransactorSession struct {
	Contract     *SettlementChainParameterRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                           // Transaction auth options to use throughout this session
}

// SettlementChainParameterRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type SettlementChainParameterRegistryRaw struct {
	Contract *SettlementChainParameterRegistry // Generic contract binding to access the raw methods on
}

// SettlementChainParameterRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SettlementChainParameterRegistryCallerRaw struct {
	Contract *SettlementChainParameterRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// SettlementChainParameterRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SettlementChainParameterRegistryTransactorRaw struct {
	Contract *SettlementChainParameterRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSettlementChainParameterRegistry creates a new instance of SettlementChainParameterRegistry, bound to a specific deployed contract.
func NewSettlementChainParameterRegistry(address common.Address, backend bind.ContractBackend) (*SettlementChainParameterRegistry, error) {
	contract, err := bindSettlementChainParameterRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SettlementChainParameterRegistry{SettlementChainParameterRegistryCaller: SettlementChainParameterRegistryCaller{contract: contract}, SettlementChainParameterRegistryTransactor: SettlementChainParameterRegistryTransactor{contract: contract}, SettlementChainParameterRegistryFilterer: SettlementChainParameterRegistryFilterer{contract: contract}}, nil
}

// NewSettlementChainParameterRegistryCaller creates a new read-only instance of SettlementChainParameterRegistry, bound to a specific deployed contract.
func NewSettlementChainParameterRegistryCaller(address common.Address, caller bind.ContractCaller) (*SettlementChainParameterRegistryCaller, error) {
	contract, err := bindSettlementChainParameterRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SettlementChainParameterRegistryCaller{contract: contract}, nil
}

// NewSettlementChainParameterRegistryTransactor creates a new write-only instance of SettlementChainParameterRegistry, bound to a specific deployed contract.
func NewSettlementChainParameterRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*SettlementChainParameterRegistryTransactor, error) {
	contract, err := bindSettlementChainParameterRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SettlementChainParameterRegistryTransactor{contract: contract}, nil
}

// NewSettlementChainParameterRegistryFilterer creates a new log filterer instance of SettlementChainParameterRegistry, bound to a specific deployed contract.
func NewSettlementChainParameterRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*SettlementChainParameterRegistryFilterer, error) {
	contract, err := bindSettlementChainParameterRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SettlementChainParameterRegistryFilterer{contract: contract}, nil
}

// bindSettlementChainParameterRegistry binds a generic wrapper to an already deployed contract.
func bindSettlementChainParameterRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SettlementChainParameterRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SettlementChainParameterRegistry.Contract.SettlementChainParameterRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.SettlementChainParameterRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.SettlementChainParameterRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SettlementChainParameterRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.contract.Transact(opts, method, params...)
}

// AdminParameterKey is a free data retrieval call binding the contract method 0x9f40b625.
//
// Solidity: function adminParameterKey() pure returns(bytes key_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCaller) AdminParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _SettlementChainParameterRegistry.contract.Call(opts, &out, "adminParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// AdminParameterKey is a free data retrieval call binding the contract method 0x9f40b625.
//
// Solidity: function adminParameterKey() pure returns(bytes key_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistrySession) AdminParameterKey() ([]byte, error) {
	return _SettlementChainParameterRegistry.Contract.AdminParameterKey(&_SettlementChainParameterRegistry.CallOpts)
}

// AdminParameterKey is a free data retrieval call binding the contract method 0x9f40b625.
//
// Solidity: function adminParameterKey() pure returns(bytes key_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCallerSession) AdminParameterKey() ([]byte, error) {
	return _SettlementChainParameterRegistry.Contract.AdminParameterKey(&_SettlementChainParameterRegistry.CallOpts)
}

// Get is a free data retrieval call binding the contract method 0x7257a3e3.
//
// Solidity: function get(bytes[] keys_) view returns(bytes32[] values_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCaller) Get(opts *bind.CallOpts, keys_ [][]byte) ([][32]byte, error) {
	var out []interface{}
	err := _SettlementChainParameterRegistry.contract.Call(opts, &out, "get", keys_)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// Get is a free data retrieval call binding the contract method 0x7257a3e3.
//
// Solidity: function get(bytes[] keys_) view returns(bytes32[] values_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistrySession) Get(keys_ [][]byte) ([][32]byte, error) {
	return _SettlementChainParameterRegistry.Contract.Get(&_SettlementChainParameterRegistry.CallOpts, keys_)
}

// Get is a free data retrieval call binding the contract method 0x7257a3e3.
//
// Solidity: function get(bytes[] keys_) view returns(bytes32[] values_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCallerSession) Get(keys_ [][]byte) ([][32]byte, error) {
	return _SettlementChainParameterRegistry.Contract.Get(&_SettlementChainParameterRegistry.CallOpts, keys_)
}

// Get0 is a free data retrieval call binding the contract method 0xd6d7d525.
//
// Solidity: function get(bytes key_) view returns(bytes32 value_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCaller) Get0(opts *bind.CallOpts, key_ []byte) ([32]byte, error) {
	var out []interface{}
	err := _SettlementChainParameterRegistry.contract.Call(opts, &out, "get0", key_)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Get0 is a free data retrieval call binding the contract method 0xd6d7d525.
//
// Solidity: function get(bytes key_) view returns(bytes32 value_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistrySession) Get0(key_ []byte) ([32]byte, error) {
	return _SettlementChainParameterRegistry.Contract.Get0(&_SettlementChainParameterRegistry.CallOpts, key_)
}

// Get0 is a free data retrieval call binding the contract method 0xd6d7d525.
//
// Solidity: function get(bytes key_) view returns(bytes32 value_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCallerSession) Get0(key_ []byte) ([32]byte, error) {
	return _SettlementChainParameterRegistry.Contract.Get0(&_SettlementChainParameterRegistry.CallOpts, key_)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SettlementChainParameterRegistry.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistrySession) Implementation() (common.Address, error) {
	return _SettlementChainParameterRegistry.Contract.Implementation(&_SettlementChainParameterRegistry.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCallerSession) Implementation() (common.Address, error) {
	return _SettlementChainParameterRegistry.Contract.Implementation(&_SettlementChainParameterRegistry.CallOpts)
}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address account_) view returns(bool isAdmin_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCaller) IsAdmin(opts *bind.CallOpts, account_ common.Address) (bool, error) {
	var out []interface{}
	err := _SettlementChainParameterRegistry.contract.Call(opts, &out, "isAdmin", account_)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address account_) view returns(bool isAdmin_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistrySession) IsAdmin(account_ common.Address) (bool, error) {
	return _SettlementChainParameterRegistry.Contract.IsAdmin(&_SettlementChainParameterRegistry.CallOpts, account_)
}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address account_) view returns(bool isAdmin_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCallerSession) IsAdmin(account_ common.Address) (bool, error) {
	return _SettlementChainParameterRegistry.Contract.IsAdmin(&_SettlementChainParameterRegistry.CallOpts, account_)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(bytes key_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCaller) MigratorParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _SettlementChainParameterRegistry.contract.Call(opts, &out, "migratorParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(bytes key_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistrySession) MigratorParameterKey() ([]byte, error) {
	return _SettlementChainParameterRegistry.Contract.MigratorParameterKey(&_SettlementChainParameterRegistry.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(bytes key_)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryCallerSession) MigratorParameterKey() ([]byte, error) {
	return _SettlementChainParameterRegistry.Contract.MigratorParameterKey(&_SettlementChainParameterRegistry.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xa224cee7.
//
// Solidity: function initialize(address[] admins_) returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryTransactor) Initialize(opts *bind.TransactOpts, admins_ []common.Address) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.contract.Transact(opts, "initialize", admins_)
}

// Initialize is a paid mutator transaction binding the contract method 0xa224cee7.
//
// Solidity: function initialize(address[] admins_) returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistrySession) Initialize(admins_ []common.Address) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.Initialize(&_SettlementChainParameterRegistry.TransactOpts, admins_)
}

// Initialize is a paid mutator transaction binding the contract method 0xa224cee7.
//
// Solidity: function initialize(address[] admins_) returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryTransactorSession) Initialize(admins_ []common.Address) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.Initialize(&_SettlementChainParameterRegistry.TransactOpts, admins_)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryTransactor) Migrate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.contract.Transact(opts, "migrate")
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistrySession) Migrate() (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.Migrate(&_SettlementChainParameterRegistry.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryTransactorSession) Migrate() (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.Migrate(&_SettlementChainParameterRegistry.TransactOpts)
}

// Set is a paid mutator transaction binding the contract method 0x1df893cc.
//
// Solidity: function set(bytes[] keys_, bytes32[] values_) returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryTransactor) Set(opts *bind.TransactOpts, keys_ [][]byte, values_ [][32]byte) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.contract.Transact(opts, "set", keys_, values_)
}

// Set is a paid mutator transaction binding the contract method 0x1df893cc.
//
// Solidity: function set(bytes[] keys_, bytes32[] values_) returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistrySession) Set(keys_ [][]byte, values_ [][32]byte) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.Set(&_SettlementChainParameterRegistry.TransactOpts, keys_, values_)
}

// Set is a paid mutator transaction binding the contract method 0x1df893cc.
//
// Solidity: function set(bytes[] keys_, bytes32[] values_) returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryTransactorSession) Set(keys_ [][]byte, values_ [][32]byte) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.Set(&_SettlementChainParameterRegistry.TransactOpts, keys_, values_)
}

// Set0 is a paid mutator transaction binding the contract method 0x23d56420.
//
// Solidity: function set(bytes key_, bytes32 value_) returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryTransactor) Set0(opts *bind.TransactOpts, key_ []byte, value_ [32]byte) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.contract.Transact(opts, "set0", key_, value_)
}

// Set0 is a paid mutator transaction binding the contract method 0x23d56420.
//
// Solidity: function set(bytes key_, bytes32 value_) returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistrySession) Set0(key_ []byte, value_ [32]byte) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.Set0(&_SettlementChainParameterRegistry.TransactOpts, key_, value_)
}

// Set0 is a paid mutator transaction binding the contract method 0x23d56420.
//
// Solidity: function set(bytes key_, bytes32 value_) returns()
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryTransactorSession) Set0(key_ []byte, value_ [32]byte) (*types.Transaction, error) {
	return _SettlementChainParameterRegistry.Contract.Set0(&_SettlementChainParameterRegistry.TransactOpts, key_, value_)
}

// SettlementChainParameterRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the SettlementChainParameterRegistry contract.
type SettlementChainParameterRegistryInitializedIterator struct {
	Event *SettlementChainParameterRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *SettlementChainParameterRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementChainParameterRegistryInitialized)
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
		it.Event = new(SettlementChainParameterRegistryInitialized)
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
func (it *SettlementChainParameterRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementChainParameterRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementChainParameterRegistryInitialized represents a Initialized event raised by the SettlementChainParameterRegistry contract.
type SettlementChainParameterRegistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*SettlementChainParameterRegistryInitializedIterator, error) {

	logs, sub, err := _SettlementChainParameterRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &SettlementChainParameterRegistryInitializedIterator{contract: _SettlementChainParameterRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *SettlementChainParameterRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _SettlementChainParameterRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementChainParameterRegistryInitialized)
				if err := _SettlementChainParameterRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) ParseInitialized(log types.Log) (*SettlementChainParameterRegistryInitialized, error) {
	event := new(SettlementChainParameterRegistryInitialized)
	if err := _SettlementChainParameterRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementChainParameterRegistryMigratedIterator is returned from FilterMigrated and is used to iterate over the raw logs and unpacked data for Migrated events raised by the SettlementChainParameterRegistry contract.
type SettlementChainParameterRegistryMigratedIterator struct {
	Event *SettlementChainParameterRegistryMigrated // Event containing the contract specifics and raw log

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
func (it *SettlementChainParameterRegistryMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementChainParameterRegistryMigrated)
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
		it.Event = new(SettlementChainParameterRegistryMigrated)
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
func (it *SettlementChainParameterRegistryMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementChainParameterRegistryMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementChainParameterRegistryMigrated represents a Migrated event raised by the SettlementChainParameterRegistry contract.
type SettlementChainParameterRegistryMigrated struct {
	Migrator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMigrated is a free log retrieval operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) FilterMigrated(opts *bind.FilterOpts, migrator []common.Address) (*SettlementChainParameterRegistryMigratedIterator, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _SettlementChainParameterRegistry.contract.FilterLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return &SettlementChainParameterRegistryMigratedIterator{contract: _SettlementChainParameterRegistry.contract, event: "Migrated", logs: logs, sub: sub}, nil
}

// WatchMigrated is a free log subscription operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) WatchMigrated(opts *bind.WatchOpts, sink chan<- *SettlementChainParameterRegistryMigrated, migrator []common.Address) (event.Subscription, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _SettlementChainParameterRegistry.contract.WatchLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementChainParameterRegistryMigrated)
				if err := _SettlementChainParameterRegistry.contract.UnpackLog(event, "Migrated", log); err != nil {
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
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) ParseMigrated(log types.Log) (*SettlementChainParameterRegistryMigrated, error) {
	event := new(SettlementChainParameterRegistryMigrated)
	if err := _SettlementChainParameterRegistry.contract.UnpackLog(event, "Migrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementChainParameterRegistryParameterSetIterator is returned from FilterParameterSet and is used to iterate over the raw logs and unpacked data for ParameterSet events raised by the SettlementChainParameterRegistry contract.
type SettlementChainParameterRegistryParameterSetIterator struct {
	Event *SettlementChainParameterRegistryParameterSet // Event containing the contract specifics and raw log

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
func (it *SettlementChainParameterRegistryParameterSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementChainParameterRegistryParameterSet)
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
		it.Event = new(SettlementChainParameterRegistryParameterSet)
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
func (it *SettlementChainParameterRegistryParameterSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementChainParameterRegistryParameterSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementChainParameterRegistryParameterSet represents a ParameterSet event raised by the SettlementChainParameterRegistry contract.
type SettlementChainParameterRegistryParameterSet struct {
	Key   common.Hash
	Value [32]byte
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterParameterSet is a free log retrieval operation binding the contract event 0x9577bebd11e9d897c6432d7db290e4a2101d3e13e93e0bb00ca291c37ff6bc54.
//
// Solidity: event ParameterSet(bytes indexed key, bytes32 indexed value)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) FilterParameterSet(opts *bind.FilterOpts, key [][]byte, value [][32]byte) (*SettlementChainParameterRegistryParameterSetIterator, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}
	var valueRule []interface{}
	for _, valueItem := range value {
		valueRule = append(valueRule, valueItem)
	}

	logs, sub, err := _SettlementChainParameterRegistry.contract.FilterLogs(opts, "ParameterSet", keyRule, valueRule)
	if err != nil {
		return nil, err
	}
	return &SettlementChainParameterRegistryParameterSetIterator{contract: _SettlementChainParameterRegistry.contract, event: "ParameterSet", logs: logs, sub: sub}, nil
}

// WatchParameterSet is a free log subscription operation binding the contract event 0x9577bebd11e9d897c6432d7db290e4a2101d3e13e93e0bb00ca291c37ff6bc54.
//
// Solidity: event ParameterSet(bytes indexed key, bytes32 indexed value)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) WatchParameterSet(opts *bind.WatchOpts, sink chan<- *SettlementChainParameterRegistryParameterSet, key [][]byte, value [][32]byte) (event.Subscription, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}
	var valueRule []interface{}
	for _, valueItem := range value {
		valueRule = append(valueRule, valueItem)
	}

	logs, sub, err := _SettlementChainParameterRegistry.contract.WatchLogs(opts, "ParameterSet", keyRule, valueRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementChainParameterRegistryParameterSet)
				if err := _SettlementChainParameterRegistry.contract.UnpackLog(event, "ParameterSet", log); err != nil {
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

// ParseParameterSet is a log parse operation binding the contract event 0x9577bebd11e9d897c6432d7db290e4a2101d3e13e93e0bb00ca291c37ff6bc54.
//
// Solidity: event ParameterSet(bytes indexed key, bytes32 indexed value)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) ParseParameterSet(log types.Log) (*SettlementChainParameterRegistryParameterSet, error) {
	event := new(SettlementChainParameterRegistryParameterSet)
	if err := _SettlementChainParameterRegistry.contract.UnpackLog(event, "ParameterSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementChainParameterRegistryUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the SettlementChainParameterRegistry contract.
type SettlementChainParameterRegistryUpgradedIterator struct {
	Event *SettlementChainParameterRegistryUpgraded // Event containing the contract specifics and raw log

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
func (it *SettlementChainParameterRegistryUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementChainParameterRegistryUpgraded)
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
		it.Event = new(SettlementChainParameterRegistryUpgraded)
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
func (it *SettlementChainParameterRegistryUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementChainParameterRegistryUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementChainParameterRegistryUpgraded represents a Upgraded event raised by the SettlementChainParameterRegistry contract.
type SettlementChainParameterRegistryUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*SettlementChainParameterRegistryUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _SettlementChainParameterRegistry.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &SettlementChainParameterRegistryUpgradedIterator{contract: _SettlementChainParameterRegistry.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *SettlementChainParameterRegistryUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _SettlementChainParameterRegistry.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementChainParameterRegistryUpgraded)
				if err := _SettlementChainParameterRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_SettlementChainParameterRegistry *SettlementChainParameterRegistryFilterer) ParseUpgraded(log types.Log) (*SettlementChainParameterRegistryUpgraded, error) {
	event := new(SettlementChainParameterRegistryUpgraded)
	if err := _SettlementChainParameterRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
