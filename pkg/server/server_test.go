package server_test

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/mocks"
	r "github.com/xmtp/xmtpd/pkg/registry"
	s "github.com/xmtp/xmtpd/pkg/server"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func NewTestServer(
	t *testing.T,
	db *sql.DB,
	registry r.NodeRegistry,
	privateKey *ecdsa.PrivateKey,
) *s.ReplicationServer {
	log := testutils.NewLog(t)

	server, err := s.NewReplicationServer(context.Background(), log, config.ServerOptions{
		SignerPrivateKey: hex.EncodeToString(crypto.FromECDSA(privateKey)),
		API: config.ApiOptions{
			Port: 0,
		},
	}, registry, db)
	require.NoError(t, err)

	return server
}

func TestCreateServer(t *testing.T) {
	dbs, dbCleanup := testutils.NewDBs(t, context.Background(), 2)
	defer dbCleanup()
	privateKey1, err := crypto.GenerateKey()
	require.NoError(t, err)
	privateKey2, err := crypto.GenerateKey()
	require.NoError(t, err)

	registry := mocks.NewMockNodeRegistry(t)
	registry.On("GetNodes").Return([]r.Node{
		{NodeID: 1, SigningKey: &privateKey1.PublicKey},
		{NodeID: 2, SigningKey: &privateKey2.PublicKey},
	}, nil)

	server1 := NewTestServer(t, dbs[0], registry, privateKey1)
	server2 := NewTestServer(t, dbs[1], registry, privateKey2)

	require.NotEqual(t, server1.Addr(), server2.Addr())
}
