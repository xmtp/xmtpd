package envelopes

import (
	"errors"

	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
)

type PayerEnvelope struct {
	proto          *message_api.PayerEnvelope
	ClientEnvelope ClientEnvelope
}

func NewPayerEnvelope(proto *message_api.PayerEnvelope) (*PayerEnvelope, error) {
	if proto == nil {
		return nil, errors.New("proto is nil")
	}

	clientEnv, err := NewClientEnvelopeFromBytes(proto.UnsignedClientEnvelope)
	if err != nil {
		return nil, err
	}
	return &PayerEnvelope{proto: proto, ClientEnvelope: *clientEnv}, nil
}

func (p *PayerEnvelope) Proto() *message_api.PayerEnvelope {
	return p.proto
}
