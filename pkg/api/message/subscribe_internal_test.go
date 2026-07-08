package message

import (
	"context"
	"errors"
	"math"
	"sync"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
)

// TestMutableSubscriptionAddRaceWithClose exercises a mutate (addTopics, on the writer goroutine)
// racing the worker's reap (closeListener) on the same listener. Both touch l.closed; the fix sets
// it under topicsMu, so -race must stay clean. Before the fix, closeListener wrote l.closed without
// the lock addTopics reads it under — a data race that could re-register a listener in
// topicListeners after its channel was already closed (leak / send-on-closed panic).
func TestMutableSubscriptionAddRaceWithClose(t *testing.T) {
	worker := &subscribeWorker{}
	for range 200 {
		ch := make(chan []*envelopes.OriginatorEnvelope, 1)
		l := &listener{
			ctx:         context.Background(),
			ch:          ch,
			topics:      make(map[string]struct{}),
			originators: make(map[uint32]struct{}),
		}
		m := &mutableSubscription{worker: worker, l: l, ch: ch}

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			m.addTopics(map[string]struct{}{"race-topic": {}})
		}()
		go func() {
			defer wg.Done()
			worker.closeListener(l)
		}()
		wg.Wait()
	}
}

// TestSubscribeSessionStaleWaveSkipsResetTopic covers the reset bug: a topic removed and re-added
// while its first catch-up wave is still in flight is owned by the NEWER wave. The stale wave's
// completion must not open the new wave's gate or announce the topic; only the owning wave does.
func TestSubscribeSessionStaleWaveSkipsResetTopic(t *testing.T) {
	ctx := context.Background()
	sess := &subscribeSession{
		svc:       &Service{ctx: ctx},
		logger:    zap.NewNop(),
		ctx:       ctx,
		keepAlive: time.Second,
		outbound:  make(chan *message_api.SubscribeResponse, 16),
		topics:    make(map[string]*topicState),
		waves:     make(map[int]*subscribeWave),
	}

	tp := topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("reset-topic"))
	st := subscribeTopic{wire: tp.Bytes(), cursorKey: string(tp.Bytes()), listenKey: tp.String()}
	key := st.cursorKey

	// Wave 0 was in flight when the client removed + re-added the topic, so wave 1 now owns its gate.
	sess.waves[0] = &subscribeWave{mutateID: 10, topics: []subscribeTopic{st}}
	sess.waves[1] = &subscribeWave{mutateID: 20, topics: []subscribeTopic{st}}
	sess.topics[key] = &topicState{
		subscribeTopic: st,
		phase:          topicGated,
		wave:           1, // owned by wave 1
		cursor:         make(db.VectorClock),
	}

	// Stale wave 0 completes: it must NOT open the gate or announce the reset topic.
	done, err := sess.handleCatchUp(catchUpBatch{wave: 0, done: true})
	require.NoError(t, err)
	require.False(t, done)
	require.Equal(t, topicGated, sess.topics[key].phase, "stale wave must leave the topic gated")
	require.Equal(t, 1, sess.topics[key].wave, "stale wave must leave the gate owned by wave 1")

	f0 := drainResponses(sess.outbound)
	require.Empty(t, topicsLiveFrames(f0), "stale wave must not announce the reset topic")
	require.Equal(t, []uint64{10}, catchupCompleteIDs(f0))
	require.NotContains(t, sess.waves, 0)

	// Owning wave 1 completes: now the gate opens and the topic is announced.
	done, err = sess.handleCatchUp(catchUpBatch{wave: 1, done: true})
	require.NoError(t, err)
	require.False(t, done)
	require.Equal(
		t,
		topicLive,
		sess.topics[key].phase,
		"owning wave should open the gate (go live)",
	)

	f1 := drainResponses(sess.outbound)
	require.Len(t, topicsLiveFrames(f1), 1, "owning wave announces the topic")
	require.Equal(t, []uint64{20}, catchupCompleteIDs(f1))
}

// TestSubscribeSessionFlushPropagatesSenderError covers the false-OK teardown bug: on a graceful
// half-close path, if the sender goroutine already died with a Send error, flush() observes
// senderDone closed and must surface that error — not return a clean nil while the wave's terminal
// frames were never written to the wire.
func TestSubscribeSessionFlushPropagatesSenderError(t *testing.T) {
	sess := &subscribeSession{
		ctx:        context.Background(),
		keepAlive:  time.Second,
		outbound:   make(chan *message_api.SubscribeResponse, 1),
		senderDone: make(chan struct{}),
	}
	// Simulate the sender having stopped early on a Send error: record its terminal status, signal done.
	wantErr := errors.New("stream send failed")
	sess.sendErr = wantErr
	close(sess.senderDone)

	require.ErrorIs(t, sess.flush(), wantErr, "flush must surface the sender error, not a false OK")
}

