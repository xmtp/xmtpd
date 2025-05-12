package indexerpoc

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	ErrReorg      = fmt.Errorf("blockchain reorganization detected")
	ErrNothingNew = fmt.Errorf("no new blocks available")
)

// LogProcessor is a function that processes logs from a specific contract and event.
type LogProcessor func(ctx context.Context, logs []types.Log) error

// ReorgProcessor is a function that handles blockchain reorganizations for a specific contract.
type ReorgProcessor func(ctx context.Context, fromBlock uint64) error

// Contract defines the interface for contracts to be indexed.
type Contract interface {
	GetName() string
	GetChainID() int64
	GetAddress() string
	GetStartBlock() uint64
	GetTopics() []string
	ProcessLogs(ctx context.Context, logs []types.Log) error
	HandleReorg(ctx context.Context, fromBlock uint64) error
}

// Source defines the interface a blockchain.
type Source interface {
	GetLogs(ctx context.Context, startBlock, endBlock uint64, filter *Filter) ([]types.Log, error)
	GetBlockByNumber(ctx context.Context, number uint64) (*types.Block, error)
	GetLatestBlockNumber(ctx context.Context) (uint64, error)
	GetBlockHash(ctx context.Context, number uint64) (common.Hash, error)
	GetNetworkName() string
	GetChainID() int64
}

// Temporary: Storage interface for persisting state and logs.
type Storage interface {
	SaveState(ctx context.Context, state *taskState) error
	GetState(ctx context.Context, contractName string, network string) (*taskState, error)
	DeleteFromBlock(
		ctx context.Context,
		contractName string,
		network string,
		blockNumber uint64,
	) error
	Begin(ctx context.Context) (Transaction, error)
}

// Temporary: Transaction represents a database transaction.
type Transaction interface {
	Commit() error
	Rollback() error
}
