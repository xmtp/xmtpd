package types

import (
	"bytes"
	cryptorand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"testing"

	mh "github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/assert"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func TestEvent_NilCid(t *testing.T) {
	t.Parallel()
	ev, err := NewEvent(nil, nil)
	assert.NoError(t, err)
	emptyHash, err := mh.Sum(nil, mh.SHA2_256, -1)
	assert.NoError(t, err)
	assert.Equal(t, ev.Cid, emptyHash)
}

func TestEvent_ValidCid(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 1000)
	_, err := cryptorand.Reader.Read(payload)
	assert.NoError(t, err)
	links := makeLinks(t, "one", "two", "three")
	ev, err := NewEvent(newRandomEnvelope(t), links)
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
	chunks := [][]byte{}
	if ev.Envelope != nil {
		head := make([]byte, 8+len(ev.ContentTopic))
		binary.BigEndian.PutUint64(head, ev.TimestampNs) // timestamp
		copy(head[8:], ev.ContentTopic)                  // topic
		chunks = append(chunks, head, ev.Message)        // message payload
	}
	for _, l := range ev.Links {
		chunks = append(chunks, l)
	}
	sum, err := mh.Sum(bytes.Join(chunks, nil), mh.SHA2_256, -1)
	assert.NoError(t, err)
	return sum
}

func newRandomEnvelope(t *testing.T) *messagev1.Envelope {
	return &messagev1.Envelope{
		ContentTopic: "topic-" + test.RandomStringLower(5),
		TimestampNs:  uint64(rand.Intn(100)),
		Message:      []byte("msg-" + test.RandomString(13)),
	}
}
