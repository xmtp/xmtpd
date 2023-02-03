package crdt

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	mh "github.com/multiformats/go-multihash"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"google.golang.org/protobuf/proto"
)

var ErrLinkTooLong = errors.New("Maximum link size exceeded")
var ErrLinksTooLong = errors.New("Total size of event links exceeds maximum limit")

// Event represents a node in the Merkle-Clock
// It captures a message and links to its preceding Events
type Event struct {
	*messagev1.Envelope
	Links []mh.Multihash // cid's of direct ancestors
	Cid   mh.Multihash   // cid is computed by hashing the links and message together
}

// NewEvent creates an event from a message and a set of links to preceding events (heads)
func NewEvent(env *messagev1.Envelope, heads []mh.Multihash) (*Event, error) {
	ev := &Event{Envelope: env, Links: heads}
	var err error
	ev.Cid, err = mh.SumStream(ev.Reader(), mh.SHA2_256, -1)
	if err != nil {
		return nil, err
	}
	return ev, nil
}

// Reader creates a chunk reader for given Event.
func (ev *Event) Reader() *chunkReader {
	// compose the chunks of the Event data
	var chunks [][]byte
	if ev.Envelope != nil {
		head := make([]byte, 8+len(ev.ContentTopic))
		binary.BigEndian.PutUint64(head, ev.TimestampNs) // timestamp
		copy(head[8:], ev.ContentTopic)                  // topic
		chunks = append(chunks, head, ev.Message)        // message payload
	}
	for _, link := range ev.Links { // links
		chunks = append(chunks, link)
	}
	return &chunkReader{chunks, 0}
}

// Compare returns an integer comparing two events based on their timestamps.
// The result will be negative if ev < ev2, and positive if ev > ev2.
// The result can only be 0 if ev and ev2 are the same event.
// TODO: total order should reflect the DAG first and foremost.
func (ev *Event) Compare(ev2 *Event) int {
	res := ev.TimestampNs - ev2.TimestampNs
	if res != 0 {
		return int(res)
	}
	return bytes.Compare(ev.Cid, ev2.Cid)
}

// chunkReader helps computing an Event CID efficiently by
// yielding the bytes composed of the various bits of the Event
// without having to concatenate them all.
// This allows passing the reader to mh.SumStream()
type chunkReader struct {
	unreadChunks [][]byte // chunks of the Event data to be hashed
	pos          int      // current position from the start of the next chunk
}

func (r *chunkReader) Read(b []byte) (n int, err error) {
	total := 0
	for len(b) > 0 && len(r.unreadChunks) > 0 {
		chunk := r.unreadChunks[0]
		n := copy(b, chunk[r.pos:])
		total += n
		b = b[n:]
		r.pos += n
		if r.pos == len(chunk) {
			r.pos = 0
			r.unreadChunks = r.unreadChunks[1:]
		}
	}
	if len(r.unreadChunks) > 0 {
		return total, nil
	}
	return total, io.EOF
}

// Storage Marshaling Support

// ByTimeKey is used as an indexing key in a ByTime index.
// The key must sort alphabeticely in the desired time order,
// It is composed of 8 bytes of the timestampNs, and cid hash after.
func ToByTimeKey(timestampNs uint64, cid []byte) []byte {
	key := make([]byte, 8+len(cid))
	binary.BigEndian.PutUint64(key, timestampNs)
	copy(key[8:], cid)
	return key
}

func FromByTimeKey(key []byte) (timestampNs uint64, cid []byte) {
	timestampNs = binary.BigEndian.Uint64(key)
	cid = key[8:]
	return timestampNs, cid
}

// Event bytes are encoded as a list of event links followed by the Envelope.
// Links are encoded as total size of the list in bytes encoded as uint16,
// followed by each link encoded as uint8 size prefix followed by the link bytes.
// The total size prefix allows efficiently skipping to the Envelope bytes.
// Envelope bytes are stored as prescribed by github.com/xmtp/proto.

func EventToBytes(ev *Event) ([]byte, error) {
	envBytes, err := proto.Marshal(ev.Envelope)
	if err != nil {
		return nil, err
	}
	linksSize := 0
	for _, link := range ev.Links {
		if len(link) > 255 {
			return nil, ErrLinkTooLong
		}
		linksSize += len(link) + 1
	}
	if linksSize > 65535 {
		return nil, ErrLinksTooLong
	}
	b := make([]byte, 2+linksSize+len(envBytes))
	binary.BigEndian.PutUint16(b, uint16(linksSize))
	pos := 2
	for _, link := range ev.Links {
		b[pos] = byte(len(link))
		copy(b[pos+1:], link)
		pos += len(link) + 1
	}
	if pos != linksSize+2 {
		panic("shouldn't happen")
	}
	copy(b[pos:], envBytes)
	return b, nil
}

func EventFromBytes(cid []byte, evBytes []byte) (*Event, error) {
	links, remainder, err := readLinks(evBytes)
	if err != nil {
		return nil, err
	}
	env, err := unmarshalEnvelope(remainder)
	if err != nil {
		return nil, err
	}
	return &Event{
		Cid:      cid,
		Links:    links,
		Envelope: env,
	}, nil
}

func EnvelopeFromBytes(evBytes []byte) (*messagev1.Envelope, error) {
	linksSize := int(binary.BigEndian.Uint16(evBytes))
	return unmarshalEnvelope(evBytes[linksSize+2:])
}

func LinksFromBytes(evBytes []byte) (cids []mh.Multihash, err error) {
	cids, _, err = readLinks(evBytes)
	return cids, err
}

func readLinks(evBytes []byte) (cids []mh.Multihash, remainder []byte, err error) {
	var links []mh.Multihash
	linksSize := int(binary.BigEndian.Uint16(evBytes))
	pos := 2
	for pos < linksSize {
		linkSize := int(evBytes[pos])
		pos++
		links = append(links, evBytes[pos:pos+linkSize])
		pos += linkSize
	}
	return links, evBytes[pos:], nil
}

func unmarshalEnvelope(envBytes []byte) (*messagev1.Envelope, error) {
	var env messagev1.Envelope
	if err := proto.Unmarshal(envBytes, &env); err != nil {
		return nil, err
	}
	return &env, nil
}
