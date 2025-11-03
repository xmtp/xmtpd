package message_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"connectrpc.com/connect"
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
	db *sql.DB,
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
	testutils.InsertGatewayEnvelopes(t, db, []queries.InsertGatewayEnvelopeParams{{
		OriginatorNodeID:     int32(nodeID),
		OriginatorSequenceID: int64(sequenceID),
		OriginatorEnvelope:   envBytes,
		Topic:                topicBytes,
	}})

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
	var (
		suite           = apiTestUtils.NewTestAPIServer(t)
		db, _           = testutils.NewDB(t, t.Context())
		installationID1 = testutils.RandomGroupID()
		installationID2 = testutils.RandomGroupID()
		installationID3 = testutils.RandomGroupID()
		installationID4 = testutils.RandomGroupID()
	)

	// Installation ID 1 has three key packages
	topic1, _ := writeKeyPackage(t, db, installationID1[:])

	// This one is totally ignored
	_, _ = writeKeyPackage(t, db, installationID1[:])

	// This one is the newest
	_, seq1 := writeKeyPackage(t, db, installationID1[:])
	topic2, seq2 := writeKeyPackage(t, db, installationID2[:])
	topic3, seq3 := writeKeyPackage(t, db, installationID3[:])

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
			resp, err := suite.ClientReplication.GetNewestEnvelope(
				context.Background(),
				connect.NewRequest(&message_api.GetNewestEnvelopeRequest{
					Topics: requestedTopicsBytes,
				}),
			)
			fmt.Printf("### DEBUG: resp %+v\n", resp)
			fmt.Printf("### DEBUG: err %+v\n", err)
			require.NoError(t, err)
			require.Equal(t, c.expectedNumReturned, len(resp.Msg.Results))

			parsedResults := parseResults(t, resp.Msg.Results)
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
					require.EqualValues(t, seq, parsedResults[i].OriginatorSequenceID())
				}
			}
		})
	}
}
