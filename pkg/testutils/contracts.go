package testutils

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/protobuf/proto"
)

const (
	LOCAL_PRIVATE_KEY = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
)

// Build an abi encoded MessageSent event struct
func BuildMessageSentEvent(
	message []byte,
) ([]byte, error) {
	gmabi, err := gm.GroupMessageBroadcasterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	inputs := gmabi.Events["MessageSent"].Inputs
	var nonIndexed abi.Arguments
	for _, input := range inputs {
		if !input.Indexed {
			nonIndexed = append(nonIndexed, input)
		}
	}

	return nonIndexed.Pack(message)
}

// BuildMessageSentLog builds a log message for a MessageSent event.
func BuildMessageSentLog(
	t *testing.T,
	groupID [16]byte,
	clientEnvelope *envelopesProto.ClientEnvelope,
	sequenceID uint64,
) types.Log {
	messageBytes, err := proto.Marshal(clientEnvelope)
	require.NoError(t, err)
	eventData, err := BuildMessageSentEvent(messageBytes)
	require.NoError(t, err)

	gmabi, err := gm.GroupMessageBroadcasterMetaData.GetAbi()
	require.NoError(t, err)

	topic0, err := utils.GetEventTopic(gmabi, "MessageSent")
	require.NoError(t, err)

	topic1 := common.BytesToHash(groupID[:])                       // indexed bytes16 groupID
	topic2 := common.BigToHash(new(big.Int).SetUint64(sequenceID)) // indexed uint64 sequenceID

	// Step 6: Assemble the log
	return types.Log{
		Topics: []common.Hash{
			topic0, // event signature
			topic1, // groupID
			topic2, // sequenceID
		},
		Data: eventData, // ABI-encoded `message` (non-indexed)
	}
}

func BuildIdentityUpdateEvent(update []byte) ([]byte, error) {
	iuabi, err := iu.IdentityUpdateBroadcasterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	inputs := iuabi.Events["IdentityUpdateCreated"].Inputs
	var nonIndexed abi.Arguments
	for _, input := range inputs {
		if !input.Indexed {
			nonIndexed = append(nonIndexed, input)
		}
	}

	return nonIndexed.Pack(update)
}

// BuildIdentityUpdateLog builds a log message for an IdentityUpdateCreated event.
func BuildIdentityUpdateLog(
	t *testing.T,
	inboxID [32]byte,
	clientEnvelope *envelopesProto.ClientEnvelope,
	sequenceID uint64,
) types.Log {
	messageBytes, err := proto.Marshal(clientEnvelope)
	require.NoError(t, err)
	eventData, err := BuildIdentityUpdateEvent(messageBytes)
	require.NoError(t, err)

	iuabi, err := iu.IdentityUpdateBroadcasterMetaData.GetAbi()
	require.NoError(t, err)

	topic0, err := utils.GetEventTopic(iuabi, "IdentityUpdateCreated")
	require.NoError(t, err)

	topic1 := common.BytesToHash(inboxID[:])                       // indexed bytes32 inboxID
	topic2 := common.BigToHash(new(big.Int).SetUint64(sequenceID)) // indexed uint64 sequenceID

	// Step 6: Assemble the log
	return types.Log{
		Topics: []common.Hash{
			topic0, // event signature
			topic1, // inboxID
			topic2, // sequenceID
		},
		Data: eventData, // ABI-encoded `message` (non-indexed)
	}
}
