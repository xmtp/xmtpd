package messagev1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func Test_Publish(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	_, err := s.Publish(ctx, &messagev1.PublishRequest{})
	require.Equal(t, err, ErrTODO)
}

func Test_Subscribe(t *testing.T) {
	s := newTestService(t)
	err := s.Subscribe(&messagev1.SubscribeRequest{}, nil)
	require.Equal(t, err, ErrTODO)
}

func Test_Query(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	_, err := s.Query(ctx, &messagev1.QueryRequest{})
	require.Equal(t, err, ErrTODO)
}

func Test_BatchQuery(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	_, err := s.BatchQuery(ctx, &messagev1.BatchQueryRequest{})
	require.Equal(t, err, ErrTODO)
}

func Test_SubscribeAll(t *testing.T) {
	s := newTestService(t)
	err := s.SubscribeAll(&messagev1.SubscribeAllRequest{}, nil)
	require.Equal(t, err, ErrTODO)
}

func newTestService(t *testing.T) *Service {
	log := test.NewLogger(t)
	s, err := NewService(log)
	require.NoError(t, err)
	return s
}
