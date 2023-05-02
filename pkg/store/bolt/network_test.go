package bolt_test

import (
	"fmt"
	"testing"

	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/node"
	ntest "github.com/xmtp/xmtpd/pkg/node/testing"
)

func Test_RandomNodeAndTopicSpraying(t *testing.T) {
	if testing.Short() {
		return
	}
	tcs := []struct {
		nodes    int
		topics   int
		messages int
	}{
		{3, 10, 300},
		{5, 3, 100},
		{10, 5, 100},
	}
	for i, tc := range tcs {
		tc := tc
		name := fmt.Sprintf("%d/%dn/%dt/%dm", i, tc.nodes, tc.topics, tc.messages)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ntest.RunRandomNodeAndTopicSpraying(t, tc.nodes, tc.topics, tc.messages,
				ntest.WithStoreMaker(func(t testing.TB, ctx context.Context) node.NodeStore {
					return newTestNodeStore(t, ctx)
				}))
		})
	}
}

func Test_WaitForPubSub(t *testing.T) {
	if testing.Short() {
		return
	}
	net := ntest.NewNetwork(t, 8,
		ntest.WithStoreMaker(func(t testing.TB, ctx context.Context) node.NodeStore {
			return newTestNodeStore(t, ctx)
		}))
	net.WaitForPubSub(t)
}
