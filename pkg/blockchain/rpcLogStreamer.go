package blockchain

import (
	"context"
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

// Struct defining all the information required to filter events from logs
type ContractConfig struct {
	ID                string
	FromBlock         uint64
	Address           common.Address
	Topics            []common.Hash
	MaxDisconnectTime time.Duration
	eventChannel      chan types.Log
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

		streamer.watchers[cfg.ID] = ContractConfig{
			ID:                cfg.ID,
			FromBlock:         cfg.FromBlock,
			Address:           cfg.Address,
			Topics:            cfg.Topics,
			MaxDisconnectTime: cfg.MaxDisconnectTime,
			eventChannel:      make(chan types.Log, 100),
			reorgChannel:      make(chan uint64, 1),
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
		isSubEnabled         = false
		backfillStartBlock   = cfg.FromBlock
		highestBackfillBlock uint64
	)

	if backfillStartBlock > 0 {
		highestBackfillBlock = backfillStartBlock - 1
	}

	defer close(cfg.eventChannel)
	defer close(cfg.reorgChannel)

	backfillEndBlock, err := r.client.BlockNumber(r.ctx)
	if err != nil {
		logger.Error("failed to get highest block", zap.Error(err))
		return
	}

	innerBackfillCh, cancelBackfill, err := r.buildBackfillChannel(cfg, backfillStartBlock)
	if err != nil {
		logger.Error("failed to retrieve historical logs", zap.Error(err))
		return
	}

	innerSubCh, err := r.buildSubscriptionChannel(cfg, subLogger)
	if err != nil {
		logger.Error("failed to subscribe to contract", zap.Error(err))
		return
	}

	logger.Info("Starting watcher")

	for {
		select {
		case <-r.ctx.Done():
			logger.Debug("Context cancelled, stopping watcher")
			return

		case reorgBlock, open := <-cfg.reorgChannel:
			if !open {
				logger.Debug("Reorg channel closed")
				return
			}
			backfillStartBlock = reorgBlock
			logger.Warn(
				"Blockchain reorg detected, resuming from block",
				zap.Uint64("fromBlock", backfillStartBlock),
			)

		case log, open := <-innerSubCh:
			if !open {
				subLogger.Fatal("subscription channel closed, closing watcher")
				return
			}

			if !isSubEnabled {
				// backfillEndBlock is a moving target.
				// Subscription is always ahead of backfill, so we update backfillEndBlock when a new log is received.
				if log.BlockNumber > backfillEndBlock {
					backfillEndBlock = log.BlockNumber
				}

				// TODO: Next PR to introduce a log buffer.
				// The buffer has to be of size lagFromHighestBlock, at minimum.

				continue
			}

			subLogger.Debug(
				"Sending log to subscription channel",
				zap.Uint64("fromBlock", log.BlockNumber),
			)

			cfg.eventChannel <- log

		case log, open := <-innerBackfillCh:
			if !open {
				continue
			}

			// This is a guard, this case shouldn't happen.
			// When the subscription is enabled the backfill is always cancelled.
			if isSubEnabled {
				continue
			}

			if log.BlockNumber > highestBackfillBlock {
				highestBackfillBlock = log.BlockNumber
			}

			cfg.eventChannel <- log

			// Check if subscription has to be enabled only after processing all logs in the last batch.
			// Duplicated logs are not a problem, lost logs are.
			if highestBackfillBlock+uint64(r.lagFromHighestBlock) >= backfillEndBlock {
				logger.Debug(
					"Backfill complete, enabling subscription mode",
					zap.Uint64("blockNumber", highestBackfillBlock),
				)

				isSubEnabled = true

				cancelBackfill()

				// TODO: Next PR to send buffered logs when switched to subscription mode.
			}
		}
	}
}

func (r *RpcLogStreamer) buildBackfillChannel(
	cfg ContractConfig,
	fromBlock uint64,
) (chan types.Log, context.CancelFunc, error) {
	var (
		innerBackfillCh = make(chan types.Log, 100)
		logger          = r.logger.Named("backfiller").
				With(zap.String("contractAddress", cfg.Address.Hex()))
	)

	ctxwc, cancel := context.WithCancel(r.ctx)

	tracing.GoPanicWrap(
		ctxwc,
		&r.wg,
		fmt.Sprintf("rpcLogStreamer-watcher-backfiller-%v", cfg.Address),
		func(ctx context.Context) {
			defer close(innerBackfillCh)

			logger.Info("Backfilling logs", zap.Uint64("fromBlock", fromBlock))

			for {
				select {
				case <-ctx.Done():
					logger.Debug("backfiller context cancelled, stopping")
					return

				default:
					logs, nextBlock, err := r.GetNextPage(ctx, cfg, fromBlock)
					if err != nil {
						logger.Error(
							"Error getting next page",
							zap.Uint64("fromBlock", fromBlock),
							zap.Error(err),
						)
						time.Sleep(sleepTimeOnError)
						continue
					}

					if nextBlock != nil {
						fromBlock = *nextBlock
					}

					if len(logs) == 0 {
						time.Sleep(sleepTimeNoLogs)
						continue
					}

					logger.Debug(
						"Got logs",
						zap.Int("numLogs", len(logs)),
						zap.Uint64("fromBlock", fromBlock),
						zap.Time("time", time.Now()),
					)

					for _, log := range logs {
						innerBackfillCh <- log
					}
				}
			}
		})

	return innerBackfillCh, cancel, nil
}

func (r *RpcLogStreamer) GetNextPage(
	ctx context.Context,
	config ContractConfig,
	fromBlock uint64,
) (logs []types.Log, nextBlock *uint64, err error) {
	contractAddress := config.Address.Hex()
	highestBlock, err := r.client.BlockNumber(ctx)
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

	toBlock := min(fromBlock+backfillBlocks, highestBlockCanProcess)

	// TODO:(nm) Use some more clever tactics to fetch the maximum number of logs at one times by parsing error messages
	// See: https://github.com/joshstevens19/rindexer/blob/master/core/src/indexer/fetch_logs.rs#L504
	logs, err = metrics.MeasureGetLogs(contractAddress, func() ([]types.Log, error) {
		return r.client.FilterLogs(
			ctx,
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

func (r *RpcLogStreamer) buildSubscriptionChannel(
	cfg ContractConfig,
	logger *zap.Logger,
) (chan types.Log, error) {
	var (
		innerSubCh = make(chan types.Log, 100)
		query      = buildSubscriptionFilterQuery(cfg)
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

func (r *RpcLogStreamer) Client() ChainClient {
	return r.client
}

func (r *RpcLogStreamer) GetEventChannel(id string) chan types.Log {
	if _, ok := r.watchers[id]; !ok {
		return nil
	}

	return r.watchers[id].eventChannel
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
	addresses := []common.Address{contractConfig.Address}
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

func buildSubscriptionFilterQuery(cfg ContractConfig) ethereum.FilterQuery {
	addresses := []common.Address{cfg.Address}
	topics := [][]common.Hash{}
	for _, topic := range cfg.Topics {
		topics = append(topics, []common.Hash{topic})
	}

	return ethereum.FilterQuery{
		Addresses: addresses,
		Topics:    topics,
	}
}
