package server_test

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
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	r "github.com/xmtp/xmtpd/pkg/registry"
	s "github.com/xmtp/xmtpd/pkg/server"
	"github.com/xmtp/xmtpd/pkg/testutils"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	networkTestUtils "github.com/xmtp/xmtpd/pkg/testutils/network"
	registryTestUtils "github.com/xmtp/xmtpd/pkg/testutils/registry"
	"github.com/xmtp/xmtpd/pkg/topic"
)

const (
	server1NodeID = uint32(100)
	server2NodeID = uint32(200)
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

	server, err := s.NewReplicationServer(s.WithContext(t.Context()),
		s.WithLogger(log),
		s.WithDB(db),
		s.WithNodeRegistry(registry),
		s.WithServerVersion(testutils.GetLatestVersion(t)),
		s.WithListenAddress(fmt.Sprintf("localhost:%d", port)),
		s.WithHTTPListenAddress(fmt.Sprintf("localhost:%d", httpPort)),
		s.WithServerOptions(&config.ServerOptions{
			Contracts: config.ContractsOptions{
				AppChain: config.AppChainOptions{
					WssURL:                 "ws://localhost:8545",
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
		}))
	require.NoError(t, err)

	return server
}

func TestCreateServer(t *testing.T) {
	ctx := t.Context()
	dbs := testutils.NewDBs(t, ctx, 2)
	privateKey1, err := crypto.GenerateKey()
	require.NoError(t, err)
	privateKey2, err := crypto.GenerateKey()
	require.NoError(t, err)

	server1Port := networkTestUtils.FindFreePort(t)
	server2Port := networkTestUtils.FindFreePort(t)

	httpServer1Port := networkTestUtils.FindFreePort(t)
	httpServer2Port := networkTestUtils.FindFreePort(t)

	nodes := []r.Node{
		registryTestUtils.CreateNode(server1NodeID, server1Port, privateKey1),
		registryTestUtils.CreateNode(server2NodeID, server2Port, privateKey2),
	}

	registry := registryTestUtils.CreateMockRegistry(t, nodes)

	server1 := NewTestServer(t, server1Port, httpServer1Port, dbs[0], registry, privateKey1)
	server2 := NewTestServer(t, server2Port, httpServer2Port, dbs[1], registry, privateKey2)

	require.NotEqual(t, server1.Addr(), server2.Addr())

	defer func() {
		server1.Shutdown(0)
		server2.Shutdown(0)
	}()

	client1 := apiTestUtils.NewReplicationAPIClient(t, server1.Addr().String())
	client2 := apiTestUtils.NewReplicationAPIClient(t, server2.Addr().String())
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
	ctx := t.Context()
	dbs := testutils.NewDBs(t, ctx, 1)
	privateKey1, err := crypto.GenerateKey()
	require.NoError(t, err)
	server1Port := networkTestUtils.FindFreePort(t)
	httpServer1Port := networkTestUtils.FindFreePort(t)

	nodeId1 := server1NodeID

	nodes := []r.Node{registryTestUtils.CreateNode(server1NodeID, server1Port, privateKey1)}
	registry := registryTestUtils.CreateMockRegistry(t, nodes)

	server1 := NewTestServer(t, server1Port, httpServer1Port, dbs[0], registry, privateKey1)
	defer func() {
		server1.Shutdown(0)
	}()

	client1 := apiTestUtils.NewReplicationAPIClient(t, server1.Addr().String())

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
