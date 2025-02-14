package upgrade_test

import (
	"fmt"
	"testing"
)

var upgradeToLatest = map[string]string{
	"0.1.4": "ghcr.io/xmtp/xmtpd:0.1.4",
}

func TestUpgradeToLatest(t *testing.T) {
	for version, image := range upgradeToLatest {
		t.Run(version, func(t *testing.T) {

			envVars := constructVariables(t)
			t.Logf("Starting old container")
			runContainer(t, fmt.Sprintf("xmtpd_test_%s", version), image, envVars)

			t.Logf("Starting new container")
			runContainer(t, "xmtpd_test_dev", "ghcr.io/xmtp/xmtpd:dev", envVars)
		})
	}

}
