// Package topic implements the topic type.
package topic

import (
	"errors"
	"fmt"
)

type TopicKind uint8

const (
	TopicKindGroupMessagesV1 TopicKind = iota
	TopicKindWelcomeMessagesV1
	TopicKindIdentityUpdatesV1
	TopicKindKeyPackagesV1
	TopicKindPayerReportsV1
	TopicKindPayerReportAttestationsV1
)

func (k TopicKind) String() string {
	switch k {
	case TopicKindGroupMessagesV1:
		return "group_messages_v1"
	case TopicKindWelcomeMessagesV1:
		return "welcome_message_v1"
	case TopicKindIdentityUpdatesV1:
		return "identity_updates_v1"
	case TopicKindKeyPackagesV1:
		return "key_packages_v1"
	case TopicKindPayerReportsV1:
		return "payer_reports_v1"
	case TopicKindPayerReportAttestationsV1:
		return "payer_report_attestations_v1"
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

// IsReserved topics can only be published to by the node itself, and not through Payers.
func (t Topic) IsReserved() bool {
	return t.kind == TopicKindPayerReportsV1 ||
		t.kind == TopicKindPayerReportAttestationsV1
}
