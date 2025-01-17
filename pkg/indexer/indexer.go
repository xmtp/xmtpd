package indexer

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/hex"
	"math/big"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/tracing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/contracts/pkg/groupmessages"
	"github.com/xmtp/xmtpd/contracts/pkg/identityupdates"
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
	client, err := blockchain.NewClient(i.ctx, cfg.RpcUrl)
	if err != nil {
		return err
	}
	builder := blockchain.NewRpcLogStreamBuilder(i.ctx, client, i.log)
	querier := queries.New(db)

	streamer, err := configureLogStream(i.ctx, builder, cfg, querier)
	if err != nil {
		return err
	}
	i.streamer = streamer

	messagesContract, err := messagesContract(cfg, client)
	if err != nil {
		return err
	}

	tracing.GoPanicWrap(
		i.ctx,
		&i.wg,
		"indexer-messages",
		func(ctx context.Context) {
			indexingLogger := i.log.Named("messages").
				With(zap.String("contractAddress", cfg.MessagesContractAddress))

			indexLogs(
				ctx,
				streamer.streamer.Client(),
				streamer.messagesChannel,
				streamer.messagesReorgChannel,
				indexingLogger,
				storer.NewGroupMessageStorer(querier, indexingLogger, messagesContract),
				streamer.messagesBlockTracker,
			)
		})

	identityUpdatesContract, err := identityUpdatesContract(cfg, client)
	if err != nil {
		return err
	}

	tracing.GoPanicWrap(
		i.ctx,
		&i.wg,
		"indexer-identities",
		func(ctx context.Context) {
			indexingLogger := i.log.Named("identity").
				With(zap.String("contractAddress", cfg.IdentityUpdatesContractAddress))
			indexLogs(
				ctx,
				streamer.streamer.Client(),
				streamer.identityUpdatesChannel,
				streamer.identityUpdatesReorgChannel,
				indexingLogger,
				storer.NewIdentityUpdateStorer(
					db,
					indexingLogger,
					identityUpdatesContract,
					validationService,
				),
				streamer.identityUpdatesBlockTracker,
			)
		})

	i.streamer.streamer.Start()
	return nil
}

type builtStreamer struct {
	streamer                    *blockchain.RpcLogStreamer
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

	messagesTracker, err := NewBlockTracker(ctx, cfg.MessagesContractAddress, querier)
	if err != nil {
		return nil, err
	}

	latestBlockNumber, _ := messagesTracker.GetLatestBlock()
	messagesChannel, messagesReorgChannel := builder.ListenForContractEvent(
		latestBlockNumber,
		common.HexToAddress(cfg.MessagesContractAddress),
		[]common.Hash{messagesTopic},
		cfg.MaxChainDisconnectTime,
	)

	identityUpdatesTopic, err := buildIdentityUpdatesTopic()
	if err != nil {
		return nil, err
	}

	identityUpdatesTracker, err := NewBlockTracker(ctx, cfg.IdentityUpdatesContractAddress, querier)
	if err != nil {
		return nil, err
	}

	latestBlockNumber, _ = identityUpdatesTracker.GetLatestBlock()
	identityUpdatesChannel, identityUpdatesReorgChannel := builder.ListenForContractEvent(
		latestBlockNumber,
		common.HexToAddress(cfg.IdentityUpdatesContractAddress),
		[]common.Hash{identityUpdatesTopic},
		cfg.MaxChainDisconnectTime,
	)

	streamer, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return &builtStreamer{
		streamer:                    streamer,
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
) {
	var (
		errStorage         storer.LogStorageError
		storedBlockNumber  uint64
		storedBlockHash    []byte
		lastBlockSeen      uint64
		reorgCheckInterval uint64 = 10 // TODO: Adapt based on blocks per batch
		reorgDetectedAt    uint64
		reorgLastCheckAt   uint64
	)

	// We don't need to listen for the ctx.Done() here, since the eventChannel will be closed when the parent context is canceled
	for event := range eventChannel {
		// When a reorg is detected, discard the events higher than the block number of the reorg
		if reorgDetectedAt > 0 {
			if event.BlockNumber >= reorgDetectedAt {
				logger.Debug("discarding event due to reorg",
					zap.Uint64("event_block", event.BlockNumber),
					zap.Uint64("reorg_block", reorgDetectedAt))
				continue
			}
			logger.Info("resumed processing after reorg",
				zap.Uint64("block_number", event.BlockNumber),
				zap.Uint64("reorg_block", reorgDetectedAt))
			reorgDetectedAt = 0
		}

		if lastBlockSeen > 0 && lastBlockSeen != event.BlockNumber {
			storedBlockNumber, storedBlockHash = blockTracker.GetLatestBlock()
		}

		lastBlockSeen = event.BlockNumber

		if storedBlockNumber > 0 &&
			event.BlockNumber > storedBlockNumber &&
			event.BlockNumber >= reorgLastCheckAt+reorgCheckInterval {
			blockByNumber, err := client.BlockByNumber(ctx, big.NewInt(int64(storedBlockNumber)))
			if err != nil {
				logger.Error("error getting block",
					zap.Uint64("blockNumber", storedBlockNumber),
					zap.Error(err),
				)
				continue
			}

			reorgLastCheckAt = event.BlockNumber

			if !bytes.Equal(storedBlockHash, blockByNumber.Hash().Bytes()) {
				logger.Warn("blockchain reorg detected",
					zap.Uint64("storedBlockNumber", storedBlockNumber),
					zap.String("storedBlockHash", hex.EncodeToString(storedBlockHash)),
					zap.String("onchainBlockHash", blockByNumber.Hash().String()),
				)
				reorgChannel <- storedBlockNumber
				reorgDetectedAt = storedBlockNumber
				continue
			}
		}

	Retry:
		for {
			errStorage = logStorer.StoreLog(ctx, event, false)
			if errStorage != nil {
				logger.Error("error storing log", zap.Error(errStorage))
				if errStorage.ShouldRetry() {
					time.Sleep(100 * time.Millisecond)
					continue Retry
				}
			} else {
				logger.Info("Stored log", zap.Uint64("blockNumber", event.BlockNumber))
				if trackerErr := blockTracker.UpdateLatestBlock(ctx, event.BlockNumber, event.BlockHash.Bytes()); trackerErr != nil {
					logger.Error("error updating block tracker", zap.Error(trackerErr))
				}
			}
			break Retry
		}
	}
	logger.Debug("finished")
}

func buildMessagesTopic() (common.Hash, error) {
	abi, err := groupmessages.GroupMessagesMetaData.GetAbi()
	if err != nil {
		return common.Hash{}, err
	}
	return utils.GetEventTopic(abi, "MessageSent")
}

func buildIdentityUpdatesTopic() (common.Hash, error) {
	abi, err := identityupdates.IdentityUpdatesMetaData.GetAbi()
	if err != nil {
		return common.Hash{}, err
	}
	return utils.GetEventTopic(abi, "IdentityUpdateCreated")
}

func messagesContract(
	cfg config.ContractsOptions,
	client *ethclient.Client,
) (*groupmessages.GroupMessages, error) {
	return groupmessages.NewGroupMessages(
		common.HexToAddress(cfg.MessagesContractAddress),
		client,
	)
}

func identityUpdatesContract(
	cfg config.ContractsOptions,
	client *ethclient.Client,
) (*identityupdates.IdentityUpdates, error) {
	return identityupdates.NewIdentityUpdates(
		common.HexToAddress(cfg.IdentityUpdatesContractAddress),
		client,
	)
}
