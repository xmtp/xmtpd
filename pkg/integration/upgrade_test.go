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
			envVars := constructVariables(t)
			t.Logf("Starting old container")
			_, err := runContainer(
				t,
				image,
				fmt.Sprintf("xmtpd_test_%s", version),
				envVars,
			)
			require.NoError(t, err, "Failed to start container version %s", version)

			t.Logf("Starting new container")
			_, err = runContainer(
				t,
				"ghcr.io/xmtp/xmtpd:dev",
				"xmtpd_test_dev",
				envVars,
			)
			require.NoError(t, err, "Failed to start dev container")
		})
	}
}

func TestLatestVersion(t *testing.T) {
	envVars := constructVariables(t)
	_, err := runContainer(
		t,
		"ghcr.io/xmtp/xmtpd:dev",
		"xmtpd_test_dev",
		envVars,
	)
	require.NoError(t, err, "Failed to start latest version container")
}
