package api

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/mocks"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/protobuf/proto"
)

func newTestService(t *testing.T) (*Service, *sql.DB, func()) {
	ctx := context.Background()
	log := testutils.NewLog(t)
	db, _, dbCleanup := testutils.NewDB(t, ctx)
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	privKeyStr := "0x" + testutils.HexEncode(crypto.FromECDSA(privKey))
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.EXPECT().GetNodes().Return([]registry.Node{
		{NodeID: 1, SigningKey: &privKey.PublicKey},
	}, nil)
	registrant, err := registrant.NewRegistrant(ctx, queries.New(db), mockRegistry, privKeyStr)
	require.NoError(t, err)

	svc, err := NewReplicationApiService(ctx, log, registrant, db)
	require.NoError(t, err)

	return svc, db, func() {
		svc.Close()
		dbCleanup()
	}
}

func TestPublishEnvelope(t *testing.T) {
	svc, db, cleanup := newTestService(t)
	defer cleanup()

	resp, err := svc.PublishEnvelope(
		context.Background(),
		&message_api.PublishEnvelopeRequest{
			PayerEnvelope: testutils.CreatePayerEnvelope(t),
		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	unsignedEnv := &message_api.UnsignedOriginatorEnvelope{}
	require.NoError(
		t,
		proto.Unmarshal(resp.GetOriginatorEnvelope().GetUnsignedOriginatorEnvelope(), unsignedEnv),
	)
	clientEnv := &message_api.ClientEnvelope{}
	require.NoError(
		t,
		proto.Unmarshal(unsignedEnv.GetPayerEnvelope().GetUnsignedClientEnvelope(), clientEnv),
	)
	require.Equal(t, uint8(0x5), clientEnv.Aad.GetTargetTopic()[0])

	// Check that the envelope was published to the database after a delay
	require.Eventually(t, func() bool {
		envs, err := queries.New(db).
			SelectGatewayEnvelopes(context.Background(), queries.SelectGatewayEnvelopesParams{})
		require.NoError(t, err)

		if len(envs) != 1 {
			return false
		}

		originatorEnv := &message_api.OriginatorEnvelope{}
		require.NoError(t, proto.Unmarshal(envs[0].OriginatorEnvelope, originatorEnv))
		return proto.Equal(originatorEnv, resp.GetOriginatorEnvelope())
	}, 500*time.Millisecond, 50*time.Millisecond)
}

func TestUnmarshalErrorOnPublish(t *testing.T) {
	svc, _, cleanup := newTestService(t)
	defer cleanup()

	envelope := testutils.CreatePayerEnvelope(t)
	envelope.UnsignedClientEnvelope = []byte("invalidbytes")
	_, err := svc.PublishEnvelope(
		context.Background(),
		&message_api.PublishEnvelopeRequest{
			PayerEnvelope: envelope,
		},
	)
	require.ErrorContains(t, err, "unmarshal")
}

func TestMismatchingOriginatorOnPublish(t *testing.T) {
	svc, _, cleanup := newTestService(t)
	defer cleanup()

	clientEnv := testutils.CreateClientEnvelope()
	clientEnv.Aad.TargetOriginator = 2
	_, err := svc.PublishEnvelope(
		context.Background(),
		&message_api.PublishEnvelopeRequest{
			PayerEnvelope: testutils.CreatePayerEnvelope(t, clientEnv),
		},
	)
	require.ErrorContains(t, err, "originator")
}

func TestMissingTopicOnPublish(t *testing.T) {
	svc, _, cleanup := newTestService(t)
	defer cleanup()

	clientEnv := testutils.CreateClientEnvelope()
	clientEnv.Aad.TargetTopic = nil
	_, err := svc.PublishEnvelope(
		context.Background(),
		&message_api.PublishEnvelopeRequest{
			PayerEnvelope: testutils.CreatePayerEnvelope(t, clientEnv),
		},
	)
	require.ErrorContains(t, err, "topic")
}

func setupQueryTest(t *testing.T, db *sql.DB) []queries.InsertGatewayEnvelopeParams {
	db_rows := []queries.InsertGatewayEnvelopeParams{
		{
			// Auto-generated ID: 1
			OriginatorNodeID:     1,
			OriginatorSequenceID: 1,
			Topic:                []byte("topicA"),
			OriginatorEnvelope: testutils.Marshal(
				t,
				testutils.CreateOriginatorEnvelope(t, 1, 1),
			),
		},
		{
			// Auto-generated ID: 2
			OriginatorNodeID:     2,
			OriginatorSequenceID: 1,
			Topic:                []byte("topicA"),
			OriginatorEnvelope: testutils.Marshal(
				t,
				testutils.CreateOriginatorEnvelope(t, 2, 1),
			),
		},
		{
			// Auto-generated ID: 3
			OriginatorNodeID:     1,
			OriginatorSequenceID: 2,
			Topic:                []byte("topicB"),
			OriginatorEnvelope: testutils.Marshal(
				t,
				testutils.CreateOriginatorEnvelope(t, 1, 2),
			),
		},
		{
			// Auto-generated ID: 4
			OriginatorNodeID:     2,
			OriginatorSequenceID: 2,
			Topic:                []byte("topicB"),
			OriginatorEnvelope: testutils.Marshal(
				t,
				testutils.CreateOriginatorEnvelope(t, 2, 2),
			),
		},
		{
			// Auto-generated ID: 5
			OriginatorNodeID:     1,
			OriginatorSequenceID: 3,
			Topic:                []byte("topicA"),
			OriginatorEnvelope: testutils.Marshal(
				t,
				testutils.CreateOriginatorEnvelope(t, 1, 3),
			),
		},
	}
	testutils.InsertGatewayEnvelopes(t, db, db_rows)
	return db_rows
}

func TestQueryAllEnvelopes(t *testing.T) {
	svc, db, cleanup := newTestService(t)
	defer cleanup()
	db_rows := setupQueryTest(t, db)

	resp, err := svc.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{},
			Limit: 0,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, db_rows, []int{0, 1, 2, 3, 4}, resp.GetEnvelopes())
}

func TestQueryPagedEnvelopes(t *testing.T) {
	svc, db, cleanup := newTestService(t)
	defer cleanup()
	db_rows := setupQueryTest(t, db)

	resp, err := svc.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{},
			Limit: 2,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, db_rows, []int{0, 1}, resp.GetEnvelopes())
}

