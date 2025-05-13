package contract

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
)

// LogProcessor is a function that processes logs from a specific contract and event.
type LogProcessor func(ctx context.Context, logs []types.Log) error

// ReorgProcessor is a function that handles blockchain reorganizations for a specific contract.
type ReorgProcessor func(ctx context.Context, reorgedBlock uint64) error

// Contract provides a standard implementation of the Contract interface
type Contract struct {
	Name           string         // Name for identification
	ChainID        int64          // Chain ID of the network where this contract exists
	Address        string         // Contract address
	StartBlock     uint64         // Block to start indexing from (required)
	Topics         []string       // Event topics to filter
	Processor      LogProcessor   // Function to process logs
	ReorgProcessor ReorgProcessor // Function to handle blockchain reorganizations
}

// GetName returns the unique name of the contract
func (c *Contract) GetName() string {
	return c.Name
}

// GetChainID returns the blockchain network's chain ID
func (c *Contract) GetChainID() int64 {
	return c.ChainID
}

// GetAddress returns the contract's address on the blockchain
func (c *Contract) GetAddress() string {
	return c.Address
}

// GetStartBlock returns the block to start indexing from
func (c *Contract) GetStartBlock() uint64 {
	return c.StartBlock
}

// GetTopics returns the event topics to filter for
func (c *Contract) GetTopics() []string {
	return c.Topics
}

// ProcessLogs processes the logs found for this contract
func (c *Contract) ProcessLogs(ctx context.Context, logs []types.Log) error {
	return c.Processor(ctx, logs)
}

// HandleReorg handles a blockchain reorganization for this contract
func (c *Contract) HandleReorg(ctx context.Context, reorgedBlock uint64) error {
	return c.ReorgProcessor(ctx, reorgedBlock)
}

// NewContract creates a new standard contract
func NewContract(
	name string,
	chainID int64,
	address string,
	startBlock uint64,
	topics []string,
	processor LogProcessor,
	reorgProcessor ReorgProcessor,
) *Contract {
	return &Contract{
		Name:           name,
		ChainID:        chainID,
		Address:        address,
		StartBlock:     startBlock,
		Topics:         topics,
		Processor:      processor,
		ReorgProcessor: reorgProcessor,
	}
}
