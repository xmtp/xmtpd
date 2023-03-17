package postgresstore_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/node"
	ntest "github.com/xmtp/xmtpd/pkg/node/testing"
	postgresstore "github.com/xmtp/xmtpd/pkg/store/postgres"
)

func Test_RandomNodeAndTopicSpraying(t *testing.T) {
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
			db, cleanup := newTestDB(t)
			defer cleanup()
			ntest.RandomNodeAndTopicSpraying(t, tc.nodes, tc.topics, tc.messages,
				ntest.WithStoreMaker(func(t testing.TB, ctx context.Context) node.NodeStore {
					store, err := postgresstore.NewNodeStore(ctx, db)
					require.NoError(t, err)
					return store
				}))
		})
	}
}
