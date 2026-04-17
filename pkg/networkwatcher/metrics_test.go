package networkwatcher

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
)

func TestRegisterMetrics_RegistersOnFreshRegistry(t *testing.T) {
	reg := prometheus.NewRegistry()
	require.NotPanics(t, func() { RegisterMetrics(reg) })
}

func TestRegisterMetrics_PanicsOnDoubleRegistration(t *testing.T) {
	reg := prometheus.NewRegistry()
	RegisterMetrics(reg)
	require.Panics(t, func() { RegisterMetrics(reg) })
}