// TestSubscribeSessionFlushPrefersDrainedSenderOverCancel covers the flush select-race finding: when
// the sender has drained (senderDone closed, sendErr nil) AND ctx is cancelled at the same instant,
// Go's select is pseudorandom — flush must deterministically prefer the completed drain and never
// report a clean graceful close as a spurious cancellation.
func TestSubscribeSessionFlushPrefersDrainedSenderOverCancel(t *testing.T) {
	for i := range 64 { // without the fix the ctx arm wins ~half the time; 64 iters reliably catches it
		ctx, cancel := context.WithCancel(context.Background())
		sess := &subscribeSession{
			ctx:        ctx,
			keepAlive:  time.Second,
			outbound:   make(chan *message_api.SubscribeResponse, 1),
			senderDone: make(chan struct{}),
		}
		close(sess.senderDone) // drained cleanly; sendErr stays nil
		cancel()
		require.NoError(t, sess.flush(),
			"a completed drain must not be reported as a cancellation (iter %d)", i)
	}
}

// TestSubscribeSessionDrainPendingRequestsHandlesClosedChannel covers the drainPendingRequests
// finding: a closed request channel must surface a transport error and, on a clean EOF, mark the
// session half-closed — so the pong-deadline caller does not then reap with a spurious DeadlineExceeded.
func TestSubscribeSessionDrainPendingRequestsHandlesClosedChannel(t *testing.T) {
	// Clean half-close: closed channel, no recvErr -> marks halfClosed, returns nil.
	clean := &subscribeSession{logger: zap.NewNop()}
	cleanCh := make(chan *message_api.SubscribeRequest)
	close(cleanCh)
	require.NoError(t, clean.drainPendingRequests(cleanCh))
	require.True(
		t,
		clean.halfClosed,
		"a clean EOF on the request channel must mark the session half-closed",
	)

	// Transport failure: closed channel WITH recvErr -> surfaces Unavailable, not a clean nil.
	failed := &subscribeSession{logger: zap.NewNop(), recvErr: errors.New("transport boom")}
	failedCh := make(chan *message_api.SubscribeRequest)
	close(failedCh)
	err := failed.drainPendingRequests(failedCh)
	require.Equal(
		t,
		connect.CodeUnavailable,
		connect.CodeOf(err),
		"a recv error must surface as Unavailable",
	)
	require.False(t, failed.halfClosed, "a transport failure is not a clean half-close")
}

// TestSubscribeSessionRemoveReaddDuringHistoryOnlyRejected covers the re-review's confirmed gap: the
// in-flight history_only overlap guard must run even on the remove+re-add (reset) path, because
// removeTopic does NOT cancel a history_only wave. Otherwise the surviving history_only wave and the
// new live wave both deliver the topic's history.
func TestSubscribeSessionRemoveReaddDuringHistoryOnlyRejected(t *testing.T) {
	ctx := context.Background()
	tp := topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("ho-topic"))
	st := subscribeTopic{wire: tp.Bytes(), cursorKey: string(tp.Bytes()), listenKey: tp.String()}

	sess := &subscribeSession{
		svc:       &Service{ctx: ctx},
		logger:    zap.NewNop(),
		ctx:       ctx,
		keepAlive: time.Second,
		topics:    make(map[string]*topicState),
		// An in-flight history_only wave for the topic (its throwaway cursors live on the wave).
		waves: map[int]*subscribeWave{
			0: {
				mutateID:    1,
				historyOnly: true,
				topics:      []subscribeTopic{st},
				cursors:     db.TopicCursors{st.cursorKey: make(db.VectorClock)},
			},
		},
	}

	// remove(T) + add(T) — the documented replay path — while T's history_only wave is still running.
	err := sess.handleMutate(&message_api.SubscribeRequest_V1_Mutate{
		MutateId: 2,
		Removes:  [][]byte{tp.Bytes()},
		Adds:     []*message_api.SubscribeRequest_V1_Mutate_Subscription{{Topic: tp.Bytes()}},
	})
	require.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err),
		"remove+re-add over an in-flight history_only wave must be rejected, not double-delivered")
}

// TestSubscribeSessionRejectsOutOfRangeCursor covers the cursor-overflow finding: a LastSeen entry
// whose value exceeds the signed DB column types would be silently dropped from the catch-up query
// and poison the live cursor (killing the topic). It must be rejected.
func TestSubscribeSessionRejectsOutOfRangeCursor(t *testing.T) {
	ctx := context.Background()
	tp := topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("ovf-topic"))
	newSess := func() *subscribeSession {
		return &subscribeSession{
			svc: &Service{ctx: ctx}, logger: zap.NewNop(), ctx: ctx, keepAlive: time.Second,
			topics: make(map[string]*topicState),
			waves:  make(map[int]*subscribeWave),
		}
	}
	mutate := func(c map[uint32]uint64) *message_api.SubscribeRequest_V1_Mutate {
		return &message_api.SubscribeRequest_V1_Mutate{
			MutateId: 1,
			Adds: []*message_api.SubscribeRequest_V1_Mutate_Subscription{
				{Topic: tp.Bytes(), LastSeen: &envelopesProto.Cursor{NodeIdToSequenceId: c}},
			},
		}
	}

	require.Equal(
		t,
		connect.CodeInvalidArgument,
		connect.CodeOf(
			newSess().handleMutate(mutate(map[uint32]uint64{100: uint64(math.MaxInt64) + 1})),
		),
		"sequence id beyond int64 must be rejected",
	)
	require.Equal(
		t,
		connect.CodeInvalidArgument,
		connect.CodeOf(
			newSess().handleMutate(mutate(map[uint32]uint64{uint32(math.MaxInt32) + 1: 5})),
		),
		"originator id beyond int32 must be rejected",
	)
}

