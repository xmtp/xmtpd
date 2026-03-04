package oracle_test

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain/oracle"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
)

func buildOracle(t *testing.T) *oracle.Oracle {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	logger := testutils.NewLog(t)

	wsURL, _ := anvil.StartAnvil(t, false)

	o, err := oracle.New(
		ctx,
		logger,
		wsURL,
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		o.Close()
	})

	return o
}

func TestOracleGetGasPrice(t *testing.T) {
	o := buildOracle(t)

	// By forcing waiting, we ensure that the gas price is fetched from the blockchain.
	time.Sleep(500 * time.Millisecond)
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
	o := buildOracle(t)

	initialPrice := o.GetGasPrice()
	require.Positive(t, initialPrice)

	var wg sync.WaitGroup
	for range 100 {
		wg.Go(func() {
			// Force the oracle to fetch a new gas price for some goroutines.
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
			price := o.GetGasPrice()
			require.Positive(t, price)
		})
	}
	wg.Wait()
}

func TestOracleGasPriceRefreshesAfterStaleness(t *testing.T) {
	o := buildOracle(t)

	// Get initial gas price
	initialPrice := o.GetGasPrice()
	t.Logf("initial gas price: %d", initialPrice)
	require.Positive(t, initialPrice)

	// Wait for staleness (250ms + buffer)
	time.Sleep(300 * time.Millisecond)

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
