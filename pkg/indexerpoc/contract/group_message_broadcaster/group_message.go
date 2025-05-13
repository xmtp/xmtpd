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
	contract *contract.Contract
	storer   *GroupMessageStorer
	logger   *zap.Logger
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
	l := logger.Named("group-message-broadcaster").
		With(zap.String("address", address))

	storer := NewGroupMessageStorer(querier, l)

	gm := GroupMessageContract{
		contract: &contract.Contract{
			Name:       name,
			ChainID:    chainID,
			Address:    address,
			StartBlock: startBlock,
			Topics:     []string{"MessageSent"},
		},
		storer: storer,
		logger: l,
	}

	return &gm
}

func (c *GroupMessageContract) GetContract() *contract.Contract {
	return c.contract
}

func (c *GroupMessageContract) GetName() string {
	return c.contract.GetName()
}

func (c *GroupMessageContract) GetChainID() int64 {
	return c.contract.GetChainID()
}

func (c *GroupMessageContract) GetAddress() string {
	return c.contract.GetAddress()
}

func (c *GroupMessageContract) GetStartBlock() uint64 {
	return c.contract.GetStartBlock()
}

func (c *GroupMessageContract) GetTopics() []string {
	return c.contract.GetTopics()
}

func (c *GroupMessageContract) ProcessLogs(ctx context.Context, logs []types.Log) error {
	// TODO: This requires more logic.
	for _, log := range logs {
		err := c.storer.StoreLog(ctx, log)
		if err != nil {
			return err
		}
	}

	return nil
}

// HandleReorg handles a blockchain reorganization for this contract
func (c *GroupMessageContract) HandleReorg(ctx context.Context, reorgedBlock uint64) error {
	// TODO: This requires more logic.
	fmt.Printf("Handling reorg for group messages from block %d\n", reorgedBlock)
	return nil
}
