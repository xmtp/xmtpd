package message

import (
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"math"
	"time"

	"connectrpc.com/connect"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
)

const (
	// subscribeSendQueueDepth buffers response batches between the writer and the sender
	// goroutine, so a slow stream.Send never parks the writer (which must stay free for the
	// liveness ping/pong reap). Small: it only smooths Send latency.
	subscribeSendQueueDepth = 8
	// subscribeCatchUpQueueDepth buffers catch-up pages from the fetcher to the writer.
	subscribeCatchUpQueueDepth = 16
	// maxActiveSubscribeTopics caps the live topic set a single stream may hold. Intentionally
	// high: a herald multiplexing many identities wants one fat connection, not a fan-out across
	// connections (which would just hit rate limits). At ~hundreds of bytes/topic this is a
	// large-but-bounded budget; exceeding it fails the Mutate with ResourceExhausted.
	maxActiveSubscribeTopics = 1_000_000
	// maxInflightSubscribeWaves caps the concurrent catch-up waves one stream may have running. Each
	// wave is a detached fetcher goroutine plus paginated DB queries, and maxActiveSubscribeTopics
	// does NOT bound it: a remove+re-add (reset) leaves the old wave running while the active-topic
	// count stays flat, so unbounded reset churn could otherwise pile up fetchers against the shared
	// DB pool. A Mutate that would start a wave past this limit is rejected with ResourceExhausted
	// (retryable once in-flight catch-ups drain); a well-behaved client batches adds into few Mutates.
	maxInflightSubscribeWaves = 256
	// maxSubscribePendingBytes caps live envelopes buffered while a topic catches up. Exceeding
	// it tears the stream down (the client reconnects from its cursors) rather than risk OOM.
	maxSubscribePendingBytes = 64 << 20 // 64 MiB
	// maxSubscribeFrameBytes bounds a single Envelopes response frame. A catch-up page or a flushed
	// pending buffer can be large; splitting keeps every frame well under the gRPC message limit so
	// the client can receive it (an unsplit multi-MB frame would abort the stream).
	maxSubscribeFrameBytes = 2 << 20 // 2 MiB
)

// subscribeTopic is the immutable identity of one topic in a stream's subscription. cursorKey
// (string of the raw topic bytes) keys the per-topic state and the catch-up query; listenKey (the
// parsed topic string) keys the worker's dispatch registration; wire is the client's topic bytes,
// echoed in TopicsLive.
type subscribeTopic struct {
	wire      []byte
	cursorKey string
	listenKey string
}

// topicPhase is a live topic's position in its gated -> live lifecycle.
type topicPhase uint8

const (
	// topicGated: a catch-up wave is replaying this topic's history; its live envelopes are buffered
	// in pending and withheld until the owning wave opens the gate.
	topicGated topicPhase = iota
	// topicLive: caught up; live envelopes are deduped against cursor and delivered immediately.
	topicLive
)

// topicState is the single source of truth for one LIVE topic on a stream: its phase, the wave that
// owns it while gated, its sparse dedup cursor, and the live envelopes buffered while it catches up.
// Collapsing what were four parallel maps (catchingUp, cursors, liveTopics, pending) into one struct
// per topic means a transition cannot advance one aspect and forget another — every mutation goes
// through a session method (gateTopic / bufferLive / flushAndGoLive / removeTopicState) that moves
// the whole struct together. history_only topics never become live, so they are NOT represented
// here; their throwaway dedup cursors live on the wave instead. Owned solely by the writer
// goroutine: no locking.
type topicState struct {
	subscribeTopic
	phase   topicPhase
	wave    int                             // the wave catching this topic up while topicGated
	cursor  db.VectorClock                  // sparse live cursor: provided floor, grown as originators are seen
	pending []*envelopes.OriginatorEnvelope // live envelopes held while topicGated
}

// subscribeWave tracks one Mutate's catch-up. Its CatchupComplete (echoing mutate_id) is emitted
// once the fetcher signals done, after the wave's TopicsLive. cursors is non-nil only for a
// history_only wave (its throwaway per-topic dedup cursors); a live wave dedups against each topic's
// own live cursor instead.
type subscribeWave struct {
	mutateID    uint64
	historyOnly bool
	topics      []subscribeTopic
	cursors     db.TopicCursors // history_only only: throwaway dedup cursors
}

// catchUpBatch is one unit handed from a fetcher goroutine to the writer: a page of raw history
// (the writer dedups + sends), a done marker (the wave finished), or an error (tear down).
type catchUpBatch struct {
	wave int
	envs []*envelopes.OriginatorEnvelope
	done bool
	err  error
}

