// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package noderegistry

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

// INodeRegistryNode is an auto generated low-level Go binding around an user-defined struct.
type INodeRegistryNode struct {
	Signer           common.Address
	IsCanonical      bool
	SigningPublicKey []byte
	HttpAddress      string
}

// INodeRegistryNodeWithId is an auto generated low-level Go binding around an user-defined struct.
type INodeRegistryNodeWithId struct {
	NodeId uint32
	Node   INodeRegistryNode
}

// NodeRegistryMetaData contains all meta data concerning the NodeRegistry contract.
var NodeRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"parameterRegistry_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"NODE_INCREMENT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addNode\",\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"signingPublicKey_\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress_\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"nodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"signer_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addToNetwork\",\"inputs\":[{\"name\":\"nodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"admin\",\"inputs\":[],\"outputs\":[{\"name\":\"admin_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"adminParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"canonicalNodesCount\",\"inputs\":[],\"outputs\":[{\"name\":\"canonicalNodesCount_\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"allNodes_\",\"type\":\"tuple[]\",\"internalType\":\"structINodeRegistry.NodeWithId[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"node\",\"type\":\"tuple\",\"internalType\":\"structINodeRegistry.Node\",\"components\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isCanonical\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signingPublicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllNodesCount\",\"inputs\":[],\"outputs\":[{\"name\":\"nodeCount_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApproved\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsCanonicalNode\",\"inputs\":[{\"name\":\"nodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"isCanonicalNode_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNode\",\"inputs\":[{\"name\":\"nodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"node_\",\"type\":\"tuple\",\"internalType\":\"structINodeRegistry.Node\",\"components\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isCanonical\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signingPublicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSigner\",\"inputs\":[{\"name\":\"nodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"signer_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"implementation_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxCanonicalNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"maxCanonicalNodes_\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxCanonicalNodesParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"migrate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migratorParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerOf\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"parameterRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeFromNetwork\",\"inputs\":[{\"name\":\"nodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBaseURI\",\"inputs\":[{\"name\":\"baseURI_\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setHttpAddress\",\"inputs\":[{\"name\":\"nodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"httpAddress_\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateAdmin\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMaxCanonicalNodes\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AdminUpdated\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BaseURIUpdated\",\"inputs\":[{\"name\":\"baseURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"HttpAddressUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MaxCanonicalNodesUpdated\",\"inputs\":[{\"name\":\"maxCanonicalNodes\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Migrated\",\"inputs\":[{\"name\":\"migrator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeAdded\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"signer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"signingPublicKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeAddedToCanonicalNetwork\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeRemovedFromCanonicalNetwork\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC721IncorrectOwner\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InsufficientApproval\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721NonexistentToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"EmptyCode\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"FailedToAddNodeToCanonicalNetwork\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToRemoveNodeFromCanonicalNetwork\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidHttpAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSigningPublicKey\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidURI\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxCanonicalNodesBelowCurrentCount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxCanonicalNodesReached\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxNodesReached\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationFailed\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertData_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"NoChange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotNodeOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterOutOfTypeBounds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroMigrator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroParameterRegistry\",\"inputs\":[]}]",
}

// NodeRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use NodeRegistryMetaData.ABI instead.
var NodeRegistryABI = NodeRegistryMetaData.ABI

// NodeRegistry is an auto generated Go binding around an Ethereum contract.
type NodeRegistry struct {
	NodeRegistryCaller     // Read-only binding to the contract
	NodeRegistryTransactor // Write-only binding to the contract
	NodeRegistryFilterer   // Log filterer for contract events
}

// NodeRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type NodeRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NodeRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NodeRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NodeRegistrySession struct {
	Contract     *NodeRegistry     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NodeRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NodeRegistryCallerSession struct {
	Contract *NodeRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// NodeRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NodeRegistryTransactorSession struct {
	Contract     *NodeRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// NodeRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type NodeRegistryRaw struct {
	Contract *NodeRegistry // Generic contract binding to access the raw methods on
}

// NodeRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NodeRegistryCallerRaw struct {
	Contract *NodeRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// NodeRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NodeRegistryTransactorRaw struct {
	Contract *NodeRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNodeRegistry creates a new instance of NodeRegistry, bound to a specific deployed contract.
func NewNodeRegistry(address common.Address, backend bind.ContractBackend) (*NodeRegistry, error) {
	contract, err := bindNodeRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NodeRegistry{NodeRegistryCaller: NodeRegistryCaller{contract: contract}, NodeRegistryTransactor: NodeRegistryTransactor{contract: contract}, NodeRegistryFilterer: NodeRegistryFilterer{contract: contract}}, nil
}

// NewNodeRegistryCaller creates a new read-only instance of NodeRegistry, bound to a specific deployed contract.
func NewNodeRegistryCaller(address common.Address, caller bind.ContractCaller) (*NodeRegistryCaller, error) {
	contract, err := bindNodeRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryCaller{contract: contract}, nil
}

// NewNodeRegistryTransactor creates a new write-only instance of NodeRegistry, bound to a specific deployed contract.
func NewNodeRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*NodeRegistryTransactor, error) {
	contract, err := bindNodeRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryTransactor{contract: contract}, nil
}

// NewNodeRegistryFilterer creates a new log filterer instance of NodeRegistry, bound to a specific deployed contract.
func NewNodeRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*NodeRegistryFilterer, error) {
	contract, err := bindNodeRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryFilterer{contract: contract}, nil
}

// bindNodeRegistry binds a generic wrapper to an already deployed contract.
func bindNodeRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NodeRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeRegistry *NodeRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeRegistry.Contract.NodeRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeRegistry *NodeRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeRegistry.Contract.NodeRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeRegistry *NodeRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeRegistry.Contract.NodeRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeRegistry *NodeRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeRegistry *NodeRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeRegistry *NodeRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeRegistry.Contract.contract.Transact(opts, method, params...)
}

// NODEINCREMENT is a free data retrieval call binding the contract method 0xfd667d1e.
//
// Solidity: function NODE_INCREMENT() view returns(uint32)
func (_NodeRegistry *NodeRegistryCaller) NODEINCREMENT(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "NODE_INCREMENT")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// NODEINCREMENT is a free data retrieval call binding the contract method 0xfd667d1e.
//
// Solidity: function NODE_INCREMENT() view returns(uint32)
func (_NodeRegistry *NodeRegistrySession) NODEINCREMENT() (uint32, error) {
	return _NodeRegistry.Contract.NODEINCREMENT(&_NodeRegistry.CallOpts)
}

