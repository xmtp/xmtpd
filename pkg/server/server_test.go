package server

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/registry"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

const WRITER_DB_CONNECTION_STRING = "postgres://postgres:xmtp@localhost:8765/postgres?sslmode=disable"

func NewTestServer(t *testing.T, registry registry.NodeRegistry) *ReplicationServer {
	log := test.NewLog(t)
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	server, err := NewReplicationServer(context.Background(), log, config.ServerOptions{
		PrivateKeyString: hex.EncodeToString(crypto.FromECDSA(privateKey)),
		API: config.ApiOptions{
			Port: 0,
		},
		DB: config.DbOptions{
			WriterConnectionString: WRITER_DB_CONNECTION_STRING,
			ReadTimeout:            time.Second * 10,
			WriteTimeout:           time.Second * 10,
			MaxOpenConns:           10,
			WaitForDB:              time.Second * 10,
		},
	}, registry)
	require.NoError(t, err)

	return server
}

func TestCreateServer(t *testing.T) {
	registry := registry.NewFixedNodeRegistry([]registry.Record{})
	server1 := NewTestServer(t, registry)
	server2 := NewTestServer(t, registry)
	require.NotEqual(t, server1.Addr(), server2.Addr())
}
