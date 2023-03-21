package testing

import (
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/api/client"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

var (
	devnet string
)

func init() {
	flag.StringVar(&devnet, "devnet", "", "run devnet tests with given number of topics and messages, format topics/messages, e.g. 30/1000")
}

func Test_DevnetSpraying(t *testing.T) {
	if len(devnet) == 0 {
		return
	}
	var topics, messages int
	n, err := fmt.Sscanf(devnet, "%d/%d", &topics, &messages)
	require.NoError(t, err)
	require.Equal(t, 2, n)
	ctx := test.NewContext(t)
	var clients []trackerNode
	for _, n := range []string{"node1", "node2", "node3"} {
		client := client.NewHTTPClient(ctx.Logger(),
			"http://localhost", "test", n,
			client.WithHeader("Host", n+".localhost"))
		clients = append(clients, client)
	}
	tracker := newConvergenceTracker(ctx, clients)
	tracker.runRandomNodeAndTopicSpraying(t, topics, messages, "-"+time.Now().Format("060102T150405"))
}
