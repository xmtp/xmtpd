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

// IdentityUpdatesMetaData contains all meta data concerning the IdentityUpdates contract.
var IdentityUpdatesMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"addIdentityUpdate\",\"inputs\":[{\"name\":\"inboxId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"update\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"IdentityUpdateCreated\",\"inputs\":[{\"name\":\"inboxId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"update\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"sequenceId\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false}]",
}

// IdentityUpdatesABI is the input ABI used to generate the binding from.
// Deprecated: Use IdentityUpdatesMetaData.ABI instead.
var IdentityUpdatesABI = IdentityUpdatesMetaData.ABI

// IdentityUpdates is an auto generated Go binding around an Ethereum contract.
type IdentityUpdates struct {
	IdentityUpdatesCaller     // Read-only binding to the contract
	IdentityUpdatesTransactor // Write-only binding to the contract
	IdentityUpdatesFilterer   // Log filterer for contract events
}

// IdentityUpdatesCaller is an auto generated read-only Go binding around an Ethereum contract.
type IdentityUpdatesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityUpdatesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IdentityUpdatesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityUpdatesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IdentityUpdatesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityUpdatesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IdentityUpdatesSession struct {
	Contract     *IdentityUpdates  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IdentityUpdatesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IdentityUpdatesCallerSession struct {
	Contract *IdentityUpdatesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// IdentityUpdatesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IdentityUpdatesTransactorSession struct {
	Contract     *IdentityUpdatesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// IdentityUpdatesRaw is an auto generated low-level Go binding around an Ethereum contract.
type IdentityUpdatesRaw struct {
	Contract *IdentityUpdates // Generic contract binding to access the raw methods on
}

// IdentityUpdatesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IdentityUpdatesCallerRaw struct {
	Contract *IdentityUpdatesCaller // Generic read-only contract binding to access the raw methods on
}

// IdentityUpdatesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IdentityUpdatesTransactorRaw struct {
	Contract *IdentityUpdatesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIdentityUpdates creates a new instance of IdentityUpdates, bound to a specific deployed contract.
func NewIdentityUpdates(address common.Address, backend bind.ContractBackend) (*IdentityUpdates, error) {
	contract, err := bindIdentityUpdates(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdates{IdentityUpdatesCaller: IdentityUpdatesCaller{contract: contract}, IdentityUpdatesTransactor: IdentityUpdatesTransactor{contract: contract}, IdentityUpdatesFilterer: IdentityUpdatesFilterer{contract: contract}}, nil
}

// NewIdentityUpdatesCaller creates a new read-only instance of IdentityUpdates, bound to a specific deployed contract.
func NewIdentityUpdatesCaller(address common.Address, caller bind.ContractCaller) (*IdentityUpdatesCaller, error) {
	contract, err := bindIdentityUpdates(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdatesCaller{contract: contract}, nil
}

// NewIdentityUpdatesTransactor creates a new write-only instance of IdentityUpdates, bound to a specific deployed contract.
func NewIdentityUpdatesTransactor(address common.Address, transactor bind.ContractTransactor) (*IdentityUpdatesTransactor, error) {
	contract, err := bindIdentityUpdates(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdatesTransactor{contract: contract}, nil
}

// NewIdentityUpdatesFilterer creates a new log filterer instance of IdentityUpdates, bound to a specific deployed contract.
func NewIdentityUpdatesFilterer(address common.Address, filterer bind.ContractFilterer) (*IdentityUpdatesFilterer, error) {
	contract, err := bindIdentityUpdates(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdatesFilterer{contract: contract}, nil
}

// bindIdentityUpdates binds a generic wrapper to an already deployed contract.
func bindIdentityUpdates(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IdentityUpdatesMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IdentityUpdates *IdentityUpdatesRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IdentityUpdates.Contract.IdentityUpdatesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IdentityUpdates *IdentityUpdatesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.IdentityUpdatesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IdentityUpdates *IdentityUpdatesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.IdentityUpdatesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IdentityUpdates *IdentityUpdatesCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IdentityUpdates.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IdentityUpdates *IdentityUpdatesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IdentityUpdates *IdentityUpdatesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.contract.Transact(opts, method, params...)
}

// AddIdentityUpdate is a paid mutator transaction binding the contract method 0xba74fc7c.
//
// Solidity: function addIdentityUpdate(bytes32 inboxId, bytes update) returns()
func (_IdentityUpdates *IdentityUpdatesTransactor) AddIdentityUpdate(opts *bind.TransactOpts, inboxId [32]byte, update []byte) (*types.Transaction, error) {
	return _IdentityUpdates.contract.Transact(opts, "addIdentityUpdate", inboxId, update)
}

// AddIdentityUpdate is a paid mutator transaction binding the contract method 0xba74fc7c.
//
// Solidity: function addIdentityUpdate(bytes32 inboxId, bytes update) returns()
func (_IdentityUpdates *IdentityUpdatesSession) AddIdentityUpdate(inboxId [32]byte, update []byte) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.AddIdentityUpdate(&_IdentityUpdates.TransactOpts, inboxId, update)
}

// AddIdentityUpdate is a paid mutator transaction binding the contract method 0xba74fc7c.
//
// Solidity: function addIdentityUpdate(bytes32 inboxId, bytes update) returns()
func (_IdentityUpdates *IdentityUpdatesTransactorSession) AddIdentityUpdate(inboxId [32]byte, update []byte) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.AddIdentityUpdate(&_IdentityUpdates.TransactOpts, inboxId, update)
}

// IdentityUpdatesIdentityUpdateCreatedIterator is returned from FilterIdentityUpdateCreated and is used to iterate over the raw logs and unpacked data for IdentityUpdateCreated events raised by the IdentityUpdates contract.
type IdentityUpdatesIdentityUpdateCreatedIterator struct {
	Event *IdentityUpdatesIdentityUpdateCreated // Event containing the contract specifics and raw log

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
func (it *IdentityUpdatesIdentityUpdateCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityUpdatesIdentityUpdateCreated)
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
		it.Event = new(IdentityUpdatesIdentityUpdateCreated)
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
func (it *IdentityUpdatesIdentityUpdateCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityUpdatesIdentityUpdateCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityUpdatesIdentityUpdateCreated represents a IdentityUpdateCreated event raised by the IdentityUpdates contract.
type IdentityUpdatesIdentityUpdateCreated struct {
	InboxId    [32]byte
	Update     []byte
	SequenceId uint64
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterIdentityUpdateCreated is a free log retrieval operation binding the contract event 0xc1a40f292090ec0435e939cdfe248e0322a88566679a90a50c4e9e5ef762dbd5.
//
// Solidity: event IdentityUpdateCreated(bytes32 inboxId, bytes update, uint64 sequenceId)
func (_IdentityUpdates *IdentityUpdatesFilterer) FilterIdentityUpdateCreated(opts *bind.FilterOpts) (*IdentityUpdatesIdentityUpdateCreatedIterator, error) {

	logs, sub, err := _IdentityUpdates.contract.FilterLogs(opts, "IdentityUpdateCreated")
	if err != nil {
		return nil, err
	}
	return &IdentityUpdatesIdentityUpdateCreatedIterator{contract: _IdentityUpdates.contract, event: "IdentityUpdateCreated", logs: logs, sub: sub}, nil
}

// WatchIdentityUpdateCreated is a free log subscription operation binding the contract event 0xc1a40f292090ec0435e939cdfe248e0322a88566679a90a50c4e9e5ef762dbd5.
//
// Solidity: event IdentityUpdateCreated(bytes32 inboxId, bytes update, uint64 sequenceId)
func (_IdentityUpdates *IdentityUpdatesFilterer) WatchIdentityUpdateCreated(opts *bind.WatchOpts, sink chan<- *IdentityUpdatesIdentityUpdateCreated) (event.Subscription, error) {

	logs, sub, err := _IdentityUpdates.contract.WatchLogs(opts, "IdentityUpdateCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityUpdatesIdentityUpdateCreated)
				if err := _IdentityUpdates.contract.UnpackLog(event, "IdentityUpdateCreated", log); err != nil {
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

// ParseIdentityUpdateCreated is a log parse operation binding the contract event 0xc1a40f292090ec0435e939cdfe248e0322a88566679a90a50c4e9e5ef762dbd5.
//
// Solidity: event IdentityUpdateCreated(bytes32 inboxId, bytes update, uint64 sequenceId)
func (_IdentityUpdates *IdentityUpdatesFilterer) ParseIdentityUpdateCreated(log types.Log) (*IdentityUpdatesIdentityUpdateCreated, error) {
	event := new(IdentityUpdatesIdentityUpdateCreated)
	if err := _IdentityUpdates.contract.UnpackLog(event, "IdentityUpdateCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
