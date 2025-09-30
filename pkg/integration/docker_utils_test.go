package integration_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/network"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

const (
	testFlag = "ENABLE_INTEGRATION_TESTS"
)

func skipIfNotEnabled() {
	if _, isSet := os.LookupEnv(testFlag); !isSet {
		fmt.Printf("Skipping integration test. %s is not set\n", testFlag)
		os.Exit(0)
	}
}

func buildDevImage() error {
	scriptPath := testutils.GetScriptPath("../../dev/docker/build")

	// Set a 5-minute timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, scriptPath)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	// Run the command and check for errors
	if err := cmd.Run(); err != nil {
		// Handle timeout separately
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("build process timed out after 5 minutes")
		} else {
			return fmt.Errorf("build process failed: %s\n", errBuf.String())
		}
	}

	return nil
}

func buildGatewayDevImage() error {
	scriptPath := testutils.GetScriptPath("../../dev/docker/build-gateway")

	// Set a 5-minute timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, scriptPath)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	// Run the command and check for errors
	if err := cmd.Run(); err != nil {
		// Handle timeout separately
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("build process timed out after 5 minutes")
		} else {
			return fmt.Errorf("build process failed: %s\n", errBuf.String())
		}
	}

	return nil
}

func pullImage(imageName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	dockerClient, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return err
	}

	reader, err := dockerClient.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = reader.Close()
	}()

	_, err = io.Copy(io.Discard, reader)
	return err
}

func MakeDockerNetwork(t *testing.T) string {
	net, err := network.New(t.Context())
	require.NoError(t, err)
	return net.Name
}
