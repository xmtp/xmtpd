package server

import (
	"crypto/ecdsa"
	"database/sql"
	"encoding/hex"
	"net"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/config"
	r "github.com/xmtp/xmtpd/pkg/registry"
	s "github.com/xmtp/xmtpd/pkg/server"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

type EnabledServices struct {
	Indexer     bool
	Sync        bool
	Replication bool
}

type TestServerCfg struct {
	GRPCListener     net.Listener
	HTTPListener     net.Listener
	Db               *sql.DB
	Registry         r.NodeRegistry
	PrivateKey       *ecdsa.PrivateKey
	ContractsOptions config.ContractsOptions
	Services         EnabledServices
}

func NewTestServer(
	t *testing.T,
	cfg TestServerCfg,
) *s.ReplicationServer {
	log := testutils.NewLog(t)

	server, err := s.NewReplicationServer(s.WithContext(t.Context()),
		s.WithLogger(log),
		s.WithDB(cfg.Db),
		s.WithNodeRegistry(cfg.Registry),
		s.WithServerVersion(testutils.GetLatestVersion(t)),
		s.WithGRPCListener(cfg.GRPCListener),
		s.WithHTTPListener(cfg.HTTPListener),
		s.WithServerOptions(&config.ServerOptions{
			Contracts: cfg.ContractsOptions,
			MlsValidation: config.MlsValidationOptions{
				GrpcAddress: "http://localhost:60051",
			},
			Signer: config.SignerOptions{
				PrivateKey: hex.EncodeToString(crypto.FromECDSA(cfg.PrivateKey)),
			},
			Sync: config.SyncOptions{
				Enable: cfg.Services.Sync,
			},
			Replication: config.ReplicationOptions{
				Enable:                cfg.Services.Replication,
				SendKeepAliveInterval: 30 * time.Second,
			},
			Indexer: config.IndexerOptions{
				Enable: cfg.Services.Indexer,
			},
		}))
	require.NoError(t, err)

	return server
}
