package storetest

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/store"
)

type TestStoreMaker func(t *testing.T) *TestStore

type TestStore struct {
	store.Store
}

func NewTestStore(s store.Store) *TestStore {
	return &TestStore{s}
}

func (s *TestStore) seed(t *testing.T, count int) []*messagev1.Envelope {
	ctx := context.Background()
	envs := make([]*messagev1.Envelope, count)
	for i := 0; i < 20; i++ {
		env := &messagev1.Envelope{
			ContentTopic: "topic",
			TimestampNs:  uint64(i + 1),
			Message:      []byte(fmt.Sprintf("msg-%d", i+1)),
		}
		err := s.InsertEnvelope(ctx, env)
		require.NoError(t, err)
		envs[i] = env
	}
	return envs
}

func (s *TestStore) query(t *testing.T, req *messagev1.QueryRequest) (*messagev1.QueryResponse, error) {
	ctx := context.Background()
	return s.QueryEnvelopes(ctx, req)
}
