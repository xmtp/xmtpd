package indexer

import (
	"context"
	"database/sql"
	"time"

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

// Start the indexer and run until the context is canceled
func StartIndexer(
	ctx context.Context,
	logger *zap.Logger,
	db *sql.DB,
	cfg config.ContractsOptions,
	validationService mlsvalidate.MLSValidationService,
) error {
	client, err := blockchain.NewClient(ctx, cfg.RpcUrl)
	if err != nil {
		return err
	}
	builder := blockchain.NewRpcLogStreamBuilder(client, logger)
	querier := queries.New(db)

	streamer, err := configureLogStream(builder, cfg)
	if err != nil {
		return err
	}

	messagesContract, err := messagesContract(cfg, client)
	if err != nil {
		return err
	}

	go indexLogs(
		ctx,
		streamer.messagesChannel,
		logger.Named("indexLogs").With(zap.String("contractAddress", cfg.MessagesContractAddress)),
		storer.NewGroupMessageStorer(querier, logger, messagesContract),
	)

	identityUpdatesContract, err := identityUpdatesContract(cfg, client)
	if err != nil {
		return err
	}

	go indexLogs(
		ctx,
		streamer.identityUpdatesChannel,
		logger.Named("indexLogs").
			With(zap.String("contractAddress", cfg.IdentityUpdatesContractAddress)),
		storer.NewIdentityUpdateStorer(db, logger, identityUpdatesContract, validationService),
	)

	return streamer.streamer.Start(ctx)
}

type builtStreamer struct {
	streamer               *blockchain.RpcLogStreamer
	messagesChannel        <-chan types.Log
	identityUpdatesChannel <-chan types.Log
}

func configureLogStream(
	builder *blockchain.RpcLogStreamBuilder,
	cfg config.ContractsOptions,
) (*builtStreamer, error) {
	messagesTopic, err := buildMessagesTopic()
	if err != nil {
		return nil, err
	}

	messagesChannel := builder.ListenForContractEvent(
		0,
		common.HexToAddress(cfg.MessagesContractAddress),
		[]common.Hash{messagesTopic},
	)

	identityUpdatesTopic, err := buildIdentityUpdatesTopic()
	if err != nil {
		return nil, err
	}

	identityUpdatesChannel := builder.ListenForContractEvent(
		0,
		common.HexToAddress(cfg.IdentityUpdatesContractAddress),
		[]common.Hash{identityUpdatesTopic},
	)

	streamer, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return &builtStreamer{
		streamer:               streamer,
		messagesChannel:        messagesChannel,
		identityUpdatesChannel: identityUpdatesChannel,
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
			}
			break Retry

		}
	}
	logger.Info("finished")
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
