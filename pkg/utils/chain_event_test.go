package utils

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
)

const (
	MessageSentSignature = "MessageSent(bytes16,bytes,uint64)"
)

func TestGetEventSignature(t *testing.T) {
	abi, _ := gm.GroupMessageBroadcasterMetaData.GetAbi()

	signature, err := GetEventSig(abi, "MessageSent")
	require.NoError(t, err)
	require.Equal(t, MessageSentSignature, signature)
}

func TestGetEventTopic(t *testing.T) {
	abi, _ := gm.GroupMessageBroadcasterMetaData.GetAbi()

	topic, err := GetEventTopic(abi, "MessageSent")
	require.NoError(t, err)
	require.Equal(t, topic, crypto.Keccak256Hash([]byte(MessageSentSignature)))
}
