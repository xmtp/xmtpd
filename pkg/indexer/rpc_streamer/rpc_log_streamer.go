// Package rpcstreamer implements a log streamer that uses the RPC node to backfill history and then switch to a subscription.
package rpcstreamer

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	c "github.com/xmtp/xmtpd/pkg/indexer/common"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"go.uber.org/zap"
)

const (
	defaultBackfillBlockPageSize = 500
	defaultExpectedLogsPerBlock  = 1
	defaultLagFromHighestBlock   = 0
	maxSubReconnectionRetries    = 10
	sleepTimeOnError             = 100 * time.Millisecond
)

var (
	ErrReorg         = errors.New("blockchain reorg detected")
	ErrEndOfBackfill = errors.New("end of backfill")
)

// ContractConfig defines all the information required to filter events from a contract.
type ContractConfig struct {
	ID                   string
	FromBlockNumber      uint64
	FromBlockHash        []byte
	Address              common.Address
	Topics               []common.Hash
	MaxDisconnectTime    time.Duration
	ExpectedLogsPerBlock uint64
	eventChannel         chan types.Log
}

type RPCLogStreamerOption func(*RPCLogStreamer) error

func WithLagFromHighestBlock(lagFromHighestBlock uint8) RPCLogStreamerOption {
	return func(streamer *RPCLogStreamer) error {
		streamer.lagFromHighestBlock = lagFromHighestBlock
		return nil
	}
}

func WithBackfillBlockPageSize(backfillBlockSize uint64) RPCLogStreamerOption {
	return func(streamer *RPCLogStreamer) error {
		if backfillBlockSize == 0 {
			return fmt.Errorf("backfillBlockPageSize must be > 0, got %d", backfillBlockSize)
		}

		streamer.backfillBlockPageSize = backfillBlockSize
		return nil
	}
}

func WithContractConfig(cfg *ContractConfig) RPCLogStreamerOption {
	return func(streamer *RPCLogStreamer) error {
		if cfg == nil {
			return fmt.Errorf("contract config is nil")
		}

		if _, ok := streamer.watchers[cfg.ID]; ok {
			streamer.logger.Error("contract config already exists", zap.String("id", cfg.ID))
			return fmt.Errorf("contract config already exists: %s", cfg.ID)
		}

		config := *cfg
		if config.ExpectedLogsPerBlock == 0 {
			config.ExpectedLogsPerBlock = defaultExpectedLogsPerBlock
		}

		streamer.watchers[cfg.ID] = &config

		return nil
	}
}

// RPCLogStreamer is a naive implementation of the ChainStreamer interface.
// It queries a remote blockchain node for log events to backfill history, and then streams new events,
// to get a complete history of events on a chain.
type RPCLogStreamer struct {
	ctx                   context.Context
	cancel                context.CancelFunc
	wg                    sync.WaitGroup
	rpcClient             blockchain.ChainClient
	wsClient              blockchain.ChainClient
	logger                *zap.Logger
	watchers              map[string]*ContractConfig
	lagFromHighestBlock   uint8
	backfillBlockPageSize uint64
}

func NewRPCLogStreamer(
	ctx context.Context,
	rpcClient blockchain.ChainClient,
	wsClient blockchain.ChainClient,
	logger *zap.Logger,
	options ...RPCLogStreamerOption,
) (*RPCLogStreamer, error) {
	ctx, cancel := context.WithCancel(ctx)

	streamLogger := logger.Named(utils.RPCLogStreamerLoggerName)

	streamer := &RPCLogStreamer{
		ctx:                   ctx,
		rpcClient:             rpcClient,
		wsClient:              wsClient,
		logger:                streamLogger,
		cancel:                cancel,
		wg:                    sync.WaitGroup{},
		watchers:              make(map[string]*ContractConfig),
		lagFromHighestBlock:   defaultLagFromHighestBlock,
		backfillBlockPageSize: defaultBackfillBlockPageSize,
	}

	for _, option := range options {
		if err := option(streamer); err != nil {
			streamLogger.Error("failed to apply option", zap.Error(err))
			return nil, err
		}
	}

	for _, w := range streamer.watchers {
		w.eventChannel = make(chan types.Log, streamer.backfillBlockPageSize*w.ExpectedLogsPerBlock)
	}

	return streamer, nil
}

