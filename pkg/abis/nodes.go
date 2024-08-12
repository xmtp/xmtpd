// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abis

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

// NodesMetaData contains all meta data concerning the Nodes contract.
var NodesMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addNode\",\"inputs\":[{\"name\":\"publicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"markNodeHealthy\",\"inputs\":[{\"name\":\"publicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"markNodeUnhealthy\",\"inputs\":[{\"name\":\"publicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"nodes\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"originatorId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isHealthy\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"NodeUpdate\",\"inputs\":[{\"name\":\"publicKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"originatorId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"isHealthy\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// NodesABI is the input ABI used to generate the binding from.
// Deprecated: Use NodesMetaData.ABI instead.
var NodesABI = NodesMetaData.ABI

// Nodes is an auto generated Go binding around an Ethereum contract.
type Nodes struct {
	NodesCaller     // Read-only binding to the contract
	NodesTransactor // Write-only binding to the contract
	NodesFilterer   // Log filterer for contract events
}

// NodesCaller is an auto generated read-only Go binding around an Ethereum contract.
type NodesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NodesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NodesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NodesSession struct {
	Contract     *Nodes            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NodesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NodesCallerSession struct {
	Contract *NodesCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// NodesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NodesTransactorSession struct {
	Contract     *NodesTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NodesRaw is an auto generated low-level Go binding around an Ethereum contract.
type NodesRaw struct {
	Contract *Nodes // Generic contract binding to access the raw methods on
}

// NodesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NodesCallerRaw struct {
	Contract *NodesCaller // Generic read-only contract binding to access the raw methods on
}

// NodesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NodesTransactorRaw struct {
	Contract *NodesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNodes creates a new instance of Nodes, bound to a specific deployed contract.
func NewNodes(address common.Address, backend bind.ContractBackend) (*Nodes, error) {
	contract, err := bindNodes(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Nodes{NodesCaller: NodesCaller{contract: contract}, NodesTransactor: NodesTransactor{contract: contract}, NodesFilterer: NodesFilterer{contract: contract}}, nil
}

// NewNodesCaller creates a new read-only instance of Nodes, bound to a specific deployed contract.
func NewNodesCaller(address common.Address, caller bind.ContractCaller) (*NodesCaller, error) {
	contract, err := bindNodes(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NodesCaller{contract: contract}, nil
}

// NewNodesTransactor creates a new write-only instance of Nodes, bound to a specific deployed contract.
func NewNodesTransactor(address common.Address, transactor bind.ContractTransactor) (*NodesTransactor, error) {
	contract, err := bindNodes(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NodesTransactor{contract: contract}, nil
}

// NewNodesFilterer creates a new log filterer instance of Nodes, bound to a specific deployed contract.
func NewNodesFilterer(address common.Address, filterer bind.ContractFilterer) (*NodesFilterer, error) {
	contract, err := bindNodes(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NodesFilterer{contract: contract}, nil
}

// bindNodes binds a generic wrapper to an already deployed contract.
func bindNodes(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NodesMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Nodes *NodesRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Nodes.Contract.NodesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Nodes *NodesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Nodes.Contract.NodesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Nodes *NodesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Nodes.Contract.NodesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Nodes *NodesCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Nodes.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Nodes *NodesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Nodes.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Nodes *NodesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Nodes.Contract.contract.Transact(opts, method, params...)
}

// Nodes is a free data retrieval call binding the contract method 0x404608bd.
//
// Solidity: function nodes(bytes ) view returns(string httpAddress, uint256 originatorId, bool isHealthy)
func (_Nodes *NodesCaller) Nodes(opts *bind.CallOpts, arg0 []byte) (struct {
	HttpAddress  string
	OriginatorId *big.Int
	IsHealthy    bool
}, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "nodes", arg0)

	outstruct := new(struct {
		HttpAddress  string
		OriginatorId *big.Int
		IsHealthy    bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.HttpAddress = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.OriginatorId = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.IsHealthy = *abi.ConvertType(out[2], new(bool)).(*bool)

	return *outstruct, err

}

// Nodes is a free data retrieval call binding the contract method 0x404608bd.
//
// Solidity: function nodes(bytes ) view returns(string httpAddress, uint256 originatorId, bool isHealthy)
func (_Nodes *NodesSession) Nodes(arg0 []byte) (struct {
	HttpAddress  string
	OriginatorId *big.Int
	IsHealthy    bool
}, error) {
	return _Nodes.Contract.Nodes(&_Nodes.CallOpts, arg0)
}

// Nodes is a free data retrieval call binding the contract method 0x404608bd.
//
// Solidity: function nodes(bytes ) view returns(string httpAddress, uint256 originatorId, bool isHealthy)
func (_Nodes *NodesCallerSession) Nodes(arg0 []byte) (struct {
	HttpAddress  string
	OriginatorId *big.Int
	IsHealthy    bool
}, error) {
	return _Nodes.Contract.Nodes(&_Nodes.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Nodes *NodesCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Nodes *NodesSession) Owner() (common.Address, error) {
	return _Nodes.Contract.Owner(&_Nodes.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Nodes *NodesCallerSession) Owner() (common.Address, error) {
	return _Nodes.Contract.Owner(&_Nodes.CallOpts)
}

// AddNode is a paid mutator transaction binding the contract method 0x5a48bd1d.
//
// Solidity: function addNode(bytes publicKey, string httpAddress) returns()
func (_Nodes *NodesTransactor) AddNode(opts *bind.TransactOpts, publicKey []byte, httpAddress string) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "addNode", publicKey, httpAddress)
}

// AddNode is a paid mutator transaction binding the contract method 0x5a48bd1d.
//
// Solidity: function addNode(bytes publicKey, string httpAddress) returns()
func (_Nodes *NodesSession) AddNode(publicKey []byte, httpAddress string) (*types.Transaction, error) {
	return _Nodes.Contract.AddNode(&_Nodes.TransactOpts, publicKey, httpAddress)
}

// AddNode is a paid mutator transaction binding the contract method 0x5a48bd1d.
//
// Solidity: function addNode(bytes publicKey, string httpAddress) returns()
func (_Nodes *NodesTransactorSession) AddNode(publicKey []byte, httpAddress string) (*types.Transaction, error) {
	return _Nodes.Contract.AddNode(&_Nodes.TransactOpts, publicKey, httpAddress)
}

// MarkNodeHealthy is a paid mutator transaction binding the contract method 0x97d6305a.
//
// Solidity: function markNodeHealthy(bytes publicKey) returns()
func (_Nodes *NodesTransactor) MarkNodeHealthy(opts *bind.TransactOpts, publicKey []byte) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "markNodeHealthy", publicKey)
}

// MarkNodeHealthy is a paid mutator transaction binding the contract method 0x97d6305a.
//
// Solidity: function markNodeHealthy(bytes publicKey) returns()
func (_Nodes *NodesSession) MarkNodeHealthy(publicKey []byte) (*types.Transaction, error) {
	return _Nodes.Contract.MarkNodeHealthy(&_Nodes.TransactOpts, publicKey)
}

// MarkNodeHealthy is a paid mutator transaction binding the contract method 0x97d6305a.
//
// Solidity: function markNodeHealthy(bytes publicKey) returns()
func (_Nodes *NodesTransactorSession) MarkNodeHealthy(publicKey []byte) (*types.Transaction, error) {
	return _Nodes.Contract.MarkNodeHealthy(&_Nodes.TransactOpts, publicKey)
}

// MarkNodeUnhealthy is a paid mutator transaction binding the contract method 0xf0da65bc.
//
// Solidity: function markNodeUnhealthy(bytes publicKey) returns()
func (_Nodes *NodesTransactor) MarkNodeUnhealthy(opts *bind.TransactOpts, publicKey []byte) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "markNodeUnhealthy", publicKey)
}

// MarkNodeUnhealthy is a paid mutator transaction binding the contract method 0xf0da65bc.
//
// Solidity: function markNodeUnhealthy(bytes publicKey) returns()
func (_Nodes *NodesSession) MarkNodeUnhealthy(publicKey []byte) (*types.Transaction, error) {
	return _Nodes.Contract.MarkNodeUnhealthy(&_Nodes.TransactOpts, publicKey)
}

// MarkNodeUnhealthy is a paid mutator transaction binding the contract method 0xf0da65bc.
//
// Solidity: function markNodeUnhealthy(bytes publicKey) returns()
func (_Nodes *NodesTransactorSession) MarkNodeUnhealthy(publicKey []byte) (*types.Transaction, error) {
	return _Nodes.Contract.MarkNodeUnhealthy(&_Nodes.TransactOpts, publicKey)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Nodes *NodesTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Nodes *NodesSession) RenounceOwnership() (*types.Transaction, error) {
	return _Nodes.Contract.RenounceOwnership(&_Nodes.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Nodes *NodesTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Nodes.Contract.RenounceOwnership(&_Nodes.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Nodes *NodesTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Nodes *NodesSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Nodes.Contract.TransferOwnership(&_Nodes.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Nodes *NodesTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Nodes.Contract.TransferOwnership(&_Nodes.TransactOpts, newOwner)
}

// NodesNodeUpdateIterator is returned from FilterNodeUpdate and is used to iterate over the raw logs and unpacked data for NodeUpdate events raised by the Nodes contract.
type NodesNodeUpdateIterator struct {
	Event *NodesNodeUpdate // Event containing the contract specifics and raw log

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
func (it *NodesNodeUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesNodeUpdate)
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
		it.Event = new(NodesNodeUpdate)
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
func (it *NodesNodeUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesNodeUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesNodeUpdate represents a NodeUpdate event raised by the Nodes contract.
type NodesNodeUpdate struct {
	PublicKey    []byte
	HttpAddress  string
	OriginatorId *big.Int
	IsHealthy    bool
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterNodeUpdate is a free log retrieval operation binding the contract event 0x94b61d0a2ca1645b5632d2a8d7037021433f9d678f3d9bfc4b69089acb64eead.
//
// Solidity: event NodeUpdate(bytes publicKey, string httpAddress, uint256 originatorId, bool isHealthy)
func (_Nodes *NodesFilterer) FilterNodeUpdate(opts *bind.FilterOpts) (*NodesNodeUpdateIterator, error) {

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "NodeUpdate")
	if err != nil {
		return nil, err
	}
	return &NodesNodeUpdateIterator{contract: _Nodes.contract, event: "NodeUpdate", logs: logs, sub: sub}, nil
}

// WatchNodeUpdate is a free log subscription operation binding the contract event 0x94b61d0a2ca1645b5632d2a8d7037021433f9d678f3d9bfc4b69089acb64eead.
//
// Solidity: event NodeUpdate(bytes publicKey, string httpAddress, uint256 originatorId, bool isHealthy)
func (_Nodes *NodesFilterer) WatchNodeUpdate(opts *bind.WatchOpts, sink chan<- *NodesNodeUpdate) (event.Subscription, error) {

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "NodeUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesNodeUpdate)
				if err := _Nodes.contract.UnpackLog(event, "NodeUpdate", log); err != nil {
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

// ParseNodeUpdate is a log parse operation binding the contract event 0x94b61d0a2ca1645b5632d2a8d7037021433f9d678f3d9bfc4b69089acb64eead.
//
// Solidity: event NodeUpdate(bytes publicKey, string httpAddress, uint256 originatorId, bool isHealthy)
func (_Nodes *NodesFilterer) ParseNodeUpdate(log types.Log) (*NodesNodeUpdate, error) {
	event := new(NodesNodeUpdate)
	if err := _Nodes.contract.UnpackLog(event, "NodeUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Nodes contract.
type NodesOwnershipTransferredIterator struct {
	Event *NodesOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *NodesOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesOwnershipTransferred)
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
		it.Event = new(NodesOwnershipTransferred)
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
func (it *NodesOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesOwnershipTransferred represents a OwnershipTransferred event raised by the Nodes contract.
type NodesOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Nodes *NodesFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*NodesOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &NodesOwnershipTransferredIterator{contract: _Nodes.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Nodes *NodesFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NodesOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesOwnershipTransferred)
				if err := _Nodes.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Nodes *NodesFilterer) ParseOwnershipTransferred(log types.Log) (*NodesOwnershipTransferred, error) {
	event := new(NodesOwnershipTransferred)
	if err := _Nodes.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
