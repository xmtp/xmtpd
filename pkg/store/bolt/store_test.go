package bolt_test

import (
	"path/filepath"
	"testing"

	mh "github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/require"
	v1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"go.etcd.io/bbolt"

	"github.com/xmtp/xmtpd/pkg/context"
	crdttest "github.com/xmtp/xmtpd/pkg/crdt/testing"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/store/bolt"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

type testNodeStore struct {
	*bolt.NodeStore
	db *bbolt.DB
}

type testStore struct {
	ns *testNodeStore
	*bolt.Store
}

func (s *testStore) Close() error {
	return s.ns.Close()
}

func newTestNodeStore(t testing.TB, ctx context.Context) *testNodeStore {
	fn := "test_" + test.RandomStringLower(13) + ".bolt"
	opts := &bolt.Options{
		DataPath: filepath.Join(t.TempDir(), fn),
	}
	db, err := bolt.NewDB(opts)
	require.NoError(t, err)
	store, err := bolt.NewNodeStore(ctx, db, opts)
	require.NoError(t, err)
	return &testNodeStore{store, db}
}

func newTestStore(t testing.TB, topic string, ns *testNodeStore) *testStore {
	store, err := ns.NewTopic(topic)
	require.NoError(t, err)
	return &testStore{ns: ns, Store: store.(*bolt.Store)}
}

func TestEvents(t *testing.T) {
	ctx := test.NewContext(t)

	crdttest.RunStoreEventTests(t, "topic", func(t *testing.T) *crdttest.TestStore {
		ns := newTestNodeStore(t, ctx)
		topic := newTestStore(t, "topic", ns)
		return crdttest.NewTestStore(ctx, topic)
	})
}

func TestQuery(t *testing.T) {
	ctx := test.NewContext(t)

	crdttest.RunStoreQueryTests(t, "topic", func(t *testing.T) *crdttest.TestStore {
		ns := newTestNodeStore(t, ctx)
		topic := newTestStore(t, "topic", ns)
		return crdttest.NewTestStore(ctx, topic)
	})
}

func Test_Basic(t *testing.T) {
	const t0 = "t0"
	ctx := test.NewContext(t)
	ns := newTestNodeStore(t, ctx)
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
