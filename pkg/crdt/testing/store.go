package crdttest

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/require"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	test "github.com/xmtp/xmtpd/pkg/testing"
	"github.com/xmtp/xmtpd/pkg/utils"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type ITestStore interface {
	crdt.Store

	Events() ([]*types.Event, error)
	Heads() []multihash.Multihash
}

type TestStoreMaker func(t *testing.T) *TestStore

type TestStore struct {
	ITestStore

	log *zap.Logger
	ctx context.Context
}

func NewTestStore(ctx context.Context, log *zap.Logger, store ITestStore) *TestStore {
	return &TestStore{
		ITestStore: store,

		log: log,
		ctx: ctx,
	}
}

func RunStoreEventTests(t *testing.T, storeMaker TestStoreMaker) {
	t.Run("add events", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)

		ev1 := s.addRandomEvent(t)
		s.requireEventsEqual(t, ev1)

		ev2 := s.addRandomEvent(t)
		s.requireEventsEqual(t, ev1, ev2)
	})

	t.Run("add existing event", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)

		ev1 := s.addRandomEvent(t)
		s.requireEventsEqual(t, ev1)

		s.addExistingEvent(t, ev1)
		s.requireEventsEqual(t, ev1)
	})

	t.Run("append events", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)

		ev1 := s.appendRandomEvent(t)
		require.NotNil(t, ev1.Cid)
		require.Nil(t, ev1.Links)
		s.requireEventsEqual(t, ev1)

		ev2 := s.appendRandomEvent(t)
		require.NotNil(t, ev2.Cid)
		require.Equal(t, []multihash.Multihash{ev1.Cid}, ev2.Links)
		s.requireEventsEqual(t, ev1, ev2)
	})

	t.Run("add remove heads", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)

		head := s.addRandomHead(t)
		s.requireEventsEqual(t, head)

		ev1 := s.newRandomEventWithHeads(t, []multihash.Multihash{head.Cid})
		require.Equal(t, []multihash.Multihash{head.Cid}, ev1.Links)
		s.addEvent(t, ev1)
		s.requireEventsEqual(t, head, ev1)

		ev2 := s.newRandomEvent(t)
		require.Nil(t, ev2.Links)
		s.addEvent(t, ev2)
		s.requireEventsEqual(t, head, ev1, ev2)
	})

	t.Run("add existing head", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)

		head := s.addRandomHead(t)
		s.requireEventsEqual(t, head)

		s.addExistingEvent(t, head)
		s.requireEventsEqual(t, head)
	})
}

