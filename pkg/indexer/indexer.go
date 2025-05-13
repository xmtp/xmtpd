package indexer

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/metrics"

	"github.com/xmtp/xmtpd/pkg/tracing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/indexer/storer"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

type Indexer struct {
	ctx      context.Context
	log      *zap.Logger
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	streamer *builtStreamer
}

func NewIndexer(
	ctx context.Context,
	log *zap.Logger,
) *Indexer {
	ctx, cancel := context.WithCancel(ctx)
	return &Indexer{
		ctx:      ctx,
		log:      log.Named("indexer"),
		cancel:   cancel,
		streamer: nil,
	}
}

func (i *Indexer) Close() {
	i.log.Debug("Closing")
	if i.streamer != nil {
		if i.streamer.messagesReorgChannel != nil {
			close(i.streamer.messagesReorgChannel)
		}
		if i.streamer.identityUpdatesReorgChannel != nil {
			close(i.streamer.identityUpdatesReorgChannel)
		}
		i.streamer.streamer.Stop()
	}
	i.cancel()
	i.wg.Wait()
	i.log.Debug("Closed")
}

func (i *Indexer) StartIndexer(
	db *sql.DB,
	cfg config.ContractsOptions,
	validationService mlsvalidate.MLSValidationService,
) error {
	// Create ChainClient for log streaming
	chainClient, err := blockchain.NewChainClient(
		i.ctx,
		cfg.AppChain.RpcURL,
		cfg.AppChain.UseDatabaseClient,
		common.HexToAddress(cfg.MessagesContractAddress),
		common.HexToAddress(cfg.IdentityUpdatesContractAddress),
	)
	if err != nil {
		return err
	}
	builder := blockchain.NewRpcLogStreamBuilder(i.ctx, chainClient, i.log)
	querier := queries.New(db)

	streamer, err := configureLogStream(i.ctx, builder, cfg, querier)
	if err != nil {
		return err
	}
	i.streamer = streamer

	tracing.GoPanicWrap(
		i.ctx,
		&i.wg,
		"indexer-messages",
		func(ctx context.Context) {
			indexingLogger := i.log.Named("messages").
				With(zap.String("contractAddress", cfg.AppChain.GroupMessageBroadcasterAddress))

			messageStorer := storer.NewGroupMessageStorer(
				querier,
				indexingLogger,
				chainClient,
			)

			indexLogs(
				ctx,
				streamer.streamer.Client(),
				streamer.messagesChannel,
				streamer.messagesReorgChannel,
				indexingLogger,
				messageStorer,
				streamer.messagesBlockTracker,
				streamer.reorgHandler,
				cfg.AppChain.GroupMessageBroadcasterAddress,
			)
		})

	tracing.GoPanicWrap(
		i.ctx,
		&i.wg,
		"indexer-identities",
		func(ctx context.Context) {
			indexingLogger := i.log.Named("identity").
				With(zap.String("contractAddress", cfg.AppChain.IdentityUpdateBroadcasterAddress))

			identityStorer := storer.NewIdentityUpdateStorer(
				db,
				indexingLogger,
				chainClient,
				validationService,
			)

			indexLogs(
				ctx,
				streamer.streamer.Client(),
				streamer.identityUpdatesChannel,
				streamer.identityUpdatesReorgChannel,
				indexingLogger,
				identityStorer,
				streamer.identityUpdatesBlockTracker,
				streamer.reorgHandler,
				cfg.AppChain.IdentityUpdateBroadcasterAddress,
			)
		})

	i.streamer.streamer.Start()
	return nil
}

type builtStreamer struct {
	streamer                    *blockchain.RpcLogStreamer
	reorgHandler                ChainReorgHandler
	messagesChannel             <-chan types.Log
	messagesReorgChannel        chan<- uint64
	identityUpdatesChannel      <-chan types.Log
	identityUpdatesReorgChannel chan<- uint64
	identityUpdatesBlockTracker *BlockTracker
	messagesBlockTracker        *BlockTracker
}