func (r *RPCLogStreamer) Start() error {
	for _, watcher := range r.watchers {
		err := r.validateWatcher(watcher)
		if err != nil {
			return err
		}

		tracing.GoPanicWrap(
			r.ctx,
			&r.wg,
			fmt.Sprintf("rpcLogStreamer-watcher-%v", watcher.Address),
			func(ctx context.Context) {
				r.watchContract(watcher)
			})
	}

	return nil
}

func (r *RPCLogStreamer) Stop() {
	r.cancel()
	r.wg.Wait()
}

func (r *RPCLogStreamer) watchContract(cfg *ContractConfig) {
	var (
		logger                  = r.logger.With(utils.ContractAddressField(cfg.Address.Hex()))
		backfillFromBlockNumber = cfg.FromBlockNumber
		backfillFromBlockHash   = cfg.FromBlockHash
		innerSubCh              = make(chan types.Log, cfg.ExpectedLogsPerBlock*10)
	)

	defer close(cfg.eventChannel)

	sub, err := r.buildSubscription(cfg, innerSubCh)
	if err != nil {
		logger.Error("failed to subscribe to contract", zap.Error(err))
		return
	}

	for {
	backfillLoop:
		for {
			select {
			case <-r.ctx.Done():
				logger.Info("context cancelled, stopping watcher")
				return

			case err := <-sub.Err():
				if err == nil {
					continue
				}

				logger.Warn("subscription error, rebuilding", zap.Error(err))
				sub, err = r.buildSubscriptionWithBackoff(cfg, innerSubCh)
				if err != nil {
					logger.Fatal("failed rebuilding subscription after max disconnect time", zap.Error(err), zap.String("max_disconnect_time", cfg.MaxDisconnectTime.String()))
				}

			default:
				response, err := r.GetNextPage(r.ctx, cfg, backfillFromBlockNumber, backfillFromBlockHash)
				if err != nil {
					switch err {
					case ErrEndOfBackfill:
						if response.NextBlockNumber != nil {
							backfillFromBlockNumber = *response.NextBlockNumber
							backfillFromBlockHash = response.NextBlockHash
						}

						if len(response.Logs) > 0 {
							for _, log := range response.Logs {
								cfg.eventChannel <- log
							}
						} else {
							cfg.eventChannel <- c.NewUpdateProgressLog(backfillFromBlockNumber, backfillFromBlockHash)
						}

						logger.Info("backfill complete, switching to subscription")
						break backfillLoop

					case ErrReorg:
						logger.Warn("reorg detected, rolled back to block", utils.BlockNumberField(*response.NextBlockNumber))

					default:
						if isBlockPageSizeError(err) {
							blockPageSize, err := extractBlockPageSize(err.Error())
							if err != nil {
								logger.Error("incorrect backfill block page size, please check your configuration", zap.Error(err))
								continue
							}

							logger.Info("adjusting backfill block page size", zap.Uint64("new_block_page_size", blockPageSize))
							r.backfillBlockPageSize = blockPageSize

							continue
						}

						logger.Error(
							"error getting next page",
							utils.BlockNumberField(backfillFromBlockNumber),
							zap.Error(err),
						)

						time.Sleep(sleepTimeOnError)
						continue
					}
				}

				if response.NextBlockNumber != nil {
					backfillFromBlockNumber = *response.NextBlockNumber
					backfillFromBlockHash = response.NextBlockHash
				}

				if len(response.Logs) == 0 {
					cfg.eventChannel <- c.NewUpdateProgressLog(backfillFromBlockNumber, backfillFromBlockHash)
					continue
				}

				logger.Debug(
					"got logs",
					utils.NumEnvelopesField(len(response.Logs)),
					utils.BlockNumberField(backfillFromBlockNumber),
					utils.TimeField(time.Now()),
				)

				for _, log := range response.Logs {
					cfg.eventChannel <- log
				}
			}
		}

		// From now on we are operating on the subscription, and we no longer check what the highest block is.
		metrics.EmitIndexerCurrentBlockLag(cfg.Address.Hex(), 0)

	subscriptionLoop:
		for {
			select {
			case <-r.ctx.Done():
				logger.Info("context cancelled, stopping watcher")
				return

			case err := <-sub.Err():
				if err == nil {
					continue
				}

				logger.Warn("subscription error, rebuilding", zap.Error(err))
				sub, err = r.buildSubscriptionWithBackoff(cfg, innerSubCh)
				if err != nil {
					logger.Fatal("failed rebuilding subscription after max disconnect time", zap.Error(err), zap.String("max_disconnect_time", cfg.MaxDisconnectTime.String()))
				}

				// Backfill everything that was missed.
				break subscriptionLoop

			case log, open := <-innerSubCh:
				if !open {
					logger.Info("subscription channel closed, closing watcher")
					return
				}

				backfillFromBlockNumber = log.BlockNumber
				backfillFromBlockHash = log.BlockHash.Bytes()

				logger.Debug(
					"received log from subscription channel",
					utils.BlockNumberField(log.BlockNumber),
				)

				metrics.EmitIndexerCurrentBlock(cfg.Address.Hex(), log.BlockNumber)
				metrics.EmitIndexerNumLogsFound(cfg.Address.Hex(), 1)
				metrics.EmitIndexerMaxBlock(cfg.Address.Hex(), log.BlockNumber)

				// TODO: Implement timelocking for logs in chains with lagFromHighestBlock > 0.

				cfg.eventChannel <- log
			}
		}
	}
}

