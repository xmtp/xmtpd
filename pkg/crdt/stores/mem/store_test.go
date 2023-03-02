package memstore_test

import (
	"testing"

	memstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
	crdttest "github.com/xmtp/xmtpd/pkg/crdt/testing"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func TestMemoryStore(t *testing.T) {
	ctx := test.NewContext(t)
	topic := "topic-" + test.RandomStringLower(13)

	crdttest.RunStoreEventTests(t, topic, func(t *testing.T) *crdttest.TestStore {
		store := memstore.New(ctx)
		return crdttest.NewTestStore(ctx, store)
	})

	crdttest.RunStoreQueryTests(t, topic, func(t *testing.T) *crdttest.TestStore {
		store := memstore.New(test.NewContext(t))
		return crdttest.NewTestStore(ctx, store)
	})
}
