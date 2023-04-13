package gateway

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/zap"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc/filters"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	opts      *Options
	ctx       context.Context
	log       *zap.Logger
	messagev1 messagev1.MessageApiServer
	metrics   *Metrics

	grpc net.Listener
	http net.Listener
}

func New(ctx context.Context, messagev1 messagev1.MessageApiServer, metrics *Metrics, opts *Options) (*Server, error) {
	err := opts.validate()
	if err != nil {
		return nil, err
	}

	log := ctx.Logger().Named("api")

	s := &Server{
		ctx:       ctx,
		log:       log,
		opts:      opts,
		messagev1: messagev1,
		metrics:   metrics,
	}

	err = s.startGRPC()
	if err != nil {
		return nil, err
	}

	err = s.startHTTP()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) startGRPC() error {
	var err error

	s.grpc, err = net.Listen("tcp", hostPortAddr(s.opts.GRPCAddress, s.opts.GRPCPort))
	if err != nil {
		return errors.Wrap(err, "creating grpc listener")
	}

	telemetry, err := NewTelemetryInterceptor(s.log, s.metrics)
	if err != nil {
		return err
	}
	unary := []grpc.UnaryServerInterceptor{
		telemetry.Unary(),
		otelgrpc.UnaryServerInterceptor(
			otelgrpc.WithInterceptorFilter(
				filters.Not(
					filters.HealthCheck(),
				),
			),
		),
	}
	stream := []grpc.StreamServerInterceptor{
		telemetry.Stream(),
		otelgrpc.StreamServerInterceptor(),
	}

	options := []grpc.ServerOption{
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(middleware.ChainUnaryServer(unary...)),
		grpc.StreamInterceptor(middleware.ChainStreamServer(stream...)),
	}
	if s.opts.MaxMsgSize > 0 {
		options = append(options, grpc.MaxRecvMsgSize(s.opts.MaxMsgSize))
	}
	grpcServer := grpc.NewServer(options...)
	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(grpcServer, healthcheck)

	messagev1.RegisterMessageApiServer(grpcServer, s.messagev1)

	go func() {
		s.log.Info("serving grpc", zap.String("address", s.grpc.Addr().String()))
		err := grpcServer.Serve(s.grpc)
		if err != nil && !isErrUseOfClosedConnection(err) {
			s.log.Error("serving grpc", zap.Error(err))
		}
	}()

	return nil
}

func (s *Server) startHTTP() error {
	conn, err := s.dialGRPC(s.ctx)
	if err != nil {
		return errors.Wrap(err, "dialing grpc server")
	}

	mux := http.NewServeMux()
	healthClient := healthgrpc.NewHealthClient(conn)
	gwmux := runtime.NewServeMux(
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
		runtime.WithStreamErrorHandler(runtime.DefaultStreamErrorHandler),
		runtime.WithIncomingHeaderMatcher(incomingHeaderMatcher),
		runtime.WithHealthzEndpoint(healthClient),
	)
	mux.Handle("/", gwmux)

	err = messagev1.RegisterMessageApiHandler(s.ctx, gwmux, conn)
	if err != nil {
		return errors.Wrap(err, "registering message handler")
	}

	addr := hostPortAddr(s.opts.HTTPAddress, s.opts.HTTPPort)
	s.http, err = net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "creating grpc-gateway listener")
	}

	server := http.Server{
		Addr:    addr,
		Handler: allowCORS(mux),
	}

	go func() {
		s.log.Info("serving http", zap.String("address", s.http.Addr().String()))
		err = server.Serve(s.http)
		if err != nil && err != http.ErrServerClosed && !isErrUseOfClosedConnection(err) {
			s.log.Error("serving http", zap.Error(err))
		}
	}()

	return nil
}

func (s *Server) Close() {
	if s.http != nil {
		err := s.http.Close()
		if err != nil {
			s.log.Error("closing http listener", zap.Error(err))
		}
	}

	if s.grpc != nil {
		err := s.grpc.Close()
		if err != nil {
			s.log.Error("closing grpc listener", zap.Error(err))
		}
	}
}

func (s *Server) HTTPListenPort() uint {
	return uint(s.http.Addr().(*net.TCPAddr).Port)
}

func (s *Server) dialGRPC(ctx context.Context) (*grpc.ClientConn, error) {
	// https://github.com/grpc/grpc/blob/master/doc/naming.md
	dialAddr := fmt.Sprintf("passthrough://localhost/%s", s.grpc.Addr().String())
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	if s.opts.MaxMsgSize > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(s.opts.MaxMsgSize),
		))
	}
	return grpc.DialContext(ctx, dialAddr, opts...)
}

func (s *Server) HTTPListenAddr() string {
	return "http://" + s.http.Addr().String()
}

func isErrUseOfClosedConnection(err error) bool {
	return strings.Contains(err.Error(), "use of closed network connection")
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

func incomingHeaderMatcher(key string) (string, bool) {
	switch strings.ToLower(key) {
	case clientVersionMetadataKey:
		return key, true
	case appVersionMetadataKey:
		return key, true
	default:
		return key, false
	}
}
