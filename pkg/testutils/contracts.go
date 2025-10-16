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
	pr "github.com/xmtp/xmtpd/pkg/abi/payerregistry"
	prm "github.com/xmtp/xmtpd/pkg/abi/payerreportmanager"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/protobuf/proto"
)

// BuildMessageSentEvent builds an abi encoded MessageSent event struct.
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

func BuildPayerReportSubmittedEvent(
	t *testing.T,
	originatorNodeID uint32,
	payerReportIndex uint64,
	startSequenceID uint64,
	endSequenceID uint64,
	endMinuteSinceEpoch uint64,
	payersMerkleRoot [32]byte,
	activeNodeIDs []uint32,
) types.Log {
	prmabi, err := prm.PayerReportManagerMetaData.GetAbi()
	require.NoError(t, err)

	topic0, err := utils.GetEventTopic(prmabi, "PayerReportSubmitted")
	require.NoError(t, err)

	inputs := prmabi.Events["PayerReportSubmitted"].Inputs
	var nonIndexed abi.Arguments
	for _, input := range inputs {
		if !input.Indexed {
			nonIndexed = append(nonIndexed, input)
		}
	}

	originatorNodeIDHash := common.BigToHash(big.NewInt(int64(originatorNodeID)))
	payerReportIndexHash := common.BigToHash(big.NewInt(int64(payerReportIndex)))
	endSequenceIDHash := common.BigToHash(big.NewInt(int64(endSequenceID)))

	data, err := nonIndexed.Pack(
		startSequenceID,
		uint32(endMinuteSinceEpoch),
		payersMerkleRoot,
		activeNodeIDs,
		activeNodeIDs,
	)
	require.NoError(t, err)

	return types.Log{
		Topics: []common.Hash{
			topic0,
			originatorNodeIDHash,
			payerReportIndexHash,
			endSequenceIDHash,
		},
		Data: data,
	}
}

func BuildPayerRegistryWithdrawalRequestedLog(
	t *testing.T,
	payer common.Address,
	amount *big.Int,
	withdrawableTimestamp uint32,
	nonce uint32,
) types.Log {
	require.Less(t, nonce, uint32(1<<24), "nonce must fit in 24 bits")

	prabi, err := pr.PayerRegistryMetaData.GetAbi()
	require.NoError(t, err)

	topic0, err := utils.GetEventTopic(prabi, "WithdrawalRequested")
	require.NoError(t, err)

	inputs := prabi.Events["WithdrawalRequested"].Inputs
	var nonIndexed abi.Arguments
	for _, input := range inputs {
		if !input.Indexed {
			nonIndexed = append(nonIndexed, input)
		}
	}

	topic1 := common.BytesToHash(payer.Bytes()) // indexed address payer

	data, err := nonIndexed.Pack(amount, withdrawableTimestamp, big.NewInt(int64(nonce)))
	require.NoError(t, err)

	return types.Log{
		Topics: []common.Hash{
			topic0, // event signature
			topic1, // payer
		},
		Data: data,
	}
}

func BuildPayerRegistryDepositLog(
	t *testing.T,
	payer common.Address,
	amount *big.Int,
) types.Log {
	prabi, err := pr.PayerRegistryMetaData.GetAbi()
	require.NoError(t, err)

	topic0, err := utils.GetEventTopic(prabi, "Deposit")
	require.NoError(t, err)

	inputs := prabi.Events["Deposit"].Inputs
	var nonIndexed abi.Arguments
	for _, input := range inputs {
		if !input.Indexed {
			nonIndexed = append(nonIndexed, input)
		}
	}

	topic1 := common.BytesToHash(payer.Bytes()) // indexed address payer

	data, err := nonIndexed.Pack(amount)
	require.NoError(t, err)

	return types.Log{
		Topics: []common.Hash{
			topic0, // event signature
			topic1,
		},
		Data: data,
	}
}

func BuildPayerRegistryUsageSettledLog(
	t *testing.T,
	payer common.Address,
	amount *big.Int,
	payerReportId common.Hash,
) types.Log {
	prabi, err := pr.PayerRegistryMetaData.GetAbi()
	require.NoError(t, err)

	topic0, err := utils.GetEventTopic(prabi, "UsageSettled")
	require.NoError(t, err)

	inputs := prabi.Events["UsageSettled"].Inputs
	var nonIndexed abi.Arguments
	for _, input := range inputs {
		if !input.Indexed {
			nonIndexed = append(nonIndexed, input)
		}
	}

	topic1 := payerReportId
	topic2 := common.BytesToHash(payer.Bytes()) // indexed address payer

	data, err := nonIndexed.Pack(amount)
	require.NoError(t, err)

	return types.Log{
		Topics: []common.Hash{
			topic0, // event signature
			topic1, // payer report ID
			topic2, // payer
		},
		Data: data,
	}
}

func BuildPayerRegistryWithdrawalCancelledLog(
	t *testing.T,
	payer common.Address,
) types.Log {
	prabi, err := pr.PayerRegistryMetaData.GetAbi()
	require.NoError(t, err)

	topic0, err := utils.GetEventTopic(prabi, "WithdrawalCancelled")
	require.NoError(t, err)

	topic1 := common.BytesToHash(payer.Bytes()) // indexed address payer

	return types.Log{
		Topics: []common.Hash{
			topic0, // event signature
			topic1, // payer
		},
		Data: []byte{}, // No non-indexed data for this event
	}
}

func BuildPayerReportSubsetSettledLog(
	t *testing.T,
	originatorNodeID uint32,
	payerReportIndex uint64,
	count uint32,
	remaining uint32,
	feesSettled *big.Int,
) types.Log {
	prmabi, err := prm.PayerReportManagerMetaData.GetAbi()
	require.NoError(t, err)

	topic0, err := utils.GetEventTopic(prmabi, "PayerReportSubsetSettled")
	require.NoError(t, err)

	inputs := prmabi.Events["PayerReportSubsetSettled"].Inputs
	var nonIndexed abi.Arguments
	for _, input := range inputs {
		if !input.Indexed {
			nonIndexed = append(nonIndexed, input)
		}
	}

	originatorNodeIDHash := common.BigToHash(big.NewInt(int64(originatorNodeID)))
	payerReportIndexHash := common.BigToHash(big.NewInt(int64(payerReportIndex)))

	data, err := nonIndexed.Pack(count, remaining, feesSettled)
	require.NoError(t, err)

	return types.Log{
		Topics: []common.Hash{
			topic0,
			originatorNodeIDHash,
			payerReportIndexHash,
		},
		Data: data,
	}
}
