package sync

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// subscribeWithFallback opens a sync stream against a peer, preferring the new
// SubscribeOriginators RPC. For server-streaming RPCs, codes.Unimplemented is
// returned via the first Recv (as trailers) rather than from the method call
// itself, so capability is probed by consuming one frame. The SubscribeOriginators
// server sends an immediate keepalive, so the probe unblocks promptly. If the
// probe returns Unimplemented (e.g. a v1.2.0 node that has not yet registered
// the new RPC), it falls back to the legacy SubscribeEnvelopes RPC. Any other
// error is propagated to the caller.
//
// The fallback decision is evaluated per stream setup; no per-peer caching is
// used, so peers upgraded between reconnects are picked up on the next attempt.
func subscribeWithFallback(
	ctx context.Context,
	client message_api.ReplicationApiClient,
	originatorNodeIDs []uint32,
	cursor *envelopes.Cursor,
	logger *zap.Logger,
	node *registry.Node,
) (envelopeRecvStream, error) {
	origStream, err := client.SubscribeOriginators(
		ctx,
		&message_api.SubscribeOriginatorsRequest{
			Filter: &message_api.SubscribeOriginatorsRequest_OriginatorFilter{
				OriginatorNodeIds: originatorNodeIDs,
				LastSeen:          cursor,
			},
		},
	)
	if err != nil {
		if status.Code(err) == codes.Unimplemented {
			return subscribeEnvelopesFallback(
				ctx,
				client,
				originatorNodeIDs,
				cursor,
				logger,
				node,
			)
		}
		return nil, err
	}

	firstFrame, err := origStream.Recv()
	if err != nil {
		if status.Code(err) == codes.Unimplemented {
			return subscribeEnvelopesFallback(
				ctx,
				client,
				originatorNodeIDs,
				cursor,
				logger,
				node,
			)
		}
		return nil, err
	}

	metrics.EmitSyncSubscribeRPC("originators", node.NodeID)
	return &originatorsStreamAdapter{
		stream:     origStream,
		firstFrame: firstFrame,
	}, nil
}

func subscribeEnvelopesFallback(
	ctx context.Context,
	client message_api.ReplicationApiClient,
	originatorNodeIDs []uint32,
	cursor *envelopes.Cursor,
	logger *zap.Logger,
	node *registry.Node,
) (envelopeRecvStream, error) {
	logger.Info(
		"peer does not support SubscribeOriginators, falling back",
		utils.OriginatorIDField(node.NodeID),
		utils.NodeHTTPAddressField(node.HTTPAddress),
	)
	metrics.EmitSyncSubscribeRPC("envelopes", node.NodeID)

	envStream, err := client.SubscribeEnvelopes(
		ctx,
		&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: originatorNodeIDs,
				LastSeen:          cursor,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return &envelopesStreamAdapter{stream: envStream}, nil
}
