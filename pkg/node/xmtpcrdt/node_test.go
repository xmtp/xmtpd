package xmtpcrdtnode_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/api/client"
	"github.com/xmtp/xmtpd/pkg/context"
	xmtpcrdtnode "github.com/xmtp/xmtpd/pkg/node/xmtpcrdt"
	"github.com/xmtp/xmtpd/pkg/otel"
	memstore "github.com/xmtp/xmtpd/pkg/store/mem"
	test "github.com/xmtp/xmtpd/pkg/testing"
	"github.com/xmtp/xmtpd/pkg/zap"
	"google.golang.org/protobuf/proto"
)

func TestNode_NewClose(t *testing.T) {
	t.Parallel()

	n := newTestNode(t)
	err := n.Close()
	require.NoError(t, err)
}

func TestNode_Publish(t *testing.T) {
	n := newTestNode(t)
	defer n.Close()
	ctx := test.NewContext(t)
	_, err := n.Publish(ctx, &messagev1.PublishRequest{})
	require.NoError(t, err)
}

func TestNode_Subscribe(t *testing.T) {
	n := newTestNode(t)
	defer n.Close()
	err := n.Subscribe(&messagev1.SubscribeRequest{}, nil)
	require.Equal(t, err, xmtpcrdtnode.ErrMissingTopic)
}

func TestNode_Query(t *testing.T) {
	n := newTestNode(t)
	defer n.Close()
	ctx := test.NewContext(t)
	_, err := n.Query(ctx, &messagev1.QueryRequest{})
	require.Equal(t, err, xmtpcrdtnode.ErrMissingTopic)
}

func TestNode_BatchQuery(t *testing.T) {
	n := newTestNode(t)
	defer n.Close()
	ctx := test.NewContext(t)
	_, err := n.BatchQuery(ctx, &messagev1.BatchQueryRequest{})
	require.NoError(t, err)
}

func TestNode_SubscribeAll(t *testing.T) {
	n := newTestNode(t)
	defer n.Close()
	ctrl := gomock.NewController(t)
	stream := xmtpcrdtnode.NewMockMessageApi_SubscribeServer(ctrl)
	stream.EXPECT().Send(&messagev1.Envelope{}).Return(nil)
	ctx := test.NewContext(t)
	ctx.Close()
	stream.EXPECT().Context().Return(ctx)
	err := n.SubscribeAll(&messagev1.SubscribeAllRequest{}, stream)
	require.NoError(t, err)
}

func TestNode_PublishSubscribeQuery_SingleNode(t *testing.T) {
	t.Parallel()

	n := newTestNode(t)
	defer n.Close()

	topic1Sub := n.subscribe(t, "topic1")
	topic1Envs := n.publishRandom(t, topic1Sub.topic, 2)
	topic1Sub.requireEventuallyCapturedEvents(t, topic1Envs)
	n.requireStoredEvents(t, "topic1", topic1Envs)

	topic2Sub := n.subscribe(t, "topic2")
	topic2Envs1 := n.publishRandom(t, topic2Sub.topic, 1)
	topic2Sub.requireEventuallyCapturedEvents(t, topic2Envs1)
	n.requireStoredEvents(t, "topic2", topic2Envs1)

	topic3Sub := n.subscribe(t, "topic3")
	topic3Sub.requireEventuallyCapturedEvents(t, nil)
	n.requireStoredEvents(t, "topic3", nil)

	topic4Sub := n.subscribe(t, "topic4")
	topic4Envs := n.publishRandom(t, topic4Sub.topic, 3)
	topic4Sub.requireEventuallyCapturedEvents(t, topic4Envs)
	n.requireStoredEvents(t, "topic4", topic4Envs)

	topic2Envs2 := n.publishRandom(t, topic2Sub.topic, 2)
	topic2Envs := append(topic2Envs1, topic2Envs2...)
	topic2Sub.requireEventuallyCapturedEvents(t, topic2Envs)
	n.requireStoredEvents(t, "topic2", topic2Envs)

	n.requireStoredEvents(t, "topic1", topic1Envs)
	n.requireStoredEvents(t, "topic2", topic2Envs)
	n.requireStoredEvents(t, "topic3", nil)
	n.requireStoredEvents(t, "topic4", topic4Envs)

}

