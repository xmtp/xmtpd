package sync

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/testutils"
	messageApiMocks "github.com/xmtp/xmtpd/pkg/testutils/mocks/message_api"
	registryTestUtils "github.com/xmtp/xmtpd/pkg/testutils/registry"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSubscribeWithFallback_UsesOriginatorsWhenAvailable(t *testing.T) {
	nodeID := uint32(200)
	node := registryTestUtils.CreateNode(nodeID, 999, testutils.RandomPrivateKey(t))
	origStream := messageApiMocks.NewMockReplicationApi_SubscribeOriginatorsClient[*message_api.SubscribeOriginatorsResponse](
		t,
	)
	client := messageApiMocks.NewMockReplicationApiClient(t)
	client.EXPECT().
		SubscribeOriginators(
			mock.Anything,
			mock.MatchedBy(func(req *message_api.SubscribeOriginatorsRequest) bool {
				f := req.GetFilter()
				return len(f.GetOriginatorNodeIds()) == 1 &&
					f.GetOriginatorNodeIds()[0] == nodeID &&
					f.GetLastSeen() != nil
			}),
		).
		Return(origStream, nil).
		Once()

	stream, err := subscribeWithFallback(
		context.Background(),
		client,
		[]uint32{nodeID},
		&envelopes.Cursor{NodeIdToSequenceId: map[uint32]uint64{}},
		testutils.NewLog(t),
		&node,
	)
	require.NoError(t, err)
	_, ok := stream.(*originatorsStreamAdapter)
	require.True(t, ok, "expected originatorsStreamAdapter, got %T", stream)
}

func TestSubscribeWithFallback_FallsBackOnUnimplemented(t *testing.T) {
	nodeID := uint32(200)
	node := registryTestUtils.CreateNode(nodeID, 999, testutils.RandomPrivateKey(t))

	envStream := messageApiMocks.NewMockReplicationApi_SubscribeEnvelopesClient[*message_api.SubscribeEnvelopesResponse](
		t,
	)
	client := messageApiMocks.NewMockReplicationApiClient(t)
	client.EXPECT().
		SubscribeOriginators(mock.Anything, mock.Anything).
		Return(nil, status.Error(codes.Unimplemented, "not implemented")).
		Once()
	client.EXPECT().
		SubscribeEnvelopes(
			mock.Anything,
			mock.MatchedBy(func(req *message_api.SubscribeEnvelopesRequest) bool {
				q := req.GetQuery()
				return len(q.GetOriginatorNodeIds()) == 1 &&
					q.GetOriginatorNodeIds()[0] == nodeID &&
					q.GetLastSeen() != nil
			}),
		).
		Return(envStream, nil).
		Once()

	core, recorded := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	stream, err := subscribeWithFallback(
		context.Background(),
		client,
		[]uint32{nodeID},
		&envelopes.Cursor{NodeIdToSequenceId: map[uint32]uint64{}},
		logger,
		&node,
	)
	require.NoError(t, err)
	_, ok := stream.(*envelopesStreamAdapter)
	require.True(t, ok, "expected envelopesStreamAdapter, got %T", stream)

	require.NotEmpty(
		t,
		recorded.FilterMessage("peer does not support SubscribeOriginators, falling back").All(),
		"expected info log about fallback",
	)
}

func TestSubscribeWithFallback_PropagatesNonUnimplementedError(t *testing.T) {
	nodeID := uint32(200)
	node := registryTestUtils.CreateNode(nodeID, 999, testutils.RandomPrivateKey(t))

	wantErr := status.Error(codes.Unavailable, "peer down")
	client := messageApiMocks.NewMockReplicationApiClient(t)
	client.EXPECT().
		SubscribeOriginators(mock.Anything, mock.Anything).
		Return(nil, wantErr).
		Once()
	// SubscribeEnvelopes must NOT be called.

	_, err := subscribeWithFallback(
		context.Background(),
		client,
		[]uint32{nodeID},
		&envelopes.Cursor{NodeIdToSequenceId: map[uint32]uint64{}},
		testutils.NewLog(t),
		&node,
	)
	require.ErrorIs(t, err, wantErr)
}

func TestSubscribeWithFallback_PropagatesEnvelopesError(t *testing.T) {
	nodeID := uint32(200)
	node := registryTestUtils.CreateNode(nodeID, 999, testutils.RandomPrivateKey(t))

	wantErr := errors.New("envelopes failed")
	client := messageApiMocks.NewMockReplicationApiClient(t)
	client.EXPECT().
		SubscribeOriginators(mock.Anything, mock.Anything).
		Return(nil, status.Error(codes.Unimplemented, "not implemented")).
		Once()
	client.EXPECT().
		SubscribeEnvelopes(mock.Anything, mock.Anything).
		Return(nil, wantErr).
		Once()

	_, err := subscribeWithFallback(
		context.Background(),
		client,
		[]uint32{nodeID},
		&envelopes.Cursor{NodeIdToSequenceId: map[uint32]uint64{}},
		testutils.NewLog(t),
		&node,
	)
	require.ErrorIs(t, err, wantErr)
}
