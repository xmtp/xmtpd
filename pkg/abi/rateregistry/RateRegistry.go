// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rateregistry

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

// IRateRegistryRates is an auto generated low-level Go binding around an user-defined struct.
type IRateRegistryRates struct {
	MessageFee          uint64
	StorageFee          uint64
	CongestionFee       uint64
	TargetRatePerMinute uint64
	StartTime           uint64
}

// RateRegistryMetaData contains all meta data concerning the RateRegistry contract.
var RateRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"parameterRegistry_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"congestionFeeParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getRates\",\"inputs\":[{\"name\":\"fromIndex_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"count_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"rates_\",\"type\":\"tuple[]\",\"internalType\":\"structIRateRegistry.Rates[]\",\"components\":[{\"name\":\"messageFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"storageFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"congestionFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"targetRatePerMinute\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRatesCount\",\"inputs\":[],\"outputs\":[{\"name\":\"count_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"implementation_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"messageFeeParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"migrate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migratorParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"parameterRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"storageFeeParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"targetRatePerMinuteParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"updateRates\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Migrated\",\"inputs\":[{\"name\":\"migrator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RatesUpdated\",\"inputs\":[{\"name\":\"messageFee\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"storageFee\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"congestionFee\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"targetRatePerMinute\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EmptyCode\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"EndIndexOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FromIndexOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationFailed\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertData_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"NoChange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterOutOfTypeBounds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroCount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroMigrator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroParameterRegistry\",\"inputs\":[]}]",
	Bin: "0x60a060405234801561000f575f5ffd5b5060405161135138038061135183398101604081905261002e9161011a565b6001600160a01b038116608081905261005a5760405163d973fd8d60e01b815260040160405180910390fd5b610062610068565b50610147565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff16156100b85760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146101175780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b5f6020828403121561012a575f5ffd5b81516001600160a01b0381168114610140575f5ffd5b9392505050565b6080516111d06101815f395f818160d8015281816102d80152818161033e0152818161039f0152818161040001526109e201526111d05ff3fe608060405234801561000f575f5ffd5b50600436106100cf575f3560e01c806363c032911161007d5780638fd3ab80116100585780638fd3ab8014610234578063ba3261d51461023c578063ed7e698614610275575f5ffd5b806363c03291146101ba5780638129fc1c146101f35780638aab82ba146101fb575f5ffd5b806345b05a43116100ad57806345b05a431461015e57806349156e091461017e5780635c60da1b14610193575f5ffd5b80630723499e146100d35780632da72291146101245780633c3821f414610154575b5f5ffd5b6100fa7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b7f988e236e2caf5758fdf811320ba1d2fca453cb71bd6049ebba876b68af5050005460405190815260200161011b565b61015c6102ae565b005b61017161016c366004610f59565b61060d565b60405161011b9190610f79565b61018661085b565b60405161011b9190611069565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546100fa565b60408051808201909152601c81527f786d74702e7261746552656769737472792e6d657373616765466565000000006020820152610186565b61015c61087b565b60408051808201909152601a81527f786d74702e7261746552656769737472792e6d69677261746f720000000000006020820152610186565b61015c6109da565b60408051808201909152601c81527f786d74702e7261746552656769737472792e73746f72616765466565000000006020820152610186565b60408051808201909152601f81527f786d74702e7261746552656769737472792e636f6e67657374696f6e466565006020820152610186565b5f7f988e236e2caf5758fdf811320ba1d2fca453cb71bd6049ebba876b68af50500090505f6103367f000000000000000000000000000000000000000000000000000000000000000061033160408051808201909152601c81527f786d74702e7261746552656769737472792e6d65737361676546656500000000602082015290565b610a47565b90505f6103977f000000000000000000000000000000000000000000000000000000000000000061033160408051808201909152601c81527f786d74702e7261746552656769737472792e73746f7261676546656500000000602082015290565b90505f6103f87f000000000000000000000000000000000000000000000000000000000000000061033160408051808201909152601f81527f786d74702e7261746552656769737472792e636f6e67657374696f6e46656500602082015290565b90505f6104277f000000000000000000000000000000000000000000000000000000000000000061033161085b565b905061043584848484610aa0565b5f429050855f016040518060a001604052808767ffffffffffffffff1681526020018667ffffffffffffffff1681526020018567ffffffffffffffff1681526020018467ffffffffffffffff1681526020018367ffffffffffffffff16815250908060018154018082558091505060019003905f5260205f2090600202015f909190919091505f820151815f015f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055506020820151815f0160086101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055506040820151815f0160106101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055506060820151815f0160186101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055506080820151816001015f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555050507fabd2140443b16c364d95086ebf21b45137b5f0af53e10a6e792c0cb3d0e2db62858585856040516105fd949392919067ffffffffffffffff948516815292841660208401529083166040830152909116606082015260800190565b60405180910390a1505050505050565b6060815f03610648576040517f047b9cec00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f988e236e2caf5758fdf811320ba1d2fca453cb71bd6049ebba876b68af505000805484106106a3576040517fea61fe7000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80546106af84866110a8565b11156106e7576040517fb6cc753100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8267ffffffffffffffff811115610700576107006110bb565b60405190808252806020026020018201604052801561077657816020015b6040805160a0810182525f808252602080830182905292820181905260608201819052608082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90920191018161071e5790505b5091505f5b83811015610853578161078e82876110a8565b8154811061079e5761079e6110e8565b5f9182526020918290206040805160a0810182526002909302909101805467ffffffffffffffff808216855268010000000000000000820481169585019590955270010000000000000000000000000000000081048516928401929092527801000000000000000000000000000000000000000000000000909104831660608301526001015490911660808201528351849083908110610840576108406110e8565b602090810291909101015260010161077b565b505092915050565b606060405180606001604052806025815260200161117660259139905090565b5f610884610c2d565b805490915060ff68010000000000000000820416159067ffffffffffffffff165f811580156108b05750825b90505f8267ffffffffffffffff1660011480156108cc5750303b155b9050811580156108da575080155b15610911576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b84547fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000016600117855583156109725784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff16680100000000000000001785555b83156109d35784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b5050505050565b610a45610a407f0000000000000000000000000000000000000000000000000000000000000000610a3b60408051808201909152601a81527f786d74702e7261746552656769737472792e6d69677261746f72000000000000602082015290565b610c55565b610c68565b565b5f5f610a538484610e73565b905067ffffffffffffffff811115610a97576040517f37f4f14800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b90505b92915050565b7f988e236e2caf5758fdf811320ba1d2fca453cb71bd6049ebba876b68af50500080545f03610acf5750610c27565b80545f908290610ae190600190611115565b81548110610af157610af16110e8565b5f9182526020918290206040805160a0810182526002909302909101805467ffffffffffffffff8082168086526801000000000000000083048216968601969096527001000000000000000000000000000000008204811693850193909352780100000000000000000000000000000000000000000000000090048216606084015260010154811660808301529092508716148015610ba757508467ffffffffffffffff16816020015167ffffffffffffffff16145b8015610bca57508367ffffffffffffffff16816040015167ffffffffffffffff16145b8015610bed57508267ffffffffffffffff16816060015167ffffffffffffffff16145b15610c24576040517fa88ee57700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50505b50505050565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00610a9a565b5f610a97610c638484610e73565b610f06565b73ffffffffffffffffffffffffffffffffffffffff8116610cb5576040517f0d626a3200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405173ffffffffffffffffffffffffffffffffffffffff8216907fa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098905f90a25f5f8273ffffffffffffffffffffffffffffffffffffffff166040515f60405180830381855af49150503d805f8114610d49576040519150601f19603f3d011682016040523d82523d5f602084013e610d4e565b606091505b509150915081610d975782816040517f68b0b16b000000000000000000000000000000000000000000000000000000008152600401610d8e929190611128565b60405180910390fd5b8051158015610dbb575073ffffffffffffffffffffffffffffffffffffffff83163b155b15610e0a576040517f626c416100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff84166004820152602401610d8e565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5473ffffffffffffffffffffffffffffffffffffffff167fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b60405160405180910390a2505050565b6040517fd6d7d5250000000000000000000000000000000000000000000000000000000081525f9073ffffffffffffffffffffffffffffffffffffffff84169063d6d7d52590610ec7908590600401611069565b602060405180830381865afa158015610ee2573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610a97919061115e565b5f73ffffffffffffffffffffffffffffffffffffffff821115610f55576040517f37f4f14800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5090565b5f5f60408385031215610f6a575f5ffd5b50508035926020909101359150565b602080825282518282018190525f918401906040840190835b8181101561101257835167ffffffffffffffff815116845267ffffffffffffffff602082015116602085015267ffffffffffffffff604082015116604085015267ffffffffffffffff606082015116606085015267ffffffffffffffff60808201511660808501525060a083019250602084019350600181019050610f92565b509095945050505050565b5f81518084528060208401602086015e5f6020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081525f610a97602083018461101d565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b80820180821115610a9a57610a9a61107b565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b81810381811115610a9a57610a9a61107b565b73ffffffffffffffffffffffffffffffffffffffff83168152604060208201525f611156604083018461101d565b949350505050565b5f6020828403121561116e575f5ffd5b505191905056fe786d74702e7261746552656769737472792e746172676574526174655065724d696e757465a26469706673582212209eeedf534b84f2f18b5935ea4b24866b91346fbb3964762855a2d36d30c8532664736f6c634300081c0033",
}

// RateRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use RateRegistryMetaData.ABI instead.
var RateRegistryABI = RateRegistryMetaData.ABI

// RateRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RateRegistryMetaData.Bin instead.
var RateRegistryBin = RateRegistryMetaData.Bin

// DeployRateRegistry deploys a new Ethereum contract, binding an instance of RateRegistry to it.
func DeployRateRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, parameterRegistry_ common.Address) (common.Address, *types.Transaction, *RateRegistry, error) {
	parsed, err := RateRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RateRegistryBin), backend, parameterRegistry_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RateRegistry{RateRegistryCaller: RateRegistryCaller{contract: contract}, RateRegistryTransactor: RateRegistryTransactor{contract: contract}, RateRegistryFilterer: RateRegistryFilterer{contract: contract}}, nil
}

// RateRegistry is an auto generated Go binding around an Ethereum contract.
type RateRegistry struct {
	RateRegistryCaller     // Read-only binding to the contract
	RateRegistryTransactor // Write-only binding to the contract
	RateRegistryFilterer   // Log filterer for contract events
}

// RateRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type RateRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RateRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RateRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RateRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RateRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RateRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RateRegistrySession struct {
	Contract     *RateRegistry     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RateRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RateRegistryCallerSession struct {
	Contract *RateRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// RateRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RateRegistryTransactorSession struct {
	Contract     *RateRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// RateRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type RateRegistryRaw struct {
	Contract *RateRegistry // Generic contract binding to access the raw methods on
}

// RateRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RateRegistryCallerRaw struct {
	Contract *RateRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// RateRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RateRegistryTransactorRaw struct {
	Contract *RateRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRateRegistry creates a new instance of RateRegistry, bound to a specific deployed contract.
func NewRateRegistry(address common.Address, backend bind.ContractBackend) (*RateRegistry, error) {
	contract, err := bindRateRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RateRegistry{RateRegistryCaller: RateRegistryCaller{contract: contract}, RateRegistryTransactor: RateRegistryTransactor{contract: contract}, RateRegistryFilterer: RateRegistryFilterer{contract: contract}}, nil
}

// NewRateRegistryCaller creates a new read-only instance of RateRegistry, bound to a specific deployed contract.
func NewRateRegistryCaller(address common.Address, caller bind.ContractCaller) (*RateRegistryCaller, error) {
	contract, err := bindRateRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RateRegistryCaller{contract: contract}, nil
}

// NewRateRegistryTransactor creates a new write-only instance of RateRegistry, bound to a specific deployed contract.
func NewRateRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*RateRegistryTransactor, error) {
	contract, err := bindRateRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RateRegistryTransactor{contract: contract}, nil
}

// NewRateRegistryFilterer creates a new log filterer instance of RateRegistry, bound to a specific deployed contract.
func NewRateRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*RateRegistryFilterer, error) {
	contract, err := bindRateRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RateRegistryFilterer{contract: contract}, nil
}

