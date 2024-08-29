package server

import (
	"context"
	"database/sql"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/xmtp/xmtpd/pkg/api"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"go.uber.org/zap"
)

type ReplicationServer struct {
	apiServer    *api.ApiServer
	ctx          context.Context
	cancel       context.CancelFunc
	log          *zap.Logger
	registrant   *registrant.Registrant
	nodeRegistry registry.NodeRegistry
	options      config.ServerOptions
	writerDB     *sql.DB
	// Can add reader DB later if needed
}

func NewReplicationServer(
	ctx context.Context,
	log *zap.Logger,
	options config.ServerOptions,
	nodeRegistry registry.NodeRegistry,
	writerDB *sql.DB,
) (*ReplicationServer, error) {
	var err error
	s := &ReplicationServer{
		options:      options,
		log:          log,
		nodeRegistry: nodeRegistry,
		writerDB:     writerDB,
	}

	s.registrant, err = registrant.NewRegistrant(
		ctx,
		queries.New(s.writerDB),
		nodeRegistry,
		options.SignerPrivateKey,
	)
	if err != nil {
		return nil, err
	}

	s.ctx, s.cancel = context.WithCancel(ctx)
	s.apiServer, err = api.NewAPIServer(ctx, s.writerDB, log, options.API.Port, s.registrant)
	if err != nil {
		return nil, err
	}
	log.Info("Replication server started", zap.Int("port", options.API.Port))
	return s, nil
}

func (s *ReplicationServer) Addr() net.Addr {
	return s.apiServer.Addr()
}

func (s *ReplicationServer) WaitForShutdown() {
	termChannel := make(chan os.Signal, 1)
	signal.Notify(termChannel, syscall.SIGINT, syscall.SIGTERM)
	<-termChannel
	s.Shutdown()
}

func (s *ReplicationServer) Shutdown() {
	s.cancel()
	if s.apiServer != nil {
		s.apiServer.Close()
	}
}
