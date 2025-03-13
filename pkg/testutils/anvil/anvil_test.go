package anvil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
)

func TestStartAnvil(t *testing.T) {
	anvilUrl, cleanup := StartAnvil(t, true)
	defer cleanup()

	client, err := blockchain.NewClient(context.Background(), anvilUrl)
	require.NoError(t, err)

	chainId, err := client.ChainID(context.Background())
	require.NoError(t, err)
	require.NotNil(t, chainId)
}

func TestCleanup(t *testing.T) {
	anvilUrl, cleanup := StartAnvil(t, false)

	client, err := blockchain.NewClient(context.Background(), anvilUrl)
	require.NoError(t, err)

	chainId, err := client.ChainID(context.Background())
	require.NoError(t, err)
	require.NotNil(t, chainId)

	cleanup()

	chainId, err = client.ChainID(context.Background())
	require.Error(t, err)
	require.Nil(t, chainId)
}

func TestStartConcurrent(t *testing.T) {
	// Start 10 anvil instances concurrently
	const numInstances = 10

	// Create channels to collect results
	type anvilInstance struct {
		url     string
		cleanup func()
	}
	results := make(chan anvilInstance, numInstances)

	// Launch goroutines to start anvil instances
	for i := 0; i < numInstances; i++ {
		go func() {
			url, cleanup := StartAnvil(t, false)
			results <- anvilInstance{url: url, cleanup: cleanup}
		}()
	}

	// Collect all instances
	instances := make([]anvilInstance, 0, numInstances)
	for range numInstances {
		instance := <-results
		instances = append(instances, instance)
		if instance.cleanup != nil {
			defer instance.cleanup()
		}
	}

	require.Len(t, instances, numInstances)

	// Verify all instances started successfully
	for i, instance := range instances {
		require.NotEmpty(t, instance.url, "Empty URL for anvil instance %d", i)

		// Verify we can connect to each instance
		client, err := blockchain.NewClient(context.Background(), instance.url)
		require.NoError(t, err, "Failed to connect to anvil instance %d", i)

		chainId, err := client.ChainID(context.Background())
		require.NoError(t, err, "Failed to get chain ID from anvil instance %d", i)
		require.NotNil(t, chainId, "Nil chain ID from anvil instance %d", i)
	}
}
