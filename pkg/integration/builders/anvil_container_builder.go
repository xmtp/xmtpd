package builders

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
)

func NewAnvilContainerBuilder(t *testing.T) *AnvilContainerBuilder {
	return &AnvilContainerBuilder{
		imageName: anvil.AnvilImage,
	}
}

type AnvilContainerBuilder struct {
	imageName     string
	containerName string
	exposedPorts  []string
	networkName   string
	networkAlias  string
}

func (b *AnvilContainerBuilder) WithImage(imageName string) *AnvilContainerBuilder {
	b.imageName = imageName
	return b
}

func (b *AnvilContainerBuilder) WithContainerName(
	name string,
) *AnvilContainerBuilder {
	b.containerName = name
	return b
}

func (b *AnvilContainerBuilder) WithPort(port string) *AnvilContainerBuilder {
	b.exposedPorts = append(b.exposedPorts, port)
	return b
}

func (b *AnvilContainerBuilder) WithNetwork(
	networkName string,
) *AnvilContainerBuilder {
	b.networkName = networkName
	return b
}

func (b *AnvilContainerBuilder) WithNetworkAlias(alias string) *AnvilContainerBuilder {
	b.networkAlias = alias
	return b
}

func (b *AnvilContainerBuilder) Build(t *testing.T) (testcontainers.Container, error) {
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	aliases := make(map[string][]string)
	if b.networkAlias != "" {
		require.NotEmpty(t, b.networkName)
		aliases[b.networkName] = []string{b.networkAlias}
	}

	req := testcontainers.ContainerRequest{
		Image:        b.imageName,
		Name:         b.containerName,
		ExposedPorts: b.exposedPorts,
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.AutoRemove = true
		},
		Networks:       []string{b.networkName},
		NetworkAliases: aliases,
		WaitingFor:     wait.ForLog("Listening on"),
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

	testcontainers.CleanupContainer(t, anvilContainer)

	return anvilContainer, err
}