type GetNextPageResponse struct {
	Logs            []types.Log
	NextBlockNumber *uint64
	NextBlockHash   []byte
}

func (r *RPCLogStreamer) GetNextPage(
	ctx context.Context,
	cfg *ContractConfig,
	fromBlockNumber uint64,
	fromBlockHash []byte,
) (GetNextPageResponse, error) {
	r.logger.Debug(
		"getting next page",
		utils.BlockNumberField(fromBlockNumber),
		utils.HashField(hex.EncodeToString(fromBlockHash)),
	)

	contractAddress := cfg.Address.Hex()

	highestBlock, err := r.rpcClient.BlockNumber(ctx)
	if err != nil {
		return GetNextPageResponse{}, err
	}

	// Do not check for reorgs at block height 0. Genesis does not have a parent.
	if fromBlockNumber > 0 {
		nextBlockHeader, err := r.rpcClient.HeaderByNumber(
			ctx,
			big.NewInt(int64(fromBlockNumber+1)),
		)
		if err != nil {
			return GetNextPageResponse{}, err
		}

		// Compare the current hash against the next block's parent hash.
		if len(fromBlockHash) == 32 &&
			!bytes.Equal(fromBlockHash, nextBlockHeader.ParentHash.Bytes()) {
			// If the current hash doesn't match the next block's parent hash,
			// move one block back and use that hash as the new starting point.
			nextBlockNumber := fromBlockNumber - 1
			nextBlockHash, err := r.rpcClient.HeaderByNumber(
				ctx,
				big.NewInt(int64(nextBlockNumber)),
			)
			if err != nil {
				return GetNextPageResponse{}, err
			}

			return GetNextPageResponse{
				Logs:            []types.Log{},
				NextBlockNumber: &nextBlockNumber,
				NextBlockHash:   nextBlockHash.Hash().Bytes(),
			}, ErrReorg
		}
	}

	metrics.EmitIndexerMaxBlock(contractAddress, highestBlock)

	if fromBlockNumber > highestBlock {
		return GetNextPageResponse{}, fmt.Errorf(
			"fromBlockNumber is higher than highestBlock: %d > %d",
			fromBlockNumber,
			highestBlock,
		)
	}

	metrics.EmitIndexerCurrentBlockLag(contractAddress, highestBlock-fromBlockNumber)

	// Define the upper bound of the block range to fetch logs from.
	// The range is inclusive, so we subtract 1 to ensure that exactly `backfillBlockPageSize` blocks
	// are queried per page. Without this, we would query one extra block per page (off-by-one error).
	toBlock := min(fromBlockNumber+r.backfillBlockPageSize-1, highestBlock)

	// TODO:(nm) Use some more clever tactics to fetch the maximum number of logs at one times by parsing error messages
	// See: https://github.com/joshstevens19/rindexer/blob/master/core/src/indexer/fetch_logs.rs#L504
	logs, err := metrics.MeasureGetLogs(contractAddress, func() ([]types.Log, error) {
		return r.rpcClient.FilterLogs(
			ctx,
			buildFilterQuery(*cfg, fromBlockNumber, toBlock),
		)
	})
	if err != nil {
		return GetNextPageResponse{}, err
	}

	metrics.EmitIndexerCurrentBlock(contractAddress, toBlock)
	metrics.EmitIndexerNumLogsFound(contractAddress, len(logs))

	if toBlock+1 > highestBlock {
		return GetNextPageResponse{
			Logs:            logs,
			NextBlockNumber: nil,
			NextBlockHash:   nil,
		}, ErrEndOfBackfill
	}

	nextBlockNumber := uint64(toBlock + 1)
	nextBlockHash, err := r.rpcClient.HeaderByNumber(ctx, big.NewInt(int64(nextBlockNumber)))
	if err != nil {
		return GetNextPageResponse{}, err
	}

	return GetNextPageResponse{
		Logs:            logs,
		NextBlockNumber: &nextBlockNumber,
		NextBlockHash:   nextBlockHash.Hash().Bytes(),
	}, nil
}