// Subscribe is the XIP-83 bidirectional mutable subscription (the decentralized binding; see the
// v3 MlsApi.Subscribe for the other). One long-lived stream the client mutates in place via
// add/remove topic deltas (no reconnect on membership change), with a WebSocket-style ping/pong so
// silent stream death is detected on both ends. In contrast to SubscribeTopics (a fixed,
// server-streaming filter set with a one-way WAITING heartbeat), this RPC is bidirectional.
//
// Concurrency model: SINGLE WRITER. The select loop is the sole owner of all mutable state (the
// per-topic vector cursors, the catch-up gate, the pending buffer, the ping bookkeeping). It is
// the only goroutine that decides WHAT to send and in what order; the actual stream.Send runs on
// one dedicated sender goroutine fed by an ordered channel, so a client that stops reading can
// never park the writer. The frame reader, the live envelope worker, and the catch-up fetchers
// are pure producers. Catch-up runs OFF the writer (in a fetcher goroutine) so a slow new topic
// never holds up live delivery for already-live topics; its live envelopes are gated in a pending
// buffer and flushed (deduped) when its catch-up completes.
func (s *Service) Subscribe(
	ctx context.Context,
	stream *connect.BidiStream[message_api.SubscribeRequest, message_api.SubscribeResponse],
) error {
	logger := s.logger.With(utils.MethodField(stream.Spec().Procedure))
	keepAlive := s.options.SendKeepAliveInterval

	// connect-go does NOT cancel the stream context when the handler returns, so derive a cancelable
	// child and cancel it on every return path. That tears down the catch-up fetchers and the
	// worker's listener promptly instead of leaking them (a fetcher parked on a full catchUpCh, or
	// the listener registration) until the outer ServeHTTP unwinds.
	streamCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	sess := &subscribeSession{
		svc:        s,
		logger:     logger,
		ctx:        streamCtx,
		keepAlive:  keepAlive,
		sub:        s.subscribeWorker.newMutableSubscription(streamCtx),
		outbound:   make(chan *message_api.SubscribeResponse, subscribeSendQueueDepth),
		senderDone: make(chan struct{}),
		catchUpCh:  make(chan catchUpBatch, subscribeCatchUpQueueDepth),
		topics:     make(map[string]*topicState),
		waves:      make(map[int]*subscribeWave),
	}
	defer sess.sub.close()

	// Sender goroutine: the SOLE caller of stream.Send, fed an ordered channel.
	go func() {
		var err error
		for resp := range sess.outbound {
			if sendErr := stream.Send(resp); sendErr != nil {
				err = sendErr
				break
			}
		}
		// Record the terminal status, THEN signal done. The write happens-before the close, and
		// every reader reads sendErr only after observing senderDone closed, so this is race-free.
		sess.sendErr = err
		close(sess.senderDone)
	}()
	// On every return path, try to stop the sender before connect finalizes the stream: a stream.Send
	// racing connect's stream Close both flush the shared HTTP/2 response writer. Graceful paths
	// already wait via flush(); this backstops error/reap returns. The wait is BOUNDED, and that bound
	// is a deliberate tradeoff: connect's stream.Send is not cancelable, so a sender wedged inside
	// Send on a non-reading client (its HTTP/2 window full) cannot be unblocked from here. We accept a
	// narrow residual window where Close may run concurrently with that wedged Send (a -race report on
	// a connection that is being torn down regardless) rather than hang teardown forever waiting on a
	// sender that may never return. Closing outbound only stops the NEXT send, not one already parked.
	defer func() {
		sess.closeOutbound()
		select {
		case <-sess.senderDone:
		case <-time.After(sess.keepAlive):
		}
	}()

	// Frame reader goroutine: pure producer of client frames.
	requestCh := make(chan *message_api.SubscribeRequest, 16)
	go func() {
		defer close(requestCh)
		for {
			req, err := stream.Receive()
			if err != nil {
				if !errors.Is(err, io.EOF) {
					sess.recvErr = err
				}
				return
			}
			select {
			case requestCh <- req:
			case <-streamCtx.Done():
				return
			}
		}
	}()

	if err := sess.send(newSubscribeStarted(uint32(keepAlive.Milliseconds()))); err != nil {
		return err
	}

	// Liveness uses two INDEPENDENT timers. pingTicker is the send-idle ping cadence — outbound
	// delivery resets it, so a stream actively receiving data is not pinged. pongDeadline is the reap
	// deadline, armed ONLY when a Ping is sent and disarmed ONLY by a matching Pong, so steady
	// outbound traffic can never postpone it. (A single shared ticker let ordinary delivery keep
	// resetting the reap, defeating silent-death detection.) A busy-but-dead client is still caught:
	// its unread window fills and the sender's send() trips its own keepAlive stall timeout.
	pingTicker := time.NewTicker(keepAlive)
	defer pingTicker.Stop()
	sess.pingTicker = pingTicker
	pongDeadline := time.NewTimer(keepAlive)
	stopTimer(pongDeadline) // not armed until a Ping is outstanding
	defer pongDeadline.Stop()
	requestChannel := requestCh // nilled on half-close so the case goes dormant

	for {
		select {
		case batch, open := <-sess.sub.ch:
			if !open {
				// The worker reaped us (ctx done, or we fell behind: channel full). The stream
				// cannot continue in order; fail so the client reconnects from its cursors.
				return connect.NewError(
					connect.CodeAborted,
					errors.New("subscription closed: consumer too slow"),
				)
			}
			if err := sess.routeLive(batch); err != nil {
				return err
			}

		case b := <-sess.catchUpCh:
			done, err := sess.handleCatchUp(b)
			if err != nil {
				return err
			}
			if done {
				return sess.flush()
			}

		case req, open := <-requestChannel:
			if !open {
				if sess.recvErr != nil {
					return connect.NewError(
						connect.CodeUnavailable,
						fmt.Errorf("stream recv failed: %w", sess.recvErr),
					)
				}
				// Clean half-close (io.EOF): finish in-flight catch-up waves, then close. If
				// nothing is in flight, drain and return now.
				if len(sess.waves) == 0 {
					return sess.flush()
				}
				sess.halfClosed = true
				requestChannel = nil // dormant: a closed channel would busy-loop this case
				continue
			}
			wasAwaitingPong := sess.awaitingPong
			if err := sess.handleRequest(req); err != nil {
				return err
			}
			if wasAwaitingPong && !sess.awaitingPong {
				stopTimer(pongDeadline) // a matching Pong arrived; disarm the reap deadline
			}

		case <-sess.senderDone:
			// The sender exited while we are still producing — only possible on a Send error, since
			// outbound is closed solely at teardown. Surface that error (sendErr is set before the
			// close we just observed).
			return sess.sendErr

		case <-pingTicker.C:
			// Ping only when idle, not already awaiting a Pong, and not half-closed (the client's
			// read side is gone — it cannot answer; we are only draining in-flight waves).
			if !sess.awaitingPong && !sess.halfClosed {
				sess.pingNonce++
				sess.awaitingPong = true
				pongDeadline.Reset(keepAlive)
				if err := sess.send(newSubscribePing(sess.pingNonce)); err != nil {
					return err
				}
			}

		case <-pongDeadline.C:
			// The reap deadline fired. A Pong may be sitting in requestChannel right as it fired;
			// process queued frames first so select fairness can't reap a stream that did answer.
			if err := sess.drainPendingRequests(requestChannel); err != nil {
				return err
			}
			if sess.awaitingPong && !sess.halfClosed {
				return connect.NewError(
					connect.CodeDeadlineExceeded,
					errors.New("no Pong within deadline"),
				)
			}

		case <-streamCtx.Done():
			return nil

		case <-s.ctx.Done():
			return connect.NewError(connect.CodeUnavailable, errors.New("service is shutting down"))
		}
	}
}

