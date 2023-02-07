package types

import (
	"bytes"
	"crypto/rand"
	"testing"

	mh "github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/assert"
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
	_, err := rand.Reader.Read(payload)
	assert.NoError(t, err)
	links := makeLinks(t, "one", "two", "three")
	ev, err := NewEvent([]byte("payload"), links)
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
	chunks := [][]byte{ev.Payload}
	for _, l := range ev.Links {
		chunks = append(chunks, l)
	}
	sum, err := mh.Sum(bytes.Join(chunks, nil), mh.SHA2_256, -1)
	assert.NoError(t, err)
	return sum
}
