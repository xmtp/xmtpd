package orbitdbnode

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	orbitdb "berty.tech/go-orbit-db"
	orbitdbiface "berty.tech/go-orbit-db/iface"
	"berty.tech/go-orbit-db/stores/operation"
	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	iface "github.com/ipfs/interface-go-ipfs-core"
	cfg "github.com/ipfs/kubo/config"
	ipfscore "github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/coreapi"
	mock "github.com/ipfs/kubo/core/mock"
	"github.com/ipfs/kubo/repo"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
	"github.com/stretchr/testify/require"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func TestOrbitDB_NewClose_Mocked(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("add and list", func(t *testing.T) {
		orbitDB, cleanup := newTestMockedOrbitDB(t)
		defer cleanup()
		testAddAndList(t, ctx, orbitDB)
	})
}

func TestOrbitDB_NewClose_NonMocked(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("add and list", func(t *testing.T) {
		orbitDB, cleanup := newTestNonMockedOrbitDB(t)
		defer cleanup()
		testAddAndList(t, ctx, orbitDB)
	})
}

func newTestMockedOrbitDB(t testing.TB) (orbitdbiface.OrbitDB, func()) {
	ctx := context.Background()

	mocknet := testingMockNet(t)
	node, nodeClean := testingMockedIPFSNode(ctx, t, mocknet)

	ipfs := testingCoreAPI(t, node)

	dataPath := filepath.Join(t.TempDir(), test.RandomStringLower(13))

	db, err := orbitdb.NewOrbitDB(ctx, ipfs, &orbitdb.NewOrbitDBOptions{
		Directory: &dataPath,
	})
	require.NoError(t, err)

	return db, func() {
		db.Close()
		nodeClean()
	}
}

func newTestNonMockedOrbitDB(t testing.TB) (orbitdbiface.OrbitDB, func()) {
	ctx := context.Background()

	node, nodeClean := testingNonMockedIPFSNode(ctx, t)

	ipfs := testingCoreAPI(t, node)

	dataPath := filepath.Join(t.TempDir(), test.RandomStringLower(13))

	db, err := orbitdb.NewOrbitDB(ctx, ipfs, &orbitdb.NewOrbitDBOptions{
		Directory: &dataPath,
	})
	require.NoError(t, err)

	return db, func() {
		db.Close()
		nodeClean()
	}
}

func testAddAndList(t testing.TB, ctx context.Context, orbitDB orbitdbiface.OrbitDB) {
	t.Helper()

	infinity := -1
	db, err := orbitDB.Log(ctx, "log database", nil)
	require.NoError(t, err)
	require.NotNil(t, db)

	defer db.Close()

	require.Equal(t, db.Type(), "eventlog")
	require.Equal(t, db.DBName(), "log database")

	// Returns 0 items when it's a fresh database.
	res := make(chan operation.Operation, 100)
	err = db.Stream(ctx, res, &orbitdb.StreamOptions{Amount: &infinity})
	require.NoError(t, err)
	require.Equal(t, len(res), 0)

	// Returns the added entry.
	op, err := db.Add(ctx, []byte("hello1"))
	require.NoError(t, err)

	ops, err := db.List(ctx, &orbitdb.StreamOptions{Amount: &infinity})

	require.NoError(t, err)
	require.Equal(t, len(ops), 1)
	item := ops[0]

	require.Equal(t, item.GetEntry().GetHash().String(), op.GetEntry().GetHash().String())

	// Returns the the added entry and the existing.
	err = db.Load(ctx, -1)
	require.NoError(t, err)

	ops, err = db.List(ctx, &orbitdb.StreamOptions{Amount: &infinity})
	require.NoError(t, err)
	require.Equal(t, len(ops), 1)

	prevHash := ops[0].GetEntry().GetHash()

	op, err = db.Add(ctx, []byte("hello2"))
	require.NoError(t, err)

	ops, err = db.List(ctx, &orbitdb.StreamOptions{Amount: &infinity})
	require.NoError(t, err)
	require.Equal(t, len(ops), 2)

	require.NotEqual(t, ops[1].GetEntry().GetHash().String(), prevHash.String())
	require.Equal(t, ops[1].GetEntry().GetHash().String(), op.GetEntry().GetHash().String())
}

