package postgresstore_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	ntest "github.com/xmtp/xmtpd/pkg/node/testing"
)

func Test_RandomNodeAndTopicSpraying(t *testing.T) {
	if testing.Short() {
		return
	}
	type testCase struct {
		nodes    int
		topics   int
		messages int
	}
	tcs := []testCase{
		{3, 10, 300},
		{5, 3, 100},
		{10, 5, 100},
	}
	if os.Getenv("CI") == "true" {
		// It seems we need to give postgres a sec to start up in CI
		// getting postgres dial errors too often.
		time.Sleep(5 * time.Second)
	}
	for i, tc := range tcs {
		tc := tc
		name := fmt.Sprintf("%d/%dn/%dt/%dm", i, tc.nodes, tc.topics, tc.messages)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ntest.RunRandomNodeAndTopicSpraying(t, tc.nodes, tc.topics, tc.messages,
				ntest.WithStoreMaker(newTestNodeStore))
		})
	}
}

func Test_WaitForPubSub(t *testing.T) {
	if testing.Short() {
		return
	}
	net := ntest.NewNetwork(t, 10,
		ntest.WithStoreMaker(newTestNodeStore))
	net.WaitForPubSub(t)
}