func configureLogStream(
	ctx context.Context,
	builder *blockchain.RpcLogStreamBuilder,
	cfg config.ContractsOptions,
	querier *queries.Queries,
) (*builtStreamer, error) {
	messagesTopic, err := buildMessagesTopic()
	if err != nil {
		return nil, err
	}

	messagesTracker, err := NewBlockTracker(
		ctx,
		cfg.AppChain.GroupMessageBroadcasterAddress,
		querier,
	)
	if err != nil {
		return nil, err
	}

	latestBlockNumber, _ := messagesTracker.GetLatestBlock()
	messagesChannel, messagesReorgChannel := builder.ListenForContractEvent(
		latestBlockNumber,
		common.HexToAddress(cfg.AppChain.GroupMessageBroadcasterAddress),
		[]common.Hash{messagesTopic},
		cfg.AppChain.MaxChainDisconnectTime,
	)

	identityUpdatesTopic, err := buildIdentityUpdatesTopic()
	if err != nil {
		return nil, err
	}

	identityUpdatesTracker, err := NewBlockTracker(
		ctx,
		cfg.AppChain.IdentityUpdateBroadcasterAddress,
		querier,
	)
	if err != nil {
		return nil, err
	}

	latestBlockNumber, _ = identityUpdatesTracker.GetLatestBlock()
	identityUpdatesChannel, identityUpdatesReorgChannel := builder.ListenForContractEvent(
		latestBlockNumber,
		common.HexToAddress(cfg.AppChain.IdentityUpdateBroadcasterAddress),
		[]common.Hash{identityUpdatesTopic},
		cfg.AppChain.MaxChainDisconnectTime,
	)

	streamer, err := builder.Build()
	if err != nil {
		return nil, err
	}

	reorgHandler := NewChainReorgHandler(ctx, streamer.Client(), querier)

	return &builtStreamer{
		streamer:                    streamer,
		reorgHandler:                reorgHandler,
		messagesChannel:             messagesChannel,
		messagesReorgChannel:        messagesReorgChannel,
		identityUpdatesChannel:      identityUpdatesChannel,
		identityUpdatesReorgChannel: identityUpdatesReorgChannel,
		identityUpdatesBlockTracker: identityUpdatesTracker,
		messagesBlockTracker:        messagesTracker,
	}, nil
}

