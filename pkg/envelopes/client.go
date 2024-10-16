package envelopes

import (
	"errors"

	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/topic"
	"google.golang.org/protobuf/proto"
)

type ClientEnvelope struct {
	proto       *message_api.ClientEnvelope
	targetTopic topic.Topic
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

	targetTopic, err := topic.ParseTopic(proto.Aad.TargetTopic)
	if err != nil {
		return nil, err
	}

	return &ClientEnvelope{proto: proto, targetTopic: *targetTopic}, nil
}

func NewClientEnvelopeFromBytes(bytes []byte) (*ClientEnvelope, error) {
	var message message_api.ClientEnvelope
	if err := proto.Unmarshal(bytes, &message); err != nil {
		return nil, err
	}
	return NewClientEnvelope(&message)
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

func (c *ClientEnvelope) Aad() *message_api.AuthenticatedData {
	return c.proto.Aad
}

func (c *ClientEnvelope) Proto() *message_api.ClientEnvelope {
	return c.proto
}

func (c *ClientEnvelope) TopicMatchesPayload() bool {
	targetTopic := c.TargetTopic()
	targetTopicKind := targetTopic.Kind()
	payload := c.proto.Payload

	switch payload.(type) {
	case *message_api.ClientEnvelope_WelcomeMessage:
		return targetTopicKind == topic.TOPIC_KIND_WELCOME_MESSAGES_V1
	case *message_api.ClientEnvelope_GroupMessage:
		return targetTopicKind == topic.TOPIC_KIND_GROUP_MESSAGES_V1
	case *message_api.ClientEnvelope_IdentityUpdate:
		return targetTopicKind == topic.TOPIC_KIND_IDENTITY_UPDATES_V1
	case *message_api.ClientEnvelope_UploadKeyPackage:
		return targetTopicKind == topic.TOPIC_KIND_KEY_PACKAGES_V1
	default:
		return false
	}
}
