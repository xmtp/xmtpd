package node

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"

	"github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/require"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	rtest "github.com/xmtp/xmtpd/pkg/crdt/testing"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
)

func Test_FetchRW(t *testing.T) {
	s := &syncer{topic: "topic"}
	var evs []*types.Event
	var cids []multihash.Multihash
	for i := 0; i < 3; i++ {
		ev, _ := types.NewEvent(&messagev1.Envelope{
			ContentTopic: s.topic,
			TimestampNs:  uint64(i),
			Message:      []byte(fmt.Sprintf("msg-%d", i)),
		}, nil)
		evs = append(evs, ev)
		cids = append(cids, ev.Cid)
	}
	t.Run("request", func(t *testing.T) {
		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		require.NoError(t, s.writeFetchRequest(w, cids))
		r := bufio.NewReader(&buf)
		code, err := r.ReadByte()
		require.NoError(t, err)
		require.Equal(t, reqFetch, code)
		topic, rcids, err := readFetchRequest(r)
		require.NoError(t, err)
		require.Equal(t, s.topic, topic)
		require.Equal(t, cids, rcids)
	})

	t.Run("response", func(t *testing.T) {
		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		require.NoError(t, writeFetchResponse(w, s.topic, evs))
		r := bufio.NewReader(&buf)
		revs, err := s.readFetchResponse(r)
		require.NoError(t, err)
		rtest.RequireEventsEqual(t, evs, revs)
	})
}
