package blockchain

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

func NewClient(ctx context.Context, rpcUrl string) (*ethclient.Client, error) {
	return ethclient.DialContext(ctx, rpcUrl)
}

// executeTransaction is a helper function that:
// - executes a transaction
// - waits for it to be mined
// - processes the event logs
func ExecuteTransaction(
	ctx context.Context,
	signer TransactionSigner,
	logger *zap.Logger,
	client *ethclient.Client,
	txFunc func(*bind.TransactOpts) (*types.Transaction, error),
	eventParser func(*types.Log) (interface{}, error),
	logHandler func(interface{}),
) error {
	if signer == nil {
		return fmt.Errorf("no signer provided")
	}

	tx, err := txFunc(&bind.TransactOpts{
		Context: ctx,
		From:    signer.FromAddress(),
		Signer:  signer.SignerFunc(),
	})
	if err != nil {
		return err
	}

	receipt, err := WaitForTransaction(
		ctx,
		logger,
		client,
		2*time.Second,
		250*time.Millisecond,
		tx.Hash(),
	)
	if err != nil {
		return err
	}

	for _, log := range receipt.Logs {
		event, err := eventParser(log)
		if err != nil {
			continue
		}
		logHandler(event)
	}

	return nil
}

// Waits for the given transaction hash to have been submitted to the chain and soft confirmed
func WaitForTransaction(
	ctx context.Context,
	logger *zap.Logger,
	client *ethclient.Client,
	timeout time.Duration,
	pollSleep time.Duration,
	hash common.Hash,
) (*types.Receipt, error) {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(timeout))
	defer cancel()

	ticker := time.NewTicker(pollSleep)
	defer ticker.Stop()

	startTime := time.Now()
	attempts := 0

	logger.Info("Waiting for transaction confirmation...",
		zap.String("tx_hash", hash.Hex()),
		zap.Duration("timeout", timeout),
		zap.Duration("poll_interval", pollSleep),
	)

	for {
		attempts++
		tx, isPending, err := client.TransactionByHash(ctx, hash)
		logger.Debug(
			"Transaction details",
			zap.String("tx_hash", hash.Hex()),
			zap.Any("tx", tx),
			zap.Bool("isPending", isPending),
			zap.Error(err),
		)

		receipt, err := client.TransactionReceipt(ctx, hash)

		if err != nil {
			if err.Error() != "not found" {
				logger.Error("Error while checking transaction receipt",
					zap.String("tx_hash", hash.Hex()),
					zap.Int("attempts", attempts),
					zap.Error(err),
				)
			} else {
				logger.Debug("Transaction not yet mined",
					zap.String("tx_hash", hash.Hex()),
					zap.Int("attempts", attempts),
					zap.Duration("elapsed", time.Since(startTime)),
				)
			}
		} else if receipt != nil {
			// Transaction found in a block
			logger.Info("Transaction confirmed",
				zap.String("tx_hash", hash.Hex()),
				zap.Uint64("block_number", receipt.BlockNumber.Uint64()),
				zap.Int("attempts", attempts),
				zap.Duration("elapsed", time.Since(startTime)),
				zap.Uint64("status", receipt.Status),
			)

			if receipt.Status == types.ReceiptStatusSuccessful {
				return receipt, nil
			}
			logger.Warn("Transaction failed",
				zap.String("tx_hash", hash.Hex()),
				zap.Uint64("block_number", receipt.BlockNumber.Uint64()),
			)
			return nil, fmt.Errorf("transaction failed: %s", hash.Hex())
		}

		select {
		case <-ctx.Done():
			logger.Error("Transaction confirmation timed out",
				zap.String("tx_hash", hash.Hex()),
				zap.Int("attempts", attempts),
				zap.Duration("elapsed", time.Since(startTime)),
				zap.Error(ctx.Err()),
			)

			return nil, fmt.Errorf("timed out waiting for transaction %s", hash.Hex())

		case <-ticker.C:
			continue
		}
	}
}
