package payer

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	metadataMocks "github.com/xmtp/xmtpd/pkg/mocks/metadata_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
	"go.uber.org/zap"
)

// Mock gRPC stream for SubscribeSyncCursor
type MockSubscribeSyncCursorClient struct {
	metadata_api.MetadataApi_SubscribeSyncCursorClient
	updates []*metadata_api.GetSyncCursorResponse
	err     error
	index   int
}

// Recv() simulates receiving cursor updates over time
func (m *MockSubscribeSyncCursorClient) Recv() (*metadata_api.GetSyncCursorResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.index < len(m.updates) {
		resp := m.updates[m.index]
		m.index++
		return resp, nil
	}
	// Simulate an open stream without new messages
	time.Sleep(50 * time.Millisecond)
	return nil, nil
}

func TestBlockUntilDesiredCursorReached_Success(t *testing.T) {
	mockClient := metadataMocks.NewMockMetadataApiClient(t)

	// Simulate a stream with sequential updates leading to the desired state
	mockStream := &MockSubscribeSyncCursorClient{
		updates: []*metadata_api.GetSyncCursorResponse{
			{
				LatestSync: &envelopes.Cursor{
					NodeIdToSequenceId: map[uint32]uint64{1: 1},
				},
			},
			{
				LatestSync: &envelopes.Cursor{
					NodeIdToSequenceId: map[uint32]uint64{1: 2},
				},
			},
			{
				LatestSync: &envelopes.Cursor{
					NodeIdToSequenceId: map[uint32]uint64{1: 3},
				},
			},
		},
	}

	// Expect `SubscribeSyncCursor` to return our mock stream
	mockClient.EXPECT().
		SubscribeSyncCursor(mock.Anything, mock.Anything).
		Return(mockStream, nil).
		Once()

	nodeRegistry := registry.NewFixedNodeRegistry([]registry.Node{
		{
			NodeID:      100,
			HttpAddress: "",
		},
		{
			NodeID:      200,
			HttpAddress: "",
		},
	})

	cm := NewClientManager(testutils.NewLog(t), nodeRegistry)

	tracker := NewNodeCursorTracker(context.Background(), zap.NewNop(), cm)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := tracker.BlockUntilDesiredCursorReached(ctx, 100, 1, 3)
	require.NoError(t, err)
}

//
//func TestBlockUntilDesiredCursorReached_StreamError(t *testing.T) {
//	mockClient := NewMockMetadataApiClient(t)
//
//	mockStream := &MockSubscribeSyncCursorClient{
//		err: status.Error(codes.Internal, "internal error"),
//	}
//
//	mockClient.EXPECT().
//		SubscribeSyncCursor(mock.Anything, mock.Anything).
//		Return(mockStream, nil).
//		Once()
//
//	clientManager := &MockClientManager{client: mockClient}
//	tracker := NewNodeCursorTracker(context.Background(), zap.NewNop(), clientManager)
//
//	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
//	defer cancel()
//
//	err := tracker.BlockUntilDesiredCursorReached(ctx, 100, 1, 3)
//	assert.Error(t, err)
//	assert.Contains(t, err.Error(), "internal error")
//}
//
//func TestBlockUntilDesiredCursorReached_ContextCancellation(t *testing.T) {
//	mockClient := NewMockMetadataApiClient(t)
//
//	mockStream := &MockSubscribeSyncCursorClient{
//		updates: []*metadata_api.GetSyncCursorResponse{
//			{LatestSync: &metadata_api.SyncCursor{NodeIdToSequenceId: map[uint32]uint64{1: 1}}},
//		},
//	}
//
//	mockClient.EXPECT().
//		SubscribeSyncCursor(mock.Anything, mock.Anything).
//		Return(mockStream, nil).
//		Once()
//
//	clientManager := &MockClientManager{client: mockClient}
//	tracker := NewNodeCursorTracker(context.Background(), zap.NewNop(), clientManager)
//
//	ctx, cancel := context.WithCancel(context.Background())
//	go func() {
//		time.Sleep(100 * time.Millisecond)
//		cancel()
//	}()
//
//	err := tracker.BlockUntilDesiredCursorReached(ctx, 100, 1, 3)
//	assert.Error(t, err)
//	assert.Equal(t, codes.Canceled, status.Code(err))
//}
//
//func TestBlockUntilDesiredCursorReached_OriginatorIDMissingInitially(t *testing.T) {
//	mockClient := NewMockMetadataApiClient(t)
//
//	mockStream := &MockSubscribeSyncCursorClient{
//		updates: []*metadata_api.GetSyncCursorResponse{
//			{LatestSync: &metadata_api.SyncCursor{NodeIdToSequenceId: map[uint32]uint64{2: 1}}},
//			{LatestSync: &metadata_api.SyncCursor{NodeIdToSequenceId: map[uint32]uint64{1: 2}}},
//			{LatestSync: &metadata_api.SyncCursor{NodeIdToSequenceId: map[uint32]uint64{1: 3}}},
//		},
//	}
//
//	mockClient.EXPECT().
//		SubscribeSyncCursor(mock.Anything, mock.Anything).
//		Return(mockStream, nil).
//		Once()
//
//	clientManager := &MockClientManager{client: mockClient}
//	tracker := NewNodeCursorTracker(context.Background(), zap.NewNop(), clientManager)
//
//	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
//	defer cancel()
//
//	err := tracker.BlockUntilDesiredCursorReached(ctx, 100, 1, 3)
//	require.NoError(t, err)
//}
