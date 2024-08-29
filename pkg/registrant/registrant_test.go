package registrant_test

import (
	"context"
	"crypto/ecdsa"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/mocks"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"google.golang.org/protobuf/proto"
)

type deps struct {
	ctx         context.Context
	db          *queries.Queries
	registry    *mocks.MockNodeRegistry
	privKey1    *ecdsa.PrivateKey
	privKey1Str string
	privKey2    *ecdsa.PrivateKey
	privKey3    *ecdsa.PrivateKey
}

func setup(t *testing.T) (deps, func()) {
	ctx := context.Background()
	mockRegistry := mocks.NewMockNodeRegistry(t)
	db, _, dbCleanup := testutils.NewDB(t, ctx)
	queries := queries.New(db)
	privKey1, err := crypto.GenerateKey()
	require.NoError(t, err)
	privKey2, err := crypto.GenerateKey()
	require.NoError(t, err)
	privKey3, err := crypto.GenerateKey()
	require.NoError(t, err)
	privKey1Str := "0x" + testutils.HexEncode(crypto.FromECDSA(privKey1))

	return deps{
		ctx:         ctx,
		db:          queries,
		registry:    mockRegistry,
		privKey1:    privKey1,
		privKey1Str: privKey1Str,
		privKey2:    privKey2,
		privKey3:    privKey3,
	}, dbCleanup
}

func setupWithRegistrant(t *testing.T) (deps, *registrant.Registrant, func()) {
	deps, cleanup := setup(t)

	deps.registry.EXPECT().GetNodes().Return([]registry.Node{
		{NodeID: 1, SigningKey: &deps.privKey1.PublicKey},
	}, nil)

	r, err := registrant.NewRegistrant(
		deps.ctx,
		deps.db,
		deps.registry,
		deps.privKey1Str,
	)
	require.NoError(t, err)

	return deps, r, cleanup
}

func TestNewRegistrantBadPrivateKey(t *testing.T) {
	deps, cleanup := setup(t)
	defer cleanup()

	_, err := registrant.NewRegistrant(deps.ctx, deps.db, deps.registry, "badkey")
	require.ErrorContains(t, err, "parse")
}

func TestNewRegistrantNotInRegistry(t *testing.T) {
	deps, cleanup := setup(t)
	defer cleanup()

	deps.registry.EXPECT().GetNodes().Return([]registry.Node{
		{NodeID: 2, SigningKey: &deps.privKey2.PublicKey},
		{NodeID: 3, SigningKey: &deps.privKey3.PublicKey},
	}, nil)

	_, err := registrant.NewRegistrant(deps.ctx, deps.db, deps.registry, deps.privKey1Str)
	require.ErrorContains(t, err, "registry")
}

func TestNewRegistrantNewDatabase(t *testing.T) {
	deps, cleanup := setup(t)
	defer cleanup()

	deps.registry.EXPECT().GetNodes().Return([]registry.Node{
		{NodeID: 2, SigningKey: &deps.privKey2.PublicKey},
		{NodeID: 3, SigningKey: &deps.privKey3.PublicKey},
		{NodeID: 1, SigningKey: &deps.privKey1.PublicKey},
	}, nil)

	_, err := registrant.NewRegistrant(deps.ctx, deps.db, deps.registry, deps.privKey1Str)
	require.NoError(t, err)
}

func TestNewRegistrantExistingDatabase(t *testing.T) {
	deps, cleanup := setup(t)
	defer cleanup()

	deps.registry.EXPECT().GetNodes().Return([]registry.Node{
		{NodeID: 5, SigningKey: &deps.privKey1.PublicKey},
	}, nil)

	_, err := deps.db.InsertNodeInfo(
		deps.ctx,
		queries.InsertNodeInfoParams{
			NodeID:    5,
			PublicKey: crypto.FromECDSAPub(&deps.privKey1.PublicKey),
		},
	)
	require.NoError(t, err)

	_, err = registrant.NewRegistrant(deps.ctx, deps.db, deps.registry, deps.privKey1Str)
	require.NoError(t, err)
}

