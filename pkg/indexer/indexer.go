package indexer

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/abis"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/indexer/storer"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

// Start the indexer and run until the context is canceled
func StartIndexer(
	ctx context.Context,
	logger *zap.Logger,
	queries *queries.Queries,
	cfg config.ContractsOptions,
) error {
	client, err := blockchain.NewClient(ctx, cfg.RpcUrl)
	if err != nil {
		return err
	}
	builder := blockchain.NewRpcLogStreamBuilder(client, logger)

	messagesTopic, err := buildMessagesTopic()
	if err != nil {
		return err
	}

	messagesChannel := builder.ListenForContractEvent(
		0,
		common.HexToAddress(cfg.MessagesContractAddress),
		[]common.Hash{messagesTopic},
	)

	streamer, err := builder.Build()
	if err != nil {
		return err
	}

	messagesContract, err := messagesContract(cfg, client)
	if err != nil {
		return err
	}

	go indexLogs(
		ctx,
		messagesChannel,
		logger.Named("indexLogs").With(zap.String("contractAddress", cfg.MessagesContractAddress)),
		storer.NewGroupMessageStorer(queries, logger, messagesContract),
	)

	return streamer.Start(ctx)
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

func messagesContract(
	cfg config.ContractsOptions,
	client *ethclient.Client,
) (*abis.GroupMessages, error) {
	return abis.NewGroupMessages(
		common.HexToAddress(cfg.MessagesContractAddress),
		client,
	)
}
