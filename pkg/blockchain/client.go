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
		Context:  ctx,
		From:     from,
		Signer:   signer.SignerFunc(),
		GasLimit: 5_000_000,
	}

	// transactions that are not simulated will always return a tx.Hash().
	// The error will be returned if the transaction fails to be mined.
	tx, err := txFunc(opts)
	if err != nil {
		return NewBlockchainError(err)
	}

	receipt, protocolErr := WaitForTransaction(
		ctx,
		logger,
		client,
		20*time.Second,
		250*time.Millisecond,
		tx.Hash(),
	)
	if protocolErr != nil {
		return protocolErr
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
) (receipt *types.Receipt, err ProtocolError) {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(timeout))
	defer cancel()

	ticker := time.NewTicker(pollSleep)
	defer ticker.Stop()

	for {
		receipt, err := client.TransactionReceipt(ctx, hash)
		if err != nil {
			if errors.Is(err, ethereum.NotFound) {
				logger.Debug("waiting for transaction", utils.HashField(hash.String()))
			} else {
				return nil, NewBlockchainError(err)
			}
		}

		if receipt != nil {
			if receipt.Status == types.ReceiptStatusSuccessful {
				return receipt, nil
			}

			if receipt.Status == types.ReceiptStatusFailed {
				tx, _, err := client.TransactionByHash(ctx, hash)
				if err != nil {
					return receipt, NewBlockchainError(
						fmt.Errorf("failed to get transaction %s: %w", hash.Hex(), err),
					)
				}

				protocolErr := getProtocolError(ctx, client, tx)
				return receipt, protocolErr
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

	var traceOut traceTransactionResult

	err := client.Client().
		CallContext(ctx, &traceOut, "debug_traceTransaction", tx.Hash(), &traceCfg)
	if err != nil {
		return NewBlockchainError(
			fmt.Errorf("failed to trace transaction %s: %w", tx.Hash().Hex(), err),
		)
	}

	if traceOut.Output == "" {
		return NewBlockchainError(
			fmt.Errorf("transaction %s reverted without reason", tx.Hash().Hex()),
		)
	}

	return NewBlockchainError(errors.New(traceOut.Output))
}
