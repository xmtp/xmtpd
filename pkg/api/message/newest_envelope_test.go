package message_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/testutils"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopesUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
	"google.golang.org/protobuf/proto"
)

func writeKeyPackage(
	t *testing.T,
	querier *queries.Queries,
	installationKey []byte,
) (topic.Topic, uint64) {
	topicObj := topic.NewTopic(topic.TopicKindKeyPackagesV1, installationKey)
	topicBytes := topicObj.Bytes()
	nodeID := uint32(100)
	sequenceID := uint64(testutils.RandomInt32())
	env := envelopesUtils.CreateOriginatorEnvelopeWithTopic(
		t,
		nodeID,
		sequenceID,
		topicBytes,
	)
	envBytes, err := proto.Marshal(env)
	require.NoError(t, err)

	_, err = querier.InsertGatewayEnvelopeV2(t.Context(), queries.InsertGatewayEnvelopeV2Params{
		OriginatorNodeID:     int32(nodeID),
		OriginatorSequenceID: int64(sequenceID),
		OriginatorEnvelope:   envBytes,
		Topic:                topicBytes,
	})

	require.NoError(t, err)

	return *topicObj, sequenceID
}

func parseResults(
	t *testing.T,
	envs []*message_api.GetNewestEnvelopeResponse_Response,
) []*envelopes.OriginatorEnvelope {
	parsedEnvelopes := make([]*envelopes.OriginatorEnvelope, len(envs))
	var err error
	for i, env := range envs {
		if env.OriginatorEnvelope == nil {
			continue
		}
		parsedEnvelopes[i], err = envelopes.NewOriginatorEnvelope(env.OriginatorEnvelope)
		require.NoError(t, err)
	}
	return parsedEnvelopes
}

func TestGetNewestEnvelope(t *testing.T) {
	api, db, _ := apiTestUtils.NewTestReplicationAPIClient(t)
	querier := queries.New(db)

	installationID1 := testutils.RandomGroupID()
	installationID2 := testutils.RandomGroupID()
	installationID3 := testutils.RandomGroupID()
	installationID4 := testutils.RandomGroupID()

	// Installation ID 1 has three key packages
	topic1, discared := writeKeyPackage(t, querier, installationID1[:])
	t.Log(topic1, discared)
	// This one is totally ignored
	_, _ = writeKeyPackage(t, querier, installationID1[:])
	// This one is the newest
	_, seq1 := writeKeyPackage(t, querier, installationID1[:])
	topic2, seq2 := writeKeyPackage(t, querier, installationID2[:])
	topic3, seq3 := writeKeyPackage(t, querier, installationID3[:])

	t.Log(topic1, seq1)
	t.Log(topic2, seq2)
	t.Log(topic3, seq3)

	// A topic that doesn't have anything in the DB
	topic4 := *topic.NewTopic(topic.TopicKindKeyPackagesV1, installationID4[:])

	cases := []struct {
		name                string
		requestedTopics     []topic.Topic
		expectedNumReturned int
		expectedTopics      []*topic.Topic
		expectedSequenceIDs []uint64
	}{
		{
			name: "all three installation IDs",
			requestedTopics: []topic.Topic{
				topic1, topic2, topic3,
			},
			expectedNumReturned: 3,
			expectedTopics:      []*topic.Topic{&topic1, &topic2, &topic3},
			expectedSequenceIDs: []uint64{seq1, seq2, seq3},
		},
		{
			name:                "only installation ID 1",
			requestedTopics:     []topic.Topic{topic1},
			expectedNumReturned: 1,
			expectedTopics:      []*topic.Topic{&topic1},
			expectedSequenceIDs: []uint64{seq1},
		},
		{
			name:                "out of order installation IDs",
			requestedTopics:     []topic.Topic{topic2, topic1},
			expectedNumReturned: 2,
			expectedTopics:      []*topic.Topic{&topic2, &topic1},
			expectedSequenceIDs: []uint64{seq2, seq1},
		},
		{
			name:                "no envelopes for a topic",
			requestedTopics:     []topic.Topic{topic1, topic2, topic3, topic4},
			expectedNumReturned: 4,
			expectedTopics:      []*topic.Topic{&topic1, &topic2, &topic3, nil},
			expectedSequenceIDs: []uint64{seq1, seq2, seq3, 0},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			requestedTopicsBytes := make([][]byte, len(c.requestedTopics))
			for i, topic := range c.requestedTopics {
				requestedTopicsBytes[i] = topic.Bytes()
			}
			resp, err := api.GetNewestEnvelope(
				context.Background(),
				&message_api.GetNewestEnvelopeRequest{
					Topics: requestedTopicsBytes,
				},
			)
			require.NoError(t, err)
			require.Equal(t, c.expectedNumReturned, len(resp.Results))
			parsedResults := parseResults(t, resp.Results)
			for i, topic := range c.expectedTopics {
				if topic == nil {
					require.Nil(t, parsedResults[i])
				} else {
					require.Equal(t, *topic, parsedResults[i].TargetTopic())
				}
			}
			for i, seq := range c.expectedSequenceIDs {
				if seq == 0 {
					require.Nil(t, parsedResults[i])
				} else {
					t.Log(seq, parsedResults[i].OriginatorSequenceID())
					require.EqualValues(t, seq, parsedResults[i].OriginatorSequenceID())
				}
			}
		})
	}
}
