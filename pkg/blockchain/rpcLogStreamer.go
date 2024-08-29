package blockchain

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

const (
	BACKFILL_BLOCKS = 1000
	// Don't index very new blocks to account for reorgs
	// Setting to 0 since we are talking about L2s with low reorg risk
	LAG_FROM_HIGHEST_BLOCK = 0
	ERROR_SLEEP_TIME       = 100 * time.Millisecond
	NO_LOGS_SLEEP_TIME     = 1 * time.Second
)

// The builder that allows you to configure contract events to listen for
type RpcLogStreamBuilder struct {
	// All the listeners
	contractConfigs []contractConfig
	logger          *zap.Logger
	ethclient       *ethclient.Client
}

func NewRpcLogStreamBuilder(client *ethclient.Client, logger *zap.Logger) *RpcLogStreamBuilder {
	return &RpcLogStreamBuilder{ethclient: client, logger: logger}
}

func (c *RpcLogStreamBuilder) ListenForContractEvent(
	fromBlock int,
	contractAddress common.Address,
	topics []common.Hash,
) <-chan types.Log {
	eventChannel := make(chan types.Log, 100)
	c.contractConfigs = append(
		c.contractConfigs,
		contractConfig{fromBlock, contractAddress, topics, eventChannel},
	)
	return eventChannel
}

func (c *RpcLogStreamBuilder) Build() (*RpcLogStreamer, error) {
	return NewRpcLogStreamer(c.ethclient, c.logger, c.contractConfigs), nil
}

// Struct defining all the information required to filter events from logs
type contractConfig struct {
	fromBlock       int
	contractAddress common.Address
	topics          []common.Hash
	channel         chan<- types.Log
}

/*
*
A RpcLogStreamer is a naive implementation of the ChainStreamer interface.
It queries a remote blockchain node for log events to backfill history, and then streams new events,
to get a complete history of events on a chain.
*
*/
type RpcLogStreamer struct {
	client   ChainClient
	watchers []contractConfig
	ctx      context.Context
	logger   *zap.Logger
}

func NewRpcLogStreamer(
	client ChainClient,
	logger *zap.Logger,
	watchers []contractConfig,
) *RpcLogStreamer {
	return &RpcLogStreamer{
		client:   client,
		watchers: watchers,
		logger:   logger.Named("rpcLogStreamer"),
	}
}

func (r *RpcLogStreamer) Start(ctx context.Context) error {
	r.ctx = ctx

	for _, watcher := range r.watchers {
		go r.watchContract(watcher)
	}
	return nil
}

func (r *RpcLogStreamer) watchContract(watcher contractConfig) {
	fromBlock := int(watcher.fromBlock)
	logger := r.logger.With(zap.String("contractAddress", watcher.contractAddress.Hex()))
	defer close(watcher.channel)
	for {
		select {
		case <-r.ctx.Done():
			logger.Info("Stopping watcher")
			return
		default:
			logs, nextBlock, err := r.getNextPage(watcher, fromBlock)
			if err != nil {
				logger.Error(
					"Error getting next page",
					zap.Int("fromBlock", fromBlock),
					zap.Error(err),
				)
				time.Sleep(ERROR_SLEEP_TIME)
				continue
			}

			logger.Info("Got logs", zap.Int("numLogs", len(logs)), zap.Int("fromBlock", fromBlock))
			if len(logs) == 0 {
				time.Sleep(NO_LOGS_SLEEP_TIME)
			}
			for _, log := range logs {
				watcher.channel <- log
			}
			if nextBlock != nil {
				fromBlock = *nextBlock
			}
		}
	}
}

func (r *RpcLogStreamer) getNextPage(
	config contractConfig,
	fromBlock int,
) (logs []types.Log, nextBlock *int, err error) {
	highestBlock, err := r.client.BlockNumber(r.ctx)
	if err != nil {
		return nil, nil, err
	}

	highestBlockCanProcess := int(highestBlock) - LAG_FROM_HIGHEST_BLOCK
	numOfBlocksToProcess := highestBlockCanProcess - fromBlock + 1

	var to int
	// Make sure we stay within a reasonable page size
	if numOfBlocksToProcess > BACKFILL_BLOCKS {
		// quick mode
		to = fromBlock + BACKFILL_BLOCKS
	} else {
		// normal mode, up to current highest block num can process
		to = highestBlockCanProcess
	}

	// TODO:(nm) Use some more clever tactics to fetch the maximum number of logs at one times by parsing error messages
	// See: https://github.com/joshstevens19/rindexer/blob/master/core/src/indexer/fetch_logs.rs#L504
	logs, err = r.client.FilterLogs(r.ctx, buildFilterQuery(config, int64(fromBlock), int64(to)))
	if err != nil {
		return nil, nil, err
	}

	nextBlockNumber := to + 1

	return logs, &nextBlockNumber, nil
}

func buildFilterQuery(
	contractConfig contractConfig,
	fromBlock int64,
	toBlock int64,
) ethereum.FilterQuery {
	addresses := []common.Address{contractConfig.contractAddress}
	topics := [][]common.Hash{}
	for _, topic := range contractConfig.topics {
		topics = append(topics, []common.Hash{topic})
	}

	return ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(toBlock),
		Addresses: addresses,
		Topics:    topics,
	}
}