// subscribeSession holds one stream's writer-owned state. All methods run on the writer goroutine,
// so there are no locks; the sender, reader, and fetcher goroutines communicate only via channels.
type subscribeSession struct {
	svc       *Service
	logger    *zap.Logger
	ctx       context.Context
	keepAlive time.Duration
	sub       *mutableSubscription
	// pingTicker is the send-idle ping cadence; send() resets it on every frame actually delivered,
	// so the cadence tracks real outbound bytes (not internal events that produce no output, e.g. a
	// fully-gated catch-up page). nil in white-box unit tests that drive session methods directly.
	pingTicker *time.Ticker
	// sendTimer bounds a single send() (reused instead of allocating time.After per call); lazily
	// created on first send so it works in white-box tests too.
	sendTimer *time.Timer
	// maxFrameBytes overrides maxSubscribeFrameBytes when > 0 (tests only).
	maxFrameBytes int

	outbound chan *message_api.SubscribeResponse
	// senderDone is closed exactly once, when the sender goroutine exits; sendErr is its terminal
	// status (nil = every queued frame was sent, non-nil = the Send error that stopped it). Together
	// they are the sender's single result contract: every reader (send, flush, the writer loop)
	// learns the outcome by observing senderDone closed and then reading sendErr — a broadcast, so
	// no reader can consume the signal out from under another (the old buffered errCh could).
	senderDone chan struct{}
	sendErr    error
	outClosed  bool
	catchUpCh  chan catchUpBatch
	recvErr    error

	// topics is every LIVE topic on this stream, keyed by cursorKey, each a small gated->live state
	// machine (see topicState). It replaces the former catchingUp / cursors / liveTopics / pending
	// maps so per-topic state can only move as a unit. Capped at maxActiveSubscribeTopics.
	topics       map[string]*topicState
	pendingBytes int // sum of proto sizes buffered across every topic's pending; bounded
	waves        map[int]*subscribeWave
	nextWave     int

	awaitingPong bool
	pingNonce    uint64
	halfClosed   bool
}

func (sess *subscribeSession) closeOutbound() {
	if !sess.outClosed {
		sess.outClosed = true
		close(sess.outbound)
	}
}

// flush drains queued frames before a GRACEFUL completion so no tail is lost. Bounded by the
// keepalive deadline / ctx; if that bound trips the drain did NOT finish, so it returns
// DeadlineExceeded rather than a false OK (the bug the v3 binding's flush originally had).
func (sess *subscribeSession) flush() error {
	sess.closeOutbound()
	select {
	case <-sess.senderDone:
		// Drained. sendErr is the sender's terminal status: nil if every queued frame went out, or
		// the Send error that stopped it early (leaving the wave's terminal TopicsLive/CatchupComplete
		// unsent) — surfaced rather than returned as a false OK.
		return sess.sendErr
	case <-sess.ctx.Done():
		// ctx fired — but the sender may have drained in the same instant (select is pseudorandom
		// when several cases are ready). Prefer a completed drain so a clean graceful close is never
		// reported as a spurious cancellation. Otherwise the stream was torn down before the drain
		// finished, so the wave's terminal TopicsLive/CatchupComplete may be unsent — surface that.
		if sess.senderDrained() {
			return sess.sendErr
		}
		return connect.NewError(
			connect.CodeCanceled,
			fmt.Errorf("flush interrupted before drain completed: %w", sess.ctx.Err()),
		)
	case <-time.After(sess.keepAlive):
		// Same priority: a drain that finished exactly as the deadline tripped is a success, not a
		// timeout.
		if sess.senderDrained() {
			return sess.sendErr
		}
		return connect.NewError(
			connect.CodeDeadlineExceeded,
			errors.New("flush timed out waiting for sender to drain"),
		)
	}
}

// senderDrained reports (without blocking) whether the sender goroutine has exited. After it returns
// true, sendErr is safe to read — the sender writes sendErr before closing senderDone.
func (sess *subscribeSession) senderDrained() bool {
	select {
	case <-sess.senderDone:
		return true
	default:
		return false
	}
}

