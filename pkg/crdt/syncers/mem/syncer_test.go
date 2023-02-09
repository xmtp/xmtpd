package memsyncer_test

import (
	"context"
	"testing"

	memstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
	memsyncer "github.com/xmtp/xmtpd/pkg/crdt/syncers/mem"
	crdttest "github.com/xmtp/xmtpd/pkg/crdt/testing"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func TestMemorySyncer(t *testing.T) {
	crdttest.RunSyncerTests(t, func(t *testing.T) *crdttest.TestSyncer {
		ctx := context.Background()
		log := test.NewLogger(t)
		store := memstore.New(log)
		syncer := memsyncer.New(log, store)
		return crdttest.NewTestSyncer(ctx, log, syncer)
	})
}
