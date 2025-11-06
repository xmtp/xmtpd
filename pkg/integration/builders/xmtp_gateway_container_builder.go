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

type XmtpdGatewayContainerBuilder struct {
	imageName     string
	containerName string
	envVars       map[string]string
	files         []testcontainers.ContainerFile
	exposedPorts  []string
	networkName   string
	networkAlias  string
	wsURL         string
	rpcURL        string
}

func NewXmtpdGatewayContainerBuilder(t *testing.T) *XmtpdGatewayContainerBuilder {
	envVars := constructVariables(t)
	envVars["XMTPD_CONTRACTS_CONFIG_FILE_PATH"] = "/cfg/anvil.json"

	return &XmtpdGatewayContainerBuilder{
		envVars: envVars,
		files: []testcontainers.ContainerFile{{
			HostFilePath:      testutils.GetScriptPath(anvilJSONRelativePath),
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

func (b *XmtpdGatewayContainerBuilder) WithNetworkAlias(
	alias string,
) *XmtpdGatewayContainerBuilder {
	b.networkAlias = alias
	return b
}

func (b *XmtpdGatewayContainerBuilder) WithWsURL(url string) *XmtpdGatewayContainerBuilder {
	b.wsURL = url
	return b
}

func (b *XmtpdGatewayContainerBuilder) WithRPCURL(url string) *XmtpdGatewayContainerBuilder {
	b.rpcURL = url
	return b
}

func (b *XmtpdGatewayContainerBuilder) Build(t *testing.T) (testcontainers.Container, error) {
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	if b.wsURL != "" {
		b.envVars["XMTPD_SETTLEMENT_CHAIN_WSS_URL"] = b.wsURL
		b.envVars["XMTPD_APP_CHAIN_WSS_URL"] = b.wsURL
	}

	if b.rpcURL != "" {
		b.envVars["XMTPD_SETTLEMENT_CHAIN_RPC_URL"] = b.rpcURL
		b.envVars["XMTPD_APP_CHAIN_RPC_URL"] = b.rpcURL
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
		WaitingFor:     wait.ForLog("started api server"),
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
