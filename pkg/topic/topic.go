package topic

import (
	"errors"
	"fmt"
)

type TopicKind uint8

const (
	TOPIC_KIND_GROUP_MESSAGES_V1 TopicKind = iota
	TOPIC_KIND_WELCOME_MESSAGES_V1
	TOPIC_KIND_IDENTITY_UPDATES_V1
	TOPIC_KIND_KEY_PACKAGES_V1
)

func (k TopicKind) String() string {
	switch k {
	case TOPIC_KIND_GROUP_MESSAGES_V1:
		return "group_messages_v1"
	case TOPIC_KIND_WELCOME_MESSAGES_V1:
		return "welcome_message_v1"
	case TOPIC_KIND_IDENTITY_UPDATES_V1:
		return "identity_updates_v1"
	case TOPIC_KIND_KEY_PACKAGES_V1:
		return "key_packages_v1"
	default:
		return "unknown"
	}
}

type Topic struct {
	kind       TopicKind
	identifier []byte
}

func NewTopic(kind TopicKind, identifier []byte) *Topic {
	return &Topic{
		kind:       kind,
		identifier: identifier,
	}
}

func (t Topic) Bytes() []byte {
	result := make([]byte, 1+len(t.identifier))
	result[0] = byte(t.kind)
	copy(result[1:], t.identifier)
	return result
}

func (t Topic) String() string {
	return fmt.Sprintf("%s/%x", t.kind.String(), t.identifier)
}

func ParseTopic(topic []byte) (*Topic, error) {
	if len(topic) < 2 {
		return nil, errors.New("topic must be at least 2 bytes long")
	}

	kind := TopicKind(topic[0])
	identifier := topic[1:]

	newTopic := NewTopic(kind, identifier)

	if newTopic.Kind().String() == "unknown" {
		return nil, fmt.Errorf("unknown topic kind %d", kind)
	}

	return newTopic, nil
}

func (t Topic) Kind() TopicKind {
	return t.kind
}

func (t Topic) Identifier() []byte {
	return t.identifier
}
