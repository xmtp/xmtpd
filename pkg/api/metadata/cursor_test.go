package metadata_test

import (
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
	metadata_apiconnect "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api/metadata_apiconnect"

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
	allRows []queries.InsertGatewayEnvelopeV3Params
)

func setupTest(
	t *testing.T,
) (metadata_apiconnect.MetadataApiClient, *testUtilsApi.APIServerTestSuite) {
	var (
		suite   = testUtilsApi.NewTestAPIServer(t)
		payerID = dbUtils.NullInt32(testutils.CreatePayer(t, suite.DB))
	)

	allRows = []queries.InsertGatewayEnvelopeV3Params{
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

	return suite.ClientMetadata, suite
}

// insertInitialRows inserts the first two rows and blocks until GetSyncCursor
// reports the subscribe worker has polled past them, so tests observe a known
// cursor state.
func insertInitialRows(
	t *testing.T,
	client metadata_apiconnect.MetadataApiClient,
	suite *testUtilsApi.APIServerTestSuite,
) {
	testutils.InsertGatewayEnvelopes(t, suite.DB, []queries.InsertGatewayEnvelopeV3Params{
		allRows[0], allRows[1],
	})
	expected := map[uint32]uint64{100: 1, 200: 1}
	require.Eventually(t, func() bool {
		resp, err := client.GetSyncCursor(
			t.Context(),
			connect.NewRequest(&metadata_api.GetSyncCursorRequest{}),
		)
		if err != nil {
			return false
		}
		return assert.ObjectsAreEqual(
			expected,
			resp.Msg.GetLatestSync().GetNodeIdToSequenceId(),
		)
	}, 5*time.Second, 5*time.Millisecond)
}

func insertAdditionalRows(
	t *testing.T,
	suite *testUtilsApi.APIServerTestSuite,
	notifyChan ...chan bool,
) {
	testutils.InsertGatewayEnvelopes(t, suite.DB, []queries.InsertGatewayEnvelopeV3Params{
		allRows[2], allRows[3], allRows[4],
	}, notifyChan...)
}

func TestGetCursorBasic(t *testing.T) {
	client, suite := setupTest(t)
	insertInitialRows(t, client, suite)

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

	require.Equal(t, expectedCursor, cursor.Msg.GetLatestSync().GetNodeIdToSequenceId())

	insertAdditionalRows(t, suite)
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

		return assert.ObjectsAreEqual(
			expectedCursor,
			cursor.Msg.GetLatestSync().GetNodeIdToSequenceId(),
		)
	}, 500*time.Millisecond, 50*time.Millisecond)
}

func TestSubscribeSyncCursorBasic(t *testing.T) {
	client, suite := setupTest(t)
	insertInitialRows(t, client, suite)

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

	require.Equal(t, expectedCursor, firstUpdate.GetLatestSync().GetNodeIdToSequenceId())

	insertAdditionalRows(t, suite)

	expectedCursor = map[uint32]uint64{
		100: 3,
		200: 2,
	}

	require.Eventually(t, func() bool {
		if stream.Receive() {
			cursor := stream.Msg()
			require.NotNil(t, cursor)
			return assert.ObjectsAreEqual(
				expectedCursor,
				cursor.GetLatestSync().GetNodeIdToSequenceId(),
			)
		}
		return false
	}, 2000*time.Millisecond, 10*time.Millisecond)
}
