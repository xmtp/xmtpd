package crdttest

import (
	"testing"

	"github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

type ITestStore interface {
	crdt.Store

	Events() ([]*types.Event, error)
	Heads() []multihash.Multihash
}

type TestStoreMaker func(t *testing.T) *TestStore

type TestStore struct {
	ITestStore
}

func NewTestStore(store ITestStore) *TestStore {
	return &TestStore{store}
}

func RunStoreTests(t *testing.T, storeMaker TestStoreMaker) {
	t.Run("add events", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)

		ev1 := s.addRandomEvent(t)
		s.requireEventsEqual(t, []*types.Event{ev1})

		ev2 := s.addRandomEvent(t)
		s.requireEventsEqual(t, []*types.Event{ev1, ev2})
	})

	t.Run("add existing event", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)

		ev1 := s.addRandomEvent(t)
		s.requireEventsEqual(t, []*types.Event{ev1})

		s.addExistingEvent(t, ev1)
		s.requireEventsEqual(t, []*types.Event{ev1})
	})

	t.Run("append events", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)

		ev1 := s.appendRandomEvent(t)
		require.NotNil(t, ev1.Cid)
		require.Nil(t, ev1.Links)
		s.requireEventsEqual(t, []*types.Event{ev1})

		ev2 := s.appendRandomEvent(t)
		require.NotNil(t, ev2.Cid)
		require.Equal(t, []multihash.Multihash{ev1.Cid}, ev2.Links)
		s.requireEventsEqual(t, []*types.Event{ev1, ev2})
	})

	t.Run("add remove heads", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)

		head := s.addRandomHead(t)
		s.requireEventsEqual(t, []*types.Event{head})

		ev1 := s.newRandomEventWithHeads(t, []multihash.Multihash{head.Cid})
		require.Equal(t, []multihash.Multihash{head.Cid}, ev1.Links)
		s.addEvent(t, ev1)
		s.requireEventsEqual(t, []*types.Event{head, ev1})

		ev2 := s.newRandomEvent(t)
		require.Nil(t, ev2.Links)
		s.addEvent(t, ev2)
		s.requireEventsEqual(t, []*types.Event{head, ev1, ev2})
	})

	t.Run("add existing head", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)

		head := s.addRandomHead(t)
		s.requireEventsEqual(t, []*types.Event{head})

		s.addExistingEvent(t, head)
		s.requireEventsEqual(t, []*types.Event{head})
	})
}

func (s *TestStore) addRandomHead(t *testing.T) *types.Event {
	t.Helper()
	ev := s.newRandomEvent(t)
	s.addHead(t, ev)
	return ev
}

func (s *TestStore) addHead(t *testing.T, head *types.Event) {
	t.Helper()
	added, err := s.AddHead(head)
	require.NoError(t, err)
	require.True(t, added)
}

func (s *TestStore) addRandomEvent(t *testing.T) *types.Event {
	t.Helper()
	ev := s.newRandomEvent(t)
	s.addEvent(t, ev)
	return ev
}

func (s *TestStore) addEvent(t *testing.T, ev *types.Event) {
	t.Helper()
	added, err := s.AddEvent(ev)
	require.NoError(t, err)
	require.True(t, added)
}

func (s *TestStore) addExistingEvent(t *testing.T, ev *types.Event) {
	t.Helper()
	added, err := s.AddEvent(ev)
	require.NoError(t, err)
	require.False(t, added)
}

func (s *TestStore) appendRandomEvent(t *testing.T) *types.Event {
	t.Helper()
	return s.appendEvent(t, newRandomEventPayload())
}

func (s *TestStore) appendEvent(t *testing.T, payload []byte) *types.Event {
	t.Helper()
	ev, err := s.AppendEvent(payload)
	require.NoError(t, err)
	return ev
}

func (s *TestStore) newRandomEvent(t *testing.T) *types.Event {
	t.Helper()
	return s.newRandomEventWithHeads(t, nil)
}

func (s *TestStore) newRandomEventWithHeads(t *testing.T, heads []multihash.Multihash) *types.Event {
	t.Helper()
	ev, err := types.NewEvent(newRandomEventPayload(), heads)
	require.NoError(t, err)
	return ev
}

func newRandomEventPayload() []byte {
	return []byte("event-" + test.RandomStringLower(13))
}

func (s *TestStore) requireEventsEqual(t *testing.T, expected []*types.Event) {
	events, err := s.Events()
	require.NoError(t, err)
	require.ElementsMatch(t, expected, events)
}