// TestSubscribeSessionRejectsTooManyInflightWaves covers the in-flight catch-up cap (round-2 review
// finding): a Mutate that would start a new wave while maxInflightSubscribeWaves are already running
// is rejected, so reset churn (remove+re-add, which never cancels the old wave) cannot pile up
// fetcher goroutines. A removes-only Mutate starts no wave and is exempt even at the cap.
func TestSubscribeSessionRejectsTooManyInflightWaves(t *testing.T) {
	ctx := context.Background()
	sess := &subscribeSession{
		svc: &Service{ctx: ctx}, logger: zap.NewNop(), ctx: ctx, keepAlive: time.Second,
		outbound: make(chan *message_api.SubscribeResponse, 1),
		topics:   make(map[string]*topicState),
		waves:    make(map[int]*subscribeWave),
		// Minimal sub so the removes-only path's worker unregister is a no-op, not a nil deref.
		sub: &mutableSubscription{
			worker: &subscribeWorker{},
			l: &listener{
				topics:      make(map[string]struct{}),
				originators: make(map[uint32]struct{}),
			},
		},
	}
	// Saturate the in-flight wave budget.
	for i := range maxInflightSubscribeWaves {
		sess.waves[i] = &subscribeWave{mutateID: uint64(i)}
	}
	sess.nextWave = maxInflightSubscribeWaves

	// Mutate ids past the saturated waves' 0..cap-1 range, so the in-flight mutate_id collision
	// check cannot preempt the cap check this test targets.
	tp := topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("over-cap"))
	err := sess.handleMutate(&message_api.SubscribeRequest_V1_Mutate{
		MutateId: maxInflightSubscribeWaves + 1,
		Adds:     []*message_api.SubscribeRequest_V1_Mutate_Subscription{{Topic: tp.Bytes()}},
	})
	require.Equal(t, connect.CodeResourceExhausted, connect.CodeOf(err),
		"a Mutate that would exceed the in-flight catch-up cap must be rejected")

	// A removes-only Mutate creates no wave, so it is not blocked by the cap.
	err = sess.handleMutate(&message_api.SubscribeRequest_V1_Mutate{
		MutateId: maxInflightSubscribeWaves + 2,
		Removes:  [][]byte{tp.Bytes()},
	})
	require.NoError(t, err, "a removes-only Mutate must not be blocked by the in-flight cap")
}

// TestSubscribeSessionMutateAddsCapRejected covers the adds-per-Mutate cap (review finding): a
// wave's merged catch-up scan flattens one floor entry per (topic, originator) pair into EVERY
// page query, so the raw adds one Mutate may carry are bounded. An over-cap Mutate must be
// rejected with ResourceExhausted BEFORE any state changes — atomically, so not even the removes
// riding the same Mutate apply — and an at-cap Mutate must succeed (wave created, topics gated).
func TestSubscribeSessionMutateAddsCapRejected(t *testing.T) {
	ctx := context.Background()
	sess := &subscribeSession{
		svc:       &Service{ctx: ctx, originatorList: stubOriginatorLister{}},
		logger:    zap.NewNop(),
		ctx:       ctx,
		keepAlive: time.Second,
		outbound:  make(chan *message_api.SubscribeResponse, 16),
		catchUpCh: make(chan catchUpBatch, subscribeCatchUpQueueDepth),
		topics:    make(map[string]*topicState),
		waves:     make(map[int]*subscribeWave),
		maxAdds:   3,
		// Minimal sub so gateTopic/removeTopic's worker (un)register is a no-op, not a nil deref.
		sub: &mutableSubscription{
			worker: &subscribeWorker{},
			l: &listener{
				topics:      make(map[string]struct{}),
				originators: make(map[uint32]struct{}),
			},
		},
	}

	// A topic already live on the stream; the over-cap Mutate below also tries to remove it.
	pre := wbTopic("cap-pre")
	sess.topics[pre.cursorKey] = &topicState{
		subscribeTopic: pre,
		phase:          topicLive,
		cursor:         make(db.VectorClock),
	}

	adds := func(names ...string) []*message_api.SubscribeRequest_V1_Mutate_Subscription {
		out := make([]*message_api.SubscribeRequest_V1_Mutate_Subscription, 0, len(names))
		for _, n := range names {
			out = append(out, &message_api.SubscribeRequest_V1_Mutate_Subscription{
				Topic: wbTopic(n).wire,
			})
		}
		return out
	}

	// 4 adds > cap of 3: rejected — and atomically, so the remove riding the same Mutate must
	// not have applied either.
	err := sess.handleMutate(&message_api.SubscribeRequest_V1_Mutate{
		MutateId: 1,
		Adds:     adds("cap-a", "cap-b", "cap-c", "cap-d"),
		Removes:  [][]byte{pre.wire},
	})
	require.Equal(t, connect.CodeResourceExhausted, connect.CodeOf(err),
		"a Mutate over the adds cap must be rejected with ResourceExhausted")
	require.Len(t, sess.topics, 1, "no topic from the rejected Mutate may be registered")
	require.Contains(t, sess.topics, pre.cursorKey,
		"the rejected Mutate's removes must not have applied (atomicity)")
	require.Empty(t, sess.waves, "the rejected Mutate must not have created a wave")
	require.Zero(t, sess.nextWave)
	require.Empty(t, drainResponses(sess.outbound), "a rejected Mutate must emit no frames")

	// Exactly at the cap: accepted — wave created, its topics gated.
	require.NoError(t, sess.handleMutate(&message_api.SubscribeRequest_V1_Mutate{
		MutateId: 2,
		Adds:     adds("cap-a", "cap-b", "cap-c"),
	}))
	require.Len(t, sess.waves, 1, "an at-cap Mutate must create its wave")
	require.Len(t, sess.topics, 4, "the at-cap Mutate's topics must be registered")
	for _, n := range []string{"cap-a", "cap-b", "cap-c"} {
		ts := sess.topics[wbTopic(n).cursorKey]
		require.NotNil(t, ts)
		require.Equal(t, topicGated, ts.phase)
	}
	// The wave's fetcher started; the stub lister fails it immediately (no DB in a white-box
	// test), which also synchronizes the detached goroutine's exit before the test returns.
	b := <-sess.catchUpCh
	require.Error(t, b.err)
}

