package blockchain

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/xmtp/xmtpd/pkg/tracing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"go.uber.org/zap"
)

const (
	backfillBlocks   = uint64(1000)
	maxSubRetries    = 10
	sleepTimeOnError = 100 * time.Millisecond
	sleepTimeNoLogs  = 100 * time.Millisecond
)

var ErrReorg = errors.New("blockchain reorg detected")

// Struct defining all the information required to filter events from logs
type ContractConfig struct {
	ID                string
	FromBlockNumber   uint64
	FromBlockHash     []byte
	Address           common.Address
	Topics            []common.Hash
	MaxDisconnectTime time.Duration
	eventChannel      chan types.Log
}

type RpcLogStreamerOption func(*RpcLogStreamer) error

func WithLagFromHighestBlock(lagFromHighestBlock uint8) RpcLogStreamerOption {
	return func(streamer *RpcLogStreamer) error {
		streamer.lagFromHighestBlock = lagFromHighestBlock
		return nil
	}
}

func WithBackfillBlockSize(backfillBlockSize uint64) RpcLogStreamerOption {
	return func(streamer *RpcLogStreamer) error {
		if backfillBlockSize == 0 {
			return fmt.Errorf("backfillBlockSize must be > 0, got %d", backfillBlockSize)
		}

		streamer.backfillBlockSize = backfillBlockSize
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

		streamer.watchers[cfg.ID] = ContractConfig{
			ID:                cfg.ID,
			FromBlockNumber:   cfg.FromBlockNumber,
			FromBlockHash:     cfg.FromBlockHash,
			Address:           cfg.Address,
			Topics:            cfg.Topics,
			MaxDisconnectTime: cfg.MaxDisconnectTime,
			eventChannel:      make(chan types.Log, 100),
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
	backfillBlockSize   uint64
}

func NewRpcLogStreamer(
	ctx context.Context,
	client ChainClient,
	logger *zap.Logger,
	options ...RpcLogStreamerOption,
) (*RpcLogStreamer, error) {
	ctx, cancel := context.WithCancel(ctx)

	streamLogger := logger.Named("rpcLogStreamer")

	streamer := &RpcLogStreamer{
		ctx:                 ctx,
		client:              client,
		logger:              streamLogger,
		cancel:              cancel,
		wg:                  sync.WaitGroup{},
		watchers:            make(map[string]ContractConfig),
		lagFromHighestBlock: 0,
		backfillBlockSize:   500,
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
			fmt.Sprintf("rpcLogStreamer-watcher-%v", watcher.Address),
			func(ctx context.Context) {
				r.watchContract(watcher)
			})
	}
}

func (r *RpcLogStreamer) Stop() {
	r.cancel()
	r.wg.Wait()
}

func (r *RpcLogStreamer) watchContract(cfg ContractConfig) {
	var (
		logger    = r.logger.With(zap.String("contractAddress", cfg.Address.Hex()))
		subLogger = logger.Named("subscription").
				With(zap.String("contractAddress", cfg.Address.Hex()))
		backfillFromBlockNumber = cfg.FromBlockNumber
		backfillFromBlockHash   = cfg.FromBlockHash
		backfillHighestBlock    = uint64(0)
	)

	defer close(cfg.eventChannel)

	innerSubCh, err := r.buildSubscriptionChannel(cfg, subLogger)
	if err != nil {
		logger.Error("failed to subscribe to contract", zap.Error(err))
		return
	}

	// TODO: Start buffering subscription logs immediately.

backfillLoop:
	for {
		select {
		case <-r.ctx.Done():
			logger.Error("Context cancelled, stopping watcher")
			return

		default:
			logs, nextBlock, nextBlockHash, highestBlock, err := r.GetNextPage(r.ctx, cfg, backfillFromBlockNumber, backfillFromBlockHash)
			if err != nil {
				logger.Error(
					"Error getting next page",
					zap.Uint64("fromBlock", backfillFromBlockNumber),
					zap.Error(err),
				)
				time.Sleep(sleepTimeOnError)
				continue
			}

			if nextBlock != nil {
				backfillFromBlockNumber = *nextBlock
				backfillFromBlockHash = nextBlockHash
				backfillHighestBlock = *nextBlock - 1
			}

			if len(logs) == 0 {
				time.Sleep(sleepTimeNoLogs)
				continue
			}

			logger.Debug(
				"Got logs",
				zap.Int("numLogs", len(logs)),
				zap.Uint64("fromBlock", backfillFromBlockNumber),
				zap.Time("time", time.Now()),
			)

			for _, log := range logs {
				cfg.eventChannel <- log
			}

			if backfillHighestBlock >= highestBlock {
				logger.Info(
					"Backfill complete, switching to websocket subscription mode",
					zap.Uint64("highestBackfillBlock", backfillHighestBlock),
				)

				break backfillLoop
			}
		}
	}

	// TODO: Apply buffered logs.

	for {
		select {
		case <-r.ctx.Done():
			logger.Error("Context cancelled, stopping watcher")
			return

		case log, open := <-innerSubCh:
			if !open {
				subLogger.Error("subscription channel closed, closing watcher")
				return
			}

			logger.Debug(
				"Received log from subscription channel",
				zap.Uint64("blockNumber", log.BlockNumber),
			)

			// TODO: Implement timelocking for logs in chains with lagFromHighestBlock > 0.

			cfg.eventChannel <- log
		}
	}
}

func (r *RpcLogStreamer) GetNextPage(
	ctx context.Context,
	cfg ContractConfig,
	fromBlockNumber uint64,
	fromBlockHash []byte,
) (logs []types.Log, nextBlockNumber *uint64, nextBlockHash []byte, highestBlock uint64, err error) {
	contractAddress := cfg.Address.Hex()

	if fromBlockNumber > 0 {
		block, err := r.client.BlockByNumber(ctx, big.NewInt(int64(fromBlockNumber+1)))
		if err != nil {
			return nil, nil, nil, 0, err
		}

		if len(fromBlockHash) == 32 &&
			!bytes.Equal(fromBlockHash, block.ParentHash().Bytes()) {
			r.logger.Error(
				"blockchain reorg detected",
				zap.Uint64("blockNumber", fromBlockNumber),
				zap.String("expectedParentHash", hex.EncodeToString(fromBlockHash)),
				zap.String("gotParentHash", block.ParentHash().Hex()),
			)

			nextBlock, err := r.client.BlockByNumber(ctx, big.NewInt(int64(fromBlockNumber-1)))
			if err != nil {
				return nil, nil, nil, 0, err
			}

			number := nextBlock.Number().Uint64()
			hash := nextBlock.Hash().Bytes()

			return nil, &number, hash, highestBlock, ErrReorg
		}
	}

	highestBlock, err = r.client.BlockNumber(ctx)
	if err != nil {
		return nil, nil, nil, 0, err
	}

	metrics.EmitIndexerMaxBlock(contractAddress, highestBlock)

	if fromBlockNumber > highestBlock {
		// TODO: Move this metric to the subscription in a subsequent PR.
		// metrics.EmitIndexerCurrentBlockLag(contractAddress, 0)
		return []types.Log{}, nil, nil, 0, nil
	}

	metrics.EmitIndexerCurrentBlockLag(contractAddress, highestBlock-fromBlockNumber)

	toBlock := min(fromBlockNumber+backfillBlocks, highestBlock)

	// TODO:(nm) Use some more clever tactics to fetch the maximum number of logs at one times by parsing error messages
	// See: https://github.com/joshstevens19/rindexer/blob/master/core/src/indexer/fetch_logs.rs#L504
	logs, err = metrics.MeasureGetLogs(contractAddress, func() ([]types.Log, error) {
		return r.client.FilterLogs(
			ctx,
			buildFilterQuery(cfg, fromBlockNumber, toBlock),
		)
	})
	if err != nil {
		return nil, nil, nil, 0, err
	}

	metrics.EmitIndexerCurrentBlock(contractAddress, toBlock)
	metrics.EmitIndexerNumLogsFound(contractAddress, len(logs))

	nextBlock, err := r.client.BlockByNumber(ctx, big.NewInt(int64(toBlock+1)))
	if err != nil {
		return nil, nil, nil, 0, err
	}

	number := nextBlock.Number().Uint64()
	hash := nextBlock.Hash().Bytes()

	return logs, &number, hash, highestBlock, nil
}

func (r *RpcLogStreamer) buildSubscriptionChannel(
	cfg ContractConfig,
	logger *zap.Logger,
) (chan types.Log, error) {
	var (
		innerSubCh = make(chan types.Log, 100)
		query      = buildBaseFilterQuery(cfg)
	)

	sub, err := r.buildSubscription(query, innerSubCh)
	if err != nil {
		logger.Error("failed to subscribe to contract", zap.Error(err))
		return nil, err
	}

	tracing.GoPanicWrap(
		r.ctx,
		&r.wg,
		fmt.Sprintf("rpcLogStreamer-watcher-subscription-%v", cfg.Address),
		func(ctx context.Context) {
			defer close(innerSubCh)
			defer sub.Unsubscribe()

			logger.Info("Subscribed to contract")

			for {
				select {
				case <-ctx.Done():
					logger.Debug("subscription context cancelled, stopping")
					return

				case err := <-sub.Err():
					if err == nil {
						continue
					}

					logger.Error("subscription error, rebuilding", zap.Error(err))

					rebuildOperation := func() (ethereum.Subscription, error) {
						sub, err = r.buildSubscription(query, innerSubCh)
						return sub, err
					}

					expBackoff := backoff.NewExponentialBackOff()
					expBackoff.InitialInterval = 1 * time.Second

					sub, err = backoff.Retry(
						r.ctx,
						rebuildOperation,
						backoff.WithBackOff(expBackoff),
						backoff.WithMaxTries(maxSubRetries),
						backoff.WithMaxElapsedTime(cfg.MaxDisconnectTime),
					)
					if err != nil {
						logger.Error(
							"failed to rebuild subscription, closing",
							zap.Error(err),
						)
						return
					}

					logger.Info("Subscription rebuilt")
				}
			}
		})

	return innerSubCh, nil
}

func (r *RpcLogStreamer) buildSubscription(
	query ethereum.FilterQuery,
	innerSubCh chan types.Log,
) (sub ethereum.Subscription, err error) {
	sub, err = r.client.SubscribeFilterLogs(
		r.ctx,
		query,
		innerSubCh,
	)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (r *RpcLogStreamer) GetEventChannel(id string) chan types.Log {
	if _, ok := r.watchers[id]; !ok {
		return nil
	}

	return r.watchers[id].eventChannel
}

func buildFilterQuery(cfg ContractConfig, from uint64, to uint64) ethereum.FilterQuery {
	query := buildBaseFilterQuery(cfg)
	query.FromBlock = new(big.Int).SetUint64(from)
	query.ToBlock = new(big.Int).SetUint64(to)

	return query
}

func buildBaseFilterQuery(cfg ContractConfig) ethereum.FilterQuery {
	return ethereum.FilterQuery{
		Addresses: []common.Address{cfg.Address},
		Topics:    [][]common.Hash{cfg.Topics},
	}
}