func RunStoreQueryTests(t *testing.T, storeMaker TestStoreMaker) {
	t.Helper()

	t.Run("all sort default", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 20)
		test.RequireProtoEqual(t, toEnvelopes(events), res.Envelopes)
	})

	t.Run("all sort ascending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_ASCENDING,
			},
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 20)
		test.RequireProtoEqual(t, toEnvelopes(events), res.Envelopes)
	})

	t.Run("all sort descending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_DESCENDING,
			},
		})
		require.NoError(t, err)
		utils.Reverse(events)
		require.Len(t, res.Envelopes, 20)
		test.RequireProtoEqual(t, toEnvelopes(events), res.Envelopes)
	})

	t.Run("limit sort default", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			PagingInfo: &messagev1.PagingInfo{
				Limit: 5,
			},
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 5)
		test.RequireProtoEqual(t, toEnvelopes(events[:5]), res.Envelopes)
	})

	t.Run("limit sort ascending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_ASCENDING,
				Limit:     5,
			},
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 5)
		test.RequireProtoEqual(t, toEnvelopes(events[:5]), res.Envelopes)
	})

	t.Run("limit sort descending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_DESCENDING,
				Limit:     5,
			},
		})
		require.NoError(t, err)
		utils.Reverse(events)
		require.Len(t, res.Envelopes, 5)
		test.RequireProtoEqual(t, toEnvelopes(events[:5]), res.Envelopes)
	})

	t.Run("start time sort default", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 10,
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 11)
		test.RequireProtoEqual(t, toEnvelopes(events[9:]), res.Envelopes)
	})

	t.Run("end time sort default", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			EndTimeNs: 10,
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 10)
		test.RequireProtoEqual(t, toEnvelopes(events[:10]), res.Envelopes)
	})

	t.Run("time range sort default", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 5,
			EndTimeNs:   15,
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 11)
		test.RequireProtoEqual(t, toEnvelopes(events[4:15]), res.Envelopes)
	})

	t.Run("start time sort ascending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 10,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_ASCENDING,
			},
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 11)
		test.RequireProtoEqual(t, toEnvelopes(events[9:]), res.Envelopes)
	})

	t.Run("end time sort ascending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			EndTimeNs: 10,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_ASCENDING,
			},
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 10)
		test.RequireProtoEqual(t, toEnvelopes(events[:10]), res.Envelopes)
	})

	t.Run("time range sort ascending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 5,
			EndTimeNs:   15,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_ASCENDING,
			},
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 11)
		test.RequireProtoEqual(t, toEnvelopes(events[4:15]), res.Envelopes)
	})

	t.Run("start time sort descending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 10,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_DESCENDING,
			},
		})
		require.NoError(t, err)
		events = events[9:]
		utils.Reverse(events)
		require.Len(t, res.Envelopes, 11)
		test.RequireProtoEqual(t, toEnvelopes(events), res.Envelopes)
	})

	t.Run("end time sort descending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			EndTimeNs: 10,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_DESCENDING,
			},
		})
		require.NoError(t, err)
		events = events[:10]
		utils.Reverse(events)
		require.Len(t, res.Envelopes, 10)
		test.RequireProtoEqual(t, toEnvelopes(events), res.Envelopes)
	})

	t.Run("time range sort descending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 5,
			EndTimeNs:   15,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_DESCENDING,
			},
		})
		require.NoError(t, err)
		events = events[4:15]
		utils.Reverse(events)
		require.Len(t, res.Envelopes, 11)
		test.RequireProtoEqual(t, toEnvelopes(events), res.Envelopes)
	})

	t.Run("limit start time sort default", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 10,
			PagingInfo: &messagev1.PagingInfo{
				Limit: 3,
			},
		})
		require.NoError(t, err)
		events = events[9:]
		require.Len(t, res.Envelopes, 3)
		test.RequireProtoEqual(t, toEnvelopes(events[:3]), res.Envelopes)
	})

	t.Run("limit end time sort default", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			EndTimeNs: 10,
			PagingInfo: &messagev1.PagingInfo{
				Limit: 3,
			},
		})
		require.NoError(t, err)
		events = events[:10]
		require.Len(t, res.Envelopes, 3)
		test.RequireProtoEqual(t, toEnvelopes(events[:3]), res.Envelopes)
	})

	t.Run("limit time range sort default", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 5,
			EndTimeNs:   15,
			PagingInfo: &messagev1.PagingInfo{
				Limit: 3,
			},
		})
		require.NoError(t, err)
		events = events[4:15]
		require.Len(t, res.Envelopes, 3)
		test.RequireProtoEqual(t, toEnvelopes(events[:3]), res.Envelopes)
	})

	t.Run("limit start time sort ascending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 10,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_ASCENDING,
				Limit:     3,
			},
		})
		require.NoError(t, err)
		events = events[9:]
		require.Len(t, res.Envelopes, 3)
		test.RequireProtoEqual(t, toEnvelopes(events[:3]), res.Envelopes)
	})

	t.Run("limit end time sort ascending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			EndTimeNs: 10,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_ASCENDING,
				Limit:     3,
			},
		})
		require.NoError(t, err)
		events = events[:10]
		require.Len(t, res.Envelopes, 3)
		test.RequireProtoEqual(t, toEnvelopes(events[:3]), res.Envelopes)
	})

	t.Run("limit time range sort ascending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 5,
			EndTimeNs:   15,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_ASCENDING,
				Limit:     3,
			},
		})
		require.NoError(t, err)
		events = events[4:15]
		require.Len(t, res.Envelopes, 3)
		test.RequireProtoEqual(t, toEnvelopes(events[:3]), res.Envelopes)
	})

	t.Run("limit start time sort descending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 10,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_DESCENDING,
				Limit:     3,
			},
		})
		require.NoError(t, err)
		events = events[9:]
		utils.Reverse(events)
		require.Len(t, res.Envelopes, 3)
		test.RequireProtoEqual(t, toEnvelopes(events[:3]), res.Envelopes)
	})

	t.Run("limit end time sort descending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			EndTimeNs: 10,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_DESCENDING,
				Limit:     3,
			},
		})
		require.NoError(t, err)
		events = events[:10]
		utils.Reverse(events)
		require.Len(t, res.Envelopes, 3)
		test.RequireProtoEqual(t, toEnvelopes(events[:3]), res.Envelopes)
	})

	t.Run("limit time range sort descending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 5,
			EndTimeNs:   15,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_DESCENDING,
				Limit:     3,
			},
		})
		require.NoError(t, err)
		events = events[5:15]
		utils.Reverse(events)
		require.Len(t, res.Envelopes, 3)
		test.RequireProtoEqual(t, toEnvelopes(events[:3]), res.Envelopes)
	})

	t.Run("cursor sort default", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			PagingInfo: &messagev1.PagingInfo{
				Cursor: &messagev1.Cursor{
					Cursor: &messagev1.Cursor_Index{
						Index: &messagev1.IndexCursor{
							SenderTimeNs: 10,
							Digest:       events[9].Cid,
						},
					},
				},
			},
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 10)
		test.RequireProtoEqual(t, toEnvelopes(events[10:]), res.Envelopes)
	})

	t.Run("cursor sort ascending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_ASCENDING,
				Cursor: &messagev1.Cursor{
					Cursor: &messagev1.Cursor_Index{
						Index: &messagev1.IndexCursor{
							SenderTimeNs: 10,
							Digest:       events[9].Cid,
						},
					},
				},
			},
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 10)
		test.RequireProtoEqual(t, toEnvelopes(events[10:]), res.Envelopes)
	})

	t.Run("cursor sort descending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		events := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_DESCENDING,
				Cursor: &messagev1.Cursor{
					Cursor: &messagev1.Cursor_Index{
						Index: &messagev1.IndexCursor{
							SenderTimeNs: 10,
							Digest:       events[9].Cid,
						},
					},
				},
			},
		})
		require.NoError(t, err)
		events = events[:9]
		utils.Reverse(events)
		require.Len(t, res.Envelopes, 9)
		test.RequireProtoEqual(t, toEnvelopes(events), res.Envelopes)
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
	added, err := s.InsertHead(s.ctx, head)
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

