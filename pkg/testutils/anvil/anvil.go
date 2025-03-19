package anvil

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	networkTestUtils "github.com/xmtp/xmtpd/pkg/testutils/network"
)

func waitForAnvil(t *testing.T, url string) {
	backgroundCtx := context.Background()
	// Create a client to connect to the Anvil instance
	client, err := blockchain.NewClient(backgroundCtx, url)
	require.NoError(t, err)

	// Try to get the chain ID to verify the connection is working
	// This will fail if Anvil is not ready yet
	ctx, cancel := context.WithTimeout(backgroundCtx, 5*time.Second)
	defer cancel()

	// Poll until we can successfully get the chain ID or timeout
	var chainID *big.Int
	for {
		select {
		case <-ctx.Done():
			t.Fatalf("Timed out waiting for Anvil to start: %v", ctx.Err())
			return
		default:
			chainID, err = client.ChainID(ctx)
			if err == nil && chainID != nil {
				// Successfully connected to Anvil
				return
			}
			// Wait a bit before trying again
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// Start an ephemeral anvil instance and return the address and a cleanup function
func StartAnvil(t *testing.T, showLogs bool) (string, func()) {
	port := networkTestUtils.FindFreePort(t)

	// we need mixed mining to work around https://github.com/xmtp/xmtpd/issues/643
	cmd := exec.Command(
		"anvil",
		"--port",
		fmt.Sprintf("%d", port),
		"--mixed-mining",
		"--block-time",
		"1",
	)
	if showLogs {
		// (Optional) You can capture stdout/stderr for logs:
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Start()
	require.NoError(t, err)
	url := fmt.Sprintf("http://localhost:%d", port)
	waitForAnvil(t, url)

	return url, func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}
}