// TestSubscribeSessionMutateCursorEntriesCapRejected covers the adds cap's companion bound: the
// add COUNT cap alone does not bound a wave's (topic, originator) floor pairs, because a single
// add's cursor is a per-originator vector that may name arbitrarily many originators. An over-cap
// Mutate must be rejected with ResourceExhausted before any state changes — atomically, like the
// adds cap — and an at-cap Mutate must succeed.
func TestSubscribeSessionMutateCursorEntriesCapRejected(t *testing.T) {
	ctx := context.Background()
	sess := &subscribeSession{
		svc:              &Service{ctx: ctx, originatorList: stubOriginatorLister{}},
		logger:           zap.NewNop(),
		ctx:              ctx,
		keepAlive:        time.Second,
		outbound:         make(chan *message_api.SubscribeResponse, 16),
		catchUpCh:        make(chan catchUpBatch, subscribeCatchUpQueueDepth),
		topics:           make(map[string]*topicState),
		waves:            make(map[int]*subscribeWave),
		maxCursorEntries: 3,
		// Minimal sub so gateTopic/removeTopic's worker (un)register is a no-op, not a nil deref.
		sub: &mutableSubscription{
			worker: &subscribeWorker{},
			l: &listener{
				topics:      make(map[string]struct{}),
				originators: make(map[uint32]struct{}),
			},
		},
	}

	// A topic already live on the stream; the over-cap Mutate below also tries to remove it.
	pre := wbTopic("ce-pre")
	sess.topics[pre.cursorKey] = &topicState{
		subscribeTopic: pre,
		phase:          topicLive,
		cursor:         make(db.VectorClock),
	}

	// ONE add whose cursor names n originators: the add count stays far under maxMutateAdds, so
	// only the cursor-entry cap can reject it.
	addWithEntries := func(name string, n int) []*message_api.SubscribeRequest_V1_Mutate_Subscription {
		c := make(map[uint32]uint64, n)
		for i := range n {
			c[uint32(100+i)] = uint64(i + 1)
		}
		return []*message_api.SubscribeRequest_V1_Mutate_Subscription{
			{Topic: wbTopic(name).wire, LastSeen: &envelopesProto.Cursor{NodeIdToSequenceId: c}},
		}
	}

	// 4 entries > cap of 3: rejected — and atomically, so the remove riding the same Mutate must
	// not have applied either.
	err := sess.handleMutate(&message_api.SubscribeRequest_V1_Mutate{
		MutateId: 1,
		Adds:     addWithEntries("ce-a", 4),
		Removes:  [][]byte{pre.wire},
	})
	require.Equal(t, connect.CodeResourceExhausted, connect.CodeOf(err),
		"a Mutate over the cursor-entry cap must be rejected with ResourceExhausted")
	require.Len(t, sess.topics, 1, "no topic from the rejected Mutate may be registered")
	require.Contains(t, sess.topics, pre.cursorKey,
		"the rejected Mutate's removes must not have applied (atomicity)")
	require.Empty(t, sess.waves, "the rejected Mutate must not have created a wave")
	require.Zero(t, sess.nextWave)
	require.Empty(t, drainResponses(sess.outbound), "a rejected Mutate must emit no frames")

	// Exactly at the cap: accepted — wave created, its topic gated.
	require.NoError(t, sess.handleMutate(&message_api.SubscribeRequest_V1_Mutate{
		MutateId: 2,
		Adds:     addWithEntries("ce-b", 3),
	}))
	require.Len(t, sess.waves, 1, "an at-cap Mutate must create its wave")
	ts := sess.topics[wbTopic("ce-b").cursorKey]
	require.NotNil(t, ts)
	require.Equal(t, topicGated, ts.phase)
	// The wave's fetcher started; the stub lister fails it immediately (no DB in a white-box
	// test), which also synchronizes the detached goroutine's exit before the test returns.
	b := <-sess.catchUpCh
	require.Error(t, b.err)
}