// NODEINCREMENT is a free data retrieval call binding the contract method 0xfd667d1e.
//
// Solidity: function NODE_INCREMENT() view returns(uint32)
func (_NodeRegistry *NodeRegistryCallerSession) NODEINCREMENT() (uint32, error) {
	return _NodeRegistry.Contract.NODEINCREMENT(&_NodeRegistry.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address admin_)
func (_NodeRegistry *NodeRegistryCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address admin_)
func (_NodeRegistry *NodeRegistrySession) Admin() (common.Address, error) {
	return _NodeRegistry.Contract.Admin(&_NodeRegistry.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address admin_)
func (_NodeRegistry *NodeRegistryCallerSession) Admin() (common.Address, error) {
	return _NodeRegistry.Contract.Admin(&_NodeRegistry.CallOpts)
}

// AdminParameterKey is a free data retrieval call binding the contract method 0x9f40b625.
//
// Solidity: function adminParameterKey() pure returns(bytes key_)
func (_NodeRegistry *NodeRegistryCaller) AdminParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "adminParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// AdminParameterKey is a free data retrieval call binding the contract method 0x9f40b625.
//
// Solidity: function adminParameterKey() pure returns(bytes key_)
func (_NodeRegistry *NodeRegistrySession) AdminParameterKey() ([]byte, error) {
	return _NodeRegistry.Contract.AdminParameterKey(&_NodeRegistry.CallOpts)
}

// AdminParameterKey is a free data retrieval call binding the contract method 0x9f40b625.
//
// Solidity: function adminParameterKey() pure returns(bytes key_)
func (_NodeRegistry *NodeRegistryCallerSession) AdminParameterKey() ([]byte, error) {
	return _NodeRegistry.Contract.AdminParameterKey(&_NodeRegistry.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_NodeRegistry *NodeRegistryCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_NodeRegistry *NodeRegistrySession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _NodeRegistry.Contract.BalanceOf(&_NodeRegistry.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_NodeRegistry *NodeRegistryCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _NodeRegistry.Contract.BalanceOf(&_NodeRegistry.CallOpts, owner)
}

// CanonicalNodesCount is a free data retrieval call binding the contract method 0xc9c02a02.
//
// Solidity: function canonicalNodesCount() view returns(uint8 canonicalNodesCount_)
func (_NodeRegistry *NodeRegistryCaller) CanonicalNodesCount(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "canonicalNodesCount")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// CanonicalNodesCount is a free data retrieval call binding the contract method 0xc9c02a02.
//
// Solidity: function canonicalNodesCount() view returns(uint8 canonicalNodesCount_)
func (_NodeRegistry *NodeRegistrySession) CanonicalNodesCount() (uint8, error) {
	return _NodeRegistry.Contract.CanonicalNodesCount(&_NodeRegistry.CallOpts)
}

// CanonicalNodesCount is a free data retrieval call binding the contract method 0xc9c02a02.
//
// Solidity: function canonicalNodesCount() view returns(uint8 canonicalNodesCount_)
func (_NodeRegistry *NodeRegistryCallerSession) CanonicalNodesCount() (uint8, error) {
	return _NodeRegistry.Contract.CanonicalNodesCount(&_NodeRegistry.CallOpts)
}

// GetAllNodes is a free data retrieval call binding the contract method 0xa1174e7d.
//
// Solidity: function getAllNodes() view returns((uint32,(address,bool,bytes,string))[] allNodes_)
func (_NodeRegistry *NodeRegistryCaller) GetAllNodes(opts *bind.CallOpts) ([]INodeRegistryNodeWithId, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "getAllNodes")

	if err != nil {
		return *new([]INodeRegistryNodeWithId), err
	}

	out0 := *abi.ConvertType(out[0], new([]INodeRegistryNodeWithId)).(*[]INodeRegistryNodeWithId)

	return out0, err

}

// GetAllNodes is a free data retrieval call binding the contract method 0xa1174e7d.
//
// Solidity: function getAllNodes() view returns((uint32,(address,bool,bytes,string))[] allNodes_)
func (_NodeRegistry *NodeRegistrySession) GetAllNodes() ([]INodeRegistryNodeWithId, error) {
	return _NodeRegistry.Contract.GetAllNodes(&_NodeRegistry.CallOpts)
}

// GetAllNodes is a free data retrieval call binding the contract method 0xa1174e7d.
//
// Solidity: function getAllNodes() view returns((uint32,(address,bool,bytes,string))[] allNodes_)
func (_NodeRegistry *NodeRegistryCallerSession) GetAllNodes() ([]INodeRegistryNodeWithId, error) {
	return _NodeRegistry.Contract.GetAllNodes(&_NodeRegistry.CallOpts)
}

// GetAllNodesCount is a free data retrieval call binding the contract method 0x50d0215f.
//
// Solidity: function getAllNodesCount() view returns(uint32 nodeCount_)
func (_NodeRegistry *NodeRegistryCaller) GetAllNodesCount(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "getAllNodesCount")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GetAllNodesCount is a free data retrieval call binding the contract method 0x50d0215f.
//
// Solidity: function getAllNodesCount() view returns(uint32 nodeCount_)
func (_NodeRegistry *NodeRegistrySession) GetAllNodesCount() (uint32, error) {
	return _NodeRegistry.Contract.GetAllNodesCount(&_NodeRegistry.CallOpts)
}

// GetAllNodesCount is a free data retrieval call binding the contract method 0x50d0215f.
//
// Solidity: function getAllNodesCount() view returns(uint32 nodeCount_)
func (_NodeRegistry *NodeRegistryCallerSession) GetAllNodesCount() (uint32, error) {
	return _NodeRegistry.Contract.GetAllNodesCount(&_NodeRegistry.CallOpts)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_NodeRegistry *NodeRegistryCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_NodeRegistry *NodeRegistrySession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _NodeRegistry.Contract.GetApproved(&_NodeRegistry.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_NodeRegistry *NodeRegistryCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _NodeRegistry.Contract.GetApproved(&_NodeRegistry.CallOpts, tokenId)
}

// GetIsCanonicalNode is a free data retrieval call binding the contract method 0xad03d0a5.
//
// Solidity: function getIsCanonicalNode(uint32 nodeId_) view returns(bool isCanonicalNode_)
func (_NodeRegistry *NodeRegistryCaller) GetIsCanonicalNode(opts *bind.CallOpts, nodeId_ uint32) (bool, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "getIsCanonicalNode", nodeId_)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsCanonicalNode is a free data retrieval call binding the contract method 0xad03d0a5.
//
// Solidity: function getIsCanonicalNode(uint32 nodeId_) view returns(bool isCanonicalNode_)
func (_NodeRegistry *NodeRegistrySession) GetIsCanonicalNode(nodeId_ uint32) (bool, error) {
	return _NodeRegistry.Contract.GetIsCanonicalNode(&_NodeRegistry.CallOpts, nodeId_)
}

// GetIsCanonicalNode is a free data retrieval call binding the contract method 0xad03d0a5.
//
// Solidity: function getIsCanonicalNode(uint32 nodeId_) view returns(bool isCanonicalNode_)
func (_NodeRegistry *NodeRegistryCallerSession) GetIsCanonicalNode(nodeId_ uint32) (bool, error) {
	return _NodeRegistry.Contract.GetIsCanonicalNode(&_NodeRegistry.CallOpts, nodeId_)
}

// GetNode is a free data retrieval call binding the contract method 0xe06f876f.
//
// Solidity: function getNode(uint32 nodeId_) view returns((address,bool,bytes,string) node_)
func (_NodeRegistry *NodeRegistryCaller) GetNode(opts *bind.CallOpts, nodeId_ uint32) (INodeRegistryNode, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "getNode", nodeId_)

	if err != nil {
		return *new(INodeRegistryNode), err
	}

	out0 := *abi.ConvertType(out[0], new(INodeRegistryNode)).(*INodeRegistryNode)

	return out0, err

}

// GetNode is a free data retrieval call binding the contract method 0xe06f876f.
//
// Solidity: function getNode(uint32 nodeId_) view returns((address,bool,bytes,string) node_)
func (_NodeRegistry *NodeRegistrySession) GetNode(nodeId_ uint32) (INodeRegistryNode, error) {
	return _NodeRegistry.Contract.GetNode(&_NodeRegistry.CallOpts, nodeId_)
}

// GetNode is a free data retrieval call binding the contract method 0xe06f876f.
//
// Solidity: function getNode(uint32 nodeId_) view returns((address,bool,bytes,string) node_)
func (_NodeRegistry *NodeRegistryCallerSession) GetNode(nodeId_ uint32) (INodeRegistryNode, error) {
	return _NodeRegistry.Contract.GetNode(&_NodeRegistry.CallOpts, nodeId_)
}

// GetSigner is a free data retrieval call binding the contract method 0x68501a3e.
//
// Solidity: function getSigner(uint32 nodeId_) view returns(address signer_)
func (_NodeRegistry *NodeRegistryCaller) GetSigner(opts *bind.CallOpts, nodeId_ uint32) (common.Address, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "getSigner", nodeId_)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetSigner is a free data retrieval call binding the contract method 0x68501a3e.
//
// Solidity: function getSigner(uint32 nodeId_) view returns(address signer_)
func (_NodeRegistry *NodeRegistrySession) GetSigner(nodeId_ uint32) (common.Address, error) {
	return _NodeRegistry.Contract.GetSigner(&_NodeRegistry.CallOpts, nodeId_)
}

// GetSigner is a free data retrieval call binding the contract method 0x68501a3e.
//
// Solidity: function getSigner(uint32 nodeId_) view returns(address signer_)
func (_NodeRegistry *NodeRegistryCallerSession) GetSigner(nodeId_ uint32) (common.Address, error) {
	return _NodeRegistry.Contract.GetSigner(&_NodeRegistry.CallOpts, nodeId_)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_NodeRegistry *NodeRegistryCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_NodeRegistry *NodeRegistrySession) Implementation() (common.Address, error) {
	return _NodeRegistry.Contract.Implementation(&_NodeRegistry.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_NodeRegistry *NodeRegistryCallerSession) Implementation() (common.Address, error) {
	return _NodeRegistry.Contract.Implementation(&_NodeRegistry.CallOpts)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_NodeRegistry *NodeRegistryCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_NodeRegistry *NodeRegistrySession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _NodeRegistry.Contract.IsApprovedForAll(&_NodeRegistry.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_NodeRegistry *NodeRegistryCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _NodeRegistry.Contract.IsApprovedForAll(&_NodeRegistry.CallOpts, owner, operator)
}

// MaxCanonicalNodes is a free data retrieval call binding the contract method 0xc18e273d.
//
// Solidity: function maxCanonicalNodes() view returns(uint8 maxCanonicalNodes_)
func (_NodeRegistry *NodeRegistryCaller) MaxCanonicalNodes(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "maxCanonicalNodes")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// MaxCanonicalNodes is a free data retrieval call binding the contract method 0xc18e273d.
//
// Solidity: function maxCanonicalNodes() view returns(uint8 maxCanonicalNodes_)
func (_NodeRegistry *NodeRegistrySession) MaxCanonicalNodes() (uint8, error) {
	return _NodeRegistry.Contract.MaxCanonicalNodes(&_NodeRegistry.CallOpts)
}

// MaxCanonicalNodes is a free data retrieval call binding the contract method 0xc18e273d.
//
// Solidity: function maxCanonicalNodes() view returns(uint8 maxCanonicalNodes_)
func (_NodeRegistry *NodeRegistryCallerSession) MaxCanonicalNodes() (uint8, error) {
	return _NodeRegistry.Contract.MaxCanonicalNodes(&_NodeRegistry.CallOpts)
}

// MaxCanonicalNodesParameterKey is a free data retrieval call binding the contract method 0x0124b882.
//
// Solidity: function maxCanonicalNodesParameterKey() pure returns(bytes key_)
func (_NodeRegistry *NodeRegistryCaller) MaxCanonicalNodesParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "maxCanonicalNodesParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// MaxCanonicalNodesParameterKey is a free data retrieval call binding the contract method 0x0124b882.
//
// Solidity: function maxCanonicalNodesParameterKey() pure returns(bytes key_)
func (_NodeRegistry *NodeRegistrySession) MaxCanonicalNodesParameterKey() ([]byte, error) {
	return _NodeRegistry.Contract.MaxCanonicalNodesParameterKey(&_NodeRegistry.CallOpts)
}

// MaxCanonicalNodesParameterKey is a free data retrieval call binding the contract method 0x0124b882.
//
// Solidity: function maxCanonicalNodesParameterKey() pure returns(bytes key_)
func (_NodeRegistry *NodeRegistryCallerSession) MaxCanonicalNodesParameterKey() ([]byte, error) {
	return _NodeRegistry.Contract.MaxCanonicalNodesParameterKey(&_NodeRegistry.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(bytes key_)
func (_NodeRegistry *NodeRegistryCaller) MigratorParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "migratorParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(bytes key_)
func (_NodeRegistry *NodeRegistrySession) MigratorParameterKey() ([]byte, error) {
	return _NodeRegistry.Contract.MigratorParameterKey(&_NodeRegistry.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(bytes key_)
func (_NodeRegistry *NodeRegistryCallerSession) MigratorParameterKey() ([]byte, error) {
	return _NodeRegistry.Contract.MigratorParameterKey(&_NodeRegistry.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_NodeRegistry *NodeRegistryCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_NodeRegistry *NodeRegistrySession) Name() (string, error) {
	return _NodeRegistry.Contract.Name(&_NodeRegistry.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_NodeRegistry *NodeRegistryCallerSession) Name() (string, error) {
	return _NodeRegistry.Contract.Name(&_NodeRegistry.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_NodeRegistry *NodeRegistryCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_NodeRegistry *NodeRegistrySession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _NodeRegistry.Contract.OwnerOf(&_NodeRegistry.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_NodeRegistry *NodeRegistryCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _NodeRegistry.Contract.OwnerOf(&_NodeRegistry.CallOpts, tokenId)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_NodeRegistry *NodeRegistryCaller) ParameterRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "parameterRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_NodeRegistry *NodeRegistrySession) ParameterRegistry() (common.Address, error) {
	return _NodeRegistry.Contract.ParameterRegistry(&_NodeRegistry.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_NodeRegistry *NodeRegistryCallerSession) ParameterRegistry() (common.Address, error) {
	return _NodeRegistry.Contract.ParameterRegistry(&_NodeRegistry.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_NodeRegistry *NodeRegistryCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_NodeRegistry *NodeRegistrySession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _NodeRegistry.Contract.SupportsInterface(&_NodeRegistry.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_NodeRegistry *NodeRegistryCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _NodeRegistry.Contract.SupportsInterface(&_NodeRegistry.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_NodeRegistry *NodeRegistryCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_NodeRegistry *NodeRegistrySession) Symbol() (string, error) {
	return _NodeRegistry.Contract.Symbol(&_NodeRegistry.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_NodeRegistry *NodeRegistryCallerSession) Symbol() (string, error) {
	return _NodeRegistry.Contract.Symbol(&_NodeRegistry.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_NodeRegistry *NodeRegistryCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_NodeRegistry *NodeRegistrySession) TokenURI(tokenId *big.Int) (string, error) {
	return _NodeRegistry.Contract.TokenURI(&_NodeRegistry.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_NodeRegistry *NodeRegistryCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _NodeRegistry.Contract.TokenURI(&_NodeRegistry.CallOpts, tokenId)
}

// AddNode is a paid mutator transaction binding the contract method 0xa0eae81d.
//
// Solidity: function addNode(address owner_, bytes signingPublicKey_, string httpAddress_) returns(uint32 nodeId_, address signer_)
func (_NodeRegistry *NodeRegistryTransactor) AddNode(opts *bind.TransactOpts, owner_ common.Address, signingPublicKey_ []byte, httpAddress_ string) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "addNode", owner_, signingPublicKey_, httpAddress_)
}

// AddNode is a paid mutator transaction binding the contract method 0xa0eae81d.
//
// Solidity: function addNode(address owner_, bytes signingPublicKey_, string httpAddress_) returns(uint32 nodeId_, address signer_)
func (_NodeRegistry *NodeRegistrySession) AddNode(owner_ common.Address, signingPublicKey_ []byte, httpAddress_ string) (*types.Transaction, error) {
	return _NodeRegistry.Contract.AddNode(&_NodeRegistry.TransactOpts, owner_, signingPublicKey_, httpAddress_)
}

// AddNode is a paid mutator transaction binding the contract method 0xa0eae81d.
//
// Solidity: function addNode(address owner_, bytes signingPublicKey_, string httpAddress_) returns(uint32 nodeId_, address signer_)
func (_NodeRegistry *NodeRegistryTransactorSession) AddNode(owner_ common.Address, signingPublicKey_ []byte, httpAddress_ string) (*types.Transaction, error) {
	return _NodeRegistry.Contract.AddNode(&_NodeRegistry.TransactOpts, owner_, signingPublicKey_, httpAddress_)
}

// AddToNetwork is a paid mutator transaction binding the contract method 0x236b6eb8.
//
// Solidity: function addToNetwork(uint32 nodeId_) returns()
func (_NodeRegistry *NodeRegistryTransactor) AddToNetwork(opts *bind.TransactOpts, nodeId_ uint32) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "addToNetwork", nodeId_)
}

// AddToNetwork is a paid mutator transaction binding the contract method 0x236b6eb8.
//
// Solidity: function addToNetwork(uint32 nodeId_) returns()
func (_NodeRegistry *NodeRegistrySession) AddToNetwork(nodeId_ uint32) (*types.Transaction, error) {
	return _NodeRegistry.Contract.AddToNetwork(&_NodeRegistry.TransactOpts, nodeId_)
}

// AddToNetwork is a paid mutator transaction binding the contract method 0x236b6eb8.
//
// Solidity: function addToNetwork(uint32 nodeId_) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) AddToNetwork(nodeId_ uint32) (*types.Transaction, error) {
	return _NodeRegistry.Contract.AddToNetwork(&_NodeRegistry.TransactOpts, nodeId_)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistryTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistrySession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.Approve(&_NodeRegistry.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.Approve(&_NodeRegistry.TransactOpts, to, tokenId)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_NodeRegistry *NodeRegistryTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_NodeRegistry *NodeRegistrySession) Initialize() (*types.Transaction, error) {
	return _NodeRegistry.Contract.Initialize(&_NodeRegistry.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_NodeRegistry *NodeRegistryTransactorSession) Initialize() (*types.Transaction, error) {
	return _NodeRegistry.Contract.Initialize(&_NodeRegistry.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_NodeRegistry *NodeRegistryTransactor) Migrate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "migrate")
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_NodeRegistry *NodeRegistrySession) Migrate() (*types.Transaction, error) {
	return _NodeRegistry.Contract.Migrate(&_NodeRegistry.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_NodeRegistry *NodeRegistryTransactorSession) Migrate() (*types.Transaction, error) {
	return _NodeRegistry.Contract.Migrate(&_NodeRegistry.TransactOpts)
}

// RemoveFromNetwork is a paid mutator transaction binding the contract method 0x8cf20c68.
//
// Solidity: function removeFromNetwork(uint32 nodeId_) returns()
func (_NodeRegistry *NodeRegistryTransactor) RemoveFromNetwork(opts *bind.TransactOpts, nodeId_ uint32) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "removeFromNetwork", nodeId_)
}

// RemoveFromNetwork is a paid mutator transaction binding the contract method 0x8cf20c68.
//
// Solidity: function removeFromNetwork(uint32 nodeId_) returns()
func (_NodeRegistry *NodeRegistrySession) RemoveFromNetwork(nodeId_ uint32) (*types.Transaction, error) {
	return _NodeRegistry.Contract.RemoveFromNetwork(&_NodeRegistry.TransactOpts, nodeId_)
}

// RemoveFromNetwork is a paid mutator transaction binding the contract method 0x8cf20c68.
//
// Solidity: function removeFromNetwork(uint32 nodeId_) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) RemoveFromNetwork(nodeId_ uint32) (*types.Transaction, error) {
	return _NodeRegistry.Contract.RemoveFromNetwork(&_NodeRegistry.TransactOpts, nodeId_)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistryTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistrySession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SafeTransferFrom(&_NodeRegistry.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SafeTransferFrom(&_NodeRegistry.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_NodeRegistry *NodeRegistryTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_NodeRegistry *NodeRegistrySession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SafeTransferFrom0(&_NodeRegistry.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SafeTransferFrom0(&_NodeRegistry.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_NodeRegistry *NodeRegistryTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_NodeRegistry *NodeRegistrySession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetApprovalForAll(&_NodeRegistry.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetApprovalForAll(&_NodeRegistry.TransactOpts, operator, approved)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string baseURI_) returns()
func (_NodeRegistry *NodeRegistryTransactor) SetBaseURI(opts *bind.TransactOpts, baseURI_ string) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "setBaseURI", baseURI_)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string baseURI_) returns()
func (_NodeRegistry *NodeRegistrySession) SetBaseURI(baseURI_ string) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetBaseURI(&_NodeRegistry.TransactOpts, baseURI_)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string baseURI_) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) SetBaseURI(baseURI_ string) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetBaseURI(&_NodeRegistry.TransactOpts, baseURI_)
}

// SetHttpAddress is a paid mutator transaction binding the contract method 0xf84ce8b9.
//
// Solidity: function setHttpAddress(uint32 nodeId_, string httpAddress_) returns()
func (_NodeRegistry *NodeRegistryTransactor) SetHttpAddress(opts *bind.TransactOpts, nodeId_ uint32, httpAddress_ string) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "setHttpAddress", nodeId_, httpAddress_)
}

// SetHttpAddress is a paid mutator transaction binding the contract method 0xf84ce8b9.
//
// Solidity: function setHttpAddress(uint32 nodeId_, string httpAddress_) returns()
func (_NodeRegistry *NodeRegistrySession) SetHttpAddress(nodeId_ uint32, httpAddress_ string) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetHttpAddress(&_NodeRegistry.TransactOpts, nodeId_, httpAddress_)
}

// SetHttpAddress is a paid mutator transaction binding the contract method 0xf84ce8b9.
//
// Solidity: function setHttpAddress(uint32 nodeId_, string httpAddress_) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) SetHttpAddress(nodeId_ uint32, httpAddress_ string) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetHttpAddress(&_NodeRegistry.TransactOpts, nodeId_, httpAddress_)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistryTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistrySession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.TransferFrom(&_NodeRegistry.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.TransferFrom(&_NodeRegistry.TransactOpts, from, to, tokenId)
}

// UpdateAdmin is a paid mutator transaction binding the contract method 0xd3b2f598.
//
// Solidity: function updateAdmin() returns()
func (_NodeRegistry *NodeRegistryTransactor) UpdateAdmin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "updateAdmin")
}

// UpdateAdmin is a paid mutator transaction binding the contract method 0xd3b2f598.
//
// Solidity: function updateAdmin() returns()
func (_NodeRegistry *NodeRegistrySession) UpdateAdmin() (*types.Transaction, error) {
	return _NodeRegistry.Contract.UpdateAdmin(&_NodeRegistry.TransactOpts)
}

// UpdateAdmin is a paid mutator transaction binding the contract method 0xd3b2f598.
//
// Solidity: function updateAdmin() returns()
func (_NodeRegistry *NodeRegistryTransactorSession) UpdateAdmin() (*types.Transaction, error) {
	return _NodeRegistry.Contract.UpdateAdmin(&_NodeRegistry.TransactOpts)
}

// UpdateMaxCanonicalNodes is a paid mutator transaction binding the contract method 0x82a5cfc3.
//
// Solidity: function updateMaxCanonicalNodes() returns()
func (_NodeRegistry *NodeRegistryTransactor) UpdateMaxCanonicalNodes(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "updateMaxCanonicalNodes")
}

// UpdateMaxCanonicalNodes is a paid mutator transaction binding the contract method 0x82a5cfc3.
//
// Solidity: function updateMaxCanonicalNodes() returns()
func (_NodeRegistry *NodeRegistrySession) UpdateMaxCanonicalNodes() (*types.Transaction, error) {
	return _NodeRegistry.Contract.UpdateMaxCanonicalNodes(&_NodeRegistry.TransactOpts)
}

// UpdateMaxCanonicalNodes is a paid mutator transaction binding the contract method 0x82a5cfc3.
//
// Solidity: function updateMaxCanonicalNodes() returns()
func (_NodeRegistry *NodeRegistryTransactorSession) UpdateMaxCanonicalNodes() (*types.Transaction, error) {
	return _NodeRegistry.Contract.UpdateMaxCanonicalNodes(&_NodeRegistry.TransactOpts)
}

// NodeRegistryAdminUpdatedIterator is returned from FilterAdminUpdated and is used to iterate over the raw logs and unpacked data for AdminUpdated events raised by the NodeRegistry contract.
type NodeRegistryAdminUpdatedIterator struct {
	Event *NodeRegistryAdminUpdated // Event containing the contract specifics and raw log

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
func (it *NodeRegistryAdminUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryAdminUpdated)
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
		it.Event = new(NodeRegistryAdminUpdated)
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
func (it *NodeRegistryAdminUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryAdminUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryAdminUpdated represents a AdminUpdated event raised by the NodeRegistry contract.
type NodeRegistryAdminUpdated struct {
	Admin common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterAdminUpdated is a free log retrieval operation binding the contract event 0x54e4612788f90384e6843298d7854436f3a585b2c3831ab66abf1de63bfa6c2d.
//
// Solidity: event AdminUpdated(address indexed admin)
func (_NodeRegistry *NodeRegistryFilterer) FilterAdminUpdated(opts *bind.FilterOpts, admin []common.Address) (*NodeRegistryAdminUpdatedIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "AdminUpdated", adminRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryAdminUpdatedIterator{contract: _NodeRegistry.contract, event: "AdminUpdated", logs: logs, sub: sub}, nil
}

// WatchAdminUpdated is a free log subscription operation binding the contract event 0x54e4612788f90384e6843298d7854436f3a585b2c3831ab66abf1de63bfa6c2d.
//
// Solidity: event AdminUpdated(address indexed admin)
func (_NodeRegistry *NodeRegistryFilterer) WatchAdminUpdated(opts *bind.WatchOpts, sink chan<- *NodeRegistryAdminUpdated, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "AdminUpdated", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryAdminUpdated)
				if err := _NodeRegistry.contract.UnpackLog(event, "AdminUpdated", log); err != nil {
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

// ParseAdminUpdated is a log parse operation binding the contract event 0x54e4612788f90384e6843298d7854436f3a585b2c3831ab66abf1de63bfa6c2d.
//
// Solidity: event AdminUpdated(address indexed admin)
func (_NodeRegistry *NodeRegistryFilterer) ParseAdminUpdated(log types.Log) (*NodeRegistryAdminUpdated, error) {
	event := new(NodeRegistryAdminUpdated)
	if err := _NodeRegistry.contract.UnpackLog(event, "AdminUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the NodeRegistry contract.
type NodeRegistryApprovalIterator struct {
	Event *NodeRegistryApproval // Event containing the contract specifics and raw log

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
func (it *NodeRegistryApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryApproval)
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
		it.Event = new(NodeRegistryApproval)
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
func (it *NodeRegistryApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryApproval represents a Approval event raised by the NodeRegistry contract.
type NodeRegistryApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_NodeRegistry *NodeRegistryFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*NodeRegistryApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryApprovalIterator{contract: _NodeRegistry.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_NodeRegistry *NodeRegistryFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *NodeRegistryApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryApproval)
				if err := _NodeRegistry.contract.UnpackLog(event, "Approval", log); err != nil {
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
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_NodeRegistry *NodeRegistryFilterer) ParseApproval(log types.Log) (*NodeRegistryApproval, error) {
	event := new(NodeRegistryApproval)
	if err := _NodeRegistry.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the NodeRegistry contract.
type NodeRegistryApprovalForAllIterator struct {
	Event *NodeRegistryApprovalForAll // Event containing the contract specifics and raw log

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
func (it *NodeRegistryApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryApprovalForAll)
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
		it.Event = new(NodeRegistryApprovalForAll)
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
func (it *NodeRegistryApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryApprovalForAll represents a ApprovalForAll event raised by the NodeRegistry contract.
type NodeRegistryApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_NodeRegistry *NodeRegistryFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*NodeRegistryApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryApprovalForAllIterator{contract: _NodeRegistry.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_NodeRegistry *NodeRegistryFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *NodeRegistryApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryApprovalForAll)
				if err := _NodeRegistry.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_NodeRegistry *NodeRegistryFilterer) ParseApprovalForAll(log types.Log) (*NodeRegistryApprovalForAll, error) {
	event := new(NodeRegistryApprovalForAll)
	if err := _NodeRegistry.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryBaseURIUpdatedIterator is returned from FilterBaseURIUpdated and is used to iterate over the raw logs and unpacked data for BaseURIUpdated events raised by the NodeRegistry contract.
type NodeRegistryBaseURIUpdatedIterator struct {
	Event *NodeRegistryBaseURIUpdated // Event containing the contract specifics and raw log

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
func (it *NodeRegistryBaseURIUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryBaseURIUpdated)
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
		it.Event = new(NodeRegistryBaseURIUpdated)
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
func (it *NodeRegistryBaseURIUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryBaseURIUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryBaseURIUpdated represents a BaseURIUpdated event raised by the NodeRegistry contract.
type NodeRegistryBaseURIUpdated struct {
	BaseURI string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBaseURIUpdated is a free log retrieval operation binding the contract event 0x6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad.
//
// Solidity: event BaseURIUpdated(string baseURI)
func (_NodeRegistry *NodeRegistryFilterer) FilterBaseURIUpdated(opts *bind.FilterOpts) (*NodeRegistryBaseURIUpdatedIterator, error) {

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "BaseURIUpdated")
	if err != nil {
		return nil, err
	}
	return &NodeRegistryBaseURIUpdatedIterator{contract: _NodeRegistry.contract, event: "BaseURIUpdated", logs: logs, sub: sub}, nil
}

// WatchBaseURIUpdated is a free log subscription operation binding the contract event 0x6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad.
//
// Solidity: event BaseURIUpdated(string baseURI)
func (_NodeRegistry *NodeRegistryFilterer) WatchBaseURIUpdated(opts *bind.WatchOpts, sink chan<- *NodeRegistryBaseURIUpdated) (event.Subscription, error) {

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "BaseURIUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryBaseURIUpdated)
				if err := _NodeRegistry.contract.UnpackLog(event, "BaseURIUpdated", log); err != nil {
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

// ParseBaseURIUpdated is a log parse operation binding the contract event 0x6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad.
//
// Solidity: event BaseURIUpdated(string baseURI)
func (_NodeRegistry *NodeRegistryFilterer) ParseBaseURIUpdated(log types.Log) (*NodeRegistryBaseURIUpdated, error) {
	event := new(NodeRegistryBaseURIUpdated)
	if err := _NodeRegistry.contract.UnpackLog(event, "BaseURIUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryHttpAddressUpdatedIterator is returned from FilterHttpAddressUpdated and is used to iterate over the raw logs and unpacked data for HttpAddressUpdated events raised by the NodeRegistry contract.
type NodeRegistryHttpAddressUpdatedIterator struct {
	Event *NodeRegistryHttpAddressUpdated // Event containing the contract specifics and raw log

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
func (it *NodeRegistryHttpAddressUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryHttpAddressUpdated)
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
		it.Event = new(NodeRegistryHttpAddressUpdated)
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
func (it *NodeRegistryHttpAddressUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryHttpAddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryHttpAddressUpdated represents a HttpAddressUpdated event raised by the NodeRegistry contract.
type NodeRegistryHttpAddressUpdated struct {
	NodeId      uint32
	HttpAddress string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterHttpAddressUpdated is a free log retrieval operation binding the contract event 0x5698a22512088407e91d125d2eb43d829d9694a71f664ab0dc2aea3a8e524712.
//
// Solidity: event HttpAddressUpdated(uint32 indexed nodeId, string httpAddress)
func (_NodeRegistry *NodeRegistryFilterer) FilterHttpAddressUpdated(opts *bind.FilterOpts, nodeId []uint32) (*NodeRegistryHttpAddressUpdatedIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "HttpAddressUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryHttpAddressUpdatedIterator{contract: _NodeRegistry.contract, event: "HttpAddressUpdated", logs: logs, sub: sub}, nil
}

// WatchHttpAddressUpdated is a free log subscription operation binding the contract event 0x5698a22512088407e91d125d2eb43d829d9694a71f664ab0dc2aea3a8e524712.
//
// Solidity: event HttpAddressUpdated(uint32 indexed nodeId, string httpAddress)
func (_NodeRegistry *NodeRegistryFilterer) WatchHttpAddressUpdated(opts *bind.WatchOpts, sink chan<- *NodeRegistryHttpAddressUpdated, nodeId []uint32) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "HttpAddressUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryHttpAddressUpdated)
				if err := _NodeRegistry.contract.UnpackLog(event, "HttpAddressUpdated", log); err != nil {
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

// ParseHttpAddressUpdated is a log parse operation binding the contract event 0x5698a22512088407e91d125d2eb43d829d9694a71f664ab0dc2aea3a8e524712.
//
// Solidity: event HttpAddressUpdated(uint32 indexed nodeId, string httpAddress)
func (_NodeRegistry *NodeRegistryFilterer) ParseHttpAddressUpdated(log types.Log) (*NodeRegistryHttpAddressUpdated, error) {
	event := new(NodeRegistryHttpAddressUpdated)
	if err := _NodeRegistry.contract.UnpackLog(event, "HttpAddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the NodeRegistry contract.
type NodeRegistryInitializedIterator struct {
	Event *NodeRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *NodeRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryInitialized)
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
		it.Event = new(NodeRegistryInitialized)
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
func (it *NodeRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryInitialized represents a Initialized event raised by the NodeRegistry contract.
type NodeRegistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_NodeRegistry *NodeRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*NodeRegistryInitializedIterator, error) {

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &NodeRegistryInitializedIterator{contract: _NodeRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_NodeRegistry *NodeRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *NodeRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryInitialized)
				if err := _NodeRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_NodeRegistry *NodeRegistryFilterer) ParseInitialized(log types.Log) (*NodeRegistryInitialized, error) {
	event := new(NodeRegistryInitialized)
	if err := _NodeRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryMaxCanonicalNodesUpdatedIterator is returned from FilterMaxCanonicalNodesUpdated and is used to iterate over the raw logs and unpacked data for MaxCanonicalNodesUpdated events raised by the NodeRegistry contract.
type NodeRegistryMaxCanonicalNodesUpdatedIterator struct {
	Event *NodeRegistryMaxCanonicalNodesUpdated // Event containing the contract specifics and raw log

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
func (it *NodeRegistryMaxCanonicalNodesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryMaxCanonicalNodesUpdated)
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
		it.Event = new(NodeRegistryMaxCanonicalNodesUpdated)
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
func (it *NodeRegistryMaxCanonicalNodesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryMaxCanonicalNodesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryMaxCanonicalNodesUpdated represents a MaxCanonicalNodesUpdated event raised by the NodeRegistry contract.
type NodeRegistryMaxCanonicalNodesUpdated struct {
	MaxCanonicalNodes uint8
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterMaxCanonicalNodesUpdated is a free log retrieval operation binding the contract event 0x581c4d2fc386422e99f02a47a9735e8936050b0c2a384b98c8a6740786d9ff76.
//
// Solidity: event MaxCanonicalNodesUpdated(uint8 maxCanonicalNodes)
func (_NodeRegistry *NodeRegistryFilterer) FilterMaxCanonicalNodesUpdated(opts *bind.FilterOpts) (*NodeRegistryMaxCanonicalNodesUpdatedIterator, error) {

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "MaxCanonicalNodesUpdated")
	if err != nil {
		return nil, err
	}
	return &NodeRegistryMaxCanonicalNodesUpdatedIterator{contract: _NodeRegistry.contract, event: "MaxCanonicalNodesUpdated", logs: logs, sub: sub}, nil
}

// WatchMaxCanonicalNodesUpdated is a free log subscription operation binding the contract event 0x581c4d2fc386422e99f02a47a9735e8936050b0c2a384b98c8a6740786d9ff76.
//
// Solidity: event MaxCanonicalNodesUpdated(uint8 maxCanonicalNodes)
func (_NodeRegistry *NodeRegistryFilterer) WatchMaxCanonicalNodesUpdated(opts *bind.WatchOpts, sink chan<- *NodeRegistryMaxCanonicalNodesUpdated) (event.Subscription, error) {

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "MaxCanonicalNodesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryMaxCanonicalNodesUpdated)
				if err := _NodeRegistry.contract.UnpackLog(event, "MaxCanonicalNodesUpdated", log); err != nil {
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

// ParseMaxCanonicalNodesUpdated is a log parse operation binding the contract event 0x581c4d2fc386422e99f02a47a9735e8936050b0c2a384b98c8a6740786d9ff76.
//
// Solidity: event MaxCanonicalNodesUpdated(uint8 maxCanonicalNodes)
func (_NodeRegistry *NodeRegistryFilterer) ParseMaxCanonicalNodesUpdated(log types.Log) (*NodeRegistryMaxCanonicalNodesUpdated, error) {
	event := new(NodeRegistryMaxCanonicalNodesUpdated)
	if err := _NodeRegistry.contract.UnpackLog(event, "MaxCanonicalNodesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryMigratedIterator is returned from FilterMigrated and is used to iterate over the raw logs and unpacked data for Migrated events raised by the NodeRegistry contract.
type NodeRegistryMigratedIterator struct {
	Event *NodeRegistryMigrated // Event containing the contract specifics and raw log

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
func (it *NodeRegistryMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryMigrated)
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
		it.Event = new(NodeRegistryMigrated)
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
func (it *NodeRegistryMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryMigrated represents a Migrated event raised by the NodeRegistry contract.
type NodeRegistryMigrated struct {
	Migrator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMigrated is a free log retrieval operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_NodeRegistry *NodeRegistryFilterer) FilterMigrated(opts *bind.FilterOpts, migrator []common.Address) (*NodeRegistryMigratedIterator, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryMigratedIterator{contract: _NodeRegistry.contract, event: "Migrated", logs: logs, sub: sub}, nil
}

// WatchMigrated is a free log subscription operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_NodeRegistry *NodeRegistryFilterer) WatchMigrated(opts *bind.WatchOpts, sink chan<- *NodeRegistryMigrated, migrator []common.Address) (event.Subscription, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryMigrated)
				if err := _NodeRegistry.contract.UnpackLog(event, "Migrated", log); err != nil {
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
func (_NodeRegistry *NodeRegistryFilterer) ParseMigrated(log types.Log) (*NodeRegistryMigrated, error) {
	event := new(NodeRegistryMigrated)
	if err := _NodeRegistry.contract.UnpackLog(event, "Migrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryNodeAddedIterator is returned from FilterNodeAdded and is used to iterate over the raw logs and unpacked data for NodeAdded events raised by the NodeRegistry contract.
type NodeRegistryNodeAddedIterator struct {
	Event *NodeRegistryNodeAdded // Event containing the contract specifics and raw log

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
func (it *NodeRegistryNodeAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryNodeAdded)
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
		it.Event = new(NodeRegistryNodeAdded)
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
func (it *NodeRegistryNodeAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryNodeAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryNodeAdded represents a NodeAdded event raised by the NodeRegistry contract.
type NodeRegistryNodeAdded struct {
	NodeId           uint32
	Owner            common.Address
	Signer           common.Address
	SigningPublicKey []byte
	HttpAddress      string
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterNodeAdded is a free log retrieval operation binding the contract event 0x9b385c30e390e1e15ab8a2e34c4caa40b3c59882c17185fcbc3f87b2bf6658a4.
//
// Solidity: event NodeAdded(uint32 indexed nodeId, address indexed owner, address indexed signer, bytes signingPublicKey, string httpAddress)
func (_NodeRegistry *NodeRegistryFilterer) FilterNodeAdded(opts *bind.FilterOpts, nodeId []uint32, owner []common.Address, signer []common.Address) (*NodeRegistryNodeAddedIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "NodeAdded", nodeIdRule, ownerRule, signerRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryNodeAddedIterator{contract: _NodeRegistry.contract, event: "NodeAdded", logs: logs, sub: sub}, nil
}

// WatchNodeAdded is a free log subscription operation binding the contract event 0x9b385c30e390e1e15ab8a2e34c4caa40b3c59882c17185fcbc3f87b2bf6658a4.
//
// Solidity: event NodeAdded(uint32 indexed nodeId, address indexed owner, address indexed signer, bytes signingPublicKey, string httpAddress)
func (_NodeRegistry *NodeRegistryFilterer) WatchNodeAdded(opts *bind.WatchOpts, sink chan<- *NodeRegistryNodeAdded, nodeId []uint32, owner []common.Address, signer []common.Address) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "NodeAdded", nodeIdRule, ownerRule, signerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryNodeAdded)
				if err := _NodeRegistry.contract.UnpackLog(event, "NodeAdded", log); err != nil {
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

// ParseNodeAdded is a log parse operation binding the contract event 0x9b385c30e390e1e15ab8a2e34c4caa40b3c59882c17185fcbc3f87b2bf6658a4.
//
// Solidity: event NodeAdded(uint32 indexed nodeId, address indexed owner, address indexed signer, bytes signingPublicKey, string httpAddress)
func (_NodeRegistry *NodeRegistryFilterer) ParseNodeAdded(log types.Log) (*NodeRegistryNodeAdded, error) {
	event := new(NodeRegistryNodeAdded)
	if err := _NodeRegistry.contract.UnpackLog(event, "NodeAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryNodeAddedToCanonicalNetworkIterator is returned from FilterNodeAddedToCanonicalNetwork and is used to iterate over the raw logs and unpacked data for NodeAddedToCanonicalNetwork events raised by the NodeRegistry contract.
type NodeRegistryNodeAddedToCanonicalNetworkIterator struct {
	Event *NodeRegistryNodeAddedToCanonicalNetwork // Event containing the contract specifics and raw log

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
func (it *NodeRegistryNodeAddedToCanonicalNetworkIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryNodeAddedToCanonicalNetwork)
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
		it.Event = new(NodeRegistryNodeAddedToCanonicalNetwork)
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
func (it *NodeRegistryNodeAddedToCanonicalNetworkIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryNodeAddedToCanonicalNetworkIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryNodeAddedToCanonicalNetwork represents a NodeAddedToCanonicalNetwork event raised by the NodeRegistry contract.
type NodeRegistryNodeAddedToCanonicalNetwork struct {
	NodeId uint32
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNodeAddedToCanonicalNetwork is a free log retrieval operation binding the contract event 0x13695734a48552c5f7d826df6e02f4094ed655e28bcedb3ccc3645997f6b47f8.
//
// Solidity: event NodeAddedToCanonicalNetwork(uint32 indexed nodeId)
func (_NodeRegistry *NodeRegistryFilterer) FilterNodeAddedToCanonicalNetwork(opts *bind.FilterOpts, nodeId []uint32) (*NodeRegistryNodeAddedToCanonicalNetworkIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "NodeAddedToCanonicalNetwork", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryNodeAddedToCanonicalNetworkIterator{contract: _NodeRegistry.contract, event: "NodeAddedToCanonicalNetwork", logs: logs, sub: sub}, nil
}

// WatchNodeAddedToCanonicalNetwork is a free log subscription operation binding the contract event 0x13695734a48552c5f7d826df6e02f4094ed655e28bcedb3ccc3645997f6b47f8.
//
// Solidity: event NodeAddedToCanonicalNetwork(uint32 indexed nodeId)
func (_NodeRegistry *NodeRegistryFilterer) WatchNodeAddedToCanonicalNetwork(opts *bind.WatchOpts, sink chan<- *NodeRegistryNodeAddedToCanonicalNetwork, nodeId []uint32) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "NodeAddedToCanonicalNetwork", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryNodeAddedToCanonicalNetwork)
				if err := _NodeRegistry.contract.UnpackLog(event, "NodeAddedToCanonicalNetwork", log); err != nil {
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

// ParseNodeAddedToCanonicalNetwork is a log parse operation binding the contract event 0x13695734a48552c5f7d826df6e02f4094ed655e28bcedb3ccc3645997f6b47f8.
//
// Solidity: event NodeAddedToCanonicalNetwork(uint32 indexed nodeId)
func (_NodeRegistry *NodeRegistryFilterer) ParseNodeAddedToCanonicalNetwork(log types.Log) (*NodeRegistryNodeAddedToCanonicalNetwork, error) {
	event := new(NodeRegistryNodeAddedToCanonicalNetwork)
	if err := _NodeRegistry.contract.UnpackLog(event, "NodeAddedToCanonicalNetwork", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryNodeRemovedFromCanonicalNetworkIterator is returned from FilterNodeRemovedFromCanonicalNetwork and is used to iterate over the raw logs and unpacked data for NodeRemovedFromCanonicalNetwork events raised by the NodeRegistry contract.
type NodeRegistryNodeRemovedFromCanonicalNetworkIterator struct {
	Event *NodeRegistryNodeRemovedFromCanonicalNetwork // Event containing the contract specifics and raw log

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
func (it *NodeRegistryNodeRemovedFromCanonicalNetworkIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryNodeRemovedFromCanonicalNetwork)
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
		it.Event = new(NodeRegistryNodeRemovedFromCanonicalNetwork)
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
func (it *NodeRegistryNodeRemovedFromCanonicalNetworkIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryNodeRemovedFromCanonicalNetworkIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryNodeRemovedFromCanonicalNetwork represents a NodeRemovedFromCanonicalNetwork event raised by the NodeRegistry contract.
type NodeRegistryNodeRemovedFromCanonicalNetwork struct {
	NodeId uint32
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNodeRemovedFromCanonicalNetwork is a free log retrieval operation binding the contract event 0x7cf9bcdd519495a485911496098851db2c18ee9a708b453dd48f2822098e16ec.
//
// Solidity: event NodeRemovedFromCanonicalNetwork(uint32 indexed nodeId)
func (_NodeRegistry *NodeRegistryFilterer) FilterNodeRemovedFromCanonicalNetwork(opts *bind.FilterOpts, nodeId []uint32) (*NodeRegistryNodeRemovedFromCanonicalNetworkIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "NodeRemovedFromCanonicalNetwork", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryNodeRemovedFromCanonicalNetworkIterator{contract: _NodeRegistry.contract, event: "NodeRemovedFromCanonicalNetwork", logs: logs, sub: sub}, nil
}

// WatchNodeRemovedFromCanonicalNetwork is a free log subscription operation binding the contract event 0x7cf9bcdd519495a485911496098851db2c18ee9a708b453dd48f2822098e16ec.
//
// Solidity: event NodeRemovedFromCanonicalNetwork(uint32 indexed nodeId)
func (_NodeRegistry *NodeRegistryFilterer) WatchNodeRemovedFromCanonicalNetwork(opts *bind.WatchOpts, sink chan<- *NodeRegistryNodeRemovedFromCanonicalNetwork, nodeId []uint32) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "NodeRemovedFromCanonicalNetwork", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryNodeRemovedFromCanonicalNetwork)
				if err := _NodeRegistry.contract.UnpackLog(event, "NodeRemovedFromCanonicalNetwork", log); err != nil {
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

// ParseNodeRemovedFromCanonicalNetwork is a log parse operation binding the contract event 0x7cf9bcdd519495a485911496098851db2c18ee9a708b453dd48f2822098e16ec.
//
// Solidity: event NodeRemovedFromCanonicalNetwork(uint32 indexed nodeId)
func (_NodeRegistry *NodeRegistryFilterer) ParseNodeRemovedFromCanonicalNetwork(log types.Log) (*NodeRegistryNodeRemovedFromCanonicalNetwork, error) {
	event := new(NodeRegistryNodeRemovedFromCanonicalNetwork)
	if err := _NodeRegistry.contract.UnpackLog(event, "NodeRemovedFromCanonicalNetwork", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the NodeRegistry contract.
type NodeRegistryTransferIterator struct {
	Event *NodeRegistryTransfer // Event containing the contract specifics and raw log

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
func (it *NodeRegistryTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryTransfer)
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
		it.Event = new(NodeRegistryTransfer)
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
func (it *NodeRegistryTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryTransfer represents a Transfer event raised by the NodeRegistry contract.
type NodeRegistryTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_NodeRegistry *NodeRegistryFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*NodeRegistryTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryTransferIterator{contract: _NodeRegistry.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_NodeRegistry *NodeRegistryFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *NodeRegistryTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryTransfer)
				if err := _NodeRegistry.contract.UnpackLog(event, "Transfer", log); err != nil {
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
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_NodeRegistry *NodeRegistryFilterer) ParseTransfer(log types.Log) (*NodeRegistryTransfer, error) {
	event := new(NodeRegistryTransfer)
	if err := _NodeRegistry.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the NodeRegistry contract.
type NodeRegistryUpgradedIterator struct {
	Event *NodeRegistryUpgraded // Event containing the contract specifics and raw log

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
func (it *NodeRegistryUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryUpgraded)
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
		it.Event = new(NodeRegistryUpgraded)
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
func (it *NodeRegistryUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryUpgraded represents a Upgraded event raised by the NodeRegistry contract.
type NodeRegistryUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_NodeRegistry *NodeRegistryFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*NodeRegistryUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryUpgradedIterator{contract: _NodeRegistry.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_NodeRegistry *NodeRegistryFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *NodeRegistryUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryUpgraded)
				if err := _NodeRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_NodeRegistry *NodeRegistryFilterer) ParseUpgraded(log types.Log) (*NodeRegistryUpgraded, error) {
	event := new(NodeRegistryUpgraded)
	if err := _NodeRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
