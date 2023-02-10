package messagev1_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	proto "github.com/xmtp/proto/v3/go/message_api/v1"
	messagev1 "github.com/xmtp/xmtpd/pkg/api/message/v1"
	"github.com/xmtp/xmtpd/pkg/crdt"
	membroadcaster "github.com/xmtp/xmtpd/pkg/crdt/broadcasters/mem"
	memstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
	memsyncer "github.com/xmtp/xmtpd/pkg/crdt/syncers/mem"
	crdttypes "github.com/xmtp/xmtpd/pkg/crdt/types"
	memsubs "github.com/xmtp/xmtpd/pkg/node/subscribers/mem"
	memtopics "github.com/xmtp/xmtpd/pkg/node/topics/mem"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func Test_Publish(t *testing.T) {
	s, cleanup := newTestService(t)
	defer cleanup()
	ctx := context.Background()
	_, err := s.Publish(ctx, &proto.PublishRequest{})
	require.NoError(t, err)
}

func Test_Subscribe(t *testing.T) {
	s, cleanup := newTestService(t)
	defer cleanup()
	err := s.Subscribe(&proto.SubscribeRequest{}, nil)
	require.Equal(t, err, messagev1.ErrMissingTopic)
}

func Test_Query(t *testing.T) {
	s, cleanup := newTestService(t)
	defer cleanup()
	ctx := context.Background()
	_, err := s.Query(ctx, &proto.QueryRequest{})
	require.Equal(t, err, messagev1.ErrMissingTopic)
}

func Test_BatchQuery(t *testing.T) {
	s, cleanup := newTestService(t)
	defer cleanup()
	ctx := context.Background()
	_, err := s.BatchQuery(ctx, &proto.BatchQueryRequest{})
	require.Equal(t, err, messagev1.ErrTODO)
}

func Test_SubscribeAll(t *testing.T) {
	s, cleanup := newTestService(t)
	defer cleanup()
	err := s.SubscribeAll(&proto.SubscribeAllRequest{}, nil)
	require.Equal(t, err, messagev1.ErrTODO)
}

func newTestService(t *testing.T) (*messagev1.Service, func()) {
	ctx := context.Background()
	log := test.NewLogger(t)
	store := memstore.New(log)
	bc := membroadcaster.New(log)
	syncer := memsyncer.New(log, store)
	subs := memsubs.New(log, 100)
	topics, err := memtopics.New(log, func(topicId string) (*crdt.Replica, error) {
		return crdt.NewReplica(ctx, log, store, bc, syncer,
			func(ev *crdttypes.Event) {
				subs.OnNewEvent(topicId, ev)
			},
		)
	})
	require.NoError(t, err)
	s, err := messagev1.New(log, topics, subs, store, bc, syncer)
	require.NoError(t, err)
	return s, func() {
		s.Close()
		topics.Close()
		subs.Close()
		syncer.Close()
		bc.Close()
		store.Close()
	}
}