func TestNewRegistrantMismatchingDatabaseNodeId(t *testing.T) {
	deps, cleanup := setup(t)
	defer cleanup()

	deps.registry.EXPECT().GetNodes().Return([]registry.Node{
		{NodeID: 7, SigningKey: &deps.privKey1.PublicKey},
	}, nil)

	_, err := deps.db.InsertNodeInfo(
		deps.ctx,
		queries.InsertNodeInfoParams{
			NodeID:    8,
			PublicKey: crypto.FromECDSAPub(&deps.privKey1.PublicKey),
		},
	)
	require.NoError(t, err)

	_, err = registrant.NewRegistrant(deps.ctx, deps.db, deps.registry, deps.privKey1Str)
	require.ErrorContains(t, err, "does not match")
}

func TestNewRegistrantMismatchingDatabasePublicKey(t *testing.T) {
	deps, cleanup := setup(t)
	defer cleanup()

	deps.registry.EXPECT().GetNodes().Return([]registry.Node{
		{NodeID: 2, SigningKey: &deps.privKey1.PublicKey},
	}, nil)

	_, err := deps.db.InsertNodeInfo(
		deps.ctx,
		queries.InsertNodeInfoParams{
			NodeID:    2,
			PublicKey: crypto.FromECDSAPub(&deps.privKey2.PublicKey),
		},
	)
	require.NoError(t, err)

	_, err = registrant.NewRegistrant(deps.ctx, deps.db, deps.registry, deps.privKey1Str)
	require.ErrorContains(t, err, "does not match")
}

func TestNewRegistrantPrivateKeyNo0x(t *testing.T) {
	deps, cleanup := setup(t)
	defer cleanup()

	deps.registry.EXPECT().GetNodes().Return([]registry.Node{
		{NodeID: 1, SigningKey: &deps.privKey1.PublicKey},
	}, nil)

	_, err := registrant.NewRegistrant(
		deps.ctx,
		deps.db,
		deps.registry,
		testutils.HexEncode(crypto.FromECDSA(deps.privKey1)),
	)
	require.NoError(t, err)
}

func TestSignStagedEnvelopeInvalidEnvelope(t *testing.T) {
	_, r, cleanup := setupWithRegistrant(t)
	defer cleanup()

	_, err := r.SignStagedEnvelope(
		queries.StagedOriginatorEnvelope{
			ID:             1,
			OriginatorTime: time.Now(),
			PayerEnvelope:  []byte{0b1},
		},
	)

	require.ErrorContains(t, err, "unmarshal")
}

func TestSignStagedEnvelopeSIDExhaustion(t *testing.T) {
	_, r, cleanup := setupWithRegistrant(t)
	defer cleanup()
	payerBytes, err := proto.Marshal(&message_api.PayerEnvelope{})
	require.NoError(t, err)

	_, err = r.SignStagedEnvelope(
		queries.StagedOriginatorEnvelope{
			ID:             0b0000000000000001000000000000000000000000000000000000000000000000,
			OriginatorTime: time.Now(),
			PayerEnvelope:  payerBytes,
		},
	)

	require.ErrorContains(t, err, "exhaustion")
}

func TestSignStagedEnvelopeSuccess(t *testing.T) {
	deps, r, cleanup := setupWithRegistrant(t)
	defer cleanup()
	payerBytes, err := proto.Marshal(
		&message_api.PayerEnvelope{UnsignedClientEnvelope: []byte{3}},
	)
	require.NoError(t, err)

	env, err := r.SignStagedEnvelope(
		queries.StagedOriginatorEnvelope{
			ID:             50,
			OriginatorTime: time.Now(),
			PayerEnvelope:  payerBytes,
		},
	)

	require.NoError(t, err)
	require.NotEmpty(t, env.GetUnsignedOriginatorEnvelope())
	require.NotEmpty(t, env.GetOriginatorSignature().Bytes)

	signingKey, err := crypto.SigToPub(
		crypto.Keccak256(env.GetUnsignedOriginatorEnvelope()),
		env.GetOriginatorSignature().Bytes,
	)
	require.NoError(t, err)
	require.True(t, signingKey.Equal(&deps.privKey1.PublicKey))

	unsignedEnv := &message_api.UnsignedOriginatorEnvelope{}
	require.NoError(t, proto.Unmarshal(env.GetUnsignedOriginatorEnvelope(), unsignedEnv))
	require.Equal(t, unsignedEnv.GetOriginatorSid(), uint64(1<<48|50))
	require.Equal(t, unsignedEnv.GetPayerEnvelope().GetUnsignedClientEnvelope()[0], uint8(3))
}
