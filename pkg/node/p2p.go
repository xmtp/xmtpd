package node

import (
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	rcmgrObs "github.com/libp2p/go-libp2p/p2p/host/resource-manager/obs"
	"github.com/prometheus/client_golang/prometheus"
)

func p2pResourceManager(metrics *Metrics) (network.ResourceManager, error) {
	// From https://github.com/libp2p/go-libp2p/tree/410248e111b1169e5cbbab455267d35f2f38baba/p2p/host/resource-manager#usage
	// Start with the default scaling limits.
	scalingLimits := rcmgr.DefaultLimits

	// Add limits around included libp2p protocols
	libp2p.SetDefaultServiceLimits(&scalingLimits)

	// Turn the scaling limits into a concrete set of limits using `.AutoScale`. This
	// scales the limits proportional to your system memory.
	scaledDefaultLimits := scalingLimits.AutoScale()

	// Tweak certain settings
	cfg := rcmgr.PartialLimitConfig{
		Protocol: map[protocol.ID]rcmgr.ResourceLimits{
			syncProtocol: {
				// Allow unlimited outbound streams
				StreamsOutbound: rcmgr.Unlimited,
			},
			// Everything else is default. The exact values will come from `scaledDefaultLimits` above.
		}}

	// Create our limits by using our cfg and replacing the default values with values from `scaledDefaultLimits`
	limits := cfg.Build(scaledDefaultLimits)

	// The resource manager expects a limiter, se we create one from our limits.
	limiter := rcmgr.NewFixedLimiter(limits)

	// (Optional if you want metrics)
	var opts []rcmgr.Option
	if metrics != nil {
		rcmgrObs.MustRegisterWith(prometheus.DefaultRegisterer)
		// The stats are emitted by the trace reporter so
		// we have to add it to the resource manager options.
		str, err := rcmgrObs.NewStatsTraceReporter()
		if err != nil {
			return nil, err
		}
		opts = append(opts, rcmgr.WithTraceReporter(str))
	}

	// Initialize the resource manager
	rm, err := rcmgr.NewResourceManager(limiter, opts...)
	if err != nil {
		return nil, err
	}

	return rm, nil
}
