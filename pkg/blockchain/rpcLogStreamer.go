package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/tracing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"go.uber.org/zap"
)

const (
	BACKFILL_BLOCKS = uint64(1000)
	// Don't index very new blocks to account for reorgs
	// Setting to 0 since we are talking about L2s with low reorg risk
	LAG_FROM_HIGHEST_BLOCK = uint64(0)
	ERROR_SLEEP_TIME       = 100 * time.Millisecond
	NO_LOGS_SLEEP_TIME     = 1 * time.Second
)

// The builder that allows you to configure contract events to listen for
type RpcLogStreamBuilder struct {
	// All the listeners
	ctx             context.Context
	contractConfigs []ContractConfig
	logger          *zap.Logger
	ethclient       *ethclient.Client
}

func NewRpcLogStreamBuilder(
	ctx context.Context,
	client *ethclient.Client,
	logger *zap.Logger,
) *RpcLogStreamBuilder {
	return &RpcLogStreamBuilder{ctx: ctx, ethclient: client, logger: logger}
}

func (c *RpcLogStreamBuilder) ListenForContractEvent(
	fromBlock uint64,
	contractAddress common.Address,
	topics []common.Hash,
	maxDisconnectTime time.Duration,
) (<-chan types.Log, chan<- uint64) {
	eventChannel := make(chan types.Log, 100)
	reorgChannel := make(chan uint64, 1)
	c.contractConfigs = append(
		c.contractConfigs,
		ContractConfig{
			fromBlock,
			contractAddress,
			topics,
			eventChannel,
			reorgChannel,
			maxDisconnectTime,
		},
	)
	return eventChannel, reorgChannel
}

func (c *RpcLogStreamBuilder) Build() (*RpcLogStreamer, error) {
	return NewRpcLogStreamer(c.ctx, c.ethclient, c.logger, c.contractConfigs), nil
}

// Struct defining all the information required to filter events from logs
type ContractConfig struct {
	FromBlock         uint64
	ContractAddress   common.Address
	Topics            []common.Hash
	EventChannel      chan<- types.Log
	reorgChannel      chan uint64
	maxDisconnectTime time.Duration
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
	watchers []ContractConfig
	ctx      context.Context
	logger   *zap.Logger
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func NewRpcLogStreamer(
	ctx context.Context,
	client ChainClient,
	logger *zap.Logger,
	watchers []ContractConfig,
) *RpcLogStreamer {
	ctx, cancel := context.WithCancel(ctx)
	return &RpcLogStreamer{
		ctx:      ctx,
		client:   client,
		watchers: watchers,
		logger:   logger.Named("rpcLogStreamer"),
		cancel:   cancel,
		wg:       sync.WaitGroup{},
	}
}

func (r *RpcLogStreamer) Start() {
	for _, watcher := range r.watchers {
		tracing.GoPanicWrap(
			r.ctx,
			&r.wg,
			fmt.Sprintf("rpcLogStreamer-watcher-%v", watcher.ContractAddress),
			func(ctx context.Context) {
				r.watchContract(watcher)
			})
	}
}

func (r *RpcLogStreamer) watchContract(watcher ContractConfig) {
	fromBlock := watcher.FromBlock
	logger := r.logger.With(zap.String("contractAddress", watcher.ContractAddress.Hex()))
	startTime := time.Now()
	defer close(watcher.EventChannel)

	for {
		select {
		case <-r.ctx.Done():
			logger.Debug("Stopping watcher")
			return
		case reorgBlock, open := <-watcher.reorgChannel:
			if !open {
				logger.Debug("Reorg channel closed")
				return
			}
			fromBlock = reorgBlock
			logger.Warn(
				"Blockchain reorg detected, resuming from block",
				zap.Uint64("fromBlock", fromBlock),
			)
		default:
			logs, nextBlock, err := r.GetNextPage(watcher, fromBlock)
			if err != nil {
				logger.Error(
					"Error getting next page",
					zap.Uint64("fromBlock", fromBlock),
					zap.Error(err),
				)
				time.Sleep(ERROR_SLEEP_TIME)

				if time.Since(startTime) > watcher.maxDisconnectTime {
					logger.Error(
						"Max disconnect time exceeded. Node might drift too far away from expected state. Shutting down...",
					)
					panic(
						"Max disconnect time exceeded. Node might drift too far away from expected state",
					)
				}
				continue
			}
			// reset self-termination timer
			startTime = time.Now()

			if len(logs) > 0 {
				logger.Debug(
					"Got logs",
					zap.Int("numLogs", len(logs)),
					zap.Uint64("fromBlock", fromBlock),
				)
			}

			if len(logs) == 0 {
				time.Sleep(NO_LOGS_SLEEP_TIME)
			}
			for _, log := range logs {
				watcher.EventChannel <- log
			}
			if nextBlock != nil {
				fromBlock = *nextBlock
			}
		}
	}
}

func (r *RpcLogStreamer) GetNextPage(
	config ContractConfig,
	fromBlock uint64,
) (logs []types.Log, nextBlock *uint64, err error) {
	contractAddress := config.ContractAddress.Hex()
	highestBlock, err := r.client.BlockNumber(r.ctx)
	if err != nil {
		return nil, nil, err
	}
	metrics.EmitIndexerMaxBlock(contractAddress, highestBlock)

	highestBlockCanProcess := highestBlock - LAG_FROM_HIGHEST_BLOCK
	if fromBlock > highestBlockCanProcess {
		metrics.EmitIndexerCurrentBlockLag(contractAddress, 0)
		return []types.Log{}, nil, nil
	}

	metrics.EmitIndexerCurrentBlockLag(contractAddress, highestBlock-fromBlock)

	toBlock := min(fromBlock+BACKFILL_BLOCKS, highestBlockCanProcess)

	// TODO:(nm) Use some more clever tactics to fetch the maximum number of logs at one times by parsing error messages
	// See: https://github.com/joshstevens19/rindexer/blob/master/core/src/indexer/fetch_logs.rs#L504
	logs, err = metrics.MeasureGetLogs(contractAddress, func() ([]types.Log, error) {
		return r.client.FilterLogs(
			r.ctx,
			buildFilterQuery(config, int64(fromBlock), int64(toBlock)),
		)
	})

	if err != nil {
		return nil, nil, err
	}

	metrics.EmitIndexerCurrentBlock(contractAddress, toBlock)
	metrics.EmitIndexerNumLogsFound(contractAddress, len(logs))

	nextBlockNumber := toBlock + 1

	return logs, &nextBlockNumber, nil
}

func (r *RpcLogStreamer) Client() ChainClient {
	return r.client
}

func buildFilterQuery(
	contractConfig ContractConfig,
	fromBlock int64,
	toBlock int64,
) ethereum.FilterQuery {
	addresses := []common.Address{contractConfig.ContractAddress}
	topics := [][]common.Hash{}
	for _, topic := range contractConfig.Topics {
		topics = append(topics, []common.Hash{topic})
	}

	return ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(toBlock),
		Addresses: addresses,
		Topics:    topics,
	}
}
func (r *RpcLogStreamer) Stop() {
	r.cancel()
	r.wg.Wait()
}
