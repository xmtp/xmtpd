package server_test

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/mocks/blockchain"
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
		//TODO(mkysel): this is not fully mocked yet
		//Payer: config.PayerOptions{
		//	Enable: true,
		//},
	}, registry, db, messagePublisher, fmt.Sprintf("localhost:%d", port), nil)
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
	cancelOnNewFunc := func() {
		close(nodesChan)
	}
	registry.On("OnNewNodes").
		Return((<-chan []r.Node)(nodesChan), r.CancelSubscription(cancelOnNewFunc))

	nodeChan := make(chan r.Node)

	cancelOnChangedFunc := func() {
		close(nodeChan)
	}

	registry.On("OnChangedNode", mock.AnythingOfType("uint32")).
		Return((<-chan r.Node)(nodeChan), r.CancelSubscription(cancelOnChangedFunc))

	registry.On("GetNode", server1NodeID).Return(&nodes[0], nil)
	registry.On("GetNode", server2NodeID).Return(&nodes[1], nil)

	server1 := NewTestServer(t, server1Port, dbs[0], registry, privateKey1)
	server2 := NewTestServer(t, server2Port, dbs[1], registry, privateKey2)
	require.NotEqual(t, server1.Addr(), server2.Addr())

	client1, cleanup1 := apiTestUtils.NewReplicationAPIClient(t, ctx, server1.Addr().String())
	defer cleanup1()
	client2, cleanup2 := apiTestUtils.NewReplicationAPIClient(t, ctx, server2.Addr().String())
	defer cleanup2()

	targetTopic := topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, []byte{1, 2, 3}).
		Bytes()

	p1, err := client1.PublishPayerEnvelopes(
		ctx,
		&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{envelopeTestUtils.CreatePayerEnvelope(
				t,
				envelopeTestUtils.CreateClientEnvelope(&envelopes.AuthenticatedData{
					TargetOriginator: server1NodeID,
					TargetTopic:      targetTopic,
					LastSeen:         &envelopes.VectorClock{},
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
				envelopeTestUtils.CreateClientEnvelope(&envelopes.AuthenticatedData{
					TargetOriginator: server2NodeID,
					TargetTopic:      targetTopic,
					LastSeen:         &envelopes.VectorClock{},
				}),
			)},
		},
	)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		q1, err := client1.QueryEnvelopes(ctx, &message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{server2NodeID},
				LastSeen:          &envelopes.VectorClock{},
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
				LastSeen:          &envelopes.VectorClock{},
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