func testingRepo(ctx context.Context, t testing.TB) repo.Repo {
	t.Helper()

	c := cfg.Config{}
	priv, pub, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	require.NoError(t, err)

	pid, err := peer.IDFromPublicKey(pub)
	require.NoError(t, err)

	privkeyb, err := crypto.MarshalPrivateKey(priv)
	require.NoError(t, err)

	c.Pubsub.Enabled = cfg.True
	// c.Swarm.ResourceMgr.Enabled = cfg.False
	c.Bootstrap = []string{}
	// TODO: random ports here
	c.Addresses.Swarm = []string{"/ip4/127.0.0.1/tcp/4001", "/ip4/127.0.0.1/udp/4001/quic"}
	c.Identity.PeerID = pid.Pretty()
	c.Identity.PrivKey = base64.StdEncoding.EncodeToString(privkeyb)

	return &repo.Mock{
		D: dsync.MutexWrap(ds.NewMapDatastore()),
		C: c,
	}
}

func testingIPFSAPIsMocked(ctx context.Context, t testing.TB, count int) ([]iface.CoreAPI, func()) {
	t.Helper()

	mn := testingMockNet(t)
	defer mn.Close()

	coreAPIs := make([]iface.CoreAPI, count)
	cleans := make([]func(), count)

	for i := 0; i < count; i++ {
		node := (*ipfscore.IpfsNode)(nil)

		node, cleans[i] = testingMockedIPFSNode(ctx, t, mn)
		coreAPIs[i] = testingCoreAPI(t, node)
	}

	return coreAPIs, func() {
		for i := 0; i < count; i++ {
			cleans[i]()
		}
	}
}

func testingIPFSAPIsNonMocked(ctx context.Context, t testing.TB, count int) ([]iface.CoreAPI, func()) {
	t.Helper()

	coreAPIs := make([]iface.CoreAPI, count)
	cleans := make([]func(), count)

	for i := 0; i < count; i++ {
		core, err := ipfscore.NewNode(ctx, &ipfscore.BuildCfg{
			Online: true,
			Repo:   testingRepo(ctx, t),
			ExtraOpts: map[string]bool{
				"pubsub": true,
			},
		})
		require.NoError(t, err)

		coreAPIs[i] = testingCoreAPI(t, core)
		cleans[i] = func() {
			core.Close()
		}
	}

	return coreAPIs, func() {
		for i := 0; i < count; i++ {
			cleans[i]()
		}
	}
}

func testingIPFSNodeWithoutPubsub(ctx context.Context, t testing.TB, m mocknet.Mocknet) (*ipfscore.IpfsNode, func()) {
	t.Helper()

	core, err := ipfscore.NewNode(ctx, &ipfscore.BuildCfg{
		Online: true,
		Repo:   testingRepo(ctx, t),
		Host:   mock.MockHostOption(m),
		ExtraOpts: map[string]bool{
			"pubsub": false,
		},
	})
	require.NoError(t, err)

	cleanup := func() { core.Close() }
	return core, cleanup
}

func testingMockedIPFSNode(ctx context.Context, t testing.TB, m mocknet.Mocknet) (*ipfscore.IpfsNode, func()) {
	t.Helper()

	core, err := ipfscore.NewNode(ctx, &ipfscore.BuildCfg{
		Online: true,
		Repo:   testingRepo(ctx, t),
		Host:   mock.MockHostOption(m),
		ExtraOpts: map[string]bool{
			"pubsub": true,
		},
	})
	require.NoError(t, err)

	cleanup := func() { core.Close() }
	return core, cleanup
}

func testingNonMockedIPFSNode(ctx context.Context, t testing.TB) (*ipfscore.IpfsNode, func()) {
	t.Helper()

	core, err := ipfscore.NewNode(ctx, &ipfscore.BuildCfg{
		Online: true,
		Repo:   testingRepo(ctx, t),
		ExtraOpts: map[string]bool{
			"pubsub": true,
		},
	})
	require.NoError(t, err)

	cleanup := func() { core.Close() }
	return core, cleanup
}

func testingCoreAPI(t testing.TB, core *ipfscore.IpfsNode) iface.CoreAPI {
	t.Helper()

	api, err := coreapi.NewCoreAPI(core)
	require.NoError(t, err)
	return api
}

func testingMockNet(t testing.TB) mocknet.Mocknet {
	mn := mocknet.New()
	t.Cleanup(func() { mn.Close() })
	return mn
}

func testingTempDir(t testing.TB, name string) (string, func()) {
	t.Helper()

	path, err := ioutil.TempDir("", name)
	require.NoError(t, err)

	cleanup := func() { os.RemoveAll(path) }
	return path, cleanup
}
