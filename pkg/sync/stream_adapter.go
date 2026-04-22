package sync

import (
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
)

// envelopeRecvStream is a narrow interface that hides the specific gRPC stream
// type used by the sync worker. It lets the reader loop stay RPC-agnostic
// while the stream setup layer picks SubscribeOriginators or SubscribeEnvelopes.
type envelopeRecvStream interface {
	// Recv returns the next batch of envelopes from the peer.
	// An empty slice with no error means "no envelopes in this frame"
	// (e.g. keepalive) — callers SHOULD continue reading.
	Recv() ([]*envelopes.OriginatorEnvelope, error)
	// CloseSend mirrors grpc.ClientStream.CloseSend so sync_worker can
	// half-close the stream during cleanup.
	CloseSend() error
}

// envelopesStreamAdapter wraps ReplicationApi_SubscribeEnvelopesClient
// so it satisfies envelopeRecvStream.
type envelopesStreamAdapter struct {
	stream message_api.ReplicationApi_SubscribeEnvelopesClient
}

func (a *envelopesStreamAdapter) Recv() ([]*envelopes.OriginatorEnvelope, error) {
	resp, err := a.stream.Recv()
	if err != nil {
		return nil, err
	}
	return resp.GetEnvelopes(), nil
}

func (a *envelopesStreamAdapter) CloseSend() error {
	return a.stream.CloseSend()
}
