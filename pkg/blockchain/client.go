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

// ExecuteTransaction is a helper function that:
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
	// Enforce the timeout with a context so that slow requests get aborted if the function has
	// run out of time
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(timeout))
	defer cancel()

	// Ticker to track polling interval
	ticker := time.NewTicker(pollSleep)
	defer ticker.Stop()

	for {
		receipt, err := client.TransactionReceipt(ctx, hash)
		if err != nil {
			if err.Error() != "not found" {
				logger.Error("waiting for transaction", zap.String("hash", hash.String()))
			}
		} else if receipt != nil {
			if receipt.Status == types.ReceiptStatusSuccessful {
				return receipt, nil
			}
			return nil, fmt.Errorf("transaction failed")
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timed out")
		case <-ticker.C:
			continue
		}
	}
}
