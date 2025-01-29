package metadata_test

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
	"testing"
	"time"

	"github.com/xmtp/xmtpd/pkg/api/message"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
	testUtilsApi "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
)

var (
	topicA = topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, []byte("topicA")).Bytes()
	topicB = topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, []byte("topicB")).Bytes()
)
var allRows []queries.InsertGatewayEnvelopeParams

func setupTest(
	t *testing.T,
) (metadata_api.MetadataApiClient, *sql.DB, testUtilsApi.ApiServerMocks, func()) {
	allRows = []queries.InsertGatewayEnvelopeParams{
		// Initial rows
		{
			OriginatorNodeID:     1,
			OriginatorSequenceID: 1,
			Topic:                topicA,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 1, 1, topicA),
			),
		},
		{
			OriginatorNodeID:     2,
			OriginatorSequenceID: 1,
			Topic:                topicA,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 2, 1, topicA),
			),
		},
		// Later rows
		{
			OriginatorNodeID:     1,
			OriginatorSequenceID: 2,
			Topic:                topicB,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 1, 2, topicB),
			),
		},
		{
			OriginatorNodeID:     2,
			OriginatorSequenceID: 2,
			Topic:                topicB,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 2, 2, topicB),
			),
		},
		{
			OriginatorNodeID:     1,
			OriginatorSequenceID: 3,
			Topic:                topicA,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 1, 3, topicA),
			),
		},
	}

	return testUtilsApi.NewTestMetadataAPIClient(t)
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
	client, db, _, cleanup := setupTest(t)
	defer cleanup()
	insertInitialRows(t, db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cursor, err := client.GetSyncCursor(ctx, &metadata_api.GetSyncCursorRequest{})

	require.NoError(t, err)
	require.NotNil(t, cursor)

	expectedCursor := map[uint32]uint64{
		1: 1,
		2: 1,
	}

	require.Equal(t, expectedCursor, cursor.LatestSync.NodeIdToSequenceId)

	insertAdditionalRows(t, db)
	require.Eventually(t, func() bool {
		expectedCursor := map[uint32]uint64{
			1: 3,
			2: 2,
		}

		cursor, err := client.GetSyncCursor(ctx, &metadata_api.GetSyncCursorRequest{})
		if err != nil {
			t.Logf("Error fetching sync cursor: %v", err)
			return false
		}
		if cursor == nil {
			t.Log("Cursor is nil")
			return false
		}

		return assert.ObjectsAreEqual(expectedCursor, cursor.LatestSync.NodeIdToSequenceId)
	}, 500*time.Millisecond, 50*time.Millisecond)
}

func TestSubscribeSyncCursorBasic(t *testing.T) {
	client, db, _, cleanup := setupTest(t)
	defer cleanup()
	insertInitialRows(t, db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.SubscribeSyncCursor(ctx, &metadata_api.GetSyncCursorRequest{})
	require.NoError(t, err)
	require.NotNil(t, stream)

	firstUpdate, err := stream.Recv()
	require.NoError(t, err)
	require.NotNil(t, firstUpdate)

	expectedCursor := map[uint32]uint64{
		1: 1,
		2: 1,
	}

	require.Equal(t, expectedCursor, firstUpdate.LatestSync.NodeIdToSequenceId)

	insertAdditionalRows(t, db)

	require.Eventually(t, func() bool {
		expectedCursor := map[uint32]uint64{
			1: 3,
			2: 2,
		}

		update, err := stream.Recv()
		if err != nil {
			t.Logf("Error receiving sync cursor update: %v", err)
			return false
		}
		if update == nil {
			t.Log("Received nil update from stream")
			return false
		}

		return assert.ObjectsAreEqual(expectedCursor, update.LatestSync.NodeIdToSequenceId)
	}, 500*time.Millisecond, 50*time.Millisecond)
}
