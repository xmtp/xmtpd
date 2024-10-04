package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func getAddressLogState(
	t *testing.T,
	querier *queries.Queries,
	address string,
	inboxId string,
) *queries.GetAddressLogsRow {
	addressLogs, err := querier.GetAddressLogs(context.Background(), []string{address})
	require.NoError(t, err)

	if len(addressLogs) == 0 {
		return nil
	}

	addressLog := addressLogs[0]
	require.Equal(t, addressLog.InboxID, inboxId)

	return &addressLog
}

func TestInsertAddressLog(t *testing.T) {
	ctx := context.Background()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	querier := queries.New(db)

	address := testutils.RandomString(20)
	inboxId := testutils.RandomInboxId()

	_, err := querier.InsertAddressLog(
		ctx,
		queries.InsertAddressLogParams{
			Address:               address,
			InboxID:               inboxId,
			AssociationSequenceID: NullInt64(1),
		},
	)
	require.NoError(t, err)

	addressLog := getAddressLogState(t, querier, address, inboxId)
	require.NotNil(t, addressLog)
	require.Equal(t, addressLog.AssociationSequenceID.Int64, int64(1))

	// Now insert a new entry with a higher sequence id
	_, err = querier.InsertAddressLog(
		ctx,
		queries.InsertAddressLogParams{
			Address:               address,
			InboxID:               inboxId,
			AssociationSequenceID: NullInt64(2),
		},
	)
	require.NoError(t, err)

	addressLog = getAddressLogState(t, querier, address, inboxId)
	require.NotNil(t, addressLog)
	require.Equal(t, addressLog.AssociationSequenceID.Int64, int64(2))

	// Try to set it back to 1. This should be a no-op
	numRows, err := querier.InsertAddressLog(
		ctx,
		queries.InsertAddressLogParams{
			Address:               address,
			InboxID:               inboxId,
			AssociationSequenceID: NullInt64(1),
		},
	)
	require.NoError(t, err)
	require.Equal(t, numRows, int64(0))

	addressLog = getAddressLogState(t, querier, address, inboxId)
	require.NotNil(t, addressLog)
	require.Equal(t, addressLog.AssociationSequenceID.Int64, int64(2))
}

func TestRevokeAddressLog(t *testing.T) {
	ctx := context.Background()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	querier := queries.New(db)

	address := testutils.RandomString(20)
	inboxId := testutils.RandomInboxId()

	_, err := querier.InsertAddressLog(
		ctx,
		queries.InsertAddressLogParams{
			Address:               address,
			InboxID:               inboxId,
			AssociationSequenceID: NullInt64(1),
		},
	)
	require.NoError(t, err)

	numRows, err := querier.RevokeAddressFromLog(
		ctx,
		queries.RevokeAddressFromLogParams{
			Address:              address,
			InboxID:              inboxId,
			RevocationSequenceID: NullInt64(2),
		},
	)
	require.NoError(t, err)
	require.Equal(t, numRows, int64(1))

	addressLog := getAddressLogState(t, querier, address, inboxId)
	require.Nil(t, addressLog)

	// Now try to associate it a second time

	numRows, err = querier.InsertAddressLog(
		ctx,
		queries.InsertAddressLogParams{
			Address:               address,
			InboxID:               inboxId,
			AssociationSequenceID: NullInt64(3),
		},
	)
	require.NoError(t, err)
	require.Equal(t, numRows, int64(1))

	addressLog = getAddressLogState(t, querier, address, inboxId)
	require.NotNil(t, addressLog)
	require.Equal(t, addressLog.AssociationSequenceID.Int64, int64(3))
}
