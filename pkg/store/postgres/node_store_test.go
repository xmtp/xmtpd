package postgresstore_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/node"
	ntest "github.com/xmtp/xmtpd/pkg/node/testing"
	postgresstore "github.com/xmtp/xmtpd/pkg/store/postgres"
)

func Test_TopicBootstrap(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()
	ntest.TestTopicBootstrap(t, func(t *testing.T, ctx context.Context) node.NodeStore {
		store, err := postgresstore.NewNodeStore(ctx, db)
		require.NoError(t, err)
		return store
	})
}