func TestNode_PublishSubscribeQuery_TwoNodes(t *testing.T) {
	t.Parallel()

	n1 := newTestNodeWithOptions(t, "node1", nil)
	defer n1.Close()

	n2 := newTestNodeWithOptions(t, "node2", nil)
	defer n2.Close()

	n1.connect(t, n2)

	n1Topic1Sub := n1.subscribe(t, "topic1")
	n1Topic1Envs := n1.publishRandom(t, n1Topic1Sub.topic, 1)
	n1Topic1Sub.requireEventuallyCapturedEvents(t, n1Topic1Envs)
	n1.requireStoredEvents(t, "topic1", n1Topic1Envs)

	n2Topic1Sub := n2.subscribe(t, "topic1")
	n2Topic1Envs := n2.publishRandom(t, n2Topic1Sub.topic, 2)
	n2Topic1Sub.requireEventuallyCapturedEvents(t, n2Topic1Envs)
	n2.requireStoredEvents(t, "topic1", append(n1Topic1Envs, n2Topic1Envs...))
}

func TestNode_PersistentPeers(t *testing.T) {
	t.Parallel()

	n1 := newTestNodeWithOptions(t, "node1", nil)
	defer n1.Close()

	// TODO: test that connect/disconnect works as expected when not a persistent peer

	n2 := newTestNodeWithOptions(t, "node2", n1.P2PListenAddresses())
	defer n2.Close()

	// TODO: remove this
	time.Sleep(time.Second)

	n1Topic1Sub := n1.subscribe(t, "topic1")
	time.Sleep(time.Second)
	n2Topic1Sub := n2.subscribe(t, "topic1")
	topic1Envs := n1.publishRandom(t, n1Topic1Sub.topic, 1)
	n1Topic1Sub.requireEventuallyCapturedEvents(t, topic1Envs)
	n2Topic1Sub.requireEventuallyCapturedEvents(t, topic1Envs)

	n1Topic2Sub := n1.subscribe(t, "topic2")
	n2Topic2Sub := n2.subscribe(t, "topic2")
	topic2Envs := n2.publishRandom(t, n2Topic2Sub.topic, 1)
	n1Topic2Sub.requireEventuallyCapturedEvents(t, topic2Envs)
	n2Topic2Sub.requireEventuallyCapturedEvents(t, topic2Envs)

	// TODO: clean up
	fmt.Println("disconnect")
	err := n1.Disconnect(n1.ctx, n2.ID())
	require.NoError(t, err)
	err = n2.Disconnect(n2.ctx, n1.ID())
	require.NoError(t, err)
	time.Sleep(1 * time.Second)
	fmt.Println("disconnected-ish")

	// Should reconnect automatically.

	n1Topic1Sub = n1.subscribe(t, "topic1")
	n2Topic1Sub = n2.subscribe(t, "topic1")
	topic1Envs = n1.publishRandom(t, n1Topic1Sub.topic, 1)
	n1Topic1Sub.requireEventuallyCapturedEvents(t, topic1Envs)
	n2Topic1Sub.requireEventuallyCapturedEvents(t, topic1Envs)

	n1Topic2Sub = n1.subscribe(t, "topic2")
	n2Topic2Sub = n2.subscribe(t, "topic2")
	topic2Envs = n2.publishRandom(t, n2Topic2Sub.topic, 1)
	n1Topic2Sub.requireEventuallyCapturedEvents(t, topic2Envs)
	n2Topic2Sub.requireEventuallyCapturedEvents(t, topic2Envs)
}

type testNode struct {
	*xmtpcrdtnode.Node
	name string

	client client.Client
	ctx    context.Context
}

func newTestNode(t *testing.T) *testNode {
	return newTestNodeWithOptions(t, "", nil)
}

func newTestNodeWithOptions(t *testing.T, name string, persistentPeers []string) *testNode {
	t.Helper()
	ctx := test.NewContext(t)
	if name != "" {
		ctx = context.WithLogger(ctx, ctx.Logger().Named(name))
	}

	node, err := xmtpcrdtnode.New(ctx, memstore.NewNodeStore(ctx), &xmtpcrdtnode.Options{
		P2P: xmtpcrdtnode.P2POptions{
			PersistentPeers: persistentPeers,
		},
		OpenTelemetry: otel.Options{
			CollectorAddress: "localhost",
			CollectorPort:    4317,
		},
	})
	require.NoError(t, err)

	client := client.NewHTTPClient(ctx.Logger(), fmt.Sprintf("http://localhost:%d", node.APIHTTPListenPort()), "test", name)

	return &testNode{
		Node:   node,
		name:   name,
		client: client,
		ctx:    ctx,
	}
}

func (n *testNode) Close() error {
	n.ctx.Close()
	n.Node.Close()
	return nil
}