// send hands one frame to the sender. It never blocks the writer indefinitely: if the sender is
// wedged on a non-reading client and the buffer fills, it fails the stream after the keepalive. It
// runs only on the writer goroutine, so the reused sendTimer needs no synchronization.
func (sess *subscribeSession) send(resp *message_api.SubscribeResponse) error {
	// A reused per-session timer rather than time.After per call (which would leak a timer for the
	// whole keepalive interval on every send — costly under high throughput / frame-splitting).
	if sess.sendTimer == nil {
		sess.sendTimer = time.NewTimer(sess.keepAlive)
	} else {
		sess.sendTimer.Reset(sess.keepAlive)
	}
	defer stopTimer(sess.sendTimer)

	select {
	case sess.outbound <- resp:
		// A frame was actually enqueued for delivery; defer the next idle Ping a full interval so
		// the ping/pong cadence tracks real outbound traffic.
		if sess.pingTicker != nil {
			sess.pingTicker.Reset(sess.keepAlive)
		}
		return nil
	case <-sess.senderDone:
		// The sender died (Send error) — stop feeding it; surface the error it recorded.
		return sess.sendErr
	case <-sess.ctx.Done():
		return nil
	case <-sess.svc.ctx.Done():
		return connect.NewError(connect.CodeUnavailable, errors.New("service is shutting down"))
	case <-sess.sendTimer.C:
		return connect.NewError(
			connect.CodeUnavailable,
			errors.New("send stalled; client not reading"),
		)
	}
}

// sendEnvelopes delivers envelopes split into frames each under maxSubscribeFrameBytes, so a large
// catch-up page or flushed pending buffer never goes out as one oversized (stream-aborting) frame.
func (sess *subscribeSession) sendEnvelopes(envs []*envelopesProto.OriginatorEnvelope) error {
	if len(envs) == 0 {
		return nil
	}
	limit := maxSubscribeFrameBytes
	if sess.maxFrameBytes > 0 {
		limit = sess.maxFrameBytes
	}
	var frame []*envelopesProto.OriginatorEnvelope
	frameBytes := 0
	flush := func() error {
		if len(frame) == 0 {
			return nil
		}
		if err := sess.send(newSubscribeEnvelopes(frame)); err != nil {
			return err
		}
		frame = nil // do NOT reuse the backing array; the sent frame still references it
		frameBytes = 0
		return nil
	}
	for _, env := range envs {
		size := proto.Size(env)
		// An envelope larger than `limit` on its own is NOT dropped: limit is a soft batching
		// target (2 MiB), an order of magnitude under the transport's hard cap (GRPCPayloadLimit,
		// 25 MiB). Such an envelope simply flushes the current frame and then goes out alone — and
		// it always fits, because it was publishable under that same 25 MiB cap. Skipping it (as the
		// old batchAndSendEnvelopes did) would silently lose a valid, deliverable message.
		if len(frame) > 0 && frameBytes+size > limit {
			if err := flush(); err != nil {
				return err
			}
		}
		frame = append(frame, env)
		frameBytes += size
	}
	return flush()
}

// routeLive splits a live batch: envelopes for a topic still catching up are buffered (gated) so
// they cannot overtake that topic's history; the rest are deduped against the live cursor and sent.
func (sess *subscribeSession) routeLive(batch []*envelopes.OriginatorEnvelope) error {
	var toSend []*envelopes.OriginatorEnvelope
	for _, env := range batch {
		ts := sess.topics[string(env.TargetTopic().Bytes())]
		if ts != nil && ts.phase == topicGated {
			if err := sess.bufferLive(ts, env); err != nil {
				return err
			}
			continue
		}
		// Live, or not (or no longer) ours: advanceLive dedups the live topics and drops the rest.
		toSend = append(toSend, env)
	}
	return sess.sendEnvelopes(sess.advanceLive(toSend))
}

// bufferLive holds a live envelope for a gated topic until its wave opens the gate, enforcing the
// global pending-bytes budget. Only valid while ts.phase == topicGated.
func (sess *subscribeSession) bufferLive(ts *topicState, env *envelopes.OriginatorEnvelope) error {
	ts.pending = append(ts.pending, env)
	sess.pendingBytes += proto.Size(env.Proto())
	if sess.pendingBytes > maxSubscribePendingBytes {
		return connect.NewError(
			connect.CodeResourceExhausted,
			errors.New("pending buffer exceeded while catching up"),
		)
	}
	return nil
}

// advanceLive dedups envs against their topics' live cursors, advancing each cursor in place, and
// returns the proto envelopes ready to send. An envelope for a topic that is not (or no longer) live
// is dropped — the per-topic analogue of advanceTopicCursors, reading the cursor from topicState.
func (sess *subscribeSession) advanceLive(
	envs []*envelopes.OriginatorEnvelope,
) []*envelopesProto.OriginatorEnvelope {
	result := make([]*envelopesProto.OriginatorEnvelope, 0, len(envs))
	for _, env := range envs {
		ts := sess.topics[string(env.TargetTopic().Bytes())]
		if ts == nil {
			sess.logger.Warn(
				"received envelope for unsubscribed topic",
				zap.Binary("topic", env.TargetTopic().Bytes()),
			)
			continue
		}
		origID := uint32(env.OriginatorNodeID())
		seqID := env.OriginatorSequenceID()
		if last, seen := ts.cursor[origID]; seen && last >= seqID {
			continue
		}
		ts.cursor[origID] = seqID
		result = append(result, env.Proto())
	}
	return result
}