func (s *TestStore) appendRandomEvent(t *testing.T) *types.Event {
	t.Helper()
	return s.appendEvent(t, newRandomEnvelope(t))
}

func (s *TestStore) appendEvent(t *testing.T, env *messagev1.Envelope) *types.Event {
	t.Helper()
	ev, err := s.AppendEvent(s.ctx, env)
	require.NoError(t, err)
	return ev
}

func (s *TestStore) newRandomEvent(t *testing.T) *types.Event {
	t.Helper()
	return s.newRandomEventWithHeads(t, nil)
}

func (s *TestStore) newRandomEventWithHeads(t *testing.T, heads []multihash.Multihash) *types.Event {
	t.Helper()
	ev, err := types.NewEvent(newRandomEnvelope(t), heads)
	require.NoError(t, err)
	return ev
}

func (s *TestStore) requireEventsEqual(t *testing.T, expected ...*types.Event) {
	t.Helper()
	events, err := s.Events()
	require.NoError(t, err)
	require.ElementsMatch(t, expected, events)
}

func (s *TestStore) seed(t *testing.T, count int) []*types.Event {
	t.Helper()
	ctx := context.Background()
	events := make([]*types.Event, count)
	for i := 0; i < 20; i++ {
		ev, err := s.AppendEvent(ctx, &messagev1.Envelope{
			ContentTopic: "topic",
			TimestampNs:  uint64(i + 1),
			Message:      []byte(fmt.Sprintf("msg-%d", i+1)),
		})
		require.NoError(t, err)
		events[i] = ev
	}
	return events
}

func (s *TestStore) query(t *testing.T, req *messagev1.QueryRequest) (*messagev1.QueryResponse, error) {
	t.Helper()
	ctx := context.Background()
	return s.Query(ctx, req)
}

func newRandomEnvelope(t *testing.T) *messagev1.Envelope {
	return &messagev1.Envelope{
		ContentTopic: "topic-" + test.RandomStringLower(5),
		TimestampNs:  uint64(rand.Intn(100)),
		Message:      []byte("msg-" + test.RandomString(13)),
	}
}

func toEnvelopes(events []*types.Event) []*messagev1.Envelope {
	envs := make([]*messagev1.Envelope, len(events))
	for i, ev := range events {
		envs[i] = ev.Envelope
	}
	return envs
}
