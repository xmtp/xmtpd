package group_message_broadcaster

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/indexerpoc/contract"
	"go.uber.org/zap"
)

// GroupMessageContract represents a contract for broadcasting group messages
// and implements the indexerpoc.Contract interface
type GroupMessageContract struct {
	contract contract.Contract
	storer   *GroupMessageStorer
}

// NewGroupMessageContract creates a new contract for the group message broadcaster
func NewGroupMessageContract(
	logger *zap.Logger,
	querier *queries.Queries,
	name string,
	chainID int64,
	address string,
	startBlock uint64,
) *GroupMessageContract {
	indexingLogger := logger.Named("group-message-broadcaster").
		With(zap.String("address", address))

	storer := NewGroupMessageStorer(querier, indexingLogger)

	gm := GroupMessageContract{
		contract: contract.Contract{
			Name:       name,
			ChainID:    chainID,
			Address:    address,
			StartBlock: startBlock,
			Topics:     []string{"MessageSent"},
		},
		storer: storer,
	}

	gm.contract.Processor = gm.createGroupMessageProcessor()
	gm.contract.ReorgProcessor = gm.createGroupMessageReorgProcessor()

	return &gm
}

// GetContract returns the underlying contract object
func (c *GroupMessageContract) GetContract() *contract.Contract {
	return &c.contract
}

// Implementation of indexerpoc.Contract interface methods

// GetName returns the unique name of the contract
func (c *GroupMessageContract) GetName() string {
	return c.contract.GetName()
}

// GetChainID returns the blockchain network's chain ID
func (c *GroupMessageContract) GetChainID() int64 {
	return c.contract.GetChainID()
}

// GetAddress returns the contract's address on the blockchain
func (c *GroupMessageContract) GetAddress() string {
	return c.contract.GetAddress()
}

// GetStartBlock returns the block to start indexing from
func (c *GroupMessageContract) GetStartBlock() uint64 {
	return c.contract.GetStartBlock()
}

// GetTopics returns the event topics to filter for
func (c *GroupMessageContract) GetTopics() []string {
	return c.contract.GetTopics()
}

// ProcessLogs processes the logs found for this contract
func (c *GroupMessageContract) ProcessLogs(ctx context.Context, logs []types.Log) error {
	return c.contract.ProcessLogs(ctx, logs)
}

// HandleReorg handles a blockchain reorganization for this contract
func (c *GroupMessageContract) HandleReorg(ctx context.Context, fromBlock uint64) error {
	return c.contract.HandleReorg(ctx, fromBlock)
}

func (c *GroupMessageContract) createGroupMessageProcessor() contract.LogProcessor {
	return func(ctx context.Context, logs []types.Log) error {
		for _, log := range logs {
			err := c.storer.StoreLog(ctx, log)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func (c *GroupMessageContract) createGroupMessageReorgProcessor() contract.ReorgProcessor {
	return func(ctx context.Context, fromBlock uint64) error {
		fmt.Printf("Handling reorg for group messages from block %d\n", fromBlock)
		return nil
	}
}
