package loadtest

import (
	"context"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"go.uber.org/zap"
)

func LoadTest(
	logger *zap.Logger,
	options *config.LoadTestOptions,
	contractsOptions *config.ContractsOptions,
) error {
	ctx := context.Background()
	chainClient, err := blockchain.NewClient(ctx, contractsOptions.RpcUrl)
	if err != nil {
		return err
	}

	signer, err := blockchain.NewPrivateKeySigner(
		options.PrivateKey,
		contractsOptions.ChainID,
	)

	if err != nil {
		return err
	}

	dbConn, err := db.NewNamespacedDB(
		ctx,
		logger,
		options.DbConnectionString,
		"loadtest",
		5*time.Second,
		1*time.Second,
	)
	if err != nil {
		return err
	}

	nonceManager := blockchain.NewSQLBackedNonceManager(dbConn, logger)
	publisher, err := blockchain.NewBlockchainPublisher(
		ctx,
		logger,
		chainClient,
		signer,
		*contractsOptions,
		nonceManager,
	)
	if err != nil {
		return err
	}

	return doLoadTest(ctx, logger, publisher, options)
}

func doLoadTest(
	ctx context.Context,
	logger *zap.Logger,
	publisher *blockchain.BlockchainPublisher,
	options *config.LoadTestOptions,
) error {
	groupID := testutils.RandomGroupID()

	// Calculate total messages to send
	totalMessages := options.MessagesPerSecond * options.Duration

	// Create a ticker to control the rate of messages
	interval := time.Second / time.Duration(options.MessagesPerSecond)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	logger.Info("Starting load test",
		zap.Int("messagesPerSecond", options.MessagesPerSecond),
		zap.Int("duration", options.Duration),
		zap.Int("totalMessages", totalMessages),
		zap.Int("messageSize", options.MessageSizeBytes))

	// Track metrics
	var sentCount int
	startTime := time.Now()

	var wg sync.WaitGroup
	var err error

	// Send messages at the specified rate
	for i := range totalMessages {
		<-ticker.C
		wg.Add(1)
		// Send a message
		go func(messageNumber int) {
			defer wg.Done()
			if _, publishError := publisher.PublishGroupMessage(ctx, groupID, testutils.RandomBytes(options.MessageSizeBytes)); publishError != nil {
				err = publishError
				logger.Error(
					"Failed to publish message",
					zap.Error(err),
					zap.Int("messageNumber", messageNumber+1),
				)
			}
		}(i)

		sentCount++
		if sentCount%options.MessagesPerSecond == 0 {
			logger.Info("Progress",
				zap.Int("messagesSent", sentCount),
				zap.Int("totalMessages", totalMessages),
				zap.Duration("elapsed", time.Since(startTime)))
		}
	}

	wg.Wait()

	elapsed := time.Since(startTime)
	actualRate := float64(sentCount) / elapsed.Seconds()

	logger.Info("Load test completed",
		zap.Int("messagesSent", sentCount),
		zap.Duration("elapsed", elapsed),
		zap.Float64("actualMessagesPerSecond", actualRate))

	return err
}
