package builders

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type GeneratorType string

const (
	GeneratorTypeIdentity GeneratorType = "identity"
	GeneratorTypeGroup    GeneratorType = "group"
	GeneratorTypeMessage  GeneratorType = "message"
)

func (e GeneratorType) String() string {
	return string(e)
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
		image:        xdbgImage,
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
