// Package oracle provides a blockchain oracle for gas price queries.
//
// The oracle fetches gas prices on-demand with a short TTL, using request
// coalescing to minimize RPC calls when multiple goroutines request the
// price simultaneously.
//
// # Chain Support
//
// The oracle works with any EVM chain:
//   - Ethereum mainnet and testnets
//   - Arbitrum One, Nova, and Orbit L3 chains
//   - Optimism, Base, and other OP Stack chains
//   - Any EVM-compatible chain
//
// # Built-in Gas Price Tracking
//
// The oracle includes gas price tracking with automatic chain detection:
//   - Arbitrum chains: Uses ArbGasInfo precompile (0x6C) for accurate L2 pricing
//   - Other chains: Uses eth_gasPrice RPC call
//
// # Usage
//
//	oracle, err := oracle.New(ctx, logger, "wss://your-rpc-endpoint")
//	if err != nil {
//	    return err
//	}
//	defer oracle.Close()
//
//	// Get gas price (fetches if stale, coalesces concurrent requests)
//	gasPrice := oracle.GetGasPrice()
//
// # Thread Safety
//
// All methods are safe for concurrent use. GetGasPrice() uses singleflight
// to coalesce concurrent requests, ensuring only one RPC call is made when
// multiple goroutines request the price simultaneously.
//
// # Staleness Handling
//
// Gas prices are cached for 250ms. If the cached value is stale, GetGasPrice()
// fetches a fresh value. If the fetch fails, it returns a safe default (0.1 gwei).
package oracle

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

type BlockchainOracle interface {
	GetGasPrice() int64
	Close()
}

const (
	// 0.1 Gwei - Arbitrum Orbit's default.
	defaultArbGasPrice       int64         = 100_000_000
	defaultEVMGasPrice       int64         = 10_000_000_000
	gasPriceBufferPercent    int64         = 10
	gasPriceMaxStaleDuration time.Duration = 250 * time.Millisecond
)

type Oracle struct {
	ctx       context.Context
	cancel    context.CancelFunc
	logger    *zap.Logger
	ethClient *ethclient.Client
	chainID   int64

	gasPriceDefaultWei  int64
	gasPriceSource      GasPriceSource
	gasPriceLastUpdated atomic.Int64
	gasPrice            atomic.Int64

	sfGroup singleflight.Group
}

func New(
	ctx context.Context,
	logger *zap.Logger,
	wsURL string,
) (*Oracle, error) {
	ethClient, err := ethclient.Dial(wsURL)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)

	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		ethClient.Close()
		cancel()
		return nil, err
	}

	logger = logger.Named(utils.BlockchainOracleLoggerName).With(
		utils.ChainIDField(chainID.Int64()),
	)

	var (
		gasPriceSource     = gasPriceSourceDefault
		gasPriceDefaultWei = defaultEVMGasPrice
	)

	if isArbChain(ctx, ethClient) {
		gasPriceSource = gasPriceSourceArbMinimum
		gasPriceDefaultWei = defaultArbGasPrice
	}

	logger.Info(
		"new blockchain oracle",
		zap.String("gas_price_source", gasPriceSource.String()),
		zap.Int64("gas_price_wei_default", gasPriceDefaultWei),
	)

	oracle := &Oracle{
		ctx:                ctx,
		cancel:             cancel,
		logger:             logger,
		ethClient:          ethClient,
		chainID:            chainID.Int64(),
		gasPriceSource:     gasPriceSource,
		gasPriceDefaultWei: gasPriceDefaultWei,
	}

	return oracle, nil
}

func (o *Oracle) Close() {
	o.cancel()
	o.ethClient.Close()
}

// GetGasPrice returns the current gas price.
// If the cached price is stale (older than 250ms), it fetches a fresh value.
// Concurrent requests are coalesced to minimize RPC calls.
// If fetching fails, returns a safe default value (0.1 gwei).
func (o *Oracle) GetGasPrice() int64 {
	if !o.isStale() {
		price := o.gasPrice.Load()
		if price > 0 {
			return price
		}
	}

	ctxwc, cancel := context.WithTimeout(o.ctx, 2*time.Second)
	defer cancel()

	price, err := o.fetchGasPrice(ctxwc)
	if err != nil {
		o.logger.Error("failed to fetch gas price", zap.Error(err))
		metrics.EmitBlockchainGasPriceDefaultFallbackTotal(o.chainID)
		return o.gasPriceDefaultWei
	}

	return price
}

func (o *Oracle) fetchGasPrice(ctx context.Context) (int64, error) {
	result, err, _ := o.sfGroup.Do("fetch-gas-price", func() (interface{}, error) {
		if err := o.updateGasPrice(ctx); err != nil {
			return o.gasPriceDefaultWei, err
		}
		return o.gasPrice.Load(), nil
	})

	return result.(int64), err
}

func (o *Oracle) isStale() bool {
	lastUpdated := o.gasPriceLastUpdated.Load()

	if lastUpdated == 0 {
		return true
	}

	return time.Since(time.UnixMilli(lastUpdated)) > gasPriceMaxStaleDuration
}

func (o *Oracle) updateGasPrice(ctx context.Context) error {
	gasPrice, err := getGasPrice(ctx, o.ethClient, o.gasPriceSource)
	if err != nil {
		o.logger.Error("failed to get gas price", zap.Error(err))
		return err
	}

	// Apply buffer to the gas price.
	bufferedPrice := gasPrice.Int64() + (gasPrice.Int64() * gasPriceBufferPercent / 100)

	now := time.Now()
	o.gasPriceLastUpdated.Store(now.UnixMilli())
	o.gasPrice.Store(bufferedPrice)

	metrics.EmitBlockchainGasPrice(o.chainID, gasPrice.Uint64())
	metrics.EmitBlockchainGasPriceUpdatesTotal(o.chainID)
	metrics.EmitBlockchainGasPriceLastUpdateTimestamp(o.chainID, now.Unix())

	return nil
}