// handleCatchUp processes one fetcher batch. Returns done=true when the writer should close (a
// half-close drain finished its last wave).
func (sess *subscribeSession) handleCatchUp(b catchUpBatch) (bool, error) {
	if b.err != nil {
		// A fetch error: fail so the client reconnects from its cursors, rather than emit a
		// misleading CatchupComplete over a history gap.
		return false, connect.NewError(
			connect.CodeUnavailable,
			fmt.Errorf("catch-up failed: %w", b.err),
		)
	}
	w := sess.waves[b.wave]
	if w == nil {
		return false, nil // wave already torn down (e.g. all its topics were removed)
	}

	// Deliver this page's history. A history_only wave dedups against its own throwaway cursors; a
	// live wave first drops pages for topics it no longer owns (removed, or reset under a newer
	// wave) — so it cannot advance a reset topic's live cursor and skip history the newer wave owes
	// — then dedups the rest against each topic's live cursor.
	var toSend []*envelopesProto.OriginatorEnvelope
	if w.historyOnly {
		toSend = advanceTopicCursors(w.cursors, b.envs, sess.logger)
	} else {
		toSend = sess.advanceLive(sess.envsOwnedByWave(b.envs, b.wave))
	}
	if len(toSend) > 0 {
		if err := sess.sendEnvelopes(toSend); err != nil {
			return false, err
		}
	}

	if !b.done {
		return false, nil
	}

	// Wave complete: open the gate for each live topic this wave still owns (flushing its buffered
	// live, deduped against the now-advanced cursor) and collect the topics to announce; then
	// CatchupComplete. flushAndGoLive is a no-op for a topic removed or reset under a newer wave, so
	// a stale wave never opens the newer wave's gate or flushes its buffer out of order.
	wire := make([][]byte, 0, len(w.topics))
	for _, t := range w.topics {
		if w.historyOnly {
			wire = append(wire, t.wire)
			continue
		}
		announced, err := sess.flushAndGoLive(t.cursorKey, b.wave)
		if err != nil {
			return false, err
		}
		if announced {
			wire = append(wire, t.wire)
		}
	}
	if len(wire) > 0 {
		if err := sess.send(newSubscribeTopicsLive(wire)); err != nil {
			return false, err
		}
	}
	if err := sess.send(newSubscribeCatchupComplete(w.mutateID)); err != nil {
		return false, err
	}
	delete(sess.waves, b.wave)
	return sess.halfClosed && len(sess.waves) == 0, nil
}

// flushAndGoLive completes a gated topic owned by `wave`: it flushes the live envelopes buffered
// during catch-up (deduped against the now-advanced live cursor) and transitions it to live,
// returning true once announced. It is a no-op returning false if the topic is gone or now owned by
// a newer wave (a reset), so a stale wave never opens the newer wave's gate or replays its buffer.
func (sess *subscribeSession) flushAndGoLive(cursorKey string, wave int) (bool, error) {
	ts := sess.topics[cursorKey]
	if ts == nil || ts.phase != topicGated || ts.wave != wave {
		return false, nil
	}
	if len(ts.pending) > 0 {
		for _, e := range ts.pending {
			sess.pendingBytes -= proto.Size(e.Proto())
		}
		buf := ts.pending
		ts.pending = nil
		if err := sess.sendEnvelopes(sess.advanceLive(buf)); err != nil {
			return false, err
		}
	}
	ts.phase = topicLive
	return true, nil
}

// handleRequest dispatches one client frame.
func (sess *subscribeSession) handleRequest(req *message_api.SubscribeRequest) error {
	v1 := req.GetV1()
	if v1 == nil {
		// Unrecognized version arm: fail rather than silently ignore, so a forward-version
		// client is not left waiting on a response (XIP-83 req 8).
		return connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("unrecognized SubscribeRequest version"),
		)
	}
	switch {
	case v1.GetPing() != nil:
		return sess.send(newSubscribePong(v1.GetPing().GetNonce()))
	case v1.GetPong() != nil:
		if v1.GetPong().GetNonce() == sess.pingNonce {
			sess.awaitingPong = false
		}
		return nil
	case v1.GetMutate() != nil:
		return sess.handleMutate(v1.GetMutate())
	default:
		return nil
	}
}

