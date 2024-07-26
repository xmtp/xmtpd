package replication

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/migrate"
	"github.com/xmtp/xmtpd/pkg/migrations/replication"
	"github.com/xmtp/xmtpd/pkg/replication/api"
	"github.com/xmtp/xmtpd/pkg/replication/registry"
	legacyServer "github.com/xmtp/xmtpd/pkg/server"
	"go.uber.org/zap"
)

type Server struct {
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

func New(ctx context.Context, log *zap.Logger, options Options, nodeRegistry registry.NodeRegistry) (*Server, error) {
	var err error
	s := &Server{
		options:      options,
		log:          log,
		nodeRegistry: nodeRegistry,
	}
	s.privateKey, err = parsePrivateKey(options.PrivateKeyString)
	if err != nil {
		return nil, err
	}
	s.writerDb, err = getWriterDb(options.DB)
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

func (s *Server) Migrate() error {
	db := bun.NewDB(s.writerDb, pgdialect.New())
	migrator := migrate.NewMigrator(db, replication.Migrations)
	err := migrator.Init(s.ctx)
	if err != nil {
		return err
	}

	group, err := migrator.Migrate(s.ctx)
	if err != nil {
		return err
	}
	if group.IsZero() {
		s.log.Info("No new migrations to run")
	}

	return nil
}

func (s *Server) Addr() net.Addr {
	return s.apiServer.Addr()
}

func (s *Server) WaitForShutdown() {
	termChannel := make(chan os.Signal, 1)
	signal.Notify(termChannel, syscall.SIGINT, syscall.SIGTERM)
	<-termChannel
	s.Shutdown()
}

func (s *Server) Shutdown() {
	s.cancel()
	if s.apiServer != nil {
		s.apiServer.Close()
	}
}

func parsePrivateKey(privateKeyString string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(privateKeyString)
}

func getWriterDb(options DbOptions) (*sql.DB, error) {
	db, err := legacyServer.NewPGXDB(options.WriterConnectionString, options.WaitForDB, options.ReadTimeout)
	if err != nil {
		return nil, err
	}

	return db, nil
}
