package indexer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/indexer/storer"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func startIndexing(t *testing.T) (*queries.Queries, context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	logger := testutils.NewLog(t)
	db, _, cleanup := testutils.NewDB(t, ctx)
	cfg := testutils.GetContractsOptions(t)
	querier := queries.New(db)

	err := StartIndexer(ctx, logger, querier, cfg)
	require.NoError(t, err)

	return querier, ctx, func() {
		cleanup()
		cancel()
	}
}

func messagePublisher(t *testing.T, ctx context.Context) *blockchain.GroupMessagePublisher {
	payerCfg := testutils.GetPayerOptions(t)
	contractsCfg := testutils.GetContractsOptions(t)
	var signer blockchain.TransactionSigner
	signer, err := blockchain.NewPrivateKeySigner(payerCfg.PrivateKey, contractsCfg.ChainID)
	require.NoError(t, err)

	client, err := blockchain.NewClient(ctx, contractsCfg.RpcUrl)
	require.NoError(t, err)

	publisher, err := blockchain.NewGroupMessagePublisher(
		testutils.NewLog(t),
		client,
		signer,
		contractsCfg,
	)
	require.NoError(t, err)

	return publisher
}

func TestStoreMessages(t *testing.T) {
	querier, ctx, cleanup := startIndexing(t)
	publisher := messagePublisher(t, ctx)
	defer cleanup()

	message := testutils.RandomBytes(32)
	groupID := testutils.RandomGroupID()
	topic := []byte(storer.BuildGroupMessageTopic(groupID))

	// Publish the message onto the blockchain
	require.NoError(t, publisher.Publish(ctx, groupID, message))

	// Poll the DB until the stored message shows up
	require.Eventually(t, func() bool {
		envelopes, err := querier.SelectGatewayEnvelopes(
			context.Background(),
			queries.SelectGatewayEnvelopesParams{
				Topic: topic,
			},
		)
		require.NoError(t, err)

		if len(envelopes) == 0 {
			return false
		}

		firstEnvelope := envelopes[0]
		require.Equal(t, firstEnvelope.OriginatorEnvelope, message)
		require.Equal(t, firstEnvelope.Topic, topic)

		return true
	}, 5*time.Second, 100*time.Millisecond, "Failed to find indexed envelope")
}