// handleMutate applies a Mutate atomically: it FULLY validates the frame (parses every topic,
// dedups adds, enforces the history_only and active-topic-cap rules) before touching any session
// state, so a malformed or over-cap Mutate cannot leave a half-applied subscription. It then
// removes first (so a topic in both removes and adds is reset), then adds with a catch-up wave
// fetched off-writer. Adds register live BEFORE catch-up starts (gate-before-fetch), so no message
// published during catch-up is missed.
func (sess *subscribeSession) handleMutate(m *message_api.SubscribeRequest_V1_Mutate) error {
	historyOnly := m.GetHistoryOnly()

	// ---- Validate (no state mutation): any failure here returns before a single change. ----

	// Parse removes up front so a malformed remove fails the whole Mutate, and so the add cap and
	// history_only checks below can account for topics this Mutate will drop.
	removes := make([]*topic.Topic, 0, len(m.GetRemoves()))
	removedKeys := make(map[string]struct{}, len(m.GetRemoves()))
	for _, w := range m.GetRemoves() {
		parsed, err := topic.ParseTopic(w)
		if err != nil {
			return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("remove: %w", err))
		}
		removes = append(removes, parsed)
		removedKeys[string(parsed.Bytes())] = struct{}{}
	}

	// Dedup adds (first cursor wins); reject malformed topics; enforce one cursor floor / one
	// in-flight catch-up per topic. A re-add of a topic that is NOT being reset by a remove in THIS
	// Mutate is special: re-gating a live topic would reset its cursor and re-deliver, so the only
	// way to replay/reset a topic is remove+re-add (which clears its floor first).
	type addReq struct {
		t        subscribeTopic
		provided db.VectorClock
	}
	order := make([]string, 0, len(m.GetAdds()))
	byKey := make(map[string]*addReq, len(m.GetAdds()))
	for _, a := range m.GetAdds() {
		parsed, err := topic.ParseTopic(a.GetTopic())
		if err != nil {
			return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("add: %w", err))
		}
		cursorKey := string(parsed.Bytes())
		if _, dup := byKey[cursorKey]; dup {
			continue
		}
		_, beingRemoved := removedKeys[cursorKey]
		if !beingRemoved {
			if _, exists := sess.topics[cursorKey]; exists {
				if historyOnly {
					// history_only needs its own cursor floor; it can't coexist with a live sub.
					return connect.NewError(
						connect.CodeInvalidArgument,
						errors.New(
							"history_only add targets a topic already subscribed on this stream",
						),
					)
				}
				// Plain re-add of an already-live topic is a no-op (idempotent): do not re-gate,
				// reset the cursor, or start a redundant wave. Replay requires remove+re-add.
				continue
			}
		}
		// A topic with an in-flight history_only catch-up can't take a second overlapping catch-up
		// (live or history_only) — even via remove+re-add: removeTopic does NOT cancel a history_only
		// wave (it only clears live state), so both waves would paginate and deliver T's history.
		if sess.hasInflightHistoryOnly(cursorKey) {
			return connect.NewError(
				connect.CodeInvalidArgument,
				errors.New("add targets a topic with an in-flight history_only catch-up"),
			)
		}
		// Validate each cursor entry against the signed DB column types. An out-of-range value would
		// be silently dropped by the catch-up query (SetPerTopicCursors) AND stored verbatim in the
		// sparse live cursor, where it would mark every real envelope from that originator as already
		// seen — permanently killing the topic on this stream. Reject instead.
		provided := make(db.VectorClock)
		for nodeID, seqID := range a.GetLastSeen().GetNodeIdToSequenceId() {
			if nodeID > math.MaxInt32 || seqID > uint64(math.MaxInt64) {
				return connect.NewError(
					connect.CodeInvalidArgument,
					fmt.Errorf(
						"cursor entry out of range (originator %d, sequence %d)",
						nodeID,
						seqID,
					),
				)
			}
			provided[nodeID] = seqID
		}
		byKey[cursorKey] = &addReq{
			t: subscribeTopic{
				wire:      a.GetTopic(),
				cursorKey: cursorKey,
				listenKey: parsed.String(),
			},
			provided: provided,
		}
		order = append(order, cursorKey)
	}

	// Cap the live set against the PROJECTED post-Mutate size: |(live \ removes) ∪ adds|. A herald
	// may hold a lot on one stream, but not unbounded.
	if !historyOnly {
		delta := 0
		for k := range removedKeys {
			if _, exists := sess.topics[k]; exists {
				delta-- // a currently-subscribed topic this Mutate drops
			}
		}
		for _, k := range order {
			_, existsNow := sess.topics[k]
			_, dropped := removedKeys[k]
			// Net-new entry unless it is already present after the removes (subscribed and kept).
			if !existsNow || dropped {
				delta++
			}
		}
		if len(sess.topics)+delta > maxActiveSubscribeTopics {
			return connect.NewError(
				connect.CodeResourceExhausted,
				fmt.Errorf("active topic limit %d exceeded", maxActiveSubscribeTopics),
			)
		}
	}

	// Bound concurrent in-flight catch-up waves (each a fetcher goroutine + DB queries). Only a
	// Mutate that would start a new wave is subject to it; removes (which never cancel an in-flight
	// wave) and no-op re-adds are exempt because they produce no order entries.
	if len(order) > 0 && len(sess.waves) >= maxInflightSubscribeWaves {
		return connect.NewError(
			connect.CodeResourceExhausted,
			fmt.Errorf("in-flight catch-up limit %d exceeded", maxInflightSubscribeWaves),
		)
	}

	// ---- Apply (validation passed; safe to mutate state). Removes first, then adds. ----
	for _, parsed := range removes {
		sess.removeTopic(parsed)
	}

	if len(order) == 0 {
		// Removes-only (or empty) Mutate: no catch-up, but confirm it applied so a client that
		// subscribed to nothing still learns the mutate took effect.
		if err := sess.send(newSubscribeCatchupComplete(m.GetMutateId())); err != nil {
			return err
		}
		if sess.halfClosed && len(sess.waves) == 0 {
			return sess.flush()
		}
		return nil
	}

	wave := &subscribeWave{mutateID: m.GetMutateId(), historyOnly: historyOnly}
	if historyOnly {
		wave.cursors = make(db.TopicCursors, len(order))
	}
	// providedCursors is the SPARSE per-topic client cursor; the fetcher fills it (with the full
	// originator set) off the writer goroutine. The persisted live cursor stays sparse (provided
	// only, grown as originators are actually seen) to bound memory at the 1M ceiling.
	providedCursors := make(db.TopicCursors, len(order))
	cursorKeys := make([]string, 0, len(order))
	for _, k := range order {
		a := byKey[k]
		wave.topics = append(wave.topics, a.t)
		cursorKeys = append(cursorKeys, k)
		providedCursors[k] = cloneVectorClock(a.provided)

		if historyOnly {
			wave.cursors[k] = cloneVectorClock(a.provided)
			continue
		}
		// Live add: gate it under this wave (seeds the sparse live cursor and registers with the
		// worker), all before the fetch starts.
		sess.gateTopic(a.t, sess.nextWave, a.provided)
	}

	sess.waves[sess.nextWave] = wave
	go sess.svc.runSubscribeCatchUp(
		sess.ctx,
		sess.nextWave,
		providedCursors,
		cursorKeys,
		sess.catchUpCh,
		sess.logger,
	)
	sess.nextWave++
	return nil
}

