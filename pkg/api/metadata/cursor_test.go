package metadata_test

import (
	"database/sql"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
	metadata_apiconnect "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api/metadata_apiconnect"

	"github.com/xmtp/xmtpd/pkg/api/message"
	dbUtils "github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
	testUtilsApi "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
)

var (
	topicA  = topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("topicA")).Bytes()
	topicB  = topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("topicB")).Bytes()
	allRows []queries.InsertGatewayEnvelopeParams
)

func setupTest(
	t *testing.T,
) (metadata_apiconnect.MetadataApiClient, *sql.DB, testUtilsApi.APIServerMocks) {
	var (
		suite   = testUtilsApi.NewTestAPIServer(t)
		payerID = dbUtils.NullInt32(testutils.CreatePayer(t, suite.DB))
	)

	allRows = []queries.InsertGatewayEnvelopeParams{
		// Initial rows
		{
			OriginatorNodeID:     100,
			OriginatorSequenceID: 1,
			Topic:                topicA,
			PayerID:              payerID,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 1, 1, topicA),
			),
		},
		{
			OriginatorNodeID:     200,
			OriginatorSequenceID: 1,
			Topic:                topicA,
			PayerID:              payerID,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 2, 1, topicA),
			),
		},
		// Later rows
		{
			OriginatorNodeID:     100,
			OriginatorSequenceID: 2,
			Topic:                topicB,
			PayerID:              payerID,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 1, 2, topicB),
			),
		},
		{
			OriginatorNodeID:     200,
			OriginatorSequenceID: 2,
			Topic:                topicB,
			PayerID:              payerID,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 2, 2, topicB),
			),
		},
		{
			OriginatorNodeID:     100,
			OriginatorSequenceID: 3,
			Topic:                topicA,
			PayerID:              payerID,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 1, 3, topicA),
			),
		},
	}

	return suite.ClientMetadata, suite.DB, suite.APIServerMocks
}

func insertInitialRows(t *testing.T, store *sql.DB) {
	testutils.InsertGatewayEnvelopes(t, store, []queries.InsertGatewayEnvelopeParams{
		allRows[0], allRows[1],
	})
	time.Sleep(message.SubscribeWorkerPollTime + 100*time.Millisecond)
}

func insertAdditionalRows(t *testing.T, store *sql.DB, notifyChan ...chan bool) {
	testutils.InsertGatewayEnvelopes(t, store, []queries.InsertGatewayEnvelopeParams{
		allRows[2], allRows[3], allRows[4],
	}, notifyChan...)
}

func TestGetCursorBasic(t *testing.T) {
	client, db, _ := setupTest(t)
	insertInitialRows(t, db)

	ctx := t.Context()

	cursor, err := client.GetSyncCursor(
		ctx,
		connect.NewRequest(&metadata_api.GetSyncCursorRequest{}),
	)

	require.NoError(t, err)
	require.NotNil(t, cursor)

	expectedCursor := map[uint32]uint64{
		100: 1,
		200: 1,
	}

	require.Equal(t, expectedCursor, cursor.Msg.LatestSync.NodeIdToSequenceId)

	insertAdditionalRows(t, db)
	require.Eventually(t, func() bool {
		expectedCursor := map[uint32]uint64{
			100: 3,
			200: 2,
		}

		cursor, err := client.GetSyncCursor(
			ctx,
			connect.NewRequest(&metadata_api.GetSyncCursorRequest{}),
		)
		if err != nil {
			t.Logf("Error fetching sync cursor: %v", err)
			return false
		}

		if cursor == nil {
			t.Log("Cursor is nil")
			return false
		}

		return assert.ObjectsAreEqual(expectedCursor, cursor.Msg.LatestSync.NodeIdToSequenceId)
	}, 500*time.Millisecond, 50*time.Millisecond)
}

func TestSubscribeSyncCursorBasic(t *testing.T) {
	client, db, _ := setupTest(t)
	insertInitialRows(t, db)

	ctx := t.Context()

	stream, err := client.SubscribeSyncCursor(
		ctx,
		connect.NewRequest(&metadata_api.GetSyncCursorRequest{}),
	)
	require.NoError(t, err)
	require.NotNil(t, stream)

	// Advance to the first update, otherwise stream.Msg() will be nil.
	shouldReceive := stream.Receive()
	require.True(t, shouldReceive)

	firstUpdate := stream.Msg()
	require.NotNil(t, firstUpdate)

	expectedCursor := map[uint32]uint64{
		100: 1,
		200: 1,
	}

	require.Equal(t, expectedCursor, firstUpdate.LatestSync.NodeIdToSequenceId)

	insertAdditionalRows(t, db)

	require.Eventually(t, func() bool {
		expectedCursor = map[uint32]uint64{
			100: 3,
			200: 2,
		}

		return assert.ObjectsAreEqual(expectedCursor, stream.Msg().LatestSync.NodeIdToSequenceId)
	}, 500*time.Millisecond, 50*time.Millisecond)

	shouldReceive = stream.Receive()
	require.True(t, shouldReceive)

	secondUpdate := stream.Msg()
	require.NotNil(t, secondUpdate)

	assert.ObjectsAreEqual(expectedCursor, secondUpdate.LatestSync.NodeIdToSequenceId)
}
