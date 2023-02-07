package memstore

import (
	"context"
	"sort"

	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	crdtmemstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type MemoryStore struct {
	crdtmemstore.MemoryStore

	log *zap.Logger

	envsByTime []*messagev1.Envelope
}

func New(log *zap.Logger) *MemoryStore {
	return &MemoryStore{
		log: log,
	}
}

func (s *MemoryStore) Close() error {
	return nil
}

func (s *MemoryStore) InsertEnvelope(ctx context.Context, env *messagev1.Envelope) error {
	i, _ := sort.Find(len(s.envsByTime), func(i int) int {
		return int(s.envsByTime[i].TimestampNs - env.TimestampNs)
	})
	if i == len(s.envsByTime) {
		s.envsByTime = append(s.envsByTime, env)
	} else {
		s.envsByTime = makeRoomAt(s.envsByTime, i)
	}
	s.envsByTime[i] = env
	return nil
}

func (s *MemoryStore) QueryEnvelopes(ctx context.Context, req *messagev1.QueryRequest) (*messagev1.QueryResponse, error) {
	return nil, nil
}

// shift events from index i to the right
// to create room at the index.
func makeRoomAt(envs []*messagev1.Envelope, i int) []*messagev1.Envelope {
	// if there's enough capacity in the slice, just shift the tail
	if len(envs) < cap(envs) {
		envs = envs[:len(envs)+1]
		copy(envs[i+1:], envs[i:])
		return envs
	}
	// figure out desired capacity of a new slice
	var newCap int
	// don't need to worry about len(events) == 0
	// because of the !found append in addEvent
	if len(envs) < 1024 {
		newCap = 2 * len(envs)
	} else {
		newCap = len(envs) + 1024
	}
	// copy events into a new slice, leaving a gap at index i
	newEvents := make([]*messagev1.Envelope, len(envs)+1, newCap)
	copy(newEvents, envs[:i])
	copy(newEvents[i+1:], envs[i:])
	return newEvents
}
