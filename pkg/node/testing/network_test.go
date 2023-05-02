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
	spray  string
	target string
)

func init() {
	flag.StringVar(&spray, "spray", "", "run spraying test with given number of topics and messages, format topics/messages, e.g. 30/1000")
	flag.StringVar(&target, "target", "devnet", "run spraying test (defined by --messages) against given target (devnet, localhost, a host address)")
}

func Test_Spraying(t *testing.T) {
	if len(spray) == 0 {
		return
	}
	var topics, messages int
	n, err := fmt.Sscanf(spray, "%d/%d", &topics, &messages)
	require.NoError(t, err)
	require.Equal(t, 2, n)

	ctx := test.NewContext(t)
	var clients []trackerNode
	switch target {
	case "devnet":
		for _, n := range []string{"node1", "node2", "node3"} {
			client := client.NewHTTPClient(ctx.Logger(),
				"http://localhost", "test", "spraying "+n,
				client.WithHeader("Host", n+".localhost"))
			clients = append(clients, client)
		}
	case "xmtp.pizza":
		for _, n := range []string{"node1", "node2", "node3", "node4"} {
			client := client.NewHTTPClient(ctx.Logger(),
				fmt.Sprintf("http://%s.xmtp.pizza", n), "test", "spraying "+n)
			clients = append(clients, client)
		}
	case "snormore.dev":
		for _, n := range []string{"node1", "node2", "node3"} {
			client := client.NewHTTPClient(ctx.Logger(),
				fmt.Sprintf("http://%s.snormore.dev", n), "test", "spraying "+n)
			clients = append(clients, client)
		}
	case "local":
		client := client.NewHTTPClient(ctx.Logger(),
			"http://localhost:5001", "test", "spraying")
		clients = append(clients, client)
	default:
		client := client.NewHTTPClient(ctx.Logger(),
			target, "test", "spraying")
		clients = append(clients, client)
	}
	tracker := newConvergenceTracker(ctx, clients)
	tracker.runRandomNodeAndTopicSpraying(t, topics, messages, "-"+time.Now().Format("060102T150405"))
}
