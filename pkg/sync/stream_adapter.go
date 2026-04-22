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

// originatorsStreamAdapter wraps ReplicationApi_SubscribeOriginatorsClient
// so it satisfies envelopeRecvStream. An empty SubscribeOriginatorsResponse
// (no envelopes oneof set) is treated as a keepalive and returns an empty
// slice so the reader loop skips it without advancing any cursor.
type originatorsStreamAdapter struct {
	stream message_api.ReplicationApi_SubscribeOriginatorsClient
}

func (a *originatorsStreamAdapter) Recv() ([]*envelopes.OriginatorEnvelope, error) {
	resp, err := a.stream.Recv()
	if err != nil {
		return nil, err
	}
	return resp.GetEnvelopes().GetEnvelopes(), nil
}

func (a *originatorsStreamAdapter) CloseSend() error {
	return a.stream.CloseSend()
}
