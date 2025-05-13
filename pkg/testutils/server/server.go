package server

import (
	"context"
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
) *s.ReplicationServer {
	log := testutils.NewLog(t)

	server, err := s.NewReplicationServer(context.Background(), log, config.ServerOptions{
		Contracts: config.ContractsOptions{
			AppChain: config.AppChainOptions{
				RpcURL:                 "http://localhost:8545",
				MaxChainDisconnectTime: 5 * time.Minute,
			},
		},
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
	}, registry, db, fmt.Sprintf("localhost:%d", port), fmt.Sprintf("localhost:%d", httpPort), testutils.GetLatestVersion(t))
	require.NoError(t, err)

	return server
}
