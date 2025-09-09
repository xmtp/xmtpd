package db_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/utils"
)

func TestIncrementUnsettledUsage(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewDB(t, ctx)

	querier := queries.New(db)
	payerID := testutils.RandomInt32()
	originatorID := testutils.RandomInt32()
	minutesSinceEpoch := utils.MinutesSinceEpochNow()

	require.NoError(t, querier.IncrementUnsettledUsage(ctx, queries.IncrementUnsettledUsageParams{
		PayerID:           payerID,
		OriginatorID:      originatorID,
		MinutesSinceEpoch: minutesSinceEpoch,
		SpendPicodollars:  100,
	}))

	unsettledUsage, err := querier.GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{
			PayerID: payerID,
		},
	)
	require.NoError(t, err)
	require.Equal(t, unsettledUsage.TotalSpendPicodollars, int64(100))

	require.NoError(t, querier.IncrementUnsettledUsage(ctx, queries.IncrementUnsettledUsageParams{
		PayerID:           payerID,
		OriginatorID:      originatorID,
		MinutesSinceEpoch: minutesSinceEpoch,
		SpendPicodollars:  100,
	}))

	unsettledUsage, err = querier.GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{
			PayerID: payerID,
		},
	)
	require.NoError(t, err)
	require.Equal(t, unsettledUsage.TotalSpendPicodollars, int64(200))
}

func TestGetUnsettledUsage(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewDB(t, ctx)

	querier := queries.New(db)
	payerID := testutils.RandomInt32()
	originatorID := testutils.RandomInt32()

	addUsage := func(minutesSinceEpoch int32, spendPicodollars int64) {
		require.NoError(
			t,
			querier.IncrementUnsettledUsage(ctx, queries.IncrementUnsettledUsageParams{
				PayerID:           payerID,
				OriginatorID:      originatorID,
				MinutesSinceEpoch: minutesSinceEpoch,
				SpendPicodollars:  spendPicodollars,
			}),
		)
	}

	addUsage(1, 100)
	addUsage(2, 200)
	addUsage(3, 300)

	unsettledUsage, err := querier.GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{
			PayerID:             payerID,
			MinutesSinceEpochGt: 2,
		},
	)
	require.NoError(t, err)
	require.Equal(t, unsettledUsage.TotalSpendPicodollars, int64(300))

	unsettledUsage, err = querier.GetPayerUnsettledUsage(
		ctx,
		queries.GetPayerUnsettledUsageParams{
			PayerID:             payerID,
			MinutesSinceEpochGt: 1,
		},
	)
	require.NoError(t, err)
	require.Equal(t, unsettledUsage.TotalSpendPicodollars, int64(500))
}
