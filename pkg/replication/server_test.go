package replication

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/replication/registry"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

const WRITER_DB_CONNECTION_STRING = "postgres://postgres:xmtp@localhost:8765/postgres?sslmode=disable"

func NewTestServer(t *testing.T, registry registry.NodeRegistry) *Server {
	log := test.NewLog(t)
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	server, err := New(context.Background(), log, Options{
		PrivateKeyString: hex.EncodeToString(crypto.FromECDSA(privateKey)),
		API: ApiOptions{
			Port: 0,
		},
		DB: DbOptions{
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
	registry := registry.NewFixedNodeRegistry([]registry.Node{})
	server1 := NewTestServer(t, registry)
	server2 := NewTestServer(t, registry)
	require.NotEqual(t, server1.Addr(), server2.Addr())
}

func TestMigrate(t *testing.T) {
	registry := registry.NewFixedNodeRegistry([]registry.Node{})
	server1 := NewTestServer(t, registry)
	require.NoError(t, server1.Migrate())
}