// TestSubscribeSessionSendEnvelopesSplitsFrames covers the frame-splitting fix: a batch larger than
// the frame limit must go out as multiple frames, never one oversized (stream-aborting) frame — and
// every split frame must carry the caller's wave tag (XIP-83 server requirement 3).
func TestSubscribeSessionSendEnvelopesSplitsFrames(t *testing.T) {
	ctx := context.Background()
	sess := &subscribeSession{
		svc: &Service{ctx: ctx}, logger: zap.NewNop(), ctx: ctx, keepAlive: time.Second,
		outbound:      make(chan *message_api.SubscribeResponse, 16),
		maxFrameBytes: 120, // each env below is ~52 bytes → 2 fit per frame
	}
	mkEnv := func(payload int) *envelopesProto.OriginatorEnvelope {
		return &envelopesProto.OriginatorEnvelope{UnsignedOriginatorEnvelope: make([]byte, payload)}
	}
	require.NoError(t, sess.sendEnvelopes([]*envelopesProto.OriginatorEnvelope{
		mkEnv(50), mkEnv(50), mkEnv(50),
	}, 5))

	var counts []int
	total := 0
	for _, f := range drainResponses(sess.outbound) {
		if env := f.GetV1().GetEnvelopes(); env != nil {
			counts = append(counts, len(env.GetEnvelopes()))
			total += len(env.GetEnvelopes())
			require.Equal(t, uint64(5), env.GetMutateId(), "every split frame carries the wave tag")
		}
	}
	require.Len(t, counts, 2, "batch must be split into 2 frames")
	require.Equal(t, 3, total, "every envelope must be delivered exactly once")
}

// TestSubscribeSessionSendEnvelopesSkipsOverCapEnvelope covers frame splitting's hard-cap corner:
// publish admission bounds what the payer sent, not the stored envelope (originator wrapper) plus
// response framing, so a lone envelope can exceed the transport's send cap — and sending it would
// abort the stream, with every reconnecting wave hitting the same row again. It must be skipped
// with a warning (like the legacy batchAndSendEnvelopes) while its neighbors are still delivered
// in order under the caller's wave tag.
func TestSubscribeSessionSendEnvelopesSkipsOverCapEnvelope(t *testing.T) {
	ctx := context.Background()
	sess := &subscribeSession{
		svc: &Service{ctx: ctx}, logger: zap.NewNop(), ctx: ctx, keepAlive: time.Second,
		outbound: make(chan *message_api.SubscribeResponse, 16),
	}
	mkEnv := func(payload int) *envelopesProto.OriginatorEnvelope {
		return &envelopesProto.OriginatorEnvelope{UnsignedOriginatorEnvelope: make([]byte, payload)}
	}
	require.NoError(t, sess.sendEnvelopes([]*envelopesProto.OriginatorEnvelope{
		mkEnv(50), mkEnv(constants.GRPCPayloadLimit + 1024), mkEnv(60),
	}, 9))

	var sizes []int
	for _, f := range drainResponses(sess.outbound) {
		env := f.GetV1().GetEnvelopes()
		require.NotNil(t, env)
		require.Equal(t, uint64(9), env.GetMutateId(), "surviving frames must keep the wave tag")
		for _, e := range env.GetEnvelopes() {
			sizes = append(sizes, len(e.GetUnsignedOriginatorEnvelope()))
		}
	}
	require.Equal(t, []int{50, 60}, sizes,
		"the over-cap envelope must be skipped; its neighbors delivered in order")
}

// TestSubscribeSessionFoldTagsAndOrders is the deterministic pin on the wave-completion fold: the
// live envelopes buffered while the wave's topics were gated go out stamped with the wave's
// mutate_id, merged into per-originator sequence order across topics, then TopicsLive, then
// CatchupComplete — and never on the live (tag-0) lane.
func TestSubscribeSessionFoldTagsAndOrders(t *testing.T) {
	t1, t2 := wbTopic("fold-t1"), wbTopic("fold-t2")
	sess := newWaveTestSession(0, 42, t1, t2)

	// One scan page: cursors advance to (100 -> 1) on t1 and (100 -> 2) on t2.
	done, err := sess.handleCatchUp(catchUpBatch{wave: 0, envs: []*envelopes.OriginatorEnvelope{
		wbEnv(t, t1, 100, 1),
		wbEnv(t, t2, 100, 2),
	}})
	require.NoError(t, err)
	require.False(t, done)
	page := drainResponses(sess.outbound)
	require.Equal(t, [][2]uint64{{100, 1}, {100, 2}}, wbEnvelopeKeys(t, page, 42))
	require.Empty(t, wbEnvelopeKeys(t, page, 0), "no scan page may ride the live tag")

	// Live envelopes for the gated topics: buffered, nothing sent.
	require.NoError(t, sess.routeLive([]*envelopes.OriginatorEnvelope{
		wbEnv(t, t1, 100, 4),
		wbEnv(t, t2, 100, 3),
		wbEnv(t, t2, 200, 1),
	}))
	require.Empty(t, drainResponses(sess.outbound), "gated live envelopes must be withheld")
	require.Positive(t, sess.pendingBytes)

	done, err = sess.handleCatchUp(catchUpBatch{wave: 0, done: true})
	require.NoError(t, err)
	require.False(t, done)

	frames := drainResponses(sess.outbound)
	require.Equal(t, [][2]uint64{{100, 3}, {100, 4}, {200, 1}}, wbEnvelopeKeys(t, frames, 42),
		"fold must merge the pending buffers into per-originator sequence order, tagged 42")
	require.Empty(t, wbEnvelopeKeys(t, frames, 0), "no fold envelope may ride the live tag")

	lastEnv, liveIdx, ccIdx := -1, -1, -1
	for i, f := range frames {
		switch {
		case f.GetV1().GetEnvelopes() != nil:
			lastEnv = i
		case f.GetV1().GetTopicsLive() != nil:
			liveIdx = i
		case f.GetV1().GetCatchupComplete() != nil:
			ccIdx = i
		}
	}
	require.True(t, lastEnv < liveIdx && liveIdx < ccIdx,
		"fold frames must precede TopicsLive, which precedes CatchupComplete (%d, %d, %d)",
		lastEnv, liveIdx, ccIdx)
	require.Equal(t, [][]byte{t1.wire, t2.wire}, topicsLiveFrames(frames))
	require.Equal(t, []uint64{42}, catchupCompleteIDs(frames))

	require.Zero(t, sess.pendingBytes, "the fold must refund every buffered byte")
	require.Equal(t, topicLive, sess.topics[t1.cursorKey].phase)
	require.Equal(t, topicLive, sess.topics[t2.cursorKey].phase)
}

