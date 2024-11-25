package indexer

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/tracing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/abis"
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

func (s *Indexer) Close() {
	s.log.Debug("Closing")
	if s.streamer != nil {
		s.streamer.streamer.Stop()
	}
	s.cancel()
	s.wg.Wait()
	s.log.Debug("Closed")
}

func (i *Indexer) StartIndexer(
	db *sql.DB,
	cfg config.ContractsOptions,
	validationService mlsvalidate.MLSValidationService) error {
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
				streamer.messagesChannel,
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
				streamer.identityUpdatesChannel, indexingLogger,
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
	identityUpdatesChannel      <-chan types.Log
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

	messagesChannel := builder.ListenForContractEvent(
		messagesTracker.GetLatestBlock(),
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

	identityUpdatesChannel := builder.ListenForContractEvent(
		identityUpdatesTracker.GetLatestBlock(),
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
		identityUpdatesChannel:      identityUpdatesChannel,
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
	eventChannel <-chan types.Log,
	logger *zap.Logger,
	logStorer storer.LogStorer,
	blockTracker IBlockTracker,
) {
	var err storer.LogStorageError
	// We don't need to listen for the ctx.Done() here, since the eventChannel will be closed when the parent context is canceled
	for event := range eventChannel {
	Retry:
		for {
			err = logStorer.StoreLog(ctx, event)
			if err != nil {
				logger.Error("error storing log", zap.Error(err))
				if err.ShouldRetry() {
					time.Sleep(100 * time.Millisecond)
					continue Retry
				}
			} else {
				logger.Info("Stored log", zap.Uint64("blockNumber", event.BlockNumber))
				if trackerErr := blockTracker.UpdateLatestBlock(ctx, event.BlockNumber); trackerErr != nil {
					logger.Error("error updating block tracker", zap.Error(trackerErr))
				}
			}
			break Retry

		}
	}
	logger.Debug("finished")
}

func buildMessagesTopic() (common.Hash, error) {
	abi, err := abis.GroupMessagesMetaData.GetAbi()
	if err != nil {
		return common.Hash{}, err
	}
	return utils.GetEventTopic(abi, "MessageSent")
}

func buildIdentityUpdatesTopic() (common.Hash, error) {
	abi, err := abis.IdentityUpdatesMetaData.GetAbi()
	if err != nil {
		return common.Hash{}, err
	}
	return utils.GetEventTopic(abi, "IdentityUpdateCreated")
}

func messagesContract(
	cfg config.ContractsOptions,
	client *ethclient.Client,
) (*abis.GroupMessages, error) {
	return abis.NewGroupMessages(
		common.HexToAddress(cfg.MessagesContractAddress),
		client,
	)
}

func identityUpdatesContract(
	cfg config.ContractsOptions,
	client *ethclient.Client,
) (*abis.IdentityUpdates, error) {
	return abis.NewIdentityUpdates(
		common.HexToAddress(cfg.IdentityUpdatesContractAddress),
		client,
	)
}
