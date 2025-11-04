package metadata_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/api/metadata"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/utils"
)

type testMessage struct {
	timestamp        time.Time
	spendPicodollars int64
	originatorID     int32
}

type testSetup struct {
	database *sql.DB
	fetcher  *metadata.PayerInfoFetcher
	payerID  int32
	baseTime time.Time
}

func setupPayerInfoTest(t *testing.T) *testSetup {
	database, _ := testutils.NewDB(t, t.Context())
	fetcher := metadata.NewPayerInfoFetcher(database)
	payerID := testutils.CreatePayer(t, database, testutils.RandomAddress().Hex())
	return &testSetup{
		database: database,
		fetcher:  fetcher,
		payerID:  payerID,
		baseTime: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
	}
}

func (ts *testSetup) insertMessages(t *testing.T, messages []testMessage) {
	ctx := t.Context()
	for i, msg := range messages {
		minutesSinceEpoch := utils.MinutesSinceEpoch(msg.timestamp)

		insertParams := queries.InsertGatewayEnvelopeParams{
			OriginatorNodeID: msg.originatorID,
			OriginatorSequenceID: int64(
				minutesSinceEpoch,
			)*10000 + int64(
				msg.originatorID,
			)*100 + int64(
				i,
			), // Ensure uniqueness
			Topic:              testutils.RandomBytes(32),
			OriginatorEnvelope: testutils.RandomBytes(100),
			PayerID:            db.NullInt32(ts.payerID),
		}

		incrementParams := queries.IncrementUnsettledUsageParams{
			PayerID:           ts.payerID,
			OriginatorID:      msg.originatorID,
			MinutesSinceEpoch: minutesSinceEpoch,
			SpendPicodollars:  msg.spendPicodollars,
			MessageCount:      1,
			SequenceID:        insertParams.OriginatorSequenceID,
		}

		numInserted, err := db.InsertGatewayEnvelopeAndIncrementUnsettledUsage(
			ctx,
			ts.database,
			insertParams,
			incrementParams,
		)
		require.NoError(t, err, msg)
		require.Equal(t, int64(1), numInserted)
	}
}

func TestPayerInfo_Granularity(t *testing.T) {
	setup := setupPayerInfoTest(t)
	ctx := t.Context()

	// Insert test data spanning multiple hours and days
	messages := []testMessage{
		{setup.baseTime, 100, 1},                       // Day 1, Hour 12
		{setup.baseTime.Add(30 * time.Minute), 200, 1}, // Day 1, Hour 12
		{setup.baseTime.Add(1 * time.Hour), 300, 1},    // Day 1, Hour 13
		{setup.baseTime.Add(12 * time.Hour), 400, 1},   // Day 2, Hour 0
		{setup.baseTime.Add(24 * time.Hour), 500, 1},   // Day 2, Hour 12
		{setup.baseTime.Add(48 * time.Hour), 600, 1},   // Day 3, Hour 12
	}
	setup.insertMessages(t, messages)

	tests := []struct {
		name     string
		groupBy  metadata.PayerInfoGroupBy
		expected []struct {
			amount    uint64
			messages  uint64
			startTime time.Time
		}
	}{
		{
			name:    "hourly_granularity",
			groupBy: metadata.PayerInfoGroupByHour,
			expected: []struct {
				amount    uint64
				messages  uint64
				startTime time.Time
			}{
				{
					300,
					2,
					setup.baseTime.Truncate(time.Hour),
				}, // Hour 12: 100+200
				{300, 1, setup.baseTime.Add(1 * time.Hour).Truncate(time.Hour)},  // Hour 13
				{400, 1, setup.baseTime.Add(12 * time.Hour).Truncate(time.Hour)}, // Day 2, Hour 0
				{500, 1, setup.baseTime.Add(24 * time.Hour).Truncate(time.Hour)}, // Day 2, Hour 12
				{600, 1, setup.baseTime.Add(48 * time.Hour).Truncate(time.Hour)}, // Day 3, Hour 12
			},
		},
		{
			name:    "daily_granularity",
			groupBy: metadata.PayerInfoGroupByDay,
			expected: []struct {
				amount    uint64
				messages  uint64
				startTime time.Time
			}{
				{600, 3, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}, // Day 1: 100+200+300
				{900, 2, time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)}, // Day 2: 400+500
				{600, 1, time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)}, // Day 3: 600
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payerInfo, err := setup.fetcher.GetPayerInfo(
				ctx,
				setup.payerID,
				tt.groupBy,
			)
			require.NoError(t, err)
			require.NotNil(t, payerInfo)
			require.Len(t, payerInfo.PeriodSummaries, len(tt.expected))

			for i, expected := range tt.expected {
				require.Equal(
					t,
					expected.amount,
					payerInfo.PeriodSummaries[i].AmountSpentPicodollars,
					"Period %d amount mismatch",
					i,
				)
				require.Equal(t, expected.messages, payerInfo.PeriodSummaries[i].NumMessages,
					"Period %d message count mismatch", i)
				require.Equal(
					t,
					uint64(expected.startTime.Unix()),
					payerInfo.PeriodSummaries[i].PeriodStartUnixSeconds,
					"Period %d start time mismatch",
					i,
				)
			}
		})
	}
}

