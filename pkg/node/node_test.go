package node_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/node"
	ntest "github.com/xmtp/xmtpd/pkg/node/testing"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func TestNode_NewClose(t *testing.T) {
	t.Parallel()

	n := ntest.NewTestNode(t)
	err := n.Close()
	require.NoError(t, err)
}

func TestNode_Publish(t *testing.T) {
	n := ntest.NewTestNode(t)
	defer n.Close()
	ctx := test.NewContext(t)
	_, err := n.Publish(ctx, &messagev1.PublishRequest{})
	require.NoError(t, err)
}

func TestNode_Subscribe_MissingTopic(t *testing.T) {
	n := ntest.NewTestNode(t)
	defer n.Close()
	err := n.Node.Subscribe(&messagev1.SubscribeRequest{}, nil)
	require.Equal(t, err, node.ErrMissingTopic)
}

func TestNode_Query(t *testing.T) {
	n := ntest.NewTestNode(t)
	defer n.Close()
	ctx := test.NewContext(t)
	_, err := n.Query(ctx, &messagev1.QueryRequest{})
	require.Equal(t, err, node.ErrMissingTopic)
}

func TestNode_Query_UnknownTopic(t *testing.T) {
	n := ntest.NewTestNode(t)
	defer n.Close()
	ctx := test.NewContext(t)
	res, err := n.Query(ctx, &messagev1.QueryRequest{
		ContentTopics: []string{"unknown-topic"},
	})
	require.NoError(t, err)
	require.Equal(t, &messagev1.QueryResponse{}, res)
}

func TestNode_BatchQuery(t *testing.T) {
	n := ntest.NewTestNode(t)
	defer n.Close()
	ctx := test.NewContext(t)
	_, err := n.BatchQuery(ctx, &messagev1.BatchQueryRequest{})
	require.NoError(t, err)
}

func TestNode_SubscribeAll(t *testing.T) {
	n := ntest.NewTestNode(t)
	defer n.Close()
	ctrl := gomock.NewController(t)
	stream := node.NewMockMessageApi_SubscribeServer(ctrl)
	stream.EXPECT().Send(&messagev1.Envelope{}).Return(nil)
	ctx := test.NewContext(t)
	ctx.Close()
	stream.EXPECT().Context().Return(ctx)
	err := n.SubscribeAll(&messagev1.SubscribeAllRequest{}, stream)
	require.NoError(t, err)
}

func TestNode_PublishSubscribeQuery_SingleNode(t *testing.T) {
	t.Parallel()

	n := ntest.NewTestNode(t)
	defer n.Close()

	topic1Sub := n.Subscribe(t, "topic1")
	topic1Envs := n.PublishRandom(t, topic1Sub.Topic, 2)
	topic1Sub.RequireEventuallyCapturedEvents(t, topic1Envs)
	n.RequireStoredEvents(t, "topic1", topic1Envs)

	topic2Sub := n.Subscribe(t, "topic2")
	topic2Envs1 := n.PublishRandom(t, topic2Sub.Topic, 1)
	topic2Sub.RequireEventuallyCapturedEvents(t, topic2Envs1)
	n.RequireStoredEvents(t, "topic2", topic2Envs1)

	topic3Sub := n.Subscribe(t, "topic3")
	topic3Sub.RequireEventuallyCapturedEvents(t, nil)

	topic4Sub := n.Subscribe(t, "topic4")
	topic4Envs := n.PublishRandom(t, topic4Sub.Topic, 3)
	topic4Sub.RequireEventuallyCapturedEvents(t, topic4Envs)
	n.RequireStoredEvents(t, "topic4", topic4Envs)

	topic2Envs2 := n.PublishRandom(t, topic2Sub.Topic, 2)
	topic2Envs := append(topic2Envs1, topic2Envs2...)
	topic2Sub.RequireEventuallyCapturedEvents(t, topic2Envs)
	n.RequireStoredEvents(t, "topic2", topic2Envs)

	n.RequireStoredEvents(t, "topic1", topic1Envs)
	n.RequireStoredEvents(t, "topic2", topic2Envs)
	n.RequireStoredEvents(t, "topic4", topic4Envs)

}

func TestNode_PublishSubscribeQuery_TwoNodes(t *testing.T) {
	t.Parallel()

	n1 := ntest.NewTestNodeWithName(t, "node1")
	defer n1.Close()

	n2 := ntest.NewTestNodeWithName(t, "node2")
	defer n2.Close()

	n1.Connect(t, n2)

	n1Topic1Sub := n1.Subscribe(t, "topic1")
	n1Topic1Envs := n1.PublishRandom(t, n1Topic1Sub.Topic, 1)
	n1Topic1Sub.RequireEventuallyCapturedEvents(t, n1Topic1Envs)
	n1.RequireStoredEvents(t, "topic1", n1Topic1Envs)

	n2Topic1Sub := n2.Subscribe(t, "topic1")
	n2Topic1Envs := n2.PublishRandom(t, n2Topic1Sub.Topic, 2)
	n2Topic1Sub.RequireEventuallyCapturedEvents(t, n2Topic1Envs)
	n2.RequireStoredEvents(t, "topic1", append(n1Topic1Envs, n2Topic1Envs...))
}

func TestNode_Fetch(t *testing.T) {
	t.Parallel()
	topic := "topic"

	n1 := ntest.NewTestNodeWithName(t, "node1")
	defer n1.Close()

	envs := n1.PublishRandom(t, topic, 3)
	n1.RequireStoredEvents(t, topic, envs)

	n2 := ntest.NewTestNodeWithName(t, "node2")
	defer n2.Close()
	n1.Connect(t, n2)

	envs = append(envs, n1.PublishRandom(t, topic, 3)...)
	n1.RequireStoredEvents(t, topic, envs)
	n2.RequireStoredEvents(t, topic, envs)
}
