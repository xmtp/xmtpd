package bolt

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/multiformats/go-multihash"
	mh "github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/require"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	v1 "github.com/xmtp/proto/v3/go/message_api/v1"
	bolt "go.etcd.io/bbolt"

	crdttest "github.com/xmtp/xmtpd/pkg/crdt/testing"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	test "github.com/xmtp/xmtpd/pkg/testing"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type testNodeStore struct {
	*NodeStore
}

type testStore struct {
	ns *testNodeStore
	*Store
}

func (s *testStore) seed(count int) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		topic := tx.Bucket(s.name)
		var prev []multihash.Multihash
		for i := 0; i < count; i++ {
			env := &messagev1.Envelope{
				ContentTopic: string(s.name),
				TimestampNs:  uint64(i + 1),
				Message:      []byte(fmt.Sprintf("msg-%d", i+1)),
			}
			ev, err := types.NewEvent(env, prev)
			if err != nil {
				return err
			}
			prev = []multihash.Multihash{ev.Cid}
			if _, err = addEvent(topic, ev); err != nil {
				return err
			}
		}
		heads := topic.Bucket(HeadsBucket)
		return heads.Put(prev[0], nil)
	})
}

func (s *testStore) Close() error {
	return s.ns.Close()
}

func newTestNodeStore(t testing.TB, log *zap.Logger) *testNodeStore {
	path := filepath.Join(t.TempDir(), "testdb.bolt")
	store, err := NewNodeStore(path, log)
	require.NoError(t, err)
	return &testNodeStore{store}
}

func newTestStore(t testing.TB, topic string, ns *testNodeStore) *testStore {
	store, err := ns.NewTopic(topic)
	require.NoError(t, err)
	return &testStore{ns: ns, Store: store.(*Store)}
}

func TestStore(t *testing.T) {
	ctx := context.Background()
	log := test.NewLogger(t)

	crdttest.RunStoreEventTests(t, "topic", func(t *testing.T) *crdttest.TestStore {
		ns := newTestNodeStore(t, log)
		topic := newTestStore(t, "topic", ns)
		return crdttest.NewTestStore(ctx, log, topic)
	})
}

func TestQuery(t *testing.T) {
	ctx := context.Background()
	log := test.NewLogger(t)

	crdttest.RunStoreQueryTests(t, "topic", func(t *testing.T) *crdttest.TestStore {
		ns := newTestNodeStore(t, log)
		topic := newTestStore(t, "topic", ns)
		return crdttest.NewTestStore(ctx, log, topic)
	})
}

func Test_Basic(t *testing.T) {
	const t0 = "t0"
	ctx := context.Background()
	log := test.NewLogger(t)
	ns := newTestNodeStore(t, log)
	defer ns.Close()
	topic := newTestStore(t, t0, ns)
	defer topic.Close()

	msg := []byte("Buh")
	ev, err := topic.AppendEvent(ctx, &v1.Envelope{ContentTopic: t0, TimestampNs: 42, Message: msg})
	require.NoError(t, err)
	evdb, err := topic.GetEvents(ctx, ev.Cid)
	require.NoError(t, err)
	require.Len(t, evdb, 1)
	require.Equal(t, msg, evdb[0].Message)
	require.Equal(t, t0, evdb[0].ContentTopic)
	require.Equal(t, uint64(42), evdb[0].TimestampNs)

	ev2, err := types.NewEvent(&v1.Envelope{Message: msg}, []mh.Multihash{ev.Cid})
	require.NoError(t, err)
	added, err := topic.InsertHead(ctx, ev2)
	require.NoError(t, err)
	require.True(t, added)
	ev2db, err := topic.GetEvents(ctx, ev2.Cid)
	require.NoError(t, err)
	require.Len(t, ev2db, 1)
	require.Equal(t, msg, ev2db[0].Message)
	require.Len(t, ev2db[0].Links, 1)
	require.Equal(t, ev.Cid, ev2db[0].Links[0])

	ev3, err := topic.AppendEvent(ctx, &v1.Envelope{Message: msg})
	require.NoError(t, err)
	ev3db, err := topic.GetEvents(ctx, ev3.Cid)
	require.NoError(t, err)
	require.Len(t, ev3db, 1)
	require.Equal(t, msg, ev3db[0].Message)
	require.EqualValues(t, []mh.Multihash{ev.Cid, ev2.Cid}, ev3db[0].Links)

	resp, err := topic.Query(ctx, &v1.QueryRequest{ContentTopics: []string{t0}})
	require.NoError(t, err)
	require.Len(t, resp.Envelopes, 3)
}
