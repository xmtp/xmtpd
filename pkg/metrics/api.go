package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var apiOpenConnections = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "xmtp_api_open_connections_gauge",
		Help: "Duration of the node publish call",
	},
	[]string{"style", "method"},
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
