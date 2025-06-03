package integration_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestXDBGRealMLSPayloads(t *testing.T) {
	envVars := constructVariables(t)
	container, err := runContainer(
		t,
		"ghcr.io/xmtp/xmtpd:dev",
		"xmtpd_test_dev",
		envVars,
	)
	require.NoError(t, err, "Failed to start latest version container")
	require.NotNil(t, container, "Failed to start latest version container")

	port, err := container.MappedPort(t.Context(), "5050/tcp")
	require.NoError(t, err)

	err = runXDBG(t, port, GeneratorTypeIdentity, 10)
	require.NoError(t, err, "Failed to execute XDBG")
	err = runXDBG(t, port, GeneratorTypeGroup, 10)
	require.NoError(t, err, "Failed to execute XDBG")
	err = runXDBG(t, port, GeneratorTypeMessage, 10)
	require.NoError(t, err, "Failed to execute XDBG")
}
