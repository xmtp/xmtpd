package storer

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/contracts/pkg/groupmessages"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
	envelopesTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
)

func buildGroupMessageStorer(t *testing.T) (*GroupMessageStorer, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	db, _, cleanup := testutils.NewDB(t, ctx)
	queryImpl := queries.New(db)
	rpcUrl, anvilCleanup := anvil.StartAnvil(t, false)
	config := testutils.NewContractsOptions(rpcUrl)
	config.MessagesContractAddress = testutils.DeployGroupMessagesContract(t, rpcUrl)

	client, err := blockchain.NewClient(ctx, config.RpcUrl)
	require.NoError(t, err)
	contract, err := groupmessages.NewGroupMessages(
		common.HexToAddress(config.MessagesContractAddress),
		client,
	)

	require.NoError(t, err)
	storer := NewGroupMessageStorer(queryImpl, testutils.NewLog(t), contract)

	return storer, func() {
		defer anvilCleanup()
		cancel()
		cleanup()
	}
}

func TestStoreGroupMessages(t *testing.T) {
	ctx := context.Background()
	storer, cleanup := buildGroupMessageStorer(t)
	defer cleanup()

	groupID := testutils.RandomGroupID()
	message := testutils.RandomBytes(30)
	sequenceID := uint64(1)

	clientEnvelope := envelopesTestUtils.CreateGroupMessageClientEnvelope(groupID, message)

	logMessage := testutils.BuildMessageSentLog(t, groupID, clientEnvelope, sequenceID)
	var err error
	err = storer.StoreLog(
		ctx,
		logMessage,
	)
	require.NoError(t, err)

	gatewayEnvelopes, err := storer.queries.SelectGatewayEnvelopes(
		ctx,
		queries.SelectGatewayEnvelopesParams{OriginatorNodeIds: []int32{0}},
	)
	require.NoError(t, err)

	require.Equal(t, len(gatewayEnvelopes), 1)

	firstEnvelope := gatewayEnvelopes[0]
	originatorEnvelope, err := envelopes.NewOriginatorEnvelopeFromBytes(
		firstEnvelope.OriginatorEnvelope,
	)
	require.NoError(t, err)
	require.True(
		t,
		originatorEnvelope.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope.TopicMatchesPayload(),
	)
	targetTopic := originatorEnvelope.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope.TargetTopic()
	require.Equal(t, targetTopic.Kind(), topic.TOPIC_KIND_GROUP_MESSAGES_V1)
	require.Equal(t, targetTopic.Identifier(), groupID[:])
	require.Equal(t, firstEnvelope.OriginatorSequenceID, int64(sequenceID))
}

func TestStoreGroupMessageDuplicate(t *testing.T) {
	ctx := context.Background()
	storer, cleanup := buildGroupMessageStorer(t)
	defer cleanup()

	var groupID [32]byte
	copy(groupID[:], testutils.RandomBytes(32))
	message := testutils.RandomBytes(30)
	sequenceID := uint64(1)

	clientEnvelope := envelopesTestUtils.CreateGroupMessageClientEnvelope(groupID, message)

	logMessage := testutils.BuildMessageSentLog(t, groupID, clientEnvelope, sequenceID)

	err := storer.StoreLog(
		ctx,
		logMessage,
	)
	require.NoError(t, err)
	// Store the log a second time
	err = storer.StoreLog(
		ctx,
		logMessage,
	)
	require.NoError(t, err)

	envelopes, queryErr := storer.queries.SelectGatewayEnvelopes(
		ctx,
		queries.SelectGatewayEnvelopesParams{OriginatorNodeIds: []int32{0}},
	)
	require.NoError(t, queryErr)

	require.Equal(t, len(envelopes), 1)
}

func TestStoreGroupMessageMalformed(t *testing.T) {
	ctx := context.Background()
	storer, cleanup := buildGroupMessageStorer(t)
	defer cleanup()

	abi, err := groupmessages.GroupMessagesMetaData.GetAbi()
	require.NoError(t, err)

	topic, err := utils.GetEventTopic(abi, "MessageSent")
	require.NoError(t, err)

	logMessage := types.Log{
		Topics: []common.Hash{topic},
		Data:   []byte("foo"),
	}

	storageErr := storer.StoreLog(ctx, logMessage)
	require.Error(t, storageErr)
	require.False(t, storageErr.ShouldRetry())
}