// hasInflightHistoryOnly reports whether the topic currently has an in-flight history_only catch-up
// wave (its cursor lives on the wave's own subscribeWave.cursors, present only while that wave is
// still in sess.waves). Used to reject a second overlapping catch-up for the same topic, which would
// double-deliver its history.
func (sess *subscribeSession) hasInflightHistoryOnly(cursorKey string) bool {
	for _, w := range sess.waves {
		if w.historyOnly {
			if _, ok := w.cursors[cursorKey]; ok {
				return true
			}
		}
	}
	return false
}

// gateTopic registers a live add: it creates the topic in the gated phase owned by `wave`, seeds its
// sparse live cursor from the client-provided floor, and registers it with the worker — all before
// the catch-up fetch starts, so no envelope published during catch-up is missed (gate-before-fetch).
func (sess *subscribeSession) gateTopic(t subscribeTopic, wave int, provided db.VectorClock) {
	sess.topics[t.cursorKey] = &topicState{
		subscribeTopic: t,
		phase:          topicGated,
		wave:           wave,
		cursor:         cloneVectorClock(provided),
	}
	sess.sub.addTopics(map[string]struct{}{t.listenKey: {}})
}

func (sess *subscribeSession) removeTopic(parsed *topic.Topic) {
	sess.sub.removeTopics(map[string]struct{}{parsed.String(): {}})
	sess.removeTopicState(string(parsed.Bytes()))
}

// removeTopicState drops a live topic's state and refunds its buffered pending bytes. It does NOT
// touch any in-flight wave: a wave only ever acts on topics it still owns (flushAndGoLive /
// envsOwnedByWave both re-check ownership), so a removed topic's pages and completion are ignored.
func (sess *subscribeSession) removeTopicState(cursorKey string) {
	ts := sess.topics[cursorKey]
	if ts == nil {
		return
	}
	for _, e := range ts.pending {
		sess.pendingBytes -= proto.Size(e.Proto())
	}
	delete(sess.topics, cursorKey)
}

// envsOwnedByWave keeps only envelopes whose topic is still being caught up by this wave. A topic
// removed (or removed and re-added under a newer wave) since the fetch began is no longer this
// wave's to deliver; dropping its pages keeps a stale wave from advancing the live cursor of a
// reset topic (which would skip history the newer wave still owes the client).
func (sess *subscribeSession) envsOwnedByWave(
	envs []*envelopes.OriginatorEnvelope,
	wave int,
) []*envelopes.OriginatorEnvelope {
	owned := func(env *envelopes.OriginatorEnvelope) bool {
		ts := sess.topics[string(env.TargetTopic().Bytes())]
		return ts != nil && ts.phase == topicGated && ts.wave == wave
	}
	// Fast path: every envelope belongs to a topic this wave still owns (the common case — no
	// concurrent remove/reset), so the page passes through without reallocation.
	allOwned := true
	for _, env := range envs {
		if !owned(env) {
			allOwned = false
			break
		}
	}
	if allOwned {
		return envs
	}
	out := make([]*envelopes.OriginatorEnvelope, 0, len(envs))
	for _, env := range envs {
		if owned(env) {
			out = append(out, env)
		}
	}
	return out
}

