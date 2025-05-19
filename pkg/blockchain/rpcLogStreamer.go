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
	"github.com/xmtp/xmtpd/pkg/metrics"
	"go.uber.org/zap"
)

const (
	BACKFILL_BLOCKS    = uint64(1000)
	ERROR_SLEEP_TIME   = 100 * time.Millisecond
	NO_LOGS_SLEEP_TIME = 100 * time.Millisecond
)

// Struct defining all the information required to filter events from logs
type ContractConfig struct {
	ID                string
	FromBlock         uint64
	ContractAddress   common.Address
	Topics            []common.Hash
	MaxDisconnectTime time.Duration
	backfillChannel   chan types.Log
	reorgChannel      chan uint64
}

type RpcLogStreamerOption func(*RpcLogStreamer) error

func WithLagFromHighestBlock(lagFromHighestBlock uint8) RpcLogStreamerOption {
	return func(streamer *RpcLogStreamer) error {
		streamer.lagFromHighestBlock = lagFromHighestBlock
		return nil
	}
}

func WithContractConfig(
	cfg ContractConfig,
) RpcLogStreamerOption {
	return func(streamer *RpcLogStreamer) error {
		if _, ok := streamer.watchers[cfg.ID]; ok {
			streamer.logger.Error("contract config already exists", zap.String("id", cfg.ID))
			return fmt.Errorf("contract config already exists: %s", cfg.ID)
		}

		backfillChannel := make(chan types.Log, 100)
		reorgChannel := make(chan uint64, 1)

		streamer.watchers[cfg.ID] = ContractConfig{
			ID:                cfg.ID,
			FromBlock:         cfg.FromBlock,
			ContractAddress:   cfg.ContractAddress,
			Topics:            cfg.Topics,
			MaxDisconnectTime: cfg.MaxDisconnectTime,
			backfillChannel:   backfillChannel,
			reorgChannel:      reorgChannel,
		}

		return nil
	}
}

/*
*
A RpcLogStreamer is a naive implementation of the ChainStreamer interface.
It queries a remote blockchain node for log events to backfill history, and then streams new events,
to get a complete history of events on a chain.
*
*/
type RpcLogStreamer struct {
	ctx                 context.Context
	cancel              context.CancelFunc
	wg                  sync.WaitGroup
	client              ChainClient
	logger              *zap.Logger
	watchers            map[string]ContractConfig
	lagFromHighestBlock uint8
}

func NewRpcLogStreamer(
	ctx context.Context,
	client ChainClient,
	logger *zap.Logger,
	chainID int,
	options ...RpcLogStreamerOption,
) (*RpcLogStreamer, error) {
	ctx, cancel := context.WithCancel(ctx)

	streamLogger := logger.Named("rpcLogStreamer").
		With(zap.Int("chainID", chainID))

	streamer := &RpcLogStreamer{
		ctx:                 ctx,
		client:              client,
		logger:              streamLogger,
		cancel:              cancel,
		wg:                  sync.WaitGroup{},
		watchers:            make(map[string]ContractConfig),
		lagFromHighestBlock: 0,
	}

	for _, option := range options {
		if err := option(streamer); err != nil {
			streamLogger.Error("failed to apply option", zap.Error(err))
			return nil, err
		}
	}

	return streamer, nil
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

func (r *RpcLogStreamer) Stop() {
	r.cancel()
	r.wg.Wait()
}

func (r *RpcLogStreamer) watchContract(watcher ContractConfig) {
	fromBlock := watcher.FromBlock
	logger := r.logger.With(zap.String("contractAddress", watcher.ContractAddress.Hex()))
	defer close(watcher.backfillChannel)
	defer close(watcher.reorgChannel)

	timer := time.NewTimer(watcher.MaxDisconnectTime)
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
			timer.Reset(watcher.MaxDisconnectTime)

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
			timer.Reset(watcher.MaxDisconnectTime)

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
				watcher.backfillChannel <- log
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

	highestBlockCanProcess := highestBlock - uint64(r.lagFromHighestBlock)
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
			buildFilterQuery(config, fromBlock, toBlock),
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

func (r *RpcLogStreamer) GetEventChannel(id string) chan types.Log {
	if _, ok := r.watchers[id]; !ok {
		return nil
	}

	return r.watchers[id].backfillChannel
}

func (r *RpcLogStreamer) GetReorgChannel(id string) chan uint64 {
	if _, ok := r.watchers[id]; !ok {
		return nil
	}

	return r.watchers[id].reorgChannel
}

func buildFilterQuery(
	contractConfig ContractConfig,
	fromBlock uint64,
	toBlock uint64,
) ethereum.FilterQuery {
	addresses := []common.Address{contractConfig.ContractAddress}
	topics := [][]common.Hash{}
	for _, topic := range contractConfig.Topics {
		topics = append(topics, []common.Hash{topic})
	}

	return ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		ToBlock:   new(big.Int).SetUint64(toBlock),
		Addresses: addresses,
		Topics:    topics,
	}
}
