// Package oracle provides a blockchain oracle that executes functions on each new block.
//
// The oracle subscribes to new block headers via WebSocket and triggers registered
// callbacks whenever a new block is received. It works with any EVM-compatible chain
// and can be extended with custom functions that need to run on each block.
//
// # Chain Support
//
// The oracle works with any EVM chain that supports eth_subscribe("newHeads"):
//   - Ethereum mainnet and testnets
//   - Arbitrum One, Nova, and Orbit L3 chains
//   - Optimism, Base, and other OP Stack chains
//   - Any EVM-compatible chain with WebSocket support
//
// # Built-in Gas Price Tracking
//
// The oracle includes gas price tracking out of the box, with automatic chain detection:
//   - Arbitrum chains: Uses ArbGasInfo precompile (0x6C) for accurate L2 pricing
//   - Other chains: Uses eth_gasPrice RPC call
//
// # Extending with Custom Functions
//
// The oracle can be extended to run additional logic on each new block by adding
// callbacks to the watch() loop. Common use cases include:
//   - Monitoring contract state changes
//   - Tracking token balances or prices
//   - Triggering time-sensitive operations
//
// # Usage
//
//	oracle, err := oracle.NewOracle(ctx, logger, "wss://your-rpc-endpoint")
//	if err != nil {
//	    return err
//	}
//	oracle.Start()
//	defer oracle.Stop()
//
//	// Get gas price (non-blocking, returns cached value)
//	gasPrice := oracle.GetGasPrice()
//
// # Thread Safety
//
// All methods are safe for concurrent use. GetGasPrice() is lock-free and suitable
// for high-frequency calls in hot paths like transaction submission.
//
// # Staleness Handling
//
// If the cached gas price becomes stale (subscription dies, no blocks received),
// GetGasPrice() returns a safe default value (0.1 gwei) rather than blocking on an RPC call.
// Metrics are emitted when fallback values are used.
package oracle

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

type BlockchainOracle interface {
	Start()
	Stop()
	GetGasPrice() uint64
}

const (
	// 0.1 Gwei - Arbitrum Orbit's default.
	defaultGasPrice  = 100_000_000
	maxStaleDuration = 30 * time.Second
)

type Oracle struct {
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	logger    *zap.Logger
	ethClient *ethclient.Client
	chainID   int64

	gasPriceSource      GasPriceSource
	gasPriceLastUpdated atomic.Int64
	gasPrice            atomic.Uint64

	running atomic.Bool
}

func New(
	ctx context.Context,
	logger *zap.Logger,
	wsURL string,
) (oracle *Oracle, err error) {
	ethClient, err := ethclient.Dial(wsURL)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)

	defer func() {
		if err != nil {
			ethClient.Close()
			cancel()
		}
	}()

	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	logger = logger.Named(utils.BlockchainOracleLoggerName).With(
		utils.ChainIDField(chainID.Int64()),
	)

	gasPriceSource := gasPriceSourceDefault
	if isArbChain(ctx, ethClient) {
		gasPriceSource = gasPriceSourceArbMinimum
	}

	logger.Info("gas price source", zap.String("source", gasPriceSource.String()))

	oracle = &Oracle{
		ctx:            ctx,
		cancel:         cancel,
		logger:         logger,
		ethClient:      ethClient,
		chainID:        chainID.Int64(),
		gasPriceSource: gasPriceSource,
	}

	err = oracle.updateGasPrice()
	if err != nil {
		return nil, err
	}

	return oracle, nil
}

// GetGasPrice is a fast non-blocking way to get the current gas price.
// If the price is stale or 0, it will return a sane default.
func (o *Oracle) GetGasPrice() uint64 {
	if time.Since(time.Unix(o.gasPriceLastUpdated.Load(), 0)) > maxStaleDuration {
		o.logger.Debug("gas price stale, using default")
		metrics.EmitBlockchainGasPriceDefaultFallbackTotal(o.chainID)
		return defaultGasPrice
	}

	price := o.gasPrice.Load()
	if price == 0 {
		o.logger.Debug("gas price is 0, using default")
		metrics.EmitBlockchainGasPriceDefaultFallbackTotal(o.chainID)
		return defaultGasPrice
	}

	return price
}

func (o *Oracle) GetChainID() int64 {
	return o.chainID
}

func (o *Oracle) Start() {
	if o.running.Swap(true) {
		o.logger.Debug("oracle already running")
		return
	}

	tracing.GoPanicWrap(
		o.ctx,
		&o.wg,
		"blockchain-oracle",
		func(ctx context.Context) {
			o.watch()
		})
}

func (o *Oracle) Stop() {
	if !o.running.Swap(false) {
		o.logger.Debug("oracle not running")
		return
	}

	o.cancel()
	o.wg.Wait()
	o.ethClient.Close()
}

func (o *Oracle) watch() {
	headCh := make(chan *types.Header)
	defer close(headCh)

	sub, err := o.ethClient.SubscribeNewHead(o.ctx, headCh)
	if err != nil {
		o.logger.Error(
			"unexpected error while creating subscription",
			zap.Error(err),
		)
		return
	}
	defer sub.Unsubscribe()

	o.logger.Info("watching for new blocks")

	for {
		select {
		case <-o.ctx.Done():
			o.logger.Debug("shutting down")
			return

		case err := <-sub.Err():
			if err == nil {
				continue
			}

			o.logger.Error("subscription error, rebuilding", zap.Error(err))

			sub.Unsubscribe()

			sub, err = o.buildSubscriptionWithBackoff(headCh)
			if err != nil {
				o.logger.Fatal(
					"failed rebuilding subscription after max disconnect time",
					zap.Error(err),
				)
			}

			// Restart the select loop.
			continue

		case head, open := <-headCh:
			if !open {
				o.logger.Error("subscription channel closed, shutting down oracle")
				return
			}

			o.logger.Debug(
				"new block received",
				utils.BlockNumberField(head.Number.Uint64()),
				utils.HashField(head.Hash().Hex()),
			)

			if err := o.updateGasPrice(); err != nil {
				o.logger.Error("failed to update gas price", zap.Error(err))
			}
		}
	}
}

func (o *Oracle) updateGasPrice() error {
	gasPrice, err := getGasPrice(o.ctx, o.ethClient, o.gasPriceSource)
	if err != nil {
		o.logger.Error("failed to get gas price", zap.Error(err))
		return err
	}

	now := time.Now().Unix()
	o.gasPriceLastUpdated.Store(now)
	o.gasPrice.Store(gasPrice.Uint64())

	metrics.EmitBlockchainGasPrice(o.chainID, gasPrice.Uint64())
	metrics.EmitBlockchainGasPriceUpdatesTotal(o.chainID)
	metrics.EmitBlockchainGasPriceLastUpdateTimestamp(o.chainID, now)

	return nil
}

func (o *Oracle) buildSubscriptionWithBackoff(
	innerSubCh chan *types.Header,
) (sub ethereum.Subscription, err error) {
	rebuildOperation := func() (ethereum.Subscription, error) {
		sub, err = o.ethClient.SubscribeNewHead(o.ctx, innerSubCh)
		return sub, err
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.InitialInterval = 1 * time.Second

	sub, err = backoff.Retry(
		o.ctx,
		rebuildOperation,
		backoff.WithBackOff(expBackoff),
	)
	if err != nil {
		o.logger.Error(
			"failed to rebuild subscription, closing",
			zap.Error(err),
		)
		return nil, err
	}

	o.logger.Info("subscription rebuilt")

	return sub, nil
}
