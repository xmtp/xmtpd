package testutils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/protobuf/proto"
	"testing"
)

const (
	LOCAL_PRIVATE_KEY = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
)

// Build an abi encoded MessageSent event struct
func BuildMessageSentEvent(
	groupID [32]byte,
	message []byte,
	sequenceID uint64,
) ([]byte, error) {
	abi, err := gm.GroupMessageBroadcasterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return abi.Events["MessageSent"].Inputs.Pack(groupID, message, sequenceID)
}

// Build a log message for a MessageSent event
func BuildMessageSentLog(
	t *testing.T,
	groupID [32]byte,
	clientEnvelope *envelopesProto.ClientEnvelope,
	sequenceID uint64,
) types.Log {
	messageBytes, err := proto.Marshal(clientEnvelope)
	require.NoError(t, err)
	eventData, err := BuildMessageSentEvent(groupID, messageBytes, sequenceID)
	require.NoError(t, err)

	abi, err := gm.GroupMessageBroadcasterMetaData.GetAbi()
	require.NoError(t, err)

	topic, err := utils.GetEventTopic(abi, "MessageSent")
	require.NoError(t, err)

	return types.Log{
		Topics: []common.Hash{topic},
		Data:   eventData,
	}
}

func BuildIdentityUpdateEvent(inboxId [32]byte, update []byte, sequenceID uint64) ([]byte, error) {
	abi, err := iu.IdentityUpdateBroadcasterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return abi.Events["IdentityUpdateCreated"].Inputs.Pack(inboxId, update, sequenceID)
}

// Build a log message for an IdentityUpdateCreated event
func BuildIdentityUpdateLog(
	t *testing.T,
	inboxId [32]byte,
	clientEnvelope *envelopesProto.ClientEnvelope,
	sequenceID uint64,
) types.Log {
	messageBytes, err := proto.Marshal(clientEnvelope)
	require.NoError(t, err)
	eventData, err := BuildIdentityUpdateEvent(inboxId, messageBytes, sequenceID)
	require.NoError(t, err)

	abi, err := iu.IdentityUpdateBroadcasterMetaData.GetAbi()
	require.NoError(t, err)

	topic, err := utils.GetEventTopic(abi, "IdentityUpdateCreated")
	require.NoError(t, err)

	return types.Log{
		Topics: []common.Hash{topic},
		Data:   eventData,
	}
}
