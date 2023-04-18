package e2e

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	nodetest "github.com/xmtp/xmtpd/pkg/node/testing"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func TestE2E(t *testing.T) {
	ctx := test.NewContext(t)
	defer ctx.Close()

	n1 := nodetest.NewNode(t, nodetest.WithName("node1"))
	defer n1.Close()

	n2 := nodetest.NewNode(t, nodetest.WithName("node2"))
	defer n2.Close()
	n2.Connect(t, n1)

	n3 := nodetest.NewNode(t, nodetest.WithName("node3"))
	defer n3.Close()
	n3.Connect(t, n1)

	e2e, err := New(ctx, &Options{
		APIURLs: []string{
			fmt.Sprintf("http://localhost:%d", n1.APIHTTPListenPort()),
			fmt.Sprintf("http://localhost:%d", n2.APIHTTPListenPort()),
			fmt.Sprintf("http://localhost:%d", n3.APIHTTPListenPort()),
		},
		MessagePerClient: 3,
		ClientsPerURL:    1,
	})
	require.NoError(t, err)

	for _, test := range e2e.Tests() {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			err := test.Run()
			require.NoError(t, err)
		})
	}
}
