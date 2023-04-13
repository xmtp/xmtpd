package crdt_test

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/context"
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
	ctx := test.NewContext(t)

	store := memstore.New(ctx)
	bc := membroadcaster.New(ctx)
	syncer := memsyncer.New(ctx, store)

	replica, err := crdt.NewReplica(ctx, nil, store, bc, syncer, nil)
	require.NoError(t, err)
	require.NotNil(t, replica)

	ctx.Close()
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
	ctx context.Context

	store  *crdttest.TestStore
	bc     *crdttest.TestBroadcaster
	syncer crdt.Syncer

	capturedEventCids  map[string]struct{}
	capturedEvents     []*types.Event
	capturedEventsLock sync.RWMutex
}

func newTestReplica(t *testing.T) *testReplica {
	t.Helper()
	ctx := test.NewContext(t)

	store := memstore.New(ctx)
	bc := membroadcaster.New(ctx)
	syncer := memsyncer.New(ctx, store)

	tr := &testReplica{
		ctx:               ctx,
		store:             crdttest.NewTestStore(ctx, store),
		bc:                crdttest.NewTestBroadcaster(ctx, bc),
		syncer:            crdttest.NewTestSyncer(ctx, syncer),
		capturedEventCids: map[string]struct{}{},
	}

	replica, err := crdt.NewReplica(ctx, nil, store, bc, syncer, func(ev *types.Event) {
		tr.capturedEventsLock.Lock()
		defer tr.capturedEventsLock.Unlock()
		if _, ok := tr.capturedEventCids[ev.Cid.String()]; ok {
			ctx.Logger().Debug("ignore duplicate event during capture", zap.Cid("event", ev.Cid))
			return
		}
		tr.capturedEventCids[ev.Cid.String()] = struct{}{}
		tr.capturedEvents = append(tr.capturedEvents, ev)
	})
	require.NoError(t, err)
	tr.Replica = replica

	return tr
}

func (r *testReplica) Close() {
	r.ctx.Close()
}

func (r *testReplica) getCapturedEvents(t *testing.T) []*types.Event {
	t.Helper()
	r.capturedEventsLock.RLock()
	defer r.capturedEventsLock.RUnlock()
	return r.capturedEvents
}

func (r *testReplica) addPeer(t *testing.T, peer *testReplica) {
	t.Helper()
	r.bc.AddPeer(t, peer.bc.ITestBroadcaster)
}

func (r *testReplica) broadcast(t *testing.T, envs []*messagev1.Envelope) []*types.Event {
	t.Helper()
	ctx := test.NewContext(t)
	events := make([]*types.Event, len(envs))
	for i, env := range envs {
		ev, err := r.Replica.BroadcastAppend(ctx, env)
		require.NoError(t, err)
		events[i] = ev
	}
	return events
}

func (r *testReplica) broadcastRandom(t *testing.T, count int) []*types.Event {
	t.Helper()
	envs := make([]*messagev1.Envelope, count)
	for i := 0; i < count; i++ {
		envs[i] = newRandomEnvelope(t)
	}
	return r.broadcast(t, envs)
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
	ctx := test.NewContext(t)
	assert.Eventually(t, func() bool {
		events, err := r.store.Events(ctx)
		require.NoError(t, err)
		return len(events) == len(expected)
	}, time.Second, 10*time.Millisecond)
	events, err := r.store.Events(ctx)
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

func newRandomEnvelope(t *testing.T) *messagev1.Envelope {
	return &messagev1.Envelope{
		ContentTopic: "topic-" + test.RandomStringLower(5),
		TimestampNs:  uint64(rand.Intn(100)),
		Message:      []byte("msg-" + test.RandomString(13)),
	}
}