func (n *testNode) connect(t *testing.T, to *testNode) {
	t.Helper()

	err := n.Connect(n.ctx, to.Address())
	require.NoError(t, err)
	log := n.ctx.Logger()

	// Wait for peers to be connected and grafted to the pubsub topic.
	// See https://github.com/libp2p/go-libp2p-pubsub/issues/331
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	attempt := 1
	var connected bool
	ctx := context.WithTimeout(n.ctx, 5*time.Second)
	defer ctx.Close()
syncLoop:
	for {
		select {
		case <-ctx.Done():
			log.Info("context closed", zap.Error(ctx.Err()))
			break syncLoop
		case <-ticker.C:
			topic := "sync-" + test.RandomStringLower(13)
			sentEnv := newRandomEnvelope(topic)
			_, err := n.client.Publish(n.ctx, &messagev1.PublishRequest{
				Envelopes: []*messagev1.Envelope{sentEnv},
			})
			require.NoError(t, err)

			func() {
				queryTicker := time.NewTicker(100 * time.Millisecond)
				defer queryTicker.Stop()
				queryCtx := context.WithTimeout(ctx, 500*time.Millisecond)
				defer queryCtx.Close()
				for {
					select {
					case <-queryCtx.Done():
						return
					case <-queryTicker.C:
						res, err := to.client.Query(n.ctx, &messagev1.QueryRequest{
							ContentTopics: []string{topic},
							PagingInfo: &messagev1.PagingInfo{
								Direction: messagev1.SortDirection_SORT_DIRECTION_DESCENDING,
								Limit:     1,
							},
						})
						require.NoError(t, err)

						if len(res.Envelopes) > 0 && proto.Equal(sentEnv, res.Envelopes[0]) {
							connected = true
							return
						}
					}
				}
			}()
			if connected {
				break syncLoop
			}

			log.Debug("waiting for p2p connectivity sync message", zap.Int("attempt", attempt))
			attempt++
		}
	}
	require.True(t, connected, fmt.Sprintf("node %s failed to connect to node %s", n.name, to.name))
}

func (n *testNode) publishRandom(t *testing.T, topic string, count int) []*messagev1.Envelope {
	t.Helper()
	envs := make([]*messagev1.Envelope, count)
	for i := 0; i < count; i++ {
		env := newRandomEnvelope(topic)
		res, err := n.client.Publish(n.ctx, &messagev1.PublishRequest{
			Envelopes: []*messagev1.Envelope{env},
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		envs[i] = env
	}
	return envs
}

func (n *testNode) subscribe(t *testing.T, topic string) *testSubscriber {
	sub := &testSubscriber{
		topic: topic,
	}
	stream, err := n.client.Subscribe(n.ctx, &messagev1.SubscribeRequest{
		ContentTopics: []string{topic},
	})
	require.NoError(t, err)
	sub.stream = stream
	go func() {
		for {
			env, err := sub.stream.Next(n.ctx)
			if err == context.Canceled {
				return
			}
			require.NoError(t, err)
			func() {
				sub.Lock()
				defer sub.Unlock()
				sub.envs = append(sub.envs, env)
			}()
		}
	}()
	return sub
}

func (n *testNode) requireStoredEvents(t *testing.T, topic string, expected []*messagev1.Envelope) {
	t.Helper()
	res, err := n.client.Query(n.ctx, &messagev1.QueryRequest{
		ContentTopics: []string{topic},
	})
	require.NoError(t, err)
	require.Len(t, res.Envelopes, len(expected))
	requireEnvelopesEqual(t, expected, res.Envelopes)
}

type testSubscriber struct {
	topic  string
	stream client.Stream
	envs   []*messagev1.Envelope
	sync.RWMutex
}

func (s *testSubscriber) requireEventuallyCapturedEvents(t *testing.T, expected []*messagev1.Envelope) {
	t.Helper()
	assert.Eventually(t, func() bool {
		s.RLock()
		defer s.RUnlock()
		return len(s.envs) == len(expected)
	}, 3*time.Second, 10*time.Millisecond)
	test.RequireProtoEqual(t, expected, s.envs)
}

func requireEnvelopesEqual(t *testing.T, actual, expected []*messagev1.Envelope) {
	t.Helper()
	expected = expected[:]
	sort.Slice(expected, func(i, j int) bool {
		d := int(expected[i].TimestampNs) - int(expected[j].TimestampNs)
		if d != 0 {
			return d < 0
		}
		return bytes.Compare(expected[i].Message, expected[j].Message) < 0
	})
	actual = actual[:]
	sort.Slice(actual, func(i, j int) bool {
		d := int(actual[i].TimestampNs) - int(actual[j].TimestampNs)
		if d != 0 {
			return d < 0
		}
		return bytes.Compare(actual[i].Message, actual[j].Message) < 0
	})
	test.RequireProtoEqual(t, expected, actual)
}

func newRandomEnvelope(topic string) *messagev1.Envelope {
	return &messagev1.Envelope{
		ContentTopic: topic,
		TimestampNs:  uint64(rand.Intn(100)),
		Message:      []byte("msg-" + test.RandomString(13)),
	}
}
