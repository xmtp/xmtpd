package metrics

import (
	"context"
	"fmt"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"net"
	"net/http"

	"github.com/pires/go-proxyproto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Server struct {
	ctx  context.Context
	log  *zap.Logger
	http net.Listener
}

func NewMetricsServer(
	ctx context.Context,
	address string,
	port int,
	log *zap.Logger,
	reg *prometheus.Registry,
) (*Server, error) {
	s := &Server{
		ctx: ctx,
		log: log.Named("metrics"),
	}

	addr := fmt.Sprintf("%s:%d", address, port)
	httpListener, err := net.Listen("tcp", addr)
	s.http = &proxyproto.Listener{Listener: httpListener}
	if err != nil {
		return nil, err
	}
	registerCollectors(reg)
	srv := http.Server{
		Addr:    addr,
		Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}),
	}

	go tracing.PanicWrap(ctx, "metrics server", func(_ context.Context) {
		s.log.Info("serving metrics http", zap.String("address", s.http.Addr().String()))
		err = srv.Serve(s.http)
		if err != nil {
			s.log.Error("serving http", zap.Error(err))
		}
	})

	return s, nil
}

func (s *Server) Close() error {
	return s.http.Close()
}

func registerCollectors(reg prometheus.Registerer) {
	//TODO: add metrics here
	var cols []prometheus.Collector
	for _, col := range cols {
		reg.MustRegister(col)
	}
}
