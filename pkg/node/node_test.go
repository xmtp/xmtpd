package node_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	messagev1 "github.com/xmtp/xmtpd/pkg/api/message/v1"
	membroadcaster "github.com/xmtp/xmtpd/pkg/crdt/broadcasters/mem"
	memstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
	memsyncer "github.com/xmtp/xmtpd/pkg/crdt/syncers/mem"
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
	bc := membroadcaster.New(log)
	syncer := memsyncer.New(log, store)
	messagev1, err := messagev1.New(log, store, bc, syncer)
	require.NoError(t, err)
	s, err := node.New(ctx, log, messagev1, &node.Options{})
	require.NoError(t, err)
	require.NotNil(t, s)
	return s, s.Close
}
