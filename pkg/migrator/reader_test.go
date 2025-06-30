package migrator_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/migrator"
	"github.com/xmtp/xmtpd/pkg/migrator/testdata"
)

func Test_GroupMessageReader(t *testing.T) {
	ctx := t.Context()

	db, cleanup := testdata.NewTestDB(t, ctx)
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
				_, ok := record.(*migrator.GroupMessage)
				require.True(t, ok)
			}
		})
	}
}

func Test_InboxLogReader(t *testing.T) {
	ctx := t.Context()

	db, cleanup := testdata.NewTestDB(t, ctx)
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
				_, ok := record.(*migrator.InboxLog)
				require.True(t, ok)
			}
		})
	}
}

func Test_InstallationReader(t *testing.T) {
	ctx := t.Context()

	db, cleanup := testdata.NewTestDB(t, ctx)
	defer cleanup()

	reader := migrator.NewInstallationReader(db)

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
			wantMaxID: 1717242270261135027,
		},
		{
			name:      "Fetch 10",
			lastID:    5,
			limit:     10,
			wantMaxID: 1717473253760941762,
		},
		{
			name:      "Fetch all",
			lastID:    10,
			limit:     9999,
			wantMaxID: 1717490371754970003, // Max ID in the test data
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			records, maxID, err := reader.Fetch(ctx, tc.lastID, tc.limit)
			require.NoError(t, err)
			require.Equal(t, tc.wantMaxID, maxID)

			for _, record := range records {
				_, ok := record.(*migrator.Installation)
				require.True(t, ok)
			}
		})
	}
}

func Test_WelcomeMessageReader(t *testing.T) {
	ctx := t.Context()

	db, cleanup := testdata.NewTestDB(t, ctx)
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
				_, ok := record.(*migrator.WelcomeMessage)
				require.True(t, ok)
			}
		})
	}
}
