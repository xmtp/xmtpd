package integration_test

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"maps"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/network"

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
	XDBG_IMAGE = "ghcr.io/xmtp/xdbg:sha-26bb960"
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

type GeneratorType string

const (
	GeneratorTypeIdentity GeneratorType = "identity"
	GeneratorTypeGroup    GeneratorType = "group"
	GeneratorTypeMessage  GeneratorType = "message"
)

func (e GeneratorType) String() string {
	return string(e)
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

type XmtpdContainerBuilder struct {
	imageName     string
	containerName string
	envVars       map[string]string
	files         []testcontainers.ContainerFile
	exposedPorts  []string
	networkName   string
}

func NewXmtpdContainerBuilder(t *testing.T) *XmtpdContainerBuilder {
	envVars := constructVariables(t)
	envVars["XMTPD_CONTRACTS_CONFIG_FILE_PATH"] = "/cfg/anvil.json"

	return &XmtpdContainerBuilder{
		envVars: envVars,
		files: []testcontainers.ContainerFile{{
			HostFilePath:      testutils.GetScriptPath("../../dev/environments/anvil.json"),
			ContainerFilePath: "/cfg/anvil.json",
			FileMode:          0o644,
		}},
	}
}

func (b *XmtpdContainerBuilder) WithImage(imageName string) *XmtpdContainerBuilder {
	b.imageName = imageName
	return b
}

func (b *XmtpdContainerBuilder) WithContainerName(name string) *XmtpdContainerBuilder {
	b.containerName = name
	return b
}

func (b *XmtpdContainerBuilder) WithEnvVars(envVars map[string]string) *XmtpdContainerBuilder {
	maps.Copy(b.envVars, envVars)
	return b
}

func (b *XmtpdContainerBuilder) WithFile(file testcontainers.ContainerFile) *XmtpdContainerBuilder {
	b.files = append(b.files, file)
	return b
}

func (b *XmtpdContainerBuilder) WithPort(port string) *XmtpdContainerBuilder {
	b.exposedPorts = append(b.exposedPorts, port)
	return b
}

func (b *XmtpdContainerBuilder) WithNetwork(networkName string) *XmtpdContainerBuilder {
	b.networkName = networkName
	return b
}

func (b *XmtpdContainerBuilder) Build(t *testing.T) (testcontainers.Container, error) {
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:        b.imageName,
		Name:         b.containerName,
		Env:          b.envVars,
		Files:        b.files,
		ExposedPorts: b.exposedPorts,
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.ExtraHosts = append(hc.ExtraHosts, "host.docker.internal:host-gateway")
		},
		Networks:   []string{b.networkName},
		WaitingFor: wait.ForLog("serving grpc"),
	}

	xmtpdContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Logger:           log.Default(),
		},
	)

	testcontainers.CleanupContainer(t, xmtpdContainer)

	return xmtpdContainer, err
}

type XdbgContainerBuilder struct {
	image        string
	networkNames []string
	targetAddr   string
	gatewayAddr  string
	genType      GeneratorType
	count        string
	dbVolumePath string
	waitStrategy wait.Strategy
}

func NewXdbgContainerBuilder() *XdbgContainerBuilder {
	return &XdbgContainerBuilder{
		image:        XDBG_IMAGE,
		dbVolumePath: "/tmp/testcontainer-xdbg-db",
		waitStrategy: wait.ForExit(),
	}
}

func (b *XdbgContainerBuilder) WithNetwork(network string) *XdbgContainerBuilder {
	b.networkNames = []string{network}
	return b
}

func (b *XdbgContainerBuilder) WithTarget(addr string) *XdbgContainerBuilder {
	b.targetAddr = addr
	return b
}

func (b *XdbgContainerBuilder) WithGatewayTarget(addr string) *XdbgContainerBuilder {
	b.gatewayAddr = addr
	return b
}

func (b *XdbgContainerBuilder) WithGeneratorType(genType GeneratorType) *XdbgContainerBuilder {
	b.genType = genType
	return b
}

