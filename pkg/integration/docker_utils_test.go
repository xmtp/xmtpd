package integration_test

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
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

const (
	testFlag   = "ENABLE_INTEGRATION_TESTS"
	XDBG_IMAGE = "ghcr.io/xmtp/xdbg:sha-78a5ac2"
)

func skipIfNotEnabled() {
	if _, isSet := os.LookupEnv(testFlag); !isSet {
		fmt.Printf("Skipping integration test. %s is not set\n", testFlag)
		os.Exit(0)
	}
}

func loadEnvFromShell() (map[string]string, error) {
	scriptPath := testutils.GetScriptPath("./scripts/load_env.sh")
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
	vars["XMTPD_PAYER_ENABLE"] = "true"

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
	imageName string,
	containerName string,
	envVars map[string]string,
) (testcontainers.Container, error) {
	ctxwc, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	envVars["XMTPD_CONTRACTS_CONFIG_FILE_PATH"] = "/cfg/anvil.json"

	req := testcontainers.ContainerRequest{
		Image: imageName,
		Name:  containerName,
		Env:   envVars,
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      testutils.GetScriptPath("../../dev/environments/anvil.json"),
				ContainerFilePath: "/cfg/anvil.json",
				FileMode:          0o644,
			},
		},
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

	return xmtpContainer, err
}

type GeneratorType string

const (
	GeneratorTypeIdentity GeneratorType = "identity"
	GeneratorTypeGroup    GeneratorType = "group"
	GeneratorTypeMessage  GeneratorType = "message"
)

func (e GeneratorType) String() string {
	return string(e)
}

func runXDBG(
	t *testing.T,
	xmtpdPort nat.Port,
	genType GeneratorType,
	count uint64,
) error {
	ctxwc, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	dbVolumePath := "/tmp/testcontainer-xdbg-db"
	_ = os.MkdirAll(dbVolumePath, 0o755)

	targetAddr := fmt.Sprintf("http://host.docker.internal:%d", xmtpdPort.Int())

	req := testcontainers.ContainerRequest{
		Image: XDBG_IMAGE,
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.ExtraHosts = append(hc.ExtraHosts, "host.docker.internal:host-gateway")
			hc.Binds = append(hc.Binds, fmt.Sprintf("%s:/root/.local/share/xdbg/", dbVolumePath))
		},
		Cmd: []string{
			"-u", targetAddr, "-p", targetAddr, "-d", "generate", "-e", genType.String(), "-a", strconv.FormatUint(count, 10),
		},
		WaitingFor: wait.ForExit(),
	}

	xdbgContainer, err := testcontainers.GenericContainer(
		ctxwc,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Logger:           log.Default(),
		},
	)
	if err != nil {
		return err
	}

	testcontainers.CleanupContainer(t, xdbgContainer)

	return handleExitedContainer(ctxwc, xdbgContainer)
}

func handleExitedContainer(
	context context.Context,
	exitedContainer testcontainers.Container,
) error {
	state, err := exitedContainer.State(context)
	if err != nil {
		return err
	}

	if state.ExitCode != 0 {
		logs, logErr := exitedContainer.Logs(context)
		if logErr != nil {
			return fmt.Errorf(
				"container exited with code %d, but failed to get logs: %v",
				state.ExitCode,
				logErr,
			)
		}
		defer func() {
			_ = logs.Close()
		}()

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, logs)

		return fmt.Errorf("container exited with code %d\nLogs:\n%s", state.ExitCode, buf.String())
	}

	return nil
}