func TestQueryEnvelopesByOriginator(t *testing.T) {
	svc, db, cleanup := newTestService(t)
	defer cleanup()
	db_rows := setupQueryTest(t, db)

	resp, err := svc.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Filter: &message_api.EnvelopesQuery_OriginatorId{
					OriginatorId: 2,
				},
				LastSeen: nil,
			},
			Limit: 0,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, db_rows, []int{1, 3}, resp.GetEnvelopes())
}

func TestQueryEnvelopesByTopic(t *testing.T) {
	svc, db, cleanup := newTestService(t)
	defer cleanup()
	db_rows := setupQueryTest(t, db)

	resp, err := svc.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Filter:   &message_api.EnvelopesQuery_Topic{Topic: []byte("topicA")},
				LastSeen: nil,
			},
			Limit: 0,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, db_rows, []int{0, 1, 4}, resp.GetEnvelopes())
}

func TestQueryEnvelopesFromGatewaySID(t *testing.T) {
	svc, db, cleanup := newTestService(t)
	defer cleanup()
	db_rows := setupQueryTest(t, db)

	resp, err := svc.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Filter:   nil,
				LastSeen: &message_api.EnvelopesQuery_GatewaySid{GatewaySid: utils.SID(1, 2)},
			},
			Limit: 0,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, db_rows, []int{2, 3, 4}, resp.GetEnvelopes())
}

func TestQueryEnvelopesWithEmptyResult(t *testing.T) {
	svc, db, cleanup := newTestService(t)
	defer cleanup()
	db_rows := setupQueryTest(t, db)

	resp, err := svc.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Filter: &message_api.EnvelopesQuery_Topic{
					Topic: []byte("topicB"),
				},
				LastSeen: &message_api.EnvelopesQuery_GatewaySid{GatewaySid: utils.SID(1, 4)},
			},
			Limit: 0,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, db_rows, []int{}, resp.GetEnvelopes())
}

func TestQueryEnvelopesWithWrongGatewaySID(t *testing.T) {
	svc, db, cleanup := newTestService(t)
	defer cleanup()
	_ = setupQueryTest(t, db)

	_, err := svc.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Filter:   nil,
				LastSeen: &message_api.EnvelopesQuery_GatewaySid{GatewaySid: utils.SID(2, 2)},
			},
			Limit: 0,
		},
	)
	require.ErrorContains(t, err, "gateway")
}

func checkRowsMatchProtos(
	t *testing.T,
	allRows []queries.InsertGatewayEnvelopeParams,
	matchingIndices []int,
	protos []*message_api.GatewayEnvelope,
) {
	require.Len(t, protos, len(matchingIndices))
	for i, p := range protos {
		row := allRows[matchingIndices[i]]
		require.Equal(t, utils.SID(1, int64(matchingIndices[i]+1)), p.GetGatewaySid())
		require.Equal(t, row.OriginatorEnvelope, testutils.Marshal(t, p.GetOriginatorEnvelope()))
	}
}
