package node_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	messagev1 "github.com/xmtp/xmtpd/pkg/api/message/v1"
	memstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
	"github.com/xmtp/xmtpd/pkg/node"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func TestNode_NewClose(t *testing.T) {
	t.Parallel()

	_, cleanup := newTestNode(t)
	defer cleanup()
}

func newTestNode(t *testing.T) (*node.Node, func()) {
	ctx := context.Background()
	log := test.NewLogger(t)
	store := memstore.New(log)
	messagev1, err := messagev1.New(log, store)
	require.NoError(t, err)
	s, err := node.New(ctx, log, messagev1, &node.Options{})
	require.NoError(t, err)
	require.NotNil(t, s)
	return s, s.Close
}
