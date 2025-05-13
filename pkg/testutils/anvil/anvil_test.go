package anvil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
)

func TestStartAnvil(t *testing.T) {
	anvilUrl := StartAnvil(t, true)

	client, err := blockchain.NewClient(context.Background(), anvilUrl)
	require.NoError(t, err)

	chainId, err := client.ChainID(context.Background())
	require.NoError(t, err)
	require.NotNil(t, chainId)
}

func TestStartConcurrent(t *testing.T) {
	// Start 10 anvil instances concurrently
	const numInstances = 10

	// Create channels to collect results
	type anvilInstance struct {
		url string
	}
	results := make(chan anvilInstance, numInstances)

	// Launch goroutines to start anvil instances
	for i := 0; i < numInstances; i++ {
		go func() {
			url := StartAnvil(t, false)
			results <- anvilInstance{url: url}
		}()
	}

	// Collect all instances
	instances := make([]anvilInstance, 0, numInstances)
	for range numInstances {
		instance := <-results
		instances = append(instances, instance)
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
