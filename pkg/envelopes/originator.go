package envelopes

import (
	"errors"

	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
)

type OriginatorEnvelope struct {
	proto                      *message_api.OriginatorEnvelope
	UnsignedOriginatorEnvelope UnsignedOriginatorEnvelope
}

func NewOriginatorEnvelope(proto *message_api.OriginatorEnvelope) (*OriginatorEnvelope, error) {
	if proto == nil {
		return nil, errors.New("proto is nil")
	}

	unsigned, err := NewUnsignedOriginatorEnvelopeFromBytes(proto.UnsignedOriginatorEnvelope)
	if err != nil {
		return nil, err
	}

	return &OriginatorEnvelope{
		proto:                      proto,
		UnsignedOriginatorEnvelope: *unsigned,
	}, nil
}

func (o *OriginatorEnvelope) Proto() *message_api.OriginatorEnvelope {
	return o.proto
}
