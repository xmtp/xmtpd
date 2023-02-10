package memsubs_test

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/node/subscribers"
	memsubs "github.com/xmtp/xmtpd/pkg/node/subscribers/mem"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func TestMemorySubscribers(t *testing.T) {
	log := test.NewLogger(t)

	subs := newTestSubscribers(t, memsubs.New(log, 100))
	defer subs.Close()

	sub1 := subs.subscribe(t, "topic")
	sub2 := subs.subscribe(t, "topic")
	otherSub := subs.subscribe(t, "other-topic")
	unsub := subs.subscribe(t, "topic")
	subs.unsubscribe(t, unsub)

	events := subs.publishMany(t, "topic", 2)
	otherEvents := subs.publishMany(t, "other-topic", 1)

	sub1.requireEventuallyCapturedEvents(t, events)
	sub2.requireEventuallyCapturedEvents(t, events)
	otherSub.requireEventuallyCapturedEvents(t, otherEvents)
	unsub.requireEventuallyCapturedEvents(t, nil)
}

type testSubscribers struct {
	subscribers.Manager

	ctx    context.Context
	cancel context.CancelFunc
}

func newTestSubscribers(t *testing.T, subs subscribers.Manager) *testSubscribers {
	t.Helper()
	s := &testSubscribers{
		Manager: subs,
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())
	return s
}

func (s *testSubscribers) Close() error {
	if s.cancel != nil {
		s.cancel()
	}
	return nil
}

func (s *testSubscribers) subscribe(t *testing.T, topic string) *testSubscriber {
	t.Helper()
	sub := &testSubscriber{
		topic: topic,
		ch:    s.Subscribe(s.ctx, topic),
	}
	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			case ev, ok := <-sub.ch:
				if !ok {
					return
				}
				sub.eventsLock.Lock()
				sub.events = append(sub.events, ev)
				sub.eventsLock.Unlock()
			}
		}
	}()
	return sub
}

func (s *testSubscribers) unsubscribe(t *testing.T, sub *testSubscriber) {
	t.Helper()
	s.Unsubscribe(sub.topic, sub.ch)
}

func (s *testSubscribers) publishMany(t *testing.T, topic string, count int) []*types.Event {
	t.Helper()
	events := make([]*types.Event, count)
	for i := 0; i < count; i++ {
		ev, err := types.NewEvent(&messagev1.Envelope{
			ContentTopic: topic,
			TimestampNs:  uint64(rand.Intn(100)),
			Message:      []byte("msg-" + test.RandomString(13)),
		}, nil)
		require.NoError(t, err)
		s.OnNewEvent(topic, ev)
		events[i] = ev
	}
	return events
}

type testSubscriber struct {
	topic      string
	ch         chan *types.Event
	events     []*types.Event
	eventsLock sync.RWMutex
}

func (s *testSubscriber) requireEventuallyCapturedEvents(t *testing.T, expected []*types.Event) {
	t.Helper()
	assert.Eventually(t, func() bool {
		s.eventsLock.RLock()
		defer s.eventsLock.RUnlock()
		return len(s.events) == len(expected) && len(s.ch) == 0
	}, time.Second, 10*time.Millisecond)
	require.ElementsMatch(t, expected, s.events)
}