// TestSubscribeSessionFoldDedupsScanDeliveredPending pins exactly-once when the same envelope
// travels both lanes (a lagging worker dispatches a row the scan already delivered): the fold's
// cursor dedup must drop it, delivering only the genuinely-new pending envelope.
func TestSubscribeSessionFoldDedupsScanDeliveredPending(t *testing.T) {
	tp := wbTopic("dual-lane")
	sess := newWaveTestSession(0, 7, tp)

	done, err := sess.handleCatchUp(catchUpBatch{wave: 0, envs: []*envelopes.OriginatorEnvelope{
		wbEnv(t, tp, 100, 1),
		wbEnv(t, tp, 100, 2),
	}})
	require.NoError(t, err)
	require.False(t, done)
	drainResponses(sess.outbound)

	// The worker lagged: it dispatches (100,2) — already delivered by the scan — plus (100,3).
	require.NoError(t, sess.routeLive([]*envelopes.OriginatorEnvelope{
		wbEnv(t, tp, 100, 2),
		wbEnv(t, tp, 100, 3),
	}))

	done, err = sess.handleCatchUp(catchUpBatch{wave: 0, done: true})
	require.NoError(t, err)
	require.False(t, done)

	frames := drainResponses(sess.outbound)
	require.Equal(t, [][2]uint64{{100, 3}}, wbEnvelopeKeys(t, frames, 7),
		"the scan-delivered (100,2) must not be resent by the fold")
	require.Equal(t, uint64(3), sess.topics[tp.cursorKey].cursor[100])
}

// TestSubscribeSessionNewOriginatorMidWaveFoldsTagged: an originator that first publishes
// mid-wave has no ceiling row and no cursor entry; its gated envelopes must flow through the
// fold (tagged) and seed the live cursor, so its subsequent live tail is neither duplicated
// nor dropped.
func TestSubscribeSessionNewOriginatorMidWaveFoldsTagged(t *testing.T) {
	tp := wbTopic("new-orig")
	sess := newWaveTestSession(0, 9, tp)

	done, err := sess.handleCatchUp(catchUpBatch{wave: 0, envs: []*envelopes.OriginatorEnvelope{
		wbEnv(t, tp, 100, 1),
	}})
	require.NoError(t, err)
	require.False(t, done)
	drainResponses(sess.outbound)

	// Originator 200 first publishes mid-wave: the wave's ceilings never saw it.
	require.NoError(t, sess.routeLive([]*envelopes.OriginatorEnvelope{
		wbEnv(t, tp, 200, 1),
		wbEnv(t, tp, 200, 2),
	}))

	done, err = sess.handleCatchUp(catchUpBatch{wave: 0, done: true})
	require.NoError(t, err)
	require.False(t, done)
	frames := drainResponses(sess.outbound)
	require.Equal(t, [][2]uint64{{200, 1}, {200, 2}}, wbEnvelopeKeys(t, frames, 9),
		"a mid-wave originator's gated envelopes fold in ascending order, tagged")
	require.Equal(t, uint64(2), sess.topics[tp.cursorKey].cursor[200])

	// The topic is live now: the originator's tail arrives on the live (tag-0) lane.
	require.NoError(t, sess.routeLive([]*envelopes.OriginatorEnvelope{
		wbEnv(t, tp, 200, 3),
	}))
	live := drainResponses(sess.outbound)
	require.Equal(t, [][2]uint64{{200, 3}}, wbEnvelopeKeys(t, live, 0))
	require.Empty(t, wbEnvelopeKeys(t, live, 9))
}

