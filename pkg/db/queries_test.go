package db_test

import (
	"context"
	"testing"

	xmtpd_db "github.com/xmtp/xmtpd/pkg/db"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func getAddressLogState(
	t *testing.T,
	querier *queries.Queries,
	address string,
	inboxID string,
) *queries.GetAddressLogsRow {
	addressLogs, err := querier.GetAddressLogs(context.Background(), []string{address})
	require.NoError(t, err)

	if len(addressLogs) == 0 {
		return nil
	}

	addressLog := addressLogs[0]
	require.Equal(t, addressLog.InboxID, inboxID)

	return &addressLog
}

func TestInsertAddressLog(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewDB(t, ctx)

	querier := queries.New(db)
	address := testutils.RandomString(20)
	inboxID := testutils.RandomInboxIDString()

	_, err := querier.InsertAddressLog(
		ctx,
		queries.InsertAddressLogParams{
			Address:               address,
			InboxID:               inboxID,
			AssociationSequenceID: xmtpd_db.NullInt64(1),
		},
	)
	require.NoError(t, err)

	addressLog := getAddressLogState(t, querier, address, inboxID)
	require.NotNil(t, addressLog)
	require.Equal(t, addressLog.AssociationSequenceID.Int64, int64(1))

	// Now insert a new entry with a higher sequence id
	_, err = querier.InsertAddressLog(
		ctx,
		queries.InsertAddressLogParams{
			Address:               address,
			InboxID:               inboxID,
			AssociationSequenceID: xmtpd_db.NullInt64(2),
		},
	)
	require.NoError(t, err)

	addressLog = getAddressLogState(t, querier, address, inboxID)
	require.NotNil(t, addressLog)
	require.Equal(t, addressLog.AssociationSequenceID.Int64, int64(2))

	// Try to set it back to 1. This should be a no-op
	numRows, err := querier.InsertAddressLog(
		ctx,
		queries.InsertAddressLogParams{
			Address:               address,
			InboxID:               inboxID,
			AssociationSequenceID: xmtpd_db.NullInt64(1),
		},
	)
	require.NoError(t, err)
	require.Equal(t, numRows, int64(0))

	addressLog = getAddressLogState(t, querier, address, inboxID)
	require.NotNil(t, addressLog)
	require.Equal(t, addressLog.AssociationSequenceID.Int64, int64(2))
}

func TestRevokeAddressLog(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewDB(t, ctx)

	querier := queries.New(db)

	address := testutils.RandomString(20)
	inboxID := testutils.RandomInboxIDString()

	_, err := querier.InsertAddressLog(
		ctx,
		queries.InsertAddressLogParams{
			Address:               address,
			InboxID:               inboxID,
			AssociationSequenceID: xmtpd_db.NullInt64(1),
		},
	)
	require.NoError(t, err)

	numRows, err := querier.RevokeAddressFromLog(
		ctx,
		queries.RevokeAddressFromLogParams{
			Address:              address,
			InboxID:              inboxID,
			RevocationSequenceID: xmtpd_db.NullInt64(2),
		},
	)
	require.NoError(t, err)
	require.Equal(t, numRows, int64(1))

	addressLog := getAddressLogState(t, querier, address, inboxID)
	require.Nil(t, addressLog)

	// Now try to associate it a second time

	numRows, err = querier.InsertAddressLog(
		ctx,
		queries.InsertAddressLogParams{
			Address:               address,
			InboxID:               inboxID,
			AssociationSequenceID: xmtpd_db.NullInt64(3),
		},
	)
	require.NoError(t, err)
	require.Equal(t, numRows, int64(1))

	addressLog = getAddressLogState(t, querier, address, inboxID)
	require.NotNil(t, addressLog)
	require.Equal(t, addressLog.AssociationSequenceID.Int64, int64(3))
}

func TestFindOrCreatePayer(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewDB(t, ctx)

	querier := queries.New(db)

	address1 := testutils.RandomString(42)
	address2 := testutils.RandomString(42)

	id1, err := querier.FindOrCreatePayer(ctx, address1)
	require.NoError(t, err)

	id2, err := querier.FindOrCreatePayer(ctx, address2)
	require.NoError(t, err)

	require.NotEqual(t, id1, id2)

	reinsertID, err := querier.FindOrCreatePayer(ctx, address1)
	require.NoError(t, err)
	require.Equal(t, id1, reinsertID)
}
