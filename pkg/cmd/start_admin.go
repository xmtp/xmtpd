package cmd

import (
	"net"
	"net/http"
	_ "net/http/pprof"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/node"
)

type AdminOptions struct {
	Address string `long:"address" description:"Admin HTTP listen address" default:"0.0.0.0"`
	Port    uint   `long:"port" description:"Admin HTTP listen port"`
}

func startAdmin(ctx context.Context, opts *AdminOptions) *node.Metrics {
	if opts.Port == 0 {
		return nil
	}

	metrics := node.NewMetrics()
	http.Handle("/metrics", promhttp.Handler())
	addr := net.JoinHostPort(opts.Address, strconv.Itoa(int(opts.Port)))
	go func() { _ = http.ListenAndServe(addr, nil) }()
	return metrics
}
