package crdt_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	crdt "github.com/xmtp/xmtpd/pkg/crdt"
	chanbroadcaster "github.com/xmtp/xmtpd/pkg/crdt/broadcasters/chan"
	memstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
	chansyncer "github.com/xmtp/xmtpd/pkg/crdt/syncers/chan"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func TestReplica_NewClose(t *testing.T) {
	log := test.NewLogger(t)
	ctx := context.Background()

	store := memstore.New(log)
	bc := chanbroadcaster.New(log)
	syncer := chansyncer.New(log)

	replica, err := crdt.NewReplica(ctx, log, store, bc, syncer, nil)
	require.NoError(t, err)
	require.NotNil(t, replica)

	err = replica.Close()
	require.NoError(t, err)
}

func TestReplica_BroadcastStore_SingleReplica(t *testing.T) {
	replica := newTestReplica(t)
	defer replica.Close()

	events := replica.broadcastRandom(t, 1)
	replica.requireEventuallyCapturedEvents(t, events)
	replica.requireEventuallyStoredEvents(t, events)
}

func TestReplica_BroadcastStore_TwoReplicas(t *testing.T) {
	replica1 := newTestReplica(t)
	defer replica1.Close()

	replica2 := newTestReplica(t)
	defer replica2.Close()

	// Add replica2 as peer of replica1, broadcast events via replica1, and
	// expect that both replicas eventually capture and store the events.
	replica1.addPeer(t, replica2)

	events1 := replica1.broadcastRandom(t, 1)

	replica1.requireEventuallyCapturedEvents(t, events1)
	replica1.requireEventuallyStoredEvents(t, events1)

	replica2.requireEventuallyCapturedEvents(t, events1)
	replica2.requireEventuallyStoredEvents(t, events1)

	// Broadcaster events via replica2, but with no peers yet, so expect that
	// replica1 captures and stores just it's originally broadcasted events,
	// and not the newly broadcasted events via replica2.
	events2 := replica2.broadcastRandom(t, 1)

	replica1.requireEventuallyCapturedEvents(t, events1)
	replica1.requireEventuallyStoredEvents(t, events1)

	events := append(events1, events2...)
	replica2.requireEventuallyCapturedEvents(t, events)
	replica2.requireEventuallyStoredEvents(t, events)

	// Add replica1 as peer of replica2, and expect that both replicas
	// eventually capture and store all events.
	// replica2.addPeer(t, replica1)

	// replica1.requireEventuallyCapturedEvents(t, events)
	// replica1.requireEventuallyStoredEvents(t, events)

	// replica2.requireEventuallyCapturedEvents(t, events)
	// replica2.requireEventuallyStoredEvents(t, events)
}

type testReplica struct {
	*crdt.Replica

	store  testStore
	bc     crdt.Broadcaster
	syncer crdt.Syncer

	capturedEvents     []*types.Event
	capturedEventsLock sync.RWMutex
}

type testStore interface {
	crdt.Store

	Events() ([]*types.Event, error)
}

func newTestReplica(t *testing.T) *testReplica {
	t.Helper()
	ctx := context.Background()
	log := test.NewLogger(t)

	store := memstore.New(log)
	bc := chanbroadcaster.New(log)
	syncer := chansyncer.New(log)

	tr := &testReplica{
		store:  store,
		bc:     bc,
		syncer: syncer,
	}

	replica, err := crdt.NewReplica(ctx, log, store, bc, syncer, func(ev *types.Event) {
		tr.capturedEventsLock.Lock()
		defer tr.capturedEventsLock.Unlock()
		tr.capturedEvents = append(tr.capturedEvents, ev)
	})
	require.NoError(t, err)
	tr.Replica = replica

	return tr
}

func (r *testReplica) CapturedEvents(t *testing.T) []*types.Event {
	t.Helper()
	r.capturedEventsLock.RLock()
	defer r.capturedEventsLock.RUnlock()
	return r.capturedEvents
}

func (r *testReplica) addPeer(t *testing.T, peer *testReplica) {
	t.Helper()
	switch bc := r.bc.(type) {
	case *chanbroadcaster.ChannelBroadcaster:
		switch peerBC := peer.bc.(type) {
		case *chanbroadcaster.ChannelBroadcaster:
			bc.AddPeer(peerBC)
		default:
			require.Fail(t, "peer broadcaster unknown")
		}
	default:
		require.Fail(t, "replica broadcaster unknown")
	}
}

func (r *testReplica) broadcast(t *testing.T, events []*types.Event) {
	t.Helper()
	for _, ev := range events {
		err := r.bc.Broadcast(ev)
		require.NoError(t, err)
	}
}

func (r *testReplica) broadcastRandom(t *testing.T, count int) []*types.Event {
	t.Helper()
	events := make([]*types.Event, count)
	for i := 0; i < count; i++ {
		ev, err := types.NewEvent([]byte("payload-"+test.RandomStringLower(13)), nil)
		require.NoError(t, err)
		events[i] = ev
	}
	r.broadcast(t, events)
	return events
}

func (r *testReplica) requireEventuallyCapturedEvents(t *testing.T, expected []*types.Event) {
	t.Helper()
	assert.Eventually(t, func() bool {
		return len(r.CapturedEvents(t)) == len(expected)
	}, time.Second, 10*time.Millisecond)
	require.Equal(t, expected, r.CapturedEvents(t))
}

func (r *testReplica) requireEventuallyStoredEvents(t *testing.T, expected []*types.Event) {
	t.Helper()
	assert.Eventually(t, func() bool {
		events, err := r.store.Events()
		require.NoError(t, err)
		return len(events) == len(expected)
	}, time.Second, 10*time.Millisecond)
	events, err := r.store.Events()
	require.NoError(t, err)
	require.Equal(t, expected, events)
}