// TestSubscribeSessionRemoveMidWaveDisposesPending pins the pending-bytes budget across a
// mid-wave remove: dropping the topic refunds its buffered bytes, the buffered envelopes are
// never delivered, and the orphaned wave still acks its own CatchupComplete without announcing
// or flushing anything.
func TestSubscribeSessionRemoveMidWaveDisposesPending(t *testing.T) {
	tp := wbTopic("remove-mid-wave")
	sess := newWaveTestSession(0, 5, tp)
	// Minimal sub so removeTopic's worker unregister is a no-op, not a nil deref.
	sess.sub = &mutableSubscription{
		worker: &subscribeWorker{},
		l: &listener{
			topics:      make(map[string]struct{}),
			originators: make(map[uint32]struct{}),
		},
	}

	require.NoError(t, sess.routeLive([]*envelopes.OriginatorEnvelope{
		wbEnv(t, tp, 100, 1),
		wbEnv(t, tp, 100, 2),
	}))
	require.Positive(t, sess.pendingBytes)

	// A removes-only Mutate is acked with its own CatchupComplete; the pending budget is refunded.
	require.NoError(t, sess.handleMutate(&message_api.SubscribeRequest_V1_Mutate{
		MutateId: 6,
		Removes:  [][]byte{tp.wire},
	}))
	require.Zero(t, sess.pendingBytes, "removing a gated topic must refund its pending bytes")
	require.Equal(t, []uint64{6}, catchupCompleteIDs(drainResponses(sess.outbound)))

	done, err := sess.handleCatchUp(catchUpBatch{wave: 0, done: true})
	require.NoError(t, err)
	require.False(t, done)
	frames := drainResponses(sess.outbound)
	for _, f := range frames {
		require.Nil(t, f.GetV1().GetEnvelopes(),
			"the removed topic's buffered envelopes must never be delivered")
	}
	require.Empty(t, topicsLiveFrames(frames), "an orphaned wave announces nothing")
	require.Equal(t, []uint64{5}, catchupCompleteIDs(frames),
		"the orphaned wave still acks its own CatchupComplete")
}

// TestSubscribeSessionDBErrorMidWaveFailsStreamNoCatchupComplete: a fetch error mid-wave must
// fail the stream (the client reconnects from its cursors) and never emit the wave's
// CatchupComplete over a history gap.
func TestSubscribeSessionDBErrorMidWaveFailsStreamNoCatchupComplete(t *testing.T) {
	tp := wbTopic("db-err")
	sess := newWaveTestSession(0, 13, tp)

	done, err := sess.handleCatchUp(catchUpBatch{wave: 0, err: errors.New("db down")})
	require.False(t, done)
	require.Equal(t, connect.CodeUnavailable, connect.CodeOf(err))
	require.Empty(t, catchupCompleteIDs(drainResponses(sess.outbound)),
		"a failed wave must never emit CatchupComplete")
}

// TestSubscribeSessionSendErrorMidFoldNoCatchupComplete: if the sender dies while the fold is
// being delivered, the error must propagate and the wave's CatchupComplete must not be emitted —
// the client must never believe it is synced past an undelivered fold.
func TestSubscribeSessionSendErrorMidFoldNoCatchupComplete(t *testing.T) {
	tp := wbTopic("send-err")
	sess := newWaveTestSession(0, 11, tp)
	// Unbuffered outbound with no reader plus a dead sender: send() can only observe senderDone.
	sess.outbound = make(chan *message_api.SubscribeResponse)
	sess.senderDone = make(chan struct{})
	wantErr := errors.New("stream send failed")
	sess.sendErr = wantErr
	close(sess.senderDone)

	require.NoError(t, sess.routeLive([]*envelopes.OriginatorEnvelope{
		wbEnv(t, tp, 100, 1),
	}))

	done, err := sess.handleCatchUp(catchUpBatch{wave: 0, done: true})
	require.False(t, done)
	require.ErrorIs(t, err, wantErr, "the sender's error must surface from the fold")
	require.Empty(t, catchupCompleteIDs(drainResponses(sess.outbound)),
		"CatchupComplete must not follow an undelivered fold")
}

// TestSubscribeSessionInflightMutateIdCollisionRejected covers the in-flight mutate_id collision
// rule (XIP-83 server requirement 3): a Mutate reusing the mutate_id of a wave still in flight is
// rejected — the two waves' replay frames and CatchupComplete acks would be indistinguishable —
// atomically, BEFORE any state changes. Even a removes-only reuse is rejected (its immediate
// CatchupComplete would be ambiguous with the in-flight wave's). Reuse after the wave's
// CatchupComplete stays legal.
func TestSubscribeSessionInflightMutateIdCollisionRejected(t *testing.T) {
	pre := wbTopic("collide-pre")
	sess := newWaveTestSession(0, 7, pre) // wave 0 in flight with mutateID 7, gating `pre`
	sess.svc = &Service{ctx: sess.ctx, originatorList: stubOriginatorLister{}}
	sess.catchUpCh = make(chan catchUpBatch, subscribeCatchUpQueueDepth)
	// Minimal sub so gateTopic/removeTopic's worker (un)register is a no-op, not a nil deref.
	sess.sub = &mutableSubscription{
		worker: &subscribeWorker{},
		l: &listener{
			topics:      make(map[string]struct{}),
			originators: make(map[uint32]struct{}),
		},
	}

	// Adds riding the in-flight id 7: rejected, and no state may have changed.
	fresh := wbTopic("collide-new")
	err := sess.handleMutate(&message_api.SubscribeRequest_V1_Mutate{
		MutateId: 7,
		Adds:     []*message_api.SubscribeRequest_V1_Mutate_Subscription{{Topic: fresh.wire}},
	})
	require.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err),
		"a Mutate reusing an in-flight mutate_id must be rejected")
	require.Len(t, sess.topics, 1, "the rejected Mutate must register no topic")
	require.NotContains(t, sess.topics, fresh.cursorKey)
	require.Len(t, sess.waves, 1, "the rejected Mutate must not create a wave")
	require.Contains(t, sess.waves, 0)
	require.Equal(t, 1, sess.nextWave)
	require.Empty(t, drainResponses(sess.outbound), "a rejected Mutate must emit no frames")

	// A removes-only reuse is rejected the same way, and its remove must not have applied.
	err = sess.handleMutate(&message_api.SubscribeRequest_V1_Mutate{
		MutateId: 7,
		Removes:  [][]byte{pre.wire},
	})
	require.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err),
		"a removes-only Mutate reusing an in-flight mutate_id must be rejected")
	require.Contains(t, sess.topics, pre.cursorKey,
		"the rejected Mutate's removes must not have applied")
	require.Empty(t, drainResponses(sess.outbound))

	// Wave 7 completes (its CatchupComplete goes out): id 7 is reusable again.
	done, err := sess.handleCatchUp(catchUpBatch{wave: 0, done: true})
	require.NoError(t, err)
	require.False(t, done)
	require.Equal(t, []uint64{7}, catchupCompleteIDs(drainResponses(sess.outbound)))

	require.NoError(t, sess.handleMutate(&message_api.SubscribeRequest_V1_Mutate{
		MutateId: 7,
		Adds:     []*message_api.SubscribeRequest_V1_Mutate_Subscription{{Topic: fresh.wire}},
	}), "reusing a mutate_id after its wave's CatchupComplete must be accepted")
	require.Contains(t, sess.waves, 1, "the reused id's wave must be created")
	// The wave's fetcher started; the stub lister fails it immediately (no DB in a white-box
	// test), which also synchronizes the detached goroutine's exit before the test returns.
	b := <-sess.catchUpCh
	require.Error(t, b.err)
}

