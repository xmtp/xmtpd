package migrator_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/migrator"
	"github.com/xmtp/xmtpd/pkg/migrator/testdata"
)

func TestGroupMessageReader(t *testing.T) {
	ctx := t.Context()

	db, _, cleanup := testdata.NewMigratorTestDB(t, ctx)
	defer cleanup()

	reader := migrator.NewGroupMessageReader(db, 0)

	cases := []struct {
		name      string
		lastID    int64
		limit     int32
		wantMaxID int64
	}{
		{
			name:      "Fetch 0",
			lastID:    0,
			limit:     0,
			wantMaxID: 0,
		},
		{
			name:      "Fetch 5",
			lastID:    0,
			limit:     5,
			wantMaxID: 5,
		},
		{
			name:      "Fetch 10",
			lastID:    5,
			limit:     10,
			wantMaxID: 15,
		},
		{
			name:      "Fetch all",
			lastID:    10,
			limit:     9999,
			wantMaxID: 19, // Max ID in the test data
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			records, err := reader.Fetch(ctx, tc.lastID, tc.limit)
			require.NoError(t, err)

			for _, record := range records {
				require.IsType(t, &migrator.GroupMessage{}, record)
			}
		})
	}
}

func TestInboxLogReader(t *testing.T) {
	ctx := t.Context()

	db, _, cleanup := testdata.NewMigratorTestDB(t, ctx)
	defer cleanup()

	reader := migrator.NewInboxLogReader(db)

	cases := []struct {
		name      string
		lastID    int64
		limit     int32
		wantMaxID int64
	}{
		{
			name:      "Fetch 0",
			lastID:    0,
			limit:     0,
			wantMaxID: 0,
		},
		{
			name:      "Fetch 5",
			lastID:    0,
			limit:     5,
			wantMaxID: 5,
		},
		{
			name:      "Fetch 10",
			lastID:    5,
			limit:     10,
			wantMaxID: 15,
		},
		{
			name:      "Fetch all",
			lastID:    10,
			limit:     9999,
			wantMaxID: 19, // Max ID in the test data
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			records, err := reader.Fetch(ctx, tc.lastID, tc.limit)
			require.NoError(t, err)

			for _, record := range records {
				require.IsType(t, &migrator.InboxLog{}, record)
			}
		})
	}
}

func TestKeyPackageReader(t *testing.T) {
	ctx := t.Context()

	db, _, cleanup := testdata.NewMigratorTestDB(t, ctx)
	defer cleanup()

	reader := migrator.NewKeyPackageReader(db, 0)

	cases := []struct {
		name      string
		lastID    int64
		limit     int32
		wantMaxID int64
	}{
		{
			name:      "Fetch 0",
			lastID:    0,
			limit:     0,
			wantMaxID: 0,
		},
		{
			name:      "Fetch 5",
			lastID:    0,
			limit:     5,
			wantMaxID: 5,
		},
		{
			name:      "Fetch 10",
			lastID:    5,
			limit:     10,
			wantMaxID: 15,
		},
		{
			name:      "Fetch all",
			lastID:    10,
			limit:     9999,
			wantMaxID: 19,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			records, err := reader.Fetch(ctx, tc.lastID, tc.limit)
			require.NoError(t, err)

			for _, record := range records {
				require.IsType(t, &migrator.KeyPackage{}, record)
			}
		})
	}
}

func TestCommitMessageReader(t *testing.T) {
	ctx := t.Context()

	db, _, cleanup := testdata.NewMigratorTestDB(t, ctx)
	defer cleanup()

	reader := migrator.NewCommitMessageReader(db)

	cases := []struct {
		name   string
		lastID int64
		limit  int32
	}{
		{
			name:   "Fetch 0",
			lastID: 0,
			limit:  0,
		},
		{
			name:   "Fetch 5",
			lastID: 0,
			limit:  5,
		},
		{
			name:   "Fetch 10",
			lastID: 5,
			limit:  10,
		},
		{
			name:   "Fetch all",
			lastID: 10,
			limit:  9999,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			records, err := reader.Fetch(ctx, tc.lastID, tc.limit)
			require.NoError(t, err)

			for _, record := range records {
				require.IsType(t, &migrator.CommitMessage{}, record)
			}
		})
	}
}

