package blockchain

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/tracing"

	"github.com/ethereum/go-ethereum/core/types"
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
	reader          AppChainReader
}

func NewRpcLogStreamBuilder(
	ctx context.Context,
	reader AppChainReader,
	logger *zap.Logger,
) *RpcLogStreamBuilder {
	return &RpcLogStreamBuilder{ctx: ctx, reader: reader, logger: logger}
}

func (c *RpcLogStreamBuilder) ListenForContractEvent(
	eventType EventType,
	fromBlock uint64,
	maxDisconnectTime time.Duration,
) (<-chan types.Log, chan<- uint64) {
	eventChannel := make(chan types.Log, 100)
	reorgChannel := make(chan uint64, 1)
	c.contractConfigs = append(
		c.contractConfigs,
		ContractConfig{
			eventType,
			fromBlock,
			eventChannel,
			reorgChannel,
			maxDisconnectTime,
		},
	)
	return eventChannel, reorgChannel
}

func (c *RpcLogStreamBuilder) Build() (*RpcLogStreamer, error) {
	return NewRpcLogStreamer(c.ctx, c.reader, c.logger, c.contractConfigs), nil
}

// Struct defining all the information required to filter events from logs
type ContractConfig struct {
	EventType         EventType
	FromBlock         uint64
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
	reader   AppChainReader
	watchers []ContractConfig
	ctx      context.Context
	logger   *zap.Logger
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func NewRpcLogStreamer(
	ctx context.Context,
	reader AppChainReader,
	logger *zap.Logger,
	watchers []ContractConfig,
) *RpcLogStreamer {
	ctx, cancel := context.WithCancel(ctx)
	return &RpcLogStreamer{
		ctx:      ctx,
		reader:   reader,
		watchers: watchers,
		logger:   logger.Named("rpcLogStreamer"),
		cancel:   cancel,
		wg:       sync.WaitGroup{},
	}
}

func (r *RpcLogStreamer) Start() {
	for _, watcher := range r.watchers {
		contractAddress, err := r.reader.ContractAddress(watcher.EventType)
		if err != nil {
			r.logger.Error("Error getting contract address", zap.Error(err))
			continue
		}
		tracing.GoPanicWrap(
			r.ctx,
			&r.wg,
			fmt.Sprintf("rpcLogStreamer-watcher-%v", contractAddress),
			func(ctx context.Context) {
				r.watchContract(watcher, contractAddress)
			})
	}
}

func (r *RpcLogStreamer) watchContract(watcher ContractConfig, contractAddress string) {
	fromBlock := watcher.FromBlock
	logger := r.logger.With(zap.String("contractAddress", contractAddress))
	defer close(watcher.EventChannel)

	timer := time.NewTimer(watcher.maxDisconnectTime)
	defer timer.Stop()

	for {
		select {
		case <-r.ctx.Done():
			logger.Debug("Stopping watcher")
			return
		case <-timer.C:
			logger.Fatal(
				"Max disconnect time exceeded. Node might drift too far away from expected state. Shutting down...",
			)

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
			timer.Reset(watcher.maxDisconnectTime)

		default:
			logs, nextBlock, err := r.GetNextPage(watcher, fromBlock)
			if err != nil {
				logger.Error(
					"Error getting next page",
					zap.Uint64("fromBlock", fromBlock),
					zap.Error(err),
				)
				time.Sleep(ERROR_SLEEP_TIME)
				continue
			}
			// reset self-termination timer
			timer.Reset(watcher.maxDisconnectTime)

			if nextBlock != nil {
				fromBlock = *nextBlock
			}

			if len(logs) == 0 {
				time.Sleep(NO_LOGS_SLEEP_TIME)
				continue
			}

			logger.Debug(
				"Got logs",
				zap.Int("numLogs", len(logs)),
				zap.Uint64("fromBlock", fromBlock),
				zap.Time("time", time.Now()),
			)
			for _, log := range logs {
				watcher.EventChannel <- log
			}
		}
	}
}

func (r *RpcLogStreamer) GetNextPage(
	config ContractConfig,
	fromBlock uint64,
) (logs []types.Log, nextBlock *uint64, err error) {
	contractAddress, err := r.reader.ContractAddress(config.EventType)
	if err != nil {
		return nil, nil, err
	}
	highestBlock, err := r.reader.BlockNumber(r.ctx)
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
		return r.reader.FilterLogs(
			r.ctx,
			config.EventType,
			fromBlock,
			toBlock,
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

func (r *RpcLogStreamer) Reader() AppChainReader {
	return r.reader
}

func (r *RpcLogStreamer) Stop() {
	r.cancel()
	r.wg.Wait()
}
