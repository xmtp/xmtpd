package crdt_test

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	crdt "github.com/xmtp/xmtpd/pkg/crdt"
	membroadcaster "github.com/xmtp/xmtpd/pkg/crdt/broadcasters/mem"
	memstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
	memsyncer "github.com/xmtp/xmtpd/pkg/crdt/syncers/mem"
	crdttest "github.com/xmtp/xmtpd/pkg/crdt/testing"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	test "github.com/xmtp/xmtpd/pkg/testing"
	"github.com/xmtp/xmtpd/pkg/zap"
)

var (
	visualize bool
)

func init() {
	flag.BoolVar(&visualize, "visualize", false, "output graphviz depiction of replicas")
}

func TestReplica_NewClose(t *testing.T) {
	log := test.NewLogger(t)
	ctx := context.Background()

	store := memstore.New(log)
	bc := membroadcaster.New(log)
	syncer := memsyncer.New(log, store)

	replica, err := crdt.NewReplica(ctx, log, store, bc, syncer, nil)
	require.NoError(t, err)
	require.NotNil(t, replica)

	err = replica.Close()
	require.NoError(t, err)
}

func TestReplica_BroadcastStore_SingleReplica(t *testing.T) {
	replica := newTestReplica(t)
	defer replica.Close()

	events := replica.broadcastRandom(t, 3)
	replica.requireEventuallyCapturedEvents(t, events)
	replica.requireEventuallyStoredEvents(t, events)

	if visualize {
		replica.visualize(os.Stdout, "replica1")
	}
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

	store  *crdttest.TestStore
	bc     *crdttest.TestBroadcaster
	syncer crdt.Syncer

	capturedEventCids  map[string]struct{}
	capturedEvents     []*types.Event
	capturedEventsLock sync.RWMutex
}

func newTestReplica(t *testing.T) *testReplica {
	t.Helper()
	ctx := context.Background()
	log := test.NewLogger(t)

	store := memstore.New(log)
	bc := membroadcaster.New(log)
	syncer := memsyncer.New(log, store)

	tr := &testReplica{
		store:             crdttest.NewTestStore(store),
		bc:                crdttest.NewTestBroadcaster(ctx, log, bc),
		syncer:            syncer,
		capturedEventCids: map[string]struct{}{},
	}

	replica, err := crdt.NewReplica(ctx, log, store, bc, syncer, func(ev *types.Event) {
		tr.capturedEventsLock.Lock()
		defer tr.capturedEventsLock.Unlock()
		if _, ok := tr.capturedEventCids[ev.Cid.String()]; ok {
			log.Debug("ignore duplicate event during capture", zap.Cid("event", ev.Cid))
			return
		}
		tr.capturedEventCids[ev.Cid.String()] = struct{}{}
		tr.capturedEvents = append(tr.capturedEvents, ev)
	})
	require.NoError(t, err)
	tr.Replica = replica

	return tr
}

func (r *testReplica) getCapturedEvents(t *testing.T) []*types.Event {
	t.Helper()
	r.capturedEventsLock.RLock()
	defer r.capturedEventsLock.RUnlock()
	return r.capturedEvents
}

func (r *testReplica) addPeer(t *testing.T, peer *testReplica) {
	t.Helper()
	r.bc.AddPeer(peer.bc.ITestBroadcaster)
}

func (r *testReplica) broadcast(t *testing.T, payloads [][]byte) []*types.Event {
	t.Helper()
	ctx := context.Background()
	events := make([]*types.Event, len(payloads))
	for i, payload := range payloads {
		ev, err := r.Replica.BroadcastAppend(ctx, payload)
		require.NoError(t, err)
		events[i] = ev
	}
	return events
}

func (r *testReplica) broadcastRandom(t *testing.T, count int) []*types.Event {
	t.Helper()
	payloads := make([][]byte, count)
	for i := 0; i < count; i++ {
		payloads[i] = []byte("event-" + test.RandomStringLower(13))
	}
	return r.broadcast(t, payloads)
}

func (r *testReplica) requireEventuallyCapturedEvents(t *testing.T, expected []*types.Event) {
	t.Helper()
	assert.Eventually(t, func() bool {
		return len(r.getCapturedEvents(t)) == len(expected)
	}, time.Second, 10*time.Millisecond)
	require.Equal(t, expected, r.getCapturedEvents(t))
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
	require.ElementsMatch(t, expected, events)
}

// visualize emits a graphviz depiction of the topic contents showing the
// graph of individual events and their links.
func (r *testReplica) visualize(w io.Writer, name string) {
	fmt.Fprintf(w, "strict digraph %s {\n", name)
	for i := len(r.capturedEvents) - 1; i >= 0; i-- {
		ev := r.capturedEvents[i]
		fmt.Fprintf(w, "\t\"%s\" [label=\"%d: \\N\"]\n", zap.ShortCid(ev.Cid), i)
		fmt.Fprintf(w, "\t\"%s\" -> { ", zap.ShortCid(ev.Cid))
		for _, l := range ev.Links {
			fmt.Fprintf(w, "\"%s\" ", zap.ShortCid(l))
		}
		fmt.Fprintf(w, "}\n")
	}
	fmt.Fprintf(w, "}\n")
}
