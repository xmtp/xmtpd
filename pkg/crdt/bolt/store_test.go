package bolt

import (
	"context"
	"os"
	"testing"

	mh "github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/require"
	v1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/crdt/tests"
	helpers "github.com/xmtp/xmtpd/pkg/testing"
)

func Test_Basic(t *testing.T) {
	withTempStore(t, func(topic crdt.TopicStore) {
		msg := []byte("Buh")
		ev, err := topic.NewEvent(&v1.Envelope{ContentTopic: t0, TimestampNs: 42, Message: msg})
		require.NoError(t, err)
		t.Logf("store event: %s", ev.Cid.String())
		evdb, err := topic.Get(ev.Cid)
		require.NoError(t, err)
		require.NotNil(t, evdb)
		require.Equal(t, msg, evdb.Message)
		require.Equal(t, t0, evdb.ContentTopic)
		require.Equal(t, uint64(42), evdb.TimestampNs)

		ev2, err := crdt.NewEvent(&v1.Envelope{Message: msg}, []mh.Multihash{ev.Cid})
		require.NoError(t, err)
		added, err := topic.AddHead(ev2)
		require.NoError(t, err)
		require.True(t, added)
		t.Logf("store event: %s", ev2.Cid.String())
		ev2db, err := topic.Get(ev2.Cid)
		require.NoError(t, err)
		require.NotNil(t, ev2db)
		require.Equal(t, msg, ev2db.Message)
		require.Len(t, ev2db.Links, 1)
		require.Equal(t, ev.Cid, ev2db.Links[0])

		ev3, err := topic.NewEvent(&v1.Envelope{Message: msg})
		require.NoError(t, err)
		t.Logf("store event: %s", ev3.Cid.String())
		ev3db, err := topic.Get(ev3.Cid)
		require.NoError(t, err)
		require.NotNil(t, ev3db)
		require.Equal(t, msg, ev3db.Message)
		require.EqualValues(t, []mh.Multihash{ev.Cid, ev2.Cid}, ev3db.Links)

		c, err := topic.Count()
		require.NoError(t, err)
		require.Equal(t, 3, c)

		envs, _, err := topic.Query(context.Background(), &v1.QueryRequest{ContentTopics: []string{t0}})
		require.NoError(t, err)
		require.Len(t, envs, 3)
	})

}

// helpers

const t0 = "t0"

func withTempStore(t *testing.T, do func(crdt.TopicStore)) {
	f, err := os.CreateTemp("", "crdt-test")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	t.Logf("temp db: %s", f.Name())
	store, err := NewStore(f.Name())
	require.NoError(t, err)
	log := helpers.NewLogger(t)
	node, err := crdt.NewNode(
		context.Background(),
		log,
		store,
		tests.NewRandomSyncer(log),
		tests.NewChanBroadcaster(log),
	)
	require.NoError(t, err)
	topic, err := store.NewTopic(t0, node)
	require.NoError(t, err)
	topics, err := store.Topics()
	require.NoError(t, err)
	require.Len(t, topics, 1)
	require.Equal(t, t0, topics[0])
	do(topic)
}
