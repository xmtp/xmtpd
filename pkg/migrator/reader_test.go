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

	reader := migrator.NewGroupMessageReader(db)

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
			records, maxID, err := reader.Fetch(ctx, tc.lastID, tc.limit)
			require.NoError(t, err)
			require.Equal(t, tc.wantMaxID, maxID)

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
			records, maxID, err := reader.Fetch(ctx, tc.lastID, tc.limit)
			require.NoError(t, err)
			require.Equal(t, tc.wantMaxID, maxID)

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

	reader := migrator.NewKeyPackageReader(db)

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
			records, maxID, err := reader.Fetch(ctx, tc.lastID, tc.limit)
			require.NoError(t, err)
			require.Equal(t, tc.wantMaxID, maxID)

			for _, record := range records {
				require.IsType(t, &migrator.KeyPackage{}, record)
			}
		})
	}
}

func TestWelcomeMessageReader(t *testing.T) {
	ctx := t.Context()

	db, _, cleanup := testdata.NewMigratorTestDB(t, ctx)
	defer cleanup()

	reader := migrator.NewWelcomeMessageReader(db)

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
			records, maxID, err := reader.Fetch(ctx, tc.lastID, tc.limit)
			require.NoError(t, err)
			require.Equal(t, tc.wantMaxID, maxID)

			for _, record := range records {
				require.IsType(t, &migrator.WelcomeMessage{}, record)
			}
		})
	}
}