/*
IndexLogs will run until the eventChannel is closed, passing each event to the logStorer.

If an event fails to be stored, and the error is retryable, it will sleep for 100ms and try again.

The only non-retriable errors should be things like malformed events or failed validations.
*/
func indexLogs(
	ctx context.Context,
	client blockchain.ChainClient,
	eventChannel <-chan types.Log,
	reorgChannel chan<- uint64,
	logger *zap.Logger,
	logStorer storer.LogStorer,
	blockTracker IBlockTracker,
	reorgHandler ChainReorgHandler,
	contractAddress string,
) {
	// L3 Orbit works with Arbitrum Elastic Block Time, which under maximum load produces a block every 0.25s.
	// With a maximum throughput of 7M gas per second and a median transaction size of roughly 200k gas,
	// checking for a reorg every 60 blocks (15 seconds) means that, theoretically, a maximum of 495 messages could be affected.
	const reorgCheckInterval = 60

	var (
		storedBlockNumber uint64
		storedBlockHash   []byte
		lastBlockSeen     uint64
		reorgCheckAt      uint64
		reorgDetectedAt   uint64
		reorgBeginsAt     uint64
		reorgFinishesAt   uint64
		reorgInProgress   bool
	)

	// We don't need to listen for the ctx.Done() here, since the eventChannel will be closed when the parent context is canceled
	for event := range eventChannel {
		now := time.Now()
		// 1.1 Handle active reorg state first
		if reorgDetectedAt > 0 {
			// Under a reorg, future events are no-op
			if event.BlockNumber >= reorgDetectedAt {
				logger.Debug("discarding future event due to reorg",
					zap.Uint64("eventBlockNumber", event.BlockNumber),
					zap.Uint64("reorgBlockNumber", reorgBeginsAt))
				continue
			}
			logger.Info("starting processing reorg",
				zap.Uint64("eventBlockNumber", event.BlockNumber),
				zap.Uint64("reorgBlockNumber", reorgBeginsAt))

			// When all future events have been discarded, it means we've reached the reorg point
			storedBlockNumber, storedBlockHash = blockTracker.GetLatestBlock()
			lastBlockSeen = event.BlockNumber
			reorgDetectedAt = 0
			reorgInProgress = true
		}

		// 1.2 Handle deactivation of reorg state
		if reorgInProgress && event.BlockNumber > reorgFinishesAt {
			logger.Info("finished processing reorg",
				zap.Uint64("eventBlockNumber", event.BlockNumber),
				zap.Uint64("reorgFinishesAt", reorgFinishesAt))
			reorgInProgress = false
		}

		// 2. Get the latest block from tracker once per block
		if lastBlockSeen > 0 && lastBlockSeen != event.BlockNumber {
			storedBlockNumber, storedBlockHash = blockTracker.GetLatestBlock()
		}
		lastBlockSeen = event.BlockNumber

		// 3. Check for reorgs, when:
		// - There are no reorgs in progress
		// - There's a stored block
		// - The event block number is greater than the stored block number
		// - The check interval has passed
		skipReorgHandling := false
		if !reorgInProgress &&
			storedBlockNumber > 0 &&
			event.BlockNumber > storedBlockNumber &&
			event.BlockNumber >= reorgCheckAt+reorgCheckInterval {
			onchainBlock, err := client.BlockByNumber(ctx, big.NewInt(int64(storedBlockNumber)))
			if err != nil {
				logger.Warn(
					"error querying block from the blockchain, proceeding with event processing",
					zap.Uint64("blockNumber", storedBlockNumber),
					zap.Error(err),
				)
				skipReorgHandling = true
			}

			if !skipReorgHandling {
				reorgCheckAt = event.BlockNumber
				logger.Debug("blockchain reorg periodic check",
					zap.Uint64("blockNumber", reorgCheckAt),
				)

				if storedBlockHash != nil &&
					!bytes.Equal(storedBlockHash, onchainBlock.Hash().Bytes()) {
					logger.Warn("blockchain reorg detected",
						zap.Uint64("storedBlockNumber", storedBlockNumber),
						zap.String("storedBlockHash", hex.EncodeToString(storedBlockHash)),
						zap.String("onchainBlockHash", onchainBlock.Hash().String()),
					)

					reorgBlockNumber, reorgBlockHash, err := reorgHandler.FindReorgPoint(
						storedBlockNumber,
					)
					if err != nil && !errors.Is(err, ErrNoBlocksFound) {
						logger.Error("reorg point not found", zap.Error(err))
						continue
					}

					reorgDetectedAt = storedBlockNumber
					reorgBeginsAt = reorgBlockNumber
					reorgFinishesAt = storedBlockNumber

					if trackerErr := blockTracker.UpdateLatestBlock(ctx, reorgBlockNumber, reorgBlockHash); trackerErr != nil {
						logger.Error("error updating block tracker", zap.Error(trackerErr))
					}

					reorgChannel <- reorgBlockNumber
					continue
				}
			}
		}

		err := retry(logger, 100*time.Millisecond, contractAddress, func() storer.LogStorageError {
			return logStorer.StoreLog(ctx, event)
		})
		if err != nil {
			continue
		}

		logger.Info("Stored log", zap.Uint64("blockNumber", event.BlockNumber))
		if trackerErr := blockTracker.UpdateLatestBlock(ctx, event.BlockNumber, event.BlockHash.Bytes()); trackerErr != nil {
			logger.Error("error updating block tracker", zap.Error(trackerErr))
		}
		metrics.EmitIndexerLogProcessingTime(time.Since(now))
	}
	logger.Debug("finished")
}

func retry(
	logger *zap.Logger,
	sleep time.Duration,
	address string,
	fn func() storer.LogStorageError,
) error {
	for {
		if err := fn(); err != nil {
			logger.Error("error storing log", zap.Error(err))
			if err.ShouldRetry() {
				metrics.EmitIndexerRetryableStorageError(address)
				time.Sleep(sleep)
				continue
			}
			return err
		}
		return nil
	}
}

func buildMessagesTopic() (common.Hash, error) {
	abi, err := gm.GroupMessageBroadcasterMetaData.GetAbi()
	if err != nil {
		return common.Hash{}, err
	}
	return utils.GetEventTopic(abi, "MessageSent")
}

func buildIdentityUpdatesTopic() (common.Hash, error) {
	abi, err := iu.IdentityUpdateBroadcasterMetaData.GetAbi()
	if err != nil {
		return common.Hash{}, err
	}
	return utils.GetEventTopic(abi, "IdentityUpdateCreated")
}

func messagesContract(
	cfg config.ContractsOptions,
	client *ethclient.Client,
) (*gm.GroupMessageBroadcaster, error) {
	return gm.NewGroupMessageBroadcaster(
		common.HexToAddress(cfg.AppChain.GroupMessageBroadcasterAddress),
		client,
	)
}

func identityUpdatesContract(
	cfg config.ContractsOptions,
	client *ethclient.Client,
) (*iu.IdentityUpdateBroadcaster, error) {
	return iu.NewIdentityUpdateBroadcaster(
		common.HexToAddress(cfg.AppChain.IdentityUpdateBroadcasterAddress),
		client,
	)
}
