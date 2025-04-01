package upgrade_test

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

const testFlag = "ENABLE_UPGRADE_TESTS"

func skipIfNotEnabled() {
	if _, isSet := os.LookupEnv(testFlag); !isSet {
		fmt.Printf("Skipping upgrade test. %s is not set\n", testFlag)
		os.Exit(0)
	}
}

func getScriptPath(scriptName string) string {
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)
	return filepath.Join(baseDir, scriptName)
}

func loadEnvFromShell() (map[string]string, error) {
	scriptPath := getScriptPath("./scripts/load_env.sh")
	cmd := exec.Command(scriptPath)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf(
			"error loading env via shell script: %v\nError: %s",
			err,
			errBuf.String(),
		)
	}

	envMap := make(map[string]string)
	scanner := bufio.NewScanner(&outBuf)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}
	return envMap, nil
}

func expandVars(vars map[string]string) {
	vars["XMTPD_REPLICATION_ENABLE"] = "true"
	vars["XMTPD_INDEXER_ENABLE"] = "true"

	dbName := testutils.GetCallerName(3) + "_" + testutils.RandomStringLower(6)

	vars["XMTPD_DB_NAME_OVERRIDE"] = dbName
}

func convertLocalhost(vars map[string]string) {
	for varKey, varValue := range vars {
		if strings.Contains(varValue, "localhost") {
			vars[varKey] = strings.ReplaceAll(varValue, "localhost", "host.docker.internal")
		}
	}
}

func constructVariables(t *testing.T) map[string]string {
	envVars, err := loadEnvFromShell()
	require.NoError(t, err)
	expandVars(envVars)
	convertLocalhost(envVars)

	return envVars
}

func buildDevImage() error {
	scriptPath := getScriptPath("../../dev/docker/build")

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

func runContainer(
	t *testing.T,
	ctx context.Context,
	imageName string,
	containerName string,
	envVars map[string]string,
) (err error) {
	ctxwc, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:        imageName,
		Name:         containerName,
		Env:          envVars,
		ExposedPorts: []string{"5050/tcp"},

		HostConfigModifier: func(hc *container.HostConfig) {
			hc.ExtraHosts = append(hc.ExtraHosts, "host.docker.internal:host-gateway")
		},
		WaitingFor: wait.ForLog(
			"serving grpc",
		), // TODO: Ideally we wait for health/liveness probe
	}

	xmtpContainer, err := testcontainers.GenericContainer(
		ctxwc,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Logger:           log.Default(),
		},
	)

	testcontainers.CleanupContainer(t, xmtpContainer)

	return err
}
