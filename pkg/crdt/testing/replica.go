package crdttest

import (
	"context"
	"fmt"
	"io"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/crdt"
	chanbroadcaster "github.com/xmtp/xmtpd/pkg/crdt/broadcasters/chan"
	memstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
	chansyncer "github.com/xmtp/xmtpd/pkg/crdt/syncers/chan"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	test "github.com/xmtp/xmtpd/pkg/testing"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type testReplica struct {
	*crdt.Replica

	store  testStore
	bc     crdt.Broadcaster
	syncer crdt.Syncer

	capturedEventCids  map[string]struct{}
	capturedEvents     []*types.Event
	capturedEventsLock sync.RWMutex
}

type testStore interface {
	crdt.Store

	Events() ([]*types.Event, error)
}

func NewTestReplica(t *testing.T) *testReplica {
	t.Helper()
	ctx := context.Background()
	log := test.NewLogger(t)

	store := memstore.New(log)
	bc := chanbroadcaster.New(log)
	syncer := chansyncer.New(log)

	tr := &testReplica{
		store:             store,
		bc:                bc,
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

func (r *testReplica) CapturedEvents(t *testing.T) []*types.Event {
	t.Helper()
	r.capturedEventsLock.RLock()
	defer r.capturedEventsLock.RUnlock()
	return r.capturedEvents
}

func (r *testReplica) AddPeer(t *testing.T, peer *testReplica) {
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

func (r *testReplica) Broadcast(t *testing.T, payloads [][]byte) []*types.Event {
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

func (r *testReplica) BroadcastRandom(t *testing.T, count int) []*types.Event {
	t.Helper()
	payloads := make([][]byte, count)
	for i := 0; i < count; i++ {
		payloads[i] = []byte("payload-" + test.RandomStringLower(13))
	}
	return r.Broadcast(t, payloads)
}

func (r *testReplica) RequireEventuallyCapturedEvents(t *testing.T, expected []*types.Event) {
	t.Helper()
	assert.Eventually(t, func() bool {
		return len(r.CapturedEvents(t)) == len(expected)
	}, time.Second, 10*time.Millisecond)
	require.Equal(t, expected, r.CapturedEvents(t))
}

func (r *testReplica) RequireEventuallyStoredEvents(t *testing.T, expected []*types.Event) {
	t.Helper()
	assert.Eventually(t, func() bool {
		events, err := r.store.Events()
		require.NoError(t, err)
		return len(events) == len(expected)
	}, time.Second, 10*time.Millisecond)
	events, err := r.store.Events()
	sort.Slice(events, func(i, j int) bool {
		return events[i].Cid.String() < events[j].Cid.String()
	})
	expected = expected[:]
	sort.Slice(expected, func(i, j int) bool {
		return expected[i].Cid.String() < expected[j].Cid.String()
	})
	require.NoError(t, err)
	require.Equal(t, expected, events)
}

// visualize emits a graphviz depiction of the topic contents showing the
// graph of individual events and their links.
func (r *testReplica) Visualize(w io.Writer, name string) {
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
