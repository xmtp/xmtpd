package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/xmtp/xmtpd/pkg/metrics"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

func NewClient(ctx context.Context, wsUrl string) (*ethclient.Client, error) {
	return ethclient.DialContext(ctx, wsUrl)
}

// ExecuteTransaction is a helper function that:
// - estimates the gas required for the transaction
// - checks if the sender has enough balance to cover the gas cost
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

	from := signer.FromAddress()

	// Step 1: Check balance before sending.
	balance, err := client.BalanceAt(ctx, from, nil)
	if err != nil {
		return fmt.Errorf("failed to check balance: %w", err)
	}
	if balance.Cmp(big.NewInt(0)) == 0 {
		return fmt.Errorf("account %s has zero balance", from.Hex())
	}

	logger.Debug(
		"Sender balance",
		zap.String("address", from.Hex()),
		zap.String("balance", balance.String()),
	)

	// Step 2: Prepare dry-run tx for gas estimation.
	opts := &bind.TransactOpts{
		Context: ctx,
		From:    from,
		Signer:  signer.SignerFunc(),
		NoSend:  true,
	}

	dryRunTx, err := txFunc(opts)
	if err != nil {
		return fmt.Errorf("failed to simulate tx (NoSend=true): %w", err)
	}

	// Step 3: Estimate gas using ethclient.EstimateGas.
	msg := ethereum.CallMsg{
		From:  from,
		To:    dryRunTx.To(),
		Data:  dryRunTx.Data(), // Not directly available — see note below
		Value: dryRunTx.Value(),
	}

	// Step 4: Fallback for GasPrice.
	gasPrice := dryRunTx.GasPrice()
	if gasPrice == nil {
		gasPrice, err = client.SuggestGasPrice(ctx)
		if err != nil {
			return fmt.Errorf("failed to get gas price: %w", err)
		}
	}

	estimatedGas, err := client.EstimateGas(ctx, msg)
	if err != nil {
		return fmt.Errorf("gas estimation failed: %w", err)
	}

	logger.Debug(
		"Gas estimation",
		zap.String("address", from.Hex()),
		zap.Uint64("gas", estimatedGas),
	)

	// Step 5: Check for balance sufficiency.
	required := new(big.Int).Mul(big.NewInt(int64(estimatedGas)), gasPrice)
	if balance.Cmp(required) < 0 {
		return fmt.Errorf(
			"insufficient funds: need %s, have %s",
			required.String(),
			balance.String(),
		)
	}

	// Step 6: Send the real tx.
	opts.NoSend = false
	opts.GasLimit = estimatedGas
	opts.GasPrice = gasPrice

	tx, err := txFunc(opts)
	if err != nil {
		return fmt.Errorf("failed to send tx: %w", err)
	}

	// Step 7: Wait for receipt.
	receipt, err := WaitForTransaction(
		ctx,
		logger,
		client,
		10*time.Second,
		250*time.Millisecond,
		tx.Hash(),
	)
	if err != nil {
		return err
	}

	// Step 8: Parse and handle logs.
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

	now := time.Now()
	defer func() {
		metrics.EmitBlockchainWaitForTransaction(time.Since(now).Seconds())
	}()

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
