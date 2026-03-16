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
		tests.NewChaosNetworkPartitionTest(),
		tests.NewChaosConnectionResetTest(),
		tests.NewChaosBandwidthThrottleTest(),
		tests.NewChaosCompoundFaultTest(),
		tests.NewChaosAttestationFaultTest(),
		tests.NewGatewayScaleTest(),
		tests.NewPayerLifecycleTest(),
		tests.NewMultiPayerTest(),
		tests.NewSettlementVerificationTest(),
		tests.NewSyncVerificationTest(),
		tests.NewSustainedLoadTest(),
		tests.NewRateRegistryChangeTest(),
		tests.NewStuckStateDetectionTest(),
	}
}
