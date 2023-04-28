package testing

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/api/client"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

type testSubscriber struct {
	Topic  string
	stream client.Stream
	envs   []*messagev1.Envelope
	sync.RWMutex
}

func (s *testSubscriber) RequireEventuallyCapturedEvents(t *testing.T, expected []*messagev1.Envelope) {
	t.Helper()
	require.Eventually(t, func() bool {
		s.RLock()
		defer s.RUnlock()
		return len(s.envs) == len(expected)
	}, 3*time.Second, 10*time.Millisecond)
	test.RequireProtoEqual(t, expected, s.envs)
}
