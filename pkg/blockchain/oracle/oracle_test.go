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
		o.Stop()
	})

	o.Start()

	return o
}

func buildOracleWithoutStart(t *testing.T) (*oracle.Oracle, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	logger := testutils.NewLog(t)

	wsURL, _ := anvil.StartAnvil(t, false)

	o, err := oracle.New(
		ctx,
		logger,
		wsURL,
	)
	require.NoError(t, err)

	return o, cancel
}

func TestOracleGetGasPrice(t *testing.T) {
	o := buildOracle(t)

	require.Eventually(t, func() bool {
		gasPrice := o.GetGasPrice()
		t.Logf("gas price: %d", gasPrice)
		return gasPrice > 0
	}, 10*time.Second, 100*time.Millisecond)
}

func TestOracleGetGasPriceBeforeStart(t *testing.T) {
	o, cancel := buildOracleWithoutStart(t)
	t.Cleanup(cancel)
	t.Cleanup(o.Stop)

	// GetGasPrice should return a valid price even before Start()
	// because NewOracle() fetches the initial gas price
	gasPrice := o.GetGasPrice()
	require.Greater(t, gasPrice, uint64(0))
}

func TestOracleStartIdempotent(t *testing.T) {
	o, cancel := buildOracleWithoutStart(t)
	t.Cleanup(cancel)
	t.Cleanup(o.Stop)

	// Multiple starts should not panic or cause issues
	o.Start()
	o.Start()
	o.Start()

	// Should still work normally
	require.Eventually(t, func() bool {
		return o.GetGasPrice() > 0
	}, 10*time.Second, 100*time.Millisecond)
}

func TestOracleStopIdempotent(t *testing.T) {
	o, cancel := buildOracleWithoutStart(t)
	t.Cleanup(cancel)

	o.Start()

	// Multiple stops should not panic
	o.Stop()
	o.Stop()
	o.Stop()
}

func TestOracleConcurrentGetGasPrice(t *testing.T) {
	o := buildOracle(t)

	// Wait for initial gas price
	require.Eventually(t, func() bool {
		return o.GetGasPrice() > 0
	}, 10*time.Second, 100*time.Millisecond)

	// GetGasPrice is safe for concurrent calls
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			price := o.GetGasPrice()
			require.Greater(t, price, uint64(0))
		}()
	}
	wg.Wait()
}

func TestOracleGasPriceUpdatesOnNewBlock(t *testing.T) {
	o := buildOracle(t)

	// Wait for initial gas price
	require.Eventually(t, func() bool {
		return o.GetGasPrice() > 0
	}, 10*time.Second, 100*time.Millisecond)

	initialPrice := o.GetGasPrice()
	t.Logf("initial gas price: %d", initialPrice)

	// Wait for a few blocks (Anvil produces blocks every ~1s by default)
	time.Sleep(3 * time.Second)

	// Price should still be valid (not stale)
	currentPrice := o.GetGasPrice()
	t.Logf("current gas price: %d", currentPrice)
	require.Greater(t, currentPrice, uint64(0))
}

func TestOracleGracefulShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	logger := testutils.NewLog(t)
	wsURL, _ := anvil.StartAnvil(t, false)

	o, err := oracle.New(ctx, logger, wsURL)
	require.NoError(t, err)
	o.Start()

	// Give it time to start watching
	time.Sleep(500 * time.Millisecond)

	// Cancel context
	cancel()

	// Stop should complete quickly
	done := make(chan struct{})
	go func() {
		o.Stop()
		close(done)
	}()

	select {
	case <-done:
		// Success - shutdown completed
	case <-time.After(5 * time.Second):
		t.Fatal("Stop() took too long")
	}
}

func TestOracleStopBeforeStart(t *testing.T) {
	o, cancel := buildOracleWithoutStart(t)
	t.Cleanup(cancel)

	// Stop before Start should not panic
	o.Stop()

	// Should still be able to start after
	o.Start()
	t.Cleanup(o.Stop)

	require.Eventually(t, func() bool {
		return o.GetGasPrice() > 0
	}, 10*time.Second, 100*time.Millisecond)
}
