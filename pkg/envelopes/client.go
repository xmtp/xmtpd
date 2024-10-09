package envelopes

import (
	"errors"

	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"google.golang.org/protobuf/proto"
)

type ClientEnvelope struct {
	proto *message_api.ClientEnvelope
}

func NewClientEnvelope(proto *message_api.ClientEnvelope) (*ClientEnvelope, error) {
	if proto == nil {
		return nil, errors.New("proto is nil")
	}

	if proto.Aad == nil {
		return nil, errors.New("aad is missing")
	}

	if proto.Payload == nil {
		return nil, errors.New("payload is missing")
	}

	// TODO:(nm) Validate topic

	return &ClientEnvelope{proto: proto}, nil
}

func NewClientEnvelopeFromBytes(bytes []byte) (*ClientEnvelope, error) {
	var message message_api.ClientEnvelope
	if err := proto.Unmarshal(bytes, &message); err != nil {
		return nil, err
	}
	return NewClientEnvelope(&message)
}

func (c *ClientEnvelope) ToBytes() ([]byte, error) {
	bytes, err := proto.Marshal(c.proto)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (c *ClientEnvelope) Aad() *message_api.AuthenticatedData {
	return c.proto.Aad
}

func (c *ClientEnvelope) Proto() *message_api.ClientEnvelope {
	return c.proto
}
