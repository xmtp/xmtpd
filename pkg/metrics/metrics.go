package metrics

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/xmtp/xmtpd/pkg/tracing"

	"github.com/pires/go-proxyproto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Server struct {
	ctx  context.Context
	log  *zap.Logger
	http net.Listener
	wg   sync.WaitGroup
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
		wg:  sync.WaitGroup{},
	}

	addr := fmt.Sprintf("%s:%d", address, port)
	httpListener, err := net.Listen("tcp", addr)
	s.http = &proxyproto.Listener{Listener: httpListener}
	if err != nil {
		return nil, err
	}
	registerCollectors(reg)
	srv := http.Server{
		Addr: addr,
		Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{
			EnableOpenMetrics: true,
			Registry:          reg,
		}),
	}

	tracing.GoPanicWrap(s.ctx, &s.wg, "metrics-server", func(ctx context.Context) {
		s.log.Info("serving metrics http", zap.String("address", s.http.Addr().String()))
		err = srv.Serve(s.http)
		if err != nil {
			s.log.Error("serving http", zap.Error(err))
		}
	})

	return s, nil
}

func (s *Server) Close() {
	s.log.Debug("Closing")
	_ = s.http.Close()
	s.wg.Wait()
	s.log.Debug("Closed")
}

func registerCollectors(reg prometheus.Registerer) {
	cols := []prometheus.Collector{
		indexerNumLogsFound,
		indexerCurrentBlock,
		indexerMaxBlock,
		indexerCurrentBlockLag,
		indexerCountRetryableStorageErrors,
		indexerGetLogsDuration,
		indexerGetLogsRequests,
		indexerLogProcessingTime,
		payerNodePublishDuration,
		payerCursorBlockTime,
		payerCurrentNonce,
		payerBanlistRetry,
		payerMessagesOriginated,
		syncOriginatorSequenceId,
		syncOutgoingSyncConnections,
		syncFailedOutgoingSyncConnections,
		syncFailedOutgoingSyncConnectionCounter,
		apiOpenConnections,
		apiIncomingNodeConnectionByVersionGauge,
		apiNodeConnectionRequestsByVersionCounter,
		apiFailedGRPCRequestsCounter,
		blockchainWaitForTransaction,
		blockchainPublishPayload,
		payerGetReaderNodeAvailableNodes,
	}

	for _, col := range cols {
		reg.MustRegister(col)
	}
}
