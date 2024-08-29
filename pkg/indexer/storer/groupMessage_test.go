package storer

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/abis"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/indexer/blockchain"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/utils"
)

func buildGroupMessageStorer(t *testing.T) (*GroupMessageStorer, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	db, _, cleanup := testutils.NewDB(t, ctx)
	queryImpl := queries.New(db)
	config := testutils.GetContractsOptions(t)
	contractAddress := config.MessagesContractAddress

	client, err := blockchain.NewClient(ctx, config.RpcUrl)
	require.NoError(t, err)
	contract, err := abis.NewGroupMessages(
		common.HexToAddress(contractAddress),
		client,
	)

	require.NoError(t, err)
	storer := NewGroupMessageStorer(queryImpl, testutils.NewLog(t), contract)

	return storer, func() {
		cancel()
		cleanup()
	}
}

func TestStoreGroupMessages(t *testing.T) {
	ctx := context.Background()
	storer, cleanup := buildGroupMessageStorer(t)
	defer cleanup()

	var groupID [32]byte
	copy(groupID[:], testutils.RandomBytes(32))
	message := testutils.RandomBytes(30)
	sequenceID := uint64(1)

	logMessage := testutils.BuildMessageSentLog(t, groupID, message, sequenceID)

	err := storer.StoreLog(
		ctx,
		logMessage,
	)
	require.NoError(t, err)

	envelopes, queryErr := storer.queries.SelectGatewayEnvelopes(
		ctx,
		queries.SelectGatewayEnvelopesParams{OriginatorNodeID: db.NullInt32(0)},
	)
	require.NoError(t, queryErr)

	require.Equal(t, len(envelopes), 1)

	firstEnvelope := envelopes[0]
	require.Equal(t, firstEnvelope.OriginatorEnvelope, message)
}

func TestStoreGroupMessageDuplicate(t *testing.T) {
	ctx := context.Background()
	storer, cleanup := buildGroupMessageStorer(t)
	defer cleanup()

	var groupID [32]byte
	copy(groupID[:], testutils.RandomBytes(32))
	message := testutils.RandomBytes(30)
	sequenceID := uint64(1)

	logMessage := testutils.BuildMessageSentLog(t, groupID, message, sequenceID)

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
		queries.SelectGatewayEnvelopesParams{OriginatorNodeID: db.NullInt32(0)},
	)
	require.NoError(t, queryErr)

	require.Equal(t, len(envelopes), 1)
}

func TestStoreGroupMessageMalformed(t *testing.T) {
	ctx := context.Background()
	storer, cleanup := buildGroupMessageStorer(t)
	defer cleanup()

	abi, err := abis.GroupMessagesMetaData.GetAbi()
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
