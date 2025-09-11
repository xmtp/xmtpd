// Package envelopes implements methods to create and manipulate different types of envelopes.
// Types implemented are ClientEnvelope, PayerEnvelope, and OriginatorEnvelope.
package envelopes

import (
	"errors"

	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/protobuf/proto"
)

type ClientEnvelope struct {
	proto       *envelopesProto.ClientEnvelope
	targetTopic topic.Topic
}

func NewClientEnvelope(proto *envelopesProto.ClientEnvelope) (*ClientEnvelope, error) {
	if proto == nil {
		return nil, errors.New("client envelope proto is nil")
	}

	if proto.Aad == nil {
		return nil, errors.New("aad is missing")
	}

	if proto.Payload == nil {
		return nil, errors.New("payload is missing")
	}

	targetTopic, err := topic.ParseTopic(proto.Aad.TargetTopic)
	if err != nil {
		return nil, err
	}

	return &ClientEnvelope{proto: proto, targetTopic: *targetTopic}, nil
}

func NewClientEnvelopeFromBytes(bytes []byte) (*ClientEnvelope, error) {
	message, err := utils.UnmarshalClientEnvelope(bytes)
	if err != nil {
		return nil, err
	}
	return NewClientEnvelope(message)
}

func (c *ClientEnvelope) Bytes() ([]byte, error) {
	bytes, err := proto.Marshal(c.proto)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (c *ClientEnvelope) TargetTopic() topic.Topic {
	return c.targetTopic
}

func (c *ClientEnvelope) Payload() interface{} {
	return c.proto.Payload
}

func (c *ClientEnvelope) Aad() *envelopesProto.AuthenticatedData {
	return c.proto.Aad
}

func (c *ClientEnvelope) Proto() *envelopesProto.ClientEnvelope {
	return c.proto
}

func (c *ClientEnvelope) TopicMatchesPayload() bool {
	targetTopic := c.TargetTopic()
	targetTopicKind := targetTopic.Kind()
	payload := c.proto.Payload

	switch payload.(type) {
	case *envelopesProto.ClientEnvelope_WelcomeMessage:
		return targetTopicKind == topic.TopicKindWelcomeMessagesV1
	case *envelopesProto.ClientEnvelope_GroupMessage:
		return targetTopicKind == topic.TopicKindGroupMessagesV1
	case *envelopesProto.ClientEnvelope_IdentityUpdate:
		return targetTopicKind == topic.TopicKindIdentityUpdatesV1
	case *envelopesProto.ClientEnvelope_UploadKeyPackage:
		return targetTopicKind == topic.TopicKindKeyPackagesV1
	case *envelopesProto.ClientEnvelope_PayerReport:
		return targetTopicKind == topic.TopicKindPayerReportsV1
	case *envelopesProto.ClientEnvelope_PayerReportAttestation:
		return targetTopicKind == topic.TopicKindPayerReportAttestationsV1
	default:
		return false
	}
}
