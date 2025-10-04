package blockchain

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xmtp/xmtpd/pkg/metrics"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"go.uber.org/zap"
)

var (
	ErrTxFailed        = fmt.Errorf("transaction failed")
	ErrTxFailedNoError = fmt.Errorf("transaction failed but no error found")
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

// ExecuteTransaction is a helper function that:
// - checks if the sender has enough balance.
// - executes a transaction.
// - waits for it to be mined.
// - if the transaction fails, it tries to get the error code from the transaction receipt.
// - processes the event logs.
func ExecuteTransaction(
	ctx context.Context,
	signer TransactionSigner,
	logger *zap.Logger,
	client *ethclient.Client,
	txFunc func(*bind.TransactOpts) (*types.Transaction, error),
	eventParser func(*types.Log) (interface{}, error),
	logHandler func(interface{}),
) ProtocolError {
	if signer == nil {
		return NewBlockchainError(fmt.Errorf("no signer provided"))
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
		"Sender balance",
		zap.String("address", from.Hex()),
		zap.String("balance", balance.String()),
	)

	opts := &bind.TransactOpts{
		Context:  ctx,
		From:     from,
		Signer:   signer.SignerFunc(),
		GasLimit: 300_000,
	}

	// transactions that are not simulated will always return a tx.Hash().
	// The error will be returned if the transaction fails to be mined.
	tx, err := txFunc(opts)
	if err != nil {
		return NewBlockchainError(err)
	}

	receipt, err := WaitForTransaction(
		ctx,
		logger,
		client,
		60*time.Second,
		250*time.Millisecond,
		tx.Hash(),
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrTxFailed):
			code, err := tryGetTxError(
				ctx,
				client,
				tx.Hash(),
				60*time.Second,
				250*time.Millisecond,
			)
			if err != nil {
				return NewBlockchainError(err)
			}
			return NewBlockchainError(errors.New(code))

		default:
			return NewBlockchainError(err)
		}
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

// WaitForTransaction waits for the given transaction hash to have been submitted to the chain and soft confirmed.
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

	now := time.Now()
	defer func() {
		metrics.EmitBlockchainWaitForTransaction(time.Since(now).Seconds())
	}()

	for {
		receipt, err := client.TransactionReceipt(ctx, hash)

		switch {
		case err == nil && receipt != nil:
			if receipt.Status == types.ReceiptStatusSuccessful {
				return receipt, nil
			}

			if receipt.Status == types.ReceiptStatusFailed {
				return receipt, ErrTxFailed
			}

		case err != nil:
			if errors.Is(err, ethereum.NotFound) {
				logger.Debug("waiting for transaction", zap.String("hash", hash.String()))
			} else {
				return nil, ErrTxFailed
			}
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timed out")
		case <-ticker.C:
			continue
		}
	}
}

func tryGetTxError(
	ctx context.Context,
	client *ethclient.Client,
	hash common.Hash,
	timeout time.Duration,
	pollSleep time.Duration,
) (string, error) {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(timeout))
	defer cancel()

	ticker := time.NewTicker(pollSleep)
	defer ticker.Stop()

	now := time.Now()
	defer func() {
		metrics.EmitBlockchainWaitForTransaction(time.Since(now).Seconds())
	}()

	type traceTransactionResult struct {
		Output string `json:"output"`
	}

	// https://geth.ethereum.org/docs/interacting-with-geth/rpc/ns-debug#debugtracetransaction
	type traceConfig struct {
		Tracer string `json:"tracer"`
	}

	tracerCfg := traceConfig{Tracer: "callTracer"}

	for {
		var result traceTransactionResult

		err := client.Client().
			CallContext(ctx, &result, "debug_traceTransaction", hash, tracerCfg)

		switch {
		case err == nil && result.Output != "":
			return result.Output, nil

		case err != nil:
			if err.Error() != "transaction not found" {
				return "", err
			}
		}

		select {
		case <-ctx.Done():
			return "", fmt.Errorf("timed out")
		case <-ticker.C:
			continue
		}
	}
}
