package node_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	messagev1 "github.com/xmtp/xmtpd/pkg/api/message/v1"
	"github.com/xmtp/xmtpd/pkg/crdt"
	membroadcaster "github.com/xmtp/xmtpd/pkg/crdt/broadcasters/mem"
	memstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
	memsyncer "github.com/xmtp/xmtpd/pkg/crdt/syncers/mem"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/node"
	memsubs "github.com/xmtp/xmtpd/pkg/node/subscribers/mem"
	memtopics "github.com/xmtp/xmtpd/pkg/node/topics/mem"
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
	subs := memsubs.New(log)
	topics, err := memtopics.New(log, func(topicId string) (*crdt.Replica, error) {
		return crdt.NewReplica(ctx, log, store, bc, syncer,
			func(ev *types.Event) {
				subs.OnNewEvent(topicId, ev)
			},
		)
	})
	require.NoError(t, err)
	messagev1, err := messagev1.New(log, topics, subs, store, bc, syncer)
	require.NoError(t, err)
	s, err := node.New(ctx, log, messagev1, &node.Options{})
	require.NoError(t, err)
	require.NotNil(t, s)
	return s, func() {
		s.Close()
		messagev1.Close()
		topics.Close()
		subs.Close()
		syncer.Close()
		bc.Close()
		store.Close()
	}
}
