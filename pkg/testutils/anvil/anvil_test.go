package anvil

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
)

func TestStartAnvil(t *testing.T) {
	_, anvilRPCURL := StartAnvil(t, true)

	client, err := blockchain.NewRPCClient(
		context.Background(),
		anvilRPCURL,
	)
	require.NoError(t, err)

	chainID, err := client.ChainID(context.Background())
	require.NoError(t, err)
	require.NotNil(t, chainID)
}

func TestStartConcurrent(t *testing.T) {
	const numInstances = 10

	for i := 0; i < numInstances; i++ {
		i := i // capture loop variable
		t.Run(fmt.Sprintf("instance-%d", i), func(t *testing.T) {
			t.Parallel()

			wsURL, rpcURL := StartAnvil(t, false)
			require.NotEmpty(t, wsURL, "Empty wsURL for anvil instance %d", i)
			require.NotEmpty(t, rpcURL, "Empty rpcURL for anvil instance %d", i)

			client, err := blockchain.NewRPCClient(
				context.Background(),
				rpcURL,
			)
			require.NoError(t, err, "Failed to connect to anvil instance %d", i)

			chainID, err := client.ChainID(context.Background())
			require.NoError(t, err, "Failed to get chain ID from anvil instance %d", i)
			require.NotNil(t, chainID, "Nil chain ID from anvil instance %d", i)
		})
	}
}
