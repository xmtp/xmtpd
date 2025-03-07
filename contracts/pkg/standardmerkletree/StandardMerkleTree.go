// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package standardmerkletree

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

// StandardMerkleTreeMetaData contains all meta data concerning the StandardMerkleTree contract.
var StandardMerkleTreeMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"multiProofVerify\",\"inputs\":[{\"name\":\"proof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"proofFlags\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"},{\"name\":\"offset\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"amounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"indices\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"error\",\"name\":\"EmptyProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidIndex\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MerkleProofInvalidMultiproof\",\"inputs\":[]}]",
	Bin: "0x6080604052348015600e575f5ffd5b506108128061001c5f395ff3fe608060405234801561000f575f5ffd5b5060043610610029575f3560e01c80635cc9c8fd1461002d575b5f5ffd5b61004061003b36600461059b565b610054565b604051901515815260200160405180910390f35b5f808667ffffffffffffffff81111561006f5761006f6106b2565b604051908082528060200260200182016040528015610098578160200160208202803683370190505b5090505f5b87811015610222575f8585838181106100b8576100b86106df565b602002919091013591506100ce9050898d610739565b8111806100da57508b81105b15610111576040517f63df817100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6101fc868684818110610126576101266106df565b905060200201358b8b8581811061013f5761013f6106df565b9050602002016020810190610154919061074c565b8a8a86818110610166576101666106df565b90506020020135604080516020810185905273ffffffffffffffffffffffffffffffffffffffff841691810191909152606081018290525f90608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152828252805160209182012090830152016040516020818303038152906040528051906020012090509392505050565b83838151811061020e5761020e6106df565b60209081029190910101525060010161009d565b5080515f0361025d576040517f668fd6f300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61026b8e8e8e8e8d8661027d565b9e9d5050505050505050505050505050565b5f8261028c8888888887610298565b14979650505050505050565b80515f90836102a8816001610739565b6102b28884610739565b146102e9576040517f3514049200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8167ffffffffffffffff811115610303576103036106b2565b60405190808252806020026020018201604052801561032c578160200160208202803683370190505b5090505f8080805b85811015610475575f87851061036e57858461034f81610786565b955081518110610361576103616106df565b6020026020010151610394565b898561037981610786565b96508151811061038b5761038b6106df565b60200260200101515b90505f8c8c848181106103a9576103a96106df565b90506020020160208101906103be91906107bd565b6103eb578e8e856103ce81610786565b96508181106103df576103df6106df565b90506020020135610442565b88861061041c5786856103fd81610786565b96508151811061040f5761040f6106df565b6020026020010151610442565b8a8661042781610786565b975081518110610439576104396106df565b60200260200101515b905061044e8282610522565b878481518110610460576104606106df565b60209081029190910101525050600101610334565b5084156104df57808b146104b5576040517f3514049200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8360018603815181106104ca576104ca6106df565b60200260200101519650505050505050610519565b85156104f757875f815181106104ca576104ca6106df565b8b8b5f818110610509576105096106df565b9050602002013596505050505050505b95945050505050565b5f81831061053c575f82815260208490526040902061054a565b5f8381526020839052604090205b90505b92915050565b5f5f83601f840112610563575f5ffd5b50813567ffffffffffffffff81111561057a575f5ffd5b6020830191508360208260051b8501011115610594575f5ffd5b9250929050565b5f5f5f5f5f5f5f5f5f5f5f5f60e08d8f0312156105b6575f5ffd5b67ffffffffffffffff8d3511156105cb575f5ffd5b6105d88e8e358f01610553565b909c509a5067ffffffffffffffff60208e013511156105f5575f5ffd5b6106058e60208f01358f01610553565b909a50985060408d0135975060608d0135965067ffffffffffffffff60808e01351115610630575f5ffd5b6106408e60808f01358f01610553565b909650945067ffffffffffffffff60a08e0135111561065d575f5ffd5b61066d8e60a08f01358f01610553565b909450925067ffffffffffffffff60c08e0135111561068a575f5ffd5b61069a8e60c08f01358f01610553565b81935080925050509295989b509295989b509295989b565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b8082018082111561054d5761054d61070c565b5f6020828403121561075c575f5ffd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461077f575f5ffd5b9392505050565b5f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036107b6576107b661070c565b5060010190565b5f602082840312156107cd575f5ffd5b8135801515811461077f575f5ffdfea2646970667358221220ef95423608055701aef4602387be244858a548c15d64c256938344678a8dac7264736f6c634300081c0033",
}

