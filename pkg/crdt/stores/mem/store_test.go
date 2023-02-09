package memstore_test

import (
	"testing"

	memstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
	crdttest "github.com/xmtp/xmtpd/pkg/crdt/testing"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func TestMemoryStore(t *testing.T) {
	crdttest.RunStoreTests(t, func(t *testing.T) *crdttest.TestStore {
		log := test.NewLogger(t)
		store := memstore.New(log)
		return crdttest.NewTestStore(store)
	})
}
