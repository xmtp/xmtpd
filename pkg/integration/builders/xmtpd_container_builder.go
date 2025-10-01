package builders

import (
	"context"
	"log"
	"maps"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

type XmtpdContainerBuilder struct {
	imageName     string
	containerName string
	envVars       map[string]string
	files         []testcontainers.ContainerFile
	exposedPorts  []string
	networkName   string
	wsUrl         string
	rpcUrl        string
	networkAlias  string
}

func NewXmtpdContainerBuilder(t *testing.T) *XmtpdContainerBuilder {
	envVars := constructVariables(t)
	envVars["XMTPD_CONTRACTS_CONFIG_FILE_PATH"] = "/cfg/anvil.json"

	return &XmtpdContainerBuilder{
		envVars: envVars,
		files: []testcontainers.ContainerFile{{
			HostFilePath:      testutils.GetScriptPath(anvilJsonRelativePath),
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

func (b *XmtpdContainerBuilder) WithNetworkAlias(alias string) *XmtpdContainerBuilder {
	b.networkAlias = alias
	return b
}

func (b *XmtpdContainerBuilder) WithWsUrl(url string) *XmtpdContainerBuilder {
	b.wsUrl = url
	return b
}

func (b *XmtpdContainerBuilder) WithRPCUrl(url string) *XmtpdContainerBuilder {
	b.rpcUrl = url
	return b
}

func (b *XmtpdContainerBuilder) Build(t *testing.T) (testcontainers.Container, error) {
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	if b.wsUrl != "" {
		b.envVars["XMTPD_SETTLEMENT_CHAIN_WSS_URL"] = b.wsUrl
		b.envVars["XMTPD_APP_CHAIN_WSS_URL"] = b.wsUrl
	}

	if b.rpcUrl != "" {
		b.envVars["XMTPD_SETTLEMENT_CHAIN_RPC_URL"] = b.rpcUrl
		b.envVars["XMTPD_APP_CHAIN_RPC_URL"] = b.rpcUrl
	}

	aliases := make(map[string][]string)
	if b.networkAlias != "" {
		require.NotEmpty(t, b.networkName)
		aliases[b.networkName] = []string{b.networkAlias}
	}

	req := testcontainers.ContainerRequest{
		Image:        b.imageName,
		Name:         b.containerName,
		Env:          b.envVars,
		Files:        b.files,
		ExposedPorts: b.exposedPorts,
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.ExtraHosts = append(hc.ExtraHosts, "host.docker.internal:host-gateway")
		},
		Networks:       []string{b.networkName},
		NetworkAliases: aliases,
		WaitingFor:     wait.ForLog("serving grpc"),
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