// StandardMerkleTreeABI is the input ABI used to generate the binding from.
// Deprecated: Use StandardMerkleTreeMetaData.ABI instead.
var StandardMerkleTreeABI = StandardMerkleTreeMetaData.ABI

// StandardMerkleTreeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StandardMerkleTreeMetaData.Bin instead.
var StandardMerkleTreeBin = StandardMerkleTreeMetaData.Bin

// DeployStandardMerkleTree deploys a new Ethereum contract, binding an instance of StandardMerkleTree to it.
func DeployStandardMerkleTree(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StandardMerkleTree, error) {
	parsed, err := StandardMerkleTreeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StandardMerkleTreeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StandardMerkleTree{StandardMerkleTreeCaller: StandardMerkleTreeCaller{contract: contract}, StandardMerkleTreeTransactor: StandardMerkleTreeTransactor{contract: contract}, StandardMerkleTreeFilterer: StandardMerkleTreeFilterer{contract: contract}}, nil
}

// StandardMerkleTree is an auto generated Go binding around an Ethereum contract.
type StandardMerkleTree struct {
	StandardMerkleTreeCaller     // Read-only binding to the contract
	StandardMerkleTreeTransactor // Write-only binding to the contract
	StandardMerkleTreeFilterer   // Log filterer for contract events
}

// StandardMerkleTreeCaller is an auto generated read-only Go binding around an Ethereum contract.
type StandardMerkleTreeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StandardMerkleTreeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StandardMerkleTreeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StandardMerkleTreeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StandardMerkleTreeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StandardMerkleTreeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StandardMerkleTreeSession struct {
	Contract     *StandardMerkleTree // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// StandardMerkleTreeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StandardMerkleTreeCallerSession struct {
	Contract *StandardMerkleTreeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// StandardMerkleTreeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StandardMerkleTreeTransactorSession struct {
	Contract     *StandardMerkleTreeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// StandardMerkleTreeRaw is an auto generated low-level Go binding around an Ethereum contract.
type StandardMerkleTreeRaw struct {
	Contract *StandardMerkleTree // Generic contract binding to access the raw methods on
}

// StandardMerkleTreeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StandardMerkleTreeCallerRaw struct {
	Contract *StandardMerkleTreeCaller // Generic read-only contract binding to access the raw methods on
}

// StandardMerkleTreeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StandardMerkleTreeTransactorRaw struct {
	Contract *StandardMerkleTreeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStandardMerkleTree creates a new instance of StandardMerkleTree, bound to a specific deployed contract.
