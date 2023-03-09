package postgresstore_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	crdttest "github.com/xmtp/xmtpd/pkg/crdt/testing"
	postgresstore "github.com/xmtp/xmtpd/pkg/store/postgres"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func TestEvents(t *testing.T) {
	ctx := test.NewContext(t)
	topic := "topic-" + test.RandomStringLower(13)
	t.Parallel()
	crdttest.RunStoreEventTests(t, topic, func(t *testing.T) *crdttest.TestStore {
		store := newTestStore(t, topic)
		return crdttest.NewTestStore(ctx, store)
	})
}

func TestQuery(t *testing.T) {
	ctx := test.NewContext(t)
	topic := "topic-" + test.RandomStringLower(13)
	t.Parallel()
	crdttest.RunStoreQueryTests(t, topic, func(t *testing.T) *crdttest.TestStore {
		store := newTestStore(t, topic)
		return crdttest.NewTestStore(ctx, store)
	})
}

type testStore struct {
	*postgresstore.Store
	db        *postgresstore.DB
	dbCleanup func()
}

func newTestStore(t testing.TB, topic string) *testStore {
	t.Helper()
	ctx := test.NewContext(t)

	db, cleanup := newTestDB(t)
	store, err := postgresstore.NewNodeStore(ctx, db)
	require.NoError(t, err)
	scopedStore, err := store.NewTopic(topic)
	require.NoError(t, err)

	return &testStore{
		Store:     scopedStore.(*postgresstore.Store),
		db:        db,
		dbCleanup: cleanup,
	}
}

func (s *testStore) Close() error {
	s.dbCleanup()
	return nil
}