// runSubscribeCatchUp paginates history for a wave's topics (off the writer goroutine) and hands
// raw pages back over catchUpCh, ending with a done marker. It resolves the originator set and fills
// providedCursors into its own query cursors here — that originator lookup is a DB round-trip on a
// cache miss, so it MUST stay off the writer goroutine (else a slow DB would stall liveness and
// live delivery and could false-reap a healthy stream). The writer owns the sparse live cursors.
// Every channel send is guarded by ctx so the fetcher cannot leak if the writer has torn down.
func (s *Service) runSubscribeCatchUp(
	ctx context.Context,
	wave int,
	providedCursors db.TopicCursors,
	cursorKeys []string,
	catchUpCh chan<- catchUpBatch,
	logger *zap.Logger,
) {
	emit := func(b catchUpBatch) bool {
		select {
		case catchUpCh <- b:
			return true
		case <-ctx.Done():
			return false
		}
	}

	// The originator set is snapshotted once. For a LIVE wave that is complete: the listener is
	// registered before this snapshot, so an originator that first publishes after the snapshot is
	// caught by the live gate (the worker dispatches by topic, independent of this list). A
	// history_only wave has no live gate, so a brand-new originator that first publishes during the
	// catch-up window is delivered on the client's next subscribe rather than this bounded sync —
	// an accepted eventual-consistency property of history_only (it is a periodic re-sync flow).
	knownOriginators, err := s.originatorList.GetOriginatorNodeIDs(ctx)
	if err != nil {
		emit(catchUpBatch{wave: wave, err: fmt.Errorf("could not get originator list: %w", err)})
		return
	}
	// queryCursors are FILLED (every originator from the provided/zero start) so catch-up covers all
	// originators; the fetcher owns and advances them for pagination.
	queryCursors := make(db.TopicCursors, len(providedCursors))
	for k, provided := range providedCursors {
		filled := cloneVectorClock(provided)
		db.FillMissingOriginators(filled, knownOriginators)
		queryCursors[k] = filled
	}

	for _, chunkKeys := range utils.ChunkSlice(cursorKeys, maxTopicsPerChunk) {
		rowsPerEntry := db.CalculateRowsPerEntry(len(chunkKeys), topicPageLimit)
		for {
			if ctx.Err() != nil {
				return
			}
			subCursors := make(db.TopicCursors, len(chunkKeys))
			for _, k := range chunkKeys {
				subCursors[k] = queryCursors[k]
			}
			rows, err := s.fetchTopicEnvelopesWithRetry(
				ctx,
				subCursors,
				topicPageLimit,
				rowsPerEntry,
			)
			if err != nil {
				emit(catchUpBatch{wave: wave, err: err})
				return
			}
			envs := unmarshalEnvelopes(rows, logger)
			// Advance the fetcher's own (filled) cursors from the RAW rows so pagination always
			// progresses even if some rows fail to unmarshal (otherwise a single bad row in a full
			// page re-fetches forever); the writer re-dedups the emitted envs against the live cursor.
			advanceCursorsFromRows(queryCursors, rows)
			if !emit(catchUpBatch{wave: wave, envs: envs}) {
				return
			}
			if int32(len(rows)) < rowsPerEntry {
				break
			}
		}
	}
	emit(catchUpBatch{wave: wave, done: true})
}

func cloneVectorClock(vc db.VectorClock) db.VectorClock {
	out := make(db.VectorClock, len(vc))
	maps.Copy(out, vc)
	return out
}

// stopTimer stops t and drains a pending fire if there is one, so a later Reset cannot observe a
// stale value. Safe on an already-stopped or already-fired-and-drained timer.
func stopTimer(t *time.Timer) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
}

// drainPendingRequests processes client frames already queued on ch without blocking. It is called
// right before the pong-deadline reap so a Pong that landed in the buffer just as the deadline fired
// (Go select picks randomly among ready cases) is still counted, instead of false-reaping a stream
// that did answer. ch may be nil (half-closed), in which case this is a no-op.
func (sess *subscribeSession) drainPendingRequests(ch <-chan *message_api.SubscribeRequest) error {
	for {
		select {
		case req, open := <-ch:
			if !open {
				// The reader closed the channel. Mirror the main loop's EOF handling: surface a
				// transport failure, else mark halfClosed for a clean EOF — so the pong-deadline
				// caller does not then reap a cleanly-finishing (or already-failed) stream with a
				// spurious DeadlineExceeded.
				if sess.recvErr != nil {
					return connect.NewError(
						connect.CodeUnavailable,
						fmt.Errorf("stream recv failed: %w", sess.recvErr),
					)
				}
				sess.halfClosed = true
				return nil
			}
			if err := sess.handleRequest(req); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

func wrapSubscribeV1(v1 *message_api.SubscribeResponse_V1) *message_api.SubscribeResponse {
	return &message_api.SubscribeResponse{Version: &message_api.SubscribeResponse_V1_{V1: v1}}
}

func newSubscribeStarted(keepaliveIntervalMs uint32) *message_api.SubscribeResponse {
	return wrapSubscribeV1(&message_api.SubscribeResponse_V1{
		Response: &message_api.SubscribeResponse_V1_Started_{
			Started: &message_api.SubscribeResponse_V1_Started{
				KeepaliveIntervalMs: keepaliveIntervalMs,
			},
		},
	})
}

func newSubscribeEnvelopes(
	envs []*envelopesProto.OriginatorEnvelope,
) *message_api.SubscribeResponse {
	return wrapSubscribeV1(&message_api.SubscribeResponse_V1{
		Response: &message_api.SubscribeResponse_V1_Envelopes_{
			Envelopes: &message_api.SubscribeResponse_V1_Envelopes{Envelopes: envs},
		},
	})
}

func newSubscribeTopicsLive(topics [][]byte) *message_api.SubscribeResponse {
	return wrapSubscribeV1(&message_api.SubscribeResponse_V1{
		Response: &message_api.SubscribeResponse_V1_TopicsLive_{
			TopicsLive: &message_api.SubscribeResponse_V1_TopicsLive{Topics: topics},
		},
	})
}

func newSubscribeCatchupComplete(mutateID uint64) *message_api.SubscribeResponse {
	return wrapSubscribeV1(&message_api.SubscribeResponse_V1{
		Response: &message_api.SubscribeResponse_V1_CatchupComplete_{
			CatchupComplete: &message_api.SubscribeResponse_V1_CatchupComplete{MutateId: mutateID},
		},
	})
}

func newSubscribePing(nonce uint64) *message_api.SubscribeResponse {
	return wrapSubscribeV1(&message_api.SubscribeResponse_V1{
		Response: &message_api.SubscribeResponse_V1_Ping{Ping: &message_api.Ping{Nonce: nonce}},
	})
}

func newSubscribePong(nonce uint64) *message_api.SubscribeResponse {
	return wrapSubscribeV1(&message_api.SubscribeResponse_V1{
		Response: &message_api.SubscribeResponse_V1_Pong{Pong: &message_api.Pong{Nonce: nonce}},
	})
}
