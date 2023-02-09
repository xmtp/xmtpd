package membroadcaster_test

import (
	"context"
	"testing"

	membroadcaster "github.com/xmtp/xmtpd/pkg/crdt/broadcasters/mem"
	crdttest "github.com/xmtp/xmtpd/pkg/crdt/testing"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func TestMemoryBroadcaster(t *testing.T) {
	crdttest.RunBroadcasterTests(t, func(t *testing.T) *crdttest.TestBroadcaster {
		ctx := context.Background()
		log := test.NewLogger(t)
		bc := membroadcaster.New(log)
		return crdttest.NewTestBroadcaster(ctx, log, bc)
	})
}
