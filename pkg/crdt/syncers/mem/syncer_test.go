package memsyncer_test

import (
	"testing"

	memstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
	memsyncer "github.com/xmtp/xmtpd/pkg/crdt/syncers/mem"
	crdttest "github.com/xmtp/xmtpd/pkg/crdt/testing"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func TestMemorySyncer(t *testing.T) {
	crdttest.RunSyncerTests(t, func(t *testing.T) *crdttest.TestSyncer {
		ctx := test.NewContext(t)
		store := memstore.New(ctx)
		syncer := memsyncer.New(ctx, store)
		return crdttest.NewTestSyncer(ctx, syncer)
	})
}