// bindRateRegistry binds a generic wrapper to an already deployed contract.
func bindRateRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RateRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RateRegistry *RateRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RateRegistry.Contract.RateRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RateRegistry *RateRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RateRegistry.Contract.RateRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RateRegistry *RateRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RateRegistry.Contract.RateRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RateRegistry *RateRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RateRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RateRegistry *RateRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RateRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RateRegistry *RateRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RateRegistry.Contract.contract.Transact(opts, method, params...)
}

// CongestionFeeParameterKey is a free data retrieval call binding the contract method 0xed7e6986.
//
// Solidity: function congestionFeeParameterKey() pure returns(bytes key_)
func (_RateRegistry *RateRegistryCaller) CongestionFeeParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "congestionFeeParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// CongestionFeeParameterKey is a free data retrieval call binding the contract method 0xed7e6986.
//
// Solidity: function congestionFeeParameterKey() pure returns(bytes key_)
func (_RateRegistry *RateRegistrySession) CongestionFeeParameterKey() ([]byte, error) {
	return _RateRegistry.Contract.CongestionFeeParameterKey(&_RateRegistry.CallOpts)
}

// CongestionFeeParameterKey is a free data retrieval call binding the contract method 0xed7e6986.
//
// Solidity: function congestionFeeParameterKey() pure returns(bytes key_)
func (_RateRegistry *RateRegistryCallerSession) CongestionFeeParameterKey() ([]byte, error) {
	return _RateRegistry.Contract.CongestionFeeParameterKey(&_RateRegistry.CallOpts)
}

// GetRates is a free data retrieval call binding the contract method 0x45b05a43.
//
// Solidity: function getRates(uint256 fromIndex_, uint256 count_) view returns((uint64,uint64,uint64,uint64,uint64)[] rates_)
func (_RateRegistry *RateRegistryCaller) GetRates(opts *bind.CallOpts, fromIndex_ *big.Int, count_ *big.Int) ([]IRateRegistryRates, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "getRates", fromIndex_, count_)

	if err != nil {
		return *new([]IRateRegistryRates), err
	}

	out0 := *abi.ConvertType(out[0], new([]IRateRegistryRates)).(*[]IRateRegistryRates)

	return out0, err

}

// GetRates is a free data retrieval call binding the contract method 0x45b05a43.
//
// Solidity: function getRates(uint256 fromIndex_, uint256 count_) view returns((uint64,uint64,uint64,uint64,uint64)[] rates_)
func (_RateRegistry *RateRegistrySession) GetRates(fromIndex_ *big.Int, count_ *big.Int) ([]IRateRegistryRates, error) {
	return _RateRegistry.Contract.GetRates(&_RateRegistry.CallOpts, fromIndex_, count_)
}

// GetRates is a free data retrieval call binding the contract method 0x45b05a43.
//
// Solidity: function getRates(uint256 fromIndex_, uint256 count_) view returns((uint64,uint64,uint64,uint64,uint64)[] rates_)
func (_RateRegistry *RateRegistryCallerSession) GetRates(fromIndex_ *big.Int, count_ *big.Int) ([]IRateRegistryRates, error) {
	return _RateRegistry.Contract.GetRates(&_RateRegistry.CallOpts, fromIndex_, count_)
}

// GetRatesCount is a free data retrieval call binding the contract method 0x2da72291.
//
// Solidity: function getRatesCount() view returns(uint256 count_)
func (_RateRegistry *RateRegistryCaller) GetRatesCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "getRatesCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRatesCount is a free data retrieval call binding the contract method 0x2da72291.
//
// Solidity: function getRatesCount() view returns(uint256 count_)
func (_RateRegistry *RateRegistrySession) GetRatesCount() (*big.Int, error) {
	return _RateRegistry.Contract.GetRatesCount(&_RateRegistry.CallOpts)
}

