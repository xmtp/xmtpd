package db_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func setupTest(t *testing.T) (context.Context, *queries.Queries) {
	ctx := context.Background()
	db, _ := testutils.NewRawDB(t, ctx)

	querier := queries.New(db)

	return ctx, querier
}

func incrementCongestion(
	t *testing.T,
	ctx context.Context,
	querier *queries.Queries,
	originatorID, minutesSinceEpoch int32,
) {
	err := querier.IncrementOriginatorCongestion(ctx, queries.IncrementOriginatorCongestionParams{
		OriginatorID:      originatorID,
		MinutesSinceEpoch: minutesSinceEpoch,
	})

	require.NoError(t, err)
}

func TestGet5MinutesOfCongestion(t *testing.T) {
	ctx, querier := setupTest(t)

	originatorID := int32(100)
	endMinute := testutils.RandomInt32()

	incrementCongestion(t, ctx, querier, originatorID, endMinute-1)
	incrementCongestion(t, ctx, querier, originatorID, endMinute-2)
	incrementCongestion(t, ctx, querier, originatorID, endMinute-10)

	congestion, err := db.Get5MinutesOfCongestion(ctx, querier, originatorID, endMinute)
	require.NoError(t, err)

	require.Equal(t, congestion[0], int32(0))
	require.Equal(t, congestion[1], int32(1))
	require.Equal(t, congestion[2], int32(1))
	require.Equal(t, congestion[3], int32(0))
	require.Equal(t, congestion[4], int32(0))
}

func TestMultipleIncrements(t *testing.T) {
	ctx, querier := setupTest(t)

	originatorID := int32(100)
	endMinute := testutils.RandomInt32()

	incrementCongestion(t, ctx, querier, originatorID, endMinute)
	incrementCongestion(t, ctx, querier, originatorID, endMinute)
	incrementCongestion(t, ctx, querier, originatorID, endMinute)

	congestion, err := db.Get5MinutesOfCongestion(ctx, querier, originatorID, endMinute)
	require.NoError(t, err)

	require.Equal(t, congestion[0], int32(3))
	require.Equal(t, congestion[1], int32(0))
	require.Equal(t, congestion[2], int32(0))
	require.Equal(t, congestion[3], int32(0))
	require.Equal(t, congestion[4], int32(0))
}

func TestNoCongestion(t *testing.T) {
	ctx, querier := setupTest(t)

	originatorID := int32(100)
	endMinute := testutils.RandomInt32()

	congestion, err := db.Get5MinutesOfCongestion(ctx, querier, originatorID, endMinute)
	require.NoError(t, err)

	require.Equal(t, congestion, [5]int32{0, 0, 0, 0, 0})
}
