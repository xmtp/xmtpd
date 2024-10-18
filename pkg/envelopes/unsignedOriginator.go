package envelopes

import (
	"errors"

	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
)

type UnsignedOriginatorEnvelope struct {
	proto         *envelopesProto.UnsignedOriginatorEnvelope
	PayerEnvelope PayerEnvelope
}

// Construct an UnsignedOriginatorEnvelope and perform validations on any child fields.
// Does not verify signatures
func NewUnsignedOriginatorEnvelope(
	proto *envelopesProto.UnsignedOriginatorEnvelope,
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
	message, err := utils.UnmarshalUnsignedEnvelope(bytes)
	if err != nil {
		return nil, err
	}
	return NewUnsignedOriginatorEnvelope(message)
}

func (u *UnsignedOriginatorEnvelope) Proto() *envelopesProto.UnsignedOriginatorEnvelope {
	return u.proto
}
