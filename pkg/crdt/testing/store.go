package crdttest

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/require"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	v1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/api"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

type ITestStore interface {
	crdt.Store

	Events(context.Context) ([]*types.Event, error)
	Heads(context.Context) ([]multihash.Multihash, error)
	InsertNewEvents(context.Context, []*types.Event) error

	Close() error
}

type TestStoreMaker func(t *testing.T) *TestStore

type TestStore struct {
	ITestStore
	ctx context.Context
}

func NewTestStore(ctx context.Context, store ITestStore) *TestStore {
	return &TestStore{
		ITestStore: store,
		ctx:        ctx,
	}
}

func (s *TestStore) addRandomHead(t *testing.T, topic string) *types.Event {
	t.Helper()
	ev := s.newRandomEvent(t, topic)
	s.addHead(t, ev)
	return ev
}

func (s *TestStore) addHead(t *testing.T, head *types.Event) {
	t.Helper()
	added, err := s.InsertHead(s.ctx, head)
	require.NoError(t, err)
	require.True(t, added)
}

func (s *TestStore) addRandomEvent(t *testing.T, topic string) *types.Event {
	t.Helper()
	ev := s.newRandomEvent(t, topic)
	s.addEvent(t, ev)
	return ev
}

func (s *TestStore) addEvent(t *testing.T, ev *types.Event) {
	t.Helper()
	added, err := s.InsertEvent(s.ctx, ev)
	require.NoError(t, err)
	require.True(t, added)
}

func (s *TestStore) addExistingEvent(t *testing.T, ev *types.Event) {
	t.Helper()
	added, err := s.InsertEvent(s.ctx, ev)
	require.NoError(t, err)
	require.False(t, added)
}

func (s *TestStore) appendRandomEvent(t *testing.T, topic string) *types.Event {
	t.Helper()
	return s.appendEvent(t, newRandomEnvelope(t, topic))
}

func (s *TestStore) appendEvent(t *testing.T, env *messagev1.Envelope) *types.Event {
	t.Helper()
	ev, err := s.AppendEvent(s.ctx, env)
	require.NoError(t, err)
	return ev
}

func (s *TestStore) newRandomEvent(t *testing.T, topic string) *types.Event {
	t.Helper()
	return s.newRandomEventWithHeads(t, topic, nil)
}

func (s *TestStore) newRandomEventWithHeads(t *testing.T, topic string, heads []multihash.Multihash) *types.Event {
	t.Helper()
	ev, err := types.NewEvent(newRandomEnvelope(t, topic), heads)
	require.NoError(t, err)
	return ev
}

func (s *TestStore) requireEventsEqual(t *testing.T, expected ...*types.Event) {
	t.Helper()
	events, err := s.Events(s.ctx)
	require.NoError(t, err)
	var exp, act []string
	for _, e := range events {
		act = append(act, DumpEvent(e))
	}
	for _, e := range expected {
		exp = append(exp, DumpEvent(e))
	}
	require.ElementsMatch(t, exp, act)
}

func (s *TestStore) requireNoEvents(t *testing.T) {
	t.Helper()
	events, err := s.Events(s.ctx)
	require.NoError(t, err)
	require.Empty(t, events)
}

func (s *TestStore) Seed(t testing.TB, topic string, count int) []*types.Event {
	t.Helper()
	if count <= 0 {
		return nil
	}
	ctx := test.NewContext(t)
	events := make([]*types.Event, 0, count)
	var prev []multihash.Multihash
	for i := 0; i < count; i++ {
		env := &messagev1.Envelope{
			ContentTopic: string(topic),
			TimestampNs:  uint64(i + 1),
			Message:      []byte(fmt.Sprintf("msg-%d", i+1)),
		}
		ev, err := types.NewEvent(env, prev)
		require.NoError(t, err)
		prev = []multihash.Multihash{ev.Cid}
		events = append(events, ev)
	}
	if len(events) > 1 {
		require.NoError(t, s.InsertNewEvents(ctx, events[:len(events)-1]))
	}
	_, err := s.InsertHead(ctx, events[len(events)-1])
	require.NoError(t, err)
	return events
}

func (s *TestStore) query(t *testing.T, topic string, modifiers ...api.QueryModifier) *messagev1.QueryResponse {
	t.Helper()
	ctx := test.NewContext(t)
	req := api.NewQuery(topic, modifiers...)
	res, err := s.Query(ctx, req)
	require.NoError(t, err)
	return res
}

func newRandomEnvelope(t *testing.T, topic string) *messagev1.Envelope {
	return &messagev1.Envelope{
		ContentTopic: topic,
		TimestampNs:  uint64(rand.Intn(100)),
		Message:      []byte("msg-" + test.RandomString(13)),
	}
}

func newRandomEnvelopeWithRandomTopic(t *testing.T) *messagev1.Envelope {
	return newRandomEnvelope(t, "topic-"+test.RandomStringLower(13))
}

func requireResultEqual(t *testing.T, res *v1.QueryResponse, from, to int) {
	t.Helper()
	var actual []int
	for _, env := range res.Envelopes {
		actual = append(actual, int(env.TimestampNs))
	}
	require.Equal(t, intRange(from, to), actual)
}

func requireResultCursor(t *testing.T, res *v1.QueryResponse, expected int) {
	t.Helper()
	require.NotNil(t, res.PagingInfo, "paging info")
	require.NotNil(t, res.PagingInfo.Cursor, "cursor")
	cursor := res.PagingInfo.Cursor.GetIndex()
	require.NotNil(t, cursor, "index cursor")
	actual := int(cursor.SenderTimeNs)
	require.Equal(t, expected, actual, "timestamp")
}

func requireNoCursor(t *testing.T, res *v1.QueryResponse) {
	t.Helper()
	require.Nil(t, res.PagingInfo, "paging info")
}

// generate a list of ints in from start to end inclusive.
// if end < start generate it in reverse.
func intRange(start, end int) (list []int) {
	if start < end {
		list = make([]int, end-start+1)
		for k := range list {
			list[k] = start + k
		}
		return list
	}
	list = make([]int, start-end+1)
	for k := range list {
		list[k] = start - k
	}
	return list
}
