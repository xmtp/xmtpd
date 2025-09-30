package builders

import (
	"context"
	"log"
	"maps"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

type CLIContainerBuilder struct {
	image        string
	networkNames []string
	cmd          []string
	envVars      map[string]string
	files        []testcontainers.ContainerFile
}

func NewCLIContainerBuilder(t *testing.T) *CLIContainerBuilder {
	return &CLIContainerBuilder{
		image: cliImage,
		files: []testcontainers.ContainerFile{{
			HostFilePath:      testutils.GetScriptPath(anvilJsonRelativePath),
			ContainerFilePath: "/cfg/anvil.json",
			FileMode:          0o644,
		}},
	}
}

func (b *CLIContainerBuilder) WithNetwork(network string) *CLIContainerBuilder {
	b.networkNames = []string{network}
	return b
}

func (b *CLIContainerBuilder) WithCmd(cmd []string) *CLIContainerBuilder {
	b.cmd = cmd
	return b
}

func (b *CLIContainerBuilder) WithEnvVars(envVars map[string]string) *CLIContainerBuilder {
	maps.Copy(b.envVars, envVars)
	return b
}

func (b *CLIContainerBuilder) Build(t *testing.T) error {
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:    b.image,
		Networks: b.networkNames,
		Env:      b.envVars,
		Files:    b.files,
		Cmd:      b.cmd,
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.ExtraHosts = append(hc.ExtraHosts, "host.docker.internal:host-gateway")
		},
	}

	cliContainer, err := testcontainers.GenericContainer(
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

	testcontainers.CleanupContainer(t, cliContainer)

	return handleExitedContainer(ctx, cliContainer)
}
