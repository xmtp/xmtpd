package upgrade_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// Guaranteed to be in order of increasing version
var (
	xmtpdVersions = []string{
		"0.3.0",
	}

	ghcrRepository = "ghcr.io/xmtp/xmtpd"
)

func TestUpgradeToLatest(t *testing.T) {
	ctx := context.Background()
	for _, version := range xmtpdVersions {
		image := fmt.Sprintf("%s:%s", ghcrRepository, version)

		t.Run(version, func(t *testing.T) {
			envVars := constructVariables(t)
			t.Logf("Starting old container")
			err := runContainer(
				t,
				ctx,
				image,
				fmt.Sprintf("xmtpd_test_%s", version),
				envVars,
			)
			require.NoError(t, err, "Failed to start container version %s", version)

			t.Logf("Starting new container")
			err = runContainer(
				t,
				ctx,
				"ghcr.io/xmtp/xmtpd:dev",
				"xmtpd_test_dev",
				envVars,
			)
			require.NoError(t, err, "Failed to start dev container")
		})
	}
}

func TestLatestVersion(t *testing.T) {
	ctx := context.Background()
	envVars := constructVariables(t)
	err := runContainer(
		t,
		ctx,
		"ghcr.io/xmtp/xmtpd:dev",
		"xmtpd_test_dev",
		envVars,
	)
	require.NoError(t, err, "Failed to start latest version container")
}
