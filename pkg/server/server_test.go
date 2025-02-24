package server_test

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/config"
	mocks "github.com/xmtp/xmtpd/pkg/mocks/registry"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	r "github.com/xmtp/xmtpd/pkg/registry"
	s "github.com/xmtp/xmtpd/pkg/server"
	"github.com/xmtp/xmtpd/pkg/testutils"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
)

const server1NodeID = uint32(100)
const server2NodeID = uint32(200)

func getNextOpenPort() (int, error) {
	// Listen on a random available port
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, fmt.Errorf("could not find open port: %w", err)
	}
	defer listener.Close()

	// Extract the port number from the listener
	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

func NewTestServer(
	t *testing.T,
	port int,
	db *sql.DB,
	registry r.NodeRegistry,
	privateKey *ecdsa.PrivateKey,
) *s.ReplicationServer {
	log := testutils.NewLog(t)

	server, err := s.NewReplicationServer(context.Background(), log, config.ServerOptions{
		Contracts: config.ContractsOptions{
			RpcUrl:                 "http://localhost:8545",
			MaxChainDisconnectTime: 5 * time.Minute,
		},
		MlsValidation: config.MlsValidationOptions{
			GrpcAddress: "http://localhost:60051",
		},
		Signer: config.SignerOptions{
			PrivateKey: hex.EncodeToString(crypto.FromECDSA(privateKey)),
		},
		API: config.ApiOptions{
			Port: port,
		},
		Sync: config.SyncOptions{
			Enable: true,
		},
		Replication: config.ReplicationOptions{
			Enable: true,
		},
	}, registry, db, fmt.Sprintf("localhost:%d", port), testutils.GetLatestVersion(t))
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

	server1Port, err := getNextOpenPort()
	require.NoError(t, err)
	server2Port, err := getNextOpenPort()
	require.NoError(t, err)

	nodes := []r.Node{
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
		}}

	registry := mocks.NewMockNodeRegistry(t)
	registry.On("GetNodes").Return(nodes, nil)

	nodesChan := make(chan []r.Node)
	registry.On("OnNewNodes").
		Return((<-chan []r.Node)(nodesChan), r.CancelSubscription(func() {}))

	nodeChan1 := make(chan r.Node)
	nodeChan2 := make(chan r.Node)
	registry.On("OnChangedNode", server1NodeID).
		Return((<-chan r.Node)(nodeChan1), r.CancelSubscription(func() {
			close(nodeChan1)
		}))
	registry.On("OnChangedNode", server2NodeID).
		Return((<-chan r.Node)(nodeChan2), r.CancelSubscription(func() {
			close(nodeChan2)
		}))

	registry.On("GetNode", server1NodeID).Return(&nodes[0], nil)
	registry.On("GetNode", server2NodeID).Return(&nodes[1], nil)

	registry.On("Stop").Return(nil)

	server1 := NewTestServer(t, server1Port, dbs[0], registry, privateKey1)
	server2 := NewTestServer(t, server2Port, dbs[1], registry, privateKey2)

	require.NotEqual(t, server1.Addr(), server2.Addr())

	defer func() {
		server1.Shutdown(0)
		server2.Shutdown(0)
	}()

	client1, cleanup1 := apiTestUtils.NewReplicationAPIClient(t, ctx, server1.Addr().String())
	defer cleanup1()
	client2, cleanup2 := apiTestUtils.NewReplicationAPIClient(t, ctx, server2.Addr().String())
	defer cleanup2()
	nodeId1 := server1NodeID
	nodeId2 := server2NodeID

	targetTopic := topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, []byte{1, 2, 3}).
		Bytes()

	p1, err := client1.PublishPayerEnvelopes(
		ctx,
		&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{envelopeTestUtils.CreatePayerEnvelope(
				t,
				nodeId1,
				envelopeTestUtils.CreateClientEnvelope(&envelopes.AuthenticatedData{
					TargetOriginator: &nodeId1,
					TargetTopic:      targetTopic,
					DependsOn:        &envelopes.Cursor{},
				}),
			)},
		},
	)
	require.NoError(t, err)
	p2, err := client2.PublishPayerEnvelopes(
		ctx,
		&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{envelopeTestUtils.CreatePayerEnvelope(
				t,
				nodeId2,
				envelopeTestUtils.CreateClientEnvelope(&envelopes.AuthenticatedData{
					TargetOriginator: &nodeId2,
					TargetTopic:      targetTopic,
					DependsOn:        &envelopes.Cursor{},
				}),
			)},
		},
	)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		q1, err := client1.QueryEnvelopes(ctx, &message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{server2NodeID},
				LastSeen:          &envelopes.Cursor{},
			},
			Limit: 10,
		})
		require.NoError(t, err)
		if len(q1.Envelopes) == 0 {
			return false
		}
		require.Len(t, q1.Envelopes, 1)
		require.Equal(t, q1.Envelopes[0], p2.OriginatorEnvelopes[0])
		return true
	}, 3000*time.Millisecond, 200*time.Millisecond)

	require.Eventually(t, func() bool {
		q2, err := client2.QueryEnvelopes(ctx, &message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{server1NodeID},
				LastSeen:          &envelopes.Cursor{},
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
	}, 3000*time.Millisecond, 200*time.Millisecond)
}

func TestReadOwnWritesGuarantee(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dbs, dbCleanup := testutils.NewDBs(t, ctx, 1)
	defer dbCleanup()
	privateKey1, err := crypto.GenerateKey()
	require.NoError(t, err)
	server1Port, err := getNextOpenPort()
	require.NoError(t, err)
	nodeId1 := server1NodeID

	nodes := []r.Node{
		{
			NodeID:        server1NodeID,
			SigningKey:    &privateKey1.PublicKey,
			HttpAddress:   fmt.Sprintf("http://localhost:%d", server1Port),
			IsHealthy:     true,
			IsValidConfig: true,
		}}

	registry := mocks.NewMockNodeRegistry(t)
	registry.On("GetNodes").Return(nodes, nil)

	nodesChan := make(chan []r.Node)
	registry.On("OnNewNodes").
		Return((<-chan []r.Node)(nodesChan), r.CancelSubscription(func() {
		}))

	registry.On("Stop").Return(nil)

	server1 := NewTestServer(t, server1Port, dbs[0], registry, privateKey1)
	defer func() {
		server1.Shutdown(0)
	}()

	client1, cleanup1 := apiTestUtils.NewReplicationAPIClient(t, ctx, server1.Addr().String())
	defer cleanup1()

	targetTopic := topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, []byte{1, 2, 3}).
		Bytes()

	_, err = client1.PublishPayerEnvelopes(
		ctx,
		&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{envelopeTestUtils.CreatePayerEnvelope(
				t,
				nodeId1,
				envelopeTestUtils.CreateClientEnvelope(&envelopes.AuthenticatedData{
					TargetTopic: targetTopic,
					DependsOn:   &envelopes.Cursor{},
				}),
			)},
		},
	)
	require.NoError(t, err)

	// query the same server immediately after writing
	// the server should return the write on the first attempt

	q1, err := client1.QueryEnvelopes(ctx, &message_api.QueryEnvelopesRequest{
		Query: &message_api.EnvelopesQuery{
			OriginatorNodeIds: []uint32{server1NodeID},
			LastSeen:          &envelopes.Cursor{},
		},
		Limit: 10,
	})
	require.NoError(t, err)
	require.Len(t, q1.Envelopes, 1)
}
