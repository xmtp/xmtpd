package topic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/utils"
)

func TestValidTopic(t *testing.T) {
	newTopic := []byte{1, 2, 3}
	topic, err := ParseTopic(newTopic)
	require.NoError(t, err)
	require.Equal(t, TOPIC_KIND_WELCOME_MESSAGES_V1, topic.Kind())
	require.Equal(t, []byte{2, 3}, topic.Identifier())
}

func TestMissingIdentifier(t *testing.T) {
	newTopic := []byte{1}
	topic, err := ParseTopic(newTopic)
	require.Error(t, err)
	require.Nil(t, topic)
}

func TestInvalidKind(t *testing.T) {
	newTopic := []byte{255, 2, 3}
	topic, err := ParseTopic(newTopic)
	require.Error(t, err)
	require.Nil(t, topic)
}

func TestTopicString(t *testing.T) {
	identifier := testutils.RandomBytes(32)

	groupMessagesTopic := NewTopic(TOPIC_KIND_GROUP_MESSAGES_V1, identifier)
	require.Equal(t, "group_messages_v1", groupMessagesTopic.Kind().String())
	require.Equal(t, identifier, groupMessagesTopic.Identifier())
	require.Equal(
		t,
		"group_messages_v1/"+utils.HexEncode(identifier),
		groupMessagesTopic.String(),
	)

	identityUpdatesTopic := NewTopic(TOPIC_KIND_IDENTITY_UPDATES_V1, identifier)
	require.Equal(t, "identity_updates_v1", identityUpdatesTopic.Kind().String())
	require.Equal(t, identifier, identityUpdatesTopic.Identifier())
	require.Equal(
		t,
		"identity_updates_v1/"+utils.HexEncode(identifier),
		identityUpdatesTopic.String(),
	)
}
