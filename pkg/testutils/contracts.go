package testutils

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/contracts/pkg/groupmessages"
	"github.com/xmtp/xmtpd/contracts/pkg/identityupdates"
	"github.com/xmtp/xmtpd/contracts/pkg/nodes"
	"github.com/xmtp/xmtpd/contracts/pkg/nodesv2"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/protobuf/proto"
)

const (
	ANVIL_LOCALNET_HOST            = "http://localhost:7545"
	ANVIL_LOCALNET_CHAIN_ID        = 31337
	LOCAL_PRIVATE_KEY              = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	NODES_CONTRACT_NAME            = "Nodes"
	NODES_V2_CONTRACT_NAME         = "NodesV2"
	GROUP_MESSAGES_CONTRACT_NAME   = "GroupMessages"
	IDENTITY_UPDATES_CONTRACT_NAME = "IdentityUpdates"
)

// Build an abi encoded MessageSent event struct
func BuildMessageSentEvent(
	groupID [32]byte,
	message []byte,
	sequenceID uint64,
) ([]byte, error) {
	abi, err := groupmessages.GroupMessagesMetaData.GetAbi()
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

	abi, err := groupmessages.GroupMessagesMetaData.GetAbi()
	require.NoError(t, err)

	topic, err := utils.GetEventTopic(abi, "MessageSent")
	require.NoError(t, err)

	return types.Log{
		Topics: []common.Hash{topic},
		Data:   eventData,
	}
}

func BuildIdentityUpdateEvent(inboxId [32]byte, update []byte, sequenceID uint64) ([]byte, error) {
	abi, err := identityupdates.IdentityUpdatesMetaData.GetAbi()
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

	abi, err := identityupdates.IdentityUpdatesMetaData.GetAbi()
	require.NoError(t, err)

	topic, err := utils.GetEventTopic(abi, "IdentityUpdateCreated")
	require.NoError(t, err)

	return types.Log{
		Topics: []common.Hash{topic},
		Data:   eventData,
	}
}

/*
*
Deploy a contract and return the contract's address. Will return a different address for each run, making it suitable for testing
*
*/
func deployContract(t *testing.T, contractName string) string {
	retryMax := 10
	var retry = 0
	var err error

	var addr common.Address

	for retry < retryMax {
		if err != nil {
			retry++
			t.Logf("Error deploying contract, retrying: %v. Attempt (%d/%d)", err, retry, retryMax)
			time.Sleep(10 * time.Millisecond)
		}

		var client *ethclient.Client
		client, err = ethclient.Dial(ANVIL_LOCALNET_HOST)
		if err != nil {
			continue
		}

		privateKey, err := crypto.HexToECDSA(LOCAL_PRIVATE_KEY)
		if err != nil {
			continue
		}

		auth, err := bind.NewKeyedTransactorWithChainID(
			privateKey,
			big.NewInt(ANVIL_LOCALNET_CHAIN_ID),
		)
		if err != nil {
			continue
		}

		switch contractName {
		case NODES_CONTRACT_NAME:
			addr, _, _, err = nodes.DeployNodes(auth, client)
		case NODES_V2_CONTRACT_NAME:
			addr, _, _, err = nodesv2.DeployNodesV2(auth, client, auth.From)
		case GROUP_MESSAGES_CONTRACT_NAME:
			addr, _, _, err = groupmessages.DeployGroupMessages(auth, client)
		case IDENTITY_UPDATES_CONTRACT_NAME:
			addr, _, _, err = identityupdates.DeployIdentityUpdates(auth, client)
		default:
			t.Fatalf("Unknown contract name: %s", contractName)
		}

		if err != nil {
			continue
		}

		break

	}

	require.NoError(t, err)
	return addr.String()

}

func DeployNodesContract(t *testing.T) string {
	return deployContract(t, NODES_CONTRACT_NAME)
}

func DeployNodesV2Contract(t *testing.T) string {
	return deployContract(t, NODES_V2_CONTRACT_NAME)
}

func DeployGroupMessagesContract(t *testing.T) string {
	return deployContract(t, GROUP_MESSAGES_CONTRACT_NAME)
}

func DeployIdentityUpdatesContract(t *testing.T) string {
	return deployContract(t, IDENTITY_UPDATES_CONTRACT_NAME)
}