func (b *XdbgContainerBuilder) WithCount(count uint64) *XdbgContainerBuilder {
	b.count = strconv.FormatUint(count, 10)
	return b
}

func (b *XdbgContainerBuilder) Build(t *testing.T) error {
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	if err := os.MkdirAll(b.dbVolumePath, 0o755); err != nil {
		return fmt.Errorf("failed to create volume directory: %w", err)
	}

	req := testcontainers.ContainerRequest{
		Image:    b.image,
		Networks: b.networkNames,
		Cmd: []string{
			"-u", b.targetAddr,
			"-p", b.gatewayAddr,
			"-d", "generate",
			"-e", b.genType.String(),
			"-a", b.count,
		},
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.ExtraHosts = append(hc.ExtraHosts, "host.docker.internal:host-gateway")
			hc.Binds = append(hc.Binds, fmt.Sprintf("%s:/root/.local/share/xdbg/", b.dbVolumePath))
		},
		WaitingFor: b.waitStrategy,
	}

	xdbgContainer, err := testcontainers.GenericContainer(
		ctx,
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

	return handleExitedContainer(ctx, xdbgContainer)
}

func MakeDockerNetwork(t *testing.T) string {
	net, err := network.New(t.Context())
	require.NoError(t, err)
	return net.Name
}

type XmtpdGatewayContainerBuilder struct {
	imageName     string
	containerName string
	envVars       map[string]string
	files         []testcontainers.ContainerFile
	exposedPorts  []string
	networkName   string
}

func NewXmtpdGatewayContainerBuilder(t *testing.T) *XmtpdGatewayContainerBuilder {
	envVars := constructVariables(t)
	envVars["XMTPD_CONTRACTS_CONFIG_FILE_PATH"] = "/cfg/anvil.json"

	return &XmtpdGatewayContainerBuilder{
		envVars: envVars,
		files: []testcontainers.ContainerFile{{
			HostFilePath:      testutils.GetScriptPath("../../dev/environments/anvil.json"),
			ContainerFilePath: "/cfg/anvil.json",
			FileMode:          0o644,
		}},
	}
}

func (b *XmtpdGatewayContainerBuilder) WithImage(imageName string) *XmtpdGatewayContainerBuilder {
	b.imageName = imageName
	return b
}

func (b *XmtpdGatewayContainerBuilder) WithContainerName(
	name string,
) *XmtpdGatewayContainerBuilder {
	b.containerName = name
	return b
}

func (b *XmtpdGatewayContainerBuilder) WithEnvVars(
	envVars map[string]string,
) *XmtpdGatewayContainerBuilder {
	maps.Copy(b.envVars, envVars)
	return b
}

func (b *XmtpdGatewayContainerBuilder) WithFile(
	file testcontainers.ContainerFile,
) *XmtpdGatewayContainerBuilder {
	b.files = append(b.files, file)
	return b
}

func (b *XmtpdGatewayContainerBuilder) WithPort(port string) *XmtpdGatewayContainerBuilder {
	b.exposedPorts = append(b.exposedPorts, port)
	return b
}

func (b *XmtpdGatewayContainerBuilder) WithNetwork(
	networkName string,
) *XmtpdGatewayContainerBuilder {
	b.networkName = networkName
	return b
}

func (b *XmtpdGatewayContainerBuilder) Build(t *testing.T) (testcontainers.Container, error) {
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:        b.imageName,
		Name:         b.containerName,
		Env:          b.envVars,
		Files:        b.files,
		ExposedPorts: b.exposedPorts,
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.ExtraHosts = append(hc.ExtraHosts, "host.docker.internal:host-gateway")
		},
		Networks:   []string{b.networkName},
		WaitingFor: wait.ForLog("serving grpc"),
	}

	xmtpdContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Logger:           log.Default(),
		},
	)

	testcontainers.CleanupContainer(t, xmtpdContainer)

	return xmtpdContainer, err
}
