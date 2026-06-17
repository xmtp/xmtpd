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

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
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

	tp := topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("over-cap"))
	err := sess.handleMutate(&message_api.SubscribeRequest_V1_Mutate{
		MutateId: 1,
		Adds:     []*message_api.SubscribeRequest_V1_Mutate_Subscription{{Topic: tp.Bytes()}},
	})
	require.Equal(t, connect.CodeResourceExhausted, connect.CodeOf(err),
		"a Mutate that would exceed the in-flight catch-up cap must be rejected")

	// A removes-only Mutate creates no wave, so it is not blocked by the cap.
	err = sess.handleMutate(&message_api.SubscribeRequest_V1_Mutate{
		MutateId: 2,
		Removes:  [][]byte{tp.Bytes()},
	})
	require.NoError(t, err, "a removes-only Mutate must not be blocked by the in-flight cap")
}

// TestAdvanceCursorsFromRowsAdvancesPastUnmarshalFailures covers the catch-up spin fix: pagination
// cursors must advance from the RAW rows even when an envelope's bytes don't unmarshal, otherwise a
// bad row in a full page re-fetches forever and the wave never completes.
func TestAdvanceCursorsFromRowsAdvancesPastUnmarshalFailures(t *testing.T) {
	tp := topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("rows-topic"))
	key := string(tp.Bytes())
	cursors := db.TopicCursors{key: make(db.VectorClock)}
	// Envelope bytes are intentionally garbage (would be dropped by unmarshalEnvelopes).
	rows := []queries.GatewayEnvelopesView{
		{
			Topic:                tp.Bytes(),
			OriginatorNodeID:     100,
			OriginatorSequenceID: 5,
			OriginatorEnvelope:   []byte("garbage"),
		},
		{
			Topic:                tp.Bytes(),
			OriginatorNodeID:     100,
			OriginatorSequenceID: 6,
			OriginatorEnvelope:   []byte("garbage"),
		},
		{
			Topic:                tp.Bytes(),
			OriginatorNodeID:     200,
			OriginatorSequenceID: 3,
			OriginatorEnvelope:   []byte("garbage"),
		},
	}
	advanceCursorsFromRows(cursors, rows)
	require.Equal(
		t,
		uint64(6),
		cursors[key][100],
		"cursor must advance to the max raw seq for an originator",
	)
	require.Equal(t, uint64(3), cursors[key][200])
}

// TestSubscribeSessionSendEnvelopesSplitsFrames covers the frame-splitting fix: a batch larger than
// the frame limit must go out as multiple frames, never one oversized (stream-aborting) frame.
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
	}))

	var counts []int
	total := 0
	for _, f := range drainResponses(sess.outbound) {
		if env := f.GetV1().GetEnvelopes(); env != nil {
			counts = append(counts, len(env.GetEnvelopes()))
			total += len(env.GetEnvelopes())
		}
	}
	require.Len(t, counts, 2, "batch must be split into 2 frames")
	require.Equal(t, 3, total, "every envelope must be delivered exactly once")
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
