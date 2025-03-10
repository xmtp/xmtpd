package upgrade_test

import (
	"context"
	"fmt"
	"testing"
)

// Guaranteed to be in order of increasing version
var (
	xmtpdVersions = []string{
		"0.1.4",
		"0.2.0",
		"0.2.1",
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
			runContainer(
				t,
				ctx,
				image,
				fmt.Sprintf("xmtpd_test_%s", version),
				envVars,
			)

			t.Logf("Starting new container")
			runContainer(
				t,
				ctx,
				"ghcr.io/xmtp/xmtpd:dev",
				"xmtpd_test_dev",
				envVars,
			)
		})
	}
}

func TestLatestVersion(t *testing.T) {
	ctx := context.Background()
	envVars := constructVariables(t)
	runContainer(
		t,
		ctx,
		"ghcr.io/xmtp/xmtpd:dev",
		"xmtpd_test_dev",
		envVars,
	)
}
