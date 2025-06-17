package server

import (
	"crypto/ecdsa"
	"database/sql"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/config"
	r "github.com/xmtp/xmtpd/pkg/registry"
	s "github.com/xmtp/xmtpd/pkg/server"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func NewTestServer(
	t *testing.T,
	port int,
	httpPort int,
	db *sql.DB,
	registry r.NodeRegistry,
	privateKey *ecdsa.PrivateKey,
	contractsOptions config.ContractsOptions,
) *s.ReplicationServer {
	log := testutils.NewLog(t)

	server, err := s.NewReplicationServer(s.WithContext(t.Context()),
		s.WithLogger(log),
		s.WithDB(db),
		s.WithNodeRegistry(registry),
		s.WithServerVersion(testutils.GetLatestVersion(t)),
		s.WithListenAddress(fmt.Sprintf("localhost:%d", port)),
		s.WithHTTPListenAddress(fmt.Sprintf("localhost:%d", httpPort)),
		s.WithServerOptions(&config.ServerOptions{
			Contracts: contractsOptions,
			MlsValidation: config.MlsValidationOptions{
				GrpcAddress: "http://localhost:60051",
			},
			Signer: config.SignerOptions{
				PrivateKey: hex.EncodeToString(crypto.FromECDSA(privateKey)),
			},
			API: config.ApiOptions{
				Port:     port,
				HTTPPort: httpPort,
			},
			Sync: config.SyncOptions{
				Enable: true,
			},
			Replication: config.ReplicationOptions{
				Enable:                true,
				SendKeepAliveInterval: 30 * time.Second,
			},
		}))
	require.NoError(t, err)

	return server
}
