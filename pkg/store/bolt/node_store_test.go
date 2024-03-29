package bolt_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/node"
	ntest "github.com/xmtp/xmtpd/pkg/node/testing"
	test "github.com/xmtp/xmtpd/pkg/testing"

	"github.com/xmtp/xmtpd/pkg/store/bolt"
)

func Test_TopicBootstrap(t *testing.T) {
	opts := &bolt.Options{
		DataPath: filepath.Join(t.TempDir(), "testdb.bolt"),
	}
	ntest.TestTopicBootstrap(t, func(t testing.TB, ctx context.Context) node.NodeStore {
		db, err := bolt.NewDB(opts)
		require.NoError(t, err)
		store, err := bolt.NewNodeStore(ctx, db, opts)
		require.NoError(t, err)
		return store

	})
}

func Test_TopicLifecyle(t *testing.T) {
	opts := &bolt.Options{
		DataPath: filepath.Join(t.TempDir(), "testdb.bolt"),
	}
	db, err := bolt.NewDB(opts)
	require.NoError(t, err)
	ctx := test.NewContext(t)
	store, err := bolt.NewNodeStore(ctx, db, opts)
	require.NoError(t, err)
	ntest.RunTopicLifecycleTest(t, store)
}
