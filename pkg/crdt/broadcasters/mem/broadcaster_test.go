package membroadcaster_test

import (
	"testing"

	membroadcaster "github.com/xmtp/xmtpd/pkg/crdt/broadcasters/mem"
	crdttest "github.com/xmtp/xmtpd/pkg/crdt/testing"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func TestMemoryBroadcaster(t *testing.T) {
	crdttest.TestBroadcaster_BroadcastNext(t, func(t *testing.T) *crdttest.TestBroadcaster {
		log := test.NewLogger(t)
		bc := membroadcaster.New(log)
		return crdttest.NewTestBroadcaster(bc)
	})
}
