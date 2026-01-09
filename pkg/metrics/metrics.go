// Package metrics implements the Prometheusmetrics for the XMTPD service.
package metrics

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"

	"github.com/pires/go-proxyproto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Server struct {
	ctx    context.Context
	logger *zap.Logger
	http   net.Listener
	wg     sync.WaitGroup
}

func NewMetricsServer(
	ctx context.Context,
	address string,
	port int,
	logger *zap.Logger,
	reg *prometheus.Registry,
) (*Server, error) {
	s := &Server{
		ctx:    ctx,
		logger: logger.Named(utils.MetricsLoggerName),
		wg:     sync.WaitGroup{},
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
		s.logger.Info("serving metrics http", zap.String("address", s.http.Addr().String()))
		if err := srv.Serve(s.http); err != nil {
			switch {
			case errors.Is(err, http.ErrServerClosed), errors.Is(err, net.ErrClosed):
				s.logger.Info("metrics server closing", zap.Error(err))
			default:
				s.logger.Error("error serving http", zap.Error(err))
			}
		}
	})

	return s, nil
}

func (s *Server) Close() {
	s.logger.Debug("closing")
	_ = s.http.Close()
	s.wg.Wait()
	s.logger.Debug("closed")
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
		indexerLogBytesIndexed,
		gatwayPublishDuration,
		gatewayCurrentNonce,
		gatewayBanlistRetry,
		gatewayMessagesOriginated,
		syncOriginatorSequenceID,
		syncOriginatorMessagesReceived,
		syncOriginatorErrorMessages,
		syncOutgoingSyncConnections,
		syncFailedOutgoingSyncConnections,
		syncFailedOutgoingSyncConnectionCounter,
		apiOpenConnections,
		apiIncomingNodeConnectionByVersionGauge,
		apiNodeConnectionRequestsByVersionCounter,
		apiFailedGRPCRequestsCounter,
		apiStagedEnvelopeProcessingDelay,
		apiWaitForGatewayPublish,
		grpcServerStartedTotal,
		grpcServerHandledTotal,
		grpcServerMsgReceivedTotal,
		grpcServerMsgSentTotal,
		grpcServerHandlingSeconds,
		blockchainBroadcastTransaction,
		blockchainWaitForTransaction,
		blockchainPublishPayload,
		blockchainGasPriceGauge,
		blockchainGasPriceUpdatesTotal,
		blockchainGasPriceDefaultFallbackTotal,
		blockchainGasPriceLastUpdateTimestamp,
		gatewayGetNodesAvailableNodes,
		migratorE2ELatency,
		migratorDestLastSequenceID,
		migratorReaderErrors,
		migratorReaderFetchDuration,
		migratorReaderNumRowsFound,
		migratorSourceLastSequenceID,
		migratorTransformerErrors,
		migratorWriterErrors,
		migratorWriterLatency,
		migratorWriterRetryAttempts,
		migratorWriterRowsMigrated,
		migratorWriterBytesMigrated,
		migratorTargetLastSequenceID,
		QueryDuration,
		QueryErrors,
	}

	for _, col := range cols {
		reg.MustRegister(col)
	}
}
