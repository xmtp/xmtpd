package node

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/api"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func TestNode_NewClose(t *testing.T) {
	t.Parallel()

	_, cleanup := newTestNode(t)
	defer cleanup()
}

func newTestNode(t *testing.T) (*Node, func()) {
	s, err := New(context.Background(), test.NewLogger(t), &Options{
		API: api.Options{
			HTTPPort: 0,
			GRPCPort: 0,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, s)
	return s, s.Close
}
