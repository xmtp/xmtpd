package blockchain

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xmtp/xmtpd/pkg/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"go.uber.org/zap"
)

var ErrTxFailed = errors.New("transaction failed")

const (
	executeTxMaxRetries  = 3
	executeTxInitialWait = 250 * time.Millisecond
	waitTxTimeout        = 20 * time.Second
	waitTxPollSleep      = 250 * time.Millisecond
)

type WebsocketClientOption func(*websocketClientConfig)

type websocketClientConfig struct {
	tcpDialer *net.Dialer
}

func WithKeepAliveConfig(cfg net.KeepAliveConfig) WebsocketClientOption {
	return func(c *websocketClientConfig) {
		c.tcpDialer = &net.Dialer{
			Timeout:         10 * time.Second,
			KeepAliveConfig: cfg,
		}
	}
}

// NewWebsocketClient creates a new websocket client that can be configured with dialer options.
// It's used mostly for subscriptions.
func NewWebsocketClient(
	ctx context.Context,
	wsURL string,
	opts ...WebsocketClientOption,
) (*ethclient.Client, error) {
	config := &websocketClientConfig{
		tcpDialer: &net.Dialer{
			Timeout: 10 * time.Second,
			KeepAliveConfig: net.KeepAliveConfig{
				Enable:   true,
				Idle:     4 * time.Second,
				Interval: 2 * time.Second,
				Count:    10,
			},
		},
	}

	for _, opt := range opts {
		opt(config)
	}

	dialer := websocket.Dialer{
		NetDialContext: config.tcpDialer.DialContext,
		Proxy:          http.ProxyFromEnvironment,
	}

	rpcClient, err := rpc.DialOptions(ctx, wsURL, rpc.WithWebsocketDialer(dialer))
	if err != nil {
		return nil, err
	}

	return ethclient.NewClient(rpcClient), nil
}

// NewRPCClient creates a new RPC client that can be used for JSON-RPC calls.
// RPC providers usually implement middleware with optimizations for HTTP JSON-RPC requests.
func NewRPCClient(ctx context.Context, rpcURL string) (*ethclient.Client, error) {
	return ethclient.DialContext(ctx, rpcURL)
}

// isUnderpricedError checks if an error is a transient "underpriced" error
// from load-balanced RPCs that should be retried.
func isUnderpricedError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "underpriced")
}

// ExecuteTransaction is a helper function that:
// - checks if the sender has enough balance.
// - executes a transaction with retry on transient underpriced errors.
// - waits for it to be mined.
// - if the transaction fails, it tries to get the error code from the transaction receipt.
// - processes the event logs.
func ExecuteTransaction(
	ctx context.Context,
	signer TransactionSigner,
	logger *zap.Logger,
	client *ethclient.Client,
	txFunc func(*bind.TransactOpts) (*types.Transaction, error),
	eventParser func(*types.Log) (any, error),
	logHandler func(any),
) ProtocolError {
	if signer == nil {
		return NewBlockchainError(errors.New("no signer provided"))
	}

	from := signer.FromAddress()

	balance, err := client.BalanceAt(ctx, from, nil)
	if err != nil {
		return NewBlockchainError(fmt.Errorf("failed to check balance: %w", err))
	}

	if balance.Cmp(big.NewInt(0)) == 0 {
		return NewBlockchainError(fmt.Errorf("account %s has zero balance", from.Hex()))
	}

	logger.Debug(
		"sender balance",
		utils.AddressField(from.Hex()),
		utils.BalanceField(balance.String()),
	)

	opts := &bind.TransactOpts{
		Context: ctx,
		From:    from,
		Signer:  signer.SignerFunc(),
	}

	tx, protocolErr := executeTransaction(ctx, logger, opts, txFunc)
	if protocolErr != nil {
		return protocolErr
	}

	receipt, protocolErr := waitForTransaction(
		ctx,
		logger,
		client,
		tx.Hash(),
	)
	if protocolErr != nil {
		return protocolErr
	}

	if eventParser != nil && logHandler != nil {
		for _, log := range receipt.Logs {
			event, err := eventParser(log)
			if err != nil {
				continue
			}
			logHandler(event)
		}
	}

	return nil
}

