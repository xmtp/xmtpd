package payer

import (
	"context"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NodeCursorTracker struct {
	ctx           context.Context
	log           *zap.Logger
	clientManager *ClientManager
}

func NewNodeCursorTracker(ctx context.Context,
	log *zap.Logger, clientManager *ClientManager) *NodeCursorTracker {
	return &NodeCursorTracker{ctx: ctx, log: log, clientManager: clientManager}
}

func (ct *NodeCursorTracker) BlockUntilDesiredCursorReached(
	ctx context.Context,
	nodeId uint32,
	desiredOriginatorId uint32,
	desiredSequenceId uint64,
) error {
	// TODO(mkysel) ideally we wouldn't create and tear down the stream for every request

	conn, err := ct.clientManager.GetClient(nodeId)
	if err != nil {
		return err
	}
	client := metadata_api.NewMetadataApiClient(conn)
	stream, err := client.SubscribeSyncCursor(ctx, &metadata_api.GetSyncCursorRequest{})
	if err != nil {
		return err
	}
	for {
		select {
		case <-ct.ctx.Done():
			// server is shutting down
			return status.Errorf(codes.Canceled, "node terminated. Cancelled wait for cursor")
		case <-ctx.Done():
			// client has shut down
			return nil
		default:
			resp, err := stream.Recv()
			if err != nil {
				if status.Code(err) == codes.Canceled {
					return nil
				}
				// TODO(mkysel): proper handling of failures
				return err
			}
			if err != nil || resp == nil || resp.LatestSync == nil {
				return status.Errorf(codes.Internal, "error getting node cursor: %v", err)
			}
			derefMap := resp.LatestSync.NodeIdToSequenceId
			seqId, exists := derefMap[desiredOriginatorId]
			if !exists {
				continue // Wait for the originator ID to appear
			}

			// Check if the sequence ID has reached the desired value
			if seqId >= desiredSequenceId {
				return nil // Desired state achieved
			}
		}
	}
}
