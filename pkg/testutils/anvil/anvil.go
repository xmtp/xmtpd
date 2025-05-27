package anvil

import (
	"bufio"
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/log"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
)

func streamContainerLogs(t *testing.T, ctx context.Context, container testcontainers.Container) {
	rc, err := container.Logs(ctx)
	if err != nil {
		t.Logf("error streaming logs: %v", err)
		return
	}
	defer func() {
		_ = rc.Close()
	}()

	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		t.Log(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		t.Logf("log scanner error: %v", err)
	}
}

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

// StartAnvil starts an ephemeral anvil instance and return the address
func StartAnvil(t *testing.T, showLogs bool) string {
	ctx := t.Context()

	req := testcontainers.ContainerRequest{
		Image:        "ghcr.io/xmtp/contracts:v0.4.0",
		ExposedPorts: []string{"8545/tcp"},
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.AutoRemove = true
		},
		WaitingFor: wait.ForLog("Listening on"),
	}

	anvilContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Logger:           log.Default(),
		},
	)
	require.NoError(t, err)

	if showLogs {
		go streamContainerLogs(t, ctx, anvilContainer)
	}

	t.Cleanup(func() {
		_ = anvilContainer.Terminate(ctx)
	})

	mappedPort, err := anvilContainer.MappedPort(ctx, "8545/tcp")
	require.NoError(t, err)

	anvilURL := fmt.Sprintf("http://localhost:%s", mappedPort.Port())
	waitForAnvil(t, anvilURL)

	return anvilURL
}
