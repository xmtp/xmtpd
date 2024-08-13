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

// GroupMessagesMetaData contains all meta data concerning the GroupMessages contract.
var GroupMessagesMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"addMessage\",\"inputs\":[{\"name\":\"groupId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"MessageSent\",\"inputs\":[{\"name\":\"groupId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"message\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"sequenceId\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false}]",
}

// GroupMessagesABI is the input ABI used to generate the binding from.
// Deprecated: Use GroupMessagesMetaData.ABI instead.
var GroupMessagesABI = GroupMessagesMetaData.ABI

// GroupMessages is an auto generated Go binding around an Ethereum contract.
type GroupMessages struct {
	GroupMessagesCaller     // Read-only binding to the contract
	GroupMessagesTransactor // Write-only binding to the contract
	GroupMessagesFilterer   // Log filterer for contract events
}

// GroupMessagesCaller is an auto generated read-only Go binding around an Ethereum contract.
type GroupMessagesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GroupMessagesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GroupMessagesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GroupMessagesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GroupMessagesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GroupMessagesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GroupMessagesSession struct {
	Contract     *GroupMessages    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GroupMessagesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GroupMessagesCallerSession struct {
	Contract *GroupMessagesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// GroupMessagesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GroupMessagesTransactorSession struct {
	Contract     *GroupMessagesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// GroupMessagesRaw is an auto generated low-level Go binding around an Ethereum contract.
type GroupMessagesRaw struct {
	Contract *GroupMessages // Generic contract binding to access the raw methods on
}

// GroupMessagesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GroupMessagesCallerRaw struct {
	Contract *GroupMessagesCaller // Generic read-only contract binding to access the raw methods on
}

// GroupMessagesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GroupMessagesTransactorRaw struct {
	Contract *GroupMessagesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGroupMessages creates a new instance of GroupMessages, bound to a specific deployed contract.
func NewGroupMessages(address common.Address, backend bind.ContractBackend) (*GroupMessages, error) {
	contract, err := bindGroupMessages(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GroupMessages{GroupMessagesCaller: GroupMessagesCaller{contract: contract}, GroupMessagesTransactor: GroupMessagesTransactor{contract: contract}, GroupMessagesFilterer: GroupMessagesFilterer{contract: contract}}, nil
}

// NewGroupMessagesCaller creates a new read-only instance of GroupMessages, bound to a specific deployed contract.
func NewGroupMessagesCaller(address common.Address, caller bind.ContractCaller) (*GroupMessagesCaller, error) {
	contract, err := bindGroupMessages(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GroupMessagesCaller{contract: contract}, nil
}

// NewGroupMessagesTransactor creates a new write-only instance of GroupMessages, bound to a specific deployed contract.
func NewGroupMessagesTransactor(address common.Address, transactor bind.ContractTransactor) (*GroupMessagesTransactor, error) {
	contract, err := bindGroupMessages(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GroupMessagesTransactor{contract: contract}, nil
}

// NewGroupMessagesFilterer creates a new log filterer instance of GroupMessages, bound to a specific deployed contract.
func NewGroupMessagesFilterer(address common.Address, filterer bind.ContractFilterer) (*GroupMessagesFilterer, error) {
	contract, err := bindGroupMessages(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GroupMessagesFilterer{contract: contract}, nil
}

// bindGroupMessages binds a generic wrapper to an already deployed contract.
func bindGroupMessages(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GroupMessagesMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GroupMessages *GroupMessagesRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GroupMessages.Contract.GroupMessagesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GroupMessages *GroupMessagesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GroupMessages.Contract.GroupMessagesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GroupMessages *GroupMessagesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GroupMessages.Contract.GroupMessagesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GroupMessages *GroupMessagesCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GroupMessages.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GroupMessages *GroupMessagesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GroupMessages.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GroupMessages *GroupMessagesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GroupMessages.Contract.contract.Transact(opts, method, params...)
}

// AddMessage is a paid mutator transaction binding the contract method 0x4dff26b5.
//
// Solidity: function addMessage(bytes32 groupId, bytes message) returns()
func (_GroupMessages *GroupMessagesTransactor) AddMessage(opts *bind.TransactOpts, groupId [32]byte, message []byte) (*types.Transaction, error) {
	return _GroupMessages.contract.Transact(opts, "addMessage", groupId, message)
}

// AddMessage is a paid mutator transaction binding the contract method 0x4dff26b5.
//
// Solidity: function addMessage(bytes32 groupId, bytes message) returns()
func (_GroupMessages *GroupMessagesSession) AddMessage(groupId [32]byte, message []byte) (*types.Transaction, error) {
	return _GroupMessages.Contract.AddMessage(&_GroupMessages.TransactOpts, groupId, message)
}

// AddMessage is a paid mutator transaction binding the contract method 0x4dff26b5.
//
// Solidity: function addMessage(bytes32 groupId, bytes message) returns()
func (_GroupMessages *GroupMessagesTransactorSession) AddMessage(groupId [32]byte, message []byte) (*types.Transaction, error) {
	return _GroupMessages.Contract.AddMessage(&_GroupMessages.TransactOpts, groupId, message)
}

// GroupMessagesMessageSentIterator is returned from FilterMessageSent and is used to iterate over the raw logs and unpacked data for MessageSent events raised by the GroupMessages contract.
type GroupMessagesMessageSentIterator struct {
	Event *GroupMessagesMessageSent // Event containing the contract specifics and raw log

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
func (it *GroupMessagesMessageSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessagesMessageSent)
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
		it.Event = new(GroupMessagesMessageSent)
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
func (it *GroupMessagesMessageSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessagesMessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessagesMessageSent represents a MessageSent event raised by the GroupMessages contract.
type GroupMessagesMessageSent struct {
	GroupId    [32]byte
	Message    []byte
	SequenceId uint64
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterMessageSent is a free log retrieval operation binding the contract event 0x91f47151424884a46811ed593aa8a02ee5012e9332a4dcf1e9236a8ed4443c3e.
//
// Solidity: event MessageSent(bytes32 groupId, bytes message, uint64 sequenceId)
func (_GroupMessages *GroupMessagesFilterer) FilterMessageSent(opts *bind.FilterOpts) (*GroupMessagesMessageSentIterator, error) {

	logs, sub, err := _GroupMessages.contract.FilterLogs(opts, "MessageSent")
	if err != nil {
		return nil, err
	}
	return &GroupMessagesMessageSentIterator{contract: _GroupMessages.contract, event: "MessageSent", logs: logs, sub: sub}, nil
}

// WatchMessageSent is a free log subscription operation binding the contract event 0x91f47151424884a46811ed593aa8a02ee5012e9332a4dcf1e9236a8ed4443c3e.
//
// Solidity: event MessageSent(bytes32 groupId, bytes message, uint64 sequenceId)
func (_GroupMessages *GroupMessagesFilterer) WatchMessageSent(opts *bind.WatchOpts, sink chan<- *GroupMessagesMessageSent) (event.Subscription, error) {

	logs, sub, err := _GroupMessages.contract.WatchLogs(opts, "MessageSent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessagesMessageSent)
				if err := _GroupMessages.contract.UnpackLog(event, "MessageSent", log); err != nil {
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
// Solidity: event MessageSent(bytes32 groupId, bytes message, uint64 sequenceId)
func (_GroupMessages *GroupMessagesFilterer) ParseMessageSent(log types.Log) (*GroupMessagesMessageSent, error) {
	event := new(GroupMessagesMessageSent)
	if err := _GroupMessages.contract.UnpackLog(event, "MessageSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
