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
	require.Greater(t, gasPrice, int64(0))
}

func TestOracleGetGasPriceAfterNew(t *testing.T) {
	o := buildOracle(t)

	// GetGasPrice should return a valid price immediately after New()
	// because it fetches lazily on first call
	gasPrice := o.GetGasPrice()
	require.Greater(t, gasPrice, int64(0))
}

func TestOracleConcurrentGetGasPrice(t *testing.T) {
	o := buildOracle(t)

	// Get initial gas price to warm up
	initialPrice := o.GetGasPrice()
	require.Greater(t, initialPrice, int64(0))

	// GetGasPrice is safe for concurrent calls
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			price := o.GetGasPrice()
			require.Greater(t, price, int64(0))
		}()
	}
	wg.Wait()
}

func TestOracleGetGasPriceRandom(t *testing.T) {
	o := buildOracle(t)

	initialPrice := o.GetGasPrice()
	require.Greater(t, initialPrice, int64(0))

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Force the oracle to fetch a new gas price for some goroutines.
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
			price := o.GetGasPrice()
			require.Greater(t, price, int64(0))
		}()
	}
	wg.Wait()
}

func TestOracleGasPriceRefreshesAfterStaleness(t *testing.T) {
	o := buildOracle(t)

	// Get initial gas price
	initialPrice := o.GetGasPrice()
	t.Logf("initial gas price: %d", initialPrice)
	require.Greater(t, initialPrice, int64(0))

	// Wait for staleness (250ms + buffer)
	time.Sleep(300 * time.Millisecond)

	// Price should still be valid after refresh
	currentPrice := o.GetGasPrice()
	t.Logf("current gas price: %d", currentPrice)
	require.Greater(t, currentPrice, int64(0))
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
	require.Greater(t, price, int64(0))

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