// TestGroupMessageReaderLowerLimit verifies that the lower limit correctly
// excludes records with IDs below the threshold.
func TestGroupMessageReaderLowerLimit(t *testing.T) {
	ctx := t.Context()

	db, _, cleanup := testdata.NewMigratorTestDB(t, ctx)
	defer cleanup()

	// Set lower limit to 10; only records with id >= 10 should be returned.
	const lowerLimit int64 = 10
	reader := migrator.NewGroupMessageReader(db, lowerLimit)

	records, err := reader.Fetch(ctx, 0, 9999)
	require.NoError(t, err)
	require.NotEmpty(t, records)

	for _, r := range records {
		require.GreaterOrEqual(t, r.GetID(), lowerLimit,
			"record id %d is below lower limit %d", r.GetID(), lowerLimit)
	}
}

// TestWelcomeMessageReaderLowerLimit verifies that welcome message readers with a
// high lower limit skip early records — this mirrors the production config where
// welcome messages start at a large offset (e.g. 150 000 000).
func TestWelcomeMessageReaderLowerLimit(t *testing.T) {
	ctx := t.Context()

	db, _, cleanup := testdata.NewMigratorTestDB(t, ctx)
	defer cleanup()

	// The test dataset contains IDs starting at 150000001. Using a limit of
	// 150000010 means only IDs >= 150000010 should be returned.
	const lowerLimit int64 = 150_000_010
	reader := migrator.NewWelcomeMessageReader(db, lowerLimit)

	records, err := reader.Fetch(ctx, 0, 9999)
	require.NoError(t, err)
	require.NotEmpty(t, records)

	for _, r := range records {
		require.GreaterOrEqual(t, r.GetID(), lowerLimit,
			"record id %d is below lower limit %d", r.GetID(), lowerLimit)
	}
}

// TestWelcomeMessageReaderLowerLimit_SkipsAll verifies that when the lower limit
// is above all available records, no records are returned.
func TestWelcomeMessageReaderLowerLimit_SkipsAll(t *testing.T) {
	ctx := t.Context()

	db, _, cleanup := testdata.NewMigratorTestDB(t, ctx)
	defer cleanup()

	// lowerLimit above all data in the test set.
	const lowerLimit int64 = 999_999_999
	reader := migrator.NewWelcomeMessageReader(db, lowerLimit)

	records, err := reader.Fetch(ctx, 0, 9999)
	require.NoError(t, err)
	require.Empty(t, records)
}

// TestKeyPackageReaderLowerLimit verifies that the key-package reader respects
// the lower limit and skips records below it.
func TestKeyPackageReaderLowerLimit(t *testing.T) {
	ctx := t.Context()

	db, _, cleanup := testdata.NewMigratorTestDB(t, ctx)
	defer cleanup()

	const lowerLimit int64 = 5
	reader := migrator.NewKeyPackageReader(db, lowerLimit)

	records, err := reader.Fetch(ctx, 0, 9999)
	require.NoError(t, err)
	require.NotEmpty(t, records)

	for _, r := range records {
		require.GreaterOrEqual(t, r.GetID(), lowerLimit,
			"record id %d is below lower limit %d", r.GetID(), lowerLimit)
	}
}

func TestWelcomeMessageReader(t *testing.T) {
	ctx := t.Context()

	db, _, cleanup := testdata.NewMigratorTestDB(t, ctx)
	defer cleanup()

	reader := migrator.NewWelcomeMessageReader(db, 0)

	cases := []struct {
		name      string
		lastID    int64
		limit     int32
		wantMaxID int64
	}{
		{
			name:      "Fetch 0",
			lastID:    0,
			limit:     0,
			wantMaxID: 0,
		},
		{
			name:      "Fetch 5",
			lastID:    0,
			limit:     5,
			wantMaxID: 5,
		},
		{
			name:      "Fetch 10",
			lastID:    5,
			limit:     10,
			wantMaxID: 15,
		},
		{
			name:      "Fetch all",
			lastID:    10,
			limit:     9999,
			wantMaxID: 19, // Max ID in the test data
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			records, err := reader.Fetch(ctx, tc.lastID, tc.limit)
			require.NoError(t, err)

			for _, record := range records {
				require.IsType(t, &migrator.WelcomeMessage{}, record)
			}
		})
	}
}
