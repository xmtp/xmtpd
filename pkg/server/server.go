package server

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/api"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/node"
	"github.com/xmtp/xmtpd/pkg/registry"
	"go.uber.org/zap"
)

type ReplicationServer struct {
	apiServer    *api.ApiServer
	ctx          context.Context
	cancel       context.CancelFunc
	log          *zap.Logger
	node         *node.Node
	nodeRegistry registry.NodeRegistry
	options      Options
	writerDB     *sql.DB
	// Can add reader DB later if needed
}

func NewReplicationServer(
	ctx context.Context,
	log *zap.Logger,
	options Options,
	nodeRegistry registry.NodeRegistry,
) (*ReplicationServer, error) {
	var err error
	s := &ReplicationServer{
		options:      options,
		log:          log,
		nodeRegistry: nodeRegistry,
	}
	s.writerDB, err = db.NewDB(
		ctx,
		options.DB.WriterConnectionString,
		options.DB.WaitForDB,
		options.DB.ReadTimeout,
	)
	if err != nil {
		return nil, err
	}

	s.node, err = node.NewNode(ctx, queries.New(s.writerDB), nodeRegistry, options.PrivateKeyString)
	if err != nil {
		return nil, err
	}

	s.ctx, s.cancel = context.WithCancel(ctx)
	s.apiServer, err = api.NewAPIServer(ctx, s.writerDB, log, s.node, options.API.Port)
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

func parsePrivateKey(privateKeyString string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(privateKeyString)
}
