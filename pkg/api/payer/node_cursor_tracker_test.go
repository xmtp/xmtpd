package payer_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/api/payer"
	metadataMocks "github.com/xmtp/xmtpd/pkg/mocks/metadata_api"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func TestCursorTrackerBasic(t *testing.T) {
	ctx := context.Background()

	desiredOriginator := uint32(1)
	desiredSequence := uint64(1)
	tracker := constructTracker(t, ctx, []*metadata_api.GetSyncCursorResponse{
		{
			LatestSync: &envelopesProto.Cursor{
				NodeIdToSequenceId: map[uint32]uint64{desiredOriginator: desiredSequence},
			},
		},
	})

	err := tracker.BlockUntilDesiredCursorReached(
		ctx,
		1,
		desiredOriginator,
		desiredSequence,
	)
	require.NoError(t, err)
}

func TestCursorTrackerClientShutsDown(t *testing.T) {
	ctx := context.Background()

	desiredOriginator := uint32(1)
	desiredSequence := uint64(1)
	tracker := constructTracker(t, ctx, []*metadata_api.GetSyncCursorResponse{
		{
			LatestSync: &envelopesProto.Cursor{
				NodeIdToSequenceId: map[uint32]uint64{desiredOriginator: desiredSequence},
			},
		},
	})

	err := tracker.BlockUntilDesiredCursorReached(
		testutils.CancelledContext(),
		1,
		desiredOriginator,
		desiredSequence+100,
	)
	require.NoError(t, err)
}

func TestCursorTrackerClientShutsDownAfterExecution(t *testing.T) {
	ctx := context.Background()

	desiredOriginator := uint32(1)
	desiredSequence := uint64(1)
	tracker := constructTracker(t, ctx, []*metadata_api.GetSyncCursorResponse{
		{
			LatestSync: &envelopesProto.Cursor{
				NodeIdToSequenceId: map[uint32]uint64{desiredOriginator: desiredSequence},
			},
		},
	})

	clientCtx, cancel := context.WithCancel(ctx)
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	err := tracker.BlockUntilDesiredCursorReached(
		clientCtx,
		1,
		desiredOriginator,
		desiredSequence+100,
	)
	require.NoError(t, err)
}

func TestCursorTrackerClientServerIsShutdown(t *testing.T) {
	ctx := context.Background()

	desiredOriginator := uint32(1)
	desiredSequence := uint64(1)

	tracker := constructTracker(
		t,
		testutils.CancelledContext(),
		[]*metadata_api.GetSyncCursorResponse{
			{
				LatestSync: &envelopesProto.Cursor{
					NodeIdToSequenceId: map[uint32]uint64{desiredOriginator: desiredSequence},
				},
			},
		},
	)

	err := tracker.BlockUntilDesiredCursorReached(
		ctx,
		1,
		desiredOriginator,
		desiredSequence+100,
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "node terminated")
}

func TestCursorTrackerClientServerShutsDown(t *testing.T) {
	ctx := context.Background()

	desiredOriginator := uint32(1)
	desiredSequence := uint64(1)

	serverCtx, cancel := context.WithCancel(ctx)

	tracker := constructTracker(
		t,
		serverCtx,
		[]*metadata_api.GetSyncCursorResponse{
			{
				LatestSync: &envelopesProto.Cursor{
					NodeIdToSequenceId: map[uint32]uint64{desiredOriginator: desiredSequence},
				},
			},
		},
	)
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	err := tracker.BlockUntilDesiredCursorReached(
		ctx,
		1,
		desiredOriginator,
		desiredSequence+100,
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "node terminated")
}

func TestCursorTrackerOriginatorDoesNotExist(t *testing.T) {
	ctx := context.Background()

	desiredOriginator := uint32(1)
	desiredSequence := uint64(1)
	tracker := constructTracker(t, ctx, []*metadata_api.GetSyncCursorResponse{
		{
			LatestSync: &envelopesProto.Cursor{
				NodeIdToSequenceId: map[uint32]uint64{desiredOriginator: desiredSequence},
			},
		},
	})

	clientCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	err := tracker.BlockUntilDesiredCursorReached(
		clientCtx,
		1,
		5,
		desiredSequence,
	)
	require.NoError(t, err)
}

func TestCursorTrackerSequenceDoesNotExist(t *testing.T) {
	ctx := context.Background()

	desiredOriginator := uint32(1)
	desiredSequence := uint64(1)
	tracker := constructTracker(t, ctx, []*metadata_api.GetSyncCursorResponse{
		{
			LatestSync: &envelopesProto.Cursor{
				NodeIdToSequenceId: map[uint32]uint64{desiredOriginator: desiredSequence},
			},
		},
	})

	clientCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	err := tracker.BlockUntilDesiredCursorReached(
		clientCtx,
		1,
		desiredOriginator,
		desiredSequence+100,
	)
	require.NoError(t, err)
}

func TestCursorTrackerThreeStages(t *testing.T) {
	ctx := context.Background()

	desiredOriginator := uint32(1)
	desiredSequence := uint64(1)
	tracker := constructTracker(t, ctx, []*metadata_api.GetSyncCursorResponse{
		{
			LatestSync: &envelopesProto.Cursor{
				NodeIdToSequenceId: map[uint32]uint64{desiredOriginator: desiredSequence},
			},
		},
		{
			LatestSync: &envelopesProto.Cursor{
				NodeIdToSequenceId: map[uint32]uint64{desiredOriginator: desiredSequence + 1},
			},
		},
		{
			LatestSync: &envelopesProto.Cursor{
				NodeIdToSequenceId: map[uint32]uint64{desiredOriginator: desiredSequence + 2},
			},
		},
	})

	err := tracker.BlockUntilDesiredCursorReached(
		testutils.CancelledContext(),
		1,
		desiredOriginator,
		desiredSequence+2,
	)
	require.NoError(t, err)
}

func constructTracker(
	t *testing.T,
	ctx context.Context,
	updates []*metadata_api.GetSyncCursorResponse,
) *payer.NodeCursorTracker {
	mockStream := &MockSubscribeSyncCursorClient{
		updates: updates,
	}

	metaMocks := metadataMocks.NewMockMetadataApiClient(t)
	metaMocks.On("SubscribeSyncCursor", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			// Capture the context from the caller
			capturedCtx := args.Get(0).(context.Context)
			mockStream.ctx = capturedCtx // Store the captured context in the mock
		}).
		Return(mockStream, nil). // ✅ Return the instance directly
		Once()

	metadataConstructor := &FixedMetadataAPIClientConstructor{
		mockClient: metaMocks,
	}
	var interf payer.MetadataAPIClientConstructor = metadataConstructor

	tracker := payer.NewNodeCursorTracker(ctx, testutils.NewLog(t), interf)
	return tracker
}
