package storer

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/abis"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func buildIdentityUpdateStorer(t *testing.T) (*IdentityUpdateStorer, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	db, _, cleanup := testutils.NewDB(t, ctx)
	queryImpl := queries.New(db)
	config := testutils.GetContractsOptions(t)
	contractAddress := config.IdentityUpdatesContractAddress

	client, err := blockchain.NewClient(ctx, config.RpcUrl)
	require.NoError(t, err)
	contract, err := abis.NewIdentityUpdates(
		common.HexToAddress(contractAddress),
		client,
	)

	require.NoError(t, err)
	storer := NewIdentityUpdateStorer(queryImpl, testutils.NewLog(t), contract)

	return storer, func() {
		cancel()
		cleanup()
	}
}

func TestStoreIdentityUpdate(t *testing.T) {
	ctx := context.Background()
	storer, cleanup := buildIdentityUpdateStorer(t)
	defer cleanup()

	// Using the RandomInboxId function, since they are both 32 bytes and we treat inbox IDs as
	// strings outside the blockchain
	inboxId := testutils.RandomGroupID()
	message := testutils.RandomBytes(30)
	sequenceID := uint64(1)

	logMessage := testutils.BuildIdentityUpdateLog(t, inboxId, message, sequenceID)

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
