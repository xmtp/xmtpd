package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	proto "github.com/xmtp/proto/v3/go/message_api/v1"
	messagev1 "github.com/xmtp/xmtpd/pkg/api/message/v1"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	opts *Options
	ctx  context.Context
	log  *zap.Logger

	grpc      net.Listener
	http      net.Listener
	messagev1 *messagev1.Service
}

func New(ctx context.Context, log *zap.Logger, opts *Options) (*Server, error) {
	err := opts.validate()
	if err != nil {
		return nil, err
	}

	log = log.Named("api")

	s := &Server{
		ctx:  ctx,
		log:  log,
		opts: opts,
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

	prometheus.EnableHandlingTimeHistogram()
	unary := []grpc.UnaryServerInterceptor{prometheus.UnaryServerInterceptor}
	stream := []grpc.StreamServerInterceptor{prometheus.StreamServerInterceptor}

	telemetry := NewTelemetryInterceptor(s.log)
	unary = append(unary, telemetry.Unary())
	stream = append(stream, telemetry.Stream())

	options := []grpc.ServerOption{
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(middleware.ChainUnaryServer(unary...)),
		grpc.StreamInterceptor(middleware.ChainStreamServer(stream...)),
		grpc.MaxRecvMsgSize(s.opts.MaxMsgSize),
	}
	grpcServer := grpc.NewServer(options...)
	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(grpcServer, healthcheck)

	s.messagev1, err = messagev1.NewService(s.log)
	if err != nil {
		return errors.Wrap(err, "creating message service")
	}
	proto.RegisterMessageApiServer(grpcServer, s.messagev1)
	prometheus.Register(grpcServer)

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

	err = proto.RegisterMessageApiHandler(s.ctx, gwmux, conn)
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
	if s.messagev1 != nil {
		s.messagev1.Close()
	}

	if s.http != nil {
		err := s.http.Close()
		if err != nil {
			s.log.Error("closing http listener", zap.Error(err))
		}
	}

	if s.http != nil {
		err := s.http.Close()
		if err != nil {
			s.log.Error("closing grpc listener", zap.Error(err))
		}
	}
}

func (s *Server) dialGRPC(ctx context.Context) (*grpc.ClientConn, error) {
	// https://github.com/grpc/grpc/blob/master/doc/naming.md
	dialAddr := fmt.Sprintf("passthrough://localhost/%s", s.grpc.Addr().String())
	return grpc.DialContext(
		ctx,
		dialAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(s.opts.MaxMsgSize),
		),
	)
}

func (s *Server) httpListenAddr() string {
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
