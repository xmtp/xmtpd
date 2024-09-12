package testutils

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/api"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	mocks "github.com/xmtp/xmtpd/pkg/mocks/registry"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
)

func NewTestService(t *testing.T) (*api.Service, *sql.DB, func()) {
	ctx := context.Background()
	log := NewLog(t)
	db, _, dbCleanup := NewDB(t, ctx)
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	privKeyStr := "0x" + HexEncode(crypto.FromECDSA(privKey))
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.EXPECT().GetNodes().Return([]registry.Node{
		{NodeID: 1, SigningKey: &privKey.PublicKey},
	}, nil)
	registrant, err := registrant.NewRegistrant(ctx, queries.New(db), mockRegistry, privKeyStr)
	require.NoError(t, err)

	svc, err := api.NewReplicationApiService(ctx, log, registrant, db)
	require.NoError(t, err)

	return svc, db, func() {
		svc.Close()
		dbCleanup()
	}
}
