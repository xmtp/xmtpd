package integration_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// Guaranteed to be in order of increasing version
var (
	xmtpdVersions = []string{}

	ghcrRepository = "ghcr.io/xmtp/xmtpd"
)

func TestUpgradeToLatest(t *testing.T) {
	for _, version := range xmtpdVersions {
		image := fmt.Sprintf("%s:%s", ghcrRepository, version)

		t.Run(version, func(t *testing.T) {
			t.Logf("Starting old container")
			_, err := NewXmtpdContainerBuilder(t).
				WithImage(image).
				Build(t)
			require.NoError(t, err)

			t.Logf("Starting new container")
			_, err = NewXmtpdContainerBuilder(t).
				WithImage("ghcr.io/xmtp/xmtpd:dev").
				Build(t)
			require.NoError(t, err, "Failed to start dev container")
		})
	}
}

func TestLatestVersion(t *testing.T) {
	_, err := NewXmtpdContainerBuilder(t).
		WithImage("ghcr.io/xmtp/xmtpd:dev").
		Build(t)

	require.NoError(t, err)

	_, err = NewXmtpdGatewayContainerBuilder(t).
		WithImage("ghcr.io/xmtp/xmtpd-gateway:dev").
		Build(t)

	require.NoError(t, err)
}
