package server_test

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
	"github.com/xmtp/xmtpd/pkg/mocks/blockchain"
	mocks "github.com/xmtp/xmtpd/pkg/mocks/registry"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	r "github.com/xmtp/xmtpd/pkg/registry"
	s "github.com/xmtp/xmtpd/pkg/server"
	"github.com/xmtp/xmtpd/pkg/testutils"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
	"github.com/xmtp/xmtpd/pkg/utils"
)

const server1NodeID = 100
const server2NodeID = 200
const server1Port = 1111
const server2Port = 2222

func NewTestServer(
	t *testing.T,
	port int,
	db *sql.DB,
	registry r.NodeRegistry,
	privateKey *ecdsa.PrivateKey,
) *s.ReplicationServer {
	log := testutils.NewLog(t)
	messagePublisher := blockchain.NewMockIBlockchainPublisher(t)

	server, err := s.NewReplicationServer(context.Background(), log, config.ServerOptions{
		Contracts: config.ContractsOptions{
			RpcUrl: "http://localhost:8545",
		},
		MlsValidation: config.MlsValidationOptions{
			GrpcAddress: "localhost:60051",
		},
		Signer: config.SignerOptions{
			PrivateKey: hex.EncodeToString(crypto.FromECDSA(privateKey)),
		},
		API: config.ApiOptions{
			Port: port,
		},
	}, registry, db, messagePublisher)
	require.NoError(t, err)

	return server
}

func TestCreateServer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dbs, dbCleanup := testutils.NewDBs(t, ctx, 2)
	defer dbCleanup()
	privateKey1, err := crypto.GenerateKey()
	require.NoError(t, err)
	privateKey2, err := crypto.GenerateKey()
	require.NoError(t, err)

	registry := mocks.NewMockNodeRegistry(t)
	registry.On("GetNodes").Return([]r.Node{
		{
			NodeID:        server1NodeID,
			SigningKey:    &privateKey1.PublicKey,
			HttpAddress:   fmt.Sprintf("http://localhost:%d", server1Port),
			IsHealthy:     true,
			IsValidConfig: true,
		},
		{
			NodeID:        server2NodeID,
			SigningKey:    &privateKey2.PublicKey,
			HttpAddress:   fmt.Sprintf("http://localhost:%d", server2Port),
			IsHealthy:     true,
			IsValidConfig: true,
		},
	}, nil)

	server1 := NewTestServer(t, server1Port, dbs[0], registry, privateKey1)
	server2 := NewTestServer(t, server2Port, dbs[1], registry, privateKey2)
	require.NotEqual(t, server1.Addr(), server2.Addr())

	client1, cleanup1 := apiTestUtils.NewAPIClient(t, ctx, server1.Addr().String())
	defer cleanup1()
	client2, cleanup2 := apiTestUtils.NewAPIClient(t, ctx, server2.Addr().String())
	defer cleanup2()

	p1, err := client1.PublishEnvelopes(ctx, &message_api.PublishEnvelopesRequest{
		PayerEnvelopes: []*message_api.PayerEnvelope{testutils.CreatePayerEnvelope(
			t,
			testutils.CreateClientEnvelope(&message_api.AuthenticatedData{
				TargetOriginator: server1NodeID,
				TargetTopic:      []byte{0x5},
				LastSeen:         &message_api.VectorClock{},
			}),
		)},
	})
	require.NoError(t, err)
	p2, err := client2.PublishEnvelopes(ctx, &message_api.PublishEnvelopesRequest{
		PayerEnvelopes: []*message_api.PayerEnvelope{testutils.CreatePayerEnvelope(
			t,
			testutils.CreateClientEnvelope(&message_api.AuthenticatedData{
				TargetOriginator: server2NodeID,
				TargetTopic:      []byte{0x5},
				LastSeen:         &message_api.VectorClock{},
			}),
		)},
	})
	utils.Unused(p2)
	require.NoError(t, err)

	// TODO(rich) synchronize test with stream setup
	// Client1 is initialized before Client2 and needs to wait for Client2 to be ready before
	// stream will succeed, and will retry a few times. If the publish on client2 goes through before
	// client1 is set up, the condition below will never succeed. Easiest solution is to implement
	// subscriptions starting from a cursor.
	// require.Eventually(t, func() bool {
	// 	q1, err := client1.QueryEnvelopes(ctx, &message_api.QueryEnvelopesRequest{
	// 		Query: &message_api.EnvelopesQuery{
	// 			OriginatorNodeIds: []uint32{server2NodeID},
	// 			LastSeen: &message_api.VectorClock{},
	// 		},
	// 		Limit: 10,
	// 	})
	// 	require.NoError(t, err)
	// 	if len(q1.Envelopes) == 0 {
	// 		return false
	// 	}
	// 	require.Len(t, q1.Envelopes, 1)
	// 	require.Equal(t, q1.Envelopes[0], p2.OriginatorEnvelopes[0])
	// 	return true
	// }, 500*time.Millisecond, 50*time.Millisecond)

	require.Eventually(t, func() bool {
		q2, err := client2.QueryEnvelopes(ctx, &message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{server1NodeID},
				LastSeen:          &message_api.VectorClock{},
			},
			Limit: 10,
		})
		require.NoError(t, err)
		if len(q2.Envelopes) == 0 {
			return false
		}
		require.Len(t, q2.Envelopes, 1)
		require.Equal(t, q2.Envelopes[0], p1.OriginatorEnvelopes[0])
		return true
	}, 500*time.Millisecond, 50*time.Millisecond)
}
