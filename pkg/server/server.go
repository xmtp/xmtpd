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
	"github.com/xmtp/xmtpd/pkg/migrations"
	"github.com/xmtp/xmtpd/pkg/registry"
	"go.uber.org/zap"
)

type ReplicationServer struct {
	options      Options
	log          *zap.Logger
	ctx          context.Context
	cancel       context.CancelFunc
	apiServer    *api.ApiServer
	nodeRegistry registry.NodeRegistry
	privateKey   *ecdsa.PrivateKey
	writerDb     *sql.DB
	// Can add reader DB later if needed
}

func NewReplicationServer(ctx context.Context, log *zap.Logger, options Options, nodeRegistry registry.NodeRegistry) (*ReplicationServer, error) {
	var err error
	s := &ReplicationServer{
		options:      options,
		log:          log,
		nodeRegistry: nodeRegistry,
	}
	s.privateKey, err = parsePrivateKey(options.PrivateKeyString)
	if err != nil {
		return nil, err
	}
	// Commenting out the DB stuff until I get the new migrations in
	s.writerDb, err = db.NewDB(ctx, options.DB.WriterConnectionString, options.DB.WaitForDB, options.DB.ReadTimeout)
	if err != nil {
		return nil, err
	}

	s.ctx, s.cancel = context.WithCancel(ctx)
	s.apiServer, err = api.NewAPIServer(ctx, log, options.API.Port)
	if err != nil {
		return nil, err
	}
	log.Info("Replication server started", zap.Int("port", options.API.Port))
	return s, nil
}

func (s *ReplicationServer) Migrate() error {
	return migrations.Migrate(s.writerDb)
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
