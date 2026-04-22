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

// For server-streaming gRPC RPCs, a server that doesn't implement the method
// returns codes.Unimplemented via trailers on the first Recv call — not from
// the client method itself. These tests exercise that real-world shape.

func TestSubscribeWithFallback_UsesOriginatorsWhenAvailable(t *testing.T) {
	nodeID := uint32(200)
	node := registryTestUtils.CreateNode(nodeID, 999, testutils.RandomPrivateKey(t))
	keepalive := &message_api.SubscribeOriginatorsResponse{}
	origStream := messageApiMocks.NewMockReplicationApi_SubscribeOriginatorsClient[*message_api.SubscribeOriginatorsResponse](
		t,
	)
	origStream.EXPECT().Recv().Return(keepalive, nil).Once()
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
	adapter, ok := stream.(*originatorsStreamAdapter)
	require.True(t, ok, "expected originatorsStreamAdapter, got %T", stream)
	// The probed keepalive frame is replayed on the first Recv so the reader
	// loop does not lose it.
	require.Same(t, keepalive, adapter.firstFrame)
}

func TestSubscribeWithFallback_ReplaysFirstFrame(t *testing.T) {
	nodeID := uint32(200)
	node := registryTestUtils.CreateNode(nodeID, 999, testutils.RandomPrivateKey(t))
	firstFrame := &message_api.SubscribeOriginatorsResponse{
		Response: &message_api.SubscribeOriginatorsResponse_Envelopes_{
			Envelopes: &message_api.SubscribeOriginatorsResponse_Envelopes{},
		},
	}
	origStream := messageApiMocks.NewMockReplicationApi_SubscribeOriginatorsClient[*message_api.SubscribeOriginatorsResponse](
		t,
	)
	origStream.EXPECT().Recv().Return(firstFrame, nil).Once()
	client := messageApiMocks.NewMockReplicationApiClient(t)
	client.EXPECT().
		SubscribeOriginators(mock.Anything, mock.Anything).
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

	// First Recv yields the already-consumed frame without calling the stream.
	envs, err := stream.Recv()
	require.NoError(t, err)
	require.Empty(t, envs)

	// Second Recv delegates to the underlying stream.
	origStream.EXPECT().Recv().Return(nil, errors.New("boom")).Once()
	_, err = stream.Recv()
	require.EqualError(t, err, "boom")
}

func TestSubscribeWithFallback_FallsBackWhenRecvReturnsUnimplemented(t *testing.T) {
	nodeID := uint32(200)
	node := registryTestUtils.CreateNode(nodeID, 999, testutils.RandomPrivateKey(t))

	origStream := messageApiMocks.NewMockReplicationApi_SubscribeOriginatorsClient[*message_api.SubscribeOriginatorsResponse](
		t,
	)
	origStream.EXPECT().
		Recv().
		Return(nil, status.Error(codes.Unimplemented, "method SubscribeOriginators not implemented")).
		Once()

	envStream := messageApiMocks.NewMockReplicationApi_SubscribeEnvelopesClient[*message_api.SubscribeEnvelopesResponse](
		t,
	)
	client := messageApiMocks.NewMockReplicationApiClient(t)
	client.EXPECT().
		SubscribeOriginators(mock.Anything, mock.Anything).
		Return(origStream, nil).
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

func TestSubscribeWithFallback_FallsBackWhenMethodReturnsUnimplemented(t *testing.T) {
	// Defensive: some gRPC transports may surface Unimplemented directly from
	// the method call instead of the first Recv. The fallback helper handles
	// both shapes.
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
		SubscribeEnvelopes(mock.Anything, mock.Anything).
		Return(envStream, nil).
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
	_, ok := stream.(*envelopesStreamAdapter)
	require.True(t, ok, "expected envelopesStreamAdapter, got %T", stream)
}

func TestSubscribeWithFallback_PropagatesNonUnimplementedMethodError(t *testing.T) {
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

func TestSubscribeWithFallback_PropagatesNonUnimplementedRecvError(t *testing.T) {
	nodeID := uint32(200)
	node := registryTestUtils.CreateNode(nodeID, 999, testutils.RandomPrivateKey(t))

	wantErr := status.Error(codes.Unavailable, "connection reset")
	origStream := messageApiMocks.NewMockReplicationApi_SubscribeOriginatorsClient[*message_api.SubscribeOriginatorsResponse](
		t,
	)
	origStream.EXPECT().Recv().Return(nil, wantErr).Once()
	client := messageApiMocks.NewMockReplicationApiClient(t)
	client.EXPECT().
		SubscribeOriginators(mock.Anything, mock.Anything).
		Return(origStream, nil).
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
	origStream := messageApiMocks.NewMockReplicationApi_SubscribeOriginatorsClient[*message_api.SubscribeOriginatorsResponse](
		t,
	)
	origStream.EXPECT().
		Recv().
		Return(nil, status.Error(codes.Unimplemented, "not implemented")).
		Once()
	client := messageApiMocks.NewMockReplicationApiClient(t)
	client.EXPECT().
		SubscribeOriginators(mock.Anything, mock.Anything).
		Return(origStream, nil).
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
