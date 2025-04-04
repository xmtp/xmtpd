package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var apiOpenConnections = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "xmtp_api_open_connections_gauge",
		Help: "Number of open API connections",
	},
	[]string{"style", "method"},
)

var apiIncomingNodeConnectionByVersionGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "xmtp_api_incoming_node_connection_by_version_gauge",
		Help: "Number of incoming node connections by version",
	},
	[]string{"version"},
)

type ApiOpenConnection struct {
	style  string
	method string
}

func NewApiOpenConnection(style string, method string) *ApiOpenConnection {
	oc := ApiOpenConnection{
		style:  style,
		method: method,
	}

	apiOpenConnections.With(prometheus.Labels{"style": oc.style, "method": oc.method}).Inc()

	return &oc
}

func (oc *ApiOpenConnection) Close() {
	apiOpenConnections.With(prometheus.Labels{"style": oc.style, "method": oc.method}).Dec()
}

type IncomingConnectionTracker struct {
	version string
}

func NewIncomingConnectionTracker(version string) *IncomingConnectionTracker {
	return &IncomingConnectionTracker{
		version: version,
	}
}

func (ct *IncomingConnectionTracker) Open() {
	apiIncomingNodeConnectionByVersionGauge.With(prometheus.Labels{"version": ct.version}).
		Inc()
}

func (ct *IncomingConnectionTracker) Close() {
	apiIncomingNodeConnectionByVersionGauge.With(prometheus.Labels{"version": ct.version}).
		Dec()
}
