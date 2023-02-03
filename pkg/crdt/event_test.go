package crdt

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"testing"

	mh "github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/assert"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
)

func Test_NilEventCID(t *testing.T) {
	ev, err := NewEvent(nil, nil)
	assert.NoError(t, err)
	emptyHash, err := mh.Sum(nil, mh.SHA2_256, -1)
	assert.NoError(t, err)
	assert.Equal(t, ev.Cid, emptyHash)
}

func Test_EventCID(t *testing.T) {
	payload := make([]byte, 1000)
	_, err := rand.Reader.Read(payload)
	assert.NoError(t, err)
	env := &messagev1.Envelope{TimestampNs: 1, ContentTopic: "topic", Message: payload}
	links := makeLinks(t, "one", "two", "three")
	ev, err := NewEvent(env, links)
	assert.NoError(t, err)
	assert.Equal(t, computeCid(t, ev), ev.Cid)
}

func makeLinks(t *testing.T, payloads ...string) (links []mh.Multihash) {
	for _, p := range payloads {
		hash, err := mh.Sum([]byte(p), mh.SHA2_256, -1)
		assert.NoError(t, err)
		links = append(links, hash)
	}
	return links
}

func computeCid(t *testing.T, ev *Event) mh.Multihash {
	timestampBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timestampBytes, ev.TimestampNs)
	chunks := [][]byte{
		timestampBytes,
		[]byte(ev.ContentTopic),
		ev.Message,
	}
	for _, l := range ev.Links {
		chunks = append(chunks, l)
	}
	sum, err := mh.Sum(bytes.Join(chunks, nil), mh.SHA2_256, -1)
	assert.NoError(t, err)
	return sum
}
