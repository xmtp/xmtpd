package anvil

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
)

func TestStartAnvil(t *testing.T) {
	anvilUrl := StartAnvil(t, true)

	client, err := blockchain.NewClient(context.Background(), blockchain.WithWebSocketURL(anvilUrl))
	require.NoError(t, err)

	chainId, err := client.ChainID(context.Background())
	require.NoError(t, err)
	require.NotNil(t, chainId)
}

func TestStartConcurrent(t *testing.T) {
	const numInstances = 10

	for i := 0; i < numInstances; i++ {
		i := i // capture loop variable
		t.Run(fmt.Sprintf("instance-%d", i), func(t *testing.T) {
			t.Parallel()

			url := StartAnvil(t, false)
			require.NotEmpty(t, url, "Empty URL for anvil instance %d", i)

			client, err := blockchain.NewClient(
				context.Background(),
				blockchain.WithWebSocketURL(url),
			)
			require.NoError(t, err, "Failed to connect to anvil instance %d", i)

			chainId, err := client.ChainID(context.Background())
			require.NoError(t, err, "Failed to get chain ID from anvil instance %d", i)
			require.NotNil(t, chainId, "Nil chain ID from anvil instance %d", i)
		})
	}
}
