package integration_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestXDBGRealMLSPayloads(t *testing.T) {
	network := MakeDockerNetwork(t)

	xmtpdContainer, err := NewXmtpdContainerBuilder(t).
		WithImage("ghcr.io/xmtp/xmtpd:dev").
		WithNetwork(network).
		Build(t)

	require.NoError(t, err)

	name, err := xmtpdContainer.Name(t.Context())
	require.NoError(t, err)

	target := fmt.Sprintf("http:/%s:5050", name)

	err = NewXdbgContainerBuilder().
		WithNetwork(network).
		WithTarget(target).
		WithGeneratorType(GeneratorTypeIdentity).
		WithCount(10).
		Build(t)
	require.NoError(t, err)

	err = NewXdbgContainerBuilder().
		WithNetwork(network).
		WithTarget(target).
		WithGeneratorType(GeneratorTypeGroup).
		WithCount(10).
		Build(t)
	require.NoError(t, err)

	err = NewXdbgContainerBuilder().
		WithNetwork(network).
		WithTarget(target).
		WithGeneratorType(GeneratorTypeMessage).
		WithCount(10).
		Build(t)
	require.NoError(t, err)
}
