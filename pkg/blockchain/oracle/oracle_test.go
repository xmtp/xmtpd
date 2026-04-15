package oracle_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain/oracle"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
)

func buildOracle(t *testing.T, opts ...oracle.Option) *oracle.Oracle {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	logger := testutils.NewLog(t)

	wsURL, _ := anvil.StartAnvil(t, false)

	o, err := oracle.New(
		ctx,
		logger,
		wsURL,
		opts...,
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		o.Close()
	})

	return o
}

func TestOracleGetGasPrice(t *testing.T) {
	o := buildOracle(t)

	// First GetGasPrice call always fetches because lastUpdated is zero,
	// so no wall-clock wait is required.
	gasPrice := o.GetGasPrice()
	require.Positive(t, gasPrice)
}

func TestOracleGetGasPriceAfterNew(t *testing.T) {
	o := buildOracle(t)

	// GetGasPrice should return a valid price immediately after New()
	// because it fetches lazily on first call
	gasPrice := o.GetGasPrice()
	require.Positive(t, gasPrice)
}

func TestOracleConcurrentGetGasPrice(t *testing.T) {
	o := buildOracle(t)

	// Get initial gas price to warm up
	initialPrice := o.GetGasPrice()
	require.Positive(t, initialPrice)

	// GetGasPrice is safe for concurrent calls
	var wg sync.WaitGroup
	for range 100 {
		wg.Go(func() {
			price := o.GetGasPrice()
			require.Positive(t, price)
		})
	}
	wg.Wait()
}

func TestOracleGetGasPriceRandom(t *testing.T) {
	// Force every GetGasPrice call to refetch by configuring a zero stale
	// window. singleflight still coalesces concurrent fetches, so this
	// exercises the racy "fetch under load" path without any wall-clock sleep.
	o := buildOracle(t, oracle.WithMaxStaleDuration(0))

	initialPrice := o.GetGasPrice()
	require.Positive(t, initialPrice)

	var wg sync.WaitGroup
	for range 100 {
		wg.Go(func() {
			price := o.GetGasPrice()
			require.Positive(t, price)
		})
	}
	wg.Wait()
}

func TestOracleGasPriceRefreshesAfterStaleness(t *testing.T) {
	// A zero stale window forces the second GetGasPrice call to fetch a
	// fresh value, deterministically exercising the refresh path.
	o := buildOracle(t, oracle.WithMaxStaleDuration(0))

	// Get initial gas price
	initialPrice := o.GetGasPrice()
	t.Logf("initial gas price: %d", initialPrice)
	require.Positive(t, initialPrice)

	// Price should still be valid after refresh
	currentPrice := o.GetGasPrice()
	t.Logf("current gas price: %d", currentPrice)
	require.Positive(t, currentPrice)
}

func TestOracleCloseIdempotent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	logger := testutils.NewLog(t)
	wsURL, _ := anvil.StartAnvil(t, false)

	o, err := oracle.New(ctx, logger, wsURL)
	require.NoError(t, err)

	// Multiple closes should not panic
	o.Close()
	o.Close()
	o.Close()
}

func TestOracleGracefulShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	logger := testutils.NewLog(t)
	wsURL, _ := anvil.StartAnvil(t, false)

	o, err := oracle.New(ctx, logger, wsURL)
	require.NoError(t, err)

	// Get a gas price to ensure connection is established
	price := o.GetGasPrice()
	require.Positive(t, price)

	// Cancel context
	cancel()

	// Close should complete quickly
	done := make(chan struct{})
	go func() {
		o.Close()
		close(done)
	}()

	select {
	case <-done:
		// Success - shutdown completed
	case <-time.After(5 * time.Second):
		t.Fatal("Close() took too long")
	}
}
