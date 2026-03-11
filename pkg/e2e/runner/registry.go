package runner

import (
	"github.com/xmtp/xmtpd/pkg/e2e/tests"
	"github.com/xmtp/xmtpd/pkg/e2e/types"
)

func AllTests() []types.Test {
	return []types.Test{
		tests.NewSmokeTest(),
		tests.NewChaosNodeDownTest(),
		tests.NewChaosLatencyTest(),
		tests.NewGatewayScaleTest(),
		tests.NewPayerLifecycleTest(),
	}
}
