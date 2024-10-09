package envelopes

import (
	"errors"

	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"google.golang.org/protobuf/proto"
)

type UnsignedOriginatorEnvelope struct {
	proto         *message_api.UnsignedOriginatorEnvelope
	PayerEnvelope PayerEnvelope
}

// Construct an UnsignedOriginatorEnvelope and perform validations on any child fields.
// Does not verify signatures
func NewUnsignedOriginatorEnvelope(
	proto *message_api.UnsignedOriginatorEnvelope,
) (*UnsignedOriginatorEnvelope, error) {
	if proto == nil {
		return nil, errors.New("proto is nil")
	}

	payer, err := NewPayerEnvelope(proto.PayerEnvelope)
	if err != nil {
		return nil, err
	}

	return &UnsignedOriginatorEnvelope{proto: proto, PayerEnvelope: *payer}, nil
}

func (u *UnsignedOriginatorEnvelope) OriginatorNodeID() uint32 {
	// Skip nil check because it is in the constructor
	return u.proto.OriginatorNodeId
}

func (u *UnsignedOriginatorEnvelope) OriginatorSequenceID() uint64 {
	// Skip nil check because it is in the constructor
	return u.proto.OriginatorSequenceId
}

func (u *UnsignedOriginatorEnvelope) OriginatorNs() int64 {
	// Skip nil check because it is in the constructor
	return u.proto.OriginatorNs
}

func NewUnsignedOriginatorEnvelopeFromBytes(bytes []byte) (*UnsignedOriginatorEnvelope, error) {
	var message message_api.UnsignedOriginatorEnvelope
	if err := proto.Unmarshal(bytes, &message); err != nil {
		return nil, err
	}
	return NewUnsignedOriginatorEnvelope(&message)
}

func (u *UnsignedOriginatorEnvelope) Proto() *message_api.UnsignedOriginatorEnvelope {
	return u.proto
}