// newWaveTestSession builds a writer-owned session with one in-flight live wave owning the given
// gated topics — the state gateTopic/handleMutate would have left, minus the fetcher goroutine
// (tests feed catchUpBatches directly).
func newWaveTestSession(wave int, mutateID uint64, topics ...subscribeTopic) *subscribeSession {
	ctx := context.Background()
	sess := &subscribeSession{
		svc:       &Service{ctx: ctx},
		logger:    zap.NewNop(),
		ctx:       ctx,
		keepAlive: time.Second,
		outbound:  make(chan *message_api.SubscribeResponse, 16),
		topics:    make(map[string]*topicState),
		waves:     make(map[int]*subscribeWave),
	}
	sess.waves[wave] = &subscribeWave{mutateID: mutateID, topics: topics}
	for _, st := range topics {
		sess.topics[st.cursorKey] = &topicState{
			subscribeTopic: st,
			phase:          topicGated,
			wave:           wave,
			cursor:         make(db.VectorClock),
		}
	}
	sess.nextWave = wave + 1
	return sess
}

// stubOriginatorLister gives handleMutate's detached fetcher goroutine a terminal originator
// lookup (an immediate error) so it exits without ever touching a DB.
type stubOriginatorLister struct{}

func (stubOriginatorLister) GetOriginatorNodeIDs(context.Context) ([]uint32, error) {
	return nil, errors.New("no originator list in white-box tests")
}

func wbTopic(name string) subscribeTopic {
	tp := topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte(name))
	return subscribeTopic{wire: tp.Bytes(), cursorKey: string(tp.Bytes()), listenKey: tp.String()}
}

func wbEnv(
	t *testing.T,
	st subscribeTopic,
	nodeID uint32,
	seqID uint64,
) *envelopes.OriginatorEnvelope {
	t.Helper()
	env, err := envelopes.NewOriginatorEnvelope(
		envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, nodeID, seqID, st.wire),
	)
	require.NoError(t, err)
	return env
}

// wbEnvelopeKeys returns the (originator, sequence) keys of envelopes carried by Envelopes frames
// stamped with the given wave tag, in frame order.
func wbEnvelopeKeys(
	t *testing.T,
	frames []*message_api.SubscribeResponse,
	tag uint64,
) [][2]uint64 {
	t.Helper()
	var keys [][2]uint64
	for _, f := range frames {
		env := f.GetV1().GetEnvelopes()
		if env == nil || env.GetMutateId() != tag {
			continue
		}
		for _, e := range env.GetEnvelopes() {
			parsed, err := envelopes.NewOriginatorEnvelope(e)
			require.NoError(t, err)
			keys = append(
				keys,
				[2]uint64{uint64(parsed.OriginatorNodeID()), parsed.OriginatorSequenceID()},
			)
		}
	}
	return keys
}

func drainResponses(ch chan *message_api.SubscribeResponse) []*message_api.SubscribeResponse {
	var out []*message_api.SubscribeResponse
	for {
		select {
		case f := <-ch:
			out = append(out, f)
		default:
			return out
		}
	}
}

func topicsLiveFrames(frames []*message_api.SubscribeResponse) [][]byte {
	var out [][]byte
	for _, f := range frames {
		if tl := f.GetV1().GetTopicsLive(); tl != nil {
			out = append(out, tl.GetTopics()...)
		}
	}
	return out
}

func catchupCompleteIDs(frames []*message_api.SubscribeResponse) []uint64 {
	var out []uint64
	for _, f := range frames {
		if cc := f.GetV1().GetCatchupComplete(); cc != nil {
			out = append(out, cc.GetMutateId())
		}
	}
	return out
}