func NewStandardMerkleTree(address common.Address, backend bind.ContractBackend) (*StandardMerkleTree, error) {
	contract, err := bindStandardMerkleTree(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StandardMerkleTree{StandardMerkleTreeCaller: StandardMerkleTreeCaller{contract: contract}, StandardMerkleTreeTransactor: StandardMerkleTreeTransactor{contract: contract}, StandardMerkleTreeFilterer: StandardMerkleTreeFilterer{contract: contract}}, nil
}

// NewStandardMerkleTreeCaller creates a new read-only instance of StandardMerkleTree, bound to a specific deployed contract.
func NewStandardMerkleTreeCaller(address common.Address, caller bind.ContractCaller) (*StandardMerkleTreeCaller, error) {
	contract, err := bindStandardMerkleTree(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StandardMerkleTreeCaller{contract: contract}, nil
}

// NewStandardMerkleTreeTransactor creates a new write-only instance of StandardMerkleTree, bound to a specific deployed contract.
func NewStandardMerkleTreeTransactor(address common.Address, transactor bind.ContractTransactor) (*StandardMerkleTreeTransactor, error) {
	contract, err := bindStandardMerkleTree(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StandardMerkleTreeTransactor{contract: contract}, nil
}

// NewStandardMerkleTreeFilterer creates a new log filterer instance of StandardMerkleTree, bound to a specific deployed contract.
func NewStandardMerkleTreeFilterer(address common.Address, filterer bind.ContractFilterer) (*StandardMerkleTreeFilterer, error) {
	contract, err := bindStandardMerkleTree(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StandardMerkleTreeFilterer{contract: contract}, nil
}

// bindStandardMerkleTree binds a generic wrapper to an already deployed contract.
func bindStandardMerkleTree(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StandardMerkleTreeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StandardMerkleTree *StandardMerkleTreeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StandardMerkleTree.Contract.StandardMerkleTreeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StandardMerkleTree *StandardMerkleTreeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StandardMerkleTree.Contract.StandardMerkleTreeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StandardMerkleTree *StandardMerkleTreeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StandardMerkleTree.Contract.StandardMerkleTreeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StandardMerkleTree *StandardMerkleTreeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StandardMerkleTree.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StandardMerkleTree *StandardMerkleTreeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StandardMerkleTree.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StandardMerkleTree *StandardMerkleTreeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StandardMerkleTree.Contract.contract.Transact(opts, method, params...)
}

// MultiProofVerify is a free data retrieval call binding the contract method 0x5cc9c8fd.
//
// Solidity: function multiProofVerify(bytes32[] proof, bool[] proofFlags, uint256 offset, bytes32 root, address[] accounts, uint256[] amounts, uint256[] indices) pure returns(bool)
func (_StandardMerkleTree *StandardMerkleTreeCaller) MultiProofVerify(opts *bind.CallOpts, proof [][32]byte, proofFlags []bool, offset *big.Int, root [32]byte, accounts []common.Address, amounts []*big.Int, indices []*big.Int) (bool, error) {
	var out []interface{}
	err := _StandardMerkleTree.contract.Call(opts, &out, "multiProofVerify", proof, proofFlags, offset, root, accounts, amounts, indices)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// MultiProofVerify is a free data retrieval call binding the contract method 0x5cc9c8fd.
//
// Solidity: function multiProofVerify(bytes32[] proof, bool[] proofFlags, uint256 offset, bytes32 root, address[] accounts, uint256[] amounts, uint256[] indices) pure returns(bool)
func (_StandardMerkleTree *StandardMerkleTreeSession) MultiProofVerify(proof [][32]byte, proofFlags []bool, offset *big.Int, root [32]byte, accounts []common.Address, amounts []*big.Int, indices []*big.Int) (bool, error) {
	return _StandardMerkleTree.Contract.MultiProofVerify(&_StandardMerkleTree.CallOpts, proof, proofFlags, offset, root, accounts, amounts, indices)
}

// MultiProofVerify is a free data retrieval call binding the contract method 0x5cc9c8fd.
//
// Solidity: function multiProofVerify(bytes32[] proof, bool[] proofFlags, uint256 offset, bytes32 root, address[] accounts, uint256[] amounts, uint256[] indices) pure returns(bool)
func (_StandardMerkleTree *StandardMerkleTreeCallerSession) MultiProofVerify(proof [][32]byte, proofFlags []bool, offset *big.Int, root [32]byte, accounts []common.Address, amounts []*big.Int, indices []*big.Int) (bool, error) {
	return _StandardMerkleTree.Contract.MultiProofVerify(&_StandardMerkleTree.CallOpts, proof, proofFlags, offset, root, accounts, amounts, indices)
}