func TestPayerInfo_MultipleOriginators(t *testing.T) {
	setup := setupPayerInfoTest(t)
	ctx := t.Context()

	// Insert messages from multiple originators in the same period
	messages := []testMessage{
		{setup.baseTime, 100, 100},
		{setup.baseTime.Add(10 * time.Minute), 200, 200},
		{setup.baseTime.Add(20 * time.Minute), 300, 300},
		{setup.baseTime.Add(30 * time.Minute), 400, 100}, // Same originator as first
	}
	setup.insertMessages(t, messages)

	// Test hourly aggregation
	payerInfo, err := setup.fetcher.GetPayerInfo(
		ctx,
		setup.payerID,
		metadata.PayerInfoGroupByHour,
	)
	require.NoError(t, err)
	require.NotNil(t, payerInfo)
	require.Len(t, payerInfo.PeriodSummaries, 1)

	// All messages should be aggregated into one hour
	require.Equal(t, uint64(1000), payerInfo.PeriodSummaries[0].AmountSpentPicodollars)
	require.Equal(t, uint64(4), payerInfo.PeriodSummaries[0].NumMessages)
}

func TestPayerInfo_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "empty_result",
			testFunc: func(t *testing.T) {
				setup := setupPayerInfoTest(t)
				ctx := t.Context()

				// Query for future time period
				payerInfo, err := setup.fetcher.GetPayerInfo(
					ctx,
					setup.payerID,
					metadata.PayerInfoGroupByHour,
				)
				require.NoError(t, err)
				require.NotNil(t, payerInfo)
				require.Len(t, payerInfo.PeriodSummaries, 0)
			},
		},
		{
			name: "non_existent_payer",
			testFunc: func(t *testing.T) {
				setup := setupPayerInfoTest(t)
				ctx := t.Context()

				// Query for non-existent payer
				payerInfo, err := setup.fetcher.GetPayerInfo(
					ctx,
					99999,
					metadata.PayerInfoGroupByHour,
				)
				require.NoError(t, err)
				require.NotNil(t, payerInfo)
				require.Len(t, payerInfo.PeriodSummaries, 0)
			},
		},
		{
			name: "chronological_ordering",
			testFunc: func(t *testing.T) {
				setup := setupPayerInfoTest(t)
				ctx := t.Context()

				// Insert messages in non-chronological order
				messages := []testMessage{
					{setup.baseTime.Add(48 * time.Hour), 300, 1}, // Day 3
					{setup.baseTime, 100, 1},                     // Day 1
					{setup.baseTime.Add(24 * time.Hour), 200, 1}, // Day 2
				}
				setup.insertMessages(t, messages)

				// Fetch daily data
				payerInfo, err := setup.fetcher.GetPayerInfo(
					ctx,
					setup.payerID,
					metadata.PayerInfoGroupByDay,
				)
				require.NoError(t, err)
				require.NotNil(t, payerInfo)
				require.Len(t, payerInfo.PeriodSummaries, 3)

				// Verify chronological ordering
				require.Equal(t, uint64(100), payerInfo.PeriodSummaries[0].AmountSpentPicodollars)
				require.Equal(t, uint64(200), payerInfo.PeriodSummaries[1].AmountSpentPicodollars)
				require.Equal(t, uint64(300), payerInfo.PeriodSummaries[2].AmountSpentPicodollars)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func TestGetPayerByAddress(t *testing.T) {
	ctx := t.Context()
	database, _ := testutils.NewDB(t, ctx)
	fetcher := metadata.NewPayerInfoFetcher(database)

	t.Run("existing_payer", func(t *testing.T) {
		// Create a payer with a known address
		address := testutils.RandomAddress().Hex()
		expectedPayerID := testutils.CreatePayer(t, database, address)

		// Test GetPayerByAddress
		payerID, err := fetcher.GetPayerByAddress(ctx, address)
		require.NoError(t, err)
		require.Equal(t, expectedPayerID, payerID)
	})

	t.Run("non_existent_payer", func(t *testing.T) {
		// Try to get a payer that doesn't exist
		nonExistentAddress := testutils.RandomAddress().Hex()

		_, err := fetcher.GetPayerByAddress(ctx, nonExistentAddress)
		require.Error(t, err)
		require.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("multiple_payers", func(t *testing.T) {
		// Create multiple payers
		address1 := testutils.RandomAddress().Hex()
		address2 := testutils.RandomAddress().Hex()
		address3 := testutils.RandomAddress().Hex()

		payerID1 := testutils.CreatePayer(t, database, address1)
		payerID2 := testutils.CreatePayer(t, database, address2)
		payerID3 := testutils.CreatePayer(t, database, address3)

		// Verify each can be looked up correctly
		foundID1, err := fetcher.GetPayerByAddress(ctx, address1)
		require.NoError(t, err)
		require.Equal(t, payerID1, foundID1)

		foundID2, err := fetcher.GetPayerByAddress(ctx, address2)
		require.NoError(t, err)
		require.Equal(t, payerID2, foundID2)

		foundID3, err := fetcher.GetPayerByAddress(ctx, address3)
		require.NoError(t, err)
		require.Equal(t, payerID3, foundID3)
	})

	t.Run("case_sensitivity", func(t *testing.T) {
		// Create a payer with lowercase address
		lowerAddress := testutils.RandomAddress().Hex()
		payerID := testutils.CreatePayer(t, database, lowerAddress)

		// Test exact match
		foundID, err := fetcher.GetPayerByAddress(ctx, lowerAddress)
		require.NoError(t, err)
		require.Equal(t, payerID, foundID)

		// Test with different address (addresses are case-sensitive in the DB)
		differentAddress := testutils.RandomAddress().Hex()
		_, err = fetcher.GetPayerByAddress(ctx, differentAddress)
		require.Error(t, err)
		require.Equal(t, sql.ErrNoRows, err)
	})
}

func TestPayerInfo_BoundaryConditions(t *testing.T) {
	setup := setupPayerInfoTest(t)
	ctx := t.Context()

	// Test messages at day/hour boundaries
	dayBoundary := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	messages := []testMessage{
		{dayBoundary.Add(-1 * time.Second), 100, 1}, // Last second of Day 1
		{dayBoundary, 200, 1},                       // First second of Day 2
		{dayBoundary.Add(1 * time.Second), 300, 1},  // Second second of Day 2
	}
	setup.insertMessages(t, messages)

	t.Run("daily_boundary", func(t *testing.T) {
		payerInfo, err := setup.fetcher.GetPayerInfo(
			ctx,
			setup.payerID,
			metadata.PayerInfoGroupByDay,
		)
		require.NoError(t, err)
		require.NotNil(t, payerInfo)
		require.Len(t, payerInfo.PeriodSummaries, 2)

		// Day 1 should have only the first message
		require.Equal(t, uint64(100), payerInfo.PeriodSummaries[0].AmountSpentPicodollars)
		require.Equal(t, uint64(1), payerInfo.PeriodSummaries[0].NumMessages)

		// Day 2 should have the other two messages
		require.Equal(t, uint64(500), payerInfo.PeriodSummaries[1].AmountSpentPicodollars)
		require.Equal(t, uint64(2), payerInfo.PeriodSummaries[1].NumMessages)
	})

	t.Run("hourly_boundary", func(t *testing.T) {
		payerInfo, err := setup.fetcher.GetPayerInfo(
			ctx,
			setup.payerID,
			metadata.PayerInfoGroupByHour,
		)
		require.NoError(t, err)
		require.NotNil(t, payerInfo)
		require.Len(t, payerInfo.PeriodSummaries, 2)

		// Hour 23 of Day 1
		require.Equal(t, uint64(100), payerInfo.PeriodSummaries[0].AmountSpentPicodollars)

		// Hour 0 of Day 2
		require.Equal(t, uint64(500), payerInfo.PeriodSummaries[1].AmountSpentPicodollars)
	})
}
