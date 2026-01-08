package metrics

import (
	"connectrpc.com/connect"
	"github.com/prometheus/client_golang/prometheus"
	"time"
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

var apiNodeConnectionRequestsByVersionCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "xmtp_api_node_connection_requests_by_version_counter",
		Help: "Number of incoming node connections by version",
	},
	[]string{"version"},
)

var apiFailedGRPCRequestsCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "xmtp_api_failed_grpc_requests_counter",
		Help: "Number of failed GRPC requests by code",
	},
	[]string{"code"},
)

type APIOpenConnection struct {
	style  string
	method string
}

func NewAPIOpenConnection(style string, method string) *APIOpenConnection {
	oc := APIOpenConnection{
		style:  style,
		method: method,
	}

	apiOpenConnections.With(prometheus.Labels{"style": oc.style, "method": oc.method}).Inc()

	return &oc
}

func (oc *APIOpenConnection) Close() {
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

func EmitNewConnectionRequestVersion(version string) {
	apiNodeConnectionRequestsByVersionCounter.With(prometheus.Labels{"version": version}).
		Inc()
}

func EmitNewFailedGRPCRequest(code connect.Code) {
	apiFailedGRPCRequestsCounter.With(prometheus.Labels{"code": code.String()}).
		Inc()
}

var apiWaitForGatewayPublish = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name: "xmtp_api_wait_for_gateway_publish_seconds",
		Help: "Time to publish a payload to the blockchain",
	},
)

func EmitApiWaitForGatewayPublish(
	duration time.Duration,
) {
	apiWaitForGatewayPublish.Observe(duration.Seconds())
}