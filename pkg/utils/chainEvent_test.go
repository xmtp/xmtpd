package utils

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/abis"
)

const (
	// Update this if event changes
	EXPECTED_MESSAGE_SENT_SIGNATURE = "MessageSent(bytes32,bytes,uint64)"
)

func TestGetEventSignature(t *testing.T) {
	abi, _ := abis.GroupMessagesMetaData.GetAbi()

	signature, err := GetEventSig(abi, "MessageSent")
	require.NoError(t, err)
	require.Equal(t, signature, EXPECTED_MESSAGE_SENT_SIGNATURE)
}

func TestGetEventTopic(t *testing.T) {
	abi, _ := abis.GroupMessagesMetaData.GetAbi()

	topic, err := GetEventTopic(abi, "MessageSent")
	require.NoError(t, err)
	require.Equal(t, topic, crypto.Keccak256Hash([]byte(EXPECTED_MESSAGE_SENT_SIGNATURE)))
}