func (r *RPCLogStreamer) buildSubscription(
	cfg *ContractConfig,
	innerSubCh chan types.Log,
) (sub ethereum.Subscription, err error) {
	query := buildBaseFilterQuery(*cfg)

	highestBlock, err := r.wsClient.BlockNumber(r.ctx)
	if err != nil {
		return nil, err
	}

	// In most implementations, FromBlock is ignored by the RPC nodes when subscribing to logs.
	// As the logs received are always starting at the highest block.
	// However, some implementations will set FromBlock to 0 if not present, which
	// would cause a failure if the RPC doesn't support big lookback ranges.
	query.FromBlock = new(big.Int).SetUint64(highestBlock)

	sub, err = r.wsClient.SubscribeFilterLogs(
		r.ctx,
		query,
		innerSubCh,
	)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

// buildSubscriptionWithBackoff builds a subscription with backoff.
// It guarantees the indexer is not disconnected from the chain for
// more than cfg.MaxDisconnectTime.
func (r *RPCLogStreamer) buildSubscriptionWithBackoff(
	cfg *ContractConfig,
	innerSubCh chan types.Log,
) (sub ethereum.Subscription, err error) {
	rebuildOperation := func() (ethereum.Subscription, error) {
		sub, err = r.buildSubscription(cfg, innerSubCh)
		return sub, err
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.InitialInterval = 1 * time.Second

	sub, err = backoff.Retry(
		r.ctx,
		rebuildOperation,
		backoff.WithBackOff(expBackoff),
		backoff.WithMaxTries(maxSubReconnectionRetries),
		backoff.WithMaxElapsedTime(cfg.MaxDisconnectTime),
	)
	if err != nil {
		r.logger.Error(
			"failed to rebuild subscription, closing",
			zap.Error(err),
		)
		return nil, err
	}

	r.logger.Info("subscription rebuilt")

	return sub, nil
}

func (r *RPCLogStreamer) validateWatcher(cfg *ContractConfig) error {
	if cfg == nil {
		return fmt.Errorf("watcher is nil")
	}

	if cfg.eventChannel == nil {
		return fmt.Errorf("event channel is nil")
	}

	testCh := make(chan types.Log, 100)
	defer close(testCh)

	sub, err := r.buildSubscription(cfg, testCh)
	if err != nil {
		return fmt.Errorf("failed to validate watcher %s: %w", cfg.Address.Hex(), err)
	}

	defer sub.Unsubscribe()

	return nil
}

func (r *RPCLogStreamer) GetEventChannel(id string) <-chan types.Log {
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

func isBlockPageSizeError(err error) bool {
	// This function has to be extended in the case XMTP chain is running on a different RPC provider.
	return strings.Contains(err.Error(), "You can make eth_getLogs requests with up to a")
}

func extractBlockPageSize(s string) (uint64, error) {
	// Example error message:
	// You can make eth_getLogs requests with up to a 500 block range. Based on your parameters, this block range should work: [0x19a788b, 0x19a7a7e]
	re, err := regexp.Compile(
		`(?i)eth_getLogs\s+requests\s+with\s+up\s+to\s+a\s+(\d+)\s+block\s+range`,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to compile regex: %w", err)
	}

	matches := re.FindStringSubmatch(s)
	if len(matches) != 2 {
		return 0, fmt.Errorf("failed to extract block range from string: %q", s)
	}

	blockPageSize, err := strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block range: %w", err)
	}

	if blockPageSize == 0 {
		return 0, fmt.Errorf("block page size is 0")
	}

	// The error message is inclusive, so we subtract 1 to get the correct block page size.
	blockPageSize--

	return blockPageSize, nil
}
