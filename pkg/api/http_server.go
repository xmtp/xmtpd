package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

type HttpRegistrationFunc func(gwmux *runtime.ServeMux, conn *grpc.ClientConn) error

func (s *ApiServer) startHTTP(
	ctx context.Context,
	log *zap.Logger,
	registrationFunc HttpRegistrationFunc,
) error {
	var err error

	client, err := s.DialGRPC(ctx)
	if err != nil {
		log.Fatal("dialing GRPC from HTTP Gateway")
		return err
	}

	health := healthgrpc.NewHealthClient(client)

	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption("application/x-protobuf", &runtime.ProtoMarshaller{}),
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
		runtime.WithStreamErrorHandler(runtime.DefaultStreamErrorHandler),
		runtime.WithHealthzEndpoint(health),
	)
	if err := registrationFunc(gwmux, client); err != nil {
		return err
	}

	gwServer := &http.Server{
		Addr:    s.httpListener.Addr().String(),
		Handler: tracing.WrapHTTPHandler(allowCORS(gwmux), nil),
	}

	tracing.GoPanicWrap(s.ctx, &s.wg, "http", func(ctx context.Context) {
		if s.httpListener == nil {
			s.log.Fatal("no http listener")
		}
		s.log.Info("serving http", zap.String("address", s.httpListener.Addr().String()))
		err = gwServer.Serve(s.httpListener)
		if err != nil && err != http.ErrServerClosed && !isErrUseOfClosedConnection(err) {
			s.log.Error("serving http", zap.Error(err))
		}
	})

	return nil
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{
		"Content-Type",
		"Accept",
		"Authorization",
		"X-Client-Version",
		"X-App-Version",
		"Baggage",
		"DNT",
		"Sec-CH-UA",
		"Sec-CH-UA-Mobile",
		"Sec-CH-UA-Platform",
		"Sentry-Trace",
		"User-Agent",
		"x-libxmtp-version",
		"x-app-version",
	}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
			preflightHandler(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}