// GetRatesCount is a free data retrieval call binding the contract method 0x2da72291.
//
// Solidity: function getRatesCount() view returns(uint256 count_)
func (_RateRegistry *RateRegistryCallerSession) GetRatesCount() (*big.Int, error) {
	return _RateRegistry.Contract.GetRatesCount(&_RateRegistry.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_RateRegistry *RateRegistryCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_RateRegistry *RateRegistrySession) Implementation() (common.Address, error) {
	return _RateRegistry.Contract.Implementation(&_RateRegistry.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_RateRegistry *RateRegistryCallerSession) Implementation() (common.Address, error) {
	return _RateRegistry.Contract.Implementation(&_RateRegistry.CallOpts)
}

// MessageFeeParameterKey is a free data retrieval call binding the contract method 0x63c03291.
//
// Solidity: function messageFeeParameterKey() pure returns(bytes key_)
func (_RateRegistry *RateRegistryCaller) MessageFeeParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "messageFeeParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// MessageFeeParameterKey is a free data retrieval call binding the contract method 0x63c03291.
//
// Solidity: function messageFeeParameterKey() pure returns(bytes key_)
func (_RateRegistry *RateRegistrySession) MessageFeeParameterKey() ([]byte, error) {
	return _RateRegistry.Contract.MessageFeeParameterKey(&_RateRegistry.CallOpts)
}

// MessageFeeParameterKey is a free data retrieval call binding the contract method 0x63c03291.
//
// Solidity: function messageFeeParameterKey() pure returns(bytes key_)
func (_RateRegistry *RateRegistryCallerSession) MessageFeeParameterKey() ([]byte, error) {
	return _RateRegistry.Contract.MessageFeeParameterKey(&_RateRegistry.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(bytes key_)
func (_RateRegistry *RateRegistryCaller) MigratorParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "migratorParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(bytes key_)
func (_RateRegistry *RateRegistrySession) MigratorParameterKey() ([]byte, error) {
	return _RateRegistry.Contract.MigratorParameterKey(&_RateRegistry.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(bytes key_)
func (_RateRegistry *RateRegistryCallerSession) MigratorParameterKey() ([]byte, error) {
	return _RateRegistry.Contract.MigratorParameterKey(&_RateRegistry.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_RateRegistry *RateRegistryCaller) ParameterRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "parameterRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_RateRegistry *RateRegistrySession) ParameterRegistry() (common.Address, error) {
	return _RateRegistry.Contract.ParameterRegistry(&_RateRegistry.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_RateRegistry *RateRegistryCallerSession) ParameterRegistry() (common.Address, error) {
	return _RateRegistry.Contract.ParameterRegistry(&_RateRegistry.CallOpts)
}

// StorageFeeParameterKey is a free data retrieval call binding the contract method 0xba3261d5.
//
// Solidity: function storageFeeParameterKey() pure returns(bytes key_)
func (_RateRegistry *RateRegistryCaller) StorageFeeParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "storageFeeParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// StorageFeeParameterKey is a free data retrieval call binding the contract method 0xba3261d5.
//
// Solidity: function storageFeeParameterKey() pure returns(bytes key_)
func (_RateRegistry *RateRegistrySession) StorageFeeParameterKey() ([]byte, error) {
	return _RateRegistry.Contract.StorageFeeParameterKey(&_RateRegistry.CallOpts)
}

// StorageFeeParameterKey is a free data retrieval call binding the contract method 0xba3261d5.
//
// Solidity: function storageFeeParameterKey() pure returns(bytes key_)
func (_RateRegistry *RateRegistryCallerSession) StorageFeeParameterKey() ([]byte, error) {
	return _RateRegistry.Contract.StorageFeeParameterKey(&_RateRegistry.CallOpts)
}

// TargetRatePerMinuteParameterKey is a free data retrieval call binding the contract method 0x49156e09.
//
// Solidity: function targetRatePerMinuteParameterKey() pure returns(bytes key_)
func (_RateRegistry *RateRegistryCaller) TargetRatePerMinuteParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "targetRatePerMinuteParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// TargetRatePerMinuteParameterKey is a free data retrieval call binding the contract method 0x49156e09.
//
// Solidity: function targetRatePerMinuteParameterKey() pure returns(bytes key_)
func (_RateRegistry *RateRegistrySession) TargetRatePerMinuteParameterKey() ([]byte, error) {
	return _RateRegistry.Contract.TargetRatePerMinuteParameterKey(&_RateRegistry.CallOpts)
}

// TargetRatePerMinuteParameterKey is a free data retrieval call binding the contract method 0x49156e09.
//
// Solidity: function targetRatePerMinuteParameterKey() pure returns(bytes key_)
func (_RateRegistry *RateRegistryCallerSession) TargetRatePerMinuteParameterKey() ([]byte, error) {
	return _RateRegistry.Contract.TargetRatePerMinuteParameterKey(&_RateRegistry.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_RateRegistry *RateRegistryTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RateRegistry.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_RateRegistry *RateRegistrySession) Initialize() (*types.Transaction, error) {
	return _RateRegistry.Contract.Initialize(&_RateRegistry.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_RateRegistry *RateRegistryTransactorSession) Initialize() (*types.Transaction, error) {
	return _RateRegistry.Contract.Initialize(&_RateRegistry.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_RateRegistry *RateRegistryTransactor) Migrate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RateRegistry.contract.Transact(opts, "migrate")
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_RateRegistry *RateRegistrySession) Migrate() (*types.Transaction, error) {
	return _RateRegistry.Contract.Migrate(&_RateRegistry.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_RateRegistry *RateRegistryTransactorSession) Migrate() (*types.Transaction, error) {
	return _RateRegistry.Contract.Migrate(&_RateRegistry.TransactOpts)
}

// UpdateRates is a paid mutator transaction binding the contract method 0x3c3821f4.
//
// Solidity: function updateRates() returns()
func (_RateRegistry *RateRegistryTransactor) UpdateRates(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RateRegistry.contract.Transact(opts, "updateRates")
}

// UpdateRates is a paid mutator transaction binding the contract method 0x3c3821f4.
//
// Solidity: function updateRates() returns()
func (_RateRegistry *RateRegistrySession) UpdateRates() (*types.Transaction, error) {
	return _RateRegistry.Contract.UpdateRates(&_RateRegistry.TransactOpts)
}

// UpdateRates is a paid mutator transaction binding the contract method 0x3c3821f4.
//
// Solidity: function updateRates() returns()
func (_RateRegistry *RateRegistryTransactorSession) UpdateRates() (*types.Transaction, error) {
	return _RateRegistry.Contract.UpdateRates(&_RateRegistry.TransactOpts)
}

// RateRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the RateRegistry contract.
type RateRegistryInitializedIterator struct {
	Event *RateRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *RateRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RateRegistryInitialized)
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
		it.Event = new(RateRegistryInitialized)
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
func (it *RateRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RateRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RateRegistryInitialized represents a Initialized event raised by the RateRegistry contract.
type RateRegistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_RateRegistry *RateRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*RateRegistryInitializedIterator, error) {

	logs, sub, err := _RateRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &RateRegistryInitializedIterator{contract: _RateRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_RateRegistry *RateRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *RateRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _RateRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RateRegistryInitialized)
				if err := _RateRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_RateRegistry *RateRegistryFilterer) ParseInitialized(log types.Log) (*RateRegistryInitialized, error) {
	event := new(RateRegistryInitialized)
	if err := _RateRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RateRegistryMigratedIterator is returned from FilterMigrated and is used to iterate over the raw logs and unpacked data for Migrated events raised by the RateRegistry contract.
type RateRegistryMigratedIterator struct {
	Event *RateRegistryMigrated // Event containing the contract specifics and raw log

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
func (it *RateRegistryMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RateRegistryMigrated)
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
		it.Event = new(RateRegistryMigrated)
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
func (it *RateRegistryMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RateRegistryMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RateRegistryMigrated represents a Migrated event raised by the RateRegistry contract.
type RateRegistryMigrated struct {
	Migrator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMigrated is a free log retrieval operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_RateRegistry *RateRegistryFilterer) FilterMigrated(opts *bind.FilterOpts, migrator []common.Address) (*RateRegistryMigratedIterator, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _RateRegistry.contract.FilterLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return &RateRegistryMigratedIterator{contract: _RateRegistry.contract, event: "Migrated", logs: logs, sub: sub}, nil
}

// WatchMigrated is a free log subscription operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_RateRegistry *RateRegistryFilterer) WatchMigrated(opts *bind.WatchOpts, sink chan<- *RateRegistryMigrated, migrator []common.Address) (event.Subscription, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _RateRegistry.contract.WatchLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RateRegistryMigrated)
				if err := _RateRegistry.contract.UnpackLog(event, "Migrated", log); err != nil {
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
func (_RateRegistry *RateRegistryFilterer) ParseMigrated(log types.Log) (*RateRegistryMigrated, error) {
	event := new(RateRegistryMigrated)
	if err := _RateRegistry.contract.UnpackLog(event, "Migrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RateRegistryRatesUpdatedIterator is returned from FilterRatesUpdated and is used to iterate over the raw logs and unpacked data for RatesUpdated events raised by the RateRegistry contract.
type RateRegistryRatesUpdatedIterator struct {
	Event *RateRegistryRatesUpdated // Event containing the contract specifics and raw log

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
func (it *RateRegistryRatesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RateRegistryRatesUpdated)
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
		it.Event = new(RateRegistryRatesUpdated)
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
func (it *RateRegistryRatesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RateRegistryRatesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RateRegistryRatesUpdated represents a RatesUpdated event raised by the RateRegistry contract.
type RateRegistryRatesUpdated struct {
	MessageFee          uint64
	StorageFee          uint64
	CongestionFee       uint64
	TargetRatePerMinute uint64
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterRatesUpdated is a free log retrieval operation binding the contract event 0xabd2140443b16c364d95086ebf21b45137b5f0af53e10a6e792c0cb3d0e2db62.
//
// Solidity: event RatesUpdated(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 targetRatePerMinute)
func (_RateRegistry *RateRegistryFilterer) FilterRatesUpdated(opts *bind.FilterOpts) (*RateRegistryRatesUpdatedIterator, error) {

	logs, sub, err := _RateRegistry.contract.FilterLogs(opts, "RatesUpdated")
	if err != nil {
		return nil, err
	}
	return &RateRegistryRatesUpdatedIterator{contract: _RateRegistry.contract, event: "RatesUpdated", logs: logs, sub: sub}, nil
}

// WatchRatesUpdated is a free log subscription operation binding the contract event 0xabd2140443b16c364d95086ebf21b45137b5f0af53e10a6e792c0cb3d0e2db62.
//
// Solidity: event RatesUpdated(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 targetRatePerMinute)
func (_RateRegistry *RateRegistryFilterer) WatchRatesUpdated(opts *bind.WatchOpts, sink chan<- *RateRegistryRatesUpdated) (event.Subscription, error) {

	logs, sub, err := _RateRegistry.contract.WatchLogs(opts, "RatesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RateRegistryRatesUpdated)
				if err := _RateRegistry.contract.UnpackLog(event, "RatesUpdated", log); err != nil {
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

// ParseRatesUpdated is a log parse operation binding the contract event 0xabd2140443b16c364d95086ebf21b45137b5f0af53e10a6e792c0cb3d0e2db62.
//
// Solidity: event RatesUpdated(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 targetRatePerMinute)
func (_RateRegistry *RateRegistryFilterer) ParseRatesUpdated(log types.Log) (*RateRegistryRatesUpdated, error) {
	event := new(RateRegistryRatesUpdated)
	if err := _RateRegistry.contract.UnpackLog(event, "RatesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RateRegistryUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the RateRegistry contract.
type RateRegistryUpgradedIterator struct {
	Event *RateRegistryUpgraded // Event containing the contract specifics and raw log

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
func (it *RateRegistryUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RateRegistryUpgraded)
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
		it.Event = new(RateRegistryUpgraded)
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
func (it *RateRegistryUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RateRegistryUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RateRegistryUpgraded represents a Upgraded event raised by the RateRegistry contract.
type RateRegistryUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_RateRegistry *RateRegistryFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*RateRegistryUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _RateRegistry.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &RateRegistryUpgradedIterator{contract: _RateRegistry.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_RateRegistry *RateRegistryFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *RateRegistryUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _RateRegistry.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RateRegistryUpgraded)
				if err := _RateRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_RateRegistry *RateRegistryFilterer) ParseUpgraded(log types.Log) (*RateRegistryUpgraded, error) {
	event := new(RateRegistryUpgraded)
	if err := _RateRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
