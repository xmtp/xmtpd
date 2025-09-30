package integration_test

import (
	"fmt"
	"testing"

	"github.com/xmtp/xmtpd/pkg/integration/builders"

	"github.com/stretchr/testify/require"
)

func TestXDBGRealMLSPayloads(t *testing.T) {
	network := MakeDockerNetwork(t)
	xmtpdAlias := "xmtpd"
	gatewayAlias := "gateway"
	anvilAlias := "anvil"

	_, err := builders.NewAnvilContainerBuilder(t).
		WithNetwork(network).
		WithNetworkAlias(anvilAlias).
		Build(t)
	require.NoError(t, err)

	wsHost := fmt.Sprintf("ws://%s:8545", anvilAlias)
	rpcHost := fmt.Sprintf("http://%s:8545", anvilAlias)

	builders.RegisterNode(t, network, rpcHost, xmtpdAlias)
	builders.EnableNode(t, network, rpcHost, 100)

	_, err = builders.NewXmtpdContainerBuilder(t).
		WithImage("ghcr.io/xmtp/xmtpd:dev").
		WithNetwork(network).
		WithNetworkAlias(xmtpdAlias).
		WithWsUrl(wsHost).
		WithRPCUrl(rpcHost).
		Build(t)
	require.NoError(t, err)

	_, err = builders.NewXmtpdGatewayContainerBuilder(t).
		WithImage("ghcr.io/xmtp/xmtpd-gateway:dev").
		WithNetwork(network).
		WithNetworkAlias(gatewayAlias).
		WithWsUrl(wsHost).
		WithRPCUrl(rpcHost).
		Build(t)
	require.NoError(t, err)

	target := fmt.Sprintf("http://%s:5050", xmtpdAlias)
	gatewayTarget := fmt.Sprintf("http:/%s:5050", gatewayAlias)

	err = builders.NewXdbgContainerBuilder().
		WithNetwork(network).
		WithTarget(target).
		WithGatewayTarget(gatewayTarget).
		WithGeneratorType(builders.GeneratorTypeIdentity).
		WithCount(10).
		Build(t)
	require.NoError(t, err)

	err = builders.NewXdbgContainerBuilder().
		WithNetwork(network).
		WithTarget(target).
		WithGatewayTarget(gatewayTarget).
		WithGeneratorType(builders.GeneratorTypeGroup).
		WithCount(10).
		Build(t)
	require.NoError(t, err)

	err = builders.NewXdbgContainerBuilder().
		WithNetwork(network).
		WithTarget(target).
		WithGatewayTarget(gatewayTarget).
		WithGeneratorType(builders.GeneratorTypeMessage).
		WithCount(10).
		Build(t)
	require.NoError(t, err)
}