func executeTransaction(
	ctx context.Context,
	logger *zap.Logger,
	opts *bind.TransactOpts,
	txFunc func(*bind.TransactOpts) (*types.Transaction, error),
) (*types.Transaction, ProtocolError) {
	var (
		tx  *types.Transaction
		err error
	)

	for attempt := 0; attempt <= executeTxMaxRetries; attempt++ {
		if ctx.Err() != nil {
			return nil, NewBlockchainError(ctx.Err())
		}

		tx, err = txFunc(opts)
		if err == nil {
			return tx, nil
		}

		if isUnderpricedError(err) {
			backoff := executeTxInitialWait * (1 << attempt)

			logger.Warn(
				"retryable transaction error, backing off",
				zap.Error(err),
				zap.Int("attempt", attempt+1),
				zap.Duration("backoff", backoff),
			)

			if attempt < executeTxMaxRetries {
				utils.RandomSleep(ctx, backoff)
			}

			continue
		}

		return nil, NewBlockchainError(err)
	}

	return nil, NewBlockchainError(fmt.Errorf("max retries reached executing transaction: %w", err))
}

// waitForTransaction waits for the given transaction hash to have been submitted to the chain and soft confirmed.
func waitForTransaction(
	ctx context.Context,
	logger *zap.Logger,
	client *ethclient.Client,
	hash common.Hash,
) (*types.Receipt, ProtocolError) {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(waitTxTimeout))
	defer cancel()

	ticker := time.NewTicker(waitTxPollSleep)
	defer ticker.Stop()

	var (
		receipt *types.Receipt
		err     error
	)

	for {
		if receipt == nil {
			receipt, err = client.TransactionReceipt(ctx, hash)
			if err != nil {
				if isTransactionNotFound(err) {
					logger.Debug("waiting for transaction", utils.HashField(hash.String()))
				} else {
					return nil, NewBlockchainError(err)
				}
			}
		}

		if receipt != nil {
			if receipt.Status == types.ReceiptStatusSuccessful {
				return receipt, nil
			}

			if receipt.Status == types.ReceiptStatusFailed {
				tx, _, err := client.TransactionByHash(ctx, hash)
				if err != nil {
					// Redundant check to handle load-balanced RPC backends that may not
					// have the tx available for tracing yet.
					if isTransactionNotFound(err) {
						logger.Debug("waiting for transaction", utils.HashField(hash.String()))
					} else {
						return receipt, NewBlockchainError(
							fmt.Errorf("failed to get transaction %s: %w", hash.Hex(), err),
						)
					}
				} else {
					protocolErr := getProtocolError(ctx, client, tx)
					return receipt, protocolErr
				}
			}
		}

		select {
		case <-ctx.Done():
			return nil, NewBlockchainError(errors.New("timed out"))
		case <-ticker.C:
			continue
		}
	}
}

type traceTransactionResult struct {
	Output string `json:"output"`
}

type traceTransactionConfig struct {
	Tracer string `json:"tracer"`
}

// getProtocolError uses debug_traceTransaction with callTracer to extract the revert reason.
// This approach works reliably on Arbitrum Orbit L3 where eth_call may not return revert data.
func getProtocolError(
	ctx context.Context,
	client *ethclient.Client,
	tx *types.Transaction,
) ProtocolError {
	traceCfg := traceTransactionConfig{
		Tracer: "callTracer",
	}

	var (
		traceOut traceTransactionResult
		ticker   = time.NewTicker(100 * time.Millisecond)
	)

	defer ticker.Stop()

	for {
		err := client.Client().
			CallContext(ctx, &traceOut, "debug_traceTransaction", tx.Hash(), &traceCfg)
		if err == nil {
			if traceOut.Output == "" {
				return NewBlockchainError(
					fmt.Errorf("transaction %s reverted without reason", tx.Hash().Hex()),
				)
			}

			return NewBlockchainError(errors.New(traceOut.Output))
		}

		if !isTransactionNotFound(err) {
			return NewBlockchainError(
				fmt.Errorf("failed to trace transaction %s: %w", tx.Hash().Hex(), err),
			)
		}

		select {
		case <-ctx.Done():
			return NewBlockchainError(errors.New("timed out"))
		case <-ticker.C:
		}
	}
}

// isTransactionNotFound checks if an error is a transaction not found,
// from a Geth and Reth clients.
//   - reth: https://github.com/paradigmxyz/reth/blob/main/crates/rpc/rpc-eth-types/src/error/mod.rs#L136
//   - geth: https://github.com/ethereum/go-ethereum/blob/master/interfaces.go#L32
func isTransactionNotFound(err error) bool {
	return errors.Is(err, ethereum.NotFound) ||
		strings.Contains(err.Error(), "transaction not found")
}
